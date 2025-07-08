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

func New(cfg *config.Config, logger *zap.Logger) (*Repository, error) {
    db, err := sqlx.Connect("postgres", cfg.RepositoryConnString())

    if err != nil {
        log.Printf("Failed to open database connection: %v\n", err)
        return nil, err
    }

    if err := migrate(db.DB, cfg.DBMigrations); err != nil {
        return nil, err
    }

    return &Repository{
        db: db,
        logger: logger,
    }, nil
}

func migrate(db *sql.DB, path string) error {
    if err := goose.SetDialect("postgres"); err != nil {
        log.Printf("Failed to set dialect for database migrations: %v\n", err)
        return err
    }


    if err := goose.Up(db, path); err != nil {
        log.Printf("Failed to run database migrations: %v\n", err)
        return err
    }

    return nil
}
