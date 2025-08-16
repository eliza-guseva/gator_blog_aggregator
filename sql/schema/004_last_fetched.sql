-- +goose Up
ALTER TABLE feeds ADD COLUMN last_feteched TIMESTAMP;

-- +goose Down
ALTER TABLE feeds DROP COLUMN last_feteched;
