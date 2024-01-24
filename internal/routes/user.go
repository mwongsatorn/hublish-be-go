package routes

import (
	"hublish-be-go/internal/handlers"
	"hublish-be-go/internal/middlewares"

	"github.com/gofiber/fiber/v2"
)

func userRouteSetup(r fiber.Router) {
	userRoutes := r.Group("/api/users")

	userRoutes.Get("/current", middlewares.RequireAuth, handlers.GetCurrentUser)
	userRoutes.Put("/settings/profile", middlewares.RequireAuth, handlers.ChangeUserProfile)
}
