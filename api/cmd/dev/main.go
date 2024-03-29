package main

import (
	"encoding/base64"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/hyoa/album/api/controller"
	"github.com/hyoa/album/api/internal/album"
	"github.com/hyoa/album/api/internal/awsinteractor"
	"github.com/hyoa/album/api/internal/cdn"
	"github.com/hyoa/album/api/internal/mailer"
	"github.com/hyoa/album/api/internal/media"
	"github.com/hyoa/album/api/internal/translator"
	"github.com/hyoa/album/api/internal/user"
	"github.com/patrickmn/go-cache"

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
	err := godotenv.Load(".env.local")

	if err != nil {
		panic(err)
	}

	translatorManager := *translator.CreateTranslator("./i18n/active.fr.toml")
	mailer := mailer.SendgridMailer{ApiKey: os.Getenv("MAILER_KEY"), Translater: translatorManager}
	s3, _ := awsinteractor.NewS3Interactor(os.Getenv("S3_ENDPOINT"), os.Getenv("AKID"), os.Getenv("ASK"))
	s3Storage := media.NewS3Storage(s3)
	converter := media.NewCloudConvert()

	cache := cache.New(5*time.Minute, 10*time.Minute)

	secret, err := base64.StdEncoding.DecodeString(os.Getenv("CDN_SECRET"))

	if err != nil {
		panic(err)
	}

	cdn, _ := cdn.NewCDNAWSInteractor(s3, *cache, string(secret))

	userManager := user.CreateUserManager(user.NewUserRepositoryDynamoDB(), &mailer)
	albumManager := album.CreateAlbumManager(album.NewAlbumRepositoryDynamoDB())
	mediaManager := media.CreateMediaManager(media.NewMediaRepositoryDynamoDB(), s3Storage, converter)

	restController := controller.CreateRestController(mediaManager)

	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowHeaders = append(config.AllowHeaders, "Authorization")

	r.Use(cors.New(config))
	r.POST("/v3/graphql", controller.GraphqlHandler(userManager, albumManager, mediaManager, &translatorManager, cdn))
	r.GET("/", playgroundHandler())
	r.POST("/v3/video/acknowledge/cloudconvert", restController.AcknowledgeCloudconvertCall)

	r.Run(":3118")
}
