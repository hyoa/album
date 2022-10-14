package media

type Storage interface {
	MediaExist(key, bucket string) (bool, error)
	SignUploadUri(key, bucket string) (string, error)
	SignDownloadUri(key, bucket string) (string, error)
}
