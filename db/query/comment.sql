-- name: FindCommentAndUserByID :one
SELECT users.id, users.username, sqlc.embed(comments) FROM comments 
INNER JOIN users 
ON comments.user_id = users.id 
WHERE comments.id = ?;

-- name: CreateComment :exec
INSERT INTO comments (id, post_id, user_id, content, created_at) VALUES (?, ?, ?, ?, ?);

-- name: UpdateComment :exec
UPDATE comments SET content=? WHERE id=?;

-- name: DeleteComment :exec
DELETE FROM comments WHERE id = ?;

-- name: FindCommentsAndUserByPostIDPaged :many
SELECT users.id, users.username, sqlc.embed(comments) FROM comments 
INNER JOIN users ON comments.user_id = users.id 
ORDER BY comments.created_at DESC
LIMIT 10 OFFSET ?;
