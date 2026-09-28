package models

// Rekap uang per mata uang - tiket 05a. TANPA Oracle.

import (
	"errors"
	"testing"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/pkg/utils"
)

// d menyusun desimal dari teks - nol float di uji ini pun.
func d(t *testing.T, s string) *apd.Decimal {
	t.Helper()
	v, err := utils.ParseDecimal(s)
	if err != nil {
		t.Fatalf("desimal %q: %v", s, err)
	}
	return v
}

// TestPremiumSatuCabangBukanEmpat - ralat yang paling mahal di tiket ini.
//
// ⛔ Bacaan yang meratakan `pyPropertiesName`/`pyPropertiesValue` membaca
// `.PREMIUM` sebagai Σ EMPAT kolom. Struktur sesungguhnya empat `WHEN` yang
// SALING MENIADAKAN atas `pyWorkPage.Type`, jadi Σ SATU kolom.
//
// Baris di bawah sengaja mengisi KEEMPAT kolom sekaligus: kalau keempatnya
// dijumlahkan, hasilnya 1000 - dan uji ini akan merah.
func TestPremiumSatuCabangBukanEmpat(t *testing.T) {
	baris := []BarisUang{{
		"GROSS_PREMIUM":              d(t, "100"),
		"GROSS_PREMIUM_REFUND":       d(t, "200"),
		"GROSS_PREMIUM_RETRO":        d(t, "300"),
		"GROSS_PREMIUM_REFUND_RETRO": d(t, "400"),
	}}
	for _, k := range []struct{ tipe, mau string }{
		{"QR", "100"},
		{"QP", "200"},
		{"TP", "300"},
		{"TR", "400"},
	} {
		r, err := RekapPerMataUang(k.tipe, baris, []string{"IDR"})
		if err != nil {
			t.Fatalf("tipe %s: %v", k.tipe, err)
		}
		if len(r) != 1 {
			t.Fatalf("tipe %s: %d rekap, mau 1", k.tipe, len(r))
		}
		if got := utils.FormatDecimal(r[0].Premium); got != k.mau {
			t.Errorf("tipe %s: PREMIUM %q, mau %q "+
				"(kalau 1000, keempat cabang dijumlahkan - itu ralatnya)",
				k.tipe, got, k.mau)
		}
	}
}

// TestBalanceKeempatCabangDariLiteral - contoh terhitung, satu per Type.
//
// Angka dipilih supaya tiap suku terlihat di hasilnya.
func TestBalanceKeempatCabangDariLiteral(t *testing.T) {
	// QR: GROSS_PREMIUM - DEDUCTION - (RI_ADMIN_FEE + BROKERAGE_FEE + TAX
	//     + PROF_COMM + CLAIM)
	//     1000 - 10 - (1 + 2 + 3 + 4 + 5) = 975
	qr := BarisUang{
		"GROSS_PREMIUM": d(t, "1000"), "DEDUCTION": d(t, "10"),
		"RI_ADMIN_FEE": d(t, "1"), "BROKERAGE_FEE": d(t, "2"),
		"TAX": d(t, "3"), "PROF_COMM": d(t, "4"), "CLAIM": d(t, "5"),
	}
	// QP: GROSS_PREMIUM_REFUND + CLAIM_AMOUNT - (DEDUCTION_REFUND
	//     + BROKERAGE_FEE_REFUND + RI_ADMIN_FEE_REFUND + TAX + PROF_COMM + CLAIM)
	//     1000 + 50 - (10 + 20 + 30 + 3 + 4 + 5) = 978
	qp := BarisUang{
		"GROSS_PREMIUM_REFUND": d(t, "1000"), "CLAIM_AMOUNT": d(t, "50"),
		"DEDUCTION_REFUND": d(t, "10"), "BROKERAGE_FEE_REFUND": d(t, "20"),
		"RI_ADMIN_FEE_REFUND": d(t, "30"),
		"TAX":                 d(t, "3"), "PROF_COMM": d(t, "4"), "CLAIM": d(t, "5"),
	}
	// TP: GROSS_PREMIUM_RETRO - DISCOUNT_PREMIUM_RETRO - RI_ADMIN_FEE_RETRO
	//     + BROKERAGE_FEE_RETRO        <- DITAMBAH, keanehan warisan
	//     1000 - 100 - 10 + 7 = 897
	tp := BarisUang{
		"GROSS_PREMIUM_RETRO": d(t, "1000"), "DISCOUNT_PREMIUM_RETRO": d(t, "100"),
		"RI_ADMIN_FEE_RETRO": d(t, "10"), "BROKERAGE_FEE_RETRO": d(t, "7"),
	}
	// TR: sejajar TP dengan kolom *_REFUND_RETRO. 2000 - 200 - 20 + 9 = 1789
	tr := BarisUang{
		"GROSS_PREMIUM_REFUND_RETRO":    d(t, "2000"),
		"DISCOUNT_PREMIUM_REFUND_RETRO": d(t, "200"),
		"RI_ADMIN_FEE_REFUND_RETRO":     d(t, "20"),
		"BROKERAGE_FEE_REFUND_RETRO":    d(t, "9"),
	}
	for _, k := range []struct {
		tipe string
		b    BarisUang
		mau  string
	}{
		{"QR", qr, "975.0000"},
		{"QP", qp, "978.0000"},
		{"TP", tp, "897.0000"},
		{"TR", tr, "1789.0000"},
	} {
		r, err := RekapPerMataUang(k.tipe, []BarisUang{k.b}, []string{"IDR"})
		if err != nil {
			t.Fatalf("tipe %s: %v", k.tipe, err)
		}
		if got := utils.FormatDecimal(r[0].Balance); got != k.mau {
			t.Errorf("tipe %s: BALANCE %q, mau %q", k.tipe, got, k.mau)
		}
	}
}

