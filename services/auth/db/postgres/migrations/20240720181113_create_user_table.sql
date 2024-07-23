-- +goose Up
-- +goose StatementBegin
CREATE TABLE "user" (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    name text,
    email text,
    role int,
    created_at timestamp,
    updated_at timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "user";
-- +goose StatementEnd
