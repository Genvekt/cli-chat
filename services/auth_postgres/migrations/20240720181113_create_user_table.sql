-- +goose Up
-- +goose StatementBegin
CREATE TABLE "user"
(
    id         BIGSERIAL NOT NULL PRIMARY KEY,
    name       text      NOT NULL,
    email      text      NOT NULL,
    role       int       NOT NULL,
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "user";
-- +goose StatementEnd
