package cdn

import (
	"context"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/aws/aws-sdk-go-v2/feature/cloudfront/sign"
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

func NewCDNAWSInteractor(s3 awsinteractor.S3Interactor) (CDNInteractor, error) {
	cache, _ := bigcache.New(context.Background(), bigcache.DefaultConfig(10*time.Minute))
	return &AwsCdn{s3Interactor: s3, cache: cache}, nil
}

type CDNInteractor interface {
	SignGetUri(key string, size MediaSize, kind MediaKind) string
}

type AwsCdn struct {
	s3Interactor awsinteractor.S3Interactor
	cache        *bigcache.BigCache
}

func (c *AwsCdn) SignGetUri(key string, size MediaSize, kind MediaKind) string {
	cacheKey := fmt.Sprintf("%s-%s", key, size)
	entry, errGet := c.cache.Get(cacheKey)

	if errGet == nil && len(entry) != 0 {
		return string(entry)
	}

	if kind == KindPhoto {
		pemString := fmt.Sprintf(`
-----BEGIN RSA PRIVATE KEY-----
%s
-----END RSA PRIVATE KEY-----`,
			os.Getenv("AWS_PK"))

		block, _ := pem.Decode([]byte(pemString))

		pk, errParse := x509.ParsePKCS1PrivateKey((block.Bytes))

		if errParse != nil {
			log.Fatal(errParse)
		}

		signer := sign.NewURLSigner(os.Getenv("KEY_PAIR_ID"), pk)

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

		url := fmt.Sprintf("%s/%s", os.Getenv("CDN_HOST"), json64)
		signedUrl, _ := signer.Sign(url, time.Now().Add(15*time.Minute))

		c.cache.Set(cacheKey, []byte(signedUrl))
		return signedUrl
	}

	signedUrl, _ := c.s3Interactor.SignGetUri(key, os.Getenv("BUCKET_VIDEO_FORMATTED"))
	c.cache.Set(cacheKey, []byte(signedUrl))

	return signedUrl
}
