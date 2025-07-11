package repository

import (
	"database/sql"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
	"github.com/t3mp14r3/shiny-umbrella/internal/config"
	"go.uber.org/zap"
	_ "github.com/lib/pq"
)

type Repository struct {
    db      *sqlx.DB
    logger  *zap.Logger
}

func New(cfg *config.Config, logger *zap.Logger) *Repository {
    db, err := sqlx.Connect("postgres", cfg.RepositoryConnString())

    if err != nil {
        log.Fatalf("Failed to open database connection: %v\n", err)
        return nil
    }

    migrate(db.DB, cfg.DBMigrations)

    return &Repository{
        db: db,
        logger: logger,
    }
}

func migrate(db *sql.DB, path string) {
    if err := goose.SetDialect("postgres"); err != nil {
        log.Fatalf("Failed to set dialect for database migrations: %v\n", err)
    }


    if err := goose.Up(db, path); err != nil {
        log.Fatalf("Failed to run database migrations: %v\n", err)
    }
}

func (r *Repository) Close() {
    r.db.Close()
}

func (r *Repository) Begin() (*sqlx.Tx, error) {
    tx, err := r.db.Beginx()

    if err != nil {
        r.logger.Error("Failed to begin new transaction!", zap.Error(err))
        return nil, err
    }
    
    return tx, nil
}

func (r *Repository) Commit(tx *sqlx.Tx) error {
    err := tx.Commit()

    if err != nil {
        r.logger.Error("Failed to commit transaction!", zap.Error(err))
        return err
    }
    
    return nil
}

func (r *Repository) Rollback(tx *sqlx.Tx) error {
    err := tx.Rollback()

    if err != nil {
        r.logger.Error("Failed to rollback transaction!", zap.Error(err))
        return err
    }
    
    return nil
}
