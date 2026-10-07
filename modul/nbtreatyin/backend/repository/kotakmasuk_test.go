package repository

// Pembaca pemegang workbasket untuk NBStatus (`[keputusan work owner 06-10-2026]`): hanya akun AKTIF
// (M_LOGIN_GO.IS_ACTIVE = '1') yang dihitung; nama workbasket dari M_WORKBASKET.NAME; nilai selalu terikat.
// Bentuk query dicoba di DEV (06-10-2026) sebelum dipasang: subquery skalar di samping COUNT/MAX ditolak Oracle
// (ORA-00937), maka satu agregat atas LEFT JOIN.

import (
	"strings"
	"testing"

	"nusantarare/inti/backend/db"
)

func TestSqlPemegangKotakMasuk(t *testing.T) {
	q := sqlPemegangKotakMasuk("UJI.M_LOGIN_GO_WORKBASKET", "UJI.M_LOGIN_GO", "UJI.M_WORKBASKET")
	satuBaris := strings.Join(strings.Fields(q), " ")
	harap := "SELECT COUNT(g.LOGIN_ID), MAX(g.NAME), MAX(w.NAME) FROM UJI.M_WORKBASKET w" +
		" LEFT JOIN UJI.M_LOGIN_GO_WORKBASKET l ON l.WORKBASKET_ID = w.WORKBASKET_ID" +
		" LEFT JOIN UJI.M_LOGIN_GO g ON g.LOGIN_ID = l.LOGIN_ID AND g.IS_ACTIVE = :1" +
		// akun divisi IT tidak dihitung (keputusan work owner 06-10-2026)
		" AND (g.DIVISION_CODE IS NULL OR g.DIVISION_CODE <> :2) WHERE w.WORKBASKET_ID = :3"
	if satuBaris != harap {
		t.Fatalf("query:\n %s\nharap:\n %s", satuBaris, harap)
	}
	if strings.Count(q, "SELECT") != 1 {
		t.Errorf("subquery di daftar SELECT agregat (ORA-00937): %s", q)
	}
	if err := db.PeriksaSQL(q); err != nil {
		t.Fatal(err)
	}
}
