package s3interactor

import (
	"bytes"
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type FileInteractor interface {
	GetJsonFile(fileName string) ([]byte, error)
	WriteJsonFile(fileName string, jsonByte []byte) error
}

type interactor struct {
	client *minio.Client
}

func NewInteractor(endpoint, keyId, keySecret string) (FileInteractor, error) {
	client, errNew := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(keyId, keySecret, ""),
		Secure: true,
	})

	if errNew != nil {
		return &interactor{}, fmt.Errorf("Unable to create client: %w", errNew)
	}

	return &interactor{
		client: client,
	}, nil
}

func (i *interactor) GetJsonFile(fileName string) ([]byte, error) {
	ctxt := context.Background()
	obj, errGet := i.client.GetObject(ctxt, "current", fileName, minio.GetObjectOptions{})

	if errGet != nil {
		return make([]byte, 0), fmt.Errorf("Unable to get file %w", errGet)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(obj)
	return buf.Bytes(), nil
}

func (i *interactor) WriteJsonFile(fileName string, jsonByte []byte) error {
	ctxt := context.Background()

	reader := bytes.NewReader(jsonByte)

	_, err := i.client.PutObject(ctxt, "current", fileName, reader, reader.Size(), minio.PutObjectOptions{})

	return err
}
