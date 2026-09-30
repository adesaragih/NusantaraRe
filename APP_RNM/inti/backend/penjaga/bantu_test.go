package penjaga

// Pembantu penjaga lintas modul - TANPA Oracle.
//
// Refactor bentuk B (30-09-2026): penjaga di paket ini dulu tinggal di
// `internal/handlers` atau `internal/repository` dan membaca berkas di folder
// yang sama ("*.go", "."). Sesudah kode pindah ke `modul/<nama>/`, pembacaan
// satu folder itu diam-diam menyempit. Pembantu di sini membaca SELURUH modul,
// dan setiap pembacanya gagal terang bila jumlah berkasnya nol.

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// akarAplikasi menunjuk folder APP_RNM dari folder paket ini
// (`inti/backend/penjaga`).
const akarAplikasi = "../../.."

// berkasHandler mengembalikan setiap berkas .go bukan-uji di folder handlers
// seluruh modul - `modul/*/backend/handlers` (`internal/` sudah tidak ada sejak
// refactor bentuk B, CLAUDE.md bab 5).
func berkasHandler(t *testing.T) []string {
	t.Helper()
	var hasil []string
	for _, pola := range []string{
		filepath.Join(akarAplikasi, "modul", "*", "backend", "handlers", "*.go"),
	} {
		cocok, err := filepath.Glob(pola)
		if err != nil {
			t.Fatal(err)
		}
		for _, b := range cocok {
			if !strings.HasSuffix(b, "_test.go") {
				hasil = append(hasil, b)
			}
		}
	}
	sort.Strings(hasil)
	if len(hasil) == 0 {
		t.Fatal("nol berkas handler terbaca; pembacanya yang rusak")
	}
	return hasil
}

// cariHandler mencari SATU berkas handler bernama `nama` di modul mana pun.
func cariHandler(t *testing.T, nama string) (string, bool) {
	t.Helper()
	var ketemu []string
	for _, b := range berkasHandler(t) {
		if filepath.Base(b) == nama {
			ketemu = append(ketemu, b)
		}
	}
	switch len(ketemu) {
	case 0:
		return "", false
	case 1:
		return ketemu[0], true
	}
	t.Fatalf("%s ada di %d folder handlers (%v); peta penjaga harus menyebut modulnya", nama, len(ketemu), ketemu)
	return "", false
}

// berkasGoProduksi mengembalikan setiap berkas .go bukan-uji di bawah
// `inti/` dan `modul/`.
func berkasGoProduksi(t *testing.T) []string {
	t.Helper()
	var hasil []string
	for _, akar := range []string{"inti", "modul"} {
		dir := filepath.Join(akarAplikasi, akar)
		if _, err := os.Stat(dir); err != nil {
			t.Fatalf("folder %s tidak terbaca: %v", dir, err)
		}
		err := filepath.Walk(dir, func(p string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() || !strings.HasSuffix(p, ".go") || strings.HasSuffix(p, "_test.go") {
				return nil
			}
			hasil = append(hasil, p)
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	}
	sort.Strings(hasil)
	if len(hasil) == 0 {
		t.Fatal("nol berkas .go terbaca; pembacanya yang rusak")
	}
	return hasil
}
