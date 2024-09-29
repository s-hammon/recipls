package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/s-hammon/recipls/internal/auth"
	"github.com/s-hammon/recipls/internal/database"
)

func (c *config) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
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

	user, err := c.DB.CreateUser(r.Context(), database.CreateUserParams{
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
		User: DBToUser(user),
	})
}

func (c *config) handleGetUserByID(w http.ResponseWriter, r *http.Request) {
	id, err := getRequestID(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id parameter")
		return
	}

	userDB, err := c.DB.GetUserByID(r.Context(), uuidToPgType(id))
	if err != nil {
		if err == pgx.ErrNoRows {
			respondError(w, http.StatusNotFound, "user not found")
			return
		}
	}

	user := DBToUser(userDB)
	respondJSON(w, http.StatusOK, user)
}

func (c *config) handlerGetUserByAPIKey(w http.ResponseWriter, r *http.Request, user database.User) {
	respondJSON(w, http.StatusOK, DBToUser(user))
}
