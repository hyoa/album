package album

type AlbumRepositoryDynamoDB struct{}

func NewAlbumRepositoryDynamoDB() AlbumRepository {
	return &AlbumRepositoryDynamoDB{}
}

func (ard *AlbumRepositoryDynamoDB) Save(album Album) error {
	return nil
}

func (ard *AlbumRepositoryDynamoDB) FindBySlug(slug string) (Album, error) {
	return Album{}, nil
}

func (ard *AlbumRepositoryDynamoDB) Search(includePrivate, includeNoMedias bool, limit, offset int, term, order string) ([]Album, error) {
	return make([]Album, 0), nil
}

func (ard *AlbumRepositoryDynamoDB) Update(a Album) error {
	return nil
}

func (ard *AlbumRepositoryDynamoDB) DeleteBySlug(slug string) error {
	return nil
}

func (ard *AlbumRepositoryDynamoDB) FindAll() ([]Album, error) {
	return make([]Album, 0), nil
}
