<?php

declare(strict_types=1);

namespace Album\Domain\Media;

use Album\Domain\Media\Exception\InvalidMediaTypeException;

class MediaEntity
{
    public const TYPE_IMAGE = 1;
    public const TYPE_VIDEO = 2;

    public const SMALL_WIDTH = 400;
    public const MEDIUM_WIDTH = 800;
    public const LARGE_WIDTH = 1024;

    protected const ENUM_TYPE_MAP = [
        MediaEntity::TYPE_IMAGE => 'image',
        MediaEntity::TYPE_VIDEO => 'video',
    ];

    public string $key;

    public string $author;

    public int $type;

    public string $folder;

    public \DateTimeImmutable $uploadDate;

    public array $mediasUri;

    public bool $visible = true;

    public function hydrate(array $data): void
    {
        if (array_key_exists('key', $data)) {
            $this->key = $data['key'];
        }

        if (array_key_exists('author', $data)) {
            $this->author = $data['author'];
        }

        if (array_key_exists('type', $data)) {
            $this->type = $data['type'];
        }

        if (array_key_exists('folder', $data)) {
            $this->folder = $data['folder'];
        }

        if (array_key_exists('uploadDate', $data)) {
            $this->uploadDate = $data['uploadDate'];
        }

        if (array_key_exists('visible', $data)) {
            $this->visible = $data['visible'];
        }
    }

    public function getMediaUri(string $size): string
    {
        return $this->mediasUri[$size];
    }

    /**
     * @throws InvalidMediaTypeException
     */
    public function getAsArray(): array
    {
        return [
            'key' => $this->key,
            'author' => $this->author,
            'uploadDate' => $this->uploadDate->format('d-m-Y'),
            'folder' => $this->folder,
            'uris' => $this->mediasUri,
            'type' => $this->getTypeAsString(),
        ];
    }

    /**
     * @throws InvalidMediaTypeException
     */
    public function getTypeAsString(): ?string
    {
        if (!array_key_exists($this->type, self::ENUM_TYPE_MAP)) {
            throw new InvalidMediaTypeException();
        }

        return self::ENUM_TYPE_MAP[$this->type];
    }

    /**
     * @throws InvalidMediaTypeException
     */
    public function setTypeFromString(string $type): void
    {
        if (!in_array($type, self::ENUM_TYPE_MAP, true)) {
            throw new InvalidMediaTypeException();
        }

        $this->type = (int) array_search($type, self::ENUM_TYPE_MAP, true);
    }
}
