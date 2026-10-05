package repository

// Uji SQL rincian Product Name (tombol View, 05-10-2026). TANPA Oracle.

import (
	"database/sql"
	"strings"
	"testing"

	"nusantarare/inti/backend/db"
)

func TestSqlRincianProduk(t *testing.T) {
	induk := sqlRincianInduk("S.M_PRODUCTNAME_LIFE")
	if err := db.PeriksaSQL(induk); err != nil {
		t.Fatal(err)
	}
	for _, mau := range []string{"FROM S.M_PRODUCTNAME_LIFE WHERE ID = :1", "TO_CHAR(BEGIN_DATE, 'YYYY-MM-DD')",
		"TO_CHAR(MAXSUMINSURED, 'TM9'"} {
		if !strings.Contains(induk, mau) {
			t.Errorf("induk tanpa %q:\n%s", mau, induk)
		}
	}
	for _, g := range GridRincianProduk {
		q := sqlRincianGrid(g, "S."+g.Tabel)
		if err := db.PeriksaSQL(q); err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(q, "WHERE PRODUCTID = :1 ORDER BY URUT") {
			t.Errorf("%s: tidak berkunci PRODUCTID / urut URUT:\n%s", g.Judul, q)
		}
	}
	// Baca-saja, kolom disebut satu per satu.
	for _, q := range []string{induk, sqlRincianGrid(GridRincianProduk[0], "S.X")} {
		if strings.Contains(q, "*") || strings.Contains(strings.ToUpper(q), "JSONDATA") {
			t.Errorf("SELECT * / JSONDATA: %s", q)
		}
	}
}

func TestSqlRateProduk(t *testing.T) {
	q := sqlRateProduk("S.RATE_LIFE")
	if err := db.PeriksaSQL(q); err != nil {
		t.Fatal(err)
	}
	for _, mau := range []string{"WHERE IDUSEDBY = :1", "ORDER BY ID DESC, RATE ASC", "FETCH FIRST 501 ROWS ONLY"} {
		if !strings.Contains(q, mau) {
			t.Errorf("rate tanpa %q:\n%s", mau, q)
		}
	}
	if strings.Contains(q, "*") || strings.Contains(q, "JSONDATA") {
		t.Errorf("SELECT * / JSONDATA: %s", q)
	}
	plan := sqlRincianGrid(GridRincianProduk[0], "S.P")
	if !strings.Contains(plan, "RIRATE, RIRATEID FROM") {
		t.Errorf("PLAN LIST tanpa kunci RIRATEID: %s", plan)
	}
}

func TestSqlRiskProduk(t *testing.T) {
	q := sqlRiskProduk("S.RIRISK_LIFE")
	if err := db.PeriksaSQL(q); err != nil {
		t.Fatal(err)
	}
	for _, mau := range []string{"SELECT ID, USEDBY, AGE, YEAR, MONTH, RISK, CONTRACT FROM S.RIRISK_LIFE",
		"WHERE IDUSEDBY = :1", "FETCH FIRST 501 ROWS ONLY"} {
		if !strings.Contains(q, mau) {
			t.Errorf("risk tanpa %q:\n%s", mau, q)
		}
	}
	if strings.Contains(q, "*") {
		t.Errorf("SELECT *: %s", q)
	}
	if induk := sqlRincianInduk("S.M"); !strings.Contains(induk, ", RIRISKID FROM S.M") {
		t.Errorf("induk tanpa kunci RIRISKID: %s", induk)
	}
}

func TestNilaiRincianProduk(t *testing.T) {
	angka := kolomRincian{"RNMSHARE", "Nusantara Re Share (%)", JenisRincianAngka}
	for masuk, mau := range map[string]string{".5": "0.5", "-.25": "-0.25", "100": "100", "1000000000": "1000000000"} {
		if got := nilaiRincian(angka, sql.NullString{String: masuk, Valid: true}); got != mau {
			t.Errorf("%q → %q, mau %q", masuk, got, mau)
		}
	}
	bayar := kolomRincian{"PAYMENT", "Premium Payment Method", JenisRincianTeks}
	if got := nilaiRincian(bayar, sql.NullString{String: "3", Valid: true}); got != "Quarterly" {
		t.Errorf("PAYMENT 3 = %q", got)
	}
	if got := nilaiRincian(angka, sql.NullString{}); got != "" {
		t.Errorf("NULL = %q", got)
	}
}
