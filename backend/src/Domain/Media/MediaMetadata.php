<?php

declare(strict_types=1);

namespace Album\Domain\Media;

class MediaMetadata
{
    private const FORMAT_IMAGE_WHITELIST = [
        'image/jpg',
        'image/png',
        'image/jpeg',
        'image/JPG',
        'image/PNG',
        'image/JPEG',
    ];

    private const FORMAT_VIDEO_WHITELIST = [
        'video/mp4',
        'video/MP4',
    ];

    public ?string $author;
    public ?string $folder;
    public ?string $album;
    public ?string $contentType;

    public function getMediaType(): int
    {
        return in_array($this->contentType, self::FORMAT_VIDEO_WHITELIST, true) ? MediaEntity::TYPE_VIDEO : MediaEntity::TYPE_IMAGE;
    }

    public function toArray(): array
    {
        return [
            'author' => $this->author,
            'folder' => $this->folder,
            'album' => $this->album,
            'contentType' => $this->contentType,
        ];
    }
}
