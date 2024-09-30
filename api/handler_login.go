package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/s-hammon/recipls/internal/auth"
	"github.com/s-hammon/recipls/internal/database"
)

const maxExpire = time.Second * 60 * 5

const (
	ErrLoginBody          = "body must contain values for 'email' and 'password'"
	ErrFetchUser          = "user not found"
	ErrInvalidPassword    = "invalid password for email"
	ErrCreateRefreshToken = "couldn't create refresh token"
	ErrWriteRefreshToken  = "could't write refresh token"
)

func (c *config) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondError(w, http.StatusBadRequest, ErrLoginBody)
		return
	}

	userDB, err := c.DB.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondError(w, http.StatusUnauthorized, ErrFetchUser)
		return
	}
	user := DBToUser(userDB)

	if err := auth.CheckHash(user.Password, params.Password); err != nil {
		respondError(w, http.StatusUnauthorized, ErrInvalidPassword)
		return
	}

	token, err := auth.MakeJWT(user.ID.String(), c.jwtSecret, maxExpire)
	if err != nil {
		respondError(w, http.StatusInternalServerError, ErrCreateJWT)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondError(w, http.StatusInternalServerError, ErrCreateRefreshToken)
		return
	}

	expiresAt := time.Now().UTC().Add(time.Hour * 24 * 14)
	if err = c.DB.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		UserID:    uuidToPgType(user.ID),
		Value:     refreshToken,
		ExpiresAt: timeToPgType(expiresAt),
	}); err != nil {
		respondError(w, http.StatusInternalServerError, ErrWriteRefreshToken)
		return
	}

	type response struct {
		ID           string `json:"id"`
		Status       string `json:"status"`
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
	respondJSON(w, http.StatusOK, response{
		ID:           user.ID.String(),
		Status:       "success",
		AccessToken:  token,
		RefreshToken: refreshToken,
	})
}
