package repository

// OQ-N13 (GILIRAN-18) - Save to RNM menyetel status cermin `OS_AKSEPTASI_KLAIM_LIFE`
// seperti `InsertJsonKlaimLife_sql` b176. TANPA Oracle.

import (
	"context"
	"os"
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

// TestStatusCerminMengikutiPenulisStatus - temuan /code-review GILIRAN-18.
//
// ⛔ Sesudah Save to RNM cermin berstatus '0'. Penolakan Admin
// (`RejectOSClaimLife_Act` langkah 5 b1970) dan akseptasi (`SaveAdjustment_Act`
// 1.6.2 b2363) menulis cermin lewat `UpdateOsAkseptasiClaimLife_sql` di Pega;
// tanpa penyelarasan ini cermin tertahan '0' dan gerbang 11.4 memblokir SETIAP
// klaim kematian berikutnya atas tertanggung itu. Karena itu ia tinggal di
// titik tunggal penulis status, `PerbaruiStatusBaris`.
func TestStatusCerminMengikutiPenulisStatus(t *testing.T) {
	for _, dariKosong := range []bool{true, false} {
		q := sqlIkutkanStatusCermin("S.L", "S.W", "S.P", dariKosong)
		if err := db.PeriksaSQL(q); err != nil {
			t.Error(err)
		}
		set := q[strings.Index(q, " SET ")+5 : strings.Index(q, "WHERE")]
		if strings.Contains(set, ",") || !strings.Contains(set, "STS_REJECT = :1") {
			t.Errorf("SET bukan satu kolom STS_REJECT: %s", set)
		}
		for _, wajib := range []string{"o.ID = :2", "p.ID = :3", "w.ID = p.CLAIM_ID", "o.CASEID ="} {
			if !strings.Contains(q, wajib) {
				t.Errorf("kunci cermin kehilangan %q:\n%s", wajib, q)
			}
		}
		punyaDari := strings.Contains(q, "o.STS_REJECT = :4")
		if punyaDari == dariKosong {
			t.Errorf("dariKosong=%v, syarat kode lama hadir=%v:\n%s", dariKosong, punyaDari, q)
		}
		if !strings.Contains(q, "o.STS_REJECT IS NULL") {
			t.Errorf("cermin NULL (disimpan sebelum N13) tidak ikut diselaraskan:\n%s", q)
		}
	}
	isi, err := os.ReadFile("klaimlife.go")
	if err != nil {
		t.Fatal(err)
	}
	badan := string(isi)[strings.Index(string(isi), "func (r *KlaimLife) PerbaruiStatusBaris("):]
	badan = badan[:strings.Index(badan, "\nfunc ")]
	if !strings.Contains(badan, "sqlIkutkanStatusCermin(") {
		t.Error("PerbaruiStatusBaris tidak menyelaraskan status cermin warisan")
	}
}
