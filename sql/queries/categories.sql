-- name: CreateCategory :one
INSERT INTO categories (id, created_at, updated_at, name)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetCategories :many
SELECT * FROM categories;

-- name: GetCategoryByID :one
SELECT * FROM categories
WHERE id = $1 LIMIT 1;

-- name: UpdateCategory :one
UPDATE categories
SET updated_at = $2, name = $3
WHERE id = $1
RETURNING *;

-- name: GetCategoryByName :one
SELECT * FROM categories
WHERE name = $1 LIMIT 1;

-- name: GetCategoryByNames :many
SELECT * FROM categories
WHERE name = ANY($1::uuid[]);

-- name: DeleteCategory :exec
DELETE FROM categories WHERE id = $1;