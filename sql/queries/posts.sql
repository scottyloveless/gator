-- name: CreatePost :one
INSERT INTO posts (created_at, updated_at, title, url, description, published_at, feed_id)
VALUES (
    NOW(),
    NOW(),
    $1,
    $2,
    $3,
    $4,
    $5
    )
RETURNING *;

-- name: GetPostsForUser :many
SELECT
    posts.*
FROM users
INNER JOIN feed_follows
    ON users.id = feed_follows.user_id
INNER JOIN feeds
    ON feed_follows.feed_id = feeds.id
INNER JOIN posts
    ON feeds.id = posts.feed_id
WHERE users.name = $1
ORDER BY posts.published_at DESC
LIMIT $2;
