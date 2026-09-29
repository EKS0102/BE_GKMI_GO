package repository

import (
	"context"
	"time"

	"BE_GKMI_NTC_GO/internal/model"
)

type RefreshTokenRepositoryInterface interface {
	Create(ctx context.Context, refreshToken model.RefreshToken) (*model.RefreshToken, error)
	GetByTokenHash(ctx context.Context, tokenHash string) (*model.RefreshToken, error)
	Revoke(ctx context.Context, tokenHash string, revokedAt time.Time) error
	SetReplacedBy(ctx context.Context, tokenHash string, replacedByTokenHash string) error
	Rotate(ctx context.Context, oldTokenHash string, newRefreshToken model.RefreshToken, revokedAt time.Time) (*model.RefreshToken, error)
}

var _ RefreshTokenRepositoryInterface = (*RefreshTokenRepository)(nil)
