-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
	INSERT INTO feed_follows (created_at, updated_at, user_id, feed_id)
	VALUES (
		$1,
		$2,
		$3,
		$4
		)
	RETURNING *
)

SELECT
    inserted_feed_follow.*,
    feeds.name AS feed_name,
    users.name AS user_name
FROM inserted_feed_follow
INNER JOIN users
    ON inserted_feed_follow.user_id = users.id
INNER JOIN feeds
    ON inserted_feed_follow.feed_id = feeds.id;

-- name: GetFeedFollowsForUser :many
SELECT
    feed_follows.user_id,
    feed_follows.feed_id,
    users.name,
    feeds.url AS feed_url,
    feeds.name AS feed_name
FROM feed_follows
INNER JOIN users
    ON feed_follows.user_id = users.id
INNER JOIN feeds
    ON feed_follows.feed_id = feeds.id
WHERE feed_follows.user_id = $1;

-- name: Unfollow :exec
DELETE FROM feed_follows
WHERE user_id = (SELECT id FROM users WHERE users.name = $1)
AND feed_id = (SELECT id FROM feeds WHERE feeds.url = $2);
