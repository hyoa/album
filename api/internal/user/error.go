package user

type InvalidPasswordError struct {
	message string
}

func (e *InvalidPasswordError) Error() string {
	return e.message
}

type UserAlreadyExistError struct {
	message string
}

func (e *UserAlreadyExistError) Error() string {
	return "User already exist"
}

type UserNotFoundError struct {
	message string
}

func (e *UserNotFoundError) Error() string {
	return "User does not exist"
}

type InvalidRoleError struct {
	message string
}

func (e *InvalidRoleError) Error() string {
	return "Role given is not valid"
}

type UnauthorizedRoleError struct {
	message string
}

func (e *UnauthorizedRoleError) Error() string {
	return "Role does not allow this action"
}
