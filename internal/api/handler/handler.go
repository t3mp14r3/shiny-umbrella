package handler

import (
	"context"
	"sync"

	"github.com/gofiber/fiber/v2"
	//    "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/t3mp14r3/shiny-umbrella/internal/api/usecase"
	"github.com/t3mp14r3/shiny-umbrella/internal/config"
	"github.com/t3mp14r3/shiny-umbrella/internal/errors"
	"go.uber.org/zap"
)

type Handler struct {
    router  *fiber.App
    usecase *usecase.UseCase
    logger  *zap.Logger
    addr    string
}

func New(cfg *config.Config, usecase *usecase.UseCase, logger *zap.Logger) *Handler {
    router := fiber.New(fiber.Config{
        DisableStartupMessage: true,
    })

    h := &Handler{
        router: router,
        usecase: usecase,
        logger: logger,
        addr: cfg.AppAddr,
    }

    usersApi := router.Group("/users")
    usersApi.Post("", h.CreateUser)
    usersApi.Put("", h.UpdateUser)
    usersApi.Get("", h.GetUsers)
    
    tournamentsApi := router.Group("/tournaments")
    tournamentsApi.Get("", h.GetTournaments)
    tournamentsApi.Post("/register", h.Register)
    tournamentsApi.Post("/score", h.Score)

    router.Use(func(c *fiber.Ctx) error {
		return SendError(c, 404, errors.New("Route not found!"))
	})

    return h
}

func (h *Handler) Run(ctx context.Context) error {
    errChan := make(chan error, 1)

    wg := &sync.WaitGroup{}
    wg.Add(1)

    go func() {
        defer wg.Done()
        if err := h.router.Listen(h.addr); err != nil {
            h.logger.Error("handler error", zap.Error(err))
            errChan <- err
        }
    }()

    var err error

    select {
        case <-ctx.Done():
            h.router.Shutdown()
        case err = <-errChan:
    }

    wg.Wait()

    return err
}
