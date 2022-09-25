package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/hyoa/album/api/controller"
	"github.com/hyoa/album/api/internal/album"
	"github.com/hyoa/album/api/internal/mailer"
	"github.com/hyoa/album/api/internal/s3interactor"
	"github.com/hyoa/album/api/internal/user"
	"github.com/hyoa/album/api/middleware"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	s3, _ := s3interactor.NewInteractor(os.Getenv("FILE_REPO_ADDR"), os.Getenv("FILE_REPO_KEY"), os.Getenv("FILE_REPO_SECRET"))

	userRepo := user.NewUserRepositoryFile(s3)

	mailer := mailer.SendgridMailer{ApiKey: os.Getenv("MAILER_KEY")}
	uc := controller.NewUserController(userRepo, &mailer)

	albumRepo := album.NewAlbumRepositoryDynamoDB()
	ac := controller.NewAlbumController(albumRepo)

	r := gin.Default()

	user := r.Group("/user")
	{
		user.POST("/signin", uc.SignIn)
		user.POST("/signup", uc.SignUp)
		user.POST("/reset-password", uc.AskResetPassword)
		user.POST("/activate", middleware.Auth(userRepo, middleware.AdminAuth), uc.Activate)
		user.POST("/invite", middleware.Auth(userRepo, middleware.AdminAuth), uc.Invite)
	}

	users := r.Group("/users")
	{
		users.GET("", middleware.Auth(userRepo, middleware.AdminAuth), uc.GetAll)
	}

	album := r.Group("/albums")
	{
		album.GET("/", ac.GetAlbums)
	}

	r.Run(":3118")
}
