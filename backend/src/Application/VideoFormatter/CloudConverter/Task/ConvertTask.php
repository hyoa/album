<?php

declare(strict_types=1);

namespace Album\Application\VideoFormatter\CloudConverter\Task;

use CloudConvert\Models\Task;

class ConvertTask
{
    const NAME = 'convertVideo';

    private function __construct()
    {
    }

    public static function createTask(): Task
    {
        return (new Task('convert', self::NAME))
            ->set('input', ImportTask::NAME)
            ->set('input_format', 'mp4')
            ->set('output_format', 'mp4')
            ->set('engine', 'ffmpeg')
            ->set('video_codec', 'x264')
            ->set('crf', '23')
            ->set('preset', 'medium')
            ->set('subtitles_mode', 'none')
            ->set('audio_codec', 'aac')
            ->set('audio_bitrate', 128)
            ->set('width', 1280)
            ->set('height', 720)
        ;
    }
}
