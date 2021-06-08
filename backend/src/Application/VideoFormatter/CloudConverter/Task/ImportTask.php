<?php

declare(strict_types=1);

namespace Album\Application\VideoFormatter\CloudConverter\Task;

use CloudConvert\Models\Task;

class ImportTask
{
    const NAME = 'importFromS3';

    private function __construct()
    {
    }

    public static function createTask(
        string $accessKeyId,
        string $secretAccessKey,
        string $bucket,
        string $region,
        string $key
    ): Task {
        return (new Task('import/s3', self::NAME))
            ->set('access_key_id', $accessKeyId)
            ->set('secret_access_key', $secretAccessKey)
            ->set('bucket', $bucket)
            ->set('region', $region)
            ->set('key', $key)
        ;
    }
}
