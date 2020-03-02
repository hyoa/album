<?php

declare(strict_types=1);

namespace Album\Tests\TestUtility;

use Album\Application\Clock\ClockInterface;

class TestClock implements ClockInterface
{
    protected \DateTimeImmutable $fixedNow;

    public function __construct(\DateTimeImmutable $now = null)
    {
        $this->fixedNow = $now ?? new \DateTimeImmutable();
    }

    public function now(): \DateTimeImmutable
    {
        return $this->fixedNow;
    }
}
