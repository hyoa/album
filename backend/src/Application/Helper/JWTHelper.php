<?php

declare(strict_types=1);

namespace Album\Application\Helper;

use Album\Application\Clock\ClockInterface;
use Lcobucci\Clock\SystemClock;
use Lcobucci\JWT\Configuration;
use Lcobucci\JWT\Signer\Hmac\Sha256;
use Lcobucci\JWT\Signer\Key\InMemory;
use Lcobucci\JWT\Token;
use Lcobucci\JWT\UnencryptedToken;
use Lcobucci\JWT\Validation\Constraint\SignedWith;
use Lcobucci\JWT\Validation\Constraint\StrictValidAt;

class JWTHelper
{
    public const TYPE_RESET_PASSWORD = 'reset_password';
    protected const DEFAULT_TIME = 3600 * 24 * 4;

    protected ClockInterface $clock;
    protected Configuration $jwtConfiguration;

    public function __construct(
        ClockInterface $clock,
        string $secret,
    ) {
        $this->clock = $clock;

        $this->jwtConfiguration = Configuration::forSymmetricSigner(
            new Sha256(),
            InMemory::base64Encoded('password', $secret)
        );

        $this->jwtConfiguration->setValidationConstraints(
            new StrictValidAt(new SystemClock($clock->now()->getTimezone())),
            new SignedWith($this->jwtConfiguration->signer(), $this->jwtConfiguration->signingKey())
        );
    }

    public function generateToken(array $data): Token
    {
        $newTokenConfiguration = $this->jwtConfiguration->builder()
                    ->issuedAt($this->clock->now())
                    ->expiresAt($this->clock->now()->modify('+ 4 days'))
                    ->canOnlyBeUsedAfter($this->clock->now())
        ;

        foreach ($data as $key => $value) {
            $newTokenConfiguration->withClaim($key, $value);
        }

        return $newTokenConfiguration->getToken($this->jwtConfiguration->signer(), $this->jwtConfiguration->verificationKey());
    }

    public function isTokenValid(string $token): bool
    {
        try {
            $parsedToken = $this->jwtConfiguration->parser()->parse($token);
            $constraints = $this->jwtConfiguration->validationConstraints();

            return $this->jwtConfiguration->validator()->validate($parsedToken, ...$constraints);
        } catch (\Exception $e) {
            return false;
        }
    }

    public function getData(string $token, string $claimName): string | int | array
    {
        $parsedToken = $this->jwtConfiguration->parser()->parse($token);

        if (!$parsedToken instanceof UnencryptedToken) {
            throw new \InvalidArgumentException();
        }

        return $parsedToken->claims()->get($claimName);
    }
}
