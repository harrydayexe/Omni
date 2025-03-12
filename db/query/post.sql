-- name: FindPostByID :one
SELECT * FROM posts WHERE id = ?;

-- name: GetUserAndPostsByIDPaged :many
SELECT users.id, users.username, sqlc.embed(posts) FROM users 
LEFT JOIN posts ON users.id = posts.user_id
WHERE users.id = ? AND posts.created_at > sqlc.arg(created_after) 
ORDER BY posts.created_at DESC
LIMIT ?;

-- name: CreatePost :exec
INSERT INTO posts (id, user_id, created_at, title, description, markdown_url) VALUES (?, ?, ?, ?, ?, ?);

-- name: UpdatePost :exec
UPDATE posts SET title=?, description=?, markdown_url=? WHERE id=?;

-- name: DeletePost :exec
DELETE FROM posts WHERE id = ?;

-- name: GetPostsPaged :many
SELECT 
    users.username,
    sqlc.embed(posts),
    CEIL(COUNT(*) OVER() / 10.0) AS total_pages
FROM posts
JOIN users ON posts.user_id = users.id
ORDER BY posts.created_at DESC
LIMIT 10 OFFSET ?;
