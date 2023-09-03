package main

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/cucumber/godog"
	gherkin_context "github.com/hyoa/album/api/gherkin/context"
	"github.com/hyoa/album/api/internal/awsinteractor"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

func TestFeatures(t *testing.T) {
	godotenv.Load("../.env.test")

	err := setupDB()

	if err != nil {
		t.Fatal("unable to start db", err)
	}

	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"./features"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}

	// tearDownDB()
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		purgeTable()
		loadFixtures()

		return ctx, nil
	})

	ctx.Step(`^I send a graphql request with payload:$`, gherkin_context.ISendAGraphqlRequestWithPayload)
	ctx.Step(`^the response should match json:$`, gherkin_context.TheResponseShouldMatchJson)
	ctx.Step(`^the response status code should be (\d+)`, gherkin_context.TheResponseStatusCodeShouldBe)
	ctx.Step(`^I query the DynamoDB table ([\w-]+) with keys:$`, gherkin_context.IQueryTheDynamoDBTableAlbumtestuserWithKeys)
	ctx.Step(`^I should have an entry in the DynamoDB query result with attributes:$`, gherkin_context.IShouldHaveAnEntryInTheDynamoDBQueryResultWithAttributes)
	ctx.Step(`^I should have (\d+) entry in the DynamoDB query result$`, gherkin_context.IShouldHaveEntryInTheDynamoDBQueryResult)
	ctx.Step(`^the response should contain an auth token with name "([^"]*)", email "([^"]*)" and role (\d+)$`, gherkin_context.TheResponseShouldContainAnAuthToken)
	ctx.Step(`^I check in the mailbox$`, gherkin_context.ICheckInTheMailbox)
	ctx.Step(`^I should have a mail that contain a password reset link for "([^"]*)" with subject "([^"]*)"$`, gherkin_context.IShouldHaveAMailThatContainAPasswordResetLinkForWithSubject)
	ctx.Step(`^I should have a mail that contain an invitation link for "([^"]*)" with subject "([^"]*)"$`, gherkin_context.IShouldHaveAMailThatContainAnInvitationLinkForWithSubject)
	ctx.Step(`^storage has key "([^"]*)"$`, gherkin_context.StorageHasKey)
	ctx.Step(`^I authenticate as an "([^"]*)"$`, gherkin_context.IAuthenticateAsAn)
	ctx.Step(`^I send a "([^"]*)" request to "([^"]*)" with payload:$`, gherkin_context.ISendARestRequestWithPayload)
}

func setupDB() error {
	db, _ := awsinteractor.NewDynamoDBInteractor()

	ouput, err := db.Client.ListTables(context.Background(), &dynamodb.ListTablesInput{})

	if err != nil {
		return err
	}

	userTableExist := false
	albumTableExist := false
	mediaTableExist := false
	for _, t := range ouput.TableNames {
		if t == os.Getenv("USER_TABLE_NAME") {
			userTableExist = true
		}

		if t == os.Getenv("ALBUM_TABLE_NAME") {
			albumTableExist = true
		}

		if t == os.Getenv("MEDIA_TABLE_NAME") {
			mediaTableExist = true
		}
	}

	if !userTableExist {
		_, errCreate := db.Client.CreateTable(context.Background(), &dynamodb.CreateTableInput{
			AttributeDefinitions: []types.AttributeDefinition{
				{AttributeName: aws.String("email"), AttributeType: types.ScalarAttributeTypeS},
			},
			KeySchema: []types.KeySchemaElement{
				{AttributeName: aws.String("email"), KeyType: types.KeyTypeHash},
			},
			TableName:   aws.String(os.Getenv("USER_TABLE_NAME")),
			BillingMode: types.BillingModePayPerRequest,
		})

		if errCreate != nil {
			panic(errCreate)
		}
	}

	if !albumTableExist {
		_, errCreate := db.Client.CreateTable(context.Background(), &dynamodb.CreateTableInput{
			AttributeDefinitions: []types.AttributeDefinition{
				{AttributeName: aws.String("slug"), AttributeType: types.ScalarAttributeTypeS},
			},
			KeySchema: []types.KeySchemaElement{
				{AttributeName: aws.String("slug"), KeyType: types.KeyTypeHash},
			},
			TableName:   aws.String(os.Getenv("ALBUM_TABLE_NAME")),
			BillingMode: types.BillingModePayPerRequest,
		})

		if errCreate != nil {
			panic(errCreate)
		}
	}

	if !mediaTableExist {
		_, errCreate := db.Client.CreateTable(context.Background(), &dynamodb.CreateTableInput{
			AttributeDefinitions: []types.AttributeDefinition{
				{AttributeName: aws.String("mediaKey"), AttributeType: types.ScalarAttributeTypeS},
				{AttributeName: aws.String("folder"), AttributeType: types.ScalarAttributeTypeS},
			},
			KeySchema: []types.KeySchemaElement{
				{AttributeName: aws.String("mediaKey"), KeyType: types.KeyTypeHash},
			},
			TableName:   aws.String(os.Getenv("MEDIA_TABLE_NAME")),
			BillingMode: types.BillingModePayPerRequest,
			GlobalSecondaryIndexes: []types.GlobalSecondaryIndex{
				{
					IndexName: aws.String("folderIndex"),
					KeySchema: []types.KeySchemaElement{
						{AttributeName: aws.String("folder"), KeyType: types.KeyTypeHash},
					},
					Projection: &types.Projection{ProjectionType: types.ProjectionTypeAll},
				},
			},
		})

		if errCreate != nil {
			panic(errCreate)
		}
	}

	return nil
}

