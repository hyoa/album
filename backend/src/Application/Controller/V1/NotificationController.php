<?php

declare(strict_types=1);

namespace Album\Application\Controller\V1;

use Album\Application\Notification\NotificationInterface;
use Symfony\Bundle\FrameworkBundle\Controller\AbstractController;
use Symfony\Component\HttpFoundation\Request;
use Symfony\Component\HttpFoundation\Response;
use Symfony\Component\HttpKernel\Exception\HttpException;
use Symfony\Component\Routing\Annotation\Route;

/**
 * @Route("/v1")
 */
class NotificationController extends AbstractController
{
    /**
     * @Route("/notification/subscribe")
     */
    public function subscribe(Request $request, NotificationInterface $notification): Response
    {
        $data = json_decode((string) $request->getContent(), true);
        $token = (string) filter_var($data['token'], FILTER_SANITIZE_STRING);
        $channel = (string) filter_var($data['channel'], FILTER_SANITIZE_STRING);

        if ($token === '' || $channel === '') {
            throw new HttpException(500);
        }

        /* @var NotificationInterface $notification */
        $notification->addTokenToChannel($token, $channel);

        return new Response(null, Response::HTTP_ACCEPTED);
    }

    /**
     * @Route("/notification/unsubscribe")
     */
    public function unsubscribe(Request $request, NotificationInterface $notification): Response
    {
        $data = json_decode((string) $request->getContent(), true);
        $token = (string) filter_var($data['token'], FILTER_SANITIZE_STRING);

        if ($token === '') {
            throw new HttpException(500);
        }

        /* @var NotificationInterface $notification */
        $notification->unsubscribeToken($token);

        return new Response(null, Response::HTTP_ACCEPTED);
    }
}
