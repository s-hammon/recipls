package api

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/s-hammon/recipls/internal/auth"
)

func (c *config) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	authToken, err := auth.GetToken("Bearer", r.Header)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "couldn't parse bearer token from auth header")
		return
	}

	refreshToken, err := c.DB.GetRefreshTokenByValue(r.Context(), authToken)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "couldn't validate refresh token")
		return
	}
	if refreshToken.ExpiresAt.Time.Before(time.Now().UTC()) {
		respondError(w, http.StatusUnauthorized, "refresh token expired; please log in")
		return
	}

	userID := uuid.UUID(refreshToken.UserID.Bytes)

	token, err := auth.MakeJWT(userID.String(), c.jwtSecret, maxExpire)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "couldn't create JWT")
		return
	}

	respondJSON(w, http.StatusOK, response{
		Token: token,
	})
}
