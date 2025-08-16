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



-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows (user_id, feed_id)
    VALUES (
        $1,
        $2
    )
    RETURNING *
)
SELECT
    inserted_feed_follow.*,
    users.name as user_name,
    feeds.name as feed_name
FROM inserted_feed_follow
JOIN users on inserted_feed_follow.user_id = users.id
JOIN feeds on inserted_feed_follow.feed_id = feeds.id;

-- name: GetFeedFollowsForUser :many
SELECT 
    users.name as user_name,
    feeds.name as feed_name
FROM feed_follows
JOIN users on feed_follows.user_id = users.id
JOIN feeds on feed_follows.feed_id = feeds.id
WHERE feed_follows.user_id = $1;


-- name: UnfollowFeed :exec
DELETE FROM feed_follows
WHERE feed_follows.user_id = $1 AND feed_follows.feed_id = $2;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET last_feteched = now(), updated_at = now()
WHERE feeds.id = $1;

-- name: GetNextFeedToFetch :one
SELECT 
    *
FROM feeds
ORDER by last_feteched NULLS FIRST
LIMIT 1;

