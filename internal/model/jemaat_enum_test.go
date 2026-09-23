package model

import "testing"

func TestJenisKelaminIsValid(t *testing.T) {
	valid := []JenisKelamin{JenisKelaminLakiLaki, JenisKelaminPerempuan}
	invalid := []JenisKelamin{"Invalid", ""}

	for _, value := range valid {
		if !value.IsValid() {
			t.Errorf("expected %q to be valid", value)
		}
	}

	for _, value := range invalid {
		if value.IsValid() {
			t.Errorf("expected %q to be invalid", value)
		}
	}
}

func TestStatusJemaatIsValid(t *testing.T) {
	valid := []StatusJemaat{StatusJemaatJemaat, StatusJemaatSimpatisan, StatusJemaatTamu}
	invalid := []StatusJemaat{"Invalid", ""}

	for _, value := range valid {
		if !value.IsValid() {
			t.Errorf("expected %q to be valid", value)
		}
	}

	for _, value := range invalid {
		if value.IsValid() {
			t.Errorf("expected %q to be invalid", value)
		}
	}
}

func TestStatusDiakoniaIsValid(t *testing.T) {
	valid := []StatusDiakonia{StatusDiakoniaYa, StatusDiakoniaTidak}
	invalid := []StatusDiakonia{"Invalid", ""}

	for _, value := range valid {
		if !value.IsValid() {
			t.Errorf("expected %q to be valid", value)
		}
	}

	for _, value := range invalid {
		if value.IsValid() {
			t.Errorf("expected %q to be invalid", value)
		}
	}
}

func TestKelompokIbadahIsValid(t *testing.T) {
	valid := []KelompokIbadah{KelompokIbadahSekolahMinggu, KelompokIbadahYouth, KelompokIbadahKompak, KelompokIbadahKoper, KelompokIbadahLansia}
	invalid := []KelompokIbadah{"Invalid", ""}

	for _, value := range valid {
		if !value.IsValid() {
			t.Errorf("expected %q to be valid", value)
		}
	}

	for _, value := range invalid {
		if value.IsValid() {
			t.Errorf("expected %q to be invalid", value)
		}
	}
}
