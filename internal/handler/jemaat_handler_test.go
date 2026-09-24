package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"BE_GKMI_NTC_GO/internal/model"
	"BE_GKMI_NTC_GO/internal/service"
)

type fakeJemaatService struct {
	getAllResult              []model.Jemaat
	getAllErr                 error
	getPaginatedResult        *model.JemaatPaginationResponse
	getPaginatedErr           error
	getPaginatedPage          int
	getPaginatedLimit         int
	getFilteredResult         *model.JemaatPaginationResponse
	getFilteredErr            error
	getFilteredSearch         string
	getFilteredJenisKelamin   string
	getFilteredStatusJemaat   string
	getFilteredStatusDiakonia string
	getFilteredKelompokIbadah string
	getFilteredSortBy         string
	getFilteredSortOrder      string
	getFilteredPage           int
	getFilteredLimit          int
	createResult              *model.Jemaat
	createErr                 error
	createRequest             model.CreateJemaatRequest
	updateResult              *model.Jemaat
	updateErr                 error
	updateID                  int
	updateRequest             model.UpdateJemaatRequest
	deleteErr                 error
	deleteID                  int
}

func (f *fakeJemaatService) GetAll(_ context.Context) ([]model.Jemaat, error) {
	return f.getAllResult, f.getAllErr
}

func (f *fakeJemaatService) GetPaginated(_ context.Context, page int, limit int) (*model.JemaatPaginationResponse, error) {
	f.getPaginatedPage = page
	f.getPaginatedLimit = limit
	return f.getPaginatedResult, f.getPaginatedErr
}

func (f *fakeJemaatService) GetFilteredPaginated(_ context.Context, search string, jenisKelamin string, statusJemaat string, statusDiakonia string, kelompokIbadah string, sortBy string, sortOrder string, page int, limit int) (*model.JemaatPaginationResponse, error) {
	f.getFilteredSearch = search
	f.getFilteredJenisKelamin = jenisKelamin
	f.getFilteredStatusJemaat = statusJemaat
	f.getFilteredStatusDiakonia = statusDiakonia
	f.getFilteredKelompokIbadah = kelompokIbadah
	f.getFilteredSortBy = sortBy
	f.getFilteredSortOrder = sortOrder
	f.getFilteredPage = page
	f.getFilteredLimit = limit
	return f.getFilteredResult, f.getFilteredErr
}

func (f *fakeJemaatService) Create(_ context.Context, request model.CreateJemaatRequest) (*model.Jemaat, error) {
	f.createRequest = request
	return f.createResult, f.createErr
}

func (f *fakeJemaatService) Update(_ context.Context, id int, request model.UpdateJemaatRequest) (*model.Jemaat, error) {
	f.updateID = id
	f.updateRequest = request
	return f.updateResult, f.updateErr
}

func (f *fakeJemaatService) Delete(_ context.Context, id int) error {
	f.deleteID = id
	return f.deleteErr
}

var _ JemaatServiceInterface = (*fakeJemaatService)(nil)

func TestJemaatGetAll(t *testing.T) {
	fakeService := &fakeJemaatService{getAllResult: []model.Jemaat{{ID: 1, NamaPanggilan: "Budi", NamaLengkap: "Budi Santoso"}, {ID: 2, NamaPanggilan: "Sari", NamaLengkap: "Sari Maria"}}}
	handler := Jemaat(nil, fakeService)
	request := httptest.NewRequest(http.MethodGet, "/api/jemaat", nil)
	recorder := httptest.NewRecorder()
	handler(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}
	var result []model.Jemaat
	if err := json.NewDecoder(recorder.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 jemaat, got %d", len(result))
	}
	if result[0].NamaLengkap != "Budi Santoso" {
		t.Fatalf("expected Budi Santoso, got %s", result[0].NamaLengkap)
	}
}

func TestJemaatGetAllServiceError(t *testing.T) {
	fakeService := &fakeJemaatService{getAllErr: errors.New("database error")}
	handler := Jemaat(nil, fakeService)
	request := httptest.NewRequest(http.MethodGet, "/api/jemaat", nil)
	recorder := httptest.NewRecorder()
	handler(recorder, request)
	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", recorder.Code)
	}
}

