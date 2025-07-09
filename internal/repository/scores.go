package repository

import (
	"context"

	"github.com/t3mp14r3/shiny-umbrella/internal/domain"
	"go.uber.org/zap"
)

func (r *Repository) CountScores(ctx context.Context, tournament_id int64, username string) (int, error) {
    var count int

    err := r.db.GetContext(ctx, &count, "SELECT COUNT(id) FROM scores WHERE tournament_id = $1 AND username = $2;", tournament_id, username)

    if err != nil {
        r.logger.Error("Failed to count score records!", zap.Error(err))
        return 0, err
    }

    return count, nil
}


func (r *Repository) CreateScore(ctx context.Context, input domain.Score) error {
    _, err := r.db.ExecContext(ctx, "INSERT INTO scores(tournament_id, username, score) VALUES($1, $2, $3);", input.TournamentID, input.Username, input.Score)

    if err != nil {
        r.logger.Error("Failed to create score record!", zap.Error(err))
        return err
    }

    return nil
}
