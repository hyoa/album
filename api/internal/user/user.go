package user

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
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
	SendMail(email, subject, body string) error
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
		return User{}, &InvalidPasswordError{message: "Password and passwordCheck does not match"}
	}

	userMatch, errorFindUser := um.userRepo.FindByEmail(email)
	if errorFindUser != nil && errorFindUser.Error() != "No user found" {
		return User{}, fmt.Errorf("An error occured while fetching user %w", errorFindUser)
	}

	if userMatch != (User{}) {
		return User{}, &UserAlreadyExistError{}
	}

	hashPassword, errHash := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if errHash != nil {
		return User{}, fmt.Errorf("Error while creating hash for password %w", errHash)
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
		return User{}, fmt.Errorf("Unable to create user %w", errorSave)
	}

	return user, nil
}

func (um *UserManager) SignIn(email, password string) (User, error) {
	user, errorFindUser := um.userRepo.FindByEmail(email)
	if errorFindUser != nil {
		return User{}, fmt.Errorf("An error occured while fetching user %w", errorFindUser)
	}

	if user == (User{}) {
		return User{}, &UserNotFoundError{}
	}

	if user.Role == RoleUnidentified {
		return User{}, &UnauthorizedRoleError{}
	}

	errCheck := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if errCheck != nil {
		return User{}, &InvalidPasswordError{message: "Password is incorrect"}
	}

	return user, nil
}

func (um *UserManager) CreateAuthJWT(u User) (string, error) {
	token, errorCreateToken := um.authTokenizer.create(authTokenInput{u: u})
	if errorCreateToken != nil {
		return "", fmt.Errorf("Unable to create token %w", errorCreateToken)
	}

	tokenAsString, errorStringify := um.authTokenizer.stringify(token)
	if errorStringify != nil {
		return "", fmt.Errorf("Unable to stringify token %w", errorStringify)
	}

	return tokenAsString, nil
}

func (um *UserManager) Auth(jwt string) (User, error) {
	token, errorDecode := um.authTokenizer.Decode(jwt)
	if errorDecode != nil {
		return User{}, fmt.Errorf("Unable to decode token %w", errorDecode)
	}

	return User{Name: token.Name, Email: token.Email}, nil
}

func (um *UserManager) ChangeRole(email string, r Role) (User, error) {
	if r != RoleUnidentified && r != RoleAdmin && r != RoleNormal {
		return User{}, &InvalidRoleError{}
	}

	user, errFind := um.userRepo.FindByEmail(email)

	if errFind != nil {
		return User{}, fmt.Errorf("unable to find user %w", errFind)
	}

	user.Role = r
	uUpdated, errUpdate := um.userRepo.Update(user)

	if errUpdate != nil {
		return User{}, fmt.Errorf("unable to update user %w", errUpdate)
	}

	return uUpdated, nil
}

func (um *UserManager) AskResetPassword(email, appUri string) (User, error) {
	resetToken, errCreate := um.resetTokenizer.create(resetTokenInput{email: email})

	if errCreate != nil {
		return User{}, fmt.Errorf("Unable to create reset token: %w", errCreate)
	}

	jwtReset, errStringify := um.resetTokenizer.stringify(resetToken)

	if errStringify != nil {
		return User{}, fmt.Errorf("Unable to stringify reset token: %w", errStringify)
	}

	u, errFind := um.userRepo.FindByEmail(email)

	if errFind != nil {
		return User{}, fmt.Errorf("Unable to find user: %w", errFind)
	}

	if u == (User{}) {
		return User{}, &UserNotFoundError{}
	}

	errSend := um.mailer.SendMail(
		u.Email,
		"Changement du mot de passe",
		fmt.Sprintf("Cliquer sur le lien suivant pour changer de mot de passe: %s?token=%s", appUri, jwtReset),
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
		"Invitation Pauline&Jules",
		fmt.Sprintf("Bonjour, nous vous invitons à voir nos albums photo à l'adresse suivante : %s", appUri),
	)

	return nil
}

func (um *UserManager) ResetPassword(newPassword, newPasswordCheck, token string) (User, error) {
	resetToken, errDecode := um.resetTokenizer.Decode(token)

	if errDecode != nil {
		return User{}, fmt.Errorf("Unable to decode token: %w", errDecode)
	}

	if newPassword != newPasswordCheck {
		return User{}, &InvalidPasswordError{message: "Password and passwordCheck does not match"}
	}

	user, errorGetUser := um.userRepo.FindByEmail(resetToken.Email)

	if errorGetUser != nil {
		return User{}, fmt.Errorf("Unable to fetch user: %w", errorGetUser)
	}

	hashPassword, errHash := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if errHash != nil {
		return User{}, fmt.Errorf("Error while creating hash for password %w", errHash)
	}

	user.Password = string(hashPassword)

	_, errorUpdate := um.userRepo.Update(user)

	return user, errorUpdate
}
