package middleware

import (
	"context"
	"net/http"
	"strings"

	"BE_GKMI_NTC_GO/internal/response"

	"github.com/golang-jwt/jwt/v5"
)

const secretKey = "gkmi-secret-key"

type contextKey string

const (
	UserIDKey   contextKey = "user_id"
	UsernameKey contextKey = "username"
	RoleKey     contextKey = "role"
)

func JWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			response.Error(
				w,
				http.StatusUnauthorized,
				"Authorization header required",
			)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)

		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(
				w,
				http.StatusUnauthorized,
				"Invalid authorization header",
			)
			return
		}

		tokenString := parts[1]

		token, err := jwt.Parse(
			tokenString,
			func(token *jwt.Token) (interface{}, error) {

				if token.Method != jwt.SigningMethodHS256 {
					return nil, jwt.ErrSignatureInvalid
				}

				return []byte(secretKey), nil
			},
		)

		if err != nil || !token.Valid {
			response.Error(
				w,
				http.StatusUnauthorized,
				"Invalid or expired token",
			)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			response.Error(
				w,
				http.StatusUnauthorized,
				"Invalid token claims",
			)
			return
		}

		ctx := context.WithValue(
			r.Context(),
			UserIDKey,
			claims["user_id"],
		)

		ctx = context.WithValue(
			ctx,
			UsernameKey,
			claims["username"],
		)

		ctx = context.WithValue(
			ctx,
			RoleKey,
			claims["role"],
		)

		next.ServeHTTP(
			w,
			r.WithContext(ctx),
		)
	})
}

func GetUserID(ctx context.Context) (int, bool) {
	value := ctx.Value(UserIDKey)

	userID, ok := value.(float64)
	if !ok {
		return 0, false
	}

	return int(userID), true
}

func GetUsername(ctx context.Context) (string, bool) {
	value := ctx.Value(UsernameKey)

	username, ok := value.(string)
	if !ok {
		return "", false
	}

	return username, true
}

func GetRole(ctx context.Context) (string, bool) {
	value := ctx.Value(RoleKey)

	role, ok := value.(string)
	if !ok {
		return "", false
	}

	return role, true
}
