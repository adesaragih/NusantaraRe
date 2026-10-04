package skemauji

// Pembantu DDL semua modul untuk uji data uji tiga modul - TANPA Oracle.
//
// Refactor bentuk B (30-09-2026): `kolomMenurutDDL` dulu dipinjam dari uji
// STRUKTUR di internal/repository, yang membaca folder migrasi Claim Life
// (ketika itu memuat semua modul). Kini dibaca dari DAFTAR modul - sumber
// yang sama dengan yang dijalankan `Pasang` - sehingga data uji tiga modul
// tetap dibandingkan dengan DDL ketiga modul.

import (
	"testing"

	"nusantarare/inti/backend/daftar"
	"nusantarare/inti/backend/migrasi"
)

// kolomMenurutDDL - kolom tiap tabel menurut CREATE TABLE dan ALTER ... ADD
// di migrasi SEMUA modul terdaftar.
func kolomMenurutDDL(t *testing.T) map[string][]string {
	t.Helper()
	hasil := map[string][]string{}
	langkah, err := migrasi.Daftar(false, daftar.SumberMigrasi()...)
	if err != nil {
		t.Fatal(err)
	}
	for _, m := range langkah {
		for _, p := range m.Pernyataan {
			if nama, kolom := migrasi.KolomCreateTable(p); nama != "" {
				hasil[nama] = kolom
				continue
			}
			if nama, kolom := migrasi.KolomAlterTambah(p); nama != "" {
				hasil[nama] = append(hasil[nama], kolom...)
				continue
			}
			// Kolom yang DIBUANG (901 menu datar: PARENT_ID) - sama dengan
			// padanannya di `inti/backend/penjaga` (temuan /code-review).
			if nama, kolom := migrasi.KolomAlterBuang(p); nama != "" {
				var sisa []string
				for _, k := range hasil[nama] {
					if k != kolom[0] {
						sisa = append(sisa, k)
					}
				}
				hasil[nama] = sisa
			}
		}
	}
	return hasil
}
