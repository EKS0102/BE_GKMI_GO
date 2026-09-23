package repository

import (
	"context"
	"fmt"
	"strings"

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
	query := `
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
ORDER BY id ASC
`

	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanJemaatRows(rows)
}

func (r *JemaatRepository) GetByID(
	ctx context.Context,
	id int,
) (*model.Jemaat, error) {

	query := `
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
WHERE id = $1
`

	var jemaat model.Jemaat

	err := r.DB.QueryRow(ctx, query, id).Scan(
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

	return &jemaat, nil
}

func (r *JemaatRepository) Create(
	ctx context.Context,
	jemaat model.Jemaat,
) (*model.Jemaat, error) {

	query := `
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
VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
RETURNING
id,
nama_panggilan,
nama_lengkap,
jenis_kelamin,
tanggal_lahir,
domisili,
status_jemaat,
status_diakonia,
kelompok_ibadah
`

	var result model.Jemaat

	err := r.DB.QueryRow(
		ctx,
		query,
		jemaat.NamaPanggilan,
		jemaat.NamaLengkap,
		jemaat.JenisKelamin,
		jemaat.TanggalLahir,
		jemaat.Domisili,
		jemaat.StatusJemaat,
		jemaat.StatusDiakonia,
		jemaat.KelompokIbadah,
	).Scan(
		&result.ID,
		&result.NamaPanggilan,
		&result.NamaLengkap,
		&result.JenisKelamin,
		&result.TanggalLahir,
		&result.Domisili,
		&result.StatusJemaat,
		&result.StatusDiakonia,
		&result.KelompokIbadah,
	)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *JemaatRepository) Update(
	ctx context.Context,
	id int,
	jemaat model.Jemaat,
) (*model.Jemaat, error) {

	query := `
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
RETURNING
id,
nama_panggilan,
nama_lengkap,
jenis_kelamin,
tanggal_lahir,
domisili,
status_jemaat,
status_diakonia,
kelompok_ibadah
`

	var result model.Jemaat

	err := r.DB.QueryRow(
		ctx,
		query,
		jemaat.NamaPanggilan,
		jemaat.NamaLengkap,
		jemaat.JenisKelamin,
		jemaat.TanggalLahir,
		jemaat.Domisili,
		jemaat.StatusJemaat,
		jemaat.StatusDiakonia,
		jemaat.KelompokIbadah,
		id,
	).Scan(
		&result.ID,
		&result.NamaPanggilan,
		&result.NamaLengkap,
		&result.JenisKelamin,
		&result.TanggalLahir,
		&result.Domisili,
		&result.StatusJemaat,
		&result.StatusDiakonia,
		&result.KelompokIbadah,
	)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *JemaatRepository) Delete(
	ctx context.Context,
	id int,
) error {

	query := `
DELETE FROM jemaat
WHERE id = $1
`

	commandTag, err := r.DB.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r *JemaatRepository) GetPaginated(
	ctx context.Context,
	page int,
	limit int,
) ([]model.Jemaat, error) {

	offset := (page - 1) * limit

	query := `
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
ORDER BY id ASC
LIMIT $1
OFFSET $2
`

	rows, err := r.DB.Query(
		ctx,
		query,
		limit,
		offset,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return scanJemaatRows(rows)
}

func (r *JemaatRepository) Count(
	ctx context.Context,
) (int, error) {

	query := `
SELECT COUNT(*)
FROM jemaat
`

	var total int

	err := r.DB.QueryRow(ctx, query).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}

func (r *JemaatRepository) SearchPaginated(
	ctx context.Context,
	search string,
	page int,
	limit int,
) ([]model.Jemaat, error) {

	offset := (page - 1) * limit
	searchPattern := "%" + strings.ToLower(search) + "%"

	query := `
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
WHERE
LOWER(nama_panggilan) LIKE $1
OR LOWER(nama_lengkap) LIKE $1
ORDER BY id ASC
LIMIT $2
OFFSET $3
`

	rows, err := r.DB.Query(
		ctx,
		query,
		searchPattern,
		limit,
		offset,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return scanJemaatRows(rows)
}

func (r *JemaatRepository) CountSearch(
	ctx context.Context,
	search string,
) (int, error) {

	searchPattern := "%" + strings.ToLower(search) + "%"

	query := `
SELECT COUNT(*)
FROM jemaat
WHERE
LOWER(nama_panggilan) LIKE $1
OR LOWER(nama_lengkap) LIKE $1
`

	var total int

	err := r.DB.QueryRow(
		ctx,
		query,
		searchPattern,
	).Scan(&total)

	if err != nil {
		return 0, err
	}

	return total, nil
}

func normalizeSort(
	sortBy string,
	sortOrder string,
) (string, string) {

	allowedSortFields := map[string]string{
		"id":              "id",
		"nama_panggilan":  "nama_panggilan",
		"nama_lengkap":    "nama_lengkap",
		"tanggal_lahir":   "tanggal_lahir",
		"jenis_kelamin":   "jenis_kelamin",
		"status_jemaat":   "status_jemaat",
		"status_diakonia": "status_diakonia",
		"kelompok_ibadah": "kelompok_ibadah",
	}

	column, ok := allowedSortFields[sortBy]
	if !ok {
		column = "id"
	}

	order := strings.ToLower(sortOrder)

	if order != "desc" {
		order = "asc"
	}

	return column, order
}

func (r *JemaatRepository) GetFilteredPaginated(
	ctx context.Context,
	search string,
	jenisKelamin string,
	statusJemaat string,
	statusDiakonia string,
	kelompokIbadah string,
	sortBy string,
	sortOrder string,
	page int,
	limit int,
) ([]model.Jemaat, error) {

	offset := (page - 1) * limit

	query := `
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
WHERE 1=1
`

	args := make([]any, 0)
	argIndex := 1

	if search != "" {
		query += fmt.Sprintf(`
AND (
LOWER(nama_panggilan) LIKE $%d
OR LOWER(nama_lengkap) LIKE $%d
)
`, argIndex, argIndex)

		args = append(
			args,
			"%"+strings.ToLower(search)+"%",
		)

		argIndex++
	}

	if jenisKelamin != "" {
		query += fmt.Sprintf(
			" AND jenis_kelamin = $%d",
			argIndex,
		)

		args = append(args, jenisKelamin)
		argIndex++
	}

	if statusJemaat != "" {
		query += fmt.Sprintf(
			" AND status_jemaat = $%d",
			argIndex,
		)

		args = append(args, statusJemaat)
		argIndex++
	}

	if statusDiakonia != "" {
		query += fmt.Sprintf(
			" AND status_diakonia = $%d",
			argIndex,
		)

		args = append(args, statusDiakonia)
		argIndex++
	}

	if kelompokIbadah != "" {
		query += fmt.Sprintf(
			" AND kelompok_ibadah = $%d",
			argIndex,
		)

		args = append(args, kelompokIbadah)
		argIndex++
	}

	sortColumn, sortDirection := normalizeSort(
		sortBy,
		sortOrder,
	)

	query += fmt.Sprintf(
		" ORDER BY %s %s LIMIT $%d OFFSET $%d",
		sortColumn,
		sortDirection,
		argIndex,
		argIndex+1,
	)

	args = append(
		args,
		limit,
		offset,
	)

	rows, err := r.DB.Query(
		ctx,
		query,
		args...,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return scanJemaatRows(rows)
}

func (r *JemaatRepository) CountFiltered(
	ctx context.Context,
	search string,
	jenisKelamin string,
	statusJemaat string,
	statusDiakonia string,
	kelompokIbadah string,
) (int, error) {

	query := `
SELECT COUNT(*)
FROM jemaat
WHERE 1=1
`

	args := make([]any, 0)
	argIndex := 1

	if search != "" {
		query += fmt.Sprintf(`
AND (
LOWER(nama_panggilan) LIKE $%d
OR LOWER(nama_lengkap) LIKE $%d
)
`, argIndex, argIndex)

		args = append(
			args,
			"%"+strings.ToLower(search)+"%",
		)

		argIndex++
	}

	if jenisKelamin != "" {
		query += fmt.Sprintf(
			" AND jenis_kelamin = $%d",
			argIndex,
		)

		args = append(args, jenisKelamin)
		argIndex++
	}

	if statusJemaat != "" {
		query += fmt.Sprintf(
			" AND status_jemaat = $%d",
			argIndex,
		)

		args = append(args, statusJemaat)
		argIndex++
	}

	if statusDiakonia != "" {
		query += fmt.Sprintf(
			" AND status_diakonia = $%d",
			argIndex,
		)

		args = append(args, statusDiakonia)
		argIndex++
	}

	if kelompokIbadah != "" {
		query += fmt.Sprintf(
			" AND kelompok_ibadah = $%d",
			argIndex,
		)

		args = append(args, kelompokIbadah)
	}

	var total int

	err := r.DB.QueryRow(
		ctx,
		query,
		args...,
	).Scan(&total)

	if err != nil {
		return 0, err
	}

	return total, nil
}

func scanJemaatRows(rows pgx.Rows) ([]model.Jemaat, error) {

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

		jemaatList = append(
			jemaatList,
			jemaat,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return jemaatList, nil
}
