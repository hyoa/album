<?php

declare(strict_types=1);

namespace Album\Application\Controller\V1;

use Album\Application\Helper\EnumErrorCodeApi;
use Album\Application\Helper\JWTHelper;
use Album\Application\Notification\NotificationInterface;
use Album\Domain\Album\AlbumEntity;
use Album\Domain\Album\AlbumManager;
use Album\Domain\Album\AlbumMediaEntity;
use Album\Domain\Album\AlbumRepositoryInterface;
use Album\Domain\Media\MediaEntity;
use Album\Domain\Media\MediaStorageInterface;
use Sensio\Bundle\FrameworkExtraBundle\Configuration\IsGranted;
use Symfony\Bundle\FrameworkBundle\Controller\AbstractController;
use Symfony\Component\HttpFoundation\JsonResponse;
use Symfony\Component\HttpFoundation\Request;
use Symfony\Component\HttpFoundation\Response;
use Symfony\Component\HttpKernel\Exception\NotFoundHttpException;
use Symfony\Component\Routing\Annotation\Route;

/**
 * @Route("/v1")
 */
class AlbumController extends AbstractController
{
    /**
     * @Route("/albums", methods={"GET"})
     * @IsGranted("ROLE_USER")
     */
    public function getAlbums(Request $request, AlbumManager $albumManager, MediaStorageInterface $mediaStorage): Response
    {
        $searchTerm = $request->query->get('search', null);

        $includePrivateAlbum = (bool) $request->query->get('private', false);
        $includeNoMedias = (bool) $request->query->get('noMedias', false);

        $limit = (int) $request->query->get('limit', 3);
        $offset = (int) $request->query->get('offset', 0);

        $albums = $albumManager->findMany($includePrivateAlbum, $includeNoMedias, $limit, $offset, $searchTerm, 'desc');

        $dataResponse = [];

        /** @var AlbumEntity $album */
        foreach ($albums as $album) {
            $favorites = [];

            if (count($album->getFavorites()) > 0) {
                /** @var AlbumMediaEntity $favorite */
                foreach ($album->getFavorites() as $favorite) {
                    $favorites[] = $mediaStorage->getUrisToAccessStore($favorite->key)['small'];
                }
            } elseif (count($album->medias) > 0) {
                /** @var AlbumMediaEntity $media */
                foreach ($album->medias as $media) {
                    if ($media->type === MediaEntity::TYPE_IMAGE) {
                        $favorites[] = $mediaStorage->getUrisToAccessStore($media->key)['small'];
                        break;
                    }
                }
            }

            $dataResponse[] = [
                'title' => $album->title,
                'slug' => $album->slug,
                'description' => $album->description,
                'favorites' => $favorites,
                'author' => $album->author,
                'medias' => $album->medias,
            ];
        }

        return new JsonResponse($dataResponse);
    }

    /**
     * @Route("/albums/autocomplete", methods={"GET"})
     * @IsGranted("ROLE_USER")
     */
    public function getAlbumsAutocomplete(Request $request, AlbumManager $albumManager): Response
    {
        $searchTerm = $request->query->get('search', null);
        $limit = (int) $request->query->get('limit', 1000);

        /** @var AlbumEntity[] $albums */
        $albums = $albumManager->findMany(true, true, $limit, 0, $searchTerm, 'desc');

        $dataResponse = [];
        foreach ($albums as $album) {
            $dataResponse[] = [
                'label' => $album->title,
                'value' => $album->slug,
            ];
        }

        return new JsonResponse($dataResponse);
    }

    /**
     * @Route("/album/{slug}", methods={"GET"})
     * @IsGranted("ROLE_USER")
     */
    public function getAlbum(AlbumManager $albumManager, string $slug): Response
    {
        $album = $albumManager->findBySlug($slug);

        if ($album === null) {
            throw new NotFoundHttpException();
        }

        $albumArray = $album->getAsArray();

        return new JsonResponse($albumArray);
    }

    /**
     * @Route("/album", methods={"POST"})
     * @IsGranted("ROLE_ADMIN")
     */
    public function addAlbum(
        Request $request,
        AlbumManager $albumManager,
        AlbumRepositoryInterface $albumRepository,
        JWTHelper $JWTHelper
    ): Response {
        $data = json_decode((string) $request->getContent(), true);
        $title = (string) filter_var($data['title'], FILTER_SANITIZE_STRING);

        if (trim($title) === '') {
            return new JsonResponse(['code' => EnumErrorCodeApi::ERROR_ALBUM_INVALID_DATA], 500);
        }

        $album = $albumRepository->findOne(['title' => $title]);
        if (!is_null($album)) {
            return new JsonResponse(['code' => EnumErrorCodeApi::ERROR_ALBUM_ALREADY_EXIST], Response::HTTP_INTERNAL_SERVER_ERROR);
        }

        $bearerToken = (string) $request->headers->get('Authorization');
        $token = str_replace('Bearer ', '', $bearerToken);
        $author = $JWTHelper->getData($token, 'name');

        $description = (string) filter_var($data['description'], FILTER_SANITIZE_STRING);
        $private = (bool) filter_var($data['private'], FILTER_SANITIZE_STRING);

        $album = $albumManager->save($title, $description, $private, (string) $author);

        return new JsonResponse([
            'title' => $album->title,
            'description' => $album->description,
            'private' => $album->private,
            'slug' => $album->slug,
            'author' => $album->author,
        ]);
    }

