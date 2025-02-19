-- name: DeleteChirpById :exec
DELETE FROM Chirps WHERE id = $1;