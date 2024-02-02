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
	userRoutes.Put("/settings/password", middlewares.RequireAuth, handlers.ChangeUserPassword)
	userRoutes.Put("/settings/email", middlewares.RequireAuth, handlers.ChangeUserEmail)

	userRoutes.Post("/:username/follow", middlewares.RequireAuth, handlers.FollowUser)
	userRoutes.Delete("/:username/follow", middlewares.RequireAuth, handlers.UnfollowUser)
}
