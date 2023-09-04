package main

import (
	"context"
	"encoding/base64"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
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
	"github.com/joho/godotenv"
	"github.com/patrickmn/go-cache"
)

var ginLambda *ginadapter.GinLambda

func init() {
	err := godotenv.Load(".env")

	if err != nil {
		panic(err)
	}

	translatorManager := *translator.CreateTranslator("./i18n/active.fr.toml")
	mailer := mailer.SendgridMailer{ApiKey: os.Getenv("MAILER_KEY"), Translater: translatorManager}
	s3, _ := awsinteractor.NewS3Interactor(os.Getenv("S3_ENDPOINT"), os.Getenv("AKID"), os.Getenv("ASK"))
	s3Storage := media.NewS3Storage(s3)
	converter := media.NewCloudConvert()

	userManager := user.CreateUserManager(user.NewUserRepositoryDynamoDB(), &mailer)
	albumManager := album.CreateAlbumManager(album.NewAlbumRepositoryDynamoDB())
	mediaManager := media.CreateMediaManager(media.NewMediaRepositoryDynamoDB(), s3Storage, converter)

	cache := cache.New(5*time.Minute, 10*time.Minute)

	secret, err := base64.StdEncoding.DecodeString(os.Getenv("CDN_SECRET"))

	if err != nil {
		panic(err)
	}

	cdn, _ := cdn.NewCDNAWSInteractor(s3, *cache, string(secret))

	restController := controller.CreateRestController(mediaManager)

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowHeaders = append(config.AllowHeaders, "Authorization")

	r.Use(cors.New(config))
	r.POST("/v3/graphql", controller.GraphqlHandler(userManager, albumManager, mediaManager, &translatorManager, cdn))
	r.POST("/v3/video/acknowledge/cloudconvert", restController.AcknowledgeCloudconvertCall)

	ginLambda = ginadapter.New(r)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(Handler)
}
