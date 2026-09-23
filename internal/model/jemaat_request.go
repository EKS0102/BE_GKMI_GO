package model

import "time"

type CreateJemaatRequest struct {
	NamaPanggilan  string         `json:"nama_panggilan"`
	NamaLengkap    string         `json:"nama_lengkap"`
	JenisKelamin   JenisKelamin   `json:"jenis_kelamin"`
	TanggalLahir   time.Time      `json:"tanggal_lahir"`
	Domisili       string         `json:"domisili"`
	StatusJemaat   StatusJemaat   `json:"status_jemaat"`
	StatusDiakonia StatusDiakonia `json:"status_diakonia"`
	KelompokIbadah KelompokIbadah `json:"kelompok_ibadah"`
}

func (r CreateJemaatRequest) Validate() string {
	if r.NamaPanggilan == "" {
		return "Nama panggilan is required"
	}

	if r.NamaLengkap == "" {
		return "Nama lengkap is required"
	}

	if !r.JenisKelamin.IsValid() {
		return "Invalid jenis kelamin"
	}

	if r.TanggalLahir.IsZero() {
		return "Tanggal lahir is required"
	}

	if r.Domisili == "" {
		return "Domisili is required"
	}

	if !r.StatusJemaat.IsValid() {
		return "Invalid status jemaat"
	}

	if !r.StatusDiakonia.IsValid() {
		return "Invalid status diakonia"
	}

	if !r.KelompokIbadah.IsValid() {
		return "Invalid kelompok ibadah"
	}

	return ""
}

type UpdateJemaatRequest struct {
	NamaPanggilan  string         `json:"nama_panggilan"`
	NamaLengkap    string         `json:"nama_lengkap"`
	JenisKelamin   JenisKelamin   `json:"jenis_kelamin"`
	TanggalLahir   time.Time      `json:"tanggal_lahir"`
	Domisili       string         `json:"domisili"`
	StatusJemaat   StatusJemaat   `json:"status_jemaat"`
	StatusDiakonia StatusDiakonia `json:"status_diakonia"`
	KelompokIbadah KelompokIbadah `json:"kelompok_ibadah"`
}

func (r UpdateJemaatRequest) Validate() string {
	if r.NamaPanggilan == "" {
		return "Nama panggilan is required"
	}

	if r.NamaLengkap == "" {
		return "Nama lengkap is required"
	}

	if !r.JenisKelamin.IsValid() {
		return "Invalid jenis kelamin"
	}

	if r.TanggalLahir.IsZero() {
		return "Tanggal lahir is required"
	}

	if r.Domisili == "" {
		return "Domisili is required"
	}

	if !r.StatusJemaat.IsValid() {
		return "Invalid status jemaat"
	}

	if !r.StatusDiakonia.IsValid() {
		return "Invalid status diakonia"
	}

	if !r.KelompokIbadah.IsValid() {
		return "Invalid kelompok ibadah"
	}

	return ""
}

func isValidJenisKelamin(value JenisKelamin) bool {
	switch value {
	case JenisKelaminLakiLaki, JenisKelaminPerempuan:
		return true
	default:
		return false
	}
}

func isValidStatusJemaat(value StatusJemaat) bool {
	switch value {
	case StatusJemaatJemaat, StatusJemaatSimpatisan, StatusJemaatTamu:
		return true
	default:
		return false
	}
}

func isValidStatusDiakonia(value StatusDiakonia) bool {
	switch value {
	case StatusDiakoniaYa, StatusDiakoniaTidak:
		return true
	default:
		return false
	}
}

func isValidKelompokIbadah(value KelompokIbadah) bool {
	switch value {
	case KelompokIbadahSekolahMinggu,
		KelompokIbadahYouth,
		KelompokIbadahKompak,
		KelompokIbadahKoper,
		KelompokIbadahLansia:
		return true
	default:
		return false
	}
}
