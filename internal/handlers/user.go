package handlers

import (
	"hublish-be-go/internal/database"
	"hublish-be-go/internal/models"
	"hublish-be-go/internal/types"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func GetCurrentUser(c *fiber.Ctx) error {
	loggedInUserID := c.Locals("user").(*jwt.Token).Claims.(*types.CustomClaims).UserID
	db := database.DB
	var foundUser models.User
	if findResult := db.First(&foundUser, "id = ?", loggedInUserID); findResult.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot find a user"})
	}

	return c.JSON(foundUser)
}
