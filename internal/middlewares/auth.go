package middlewares

import (
	"hublish-be-go/internal/types"
	"os"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

var RequireAuth = jwtware.New(jwtware.Config{
	SigningKey:   jwtware.SigningKey{Key: []byte(os.Getenv("ACCESSTOKEN_KEY"))},
	Claims:       &types.CustomClaims{},
	ErrorHandler: requireAuthErrorHandler,
})

var IsLoggedIn = jwtware.New(jwtware.Config{
	SigningKey:     jwtware.SigningKey{Key: []byte(os.Getenv("ACCESSTOKEN_KEY"))},
	Claims:         &types.CustomClaims{},
	SuccessHandler: isLoggedInSuccessHandler,
	ErrorHandler:   isLoggedInErrorHandler,
})

func requireAuthErrorHandler(c *fiber.Ctx, e error) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token expired."})
}

func isLoggedInSuccessHandler(c *fiber.Ctx) error {
	c.Locals("isLoggedIn", true)
	return c.Next()
}

func isLoggedInErrorHandler(c *fiber.Ctx, e error) error {
	c.Locals("isLoggedIn", false)
	return c.Next()
}
