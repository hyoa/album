package album

type AlbumRepository interface {
	Save(a Album) error
	FindBySlug(slug string) (Album, error)
	Search(includePrivate, includeNoMedias bool, limit, offset int, term, order string) ([]Album, error)
	Update(a Album) error
	DeleteBySlug(slug string) error
	FindAll() ([]Album, error)
}
