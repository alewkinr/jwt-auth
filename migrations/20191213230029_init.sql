-- +goose Up
-- +goose StatementBegin
CREATE TABLE public.user (
   id serial PRIMARY KEY,
   login VARCHAR (255) UNIQUE NOT NULL,
   password VARCHAR (60) NOT NULL,
   role VARCHAR (10) NOT NULL,
   created_on TIMESTAMP NOT NULL,
   last_login TIMESTAMP
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE user;
-- +goose StatementEnd
