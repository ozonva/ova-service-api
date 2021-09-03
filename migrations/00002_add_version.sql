-- +goose Up
-- +goose StatementBegin
ALTER TABLE services
  ADD COLUMN version INTEGER NOT NULL DEFAULT 1;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE services
  DROP COLUMN version;
-- +goose StatementEnd
