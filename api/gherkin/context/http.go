package gherkin_context

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http/httptest"
	"reflect"

	"github.com/cucumber/godog"
	"github.com/gin-gonic/gin"
	"github.com/hyoa/album/api/controller"
	"github.com/hyoa/album/api/gherkin/mock"
	"github.com/hyoa/album/api/internal/album"
	"github.com/hyoa/album/api/internal/media"
	"github.com/hyoa/album/api/internal/translator"
	"github.com/hyoa/album/api/internal/user"
)

type testHttpKey struct{}

func setUpRouter(storage *mock.Storage) *gin.Engine {
	translatorManager := *translator.CreateTranslator("../i18n/active.fr.toml")
	mailer := mock.Mailer{Translator: translatorManager}
	converter := mock.VideoConverter{}

	userManager := user.CreateUserManager(user.NewUserRepositoryDynamoDB(), &mailer)
	albumManager := album.CreateAlbumManager(album.NewAlbumRepositoryDynamoDB())
	mediaManager := media.CreateMediaManager(media.NewMediaRepositoryDynamoDB(), storage, &converter)

	restController := controller.CreateRestController(mediaManager)

	router := gin.Default()
	router.POST("/v3/query", controller.GraphqlHandler(userManager, albumManager, mediaManager, &translatorManager, &mock.CDN{}))
	router.POST("/v3/video/acknowledge/cloudconvert", restController.AcknowledgeCloudconvertCall)

	return router
}

func StorageHasKey(ctx context.Context, key string) (context.Context, error) {
	return context.WithValue(ctx, testHttpKey{}, key), nil
}

func ISendAGraphqlRequestWithPayload(ctx context.Context, arg1 *godog.DocString) (context.Context, error) {
	key, _ := ctx.Value(testHttpKey{}).(string)
	jwt, okJwt := ctx.Value("AuthJwt").(string)

	storageMock := &mock.Storage{Keys: append(make([]string, 0), key)}

	if arg1.Content == "" {
		return ctx, errors.New("no payload")
	}

	jsonData := map[string]string{
		"query": arg1.Content,
	}

	jsonValue, _ := json.Marshal(jsonData)
	req := httptest.NewRequest("POST", "/v3/query", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	if okJwt {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwt))
	}

	r := setUpRouter(storageMock)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	return context.WithValue(ctx, testHttpKey{}, w), nil
}

func ISendARestRequestWithPayload(ctx context.Context, method, url string, payload *godog.DocString) (context.Context, error) {
	key, _ := ctx.Value(testHttpKey{}).(string)
	jwt, okJwt := ctx.Value("AuthJwt").(string)

	storageMock := &mock.Storage{Keys: append(make([]string, 0), key)}

	if payload.Content == "" {
		return ctx, errors.New("no payload")
	}

	jsonStr := []byte(payload.Content)

	req := httptest.NewRequest(method, url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	if okJwt {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwt))
	}

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

func IAuthenticateAsAn(ctx context.Context, role string) (context.Context, error) {
	return context.WithValue(ctx, "AuthJwt", "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJOYW1lIjoiY2hlY2siLCJFbWFpbCI6InRvdG9AdG90by5jb20iLCJSb2xlIjo5LCJleHAiOjE4MjMwMzI4MTEsImlhdCI6MTY2MjcyNTA3OCwiaXNzIjoiYXBpdjMifQ.WUAIpDnWW3DOgToq_VNGqxhZu3X6bhDKVMbChVhBbZFLdEEALxYIZmvPVo02fbtmAlaV_pAO8tLCJjd3qZUf7w"), nil
}
