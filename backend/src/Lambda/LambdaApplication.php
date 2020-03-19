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
        $this->kernel = new Kernel($_SERVER['APP_ENV'], (bool) $_SERVER['APP_DEBUG']);
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
        $this->run($context, $event);

        return true;
    }
}

require dirname(__DIR__).'/../config/bootstrap.php';

return new LambdaApplication();
