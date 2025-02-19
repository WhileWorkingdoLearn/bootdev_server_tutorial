-- name: GetUserByEmail :one
SELECT * FROM Users WHERE email = $1;