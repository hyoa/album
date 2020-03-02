<?php

declare(strict_types=1);

namespace Album\Application\Clock;

interface ClockInterface
{
    public function now(): \DateTimeImmutable;
}
