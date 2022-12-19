package cdn

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

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
	Key             string `json:"key"`
	Kind            string `json:"type"`
	Size            string `json:"size"`
	ApplicationName string `json:"applicationName"`
}

type mediaTokenCustomClaims struct {
	mediaToken
	jwt.StandardClaims
}

func SignGetUri(key string, size MediaSize, kind MediaKind) string {
	t := mediaToken{}

	// var path string
	if kind == KindPhoto {
		t.Size = string(size)
		// path = fmt.Sprintf(os.Getenv("CDN_IMAGE_PATH"), string(size), key)
	}
	// } else {
	// 	// path = fmt.Sprintf(os.Getenv("CDN_VIDEO_PATH"), key)
	// }

	t.Key = key
	t.Kind = string(kind)
	t.ApplicationName = os.Getenv("APP_NAME")

	// secret := []byte(os.Getenv("CDN_SECRET"))
	// claims := mediaTokenCustomClaims{
	// 	t,
	// 	jwt.StandardClaims{
	// 		ExpiresAt: time.Now().Unix() + 3600,
	// 		IssuedAt:  time.Now().Unix(),
	// 		Issuer:    os.Getenv("APPLICATION_NAME"),
	// 	},
	// }

	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// jwtSigned, _ := token.SignedString(secret)

	json, _ := json.Marshal(t)

	return fmt.Sprintf("%s/%s", os.Getenv("CDN_HOST"), encrypt(json))
}

func encrypt(content []byte) string {
	key := []byte("6368616e676520746869732070617373")

	// Create the AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	plaintext, _ := pkcs7Pad(content, block.BlockSize())
	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}
	bm := cipher.NewCBCEncrypter(block, iv)
	bm.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)

	return fmt.Sprintf("%x", ciphertext)
}

func pkcs7Pad(b []byte, blocksize int) ([]byte, error) {
	if blocksize <= 0 {
		return nil, errors.New("invalid blocksize")
	}
	if len(b) == 0 {
		return nil, errors.New("invalid PKCS7 data (empty or not padded)")
	}
	n := blocksize - (len(b) % blocksize)
	pb := make([]byte, len(b)+n)
	copy(pb, b)
	copy(pb[len(b):], bytes.Repeat([]byte{byte(n)}, n))
	return pb, nil
}
