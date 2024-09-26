// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package database

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Category struct {
	ID        pgtype.UUID
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
	Name      string
}

type Recipe struct {
	ID           pgtype.UUID
	CreatedAt    pgtype.Timestamp
	UpdatedAt    pgtype.Timestamp
	Title        string
	Description  string
	Ingredients  string
	Instructions string
	CategoryID   pgtype.UUID
	UserID       pgtype.UUID
	Difficulty   pgtype.Int4
}

type User struct {
	ID        pgtype.UUID
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
	Name      string
	ApiKey    string
}
