<?php

declare(strict_types=1);

namespace Album\Application\Clock;

class SystemClock implements ClockInterface
{
    public function now(): \DateTimeImmutable
    {
        return new \DateTimeImmutable();
    }
}
