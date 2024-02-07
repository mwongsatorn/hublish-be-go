package routes

import (
	"hublish-be-go/internal/handlers"
	"hublish-be-go/internal/middlewares"

	"github.com/gofiber/fiber/v2"
)

func articleRouteSetup(r fiber.Router) {
	articleRoutes := r.Group("/api/articles")

	articleRoutes.Post("/", middlewares.RequireAuth, handlers.CreateArticle)
	articleRoutes.Get("/:slug", handlers.GetArticle)
	articleRoutes.Put("/:slug", middlewares.RequireAuth, handlers.EditArticle)
	articleRoutes.Delete("/:slug", middlewares.RequireAuth, handlers.DeleteArticle)

	articleRoutes.Post("/:slug/favourite", middlewares.RequireAuth, handlers.FavouriteArticle)
	articleRoutes.Delete("/:slug/favourite", middlewares.RequireAuth, handlers.UnfavouriteArticle)

	articleRoutes.Post("/:slug/comments", middlewares.RequireAuth, handlers.AddComment)
	articleRoutes.Delete("/:slug/comments/:comment_id", middlewares.RequireAuth, handlers.DeleteComment)
}
