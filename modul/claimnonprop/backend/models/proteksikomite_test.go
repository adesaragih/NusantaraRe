package models

import "testing"

// ProteksiSendKomiteCNP_Act langkah 2 menyalin `ClaimData.CNPCircumtances` ke `.DataCommitteeTreaty.CircumCauseOfLoss`
// BARIS akseptasi (kolom KOMITE_CIRCUM_CAUSE_OF_LOSS) - dibaca CreateChildKomiteCNP_Act ke `Komite.CircumtansesCouseOfLoss`
// layar komite "Circumstances".
func TestProteksiKomiteMenyalinCircumstancesKeBarisAkseptasi(t *testing.T) {
	h := HalamanBaru()
	h.Setel(CD+"CNPCircumtances", "UJI KRONOLOGI")
	h.SetelDaftar(DaftarAdjustment, []Baris{{"PaymentType": "1"}})
	if _, err := ProteksiKirimKomite(konteksUji(), h, 1, "UJI-EMAIL"); err != nil {
		t.Fatal(err)
	}
	if v := h.AmbilDaftar(DaftarAdjustment)[0]["DataCommitteeTreaty.CircumCauseOfLoss"]; v != "UJI KRONOLOGI" {
		t.Fatalf("baris akseptasi DataCommitteeTreaty.CircumCauseOfLoss = %q", v)
	}
	if _, ada := h.Nilai[JalurAnak(DaftarAdjustment, 1, "DataCommitteeTreaty.CircumCauseOfLoss")]; ada {
		t.Fatal("ditulis ke jalur halaman, bukan ke baris akseptasi")
	}
}
