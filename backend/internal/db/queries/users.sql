-- name: CreateUser :one
INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1 LIMIT 1;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByUserID :one
SELECT * FROM users WHERE user_id = $1;

-- name: UpdateUsername :exec
UPDATE users
SET username = $1, updated_at = CURRENT_TIMESTAMP
WHERE id = $2;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;