package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/s-hammon/recipls/internal/database"
)

const maxLimit = 1000

type UserForMetrics struct {
	CreatedAt        time.Time `json:"created_at"`
	Name             string    `json:"name"`
	RecipesPublished int       `json:"recipes_published"`
}

type RecipeForMetrics struct {
	CreatedAt  time.Time `json:"created_at"`
	Title      string    `json:"title"`
	Difficulty int       `json:"difficulty"`
	Steps      int       `json:"steps"`
	Category   string    `json:"category"`
}

func (a *apiConfig) handlerGetMetrics(w http.ResponseWriter, r *http.Request, user database.User) {
	type response struct {
		Users   []UserForMetrics
		Recipes []RecipeForMetrics
	}

	limit := 100
	reqLimit := r.URL.Query().Get("limit")
	if reqLimit != "" {
		intLimit, err := strconv.Atoi(reqLimit)
		if err != nil || intLimit < 1 {
			respondError(w, http.StatusBadRequest, "limit must be a positive, non-zero integer")
			return
		}
		if intLimit > maxLimit {
			respondError(w, http.StatusBadRequest, "limit must be no greater than 1000")
			return
		}
		limit = intLimit
	}

	usersDB, err := a.DB.GetUsersWithLimit(r.Context(), int32(limit))
	if err != nil {
		respondError(w, http.StatusInternalServerError, "error fetching users")
		return
	}

	recipesDB, err := a.DB.GetRecipesWithLimit(r.Context(), int32(limit))
	if err != nil {
		respondError(w, http.StatusInternalServerError, "error fetching recipes")
		return
	}

	users := []UserForMetrics{}
	if len(usersDB) != 0 {
		for _, u := range usersDB {
			userRecipes, err := a.DB.GetRecipesByUser(r.Context(), u.ID)
			if err != nil {
				respondError(w, http.StatusInternalServerError, "error fetching user's recipes")
				return
			}
			user := dbToUser(u)
			users = append(users, user.toMetrics(len(userRecipes)))
		}
	}

	recipes := []RecipeForMetrics{}
	if len(recipesDB) != 0 {
		for _, p := range recipesDB {
			category, err := a.DB.GetCategoryByID(r.Context(), p.CategoryID)
			if err != nil {
				respondError(w, http.StatusInternalServerError, "error fetching recipe's category")
				return
			}
			recipe := dbToRecipe(p)
			recipes = append(recipes, recipe.toMetrics(category.Name))
		}
	}

	respondJSON(w, http.StatusOK, response{
		Users:   users,
		Recipes: recipes,
	})
}
