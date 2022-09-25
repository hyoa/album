package media

type Storage interface {
	MediaExist(key string) (bool, error)
	SignUploadUri(key string) (string, error)
}
