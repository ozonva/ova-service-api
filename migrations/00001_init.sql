-- +goose Up
-- +goose StatementBegin
CREATE TABLE Services
(
  ID UUID PRIMARY KEY,
  UserID BIGINT NOT NULL,
  Description VARCHAR(4000) NULL,
  ServiceName VARCHAR(1000) NULL,
  ServiceAddress VARCHAR(1000) NULL,
  WhenLocal TIMESTAMP NULL,
  WhenUTC TIMESTAMP NULL,
  CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  UpdatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE Services;
-- +goose StatementEnd
