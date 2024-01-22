package routes

import (
	"hublish-be-go/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

func authRoutesSetup(r fiber.Router) {
	authRoutes := r.Group("/api/auth")
	authRoutes.Post("/signup", handlers.SignUpUser)
	authRoutes.Post("/login", handlers.LogInUser)
	authRoutes.Get("/refresh", handlers.RefreshAccessToken)
	authRoutes.Delete("/logout", handlers.LogOutUser)
}
