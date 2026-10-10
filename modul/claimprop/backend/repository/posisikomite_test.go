package repository

import (
	"strings"
	"testing"
)

// Work owner 09-10-2026 ("posisition itu udh di setting?"): kasus komite TKMT- lahir dengan T_WORK_CLAIM.POSITION
// kosong. POSITION = KomiteID tingkat pertama (workbasket sejak migrasi 537), seperti POSITION klaim = workbasket
// pemegang tahapnya.
func TestPosisiKasusKomiteLahirDiTingkatPertama(t *testing.T) {
	tangga := []AnggotaTangga{{Urut: 2, OperatorID: "UJI-WB-2"}, {Urut: 1, OperatorID: "UJI-WB-1"},
		{Urut: 3, OperatorID: "UJI-WB-3"}}
	if p := posisiAwal(tangga); p != "UJI-WB-1" {
		t.Fatalf("POSITION kasus komite = %q, mau KomiteID tingkat pertama UJI-WB-1", p)
	}
	if p := posisiAwal(nil); p != "" {
		t.Fatalf("tangga kosong: POSITION %q, mau kosong", p)
	}
	q := sqlSisipKasusKomite("UJI.T_WORK_CLAIM")
	if !strings.Contains(q, "POSITION") || strings.Count(q, ":") != 9 {
		t.Fatalf("kelahiran kasus komite harus menulis POSITION (9 bind): %s", q)
	}
}
