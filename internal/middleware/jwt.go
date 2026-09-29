package middleware

import (
	"context"
	"net/http"
	"strings"

	"BE_GKMI_NTC_GO/internal/response"
	"BE_GKMI_NTC_GO/internal/token"

	"github.com/golang-jwt/jwt/v5"
)

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
			response.Error(w, http.StatusUnauthorized, "Authorization header required")
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)

		if len(parts) != 2 || parts[0] != "Bearer" || parts[1] == "" {
			response.Error(w, http.StatusUnauthorized, "Invalid authorization header")
			return
		}

		parsedToken, err := token.ParseAccessToken(parts[1])
		if err != nil || !parsedToken.Valid {
			response.Error(w, http.StatusUnauthorized, "Invalid or expired token")
			return
		}

		claims, ok := parsedToken.Claims.(jwt.MapClaims)
		if !ok {
			response.Error(w, http.StatusUnauthorized, "Invalid token claims")
			return
		}

		userID, ok := claims["user_id"].(float64)
		if !ok {
			response.Error(w, http.StatusUnauthorized, "Invalid user ID claim")
			return
		}

		username, ok := claims["username"].(string)
		if !ok {
			response.Error(w, http.StatusUnauthorized, "Invalid username claim")
			return
		}

		role, ok := claims["role"].(string)
		if !ok {
			response.Error(w, http.StatusUnauthorized, "Invalid role claim")
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		ctx = context.WithValue(ctx, UsernameKey, username)
		ctx = context.WithValue(ctx, RoleKey, role)

		next.ServeHTTP(w, r.WithContext(ctx))
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
