package api

import (
	"time"

	"github.com/google/uuid"
	"github.com/s-hammon/recipls/internal/database"
)

type Category struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
}

func DBToCategory(category database.Category) Category {
	return Category{
		ID:        category.ID.Bytes,
		CreatedAt: category.CreatedAt.Time,
		UpdatedAt: category.UpdatedAt.Time,
		Name:      category.Name,
	}
}
