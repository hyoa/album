package main

import (
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/hyoa/album/api/controller"
	"github.com/hyoa/album/api/internal/album"
	"github.com/hyoa/album/api/internal/awsinteractor"
	"github.com/hyoa/album/api/internal/mailer"
	"github.com/hyoa/album/api/internal/media"
	"github.com/hyoa/album/api/internal/user"

	"github.com/joho/godotenv"

	"github.com/99designs/gqlgen/graphql/playground"
)

// Defining the Playground handler
func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func main() {
	err := godotenv.Load(".env.test")

	if err != nil {
		panic(err)
	}

	mailer := mailer.SendgridMailer{ApiKey: os.Getenv("MAILER_KEY")}
	s3, _ := awsinteractor.NewS3Interactor(os.Getenv("S3_ENDPOINT"), os.Getenv("AKID"), os.Getenv("ASK"))
	s3Storage := media.NewS3Storage(s3)
	converter := media.NewCloudConvert()

	userManager := user.CreateUserManager(user.NewUserRepositoryDynamoDB(), &mailer)
	albumManager := album.CreateAlbumManager(album.NewAlbumRepositoryDynamoDB())
	mediaManager := media.CreateMediaManager(media.NewMediaRepositoryDynamoDB(), s3Storage, converter)

	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowHeaders = append(config.AllowHeaders, "Authorization")

	r.Use(cors.New(config))
	r.POST("/v3/graphql", controller.GraphqlHandler(userManager, albumManager, mediaManager))
	r.GET("/", playgroundHandler())

	r.Run(":3118")
}
