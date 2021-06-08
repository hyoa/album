<?php

declare(strict_types=1);

namespace Album\Application\VideoFormatter\CloudConverter\Task;

use CloudConvert\Models\Task;

class ExportTask
{
    const NAME = 'exportToS3';

    private function __construct()
    {
    }

    public static function createTask(
        string $accessKeyId,
        string $secretAccessKey,
        string $bucket,
        string $region
    ): Task {
        return (new Task('export/s3', self::NAME))
            ->set('input', ConvertTask::NAME)
            ->set('access_key_id', $accessKeyId)
            ->set('secret_access_key', $secretAccessKey)
            ->set('bucket', $bucket)
            ->set('region', $region)
        ;
    }
}
