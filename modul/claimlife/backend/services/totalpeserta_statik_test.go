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
		filepath.Join("..", "..", "..", "..", "modul", "*", "backend", "migrations", "*.sql"),
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
	diperiksa := 0
	for _, jalur := range masuk {
		nama := filepath.Base(jalur)
		isi, err := os.ReadFile(jalur)
		if err != nil {
			t.Fatalf("membaca %s: %v", nama, err)
		}
		for _, cocok := range kolomTotalTerlarang(jalur, string(isi)) {
			t.Errorf("%s memuat kolom %q. Keenam total peserta DIHITUNG saat "+
				"dibaca (models.HitungTotalPeserta, bukti di SavePesertaClaim.xml "+
				"b4221-b4743); menyimpannya membuat angka uang yang dapat basi "+
				"terhadap baris adjustment-nya sendiri.", relAkar(jalur), cocok)
		}
		if milikClaimLife(jalur) {
			diperiksa++
		}
	}
	if diperiksa < 10 {
		t.Fatalf("hanya %d migrasi Claim Life terbaca; pembacanya yang rusak", diperiksa)
	}
}

// kolomTotalTerlarang - kolom TOTAL_* yang dilarang di migrasi itu (lingkup 10-10-2026, pindai_bantu_test.go).
//
// ⛔ Keenam total peserta adalah turunan KLAIM. Migrasi Claim Life diperiksa utuh; migrasi modul lain hanya
// PERNYATAAN yang menyebut tabel Claim Life (mis. `ALTER TABLE T_CLAIMLF_... ADD TOTAL_...`). Kolom TOTAL_* di tabel
// modul lain (nbfacin 188/193/199) adalah keputusan modul itu, bukan turunan peserta klaim.
func kolomTotalTerlarang(jalur, isi string) []string {
	var bersih strings.Builder
	for _, baris := range strings.Split(isi, "\n") {
		if !strings.HasPrefix(strings.TrimSpace(baris), "--") {
			bersih.WriteString(baris)
			bersih.WriteString("\n")
		}
	}
	if milikClaimLife(jalur) {
		return polaKolomTotal.FindAllString(isi, -1)
	}
	var hasil []string
	for _, pernyataan := range strings.Split(bersih.String(), ";") {
		if polaTabelClaimLife.MatchString(pernyataan) {
			hasil = append(hasil, polaKolomTotal.FindAllString(pernyataan, -1)...)
		}
	}
	return hasil
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

// TestLingkupKolomTotalMigrasi membuktikan lingkup kolomTotalTerlarang (10-10-2026) tanpa menanam berkas migrasi
// sungguhan: migrasi Claim Life tetap merah untuk TOTAL_* di mana pun; migrasi modul lain merah HANYA pada pernyataan
// yang menyebut tabel Claim Life; komentar tidak dihitung.
func TestLingkupKolomTotalMigrasi(t *testing.T) {
	cl := "../../../../modul/claimlife/backend/migrations/024_uji.sql"
	lain := "../../../../modul/nbfacin/backend/migrations/188_uji.sql"
	for _, k := range []struct {
		jalur, isi string
		mau        int
	}{
		{cl, "ALTER TABLE {skema}.T_TABEL_APA_PUN ADD TOTAL_SI NUMBER", 1},
		{cl, "ALTER TABLE {skema}.T_CLAIMLF_ADJUSTMENT ADD (TOTAL_SI NUMBER, TOTAL_PREMI NUMBER)", 2},
		{lain, "ALTER TABLE {skema}.T_NB_FACIN ADD TOTAL_SI NUMBER", 0},
		{lain, "ALTER TABLE {skema}.T_NB_FACIN ADD TOTAL_SI NUMBER;\nALTER TABLE {skema}.T_GENERAL_CLAIM ADD TOTAL_CLAIM NUMBER", 1},
		{lain, "-- ALTER TABLE T_GENERAL_CLAIM ADD TOTAL_X NUMBER\nALTER TABLE {skema}.T_NB_FACIN ADD TOTAL_SI NUMBER", 0},
	} {
		if got := kolomTotalTerlarang(k.jalur, k.isi); len(got) != k.mau {
			t.Errorf("%s %q: %v, mau %d", k.jalur, k.isi, got, k.mau)
		}
	}
}
