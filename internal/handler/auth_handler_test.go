package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"BE_GKMI_NTC_GO/internal/model"
)

type fakeAuthService struct {
	user                 *model.User
	refreshToken         string
	err                  error
	refreshUser          *model.User
	newRefreshToken      string
	refreshErr           error
	loginCalled          bool
	refreshCalled        bool
	receivedRefreshToken string
	logoutCalled         bool
	receivedLogoutToken  string
}

func (f *fakeAuthService) Login(ctx context.Context, request model.LoginRequest) (*model.User, string, error) {
	f.loginCalled = true
	if f.err != nil {
		return nil, "", f.err
	}
	return f.user, f.refreshToken, nil
}

func (f *fakeAuthService) Refresh(ctx context.Context, refreshToken string) (*model.User, string, error) {
	f.refreshCalled = true
	f.receivedRefreshToken = refreshToken
	if f.refreshErr != nil {
		return nil, "", f.refreshErr
	}
	return f.refreshUser, f.newRefreshToken, nil
}

func (f *fakeAuthService) Logout(ctx context.Context, refreshToken string) error {
	f.logoutCalled = true
	f.receivedLogoutToken = refreshToken
	return f.refreshErr
}

var _ AuthServiceInterface = (*fakeAuthService)(nil)

func TestLoginSuccess(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret")
	fakeService := &fakeAuthService{
		user:         &model.User{ID: 10, Username: "eko", Role: "viewer"},
		refreshToken: "refresh-token-test",
	}
	handler := Login(fakeService)

	request := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(`{"username":"eko","password":"rahasia123"}`))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", recorder.Code, recorder.Body.String())
	}
	if !strings.Contains(recorder.Body.String(), `"refresh_token":"refresh-token-test"`) {
		t.Fatalf("expected refresh token in response, got %s", recorder.Body.String())
	}
	if !strings.Contains(recorder.Body.String(), `"token_type":"Bearer"`) {
		t.Fatalf("expected Bearer token type, got %s", recorder.Body.String())
	}
	if !strings.Contains(recorder.Body.String(), `"access_token":`) {
		t.Fatalf("expected access token in response, got %s", recorder.Body.String())
	}
	if !fakeService.loginCalled {
		t.Fatal("expected Login to be called")
	}
}

func TestLoginMethodNotAllowed(t *testing.T) {
	fakeService := &fakeAuthService{}
	handler := Login(fakeService)

	request := httptest.NewRequest(http.MethodGet, "/api/auth/login", nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", recorder.Code)
	}
}

func TestLoginInvalidBody(t *testing.T) {
	fakeService := &fakeAuthService{}
	handler := Login(fakeService)

	request := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(`{"username":`))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", recorder.Code)
	}
	if fakeService.loginCalled {
		t.Fatal("Login should not be called for invalid JSON")
	}
}

func TestLoginInvalidCredentials(t *testing.T) {
	fakeService := &fakeAuthService{err: errors.New("invalid credentials")}
	handler := Login(fakeService)

	request := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(`{"username":"eko","password":"salah"}`))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", recorder.Code)
	}
}

func TestLoginRepositoryError(t *testing.T) {
	fakeService := &fakeAuthService{err: errors.New("database error")}
	handler := Login(fakeService)

	request := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(`{"username":"eko","password":"rahasia123"}`))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", recorder.Code)
	}
}

