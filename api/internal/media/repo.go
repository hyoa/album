package media

type MediaRepository interface {
	Save(media Media) error
	FindByFolder(folder string) ([]Media, error)
	FindFoldersName(name string) ([]string, error)
	FindAll() ([]Media, error)
	FindManyByKeys(keys []string) ([]Media, error)
	FindByKey(key string) (Media, error)
}
