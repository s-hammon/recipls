package api

import (
	"time"

	"github.com/google/uuid"
	"github.com/s-hammon/recipls/internal/database"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	ApiKey    string    `json:"api_key"`
}

func DBToUser(user database.User) User {
	return User{
		ID:        user.ID.Bytes,
		CreatedAt: user.CreatedAt.Time,
		UpdatedAt: user.UpdatedAt.Time,
		Name:      user.Name,
		Email:     user.Email,
		Password:  user.Password,
		ApiKey:    user.ApiKey,
	}
}

func (u *User) toMetrics(recipesPublished int) UserForMetrics {
	return UserForMetrics{
		CreatedAt:        u.CreatedAt,
		Name:             u.Name,
		RecipesPublished: recipesPublished,
	}
}
