-- name: DeleteUserById :exec
DELETE FROM Users WHERE id = $1;