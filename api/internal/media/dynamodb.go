package media

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	dynamodbinteractor "github.com/hyoa/album/api/internal/dynamodbInteractor"
)

type MediaRepositoryDynamoDB struct {
	client *dynamodb.Client
	table  *string
}

type mediaModel struct {
	Key        string `dynamodbav:"mediaKey"`
	Author     string `dynamodbav:"author"`
	Kind       int    `dynamodbav:"mediaType"`
	Folder     string `dynamodbav:"folder"`
	UploadDate int    `dynamodbav:"uploadDate"`
	Visible    bool   `dynamodbav:"visible"`
}

func NewMediaRepositoryDynamoDB() MediaRepository {
	db, _ := dynamodbinteractor.NewInteractor()

	return &MediaRepositoryDynamoDB{
		client: db.Client,
		table:  aws.String(os.Getenv("MEDIA_TABLE_NAME")),
	}
}

func (mrd *MediaRepositoryDynamoDB) Save(media Media) error {
	var kind int
	if media.Kind == KindPhoto {
		kind = 1
	} else {
		kind = 2
	}

	data, errMarshal := attributevalue.MarshalMap(mediaModel{
		Key:        media.Key,
		Author:     media.Author,
		Kind:       kind,
		Folder:     media.Folder,
		UploadDate: media.UploadDate,
		Visible:    media.Visible,
	})

	if errMarshal != nil {
		return errMarshal
	}

	_, errSave := mrd.client.PutItem(context.Background(), &dynamodb.PutItemInput{
		TableName: mrd.table,
		Item:      data,
	})

	return errSave
}

func (mrd *MediaRepositoryDynamoDB) FindByFolder(folder string) ([]Media, error) {
	output, errQuery := mrd.client.Query(context.Background(), &dynamodb.QueryInput{
		TableName:              mrd.table,
		KeyConditionExpression: aws.String("#f = :f"),
		ExpressionAttributeNames: map[string]string{
			"#f": "folder",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":f": &types.AttributeValueMemberS{Value: folder},
		},
		IndexName: aws.String("folderIndex"),
	})

	if errQuery != nil {
		return make([]Media, 0), errQuery
	}

	var items []mediaModel
	errUnmarshal := attributevalue.UnmarshalListOfMaps(output.Items, &items)

	if errUnmarshal != nil {
		return make([]Media, 0), errUnmarshal
	}

	var medias []Media
	for _, m := range items {
		var kind MediaKind
		if m.Kind == 1 {
			kind = KindPhoto
		} else {
			kind = KindVideo
		}

		medias = append(medias, Media{
			Key:        m.Key,
			Author:     m.Author,
			Kind:       kind,
			Folder:     m.Folder,
			UploadDate: m.UploadDate,
			Visible:    m.Visible,
		})
	}

	return medias, nil
}

func (mrd *MediaRepositoryDynamoDB) FindFoldersName(name string) ([]string, error) {
	var output *dynamodb.ScanOutput
	var errScan error

	if name != "" {
		exprContain := expression.Name("folder").Contains(name)
		expr, errBuild := expression.NewBuilder().WithFilter(exprContain).Build()

		if errBuild != nil {
			return make([]string, 0), errBuild
		}

		output, errScan = mrd.client.Scan(context.Background(), &dynamodb.ScanInput{
			TableName:                 mrd.table,
			FilterExpression:          expr.Filter(),
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			IndexName:                 aws.String("folderIndex"),
		})
	} else {
		output, errScan = mrd.client.Scan(context.Background(), &dynamodb.ScanInput{
			TableName: mrd.table,
		})
	}

	if errScan != nil {
		return make([]string, 0), nil
	}

	var items []mediaModel
	errUnmarshal := attributevalue.UnmarshalListOfMaps(output.Items, &items)

	if errUnmarshal != nil {
		return make([]string, 0), errUnmarshal
	}

	foldersNames := make(map[string]bool)
	for _, m := range items {
		if _, ok := foldersNames[m.Folder]; !ok {
			foldersNames[m.Folder] = true
		}
	}

	var names []string
	for f := range foldersNames {
		names = append(names, f)
	}

	return names, nil
}

func (mrd *MediaRepositoryDynamoDB) FindAll() ([]Media, error) {
	return make([]Media, 0), nil
}

func (mrd *MediaRepositoryDynamoDB) FindManyByKeys(keys []string) ([]Media, error) {
	var medias []Media

	for _, k := range keys {
		media, errGet := mrd.FindByKey(k)

		if errGet == nil {
			medias = append(medias, media)
		}
	}

	return medias, nil
}

func (mrd *MediaRepositoryDynamoDB) FindByKey(key string) (Media, error) {
	output, errQuery := mrd.client.Query(context.Background(), &dynamodb.QueryInput{
		TableName:              mrd.table,
		KeyConditionExpression: aws.String("#k = :k"),
		ExpressionAttributeNames: map[string]string{
			"#k": "mediaKey",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":k": &types.AttributeValueMemberS{Value: key},
		},
	})

	if errQuery != nil {
		return Media{}, errQuery
	}

	if len(output.Items) > 0 {
		var items []mediaModel
		errUnmarshal := attributevalue.UnmarshalListOfMaps(output.Items, &items)

		if errUnmarshal != nil {
			return Media{}, errUnmarshal
		}

		var kind MediaKind
		if items[0].Kind == 1 {
			kind = KindPhoto
		} else {
			kind = KindVideo
		}

		return Media{
			Key:        items[0].Key,
			Author:     items[0].Author,
			Folder:     items[0].Folder,
			UploadDate: items[0].UploadDate,
			Visible:    items[0].Visible,
			Kind:       kind,
		}, nil
	}

	return Media{}, nil
}
