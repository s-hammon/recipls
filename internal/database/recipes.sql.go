// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: recipes.sql

package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createRecipe = `-- name: CreateRecipe :one
INSERT INTO recipes (id, created_at, updated_at, title, description, ingredients, instructions, category_id, user_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING id, created_at, title
`

type CreateRecipeParams struct {
	ID           pgtype.UUID
	CreatedAt    pgtype.Timestamp
	UpdatedAt    pgtype.Timestamp
	Title        string
	Description  string
	Ingredients  string
	Instructions string
	CategoryID   pgtype.UUID
	UserID       pgtype.UUID
}

type CreateRecipeRow struct {
	ID        pgtype.UUID
	CreatedAt pgtype.Timestamp
	Title     string
}

func (q *Queries) CreateRecipe(ctx context.Context, arg CreateRecipeParams) (CreateRecipeRow, error) {
	row := q.db.QueryRow(ctx, createRecipe,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Title,
		arg.Description,
		arg.Ingredients,
		arg.Instructions,
		arg.CategoryID,
		arg.UserID,
	)
	var i CreateRecipeRow
	err := row.Scan(&i.ID, &i.CreatedAt, &i.Title)
	return i, err
}

const getRecipeByID = `-- name: GetRecipeByID :one
SELECT id, created_at, updated_at, title, description, ingredients, instructions, category_id, user_id FROM recipes
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetRecipeByID(ctx context.Context, id pgtype.UUID) (Recipe, error) {
	row := q.db.QueryRow(ctx, getRecipeByID, id)
	var i Recipe
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Title,
		&i.Description,
		&i.Ingredients,
		&i.Instructions,
		&i.CategoryID,
		&i.UserID,
	)
	return i, err
}

const getRecipes = `-- name: GetRecipes :many
SELECT id, created_at, updated_at, title, description, ingredients, instructions, category_id, user_id FROM recipes
`

func (q *Queries) GetRecipes(ctx context.Context) ([]Recipe, error) {
	rows, err := q.db.Query(ctx, getRecipes)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Recipe
	for rows.Next() {
		var i Recipe
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Title,
			&i.Description,
			&i.Ingredients,
			&i.Instructions,
			&i.CategoryID,
			&i.UserID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getRecipesWithLimit = `-- name: GetRecipesWithLimit :many
SELECT id, created_at, updated_at, title, description, ingredients, instructions, category_id, user_id FROM recipes
LIMIT $1
`

func (q *Queries) GetRecipesWithLimit(ctx context.Context, limit int32) ([]Recipe, error) {
	rows, err := q.db.Query(ctx, getRecipesWithLimit, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Recipe
	for rows.Next() {
		var i Recipe
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Title,
			&i.Description,
			&i.Ingredients,
			&i.Instructions,
			&i.CategoryID,
			&i.UserID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateRecipe = `-- name: UpdateRecipe :one
UPDATE recipes
SET updated_at = $2, title = $3, description = $4, ingredients = $5, instructions = $6, category_id = $7
WHERE id = $1
RETURNING id, updated_at, title
`

type UpdateRecipeParams struct {
	ID           pgtype.UUID
	UpdatedAt    pgtype.Timestamp
	Title        string
	Description  string
	Ingredients  string
	Instructions string
	CategoryID   pgtype.UUID
}

type UpdateRecipeRow struct {
	ID        pgtype.UUID
	UpdatedAt pgtype.Timestamp
	Title     string
}

func (q *Queries) UpdateRecipe(ctx context.Context, arg UpdateRecipeParams) (UpdateRecipeRow, error) {
	row := q.db.QueryRow(ctx, updateRecipe,
		arg.ID,
		arg.UpdatedAt,
		arg.Title,
		arg.Description,
		arg.Ingredients,
		arg.Instructions,
		arg.CategoryID,
	)
	var i UpdateRecipeRow
	err := row.Scan(&i.ID, &i.UpdatedAt, &i.Title)
	return i, err
}
