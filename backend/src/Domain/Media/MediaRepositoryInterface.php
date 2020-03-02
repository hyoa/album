<?php

declare(strict_types=1);

namespace Album\Domain\Media;

interface MediaRepositoryInterface
{
    public function insert(MediaEntity $mediaEntity): void;

    public function findByFolder(string $folder): array;

    public function findFolders(?string $term = null): array;

    public function update(array $condition, array $change): void;

    public function find(): array;

    public function findOne(array $term): ?MediaEntity;
}
