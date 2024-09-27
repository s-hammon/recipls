-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name, email, password, api_key)
VALUES ($1, $2, $3, $4, $5, $6, encode(sha256(random()::text::bytea), 'hex'))
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByAPIKey :one
SELECT * FROM users
WHERE api_key = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: UpdateUser :one
UPDATE users
SET updated_at = $2, name = $3, email = $4, password = $5
WHERE id = $1
RETURNING id, updated_at, name;