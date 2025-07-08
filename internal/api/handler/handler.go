package handler

import (
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

func New(cfg *config.Config, usecase *usecase.UseCase, logger *zap.Logger) (*Handler, error) {
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

    router.Get("/test", h.Test)

    router.Use(func(c *fiber.Ctx) error {
		return SendError(c, 404, errors.New("Route not found!"))
	})

    return h, nil

}

func (h *Handler) Run() error {
    return h.router.Listen(h.addr)
}
