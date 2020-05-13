<?php

declare(strict_types=1);

namespace Album\Application\Controller\V2;

use Album\Domain\Album\AlbumManager;
use Album\Domain\Media\MediaStorageInterface;
use Ausi\SlugGenerator\SlugGenerator;
use Sensio\Bundle\FrameworkExtraBundle\Configuration\IsGranted;
use Symfony\Bundle\FrameworkBundle\Controller\AbstractController;
use Symfony\Component\HttpFoundation\JsonResponse;
use Symfony\Component\HttpFoundation\Request;
use Symfony\Component\HttpFoundation\Response;
use Symfony\Component\Routing\Annotation\Route;

/**
 * @Route("/v2")
 */
class MediaController extends AbstractController
{
    /**
     * @Route("/medias/signed-uri", methods={"POST"})
     * @IsGranted("ROLE_ADMIN")
     */
    public function getSignedUri(Request $request, MediaStorageInterface $mediaStorage): Response
    {
        $data = json_decode((string) $request->getContent(), true);
        $signedUris = [];

        foreach ($data as $item) {
            $file = $item['file'];
            $type = $item['type'];

            $location = MediaStorageInterface::LOCATION_MEDIAS;

            if (strpos($type, 'video') !== false) {
                $location = MediaStorageInterface::LOCATION_RAW_VIDEOS;
            }

            $signedUris[] = [
                'uri' => $mediaStorage->generateSignedUri($file, $location, 'PutObject'),
                'key' => $file,
            ];
        }

        return new JsonResponse($signedUris);
    }

    /**
     * @Route("/medias/ingest", methods={"POST"})
     * @IsGranted("ROLE_ADMIN")
     */
    public function ingest(Request $request, AlbumManager $albumManager, MediaStorageInterface $mediaStorage): Response
    {
        $data = json_decode((string) $request->getContent(), true);

        if (isset($data['album'])) {
            $albumName = $data['album']['title'] ?? null;

            if ($albumName !== null && trim($albumName) !== '') {
                $slug = (new SlugGenerator())->generate((string) $albumName);
                $album = $albumManager->findBySlug($slug);

                if ($album === null) {
                    $albumManager->save((string) $albumName, null, true, (string) $data['album']['author'] ?? null);
                }
            }
        }

        $signedUris = [];
        foreach ($data['medias'] as $item) {
            $file = $item['file'];
            $type = $item['type'];
            $metadata = $item['metadata'];

            $location = MediaStorageInterface::LOCATION_MEDIAS;

            if (strpos($type, 'video') !== false) {
                $location = MediaStorageInterface::LOCATION_RAW_VIDEOS;
            }

            $signedUris[] = [
                'uri' => $mediaStorage->generateSignedUri($file, $location, 'PutObject', $metadata),
                'key' => $file,
            ];
        }

        return new JsonResponse($signedUris);
    }
}
