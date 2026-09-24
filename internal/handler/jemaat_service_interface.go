package handler

import (
	"context"

	"BE_GKMI_NTC_GO/internal/model"
	"BE_GKMI_NTC_GO/internal/service"
)

type JemaatServiceInterface interface {
	GetAll(ctx context.Context) ([]model.Jemaat, error)
	GetPaginated(ctx context.Context, page int, limit int) (*model.JemaatPaginationResponse, error)
	GetFilteredPaginated(ctx context.Context, search string, jenisKelamin string, statusJemaat string, statusDiakonia string, kelompokIbadah string, sortBy string, sortOrder string, page int, limit int) (*model.JemaatPaginationResponse, error)
	Create(ctx context.Context, request model.CreateJemaatRequest) (*model.Jemaat, error)
	Update(ctx context.Context, id int, request model.UpdateJemaatRequest) (*model.Jemaat, error)
	Delete(ctx context.Context, id int) error
}

var _ JemaatServiceInterface = (*service.JemaatService)(nil)
