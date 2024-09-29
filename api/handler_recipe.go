package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/s-hammon/recipls/internal/database"
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
			respondError(w, http.StatusNotFound, "invalid category")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	difficulty, err := strconv.Atoi(params.Difficulty)
	if err != nil || (difficulty < 0 || difficulty > 5) {
		respondError(w, http.StatusBadRequest, "difficulty must be c string integer between 1 and 5")
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
		slog.Info("got user_id in request", "user_id", reqUserId)
		userID = reqUserId
	}

	var dbRecipe []database.Recipe
	var err error
	switch userID {
	case "":
		dbRecipe, err = c.DB.GetRecipesWithLimit(r.Context(), 100)
		if err != nil {
			if err == pgx.ErrNoRows {
				respondError(w, http.StatusNotFound, "no recipes found")
				return
			}
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
	default:
		id, err := uuid.Parse(userID)
		if err != nil {
			respondError(w, http.StatusBadRequest, "couldn't parse user UUID")
			return
		}
		dbRecipe, err = c.DB.GetRecipesByUser(r.Context(), uuidToPgType(id))
		if err != nil {
			if err == pgx.ErrNoRows {
				respondError(w, http.StatusNotFound, "no recipes found")
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
		respondError(w, http.StatusBadRequest, "invalid id parameter")
		return
	}

	dbRecipe, err := c.DB.GetRecipeByID(r.Context(), uuidToPgType(id))
	if err != nil {
		if err == pgx.ErrNoRows {
			fmt.Printf("recipe not found: %v\n", id)
			respondError(w, http.StatusNotFound, "recipe not found")
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
		respondError(w, http.StatusBadRequest, "invalid id parameter")
	}

	recipe, err := c.DB.GetRecipeByID(r.Context(), uuidToPgType(id))
	if err != nil {
		if err == pgx.ErrNoRows {
			respondError(w, http.StatusNotFound, "recipe not found")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if !belongsToUser(user, recipe) {
		respondError(w, http.StatusForbidden, "you do not have permission to update this recipe")
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
			respondError(w, http.StatusNotFound, "invalid category_id")
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
		respondError(w, http.StatusBadRequest, "invalid id parameter")
		return
	}

	recipe, err := c.DB.GetRecipeByID(r.Context(), uuidToPgType(id))
	if err != nil {
		if err == pgx.ErrNoRows {
			respondError(w, http.StatusNotFound, "recipe not found")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if !belongsToUser(user, recipe) {
		respondError(w, http.StatusForbidden, "you do not have permission to delete this recipe")
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
