package handlers

import (
	"errors"
	"strconv"

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

func GetUserProfile(c *fiber.Ctx) error {

	loggedInUserID := "00000000-0000-0000-0000-000000000000"
	targetUsername := c.Params("username")
	if c.Locals("isLoggedIn") == true {
		loggedInUserID = c.Locals("user").(*jwt.Token).Claims.(*types.CustomClaims).UserID
	}

	var targetUser models.User
	db := database.DB
	findTargetUserResult := db.Where("username = ?", targetUsername).First(&targetUser)
	if findTargetUserResult.Error != nil && !errors.Is(findTargetUserResult.Error, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot find a user."})
	}
	if errors.Is(findTargetUserResult.Error, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No user found."})
	}

	var userProfile types.UserQuery
	if findUserProfileResult := db.Table("users u").
		Select([]string{"u.*", "CASE WHEN f.id IS NOT NULL THEN true ELSE false END AS followed"}).
		Joins("LEFT JOIN follows f on f.following_id = u.id AND f.follower_id = ?", loggedInUserID).
		Where("u.username = ?", targetUser.Username).
		First(&userProfile); findUserProfileResult.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot find a user profile"})
	}

	return c.JSON(userProfile)
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

func FollowUser(c *fiber.Ctx) error {

	loggedInUserID := c.Locals("user").(*jwt.Token).Claims.(*types.CustomClaims).UserID
	targetUsername := c.Params("username")

	db := database.DB
	var userToFollow models.User
	findTargetUserResult := db.First(&userToFollow, "username = ?", targetUsername)
	if findTargetUserResult.Error != nil && !errors.Is(findTargetUserResult.Error, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot find a user."})
	}
	if errors.Is(findTargetUserResult.Error, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No user found."})
	}
	if loggedInUserID == userToFollow.ID {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "You cannot follow yourself."})
	}

	isFollowingResult := db.
		Where(models.Follow{FollowingID: userToFollow.ID, FollowerID: loggedInUserID}).
		First(&models.Follow{})
	if isFollowingResult.Error != nil && !errors.Is(isFollowingResult.Error, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot find a follow relation."})
	}
	if isFollowingResult.RowsAffected != 0 {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "You've already followed this user."})
	}

	err := db.Transaction(func(tx *gorm.DB) error {

		if followResult := tx.Create(&models.Follow{
			FollowingID: userToFollow.ID,
			FollowerID:  loggedInUserID,
		}); followResult.Error != nil {
			return errors.New("Cannot follow a user: Error on create follow relation.")
		}

		if updateLoggedInUserResult := tx.Model(&models.User{}).
			Where("id = ?", loggedInUserID).
			Updates(map[string]interface{}{
				"following_count": gorm.Expr("following_count + 1"),
			}); updateLoggedInUserResult.Error != nil {
			return errors.New("Cannot follow a user: Error on update a logged-in user.")
		}

		userToFollow.FollowerCount += 1
		if updateTargetUserResult := tx.Select("follower_count").
			Save(userToFollow); updateTargetUserResult.Error != nil {
			return errors.New("Cannot follow a user: Error on update a target user.")
		}

		return nil
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(userToFollow)
}

func UnfollowUser(c *fiber.Ctx) error {

	loggedInUserID := c.Locals("user").(*jwt.Token).Claims.(*types.CustomClaims).UserID
	targetUsername := c.Params("username")

	db := database.DB
	var userToUnfollow models.User
	findTargetUserResult := db.First(&userToUnfollow, "username = ?", targetUsername)
	if findTargetUserResult.Error != nil && !errors.Is(findTargetUserResult.Error, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot find a user."})
	}
	if errors.Is(findTargetUserResult.Error, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No user found."})
	}
	if loggedInUserID == userToUnfollow.ID {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "You cannot unfollow yourself."})
	}

	isFollowingResult := db.
		Where(models.Follow{FollowingID: userToUnfollow.ID, FollowerID: loggedInUserID}).
		First(&models.Follow{})
	if isFollowingResult.Error != nil && !errors.Is(isFollowingResult.Error, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot find a follow relation."})
	}
	if isFollowingResult.RowsAffected == 0 {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "You haven't followed this user yet."})
	}

	err := db.Transaction(func(tx *gorm.DB) error {

		if unfollowResult := tx.
			Where("follower_id = ? AND following_id = ?", loggedInUserID, userToUnfollow.ID).
			Delete(&models.Follow{}); unfollowResult.Error != nil {
			return errors.New("Cannot follow a user: Error on delete follow relation.")
		}

		if updateLoggedInUserResult := tx.Model(&models.User{}).
			Where("id = ?", loggedInUserID).
			Updates(map[string]interface{}{
				"following_count": gorm.Expr("following_count - 1"),
			}); updateLoggedInUserResult.Error != nil {
			return errors.New("Cannot follow a user: Error on update a logged-in user.")
		}

		userToUnfollow.FollowerCount -= 1
		if updateTargetUserResult := tx.Select("follower_count").
			Save(userToUnfollow); updateTargetUserResult.Error != nil {
			return errors.New("Cannot follow a user: Error on update a target user.")
		}

		return nil
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(userToUnfollow)

}

func GetUserFollowers(c *fiber.Ctx) error {

	loggedInUserID := "00000000-0000-0000-0000-000000000000"
	targetUsername := c.Params("username")
	if c.Locals("isLoggedIn") == true {
		loggedInUserID = c.Locals("user").(*jwt.Token).Claims.(*types.CustomClaims).UserID
	}

	var targetUser models.User
	db := database.DB
	findTargetUserResult := db.Where("username = ?", targetUsername).First(&targetUser)
	if findTargetUserResult.Error != nil && !errors.Is(findTargetUserResult.Error, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot find a user."})
	}
	if errors.Is(findTargetUserResult.Error, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No user found."})
	}

	var userFollowers []types.ShortUserQuery
	if findUserFollowersResult := db.Table("follows f").
		Select([]string{"u.id", "u.username", "u.name", "u.bio", "u.image",
			"CASE WHEN f_log.id IS NOT NULL THEN true ELSE false END AS followed"}).
		Joins("JOIN users u ON u.id = f.follower_id").
		Joins("LEFT JOIN follows f_log ON f_log.following_id = f.follower_id AND f_log.follower_id = ?", loggedInUserID).
		Where("f.following_id = ?", targetUser.ID).
		Find(&userFollowers); findUserFollowersResult.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot find user followers."})
	}

	return c.JSON(userFollowers)
}

