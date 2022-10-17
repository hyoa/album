package controller

import (
	"context"
	"errors"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/gin-gonic/gin"
	"github.com/hyoa/album/api/graph"
	"github.com/hyoa/album/api/graph/generated"
	"github.com/hyoa/album/api/graph/model"
	"github.com/hyoa/album/api/internal/album"
	"github.com/hyoa/album/api/internal/media"
	"github.com/hyoa/album/api/internal/user"
)

func GraphqlHandler(userManager user.UserManager, albumManager album.AlbumManager, mediaManager media.MediaManager) gin.HandlerFunc {
	r := &graph.Resolver{}
	r.UserManager = userManager
	r.AlbumManager = albumManager
	r.MediaManager = mediaManager

	c := generated.Config{Resolvers: r}
	c.Directives.HasRole = func(ctx context.Context, obj interface{}, next graphql.Resolver, role model.Role) (res interface{}, err error) {
		authTokenContext := ctx.Value("AuthToken")

		userRole := user.RoleUnidentified

		token, ok := authTokenContext.(string)
		if ok {
			tokenizer := user.CreateAuthTokenizer()
			authToken, _ := tokenizer.Decode(strings.Replace(token, "Bearer ", "", -1))

			userRole = authToken.Role
		}

		if (role == model.RoleAdmin && userRole != user.RoleAdmin) || (role == model.RoleNormal && userRole != user.RoleNormal && userRole != user.RoleAdmin) {
			return nil, errors.New("access denied")
		}

		return next(ctx)
	}

	h := handler.NewDefaultServer(generated.NewExecutableSchema(c))

	return func(c *gin.Context) {
		v, ok := c.Request.Header["Authorization"]

		var ctx context.Context
		if ok && len(v) > 0 {
			ctx = context.WithValue(c.Request.Context(), "AuthToken", v[0])
		} else {
			ctx = c
		}

		h.ServeHTTP(c.Writer, c.Request.WithContext(ctx))
	}
}
