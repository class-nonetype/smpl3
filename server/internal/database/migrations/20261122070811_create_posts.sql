-- +goose Up
CREATE TABLE posts (
    id         uuid PRIMARY KEY DEFAULT uuidv7(),
    author_id  uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    title      text NOT NULL,
    body       text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT posts_title_length CHECK (char_length(title) BETWEEN 1 AND 200),
    CONSTRAINT posts_body_length CHECK (char_length(body) <= 20000)
);

CREATE INDEX posts_author_id_idx ON posts (author_id);

-- +goose Down
DROP TABLE posts;