func TestJemaatGetPaginated(t *testing.T) {
	fakeService := &fakeJemaatService{getPaginatedResult: &model.JemaatPaginationResponse{Items: []model.Jemaat{{ID: 3, NamaLengkap: "Andi Saputra"}}, Page: 2, Limit: 5, Total: 12, TotalPages: 3}}
	handler := Jemaat(nil, fakeService)
	request := httptest.NewRequest(http.MethodGet, "/api/jemaat?page=2&limit=5", nil)
	recorder := httptest.NewRecorder()
	handler(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}
	if fakeService.getPaginatedPage != 2 {
		t.Fatalf("expected page 2, got %d", fakeService.getPaginatedPage)
	}
	if fakeService.getPaginatedLimit != 5 {
		t.Fatalf("expected limit 5, got %d", fakeService.getPaginatedLimit)
	}
}

func TestJemaatGetPaginatedInvalidPage(t *testing.T) {
	fakeService := &fakeJemaatService{}
	handler := Jemaat(nil, fakeService)
	request := httptest.NewRequest(http.MethodGet, "/api/jemaat?page=abc&limit=5", nil)
	recorder := httptest.NewRecorder()
	handler(recorder, request)
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", recorder.Code)
	}
}

func TestJemaatGetPaginatedValidationError(t *testing.T) {
	fakeService := &fakeJemaatService{getPaginatedErr: &service.ValidationError{Message: "page must be greater than or equal to 1"}}
	handler := Jemaat(nil, fakeService)
	request := httptest.NewRequest(http.MethodGet, "/api/jemaat?page=0&limit=5", nil)
	recorder := httptest.NewRecorder()
	handler(recorder, request)
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", recorder.Code)
	}
}

func TestJemaatGetFilteredPaginated(t *testing.T) {
	fakeService := &fakeJemaatService{getFilteredResult: &model.JemaatPaginationResponse{Items: []model.Jemaat{{ID: 4, NamaLengkap: "Sinta Maria"}}, Page: 1, Limit: 2, Total: 1, TotalPages: 1}}
	handler := Jemaat(nil, fakeService)
	request := httptest.NewRequest(http.MethodGet, "/api/jemaat?search=maria&jenis_kelamin=Perempuan&status_jemaat=Jemaat&status_diakonia=Ya&kelompok_ibadah=Youth&sort_by=nama_lengkap&sort_order=desc&page=1&limit=2", nil)
	recorder := httptest.NewRecorder()
	handler(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}
	if fakeService.getFilteredSearch != "maria" {
		t.Fatalf("expected search maria, got %s", fakeService.getFilteredSearch)
	}
	if fakeService.getFilteredJenisKelamin != "Perempuan" {
		t.Fatalf("expected jenis_kelamin Perempuan, got %s", fakeService.getFilteredJenisKelamin)
	}
	if fakeService.getFilteredStatusJemaat != "Jemaat" {
		t.Fatalf("expected status_jemaat Jemaat, got %s", fakeService.getFilteredStatusJemaat)
	}
	if fakeService.getFilteredStatusDiakonia != "Ya" {
		t.Fatalf("expected status_diakonia Ya, got %s", fakeService.getFilteredStatusDiakonia)
	}
	if fakeService.getFilteredKelompokIbadah != "Youth" {
		t.Fatalf("expected kelompok_ibadah Youth, got %s", fakeService.getFilteredKelompokIbadah)
	}
	if fakeService.getFilteredSortBy != "nama_lengkap" {
		t.Fatalf("expected sort_by nama_lengkap, got %s", fakeService.getFilteredSortBy)
	}
	if fakeService.getFilteredSortOrder != "desc" {
		t.Fatalf("expected sort_order desc, got %s", fakeService.getFilteredSortOrder)
	}
	if fakeService.getFilteredPage != 1 {
		t.Fatalf("expected page 1, got %d", fakeService.getFilteredPage)
	}
	if fakeService.getFilteredLimit != 2 {
		t.Fatalf("expected limit 2, got %d", fakeService.getFilteredLimit)
	}
}

func TestJemaatGetFilteredPaginatedValidationError(t *testing.T) {
	fakeService := &fakeJemaatService{getFilteredErr: &service.ValidationError{Message: "Invalid sort_by"}}
	handler := Jemaat(nil, fakeService)
	request := httptest.NewRequest(http.MethodGet, "/api/jemaat?sort_by=nama_salah", nil)
	recorder := httptest.NewRecorder()
	handler(recorder, request)
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", recorder.Code)
	}
}

func TestJemaatGetFilteredPaginatedServiceError(t *testing.T) {
	fakeService := &fakeJemaatService{getFilteredErr: errors.New("database error")}
	handler := Jemaat(nil, fakeService)
	request := httptest.NewRequest(http.MethodGet, "/api/jemaat?search=maria", nil)
	recorder := httptest.NewRecorder()
	handler(recorder, request)
	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", recorder.Code)
	}
}

