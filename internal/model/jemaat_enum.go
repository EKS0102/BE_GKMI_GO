package model

type JenisKelamin string

const (
	JenisKelaminLakiLaki  JenisKelamin = "Laki-Laki"
	JenisKelaminPerempuan JenisKelamin = "Perempuan"
)

func (j JenisKelamin) IsValid() bool {
	switch j {
	case JenisKelaminLakiLaki, JenisKelaminPerempuan:
		return true
	default:
		return false
	}
}

type StatusJemaat string

const (
	StatusJemaatJemaat     StatusJemaat = "Jemaat"
	StatusJemaatSimpatisan StatusJemaat = "Simpatisan"
	StatusJemaatTamu       StatusJemaat = "Tamu"
)

func (s StatusJemaat) IsValid() bool {
	switch s {
	case StatusJemaatJemaat,
		StatusJemaatSimpatisan,
		StatusJemaatTamu:
		return true
	default:
		return false
	}
}

type StatusDiakonia string

const (
	StatusDiakoniaYa    StatusDiakonia = "Ya"
	StatusDiakoniaTidak StatusDiakonia = "Tidak"
)

func (s StatusDiakonia) IsValid() bool {
	switch s {
	case StatusDiakoniaYa, StatusDiakoniaTidak:
		return true
	default:
		return false
	}
}

type KelompokIbadah string

const (
	KelompokIbadahSekolahMinggu KelompokIbadah = "Sekolah Minggu"
	KelompokIbadahYouth         KelompokIbadah = "Youth"
	KelompokIbadahKompak        KelompokIbadah = "Kompak"
	KelompokIbadahKoper         KelompokIbadah = "Koper"
	KelompokIbadahLansia        KelompokIbadah = "Lansia"
)

func (k KelompokIbadah) IsValid() bool {
	switch k {
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

type SortBy string

const (
	SortByID             SortBy = "id"
	SortByNamaPanggilan  SortBy = "nama_panggilan"
	SortByNamaLengkap    SortBy = "nama_lengkap"
	SortByTanggalLahir   SortBy = "tanggal_lahir"
	SortByJenisKelamin   SortBy = "jenis_kelamin"
	SortByStatusJemaat   SortBy = "status_jemaat"
	SortByStatusDiakonia SortBy = "status_diakonia"
	SortByKelompokIbadah SortBy = "kelompok_ibadah"
)

func (s SortBy) IsValid() bool {
	switch s {
	case SortByID,
		SortByNamaPanggilan,
		SortByNamaLengkap,
		SortByTanggalLahir,
		SortByJenisKelamin,
		SortByStatusJemaat,
		SortByStatusDiakonia,
		SortByKelompokIbadah:
		return true
	default:
		return false
	}
}

type SortOrder string

const (
	SortOrderAsc  SortOrder = "asc"
	SortOrderDesc SortOrder = "desc"
)

func (s SortOrder) IsValid() bool {
	switch s {
	case SortOrderAsc, SortOrderDesc:
		return true
	default:
		return false
	}
}
