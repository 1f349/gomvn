-- name: CountUsers :one
SELECT count(*)
FROM users;

-- name: IsAdmin :one
SELECT 1
FROM users
WHERE admin = 1
  AND token_hash = ?;

-- name: GetAllUsers :many
SELECT id, name, admin, created_at, updated_at
FROM users;

-- name: CreateUser :execlastid
INSERT INTO users (name, admin, token_hash, created_at, updated_at)
VALUES (?, ?, ?, ?, ?);