func TestJemaatPost(t *testing.T) {
	fakeService := &fakeJemaatService{createResult: &model.Jemaat{ID: 14, NamaPanggilan: "Eko", NamaLengkap: "Eko Siswanto"}}
	handler := Jemaat(nil, fakeService)
	body := `{"nama_panggilan":"Eko","nama_lengkap":"Eko Siswanto","jenis_kelamin":"Laki-Laki","tanggal_lahir":"1994-01-15T00:00:00Z","domisili":"Cikarang","status_jemaat":"Jemaat","status_diakonia":"Tidak","kelompok_ibadah":"Youth"}`
	request := httptest.NewRequest(http.MethodPost, "/api/jemaat", strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	handler(recorder, request)
	if recorder.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", recorder.Code)
	}
	if fakeService.createRequest.NamaPanggilan != "Eko" {
		t.Fatalf("expected Eko, got %s", fakeService.createRequest.NamaPanggilan)
	}
}

func TestJemaatPostInvalidJSON(t *testing.T) {
	fakeService := &fakeJemaatService{}
	handler := Jemaat(nil, fakeService)
	request := httptest.NewRequest(http.MethodPost, "/api/jemaat", strings.NewReader(`{"nama_panggilan":`))
	recorder := httptest.NewRecorder()
	handler(recorder, request)
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", recorder.Code)
	}
}

func TestJemaatPostValidationError(t *testing.T) {
	fakeService := &fakeJemaatService{}
	handler := Jemaat(nil, fakeService)
	body := `{"nama_panggilan":"","nama_lengkap":"Eko Siswanto","jenis_kelamin":"Laki-Laki","tanggal_lahir":"1994-01-15T00:00:00Z","domisili":"Cikarang","status_jemaat":"Jemaat","status_diakonia":"Tidak","kelompok_ibadah":"Youth"}`
	request := httptest.NewRequest(http.MethodPost, "/api/jemaat", strings.NewReader(body))
	recorder := httptest.NewRecorder()
	handler(recorder, request)
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", recorder.Code)
	}
}

func TestJemaatPostServiceError(t *testing.T) {
	fakeService := &fakeJemaatService{createErr: errors.New("database error")}
	handler := Jemaat(nil, fakeService)
	body := `{"nama_panggilan":"Eko","nama_lengkap":"Eko Siswanto","jenis_kelamin":"Laki-Laki","tanggal_lahir":"1994-01-15T00:00:00Z","domisili":"Cikarang","status_jemaat":"Jemaat","status_diakonia":"Tidak","kelompok_ibadah":"Youth"}`
	request := httptest.NewRequest(http.MethodPost, "/api/jemaat", strings.NewReader(body))
	recorder := httptest.NewRecorder()
	handler(recorder, request)
	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", recorder.Code)
	}
}

func TestJemaatPut(t *testing.T) {
	fakeService := &fakeJemaatService{updateResult: &model.Jemaat{ID: 5, NamaPanggilan: "Doni Baru", NamaLengkap: "Doni Wijaya Baru"}}
	handler := Jemaat(nil, fakeService)
	body := `{"nama_panggilan":"Doni Baru","nama_lengkap":"Doni Wijaya Baru","jenis_kelamin":"Laki-Laki","tanggal_lahir":"1992-03-15T00:00:00Z","domisili":"Jakarta","status_jemaat":"Jemaat","status_diakonia":"Tidak","kelompok_ibadah":"Youth"}`
	request := httptest.NewRequest(http.MethodPut, "/api/jemaat/5", strings.NewReader(body))
	recorder := httptest.NewRecorder()
	handler(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}
	if fakeService.updateID != 5 {
		t.Fatalf("expected ID 5, got %d", fakeService.updateID)
	}
	if fakeService.updateRequest.NamaPanggilan != "Doni Baru" {
		t.Fatalf("expected Doni Baru, got %s", fakeService.updateRequest.NamaPanggilan)
	}
}

func TestJemaatPutInvalidID(t *testing.T) {
	fakeService := &fakeJemaatService{}
	handler := Jemaat(nil, fakeService)
	request := httptest.NewRequest(http.MethodPut, "/api/jemaat/abc", strings.NewReader(`{}`))
	recorder := httptest.NewRecorder()
	handler(recorder, request)
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", recorder.Code)
	}
}

