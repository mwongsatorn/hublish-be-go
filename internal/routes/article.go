package routes

import (
	"hublish-be-go/internal/handlers"
	"hublish-be-go/internal/middlewares"

	"github.com/gofiber/fiber/v2"
)

func articleRouteSetup(r fiber.Router) {
	articleRoutes := r.Group("/api/articles")

	articleRoutes.Get("/", middlewares.IsLoggedIn, handlers.SearchArticles)
	articleRoutes.Post("/", middlewares.RequireAuth, handlers.CreateArticle)

	articleRoutes.Get("/feed", middlewares.RequireAuth, handlers.GetUserFeedArticles)

	articleRoutes.Get("/:slug", middlewares.IsLoggedIn, handlers.GetArticle)
	articleRoutes.Put("/:slug", middlewares.RequireAuth, handlers.EditArticle)
	articleRoutes.Delete("/:slug", middlewares.RequireAuth, handlers.DeleteArticle)

	articleRoutes.Post("/:slug/favourite", middlewares.RequireAuth, handlers.FavouriteArticle)
	articleRoutes.Delete("/:slug/favourite", middlewares.RequireAuth, handlers.UnfavouriteArticle)

	articleRoutes.Post("/:slug/comments", middlewares.RequireAuth, handlers.AddComment)
	articleRoutes.Delete("/:slug/comments/:comment_id", middlewares.RequireAuth, handlers.DeleteComment)

	articleRoutes.Get("/:slug/comments", handlers.GetComments)

	articleRoutes.Get("/:username/created", middlewares.IsLoggedIn, handlers.GetUserCreatedArticles)
	articleRoutes.Get("/:username/favourite", middlewares.IsLoggedIn, handlers.GetUserFavouriteArticles)
}
