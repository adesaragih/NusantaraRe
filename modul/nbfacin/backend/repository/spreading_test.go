package repository

import (
	"strings"
	"testing"
)

// TestSQLSpreading - tiket 48: pohon urut SEQ_NO lokasi / item / coverage; baris spreading dan hapus dibatasi coverage case
// :1; sisip 14 kolom (angka lewat TO_NUMBER); master treaty = GetTreatyName_SQL (bizcode :1 / :2, tanggal :3 / :4),
// KAPASITAS_TREATY baris pertama urut MAX_LIMIT_IDR, kurs TREATYEXCHANGEYEARLY.
func TestSQLSpreading(t *testing.T) {
	tb := tabelSpreading{work: "U.W", general: "U.G", loc: "U.L", prop: "U.P", risk: "U.R", item: "U.I", cov: "U.V", spread: "U.S"}
	cek := func(nama, q string, harus ...string) {
		t.Helper()
		for _, h := range harus {
			if !strings.Contains(q, h) {
				t.Errorf("%s tanpa %q: %s", nama, h, q)
			}
		}
		if strings.Contains(strings.ToUpper(q), "COMMIT") || strings.Contains(strings.ToUpper(q), "CALL ") {
			t.Errorf("%s memuat COMMIT / CALL: %s", nama, q)
		}
	}
	cek("pohon", sqlPohonSpreading(tb), "FROM U.L l", "LEFT JOIN U.V v ON v.PARENT_ID = i.ID AND "+syaratIndukCoverage,
		"WHERE l.PARENT_ID = :1", "ORDER BY l.SEQ_NO, i.SEQ_NO, v.SEQ_NO", "v.OLDID", "p.IS_TOP_RISK")
	if n := len(kolomPohonSpreading); n != 23 {
		t.Errorf("kolom pohon = %d", n)
	}
	cek("baris", sqlBacaBarisSpreading(tb), "FROM U.S s WHERE s.PARENT_ID IN (SELECT v.ID FROM U.V v", "l.PARENT_ID = :1",
		"ORDER BY s.PARENT_ID, s.SEQ_NO")
	cek("hapus", sqlHapusSpreading("U.S", "U.V", "U.I", "U.P", "U.L"), "DELETE FROM U.S WHERE PARENT_ID IN (SELECT v.ID FROM U.V v")
	cek("share", sqlUbahShare("U.G"), "UPDATE U.G SET PERCENT_SHARE = "+angkaMasuk(":1")+" WHERE ID = :2")
	cek("nr", sqlUbahNR("U.V"), "TSI_NUSANTARA_RE = "+angkaMasuk(":1"), "PREMI_NUSANTARA_RE = "+angkaMasuk(":2"), "WHERE ID = :3")
	sisip := sqlSisipSpreading("U.S")
	cek("sisip", sisip, "INSERT INTO U.S (ID, PARENT_ID, PARENT_TABLE, SRC_PATH, SEQ_NO, ROW_UID, TREATY_TYPE, TREATY_NAME,",
		angkaMasuk(":9"), angkaMasuk(":13"), ":14)")
	if strings.Contains(sisip, ":15") {
		t.Errorf("sisip lebih dari 14 bind: %s", sisip)
	}
	cek("treaty", sqlDaftarTreaty("U.PA", "U.RT", "U.TB", "U.TC"), "FROM U.PA a", "TREATYDESCNAME = 'TREATY LIMIT'",
		"b.BIZCODE = :1", "b.BIZCODE = :2", "TO_DATE(:3, 'DD/MM/RRRR')", "TO_DATE(:4, 'DD/MM/RRRR')", "ORDER BY CARI5 ASC")
	cek("kapasitas", sqlKapasitasTreaty("U.KT"), "LIMIT_BOTTOM_IDR <= "+angkaMasuk(":1"), "MAX_LIMIT_IDR >= "+angkaMasuk(":2"),
		"STARTDATE <= TO_DATE(:3, 'DD/MM/YYYY')", "ENDDATE >= TO_DATE(:4, 'DD/MM/YYYY')", "ORDER BY MAX_LIMIT_IDR ASC FETCH FIRST 1 ROWS ONLY")
	cek("kurs terbaru", sqlKursTerbaru("U.X"), "x.IDCURRENCY = :1 ORDER BY x.STARTDATE DESC) WHERE ROWNUM = 1")
	cek("kurs pada", sqlKursPada("U.X"), "x.IDCURRENCY = :1 AND :2 BETWEEN x.STARTDATE AND x.ENDDATE")
	cek("id mata uang", sqlIDMataUang("U.C"), "WHERE c.CURRENCY = :1")
	if q := sqlCatatanJenisTreaty("U.RT"); q != "SELECT NOTE FROM U.RT WHERE TO_CHAR(ID) = :1" {
		t.Errorf("catatan jenis treaty: %s", q)
	}
}
