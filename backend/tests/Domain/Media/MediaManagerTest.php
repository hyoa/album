<?php

declare(strict_types=1);

namespace Album\Tests\Domain\Media;

use Album\Application\Clock\ClockInterface;
use Album\Domain\Media\MediaEntity;
use Album\Domain\Media\MediaManager;
use Album\Domain\Media\MediaRepositoryInterface;
use Album\Domain\Media\MediaStorageInterface;
use Album\Tests\TestUtility\TestClock;
use PHPUnit\Framework\TestCase;
use Prophecy\Argument;
use Prophecy\Prophecy\ObjectProphecy;
use Prophecy\Prophet;

/**
 * @group unit
 */
class MediaManagerTest extends TestCase
{
    protected Prophet $prophet;

    /** @var ObjectProphecy|MediaRepositoryInterface */
    protected ObjectProphecy $mediaRepositoryMock;

    /** @var ObjectProphecy|MediaStorageInterface */
    protected ObjectProphecy $mediaStorageMock;

    public function setUp(): void
    {
        $this->prophet = new Prophet();
        $this->mediaRepositoryMock = $this->prophet->prophesize(MediaRepositoryInterface::class);
        $this->mediaStorageMock = $this->prophet->prophesize(MediaStorageInterface::class);
    }

    public function testSaveMediaShouldCreateMediaInDatabase(): void
    {
        $this->mediaRepositoryMock->insert(Argument::any())->shouldBeCalledTimes(1);

        $clock = new TestClock();
        $manager = $this->buildMediaManager($clock);
        $toAssert = $manager->save('key', 'author', MediaEntity::TYPE_IMAGE, 'folder');

        self::assertEquals('key', $toAssert->key);
        self::assertEquals('author', $toAssert->author);
        self::assertEquals(MediaEntity::TYPE_IMAGE, $toAssert->type);
        self::assertEquals('folder', $toAssert->folder);
        self::assertEquals($clock->now(), $toAssert->uploadDate);
    }

    public function testFindByFolderShouldReturnAllMediaInAFolder(): void
    {
        $media1 = new MediaEntity();
        $media1->folder = 'folder';
        $media1->key = 'one';
        $media1->type = MediaEntity::TYPE_IMAGE;

        $media2 = new MediaEntity();
        $media2->folder = 'folder';
        $media2->key = 'two';
        $media2->type = MediaEntity::TYPE_VIDEO;

        $this->mediaRepositoryMock->findByFolder('folder')->willReturn([$media1, $media2]);
        $this->mediaStorageMock->getUrisToAccessStore(Argument::any(), Argument::any())->willReturn([
            'small' => 'http://small',
            'medium' => 'http://medium',
            'original' => 'http://original',
        ]);
        $this->mediaStorageMock->getUrisToAccessStore(Argument::any(), Argument::any())->shouldBeCalledTimes(2);

        $manager = $this->buildMediaManager();

        /** @var MediaEntity[] $toAssert */
        $toAssert = $manager->findByFolder('folder');

        self::assertCount(2, $toAssert);
        self::assertEquals('http://small', $toAssert[0]->getMediaUri('small'));
    }

    public function testFindFoldersShouldReturnAllFolders(): void
    {
        $this->mediaRepositoryMock->findFolders(Argument::any())->willReturn(['jedi', 'sith']);

        $manager = $this->buildMediaManager();

        $toAssert = $manager->findFolders();
        self::assertSame(['jedi', 'sith'], $toAssert);
    }

    public function testDeleteFolderShouldRemoveFolderOfItsMedia(): void
    {
        $this->mediaRepositoryMock->update(['folder' => 'folder'], ['folder' => 'none'])->shouldBeCalledTimes(1);

        $manager = $this->buildMediaManager();

        self::assertTrue($manager->deleteFolder('folder'));
    }

    public function tearDown(): void
    {
        $this->prophet->checkPredictions();
        parent::tearDown();
    }

    protected function buildMediaManager(ClockInterface $clock = null): MediaManager
    {
        if (null === $clock) {
            $clock = new TestClock();
        }

        return new MediaManager(
            $clock,
            $this->mediaRepositoryMock->reveal(),
            $this->mediaStorageMock->reveal()
        );
    }
}
