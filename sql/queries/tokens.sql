-- name: CreateRefreshToken :exec
INSERT INTO tokens (user_id, value, expires_at)
VALUES ($1, $2, $3);

-- name: GetRefreshTokenByValue :one
SELECT * FROM tokens
WHERE value = $1;

-- name: DeleteRefreshTokenByValue :exec
DELETE FROM tokens
WHERE value = $1;