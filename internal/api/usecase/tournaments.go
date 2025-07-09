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
