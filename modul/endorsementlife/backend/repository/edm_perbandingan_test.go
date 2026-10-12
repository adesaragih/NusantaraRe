package repository

// Bentuk SQL perbandingan summary (08-10-2026) - tanpa Oracle.

import (
	"strings"
	"testing"

	"nusantarare/modul/endorsementlife/backend/models"
)

func TestSQLPerbandinganPerMataUangDanKelas(t *testing.T) {
	for _, tipe := range []string{models.TypeQR, models.TypeQP, models.TypeTP, models.TypeTR} {
		q := sqlPerbandingan(uPeserta, tipe)
		if n := penampungUnik(t, "sqlPerbandingan "+tipe, q); n != 1 {
			t.Fatalf("%s: %d penampung", tipe, n)
		}
		kel := kelompokStatus()
		if !strings.Contains(q, "WHERE d.PREMIUM_LIST_ID = :1 GROUP BY d.CURRENCY, "+kel+" ORDER BY d.CURRENCY") {
			t.Errorf("%s: pengelompokan: %s", tipe, q)
		}
		if strings.Contains(q, "SELECT *") || strings.Contains(q, "d.*") {
			t.Errorf("%s: SELECT *", tipe)
		}
		// Ekspresi rekap SAMA dengan sqlSisipRekap: premi/balance per Type dan setiap SUM kolom.
		premi, balance := rumusPremiBalance(tipe)
		sisip := sqlSisipRekap(uRekap, uPeserta, tipe)
		for _, e := range []string{premi, balance} {
			if !strings.Contains(q, "TO_CHAR("+e+", 'TM9'") || !strings.Contains(sisip, e) {
				t.Errorf("%s: rumus %s tidak sama dengan rekap", tipe, e)
			}
		}
		for _, k := range kolomJumlahRekap {
			if !strings.Contains(q, "SUM("+nv(k[1])+")") {
				t.Errorf("%s: kolom %s tidak dijumlah", tipe, k[0])
			}
		}
	}
	// Setiap kolom rekap kasus punya ekspresi nyata (bukan "0").
	for _, e := range ekspresiRekap(models.TypeQR) {
		if e[1] == "0" {
			t.Errorf("kolom rekap %s tanpa ekspresi", e[0])
		}
	}
	if len(ekspresiRekap(models.TypeQR)) != len(models.KolomRekapKasus)-1 {
		t.Error("cacah kolom perbandingan != kolom rekap")
	}
	kel := kelompokStatus()
	for _, s := range []string{"'Delete', 'Batal') THEN 'D'", "= 'New' THEN 'N'", "ELSE 'O' END"} {
		if !strings.Contains(kel, s) {
			t.Errorf("kelas EDM_STATUS tanpa %q: %s", s, kel)
		}
	}
}
