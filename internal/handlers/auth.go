package handlers

import (
	"hublish-be-go/internal/database"
	"hublish-be-go/internal/models"
	"hublish-be-go/internal/validator"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"fmt"
)

func SignUpUser(c *fiber.Ctx) error {
	req := new(validator.SignUpRequest)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Something went wrong."})
	}

	if err := validator.V.Struct(req); err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User information is not valid."})
	}

	db := database.DB
	var foundUser models.User
	findResult := db.Where("username = ?", req.Username).Or("email = ?", req.Email).First(&foundUser)

	if findResult.Error == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Username or Email is already used."})
	}

	if findResult.Error == gorm.ErrRecordNotFound {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Something went wrong."})
		}
		newUser := &models.User{
			Username: req.Username,
			Password: string(hashedPassword),
			Email:    req.Email,
		}

		createResult := db.Select([]string{"Username", "Password", "Email"}).Create(newUser)
		if createResult.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Something went wrong."})
		}
		return c.Status(fiber.StatusCreated).JSON(newUser)
	}
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Something went wrong."})
}
