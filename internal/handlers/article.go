package handlers

import (
	"errors"
	"hublish-be-go/internal/database"
	"hublish-be-go/internal/models"
	"hublish-be-go/internal/types"
	"hublish-be-go/internal/utils"
	"hublish-be-go/internal/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

	createdArticle := models.Article{
		Title:    body.Title,
		Content:  body.Content,
		Tags:     body.Tags,
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

func EditArticle(c *fiber.Ctx) error {

	body := new(types.EditArticleRequestBody)
	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot parse a request body."})
	}

	if err := validator.V.Struct(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Article information is not valid."})
	}

	articleSlug := c.Params("slug")
	loggedInUserID := c.Locals("user").(*jwt.Token).Claims.(*types.CustomClaims).UserID

	if body.Title != nil {
		body.Slug = utils.GenerateSlug(*body.Title)
	}

	var editedArticle models.Article
	db := database.DB
	updateResult := db.Model(&editedArticle).
		Clauses(clause.Returning{}).
		Where("author_id = ? AND slug = ?", loggedInUserID, articleSlug).
		Updates(body)

	if updateResult.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot update an article."})
	}

	if updateResult.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No article found."})
	}

	res, err := utils.ResponseOmitFilter(editedArticle, []string{"author"})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot make a response object"})
	}

	return c.JSON(res)

}

func GetArticle(c *fiber.Ctx) error {

	articleSlug := c.Params("slug")
	loggedInUserID := "00000000-0000-0000-0000-000000000000"
	if c.Locals("isLoggedIn") == true {
		loggedInUserID = c.Locals("user").(*jwt.Token).Claims.(*types.CustomClaims).UserID
	}

	var foundArticle types.ArticleQuery
	db := database.DB
	findArticleResult := db.
		Table("articles AS a").
		Select([]string{"a.*", "u.id as aid", "u.username", "u.name", "u.bio", "u.image",
			"CASE WHEN f.id IS NOT NULL THEN true ELSE false END AS favourited"}).
		Joins("JOIN users u ON u.id = a.author_id").
		Joins("LEFT JOIN favourites f ON f.article_id = a.id AND f.user_id = ?", loggedInUserID).
		Where("a.slug = ?", articleSlug).
		Take(&foundArticle)

	if findArticleResult.Error != nil && !errors.Is(findArticleResult.Error, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot find an article."})
	}

	if errors.Is(findArticleResult.Error, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No article found."})
	}

	return c.JSON(foundArticle)
}

func DeleteArticle(c *fiber.Ctx) error {

	articleSlug := c.Params("slug")
	loggedInUserID := c.Locals("user").(*jwt.Token).Claims.(*types.CustomClaims).UserID

	db := database.DB
	deleteResult := db.Where("slug = ? AND author_id = ?", articleSlug, loggedInUserID).Delete(&models.Article{})

	if deleteResult.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot delete an article."})
	}

	if deleteResult.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No article found."})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func FavouriteArticle(c *fiber.Ctx) error {

	articleSlug := c.Params("slug")
	loggedInUserID := c.Locals("user").(*jwt.Token).Claims.(*types.CustomClaims).UserID

	db := database.DB
	var foundArticle models.Article
	findArticleResult := db.Where("slug = ?", articleSlug).First(&foundArticle)
	if findArticleResult.Error != nil && !errors.Is(findArticleResult.Error, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot find an article."})
	}
	if errors.Is(findArticleResult.Error, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No article found."})
	}

	findFavouriteResult := db.Where("user_id = ? AND article_id = ?", loggedInUserID, foundArticle.ID).First(&models.Favourite{})
	if findFavouriteResult.Error != nil && !errors.Is(findFavouriteResult.Error, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot find a favourite relation."})
	}

	if findFavouriteResult.RowsAffected > 0 {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "You've already favourited this article."})
	}

	err := db.Transaction(func(tx *gorm.DB) error {

		newFavouriteRelation := models.Favourite{
			ArticleID: foundArticle.ID,
			UserID:    loggedInUserID,
		}
		if createFavouriteRelationResult := tx.Create(&newFavouriteRelation); createFavouriteRelationResult.Error != nil {
			return errors.New("Cannot favourite this article: Error on create favourite relation.")
		}

		foundArticle.FavouriteCount += 1
		if updateArticleResult := tx.Select("favourite_count").
			Save(&foundArticle); updateArticleResult.Error != nil {
			return errors.New("Cannot favourite this article: Error on update article's favourite count.")
		}

		return nil

	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := utils.ResponseOmitFilter(foundArticle, []string{"author"})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot make a response object."})
	}

	return c.JSON(res)
}

