package media

import (
	"github.com/hyoa/album/api/internal/awsinteractor"
)

type Storage interface {
	MediaExist(key, bucket string) (bool, error)
	SignUploadUri(key, bucket string) (string, error)
	SignDownloadUri(key, bucket string) (string, error)
}

type S3Storage struct {
	interactor awsinteractor.S3Interactor
}

func NewS3Storage(interactor awsinteractor.S3Interactor) Storage {
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

func (s *S3Storage) SignDownloadUri(key, bucket string) (string, error) {
	return s.interactor.SignGetUri(key, bucket)
}
