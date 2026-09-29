package repository

import (
	"database/sql"
	"errors"
	"regexp"
	"strings"
	"testing"

	"nusantarare/internal/models"
)

// AC 47: tidak ada identitas mata uang di teks SQL; saringan di-bind.
func TestSQLKursTCO(t *testing.T) {
	q := sqlDaftarKursTCO("S.K")
	if !strings.Contains(q, "WHERE QUARTER = :1 AND IDCURRENCY = :2") {
		t.Errorf("saringan: %s", q)
	}
	if regexp.MustCompile(`'[^']*'`).MatchString(q) {
		t.Errorf("literal di SQL kurs: %s", q)
	}
	if err := PeriksaSQL(q); err != nil {
		t.Error(err)
	}
	// Master kurs dan master mata uang terdaftar sebagai dibaca-saja - penjaga
	// TestTCOWarisanHanyaDibaca menolak setiap penulis yang menyebutnya.
	for _, m := range []string{MasterKursTahunanTCO, MasterMataUangTCO} {
		ada := false
		for _, x := range masterDibacaSajaTCO {
			ada = ada || x == m
		}
		if !ada {
			t.Errorf("%s tidak terdaftar dibaca-saja", m)
		}
	}
}

func nsKurs(v ...string) [6]sql.NullString {
	var n [6]sql.NullString
	for i, s := range v {
		n[i] = sql.NullString{String: s, Valid: s != ""}
	}
	return n
}

func TestUraiBarisKursTCO(t *testing.T) {
	k, err := uraiBarisKursTCO(nsKurs("15500,25", "20260101T000000.000 GMT", "20261231T000000.000 GMT", "UJI-USD", "USD", "0"))
	if err != nil || k.ToIDR.Text('f') != "15500.25" || k.Mulai.Format("2006-01-02") != "2026-01-01" ||
		k.Akhir.Format("2006-01-02") != "2026-12-31" || k.IDCurrency != "UJI-USD" || k.Quarter != "0" {
		t.Errorf("urai: %+v %v", k, err)
	}
	for _, buruk := range [][6]sql.NullString{
		nsKurs("", "20260101T000000.000 GMT", "20261231T000000.000 GMT"),
		nsKurs("15500", "01/01/2026", "20261231T000000.000 GMT"),
		nsKurs("15500", "20260101T000000.000 GMT", ""),
	} {
		if _, err := uraiBarisKursTCO(buruk); !errors.Is(err, models.ErrKursTakTerurai) {
			t.Errorf("%v: %v", buruk, err)
		}
	}
}