func TestRefreshSuccess(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret")
	fakeService := &fakeAuthService{
		refreshUser:     &model.User{ID: 10, Username: "eko", Role: "viewer"},
		newRefreshToken: "new-refresh-token",
	}
	handler := Refresh(fakeService)

	request := httptest.NewRequest(http.MethodPost, "/api/auth/refresh", strings.NewReader(`{"refresh_token":"old-refresh-token"}`))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", recorder.Code, recorder.Body.String())
	}
	if !fakeService.refreshCalled {
		t.Fatal("expected Refresh service to be called")
	}
	if fakeService.receivedRefreshToken != "old-refresh-token" {
		t.Fatalf("expected old-refresh-token, got %s", fakeService.receivedRefreshToken)
	}
	if !strings.Contains(recorder.Body.String(), `"refresh_token":"new-refresh-token"`) {
		t.Fatalf("expected new refresh token in response, got %s", recorder.Body.String())
	}
	if !strings.Contains(recorder.Body.String(), `"token_type":"Bearer"`) {
		t.Fatalf("expected Bearer token type, got %s", recorder.Body.String())
	}
	if !strings.Contains(recorder.Body.String(), `"access_token":`) {
		t.Fatalf("expected access token in response, got %s", recorder.Body.String())
	}
}

func TestRefreshMethodNotAllowed(t *testing.T) {
	fakeService := &fakeAuthService{}
	handler := Refresh(fakeService)

	request := httptest.NewRequest(http.MethodGet, "/api/auth/refresh", nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", recorder.Code)
	}
	if fakeService.refreshCalled {
		t.Fatal("Refresh service should not be called")
	}
}

func TestRefreshInvalidBody(t *testing.T) {
	fakeService := &fakeAuthService{}
	handler := Refresh(fakeService)

	request := httptest.NewRequest(http.MethodPost, "/api/auth/refresh", strings.NewReader(`{"refresh_token":`))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", recorder.Code)
	}
	if fakeService.refreshCalled {
		t.Fatal("Refresh service should not be called for invalid JSON")
	}
}

func TestRefreshEmptyToken(t *testing.T) {
	fakeService := &fakeAuthService{refreshErr: errors.New("refresh token is required")}
	handler := Refresh(fakeService)

	request := httptest.NewRequest(http.MethodPost, "/api/auth/refresh", strings.NewReader(`{"refresh_token":""}`))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", recorder.Code)
	}
	if !fakeService.refreshCalled {
		t.Fatal("expected Refresh service to be called")
	}
}

func TestRefreshInvalidToken(t *testing.T) {
	fakeService := &fakeAuthService{refreshErr: errors.New("invalid refresh token")}
	handler := Refresh(fakeService)

	request := httptest.NewRequest(http.MethodPost, "/api/auth/refresh", strings.NewReader(`{"refresh_token":"invalid-token"}`))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", recorder.Code)
	}
}
func TestLogoutSuccess(t *testing.T) {
	fakeService := &fakeAuthService{}
	requestBody := `{"refresh_token":"refresh-token-123"}`
	request := httptest.NewRequest(http.MethodPost, "/api/auth/logout", strings.NewReader(requestBody))
	responseRecorder := httptest.NewRecorder()

	Logout(fakeService).ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", responseRecorder.Code)
	}
	if !fakeService.logoutCalled {
		t.Fatal("expected Logout service to be called")
	}
	if fakeService.receivedLogoutToken != "refresh-token-123" {
		t.Fatalf("expected refresh token to be passed to service, got %q", fakeService.receivedLogoutToken)
	}
}

func TestLogoutMethodNotAllowed(t *testing.T) {
	fakeService := &fakeAuthService{}
	request := httptest.NewRequest(http.MethodGet, "/api/auth/logout", nil)
	responseRecorder := httptest.NewRecorder()

	Logout(fakeService).ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", responseRecorder.Code)
	}
	if fakeService.logoutCalled {
		t.Fatal("expected Logout service not to be called")
	}
}

func TestLogoutInvalidBody(t *testing.T) {
	fakeService := &fakeAuthService{}
	request := httptest.NewRequest(http.MethodPost, "/api/auth/logout", strings.NewReader("{invalid-json"))
	responseRecorder := httptest.NewRecorder()

	Logout(fakeService).ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", responseRecorder.Code)
	}
	if fakeService.logoutCalled {
		t.Fatal("expected Logout service not to be called")
	}
}
