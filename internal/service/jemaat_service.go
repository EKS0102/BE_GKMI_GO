package service

import (
	"context"
	"fmt"

	"BE_GKMI_NTC_GO/internal/model"
	"BE_GKMI_NTC_GO/internal/repository"

	"github.com/jackc/pgx/v5"
)

type JemaatService struct {
	JemaatRepository repository.JemaatRepositoryInterface
}

func NewJemaatService(
	jemaatRepository repository.JemaatRepositoryInterface,
) *JemaatService {
	return &JemaatService{
		JemaatRepository: jemaatRepository,
	}
}

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

func (s *JemaatService) GetAll(
	ctx context.Context,
) ([]model.Jemaat, error) {
	return s.JemaatRepository.GetAll(ctx)
}

func (s *JemaatService) GetPaginated(
	ctx context.Context,
	page int,
	limit int,
) (*model.JemaatPaginationResponse, error) {

	if page < 1 {
		return nil, &ValidationError{
			Message: "Page must be greater than 0",
		}
	}

	if limit < 1 {
		return nil, &ValidationError{
			Message: "Limit must be greater than 0",
		}
	}

	if limit > 100 {
		return nil, &ValidationError{
			Message: "Limit must not exceed 100",
		}
	}

	jemaatList, err := s.JemaatRepository.GetPaginated(
		ctx,
		page,
		limit,
	)
	if err != nil {
		return nil, err
	}

	total, err := s.JemaatRepository.Count(ctx)
	if err != nil {
		return nil, err
	}

	totalPages := 0

	if total > 0 {
		totalPages = (total + limit - 1) / limit
	}

	return &model.JemaatPaginationResponse{
		Items:      jemaatList,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

func (s *JemaatService) SearchPaginated(
	ctx context.Context,
	search string,
	page int,
	limit int,
) (*model.JemaatPaginationResponse, error) {

	if page < 1 {
		return nil, &ValidationError{
			Message: "Page must be greater than 0",
		}
	}

	if limit < 1 {
		return nil, &ValidationError{
			Message: "Limit must be greater than 0",
		}
	}

	if limit > 100 {
		return nil, &ValidationError{
			Message: "Limit must not exceed 100",
		}
	}

	jemaatList, err := s.JemaatRepository.SearchPaginated(
		ctx,
		search,
		page,
		limit,
	)
	if err != nil {
		return nil, err
	}

	total, err := s.JemaatRepository.CountSearch(
		ctx,
		search,
	)
	if err != nil {
		return nil, err
	}

	totalPages := 0

	if total > 0 {
		totalPages = (total + limit - 1) / limit
	}

	return &model.JemaatPaginationResponse{
		Items:      jemaatList,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

func (s *JemaatService) GetFilteredPaginated(
	ctx context.Context,
	search string,
	jenisKelamin string,
	statusJemaat string,
	statusDiakonia string,
	kelompokIbadah string,
	sortBy string,
	sortOrder string,
	page int,
	limit int,
) (*model.JemaatPaginationResponse, error) {

	if page < 1 {
		return nil, &ValidationError{
			Message: "Page must be greater than 0",
		}
	}

	if limit < 1 {
		return nil, &ValidationError{
			Message: "Limit must be greater than 0",
		}
	}

	if limit > 100 {
		return nil, &ValidationError{
			Message: "Limit must not exceed 100",
		}
	}

	if !model.SortBy(sortBy).IsValid() {
		return nil, &ValidationError{
			Message: "Invalid sort_by",
		}
	}

	if !model.SortOrder(sortOrder).IsValid() {
		return nil, &ValidationError{
			Message: "Invalid sort_order",
		}
	}

	if jenisKelamin != "" && !model.JenisKelamin(jenisKelamin).IsValid() {
		return nil, &ValidationError{
			Message: "Invalid jenis_kelamin",
		}
	}

	if statusJemaat != "" && !model.StatusJemaat(statusJemaat).IsValid() {
		return nil, &ValidationError{
			Message: "Invalid status_jemaat",
		}
	}

	if statusDiakonia != "" && !model.StatusDiakonia(statusDiakonia).IsValid() {
		return nil, &ValidationError{
			Message: "Invalid status_diakonia",
		}
	}

	if kelompokIbadah != "" && !model.KelompokIbadah(kelompokIbadah).IsValid() {
		return nil, &ValidationError{
			Message: "Invalid kelompok_ibadah",
		}
	}
	jemaatList, err := s.JemaatRepository.GetFilteredPaginated(
		ctx,
		search,
		jenisKelamin,
		statusJemaat,
		statusDiakonia,
		kelompokIbadah,
		sortBy,
		sortOrder,
		page,
		limit,
	)

	if err != nil {
		return nil, err
	}

	total, err := s.JemaatRepository.CountFiltered(
		ctx,
		search,
		jenisKelamin,
		statusJemaat,
		statusDiakonia,
		kelompokIbadah,
	)

	if err != nil {
		return nil, err
	}

	totalPages := 0

	if total > 0 {
		totalPages = (total + limit - 1) / limit
	}

	return &model.JemaatPaginationResponse{
		Items:      jemaatList,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

func (s *JemaatService) Create(
	ctx context.Context,
	request model.CreateJemaatRequest,
) (*model.Jemaat, error) {

	jemaat := model.Jemaat{
		NamaPanggilan:  request.NamaPanggilan,
		NamaLengkap:    request.NamaLengkap,
		JenisKelamin:   request.JenisKelamin,
		TanggalLahir:   request.TanggalLahir,
		Domisili:       request.Domisili,
		StatusJemaat:   request.StatusJemaat,
		StatusDiakonia: request.StatusDiakonia,
		KelompokIbadah: request.KelompokIbadah,
	}

	result, err := s.JemaatRepository.Create(
		ctx,
		jemaat,
	)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *JemaatService) Update(
	ctx context.Context,
	id int,
	request model.UpdateJemaatRequest,
) (*model.Jemaat, error) {

	jemaat := model.Jemaat{
		ID:             id,
		NamaPanggilan:  request.NamaPanggilan,
		NamaLengkap:    request.NamaLengkap,
		JenisKelamin:   request.JenisKelamin,
		TanggalLahir:   request.TanggalLahir,
		Domisili:       request.Domisili,
		StatusJemaat:   request.StatusJemaat,
		StatusDiakonia: request.StatusDiakonia,
		KelompokIbadah: request.KelompokIbadah,
	}

	result, err := s.JemaatRepository.Update(
		ctx,
		id,
		jemaat,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("jemaat not found")
		}

		return nil, err
	}

	return result, nil
}

func (s *JemaatService) Delete(
	ctx context.Context,
	id int,
) error {

	err := s.JemaatRepository.Delete(
		ctx,
		id,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("jemaat not found")
		}

		return err
	}

	return nil
}
