package user_test

import (
	"errors"
	"testing"

	_user "github.com/hyoa/album/api/internal/user"
	_mocks "github.com/hyoa/album/api/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

/**
*
* Create
*
**/

type mocks struct {
	userRepo *_mocks.UserRepo
	mailer   *_mocks.Mailer
}

func getUseCaseWithMocks() (_user.UserManager, mocks) {
	mocksList := mocks{
		userRepo: new(_mocks.UserRepo),
		mailer:   new(_mocks.Mailer),
	}

	return _user.CreateUserManager(
		mocksList.userRepo,
		mocksList.mailer,
	), mocksList
}

func TestItShouldCreateAnUserIfEverythingIsOk(t *testing.T) {
	useCase, mocks := getUseCaseWithMocks()

	mocks.userRepo.On("FindByEmail", "email").Return(_user.User{}, nil)
	mocks.userRepo.On("Save", mock.AnythingOfType("user.User")).Return(nil)

	user, err := useCase.Create("name", "email", "password", "password")

	assert.Nil(t, err)
	assert.Equal(t, "name", user.Name)
	assert.Equal(t, "email", user.Email)
}

func TestItShouldNotCreateAnUserIfPasswordDoesNotMatch(t *testing.T) {
	useCase, _ := getUseCaseWithMocks()

	user, err := useCase.Create("name", "email", "password", "password2")

	assert.NotNil(t, err)
	assert.IsType(t, &_user.InvalidPasswordError{}, err)
	assert.Equal(t, _user.User{}, user)
}

func TestItShouldNotCreateAnUserIfTheEmailAlreadyExist(t *testing.T) {
	useCase, mocks := getUseCaseWithMocks()

	mocks.userRepo.On("FindByEmail", "email").Return(_user.User{Email: "email"}, nil)

	user, err := useCase.Create("name", "email", "password", "password")

	assert.NotNil(t, err)
	assert.IsType(t, &_user.UserAlreadyExistError{}, err)
	assert.Equal(t, _user.User{}, user)
}

/**
*
* SignIn
*
**/

func TestItShouldAuthIfEmailAndPasswordAreValid(t *testing.T) {
	useCase, mocks := getUseCaseWithMocks()

	hashPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)

	mocks.userRepo.On("FindByEmail", "email").Return(_user.User{Name: "name", Email: "email", Password: string(hashPassword), Role: _user.RoleNormal}, nil)

	user, err := useCase.SignIn("email", "password")

	assert.Nil(t, err)
	assert.Equal(t, "name", user.Name)
	assert.Equal(t, "email", user.Email)
}

func TestItShouldNotAuthUserIfUserDoesNotExist(t *testing.T) {
	useCase, mocks := getUseCaseWithMocks()

	mocks.userRepo.On("FindByEmail", "email2").Return(_user.User{}, nil)

	user, err := useCase.SignIn("email2", "password")

	assert.NotNil(t, err)
	assert.Equal(t, _user.User{}, user)
}

func TestItShouldNotAuthUserIfPasswordIsNotValid(t *testing.T) {
	useCase, mocks := getUseCaseWithMocks()

	hashPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)

	mocks.userRepo.On("FindByEmail", "email").Return(_user.User{Name: "name", Email: "email", Password: string(hashPassword)}, nil)

	user, err := useCase.SignIn("email", "password2")

	assert.NotNil(t, err)
	assert.Equal(t, _user.User{}, user)
}

/**
*
* CreateAuthJWT
*
**/
func TestItShouldReturnAnAuthJWT(t *testing.T) {
	useCase, _ := getUseCaseWithMocks()

	user := _user.User{Email: "email", Name: "name"}

	_, err := useCase.CreateAuthJWT(user)

	assert.Nil(t, err)
}

/**
*
* Auth
*
**/

func TestItShouldAuthUser(t *testing.T) {
	useCase, _ := getUseCaseWithMocks()

	jwt := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJOYW1lIjoibmFtZSIsIkVtYWlsIjoiZW1haWwiLCJSb2xlIjowLCJleHAiOjE5NTc3MjI5MTIsImlhdCI6MTY1NzcxOTMxMiwiaXNzIjoiYXBpdjIifQ.hsaebn8FSaiJMe3ZjoMQnGnVfPIe4Y99JzMPZ7dXEAU"
	user := _user.User{Email: "email", Name: "name"}

	userToValidate, err := useCase.Auth(jwt)

	assert.Equal(t, user, userToValidate)
	assert.Nil(t, err)
}

