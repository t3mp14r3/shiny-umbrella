package domain

import "time"

type User struct {
    Username    string  `db:"username" json:"username"`
    Balance     int     `db:"balance" json:"balance"`
}

type Tournament struct {
    ID          int64 `db:"id" json:"id"`
    Price       int `db:"price" json:"price"`
    MinUsers    int `db:"min_users" json:"min_users"`
    MaxUsers    int `db:"max_users" json:"max_users"`
    Bets        int `db:"bets" json:"bets"`
    StartsAt    time.Time   `db:"starts_at" json:"starts_at"`
    Duration    int64   `db:"duration" json:"duration"`
    Rewards     []uint8 `db:"rewards" json:"-"`
    Participants    int `db:"participants" json:"participants"`
    Registered  bool    `db:"registered" json:"registered"`
}

type Automatic struct {
    ID          int64 `db:"id" json:"id"`
    Price       int `db:"price" json:"price"`
    MinUsers    int `db:"min_users" json:"min_users"`
    MaxUsers    int `db:"max_users" json:"max_users"`
    Bets        int `db:"bets" json:"bets"`
    Duration    int64   `db:"duration" json:"duration"`
    Repeat      int64   `db:"repeat" json:"repeat"`
    Rewards     []uint8 `db:"rewards" json:"-"`
}

type Reward struct {
    ID              int64 `db:"id" json:"id,omitempty"`
    TournamentID    int64 `db:"tournament_id" json:"tournament_id,omitempty"`
    Place           int   `db:"place" json:"place"`
    Prize           int   `db:"prize" json:"prize"`
}

type Registration struct {
    ID              int64 `db:"id" json:"id,omitempty"`
    TournamentID    int64 `db:"tournament_id" json:"tournament_id,omitempty"`
    Username        string  `db:"username" json:"username"`
    RegisteredAt    time.Time   `db:"registered_at" json:"registered_at"`
}

type Score struct {
    ID              int64 `db:"id" json:"id,omitempty"`
    TournamentID    int64 `db:"tournament_id" json:"tournament_id,omitempty"`
    Username        string  `db:"username" json:"username"`
    Score           int     `db:"score" json:"score"`
    PlacedAt        time.Time   `db:"placed_at" json:"placed_at"`
}
