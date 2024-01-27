package handlers

import (
	"hublish-be-go/internal/database"
	"hublish-be-go/internal/models"
	"hublish-be-go/internal/types"
	"hublish-be-go/internal/utils"
	"hublish-be-go/internal/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func CreateArticle(c *fiber.Ctx) error {

	body := new(types.CreateArticleRequestBody)
	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot parse a request body."})
	}

	if err := validator.V.Struct(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Article information is not valid."})
	}

	loggedInUserID := c.Locals("user").(*jwt.Token).Claims.(*types.CustomClaims).UserID

	slug := utils.GenerateSlug(body.Title)
	var tags []string
	if body.Tags == nil {
		tags = []string{}
	} else {
		tags = *body.Tags
	}

	createdArticle := models.Article{
		Title:    body.Title,
		Content:  body.Content,
		Tags:     tags,
		Slug:     slug,
		AuthorID: loggedInUserID,
	}
	db := database.DB
	if createResult := db.Create(&createdArticle); createResult.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot create an article."})
	}

	res, err := utils.ResponseOmitFilter(createdArticle, []string{"author"})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot make a response object"})
	}

	return c.Status(fiber.StatusCreated).JSON(res)
	
}