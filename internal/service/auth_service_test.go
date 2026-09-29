package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"BE_GKMI_NTC_GO/internal/model"
	"BE_GKMI_NTC_GO/internal/security"
	"BE_GKMI_NTC_GO/internal/token"
)

type fakeUserRepository struct {
	user     *model.User
	err      error
	username string
	userID   int
}

func (f *fakeUserRepository) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	f.username = username
	if f.err != nil {
		return nil, f.err
	}
	return f.user, nil
}

func (f *fakeUserRepository) GetUserByID(ctx context.Context, id int) (*model.User, error) {
	f.userID = id
	if f.err != nil {
		return nil, f.err
	}
	if f.user == nil || f.user.ID != id {
		return nil, errors.New("user not found")
	}
	return f.user, nil
}

type fakeRefreshTokenRepository struct {
	refreshTokens     map[string]*model.RefreshToken
	createdToken      *model.RefreshToken
	getTokenHash      string
	revokedTokenHash  string
	revokedAt         time.Time
	replacedTokenHash string
	replacedByHash    string
	createErr         error
	getErr            error
	revokeErr         error
	replaceErr        error
}

func newFakeRefreshTokenRepository() *fakeRefreshTokenRepository {
	return &fakeRefreshTokenRepository{
		refreshTokens: make(map[string]*model.RefreshToken),
	}
}

func (f *fakeRefreshTokenRepository) Create(ctx context.Context, refreshToken model.RefreshToken) (*model.RefreshToken, error) {
	if f.createErr != nil {
		return nil, f.createErr
	}
	copyToken := refreshToken
	f.createdToken = &copyToken
	f.refreshTokens[refreshToken.TokenHash] = &copyToken
	return &copyToken, nil
}

func (f *fakeRefreshTokenRepository) GetByTokenHash(ctx context.Context, tokenHash string) (*model.RefreshToken, error) {
	f.getTokenHash = tokenHash
	if f.getErr != nil {
		return nil, f.getErr
	}
	refreshToken, ok := f.refreshTokens[tokenHash]
	if !ok {
		return nil, errors.New("refresh token not found")
	}
	copyToken := *refreshToken
	return &copyToken, nil
}

func (f *fakeRefreshTokenRepository) Revoke(ctx context.Context, tokenHash string, revokedAt time.Time) error {
	if f.revokeErr != nil {
		return f.revokeErr
	}
	f.revokedTokenHash = tokenHash
	f.revokedAt = revokedAt
	if refreshToken, ok := f.refreshTokens[tokenHash]; ok {
		refreshToken.RevokedAt = &revokedAt
	}
	return nil
}

func (f *fakeRefreshTokenRepository) SetReplacedBy(ctx context.Context, tokenHash string, replacedByTokenHash string) error {
	if f.replaceErr != nil {
		return f.replaceErr
	}
	f.replacedTokenHash = tokenHash
	f.replacedByHash = replacedByTokenHash
	if refreshToken, ok := f.refreshTokens[tokenHash]; ok {
		refreshToken.ReplacedByTokenHash = replacedByTokenHash
	}
	return nil
}

func (f *fakeRefreshTokenRepository) Rotate(ctx context.Context, oldTokenHash string, newRefreshToken model.RefreshToken, revokedAt time.Time) (*model.RefreshToken, error) {
	if f.createErr != nil {
		return nil, f.createErr
	}
	if f.revokeErr != nil {
		return nil, f.revokeErr
	}
	if f.replaceErr != nil {
		return nil, f.replaceErr
	}

	copyToken := newRefreshToken
	f.createdToken = &copyToken
	f.refreshTokens[newRefreshToken.TokenHash] = &copyToken
	f.revokedTokenHash = oldTokenHash
	f.revokedAt = revokedAt
	f.replacedTokenHash = oldTokenHash
	f.replacedByHash = newRefreshToken.TokenHash

	if refreshToken, ok := f.refreshTokens[oldTokenHash]; ok {
		refreshToken.RevokedAt = &revokedAt
		refreshToken.ReplacedByTokenHash = newRefreshToken.TokenHash
	}

	return &copyToken, nil
}
func createTestUser(t *testing.T, active bool) *model.User {
	passwordHash, err := security.HashPassword("rahasia123")
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}
	return &model.User{
		ID:           10,
		Username:     "eko",
		PasswordHash: passwordHash,
		Role:         "viewer",
		IsActive:     active,
	}
}

func createTestRefreshToken(userID int, tokenValue string, expiresAt time.Time) (*model.RefreshToken, string) {
	tokenHash := token.HashRefreshToken(tokenValue)
	return &model.RefreshToken{
		ID:        1,
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now().Add(-1 * time.Hour),
	}, tokenHash
}

