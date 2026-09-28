package repository

import (
	"regexp"
	"strings"
	"testing"
	"time"
)

func urutanPenampung(q string) string {
	var got []string
	for _, m := range regexp.MustCompile(`:(\d)`).FindAllStringSubmatch(q, -1) {
		got = append(got, m[1])
	}
	return strings.Join(got, "")
}

// TestSQLGandaBerurutPosisi - driver mengikat berposisi menurut kemunculan.
func TestSQLGandaBerurutPosisi(t *testing.T) {
	for nama, u := range map[string]struct{ q, mau string }{
		"dob":    {sqlDOBSumberKosong("S"), "123"},
		"death":  {sqlStatusWarisanTerakhir("L", "S"), "1234"},
		"health": {sqlAdaWarisanSamaDOL("L", "S"), "12345"},
	} {
		if err := PeriksaSQL(u.q); err != nil {
			t.Errorf("%s: %v", nama, err)
		}
		if got := urutanPenampung(u.q); got != u.mau {
			t.Errorf("%s: urutan penampung %q, mau %q\n%s", nama, got, u.mau, u.q)
		}
	}
}

// TestSQLGandaTidakMengeluarkanNama - nama dan DOB hanya di WHERE.
func TestSQLGandaTidakMengeluarkanNama(t *testing.T) {
	for _, q := range []string{sqlStatusWarisanTerakhir("L", "S"), sqlAdaWarisanSamaDOL("L", "S"), sqlDOBSumberKosong("S")} {
		pilih := q[:strings.Index(q, "FROM")]
		if strings.Contains(pilih, "NAME_OF_INSURED") || regexp.MustCompile(`\bDOB\b`).MatchString(strings.ReplaceAll(pilih, "DOB IS NULL", "")) {
			t.Errorf("SELECT mengeluarkan data pribadi:\n%s", q)
		}
	}
}

// TestSQLGandaVerbatimKunci - kelima kunci SQL warisan, dan urutannya.
func TestSQLGandaVerbatimKunci(t *testing.T) {
	q := sqlStatusWarisanTerakhir("L", "S")
	for _, k := range []string{"o.CEDINGCO", "o.NAME_OF_INSURED", "o.DOB", "o.CERTIFICATE_NO", "o.PL_NUMBER",
		"ORDER BY o.ACCEPTATION_DATE DESC", "FETCH FIRST 1 ROWS ONLY"} {
		if !strings.Contains(q, k) {
			t.Errorf("kunci %q hilang:\n%s", k, q)
		}
	}
	if !strings.Contains(sqlAdaWarisanSamaDOL("L", "S"), "o.LAPSE_DATE = TO_DATE(:5") {
		t.Error("health tidak menyaring LAPSE_DATE = DOL")
	}
}

// TestStatusBarisKodeLamaKosongMenjadiISNULL - `= NULL` tidak pernah benar.
func TestStatusBarisKodeLamaKosongMenjadiISNULL(t *testing.T) {
	q := sqlStatusBarisAdjustment("A", true)
	if !strings.Contains(q, "STS_REJECT IS NULL") || strings.Contains(q, ":5") {
		t.Errorf("kode lama kosong:\n%s", q)
	}
	if n := len(argStatusBarisAdjustment("0", "", time.Time{}, "ADJ", "")); n != 4 {
		t.Errorf("argumen kode lama kosong = %d, mau 4 (penampung :1-:4)", n)
	}
	q = sqlStatusBarisAdjustment("A", false)
	if !strings.Contains(q, "STS_REJECT = :5") {
		t.Errorf("kode lama terisi:\n%s", q)
	}
	if n := len(argStatusBarisAdjustment("1", "N", time.Time{}, "ADJ", "0")); n != 5 {
		t.Errorf("argumen kode lama terisi = %d, mau 5", n)
	}
	for _, lamaKosong := range []bool{true, false} {
		if err := PeriksaSQL(sqlStatusBarisAdjustment("A", lamaKosong)); err != nil {
			t.Error(err)
		}
	}
}
