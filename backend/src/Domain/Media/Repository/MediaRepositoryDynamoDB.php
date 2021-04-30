<?php

declare(strict_types=1);

namespace Album\Domain\Media\Repository;

use Album\Application\Repository\AbstractDynamoDBRepository;
use Album\Domain\Media\MediaEntity;
use Album\Domain\Media\MediaRepositoryInterface;

class MediaRepositoryDynamoDB extends AbstractDynamoDBRepository implements MediaRepositoryInterface
{
    public function insert(MediaEntity $mediaEntity): void
    {
        $this->putItem([
            'mediaKey' => $mediaEntity->key,
            'mediaType' => $mediaEntity->type,
            'author' => $mediaEntity->author,
            'folder' => $mediaEntity->folder,
            'uploadDate' => $mediaEntity->uploadDate->getTimestamp(),
            'visible' => $mediaEntity->visible,
        ]);
    }

    public function findByFolder(string $folder): array
    {
        $medias = [];

        $items = $this->query(
            'folder = :folder',
            [
                ':folder' => $folder,
            ],
            'folderIndex'
        );

        foreach ($items as $item) {
            $media = new MediaEntity();
            $media->hydrate($this->transformDocumentToNormalData($item));
            $medias[] = $media;
        }

        return $medias;
    }

    public function findFolders(?string $term = null): array
    {
        $folders = [];

        if ($term !== null) {
            $items = $this->scan(
                'contains(folder, :folder)',
                [
                    ':folder' => $term,
                ]
            );
        } else {
            $items = $this->scan();
        }

        foreach ($items as $item) {
            if (!in_array($item['folder']['S'], $folders, true)) {
                $folders[] = $item['folder']['S'];
            }
        }

        return $folders;
    }

    public function find(): array
    {
        $medias = [];
        $items = $this->scan();

        foreach ($items as $item) {
            $media = new MediaEntity();
            $media->hydrate($this->transformDocumentToNormalData($item));
            $medias[] = $media;
        }

        return $medias;
    }

    public function findOne(array $term): ?MediaEntity
    {
        $item = $this->getItem(['mediaKey' => $term['key']]);

        if ($item !== null) {
            $media = new MediaEntity();
            $media->hydrate($this->transformDocumentToNormalData($item));

            return $media;
        }

        return null;
    }

    public function update(array $condition, array $change): void
    {
        if (array_key_exists('folder', $condition)) {
            $itemsToRename = $this->query(
                'folder = :folder',
                [':folder' => $condition['folder']],
                'folderIndex'
            );

            foreach ($itemsToRename as $item) {
                $this->updateOneItem(
                    'set folder = :folder',
                    ['mediaKey' => $item['mediaKey']['S']],
                    [':folder' => $change['folder']]
                );
            }
        } elseif (array_key_exists('key', $condition)) {
            $keysExpressions = [];
            $eav = [];
            $i = 0;
            foreach ($condition['key']['$in'] as $expression) {
                $key = ':key'.$i;
                $keysExpressions[] = $key;
                $eav[$key] = $expression;
                $i++;
            }

            $itemsToUpdate = $this->scan(
                'mediaKey IN ('.implode(',', $keysExpressions).')',
                $eav
            );

            if (array_key_exists('folder', $change)) {
                foreach ($itemsToUpdate as $item) {
                    $this->updateOneItem(
                        'set folder = :folder',
                        ['mediaKey' => $item['mediaKey']['S']],
                        [':folder' => $change['folder']]
                    );
                }
            }

            if (array_key_exists('visible', $change)) {
                foreach ($itemsToUpdate as $item) {
                    $this->updateOneItem(
                        'set visible = :visible',
                        ['mediaKey' => $item['mediaKey']['S']],
                        [':visible' => $change['visible']]
                    );
                }
            }
        }
    }

    protected function transformDocumentToNormalData(array $document): array
    {
        $data = [];

        if (array_key_exists('mediaKey', $document)) {
            $data['key'] = $document['mediaKey']['S'];
        }

        if (array_key_exists('author', $document)) {
            $data['author'] = $document['author']['S'];
        }

        if (array_key_exists('mediaType', $document)) {
            $data['type'] = (int) $document['mediaType']['N'];
        }

        if (array_key_exists('folder', $document)) {
            $data['folder'] = $document['folder']['S'];
        }

        if (array_key_exists('uploadDate', $document)) {
            $date = (new \DateTimeImmutable())->setTimestamp((int) $document['uploadDate']['N']);
            $data['uploadDate'] = $date;
        }

        if (array_key_exists('visible', $document)) {
            $data['visible'] = $document['visible']['BOOL'];
        }

        return $data;
    }

    protected function getTableName(): string
    {
        return 'media';
    }
}
