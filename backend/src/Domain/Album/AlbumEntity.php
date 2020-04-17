<?php

declare(strict_types=1);

namespace Album\Domain\Album;

use Album\Domain\Media\MediaEntity;

class AlbumEntity
{
    public string $title;

    public ?string $description;

    public bool $private;

    public string $author;

    public \DateTimeImmutable $creationDate;

    public string $id;

    public string $slug;

    public array $medias = [];

    public function hydrate(array $data): self
    {
        if (array_key_exists('id', $data)) {
            $this->id = $data['id'];
        }

        if (array_key_exists('title', $data)) {
            $this->title = $data['title'];
        }

        if (array_key_exists('slug', $data)) {
            $this->slug = $data['slug'];
        }

        if (array_key_exists('description', $data)) {
            $this->description = $data['description'] !== '' ? $data['description'] : null;
        }

        if (array_key_exists('private', $data)) {
            $this->private = $data['private'];
        }

        if (array_key_exists('author', $data)) {
            $this->author = $data['author'];
        }

        if (array_key_exists('medias', $data)) {
            $this->medias = $data['medias'];
        }

        if (array_key_exists('creationDate', $data)) {
            $this->creationDate = $data['creationDate'];
        }

        return $this;
    }

    public function addMedia(MediaEntity $mediaEntity): self
    {
        $this->medias[] = $mediaEntity;

        return $this;
    }

    public function getAsArray(): array
    {
        return [
            'title' => $this->title,
            'description' => $this->description,
            'private' => $this->private,
            'author' => $this->author,
            'slug' => $this->slug,
            'creationDate' => $this->creationDate->format('d-m-Y'),
            'medias' => array_map(function (AlbumMediaEntity $mediaEntity): array {
                return [
                    'key' => $mediaEntity->key,
                    'author' => $mediaEntity->author,
                    'type' => $mediaEntity->getTypeAsString(),
                    'uris' => [
                        'small' => $mediaEntity->getMediaUri('small'),
                        'medium' => $mediaEntity->getMediaUri('medium'),
                        'original' => $mediaEntity->getMediaUri('original'),
                    ],
                    'favorite' => $mediaEntity->isFavorite,
                ];
            }, $this->medias),
        ];
    }

    public function getFavorites(): array
    {
        return array_filter($this->medias, function (AlbumMediaEntity $mediaEntity): bool {
            return $mediaEntity->isFavorite;
        });
    }
}
