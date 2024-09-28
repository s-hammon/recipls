package main

import (
	"html/template"
	"log/slog"
	"net/http"

	"github.com/s-hammon/recipls/internal/database"
)

func (a *apiConfig) renderHomeTemplate(w http.ResponseWriter, r *http.Request, user database.User) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var recipes []Recipe
	recipesDB, err := a.DB.GetRecipesByUser(r.Context(), user.ID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "error fetching recipes")
		return
	}
	if len(recipesDB) == 0 {
		slog.Warn("WARN: no recipes found", "user_id", user.ID)
	}

	for _, r := range recipesDB {
		recipes = append(recipes, dbToRecipe(r))
	}

	tmpl := getTemplate("home.html", nil)
	data := struct {
		User    User
		Recipes []Recipe
	}{
		dbToUser(user),
		recipes,
	}

	if err := tmpl.Execute(w, data); err != nil {
		respondError(w, http.StatusInternalServerError, nil)
	}
}

func (a *apiConfig) renderRecipeTemplate(w http.ResponseWriter, r *http.Request) {
	id, err := getRequestID(r)
	if err != nil {
		respondError(w, http.StatusNotFound, "recipe not found 😔")
		return
	}

	recipeDB, err := a.DB.GetRecipeByID(r.Context(), uuidToPgType(id))
	if err != nil {
		respondError(w, http.StatusInternalServerError, "error getting recipe")
		return
	}
	recipe := dbToRecipe(recipeDB)

	userDB, err := a.DB.GetUserByID(r.Context(), uuidToPgType(recipe.UserID))
	if err != nil {
		respondError(w, http.StatusInternalServerError, "error getting user")
		return
	}
	user := dbToUser(userDB)

	tmpl := getTemplate("recipe.html", template.FuncMap{"splitLines": splitLines})
	data := struct {
		Recipe Recipe
		User   User
	}{
		Recipe: recipe,
		User:   user,
	}

	if err := tmpl.Execute(w, data); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
	}
}

func (a *apiConfig) renderLoginTemplate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, nil)
	}

	tmpl := getTemplate("login.html", nil)
	if err := tmpl.Execute(w, nil); err != nil {
		respondError(w, http.StatusNotFound, "couldn't find login page")
	}
}
