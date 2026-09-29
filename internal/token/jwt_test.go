package token

import (
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestGenerateAccessToken(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-key")
	tokenString, err := GenerateAccessToken(10, "eko", "viewer")
	if err != nil {
		t.Fatalf("GenerateAccessToken() error = %v", err)
	}
	if tokenString == "" {
		t.Fatal("expected token, got empty string")
	}
	token, err := ParseAccessToken(tokenString)
	if err != nil {
		t.Fatalf("ParseAccessToken() error = %v", err)
	}
	if !token.Valid {
		t.Fatal("expected token to be valid")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal("expected MapClaims")
	}
	if claims["user_id"] != float64(10) {
		t.Fatalf("user_id = %v, want 10", claims["user_id"])
	}
	if claims["username"] != "eko" {
		t.Fatalf("username = %v, want eko", claims["username"])
	}
	if claims["role"] != "viewer" {
		t.Fatalf("role = %v, want viewer", claims["role"])
	}
	if _, ok := claims["exp"]; !ok {
		t.Fatal("expected exp claim")
	}
	if _, ok := claims["iat"]; !ok {
		t.Fatal("expected iat claim")
	}
}

func TestGenerateAccessTokenRequiresSecret(t *testing.T) {
	previous := os.Getenv("JWT_SECRET")
	os.Unsetenv("JWT_SECRET")
	t.Cleanup(func() {
		if previous == "" {
			os.Unsetenv("JWT_SECRET")
		} else {
			os.Setenv("JWT_SECRET", previous)
		}
	})
	_, err := GenerateAccessToken(10, "eko", "viewer")
	if err == nil {
		t.Fatal("expected error when JWT_SECRET is missing")
	}
}

func TestParseAccessTokenRejectsInvalidSignature(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-key")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"user_id": 10, "username": "eko", "role": "viewer", "exp": time.Now().Add(15 * time.Minute).Unix()})
	tokenString, err := token.SignedString([]byte("wrong-secret"))
	if err != nil {
		t.Fatalf("SignedString() error = %v", err)
	}
	parsedToken, err := ParseAccessToken(tokenString)
	if err == nil || parsedToken.Valid {
		t.Fatal("expected invalid signature error")
	}
}

func TestParseAccessTokenRejectsExpiredToken(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-key")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"user_id": 10, "username": "eko", "role": "viewer", "exp": time.Now().Add(-1 * time.Minute).Unix()})
	tokenString, err := token.SignedString([]byte("test-secret-key"))
	if err != nil {
		t.Fatalf("SignedString() error = %v", err)
	}
	parsedToken, err := ParseAccessToken(tokenString)
	if err == nil || parsedToken.Valid {
		t.Fatal("expected expired token error")
	}
}
