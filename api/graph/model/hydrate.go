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
			Urls: &Urls{
				Small:  cdn.SignGetUri(album.Medias[j].Key, _cdn.SizeSmall, _cdn.MediaKind(string(album.Medias[j].Kind))),
				Medium: cdn.SignGetUri(album.Medias[j].Key, _cdn.SizeMedium, _cdn.MediaKind(string(album.Medias[j].Kind))),
				Large:  cdn.SignGetUri(album.Medias[j].Key, _cdn.SizeLarge, _cdn.MediaKind(string(album.Medias[j].Kind))),
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

func HydrateMedia(media media.Media, cdn _cdn.CDNInteractor) *Media {
	return &Media{
		Key:    media.Key,
		Author: media.Author,
		Kind:   MediaTypeReverse[string(media.Kind)],
		Urls: &Urls{
			Small:  cdn.SignGetUri(media.Key, _cdn.SizeSmall, _cdn.MediaKind(string(media.Kind))),
			Medium: cdn.SignGetUri(media.Key, _cdn.SizeMedium, _cdn.MediaKind(string(media.Kind))),
			Large:  cdn.SignGetUri(media.Key, _cdn.SizeLarge, _cdn.MediaKind(string(media.Kind))),
		},
	}
}
