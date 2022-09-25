package user

import (
	"encoding/json"
	"fmt"

	"github.com/hyoa/album/api/internal/s3interactor"
)

type UserRepositoryFile struct {
	s3interactor s3interactor.FileInteractor
}

func NewUserRepositoryFile(interactor s3interactor.FileInteractor) UserRepo {
	return &UserRepositoryFile{
		s3interactor: interactor,
	}
}

func (urf *UserRepositoryFile) Save(u User) error {
	users, errGet := urf.getUsers()

	if errGet != nil {
		return errGet
	}

	users = append(users, u)

	urf.saveUsers(users)

	return nil
}

func (urf *UserRepositoryFile) FindByEmail(e string) (User, error) {
	users, errGet := urf.getUsers()

	if errGet != nil {
		return User{}, errGet
	}

	for _, u := range users {
		if u.Email == e {
			return u, nil
		}
	}

	return User{}, nil
}

func (urf *UserRepositoryFile) Update(u User) (User, error) {
	users, errGet := urf.getUsers()

	if errGet != nil {
		return User{}, errGet
	}

	found := false
	for k, user := range users {
		if u.Email == user.Email {
			users[k] = u
			found = true
			break
		}
	}

	if found {
		urf.saveUsers(users)
		return u, nil
	}

	return User{}, &UserNotFoundError{}
}

func (urf *UserRepositoryFile) FindAll() ([]User, error) {
	users, errGet := urf.getUsers()

	if errGet != nil {
		return make([]User, 0), errGet
	}
	return users, nil
}

func (repo *UserRepositoryFile) getUsers() ([]User, error) {
	b, errGet := repo.s3interactor.GetJsonFile("users.json")

	if errGet != nil {
		return make([]User, 0), fmt.Errorf("Unable to get file: %w", errGet)
	}

	var users []User

	json.Unmarshal(b, &users)

	return users, nil
}

func (repo *UserRepositoryFile) saveUsers(u []User) error {
	b, errMarshal := json.Marshal(u)

	if errMarshal != nil {
		return fmt.Errorf("Unable to marshal users: %w", errMarshal)
	}

	repo.s3interactor.WriteJsonFile("users.json", b)

	return nil
}