func UnfavouriteArticle(c *fiber.Ctx) error {

	articleSlug := c.Params("slug")
	loggedInUserID := c.Locals("user").(*jwt.Token).Claims.(*types.CustomClaims).UserID

	db := database.DB
	var foundArticle models.Article
	findArticleResult := db.Where("slug = ?", articleSlug).First(&foundArticle)
	if findArticleResult.Error != nil && !errors.Is(findArticleResult.Error, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot find an article."})
	}
	if errors.Is(findArticleResult.Error, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No article found."})
	}

	findFavouriteResult := db.Where("user_id = ? AND article_id = ?", loggedInUserID, foundArticle.ID).First(&models.Favourite{})
	if findFavouriteResult.Error != nil && !errors.Is(findFavouriteResult.Error, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot find a favourite relation."})
	}

	if errors.Is(findFavouriteResult.Error, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "You haven't favourited this article yet."})
	}

	err := db.Transaction(func(tx *gorm.DB) error {

		if deleteFavouriteRelationResult := tx.Where("user_id = ? AND article_id = ?", loggedInUserID, foundArticle.ID).
			Delete(&models.Favourite{}); deleteFavouriteRelationResult.Error != nil {
			return errors.New("Cannot unfavourite this article: Error on delete favourite relation.")
		}

		foundArticle.FavouriteCount -= 1
		if updateArticleResult := tx.Select("favourite_count").
			Save(&foundArticle); updateArticleResult.Error != nil {
			return errors.New("Cannot unfavourite this article: Error on update article's favourite count.")
		}

		return nil

	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := utils.ResponseOmitFilter(foundArticle, []string{"author"})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot make a response object."})
	}

	return c.JSON(res)
}

func AddComment(c *fiber.Ctx) error {

	articleSlug := c.Params("slug")
	loggedInUserID := c.Locals("user").(*jwt.Token).Claims.(*types.CustomClaims).UserID

	req := new(types.AddCommentRequestBody)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot parse a body request."})
	}
	if err := validator.V.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Comment is not valid."})
	}

	db := database.DB
	var foundArticle models.Article
	findArticleResult := db.Where("slug = ?", articleSlug).First(&foundArticle)
	if findArticleResult.Error != nil && !errors.Is(findArticleResult.Error, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot find an article."})
	}
	if errors.Is(findArticleResult.Error, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No article found."})
	}

	newComment := models.Comment{
		Body:            req.Body,
		CommentAuthorID: loggedInUserID,
		ArticleID:       foundArticle.ID,
	}
	if addCommentResult := db.Create(&newComment); addCommentResult.Error != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot create a comment."})
	}

	res, err := utils.ResponseOmitFilter(newComment, []string{"article", "user"})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot make a response object."})
	}

	return c.Status(fiber.StatusCreated).JSON(res)
}

func DeleteComment(c *fiber.Ctx) error {

	articleSlug := c.Params("slug")
	commentID := c.Params("comment_id")
	loggedInUserID := c.Locals("user").(*jwt.Token).Claims.(*types.CustomClaims).UserID

	db := database.DB
	var foundArticle models.Article
	findArticleResult := db.Where("slug = ?", articleSlug).First(&foundArticle)
	if findArticleResult.Error != nil && !errors.Is(findArticleResult.Error, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot find an article."})
	}
	if errors.Is(findArticleResult.Error, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No article found."})
	}

	var foundComment models.Comment
	findCommentResult := db.Where("id = ? AND \"commentAuthor_id\" = ?", commentID, loggedInUserID).First(&foundComment)
	if findCommentResult.Error != nil && !errors.Is(findCommentResult.Error, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot find a comment"})
	}
	if errors.Is(findCommentResult.Error, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No comment found."})
	}

	if deleteCommentResult := db.Delete(&models.Comment{
		CommonFields:    models.CommonFields{ID: foundComment.ID},
		CommentAuthorID: loggedInUserID,
		ArticleID:       foundArticle.ID,
	}); deleteCommentResult.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot delete a comment"})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func GetComments(c *fiber.Ctx) error {

	articleSlug := c.Params("slug")

	db := database.DB
	var foundArticle models.Article
	findArticleResult := db.Where("slug = ?", articleSlug).First(&foundArticle)
	if findArticleResult.Error != nil && !errors.Is(findArticleResult.Error, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot find an article."})
	}
	if errors.Is(findArticleResult.Error, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No article found."})
	}

	var comments []types.CommentQuery
	if findCommentsResult := db.Table("comments c").
		Select([]string{"c.*", "u.id as caid", "u.username", "u.name", "u.image"}).
		Joins("JOIN users u ON \"commentAuthor_id\" = u.id").
		Where("c.article_id = ?", foundArticle.ID).
		Find(&comments); findCommentsResult.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot find comments"})
	}

	return c.JSON(comments)

}