func TestItShouldNotAuthUserIfDecodeFail(t *testing.T) {
	useCase, _ := getUseCaseWithMocks()

	userToValidate, err := useCase.Auth("jwt")

	assert.Equal(t, _user.User{}, userToValidate)
	assert.NotNil(t, err)
}

/**
*
* ChangeRole
*
**/

func TestItShouldChangeTheRoleOfAnUser(t *testing.T) {
	useCase, mocks := getUseCaseWithMocks()

	user := _user.User{Email: "email", Name: "name", Role: _user.RoleUnidentified}
	userExpected := _user.User{Email: "email", Name: "name", Role: _user.RoleAdmin}

	mocks.userRepo.On("Update", userExpected).Return(userExpected, nil)
	mocks.userRepo.On("FindByEmail", "email").Return(user, nil)

	_, err := useCase.ChangeRole("email", _user.RoleAdmin)

	assert.Nil(t, err)
}

func TestItShouldNotUpdateRoleIfRoleIsInvalid(t *testing.T) {
	useCase, mocks := getUseCaseWithMocks()

	user := _user.User{Email: "email", Name: "name", Role: _user.RoleUnidentified}

	mocks.userRepo.On("FindByEmail", "email").Return(user, nil)

	_, err := useCase.ChangeRole("email", 100)

	assert.NotNil(t, err)
}

/**
*
* AskResetPassword
*
**/

func TestItShouldSendAnEmailToResetPassword(t *testing.T) {
	useCase, mocks := getUseCaseWithMocks()

	user := _user.User{Email: "email", Name: "name", Role: _user.RoleUnidentified}

	mocks.userRepo.On("FindByEmail", "email").Return(user, nil)

	mocks.mailer.On(
		"SendMail",
		user.Email,
		"MailSubjectPasswordChange",
		"MailBodyPasswordChange",
		mock.Anything,
	).Return(nil)

	_, err := useCase.AskResetPassword("email", "localhost")

	assert.Nil(t, err)
}

func TestItShouldNotSendAnEmailIfUserGetFailed(t *testing.T) {
	useCase, mocks := getUseCaseWithMocks()

	mocks.userRepo.On("FindByEmail", "email").Return(_user.User{}, errors.New("failed"))

	_, err := useCase.AskResetPassword("email", "")

	assert.NotNil(t, err)
}

func TestItShouldNotSendAnEmailIfUserDoesNotExist(t *testing.T) {
	useCase, mocks := getUseCaseWithMocks()

	mocks.userRepo.On("FindByEmail", "email").Return(_user.User{}, nil)

	user, err := useCase.AskResetPassword("email", "")

	assert.Nil(t, err)
	assert.Equal(t, "", user.Name)
}

/**
*
* Invite
*
**/

func TestItShouldSendAnInvitationToTheRequestedEmail(t *testing.T) {
	useCase, mocks := getUseCaseWithMocks()

	mocks.mailer.On(
		"SendMail",
		"email@email.com",
		"MailSubjectInvite",
		"MailBodyInvite",
		mock.Anything,
	).Return(nil)

	err := useCase.Invite(_user.User{}, "email@email.com", "localhost")

	assert.Nil(t, err)
}

/**
*
* ResetPassword
*
**/

func TestItShouldResetThePassword(t *testing.T) {
	useCase, mocks := getUseCaseWithMocks()

	mocks.userRepo.On("FindByEmail", "email").Return(_user.User{Email: "email"}, nil)
	mocks.userRepo.On("Update", mock.AnythingOfType("User")).Return(_user.User{Email: "email"}, nil)

	jwt := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJOYW1lIjoibmFtZSIsIkVtYWlsIjoiZW1haWwiLCJSb2xlIjowLCJleHAiOjE5NTc3MjI5MTIsImlhdCI6MTY1NzcxOTMxMiwiaXNzIjoiYXBpdjIifQ.hsaebn8FSaiJMe3ZjoMQnGnVfPIe4Y99JzMPZ7dXEAU"

	_, err := useCase.ResetPassword("password", "password", jwt)
	assert.Nil(t, err)
}
