<?php

declare(strict_types=1);

namespace Album\Application\Notification;

interface NotificationInterface
{
    const CHANNEL_ALBUM = 'album';
    const CHANNEL_ADMIN = 'admin';

    public function addTokenToChannel(string $token, string $channel): void;

    public function unsubscribeToken(string $token): void;

    public function sendMessageToChannel(string $title, string $message, string $link, string $channel): void;
}
