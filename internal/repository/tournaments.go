package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/t3mp14r3/shiny-umbrella/internal/domain"
	"go.uber.org/zap"
)

func (r *Repository) GetTournaments(ctx context.Context) ([]domain.Tournament, error) {
    var out []domain.Tournament

    err := r.db.SelectContext(ctx, &out, "SELECT id, price, min_users, max_users, bets, starts_at, EXTRACT(EPOCH FROM duration)::BIGINT AS duration FROM tournaments;")
   
    if errors.Is(err, sql.ErrNoRows) {
        return []domain.Tournament{}, nil
    } else if err != nil {
        r.logger.Error("Failed to select tournament records!", zap.Error(err))
        return nil, err
    }

    return out, nil
}
