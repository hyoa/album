package gherkin_context

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/cucumber/godog"
	dynamodbinteractor "github.com/hyoa/album/api/internal/dynamodbInteractor"
)

type testDynamoDBKey struct{}

func IQueryTheDynamoDBTableAlbumtestuserWithKeys(ctx context.Context, tableName string, arg1 *godog.Table) (context.Context, error) {
	type key struct {
		name, value string
	}

	rows := arg1.Rows[1:]

	var keys []key
	for _, r := range rows {
		keys = append(keys, key{name: r.Cells[0].Value, value: r.Cells[1].Value})
	}

	db, _ := dynamodbinteractor.NewInteractor()

	keysCondExp := make([]string, 0)
	eav := make(map[string]types.AttributeValue)

	for _, k := range keys {
		keysCondExp = append(keysCondExp, fmt.Sprintf("%s = :%s", k.name, k.name))
		eav[fmt.Sprintf(":%s", k.name)] = &types.AttributeValueMemberS{Value: k.value}
	}

	output, errQuery := db.Client.Query(context.Background(), &dynamodb.QueryInput{
		TableName:                 aws.String(tableName),
		ExpressionAttributeValues: eav,
		KeyConditionExpression:    aws.String(strings.Join(keysCondExp, ",")),
	})

	return context.WithValue(ctx, testDynamoDBKey{}, output), errQuery
}

func IShouldHaveAnEntryInTheDynamoDBQueryResultWithAttributes(ctx context.Context, arg1 *godog.Table) error {
	type entry struct {
		name, value, kind, condition string
	}

	rows := arg1.Rows[1:]

	var entries []entry
	for _, r := range rows {
		entries = append(entries, entry{
			name:      r.Cells[0].Value,
			value:     r.Cells[1].Value,
			kind:      r.Cells[2].Value,
			condition: r.Cells[3].Value,
		})
	}

	res, _ := ctx.Value(testDynamoDBKey{}).(*dynamodb.QueryOutput)

	foundCount := 0
	for _, e := range entries {
		for _, it := range res.Items {
			if v, ok := it[e.name]; ok {
				switch e.kind {
				case "string":
					var vString string
					attributevalue.Unmarshal(v, &vString)

					if e.condition == "equal" && vString == e.value {
						foundCount++
					} else if e.condition == "notEmpty" && vString != "" {
						foundCount++
					}
				case "int":
					var vInt int
					eInt, _ := strconv.Atoi(e.value)
					attributevalue.Unmarshal(v, &vInt)

					if e.condition == "equal" && vInt == eInt {
						foundCount++
					} else if e.condition == "notEmpty" && vInt >= 0 {
						foundCount++
					}
				case "bool":
					var vBool bool
					eBool, _ := strconv.ParseBool(e.value)
					attributevalue.Unmarshal(v, &vBool)

					if e.condition == "equal" && vBool == eBool {
						foundCount++
					}
				default:
					log.Panic("Wrong type")
				}
			}
		}
	}

	if foundCount != len(entries) {
		return fmt.Errorf("query result does not contain entry")
	}

	return nil
}

func IShouldHaveEntryInTheDynamoDBQueryResult(ctx context.Context, count int) error {
	res, _ := ctx.Value(testDynamoDBKey{}).(*dynamodb.QueryOutput)

	if len(res.Items) != count {
		return fmt.Errorf("Result count is invalid. Expected %d got %d", count, len(res.Items))
	}

	return nil
}
