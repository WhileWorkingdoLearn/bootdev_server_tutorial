-- +goose Up
ALTER TABLE Users ADD COLUMN hashed_password TEXT NOT NULL;
ALTER TABLE Users ALTER COLUMN hashed_password SET DEFAULT 'unset';

-- +goose Down
ALTER TABLE Users
DROP COLUMN hashed_password;