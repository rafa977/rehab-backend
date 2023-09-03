-- name: CreateAccount :one
INSERT INTO accounts (user_id, username)
VALUES ($1, $2)
RETURNING *;

-- name: GetAccount :one
SELECT *
FROM accounts
WHERE user_id = $1
LIMIT 1;


-- name: CreateRoles
INSERT INTO roles (id, title, created_at) VALUES (1,'administrator', NOW());
INSERT INTO roles (id, title, created_at) VALUES (2,'user', NOW());