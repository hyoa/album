package model

import "github.com/hyoa/album/api/internal/album"

func HydrateAlbum(album album.Album) *Album {
	var medias []*MediaAlbum

	for j := range album.Medias {
		medias = append(medias, &MediaAlbum{
			Key:      album.Medias[j].Key,
			Author:   album.Medias[j].Author,
			Favorite: &album.Medias[j].Favorite,
			Kind:     MediaTypeReverse[string(album.Medias[j].Kind)],
			Urls: &Urls{
				Small:  "small",
				Medium: "medium",
				Large:  "large",
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
