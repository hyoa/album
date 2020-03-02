<?php

declare(strict_types=1);

namespace Album\Domain\Album;

interface AlbumRepositoryInterface
{
    public function insert(AlbumEntity $albumEntity): string;

    public function updateOne(AlbumEntity $albumEntity): void;

    public function findOne(array $data): ?AlbumEntity;

    public function find(bool $includePrivate, bool $includeNoMedia, int $limit, int $offset, ?string $term, string $order = 'asc'): array;

    public function deleteOne(array $filters): void;
}
