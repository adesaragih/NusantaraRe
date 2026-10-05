package repository

import (
	"strings"
	"testing"

	"nusantarare/modul/bordereaux/backend/models"
)

// Saringan daftar: hanya filter yang terisi, penampung berurutan, Contains memakai pola LIKE ber-ESCAPE, tanggal rentang.
func TestSaringanDaftar(t *testing.T) {
	where, args := saringan(models.Filter{}, 1)
	if where != "1 = 1" || len(args) != 0 {
		t.Errorf("kosong %q %v", where, args)
	}
	where, args = saringan(models.Filter{BdxID: "bdx_1", Type: "premium", ReportStart: "01-07-2026", ReportEnd: "30-09-2026", Status: "Accept"}, 1)
	mau := `UPPER(BDX_ID) LIKE :1 ESCAPE '\' AND UPPER(TYPE) = UPPER(:2) AND BDXREPORT_START >= TO_DATE(:3, 'DD-MM-YYYY') AND ` +
		`BDXREPORT_END <= TO_DATE(:4, 'DD-MM-YYYY') AND UPPER(STATUSAKSEP) = UPPER(:5)`
	if where != mau || len(args) != 5 || args[0] != `%BDX\_1%` {
		t.Errorf("where %q\nargs %v", where, args)
	}
}

func TestPecahAngka(t *testing.T) {
	for masuk, mau := range map[string][2]any{"1234.50": {"123450", 2}, "-0.5": {"-5", 1}, "0": {"0", 0}, "10": {"10", 0}} {
		k, s, err := PecahAngka(masuk)
		if err != nil || k != mau[0] || s != mau[1] {
			t.Errorf("%s = %v %v %v", masuk, k, s, err)
		}
	}
	if k, _, _ := PecahAngka(""); k != nil {
		t.Error("kosong = NULL")
	}
	if _, _, err := PecahAngka("abc"); err == nil {
		t.Error("bukan angka")
	}
}

// Ekspresi baca: angka lewat TO_CHAR TM9, tanggal DD-MM-YYYY.
func TestEkspresiBaca(t *testing.T) {
	if e := ekspresiBaca(models.KolomCSV{Kolom: "PREMIUM_RNM", Jenis: models.Angka}); !strings.Contains(e, "TO_CHAR(PREMIUM_RNM, 'TM9'") {
		t.Error(e)
	}
	if e := ekspresiBaca(models.KolomCSV{Kolom: "POI_START", Jenis: models.Tanggal}); e != "TO_CHAR(POI_START, 'DD-MM-YYYY')" {
		t.Error(e)
	}
}
