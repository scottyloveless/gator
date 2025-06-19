-- name: CreateFeed :one
INSERT INTO feeds (created_at, updated_at, name, url, user_id)
VALUES (
	$1,
	$2,
	$3,
	$4,
	$5
	)
RETURNING *;

-- name: ListFeeds :many
SELECT feeds.*, users.name AS username
FROM feeds
LEFT JOIN users
ON feeds.user_id = users.id;

-- name: GetFeedByURL :one
SELECT *
FROM feeds
WHERE url = $1;
