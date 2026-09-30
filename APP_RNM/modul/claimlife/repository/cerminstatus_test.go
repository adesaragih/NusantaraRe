package repository

// OQ-N13 (GILIRAN-18) - Save to RNM menyetel status cermin `OS_AKSEPTASI_KLAIM_LIFE`
// seperti `InsertJsonKlaimLife_sql` b176. TANPA Oracle.

import (
	"context"
	"strings"
	"testing"

	"nusantarare/inti/db"
)

// TestSetelCerminOutstandingHanyaKolomStatus - larangan "tabel warisan
// baca-saja di Save to RNM" dicabut untuk kolom `STS_REJECT` SAJA.
//
// ⛔ Satu kolom di SET, dikunci `ID` + `CASEID` (pola `HapusCerminBelumDisimpan`)
// dan `STS_REJECT IS NULL`: baris era Pega dan baris yang sudah berkeputusan
// (mis. ditolak Komite) tidak ditimpa.
func TestSetelCerminOutstandingHanyaKolomStatus(t *testing.T) {
	q := sqlSetelCerminOutstanding("S.L")
	if q != "UPDATE S.L SET STS_REJECT = :1 WHERE ID = :2 AND CASEID = :3 AND STS_REJECT IS NULL" {
		t.Errorf("penyetel cermin berubah:\n%s", q)
	}
	if err := db.PeriksaSQL(q); err != nil {
		t.Error(err)
	}
	if strings.Contains(strings.ToUpper(q), "COMMIT") {
		t.Errorf("penyetel cermin memuat COMMIT:\n%s", q)
	}
	set := q[strings.Index(q, " SET ")+5 : strings.Index(q, " WHERE ")]
	if strings.Contains(set, ",") {
		t.Errorf("SET menulis lebih dari satu kolom: %s", set)
	}
}

// TestSetelCerminOutstandingMenolakKunciKosong - tanpa kunci, UPDATE dapat
// menyentuh baris warisan yang bukan milik klaim ini.
func TestSetelCerminOutstandingMenolakKunciKosong(t *testing.T) {
	r := NewKlaimLife(nil)
	for nama, k := range map[string][2]string{
		"tanpa adjustment": {" ", "CLM-1"},
		"tanpa CASEID":     {"ADJ-1", ""},
	} {
		if _, err := r.SetelCerminOutstanding(context.Background(), nil, k[0], k[1]); err == nil {
			t.Errorf("%s: diterima", nama)
		}
	}
}
