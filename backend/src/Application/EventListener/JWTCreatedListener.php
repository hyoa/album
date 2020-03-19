<?php

declare(strict_types=1);

namespace Album\Application\EventListener;

use Album\Application\Security\UserSecurity;
use Lexik\Bundle\JWTAuthenticationBundle\Event\JWTCreatedEvent;

class JWTCreatedListener
{
    public function onJWTCreated(JWTCreatedEvent $event): void
    {
        $payload = $event->getData();

        /** @var UserSecurity $user */
        $user = $event->getUser();

        $payload['name'] = $user->name;
        $payload['role'] = $user->role;

        $event->setData($payload);
    }
}
