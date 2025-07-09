package usecase

import (
	"context"
	"encoding/json"
	"time"

	"github.com/t3mp14r3/shiny-umbrella/internal/domain"
	"github.com/t3mp14r3/shiny-umbrella/internal/errors"
)

type TournamentOutput struct {
    ID          int64 `json:"id"`
    Price       int `json:"price"`
    MinUsers    int `json:"min_users"`
    MaxUsers    int `json:"max_users"`
    Bets        int `json:"bets"`
    Status      string   `json:"status"`
    StartsAt    string   `json:"starts_at"`
    EndsAt      string   `json:"ends_at"`
    Rewards     []domain.Reward `json:"rewards"`
    Participants    int `json:"participants"`
    Registered  bool    `json:"registered"`
}

func (u *UseCase) GetTournaments(ctx context.Context, username ...string) ([]TournamentOutput, error) {
    var list []domain.Tournament
    var err error

    if len(username) > 0 {
        list, err = u.repo.GetTournaments(ctx, username[0])
    } else {
        list, err = u.repo.GetTournaments(ctx)
    }

    if err != nil {
        return nil, errors.ErrorSomethingWentWrong
    }

    var out []TournamentOutput
    now := time.Now()

    for _, t := range list {
        var status string
        rewards := []domain.Reward{}

        endsAt := t.StartsAt.Add(time.Duration(t.Duration * int64(time.Second)))

        if t.StartsAt.After(now) {
            status = "Upcoming"
        } else if endsAt.Before(now) && t.Participants < t.MinUsers {
            status = "Canceled"
        } else if endsAt.Before(now) {
            status = "Ended"
        } else {
            status = "Active"
        }

        json.Unmarshal([]byte(t.Rewards), &rewards)

        out = append(out, TournamentOutput{
            ID: t.ID,
            Price: t.Price,
            MinUsers: t.MinUsers,
            MaxUsers: t.MaxUsers,
            Bets: t.Bets,
            Status: status,
            StartsAt: t.StartsAt.Format("02.01.2006 15:04:05"),
            EndsAt: endsAt.Format("02.01.2006 15:04:05"),
            Rewards: rewards,
            Participants: t.Participants,
            Registered: t.Registered,
        })
    }

    return out, nil
}

type RegisterInput struct {
    Username        string  `json:"username"`
    TournamentID    int64   `json:"tournament_id"`
}

func (u *UseCase) Register(ctx context.Context, input RegisterInput) error {
    t, err := u.repo.GetTournament(ctx, input.TournamentID)
   
    if err != nil {
        return errors.ErrorSomethingWentWrong
    } else if err == nil && t == nil {
        return errors.ErrorTournamentNotFound
    }

    count, err := u.repo.CountRegistrations(ctx, input.TournamentID)
    
    if err != nil {
        return errors.ErrorSomethingWentWrong
    }

    user, err := u.repo.GetUser(ctx, input.Username)
    
    if err != nil {
        return errors.ErrorSomethingWentWrong
    }

    now := time.Now()
    endsAt := t.StartsAt.Add(time.Duration(t.Duration * int64(time.Second)))

    if endsAt.Before(now) {
        return errors.ErrorTournamentEnded
    } else if t.MaxUsers == count {
        return errors.ErrorTournamentMaxed
    } else if t.Price > user.Balance {
        return errors.ErrorNotEnoughFunds
    }

    tx, err := u.repo.Begin()

    if err != nil {
        return errors.ErrorSomethingWentWrong
    }

    _, err = u.repo.UpdateUserTx(ctx, tx, domain.User{
        Username: input.Username,
        Balance: user.Balance - t.Price,
    })
    
    if err != nil {
        u.repo.Rollback(tx)
        return errors.ErrorSomethingWentWrong
    }

    _, err = u.repo.CreateRegistrationTx(ctx, tx, domain.Registration{
        TournamentID: input.TournamentID,
        Username: input.Username,
    })
    
    if err != nil {
        u.repo.Rollback(tx)
        return errors.ErrorSomethingWentWrong
    }

    err = u.repo.Commit(tx)
    
    if err != nil {
        return errors.ErrorSomethingWentWrong
    }

    return nil
}

type ScoreInput struct {
    Username        string  `json:"username"`
    TournamentID    int64   `json:"tournament_id"`
    Score           int     `json:"score"`
}

func (u *UseCase) Score(ctx context.Context, input ScoreInput) error {
    r, err := u.repo.GetRegistration(ctx, input.TournamentID, input.Username)
   
    if err != nil {
        return errors.ErrorSomethingWentWrong
    } else if err == nil && r == nil {
        return errors.ErrorNotRegistered
    }
    
    t, err := u.repo.GetTournament(ctx, input.TournamentID)
   
    if err != nil {
        return errors.ErrorSomethingWentWrong
    } else if err == nil && r == nil {
        return errors.ErrorTournamentNotFound
    }

    count, err := u.repo.CountScores(ctx, input.TournamentID, input.Username)
    
    if err != nil {
        return errors.ErrorSomethingWentWrong
    }

    if count >= t.Bets {
        return errors.ErrorMaximumBets
    }

    err = u.repo.CreateScore(ctx, domain.Score{
        TournamentID: input.TournamentID,
        Username: input.Username,
        Score: input.Score,
    })
    
    if err != nil {
        return errors.ErrorSomethingWentWrong
    }

    return nil
}
