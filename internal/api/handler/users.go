package handler

import (
	"context"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/t3mp14r3/shiny-umbrella/internal/domain"
    "github.com/t3mp14r3/shiny-umbrella/internal/errors"
)

type createUserInput struct {
    Username    string  `json:"username"`
    Balance     int     `json:"balance"`
}

func (i createUserInput) Valid() bool {
    if len(i.Username) > 0 && i.Balance > 0 {
        return true
    }
    return false
}

func (h *Handler) CreateUser(c *fiber.Ctx) error {
    ctx, _ := context.WithCancel(c.Context())
    var input createUserInput

    err := c.BodyParser(&input)

    if err != nil || !input.Valid() {
        return SendError(c, http.StatusBadRequest, errors.New("Incorrect request data!"))
    }

    out, err := h.usecase.CreateUser(ctx, domain.User{
        Username: input.Username,
        Balance: input.Balance,
    })

    if err != nil {
        return SendError(c, errors.Codes[err], err)
    }

    return c.Status(http.StatusOK).JSON(out)
}

type updateUserInput struct {
    Username    string  `json:"username"`
    Balance     int     `json:"balance"`
}

func (i updateUserInput) Valid() bool {
    if len(i.Username) > 0 && i.Balance > 0 {
        return true
    }
    return false
}

func (h *Handler) UpdateUser(c *fiber.Ctx) error {
    ctx, _ := context.WithCancel(c.Context())
    var input updateUserInput

    err := c.BodyParser(&input)

    if err != nil || !input.Valid() {
        return SendError(c, http.StatusBadRequest, errors.New("Incorrect request data!"))
    }

    out, err := h.usecase.UpdateUser(ctx, domain.User{
        Username: input.Username,
        Balance: input.Balance,
    })

    if err != nil {
        return SendError(c, errors.Codes[err], err)
    }

    return c.Status(http.StatusOK).JSON(out)
}

func (h *Handler) GetUsers(c *fiber.Ctx) error {
    ctx, _ := context.WithCancel(c.Context())

    out, err := h.usecase.GetUsers(ctx)

    if err != nil {
        return SendError(c, errors.Codes[err], err)
    }

    return c.Status(http.StatusOK).JSON(out)
}
