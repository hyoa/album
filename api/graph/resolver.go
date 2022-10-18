package graph

//go:generate go run github.com/99designs/gqlgen

import (
	"github.com/hyoa/album/api/internal/album"
	"github.com/hyoa/album/api/internal/media"
	"github.com/hyoa/album/api/internal/translator"
	"github.com/hyoa/album/api/internal/user"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	UserManager  user.UserManager
	AlbumManager album.AlbumManager
	MediaManager media.MediaManager
	Translator   translator.Translator
}
