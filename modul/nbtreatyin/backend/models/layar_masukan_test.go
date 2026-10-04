package models

// Uji audit silang putaran 3 bab 7 (W3, W4, W5): medan admin yang diterima dari
// kiriman layar mengikuti sel / wadah yang TAMPIL dan TERBUKA di
// `Section/DetailPolicyTreatyIn` (dibaca ulang 04-10-2026); baris angsuran dari
// action set server; nilai bawaan sel. Harapan dari langkah XML yang dikutip,
// dihitung tangan.

import (
	"errors"
	"testing"
	"time"
)

// W5 - grid `.ListInstallment` (S45) `readOnly`: baris TIDAK PERNAH dari layar.
// Penulisnya hanya action set sel terbuka:
//
//	.Installment  change -> refresh FillPaymentInstallment(Installment=.Installment)
//	              3.2 pct += @Math.divide(100,n,4); 3.5 .Premium = @Math.divide(premi*pct,100,4),
//	              .DueDate = @CurrentDateTime(), .PaymentTotal = .Premium
//	sel uang      change -> refresh CountOGPONP_Act -> langkah 9 CountSpreading_Act per
//	              baris, langkah 10 SetValidateInstallment_Act 3.1:
//	              .Premium = BalanceDueTo * @divide(.InstallmentPercentage,100,4)
//
// Hitung tangan: BalanceDueTo 1000 -> Installment 4: 25% x 1000 = 250; uang
// berubah dengan BalanceDueTo 2000: 2000 x 0.5 = 1000.
func TestAngsuranDariServer(t *testing.T) {
	sekarang := time.Date(2026, 10, 3, 9, 0, 0, 0, time.UTC)
	tersimpan := func() *Halaman {
		h := HalamanBaru()
		h.Setel("PolicyTreatyIn.BalanceDueTo", "1000")
		h.Setel("PolicyTreatyIn.NetPremium", "1000")
		h.Setel("PolicyTreatyIn.PremiOgp", "1000")
		h.Setel("PolicyTreatyIn.Installment", "2")
		h.SetelDaftar(DaftarAngsuran, []Baris{
			{"InstallmentNo": "1", "InstallmentPercentage": "50", "Premium": "500", "PaymentTotal": "500", "DueDate": "2026-09-01"},
			{"InstallmentNo": "2", "InstallmentPercentage": "50", "Premium": "500", "PaymentTotal": "500", "DueDate": "2026-09-01"},
		})
		h.SetelDaftar(DaftarSpreading, []Baris{{"TreatyType": "UJI-A", "SharePercentage": "100", "ClaimPercentage": "100", "PremiumSpreaded": "1000"}})
		return h
	}
	kirim := func(isi map[string]string) *Halaman {
		m := HalamanBaru()
		for j, v := range isi {
			m.Setel("PolicyTreatyIn."+j, v)
		}
		m.SetelDaftar(DaftarAngsuran, []Baris{{"InstallmentNo": "1", "InstallmentPercentage": "100", "Premium": "1", "PaymentTotal": "1"}})
		return m
	}
	gabung := func(t *testing.T, h *Halaman, isi map[string]string) PemicuLayar {
		t.Helper()
		p, err := GabungMasukanLayar(h, kirim(isi), PosisiAdmin, nil)
		if err != nil {
			t.Fatal(err)
		}
		return p
	}
	t.Run("baris palsu, Installment dan uang tetap: baris server", func(t *testing.T) {
		h := tersimpan()
		p := gabung(t, h, map[string]string{"Installment": "2", "PremiOgp": "1000"})
		if err := TerapkanPemicu(h, p, sekarang); err != nil {
			t.Fatal(err)
		}
		b := h.AmbilDaftar(DaftarAngsuran)
		if len(b) != 2 || !samaNilai(b[0]["Premium"], "500") || b[0]["DueDate"] != "2026-09-01" {
			t.Fatalf("ListInstallment %v", b)
		}
	})
	t.Run("Installment berubah: FillPaymentInstallment server", func(t *testing.T) {
		h := tersimpan()
		p := gabung(t, h, map[string]string{"Installment": "4"})
		if err := TerapkanPemicu(h, p, sekarang); err != nil {
			t.Fatal(err)
		}
		b := h.AmbilDaftar(DaftarAngsuran)
		if len(b) != 4 {
			t.Fatalf("FillPaymentInstallment(4): %d baris", len(b))
		}
		for i, x := range b {
			if !samaNilai(x["InstallmentPercentage"], "25") || !samaNilai(x["Premium"], "250") ||
				!samaNilai(x["PaymentTotal"], "250") || x["DueDate"] != "2026-10-03" {
				t.Errorf("baris %d: %v", i+1, x)
			}
		}
	})
	t.Run("sel uang berubah: CountOGPONP_Act langkah 9 dan 10", func(t *testing.T) {
		h := tersimpan()
		p := gabung(t, h, map[string]string{"Installment": "2", "PremiOgp": "2000"})
		// services.turunkan (CountNetPremi_act) menghitung NetPremium / BalanceDueTo lebih dulu.
		h.Setel("PolicyTreatyIn.NetPremium", "2000")
		h.Setel("PolicyTreatyIn.BalanceDueTo", "2000")
		if err := TerapkanPemicu(h, p, sekarang); err != nil {
			t.Fatal(err)
		}
		b := h.AmbilDaftar(DaftarAngsuran)
		if len(b) != 2 || !samaNilai(b[0]["Premium"], "1000") || !samaNilai(b[1]["PaymentTotal"], "1000") || b[0]["DueDate"] != "2026-09-01" {
			t.Fatalf("SetValidateInstallment 3.1: %v", b)
		}
		if got := h.AmbilDaftar(DaftarSpreading)[0]["PremiumSpreaded"]; !samaNilai(got, "2000") {
			t.Fatalf("CountOGPONP_Act 9 -> CountSpreading_Act: PremiumSpreaded %q, harap 2000", got)
		}
	})
}

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

