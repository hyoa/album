<?php

declare(strict_types=1);

namespace Album\Tests\Domain\Album;

use Album\Application\Clock\ClockInterface;
use Album\Domain\Album\AlbumEntity;
use Album\Domain\Album\AlbumManager;
use Album\Domain\Album\AlbumMediaEntity;
use Album\Domain\Album\AlbumRepositoryInterface;
use Album\Domain\Media\MediaEntity;
use Album\Domain\Media\MediaStorageInterface;
use Album\Tests\TestUtility\TestClock;
use PHPUnit\Framework\TestCase;
use Prophecy\Argument;
use Prophecy\Prophecy\ObjectProphecy;

class AlbumManagerTest extends TestCase
{
    /** @var ObjectProphecy|AlbumRepositoryInterface */
    protected ObjectProphecy $albumRepositoryMock;

    /** @var ObjectProphecy|MediaStorageInterface */
    protected ObjectProphecy $mediaStorageMock;

    public function setUp(): void
    {
        $this->albumRepositoryMock = $this->prophesize(AlbumRepositoryInterface::class);
        $this->mediaStorageMock = $this->prophesize(MediaStorageInterface::class);
    }

    public function testSaveShouldSaveAlbumInDatabase(): void
    {
        $clock = new TestClock();
        $this->albumRepositoryMock->insert(Argument::any())->willReturn('id');
        $manager = $this->buildAlbumManager($clock);

        $toAssert = $manager->save('title album', 'description', true, 'yoda');
        self::assertEquals('title album', $toAssert->title);
        self::assertEquals('description', $toAssert->description);
        self::assertEquals('id', $toAssert->id);
        self::assertEquals('title-album', $toAssert->slug);
        self::assertEquals('yoda', $toAssert->author);
        self::assertTrue($toAssert->private);
        self::assertEquals($clock->now(), $toAssert->creationDate);
    }

    public function testUpdateShouldUpdateAnAlbum(): void
    {
        $clock = new TestClock();
        $album = new AlbumEntity();
        $album->hydrate(
            [
                'id' => 'id',
                'title' => 'title album',
                'slug' => 'title-album',
                'description' => 'description',
                'private' => true,
                'author' => 'yoda',
            ]
        );
        $album->creationDate = $clock->now();

        $data = [
            'title' => 'title modified',
            'description' => 'description modified',
            'private' => false,
        ];

        $this->albumRepositoryMock->updateOne(Argument::any())->shouldBeCalledTimes(1);
        $manager = $this->buildAlbumManager($clock);
        $toAssert = $manager->updateOne($album, $data);
        self::assertEquals('title modified', $toAssert->title);
        self::assertEquals('description modified', $toAssert->description);
        self::assertEquals('id', $toAssert->id);
        self::assertEquals('title-modified', $toAssert->slug);
        self::assertEquals('yoda', $toAssert->author);
        self::assertFalse($toAssert->private);
        self::assertEquals($clock->now(), $toAssert->creationDate);
    }

    public function testAddMediaShouldAddMediaToAnAlbum(): void
    {
        $clock = new TestClock();

        $media = new MediaEntity();
        $media->key = 'one';
        $media->type = MediaEntity::TYPE_IMAGE;

        $media2 = new MediaEntity();
        $media2->key = 'two';
        $media2->type = MediaEntity::TYPE_VIDEO;

        $album = new AlbumEntity();
        $album->hydrate(
            [
                'id' => 'id',
                'title' => 'title album',
                'slug' => 'title-album',
                'description' => 'description',
                'private' => true,
                'author' => 'yoda',
            ]
        );
        $album->creationDate = $clock->now();

        $this->albumRepositoryMock->updateOne(Argument::any())->shouldBeCalledTimes(1);
        $this->mediaStorageMock->getUrisToAccessStore(Argument::any(), Argument::any())->willReturn([]);

        $manager = $this->buildAlbumManager($clock);

        $toAssert = $manager->addMedias($album, [$media, $media2]);
        self::assertCount(2, $toAssert->medias);
        self::assertEquals($media, $toAssert->medias[0]);
        self::assertEquals($media2, $toAssert->medias[1]);
    }

