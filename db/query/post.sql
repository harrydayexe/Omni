-- name: FindPostByID :one
SELECT * FROM Posts WHERE id = ?;

-- name: GetUserAndPostsByIDPaged :one
SELECT sqlc.embed(Users), sqlc.embed(Posts) FROM Users 
LEFT JOIN Posts ON Users.id = Posts.user_id
WHERE Users.id = ? AND Posts.created_at > sqlc.arg(created_after) 
ORDER BY Posts.created_at ASC
LIMIT ?;

-- name: CreatePost :exec
INSERT INTO Posts (id, user_id, created_at, title, description, markdown_url) VALUES (?, ?, ?, ?, ?, ?);

-- name: UpdatePost :exec
UPDATE Posts SET title=?, description=?, markdown_url=? WHERE id=?;

-- name: DeletePost :exec
DELETE FROM Posts WHERE id = ?;
