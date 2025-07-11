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

    err := r.db.SelectContext(ctx, &out, "SELECT a.id, a.price, a.min_users, a.max_users, a.bets, EXTRACT(EPOCH FROM duration)::BIGINT AS duration, EXTRACT(EPOCH FROM repeat)::BIGINT AS repeat, (SELECT json_agg(json_build_object('place', r.place, 'prize', r.prize)) FROM automatic_rewards r WHERE r.automatic_id= a.id) AS rewards FROM automatic a;")
   
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

    err := r.db.GetContext(ctx, &out, "SELECT a.id, a.price, a.min_users, a.max_users, a.bets, EXTRACT(EPOCH FROM duration)::BIGINT AS duration, EXTRACT(EPOCH FROM repeat)::BIGINT AS repeat, (SELECT json_agg(json_build_object('place', r.place, 'prize', r.prize)) FROM automatic_rewards r WHERE r.automatic_id= a.id) AS rewards FROM automatic a WHERE id = $1;", id)
   
    if errors.Is(err, sql.ErrNoRows) {
        return nil, nil
    } else if err != nil {
        r.logger.Error("Failed to get automatic record!", zap.Error(err))
        return nil, err
    }

    return &out, nil
}
