package handler

import "github.com/gofiber/fiber/v2"

func SendError(c *fiber.Ctx, code int, msg error) error {
    return c.Status(code).JSON(fiber.Map{
        "error": msg.Error(),
    })
}
