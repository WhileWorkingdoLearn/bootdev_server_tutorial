-- name: UpdateEmailAndPassword :exec
UPDATE Users SET
updated_at =  NOW(),
email = $1,
hashed_password = $2
WHERE id = $3;