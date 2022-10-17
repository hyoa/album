package album_test

import (
	"errors"
	"testing"

	"github.com/hyoa/album/api/internal/album"
	_album "github.com/hyoa/album/api/internal/album"
	_mocks "github.com/hyoa/album/api/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mocks struct {
	albumRepo *_mocks.AlbumRepository
}

func getManagerWithMocks() (_album.AlbumManager, mocks) {
	list := mocks{
		albumRepo: new(_mocks.AlbumRepository),
	}

	return _album.CreateAlbumManager(list.albumRepo), list
}

func TestItShouldCreateAnAlbumIfEverythingIsOk(t *testing.T) {
	manager, mocks := getManagerWithMocks()

	album := _album.Album{
		Title:       "title",
		Description: "description",
		Private:     false,
		Author:      "me",
		Id:          "1",
		Slug:        "title",
		Medias:      make([]_album.Media, 0),
	}
	mocks.albumRepo.On("FindBySlug", album.Slug).Return(_album.Album{}, nil)
	mocks.albumRepo.On("Save", mock.AnythingOfType("album.Album")).Return(nil)

	albumToAssert, err := manager.Create("title", "me", "description", false)
	assert.Equal(t, album.Slug, albumToAssert.Slug)
	assert.NotEqual(t, "", albumToAssert.Id)
	assert.NotEqual(t, 0, albumToAssert.CreationDate)
	assert.Nil(t, err)
}

func TestItShouldCreateAnAlbumWithDifferentSlugIfItExist(t *testing.T) {
	manager, mocks := getManagerWithMocks()

	mocks.albumRepo.On("FindBySlug", "title").Return(_album.Album{Slug: "title"}, nil)
	mocks.albumRepo.On("Save", mock.AnythingOfType("album.Album")).Return(nil)

	albumToAssert, err := manager.Create("title", "me", "description", false)
	assert.Contains(t, albumToAssert.Slug, "title-")
	assert.NotEqual(t, "", albumToAssert.Id)
	assert.NotEqual(t, 0, albumToAssert.CreationDate)
	assert.Nil(t, err)
}

func TestItShouldNotCreateAnAlbumIfFindSlugFail(t *testing.T) {
	manager, mocks := getManagerWithMocks()

	mocks.albumRepo.On("FindBySlug", "title").Return(_album.Album{}, errors.New("fail"))

	albumToAssert, err := manager.Create("title", "me", "description", false)
	assert.Equal(t, "", albumToAssert.Slug)
	assert.NotNil(t, err)
}

func TestItShouldNotCreateAnAlbumIfSaveFail(t *testing.T) {
	manager, mocks := getManagerWithMocks()

	mocks.albumRepo.On("FindBySlug", "title").Return(_album.Album{}, nil)
	mocks.albumRepo.On("Save", mock.AnythingOfType("album.Album")).Return(errors.New("fail"))

	albumToAssert, err := manager.Create("title", "me", "description", false)
	assert.Equal(t, "", albumToAssert.Slug)
	assert.NotNil(t, err)
}

func TestItShouldEditAnAlbumIfEverythingIsOK(t *testing.T) {
	manager, mocks := getManagerWithMocks()

	albumFind := album.Album{
		Title:       "title",
		Description: "description",
		Private:     false,
		Slug:        "title",
	}

	mocks.albumRepo.On("FindBySlug", "title").Return(albumFind, nil)
	mocks.albumRepo.On("Save", mock.AnythingOfType("album.Album")).Return(nil)

	albumToAssert, err := manager.Edit(albumFind.Title, "new description", albumFind.Slug, albumFind.Private)

	assert.Equal(t, "new description", albumToAssert.Description)
	assert.Nil(t, err)
}

func TestItShoulNotEditAnAlbumIfFindBySlugFail(t *testing.T) {
	manager, mocks := getManagerWithMocks()

	mocks.albumRepo.On("FindBySlug", "title").Return(album.Album{}, errors.New("fail"))

	albumToAssert, err := manager.Edit("title", "new description", "title", false)
	assert.Equal(t, "", albumToAssert.Slug)
	assert.NotNil(t, err)
}

func TestItShoulNotEditAnAlbumIfFindBySlugReturnNothing(t *testing.T) {
	manager, mocks := getManagerWithMocks()

	mocks.albumRepo.On("FindBySlug", "title").Return(album.Album{}, nil)

	albumToAssert, err := manager.Edit("title", "new description", "title", false)
	assert.Equal(t, "", albumToAssert.Slug)
	assert.NotNil(t, err)
}

