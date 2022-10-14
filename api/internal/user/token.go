package user

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

type authToken struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  Role   `json:"role"`
}

type resetToken struct {
	Email string `json:"email"`
}

type token interface {
	authToken | resetToken
}

type authTokenInput struct {
	u User
}

type resetTokenInput struct {
	email string
}

type tokenInput interface {
	authTokenInput | resetTokenInput
}

type authTokenCustomClaims struct {
	authToken
	jwt.StandardClaims
}

type resetTokenCustomClaims struct {
	resetToken
	jwt.StandardClaims
}

type Tokenizer[T token, I tokenInput] interface {
	create(i I) (T, error)
	stringify(t T) (string, error)
	Decode(t string) (T, error)
}

type authTokenizer struct{}

func CreateAuthTokenizer() Tokenizer[authToken, authTokenInput] {
	return authTokenizer{}
}

func (a authTokenizer) create(i authTokenInput) (authToken, error) {
	return authToken{
		Name:  i.u.Name,
		Email: i.u.Email,
		Role:  i.u.Role,
	}, nil
}

func (a authTokenizer) stringify(token authToken) (string, error) {

	key := []byte("secret")
	claims := authTokenCustomClaims{
		token,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + 3600,
			IssuedAt:  time.Now().Unix(),
			Issuer:    "apiv3",
		},
	}

	tokenEncoded := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	ss, err := tokenEncoded.SignedString(key)

	return ss, err
}

func (a authTokenizer) Decode(token string) (authToken, error) {
	tokenParsed, errParse := jwt.ParseWithClaims(token, &authTokenCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})

	if errParse != nil {
		return authToken{}, errParse
	}

	if claims, ok := tokenParsed.Claims.(*authTokenCustomClaims); ok && tokenParsed.Valid {
		return authToken{Name: claims.Name, Email: claims.Email, Role: claims.Role}, nil
	}

	return authToken{}, errParse
}

type resetTokenizer struct{}

func CreateResetTokenizer() Tokenizer[resetToken, resetTokenInput] {
	return resetTokenizer{}
}

func (r resetTokenizer) create(i resetTokenInput) (resetToken, error) {
	return resetToken{
		Email: i.email,
	}, nil
}

func (r resetTokenizer) Decode(token string) (resetToken, error) {
	tokenParsed, err := jwt.ParseWithClaims(token, &resetTokenCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})

	if err != nil {
		return resetToken{}, err
	}

	if claims, ok := tokenParsed.Claims.(*resetTokenCustomClaims); ok && tokenParsed.Valid {
		return resetToken{Email: claims.Email}, nil
	}

	return resetToken{}, errors.New("reset token is invalid")
}

func (r resetTokenizer) stringify(token resetToken) (string, error) {
	key := []byte("secret")
	claims := resetTokenCustomClaims{
		token,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + 3600,
			IssuedAt:  time.Now().Unix(),
			Issuer:    "apiv3",
		},
	}

	resetT := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	return resetT.SignedString(key)
}
