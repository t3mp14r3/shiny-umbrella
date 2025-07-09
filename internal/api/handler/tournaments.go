package handler

import (
	"context"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/t3mp14r3/shiny-umbrella/internal/api/usecase"
	"github.com/t3mp14r3/shiny-umbrella/internal/errors"
)

func (h *Handler) GetTournaments(c *fiber.Ctx) error {
    username, ok := c.Queries()["username"]

    ctx, _ := context.WithCancel(c.Context())

    var out []usecase.TournamentOutput
    var err error

    if ok {
        out, err = h.usecase.GetTournaments(ctx, username)
    } else {
        out, err = h.usecase.GetTournaments(ctx)
    }

    if err != nil {
        return SendError(c, errors.Codes[err], err)
    }

    return c.Status(http.StatusOK).JSON(out)
}

type registerInput struct {
    Username        string  `json:"username"`
    TournamentID    int64   `json:"tournament_id"`
}

func (i registerInput) Valid() bool {
    if len(i.Username) > 0 && i.TournamentID >= 0 {
        return true
    }
    return false
}

func (h *Handler) Register(c *fiber.Ctx) error {
    ctx, _ := context.WithCancel(c.Context())
    var input registerInput
    
    err := c.BodyParser(&input)

    if err != nil || !input.Valid() {
        return SendError(c, http.StatusBadRequest, errors.New("Incorrect request data!"))
    }

    err = h.usecase.Register(ctx, usecase.RegisterInput{
        TournamentID: input.TournamentID,
        Username: input.Username,
    })

    if err != nil {
        return SendError(c, errors.Codes[err], err)
    }

    return c.Status(http.StatusOK).JSON(map[string]string{})
}

type scoreInput struct {
    Username        string  `json:"username"`
    TournamentID    int64   `json:"tournament_id"`
    Score           int     `json:"score"`
}

func (i scoreInput) Valid() bool {
    if len(i.Username) > 0 && i.TournamentID >= 0 && i.Score > 0 {
        return true
    }
    return false
}

func (h *Handler) Score(c *fiber.Ctx) error {
    ctx, _ := context.WithCancel(c.Context())
    var input scoreInput
    
    err := c.BodyParser(&input)

    if err != nil || !input.Valid() {
        return SendError(c, http.StatusBadRequest, errors.New("Incorrect request data!"))
    }

    err = h.usecase.Score(ctx, usecase.ScoreInput{
        TournamentID: input.TournamentID,
        Username: input.Username,
        Score: input.Score,
    })

    if err != nil {
        return SendError(c, errors.Codes[err], err)
    }

    return c.Status(http.StatusOK).JSON(map[string]string{})
}
