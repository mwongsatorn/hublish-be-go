package routes

import (
	"hublish-be-go/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app fiber.Router) {
	app.Get("/", handlers.HomeHanlder)
	AuthRoutesSetup(app)
}
