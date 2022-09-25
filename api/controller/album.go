package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hyoa/album/api/internal/album"
)

type albumController struct {
	albumManager album.AlbumManager
}

func NewAlbumController(albumRepo album.AlbumRepository) albumController {
	return albumController{
		albumManager: album.CreateAlbumManager(albumRepo),
	}
}

type GetAlbumRequest struct {
	Term                string `form:"search"`
	IncludePrivateAlbum bool   `form:"private"`
	IncludeNoMedias     bool   `form:"noMedias"`
	Limit               int    `form:"limit"`
	Offset              int    `form:"offset"`
}

type GetAlbumsAutocompleteRequest struct {
	Term string `form:"search"`
}

func (ac *albumController) GetAlbums(ctx *gin.Context) {
	var payload GetAlbumRequest

	if ctx.ShouldBindQuery(&payload) != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Wrong payload"})
		return
	}

	albums, errSearch := ac.albumManager.Search(payload.IncludePrivateAlbum, payload.IncludeNoMedias, payload.Limit, payload.Offset, payload.Term, "desc")

	if errSearch != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": errSearch.Error()})
		return
	}

	ctx.JSON(http.StatusOK, albums)
}

func (ac *albumController) GetAlbumsAutocomplete(ctx *gin.Context) {
	type responseEntity struct {
		Label string `json:"label"`
		Value string `json:"value"`
	}

	var payload GetAlbumsAutocompleteRequest

	if ctx.ShouldBindQuery(&payload) != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Wrong payload"})
		return
	}

	albums, errSearch := ac.albumManager.Search(true, true, 1000, 0, payload.Term, "desc")

	if errSearch != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": errSearch.Error()})
		return
	}

	response := make([]responseEntity, 0)
	for _, a := range albums {
		response = append(response, responseEntity{
			Label: a.Title,
			Value: a.Slug,
		})
	}

	ctx.JSON(http.StatusOK, response)
}

func (ac *albumController) GetAlbum(ctx *gin.Context) {
	album, errSearch := ac.albumManager.GetBySlug(ctx.Param("slug"))

	if errSearch != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": errSearch.Error()})
		return
	}

	if album.Slug == "" {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "No album found"})
		return
	}

	ctx.JSON(http.StatusOK, album)
}

func (ac *albumController) CreateAlbum(ctx *gin.Context) {

}

func (ac *albumController) EditAlbum(ctx *gin.Context) {

}

func (ac *albumController) UpdateMedias(ctx *gin.Context) {

}

func (ac *albumController) DeleteAlbum(ctx *gin.Context) {

}

func (ac *albumController) AdminResume(ctx *gin.Context) {

}

func (ac *albumController) ToggleFavorite(ctx *gin.Context) {

}
