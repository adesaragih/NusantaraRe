package repository

// Kontrak hilir PremiumList → Claim Life atas `M_LIFE_PREMIUM_DETAIL` - tiket 08.
// TANPA Oracle.
//
// ⛔ Kenapa uji ini ada: sejak pl2 tabel itu punya DUA pemilik di repo ini -
// PENULISnya (polis_warisan.go, modul PremiumList Life) dan PEMBACAnya
// (pesertapolis.go, modul Claim Life). Pelajaran envelope `galat` dan
// `IsCheck`: masing-masing sisi benar menurut dirinya sendiri, dan hanya
// PERTEMUANNYA yang dapat salah. Maka yang diuji di sini pertemuannya:
// setiap kolom yang dibaca Claim Life harus diisi penulis PremiumList.

import (
	"regexp"
	"sort"
	"strings"
	"testing"
)

// kolomDitulisWarisan adalah seluruh kolom yang diisi `PesertaWarisan.Ganti`.
func kolomDitulisWarisan() map[string]bool {
	ada := map[string]bool{"ID": true}
	for _, k := range kolomPesertaWarisan {
		ada[k.Kolom] = true
	}
	return ada
}

// TestKolomBacaClaimLifeDiisiPenulisPremiumList - kontrak dua sisi pl2.
//
// Dua pembaca Claim Life: `kolomSalin` (salinan ke klaim) dan daftar pilih
// `sqlCariPeserta` (layar Find Insured, padanan `GetPesertaClaim_sql1`).
func TestKolomBacaClaimLifeDiisiPenulisPremiumList(t *testing.T) {
	ditulis := kolomDitulisWarisan()

	q, _ := sqlCariPeserta("S.M", "UJI-PL", "", "", 10)
	m := regexp.MustCompile(`(?s)SELECT (.*?)\s+FROM`).FindStringSubmatch(q)
	if m == nil {
		t.Fatal("daftar pilih sqlCariPeserta tidak terbaca; pembacanya yang rusak")
	}
	for nama, kolom := range map[string][]string{
		"kolomSalin":     NamaKolomSalinPeserta(),
		"sqlCariPeserta": namaKolomDaftarPilih(m[1]),
	} {
		if len(kolom) < 5 {
			t.Fatalf("%s: hanya %d kolom terbaca; pembacanya yang rusak", nama, len(kolom))
		}
		// ⚠️ Instrumen diuji atas jawaban yang sudah diketahui: ujung pertama
		// dan terakhir tiap daftar harus terbaca, termasuk yang ber-TO_CHAR.
		mau := map[string][2]string{
			"kolomSalin":     {"ID", "CLAIM_AMOUNT"},
			"sqlCariPeserta": {"PL_NUMBER", "EDMSTATUS"},
		}[nama]
		if kolom[0] != mau[0] || kolom[len(kolom)-1] != mau[1] {
			t.Fatalf("%s terbaca %s…%s, mau %s…%s; pembacanya yang rusak",
				nama, kolom[0], kolom[len(kolom)-1], mau[0], mau[1])
		}
		var hilang []string
		for _, k := range kolom {
			if !ditulis[k] {
				hilang = append(hilang, k)
			}
		}
		sort.Strings(hilang)
		if len(hilang) > 0 {
			t.Errorf("%s (Claim Life) membaca %v dari M_LIFE_PREMIUM_DETAIL, "+
				"tetapi penulis PremiumList tidak mengisinya", nama, hilang)
		}
	}
}

// TestPesertaNBTetapHidupDiJalurBacaKlaim - AC 26 spec Claim Life, tiket 08.
//
// ⛔ Jalur NB tidak mengisi `EDMSTATUS` (`d.EDM_STATUS` kosong di unggahan
// NB), jadi baris warisannya ber-`EDMSTATUS` NULL. Penyaring naif
// `NOT IN ('Delete','Batal')` membuang NULL di Oracle - dan seluruh peserta
// NB hilang dari layar klaim.
func TestPesertaNBTetapHidupDiJalurBacaKlaim(t *testing.T) {
	var sumber string
	for _, k := range kolomPesertaWarisan {
		if k.Kolom == "EDMSTATUS" {
			sumber = k.Sumber
		}
	}
	if sumber != "d.EDM_STATUS" {
		t.Fatalf("EDMSTATUS bersumber %q, mau d.EDM_STATUS (TempValue.EDMStatus)", sumber)
	}
	// NULL tiba di pembaca sebagai teks kosong.
	if !PesertaHidup("") {
		t.Error("peserta ber-EDMSTATUS kosong (NB) disaring keluar dari jalur baca klaim")
	}
	if !strings.Contains(penyaringHidup, "IS NULL") {
		t.Errorf("penyaring hidup %q tidak meloloskan NULL", penyaringHidup)
	}
}

// TestPenulisWarisanMengisiKunciBacaKlaim - kunci baca = kunci tulis.
//
// Claim Life membaca berkunci `PL_NUMBER` (+ `CERTIFICATE_NO`); penulis
// mengisi `PL_NUMBER` dengan nomor yang BARU terbit di transaksi yang sama,
// bukan dari kolom sumber yang mungkin masih kosong.
func TestPenulisWarisanMengisiKunciBacaKlaim(t *testing.T) {
	for _, k := range kolomPesertaWarisan {
		switch k.Kolom {
		case "PL_NUMBER":
			if k.Sumber != "" {
				t.Errorf("PL_NUMBER bersumber %q; ia harus nomor yang dirakit nilaiSalinWarisan", k.Sumber)
			}
		case "CERTIFICATE_NO":
			if k.Sumber != "d.CERTIFICATE_NO" {
				t.Errorf("CERTIFICATE_NO bersumber %q", k.Sumber)
			}
		}
	}
	arg := nilaiSalinWarisan("UJI-PL-9", "UJI-W", BarisWarisan{})
	for i, k := range kolomPesertaWarisan {
		if k.Kolom == "PL_NUMBER" && arg[i] != "UJI-PL-9" {
			t.Errorf("PL_NUMBER dikirim %v", arg[i])
		}
	}
}
