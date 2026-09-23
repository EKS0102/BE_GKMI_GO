package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const secretKey = "gkmi-secret-key"

func GenerateAccessToken(userID int, username string, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(15 * time.Minute).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secretKey))
}
