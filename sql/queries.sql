-- name: CreateAccount :one
INSERT INTO accounts (user_id, username)
VALUES ($1, $2)
RETURNING *;

-- name: GetAccount :one
SELECT *
FROM accounts
WHERE user_id = $1
LIMIT 1;
