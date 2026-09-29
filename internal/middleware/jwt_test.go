package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"BE_GKMI_NTC_GO/internal/token"
)

func TestJWTMissingAuthorization(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	recorder := httptest.NewRecorder()

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("next handler should not be called")
	})

	JWT(next).ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusUnauthorized)
	}
}

func TestJWTInvalidAuthorizationFormat(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	request.Header.Set("Authorization", "Basic abc123")
	recorder := httptest.NewRecorder()

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("next handler should not be called")
	})

	JWT(next).ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusUnauthorized)
	}
}

func TestJWTEmptyBearerToken(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	request.Header.Set("Authorization", "Bearer ")
	recorder := httptest.NewRecorder()

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("next handler should not be called")
	})

	JWT(next).ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusUnauthorized)
	}
}

func TestJWTInvalidToken(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-key")

	request := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	request.Header.Set("Authorization", "Bearer invalid-token")
	recorder := httptest.NewRecorder()

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("next handler should not be called")
	})

	JWT(next).ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusUnauthorized)
	}
}

func TestJWTValidToken(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-key")

	accessToken, err := token.GenerateAccessToken(10, "eko", "viewer")
	if err != nil {
		t.Fatalf("GenerateAccessToken() error = %v", err)
	}

	request := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	request.Header.Set("Authorization", "Bearer "+accessToken)
	recorder := httptest.NewRecorder()

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, ok := GetUserID(r.Context())
		if !ok {
			t.Fatal("user ID not found in context")
		}
		if userID != 10 {
			t.Fatalf("user ID = %d, want 10", userID)
		}

		username, ok := GetUsername(r.Context())
		if !ok {
			t.Fatal("username not found in context")
		}
		if username != "eko" {
			t.Fatalf("username = %q, want %q", username, "eko")
		}

		role, ok := GetRole(r.Context())
		if !ok {
			t.Fatal("role not found in context")
		}
		if role != "viewer" {
			t.Fatalf("role = %q, want %q", role, "viewer")
		}

		w.WriteHeader(http.StatusOK)
	})

	JWT(next).ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
}

func TestJWTExpiredToken(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-key")

	expiredToken := strings.TrimSpace("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjF9.invalid")

	request := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	request.Header.Set("Authorization", "Bearer "+expiredToken)
	recorder := httptest.NewRecorder()

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("next handler should not be called")
	})

	JWT(next).ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusUnauthorized)
	}
}
