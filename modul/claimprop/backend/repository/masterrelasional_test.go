package repository

// Master treaty Choose Master dari tabel relasional, bukan JSON (perintah work owner 10-10-2026 "perbaiki cek choose
// master. jangan ambil dari json"): `MasterTreaty` / `bacaMaster` membaca TREATY_IN / TREATY_IN_EDM dan view
// TREATYINDETAILJOINEDM (+ Bordereaux / Bordereaux Note / Accounting Mode dari tabel flat T_TREATY_REVISION), dan tidak
// menyebut M_TREATY_IN maupun JSON sama sekali.

import (
	"os"
	"strings"
	"testing"
)

func TestMasterTreatyTidakMembacaJSON(t *testing.T) {
	isi, err := os.ReadFile("acuan.go")
	if err != nil {
		t.Fatal(err)
	}
	s := string(isi)
	awal := strings.Index(s, "func (a *Acuan) MasterTreaty(")
	if awal < 0 {
		t.Fatal("MasterTreaty tidak ditemukan")
	}
	sisa := s[awal:]
	akhir := strings.Index(sisa, "// tglPega")
	if akhir < 0 {
		t.Fatal("ujung bacaMaster tidak ditemukan")
	}
	badan := sisa[:akhir]
	for _, w := range []string{`"TREATY_IN", "TREATY_IN_EDM"`, `a.q("TREATYINDETAILJOINEDM")`, "GROUP BY TREATYTYPE, TREATYGROUPID, RNM_SHARE",
		`a.q("T_TREATY_REVISION")`, "SELECT BORDEAUX, BORDEREAUXNOTE, ACCOUNTINGMODE FROM %s WHERE MASTERID = :1"} {
		if !strings.Contains(badan, w) {
			t.Errorf("MasterTreaty tanpa %q", w)
		}
	}
	for _, larang := range []string{"JSON", "M_TREATY_IN\"", "CashLoss", "SpreadingList"} {
		if strings.Contains(badan, larang) {
			t.Errorf("MasterTreaty masih menyebut %q", larang)
		}
	}
}
