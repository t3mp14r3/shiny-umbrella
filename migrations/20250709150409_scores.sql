-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS scores(
    id              BIGSERIAL   PRIMARY KEY,
    tournament_id   BIGINT  NOT NULL,
    username        TEXT    NOT NULL,
    score           INT     NOT NULL,
    placed_at       TIMESTAMP   NOT NULL DEFAULT NOW()
);

ALTER TABLE scores
      ADD CONSTRAINT fk_scores_tournament_id FOREIGN KEY (tournament_id) 
          REFERENCES tournaments (id);

ALTER TABLE scores
      ADD CONSTRAINT fk_scores_username FOREIGN KEY (username) 
          REFERENCES users (username);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE scores;
-- +goose StatementEnd
