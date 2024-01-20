package utils

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func SetCookie(c *fiber.Ctx, name string, value string, expiration time.Time) {
	c.Cookie(buildCookie(name, value, expiration))
}

func ClearCookie(c *fiber.Ctx, name string) {
	c.Cookie(buildCookie(name, "", time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)))
}

func buildCookie(name string, value string, expires time.Time) *fiber.Cookie {
	cookie := new(fiber.Cookie)
	*cookie = fiber.Cookie{
		Name:     name,
		Value:    value,
		HTTPOnly: true,
		Expires:  expires,
	}
	return cookie
}
