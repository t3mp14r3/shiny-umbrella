package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/t3mp14r3/shiny-umbrella/internal/domain"
	"go.uber.org/zap"
)

func (r *Repository) GetAutomatics(ctx context.Context) ([]domain.Automatic, error) {
    var out []domain.Automatic

    err := r.db.SelectContext(ctx, &out, "SELECT id, price, min_users, max_users, bets, starts_at, EXTRACT(EPOCH FROM duration)::BIGINT AS duration, EXTRACT(EPOCH FROM repeat)::BIGINT AS repeat FROM automatic;")
   
    if errors.Is(err, sql.ErrNoRows) {
        return []domain.Automatic{}, nil
    } else if err != nil {
        r.logger.Error("Failed to select automatic records!", zap.Error(err))
        return nil, err
    }

    return out, nil
}

func (r *Repository) GetAutomatic(ctx context.Context, id int64) (*domain.Automatic, error) {
    var out domain.Automatic

    err := r.db.GetContext(ctx, &out, "SELECT id, price, min_users, max_users, bets, starts_at, EXTRACT(EPOCH FROM duration)::BIGINT AS duration, EXTRACT(EPOCH FROM repeat)::BIGINT AS repeat FROM automatic WHERE id = $1;", id)
   
    if errors.Is(err, sql.ErrNoRows) {
        return nil, nil
    } else if err != nil {
        r.logger.Error("Failed to get automatic record!", zap.Error(err))
        return nil, err
    }

    return &out, nil
}
