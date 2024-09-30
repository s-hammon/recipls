package api

import (
	"net/http"
)

const ErrFetchCategories = "couldn't fetch categories"

func (c *config) handlerGetCategories(w http.ResponseWriter, r *http.Request) {
	categoriesDB, err := c.DB.GetCategories(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, ErrFetchCategories)
		return
	}

	var categories []Category
	for _, i := range categoriesDB {
		categories = append(categories, DBToCategory(i))
	}

	respondJSON(w, http.StatusOK, categories)
}
