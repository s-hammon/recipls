package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/s-hammon/recipls/internal/database"
)

func (a *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	userUUID := uuid.New().String()
	id := pgtype.UUID{}
	if err := id.Scan(userUUID); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	nowDT := time.Now().UTC()
	createdAt := pgtype.Timestamp{}
	if err := createdAt.Scan(nowDT); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	updatedAt := pgtype.Timestamp{}
	if err := updatedAt.Scan(nowDT); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	userParams := database.CreateUserParams{
		ID:        id,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		Name:      params.Name,
	}

	fmt.Printf("creating user with id: %b\ncreated_at: %v\n", userParams.ID.Bytes, userParams.CreatedAt.Time)

	user, err := a.DB.CreateUser(r.Context(), userParams)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, user)
}

func (a *apiConfig) handleGetUserByAPIKey(w http.ResponseWriter, r *http.Request, user database.User) {
	respondJSON(w, http.StatusOK, dbToUser(user))
}
