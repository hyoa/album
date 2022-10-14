package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"

	"github.com/hyoa/album/api/graph/generated"
	"github.com/hyoa/album/api/graph/model"
	"github.com/hyoa/album/api/internal/album"
	"github.com/hyoa/album/api/internal/media"
	"github.com/hyoa/album/api/internal/user"
)

// CreateUser is the resolver for the createUser field.
func (r *mutationResolver) CreateUser(ctx context.Context, input model.CreateInput) (*model.User, error) {
	user, errSign := r.UserManager.Create(input.Name, input.Email, input.Password, input.PasswordCheck)

	if errSign != nil {
		return &model.User{}, errSign
	}

	return &model.User{
		Name:       user.Name,
		Email:      user.Email,
		CreateDate: int(user.CreateDate),
		Role:       model.RoleReverse[int(user.Role)],
	}, nil
}

// UpdateUser is the resolver for the updateUser field.
func (r *mutationResolver) UpdateUser(ctx context.Context, input model.UpdateInput) (*model.User, error) {
	if input.Role != nil {
		user, errUpdateRole := r.UserManager.ChangeRole(input.Email, user.Role(input.Role.Int()))

		if errUpdateRole != nil {
			return &model.User{}, errUpdateRole
		}

		return &model.User{
			Name:       user.Name,
			Email:      user.Email,
			CreateDate: int(user.CreateDate),
			Role:       model.RoleReverse[int(user.Role)],
		}, nil
	}

	return &model.User{}, errors.New("no modification provided")
}

// ResetPassword is the resolver for the resetPassword field.
func (r *mutationResolver) ResetPassword(ctx context.Context, input *model.ResetPasswordInput) (*model.User, error) {
	user, errReset := r.UserManager.ResetPassword(input.Password, input.PasswordCheck, input.TokenValidation)

	if errReset != nil {
		return &model.User{}, errReset
	}

	return &model.User{
		Name:       user.Name,
		Email:      user.Email,
		CreateDate: int(user.CreateDate),
		Role:       model.RoleReverse[int(user.Role)],
	}, nil
}

// AskResetPassword is the resolver for the askResetPassword field.
func (r *mutationResolver) AskResetPassword(ctx context.Context, input model.AskResetPasswordInput) (*model.User, error) {
	user, err := r.UserManager.AskResetPassword(input.Email, "localhost:3118")

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

// Invite is the resolver for the invite field.
func (r *mutationResolver) Invite(ctx context.Context, input *model.InviteInput) (*model.Invitation, error) {
	err := r.UserManager.Invite(user.User{}, input.Email, "localhost:3118")

	return &model.Invitation{Email: input.Email}, err
}

// CreateAlbum is the resolver for the createAlbum field.
func (r *mutationResolver) CreateAlbum(ctx context.Context, input model.CreateAlbumInput) (*model.Album, error) {
	if input.Description == nil {
		*input.Description = ""
	}

	album, err := r.AlbumManager.Create(input.Title, input.Author, *input.Description, input.Private)

	if err != nil {
		return &model.Album{}, err
	}

	return model.HydrateAlbum(album), nil
}

// UpdateAlbum is the resolver for the updateAlbum field.
func (r *mutationResolver) UpdateAlbum(ctx context.Context, input model.UpdateAlbumInput) (*model.Album, error) {
	album, err := r.AlbumManager.Edit(input.Title, input.Description, input.Slug, input.Private)

	if err != nil {
		return &model.Album{}, err
	}

	return model.HydrateAlbum(album), nil
}

// DeleteAlbum is the resolver for the deleteAlbum field.
func (r *mutationResolver) DeleteAlbum(ctx context.Context, input model.DeleteAlbumInput) (*model.ActionResult, error) {
	err := r.AlbumManager.Delete(input.Slug)

	res := true
	if err != nil {
		res = false
	}

	return &model.ActionResult{Success: &res}, err
}

// UpdateAlbumMedias is the resolver for the updateAlbumMedias field.
func (r *mutationResolver) UpdateAlbumMedias(ctx context.Context, input model.UpdateAlbumMediasInput) (*model.Album, error) {
	action := model.ActionUpdateAlbumMediasToString[*input.Action]

	var medias []album.Media
	for _, m := range input.Medias {

		medias = append(medias, album.Media{
			Key:    m.Key,
			Author: m.Author,
			Kind:   album.MediaKind(model.MediaTypeToString[m.Kind]),
		})
	}

	album, err := r.AlbumManager.UpdateMedias(input.Slug, medias, album.UpdateMediaKind(action))

	if err != nil {
		return &model.Album{}, err
	}

	return model.HydrateAlbum(album), nil
}

// UpdateAlbumFavorite is the resolver for the updateAlbumFavorite field.
func (r *mutationResolver) UpdateAlbumFavorite(ctx context.Context, input model.UpdateAlbumFavoriteInput) (*model.Album, error) {
	album, err := r.AlbumManager.ToggleFavorite(input.Slug, input.MediaKey)

	if err != nil {
		return &model.Album{}, err
	}

	return model.HydrateAlbum(album), nil
}

// Ingest is the resolver for the ingest field.
func (r *mutationResolver) Ingest(ctx context.Context, input model.PutIngestInput) ([]*model.PutIngestMediaOutput, error) {
	var mediasIngest []*model.PutIngestMediaOutput
	for _, m := range input.Medias {
		media, err := r.MediaManager.Ingest(m.Key, m.Author, m.Folder, media.MediaKind(model.MediaTypeToString[m.Kind]))

		var status model.PutIngestMediaStatus

		if err != nil {
			status = model.PutIngestMediaStatusFailed
		} else if media.Key == "" {
			status = model.PutIngestMediaStatusAlreadyExist
		} else {
			status = model.PutIngestMediaStatusSuccess
		}

		mediasIngest = append(mediasIngest, &model.PutIngestMediaOutput{
			Key:    m.Key,
			Status: status,
		})

	}

	return mediasIngest, nil
}

// ChangeMediasFolder is the resolver for the changeMediasFolder field.
func (r *mutationResolver) ChangeMediasFolder(ctx context.Context, input *model.ChangeMediasFolderInput) (*model.Folder, error) {
	medias, err := r.MediaManager.ChangeMediasFolder(input.Keys, input.FolderName)

	if err != nil {
		return &model.Folder{}, err
	}

	var mediasToReturn []*model.Media

	for _, m := range medias {
		mediasToReturn = append(mediasToReturn, model.HydrateMedia(m))
	}

	return &model.Folder{
		Name:   input.FolderName,
		Medias: mediasToReturn,
	}, nil
}

// ChangeFolderName is the resolver for the changeFolderName field.
func (r *mutationResolver) ChangeFolderName(ctx context.Context, input *model.ChangeFolderNameInput) (*model.Folder, error) {
	medias, err := r.MediaManager.ChangeFolderName(input.OldName, input.NewName)

	if err != nil {
		return &model.Folder{}, err
	}

	var mediasToReturn []*model.Media

	for _, m := range medias {
		mediasToReturn = append(mediasToReturn, model.HydrateMedia(m))
	}

	return &model.Folder{
		Name:   input.NewName,
		Medias: mediasToReturn,
	}, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
