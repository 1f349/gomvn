-- name: CountUsers :one
SELECT count(*)
FROM users;

-- name: IsAdmin :one
SELECT 1
FROM users
WHERE admin = 1
  AND token_hash = ?;

-- name: IsValid :one
SELECT 1
FROM users
WHERE token_hash = ?;

-- name: GetAllUsers :many
SELECT id, name
FROM users;

-- name: CreateUser :execlastid
INSERT INTO users (name, admin, token_hash)
VALUES (?, ?, ?);

-- name: RefreshUserToken :exec
UPDATE users
SET token_hash =?
WHERE id = ?
  AND token_hash = ?;

-- name: CheckUserDetails :one
SELECT 1
FROM users
WHERE name = ?
  AND token_hash = ?;
