-- +goose Up
CREATE TABLE feeds (
	id INTEGER PRIMARY KEY,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	name TEXT NOT NULL,
	url TEXT NOT NULL UNIQUE,
	user_id TEXT,
	CONSTRAINT fk_user_id
	FOREIGN KEY (user_id)
	REFERENCES users(id)
	ON DELETE CASCADE
);

-- +goose Down
DROP TABLE feeds;
