package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
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

	nowDT := time.Now().UTC()
	userUUID, err := uuid.NewV4()
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	userParams := database.CreateUserParams{
		ID:        pgtype.UUID{Bytes: [16]byte(userUUID.Bytes())},
		CreatedAt: pgtype.Timestamp{Time: nowDT},
		UpdatedAt: pgtype.Timestamp{Time: nowDT},
		Name:      params.Name,
	}

	user, err := a.DB.CreateUser(r.Context(), userParams)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, user)
}