func TestAuthServiceLoginSuccess(t *testing.T) {
	userRepository := &fakeUserRepository{user: createTestUser(t, true)}
	refreshTokenRepository := newFakeRefreshTokenRepository()
	authService := NewAuthService(userRepository, refreshTokenRepository)

	user, refreshToken, err := authService.Login(context.Background(), model.LoginRequest{Username: "eko", Password: "rahasia123"})
	if err != nil {
		t.Fatalf("Login() error = %v", err)
	}
	if user == nil || user.Username != "eko" {
		t.Fatal("expected eko user")
	}
	if refreshToken == "" {
		t.Fatal("expected refresh token")
	}
	if refreshTokenRepository.createdToken == nil {
		t.Fatal("expected refresh token to be stored")
	}
	if refreshTokenRepository.createdToken.TokenHash == refreshToken {
		t.Fatal("refresh token must not be stored as plaintext")
	}
	if refreshTokenRepository.createdToken.ExpiresAt.Before(time.Now().Add(6 * 24 * time.Hour)) {
		t.Fatal("expected refresh token expiry around 7 days")
	}
}

func TestAuthServiceLoginRepositoryError(t *testing.T) {
	userRepository := &fakeUserRepository{err: errors.New("database error")}
	refreshTokenRepository := newFakeRefreshTokenRepository()
	authService := NewAuthService(userRepository, refreshTokenRepository)

	user, refreshToken, err := authService.Login(context.Background(), model.LoginRequest{Username: "eko", Password: "rahasia123"})
	if err == nil || user != nil || refreshToken != "" {
		t.Fatal("expected login repository error")
	}
}

func TestAuthServiceLoginInactiveUser(t *testing.T) {
	userRepository := &fakeUserRepository{user: createTestUser(t, false)}
	refreshTokenRepository := newFakeRefreshTokenRepository()
	authService := NewAuthService(userRepository, refreshTokenRepository)

	user, refreshToken, err := authService.Login(context.Background(), model.LoginRequest{Username: "eko", Password: "rahasia123"})
	if err == nil || user != nil || refreshToken != "" {
		t.Fatal("expected inactive user error")
	}
	if refreshTokenRepository.createdToken != nil {
		t.Fatal("refresh token must not be created")
	}
}

func TestAuthServiceLoginInvalidPassword(t *testing.T) {
	userRepository := &fakeUserRepository{user: createTestUser(t, true)}
	refreshTokenRepository := newFakeRefreshTokenRepository()
	authService := NewAuthService(userRepository, refreshTokenRepository)

	user, refreshToken, err := authService.Login(context.Background(), model.LoginRequest{Username: "eko", Password: "salah"})
	if err == nil || user != nil || refreshToken != "" {
		t.Fatal("expected invalid password error")
	}
	if refreshTokenRepository.createdToken != nil {
		t.Fatal("refresh token must not be created")
	}
}

func TestAuthServiceLoginRefreshTokenRepositoryError(t *testing.T) {
	userRepository := &fakeUserRepository{user: createTestUser(t, true)}
	refreshTokenRepository := newFakeRefreshTokenRepository()
	refreshTokenRepository.createErr = errors.New("refresh token database error")
	authService := NewAuthService(userRepository, refreshTokenRepository)

	user, refreshToken, err := authService.Login(context.Background(), model.LoginRequest{Username: "eko", Password: "rahasia123"})
	if err == nil || user != nil || refreshToken != "" {
		t.Fatal("expected refresh token repository error")
	}
}

func TestAuthServiceRefreshSuccess(t *testing.T) {
	userRepository := &fakeUserRepository{user: createTestUser(t, true)}
	refreshTokenRepository := newFakeRefreshTokenRepository()
	oldToken := "old-refresh-token"
	storedToken, oldHash := createTestRefreshToken(10, oldToken, time.Now().Add(7*24*time.Hour))
	refreshTokenRepository.refreshTokens[oldHash] = storedToken
	authService := NewAuthService(userRepository, refreshTokenRepository)

	user, newToken, err := authService.Refresh(context.Background(), oldToken)
	if err != nil {
		t.Fatalf("Refresh() error = %v", err)
	}
	if user == nil || user.ID != 10 {
		t.Fatal("expected user 10")
	}
	if newToken == "" {
		t.Fatal("expected new refresh token")
	}
	if newToken == oldToken {
		t.Fatal("expected token rotation")
	}
	if refreshTokenRepository.getTokenHash != oldHash {
		t.Fatal("expected lookup by old token hash")
	}
	if refreshTokenRepository.createdToken == nil {
		t.Fatal("expected new refresh token to be created")
	}
	newHash := token.HashRefreshToken(newToken)
	if refreshTokenRepository.createdToken.TokenHash != newHash {
		t.Fatal("stored hash does not match new refresh token")
	}
	if refreshTokenRepository.createdToken.UserID != 10 {
		t.Fatal("new refresh token belongs to wrong user")
	}
	if refreshTokenRepository.revokedTokenHash != oldHash {
		t.Fatal("old refresh token was not revoked")
	}
	if refreshTokenRepository.revokedAt.IsZero() {
		t.Fatal("expected revoked_at")
	}
	if refreshTokenRepository.replacedTokenHash != oldHash {
		t.Fatal("expected old token to be marked as replaced")
	}
	if refreshTokenRepository.replacedByHash != newHash {
		t.Fatal("expected replaced_by_token_hash to point to new token")
	}
}