// W3 - `.QuotationData.ProportionalType` (`pyReadOnly=true`) dan
// `.IsNewPolicyNonProp` (bukan sel) hanya diubah DataTransform
// `TreatyEnableDisableInput` (tombol "Enable / Disable Input Type", pyVisible
// `.TreatyType='XOL'`): langkah 1 ProportionalType = "NonProportional"; 2-4
// WHEN berurutan -> IsNewPolicyNonProp selalu "0" (`TestEnableDisableSelaluNol`).
// Kiriman diterima hanya bila TreatyType XOL DAN sama dengan hasil DT atas
// halaman server; tombol tak tampil = medan terkunci (diabaikan); nilai lain
// ditolak (`GalatKiriman`, 422).
func TestEnableDisableHanyaHasilTombol(t *testing.T) {
	server := func(treaty, nonProp, jenis string) *Halaman {
		h := HalamanBaru()
		h.Setel("PolicyTreatyIn.TreatyType", treaty)
		h.Setel("PolicyTreatyIn.IsNewPolicyNonProp", nonProp)
		h.Setel("PolicyTreatyIn.QuotationData.ProportionalType", jenis)
		return h
	}
	kiriman := func(nonProp, jenis string) *Halaman {
		m := HalamanBaru()
		m.Setel("PolicyTreatyIn.IsNewPolicyNonProp", nonProp)
		m.Setel("PolicyTreatyIn.QuotationData.ProportionalType", jenis)
		return m
	}
	h := server("XOL", "1", JenisNonProporsional)
	if _, err := GabungMasukanLayar(h, kiriman("0", JenisNonProporsional), PosisiAdmin, nil); err != nil {
		t.Fatalf("XOL, hasil DT: %v", err)
	}
	if h.Ambil("PolicyTreatyIn.IsNewPolicyNonProp") != "0" {
		t.Fatal("XOL: hasil TreatyEnableDisableInput diterima")
	}
	h = server("XOL", "0", JenisProporsional)
	if _, err := GabungMasukanLayar(h, kiriman("0", JenisNonProporsional), PosisiAdmin, nil); err != nil ||
		h.Ambil("PolicyTreatyIn.QuotationData.ProportionalType") != JenisNonProporsional {
		t.Fatalf("XOL Proportional -> NonProportional (DT langkah 1): %v %q", err, h.Ambil("PolicyTreatyIn.QuotationData.ProportionalType"))
	}
	h = server("UJI-QS", "1", JenisNonProporsional)
	if _, err := GabungMasukanLayar(h, kiriman("0", JenisProporsional), PosisiAdmin, nil); err != nil {
		t.Fatalf("bukan XOL: tombol tak tampil, kiriman diabaikan tanpa galat: %v", err)
	}
	if h.Ambil("PolicyTreatyIn.IsNewPolicyNonProp") != "1" || h.Ambil("PolicyTreatyIn.QuotationData.ProportionalType") != JenisNonProporsional {
		t.Fatal("bukan XOL: medan terkunci, nilai server bertahan")
	}
	for _, k := range []*Halaman{kiriman("0", JenisProporsional), kiriman("1", JenisNonProporsional)} {
		h = server("XOL", "0", JenisNonProporsional)
		_, err := GabungMasukanLayar(h, k, PosisiAdmin, nil)
		var g *GalatKiriman
		if !errors.As(err, &g) {
			t.Errorf("XOL, kiriman %v bukan hasil DT: galat %v, harap GalatKiriman", k.Nilai, err)
		}
	}
}

