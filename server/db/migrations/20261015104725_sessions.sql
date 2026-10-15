-- +goose Up
CREATE TABLE sessions (
  id varchar(255) PRIMARY KEY NOT NULL,
  user_email varchar(255) NOT NULL,
  refresh_token varchar(512) NOT NULL,
  is_revoked bool NOT NULL DEFAULT false,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  expires_at TIMESTAMP
);

-- +goose Down
DROP TABLE sessions;
