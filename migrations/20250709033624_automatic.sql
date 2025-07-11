-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS automatic(
    id          BIGSERIAL   PRIMARY KEY,
    price       INT         NOT NULL,
    min_users   INT         NOT NULL,
    max_users   INT         NOT NULL,
    bets        INT         NOT NULL,
    duration    INTERVAL    NOT NULL,
    repeat      INTERVAL    NOT NULL
);

CREATE OR REPLACE FUNCTION automatic_trigger_function() RETURNS TRIGGER AS $$
BEGIN
    IF (TG_OP = 'INSERT') THEN
        PERFORM pg_notify('inserts', NEW.id::TEXT);
        RETURN NEW;
    ELSIF (TG_OP = 'UPDATE') THEN
        PERFORM pg_notify('updates', NEW.id::TEXT);
        RETURN NEW;
    ELSIF (TG_OP = 'DELETE') THEN
        PERFORM pg_notify('deletes', OLD.id::TEXT);
        RETURN OLD;
    END IF;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER automatic_trigger AFTER INSERT OR UPDATE OR DELETE ON automatic
FOR EACH ROW
EXECUTE FUNCTION automatic_trigger_function();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE automatic;
-- +goose StatementEnd
