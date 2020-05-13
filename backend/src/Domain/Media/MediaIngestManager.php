<?php

declare(strict_types=1);

namespace Album\Domain\Media;

use Album\Application\Clock\ClockInterface;
use Album\Application\VideoFormatter\VideoFormatterInterface;
use Album\Domain\Album\AlbumManager;
use Album\Domain\Album\AlbumMediaEntity;
use Album\Domain\Media\Exception\InvalidVideoSizeException;
use Ausi\SlugGenerator\SlugGenerator;

class MediaIngestManager
{
    protected MediaRepositoryInterface $mediaRepository;

    protected ClockInterface $clock;

    protected MediaStorageInterface $mediaStorage;

    protected string $mediaStorageLocation;

    protected string $videoRawsStorageLocation;

    protected VideoFormatterInterface $videoFormatter;

    protected AlbumManager $albumManager;

    public function __construct(
        MediaRepositoryInterface $mediaRepository,
        ClockInterface $clock,
        MediaStorageInterface $mediaStorage,
        VideoFormatterInterface $videoFormatter,
        AlbumManager $albumManager,
        string $mediaStorageLocation,
        string $videoRawStorageLocation
    ) {
        $this->clock = $clock;
        $this->mediaStorage = $mediaStorage;
        $this->mediaRepository = $mediaRepository;
        $this->mediaStorageLocation = $mediaStorageLocation;
        $this->videoRawsStorageLocation = $videoRawStorageLocation;
        $this->videoFormatter = $videoFormatter;
        $this->albumManager = $albumManager;
    }

    public function ingest(string $key): bool
    {
        if (!is_null($this->mediaRepository->findOne(['key' => $key]))) {
            return false;
        }

        $metaData = $this->mediaStorage->getMediaMetadata($key, $this->mediaStorageLocation);

        $media = new MediaEntity();
        $media->key = $key;
        $media->author = $metaData->author ?? 'none';
        $media->folder = $metaData->folder ?? 'sans dossier';
        $media->uploadDate = $this->clock->now();
        $media->type = $metaData->getMediaType();

        if ($metaData->album !== null) {
            $album = $this->albumManager->findBySlug((new SlugGenerator())->generate($metaData->album));

            if ($album !== null) {
                $albumMediaEntity = new AlbumMediaEntity();
                $albumMediaEntity->key = $media->key;
                $albumMediaEntity->author = $media->author;
                $albumMediaEntity->folder = $media->folder;
                $albumMediaEntity->uploadDate = $media->uploadDate;
                $albumMediaEntity->type = $media->type;

                $this->albumManager->addMedias($album, [$albumMediaEntity]);
            }
        }

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
        $metadata = $this->mediaStorage->getMediaMetadata($key, $this->videoRawsStorageLocation);

        $this->mediaStorage->putObject($key, $this->mediaStorageLocation, $pathVideoFormatted, 'video/mp4', $metadata);

        return true;
    }
}
