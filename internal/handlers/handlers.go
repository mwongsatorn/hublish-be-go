package handlers

import "github.com/gofiber/fiber/v2"

func HomeHanlder(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "hello from fiber"})
}
