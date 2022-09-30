package media

import "github.com/hyoa/album/api/internal/s3interactor"

type S3Storage struct {
	interactor s3interactor.S3Interactor
}

func NewS3Storage(interactor s3interactor.S3Interactor) Storage {
	return &S3Storage{
		interactor: interactor,
	}
}

func (s *S3Storage) MediaExist(key, bucket string) (bool, error) {
	return s.interactor.MediaExist(key, bucket)
}

func (s *S3Storage) SignUploadUri(key, bucket string) (string, error) {
	return s.interactor.SignPutUri(key, bucket)
}
