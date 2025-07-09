package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/t3mp14r3/shiny-umbrella/internal/domain"
	"go.uber.org/zap"
)

func (r *Repository) CreateUser(ctx context.Context, user domain.User) (*domain.User, error) {
    var out domain.User

    err := r.db.GetContext(ctx, &out, "INSERT INTO users(username, balance) VALUES($1, $2) RETURNING *;", user.Username, user.Balance)
    
    if err != nil {
        r.logger.Error("Failed to create new user record!", zap.Error(err))
        return nil, err
    }

    return &out, nil
}

func (r *Repository) CreateUserTx(ctx context.Context, tx *sqlx.Tx, user domain.User) (*domain.User, error) {
    var out domain.User

    err := tx.GetContext(ctx, &out, "INSERT INTO users(username, balance) VALUES($1, $2) RETURNING *;", user.Username, user.Balance)
    
    if err != nil {
        r.logger.Error("Failed to create new user record!", zap.Error(err))
        return nil, err
    }

    return &out, nil
}

func (r *Repository) GetUser(ctx context.Context, username string) (*domain.User, error) {
    var out domain.User

    err := r.db.GetContext(ctx, &out, "SELECT * FROM users WHERE username = $1;", username)

    if errors.Is(err, sql.ErrNoRows) {
        return nil, nil
    } else if err != nil {
        r.logger.Error("Failed to select user record!", zap.Error(err))
        return nil, err
    }

    return &out, nil
}

func (r *Repository) GetUsers(ctx context.Context) ([]domain.User, error) {
    var out []domain.User

    err := r.db.SelectContext(ctx, &out, "SELECT * FROM users;")

    if errors.Is(err, sql.ErrNoRows) {
        return []domain.User{}, nil
    } else if err != nil {
        r.logger.Error("Failed to select user records!", zap.Error(err))
        return nil, err
    }

    return out, nil
}

func (r *Repository) UpdateUser(ctx context.Context, user domain.User) (*domain.User, error) {
    var out domain.User

    err := r.db.GetContext(ctx, &out, "UPDATE users SET balance = $1 WHERE username = $2 RETURNING *;", user.Balance, user.Username)
    
    if errors.Is(err, sql.ErrNoRows) {
        return nil, nil
    } else if err != nil {
        r.logger.Error("Failed to update user record!", zap.Error(err))
        return nil, err
    }

    return &out, nil
}

func (r *Repository) UpdateUserTx(ctx context.Context, tx *sqlx.Tx, user domain.User) (*domain.User, error) {
    var out domain.User

    err := tx.GetContext(ctx, &out, "UPDATE users SET balance = $1 WHERE username = $2 RETURNING *;", user.Balance, user.Username)
    
    if errors.Is(err, sql.ErrNoRows) {
        return nil, nil
    } else if err != nil {
        r.logger.Error("Failed to update user record!", zap.Error(err))
        return nil, err
    }

    return &out, nil
}
