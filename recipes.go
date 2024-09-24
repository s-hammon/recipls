package main

import (
	"time"

	"github.com/google/uuid"
	"github.com/s-hammon/recipls/internal/database"
)

type Recipe struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Ingredients  string    `json:"ingredients"`
	Instructions string    `json:"instructions"`
	CategoryID   uuid.UUID `json:"category_id"`
	UserID       uuid.UUID `json:"user_id"`
}

func dbToRecipe(recipe database.Recipe) Recipe {
	return Recipe{
		ID:           recipe.ID.Bytes,
		CreatedAt:    recipe.CreatedAt.Time,
		UpdatedAt:    recipe.UpdatedAt.Time,
		Title:        recipe.Title,
		Description:  recipe.Description,
		Ingredients:  recipe.Ingredients,
		Instructions: recipe.Instructions,
		CategoryID:   recipe.CategoryID.Bytes,
		UserID:       recipe.UserID.Bytes,
	}
}
