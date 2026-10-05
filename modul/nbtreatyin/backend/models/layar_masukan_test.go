package models

// Uji audit silang putaran 3 bab 7 (W4): medan admin yang diterima dari kiriman
// layar mengikuti sel / wadah yang TAMPIL dan TERBUKA di
// `Section/DetailPolicyTreatyIn` (dibaca ulang 04-10-2026) - matriks fungsi
// murni `GabungMasukanLayar` yang lebih lebar dari uji seam HTTP. Angsuran
// (W5), Enable / Disable (W3), dan nilai bawaan sel diuji lewat seam HTTP
// (`handlers/masukanlayar_test.go`) - C10 tinjauan P3. Harapan dari langkah XML
// yang dikutip, dihitung tangan.

import "testing"

// W4 - medan admin dari sel / wadah TERSEMBUNYI tidak diterima dari layar:
//
//	S19 `.IsNewPolicyNonProp != 1 && .IsNewPolicyListFormat != 1`  medan uang, .Installment
//	.IDCurrency  pyCondition `.IsNewPolicyNonProp != 1`
//	S7  `.ClaimType != 'XOL Retro'`  .FlagPPH, .TypeTax (+ pyCondition `.FlagPPH = true`)
//	.FlagRetroTreaty  pyCondition `.ClaimType != 'XOL Retro'`
//	S14 `.QuotationData.ProportionalType = 'Proportional'`  .Quartal, .YearOfQuartal
//	.QuotationData.IsSurveyReport  pyCondition `pyWorkPage.Quotation.ProportionalType != 'NonProportional'`
func TestMedanAdminTersembunyiTidakDiterima(t *testing.T) {
	for _, tt := range []struct {
		nama    string
		server  map[string]string // halaman server sebelum digabung
		kiriman map[string]string
		terima  []string // diterima
		tolak   []string // diabaikan (nilai server bertahan)
	}{
		{"polis NonProp baru: wadah uang tersembunyi",
			map[string]string{"PolicyTreatyIn.IsNewPolicyNonProp": "1", "PolicyTreatyIn.PremiOgp": "3000", "PolicyTreatyIn.Deduction1": "306.6"},
			map[string]string{"PremiOgp": "1", "Deduction1": "2", "Deduction2": "3", "Installment": "4", "GrossPremium": "5", "IDCurrency": "UJI-ID", "Remark": "UJI-R"},
			[]string{"Remark"}, []string{"PremiOgp", "Deduction1", "Deduction2", "Installment", "GrossPremium", "IDCurrency"}},
		{"IsNewPolicyListFormat 1: wadah uang admin tersembunyi",
			map[string]string{"PolicyTreatyIn.IsNewPolicyListFormat": "1"},
			map[string]string{"PremiOgp": "1", "IDCurrency": "UJI-ID"},
			[]string{"IDCurrency"}, []string{"PremiOgp"}},
		{"ClaimType XOL Retro (kiriman): FlagPPH, TypeTax, FlagRetroTreaty tersembunyi",
			map[string]string{},
			map[string]string{"ClaimType": KlaimXOLRetro, "FlagPPH": "true", "TypeTax": "UJI-T", "FlagRetroTreaty": "true"},
			[]string{"ClaimType"}, []string{"FlagPPH", "TypeTax", "FlagRetroTreaty"}},
		{"FlagPPH bukan true: TypeTax tersembunyi",
			map[string]string{"PolicyTreatyIn.FlagPPH": "false"},
			map[string]string{"TypeTax": "UJI-T"},
			nil, []string{"TypeTax"}},
		{"FlagPPH true (kiriman): TypeTax tampil",
			map[string]string{},
			map[string]string{"FlagPPH": "true", "TypeTax": "UJI-T"},
			[]string{"FlagPPH", "TypeTax"}, nil},
		{"QuotationData NonProportional: Quartal / U/Y tersembunyi",
			map[string]string{"PolicyTreatyIn.QuotationData.ProportionalType": JenisNonProporsional},
			map[string]string{"Quartal": "4", "YearOfQuartal": "2026"},
			nil, []string{"Quartal", "YearOfQuartal"}},
		{"QuotationData Proportional: Quartal / U/Y tampil",
			map[string]string{"PolicyTreatyIn.QuotationData.ProportionalType": JenisProporsional},
			map[string]string{"Quartal": "4", "YearOfQuartal": "2026"},
			[]string{"Quartal", "YearOfQuartal"}, nil},
		{"Quotation NonProportional: Survey Report tersembunyi",
			map[string]string{"Quotation.ProportionalType": JenisNonProporsional, "PolicyTreatyIn.IsNewPolicyNonProp": "1"},
			map[string]string{"QuotationData.IsSurveyReport": "UJI-Yes", "StartDate": "2026-10-01"},
			[]string{"StartDate"}, []string{"QuotationData.IsSurveyReport"}},
	} {
		h := HalamanBaru()
		for j, v := range tt.server {
			h.Setel(j, v)
		}
		lama := h.Salin()
		m := HalamanBaru()
		for j, v := range tt.kiriman {
			m.Setel("PolicyTreatyIn."+j, v)
		}
		if _, err := GabungMasukanLayar(h, m, PosisiAdmin, nil); err != nil {
			t.Fatalf("%s: %v", tt.nama, err)
		}
		for _, j := range tt.terima {
			if got := h.Ambil("PolicyTreatyIn." + j); got != tt.kiriman[j] {
				t.Errorf("%s: %s = %q, harap diterima %q", tt.nama, j, got, tt.kiriman[j])
			}
		}
		for _, j := range tt.tolak {
			if got, w := h.Ambil("PolicyTreatyIn."+j), lama.Ambil("PolicyTreatyIn."+j); got != w {
				t.Errorf("%s: %s tersembunyi = %q, harap nilai server %q", tt.nama, j, got, w)
			}
		}
	}
}

// `.FlagPPH` change -> runActivity `RemoveTypeTax_ACT` (langkah 2: FlagPPH==false ->
// Property-Remove .TypeTax). Dengan TypeTax tersembunyi tak diterima (W4), sel
// FlagPPH yang dilepas di layar menghapus TypeTax server - bukan nilai lama
// yang tertinggal.
func TestFlagPPHDilepasMenghapusTypeTax(t *testing.T) {
	h := HalamanBaru()
	h.Setel("PolicyTreatyIn.FlagPPH", "true")
	h.Setel("PolicyTreatyIn.TypeTax", TypeTaxInclusive)
	m := HalamanBaru()
	m.Setel("PolicyTreatyIn.FlagPPH", "false")
	if _, err := GabungMasukanLayar(h, m, PosisiAdmin, nil); err != nil {
		t.Fatal(err)
	}
	if got := h.Ambil("PolicyTreatyIn.TypeTax"); got != "" {
		t.Fatalf("TypeTax %q, harap terhapus (RemoveTypeTax_ACT)", got)
	}
}
