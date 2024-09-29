package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/s-hammon/recipls/internal/auth"
	"github.com/s-hammon/recipls/internal/database"
)

const maxExpire = time.Second * 60 * 5

func (c *config) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request")
		return
	}

	userDB, err := c.DB.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondError(w, http.StatusNotFound, "user not found")
		return
	}
	user := DBToUser(userDB)

	if err := auth.CheckHash(user.Password, params.Password); err != nil {
		respondError(w, http.StatusUnauthorized, "invalid password")
		return
	}

	token, err := auth.MakeJWT(user.ID.String(), c.jwtSecret, maxExpire)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "couldn't create JWT")
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "couldn't create refresh token")
		return
	}

	expiresAt := time.Now().UTC().Add(time.Hour * 24 * 14)
	if err = c.DB.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		UserID:    uuidToPgType(user.ID),
		Value:     refreshToken,
		ExpiresAt: timeToPgType(expiresAt),
	}); err != nil {
		respondError(w, http.StatusInternalServerError, "couldn't write refresh token")
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
