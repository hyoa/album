<?php

declare(strict_types=1);

namespace Album\Domain\Album;

use Album\Application\Clock\ClockInterface;
use Album\Domain\Media\MediaEntity;
use Album\Domain\Media\MediaStorageInterface;
use Ausi\SlugGenerator\SlugGenerator;

class AlbumManager
{
    protected ClockInterface $clock;
    protected AlbumRepositoryInterface $albumRepository;
    protected MediaStorageInterface $mediaStorage;

    public function __construct(
        ClockInterface $clock,
        AlbumRepositoryInterface $albumRepository,
        MediaStorageInterface $mediaStorage
    ) {
        $this->clock = $clock;
        $this->albumRepository = $albumRepository;
        $this->mediaStorage = $mediaStorage;
    }

    public function save(string $title, string $description, bool $private, string $author): AlbumEntity
    {
        $album = new AlbumEntity();
        $album->hydrate(
            [
                'title' => $title,
                'description' => $description,
                'private' => $private,
                'author' => $author,
            ]
        );
        $album->creationDate = $this->clock->now();

        $slug = (new SlugGenerator())->generate($title);
        $album->slug = $slug;

        $id = $this->albumRepository->insert($album);
        $album->id = $id;

        return $this->attachUrisToMedias($album);
    }

    public function updateOne(AlbumEntity $albumEntity, array $data): AlbumEntity
    {
        if (array_key_exists('title', $data)) {
            $albumEntity->title = $data['title'];
            $albumEntity->slug = (new SlugGenerator())->generate($albumEntity->title);
        }

        if (array_key_exists('description', $data)) {
            $albumEntity->description = $data['description'];
        }

        if (array_key_exists('private', $data)) {
            $albumEntity->private = $data['private'];
        }

        $this->albumRepository->updateOne($albumEntity);

        return $this->attachUrisToMedias($albumEntity);
    }

    public function addMedias(AlbumEntity $albumEntity, array $medias): AlbumEntity
    {
        foreach ($medias as $media) {
            $albumEntity->addMedia($media);
        }

        $this->albumRepository->updateOne($albumEntity);

        return $this->attachUrisToMedias($albumEntity);
    }

    public function removeMedias(AlbumEntity $albumEntity, array $mediasToRemove): AlbumEntity
    {
        $mediasToSave = [];

        /** @var MediaEntity $mediaToRemove */
        foreach ($mediasToRemove as $mediaToRemove) {
            /** @var MediaEntity $media */
            foreach ($albumEntity->medias as $media) {
                if ($media->key !== $mediaToRemove->key) {
                    $mediasToSave[] = $media;
                }
            }
        }

        $albumEntity->medias = $mediasToSave;

        $this->albumRepository->updateOne($albumEntity);

        return $this->attachUrisToMedias($albumEntity);
    }

    public function findBySlug(string $slug): ?AlbumEntity
    {
        $album = $this->albumRepository->findOne(['slug' => $slug]);

        return is_null($album) ? null : $this->attachUrisToMedias($album);
    }

    public function findMany(
        bool $includePrivateAlbum = true,
        bool $includeNoMedias = true,
        int $limit = 10,
        int $offset = 0,
        string $term = null,
        string $order = 'asc'
    ): array {
        return $this->albumRepository->find($includePrivateAlbum, $includeNoMedias, $limit, $offset, $term, $order);
    }

    public function deleteOne(AlbumEntity $albumEntity): void
    {
        $this->albumRepository->deleteOne(['slug' => $albumEntity->slug]);
    }

    public function getAdminResume(): array
    {
        $resume = [
            'publicCount' => 0,
            'privateCount' => 0,
        ];

        /** @var AlbumEntity[] $albums */
        $albums = $this->albumRepository->find(true, true, 100000, 0, null);

        foreach ($albums as $album) {
            if ($album->private) {
                $resume['privateCount']++;
            } else {
                $resume['publicCount']++;
            }
        }

        return $resume;
    }

    public function toggleFavorite(AlbumEntity $album, string $favoriteKey, bool $isFavorite): AlbumEntity
    {
        $medias = [];

        /** @var AlbumMediaEntity $media */
        foreach ($album->medias as $media) {
            if ($media->key === $favoriteKey) {
                $media->isFavorite = $isFavorite;
            }

            $medias[] = $media;
        }

        $album->medias = $medias;

        $this->albumRepository->updateOne($album);

        return $this->attachUrisToMedias($album);
    }

    protected function attachUrisToMedias(AlbumEntity $album): AlbumEntity
    {
        $mediasUpdated = [];
        /** @var MediaEntity $media */
        foreach ($album->medias as $media) {
            $mediaUris = $this->mediaStorage->getUrisToAccessStore($media->key, $media->type);
            $media->mediasUri = $mediaUris;
            $mediasUpdated[] = $media;
        }

        $album->medias = $mediasUpdated;

        return $album;
    }
}