func TestItShoulNotEditAnAlbumIfSaveFail(t *testing.T) {
	manager, mocks := getManagerWithMocks()

	mocks.albumRepo.On("FindBySlug", "title").Return(album.Album{Slug: "slug"}, nil)
	mocks.albumRepo.On("Save", mock.AnythingOfType("album.Album")).Return(errors.New("fail"))

	albumToAssert, err := manager.Edit("title", "new description", "title", false)
	assert.Equal(t, "", albumToAssert.Slug)
	assert.NotNil(t, err)
}

func TestItShouldUpdateMediasOnInsertIfEverythingIsOk(t *testing.T) {
	manager, mocks := getManagerWithMocks()

	mediasToAdd := make([]_album.Media, 0)
	mediasToAdd = append(mediasToAdd, _album.Media{
		Key:    "1",
		Author: "me",
		Kind:   album.KindPhoto,
	})

	mocks.albumRepo.On("FindBySlug", "slug").Return(album.Album{Slug: "slug"}, nil)
	mocks.albumRepo.On("Save", mock.AnythingOfType("album.Album")).Return(nil)

	albumToAssert, err := manager.UpdateMedias("slug", mediasToAdd, _album.UpdateMediaAdd)

	assert.Equal(t, 1, len(albumToAssert.Medias))
	assert.Equal(t, "1", albumToAssert.Medias[0].Key)
	assert.Nil(t, err)
}

func TestItShouldUpdateMediasOnRemoveIfEverythingIsOk(t *testing.T) {
	manager, mocks := getManagerWithMocks()

	mediasToRemove := make([]_album.Media, 0)
	mediasToRemove = append(mediasToRemove, _album.Media{
		Key:    "1",
		Author: "me",
		Kind:   album.KindPhoto,
	})

	albumFind := album.Album{
		Title:       "title",
		Description: "description",
		Private:     false,
		Slug:        "title",
		Medias: []_album.Media{
			{Key: "1"},
			{Key: "2"},
		},
	}

	mocks.albumRepo.On("FindBySlug", "slug").Return(albumFind, nil)
	mocks.albumRepo.On("Save", mock.AnythingOfType("album.Album")).Return(nil)

	albumToAssert, err := manager.UpdateMedias("slug", mediasToRemove, _album.UpdateMediaRemove)

	assert.Equal(t, 1, len(albumToAssert.Medias))
	assert.Equal(t, "2", albumToAssert.Medias[0].Key)
	assert.Nil(t, err)
}

func TestItShouldDeleteAnAlbumIfEverythingIsOk(t *testing.T) {
	manager, mocks := getManagerWithMocks()

	mocks.albumRepo.On("DeleteBySlug", "slug").Return(nil)
	err := manager.Delete("slug")

	assert.Nil(t, err)
}

func TestItShouldSearchWithParameters(t *testing.T) {
	manager, mocks := getManagerWithMocks()

	albums := []_album.Album{
		{Slug: "slug1"},
		{Slug: "slug2"},
	}

	mocks.albumRepo.On("Search", false, false, 10, 0, "", "desc").Return(albums, nil)

	albumsToAssert, err := manager.Search(false, false, 10, 0, "", "desc")

	assert.Equal(t, 2, len(albumsToAssert))
	assert.Nil(t, err)
}

func TestItShouldGetBySlug(t *testing.T) {
	manager, mocks := getManagerWithMocks()

	mocks.albumRepo.On("FindBySlug", "slug").Return(_album.Album{Slug: "slug"}, nil)

	albumToAssert, err := manager.GetBySlug("slug")

	assert.Equal(t, "slug", albumToAssert.Slug)
	assert.Nil(t, err)
}

func TestItShouldGetAll(t *testing.T) {
	manager, mocks := getManagerWithMocks()

	albums := []_album.Album{
		{Slug: "slug1"},
		{Slug: "slug2"},
	}

	mocks.albumRepo.On("Search", true, true, 1000000, 0, "", "desc").Return(albums, nil)

	albumsToAssert, err := manager.GetAll()

	assert.Equal(t, 2, len(albumsToAssert))
	assert.Nil(t, err)
}

func TestItShouldToggleFavorite(t *testing.T) {
	manager, mocks := getManagerWithMocks()

	album := _album.Album{
		Slug: "slug",
		Medias: []_album.Media{
			{Key: "1", Favorite: false},
			{Key: "2", Favorite: false},
		},
	}

	mocks.albumRepo.On("FindBySlug", "slug").Return(album, nil)
	mocks.albumRepo.On("Save", mock.AnythingOfType("album.Album")).Return(nil)

	albumToAssert, err := manager.ToggleFavorite("slug", "1")

	for _, m := range albumToAssert.Medias {
		if m.Key == "1" {
			assert.True(t, m.Favorite)
			break
		}
	}
	assert.Nil(t, err)
}
