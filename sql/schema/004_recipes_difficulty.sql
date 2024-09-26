-- +goose Up
ALTER TABLE recipes
  ADD COLUMN difficulty INTEGER DEFAULT 1;

-- +goose Down
ALTER TABLE recipes
  DROP COLUMN difficulty;