func GetUserFollowings(c *fiber.Ctx) error {

	loggedInUserID := "00000000-0000-0000-0000-000000000000"
	targetUsername := c.Params("username")
	if c.Locals("isLoggedIn") == true {
		loggedInUserID = c.Locals("user").(*jwt.Token).Claims.(*types.CustomClaims).UserID
	}

	var targetUser models.User
	db := database.DB
	findTargetUserResult := db.Where("username = ?", targetUsername).First(&targetUser)
	if findTargetUserResult.Error != nil && !errors.Is(findTargetUserResult.Error, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot find a user."})
	}
	if errors.Is(findTargetUserResult.Error, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No user found."})
	}

	var userFollowings []types.ShortUserQuery
	if findUserFollowingsResult := db.Table("follows f").
		Select([]string{"u.id", "u.username", "u.name", "u.bio", "u.image",
			"CASE WHEN f_log.id IS NOT NULL THEN true ELSE false END AS followed"}).
		Joins("JOIN users u ON u.id = f.following_id").
		Joins("LEFT JOIN follows f_log ON f_log.following_id = f.following_id AND f_log.follower_id = ?", loggedInUserID).
		Where("f.follower_id = ?", targetUser.ID).
		Find(&userFollowings); findUserFollowingsResult.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot find user followers."})
	}

	return c.JSON(userFollowings)
}

func SearchUsers(c *fiber.Ctx) error {

	loggedInUserID := "00000000-0000-0000-0000-000000000000"
	if c.Locals("isLoggedIn") == true {
		loggedInUserID = c.Locals("user").(*jwt.Token).Claims.(*types.CustomClaims).UserID
	}

	query := c.Query("query")
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Query is not valid"})
	}
	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Query is not valid"})
	}

	db := database.DB

	sqlQuery := "%" + query + "%"
	var totalResults int64
	if countTotalResults := db.Model(&models.User{}).Where("username LIKE ? OR name LIKE ?", sqlQuery, sqlQuery).
		Count(&totalResults); countTotalResults.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot count total results."})
	}

	var foundUsers []types.ShortUserQuery
	if findUsersResult := db.Table("users u").
		Select([]string{"u.id", "u.username", "u.name", "u.bio", "u.image",
			"CASE WHEN f.id IS NOT NULL THEN true ELSE false END AS followed"}).
		Joins("LEFT JOIN follows f ON f.following_id = u.id AND f.follower_id = ?", loggedInUserID).
		Where("username LIKE ? OR name LIKE ?", sqlQuery, sqlQuery).
		Offset(limit * (page - 1)).
		Limit(limit).Find(&foundUsers); findUsersResult.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot find Users."})
	}

	res := types.SearchQuery[types.ShortUserQuery]{
		TotalResults: int(totalResults),
		TotalPages:   (int(totalResults) + limit - 1) / limit,
		Page:         page,
		Results:      foundUsers,
	}

	return c.JSON(res)
}
