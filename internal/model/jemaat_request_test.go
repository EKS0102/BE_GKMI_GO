package model

import (
	"testing"
	"time"
)

func validJemaatRequest() CreateJemaatRequest {
	return CreateJemaatRequest{
		NamaPanggilan:  "Budi",
		NamaLengkap:    "Budi Santoso",
		JenisKelamin:   JenisKelaminLakiLaki,
		TanggalLahir:   time.Date(2000, 5, 15, 0, 0, 0, 0, time.UTC),
		Domisili:       "Salatiga",
		StatusJemaat:   StatusJemaatJemaat,
		StatusDiakonia: StatusDiakoniaYa,
		KelompokIbadah: KelompokIbadahYouth,
	}
}

func TestCreateJemaatRequestValidate(t *testing.T) {
	tests := []struct {
		name    string
		request CreateJemaatRequest
		want    string
	}{
		{
			name:    "valid request",
			request: validJemaatRequest(),
			want:    "",
		},
		{
			name: "nama panggilan required",
			request: func() CreateJemaatRequest {
				r := validJemaatRequest()
				r.NamaPanggilan = ""
				return r
			}(),
			want: "Nama panggilan is required",
		},
		{
			name: "nama lengkap required",
			request: func() CreateJemaatRequest {
				r := validJemaatRequest()
				r.NamaLengkap = ""
				return r
			}(),
			want: "Nama lengkap is required",
		},
		{
			name: "invalid jenis kelamin",
			request: func() CreateJemaatRequest {
				r := validJemaatRequest()
				r.JenisKelamin = "Tidak Valid"
				return r
			}(),
			want: "Invalid jenis kelamin",
		},
		{
			name: "tanggal lahir required",
			request: func() CreateJemaatRequest {
				r := validJemaatRequest()
				r.TanggalLahir = time.Time{}
				return r
			}(),
			want: "Tanggal lahir is required",
		},
		{
			name: "invalid status jemaat",
			request: func() CreateJemaatRequest {
				r := validJemaatRequest()
				r.StatusJemaat = "Tidak Valid"
				return r
			}(),
			want: "Invalid status jemaat",
		},
		{
			name: "invalid status diakonia",
			request: func() CreateJemaatRequest {
				r := validJemaatRequest()
				r.StatusDiakonia = "Tidak Valid"
				return r
			}(),
			want: "Invalid status diakonia",
		},
		{
			name: "invalid kelompok ibadah",
			request: func() CreateJemaatRequest {
				r := validJemaatRequest()
				r.KelompokIbadah = "Tidak Valid"
				return r
			}(),
			want: "Invalid kelompok ibadah",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.request.Validate()

			if got != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, got)
			}
		})
	}
}

func TestSortByIsValid(t *testing.T) {
	tests := []struct {
		name   string
		sortBy SortBy
		want   bool
	}{
		{
			name:   "id",
			sortBy: SortByID,
			want:   true,
		},
		{
			name:   "nama panggilan",
			sortBy: SortByNamaPanggilan,
			want:   true,
		},
		{
			name:   "nama lengkap",
			sortBy: SortByNamaLengkap,
			want:   true,
		},
		{
			name:   "tanggal lahir",
			sortBy: SortByTanggalLahir,
			want:   true,
		},
		{
			name:   "jenis kelamin",
			sortBy: SortByJenisKelamin,
			want:   true,
		},
		{
			name:   "status jemaat",
			sortBy: SortByStatusJemaat,
			want:   true,
		},
		{
			name:   "status diakonia",
			sortBy: SortByStatusDiakonia,
			want:   true,
		},
		{
			name:   "kelompok ibadah",
			sortBy: SortByKelompokIbadah,
			want:   true,
		},
		{
			name:   "invalid sort field",
			sortBy: SortBy("nama_salah"),
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.sortBy.IsValid()

			if got != tt.want {
				t.Fatalf("expected %v, got %v", tt.want, got)
			}
		})
	}
}

func TestSortOrderIsValid(t *testing.T) {
	tests := []struct {
		name      string
		sortOrder SortOrder
		want      bool
	}{
		{
			name:      "asc",
			sortOrder: SortOrderAsc,
			want:      true,
		},
		{
			name:      "desc",
			sortOrder: SortOrderDesc,
			want:      true,
		},
		{
			name:      "invalid sort order",
			sortOrder: SortOrder("random"),
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.sortOrder.IsValid()

			if got != tt.want {
				t.Fatalf("expected %v, got %v", tt.want, got)
			}
		})
	}
}
