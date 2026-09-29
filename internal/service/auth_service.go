package service

import (
	"context"
	"errors"
	"time"

	"BE_GKMI_NTC_GO/internal/model"
	"BE_GKMI_NTC_GO/internal/repository"
	"BE_GKMI_NTC_GO/internal/security"
	"BE_GKMI_NTC_GO/internal/token"
)

type AuthService struct {
	UserRepository         repository.UserRepositoryInterface
	RefreshTokenRepository repository.RefreshTokenRepositoryInterface
}

func NewAuthService(userRepository repository.UserRepositoryInterface, refreshTokenRepository repository.RefreshTokenRepositoryInterface) *AuthService {
	return &AuthService{
		UserRepository:         userRepository,
		RefreshTokenRepository: refreshTokenRepository,
	}
}

func (s *AuthService) Login(ctx context.Context, request model.LoginRequest) (*model.User, string, error) {
	user, err := s.UserRepository.GetUserByUsername(ctx, request.Username)
	if err != nil {
		return nil, "", err
	}
	if !user.IsActive {
		return nil, "", errors.New("user is inactive")
	}
	if !security.VerifyPassword(request.Password, user.PasswordHash) {
		return nil, "", errors.New("invalid credentials")
	}

	refreshToken, err := token.GenerateRefreshToken()
	if err != nil {
		return nil, "", err
	}

	refreshTokenHash := token.HashRefreshToken(refreshToken)
	now := time.Now()
	refreshTokenModel := model.RefreshToken{
		UserID:    user.ID,
		TokenHash: refreshTokenHash,
		ExpiresAt: now.Add(7 * 24 * time.Hour),
		CreatedAt: now,
	}

	_, err = s.RefreshTokenRepository.Create(ctx, refreshTokenModel)
	if err != nil {
		return nil, "", err
	}

	return user, refreshToken, nil
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	if refreshToken == "" {
		return errors.New("refresh token is required")
	}

	tokenHash := token.HashRefreshToken(refreshToken)
	storedToken, err := s.RefreshTokenRepository.GetByTokenHash(ctx, tokenHash)
	if err != nil {
		return errors.New("invalid refresh token")
	}

	if storedToken.RevokedAt != nil {
		return errors.New("refresh token has already been revoked")
	}

	if err := s.RefreshTokenRepository.Revoke(ctx, tokenHash, time.Now()); err != nil {
		return err
	}

	return nil
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (*model.User, string, error) {
	if refreshToken == "" {
		return nil, "", errors.New("refresh token is required")
	}

	tokenHash := token.HashRefreshToken(refreshToken)
	storedToken, err := s.RefreshTokenRepository.GetByTokenHash(ctx, tokenHash)
	if err != nil {
		return nil, "", errors.New("invalid refresh token")
	}

	if storedToken.RevokedAt != nil {
		return nil, "", errors.New("refresh token has been revoked")
	}

	if !storedToken.ExpiresAt.After(time.Now()) {
		return nil, "", errors.New("refresh token has expired")
	}

	user, err := s.UserRepository.GetUserByID(ctx, storedToken.UserID)
	if err != nil {
		return nil, "", err
	}

	if !user.IsActive {
		return nil, "", errors.New("user is inactive")
	}

	newRefreshToken, err := token.GenerateRefreshToken()
	if err != nil {
		return nil, "", err
	}

	newRefreshTokenHash := token.HashRefreshToken(newRefreshToken)
	now := time.Now()
	newRefreshTokenModel := model.RefreshToken{
		UserID:    user.ID,
		TokenHash: newRefreshTokenHash,
		ExpiresAt: now.Add(7 * 24 * time.Hour),
		CreatedAt: now,
	}

	_, err = s.RefreshTokenRepository.Rotate(ctx, tokenHash, newRefreshTokenModel, now)
	if err != nil {
		return nil, "", err
	}

	return user, newRefreshToken, nil
}
