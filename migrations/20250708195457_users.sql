-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users(
    username    TEXT    NOT NULL PRIMARY KEY,
    balance     INT     NOT NULL DEFAULT 100
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
