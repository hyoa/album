<?php

declare(strict_types=1);

namespace Album\Application\Notification;

use Album\Application\Notification\Exception\NotificationSendException;
use Album\Application\Notification\Exception\NotificationSubscriptionException;
use Symfony\Contracts\HttpClient\HttpClientInterface;

class FireBaseNotification implements NotificationInterface
{
    protected HttpClientInterface $httpClient;

    protected string $firebaseKey;

    protected string $channelSuffix;

    public function __construct(HttpClientInterface $httpClient, string $firebaseKey, string $channelSuffix)
    {
        $this->httpClient = $httpClient;
        $this->firebaseKey = $firebaseKey;
        $this->channelSuffix = $channelSuffix;
    }

    public function addTokenToChannel(string $token, string $channel): void
    {
        $url = sprintf(
            'https://iid.googleapis.com/iid/v1/%s/rel/topics/%s-%s',
            $token,
            $channel,
            (string) getenv('FIREBASE_CHANNEL_SUFFIX')
        );

        $response = $this->httpClient->request(
            'POST',
            $url,
            [
                'headers' => [
                    'Content-Type' => 'application/json',
                    'Authorization' => sprintf('key=%s', $this->firebaseKey),
                ],
            ]
        );

        if (200 !== $response->getStatusCode()) {
            throw new NotificationSubscriptionException(sprintf('Invalid subscription for token %s on channel %s, %s', $token, $channel, $response->getContent()));
        }
    }

    public function unsubscribeToken(string $token): void
    {
        $url = sprintf(
            'https://iid.googleapis.com/v1/web/iid/%s',
            $token
        );

        $response = $this->httpClient->request(
            'DELETE',
            $url,
            [
                'headers' => [
                    'Content-Type' => 'application/json',
                    'Authorization' => sprintf('key=%s', $this->firebaseKey),
                ],
            ]
        );

        if (200 !== $response->getStatusCode()) {
            throw new NotificationSubscriptionException(sprintf('Invalid unsubscription for token %s', $token));
        }
    }

    public function sendMessageToChannel(string $title, string $message, string $link, string $channel): void
    {
        $parameters = [
            'notification' => [
                'title' => $title,
                'body' => $message,
                'click_action' => $link,
                'icon' => 'http://url-to-an-icon/icon.png',
            ],
            'to' => sprintf('/topics/%s-%s', $channel, $this->channelSuffix),
        ];

        $response = $this->httpClient->request(
            'POST',
            'https://fcm.googleapis.com/fcm/send',
            [
                'headers' => [
                    'Content-Type' => 'application/json',
                    'Authorization' => sprintf('key=%s', $this->firebaseKey),
                ],
                'json' => $parameters,
            ]
        );

        if (200 !== $response->getStatusCode()) {
            throw new NotificationSendException(sprintf('Invalid message send %s', $response->getContent()));
        }
    }
}
