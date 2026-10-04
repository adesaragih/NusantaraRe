package backend

// Migrasi 760 - perbaikan MARKETINGOFFICER_LOG. TANPA Oracle: bentuk SQL-nya.

import (
	"reflect"
	"strings"
	"testing"

	"nusantarare/inti/backend/migrasi"
)

func langkah(t *testing.T, nama string) []string {
	t.Helper()
	p, err := migrasi.PernyataanLangkah(berkasMigrasi, nama)
	if err != nil {
		t.Fatal(err)
	}
	return p
}

// Dua kolom ditambah TANPA default (baris log lama tidak boleh ikut mendapat waktu migrasi), lalu default
// SYSTIMESTAMP untuk baris baru. Trigger warisan TIDAK disentuh.
func TestMigrasi760TambahKolomLogTanpaMenyentuhTrigger(t *testing.T) {
	p := langkah(t, "760_marketingofficer_log.sql")
	mau := []string{
		"ALTER TABLE {skema}.MARKETINGOFFICER_LOG ADD (\n  LOG_TIME    TIMESTAMP(6),\n  AKSES_LOGIN VARCHAR2(50)\n)",
		"ALTER TABLE {skema}.MARKETINGOFFICER_LOG MODIFY (LOG_TIME DEFAULT SYSTIMESTAMP)",
	}
	if !reflect.DeepEqual(p, mau) {
		t.Errorf("760:\n dapat %q\n mau   %q", p, mau)
	}
	if tabel, kolom := migrasi.KolomAlterTambah(p[0]); tabel != "MARKETINGOFFICER_LOG" || !reflect.DeepEqual(kolom, []string{"LOG_TIME", "AKSES_LOGIN"}) {
		t.Errorf("kolom tambahan terbaca %s %v", tabel, kolom)
	}
	mundur := langkah(t, "760_marketingofficer_log_down.sql")
	if !reflect.DeepEqual(mundur, []string{"ALTER TABLE {skema}.MARKETINGOFFICER_LOG DROP (LOG_TIME, AKSES_LOGIN)"}) {
		t.Errorf("760_down: %q", mundur)
	}
	for _, s := range append(p, mundur...) {
		atas := strings.ToUpper(s)
		if strings.Contains(atas, "COMMIT") || strings.Contains(atas, "TRIGGER") {
			t.Errorf("760 menyentuh trigger atau memuat COMMIT: %s", s)
		}
	}
}
