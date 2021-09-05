-- +goose Up
-- +goose StatementBegin
CREATE TABLE Services
(
  id UUID PRIMARY KEY,
  user_id BIGINT NOT NULL,
  description VARCHAR(4000) NULL,
  service_name VARCHAR(1000) NULL,
  service_address VARCHAR(1000) NULL,
  when_local TIMESTAMP NULL,
  when_utc TIMESTAMP NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE Services;
-- +goose StatementEnd
