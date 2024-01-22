package utils

import (
	"hublish-be-go/internal/types"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWTToken(user_id string, exp time.Time, key string) (string, error) {
	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS256, types.CustomClaims{
		UserID: user_id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	})

	token, err := rawToken.SignedString([]byte(key))
	return token, err
}
