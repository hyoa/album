<?php

declare(strict_types=1);

namespace Album\Domain\Media;

interface MediaStorageInterface
{
    const LOCATION_MEDIAS = 'medias';
    const LOCATION_RAW_VIDEOS = 'raw_videos';

    public function getUrisToAccessStore(string $key, int $type = MediaEntity::TYPE_IMAGE): array;

    public function getUriToAccessImageProxy(string $key, ?int $width = null, ?int $height = null): string;

    public function generateSignedUri(string $key, string $location, string $commandType, array $metadata = []): string;

    public function putObject(string $key, string $location, string $videoPath, string $contentType, MediaMetadata $mediaMetadata): void;

    public function getFileSize(string $key, string $location): int;

    public function getMediaMetadata(string $key, string $location): MediaMetadata;
}
