package repository

import "testing"

// TestSQLTableOfLimit - tiket 40: BIZCODE persis, TANPA TAHUN (A161), CATEGORY hanya bila diisi (filter RD dibuang bila
// kosong), DISTINCT atas kolom laporan RD, urut CATEGORY lalu DESCRIPTION; BUSINESS dicocokkan NOTE + group.
func TestSQLTableOfLimit(t *testing.T) {
	if q := sqlTableOfLimit("UJI.T", true); q != "SELECT DESCRIPTION, PCTLIMIT FROM (SELECT DISTINCT BIZCODE, CATEGORY, DESCRIPTION, PCTLIMIT, NOTE FROM UJI.T WHERE BIZCODE = :1 AND CATEGORY = :2) ORDER BY CATEGORY, DESCRIPTION, PCTLIMIT, NOTE FETCH FIRST :3 ROWS ONLY" {
		t.Errorf("dengan kategori: %q", q)
	}
	if q := sqlTableOfLimit("UJI.T", false); q != "SELECT DESCRIPTION, PCTLIMIT FROM (SELECT DISTINCT BIZCODE, CATEGORY, DESCRIPTION, PCTLIMIT, NOTE FROM UJI.T WHERE BIZCODE = :1) ORDER BY CATEGORY, DESCRIPTION, PCTLIMIT, NOTE FETCH FIRST :2 ROWS ONLY" {
		t.Errorf("tanpa kategori: %q", q)
	}
	if q := sqlKodeBisnis("UJI.B"); q != "SELECT ID FROM UJI.B WHERE NOTE = :1 AND BUSINESSGROUPID = :2 ORDER BY ID FETCH FIRST 2 ROWS ONLY" {
		t.Errorf("kode bisnis: %q", q)
	}
	if BatasTableOfLimit != 500 || TabelTableOfLimitWarisan != "TABLEOFLIMIT" {
		t.Error("parameter RD / nama tabel DDL berubah")
	}
}
