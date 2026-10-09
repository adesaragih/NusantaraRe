package models_test

import (
	"testing"

	"nusantarare/inti/backend/kontrak"
	"nusantarare/modul/komiteclaimprop/backend/models"
)

// Laporan work owner 09-10-2026 pada layar komite: "perbaiki itu yang tampilin kode kode" - grid klaim menampilkan
// CurrencyID (10026), TreatyType (10035), dan Type estimasi (1). Yang tampil: nama mata uang, nama treaty, label
// Type; kode hanya bila namanya kosong.
func TestGridKomiteMenampilkanNamaBukanKode(t *testing.T) {
	kl := kontrak.KlaimTreaty{Adjustment: 1, Daftar: map[string][]map[string]string{
		"ClaimData.AdjustmentList":  {{"ID": "UJI-ADJ-1", "Type": "2", "Currency": "UJI-USD", "TreatyName": "UJI QS"}},
		"ClaimData.InterestList":    {{"CurrencyID": "UJI-1", "Currency": "UJI-IDR"}, {"CurrencyID": "UJI-2"}},
		"ClaimData.ListClaimAmount": {{"CurrencyID": "UJI-1", "Currency": "UJI-IDR"}},
		"ClaimData.SpreadingRisk": {{"CurrencyID": "UJI-1", "Currency": "UJI-IDR", "TreatyType": "UJI-10035",
			"TreatyName": "UJI QUOTA SHARE"}},
		"ClaimData.EstimationList": {{"CurrencyID": "UJI-2", "Currency": "UJI-USD", "Type": "1"}},
	}}
	ly := models.SusunLayar(kasusUji(1), kl, nil, "UJI-K1", nil)
	sel := func(judul string, baris int, label string) string {
		t.Helper()
		for _, b := range ly.Bagian {
			for _, g := range b.Grid {
				if g.Judul != judul {
					continue
				}
				for _, k := range g.Kolom {
					if k.Label == label {
						return g.Baris[baris][k.Properti]
					}
				}
				t.Fatalf("grid %q tanpa kolom %q", judul, label)
			}
		}
		t.Fatalf("grid %q tidak ada", judul)
		return ""
	}
	cek := func(judul string, baris int, label, mau string) {
		t.Helper()
		if v := sel(judul, baris, label); v != mau {
			t.Errorf("%s baris %d kolom %s = %q, mau %q", judul, baris+1, label, v, mau)
		}
	}
	cek("Insured Interests 100 %", 0, "Currency", "UJI-IDR")
	cek("Insured Interests 100 %", 1, "Currency", "UJI-2") // nama kosong: kodenya
	cek("Count Claim Amount", 0, "Currency", "UJI-IDR")
	cek("Loss Allocation", 0, "Curr", "UJI-IDR")
	cek("Loss Allocation", 0, "Treaty Type", "UJI QUOTA SHARE")
	cek("Estimation List", 0, "Currency", "UJI-USD")
	cek("Estimation List", 0, "Type", "Claim")
	cek("History Adjustment", 0, "Payment Type", "Adjuster Fee")
}
