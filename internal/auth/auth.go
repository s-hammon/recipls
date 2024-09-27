package auth

import (
	"errors"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

const (
	AccessTokenType = "Bearer"
	APIKeyTokenType = "ApiKey"
)

var (
	ErrMissingAuthHeader = errors.New("no auth header in request")
	ErrInvalidAuthHeader = errors.New("invalid auth header format")
)

func HashPassword(password string) (string, error) {
	pwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(pwd), nil
}

func CheckHash(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func GetToken(tokenType string, headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", ErrMissingAuthHeader
	}

	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) < 2 || splitAuth[0] != tokenType {
		return "", ErrInvalidAuthHeader
	}

	return splitAuth[1], nil
}
