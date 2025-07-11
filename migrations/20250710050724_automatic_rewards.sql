-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS automatic_rewards(
    id              BIGSERIAL   PRIMARY KEY,
    automatic_id    BIGINT  NOT NULL,
    place           INT     NOT NULL,
    prize           INT     NOT NULL
);

ALTER TABLE automatic_rewards
      ADD CONSTRAINT fk_automatic_rewards_automatic_id FOREIGN KEY (automatic_id) 
          REFERENCES automatic (id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE automatic_rewards;
-- +goose StatementEnd
