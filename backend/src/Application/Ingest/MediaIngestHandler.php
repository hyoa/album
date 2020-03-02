<?php

declare(strict_types=1);

namespace Album\Application\Ingest;

use Album\Domain\Media\MediaIngestManager;
use Album\Lambda\LambdaHandlerInterface;
use Bref\Context\Context;

class MediaIngestHandler implements LambdaHandlerInterface
{
    protected MediaIngestManager $ingestManager;

    public function __construct(MediaIngestManager $ingestManager)
    {
        $this->ingestManager = $ingestManager;
    }

    public function __invoke(Context $context, array $payload): array
    {
        $keys = [];
        foreach ($payload['Records'] as $media) {
            $key = $media['s3']['object']['key'];
            $keys[] = $key;
            $this->ingestManager->ingest($key);
        }

        return $keys;
    }
}
