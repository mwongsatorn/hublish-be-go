package handlers

import (
	"hublish-be-go/internal/database"
	"hublish-be-go/internal/models"
	"hublish-be-go/internal/types"
	"hublish-be-go/internal/utils"
	"hublish-be-go/internal/validator"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"errors"
	"fmt"
	"slices"
	"time"
)

var (
	accessTokenKey  = os.Getenv("ACCESSTOKEN_KEY")
	refreshTokenKey = os.Getenv("REFRESHTOKEN_KEY")
)

func SignUpUser(c *fiber.Ctx) error {
	req := new(types.SignUpRequestBody)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot parse a body request."})
	}

	if err := validator.V.Struct(req); err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User information is not valid."})
	}

	db := database.DB
	var foundUser models.User
	findResult := db.Where("username = ?", req.Username).Or("email = ?", req.Email).First(&foundUser)

	if findResult.Error != nil && !errors.Is(findResult.Error, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot find a user."})
	}

	if findResult.Error == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Username or Email is already used."})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot hash a password."})
	}
	newUser := &models.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Email:    req.Email,
	}
	if createResult := db.Select([]string{"Username", "Password", "Email"}).Create(newUser); createResult.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot create a user."})
	}

	return c.Status(fiber.StatusCreated).JSON(newUser)
}

func LogInUser(c *fiber.Ctx) error {

	req := new(types.LoginRequestBody)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot parse a body request."})
	}

	db := database.DB
	var foundUser models.User
	findResult := db.Where("username = ?", req.Username).First(&foundUser)
	if findResult.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Username or Password is incorrect"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(req.Password)); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Username or Password is incorrect"})
	}

	accessToken, err := utils.GenerateJWTToken(foundUser.ID, time.Now().Add(time.Minute*15), accessTokenKey)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot generate a access token."})
	}

	refreshToken, err := utils.GenerateJWTToken(foundUser.ID, time.Now().Add(time.Hour*24), refreshTokenKey)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot generate a refresh token."})
	}

	foundUser.RefreshTokens = append(foundUser.RefreshTokens, refreshToken)
	db.Select("refresh_tokens").Save(&foundUser)

	utils.SetCookie(c, "refreshToken", refreshToken, time.Now().Add(time.Hour*24))

	loginRes := struct {
		models.User
		AccessToken string `json:"accessToken"`
	}{
		foundUser,
		accessToken,
	}
	return c.Status(fiber.StatusOK).JSON(loginRes)
}

func RefreshAccessToken(c *fiber.Ctx) error {

	refreshToken := c.Cookies("refreshToken")

	if refreshToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Refresh token required"})
	}
	utils.ClearCookie(c, "refreshToken")

	token, err := jwt.ParseWithClaims(refreshToken, &types.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(refreshTokenKey), nil
	})
	claims, ok := token.Claims.(*types.CustomClaims)

	if (err != nil && !errors.Is(err, jwt.ErrTokenInvalidClaims)) || !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot parse a refresh token."})
	}

	db := database.DB
	var foundUser models.User
	findResult := db.Where("? = ANY(refresh_tokens) AND id = ?", refreshToken, claims.UserID).First(&foundUser)
	if findResult.Error != nil && !errors.Is(findResult.Error, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Something went wrong."})
	}
	if errors.Is(findResult.Error, gorm.ErrRecordNotFound) {
		if clearTokensResult := db.Model(&models.User{}).Where("id = ?", claims.UserID).Update("refresh_tokens", "{}"); clearTokensResult.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot clear tokens."})
		}
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Reuse refresh token detected."})
	}

	foundUser.RefreshTokens = slices.DeleteFunc(foundUser.RefreshTokens, func(element string) bool {
		return element == refreshToken
	})

	if !token.Valid {
		if removeTokenResult := db.Select("refresh_tokens").Save(&foundUser); removeTokenResult.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot remove token from the list."})
		}
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token expired."})
	}

	newAccessToken, err := utils.GenerateJWTToken(claims.UserID, time.Now().Add(time.Minute*15), accessTokenKey)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot generate a access token."})
	}

	newRefreshToken, err := utils.GenerateJWTToken(claims.UserID, time.Now().Add(time.Hour*24*3), refreshTokenKey)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot generate a refrest token."})
	}

	foundUser.RefreshTokens = append(foundUser.RefreshTokens, newRefreshToken)
	if updateTokenResult := db.Select("refresh_tokens").Save(&foundUser); updateTokenResult.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Update token list error."})
	}

	utils.SetCookie(c, "refreshToken", newRefreshToken, time.Now().Add(time.Hour*24))

	return c.JSON(fiber.Map{"accessToken": newAccessToken})
}

func LogOutUser(c *fiber.Ctx) error {

	refreshToken := c.Cookies("refreshToken")
	utils.ClearCookie(c, "refreshToken")
	if refreshToken == "" {
		return c.SendStatus(fiber.StatusNoContent)
	}

	token, err := jwt.ParseWithClaims(refreshToken, &types.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(refreshTokenKey), nil
	})
	claims, ok := token.Claims.(*types.CustomClaims)
	if !token.Valid || err != nil || !ok {
		return c.SendStatus(fiber.StatusNoContent)
	}

	db := database.DB
	var foundUser models.User
	findResult := db.Where("id = ? AND ? = ANY(refresh_tokens)", claims.UserID, refreshToken).First(&foundUser)
	if findResult.Error != nil {
		c.SendStatus(fiber.StatusNoContent)
	}

	foundUser.RefreshTokens = slices.DeleteFunc(foundUser.RefreshTokens, func(element string) bool {
		return element == refreshToken
	})
	db.Select("refresh_tokens").Save(&foundUser)

	return c.SendStatus(fiber.StatusNoContent)
}
