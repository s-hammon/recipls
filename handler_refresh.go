package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/s-hammon/recipls/internal/auth"
)

func (a *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	authToken, err := auth.GetToken("Bearer", r.Header)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "couldn't parse bearer token from auth header")
		return
	}

	refreshToken, err := a.DB.GetRefreshTokenByValue(r.Context(), authToken)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "couldn't validate refresh token")
		return
	}

	userID, err := uuid.FromBytes(refreshToken.UserID.Bytes[:])
	if err != nil {
		respondError(w, http.StatusUnauthorized, "couldn't validate refresh token")
	}
	token, err := auth.MakeJWT(userID.String(), a.jwtSecret, maxExpire)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "couldn't create JWT")
		return
	}

	respondJSON(w, http.StatusOK, response{
		Token: token,
	})
}
