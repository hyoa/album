package album

import (
	"context"
	"os"
	"sort"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	dynamodbinteractor "github.com/hyoa/album/api/internal/dynamodbInteractor"
)

type AlbumRepositoryDynamoDB struct {
	client *dynamodb.Client
	table  *string
}

type albumModel struct {
	Title        string       `dynamodbav:"title"`
	Description  string       `dynamodbav:"description"`
	Private      bool         `dynamodbav:"isPrivate"`
	Author       string       `dynamodbav:"author"`
	CreationDate int          `dynamodbav:"creationDate"`
	Slug         string       `dynamodbav:"slug"`
	Medias       []mediaModel `dynamodbav:"medias"`
}

type mediaModel struct {
	Key      string `dynamodbav:"mediaKey"`
	Author   string `dynamodbav:"author"`
	Kind     int    `dynamodbav:"mediaType"`
	Favorite bool   `dynamodbav:"favorite"`
}

func NewAlbumRepositoryDynamoDB() AlbumRepository {
	db, _ := dynamodbinteractor.NewInteractor()

	return &AlbumRepositoryDynamoDB{
		client: db.Client,
		table:  aws.String(os.Getenv("ALBUM_TABLE_NAME")),
	}
}

func (ard *AlbumRepositoryDynamoDB) Save(album Album) error {
	var medias []mediaModel

	for _, m := range album.Medias {
		var kind int
		if m.Kind == KindPhoto {
			kind = 1
		} else {
			kind = 2
		}

		medias = append(medias, mediaModel{
			Key:      m.Key,
			Author:   m.Author,
			Kind:     kind,
			Favorite: m.Favorite,
		})
	}

	data, errMarshal := attributevalue.MarshalMap(albumModel{
		Title:        album.Title,
		Description:  album.Description,
		Author:       album.Author,
		Private:      album.Private,
		Slug:         album.Slug,
		CreationDate: album.CreationDate,
		Medias:       medias,
	})

	if errMarshal != nil {
		return errMarshal
	}

	_, errSave := ard.client.PutItem(context.Background(), &dynamodb.PutItemInput{
		TableName: ard.table,
		Item:      data,
	})

	return errSave
}

func (ard *AlbumRepositoryDynamoDB) FindBySlug(slug string) (Album, error) {
	output, errQuery := ard.client.Query(context.Background(), &dynamodb.QueryInput{
		TableName:              ard.table,
		KeyConditionExpression: aws.String("#s = :s"),
		ExpressionAttributeNames: map[string]string{
			"#s": "slug",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":s": &types.AttributeValueMemberS{Value: slug},
		},
	})

	if errQuery != nil {
		return Album{}, errQuery
	}

	var items []albumModel

	if len(output.Items) > 0 {
		errUnmarshal := attributevalue.UnmarshalListOfMaps(output.Items, &items)

		if errUnmarshal != nil {
			return Album{}, errUnmarshal
		}

		var medias []Media
		for _, m := range items[0].Medias {

			var kind MediaKind
			if m.Kind == 1 {
				kind = KindPhoto
			} else {
				kind = KindVideo
			}

			medias = append(medias, Media{
				Key:      m.Key,
				Author:   m.Author,
				Kind:     kind,
				Favorite: m.Favorite,
			})
		}
		return Album{
			Title:        items[0].Title,
			Description:  items[0].Description,
			Slug:         items[0].Slug,
			Author:       items[0].Author,
			CreationDate: items[0].CreationDate,
			Private:      items[0].Private,
			Medias:       medias,
		}, nil
	}

	return Album{}, nil
}

func (ard *AlbumRepositoryDynamoDB) Search(includePrivate, includeNoMedias bool, limit, offset int, term, order string) ([]Album, error) {
	var globalExpr expression.ConditionBuilder

	if !includePrivate {
		globalExpr = expression.Name("isPrivate").Equal(expression.Value(false))
	} else {
		globalExpr = expression.Name("isPrivate").Equal(expression.Value(false)).Or(expression.Name("isPrivate").Equal(expression.Value(true)))
	}

	if !includeNoMedias {
		exprIncludeNoMedias := expression.Name("medias").Size().GreaterThan(expression.Value(0))
		globalExpr = globalExpr.And(exprIncludeNoMedias)
	}

	if term != "" {
		globalExpr = globalExpr.And(expression.Name("title").Contains(term).Or(expression.Name("description").Contains(term)))
	}

	expr, errBuild := expression.NewBuilder().WithFilter(globalExpr).Build()

	if errBuild != nil {
		return make([]Album, 0), errBuild
	}

	output, errScan := ard.client.Scan(context.Background(), &dynamodb.ScanInput{
		TableName:                 ard.table,
		FilterExpression:          expr.Filter(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})

	if errScan != nil {
		return make([]Album, 0), errScan
	}

	if len(output.Items) == 0 {
		return make([]Album, 0), nil
	}

	if len(output.Items) < limit {
		limit = len(output.Items)
	}

	itemsWithLimitOffset := output.Items[offset:limit]
	var items []albumModel
	errUnmarshal := attributevalue.UnmarshalListOfMaps(itemsWithLimitOffset, &items)

	if errUnmarshal != nil {
		return make([]Album, 0), errUnmarshal
	}

	var albums []Album
	for _, item := range items {
		var medias []Media

		for _, media := range item.Medias {
			var kind MediaKind
			if media.Kind == 1 {
				kind = KindPhoto
			} else {
				kind = KindVideo
			}

			medias = append(medias, Media{
				Key:      media.Key,
				Author:   media.Author,
				Kind:     kind,
				Favorite: media.Favorite,
			})
		}

		albums = append(albums, Album{
			Title:        item.Title,
			Description:  item.Description,
			Private:      item.Private,
			Author:       item.Author,
			CreationDate: item.CreationDate,
			Slug:         item.Slug,
			Medias:       medias,
		})
	}

	sort.SliceStable(albums, func(i, j int) bool {
		if order == "asc" {
			return albums[i].CreationDate < albums[j].CreationDate
		} else {
			return albums[i].CreationDate > albums[j].CreationDate
		}
	})

	return albums, nil
}

func (ard *AlbumRepositoryDynamoDB) Update(a Album) error {
	return nil
}

func (ard *AlbumRepositoryDynamoDB) DeleteBySlug(slug string) error {
	_, err := ard.client.DeleteItem(context.Background(), &dynamodb.DeleteItemInput{
		TableName: ard.table,
		Key: map[string]types.AttributeValue{
			"slug": &types.AttributeValueMemberS{Value: slug},
		},
	})

	return err
}

func (ard *AlbumRepositoryDynamoDB) FindAll() ([]Album, error) {
	return make([]Album, 0), nil
}
