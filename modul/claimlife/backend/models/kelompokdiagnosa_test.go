package models_test

import (
	"testing"

	"nusantarare/modul/claimlife/backend/models"
)

// Daftar GROUP DIAGNOSE dari work owner 12-10-2026 - VERBATIM, berurutan, satu sumber.
func TestKelompokDiagnosaVerbatimBerurutan(t *testing.T) {
	mau := []string{
		"Certain infectious and parasitic diseases",
		"Neoplasms",
		"Diseases of the blood and blood-forming organs and certain disorders involving the immune mechanism",
		"Endocrine, nutritional and metabolic diseases",
		"Mental and behavioural disorders",
		"Diseases of the nervous system",
		"Diseases of the eye and adnexa",
		"Diseases of the ear and mastoid process",
		"Diseases of the circulatory system",
		"Diseases of the respiratory system",
		"Diseases of the digestive system",
		"Diseases of the skin and subcutaneous tissue",
		"Diseases of the musculoskeletal system and connective tissue",
		"Diseases of the genitourinary system",
		"Pregnancy, childbirth and the puerperium",
		"Certain conditions originating in the perinatal period",
		"Congenital malformations, deformations and chromosomal abnormalities",
		"Symptoms, signs and abnormal clinical and laboratory findings, not elsewhere classified",
		"Injury, poisoning and certain other consequences of external causes",
		"External causes of morbidity and mortality",
		"Factors influencing health status and contact with health services",
		"Codes for special purposes",
		"Unknown cause of death/ Sudden death",
		"Unknown cause of death/ Death on arrival",
		"Multi organ failure",
		"Multi organ disease",
	}
	got := models.KelompokDiagnosa()
	if len(got) != 26 || len(got) != len(mau) {
		t.Fatalf("jumlah %d, mau 26", len(got))
	}
	terpanjang := 0
	for i := range mau {
		if got[i] != mau[i] {
			t.Errorf("urutan %d = %q, mau %q", i+1, got[i], mau[i])
		}
		terpanjang = max(terpanjang, len(got[i]))
	}
	// Disimpan apa adanya di GROUP_DIAGNOSE VARCHAR2(255) - tidak pernah dipotong.
	if terpanjang > models.BatasGroupDiagnose || terpanjang != 99 {
		t.Errorf("teks terpanjang %d byte (batas kolom %d, mau 99)", terpanjang, models.BatasGroupDiagnose)
	}
	// SALINAN: mengubah hasilnya tidak mengubah sumbernya.
	got[0] = "UJI-UBAH"
	if models.KelompokDiagnosa()[0] != mau[0] {
		t.Error("KelompokDiagnosa mengembalikan slice sumber, bukan salinan")
	}
}

func TestKelompokDiagnosaSah(t *testing.T) {
	for _, u := range []struct {
		nama, nilai, tersimpan string
		mau                    bool
	}{
		{"kosong", "", "Neoplasms", true},
		{"di daftar", "Neoplasms", "", true},
		{"di luar daftar", "UJI-BARU", "", false},
		{"warisan tidak diubah", "UJI-WARISAN", "UJI-WARISAN", true},
		{"warisan diganti nilai asing", "UJI-LAIN", "UJI-WARISAN", false},
		{"warisan diganti nilai daftar", "Multi organ failure", "UJI-WARISAN", true},
		{"beda huruf besar", "neoplasms", "", false},
	} {
		if got := models.KelompokDiagnosaSah(u.nilai, u.tersimpan); got != u.mau {
			t.Errorf("%s: %v, mau %v", u.nama, got, u.mau)
		}
	}
}
