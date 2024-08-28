-- +goose Up
-- +goose StatementBegin
ALTER TABLE "user"
    ADD password_hash TEXT NOT NULL DEFAULT '';

DELETE FROM "user"
    WHERE password_hash == '';

ALTER TABLE "user"
    ALTER COLUMN password_hash DROP DEFAULT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "user"
    DROP COLUMN IF EXISTS password_hash;
-- +goose StatementEnd
