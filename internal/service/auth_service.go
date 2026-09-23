package service

import (
	"context"
	"errors"

	"BE_GKMI_NTC_GO/internal/model"
	"BE_GKMI_NTC_GO/internal/repository"
	"BE_GKMI_NTC_GO/internal/security"
)

type AuthService struct {
	UserRepository *repository.UserRepository
}

func NewAuthService(
	userRepository *repository.UserRepository,
) *AuthService {
	return &AuthService{
		UserRepository: userRepository,
	}
}

func (s *AuthService) Login(
	ctx context.Context,
	request model.LoginRequest,
) (*model.User, error) {
	user, err := s.UserRepository.GetUserByUsername(
		ctx,
		request.Username,
	)
	if err != nil {
		return nil, err
	}

	if !user.IsActive {
		return nil, errors.New("user is inactive")
	}

	if !security.VerifyPassword(
		request.Password,
		user.PasswordHash,
	) {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}