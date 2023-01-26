package user

import (
	"errors"
	"fmt"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserError struct {
	i18nId string
	err    error
}

func (e *UserError) Error() string {
	return fmt.Sprintf("%v", e.err)
}

func (e *UserError) Unwrap() error {
	return e.err
}

func (e *UserError) I18NID() string {
	return e.i18nId
}

type userError string

func (e userError) Error() string {
	return string(e)
}

const ErrInvalidPassword = userError("password and passwordCheck does not match")
const ErrUserAlreadyExist = userError("user already exists")
const ErrUnauthorized = userError("role does not allow this action")
const ErrAuthentication = userError("email and/or password is invalid")

const ErrUserNotFound = userError("user does not exist")
const ErrInvalidRole = userError("invalid role")
const ErrPasswordGeneration = userError("password generation failed")
const ErrGetUser = userError("cannot get user")
const ErrSaveUSer = userError("cannot save user")
const ErrCreateAuthJWT = userError("cannot create auth jwt")
const ErrAuth = userError("auth failed")
const ErrUpdateUser = userError("cannot update user")
const ErrResetPassword = userError("cannot reset password")
const ErrAskResetPassword = userError("cannot ask reset password")

type Role int

const (
	RoleUnidentified Role = iota
	RoleNormal            = 1
	RoleAdmin             = 9
)

type User struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	CreateDate int64  `json:"createDate"`
	Password   string `json:"password"`
	Role       Role   `json:"role"`
}

type UserNoConfidentialData struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	CreateDate int64  `json:"createDate"`
	Role       Role   `json:"role"`
}

type UserManager struct {
	userRepo       UserRepo
	mailer         Mailer
	authTokenizer  authTokenizer
	resetTokenizer resetTokenizer
}

type Mailer interface {
	SendMail(email, subjectKey, bodyKey string, bodyData map[string]interface{}) error
}

func CreateUserManager(ur UserRepo, m Mailer) UserManager {

	return UserManager{
		userRepo:       ur,
		mailer:         m,
		authTokenizer:  CreateAuthTokenizer().(authTokenizer),
		resetTokenizer: CreateResetTokenizer().(resetTokenizer),
	}
}

func (um *UserManager) Create(name, email, password, passwordCheck string) (User, error) {
	if password != passwordCheck {
		return User{}, createErrorWithI18N(ErrInvalidPassword)
	}

	userMatch, errorFindUser := um.userRepo.FindByEmail(email)
	if errorFindUser != nil && errorFindUser.Error() != "No user found" {
		return User{}, ErrGetUser
	}

	if userMatch != (User{}) {
		return User{}, createErrorWithI18N(ErrUserAlreadyExist)
	}

	hashPassword, errHash := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if errHash != nil {
		return User{}, ErrPasswordGeneration
	}

	user := User{
		Name:       name,
		Email:      email,
		CreateDate: time.Now().Unix(),
		Password:   string(hashPassword),
		Role:       RoleUnidentified,
	}

	errorSave := um.userRepo.Save(user)
	if errorSave != nil {
		return User{}, ErrSaveUSer
	}

	return user, nil
}

func (um *UserManager) SignIn(email, password string) (User, error) {
	user, errorFindUser := um.userRepo.FindByEmail(email)
	if errorFindUser != nil {
		return User{}, ErrGetUser
	}

	if user == (User{}) {
		return User{}, createErrorWithI18N(ErrAuthentication)
	}

	if user.Role == RoleUnidentified {
		return User{}, createErrorWithI18N(ErrUnauthorized)
	}

	errCheck := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if errCheck != nil {
		return User{}, createErrorWithI18N(ErrInvalidPassword)
	}

	return user, nil
}

