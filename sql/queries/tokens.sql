-- name: CreateRefreshToken :exec
INSERT INTO tokens (user_id, value, expires_at)
VALUES ($1, $2, $3);