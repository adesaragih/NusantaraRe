package models_test

// Keputusan work owner 06-10-2026: Insured Name diisi dari Ceding Company kontrak yang dipilih (tombol Choose).
// Di Pega `InputPolicyTreatyInDetail_preACT` langkah 3 menyalin `.InsuredName = Quotation.InsuredName`, tetapi
// `Quotation.InsuredName` tidak diisi rule NB Treaty In mana pun (dibawa kasus portal SFA, di luar korpus); dasar
// pilihan: RDB `TreatyRealizationCheckDuplicate` mencocokkan INSUREDNAME produksi dengan CedingCoName.

import (
	"testing"

	"nusantarare/modul/nbtreatyin/backend/models"
)

func TestPilihBisnisMengisiInsuredNameDariCeding(t *testing.T) {
	h := models.HalamanBaru()
	models.TerapkanDetailKontrak(h, models.BarisKontrak{"CEDING": "UJI-CEDING-A", "CEDINGID": "UJI-1"})
	if got := h.Ambil(models.HalamanQuotation + ".InsuredName"); got != "UJI-CEDING-A" {
		t.Errorf("Quotation.InsuredName = %q, harap UJI-CEDING-A", got)
	}
	if got := h.Ambil(models.HalamanPolis + ".InsuredName"); got != "UJI-CEDING-A" {
		t.Errorf("PolicyTreatyIn.InsuredName = %q, harap UJI-CEDING-A", got)
	}
}

func TestPilihBisnisLainMenggantiInsuredName(t *testing.T) {
	h := models.HalamanBaru()
	models.TerapkanDetailKontrak(h, models.BarisKontrak{"CEDING": "UJI-CEDING-A"})
	models.TerapkanDetailKontrak(h, models.BarisKontrak{"CEDING": "UJI-CEDING-B"})
	if got := h.Ambil(models.HalamanPolis + ".InsuredName"); got != "UJI-CEDING-B" {
		t.Errorf("PolicyTreatyIn.InsuredName = %q sesudah memilih kontrak lain, harap UJI-CEDING-B", got)
	}
}

func TestPilihBisnisTanpaCedingMempertahankanInsuredName(t *testing.T) {
	h := models.HalamanBaru()
	h.Setel(models.HalamanQuotation+".InsuredName", "UJI-TERTANGGUNG-LAMA")
	models.TerapkanDetailKontrak(h, models.BarisKontrak{"CEDING": ""})
	if got := h.Ambil(models.HalamanPolis + ".InsuredName"); got != "UJI-TERTANGGUNG-LAMA" {
		t.Errorf("PolicyTreatyIn.InsuredName = %q, harap nilai Quotation yang ada (UJI-TERTANGGUNG-LAMA)", got)
	}
}
