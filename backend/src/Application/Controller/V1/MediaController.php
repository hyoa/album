<?php

declare(strict_types=1);

namespace Album\Application\Controller\V1;

use Album\Domain\Media\MediaEntity;
use Album\Domain\Media\MediaManager;
use Album\Domain\Media\MediaStorageInterface;
use Sensio\Bundle\FrameworkExtraBundle\Configuration\IsGranted;
use Symfony\Bundle\FrameworkBundle\Controller\AbstractController;
use Symfony\Component\HttpFoundation\JsonResponse;
use Symfony\Component\HttpFoundation\Request;
use Symfony\Component\HttpFoundation\Response;
use Symfony\Component\Routing\Annotation\Route;

/**
 * @Route("/v1")
 */
class MediaController extends AbstractController
{
    /**
     * @Route("/medias/folder/{folderName}", methods={"GET"})
     * @IsGranted("ROLE_ADMIN")
     */
    public function getByFolder(MediaManager $mediaManager, string $folderName): Response
    {
        $mediasCollection = $mediaManager->findByFolder($folderName);

        $mediasVisible = array_filter(
            $mediasCollection,
            fn (MediaEntity $mediaEntity) => $mediaEntity->visible
        );

        $medias = array_map(function (MediaEntity $mediaEntity): array {
            return $mediaEntity->getAsArray();
        }, $mediasVisible);

        return new JsonResponse($medias);
    }

    /**
     * @Route("/medias/folders", methods={"GET"})
     * @IsGranted("ROLE_ADMIN")
     */
    public function getFolders(MediaManager $mediaManager): Response
    {
        return new JsonResponse($mediaManager->findFolders());
    }

    /**
     * @Route("/medias/folders/autocomplete", methods={"GET"})
     * @IsGranted("ROLE_ADMIN")
     */
    public function getFoldersAutocomplete(Request $request, MediaManager $mediaManager): Response
    {
        $folders = array_map(function ($folder): array {
            return [
                'label' => $folder,
                'value' => $folder,
            ];
        }, $mediaManager->findFolders($request->query->get('search', null)));

        return new JsonResponse($folders);
    }

    /**
     * @Route("/media/signed-uri", methods={"POST"})
     * @IsGranted("ROLE_ADMIN")
     */
    public function getSignedUri(Request $request, MediaStorageInterface $mediaStorage, string $bucketVideoInput): Response
    {
        $data = json_decode((string) $request->getContent(), true);
        $file = $data['file'];
        $type = $data['type'];

        $location = MediaStorageInterface::LOCATION_MEDIAS;

        if (strpos($type, 'video') !== false) {
            $location = MediaStorageInterface::LOCATION_RAW_VIDEOS;
        }

        $signedUri = $mediaStorage->generateSignedUri($file, $location, 'PutObject');

        return new JsonResponse(['uri' => $signedUri]);
    }

    /**
     * @Route("/media/folder/{folderName}", methods={"DELETE"})
     * @IsGranted("ROLE_ADMIN")
     */
    public function deleteFolder(MediaManager $mediaManager, string $folderName): Response
    {
        $mediaManager->deleteFolder($folderName);

        return new Response(null, Response::HTTP_ACCEPTED);
    }

    /**
     * @Route("/medias/folder/name", methods={"POST"})
     * @IsGranted("ROLE_ADMIN")
     */
    public function updateFolderName(Request $request, MediaManager $mediaManager): Response
    {
        $data = json_decode((string) $request->getContent(), true);
        $folderToUpdate = $data['folderToUpdate'];
        $newFolderName = $data['newFolderName'];

        $mediaManager->updateFolderName($folderToUpdate, $newFolderName);

        return new Response(null, Response::HTTP_ACCEPTED);
    }

    /**
     * @Route("/medias/resume")
     * @IsGranted("ROLE_ADMIN")
     */
    public function getAdminResume(MediaManager $mediaManager): Response
    {
        return new JsonResponse($mediaManager->getAdminResume());
    }

    /**
     * @Route("/medias/many/folder/name")
     * @IsGranted("ROLE_ADMIN")
     */
    public function updateFolderNameForManyMedias(Request $request, MediaManager $mediaManager): Response
    {
        $data = json_decode((string) $request->getContent(), true);
        $folderName = $data['folderName'];
        $medias = $data['medias'];

        $mediaManager->updateFolderNameForMedias($folderName, $medias);

        return new Response(null, Response::HTTP_ACCEPTED);
    }
}