// TestBrokerageRetroDitambahBukanDikurangi - keanehan warisan nomor 1.
//
// ⛔ VERBATIM. `BROKERAGE_FEE_RETRO` MENAMBAH saldo di cabang TP/TR,
// sedangkan `BROKERAGE_FEE` MENGURANGI di cabang QR. Menormalkan tandanya
// "supaya konsisten" mengubah angka uang yang sudah beredar.
func TestBrokerageRetroDitambahBukanDikurangi(t *testing.T) {
	tambah, kurang, err := SukuBalanceUntukTipe("TP")
	if err != nil {
		t.Fatal(err)
	}
	if !memuat(tambah, "BROKERAGE_FEE_RETRO") {
		t.Error("BROKERAGE_FEE_RETRO tidak di sisi TAMBAH cabang TP")
	}
	if memuat(kurang, "BROKERAGE_FEE_RETRO") {
		t.Error("BROKERAGE_FEE_RETRO di sisi KURANG - itu 'perbaikan' yang mengubah uang")
	}
	// Dan di cabang QR ia DIKURANGI, dengan nama tanpa akhiran retro.
	_, kurangQR, _ := SukuBalanceUntukTipe("QR")
	if !memuat(kurangQR, "BROKERAGE_FEE") {
		t.Error("BROKERAGE_FEE tidak di sisi KURANG cabang QR")
	}
}

// TestCabangQPMemakaiTaxBukanTaxRefund - keanehan warisan nomor 2.
func TestCabangQPMemakaiTaxBukanTaxRefund(t *testing.T) {
	_, kurang, err := SukuBalanceUntukTipe("QP")
	if err != nil {
		t.Fatal(err)
	}
	for _, k := range []string{"TAX", "PROF_COMM", "CLAIM"} {
		if !memuat(kurang, k) {
			t.Errorf("cabang QP tidak mengurangi %q", k)
		}
	}
	// ⛔ Padanan *_REFUND-nya TIDAK dipakai, walau tiga biaya lain memakainya.
	if memuat(kurang, "TAX_REFUND") {
		t.Error("cabang QP memakai TAX_REFUND - korpus memakai TAX")
	}
}

func memuat(s []string, k string) bool {
	for _, x := range s {
		if x == k {
			return true
		}
	}
	return false
}

// TestCommissionHanyaCOMM - tanpa cabang Type, dan tanpa PROF/OVR.
//
// ⛔ Menggabungkan `PROF_COMM` atau `OVR_COMM` ke sini menjumlahkan komisi
// dua kali: keduanya parameter summary tersendiri.
func TestCommissionHanyaCOMM(t *testing.T) {
	b := BarisUang{
		"COMM": d(t, "11"), "PROF_COMM": d(t, "22"), "OVR_COMM": d(t, "33"),
	}
	for _, tipe := range []string{"QR", "QP", "TP", "TR"} {
		r, err := RekapPerMataUang(tipe, []BarisUang{b}, []string{"IDR"})
		if err != nil {
			t.Fatal(err)
		}
		if got := utils.FormatDecimal(r[0].Commission); got != "11" {
			t.Errorf("tipe %s: COMMISSION %q, mau 11", tipe, got)
		}
	}
}

