package usecase

import (
	"context"

	"github.com/t3mp14r3/shiny-umbrella/internal/domain"
	"github.com/t3mp14r3/shiny-umbrella/internal/errors"
)

func (u *UseCase) CreateUser(ctx context.Context, user domain.User) (*domain.User, error) {
    check, err := u.repo.GetUser(ctx, user.Username)

    if err == nil && check != nil {
        return nil, errors.ErrorUsernameInUse
    } else if err != nil {
        return nil, errors.ErrorSomethingWentWrong
    }

    return u.repo.CreateUser(ctx, user)
}

func (u *UseCase) UpdateUser(ctx context.Context, user domain.User) (*domain.User, error) {
    out, err := u.repo.UpdateUser(ctx, user)

    if err == nil && out == nil {
        return nil, errors.ErrorUserNotFound
    } else if err != nil {
        return nil, errors.ErrorSomethingWentWrong
    }

    return out, nil
}
