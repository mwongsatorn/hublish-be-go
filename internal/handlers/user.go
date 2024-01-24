package handlers

import (
	"hublish-be-go/internal/database"
	"hublish-be-go/internal/models"
	"hublish-be-go/internal/types"
	"hublish-be-go/internal/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm/clause"
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

func ChangeUserProfile(c *fiber.Ctx) error {

	profile := new(types.UpdateProfileRequestBody)
	if err := c.BodyParser(profile); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot parse a request body."})
	}
	if err := validator.V.Struct(profile); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Profile information is not valid."})
	}

	loggedInUserID := c.Locals("user").(*jwt.Token).Claims.(*types.CustomClaims).UserID

	db := database.DB
	var updatedUser models.User
	if updateResult := db.Model(&updatedUser).
		Clauses(clause.Returning{}).
		Where("id = ?", loggedInUserID).
		Updates(profile); updateResult.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot update a user profile."})
	}

	return c.JSON(updatedUser)
}
