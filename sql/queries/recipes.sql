-- name: CreateRecipe :one
INSERT INTO recipes (id, created_at, updated_at, title, description, difficulty, ingredients, instructions, category_id, user_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING id, created_at, title;

-- name: GetRecipes :many
SELECT * FROM recipes;

-- name: GetRecipesWithLimit :many
SELECT * FROM recipes
ORDER BY created_at DESC
LIMIT $1;

-- name: GetRecipeByID :one
SELECT * FROM recipes
WHERE id = $1 LIMIT 1;

-- name: UpdateRecipe :one
UPDATE recipes
SET updated_at = $2, title = $3, description = $4, ingredients = $5, instructions = $6, category_id = $7
WHERE id = $1
RETURNING id, updated_at, title;