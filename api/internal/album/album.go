package album

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
)

type MediaKind string
type UpdateMediaKind string

const (
	KindPhoto         MediaKind       = "photo"
	KindVideo         MediaKind       = "video"
	UpdateMediaAdd    UpdateMediaKind = "add"
	UpdateMediaRemove UpdateMediaKind = "remove"
)

type Media struct {
	Key      string    `json:"key"`
	Author   string    `json:"author"`
	Kind     MediaKind `json:"kind"`
	Favorite bool      `json:"favorite"`
}

type Album struct {
	Title        string  `json:"title"`
	Description  string  `json:"description"`
	Private      bool    `json:"private"`
	Author       string  `json:"author"`
	CreationDate int     `json:"creationDate"`
	Id           string  `json:"id"`
	Slug         string  `json:"slug"`
	Medias       []Media `json:"medias"`
}

type AlbumManager struct {
	repository AlbumRepository
}

func CreateAlbumManager(repo AlbumRepository) AlbumManager {
	return AlbumManager{
		repository: repo,
	}
}

type albumError string

func (e albumError) Error() string {
	return string(e)
}

const ErrGetAlbum = albumError("Cannot get album")
const ErrSaveAlbum = albumError("Cannot save album")

func (am *AlbumManager) Create(title, author, description string, private bool) (Album, error) {
	slugV := slug.MakeLang(title, "fr")
	albumExist, errFind := am.repository.FindBySlug(slugV)

	if errFind != nil {
		return Album{}, fmt.Errorf("%w", ErrGetAlbum)
	}

	if albumExist.Slug != "" {
		slugV = slug.MakeLang(fmt.Sprintf("%s %s", title, uuid.NewString()), "fr")
	}

	albumToSave := Album{
		Id:           uuid.NewString(),
		Title:        title,
		Author:       author,
		Description:  description,
		Private:      private,
		Medias:       make([]Media, 0),
		CreationDate: int(time.Now().Unix()),
		Slug:         slugV,
	}

	errSave := am.repository.Save(albumToSave)

	if errSave != nil {
		return Album{}, fmt.Errorf(" %w", ErrSaveAlbum)
	}

	return albumToSave, nil
}

func (am *AlbumManager) Edit(title, description, slug string, private bool) (Album, error) {
	albumFound, errFind := am.repository.FindBySlug(slug)

	if errFind != nil {
		return Album{}, fmt.Errorf("Unable to find by slug: %w", errFind)
	}

	if albumFound.Slug == "" {
		return Album{}, fmt.Errorf("cannot get album for slug %s: %w", slug, ErrSaveAlbum)
	}

	albumFound.Title = title
	albumFound.Description = description
	albumFound.Private = private

	errSave := am.repository.Save(albumFound)

	if errSave != nil {
		return Album{}, fmt.Errorf("cannot save album with slug %s: %w", slug, ErrSaveAlbum)
	}

	return albumFound, nil
}

func (am *AlbumManager) UpdateMedias(slug string, medias []Media, updateKind UpdateMediaKind) (Album, error) {
	a, errFind := am.repository.FindBySlug(slug)

	if errFind != nil {
		return Album{}, fmt.Errorf("%w", ErrGetAlbum)
	}

	if a.Slug == "" {
		return Album{}, nil
	}

	if updateKind == UpdateMediaAdd {
		a.addMedias(medias)
	} else {
		a.removeMedias(medias)
	}

	errSave := am.repository.Save(a)

	if errSave != nil {
		return Album{}, fmt.Errorf("%w", ErrSaveAlbum)
	}

	return a, nil
}

func (a *Album) addMedias(medias []Media) {
	for _, m := range medias {
		found := false
		for _, mA := range a.Medias {
			if m.Key == mA.Key {
				found = true
			}
		}

		if !found {
			a.Medias = append(a.Medias, m)
		}
	}
}

func (a *Album) removeMedias(medias []Media) {
	for _, m := range medias {
		for k, mA := range a.Medias {
			if m.Key == mA.Key {
				a.Medias[k] = a.Medias[len(a.Medias)-1]
				a.Medias = a.Medias[:len(a.Medias)-1]
				break
			}
		}

	}
}

func (am *AlbumManager) Delete(slug string) error {
	return am.repository.DeleteBySlug(slug)
}

func (am *AlbumManager) Search(includePrivate, includeNoMedias bool, limit, offset int, term, order string) ([]Album, error) {
	return am.repository.Search(includePrivate, includeNoMedias, limit, offset, term, order)
}

func (am *AlbumManager) GetBySlug(slug string) (Album, error) {
	return am.repository.FindBySlug(slug)
}

func (am *AlbumManager) GetAll() ([]Album, error) {
	return am.repository.Search(true, true, 1000000, 0, "", "desc")
}

func (am *AlbumManager) ToggleFavorite(slug, mediaKey string) (Album, error) {
	album, errFind := am.repository.FindBySlug(slug)

	if errFind != nil {
		return Album{}, fmt.Errorf("%w", ErrGetAlbum)
	}

	if album.Slug == "" {
		return Album{}, nil
	}

	for k, m := range album.Medias {
		if m.Key == mediaKey {
			album.Medias[k].Favorite = !m.Favorite
			break
		}
	}

	errSave := am.repository.Save(album)

	if errSave != nil {
		return Album{}, fmt.Errorf("%w", ErrSaveAlbum)
	}

	return album, nil
}
