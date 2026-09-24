package repository

import (
	"context"

	"BE_GKMI_NTC_GO/internal/model"
)

type JemaatRepositoryInterface interface {
	GetAll(ctx context.Context) ([]model.Jemaat, error)
	GetPaginated(ctx context.Context, page int, limit int) ([]model.Jemaat, error)
	Count(ctx context.Context) (int, error)
	SearchPaginated(ctx context.Context, search string, page int, limit int) ([]model.Jemaat, error)
	CountSearch(ctx context.Context, search string) (int, error)
	GetFilteredPaginated(ctx context.Context, search string, jenisKelamin string, statusJemaat string, statusDiakonia string, kelompokIbadah string, sortBy string, sortOrder string, page int, limit int) ([]model.Jemaat, error)
	CountFiltered(ctx context.Context, search string, jenisKelamin string, statusJemaat string, statusDiakonia string, kelompokIbadah string) (int, error)
	Create(ctx context.Context, jemaat model.Jemaat) (*model.Jemaat, error)
	Update(ctx context.Context, id int, jemaat model.Jemaat) (*model.Jemaat, error)
	Delete(ctx context.Context, id int) error
}

var _ JemaatRepositoryInterface = (*JemaatRepository)(nil)
