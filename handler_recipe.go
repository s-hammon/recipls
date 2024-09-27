package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/s-hammon/recipls/internal/database"
)

func (a *apiConfig) handlerCreateRecipe(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Title        string    `json:"title"`
		Description  string    `json:"description"`
		Difficulty   int       `json:"difficulty"`
		Ingredients  string    `json:"ingredients"`
		Instructions string    `json:"instructions"`
		CategoryID   uuid.UUID `json:"category_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondError(w, http.StatusBadRequest, fmt.Sprintf("invalid request: %v", err))
		return
	}

	category, err := a.DB.GetCategoryByID(r.Context(), uuidToPgType(params.CategoryID))
	if err != nil {
		if err == pgx.ErrNoRows {
			respondError(w, http.StatusNotFound, "invalid category_id")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	recipe, err := a.DB.CreateRecipe(r.Context(), database.CreateRecipeParams{
		ID:           uuidToPgType(uuid.New()),
		CreatedAt:    timeToPgType(time.Now().UTC()),
		UpdatedAt:    timeToPgType(time.Now().UTC()),
		Title:        params.Title,
		Description:  params.Description,
		Ingredients:  params.Ingredients,
		Instructions: params.Instructions,
		CategoryID:   category.ID,
		UserID:       user.ID,
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	type response struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		Title     string    `json:"title"`
	}
	resp := response{
		ID:        recipe.ID.Bytes,
		CreatedAt: recipe.CreatedAt.Time,
		Title:     recipe.Title,
	}

	respondJSON(w, http.StatusCreated, resp)
}

func (a *apiConfig) handlerGetRecipeByID(w http.ResponseWriter, r *http.Request) {
	id, err := getRequestID(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id parameter")
		return
	}

	dbRecipe, err := a.DB.GetRecipeByID(r.Context(), uuidToPgType(id))
	if err != nil {
		if err == pgx.ErrNoRows {
			fmt.Printf("recipe not found: %v\n", id)
			respondError(w, http.StatusNotFound, "recipe not found")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	recipe := dbToRecipe(dbRecipe)
	respondJSON(w, http.StatusOK, recipe)
}
