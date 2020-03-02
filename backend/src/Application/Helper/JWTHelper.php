<?php

declare(strict_types=1);

namespace Album\Application\Helper;

use Album\Application\Clock\ClockInterface;
use Lcobucci\JWT\Builder;
use Lcobucci\JWT\Parser;
use Lcobucci\JWT\Signer\Hmac\Sha256;
use Lcobucci\JWT\Token;
use Lcobucci\JWT\ValidationData;

class JWTHelper
{
    public const TYPE_RESET_PASSWORD = 'reset_password';
    protected const DEFAULT_TIME = 3600 * 24 * 4;

    protected ClockInterface $clock;
    protected string $secret;

    public function __construct(ClockInterface $clock, string $secret)
    {
        $this->clock = $clock;
        $this->secret = $secret;
    }

    public function generateToken(array $data): Token
    {
        $signer = new Sha256();
        $token = (new Builder())
            ->setIssuedAt($this->clock->now()->getTimestamp())
            ->setExpiration($this->clock->now()->getTimestamp() + self::DEFAULT_TIME);

        foreach ($data as $key => $value) {
            $token->set($key, $value);
        }

        $token->sign($signer, $this->secret);

        return $token->getToken();
    }

    public function isTokenValid(string $token): bool
    {
        try {
            $token = (new Parser())->parse($token);

            $signer = new Sha256();
            $validationData = new ValidationData();

            return $token->verify($signer, $this->secret) && $token->validate($validationData);
        } catch (\Exception $e) {
            return false;
        }
    }

    /**
     * @return string|int
     */
    public function getData(string $token, string $dataName)
    {
        $token = (new Parser())->parse($token);

        return $token->getClaim($dataName);
    }
}
