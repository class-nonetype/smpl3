-- +goose Up
CREATE TABLE users (
    id              uuid PRIMARY KEY DEFAULT uuidv7(),
    username        text NOT NULL UNIQUE,
    email           text NOT NULL,
    password        text NOT NULL,
    created_at      timestamptz NOT NULL DEFAULT now(),
    updated_at      timestamptz NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX users_email_key ON users (lower(email));

-- +goose Down
DROP TABLE users;
