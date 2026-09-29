package handler

import (
	"context"

	"BE_GKMI_NTC_GO/internal/model"
	"BE_GKMI_NTC_GO/internal/service"
)

type AuthServiceInterface interface {
	Login(ctx context.Context, request model.LoginRequest) (*model.User, string, error)
	Refresh(ctx context.Context, refreshToken string) (*model.User, string, error)
	Logout(ctx context.Context, refreshToken string) error
}

var _ AuthServiceInterface = (*service.AuthService)(nil)
