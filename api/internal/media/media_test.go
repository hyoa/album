package media_test

import (
	"testing"

	_media "github.com/hyoa/album/api/internal/media"
	_mocks "github.com/hyoa/album/api/mocks"
	"github.com/stretchr/testify/assert"
)

type mocks struct {
	mediaRepo *_mocks.MediaRepository
	storage   *_mocks.Storage
}

func getManagerWithMocks() (_media.MediaManager, mocks) {
	list := mocks{
		mediaRepo: new(_mocks.MediaRepository),
		storage:   new(_mocks.Storage),
	}

	return _media.CreateMediaManager(list.mediaRepo, list.storage), list
}

func TestItShouldGetMediaByFolders(t *testing.T) {
	manager, mocks := getManagerWithMocks()

	mocks.mediaRepo.On("FindByFolder", "foo").Return([]_media.Media{{Key: "1", Folder: "foo"}}, nil)

	mediasToAssert, err := manager.GetMediasByFolder("foo")

	assert.Len(t, mediasToAssert, 1)
	assert.Nil(t, err)
}

func TestItShouldGetFoldersNames(t *testing.T) {
	manager, mocks := getManagerWithMocks()

	mocks.mediaRepo.On("FindFoldersName", "").Return([]string{"foo"}, nil)

	foldersNamesToAssert, err := manager.GetFoldersNames("")

	assert.Len(t, foldersNamesToAssert, 1)
	assert.Equal(t, "foo", foldersNamesToAssert[0])
	assert.Nil(t, err)
}

func TestItShouldDeleteAFolder(t *testing.T) {
	manager, mocks := getManagerWithMocks()

	media := _media.Media{Key: "1", Folder: "foo", Author: "author", Kind: _media.KindPhoto}
	mediaAltered := media
	mediaAltered.Folder = "none"

	mocks.mediaRepo.On("FindByFolder", "foo").Return([]_media.Media{media}, nil)
	mocks.mediaRepo.On("Save", mediaAltered).Return(nil)

	err := manager.DeleteFolder("foo")
	assert.Nil(t, err)
}

func TestItShouldGetAllMedias(t *testing.T) {
	manager, mocks := getManagerWithMocks()

	media := _media.Media{Key: "1", Folder: "foo", Author: "author", Kind: _media.KindPhoto}

	mocks.mediaRepo.On("FindAll").Return([]_media.Media{media}, nil)

	mediasToAssert, err := manager.GetAll()

	assert.Len(t, mediasToAssert, 1)
	assert.Nil(t, err)
}

func TestItShouldChangeMediasFolder(t *testing.T) {
	manager, mocks := getManagerWithMocks()

	media := _media.Media{Key: "1", Folder: "foo", Author: "author", Kind: _media.KindPhoto}
	mediaAltered := media
	mediaAltered.Folder = "bar"

	mocks.mediaRepo.On("FindManyByKeys", []string{"key"}).Return([]_media.Media{media}, nil)
	mocks.mediaRepo.On("Save", mediaAltered).Return(nil)

	mediasToAssert, err := manager.ChangeMediasFolder([]string{"key"}, "bar")

	assert.Len(t, mediasToAssert, 1)
	assert.Equal(t, "bar", mediasToAssert[0].Folder)
	assert.Nil(t, err)
}

func TestItShouldChangeFolderName(t *testing.T) {
	manager, mocks := getManagerWithMocks()

	media := _media.Media{Key: "1", Folder: "foo", Author: "author", Kind: _media.KindPhoto}
	mediaAltered := media
	mediaAltered.Folder = "bar"

	mocks.mediaRepo.On("FindByFolder", "foo").Return([]_media.Media{media}, nil)
	mocks.mediaRepo.On("Save", mediaAltered).Return(nil)

	mediasToAssert, err := manager.ChangeFolderName("foo", "bar")

	assert.Len(t, mediasToAssert, 1)
	assert.Equal(t, "bar", mediasToAssert[0].Folder)
	assert.Nil(t, err)
}

func TestItShouldIngestMedia(t *testing.T) {
	manager, mocks := getManagerWithMocks()

	mocks.mediaRepo.On("FindByKey", "key").Return(_media.Media{}, nil)
	mocks.storage.On("MediaExist", "key").Return(true, nil)
	mocks.mediaRepo.On("Save", _media.Media{Key: "key", Author: "author", Kind: _media.KindPhoto, Folder: "folder"}).Return(nil)

	mediaToAssert, err := manager.Ingest("key", "author", "folder", _media.KindPhoto)

	assert.Equal(t, "key", mediaToAssert.Key)
	assert.Nil(t, err)
}

func TestItShouldReturnSignedUriForUpload(t *testing.T) {
	manager, mocks := getManagerWithMocks()

	mocks.storage.On("SignUploadUri", "key").Return("uri", nil)

	uriToAssert, err := manager.GetUploadSignedUri("key")
	assert.Equal(t, "uri", uriToAssert)
	assert.Nil(t, err)
}
