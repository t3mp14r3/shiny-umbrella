package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/t3mp14r3/shiny-umbrella/internal/domain"
	"go.uber.org/zap"
)

func (r *Repository) GetTournaments(ctx context.Context, username ...string) ([]domain.Tournament, error) {
    var out []domain.Tournament
    var err error

    if len(username) > 0 {
        err = r.db.SelectContext(ctx, &out, "SELECT t.id, t.price, t.min_users, t.max_users, t.bets, t.starts_at, EXTRACT(EPOCH FROM t.duration)::BIGINT AS duration, (SELECT json_agg(json_build_object('place', r.place, 'prize', r.prize)) FROM rewards r WHERE r.tournament_id = t.id) AS rewards, COUNT(DISTINCT r.username) AS participants, (SELECT COUNT(r.id) > 0 FROM registrations r WHERE r.tournament_id = t.id AND r.username = $1) AS registered FROM tournaments t LEFT JOIN registrations r ON t.id = r.tournament_id GROUP BY t.id;", username[0])
    } else {
        err = r.db.SelectContext(ctx, &out, "SELECT t.id, t.price, t.min_users, t.max_users, t.bets, t.starts_at, EXTRACT(EPOCH FROM t.duration)::BIGINT AS duration, (SELECT json_agg(json_build_object('place', r.place, 'prize', r.prize)) FROM rewards r WHERE r.tournament_id = t.id) AS rewards, COUNT(DISTINCT r.username) AS participants FROM tournaments t LEFT JOIN registrations r ON t.id = r.tournament_id GROUP BY t.id;")
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

func (r *Repository) GetTournament(ctx context.Context, id int64) (*domain.Tournament, error) {
    var out domain.Tournament
    var err error

    err = r.db.GetContext(ctx, &out, "SELECT t.id, t.price, t.min_users, t.max_users, t.bets, t.starts_at, EXTRACT(EPOCH FROM t.duration)::BIGINT AS duration FROM tournaments t WHERE t.id = $1;", id)
   
    if errors.Is(err, sql.ErrNoRows) {
        return nil, nil
    } else if err != nil {
        r.logger.Error("Failed to get a tournament record!", zap.Error(err))
        return nil, err
    }

    return &out, nil
}

func (r *Repository) CountRegistrations(ctx context.Context, id int64) (int, error) {
    var count int
    var err error

    err = r.db.GetContext(ctx, &count, "SELECT COUNT(r.id) FROM registrations r WHERE r.tournament_id = $1;", id)
   
    if err != nil {
        r.logger.Error("Failed to count registration records!", zap.Error(err))
        return 0, err
    }

    return count, nil
}

func (r *Repository) CreateRegistrationTx(ctx context.Context, tx *sqlx.Tx, input domain.Registration) (*domain.Registration, error) {
    var out domain.Registration

    err := tx.GetContext(ctx, &out, "INSERT INTO registrations(tournament_id, username) VALUES($1, $2) RETURNING *;", input.TournamentID, input.Username)

    if err != nil {
        r.logger.Error("Failed to create a registration record!", zap.Error(err))
        return nil, err
    }

    return &out, nil
}

func (r *Repository) GetRegistration(ctx context.Context, tournament_id int64, username string) (*domain.Registration, error) {
    var out domain.Registration

    err := r.db.GetContext(ctx, &out, "SELECT * FROM registrations WHERE tournament_id = $1 AND username = $2;", tournament_id, username)
   
    if errors.Is(err, sql.ErrNoRows) {
        return nil, nil
    } else if err != nil {
        r.logger.Error("Failed to get registration record!", zap.Error(err))
        return nil, err
    }

    return &out, nil
}
