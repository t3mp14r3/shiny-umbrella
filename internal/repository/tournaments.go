package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/t3mp14r3/shiny-umbrella/internal/domain"
	"go.uber.org/zap"
)

func (r *Repository) GetTournaments(ctx context.Context, username ...string) ([]domain.Tournament, error) {
    var out []domain.Tournament
    var err error

    if len(username) > 0 {
        err = r.db.SelectContext(ctx, &out, "SELECT t.id, t.price, t.min_users, t.max_users, t.bets, t.starts_at, EXTRACT(EPOCH FROM t.duration)::BIGINT AS duration, (SELECT json_agg(json_build_object('place', r.place, 'prize', r.prize)) FROM rewards r WHERE r.tournament_id = t.id) AS rewards, COUNT(DISTINCT s.username) AS participants, (SELECT COUNT(s.id) > 0 FROM scores s WHERE s.tournament_id = t.id AND s.username = $1) AS registered FROM tournaments t LEFT JOIN scores s ON t.id = s.tournament_id GROUP BY t.id;", username[0])
    } else {
        err = r.db.SelectContext(ctx, &out, "SELECT t.id, t.price, t.min_users, t.max_users, t.bets, t.starts_at, EXTRACT(EPOCH FROM t.duration)::BIGINT AS duration, (SELECT json_agg(json_build_object('place', r.place, 'prize', r.prize)) FROM rewards r WHERE r.tournament_id = t.id) AS rewards, COUNT(DISTINCT s.username) AS participants FROM tournaments t LEFT JOIN scores s ON t.id = s.tournament_id GROUP BY t.id;")
    }
   
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
