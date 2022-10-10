package dynamodbinteractor

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type DynamoDBInteractor struct {
	Client *dynamodb.Client
}

func NewInteractor() (DynamoDBInteractor, error) {
	if os.Getenv("DYNAMODB_ENDPOINT") != "" {
		customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{URL: os.Getenv("DYNAMODB_ENDPOINT"), PartitionID: "aws", SigningRegion: "eu-west-3"}, nil
		})

		cfg, errLoad := config.LoadDefaultConfig(
			context.TODO(),
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(os.Getenv("AKID"), os.Getenv("ASK"), "")),
			config.WithRegion("eu-west-3"),
			config.WithEndpointResolverWithOptions(customResolver),
		)

		return DynamoDBInteractor{Client: dynamodb.NewFromConfig(cfg)}, errLoad
	}

	cfg, errLoad := config.LoadDefaultConfig(
		context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(os.Getenv("AKID"), os.Getenv("ASK"), "")),
		config.WithRegion("eu-west-2"),
	)

	return DynamoDBInteractor{Client: dynamodb.NewFromConfig(cfg)}, errLoad
}
