package routes

import (
	"hublish-be-go/internal/handlers"
	"hublish-be-go/internal/middlewares"

	"github.com/gofiber/fiber/v2"
)

func userRouteSetup(r fiber.Router) {
	userRoutes := r.Group("/api/users")

	userRoutes.Get("/", middlewares.IsLoggedIn, handlers.SearchUsers)

	userRoutes.Get("/current", middlewares.RequireAuth, handlers.GetCurrentUser)
	userRoutes.Get("/:username/profile", middlewares.IsLoggedIn, handlers.GetUserProfile)

	userRoutes.Put("/settings/profile", middlewares.RequireAuth, handlers.ChangeUserProfile)
	userRoutes.Put("/settings/password", middlewares.RequireAuth, handlers.ChangeUserPassword)
	userRoutes.Put("/settings/email", middlewares.RequireAuth, handlers.ChangeUserEmail)

	userRoutes.Post("/:username/follow", middlewares.RequireAuth, handlers.FollowUser)
	userRoutes.Delete("/:username/follow", middlewares.RequireAuth, handlers.UnfollowUser)

	userRoutes.Get("/:username/followers", middlewares.IsLoggedIn, handlers.GetUserFollowers)
	userRoutes.Get("/:username/followings", middlewares.IsLoggedIn, handlers.GetUserFollowings)
}
