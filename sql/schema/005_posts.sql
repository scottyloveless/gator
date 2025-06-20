-- +goose Up
CREATE TABLE posts (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    title TEXT NOT NULL,
    url TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL,
    published_at TIMESTAMP NOT NULL,
    feed_id INTEGER NOT NULL,
    CONSTRAINT fk_feed_id
    FOREIGN KEY (feed_id)
    REFERENCES feeds (id)
    ON DELETE CASCADE
);


-- +goose Down
DROP TABLE posts;
