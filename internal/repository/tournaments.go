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

    err := r.db.SelectContext(ctx, &out, "SELECT t.id, t.price, t.min_users, t.max_users, t.bets, t.starts_at, EXTRACT(EPOCH FROM t.duration)::BIGINT AS duration, json_agg(json_build_object('place', r.place, 'prize', r.prize)) AS rewards FROM tournaments t INNER JOIN rewards r ON t.id = r.tournament_id GROUP BY t.id;")
   
    if errors.Is(err, sql.ErrNoRows) {
        return []domain.Tournament{}, nil
    } else if err != nil {
        r.logger.Error("Failed to select tournament records!", zap.Error(err))
        return nil, err
    }

    return out, nil
}

func (r *Repository) CreateTournament(ctx context.Context, input domain.Tournament) (*domain.Tournament, error) {
    var out domain.Tournament

    err := r.db.GetContext(ctx, &out, "INSERT INTO tournaments(price, min_users, max_users, bets, starts_at, duration) VALUES($1, $2, $3, $4, $5, $6) RETURNING id, price, min_users, max_users, bets, starts_at, EXTRACT(EPOCH FROM duration)::BIGINT AS duration;",
        input.Price,
        input.MinUsers,
        input.MaxUsers,
        input.Bets,
        input.StartsAt,
        input.Duration,
    )
   
    if errors.Is(err, sql.ErrNoRows) {
        return nil, nil
    } else if err != nil {
        r.logger.Error("Failed to create a tournament record!", zap.Error(err))
        return nil, err
    }

    return &out, nil
}
