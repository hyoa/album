package gherkin_context

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http/httptest"
	"reflect"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/cucumber/godog"
	"github.com/gin-gonic/gin"
	"github.com/hyoa/album/api/gherkin/mock"
	"github.com/hyoa/album/api/graph"
	"github.com/hyoa/album/api/graph/generated"
	"github.com/hyoa/album/api/internal/album"
	"github.com/hyoa/album/api/internal/media"
	"github.com/hyoa/album/api/internal/user"
)

type testHttpKey struct{}

func setUpRouter(storage *mock.Storage) *gin.Engine {
	router := gin.Default()
	router.POST("/query", graphqlHandler(storage))

	return router
}

func graphqlHandler(storage *mock.Storage) gin.HandlerFunc {
	r := &graph.Resolver{}
	mailer := mock.Mailer{}

	r.UserManager = user.CreateUserManager(user.NewUserRepositoryDynamoDB(), &mailer)
	r.AlbumManager = album.CreateAlbumManager(album.NewAlbumRepositoryDynamoDB())
	r.MediaManager = media.CreateMediaManager(media.NewMediaRepositoryDynamoDB(), storage)

	h := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: r}))

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func StorageHasKey(ctx context.Context, key string) (context.Context, error) {
	return context.WithValue(ctx, testHttpKey{}, key), nil
}

func ISendAGraphqlRequestWithPayload(ctx context.Context, arg1 *godog.DocString) (context.Context, error) {
	key, _ := ctx.Value(testHttpKey{}).(string)

	storageMock := &mock.Storage{Keys: append(make([]string, 0), key)}

	if arg1.Content == "" {
		return ctx, errors.New("no payload")
	}

	jsonData := map[string]string{
		"query": arg1.Content,
	}

	jsonValue, _ := json.Marshal(jsonData)
	req := httptest.NewRequest("POST", "/query", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	r := setUpRouter(storageMock)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	return context.WithValue(ctx, testHttpKey{}, w), nil
}

func TheResponseShouldMatchJson(ctx context.Context, arg1 *godog.DocString) error {
	res, _ := ctx.Value(testHttpKey{}).(*httptest.ResponseRecorder)

	var expected, actual interface{}

	if errExpected := json.Unmarshal([]byte(arg1.Content), &expected); errExpected != nil {

		return errExpected
	}

	if errActual := json.Unmarshal(res.Body.Bytes(), &actual); errActual != nil {

		return errActual
	}

	if !reflect.DeepEqual(expected, actual) {
		return fmt.Errorf("expected JSON does not match actual, %v vs. %v", expected, actual)
	}

	return nil
}

func TheResponseStatusCodeShouldBe(ctx context.Context, statusCode int) error {
	res, _ := ctx.Value(testHttpKey{}).(*httptest.ResponseRecorder)

	if statusCode != res.Code {
		return fmt.Errorf("status code should %d, got %d", statusCode, res.Code)
	}

	return nil
}

func TheResponseShouldContainAnAuthToken(ctx context.Context, name, email string, role int) error {
	res, _ := ctx.Value(testHttpKey{}).(*httptest.ResponseRecorder)

	type authResponse struct {
		Data struct {
			Auth struct {
				Token string `json:"token"`
			} `json:"auth"`
		} `json:"data"`
	}

	var auth authResponse
	if errDecode := json.Unmarshal(res.Body.Bytes(), &auth); errDecode != nil {
		return errDecode
	}

	if auth.Data.Auth.Token == "" {
		return errors.New("token cannot be empty")
	}

	tokenizer := user.CreateAuthTokenizer()
	t, _ := tokenizer.Decode(auth.Data.Auth.Token)

	if t.Email != email || t.Name != name || t.Role != user.Role(role) {
		return errors.New("token does not contain valid data")
	}

	return nil
}
