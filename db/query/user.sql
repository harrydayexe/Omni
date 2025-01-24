-- name: GetUserByID :one
SELECT * FROM users WHERE id = ?;

-- name: CreateUser :exec
INSERT INTO users (id, username) VALUES (?, ?);

-- name: UpdateUser :exec
UPDATE users SET username = ? WHERE id = ?;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = ?;
