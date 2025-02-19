-- +goose Up
ALTER TABLE Users ADD COLUMN is_chirpy_red BOOLEAN NOT NULL DEFAULT false;


-- +goose Down
ALTER TABLE Users
DROP COLUMN is_chirpy_red;