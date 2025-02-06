-- name: FindPostByID :one
SELECT * FROM posts WHERE id = ?;

-- name: GetUserAndPostsByIDPaged :many
SELECT users.id, users.username, sqlc.embed(posts) FROM users 
LEFT JOIN posts ON users.id = posts.user_id
WHERE users.id = ? AND posts.created_at > sqlc.arg(created_after) 
ORDER BY posts.created_at ASC
LIMIT ?;

-- name: CreatePost :exec
INSERT INTO posts (id, user_id, created_at, title, description, markdown_url) VALUES (?, ?, ?, ?, ?, ?);

-- name: UpdatePost :exec
UPDATE posts SET title=?, description=?, markdown_url=? WHERE id=?;

-- name: DeletePost :exec
DELETE FROM posts WHERE id = ?;
