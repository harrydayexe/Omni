-- name: FindCommentAndUserByID :one
SELECT sqlc.embed(Comments), sqlc.embed(Users) FROM Comments 
INNER JOIN Users 
ON Comments.user_id = Users.id 
WHERE Comments.id = ?;

-- name: CreateComment :exec
INSERT INTO Comments (id, post_id, user_id, content, created_at) VALUES (?, ?, ?, ?, ?);

-- name: UpdateComment :exec
UPDATE Comments SET content=? WHERE id=?;

-- name: DeleteComment :exec
DELETE FROM Comments WHERE id = ?;

-- name: FindCommentsAndUserByPostIDPaged :many
SELECT sqlc.embed(Comments), sqlc.embed(Users) FROM Comments 
INNER JOIN Users 
ON Comments.user_id = Users.id 
WHERE Comments.post_id = ? AND Comments.created_at > sqlc.arg(created_after) 
ORDER BY Comments.created_at ASC
LIMIT ?;
