package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) Test(c *fiber.Ctx) error {
    msg, err := h.usecase.Test()

    if err != nil {
        c.Status(http.StatusInternalServerError).Send([]byte(`{"msg":"ok"}`))
    }

    return c.JSON(fiber.Map{"msg": msg})
}
