package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hyoa/album/api/internal/user"
)

type userController struct {
	userManager user.UserManager
}

func NewUserController(userRepo user.UserRepo, mailer user.Mailer) userController {
	return userController{
		userManager: user.CreateUserManager(userRepo, mailer),
	}
}

type SignUpBody struct {
	Name          string `json:"name"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	PasswordCheck string `json:"passwordCheck"`
}

type SignInBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ActivateBody struct {
	Email string `json:"email"`
}

type AskResetPasswordBody struct {
	Email string `json:"email"`
}

type InviteBody struct {
	Email string `json:"email"`
}

func (uc *userController) SignUp(ctx *gin.Context) {
	var requestBody SignUpBody
	if err := ctx.BindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, "Unable to parse body")
	}

	_, err := uc.userManager.Create(requestBody.Name, requestBody.Email, requestBody.Password, requestBody.PasswordCheck)

	if err != nil {
		switch t := err.(type) {
		default:
			ctx.JSON(http.StatusInternalServerError, map[string]string{"error": t.Error()})
		}
	} else {
		ctx.JSON(http.StatusOK, "success")
	}
}

func (uc *userController) SignIn(ctx *gin.Context) {
	var requestBody SignInBody
	if err := ctx.BindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, "Unable to parse body")
	}

	user, errSignin := uc.userManager.SignIn(requestBody.Email, requestBody.Password)

	if errSignin != nil {
		ctx.JSON(http.StatusForbidden, map[string]string{"error": "Invalid email or password"})
		return
	}

	token, errAuth := uc.userManager.CreateAuthJWT(user)

	if errAuth != nil {
		ctx.JSON(http.StatusForbidden, map[string]string{"error": errAuth.Error()})
	}

	ctx.JSON(http.StatusOK, map[string]string{"token": token})
}

func (uc *userController) Activate(ctx *gin.Context) {
	var requestBody ActivateBody
	if err := ctx.BindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, "Unable to parse body")
	}

	err := uc.userManager.ChangeRole(requestBody.Email, user.RoleNormal)

	if err != nil {
		ctx.JSON(http.StatusForbidden, map[string]string{"error": err.Error()})
	}

	ctx.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

func (uc *userController) AskResetPassword(ctx *gin.Context) {
	var requestBody AskResetPasswordBody
	if err := ctx.BindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, "Unable to parse body")
		return
	}

	err := uc.userManager.AskResetPassword(requestBody.Email, "localhost")

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "sent"})
}

func (uc *userController) Invite(ctx *gin.Context) {
	var requestBody InviteBody
	if err := ctx.BindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, "Unable to parse body")
		return
	}

	uAny, _ := ctx.Get("user")
	var u user.User

	switch v := uAny.(type) {
	case user.User:
		u = v
	default:
		ctx.JSON(http.StatusInternalServerError, "Cannot do anything")
		return
	}

	uc.userManager.Invite(u, requestBody.Email, "localhost")
}

func (uc *userController) GetAll(ctx *gin.Context) {
	users, err := uc.userManager.GetUsers()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	usersNoConfidential := make([]user.UserNoConfidentialData, 0)
	for _, u := range users {
		usersNoConfidential = append(usersNoConfidential, user.UserNoConfidentialData{
			Name:       u.Name,
			Email:      u.Email,
			CreateDate: u.CreateDate,
			Role:       u.Role,
		})
	}

	ctx.JSON(http.StatusOK, usersNoConfidential)
}
