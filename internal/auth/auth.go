package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const (
	AccessTokenType = "Bearer"
	APIKeyTokenType = "ApiKey"
)

var (
	ErrMissingAuthHeader = errors.New("no auth header in request")
	ErrInvalidAuthHeader = errors.New("invalid auth header format")
	ErrExpiredJWT        = errors.New("JWT has expired")
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

func MakeRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}

func MakeJWT(userID, tokenSecret string, expiresIn time.Duration) (string, error) {
	signKey := []byte(tokenSecret)
	claims := jwt.RegisteredClaims{
		Issuer:    "recipls",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(signKey)
}

func ValidateJWT(tokenString, tokenSecret string) (string, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(t *jwt.Token) (interface{}, error) { return []byte(tokenSecret), nil },
	)
	if err != nil {
		return "", err
	}

	userID, err := token.Claims.GetSubject()
	if err != nil {
		return "", err
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return "", err
	}
	if issuer != string("recipls") {
		return "", errors.New("invalid issuer")
	}

	expiry, err := token.Claims.GetExpirationTime()
	if err != nil {
		return "", err
	}
	if expiry.Time.Before(time.Now().UTC()) {
		return "", ErrExpiredJWT
	}

	return userID, nil
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
