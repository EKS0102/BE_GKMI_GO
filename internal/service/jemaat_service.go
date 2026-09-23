package service

import (
	"context"

	"BE_GKMI_NTC_GO/internal/model"
	"BE_GKMI_NTC_GO/internal/repository"
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

func (s *JemaatService) Create(
	ctx context.Context,
	request model.CreateJemaatRequest,
) (*model.Jemaat, error) {
	if errMessage := request.Validate(); errMessage != "" {
		return nil, &ValidationError{
			Message: errMessage,
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

	err := s.JemaatRepository.Create(ctx, jemaat)
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
	if errMessage := request.Validate(); errMessage != "" {
		return nil, &ValidationError{
			Message: errMessage,
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

	err := s.JemaatRepository.Update(ctx, jemaat)
	if err != nil {
		return nil, err
	}

	return jemaat, nil
}

func (s *JemaatService) Delete(
	ctx context.Context,
	id int,
) error {
	return s.JemaatRepository.Delete(ctx, id)
}

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}
