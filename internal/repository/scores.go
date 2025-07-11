package repository

import (
	"context"
	"log"

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

func (r *Repository) Calculate(ctx context.Context, id int64) error {
    log.Println("calculating new result")

    tournament, err := r.GetTournament(ctx, id)
    
    if err != nil {
        return err
    }

    count, err := r.CountRegistrations(ctx, id)
    
    if err != nil {
        return err
    }
    
    if count < tournament.MinUsers {
        _, err := r.db.ExecContext(ctx, "UPDATE users SET balance = balance + $1 WHERE username IN (SELECT username FROM registrations WHERE tournament_id = $2);", tournament.Price, tournament.ID)
    
        if err != nil {
            r.logger.Error("Failed to return money from registering!", zap.Error(err))
            return err
        }
    } else {
        var rewards []int
        var usernames []struct{
            Username    string  `db:"username"`
            Score       int     `db:"score"`
        }

        err := r.db.SelectContext(ctx, &rewards, "SELECT prize FROM rewards WHERE tournament_id = $1 ORDER BY prize DESC;", tournament.ID)
        
        if err != nil {
            r.logger.Error("Failed to get reward prizes for tournament!", zap.Error(err))
            return err
        }
        
        err = r.db.SelectContext(ctx, &usernames, "SELECT DISTINCT ON (username) username, score FROM scores WHERE tournament_id = $1 ORDER BY username, score DESC LIMIT $2;", tournament.ID, len(rewards))
        
        if err != nil {
            r.logger.Error("Failed to get winner usernames!", zap.Error(err))
            return err
        }
   
        tx, err := r.db.Beginx()

        if err != nil {
            r.logger.Error("Failed to open a new transaction!", zap.Error(err))
            return err
        }

        log.Println(rewards)
        log.Println(usernames)

        if len(usernames) < len(rewards) {
            rewards = rewards[:len(usernames)]
        }

        for i, reward := range rewards {
            _, err := tx.ExecContext(ctx, "UPDATE users SET balance = balance + $1 WHERE username = $2;", reward, usernames[i].Username)
        
            if err != nil {
                tx.Rollback()
                r.logger.Error("Failed to update winner's balance!", zap.Error(err))
                return err
            }
        }

        err = tx.Commit()
        
        if err != nil {
            r.logger.Error("Failed to commit a transaction!", zap.Error(err))
            return err
        }
    }

    return nil
}
