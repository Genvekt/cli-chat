-- +goose Up
-- +goose StatementBegin
ALTER TABLE "user"
    ADD CONSTRAINT user_unique_name UNIQUE (name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "user"
    DROP CONSTRAINT IF EXISTS user_unique_name;
-- +goose StatementEnd
