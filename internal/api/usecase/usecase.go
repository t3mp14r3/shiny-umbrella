package usecase

import (
	"github.com/t3mp14r3/shiny-umbrella/internal/config"
	"github.com/t3mp14r3/shiny-umbrella/internal/repository"
	"go.uber.org/zap"
)

type UseCase struct {
    repo    *repository.Repository
    logger  *zap.Logger
    addr    string
}

func New(cfg *config.Config, repo *repository.Repository, logger *zap.Logger) (*UseCase, error) {
    u := &UseCase{
        repo: repo,
        logger: logger,
        addr: cfg.AppAddr,
    }

    return u, nil
}
