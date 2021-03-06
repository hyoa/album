<?php

declare(strict_types=1);

namespace Album\Application\VideoFormatter;

interface VideoFormatterInterface
{
    public function run(string $key, int $width = 720, int $height = 480): bool;
}
