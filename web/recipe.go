package web

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/s-hammon/recipls/api"
)

const tmplRecipe = "recipe.html"

func (c *config) renderRecipeTemplate(w http.ResponseWriter, r *http.Request) {
	id, err := getRequestID(r)
	if err != nil {
		slog.Error("couldn't find message", "error", "recipe not found 😔")
		respondError(w, http.StatusNotFound, "recipe not found 😔")
		return
	}

	recipeEndpoint := fmt.Sprintf("/recipes/%s", id.String())
	recipe, err := fetchRecord[api.Recipe](c.client, recipeEndpoint, nil)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	userEndpoint := fmt.Sprintf("/users/%s", recipe.UserID.String())
	user, err := fetchRecord[api.User](c.client, userEndpoint, nil)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	tmpl := getTemplate(tmplRecipe, template.FuncMap{"splitLines": splitLines})
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