func (um *UserManager) CreateAuthJWT(u User) (string, error) {
	token, errorCreateToken := um.authTokenizer.create(authTokenInput{u: u})
	if errorCreateToken != nil {
		return "", ErrCreateAuthJWT
	}

	tokenAsString, errorStringify := um.authTokenizer.stringify(token)
	if errorStringify != nil {
		return "", ErrCreateAuthJWT
	}

	return tokenAsString, nil
}

func (um *UserManager) Auth(jwt string) (User, error) {
	token, errorDecode := um.authTokenizer.Decode(jwt)
	if errorDecode != nil {
		return User{}, ErrAuth
	}

	return User{Name: token.Name, Email: token.Email}, nil
}

func (um *UserManager) ChangeRole(email string, r Role) (User, error) {
	if r != RoleUnidentified && r != RoleAdmin && r != RoleNormal {
		return User{}, ErrInvalidRole
	}

	user, errFind := um.userRepo.FindByEmail(email)

	if errFind != nil {
		return User{}, ErrGetUser
	}

	user.Role = r
	uUpdated, errUpdate := um.userRepo.Update(user)

	if errUpdate != nil {
		return User{}, ErrUpdateUser
	}

	return uUpdated, nil
}

func (um *UserManager) AskResetPassword(email, appUri string) (User, error) {
	resetToken, errCreate := um.resetTokenizer.create(resetTokenInput{email: email})

	if errCreate != nil {
		return User{}, ErrAskResetPassword
	}

	jwtReset, errStringify := um.resetTokenizer.stringify(resetToken)

	if errStringify != nil {
		return User{}, ErrAskResetPassword
	}

	u, errFind := um.userRepo.FindByEmail(email)

	if errFind != nil {
		return User{}, ErrGetUser
	}

	//We don't send the email, but we don't return any error. No need to tell anyone that the user does not exist
	if u == (User{}) {
		return User{}, nil
	}

	errSend := um.mailer.SendMail(
		u.Email,
		"MailSubjectPasswordChange",
		"MailBodyPasswordChange",
		map[string]interface{}{"uri": appUri, "token": jwtReset},
	)

	return u, errSend
}

func (um *UserManager) GetUsers() ([]User, error) {
	return um.userRepo.FindAll()
}

func (um *UserManager) GetUser(email string) (User, error) {
	return um.userRepo.FindByEmail(email)
}

func (um *UserManager) Invite(u User, toInviteEmail, appUri string) error {
	um.mailer.SendMail(
		toInviteEmail,
		"MailSubjectInvite",
		"MailBodyInvite",
		map[string]interface{}{"Uri": appUri, "AppName": os.Getenv("APP_NAME")},
	)

	return nil
}

func (um *UserManager) ResetPassword(newPassword, newPasswordCheck, token string) (User, error) {
	resetToken, errDecode := um.resetTokenizer.Decode(token)

	if errDecode != nil {
		return User{}, ErrResetPassword
	}

	if newPassword != newPasswordCheck {
		return User{}, createErrorWithI18N(ErrInvalidPassword)
	}

	user, errorGetUser := um.userRepo.FindByEmail(resetToken.Email)

	if errorGetUser != nil {
		return User{}, ErrGetUser
	}

	hashPassword, errHash := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if errHash != nil {
		return User{}, ErrResetPassword
	}

	user.Password = string(hashPassword)

	_, errorUpdate := um.userRepo.Update(user)

	if errorUpdate != nil {
		return User{}, ErrUpdateUser
	}
	return user, nil
}

func createErrorWithI18N(err error) error {
	var key string
	if errors.Is(err, ErrInvalidPassword) {
		key = "CreateUserPasswordCheckError"
	} else if errors.Is(err, ErrUserAlreadyExist) {
		key = "CreateUserAlreadyExist"
	} else if errors.Is(err, ErrUnauthorized) {
		key = "CreateUserRoleUnidentified"
	} else if errors.Is(err, ErrAuthentication) {
		key = "CreateUserAuthNotFoundError"
	}

	return &UserError{i18nId: key, err: err}
}
