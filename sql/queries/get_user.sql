-- name: GetUserByID :one
SELECT * FROM Users WHERE id = $1;