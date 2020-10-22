<?php

declare(strict_types=1);

namespace Album\Application\VideoFormatter;

use FFMpeg\Coordinate\Dimension;
use FFMpeg\FFMpeg;
use FFMpeg\Filters\Video\VideoFilters;

class FFMPEGFormatter implements VideoFormatterInterface
{
    public function run(string $url, string $key, int $width = 720, int $height = 480): string
    {
        $ffmpeg = FFMpeg::create(['timeout' => 900]);
        $video = $ffmpeg->open($url);

        $format = new \FFMpeg\Format\Video\X264();
        $format->setAudioCodec('aac');

        /** @var VideoFilters $filters */
        $filters = $video->filters();
        $filters->resize(new Dimension($width, $height));
        $filters->synchronize();

        $pathToSave = '/tmp/'.$key;
        $video->save($format, $pathToSave);

        return $pathToSave;
    }
}
