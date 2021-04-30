<?php

declare(strict_types=1);

namespace Album\Application\VideoFormatter\CloudConverter;

use Album\Application\Storage\S3Storage;
use Album\Application\VideoFormatter\CloudConverter\Task\ConvertTask;
use Album\Application\VideoFormatter\CloudConverter\Task\ExportTask;
use Album\Application\VideoFormatter\CloudConverter\Task\ImportTask;
use Album\Application\VideoFormatter\VideoFormatterInterface;
use CloudConvert\CloudConvert;
use CloudConvert\Models\Job;

class Formatter implements VideoFormatterInterface
{
    protected CloudConvert $transcoder;
    protected S3Storage $s3Storage;

    public function __construct(
        S3Storage $s3Storage,
        string $apiKey,
        bool $sandbox
    ) {
        dump($sandbox);
        $this->transcoder = new CloudConvert([
            'api_key' => $apiKey,
            'sandbox' => $sandbox,
        ]);
        $this->s3Storage = $s3Storage;
    }

    public function run(string $key, int $width = 720, int $height = 480): bool
    {
        $credentials = $this->s3Storage->getCredentials(S3Storage::LOCATION_RAW_VIDEOS);
        $importTask = ImportTask::createTask(
            $credentials['accessKeyId'],
            $credentials['secretAccessKey'],
            $this->s3Storage->getBucket(S3Storage::LOCATION_RAW_VIDEOS),
            $this->s3Storage->getMediaStorageRegion(),
            $key
        );

        $convertTask = ConvertTask::createTask();

        $credentials = $this->s3Storage->getCredentials(S3Storage::LOCATION_MEDIAS);
        $exportTask = ExportTask::createTask(
            $credentials['accessKeyId'],
            $credentials['secretAccessKey'],
            $this->s3Storage->getBucket(S3Storage::LOCATION_MEDIAS),
            $this->s3Storage->getMediaStorageRegion()
        );

        $job = (new Job())
            ->addTask($importTask)
            ->addTask($convertTask)
            ->addTask($exportTask)
        ;

        $this->transcoder->jobs()->create($job);

        return false;
    }
}
