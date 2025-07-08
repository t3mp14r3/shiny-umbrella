-- +goose Up
-- +goose StatementBegin
CREATE TABLE test(
    test    TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE test;
-- +goose StatementEnd
