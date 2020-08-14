-- +goose Up
-- +goose StatementBegin
ALTER TABLE "user"
    ALTER COLUMN phone SET NOT NULL,
    ALTER COLUMN status SET NOT NULL,
    DROP CONSTRAINT user_login_key cascade,
    ADD CONSTRAINT user_phone_pk
    		UNIQUE (phone);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
ALTER TABLE "user"
    ALTER COLUMN phone DROP NOT NULL,
    ALTER COLUMN status DROP NOT NULL,
    DROP CONSTRAINT user_phone_pk;
CREATE UNIQUE INDEX user_login_key
    on "user"(name);

-- +goose StatementEnd
