-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS registrations(
    id              BIGSERIAL   PRIMARY KEY,
    tournament_id   BIGINT  NOT NULL,
    username        TEXT    NOT NULL,
    registered_at   TIMESTAMP   NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE registrations;
-- +goose StatementEnd