    public function testRemoveMediaShouldRemoveMediaFromAlbum(): void
    {
        $clock = new TestClock();

        $media = new MediaEntity();
        $media->key = 'one';
        $media->type = MediaEntity::TYPE_IMAGE;

        $media2 = new MediaEntity();
        $media2->key = 'two';
        $media2->type = MediaEntity::TYPE_VIDEO;

        $album = new AlbumEntity();
        $album->hydrate(
            [
                'id' => 'id',
                'title' => 'title album',
                'slug' => 'title-album',
                'description' => 'description',
                'private' => true,
                'author' => 'yoda',
            ]
        );
        $album->creationDate = $clock->now();
        $album->addMedia($media);
        $album->addMedia($media2);

        $this->albumRepositoryMock->updateOne(Argument::any())->shouldBeCalledTimes(1);
        $this->mediaStorageMock->getUrisToAccessStore(Argument::any(), Argument::any())->willReturn([]);

        $manager = $this->buildAlbumManager($clock);

        $toAssert = $manager->removeMedias($album, [$media2]);
        self::assertCount(1, $toAssert->medias);
        self::assertNotContains($media2, $toAssert->medias);
    }

    public function testDeleteAlbumShouldRemoveAnAlbum(): void
    {
        $this->albumRepositoryMock->deleteOne(Argument::any())->shouldBeCalled();

        $manager = $this->buildAlbumManager();

        $album = new AlbumEntity();
        $album->slug = 'slug';
        $manager->deleteOne($album);
    }

    public function testAddFavoriteShouldAddAFavoriteToAnAlbum(): void
    {
        $clock = new TestClock();

        $media = new AlbumMediaEntity();
        $media->key = 'one';
        $media->type = MediaEntity::TYPE_IMAGE;

        $media2 = new AlbumMediaEntity();
        $media2->key = 'two';
        $media2->type = MediaEntity::TYPE_VIDEO;

        $album = new AlbumEntity();
        $album->hydrate(
            [
                'id' => 'id',
                'title' => 'title album',
                'slug' => 'title-album',
                'description' => 'description',
                'private' => true,
                'author' => 'yoda',
                'medias' => [$media, $media2],
            ]
        );

        $this->albumRepositoryMock->updateOne(Argument::any())->shouldBeCalled();
        $this->mediaStorageMock->getUrisToAccessStore(Argument::any(), Argument::any())->willReturn([]);

        $manager = $this->buildAlbumManager($clock);
        $toAssert = $manager->toggleFavorite($album, 'one', true);

        self::assertCount(1, $toAssert->getFavorites());
        self::assertEquals('one', $toAssert->getFavorites()[0]->key);
    }

    public function testAddFavoriteShouldNotRemoveExistentFavorite(): void
    {
        $clock = new TestClock();

        $media = new AlbumMediaEntity();
        $media->key = 'one';
        $media->type = MediaEntity::TYPE_IMAGE;

        $media2 = new AlbumMediaEntity();
        $media2->key = 'two';
        $media2->type = MediaEntity::TYPE_VIDEO;
        $media2->isFavorite = true;

        $album = new AlbumEntity();
        $album->hydrate(
            [
                'id' => 'id',
                'title' => 'title album',
                'slug' => 'title-album',
                'description' => 'description',
                'private' => true,
                'author' => 'yoda',
                'medias' => [$media, $media2],
            ]
        );

        $this->albumRepositoryMock->updateOne(Argument::any())->shouldBeCalled();
        $this->mediaStorageMock->getUrisToAccessStore(Argument::any(), Argument::any())->willReturn([]);

        $manager = $this->buildAlbumManager($clock);

        $toAssert = $manager->toggleFavorite($album, 'one', true);
        self::assertCount(2, $toAssert->getFavorites());
    }

    public function testRemoveFavoriteShouldRemoveTheSelectedFavorite(): void
    {
        $clock = new TestClock();

        $media = new AlbumMediaEntity();
        $media->key = 'one';
        $media->type = MediaEntity::TYPE_IMAGE;

        $media2 = new AlbumMediaEntity();
        $media2->key = 'two';
        $media2->type = MediaEntity::TYPE_VIDEO;
        $media2->isFavorite = true;

        $album = new AlbumEntity();
        $album->hydrate(
            [
                'id' => 'id',
                'title' => 'title album',
                'slug' => 'title-album',
                'description' => 'description',
                'private' => true,
                'author' => 'yoda',
                'medias' => [$media, $media2],
            ]
        );

        $this->albumRepositoryMock->updateOne(Argument::any())->shouldBeCalled();
        $this->mediaStorageMock->getUrisToAccessStore(Argument::any(), Argument::any())->willReturn([]);

        $manager = $this->buildAlbumManager($clock);

        $toAssert = $manager->toggleFavorite($album, 'two', false);
        self::assertCount(0, $toAssert->getFavorites());
    }

    protected function buildAlbumManager(ClockInterface $clock = null): AlbumManager
    {
        if ($clock === null) {
            $clock = new TestClock();
        }

        return new AlbumManager($clock, $this->albumRepositoryMock->reveal(), $this->mediaStorageMock->reveal());
    }
}
