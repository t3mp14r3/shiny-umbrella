package domain

import "time"

type User struct {
    Username    string  `db:"username" json:"username"`
    Balance     int     `db:"balance" json:"balance"`
}

type Tournament struct {
    ID          int `db:"id" json:"id"`
    Price       int `db:"price" json:"price"`
    MinUsers    int `db:"min_users" json:"min_users"`
    MaxUsers    int `db:"max_users" json:"max_users"`
    Bets        int `db:"bets" json:"bets"`
    Canceled    bool    `db:"canceled" json:"canceled"`
    StartsAt    time.Time   `db:"starts_at" json:"starts_at"`
    Duration    int64   `db:"duration" json:"duration"`
}
