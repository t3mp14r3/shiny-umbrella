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
