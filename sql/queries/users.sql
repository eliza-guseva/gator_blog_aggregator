-- name: CreateUser :one
INSERT INTO users (name)
VALUES (
    $1
)
RETURNING *;    

-- name: GetUser :one
SELECT * FROM users where users.name = $1;

-- name: TruncateUsers :exec
DELETE FROM users;

-- name: GetUsers :many
SELECT name FROM users;
