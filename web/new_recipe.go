package web

import (
	"html/template"
	"net/http"

	"github.com/s-hammon/recipls/api"
	"github.com/s-hammon/recipls/internal/database"
)

const tmplNewRecipe = "new_recipe.html"

func (c *config) renderNewRecipeTemplate(w http.ResponseWriter, r *http.Request, user database.User) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, nil)
		return
	}

	categories, err := fetchRecord[[]api.Category](c.client, "/categories", nil)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
	}

	tmpl := getTemplate(tmplNewRecipe, template.FuncMap{
		"seq": seq,
	})
	data := struct {
		Categories []api.Category
	}{categories}

	if err := tmpl.Execute(w, data); err != nil {
		respondError(w, http.StatusInternalServerError, "couldn't render template")
	}
}

func seq(start, end int) []int {
	s := make([]int, end-start+1)
	for i := start; i <= end; i++ {
		s[i-start] = i
	}

	return s
}
