package usecase

import (
	"github.com/t3mp14r3/shiny-umbrella/internal/repository"
	"go.uber.org/zap"
)

type UseCase struct {
    repo    *repository.Repository
    logger  *zap.Logger
}

func New(repo *repository.Repository, logger *zap.Logger) *UseCase {
    return &UseCase{
        repo: repo,
        logger: logger,
    }
}
