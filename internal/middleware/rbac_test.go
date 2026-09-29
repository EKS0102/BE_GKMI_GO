package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"BE_GKMI_NTC_GO/internal/response"
)

func TestRequireRoleAllowed(t *testing.T) {
	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		response.JSON(w, http.StatusOK, map[string]string{"message": "allowed"})
	})

	request := httptest.NewRequest(http.MethodGet, "/test", nil)
	request = request.WithContext(contextWithRole(request.Context(), "admin"))
	recorder := httptest.NewRecorder()

	RequireRole("admin", "staff")(next).ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}
	if !nextCalled {
		t.Fatal("expected next handler to be called")
	}
}

func TestRequireRoleForbidden(t *testing.T) {
	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	})

	request := httptest.NewRequest(http.MethodGet, "/test", nil)
	request = request.WithContext(contextWithRole(request.Context(), "viewer"))
	recorder := httptest.NewRecorder()

	RequireRole("admin", "staff")(next).ServeHTTP(recorder, request)

	if recorder.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", recorder.Code)
	}
	if nextCalled {
		t.Fatal("expected next handler not to be called")
	}
}

func TestRequireRoleMissingRole(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	request := httptest.NewRequest(http.MethodGet, "/test", nil)
	recorder := httptest.NewRecorder()

	RequireRole("admin")(next).ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", recorder.Code)
	}
}

func contextWithRole(ctx context.Context, role string) context.Context {
	return context.WithValue(ctx, RoleKey, role)
}
