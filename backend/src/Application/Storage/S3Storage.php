<?php

declare(strict_types=1);

namespace Album\Application\Storage;

use Album\Application\Clock\ClockInterface;
use Album\Domain\Media\MediaEntity;
use Album\Domain\Media\MediaStorageInterface;
use Aws\CloudFront\UrlSigner;
use Aws\S3\S3Client;

class S3Storage implements MediaStorageInterface
{
    protected S3Client $s3Client;

    protected ClockInterface $clock;

    protected string $proxyImage;

    protected string $keyPairId;

    protected string $awsPk;

    protected string $mediaStorageLocation;

    protected string $videoRawStorageLocation;

    public function __construct(
        ClockInterface $clock,
        S3Client $s3Client,
        string $proxyImage,
        string $keyPairId,
        string $awsPk,
        string $mediaStorageLocation,
        string $videoRawStorageLocation
    ) {
        $this->clock = $clock;
        $this->s3Client = $s3Client;
        $this->proxyImage = $proxyImage;
        $this->keyPairId = $keyPairId;
        $this->awsPk = $awsPk;
        $this->mediaStorageLocation = $mediaStorageLocation;
        $this->videoRawStorageLocation = $videoRawStorageLocation;
    }

    public function generateSignedUri(string $key, string $location, string $commandType): string
    {
        $cmd = $this->s3Client->getCommand($commandType, [
            'Bucket' => $this->getBucket($location),
            'Key' => $key,
        ]);

        $request = $this->s3Client->createPresignedRequest($cmd, '+10 minutes');

        return (string) $request->getUri();
    }

    public function getUrisToAccessStore(string $key, int $type = MediaEntity::TYPE_IMAGE): array
    {
        if ($type === MediaEntity::TYPE_VIDEO) {
            $uri = $this->generateSignedUri($key, self::LOCATION_MEDIAS, 'GetObject');

            return [
                'small' => $uri,
                'original' => $uri,
                'medium' => $uri,
            ];
        }

        return [
            'small' => $this->getUriToAccessImageProxy($key, MediaEntity::SMALL_WIDTH),
            'medium' => $this->getUriToAccessImageProxy($key, MediaEntity::MEDIUM_WIDTH),
            'original' => $this->getUriToAccessImageProxy($key, MediaEntity::LARGE_WIDTH),
        ];
    }

    public function getUriToAccessImageProxy(string $key, ?int $width = null, ?int $height = null): string
    {
        $this->writePkIntoFile();
        $signer = new UrlSigner($this->keyPairId, '/tmp/aws-key.pem');
        $data = [
            'bucket' => $this->getBucket(self::LOCATION_MEDIAS),
            'key' => $key,
            'edits' => [
                'rotate' => null,
            ],
        ];

        if (is_int($width) || is_int($height)) {
            $data['edits']['resize'] = [];
            $data['edits']['resize']['fit'] = 'cover';

            if (is_int($width)) {
                $data['edits']['resize']['width'] = $width;
            }

            if (is_int($height)) {
                $data['edits']['resize']['height'] = $height;
            }
        }

        $json = (string) json_encode($data);

        $uri = $this->proxyImage.'/'.base64_encode($json);

        $time = $this->clock->now()->getTimestamp() + 600;

        return $signer->getSignedUrl($uri, $time);
    }

    public function putObject(string $key, string $location, string $videoPath, string $contentType): void
    {
        $this->s3Client->putObject([
            'Bucket' => $location,
            'Key' => $key,
            'Body' => file_get_contents($videoPath),
            'ContentType' => $contentType,
        ]);
    }

    public function getFileSize(string $key, string $location): int
    {
        $result = $this->s3Client->headObject([
            'Key' => $key,
            'Bucket' => $location,
        ]);

        return (int) $result['ContentLength'];
    }

    protected function writePkIntoFile(): void
    {
        $env = $this->awsPk;
        $key = <<<EOT
-----BEGIN RSA PRIVATE KEY-----
{$env}
-----END RSA PRIVATE KEY-----
EOT;

        if (!file_exists('/tmp/aws-key.pem')) {
            file_put_contents('/tmp/aws-key.pem', $key);
        }
    }

    protected function getBucket(string $location): string
    {
        switch ($location) {
            case self::LOCATION_MEDIAS:
                return $this->mediaStorageLocation;
            case self::LOCATION_RAW_VIDEOS:
                return $this->videoRawStorageLocation;
            default:
                throw new \Exception('Invalid storage location');
        }
    }
}
