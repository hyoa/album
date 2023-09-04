package cdn

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/patrickmn/go-cache"

	"github.com/hyoa/album/api/internal/awsinteractor"
)

type MediaKind string
type MediaSize string

const (
	KindPhoto MediaKind = "photo"
	KindVideo MediaKind = "video"

	SizeSmall  MediaSize = "small"
	SizeMedium MediaSize = "medium"
	SizeLarge  MediaSize = "large"
)

type Edits struct {
	Rotate interface{} `json:"rotate"`
	Resize Resize      `json:"resize"`
}

type Resize struct {
	Fit   string `json:"fit"`
	Width int    `json:"width"`
}

type CDNData struct {
	Bucket string `json:"bucket"`
	Key    string `json:"key"`
	Edits  Edits  `json:"edits"`
}

func NewCDNAWSInteractor(s3 awsinteractor.S3Interactor, cache cache.Cache, secret string) (CDNInteractor, error) {
	return &AwsCdn{s3Interactor: s3, cache: cache, secret: secret}, nil
}

type CDNInteractor interface {
	SignGetUri(key string, size MediaSize, kind MediaKind) string
}

type AwsCdn struct {
	s3Interactor awsinteractor.S3Interactor
	cache        cache.Cache
	secret       string
}

func (c *AwsCdn) SignGetUri(key string, size MediaSize, kind MediaKind) string {
	cacheKey := fmt.Sprintf("%s-%s", key, size)
	entry, exist := c.cache.Get(cacheKey)

	if exist {
		return entry.(string)
	}

	if kind == KindPhoto {
		var width int
		switch size {
		case SizeSmall:
			width = 400
		case SizeMedium:
			width = 800
		case SizeLarge:
			width = 1024
		}

		data := CDNData{
			Bucket: os.Getenv("BUCKET_IMAGE"),
			Key:    key,
			Edits: Edits{
				Rotate: nil,
				Resize: Resize{
					Fit:   "cover",
					Width: width,
				},
			},
		}

		json, _ := json.Marshal(data)
		json64 := base64.StdEncoding.EncodeToString([]byte(json))

		path := fmt.Sprintf("/%s", json64)

		h := hmac.New(sha256.New, []byte(c.secret))
		h.Write([]byte(path))
		signature := hex.EncodeToString(h.Sum(nil))

		signedUrl := fmt.Sprintf("%s/%s?signature=%s", os.Getenv("CDN_HOST"), json64, signature)

		c.cache.Set(cacheKey, signedUrl, time.Minute*10)

		return signedUrl
	}

	signedUrl, _ := c.s3Interactor.SignGetUri(key, os.Getenv("BUCKET_VIDEO_FORMATTED"))
	c.cache.Set(cacheKey, signedUrl, time.Minute*10)

	return signedUrl
}