// TestSatuBarisRekapPerMataUang - AC 37.
func TestSatuBarisRekapPerMataUang(t *testing.T) {
	baris := []BarisUang{
		{"GROSS_PREMIUM": d(t, "100")},
		{"GROSS_PREMIUM": d(t, "200")},
		{"GROSS_PREMIUM": d(t, "50")},
	}
	r, err := RekapPerMataUang("QR", baris, []string{"IDR", "USD", "IDR"})
	if err != nil {
		t.Fatal(err)
	}
	if len(r) != 2 {
		t.Fatalf("%d rekap, mau 2 - satu per mata uang", len(r))
	}
	// Urutannya urutan kemunculan, supaya jawabannya tidak berubah-ubah.
	if r[0].Currency != "IDR" || r[1].Currency != "USD" {
		t.Errorf("urutan mata uang %q lalu %q", r[0].Currency, r[1].Currency)
	}
	if got := utils.FormatDecimal(r[0].Premium); got != "150" {
		t.Errorf("PREMIUM IDR %q, mau 150", got)
	}
	if r[0].CacahBaris != 2 || r[1].CacahBaris != 1 {
		t.Errorf("cacah baris %d dan %d, mau 2 dan 1", r[0].CacahBaris, r[1].CacahBaris)
	}
}

// TestPembulatanHanyaDiAkhir - bukan per baris.
//
// ⛔ Tiga baris ber-`0.00005` berjumlah `0.00015`. Membulatkan PER BARIS ke
// empat angka memberi `0.0001 * 3 = 0.0003`; menjumlah lalu membulatkan
// memberi `0.0002`. Selisihnya tumbuh bersama cacah peserta, dan ia uang.
func TestPembulatanHanyaDiAkhir(t *testing.T) {
	var baris []BarisUang
	var cur []string
	for i := 0; i < 3; i++ {
		baris = append(baris, BarisUang{"GROSS_PREMIUM": d(t, "0.00005")})
		cur = append(cur, "IDR")
	}
	r, err := RekapPerMataUang("QR", baris, cur)
	if err != nil {
		t.Fatal(err)
	}
	if got := utils.FormatDecimal(r[0].Balance); got != "0.0002" {
		t.Errorf("BALANCE %q, mau 0.0002 - pembulatan per baris memberi 0.0003", got)
	}
	// PREMIUM tidak dibulatkan sama sekali: hanya BALANCE yang `@divide`-nya.
	if got := utils.FormatDecimal(r[0].Premium); got != "0.00015" {
		t.Errorf("PREMIUM %q, mau 0.00015 - ia tidak dibulatkan", got)
	}
}

// TestTipeDiLuarEmpatDitolak - brief GILIRAN-9 menuntutnya.
func TestTipeDiLuarEmpatDitolak(t *testing.T) {
	for _, tipe := range []string{"", "qr", "FAC", "QRQP", "X"} {
		if _, err := RekapPerMataUang(tipe, nil, nil); !errors.Is(err, ErrTipePLTanpaCabang) {
			t.Errorf("tipe %q diterima; mau ditolak (%v)", tipe, err)
		}
	}
}

// TestBarisDanMataUangHarusSejajar - salah panjang adalah salah rekap.
func TestBarisDanMataUangHarusSejajar(t *testing.T) {
	b := []BarisUang{{"GROSS_PREMIUM": d(t, "1")}}
	if _, err := RekapPerMataUang("QR", b, []string{"IDR", "USD"}); err == nil {
		t.Error("panjang tak sejajar diterima; mau ditolak")
	}
}

// TestKolomAbsenDibacaNolBukanGalat.
//
// ⚠️ Rule aslinya menjumlahkan properti kosong sebagai nol, dan baris yang
// tidak punya kolom retro memang bukan baris yang rusak.
func TestKolomAbsenDibacaNolBukanGalat(t *testing.T) {
	r, err := RekapPerMataUang("TP", []BarisUang{{}}, []string{"IDR"})
	if err != nil {
		t.Fatalf("baris tanpa kolom ditolak: %v", err)
	}
	if got := utils.FormatDecimal(r[0].Balance); got != "0.0000" {
		t.Errorf("BALANCE %q, mau 0.0000", got)
	}
	// Ke-32 kolom tetap ada di rekap, bernilai nol - bukan hilang.
	if len(r[0].Jumlah) != len(KolomUangUnggah) {
		t.Errorf("%d kolom di rekap, mau %d", len(r[0].Jumlah), len(KolomUangUnggah))
	}
}