func TestJemaatPutInvalidJSON(t *testing.T) {
	fakeService := &fakeJemaatService{}
	handler := Jemaat(nil, fakeService)
	request := httptest.NewRequest(http.MethodPut, "/api/jemaat/5", strings.NewReader(`{"nama_panggilan":`))
	recorder := httptest.NewRecorder()
	handler(recorder, request)
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", recorder.Code)
	}
}

func TestJemaatPutValidationError(t *testing.T) {
	fakeService := &fakeJemaatService{}
	handler := Jemaat(nil, fakeService)
	body := `{"nama_panggilan":"","nama_lengkap":"Doni Wijaya","jenis_kelamin":"Laki-Laki","tanggal_lahir":"1992-03-15T00:00:00Z","domisili":"Jakarta","status_jemaat":"Jemaat","status_diakonia":"Tidak","kelompok_ibadah":"Youth"}`
	request := httptest.NewRequest(http.MethodPut, "/api/jemaat/5", strings.NewReader(body))
	recorder := httptest.NewRecorder()
	handler(recorder, request)
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", recorder.Code)
	}
}

func TestJemaatPutNotFound(t *testing.T) {
	fakeService := &fakeJemaatService{updateErr: errors.New("jemaat not found")}
	handler := Jemaat(nil, fakeService)
	body := `{"nama_panggilan":"Doni","nama_lengkap":"Doni Wijaya","jenis_kelamin":"Laki-Laki","tanggal_lahir":"1992-03-15T00:00:00Z","domisili":"Jakarta","status_jemaat":"Jemaat","status_diakonia":"Tidak","kelompok_ibadah":"Youth"}`
	request := httptest.NewRequest(http.MethodPut, "/api/jemaat/999", strings.NewReader(body))
	recorder := httptest.NewRecorder()
	handler(recorder, request)
	if recorder.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", recorder.Code)
	}
}

func TestJemaatPutServiceError(t *testing.T) {
	fakeService := &fakeJemaatService{updateErr: errors.New("database error")}
	handler := Jemaat(nil, fakeService)
	body := `{"nama_panggilan":"Doni","nama_lengkap":"Doni Wijaya","jenis_kelamin":"Laki-Laki","tanggal_lahir":"1992-03-15T00:00:00Z","domisili":"Jakarta","status_jemaat":"Jemaat","status_diakonia":"Tidak","kelompok_ibadah":"Youth"}`
	request := httptest.NewRequest(http.MethodPut, "/api/jemaat/5", strings.NewReader(body))
	recorder := httptest.NewRecorder()
	handler(recorder, request)
	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", recorder.Code)
	}
}

func TestJemaatDelete(t *testing.T) {
	fakeService := &fakeJemaatService{}
	handler := Jemaat(nil, fakeService)
	request := httptest.NewRequest(http.MethodDelete, "/api/jemaat/5", nil)
	recorder := httptest.NewRecorder()
	handler(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}
	if fakeService.deleteID != 5 {
		t.Fatalf("expected ID 5, got %d", fakeService.deleteID)
	}
}

func TestJemaatDeleteInvalidID(t *testing.T) {
	fakeService := &fakeJemaatService{}
	handler := Jemaat(nil, fakeService)
	request := httptest.NewRequest(http.MethodDelete, "/api/jemaat/abc", nil)
	recorder := httptest.NewRecorder()
	handler(recorder, request)
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", recorder.Code)
	}
}

func TestJemaatDeleteNotFound(t *testing.T) {
	fakeService := &fakeJemaatService{deleteErr: errors.New("jemaat not found")}
	handler := Jemaat(nil, fakeService)
	request := httptest.NewRequest(http.MethodDelete, "/api/jemaat/999", nil)
	recorder := httptest.NewRecorder()
	handler(recorder, request)
	if recorder.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", recorder.Code)
	}
}

func TestJemaatDeleteServiceError(t *testing.T) {
	fakeService := &fakeJemaatService{deleteErr: errors.New("database error")}
	handler := Jemaat(nil, fakeService)
	request := httptest.NewRequest(http.MethodDelete, "/api/jemaat/5", nil)
	recorder := httptest.NewRecorder()
	handler(recorder, request)
	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", recorder.Code)
	}
}

func TestJemaatMethodNotAllowed(t *testing.T) {
	fakeService := &fakeJemaatService{}
	handler := Jemaat(nil, fakeService)
	request := httptest.NewRequest(http.MethodPatch, "/api/jemaat", nil)
	recorder := httptest.NewRecorder()
	handler(recorder, request)
	if recorder.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", recorder.Code)
	}
}
