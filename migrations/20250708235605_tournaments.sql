-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tournaments(
    id          BIGSERIAL   PRIMARY KEY,
    price       INT         NOT NULL,
    min_users   INT         NOT NULL,
    max_users   INT         NOT NULL,
    bets        INT         NOT NULL,
    starts_at   TIMESTAMP   NOT NULL,
    duration    INTERVAL    NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE tournaments;
-- +goose StatementEnd
