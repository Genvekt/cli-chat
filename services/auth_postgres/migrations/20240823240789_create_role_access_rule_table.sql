-- +goose Up
-- +goose StatementBegin
CREATE TABLE "role_access_rule"
(
    role     int  NOT NULL,
    endpoint text NOT NULL,
    PRIMARY KEY (role, endpoint)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "role_access_rule";
-- +goose StatementEnd
