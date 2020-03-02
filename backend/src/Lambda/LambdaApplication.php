<?php

declare(strict_types=1);

namespace Album\Lambda;

use Album\Kernel;
use Bref\Context\Context;
use Bref\Event\Handler;

class LambdaApplication implements Handler
{
    protected Kernel $kernel;

    public function __construct()
    {
        $this->kernel = new Kernel('dev', true);
    }

    public function run(Context $context, array $payload): array
    {
        $this->kernel->boot();

        $container = $this->kernel->getContainer();

        /** @var LambdaHandler $handler */
        $handler = $container->get('Album\Lambda\LambdaHandler');

        return $handler($context, $payload);
    }

    public function handle($event, Context $context)
    {
        return $this->run($context, $event);
    }
}

require dirname(__DIR__).'/../config/bootstrap.php';

return new LambdaApplication();
