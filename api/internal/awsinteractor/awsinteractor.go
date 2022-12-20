package awsinteractor

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	awsCredential "github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/minio/minio-go/v7"
	minioCredential "github.com/minio/minio-go/v7/pkg/credentials"
)

type DynamoDBInteractor struct {
	Client *dynamodb.Client
}

func NewDynamoDBInteractor() (DynamoDBInteractor, error) {
	if os.Getenv("DYNAMODB_ENDPOINT") != "" {
		customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{URL: os.Getenv("DYNAMODB_ENDPOINT"), PartitionID: "aws", SigningRegion: "eu-west-3"}, nil
		})

		cfg, errLoad := config.LoadDefaultConfig(
			context.TODO(),
			config.WithCredentialsProvider(awsCredential.NewStaticCredentialsProvider(os.Getenv("AKID"), os.Getenv("ASK"), "")),
			config.WithRegion("eu-west-3"),
			config.WithEndpointResolverWithOptions(customResolver),
		)

		return DynamoDBInteractor{Client: dynamodb.NewFromConfig(cfg)}, errLoad
	}

	cfg, errLoad := config.LoadDefaultConfig(
		context.TODO(),
		config.WithCredentialsProvider(awsCredential.NewStaticCredentialsProvider(os.Getenv("AKID"), os.Getenv("ASK"), "")),
		config.WithRegion("eu-west-2"),
	)

	return DynamoDBInteractor{Client: dynamodb.NewFromConfig(cfg)}, errLoad
}

type S3Interactor interface {
	GetJsonFile(fileName string) ([]byte, error)
	WriteJsonFile(fileName string, jsonByte []byte) error
	SignPutUri(key, bucket string) (string, error)
	SignGetUri(key, bucket string) (string, error)
	MediaExist(key, bucket string) (bool, error)
}

type interactor struct {
	client *minio.Client
}

func NewS3Interactor(endpoint, keyId, keySecret string) (S3Interactor, error) {
	client, errNew := minio.New(endpoint, &minio.Options{
		Creds:  minioCredential.NewStaticV4(keyId, keySecret, ""),
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

func (i *interactor) SignPutUri(key, bucket string) (string, error) {
	u, err := i.client.PresignedPutObject(context.Background(), bucket, key, 15*time.Minute)
	return u.String(), err
}

func (i *interactor) SignGetUri(key, bucket string) (string, error) {
	u, err := i.client.PresignedGetObject(context.Background(), bucket, key, 15*time.Minute, url.Values{})
	return u.String(), err
}

func (i *interactor) MediaExist(key, bucket string) (bool, error) {
	u, err := i.client.StatObject(context.Background(), bucket, key, minio.GetObjectOptions{})

	if err != nil {
		return false, err
	}

	if u.Key != "" {
		return true, nil
	}

	return false, nil
}
