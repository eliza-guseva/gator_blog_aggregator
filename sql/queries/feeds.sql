-- name: CreateFeed :one
INSERT INTO feeds (name, url, user_id)
VALUES (
    $1,
    $2,
    $3
)
RETURNING *;

-- name: GetFeed :one
SELECT * FROM feeds where feeds.url = $1;

-- name: GetFeeds :many
SELECT 
    feeds.name,
    feeds.url,
    users.name as user_name
FROM feeds
INNER JOIN users ON feeds.user_id = users.id;
