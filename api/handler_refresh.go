package api

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/s-hammon/recipls/internal/auth"
)

const (
	ErrParseHeader          = "couldn't parse bearer token from auth header"
	ErrValidateRefreshToken = "couldn't validate refresh token"
	ErrExpiredRefreshToken  = "refresh token expired"
	ErrCreateJWT            = "couldn't create JWT"
)

func (c *config) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"refresh_token"`
	}

	authToken, err := auth.GetToken("Bearer", r.Header)
	if err != nil {
		respondError(w, http.StatusUnauthorized, ErrParseHeader)
		return
	}

	refreshToken, err := c.DB.GetRefreshTokenByValue(r.Context(), authToken)
	if err != nil {
		respondError(w, http.StatusUnauthorized, ErrValidateJWT)
		return
	}
	if refreshToken.ExpiresAt.Time.Before(time.Now().UTC()) {
		respondError(w, http.StatusUnauthorized, ErrExpiredRefreshToken)
		return
	}

	userID := uuid.UUID(refreshToken.UserID.Bytes)

	token, err := auth.MakeJWT(userID.String(), c.jwtSecret, maxExpire)
	if err != nil {
		respondError(w, http.StatusInternalServerError, ErrCreateJWT)
		return
	}

	respondJSON(w, http.StatusOK, response{
		Token: token,
	})
}
