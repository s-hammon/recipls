package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/s-hammon/recipls/internal/auth"
	"github.com/s-hammon/recipls/internal/database"
)

func (a *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if params.Email == "" || params.Name == "" || params.Password == "" {
		respondError(w, http.StatusBadRequest, "please provide name, email, & password")
		return
	}

	pwd, err := auth.HashPassword(params.Password)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "couldn't create user")
		return
	}

	user, err := a.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuidToPgType(uuid.New()),
		CreatedAt: timeToPgType(time.Now().UTC()),
		UpdatedAt: timeToPgType(time.Now().UTC()),
		Name:      params.Name,
		Email:     params.Email,
		Password:  pwd,
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	type response struct {
		User
	}
	respondJSON(w, http.StatusCreated, response{
		User: dbToUser(user),
	})
}

func (a *apiConfig) handleGetUserByAPIKey(w http.ResponseWriter, r *http.Request, user database.User) {
	respondJSON(w, http.StatusOK, dbToUser(user))
}
