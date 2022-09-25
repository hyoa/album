package user

type UserRepo interface {
	Save(u User) error
	FindByEmail(e string) (User, error)
	Update(u User) (User, error)
	FindAll() ([]User, error)
}
