package usecase

import (
	"context"
	"time"

	"github.com/t3mp14r3/shiny-umbrella/internal/errors"
)

type TournamentOutput struct {
    ID          int `json:"id"`
    Price       int `json:"price"`
    MinUsers    int `json:"min_users"`
    MaxUsers    int `json:"max_users"`
    Bets        int `json:"bets"`
    Status      string   `json:"status"`
    StartsAt    string   `json:"starts_at"`
    EndsAt      string   `json:"ends_at"`
}

func (u *UseCase) GetTournaments(ctx context.Context) ([]TournamentOutput, error) {
    list, err := u.repo.GetTournaments(ctx)

    if err != nil {
        return nil, errors.ErrorSomethingWentWrong
    }

    var out []TournamentOutput
    now := time.Now()

    for _, t := range list {
        var status string
        endsAt := t.StartsAt.Add(time.Duration(t.Duration * int64(time.Second)))

        if t.Canceled {
            status = "Canceled"
        } else if t.StartsAt.After(now) {
            status = "Upcoming"
        } else if endsAt.Before(now) {
            status = "Ended"
        } else {
            status = "Active"
        }

        out = append(out, TournamentOutput{
            ID: t.ID,
            Price: t.Price,
            MinUsers: t.MinUsers,
            MaxUsers: t.MaxUsers,
            Bets: t.Bets,
            Status: status,
            StartsAt: t.StartsAt.Format("02.01.2006 15:04:05"),
            EndsAt: endsAt.Format("02.01.2006 15:04:05"),
        })
    }

    return out, nil
}
