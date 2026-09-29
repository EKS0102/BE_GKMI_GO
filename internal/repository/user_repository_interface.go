package repository

import (
	"context"

	"BE_GKMI_NTC_GO/internal/model"
)

type UserRepositoryInterface interface {
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	GetUserByID(ctx context.Context, id int) (*model.User, error)
}

var _ UserRepositoryInterface = (*UserRepository)(nil)
