-- name: GetUserByID :one
SELECT * FROM Users WHERE id = ?;

-- name: CreateUser :exec
INSERT INTO Users (id, username) VALUES (?, ?);

-- name: UpdateUser :exec
UPDATE Users SET username = ? WHERE id = ?;

-- name: DeleteUser :exec
DELETE FROM Users WHERE id = ?;
