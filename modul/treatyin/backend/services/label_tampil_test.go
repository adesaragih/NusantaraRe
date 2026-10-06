package services_test

// Bukti ketiga penerjemah tampil yang lahir dari dokumen desain 5 Oktober
// 2026.
//
// ⛔ Tiap padanan di bawah menyebut GAMBARNYA. Padanan tanpa gambar dan tanpa
// ekspor tidak dipasang — ia ditanyakan.

import (
	"testing"

	"nusantarare/modul/treatyin/backend/services"
)

// ⭐ `dd/mm/yyyy` untuk MEDAN, `dd/mm/yy` untuk baris grid terlipat.
//
// Keduanya ada di gambar 01 pada layar yang SAMA: `Commencement 01/01/2025`
// dan grid Rate of Exchange `01/01/25`.
func TestDuaBentukTanggal(t *testing.T) {
	for _, p := range []struct{ masuk, pendek, panjang string }{
		{"20250101", "01/01/25", "01/01/2025"},
		{"20251231", "31/12/25", "31/12/2025"},
		{"20240118", "18/01/24", "18/01/2024"},
	} {
		if got := services.TanggalTampil(p.masuk); got != p.pendek {
			t.Errorf("TanggalTampil(%q) = %q, mau %q", p.masuk, got, p.pendek)
		}
		if got := services.TanggalTampilPanjang(p.masuk); got != p.panjang {
			t.Errorf("TanggalTampilPanjang(%q) = %q, mau %q", p.masuk, got, p.panjang)
		}
	}
}

// ⛔ Yang BUKAN delapan angka lewat apa adanya — aturan yang sama pada
// keduanya. Kolomnya `VARCHAR2` dan nullable.
func TestTanggalPanjangTidakMengarang(t *testing.T) {
	for _, aneh := range []string{"", "   ", "2025", "20250101T00", "bukan tanggal", "2025-01-01"} {
		if got := services.TanggalTampilPanjang(aneh); got != aneh {
			t.Errorf("TanggalTampilPanjang(%q) = %q, mau apa adanya", aneh, got)
		}
	}
}

// ⭐ Ketiga label `Accounting Mode` yang TERBACA.
func TestCaraPembukuanTampil(t *testing.T) {
	for _, p := range []struct{ tersimpan, mau, sumber string }{
		{"underwriting", "Underwriting Year", "gambar 01, kontrak 1001846"},
		{"loss", "Loss Occuring", "gambar 26, kontrak 1001841"},
		{"accounting", "Accounting Year", "pasangan yang pemilik proses nyatakan"},
	} {
		if got := services.CaraPembukuanTampil(p.tersimpan); got != p.mau {
			t.Errorf("CaraPembukuanTampil(%q) = %q, mau %q (%s)", p.tersimpan, got, p.mau, p.sumber)
		}
	}
}

// ⛔ `risk` TIDAK diterjemahkan, dan uji ini yang menjaganya tetap begitu.
//
// Labelnya tidak ada di satu pun dari 43 gambar dan tidak ada di korpus.
// Satu-satunya "Risk Attaching" di `D:\XML_NURE` ada di MEMO berbahasa
// Indonesia pada `Claim Non Prop\Activity\SetEndDate_Act.xml` — catatan
// pengembang di modul lain, bukan label properti ini.
//
// ⚠️ Uji ini SENGAJA merah pada hari seseorang memasang tebakan. Yang boleh
// mengubahnya adalah jawaban pemilik proses, dan hari itu komentar ini ikut
// dicabut.
func TestRiskTidakDitebak(t *testing.T) {
	if got := services.CaraPembukuanTampil("risk"); got != "risk" {
		t.Errorf("CaraPembukuanTampil(\"risk\") = %q — label itu belum diketahui "+
			"dan tidak boleh dikarang; lihat PERTANYAAN-TERBUKA-LAYAR-PEGA.md", got)
	}
}

// ⭐ `reporting` -> "Reporting" (gambar 01), dan `nonreporting` apa adanya.
//
// ⛔ "Non Reporting" terdengar jelas dan tetap tebakan: ia tidak muncul di
// satu pun dari 43 gambar. 687 kontrak memakainya.
func TestBordereauxTampil(t *testing.T) {
	if got := services.BordereauxTampil("reporting"); got != "Reporting" {
		t.Errorf("BordereauxTampil(%q) = %q, mau Reporting (gambar 01)", "reporting", got)
	}
	if got := services.BordereauxTampil("nonreporting"); got != "nonreporting" {
		t.Errorf("BordereauxTampil(%q) = %q — labelnya belum terlihat di gambar mana pun",
			"nonreporting", got)
	}
}

// ⛔ Nilai di luar himpunan yang terukur lewat APA ADANYA, kedua penerjemah.
func TestPenerjemahLabelTidakMengarang(t *testing.T) {
	for _, aneh := range []string{"", "  ", "UNDERWRITING", "Loss Occuring", "entah"} {
		if got := services.CaraPembukuanTampil(aneh); got != aneh {
			t.Errorf("CaraPembukuanTampil(%q) = %q, mau apa adanya", aneh, got)
		}
		if got := services.BordereauxTampil(aneh); got != aneh {
			t.Errorf("BordereauxTampil(%q) = %q, mau apa adanya", aneh, got)
		}
	}
}

// ===========================================================================
// Cap waktu Pega — cacat nyata grid `Rate of Exchange`, 6 Oktober 2026
// ===========================================================================
//
// `TREATYEXCHANGEYEARLY.STARTDATE` menyimpan cap waktu PENUH
// (`20250701T075400.000 GMT`), bukan tanggal polos. Sebelum perbaikan ini
// gridnya menampilkan teks itu apa adanya, sebab `TanggalTampil` menuntut
// persis delapan digit.
//
// ⚠️ Bedanya tidak terlihat sampai sumbernya berganti: tabel pendaratan
// menyimpan `YYYYMMDD` polos, `TREATYEXCHANGEYEARLY` tidak.
func TestCapWaktuPegaMenjadiTanggal(t *testing.T) {
	for _, u := range []struct{ masuk, mau string }{
		{"20250701T075400.000 GMT", "01/07/25"},
		{"20260630T075500.000 GMT", "30/06/26"},
		{"20250101", "01/01/25"}, // bentuk polos tetap bekerja
	} {
		if got := services.TanggalTampil(u.masuk); got != u.mau {
			t.Errorf("TanggalTampil(%q) = %q, mau %q", u.masuk, got, u.mau)
		}
	}
	if got := services.TanggalTampilPanjang("20250701T075400.000 GMT"); got != "01/07/2025" {
		t.Errorf("TanggalTampilPanjang cap waktu = %q, mau 01/07/2025", got)
	}
}

// ⛔ Dan ia TIDAK mengarang: cap waktu TERPOTONG lewat apa adanya.
//
// Penjaga ini pasangan dari yang di atas. Melonggarkan penerimaan sampai
// "apa pun yang memuat T" akan membuat nilai rusak terbaca sebagai tanggal
// yang masuk akal — dan nilai yang aneh harus tetap terlihat aneh.
func TestCapWaktuTerpotongTidakDikarang(t *testing.T) {
	for _, masuk := range []string{
		"20250101T00",     // jam tidak lengkap
		"20250101Txxxxxx", // jam bukan angka
		"2025070",         // kurang dari delapan
		"",                //
	} {
		if got := services.TanggalDelapanDigit(masuk); got != masuk {
			t.Errorf("TanggalDelapanDigit(%q) = %q, mau apa adanya", masuk, got)
		}
	}
}
