package handlers

import (
	"errors"

	"hublish-be-go/internal/database"
	"hublish-be-go/internal/models"
	"hublish-be-go/internal/types"
	"hublish-be-go/internal/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
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

	profile := new(types.ChangeProfileRequestBody)
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

func ChangeUserPassword(c *fiber.Ctx) error {

	body := new(types.ChangePasswordRequestBody)
	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot parse a request body"})
	}

	if err := validator.V.Struct(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Request body is invalid."})
	}

	loggedInUserID := c.Locals("user").(*jwt.Token).Claims.(*types.CustomClaims).UserID
	db := database.DB
	var foundUser models.User
	if findResult := db.First(&foundUser, "id = ?", loggedInUserID); findResult.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot find a user."})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(body.CurrentPassword)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Your current password does not match."})
	}

	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.NewPassword), 10)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot hash a password."})
	}

	foundUser.Password = string(newHashedPassword)
	if updateResult := db.Select("password").Save(foundUser); updateResult.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot update a user password."})
	}
	return c.JSON(foundUser)
}

func ChangeUserEmail(c *fiber.Ctx) error {

	body := new(types.ChangeEmailRequestBody)
	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot parse a request body"})
	}

	loggedInUserID := c.Locals("user").(*jwt.Token).Claims.(*types.CustomClaims).UserID
	db := database.DB
	var foundUser models.User
	if findUserResult := db.First(&foundUser, "id = ?", loggedInUserID); findUserResult.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot find a user."})
	}
	if err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(body.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Your password does not match."})
	}

	var foundEmail models.User
	findEmailResult := db.Where("email = ?", body.NewEmail).First(&foundEmail)
	if findEmailResult.Error != nil && !errors.Is(findEmailResult.Error, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot find an email."})
	}
	if findEmailResult.Error == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "This email is already used."})
	}

	foundUser.Email = body.NewEmail
	if updateEmailResult := db.Select("email").Save(foundUser); updateEmailResult.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot update user email."})
	}

	return c.JSON(foundUser)
}
