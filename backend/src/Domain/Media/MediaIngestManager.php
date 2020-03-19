<?php

declare(strict_types=1);

namespace Album\Domain\Media;

use Album\Application\Clock\ClockInterface;
use Album\Application\VideoFormatter\VideoFormatterInterface;
use Album\Domain\Media\Exception\InvalidVideoSizeException;

class MediaIngestManager
{
    private const FORMAT_IMAGE_WHITELIST = [
        'jpg',
        'png',
        'jpeg',
        'JPG',
        'PNG',
        'JPEG',
    ];

    private const FORMAT_VIDEO_WHITELIST = [
        'mp4',
        'MP4',
    ];

    protected MediaRepositoryInterface $mediaRepository;

    protected ClockInterface $clock;

    protected MediaStorageInterface $mediaStorage;

    protected string $mediaStorageLocation;

    protected string $videoRawsStorageLocation;

    protected VideoFormatterInterface $videoFormatter;

    public function __construct(
        MediaRepositoryInterface $mediaRepository,
        ClockInterface $clock,
        MediaStorageInterface $mediaStorage,
        VideoFormatterInterface $videoFormatter,
        string $mediaStorageLocation,
        string $videoRawStorageLocation
    ) {
        $this->clock = $clock;
        $this->mediaStorage = $mediaStorage;
        $this->mediaRepository = $mediaRepository;
        $this->mediaStorageLocation = $mediaStorageLocation;
        $this->videoRawsStorageLocation = $videoRawStorageLocation;
        $this->videoFormatter = $videoFormatter;
    }

    public function ingest(string $key): bool
    {
        if (!is_null($this->mediaRepository->findOne(['key' => $key]))) {
            return false;
        }

        $metaData = $this->getMetaData($key);

        $media = new MediaEntity();
        $media->key = $key;
        $media->author = $metaData['author'];
        $media->folder = $metaData['folder'];
        $media->uploadDate = $this->clock->now();
        $media->type = $metaData['type'];

        $this->mediaRepository->insert($media);

        return true;
    }

    public function ingestVideo(string $key): bool
    {
        $media = $this->mediaRepository->findOne(['key' => $key]);

        if (!is_null($media)) {
            return false;
        }

        $fileSize = $this->mediaStorage->getFileSize($key, $this->videoRawsStorageLocation);

        if ($fileSize > 200000000) {
            throw new InvalidVideoSizeException('Video size cannot exceed 200Mo');
        }

        $mediaUris = $this->mediaStorage->generateSignedUri($key, MediaStorageInterface::LOCATION_RAW_VIDEOS, 'GetObject');
        $pathVideoFormatted = $this->videoFormatter->run($mediaUris, $key);
        $this->mediaStorage->putObject($key, $this->mediaStorageLocation, $pathVideoFormatted, 'video/mp4');

        return true;
    }

    protected function getMetaData(string $fileName): array
    {
        $regex = '/^([a-zA-Z\d]+)_([A-Za-zÀ-ÿ\d\-]+)_(\w+)\.([a-z0-9]+)$/';
        preg_match_all($regex, $fileName, $matches, PREG_SET_ORDER, 0);

        if (!isset($matches[0][1], $matches[0][2])) {
            return [
                'author' => 'none',
                'folder' => 'sans fichier',
                'type' => MediaEntity::TYPE_IMAGE,
            ];
        }

        $type = MediaEntity::TYPE_IMAGE;
        if (in_array($matches[0][4], self::FORMAT_VIDEO_WHITELIST, true)) {
            $type = MediaEntity::TYPE_VIDEO;
        }

        return [
            'author' => $matches[0][1],
            'folder' => str_replace('-', ' ', $matches[0][2]),
            'type' => $type,
        ];
    }
}
