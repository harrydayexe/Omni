-- name: GetUserByID :one
SELECT id, username FROM users WHERE id = ?;

-- name: GetUserByUsername :one
SELECT id FROM users WHERE username = ?;

-- name: GetPasswordByID :one
SELECT password FROM users WHERE id = ?;

-- name: CreateUser :exec
INSERT INTO users (id, username, password) VALUES (?, ?, ?);

-- name: UpdateUser :exec
UPDATE users SET username = ? WHERE id = ?;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = ?;
