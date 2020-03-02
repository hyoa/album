<?php

declare(strict_types=1);

namespace Album\Domain\Album;

use Album\Domain\Media\MediaEntity;

class AlbumMediaEntity extends MediaEntity
{
    public bool $isFavorite = false;

    public function hydrate(array $data): void
    {
        parent::hydrate($data);

        if (array_key_exists('favorite', $data)) {
            $this->isFavorite = $data['favorite'];
        }
    }
}
