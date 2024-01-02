package routes

import (
	"github.com/gofiber/fiber/v2"
	"hublish-be-go/internal/handlers"
)

func RegisterRoutes(app fiber.Router) {
	app.Get("/", handlers.HomeHanlder)
}