// W5 - `pyDefaultValue` sel terbuka `Section/DetailPolicyTreatyIn` (dua
// satu-satunya di layar NB selain label mati `1=2`): `.QuotationData.IsSurveyReport`
// "No" (pyCondition `pyWorkPage.Quotation.ProportionalType != 'NonProportional'`)
// dan `.TypeTax` "Inclusive" (wadah S7 `.ClaimType != 'XOL Retro'`, pyCondition
// `.FlagPPH = true`) - dipakai bila medan kosong saat sel dirender.
func TestNilaiBawaanSel(t *testing.T) {
	h := HalamanBaru()
	TerapkanNilaiBawaanSel(h)
	if h.Ambil("PolicyTreatyIn.QuotationData.IsSurveyReport") != "No" || h.Ambil("PolicyTreatyIn.TypeTax") != "" {
		t.Fatalf("IsSurveyReport %q, TypeTax %q (FlagPPH kosong: sel tersembunyi)",
			h.Ambil("PolicyTreatyIn.QuotationData.IsSurveyReport"), h.Ambil("PolicyTreatyIn.TypeTax"))
	}
	h.Setel("PolicyTreatyIn.FlagPPH", "true")
	h.Setel("PolicyTreatyIn.QuotationData.IsSurveyReport", "UJI-Yes")
	TerapkanNilaiBawaanSel(h)
	if h.Ambil("PolicyTreatyIn.TypeTax") != TypeTaxInclusive || h.Ambil("PolicyTreatyIn.QuotationData.IsSurveyReport") != "UJI-Yes" {
		t.Fatal("FlagPPH true: TypeTax bawaan Inclusive; isian yang ada tidak ditimpa")
	}
	n := HalamanBaru()
	n.Setel("Quotation.ProportionalType", JenisNonProporsional)
	n.Setel("PolicyTreatyIn.ClaimType", KlaimXOLRetro)
	n.Setel("PolicyTreatyIn.FlagPPH", "true")
	TerapkanNilaiBawaanSel(n)
	if n.Ambil("PolicyTreatyIn.QuotationData.IsSurveyReport") != "" || n.Ambil("PolicyTreatyIn.TypeTax") != "" {
		t.Fatal("sel tersembunyi: tanpa nilai bawaan")
	}
}
