package token

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func getSecretKey() (string, error) {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		return "", os.ErrNotExist
	}
	return secretKey, nil
}

func GenerateAccessToken(userID int, username string, role string) (string, error) {
	secretKey, err := getSecretKey()
	if err != nil {
		return "", err
	}
	now := time.Now()
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"role":     role,
		"exp":      now.Add(15 * time.Minute).Unix(),
		"iat":      now.Unix(),
	}
	signedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return signedToken.SignedString([]byte(secretKey))
}

func ParseAccessToken(tokenString string) (*jwt.Token, error) {
	secretKey, err := getSecretKey()
	if err != nil {
		return nil, err
	}
	return jwt.Parse(tokenString, func(parsedToken *jwt.Token) (interface{}, error) {
		if parsedToken.Method != jwt.SigningMethodHS256 {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secretKey), nil
	})
}
