<?php

declare(strict_types=1);

namespace Album\Tests\Application\Controller;

use Album\Domain\Media\MediaEntity;
use Album\Domain\Media\MediaRepositoryInterface;
use Prophecy\Prophecy\ObjectProphecy;
use Symfony\Component\HttpFoundation\Response;

/**
 * @group ui
 */
class MediaControllerTest extends AbstractControllerTest
{
    /** @var ObjectProphecy|MediaRepositoryInterface */
    protected ObjectProphecy $mediaRepositoryMock;

    public function setUp(): void
    {
        parent::setUp();
        $this->mediaRepositoryMock = $this->prophet->prophesize(MediaRepositoryInterface::class);
    }

    public function testFindMediasByFolderShouldReturnAllMediasWithTheGivenFolder(): void
    {
        $response = $this->makeApiCall('GET', '/v1/medias/folder/jedi', [], self::JWT_ADMIN);

        $toAssert = json_decode((string) $response->getContent(), true);

        self::assertEquals(200, $response->getStatusCode());
        self::assertCount(3, $toAssert);

        foreach ($toAssert as $media) {
            if ($media['type'] === MediaEntity::TYPE_IMAGE) {
                self::assertStringContainsString((string) getenv('PROXY_IMAGE'), $media['uris']['small']);
            } elseif ($media['type'] === MediaEntity::TYPE_VIDEO) {
                self::assertStringContainsString((string) getenv('ALBUM_BUCKET'), $media['uris']['small']);
            }
        }
    }

    public function testFindFoldersShouldReturnAListOfAllFolders(): void
    {
        $response = $this->makeApiCall('GET', '/v1/medias/folders', [], self::JWT_ADMIN);
        $toAssert = json_decode((string) $response->getContent(), true);

        self::assertEquals(200, $response->getStatusCode());
        self::assertCount(4, $toAssert);
        self::assertContains('jedi', $toAssert);
        self::assertContains('jedi school', $toAssert);
        self::assertContains('sith', $toAssert);
        self::assertContains('hoth', $toAssert);
    }

    public function testSignUriShouldReturnAPreSignedUriForAnImage(): void
    {
        $data = [
            'file' => 'example.jpg',
            'type' => 'image/jpeg',
        ];

        $response = $this->makeApiCall('POST', '/v1/media/signed-uri', $data, self::JWT_ADMIN);

        $toAssert = json_decode((string) $response->getContent(), true);

        self::assertEquals(200, $response->getStatusCode());
        self::assertStringContainsString('aws', $toAssert['uri']);
        self::assertStringContainsString('medias', $toAssert['uri']);
    }

    public function testSignUriShouldReturnAPreSignedUriForAVideo(): void
    {
        $data = [
            'file' => 'video.mp4',
            'type' => 'video/mp4',
        ];

        $response = $this->makeApiCall('POST', '/v1/media/signed-uri', $data, self::JWT_ADMIN);

        $toAssert = json_decode((string) $response->getContent(), true);

        self::assertEquals(200, $response->getStatusCode());
        self::assertStringContainsString('aws', $toAssert['uri']);
        self::assertStringContainsString('video', $toAssert['uri']);
    }

    public function testDeleteFolderShouldReturnASuccess(): void
    {
        $response = $this->makeApiCall('DELETE', '/v1/media/folder/jedi', [], self::JWT_ADMIN);
        self::assertEquals(Response::HTTP_ACCEPTED, $response->getStatusCode());

        $toAssert = $this->query(
            'local-media',
            'folder = :folder',
            [':folder' => 'none'],
            'folderIndex'
        );

        self::assertCount(3, $toAssert);
    }

    public function testChangeTheNameOfAFolderShouldUpdateFolderOfMedia(): void
    {
        $data = [
            'folderToUpdate' => 'jedi',
            'newFolderName' => 'jedi updated',
        ];

        $response = $this->makeApiCall('POST', '/v1/medias/folder/name', $data, self::JWT_ADMIN);
        self::assertEquals(Response::HTTP_ACCEPTED, $response->getStatusCode());

        $toAssert = $this->query(
            'local-media',
            'folder = :folder',
            [':folder' => $data['newFolderName']],
            'folderIndex'
        );

        self::assertCount(3, $toAssert);
    }

    public function testGetAdminResumeShouldReturnCountOfVideoAndImage(): void
    {
        $response = $this->makeApiCall('GET', '/v1/medias/resume', [], self::JWT_ADMIN);
        $toAssert = json_decode((string) $response->getContent(), true);

        self::assertEquals(200, $response->getStatusCode());
        self::assertCount(2, $toAssert);

        self::assertArrayHasKey('videosCount', $toAssert);
        self::assertArrayHasKey('imagesCount', $toAssert);
        self::assertEquals(1, $toAssert['videosCount']);
        self::assertEquals(5, $toAssert['imagesCount']);
    }

    public function testUpdateManyMediasFolderShouldChangeTheirFolder(): void
    {
        $data = [
            'folderName' => 'new jedi',
            'medias' => ['3.jpg', '2.jpg'],
        ];

        $response = $this->makeApiCall('POST', '/v1/medias/many/folder/name', $data, self::JWT_ADMIN);
        self::assertEquals(Response::HTTP_ACCEPTED, $response->getStatusCode());

        $toAssert = $this->query(
            'local-media',
            'folder = :folder',
            [':folder' => $data['folderName']],
            'folderIndex'
        );

        self::assertCount(2, $toAssert);
    }
}
