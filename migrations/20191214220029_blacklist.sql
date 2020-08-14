-- +goose Up
-- +goose StatementBegin
CREATE TABLE blacklist (
   id serial PRIMARY KEY,
   user_id int8,
   token VARCHAR (255) NOT NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE blacklist;
-- +goose StatementEnd
