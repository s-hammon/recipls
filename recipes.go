package main

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/s-hammon/recipls/app"
	"github.com/s-hammon/recipls/internal/database"
)

type Recipe struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Difficulty   string    `json:"difficulty"`
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
		Difficulty:   getDifficultyString(int(recipe.Difficulty.Int32)),
		Ingredients:  recipe.Ingredients,
		Instructions: recipe.Instructions,
		CategoryID:   recipe.CategoryID.Bytes,
		UserID:       recipe.UserID.Bytes,
	}
}

func (a *apiConfig) UpdateRSS() error {
	recipesDB, err := a.DB.GetRecipesWithLimit(context.Background(), 100)
	if err != nil {
		return err
	}

	items := []app.Item{}
	for _, recipe := range recipesDB {
		r := dbToRecipe(recipe)
		items = append(items, app.Item{
			Title:       r.Title,
			Link:        xmlDomain + "/recipes/" + r.ID.String(),
			Description: r.Description,
			PubDate:     r.CreatedAt,
		})
	}

	return a.App.AddItems(items)
}

func (a *apiConfig) rssUpdateWorker(requestInterval time.Duration) {
	ticker := time.NewTicker(requestInterval)

	for ; ; <-ticker.C {
		if err := a.UpdateRSS(); err != nil {
			log.Println("couldn't update RSS feed: ", err)
			continue
		}

		log.Println("RSS feed updated")
	}
}
