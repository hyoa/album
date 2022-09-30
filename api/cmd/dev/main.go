package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/hyoa/album/api/graph"
	"github.com/hyoa/album/api/graph/generated"
	"github.com/hyoa/album/api/internal/album"
	"github.com/hyoa/album/api/internal/mailer"
	"github.com/hyoa/album/api/internal/media"
	"github.com/hyoa/album/api/internal/s3interactor"
	"github.com/hyoa/album/api/internal/user"

	"github.com/joho/godotenv"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

// Defining the Graphql handler
func graphqlHandler() gin.HandlerFunc {
	r := &graph.Resolver{}
	mailer := mailer.SendgridMailer{ApiKey: os.Getenv("MAILER_KEY")}
	s3, _ := s3interactor.NewInteractor(os.Getenv("S3_ENDPOINT"), os.Getenv("AKID"), os.Getenv("ASK"))

	r.UserManager = user.CreateUserManager(user.NewUserRepositoryDynamoDB(), &mailer)
	r.AlbumManager = album.CreateAlbumManager(album.NewAlbumRepositoryDynamoDB())
	r.MediaManager = media.CreateMediaManager(media.NewMediaRepositoryDynamoDB(), media.NewS3Storage(s3))

	h := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: r}))

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
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
	godotenv.Load()
	// s3, _ := s3interactor.NewInteractor(os.Getenv("FILE_REPO_ADDR"), os.Getenv("FILE_REPO_KEY"), os.Getenv("FILE_REPO_SECRET"))

	// userRepo := user.NewUserRepositoryFile(s3)

	// mailer := mailer.SendgridMailer{ApiKey: os.Getenv("MAILER_KEY")}
	// uc := controller.NewUserController(userRepo, &mailer)

	// albumRepo := album.NewAlbumRepositoryDynamoDB()
	// ac := controller.NewAlbumController(albumRepo)

	// r := gin.Default()

	// user := r.Group("/user")
	// {
	// 	user.POST("/signin", uc.SignIn)
	// 	user.POST("/signup", uc.SignUp)
	// 	user.POST("/reset-password", uc.AskResetPassword)
	// 	user.POST("/activate", middleware.Auth(userRepo, middleware.AdminAuth), uc.Activate)
	// 	user.POST("/invite", middleware.Auth(userRepo, middleware.AdminAuth), uc.Invite)
	// }

	// users := r.Group("/users")
	// {
	// 	users.GET("", middleware.Auth(userRepo, middleware.AdminAuth), uc.GetAll)
	// }

	// album := r.Group("/albums")
	// {
	// 	album.GET("/", ac.GetAlbums)
	// }

	r := gin.Default()
	r.POST("/query", graphqlHandler())
	r.GET("/", playgroundHandler())

	r.Run(":3118")
}
