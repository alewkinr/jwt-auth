-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.ea_db_sync (
    tbl_title VARCHAR(50),
    last_sync_id INT NOT NULL DEFAULT 0,
    last_sync_ts TIMESTAMP
);

-- update user table so it fit 'psychologist' role
ALTER TABLE public.user
    ALTER role TYPE varchar(15);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE public.ea_db_sync cascade;
-- +goose StatementEnd
