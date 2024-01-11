package routes

import (
	"hublish-be-go/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutesSetup(r fiber.Router) {
	authRoutes := r.Group("/api/auth")
	authRoutes.Post("/signup", handlers.SignUpUser)
}
