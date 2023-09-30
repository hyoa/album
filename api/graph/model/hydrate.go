package model

import (
	"github.com/hyoa/album/api/internal/album"
	_cdn "github.com/hyoa/album/api/internal/cdn"
	"github.com/hyoa/album/api/internal/media"
)

func HydrateAlbum(album album.Album, cdn _cdn.CDNInteractor) *Album {
	var medias []*MediaAlbum

	for j := range album.Medias {
		medias = append(medias, &MediaAlbum{
			Key:      album.Medias[j].Key,
			Author:   album.Medias[j].Author,
			Favorite: &album.Medias[j].Favorite,
			Kind:     MediaTypeReverse[string(album.Medias[j].Kind)],
		})
	}

	var favorites []*MediaAlbum

	for _, m := range medias {
		if m.Favorite != nil && *m.Favorite {
			favorites = append(favorites, m)
		}
	}

	if len(favorites) == 0 && len(medias) > 0 {
		favorites = append(favorites, medias[0])
	}

	return &Album{
		Title:        album.Title,
		Description:  &album.Description,
		Private:      &album.Private,
		Author:       album.Author,
		CreationDate: album.CreationDate,
		ID:           album.Id,
		Slug:         album.Slug,
		Medias:       medias,
		Favorites:    favorites,
	}
}

func HydrateMedia(media media.Media, cdn _cdn.CDNInteractor) *Media {
	return &Media{
		Key:    media.Key,
		Author: media.Author,
		Kind:   MediaTypeReverse[string(media.Kind)],
	}
}
