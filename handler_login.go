package main

import (
	"encoding/json"
	"net/http"

	"github.com/s-hammon/recipls/internal/auth"
)

// const maxExpire = time.Second * 60 * 60 * 24

func (a *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request")
		return
	}

	user, err := a.DB.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondError(w, http.StatusNotFound, "user not found")
		return
	}

	if err := auth.CheckHash(user.Password, params.Password); err != nil {
		respondError(w, http.StatusUnauthorized, "invalid password")
		return
	}

	respondJSON(w, http.StatusAccepted, "login successful")
}
