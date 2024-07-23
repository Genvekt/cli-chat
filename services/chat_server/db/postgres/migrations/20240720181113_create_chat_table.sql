-- +goose Up
-- +goose StatementBegin
CREATE TABLE "chat" (
    id BIGSERIAL NOT NULL PRIMARY KEY
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "chat";
-- +goose StatementEnd
