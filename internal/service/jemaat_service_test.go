package service

import (
	"context"
	"errors"
	"testing"

	"BE_GKMI_NTC_GO/internal/model"
	"BE_GKMI_NTC_GO/internal/repository"
)

type fakeJemaatRepository struct {
	items              []model.Jemaat
	total              int
	err                error
	searchItems        []model.Jemaat
	searchTotal        int
	searchErr          error
	lastSearch         string
	lastSearchPage     int
	lastSearchLimit    int
	filteredItems      []model.Jemaat
	filteredTotal      int
	filteredErr        error
	countFilteredErr   error
	lastFilterSearch   string
	lastJenisKelamin   string
	lastStatusJemaat   string
	lastStatusDiakonia string
	lastKelompokIbadah string
	lastSortBy         string
	lastSortOrder      string
	lastFilterPage     int
	lastFilterLimit    int
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

func (f *fakeJemaatRepository) SearchPaginated(ctx context.Context, search string, page int, limit int) ([]model.Jemaat, error) {
	f.lastSearch = search
	f.lastSearchPage = page
	f.lastSearchLimit = limit
	if f.searchErr != nil {
		return nil, f.searchErr
	}
	return f.searchItems, nil
}

func (f *fakeJemaatRepository) CountSearch(context.Context, string) (int, error) {
	if f.searchErr != nil {
		return 0, f.searchErr
	}
	return f.searchTotal, nil
}

func (f *fakeJemaatRepository) GetFilteredPaginated(ctx context.Context, search string, jenisKelamin string, statusJemaat string, statusDiakonia string, kelompokIbadah string, sortBy string, sortOrder string, page int, limit int) ([]model.Jemaat, error) {
	f.lastFilterSearch = search
	f.lastJenisKelamin = jenisKelamin
	f.lastStatusJemaat = statusJemaat
	f.lastStatusDiakonia = statusDiakonia
	f.lastKelompokIbadah = kelompokIbadah
	f.lastSortBy = sortBy
	f.lastSortOrder = sortOrder
	f.lastFilterPage = page
	f.lastFilterLimit = limit
	if f.filteredErr != nil {
		return nil, f.filteredErr
	}
	return f.filteredItems, nil
}

func (f *fakeJemaatRepository) CountFiltered(context.Context, string, string, string, string, string) (int, error) {
	if f.countFilteredErr != nil {
		return 0, f.countFilteredErr
	}
	return f.filteredTotal, nil
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
		items: []model.Jemaat{{ID: 1}, {ID: 2}},
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
	fake := &fakeJemaatRepository{err: expectedErr}
	service := NewJemaatService(fake)
	_, err := service.GetPaginated(context.Background(), 1, 10)
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected repository error, got %v", err)
	}
}

func TestJemaatServiceSearchPaginated(t *testing.T) {
	fake := &fakeJemaatRepository{
		searchItems: []model.Jemaat{{ID: 12, NamaPanggilan: "Maria", NamaLengkap: "Maria Elisabeth"}},
		searchTotal: 2,
	}
	service := NewJemaatService(fake)
	result, err := service.SearchPaginated(context.Background(), "maria", 2, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fake.lastSearch != "maria" {
		t.Fatalf("expected search maria, got %s", fake.lastSearch)
	}
	if fake.lastSearchPage != 2 {
		t.Fatalf("expected page 2, got %d", fake.lastSearchPage)
	}
	if fake.lastSearchLimit != 1 {
		t.Fatalf("expected limit 1, got %d", fake.lastSearchLimit)
	}
	if result.Page != 2 {
		t.Fatalf("expected page 2, got %d", result.Page)
	}
	if result.Limit != 1 {
		t.Fatalf("expected limit 1, got %d", result.Limit)
	}
	if result.Total != 2 {
		t.Fatalf("expected total 2, got %d", result.Total)
	}
	if result.TotalPages != 2 {
		t.Fatalf("expected total pages 2, got %d", result.TotalPages)
	}
	if len(result.Items) != 1 || result.Items[0].ID != 12 {
		t.Fatal("unexpected search items returned")
	}
}

func TestJemaatServiceSearchPaginatedRepositoryError(t *testing.T) {
	expectedErr := errors.New("database error")
	fake := &fakeJemaatRepository{searchErr: expectedErr}
	service := NewJemaatService(fake)
	_, err := service.SearchPaginated(context.Background(), "maria", 1, 10)
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected repository error, got %v", err)
	}
}

