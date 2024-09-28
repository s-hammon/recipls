package main

import (
	"html/template"
	"net/http"
)

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
