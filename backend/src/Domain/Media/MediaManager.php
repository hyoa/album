<?php

declare(strict_types=1);

namespace Album\Domain\Media;

use Album\Application\Clock\ClockInterface;

class MediaManager
{
    protected ClockInterface $clock;

    protected MediaRepositoryInterface $mediaRepository;

    protected MediaStorageInterface $mediaStorage;

    public function __construct(ClockInterface $clock, MediaRepositoryInterface $mediaRepository, MediaStorageInterface $mediaStorage)
    {
        $this->clock = $clock;
        $this->mediaRepository = $mediaRepository;
        $this->mediaStorage = $mediaStorage;
    }

    public function save(string $key, string $author, int $type, string $folder, bool $visible = true): MediaEntity
    {
        $media = new MediaEntity();
        $media->hydrate(['key' => $key, 'author' => $author, 'type' => $type, 'folder' => $folder, 'visible' => $visible]);
        $media->uploadDate = $this->clock->now();

        $this->mediaRepository->insert($media);

        return $media;
    }

    public function findByFolder(string $folder): array
    {
        /** @var MediaEntity[] $medias */
        $medias = $this->mediaRepository->findByFolder($folder);

        foreach ($medias as $media) {
            $uris = $this->mediaStorage->getUrisToAccessStore($media->key, $media->type);
            $media->mediasUri = $uris;
        }

        return $medias;
    }

    public function findFolders(?string $term = null): array
    {
        return $this->mediaRepository->findFolders($term);
    }

    public function deleteFolder(string $folderName): bool
    {
        try {
            $this->mediaRepository->update(['folder' => $folderName], ['folder' => 'none']);

            return true;
        } catch (\Exception $e) {
            return false;
        }
    }

    public function updateFolderName(string $folderNameToUpload, string $newFolderName): void
    {
        $this->mediaRepository->update(['folder' => $folderNameToUpload], ['folder' => $newFolderName]);
    }

    public function getAdminResume(): array
    {
        /** @var MediaEntity[] $medias */
        $medias = $this->mediaRepository->find();

        $resume = [
            'imagesCount' => 0,
            'videosCount' => 0,
        ];

        foreach ($medias as $media) {
            if ($media->type === MediaEntity::TYPE_IMAGE) {
                $resume['imagesCount']++;
            } else {
                $resume['videosCount']++;
            }
        }

        return $resume;
    }

    public function updateFolderNameForMedias(string $folderName, array $mediasKeys): void
    {
        $this->mediaRepository->update(['key' => ['$in' => $mediasKeys]], ['folder' => $folderName]);
    }
}
