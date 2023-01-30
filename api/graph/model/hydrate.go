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

	return &Album{
		Title:        album.Title,
		Description:  &album.Description,
		Private:      &album.Private,
		Author:       album.Author,
		CreationDate: album.CreationDate,
		ID:           album.Id,
		Slug:         album.Slug,
		Medias:       medias,
	}
}

func HydrateMedia(media media.Media, cdn _cdn.CDNInteractor) *Media {
	return &Media{
		Key:    media.Key,
		Author: media.Author,
		Kind:   MediaTypeReverse[string(media.Kind)],
	}
}
