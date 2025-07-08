package domain

type User struct {
    Username    string  `db:"username" json:"username"`
    Balance     int     `db:"balance" json:"balance"`
}
