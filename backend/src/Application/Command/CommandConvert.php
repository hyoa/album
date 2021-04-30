<?php

declare(strict_types=1);

namespace Album\Application\Command;

use Album\Domain\Media\MediaIngestManager;
use Symfony\Component\Console\Command\Command;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Output\OutputInterface;

class CommandConvert extends Command
{
    protected static $defaultName = 'album:convert';
    protected MediaIngestManager $mediaIngestManager;

    public function __construct(
        MediaIngestManager $mediaIngestManager,
        string $name = null
    ) {
        parent::__construct($name);
        $this->mediaIngestManager = $mediaIngestManager;
    }

    protected function execute(InputInterface $input, OutputInterface $output): int
    {
        $this->mediaIngestManager->ingestVideo('Jules_video_videolarge.mp4');

        return 0;
    }
}
