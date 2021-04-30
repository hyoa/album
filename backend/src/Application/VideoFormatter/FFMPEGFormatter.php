<?php

declare(strict_types=1);

namespace Album\Application\VideoFormatter;

use Album\Domain\Media\Exception\InvalidVideoSizeException;
use Album\Domain\Media\MediaStorageInterface;
use FFMpeg\Coordinate\Dimension;
use FFMpeg\FFMpeg;
use FFMpeg\Filters\Video\VideoFilters;

class FFMPEGFormatter
{
    public function __construct(
        protected MediaStorageInterface $mediaStorage,
        protected string $mediaStorageLocation,
        protected string $videoRawStorageLocation
    ) {
    }

    public function run(string $key, int $width = 720, int $height = 480): bool
    {
        $fileSize = $this->mediaStorage->getFileSize($key, $this->videoRawStorageLocation);

        if ($fileSize > 200000000) {
            throw new InvalidVideoSizeException('Video size cannot exceed 200Mo');
        }

        $mediaUris = $this->mediaStorage->generateSignedUri($key, MediaStorageInterface::LOCATION_RAW_VIDEOS, 'GetObject');
        $pathVideoFormatted = $this->formatVideo($mediaUris, $key);

        $metadata = $this->mediaStorage->getMediaMetadata($key, $this->videoRawStorageLocation);

        $this->mediaStorage->putObject($key, $this->mediaStorageLocation, $pathVideoFormatted, 'video/mp4', $metadata);

        return true;
    }

    protected function formatVideo(string $url, string $key, int $width = 720, int $height = 480): string
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
