package model

import "time"

type Jemaat struct {
	ID             int            `json:"id"`
	NamaPanggilan  string         `json:"nama_panggilan"`
	NamaLengkap    string         `json:"nama_lengkap"`
	JenisKelamin   JenisKelamin   `json:"jenis_kelamin"`
	TanggalLahir   time.Time      `json:"tanggal_lahir"`
	Domisili       string         `json:"domisili"`
	StatusJemaat   StatusJemaat   `json:"status_jemaat"`
	StatusDiakonia StatusDiakonia `json:"status_diakonia"`
	KelompokIbadah KelompokIbadah `json:"kelompok_ibadah"`
}
