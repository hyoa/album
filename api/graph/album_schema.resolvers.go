package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.36

import (
	"context"

	"github.com/hyoa/album/api/graph/generated"
	"github.com/hyoa/album/api/graph/model"
	_cdn "github.com/hyoa/album/api/internal/cdn"
)

// Urls is the resolver for the urls field.
func (r *mediaAlbumResolver) Urls(ctx context.Context, obj *model.MediaAlbum) (*model.Urls, error) {
	urls := &model.Urls{
		Small:  r.CDN.SignGetUri(obj.Key, _cdn.SizeSmall, _cdn.MediaKind(model.MediaTypeToString[obj.Kind])),
		Medium: r.CDN.SignGetUri(obj.Key, _cdn.SizeMedium, _cdn.MediaKind(model.MediaTypeToString[obj.Kind])),
		Large:  r.CDN.SignGetUri(obj.Key, _cdn.SizeLarge, _cdn.MediaKind(model.MediaTypeToString[obj.Kind])),
	}

	return urls, nil
}

// MediaAlbum returns generated.MediaAlbumResolver implementation.
func (r *Resolver) MediaAlbum() generated.MediaAlbumResolver { return &mediaAlbumResolver{r} }

type mediaAlbumResolver struct{ *Resolver }
