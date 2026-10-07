package repository

// Rekap tersimpan dibaca saat Confirm - keputusan work owner 03-10-2026.

import (
	"fmt"
	"strings"
	"testing"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/premiumlistlife/backend/models"
)

func TestSqlBacaRekapTeksDesimalDanTerurut(t *testing.T) {
	q := sqlBacaRekap("SKEMAUJI.T_PREMIUM_LIST_SUMMARY")
	if err := db.PeriksaSQL(q); err != nil {
		t.Fatal(err)
	}
	for _, k := range append([]string{"BALANCE", "PREMIUM", "COMMISSION"}, models.KolomJumlahSummary...) {
		if !strings.Contains(q, fmt.Sprintf(db.FmtDesimal, k)) {
			t.Errorf("kolom %s tidak dibaca sebagai teks desimal TM9", k)
		}
	}
	for _, potong := range []string{"WHERE PREMIUM_LIST_ID = :1", "ORDER BY CURRENCY"} {
		if !strings.Contains(q, potong) {
			t.Errorf("query tanpa %q", potong)
		}
	}
}

// TestSqlTulisNomorRekap - PL_NUMBER rekap (064), satu polis saja.
func TestSqlTulisNomorRekap(t *testing.T) {
	q := sqlTulisNomorRekap("SKEMAUJI.T_PREMIUM_LIST_SUMMARY")
	if err := db.PeriksaSQL(q); err != nil {
		t.Fatal(err)
	}
	for _, potong := range []string{"SET PL_NUMBER = :1", "WHERE PREMIUM_LIST_ID = :2"} {
		if !strings.Contains(q, potong) {
			t.Errorf("query tanpa %q:\n%s", potong, q)
		}
	}
}

// TestSqlTulisWPC - WPC baris utama header saja.
func TestSqlTulisWPC(t *testing.T) {
	q := sqlTulisWPC("SKEMAUJI.T_PREMIUM_LIST")
	if err := db.PeriksaSQL(q); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(q, "SET WPC = :1 WHERE ID = :2") {
		t.Errorf("query WPC:\n%s", q)
	}
}
