-- +goose Up
-- +goose StatementBegin
ALTER TABLE public.user RENAME COLUMN login TO email;
ALTER TABLE public.user
  ADD name VARCHAR (1024) NULL,
  ADD phone VARCHAR(16) NULL,
  ADD status VARCHAR(64) NULL;

--UPDATE public.user SET name = 'Админ Админович', phone = '+79779814962', status='active' WHERE id = 1;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE public.user RENAME COLUMN email TO login;
ALTER TABLE public.user
  DROP COLUMN name,
  DROP COLUMN phone,
  DROP COLUMN status;
-- +goose StatementEnd
