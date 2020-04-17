<?php

declare(strict_types=1);

namespace Album\Tests\Application\Controller;

use Symfony\Component\HttpFoundation\Response;

/**
 * @group ui
 */
class AlbumControllerTest extends AbstractControllerTest
{
    public function testShouldGetTheLast3ThatContainMediasAndArePublic(): void
    {
        $response = $this->makeApiCall('GET', '/v1/albums', [], self::JWT_USER);
        $result = json_decode((string) $response->getContent(), true);

        self::assertEquals(200, $response->getStatusCode());
        self::assertCount(3, $result);
        self::assertEquals('My droids', $result[0]['title']);
        self::assertCount(1, $result[0]['favorites']);

        foreach ($result[0]['favorites'] as $favorite) {
            self::assertStringContainsString((string) getenv('IMG_PROXY'), $favorite);
        }
    }

    public function testShouldGetAlbumThatMatchSearchTerm(): void
    {
        $response = $this->makeApiCall('GET', '/v1/albums?search=everywhere', [], self::JWT_USER);
        $result = json_decode((string) $response->getContent(), true);

        self::assertEquals(200, $response->getStatusCode());
        self::assertCount(2, $result);
        self::assertEquals('Tatooine', $result[0]['title']);
        self::assertEquals('The Clone war', $result[1]['title']);
    }

    public function testShouldGetAnAlbumThatMatchTheSlug(): void
    {
        $response = $this->makeApiCall('GET', '/v1/album/jedi', [], self::JWT_USER);
        $result = json_decode((string) $response->getContent(), true);

        self::assertEquals(200, $response->getStatusCode());
        self::assertEquals('Jedi', $result['title']);
        self::assertCount(1, $result['medias']);

        foreach ($result['medias'] as $media) {
            self::assertStringContainsString((string) getenv('IMG_PROXY'), $media['uris']['small']);
            self::assertStringContainsString((string) getenv('IMG_PROXY'), $media['uris']['medium']);
            self::assertStringContainsString((string) getenv('IMG_PROXY'), $media['uris']['original']);
            self::assertArrayHasKey('type', $media);
        }
    }

    public function testShouldCreateAnAlbum(): void
    {
        $data = [
            'title' => 'Sith are the best',
            'description' => 'Come an join the Sith',
            'private' => true,
        ];

        $response = $this->makeApiCall('POST', '/v1/album', $data, self::JWT_ADMIN);
        $result = json_decode((string) $response->getContent(), true);

        self::assertEquals(200, $response->getStatusCode());
        self::assertEquals('sith-are-the-best', $result['slug']);
        self::assertTrue($result['private']);

        $toAssert = $this->findOneInDatabase('local-album', ['slug' => $result['slug']]);
        self::assertCount(7, $toAssert);
    }

    public function testShouldEditAnAlbum(): void
    {
        $data = [
            'title' => 'Jedi - fixed',
            'private' => false,
        ];

        $response = $this->makeApiCall('POST', '/v1/album/jedi', $data, self::JWT_ADMIN);
        $result = json_decode((string) $response->getContent(), true);

        self::assertEquals(200, $response->getStatusCode());
        self::assertEquals('Jedi - fixed', $result['title']);
        self::assertFalse($result['private']);

        $toAssert = $this->findOneInDatabase('local-album', ['slug' => $result['slug']]);
        self::assertEquals('Jedi - fixed', $toAssert['title']['S']);
        self::assertFalse($toAssert['isPrivate']['BOOL']);
    }

    public function testShouldAddMediasToAnAlbum(): void
    {
        $data = [
            [
                'key' => 'two',
                'author' => 'kenobi',
                'type' => 'image',
            ],
        ];

        $response = $this->makeApiCall('POST', '/v1/album/jedi/medias/add', $data, self::JWT_ADMIN);
        $result = json_decode((string) $response->getContent(), true);

        self::assertEquals(200, $response->getStatusCode());
        self::assertCount(2, $result['medias']);
        self::assertEquals('kenobi', $result['medias'][1]['author']);

        $toAssert = $this->findOneInDatabase('local-album', ['slug' => $result['slug']]);
        self::assertCount(2, $toAssert['medias']['L']);
        self::assertEquals('kenobi', $toAssert['medias']['L'][1]['M']['author']['S']);
    }

