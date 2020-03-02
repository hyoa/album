<?php

declare(strict_types=1);

namespace Album\Lambda;

use Bref\Context\Context;

interface LambdaHandlerInterface
{
    public function __invoke(Context $context, array $payload): array;
}
