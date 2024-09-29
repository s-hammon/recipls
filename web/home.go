package web

import (
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/s-hammon/recipls/api"
	"github.com/s-hammon/recipls/internal/database"
)

const tmplHome = "home.html"

func (c *config) renderHomeTemplate(w http.ResponseWriter, r *http.Request, user database.User) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	userID := uuid.UUID(user.ID.Bytes)
	queryParams := queryParams{"user_id": userID.String()}

	recipes, err := fetchRecord[[]api.Recipe](c.client, "/recipes", queryParams)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
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
		slog.Error("couldn't serve template", "tempalte", tmplHome, "error", err)
		respondError(w, http.StatusInternalServerError, nil)
	}
}