    public function testShouldRemoveMediasToAnAlbum(): void
    {
        $data = [
            [
                'key' => 'one',
                'author' => 'yoda',
                'type' => 'image',
            ],
        ];

        $response = $this->makeApiCall('POST', '/v1/album/jedi/medias/remove', $data, self::JWT_ADMIN);
        $result = json_decode((string) $response->getContent(), true);

        self::assertEquals(200, $response->getStatusCode());
        self::assertCount(0, $result['medias']);

        $toAssert = $this->findOneInDatabase('local-album', ['slug' => $result['slug']]);
        self::assertCount(0, $toAssert['medias']['L']);
    }

    public function testShouldGetAllAlbumsIfUserIsAnAdmin(): void
    {
        $response = $this->makeApiCall('GET', '/v1/albums?private=1&noMedias=1&limit=100', [], self::JWT_ADMIN);
        $result = json_decode((string) $response->getContent(), true);

        self::assertEquals(200, $response->getStatusCode());
        self::assertCount(6, $result);
    }

    public function testShouldReturnListOfAlbumsForAutocompleteThatMatchTerm(): void
    {
        $response = $this->makeApiCall('GET', '/v1/albums/autocomplete?search=everywhere', [], self::JWT_USER);
        $result = json_decode((string) $response->getContent(), true);

        self::assertEquals(200, $response->getStatusCode());
        self::assertCount(3, $result);
        self::assertEquals('Hoth', $result[0]['label']);
        self::assertEquals('hoth', $result[0]['value']);
        self::assertEquals('Tatooine', $result[1]['label']);
        self::assertEquals('tatooine', $result[1]['value']);
        self::assertEquals('The Clone war', $result[2]['label']);
        self::assertEquals('the-clone-war', $result[2]['value']);
    }

    public function testShouldReturnAdminResume(): void
    {
        $response = $this->makeApiCall('GET', '/v1/albums/resume', [], self::JWT_ADMIN);
        self::assertEquals(200, $response->getStatusCode());
        $toAssert = json_decode((string) $response->getContent(), true);

        self::assertArrayHasKey('publicCount', $toAssert);
        self::assertArrayHasKey('privateCount', $toAssert);
        self::assertEquals(5, $toAssert['publicCount']);
        self::assertEquals(1, $toAssert['privateCount']);
    }

    public function testShouldAddFavoriteToAlbum(): void
    {
        $data = [
            'favorite' => 'one',
        ];

        $response = $this->makeApiCall('PUT', '/v1/album/the-clone-war/favorite/add', $data, self::JWT_ADMIN);

        self::assertEquals(Response::HTTP_ACCEPTED, $response->getStatusCode());

        $toAssert = $this->findOneInDatabase('local-album', ['slug' => 'the-clone-war']);

        $favoriteCount = 0;
        foreach ($toAssert['medias']['L'] as $media) {
            if ($media['M']['favorite']['BOOL']) {
                $favoriteCount++;
            }
        }

        self::assertEquals(2, $favoriteCount);
    }

    public function testShouldRemoveAFavoriteFromAlbum(): void
    {
        $data = [
            'favorite' => 'two',
        ];

        $response = $this->makeApiCall('PUT', '/v1/album/the-clone-war/favorite/remove', $data, self::JWT_ADMIN);

        self::assertEquals(Response::HTTP_ACCEPTED, $response->getStatusCode());

        $toAssert = $this->findOneInDatabase('local-album', ['slug' => 'the-clone-war']);

        $favoriteCount = 0;
        foreach ($toAssert['medias']['L'] as $media) {
            if ($media['M']['favorite']['BOOL']) {
                $favoriteCount++;
            }
        }

        self::assertEquals(0, $favoriteCount);
    }

    public function testShouldDisplayMoreArticleThanTheFirst3(): void
    {
        $response = $this->makeApiCall('GET', '/v1/albums?limit=3&offset=3', [], self::JWT_USER);
        $toAssert = json_decode((string) $response->getContent(), true);

        self::assertEquals(200, $response->getStatusCode());
        self::assertEquals(1, count($toAssert));

        self::assertEquals('Sith', $toAssert[0]['title']);
    }

    public function testShouldNotCreateAnAlbumIfTitleIsAlreadyUsed(): void
    {
        $data = [
            'title' => 'The Clone war',
            'description' => 'Come an join the Sith',
            'private' => true,
        ];

        $toAssert = $this->makeApiCall('POST', '/v1/album', $data, self::JWT_ADMIN);
        self::assertEquals(500, $toAssert->getStatusCode());
    }
}
