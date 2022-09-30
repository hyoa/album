package cdn

import (
	"encoding/base64"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
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

type mediaToken struct {
	Key  string `json:"key"`
	Kind string `json:"type"`
	Size string `json:"size"`
}

type mediaTokenCustomClaims struct {
	mediaToken
	jwt.StandardClaims
}

func SignGetUri(key string, size MediaSize, kind MediaKind) string {
	t := mediaToken{}

	var path string
	if kind == KindPhoto {
		t.Size = string(size)
		path = fmt.Sprintf(os.Getenv("CDN_IMAGE_PATH"), string(size), key)
	} else {
		path = fmt.Sprintf(os.Getenv("CDN_VIDEO_PATH"), key)
	}

	t.Key = key
	t.Kind = string(kind)

	secret := []byte(os.Getenv("CDN_SECRET"))
	claims := mediaTokenCustomClaims{
		t,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + 3600,
			IssuedAt:  time.Now().Unix(),
			Issuer:    os.Getenv("APPLICATION_NAME"),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwtSigned, _ := token.SignedString(secret)

	return fmt.Sprintf("%s%s?key=%s", os.Getenv("CDN_HOST"), path, base64.StdEncoding.EncodeToString([]byte(jwtSigned)))
}
