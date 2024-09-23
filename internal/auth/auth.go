package auth

import (
	"errors"
	"net/http"
	"strings"
)

const (
	AccessTokenType = "Bearer"
	APIKeyTokenType = "ApiKey"
)

var (
	ErrMissingAuthHeader = errors.New("no auth header in request")
	ErrInvalidAuthHeader = errors.New("invalid auth header format")
)

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
