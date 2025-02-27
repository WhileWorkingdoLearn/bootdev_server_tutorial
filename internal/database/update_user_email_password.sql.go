// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: update_user_email_password.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const updateEmailAndPassword = `-- name: UpdateEmailAndPassword :exec
UPDATE Users SET
updated_at =  NOW(),
email = $1,
hashed_password = $2
WHERE id = $3
`

type UpdateEmailAndPasswordParams struct {
	Email          string
	HashedPassword string
	ID             uuid.UUID
}

func (q *Queries) UpdateEmailAndPassword(ctx context.Context, arg UpdateEmailAndPasswordParams) error {
	_, err := q.db.ExecContext(ctx, updateEmailAndPassword, arg.Email, arg.HashedPassword, arg.ID)
	return err
}
