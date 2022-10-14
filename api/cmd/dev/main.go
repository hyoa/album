package main

import (
	"context"
	"errors"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/hyoa/album/api/graph"
	"github.com/hyoa/album/api/graph/generated"
	"github.com/hyoa/album/api/graph/model"
	"github.com/hyoa/album/api/internal/album"
	"github.com/hyoa/album/api/internal/mailer"
	"github.com/hyoa/album/api/internal/media"
	"github.com/hyoa/album/api/internal/media/impl"
	"github.com/hyoa/album/api/internal/s3interactor"
	"github.com/hyoa/album/api/internal/user"

	"github.com/joho/godotenv"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

// Defining the Graphql handler
func graphqlHandler() gin.HandlerFunc {
	r := &graph.Resolver{}
	mailer := mailer.SendgridMailer{ApiKey: os.Getenv("MAILER_KEY")}
	s3, _ := s3interactor.NewInteractor(os.Getenv("S3_ENDPOINT"), os.Getenv("AKID"), os.Getenv("ASK"))
	s3Storage := media.NewS3Storage(s3)
	converter := impl.NewCloudConvert()

	r.UserManager = user.CreateUserManager(user.NewUserRepositoryDynamoDB(), &mailer)
	r.AlbumManager = album.CreateAlbumManager(album.NewAlbumRepositoryDynamoDB())
	r.MediaManager = media.CreateMediaManager(media.NewMediaRepositoryDynamoDB(), s3Storage, converter)

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

// Defining the Playground handler
func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		panic(err)
	}

	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowHeaders = append(config.AllowHeaders, "Authorization")

	r.Use(cors.New(config))
	r.POST("/v3/graphql", graphqlHandler())
	r.GET("/", playgroundHandler())

	r.Run(":3118")
}
