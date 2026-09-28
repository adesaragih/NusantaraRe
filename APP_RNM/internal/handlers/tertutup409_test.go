package handlers

// Penjaga: SETIAP rute pengubah Claim Life menjawab kasus tertutup dengan 409.
//
// ⛔ Lahir 28-09-2026 (GILIRAN-11, saat menyusun panduan uji layar): enam
// handler - akseptasi, hapus, komite, putaran, tahap, tolak - tidak mengenal
// `ErrKasusSudahTertutup`, sehingga layanan yang MENOLAK kasus tertutup dengan
// benar (penjaga butir bb) sampai ke layar sebagai 500 "gagal ...". Rute DOL
// pernah punya cacat yang sama. Penjaga layanan tidak melihatnya: yang salah
// terjemahannya, bukan penolakannya.

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// dikecualikanTertutup - handler rute pengubah yang memang tidak menjumpai
// kasus tertutup, beserta alasannya.
var dikecualikanTertutup = map[string]string{
	"daftarKlaim": "pendaftaran MELAHIRKAN kasus; belum ada kasus untuk ditutup",
}

func TestSetiapHandlerPengubahMenjawabKasusTertutup409(t *testing.T) {
	isi, err := os.ReadFile("handlers.go")
	if err != nil {
		t.Fatal(err)
	}
	pola := regexp.MustCompile(`mux\.HandleFunc\(\s*"(POST|PUT|DELETE) /api/klaim-life[^"]*",\s*([a-zA-Z]+)\(`)
	cocok := pola.FindAllStringSubmatch(string(isi), -1)
	if len(cocok) < 10 {
		t.Fatalf("hanya %d rute pengubah terbaca; pembacanya yang rusak", len(cocok))
	}
	berkas, _ := filepath.Glob("*.go")
	sumber := map[string]string{}
	for _, b := range berkas {
		if strings.HasSuffix(b, "_test.go") {
			continue
		}
		teks, err := os.ReadFile(b)
		if err != nil {
			t.Fatal(err)
		}
		sumber[b] = string(teks)
	}
	for _, m := range cocok {
		nama := m[2]
		if _, ada := dikecualikanTertutup[nama]; ada {
			continue
		}
		var pemilik string
		for b, teks := range sumber {
			if strings.Contains(teks, "func "+nama+"(") {
				pemilik = b
			}
		}
		if pemilik == "" {
			t.Errorf("%s: fungsi handler tidak ditemukan", nama)
			continue
		}
		if !strings.Contains(sumber[pemilik], "services.ErrKasusSudahTertutup") {
			t.Errorf("%s (%s) tidak menerjemahkan ErrKasusSudahTertutup; kasus tertutup "+
				"akan sampai ke layar sebagai 500", nama, pemilik)
		}
	}
}
