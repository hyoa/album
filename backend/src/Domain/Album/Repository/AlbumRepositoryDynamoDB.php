<?php

declare(strict_types=1);

namespace Album\Domain\Album\Repository;

use Album\Application\Repository\AbstractDynamoDBRepository;
use Album\Domain\Album\AlbumEntity;
use Album\Domain\Album\AlbumMediaEntity;
use Album\Domain\Album\AlbumRepositoryInterface;
use Ausi\SlugGenerator\SlugGenerator;
use Aws\DynamoDb\Marshaler;

class AlbumRepositoryDynamoDB extends AbstractDynamoDBRepository implements AlbumRepositoryInterface
{
    public function insert(AlbumEntity $albumEntity): string
    {
        $this->putItem([
            'title' => $albumEntity->title,
            'description' => $albumEntity->description,
            'author' => $albumEntity->author,
            'isPrivate' => $albumEntity->private,
            'creationDate' => $albumEntity->creationDate->getTimestamp(),
            'slug' => $albumEntity->slug,
            'medias' => array_map(function (AlbumMediaEntity $media): array {
                return [
                    'mediaKey' => $media->key,
                    'author' => $media->author,
                    'mediaType' => $media->type,
                    'favorite' => $media->isFavorite,
                ];
            }, $albumEntity->medias),
        ]);

        return $albumEntity->slug;
    }

    public function updateOne(AlbumEntity $albumEntity): void
    {
        $this->updateOneItem(
            'set description = :description, author = :author, isPrivate = :isPrivate, medias = :medias, title = :title',
            ['slug' => $albumEntity->slug],
            [
                ':title' => $albumEntity->title,
                ':description' => $albumEntity->description,
                ':author' => $albumEntity->author,
                ':isPrivate' => $albumEntity->private,
                ':medias' => array_map(function (AlbumMediaEntity $media): array {
                    return [
                        'mediaKey' => $media->key,
                        'author' => $media->author,
                        'mediaType' => $media->type,
                        'favorite' => $media->isFavorite,
                    ];
                }, $albumEntity->medias),
            ]
        );
    }

    public function findOne(array $data): ?AlbumEntity
    {
        if (array_key_exists('slug', $data)) {
            $slug = $data['slug'];
        } elseif (array_key_exists('title', $data)) {
            $slug = (new SlugGenerator())->generate($data['title']);
        } else {
            throw new \UnexpectedValueException();
        }

        $item = $this->getItem(['slug' => $slug]);

        if (is_null($item)) {
            return null;
        }

        $album = new AlbumEntity();
        $album->hydrate($this->transformDocumentToNormalData($item));

        return $album;
    }

    public function find(bool $includePrivate, bool $includeNoMedia, int $limit, int $offset, ?string $term, string $order = 'asc'): array
    {
        $filterExpressions = [];
        $expressionAttributeValues = [];

        if (!$includePrivate) {
            $filterExpressions[] = 'isPrivate = :isPrivate';
            $expressionAttributeValues[':isPrivate'] = false;
        }

        if (!$includeNoMedia) {
            $filterExpressions[] = 'size(medias) > :mediasSize';
            $expressionAttributeValues[':mediasSize'] = 0;
        }

        if ($term !== null) {
            $filterExpressions[] = '(contains(title, :term) OR contains(description, :term))';
            $expressionAttributeValues[':term'] = $term;
        }

        $filterExpression = implode(' AND ', $filterExpressions);

        $items = $this->scan($filterExpression, $expressionAttributeValues);

        $albums = array_map(function (array $item): AlbumEntity {
            return (new AlbumEntity())
                        ->hydrate($this->transformDocumentToNormalData($item));
        }, $items);

        if ($order === 'asc') {
            usort($albums, fn (AlbumEntity $album1, AlbumEntity $album2) => $album1->creationDate->getTimestamp() <=> $album2->creationDate->getTimestamp());
        } else {
            usort($albums, fn (AlbumEntity $album1, AlbumEntity $album2) => $album2->creationDate->getTimestamp() <=> $album1->creationDate->getTimestamp());
        }

        return array_slice($albums, $offset, $limit);
    }

    public function deleteOne(array $filters): void
    {
        // TODO: Implement deleteOne() method.
    }

    public function transformDocumentToNormalData(array $document): array
    {
        $document = (array) (new Marshaler())->unmarshalItem($document);
        $data = [];

        if (array_key_exists('title', $document)) {
            $data['title'] = $document['title'];
        }

        if (array_key_exists('description', $document)) {
            $data['description'] = $document['description'];
        }

        if (array_key_exists('author', $document)) {
            $data['author'] = $document['author'];
        }

        if (array_key_exists('isPrivate', $document)) {
            $data['private'] = $document['isPrivate'];
        }

        if (array_key_exists('medias', $document)) {
            $mediasArray = (array) $document['medias'];
            $medias = [];

            foreach ($mediasArray as $item) {
                $media = new AlbumMediaEntity();
                $media->key = $item['mediaKey'];
                $media->type = $item['mediaType'];

                if (isset($item['author'])) {
                    $media->author = $item['author'];
                }

                if (isset($item['favorite'])) {
                    $media->isFavorite = $item['favorite'];
                }

                $medias[] = $media;
            }

            $data['medias'] = $medias;
        }

        if (array_key_exists('slug', $document)) {
            $data['slug'] = $document['slug'];
        }

        if (array_key_exists('creationDate', $document)) {
            $date = (new \DateTimeImmutable())->setTimestamp($document['creationDate']);
            $data['creationDate'] = $date;
        }

        return $data;
    }

    protected function getTableName(): string
    {
        return 'album';
    }
}