    /**
     * @Route("/album/{slug}", methods={"POST"})
     * @IsGranted("ROLE_ADMIN")
     */
    public function editAlbum(Request $request, AlbumManager $albumManager, string $slug): Response
    {
        $data = json_decode((string) $request->getContent(), true);

        $album = $albumManager->findBySlug($slug);

        if ($album === null) {
            throw new NotFoundHttpException();
        }

        $album = $albumManager->updateOne($album, $data);

        return new JsonResponse([
            'title' => $album->title,
            'description' => $album->description,
            'private' => $album->private,
            'slug' => $album->slug,
            'author' => $album->author,
        ]);
    }

    /**
     * @Route("/album/{slug}/medias/{type}", methods={"POST"})
     * @IsGranted("ROLE_ADMIN")
     */
    public function updateMedias(
        Request $request,
        AlbumManager $albumManager,
        NotificationInterface $notification,
        string $slug,
        string $type
    ): Response {
        $data = json_decode((string) $request->getContent(), true);
        $dispatchNotification = false;
        $album = $albumManager->findBySlug($slug);

        if ($album === null) {
            throw new NotFoundHttpException();
        }

        $medias = [];
        foreach ($data as $item) {
            $media = new AlbumMediaEntity();
            $media->key = $item['key'];
            $media->author = $item['author'];
            $media->setTypeFromString($item['type']);

            $medias[] = $media;
        }

        if ($type === 'add' && count($medias) > 0 && count($album->medias) === 0) {
            $dispatchNotification = true;
        }

        if ($type === 'add') {
            $album = $albumManager->addMedias($album, $medias);
        } else {
            $album = $albumManager->removeMedias($album, $medias);
        }

        if ($dispatchNotification) {
            $notification->sendMessageToChannel(
                'Un nouvel album vient d\'être créé !',
                sprintf('%s vient d\'être ajouté par %s', $album->title, $album->author),
                sprintf('%s/#/album/%s', (string) getenv('APP_URI'), $album->slug),
                NotificationInterface::CHANNEL_ALBUM
            );
        }

        return new JsonResponse($album->getAsArray());
    }

    /**
     * @Route("/album/{slug}", methods={"DELETE"})
     * @IsGranted("ROLE_ADMIN")
     */
    public function delete(AlbumManager $albumManager, string $slug): Response
    {
        $album = $albumManager->findBySlug($slug);

        if ($album === null) {
            throw new NotFoundHttpException();
        }

        $albumManager->deleteOne($album);

        return new Response(null, Response::HTTP_ACCEPTED);
    }

    /**
     * @Route("/albums/resume", methods={"GET"})
     * @IsGranted("ROLE_ADMIN")
     */
    public function adminResume(AlbumManager $albumManager): Response
    {
        return new JsonResponse($albumManager->getAdminResume());
    }

    /**
     * @Route("/album/{slug}/favorite/add", methods={"PUT"})
     * @IsGranted("ROLE_ADMIN")
     */
    public function addFavorite(Request $request, AlbumManager $albumManager, string $slug): Response
    {
        $data = json_decode((string) $request->getContent(), true);

        if (!isset($data['favorite'])) {
            return new JsonResponse(['error' => 'favorite should be set'], Response::HTTP_INTERNAL_SERVER_ERROR);
        }

        $album = $albumManager->findBySlug($slug);

        if ($album === null) {
            throw new NotFoundHttpException();
        }

        $albumManager->toggleFavorite($album, $data['favorite'], true);

        return new Response(null, Response::HTTP_ACCEPTED);
    }

    /**
     * @Route("/album/{slug}/favorite/remove", methods={"PUT"})
     * @IsGranted("ROLE_ADMIN")
     */
    public function removeFavorite(Request $request, AlbumManager $albumManager, string $slug): Response
    {
        $data = json_decode((string) $request->getContent(), true);

        if (!isset($data['favorite'])) {
            return new JsonResponse(['error' => 'favorite should be set'], Response::HTTP_INTERNAL_SERVER_ERROR);
        }

        $album = $albumManager->findBySlug($slug);

        if ($album === null) {
            throw new NotFoundHttpException();
        }

        $albumManager->toggleFavorite($album, $data['favorite'], false);

        return new Response(null, Response::HTTP_ACCEPTED);
    }
}
