package main

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	ApiKey    string    `json:"api_key"`
}
