package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/hyoa/album/api/graph/generated"
	"github.com/hyoa/album/api/graph/model"
	"github.com/hyoa/album/api/internal/cdn"
	"github.com/hyoa/album/api/internal/media"
)

// User is the resolver for the user field.
func (r *queryResolver) User(ctx context.Context, input model.GetUserInput) (*model.User, error) {
	user, err := r.UserManager.GetUser(input.Email)

	if err != nil {
		return &model.User{}, err
	}

	return &model.User{
		Name:       user.Name,
		Email:      user.Email,
		Role:       model.RoleReverse[int(user.Role)],
		CreateDate: int(user.CreateDate),
	}, nil
}

// Users is the resolver for the users field.
func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {
	users, err := r.UserManager.GetUsers()

	usersModel := make([]*model.User, 0)
	for k := range users {

		usersModel = append(usersModel, &model.User{
			Name:       users[k].Name,
			Email:      users[k].Email,
			Role:       model.RoleReverse[int(users[k].Role)],
			CreateDate: int(users[k].CreateDate),
		})
	}

	return usersModel, err
}

// Auth is the resolver for the auth field.
func (r *queryResolver) Auth(ctx context.Context, input *model.AuthInput) (*model.Auth, error) {
	user, errSign := r.UserManager.SignIn(input.Email, input.Password)

	if errSign != nil {
		return &model.Auth{}, errSign
	}

	jwt, errJwt := r.UserManager.CreateAuthJWT(user)

	if errJwt != nil {
		return &model.Auth{}, errJwt
	}

	return &model.Auth{Token: jwt}, nil
}

// Albums is the resolver for the albums field.
func (r *queryResolver) Albums(ctx context.Context, input model.GetAlbumsInput) ([]*model.Album, error) {
	includePrivate := false
	includeNoMedias := false
	limit := 100
	offset := 0
	term := ""
	order := "asc"

	if input.IncludePrivate != nil {
		includePrivate = *input.IncludePrivate
	}

	if input.IncludeNoMedias != nil {
		includeNoMedias = *input.IncludeNoMedias
	}

	if input.Limit != nil && *input.Limit != 0 {
		limit = *input.Limit
	}

	if input.Offset != nil {
		offset = *input.Offset
	}

	if input.Term != nil {
		term = *input.Term
	}

	if input.Order != nil {
		order = *input.Order
	}

	albums, err := r.AlbumManager.Search(
		includePrivate,
		includeNoMedias,
		limit,
		offset,
		term,
		order,
	)

	if err != nil {
		return make([]*model.Album, 0), err
	}

	var albumsModel []*model.Album
	for k := range albums {
		albumsModel = append(albumsModel, model.HydrateAlbum(albums[k]))
	}

	return albumsModel, nil
}

// Album is the resolver for the album field.
func (r *queryResolver) Album(ctx context.Context, input model.GetAlbumInput) (*model.Album, error) {
	album, err := r.AlbumManager.GetBySlug(input.Slug)

	if err != nil {
		return &model.Album{}, err
	}

	return model.HydrateAlbum(album), nil
}

// Folders is the resolver for the folders field.
func (r *queryResolver) Folders(ctx context.Context, input model.GetFoldersInput) ([]*model.Folder, error) {
	var name string

	if input.Name != nil {
		name = *input.Name
	}

	foldersName, err := r.MediaManager.GetFolders(name)

	if err != nil {
		return make([]*model.Folder, 0), err
	}

	var folders []*model.Folder
	for _, f := range foldersName {
		medias, _ := r.MediaManager.GetMediasByFolder(f)

		var mediasModel []*model.Media
		for _, m := range medias {
			mediasModel = append(mediasModel, model.HydrateMedia(m))
		}

		folders = append(folders, &model.Folder{
			Name:   f,
			Medias: mediasModel,
		})
	}

	return folders, err
}

// Folder is the resolver for the folder field.
func (r *queryResolver) Folder(ctx context.Context, input model.GetFolderInput) (*model.Folder, error) {
	medias, err := r.MediaManager.GetMediasByFolder(input.Name)

	if err != nil {
		return &model.Folder{}, err
	}

	var mediasModel []*model.Media
	for _, m := range medias {

		mediasModel = append(mediasModel, &model.Media{
			Key:    m.Key,
			Author: m.Author,
			Kind:   model.MediaTypeReverse[string(m.Kind)],
			Folder: m.Folder,
			Urls: &model.Urls{
				Small:  cdn.SignGetUri(m.Key, cdn.SizeSmall, cdn.MediaKind(string(m.Kind))),
				Medium: cdn.SignGetUri(m.Key, cdn.SizeMedium, cdn.MediaKind(string(m.Kind))),
				Large:  cdn.SignGetUri(m.Key, cdn.SizeLarge, cdn.MediaKind(string(m.Kind))),
			},
		})
	}

	return &model.Folder{Name: input.Name, Medias: mediasModel}, nil
}

// Ingest is the resolver for the ingest field.
func (r *queryResolver) Ingest(ctx context.Context, input model.GetIngestInput) ([]*model.GetIngestMediaOutput, error) {
	var medias []*model.GetIngestMediaOutput

	for _, m := range input.Medias {
		uri, err := r.MediaManager.GetUploadSignedUri(m.Key, media.MediaKind(model.MediaTypeToString[m.Kind]))

		if err != nil {
			continue
		}

		medias = append(medias, &model.GetIngestMediaOutput{Key: m.Key, SignedURI: uri})
	}

	return medias, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
