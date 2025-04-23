-- name: GetPostsPaged :many
SELECT 
    users.username,
    sqlc.embed(posts),
    CEIL(COUNT(*) OVER() / 10.0) AS total_pages
FROM posts
JOIN users ON posts.user_id = users.id
ORDER BY posts.created_at DESC
LIMIT 10 OFFSET ?;

