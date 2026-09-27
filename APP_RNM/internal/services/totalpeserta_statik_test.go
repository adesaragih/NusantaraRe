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
	dir := filepath.Join("..", "repository", "migrations")
	masuk, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("membaca %s: %v", dir, err)
	}
	if len(masuk) == 0 {
		t.Fatal("nol berkas migrasi - penjaga ini tidak menjaga apa pun")
	}
	for _, e := range masuk {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".sql") {
			continue
		}
		isi, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			t.Fatalf("membaca %s: %v", e.Name(), err)
		}
		for _, cocok := range polaKolomTotal.FindAllString(string(isi), -1) {
			t.Errorf("%s memuat kolom %q. Keenam total peserta DIHITUNG saat "+
				"dibaca (models.HitungTotalPeserta, bukti di SavePesertaClaim.xml "+
				"b4221-b4743); menyimpannya membuat angka uang yang dapat basi "+
				"terhadap baris adjustment-nya sendiri.", e.Name(), cocok)
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
