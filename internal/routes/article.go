package routes

import (
	"hublish-be-go/internal/handlers"
	"hublish-be-go/internal/middlewares"

	"github.com/gofiber/fiber/v2"
)

func articleRouteSetup(r fiber.Router) {
	articleRoutes := r.Group("/api/articles")

	articleRoutes.Post("/", middlewares.RequireAuth, handlers.CreateArticle)
}
