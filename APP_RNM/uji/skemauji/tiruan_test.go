package skemauji

// Penjaga tiruan tabel warisan - TANPA Oracle.

import (
	"regexp"
	"testing"

	"nusantarare/modul/claimlife/backend/repository"
)

// TestTiruanPesertaPolisMemuatSetiapKolomSalin - tiruan tidak boleh
// tertinggal dari pembacanya.
//
// ⛔ Ditemukan GILIRAN-14: `kolomSalin` membaca SHARE_NUSANTARA_RE_GROSS, AGE,
// ENTRY_AGE, dan CURRENT_AGE sejak lanjutan 10, tetapi tiruan ini tidak
// pernah memuatnya. Setiap uji bertag db yang mendaftarkan klaim akan gagal
// ORA-00904 - dan tidak satu pun yang tahu, sebab uji itu selalu melewati
// tanpa ORACLE_DSN.
func TestTiruanPesertaPolisMemuatSetiapKolomSalin(t *testing.T) {
	ddl := ddlTiruanPesertaPolis("SKEMAUJI")
	kolom := repository.NamaKolomSalinPeserta()
	if len(kolom) < 20 {
		t.Fatalf("hanya %d kolom salin terbaca; pembacanya yang rusak", len(kolom))
	}
	for _, k := range kolom {
		if !regexp.MustCompile(`(?m)^\s*` + k + `\s`).MatchString(ddl) {
			t.Errorf("tiruan %s tidak memuat kolom %s yang dibaca kolomSalin",
				namaTabelPesertaPolis, k)
		}
	}
}