func tearDownDB() error {
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{URL: "http://localhost:9966", PartitionID: "aws", SigningRegion: "eu-west-3"}, nil
	})

	cfg, errLoad := config.LoadDefaultConfig(
		context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("not-a-real-key", "@@not-a-real-secre", "")),
		config.WithRegion("eu-west-3"),
		config.WithEndpointResolverWithOptions(customResolver),
	)

	if errLoad != nil {
		log.Fatalf("unable to load SDK config, %v", errLoad)
	}

	svc := dynamodb.NewFromConfig(cfg)

	_, errDelete := svc.DeleteTable(context.Background(), &dynamodb.DeleteTableInput{
		TableName: aws.String(os.Getenv("USER_TABLE_NAME")),
	})

	return errDelete
}

func purgeTable() error {
	type table struct {
		name, key string
	}
	tables := []table{
		{name: os.Getenv("USER_TABLE_NAME"), key: "email"},
		{name: os.Getenv("ALBUM_TABLE_NAME"), key: "slug"},
		{name: os.Getenv("MEDIA_TABLE_NAME"), key: "mediaKey"},
	}

	for _, t := range tables {
		deleteItemsForTable(t.key, t.name)
	}

	return nil
}

func deleteItemsForTable(key, tableName string) error {
	db, _ := awsinteractor.NewDynamoDBInteractor()

	p := dynamodb.NewScanPaginator(db.Client, &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	})

	for p.HasMorePages() {
		out, err := p.NextPage(context.TODO())
		if err != nil {
			panic(err)
		}

		for _, item := range out.Items {
			_, err = db.Client.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
				TableName: aws.String(tableName),
				Key: map[string]types.AttributeValue{
					key: item[key],
				},
			})
			if err != nil {
				panic(err)
			}
		}
	}

	return nil
}

func loadFixtures() error {
	loadUserFixture()
	loadAlbumFixture()
	loadMediaFixture()

	return nil
}

func loadUserFixture() error {
	type user struct {
		Name             string `yaml:"name" dynamodbav:"name"`
		Email            string `yaml:"email" dynamodbav:"email"`
		RegistrationDate int    `yaml:"registrationDate" dynamodbav:"registrationDate"`
		Password         string `yaml:"password" dynamodbav:"password"`
		Role             int    `yaml:"userRole" dynamodbav:"userRole"`
	}

	var users []user
	yamlFile, _ := ioutil.ReadFile("./fixtures/user.yaml")

	yaml.Unmarshal(yamlFile, &users)

	db, _ := awsinteractor.NewDynamoDBInteractor()

	for _, u := range users {
		data, _ := attributevalue.MarshalMap(u)

		db.Client.PutItem(context.Background(), &dynamodb.PutItemInput{
			TableName: aws.String(os.Getenv("USER_TABLE_NAME")),
			Item:      data,
		})
	}

	return nil
}

func loadAlbumFixture() error {
	type media struct {
		Author   string `yaml:"author" dynamodbav:"author"`
		Key      string `yaml:"mediaKey" dynamodbav:"mediaKey"`
		Favorite bool   `yaml:"favorite" dynamodbav:"favorite"`
		Kind     int    `yaml:"mediaType" dynamodbav:"mediaType"`
	}

	type album struct {
		Slug          string  `yaml:"slug" dynamodbav:"slug"`
		Title         string  `yaml:"title" dynamodbav:"title"`
		Author        string  `yaml:"author" dynamodbav:"author"`
		CreattionDate int     `yaml:"creationDate" dynamodbav:"creationDate"`
		Description   string  `yaml:"description" dynamodbav:"description"`
		IsPrivate     bool    `yaml:"isPrivate" dynamodbav:"isPrivate"`
		Medias        []media `yaml:"medias" dynamodbav:"medias"`
	}

	var albums []album
	yamlFile, _ := ioutil.ReadFile("./fixtures/album.yaml")

	yaml.Unmarshal(yamlFile, &albums)

	db, _ := awsinteractor.NewDynamoDBInteractor()

	for _, a := range albums {
		data, _ := attributevalue.MarshalMap(a)

		db.Client.PutItem(context.Background(), &dynamodb.PutItemInput{
			TableName: aws.String(os.Getenv("ALBUM_TABLE_NAME")),
			Item:      data,
		})
	}

	return nil
}

func loadMediaFixture() error {
	type media struct {
		Author     string `yaml:"author" dynamodbav:"author"`
		Folder     string `yaml:"folder" dynamodbav:"folder"`
		UploadDate int    `yaml:"uploadDate" dynamodbav:"uploadDate"`
		Visible    bool   `yaml:"visible" dynamodbav:"visible"`
		Key        string `yaml:"mediaKey" dynamodbav:"mediaKey"`
		Kind       int    `yaml:"mediaType" dynamodbav:"mediaType"`
	}

	var medias []media
	yamlFile, _ := ioutil.ReadFile("./fixtures/media.yaml")

	yaml.Unmarshal(yamlFile, &medias)

	db, _ := awsinteractor.NewDynamoDBInteractor()

	for _, m := range medias {
		data, _ := attributevalue.MarshalMap(m)

		db.Client.PutItem(context.Background(), &dynamodb.PutItemInput{
			TableName: aws.String(os.Getenv("MEDIA_TABLE_NAME")),
			Item:      data,
		})
	}

	return nil
}
