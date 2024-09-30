package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/s-hammon/recipls/internal/database"
)

const (
	ErrInvalidCategory       = "invalid category"
	ErrInvalidDifficulty     = "difficuly must be a string integer between 1 and 5"
	ErrFetchRecipe           = "couldn't fine recipe"
	ErrFetchRecipes          = "no recipes found"
	ErrParseUserID           = "couldn't parse user ID"
	ErrRequestID             = "invalid id parameter"
	ErrForbiddenEditRecipe   = "you do not have permission to edit this recipe"
	ErrForbiddenDeleteRecipe = "you do not have permission to delete this recipe"
)

func (c *config) handlerCreateRecipe(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Title        string `json:"title"`
		Description  string `json:"description"`
		Difficulty   string `json:"difficulty"`
		Ingredients  string `json:"ingredients"`
		Instructions string `json:"instructions"`
		Category     string `json:"category"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondError(w, http.StatusBadRequest, fmt.Sprintf("invalid request: %v", err))
		return
	}

	category, err := c.DB.GetCategoryByName(r.Context(), params.Category)
	if err != nil {
		if err == pgx.ErrNoRows {
			respondError(w, http.StatusNotFound, ErrInvalidCategory)
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	difficulty, err := strconv.Atoi(params.Difficulty)
	if err != nil || (difficulty < 0 || difficulty > 5) {
		respondError(w, http.StatusBadRequest, ErrInvalidDifficulty)
		return
	}

	recipe, err := c.DB.CreateRecipe(r.Context(), database.CreateRecipeParams{
		ID:           uuidToPgType(uuid.New()),
		CreatedAt:    timeToPgType(time.Now().UTC()),
		UpdatedAt:    timeToPgType(time.Now().UTC()),
		Title:        params.Title,
		Description:  params.Description,
		Difficulty:   intToPgType(difficulty),
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

func (c *config) handlerGetRecipes(w http.ResponseWriter, r *http.Request) {
	userID := ""
	reqUserId := r.URL.Query().Get("user_id")
	if reqUserId != "" {
		userID = reqUserId
	}

	var dbRecipe []database.Recipe
	var err error
	switch userID {
	case "":
		dbRecipe, err = c.DB.GetRecipesWithLimit(r.Context(), 100)
		if err != nil {
			if err == pgx.ErrNoRows {
				respondError(w, http.StatusNotFound, ErrFetchRecipes)
				return
			}
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
	default:
		id, err := uuid.Parse(userID)
		if err != nil {
			respondError(w, http.StatusBadRequest, ErrParseUserID)
			return
		}
		dbRecipe, err = c.DB.GetRecipesByUser(r.Context(), uuidToPgType(id))
		if err != nil {
			if err == pgx.ErrNoRows {
				respondError(w, http.StatusNotFound, ErrFetchRecipes)
				return
			}
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	var recipes []Recipe
	for _, r := range dbRecipe {
		recipes = append(recipes, DBToRecipe(r))
	}

	respondJSON(w, http.StatusOK, recipes)
}

func (c *config) handlerGetRecipeByID(w http.ResponseWriter, r *http.Request) {
	id, err := getRequestID(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, ErrRequestID)
		return
	}

	dbRecipe, err := c.DB.GetRecipeByID(r.Context(), uuidToPgType(id))
	if err != nil {
		if err == pgx.ErrNoRows {
			respondError(w, http.StatusNotFound, ErrFetchRecipe)
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	recipe := DBToRecipe(dbRecipe)
	respondJSON(w, http.StatusOK, recipe)
}

func (c *config) handlerUpdateRecipe(w http.ResponseWriter, r *http.Request, user database.User) {
	id, err := getRequestID(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, ErrRequestID)
	}

	recipe, err := c.DB.GetRecipeByID(r.Context(), uuidToPgType(id))
	if err != nil {
		if err == pgx.ErrNoRows {
			respondError(w, http.StatusNotFound, ErrFetchRecipe)
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if !belongsToUser(user, recipe) {
		respondError(w, http.StatusForbidden, ErrForbiddenEditRecipe)
		return
	}

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

	category, err := c.DB.GetCategoryByID(r.Context(), uuidToPgType(params.CategoryID))
	if err != nil {
		if err == pgx.ErrNoRows {
			respondError(w, http.StatusNotFound, ErrInvalidCategory)
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if _, err = c.DB.UpdateRecipe(r.Context(), database.UpdateRecipeParams{
		ID:           uuidToPgType(id),
		UpdatedAt:    timeToPgType(time.Now().UTC()),
		Title:        params.Title,
		Description:  params.Description,
		Difficulty:   intToPgType(params.Difficulty),
		Ingredients:  params.Ingredients,
		Instructions: params.Instructions,
		CategoryID:   category.ID,
	}); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusNoContent, nil)
}

func (c *config) handlerDeleteRecipe(w http.ResponseWriter, r *http.Request, user database.User) {
	id, err := getRequestID(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, ErrRequestID)
		return
	}

	recipe, err := c.DB.GetRecipeByID(r.Context(), uuidToPgType(id))
	if err != nil {
		if err == pgx.ErrNoRows {
			respondError(w, http.StatusNotFound, ErrFetchRecipe)
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if !belongsToUser(user, recipe) {
		respondError(w, http.StatusForbidden, ErrForbiddenDeleteRecipe)
		return
	}

	if err := c.DB.DeleteRecipe(r.Context(), uuidToPgType(id)); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusNoContent, nil)
}

func belongsToUser(user database.User, recipe database.Recipe) bool {
	return user.ID == recipe.UserID
}
