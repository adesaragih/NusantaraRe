package services

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// Penjaga statik atas keputusan bc: keenam total peserta adalah TURUNAN.
//
// Turunan yang disimpan adalah turunan yang suatu hari berbeda dari
// sumbernya, dan bedanya tidak akan terlihat sampai seseorang menjumlahkan
// ulang dengan tangan. Kedua uji di bawah menjaga dua pintu masuknya:
// kolom di basis data, dan penulisan dari repository.

// polaKolomTotal mencari kolom bernama TOTAL_* di migrasi.
//
// ⚠️ Underscore-nya wajib. Tanpa itu pola ini akan menuduh kata "total" di
// dalam komentar migrasi mana pun - dan penjaga yang menuduh hal yang benar
// akan dilonggarkan orang, bukan dipatuhi.
var polaKolomTotal = regexp.MustCompile(`(?i)\bTOTAL_[A-Z_]+`)

func TestMigrasiTidakMenyimpanTotalPeserta(t *testing.T) {
	// Refactor bentuk B (30-09-2026): dulu satu folder,
	// `internal/repository/migrations`, yang ketika itu memuat migrasi semua
	// modul. Kini folder `migrations/` SETIAP modul, Claim Life termasuk.
	var masuk []string
	for _, pola := range []string{
		filepath.Join("..", "..", "..", "modul", "*", "migrations", "*.sql"),
	} {
		cocok, err := filepath.Glob(pola)
		if err != nil {
			t.Fatalf("membaca %s: %v", pola, err)
		}
		masuk = append(masuk, cocok...)
	}
	if len(masuk) == 0 {
		t.Fatal("nol berkas migrasi - penjaga ini tidak menjaga apa pun")
	}
	for _, jalur := range masuk {
		nama := filepath.Base(jalur)
		isi, err := os.ReadFile(jalur)
		if err != nil {
			t.Fatalf("membaca %s: %v", nama, err)
		}
		for _, cocok := range polaKolomTotal.FindAllString(string(isi), -1) {
			t.Errorf("%s memuat kolom %q. Keenam total peserta DIHITUNG saat "+
				"dibaca (models.HitungTotalPeserta, bukti di SavePesertaClaim.xml "+
				"b4221-b4743); menyimpannya membuat angka uang yang dapat basi "+
				"terhadap baris adjustment-nya sendiri.", nama, cocok)
		}
	}
}

// polaTulisTotal mencari penulisan medan Total pada peserta.
//
// ⚠️ Dipersempit ke `=` yang BUKAN `==`. Pelajaran polaPenulisStatus:
// `.Total ==` di dalam sebuah perbandingan sempat lolos sebagai "penulisan"
// hanya karena ada tanda kurung yang menyusul.
var polaTulisTotal = regexp.MustCompile(`\.Total\s*=([^=]|$)|Total:\s*models\.TotalPeserta|Total:\s*total`)

func TestRepositoryTidakMenulisTotalPeserta(t *testing.T) {
	akar := filepath.Join("..", "repository")
	ketemu := 0
	err := filepath.Walk(akar, func(jalur string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(jalur, ".go") {
			return nil
		}
		isi, err := os.ReadFile(jalur)
		if err != nil {
			return err
		}
		ketemu++
		for i, baris := range strings.Split(string(isi), "\n") {
			if strings.HasPrefix(strings.TrimSpace(baris), "//") {
				continue
			}
			if polaTulisTotal.MatchString(baris) {
				t.Errorf("%s:%d menulis total peserta: %s\n"+
					"Total peserta dirakit di services.KlaimLife.Ambil dari "+
					"baris adjustment; repository membaca kolom, bukan "+
					"menghitung turunan.", jalur, i+1, strings.TrimSpace(baris))
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("menelusuri %s: %v", akar, err)
	}
	if ketemu == 0 {
		t.Fatal("nol berkas Go ditelusuri - penjaga ini tidak menjaga apa pun")
	}
}
