-- name: UpdateUserStatus :exec
UPDATE Users SET
updated_at =  NOW(),
is_chirpy_red = $2
WHERE id = $1;