func TestAuthServiceRefreshEmptyToken(t *testing.T) {
	userRepository := &fakeUserRepository{user: createTestUser(t, true)}
	refreshTokenRepository := newFakeRefreshTokenRepository()
	authService := NewAuthService(userRepository, refreshTokenRepository)

	user, newToken, err := authService.Refresh(context.Background(), "")
	if err == nil {
		t.Fatal("expected error")
	}
	if user != nil || newToken != "" {
		t.Fatal("expected empty result")
	}
}

func TestAuthServiceRefreshTokenNotFound(t *testing.T) {
	userRepository := &fakeUserRepository{user: createTestUser(t, true)}
	refreshTokenRepository := newFakeRefreshTokenRepository()
	refreshTokenRepository.getErr = errors.New("refresh token not found")
	authService := NewAuthService(userRepository, refreshTokenRepository)

	user, newToken, err := authService.Refresh(context.Background(), "unknown-token")
	if err == nil {
		t.Fatal("expected error")
	}
	if user != nil || newToken != "" {
		t.Fatal("expected empty result")
	}
}

func TestAuthServiceRefreshRevokedToken(t *testing.T) {
	userRepository := &fakeUserRepository{user: createTestUser(t, true)}
	refreshTokenRepository := newFakeRefreshTokenRepository()
	oldToken := "revoked-refresh-token"
	storedToken, oldHash := createTestRefreshToken(10, oldToken, time.Now().Add(7*24*time.Hour))
	revokedAt := time.Now().Add(-1 * time.Hour)
	storedToken.RevokedAt = &revokedAt
	refreshTokenRepository.refreshTokens[oldHash] = storedToken
	authService := NewAuthService(userRepository, refreshTokenRepository)

	user, newToken, err := authService.Refresh(context.Background(), oldToken)
	if err == nil {
		t.Fatal("expected revoked token error")
	}
	if user != nil || newToken != "" {
		t.Fatal("expected empty result")
	}
}

func TestAuthServiceRefreshExpiredToken(t *testing.T) {
	userRepository := &fakeUserRepository{user: createTestUser(t, true)}
	refreshTokenRepository := newFakeRefreshTokenRepository()
	oldToken := "expired-refresh-token"
	storedToken, oldHash := createTestRefreshToken(10, oldToken, time.Now().Add(-1*time.Minute))
	refreshTokenRepository.refreshTokens[oldHash] = storedToken
	authService := NewAuthService(userRepository, refreshTokenRepository)

	user, newToken, err := authService.Refresh(context.Background(), oldToken)
	if err == nil {
		t.Fatal("expected expired token error")
	}
	if user != nil || newToken != "" {
		t.Fatal("expected empty result")
	}
}

func TestAuthServiceRefreshInactiveUser(t *testing.T) {
	userRepository := &fakeUserRepository{user: createTestUser(t, false)}
	refreshTokenRepository := newFakeRefreshTokenRepository()
	oldToken := "inactive-user-token"
	storedToken, oldHash := createTestRefreshToken(10, oldToken, time.Now().Add(7*24*time.Hour))
	refreshTokenRepository.refreshTokens[oldHash] = storedToken
	authService := NewAuthService(userRepository, refreshTokenRepository)

	user, newToken, err := authService.Refresh(context.Background(), oldToken)
	if err == nil {
		t.Fatal("expected inactive user error")
	}
	if user != nil || newToken != "" {
		t.Fatal("expected empty result")
	}
	if refreshTokenRepository.createdToken != nil {
		t.Fatal("new refresh token must not be created")
	}
}

func TestAuthServiceRefreshUserRepositoryError(t *testing.T) {
	userRepository := &fakeUserRepository{err: errors.New("user database error")}
	refreshTokenRepository := newFakeRefreshTokenRepository()
	oldToken := "user-error-token"
	storedToken, oldHash := createTestRefreshToken(10, oldToken, time.Now().Add(7*24*time.Hour))
	refreshTokenRepository.refreshTokens[oldHash] = storedToken
	authService := NewAuthService(userRepository, refreshTokenRepository)

	user, newToken, err := authService.Refresh(context.Background(), oldToken)
	if err == nil {
		t.Fatal("expected user repository error")
	}
	if user != nil || newToken != "" {
		t.Fatal("expected empty result")
	}
}
