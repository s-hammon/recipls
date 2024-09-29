package web

import (
	"log/slog"
	"net/http"

	"github.com/s-hammon/recipls/api"
	"github.com/s-hammon/recipls/internal/database"
)

const tmplHome = "home.html"

func (c *config) renderHomeTemplate(w http.ResponseWriter, r *http.Request, user database.User) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var recipes []api.Recipe
	recipesDB, err := c.DB.GetRecipesByUser(r.Context(), user.ID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "error fetching recipes")
		return
	}
	if len(recipesDB) == 0 {
		slog.Warn("WARN: no recipes found", "user_id", user.ID)
	}

	for _, r := range recipesDB {
		recipes = append(recipes, api.DBToRecipe(r))
	}

	tmpl := getTemplate(tmplHome, nil)
	data := struct {
		User    api.User
		Recipes []api.Recipe
	}{
		api.DBToUser(user),
		recipes,
	}

	if err := tmpl.Execute(w, data); err != nil {
		respondError(w, http.StatusInternalServerError, nil)
	}
}
