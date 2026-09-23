package repository

import (
	"context"
	"testing"

	"BE_GKMI_NTC_GO/internal/database"
)

func TestJemaatRepositoryGetAll(t *testing.T) {
	db, err := database.Connect()
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}
	defer db.Close(context.Background())

	repo := NewJemaatRepository(db)

	jemaatList, err := repo.GetAll(context.Background())
	if err != nil {
		t.Fatalf("failed to get jemaat: %v", err)
	}

	t.Logf("total jemaat: %d", len(jemaatList))

	for _, jemaat := range jemaatList {
		t.Logf(
			"id=%d nama=%s tanggal_lahir=%s jenis_kelamin=%s status=%s",
			jemaat.ID,
			jemaat.NamaLengkap,
			jemaat.TanggalLahir.Format("2006-01-02"),
			jemaat.JenisKelamin,
			jemaat.StatusJemaat,
		)
	}
}
