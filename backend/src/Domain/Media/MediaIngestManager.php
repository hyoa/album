<?php

declare(strict_types=1);

namespace Album\Domain\Media;

use Album\Application\Clock\ClockInterface;
use Album\Application\VideoFormatter\VideoFormatterInterface;
use Album\Domain\Album\AlbumManager;
use Album\Domain\Album\AlbumMediaEntity;
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
        if (!is_null($media = $this->mediaRepository->findOne(['key' => $key]))) {
            if ($media?->type === MediaEntity::TYPE_VIDEO && $media->visible === false) {
                $this->mediaRepository->update(['key' => ['$in' => [$key]]], ['visible' => true]);

                return true;
            }

            return false;
        }

        $this->createMedia($key, $this->mediaStorageLocation);

        return true;
    }

    public function ingestVideo(string $key): bool
    {
        $media = $this->mediaRepository->findOne(['key' => $key]);

        if ($media?->visible) {
            return false;
        }

        $this->createMedia($key, $this->videoRawsStorageLocation, false);

        return $this->videoFormatter->run($key);
    }

    protected function createMedia(string $key, string $location, bool $visible = true): void
    {
        $metaData = $this->mediaStorage->getMediaMetadata($key, $location);

        $media = new MediaEntity();
        $media->key = $key;
        $media->author = $metaData->author ?? 'none';
        $media->folder = $metaData->folder ?? 'sans dossier';
        $media->uploadDate = $this->clock->now();
        $media->type = $metaData->getMediaType();
        $media->visible = $visible;

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
    }
}
