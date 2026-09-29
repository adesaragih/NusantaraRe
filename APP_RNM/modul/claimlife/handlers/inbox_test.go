package handlers

// Uji pintu kotak masuk - F0.4.
//
// Yang diuji di sini PENGURAIAN PARAMETER dan pemetaan galatnya; gerbang
// perannya diuji di `services`.

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func permintaanKueri(kueri string) *http.Request {
	return httptest.NewRequest(http.MethodGet, "/api/klaim-life?"+kueri, nil)
}

func TestBilanganKueriMenolakYangBukanBilangan(t *testing.T) {
	// ⛔ Parameter yang ADA tetapi bukan bilangan DITOLAK, tidak diam-diam
	// jatuh ke bawaan: `?ukuran=limapuluh` yang dijawab 50 menyembunyikan
	// salah ketik, dan pemanggil mengira ia mendapat apa yang ia minta.
	for _, buruk := range []string{"tahap=satu", "halaman=x", "ukuran=limapuluh"} {
		nama := strings.SplitN(buruk, "=", 2)[0]
		if _, sah := bilanganKueri(permintaanKueri(buruk), nama, 7); sah {
			t.Errorf("%q diterima; mau ditolak", buruk)
		}
	}
}

func TestBilanganKueriKosongJatuhKeBawaan(t *testing.T) {
	// Parameter yang TIDAK DISEBUT berarti pemanggil meminta bawaannya.
	n, sah := bilanganKueri(permintaanKueri(""), "ukuran", 50)
	if !sah || n != 50 {
		t.Errorf("kosong = (%d, %v), mau (50, true)", n, sah)
	}
}

func TestBilanganKueriMembacaNilainya(t *testing.T) {
	n, sah := bilanganKueri(permintaanKueri("halaman=3"), "halaman", 1)
	if !sah || n != 3 {
		t.Errorf("halaman=3 = (%d, %v), mau (3, true)", n, sah)
	}
}

// Rute GET koleksi terdaftar, dan ia BUKAN rute yang sama dengan POST-nya.
func TestRuteKotakMasukTerdaftar(t *testing.T) {
	isi, err := os.ReadFile("handlers.go")
	if err != nil {
		t.Fatal(err)
	}
	teks := string(isi)
	for _, wajib := range []string{
		`"GET /api/klaim-life"`,
		`"POST /api/klaim-life"`,
	} {
		if !strings.Contains(teks, wajib) {
			t.Errorf("rute %s tidak terdaftar", wajib)
		}
	}
}

// Kotak masuk TIDAK menyebut uang di kabelnya (ADR-U-0003).
//
// ⛔ Penjaga bentuk. Nilai klaim di daftar akan menggoda seseorang
// menjumlahkannya di layar - dan penjumlahan uang di layar adalah tempat
// float64 masuk kembali lewat pintu belakang.
func TestKotakMasukTanpaMedanUang(t *testing.T) {
	isi, err := os.ReadFile("inbox.go")
	if err != nil {
		t.Fatal(err)
	}
	blok := string(isi)
	awal := strings.Index(blok, "type barisInboxJSON struct")
	akhir := strings.Index(blok[awal:], "\n}") + awal
	medan := blok[awal:akhir]
	// ⚠️ Yang dilarang NILAInya, bukan mata uangnya. Ronde pertama
	// melarang "Uang" dan menuduh `MataUang` - kode mata uang adalah
	// TEKS dan memang milik daftar ini (`InputOSClaimLife` punya
	// kolom CURRENCY). Penjaga yang menuduh kolom yang benar akan
	// dilonggarkan orang sampai ia tidak menjaga apa pun.
	for _, terlarang := range []string{
		"Amount", "claimAmount", "Money", "Jumlah", "NilaiKlaim",
	} {
		if strings.Contains(medan, terlarang) {
			t.Errorf("barisInboxJSON memuat %q; uang tidak ikut di daftar", terlarang)
		}
	}
}
