package repository

import (
	"strings"
	"testing"

	"nusantarare/modul/endorsementlife/backend/models"
)

func TestSQLCSV(t *testing.T) {
	a := sqlAcuanCSV(uPeserta)
	if penampungUnik(t, "sqlAcuanCSV", a) != 2 || !strings.Contains(a, "(d.EDM_STATUS IS NULL OR d.EDM_STATUS <> :2)") ||
		!strings.Contains(a, "ORDER BY d.CERTIFICATE_NO, d.NAME_OF_INSURED, d.ID FETCH FIRST 1 ROWS ONLY") {
		t.Errorf("acuan bukan peserta pertama urutan grid di luar New: %s", a)
	}
	if !strings.Contains(sqlPeserta(uPeserta, nil), "ORDER BY d.CERTIFICATE_NO, d.NAME_OF_INSURED, d.ID") {
		t.Error("urutan grid berubah; acuan CSV harus mengikutinya")
	}
	h := sqlHapusPesertaBaru(uPeserta)
	if penampungUnik(t, "sqlHapusPesertaBaru", h) != 2 || !strings.HasSuffix(h, "WHERE d.PREMIUM_LIST_ID = :1 AND d.EDM_STATUS = :2") {
		t.Errorf("hapus New: %s", h)
	}
	s := sqlSisipPesertaCSV(uPeserta)
	tersimpan := models.KolomCSVTersimpan()
	desimal := 0
	for _, c := range tersimpan {
		if desimalCSV(c) {
			desimal++
		}
	}
	if n := penampungUnik(t, "sqlSisipPesertaCSV", s); n != 5+len(tersimpan)+desimal {
		t.Fatalf("%d penampung, mau %d", n, 5+len(tersimpan)+desimal)
	}
	if c := strings.Count(s, "(TO_NUMBER(:"); c != desimal || desimal != 43 {
		t.Errorf("%d desimal kebal NLS dari %d (mau 43: 42 uang + FACTOR)", c, desimal)
	}
	for in, mau := range map[string][2]any{"": {nil, 0}, "0": {"0", 0}, "-10.50": {"-1050", 2}, "0.00000001": {"1", 8}, "1500": {"1500", 0}, "-0.0": {"0", 0}} {
		k, sk := pecahDesimal(in)
		if k != mau[0] || sk != mau[1] {
			t.Errorf("pecahDesimal(%q) = %v, %d; mau %v, %v", in, k, sk, mau[0], mau[1])
		}
	}
	if len(tersimpan) != 69 {
		t.Errorf("%d kolom tersimpan, mau 72 - UW_STATUS, SUM_AT_RISK, REMAINING_PERIOD", len(tersimpan))
	}
	for _, c := range tersimpan {
		if c.Jenis == models.CSVTanggal && !strings.Contains(s, "TO_DATE(:") {
			t.Errorf("%s tanpa TO_DATE", c.Kolom)
		}
	}
	if c := strings.Count(s, "TO_DATE("); c != 9 {
		t.Errorf("%d kolom DATE, mau 9 (STNC/WPC teks)", c)
	}
}

// TestKolomCSVAdaDiDDL - setiap kolom tujuan CSV adalah kolom DDL 052.
func TestKolomCSVAdaDiDDL(t *testing.T) {
	ddl := map[string]bool{}
	for _, k := range kolomDDL(t, "052_t_premium_list_detail.sql") {
		ddl[k] = true
	}
	for _, c := range models.PemetaanCSV {
		switch {
		case c.Kolom != "" && !ddl[c.Kolom]:
			t.Errorf("%s → %s bukan kolom T_PREMIUM_LIST_DETAIL", c.Properti, c.Kolom)
		case c.Kolom == "" && ddl[c.Properti]:
			t.Errorf("%s punya kolom di DDL tetapi tidak disimpan", c.Properti)
		}
	}
}
