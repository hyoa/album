package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/hyoa/album/api/graph/generated"
	"github.com/hyoa/album/api/graph/model"

	_cdn "github.com/hyoa/album/api/internal/cdn"
)

// Urls is the resolver for the urls field.
func (r *mediaResolver) Urls(ctx context.Context, media *model.Media) (*model.Urls, error) {
	urls := &model.Urls{
		Small:  r.CDN.SignGetUri(media.Key, _cdn.SizeSmall, _cdn.MediaKind(model.MediaTypeToString[media.Kind])),
		Medium: r.CDN.SignGetUri(media.Key, _cdn.SizeMedium, _cdn.MediaKind(model.MediaTypeToString[media.Kind])),
		Large:  r.CDN.SignGetUri(media.Key, _cdn.SizeLarge, _cdn.MediaKind(model.MediaTypeToString[media.Kind])),
	}

	return urls, nil
}

// Media returns generated.MediaResolver implementation.
func (r *Resolver) Media() generated.MediaResolver { return &mediaResolver{r} }

type mediaResolver struct{ *Resolver }
