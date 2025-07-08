package handler

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) Test(c *fiber.Ctx) error {
    msg, err := h.usecase.Test()

    if err != nil {
        c.Status(http.StatusInternalServerError).Send([]byte(`{"msg":"ok"}`))
    }

    return c.Send([]byte(fmt.Sprintf(`{"msg":"%s"}`, msg)))
}
