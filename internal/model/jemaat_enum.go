package model

type JenisKelamin string

const (
	JenisKelaminLakiLaki  JenisKelamin = "Laki-Laki"
	JenisKelaminPerempuan JenisKelamin = "Perempuan"
)

type StatusJemaat string

const (
	StatusJemaatJemaat     StatusJemaat = "Jemaat"
	StatusJemaatSimpatisan StatusJemaat = "Simpatisan"
	StatusJemaatTamu       StatusJemaat = "Tamu"
)

type StatusDiakonia string

const (
	StatusDiakoniaYa    StatusDiakonia = "Ya"
	StatusDiakoniaTidak StatusDiakonia = "Tidak"
)

type KelompokIbadah string

const (
	KelompokIbadahSekolahMinggu KelompokIbadah = "Sekolah Minggu"
	KelompokIbadahYouth         KelompokIbadah = "Youth"
	KelompokIbadahKompak        KelompokIbadah = "Kompak"
	KelompokIbadahKoper         KelompokIbadah = "Koper"
	KelompokIbadahLansia        KelompokIbadah = "Lansia"
)

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

type SortOrder string

const (
	SortOrderAsc  SortOrder = "asc"
	SortOrderDesc SortOrder = "desc"
)
