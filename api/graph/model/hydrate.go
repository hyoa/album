package model

import (
	"github.com/hyoa/album/api/internal/album"
	"github.com/hyoa/album/api/internal/cdn"
	"github.com/hyoa/album/api/internal/media"
)

func HydrateAlbum(album album.Album) *Album {
	var medias []*MediaAlbum

	for j := range album.Medias {
		medias = append(medias, &MediaAlbum{
			Key:      album.Medias[j].Key,
			Author:   album.Medias[j].Author,
			Favorite: &album.Medias[j].Favorite,
			Kind:     MediaTypeReverse[string(album.Medias[j].Kind)],
			Urls: &Urls{
				Small:  cdn.SignGetUri(album.Medias[j].Key, cdn.SizeSmall, cdn.MediaKind(string(album.Medias[j].Kind))),
				Medium: cdn.SignGetUri(album.Medias[j].Key, cdn.SizeMedium, cdn.MediaKind(string(album.Medias[j].Kind))),
				Large:  cdn.SignGetUri(album.Medias[j].Key, cdn.SizeLarge, cdn.MediaKind(string(album.Medias[j].Kind))),
			},
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

func HydrateMedia(media media.Media) *Media {
	return &Media{
		Key:    media.Key,
		Author: media.Author,
		Kind:   MediaTypeReverse[string(media.Kind)],
		Urls: &Urls{
			Small:  cdn.SignGetUri(media.Key, cdn.SizeSmall, cdn.MediaKind(string(media.Kind))),
			Medium: cdn.SignGetUri(media.Key, cdn.SizeMedium, cdn.MediaKind(string(media.Kind))),
			Large:  cdn.SignGetUri(media.Key, cdn.SizeLarge, cdn.MediaKind(string(media.Kind))),
		},
	}
}
