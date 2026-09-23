package service

import (
	"context"
	"fmt"

	"BE_GKMI_NTC_GO/internal/model"
	"BE_GKMI_NTC_GO/internal/repository"

	"github.com/jackc/pgx/v5"
)

type JemaatService struct {
	JemaatRepository *repository.JemaatRepository
}

func NewJemaatService(
	jemaatRepository *repository.JemaatRepository,
) *JemaatService {
	return &JemaatService{
		JemaatRepository: jemaatRepository,
	}
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
	if search == "" {
		return nil, &ValidationError{
			Message: "Search is required",
		}
	}

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

func (s *JemaatService) Create(
	ctx context.Context,
	request model.CreateJemaatRequest,
) (*model.Jemaat, error) {
	if message := request.Validate(); message != "" {
		return nil, &ValidationError{
			Message: message,
		}
	}

	jemaat := &model.Jemaat{
		NamaPanggilan:  request.NamaPanggilan,
		NamaLengkap:    request.NamaLengkap,
		JenisKelamin:   request.JenisKelamin,
		TanggalLahir:   request.TanggalLahir,
		Domisili:       request.Domisili,
		StatusJemaat:   request.StatusJemaat,
		StatusDiakonia: request.StatusDiakonia,
		KelompokIbadah: request.KelompokIbadah,
	}

	err := s.JemaatRepository.Create(
		ctx,
		jemaat,
	)
	if err != nil {
		return nil, err
	}

	return jemaat, nil
}

func (s *JemaatService) Update(
	ctx context.Context,
	id int,
	request model.UpdateJemaatRequest,
) (*model.Jemaat, error) {
	if message := request.Validate(); message != "" {
		return nil, &ValidationError{
			Message: message,
		}
	}

	jemaat := &model.Jemaat{
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

	err := s.JemaatRepository.Update(
		ctx,
		jemaat,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, &NotFoundError{
				Message: "Jemaat not found",
			}
		}

		return nil, err
	}

	return jemaat, nil
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
			return &NotFoundError{
				Message: "Jemaat not found",
			}
		}

		return err
	}

	return nil
}

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

type NotFoundError struct {
	Message string
}

func (e *NotFoundError) Error() string {
	return e.Message
}

func NewValidationError(message string) error {
	return &ValidationError{
		Message: message,
	}
}

func ValidationErrorMessage(err error) string {
	if validationErr, ok := err.(*ValidationError); ok {
		return validationErr.Message
	}

	return fmt.Sprintf("%v", err)
}
