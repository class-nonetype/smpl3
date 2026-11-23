-- name: CreatePost :one
INSERT INTO posts (author_id, title, body)
VALUES ($1, $2, $3)
RETURNING id;

-- name: GetPost :one
SELECT p.id, p.title, p.body, p.created_at, p.updated_at,
       u.id AS author_id, u.username AS author_username
FROM posts p
JOIN users u ON u.id = p.author_id
WHERE p.id = $1;

-- name: ListPosts :many
SELECT p.id, p.title, p.body, p.created_at, p.updated_at,
       u.id AS author_id, u.username AS author_username
FROM posts p
JOIN users u ON u.id = p.author_id
WHERE sqlc.narg(before)::uuid IS NULL OR p.id < sqlc.narg(before)::uuid
ORDER BY p.id DESC
LIMIT sqlc.arg(page_size);

-- name: DeletePost :execrows
DELETE FROM posts
WHERE id = $1 AND author_id = $2;