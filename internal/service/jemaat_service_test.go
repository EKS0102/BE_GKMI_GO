package service

import (
	"context"
	"errors"
	"testing"

	"BE_GKMI_NTC_GO/internal/model"
	"BE_GKMI_NTC_GO/internal/repository"
)

type fakeJemaatRepository struct {
	items []model.Jemaat
	total int
	err   error
}

func (f *fakeJemaatRepository) GetAll(context.Context) ([]model.Jemaat, error) { return nil, nil }

func (f *fakeJemaatRepository) GetPaginated(context.Context, int, int) ([]model.Jemaat, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.items, nil
}

func (f *fakeJemaatRepository) Count(context.Context) (int, error) {
	if f.err != nil {
		return 0, f.err
	}
	return f.total, nil
}

func (f *fakeJemaatRepository) SearchPaginated(context.Context, string, int, int) ([]model.Jemaat, error) {
	return nil, nil
}

func (f *fakeJemaatRepository) CountSearch(context.Context, string) (int, error) { return 0, nil }

func (f *fakeJemaatRepository) GetFilteredPaginated(context.Context, string, string, string, string, string, string, string, int, int) ([]model.Jemaat, error) {
	return nil, nil
}

func (f *fakeJemaatRepository) CountFiltered(context.Context, string, string, string, string, string) (int, error) {
	return 0, nil
}

func (f *fakeJemaatRepository) Create(context.Context, model.Jemaat) (*model.Jemaat, error) {
	return nil, nil
}

func (f *fakeJemaatRepository) Update(context.Context, int, model.Jemaat) (*model.Jemaat, error) {
	return nil, nil
}

func (f *fakeJemaatRepository) Delete(context.Context, int) error { return nil }

var _ repository.JemaatRepositoryInterface = (*fakeJemaatRepository)(nil)

func TestJemaatServiceGetPaginatedInvalidPage(t *testing.T) {
	service := NewJemaatService(nil)

	_, err := service.GetPaginated(context.Background(), 0, 10)

	if err == nil {
		t.Fatal("expected error for page < 1")
	}

	if err.Error() != "Page must be greater than 0" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestJemaatServiceGetPaginated(t *testing.T) {
	fake := &fakeJemaatRepository{
		items: []model.Jemaat{
			{ID: 1},
			{ID: 2},
		},
		total: 5,
	}

	service := NewJemaatService(fake)

	result, err := service.GetPaginated(context.Background(), 2, 2)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Page != 2 {
		t.Fatalf("expected page 2, got %d", result.Page)
	}

	if result.Limit != 2 {
		t.Fatalf("expected limit 2, got %d", result.Limit)
	}

	if result.Total != 5 {
		t.Fatalf("expected total 5, got %d", result.Total)
	}

	if result.TotalPages != 3 {
		t.Fatalf("expected total pages 3, got %d", result.TotalPages)
	}

	if len(result.Items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(result.Items))
	}

	if result.Items[0].ID != 1 || result.Items[1].ID != 2 {
		t.Fatal("unexpected items returned")
	}
}

func TestJemaatServiceGetPaginatedRepositoryError(t *testing.T) {
	expectedErr := errors.New("database error")

	fake := &fakeJemaatRepository{
		err: expectedErr,
	}

	service := NewJemaatService(fake)

	_, err := service.GetPaginated(context.Background(), 1, 10)

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected repository error, got %v", err)
	}
}