func TestJemaatServiceGetFilteredPaginated(t *testing.T) {
	fake := &fakeJemaatRepository{
		filteredItems: []model.Jemaat{{ID: 4, NamaPanggilan: "Sinta", NamaLengkap: "Sinta Maria"}, {ID: 12, NamaPanggilan: "Maria", NamaLengkap: "Maria Elisabeth"}},
		filteredTotal: 2,
	}
	service := NewJemaatService(fake)
	result, err := service.GetFilteredPaginated(context.Background(), "maria", "Perempuan", "Jemaat", "Ya", "Kompak", "nama_lengkap", "asc", 1, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fake.lastFilterSearch != "maria" {
		t.Fatalf("expected search maria, got %s", fake.lastFilterSearch)
	}
	if fake.lastJenisKelamin != "Perempuan" {
		t.Fatalf("expected jenis_kelamin Perempuan, got %s", fake.lastJenisKelamin)
	}
	if fake.lastStatusJemaat != "Jemaat" {
		t.Fatalf("expected status_jemaat Jemaat, got %s", fake.lastStatusJemaat)
	}
	if fake.lastStatusDiakonia != "Ya" {
		t.Fatalf("expected status_diakonia Ya, got %s", fake.lastStatusDiakonia)
	}
	if fake.lastKelompokIbadah != "Kompak" {
		t.Fatalf("expected kelompok_ibadah Kompak, got %s", fake.lastKelompokIbadah)
	}
	if fake.lastSortBy != "nama_lengkap" {
		t.Fatalf("expected sort_by nama_lengkap, got %s", fake.lastSortBy)
	}
	if fake.lastSortOrder != "asc" {
		t.Fatalf("expected sort_order asc, got %s", fake.lastSortOrder)
	}
	if fake.lastFilterPage != 1 {
		t.Fatalf("expected page 1, got %d", fake.lastFilterPage)
	}
	if fake.lastFilterLimit != 10 {
		t.Fatalf("expected limit 10, got %d", fake.lastFilterLimit)
	}
	if result.Total != 2 {
		t.Fatalf("expected total 2, got %d", result.Total)
	}
	if result.TotalPages != 1 {
		t.Fatalf("expected total pages 1, got %d", result.TotalPages)
	}
	if len(result.Items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(result.Items))
	}
	if result.Items[0].ID != 4 || result.Items[1].ID != 12 {
		t.Fatal("unexpected filtered items returned")
	}
}

func TestJemaatServiceGetFilteredPaginatedInvalidJenisKelamin(t *testing.T) {
	service := NewJemaatService(&fakeJemaatRepository{})
	_, err := service.GetFilteredPaginated(context.Background(), "", "Invalid", "", "", "", "id", "asc", 1, 10)
	if err == nil {
		t.Fatal("expected validation error")
	}
	if err.Error() != "Invalid jenis_kelamin" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestJemaatServiceGetFilteredPaginatedInvalidStatusJemaat(t *testing.T) {
	service := NewJemaatService(&fakeJemaatRepository{})
	_, err := service.GetFilteredPaginated(context.Background(), "", "", "Invalid", "", "", "id", "asc", 1, 10)
	if err == nil {
		t.Fatal("expected validation error")
	}
	if err.Error() != "Invalid status_jemaat" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestJemaatServiceGetFilteredPaginatedInvalidStatusDiakonia(t *testing.T) {
	service := NewJemaatService(&fakeJemaatRepository{})
	_, err := service.GetFilteredPaginated(context.Background(), "", "", "", "Invalid", "", "id", "asc", 1, 10)
	if err == nil {
		t.Fatal("expected validation error")
	}
	if err.Error() != "Invalid status_diakonia" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestJemaatServiceGetFilteredPaginatedInvalidKelompokIbadah(t *testing.T) {
	service := NewJemaatService(&fakeJemaatRepository{})
	_, err := service.GetFilteredPaginated(context.Background(), "", "", "", "", "Invalid", "id", "asc", 1, 10)
	if err == nil {
		t.Fatal("expected validation error")
	}
	if err.Error() != "Invalid kelompok_ibadah" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestJemaatServiceGetFilteredPaginatedInvalidSortBy(t *testing.T) {
	service := NewJemaatService(&fakeJemaatRepository{})
	_, err := service.GetFilteredPaginated(context.Background(), "", "", "", "", "", "nama_salah", "asc", 1, 10)
	if err == nil {
		t.Fatal("expected validation error")
	}
	if err.Error() != "Invalid sort_by" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestJemaatServiceGetFilteredPaginatedInvalidSortOrder(t *testing.T) {
	service := NewJemaatService(&fakeJemaatRepository{})
	_, err := service.GetFilteredPaginated(context.Background(), "", "", "", "", "", "id", "random", 1, 10)
	if err == nil {
		t.Fatal("expected validation error")
	}
	if err.Error() != "Invalid sort_order" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestJemaatServiceGetFilteredPaginatedInvalidPage(t *testing.T) {
	service := NewJemaatService(&fakeJemaatRepository{})
	_, err := service.GetFilteredPaginated(context.Background(), "", "", "", "", "", "id", "asc", 0, 10)
	if err == nil {
		t.Fatal("expected validation error")
	}
	if err.Error() != "Page must be greater than 0" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestJemaatServiceGetFilteredPaginatedInvalidLimit(t *testing.T) {
	service := NewJemaatService(&fakeJemaatRepository{})
	_, err := service.GetFilteredPaginated(context.Background(), "", "", "", "", "", "id", "asc", 1, 0)
	if err == nil {
		t.Fatal("expected validation error")
	}
	if err.Error() != "Limit must be greater than 0" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestJemaatServiceGetFilteredPaginatedRepositoryError(t *testing.T) {
	expectedErr := errors.New("database error")
	fake := &fakeJemaatRepository{filteredErr: expectedErr}
	service := NewJemaatService(fake)
	_, err := service.GetFilteredPaginated(context.Background(), "", "", "", "", "", "id", "asc", 1, 10)
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected repository error, got %v", err)
	}
}

func TestJemaatServiceGetFilteredPaginatedCountError(t *testing.T) {
	expectedErr := errors.New("count database error")
	fake := &fakeJemaatRepository{countFilteredErr: expectedErr}
	service := NewJemaatService(fake)
	_, err := service.GetFilteredPaginated(context.Background(), "", "", "", "", "", "id", "asc", 1, 10)
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected count repository error, got %v", err)
	}
}
