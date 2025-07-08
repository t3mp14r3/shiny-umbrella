package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/t3mp14r3/shiny-umbrella/internal/api/usecase"
	"github.com/t3mp14r3/shiny-umbrella/internal/config"
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
        logger: logger,
        addr: cfg.AppAddr,
    }

    router.Get("/test", h.Test)

    return h, nil

}

func (h *Handler) Run() error {
    return h.router.Listen(h.addr)
}
