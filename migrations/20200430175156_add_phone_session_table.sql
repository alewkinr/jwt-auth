-- +goose Up
-- +goose StatementBegin
CREATE TABLE sessions (
  session_id uuid PRIMARY KEY,
  users_phone varchar(20) NOT NULL,
  verification_code int NOT NULL,
  is_verified bool NOT NULL DEFAULT false,
  expires_at timestamp with time zone NOT NULL,
  created_at timestamp with time zone NOT NULL
);

CREATE INDEX users_phone_idx
    ON sessions (users_phone);

CREATE INDEX expires_at_idx
    ON sessions (expires_at);


-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE sessions CASCADE;
DROP INDEX users_phone_idx, expires_at_idx CASCADE;
-- +goose StatementEnd
