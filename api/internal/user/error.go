package user

/**
* ERRORS WITH TRANSLATIONS
**/

type ErrorWithTranslation interface {
	ID() string
}

type InvalidPasswordError struct {
}

func (e *InvalidPasswordError) Error() string {
	return "password and passwordCheck does not match"
}

func (e *InvalidPasswordError) ID() string {
	return "CreateUserPasswordCheckError"
}

type UserAlreadyExistError struct {
}

func (e *UserAlreadyExistError) Error() string {
	return "user already exist"
}

func (e *UserAlreadyExistError) ID() string {
	return "CreateUserAlreadyExist"
}

type UnauthorizedRoleError struct {
}

func (e *UnauthorizedRoleError) Error() string {
	return "role does not allow this action"
}

func (e *UnauthorizedRoleError) ID() string {
	return "CreateUserRoleUnidentified"
}

type AuthenticationNotFoundError struct {
}

func (e *AuthenticationNotFoundError) Error() string {
	return "email or password is invalid"
}

func (e *AuthenticationNotFoundError) ID() string {
	return "CreateUserAuthNotFoundError"
}

/**
* ERRORS WITHOUT TRANSLATION
**/

type UserNotFoundError struct {
	message string
}

func (e *UserNotFoundError) Error() string {
	return "user does not exist"
}

type UserFetchError struct {
	message string
}

func (e *UserFetchError) Error() string {
	return e.message
}

type InvalidRoleError struct {
	message string
}

func (e *InvalidRoleError) Error() string {
	return "Role given is not valid"
}

type PasswordGenerationError struct {
	message string
}

func (e *PasswordGenerationError) Error() string {
	return e.message
}

type UserCreationError struct {
	message string
}

func (e *UserCreationError) Error() string {
	return e.message
}

type CreateAuthJwtError struct {
	message string
}

func (e *CreateAuthJwtError) Error() string {
	return e.message
}

type AuthError struct {
	message string
}

func (e *AuthError) Error() string {
	return e.message
}

type UserUpdateError struct {
	message string
}

func (e *UserUpdateError) Error() string {
	return e.message
}

type AskResetPasswordError struct {
	message string
}

func (e *AskResetPasswordError) Error() string {
	return e.message
}

type ResetPasswordError struct {
	message string
}

func (e *ResetPasswordError) Error() string {
	return e.message
}
