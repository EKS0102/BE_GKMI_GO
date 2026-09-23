package repository

import (
	"context"

	"BE_GKMI_NTC_GO/internal/model"

	"github.com/jackc/pgx/v5"
)

type JemaatRepository struct {
	DB *pgx.Conn
}

func NewJemaatRepository(db *pgx.Conn) *JemaatRepository {
	return &JemaatRepository{
		DB: db,
	}
}

func (r *JemaatRepository) GetAll(ctx context.Context) ([]model.Jemaat, error) {
	rows, err := r.DB.Query(ctx, `
SELECT
id,
nama_panggilan,
nama_lengkap,
jenis_kelamin,
tanggal_lahir,
domisili,
status_jemaat,
status_diakonia,
kelompok_ibadah
FROM jemaat
ORDER BY id
`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	jemaatList := make([]model.Jemaat, 0)

	for rows.Next() {
		var jemaat model.Jemaat

		err := rows.Scan(
			&jemaat.ID,
			&jemaat.NamaPanggilan,
			&jemaat.NamaLengkap,
			&jemaat.JenisKelamin,
			&jemaat.TanggalLahir,
			&jemaat.Domisili,
			&jemaat.StatusJemaat,
			&jemaat.StatusDiakonia,
			&jemaat.KelompokIbadah,
		)
		if err != nil {
			return nil, err
		}

		jemaatList = append(jemaatList, jemaat)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return jemaatList, nil
}

func (r *JemaatRepository) Create(
	ctx context.Context,
	jemaat *model.Jemaat,
) error {
	err := r.DB.QueryRow(ctx, `
INSERT INTO jemaat (
nama_panggilan,
nama_lengkap,
jenis_kelamin,
tanggal_lahir,
domisili,
status_jemaat,
status_diakonia,
kelompok_ibadah
)
VALUES (
$1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING id
`,
		jemaat.NamaPanggilan,
		jemaat.NamaLengkap,
		jemaat.JenisKelamin,
		jemaat.TanggalLahir,
		jemaat.Domisili,
		jemaat.StatusJemaat,
		jemaat.StatusDiakonia,
		jemaat.KelompokIbadah,
	).Scan(&jemaat.ID)

	return err
}

func (r *JemaatRepository) Update(
	ctx context.Context,
	jemaat *model.Jemaat,
) error {
	result, err := r.DB.Exec(ctx, `
UPDATE jemaat
SET
nama_panggilan = $1,
nama_lengkap = $2,
jenis_kelamin = $3,
tanggal_lahir = $4,
domisili = $5,
status_jemaat = $6,
status_diakonia = $7,
kelompok_ibadah = $8
WHERE id = $9
`,
		jemaat.NamaPanggilan,
		jemaat.NamaLengkap,
		jemaat.JenisKelamin,
		jemaat.TanggalLahir,
		jemaat.Domisili,
		jemaat.StatusJemaat,
		jemaat.StatusDiakonia,
		jemaat.KelompokIbadah,
		jemaat.ID,
	)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r *JemaatRepository) Delete(
	ctx context.Context,
	id int,
) error {
	result, err := r.DB.Exec(ctx, `
DELETE FROM jemaat
WHERE id = $1
`, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
