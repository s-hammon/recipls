package web

import (
	"html/template"
	"net/http"

	"github.com/s-hammon/recipls/api"
	"github.com/s-hammon/recipls/internal/database"
)

func (c *config) renderRecipeTemplate(w http.ResponseWriter, r *http.Request) {
	id, err := getRequestID(r)
	if err != nil {
		respondError(w, http.StatusNotFound, "recipe not found 😔")
		return
	}

	recipeDB, err := c.DB.GetRecipeByID(r.Context(), uuidToPgType(id))
	if err != nil {
		respondError(w, http.StatusInternalServerError, "error getting recipe")
		return
	}
	recipe := api.DBToRecipe(recipeDB)

	userDB, err := c.DB.GetUserByID(r.Context(), uuidToPgType(recipe.UserID))
	if err != nil {
		respondError(w, http.StatusInternalServerError, "error getting user")
		return
	}
	user := api.DBToUser(userDB)

	tmpl := getTemplate("recipe.html", template.FuncMap{"splitLines": splitLines})
	data := struct {
		Recipe api.Recipe
		User   api.User
	}{
		Recipe: recipe,
		User:   user,
	}

	if err := tmpl.Execute(w, data); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
	}
}

func (c *config) renderNewRecipeTemplate(w http.ResponseWriter, r *http.Request, user database.User) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, nil)
		return
	}

	categories, err := c.DB.GetCategories(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, "couldn't fetch categories")
		return
	}

	tmpl := getTemplate("new_recipe.html", template.FuncMap{
		"seq": seq,
	})
	data := struct {
		Categories []database.Category
	}{categories}

	if err := tmpl.Execute(w, data); err != nil {
		respondError(w, http.StatusInternalServerError, "couldn't render template")
	}
}

func (c *config) renderLoginTemplate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, nil)
		return
	}

	tmpl := getTemplate("login.html", nil)
	if err := tmpl.Execute(w, nil); err != nil {
		respondError(w, http.StatusNotFound, "couldn't find login page")
	}
}

func seq(start, end int) []int {
	s := make([]int, end-start+1)
	for i := start; i <= end; i++ {
		s[i-start] = i
	}

	return s
}
