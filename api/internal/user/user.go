package user

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

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
		return User{}, &InvalidPasswordError{}
	}

	userMatch, errorFindUser := um.userRepo.FindByEmail(email)
	if errorFindUser != nil && errorFindUser.Error() != "No user found" {
		return User{}, &UserFetchError{message: errorFindUser.Error()}
	}

	if userMatch != (User{}) {
		return User{}, &UserAlreadyExistError{}
	}

	hashPassword, errHash := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if errHash != nil {
		return User{}, &PasswordGenerationError{message: errHash.Error()}
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
		return User{}, &UserCreationError{message: errorSave.Error()}
	}

	return user, nil
}

func (um *UserManager) SignIn(email, password string) (User, error) {
	user, errorFindUser := um.userRepo.FindByEmail(email)
	if errorFindUser != nil {
		return User{}, &UserFetchError{message: errorFindUser.Error()}
	}

	if user == (User{}) {
		return User{}, &AuthenticationNotFoundError{}
	}

	if user.Role == RoleUnidentified {
		return User{}, &UnauthorizedRoleError{}
	}

	errCheck := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if errCheck != nil {
		return User{}, &InvalidPasswordError{}
	}

	return user, nil
}

func (um *UserManager) CreateAuthJWT(u User) (string, error) {
	token, errorCreateToken := um.authTokenizer.create(authTokenInput{u: u})
	if errorCreateToken != nil {
		return "", &CreateAuthJwtError{message: fmt.Sprintf("Unable to create token %s", errorCreateToken.Error())}
	}

	tokenAsString, errorStringify := um.authTokenizer.stringify(token)
	if errorStringify != nil {
		return "", &CreateAuthJwtError{message: fmt.Sprintf("Unable to stringify token %s", errorStringify.Error())}
	}

	return tokenAsString, nil
}

func (um *UserManager) Auth(jwt string) (User, error) {
	token, errorDecode := um.authTokenizer.Decode(jwt)
	if errorDecode != nil {
		return User{}, &AuthError{message: fmt.Sprintf("Unable to decode token %s", errorDecode.Error())}
	}

	return User{Name: token.Name, Email: token.Email}, nil
}

func (um *UserManager) ChangeRole(email string, r Role) (User, error) {
	if r != RoleUnidentified && r != RoleAdmin && r != RoleNormal {
		return User{}, &InvalidRoleError{}
	}

	user, errFind := um.userRepo.FindByEmail(email)

	if errFind != nil {
		return User{}, &UserFetchError{message: errFind.Error()}
	}

	user.Role = r
	uUpdated, errUpdate := um.userRepo.Update(user)

	if errUpdate != nil {
		return User{}, &UserUpdateError{message: fmt.Sprintf("unable to update user %s", errUpdate.Error())}
	}

	return uUpdated, nil
}

func (um *UserManager) AskResetPassword(email, appUri string) (User, error) {
	resetToken, errCreate := um.resetTokenizer.create(resetTokenInput{email: email})

	if errCreate != nil {
		return User{}, &AskResetPasswordError{message: fmt.Sprintf("Unable to create reset token: %s", errCreate.Error())}
	}

	jwtReset, errStringify := um.resetTokenizer.stringify(resetToken)

	if errStringify != nil {
		return User{}, &AskResetPasswordError{message: fmt.Sprintf("Unable to stringify reset token: %s", errStringify.Error())}
	}

	u, errFind := um.userRepo.FindByEmail(email)

	if errFind != nil {
		return User{}, &UserFetchError{message: errFind.Error()}
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
		map[string]interface{}{"Uri": appUri, "AppName": "Pauline&Jules"},
	)

	return nil
}

func (um *UserManager) ResetPassword(newPassword, newPasswordCheck, token string) (User, error) {
	resetToken, errDecode := um.resetTokenizer.Decode(token)

	if errDecode != nil {
		return User{}, &ResetPasswordError{message: fmt.Sprintf("Unable to decode token: %s", errDecode.Error())}
	}

	if newPassword != newPasswordCheck {
		return User{}, &InvalidPasswordError{}
	}

	user, errorGetUser := um.userRepo.FindByEmail(resetToken.Email)

	if errorGetUser != nil {
		return User{}, &UserFetchError{message: errorGetUser.Error()}
	}

	hashPassword, errHash := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if errHash != nil {
		return User{}, &ResetPasswordError{message: fmt.Sprintf("Error while creating hash for password %s", errHash.Error())}
	}

	user.Password = string(hashPassword)

	_, errorUpdate := um.userRepo.Update(user)

	if errorUpdate != nil {
		return User{}, &UserUpdateError{message: errorUpdate.Error()}
	}
	return user, nil
}
