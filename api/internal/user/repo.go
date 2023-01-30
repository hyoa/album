package user

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/hyoa/album/api/internal/awsinteractor"
)

type UserRepo interface {
	Save(u User) error
	FindByEmail(e string) (User, error)
	Update(u User) (User, error)
	FindAll() ([]User, error)
}

type UserRepositoryDynamoDB struct {
	client *dynamodb.Client
	table  *string
}

type userModel struct {
	Email      string `dynamodbav:"email"`
	Name       string `dynamodbav:"name"`
	Password   string `dynamodbav:"password"`
	CreateDate int    `dynamodbav:"registrationDate"`
	Role       int    `dynamodbav:"userRole"`
}

func NewUserRepositoryDynamoDB() UserRepo {
	db, _ := awsinteractor.NewDynamoDBInteractor()

	return &UserRepositoryDynamoDB{
		client: db.Client,
		table:  aws.String(os.Getenv("USER_TABLE_NAME")),
	}
}

func (urd *UserRepositoryDynamoDB) Save(u User) error {
	data, errMarshal := attributevalue.MarshalMap(userModel{
		Email:      u.Email,
		Name:       u.Name,
		Password:   u.Password,
		CreateDate: int(u.CreateDate),
		Role:       int(u.Role),
	})

	if errMarshal != nil {
		return errMarshal
	}

	_, errSave := urd.client.PutItem(context.Background(), &dynamodb.PutItemInput{
		TableName: urd.table,
		Item:      data,
	})

	return errSave
}

func (urd *UserRepositoryDynamoDB) FindByEmail(e string) (User, error) {
	output, errQuery := urd.client.Query(context.Background(), &dynamodb.QueryInput{
		TableName:              urd.table,
		KeyConditionExpression: aws.String("#e = :e"),
		ExpressionAttributeNames: map[string]string{
			"#e": "email",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":e": &types.AttributeValueMemberS{Value: e},
		},
	})

	if errQuery != nil {
		return User{}, errQuery
	}

	var items []userModel

	if len(output.Items) > 0 {
		errUnmarshal := attributevalue.UnmarshalListOfMaps(output.Items, &items)

		if errUnmarshal != nil {
			return User{}, errUnmarshal
		}

		return User{
			Name:       items[0].Name,
			Email:      items[0].Email,
			Role:       Role(items[0].Role),
			Password:   items[0].Password,
			CreateDate: int64(items[0].CreateDate),
		}, nil
	}

	return User{}, nil
}

func (urd *UserRepositoryDynamoDB) Update(u User) (User, error) {
	return u, urd.Save(u)
}

func (urd *UserRepositoryDynamoDB) FindAll() ([]User, error) {
	output, errScan := urd.client.Scan(context.Background(), &dynamodb.ScanInput{TableName: urd.table})

	if errScan != nil {
		return make([]User, 0), errScan
	}

	items := make([]userModel, 0)
	errUnmarshal := attributevalue.UnmarshalListOfMaps(output.Items, &items)

	if errUnmarshal != nil {
		return make([]User, 0), errUnmarshal
	}

	var users []User
	for k := range items {
		users = append(users, User{
			Name:       items[k].Name,
			Email:      items[k].Email,
			Password:   items[k].Password,
			Role:       Role(items[k].Role),
			CreateDate: int64(items[k].CreateDate),
		})
	}

	return users, nil
}
