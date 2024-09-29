package api

import (
	"net/http"
)

func (c *config) handlerGetCategories(w http.ResponseWriter, r *http.Request) {
	categoriesDB, err := c.DB.GetCategories(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, "couldn't fetch categories")
		return
	}

	var categories []Category
	for _, i := range categoriesDB {
		categories = append(categories, DBToCategory(i))
	}

	respondJSON(w, http.StatusOK, categories)
}
