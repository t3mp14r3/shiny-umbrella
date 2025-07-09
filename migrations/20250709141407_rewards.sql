-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS rewards(
    id              BIGSERIAL   PRIMARY KEY,
    tournament_id   BIGINT  NOT NULL,
    place           INT     NOT NULL,
    prize           INT     NOT NULL
);

ALTER TABLE rewards
      ADD CONSTRAINT fk_rewards_tournament_id FOREIGN KEY (tournament_id) 
          REFERENCES tournaments (id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE rewards;
-- +goose StatementEnd
