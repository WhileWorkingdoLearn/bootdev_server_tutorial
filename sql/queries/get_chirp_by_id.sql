-- name: GetChirpById :one
SELECT * FROM Chirps WHERE id = $1;