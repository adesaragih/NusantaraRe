//go:build db

package handlers_test

// `On Retention` → `OutwardList` terhadap skema uji Oracle NYATA (paket 9):
// `BrowseReinstypeOR_SQL` b84 atas `TREATYCONTRACT_LIFE` × `TREATYYEAR_LIFE`.
//
// ⚠️ Kedua tabel dipakai bersama (Master Contract Retro Life menirunya,
// `uji/skemauji` meniru `TREATYYEAR_LIFE` untuk Claim Life - bentuk sama): dibuat
// HANYA bila belum ada; baris uji ber-ID `UJI-MPNL-*` dibersihkan.

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
)

var ddlKontrakTreaty = []struct{ nama, ddl string }{
	{"TREATYYEAR_LIFE", `ID VARCHAR2(100), TREATYYEAR VARCHAR2(100), UNDERWRITINGYEAR VARCHAR2(100),
		USERID VARCHAR2(100), TGLUPDATE DATE, STARTDATE DATE, ENDDATE DATE`},
	{"TREATYCONTRACT_LIFE", `ID VARCHAR2(100), IDTREATYYEAR VARCHAR2(100), REINSTYPEID VARCHAR2(100),
		REINSTYPENAME VARCHAR2(1000), USERID VARCHAR2(100), TGLUPDATE DATE, IDR NUMBER, USD NUMBER, B_IDR NUMBER,
		B_USD NUMBER, IDR_SELISIH NUMBER, USD_SELISIH NUMBER, TREATYENDDATE DATE, TREATYSTARTDATE DATE`},
}

func (u *ujiDB) pasangKontrakTreaty(t *testing.T) {
	t.Helper()
	for _, d := range ddlKontrakTreaty {
		if u.ada(t, "TABLE", d.nama) {
			continue
		}
		u.exec(t, fmt.Sprintf(`CREATE TABLE %s.%s (%s)`, u.skema, d.nama, d.ddl))
		nama := d.nama
		t.Cleanup(func() { _, _ = u.mentah.ExecContext(u.ctx, fmt.Sprintf(`DROP TABLE %s.%s PURGE`, u.skema, nama)) })
	}
	bersih := func() {
		for _, tabel := range []string{"TREATYCONTRACT_LIFE", "TREATYYEAR_LIFE"} {
			_, _ = u.mentah.ExecContext(u.ctx, fmt.Sprintf(`DELETE FROM %s.%s WHERE ID LIKE 'UJI-MPNL-%%'`, u.skema, tabel))
		}
	}
	bersih()
	t.Cleanup(bersih)
	u.exec(t, `INSERT INTO {s}.TREATYYEAR_LIFE (ID, TREATYYEAR, UNDERWRITINGYEAR) VALUES ('UJI-MPNL-TY1', '2025', '2025')`)
	for _, k := range [][4]string{
		{"UJI-MPNL-TC1", "10200", "UJI OR", "2036-12-31"},
		{"UJI-MPNL-TC2", "10196", "UJI QS", "2036-12-31"},
		{"UJI-MPNL-TC3", "10200", "UJI OR LAMA", "2026-06-30"},
	} {
		u.exec(t, `INSERT INTO {s}.TREATYCONTRACT_LIFE (ID, IDTREATYYEAR, REINSTYPEID, REINSTYPENAME, TREATYSTARTDATE, TREATYENDDATE)
			VALUES (:1, 'UJI-MPNL-TY1', :2, :3, DATE '2025-01-01', TO_DATE(:4, 'YYYY-MM-DD'))`, k[0], k[1], k[2], k[3])
	}
}

func TestDBOnRetentionOutwardListDariKontrakOR(t *testing.T) {
	u := pasangDB(t)
	u.isiMasterUji(t)
	u.pasangKontrakTreaty(t)
	badan := strings.Replace(badanLengkap, `{"umum":`, `{"hitungOutward":true,"umum":`, 1)
	kode, jawab := u.kirim(t, "POST", pre+"/produk", badan)
	if kode != http.StatusOK {
		t.Fatalf("simpan: %d %s", kode, jawab)
	}
	// Periode BEGIN–MATURE badanLengkap di dalam UJI-MPNL-TC1 saja (TC2 bukan OR, TC3 berakhir lebih dulu).
	if !strings.Contains(jawab, `"reinsTypeName":"UJI OR"`) || strings.Contains(jawab, "UJI QS") {
		t.Errorf("OutwardList: %s", jawab)
	}
	if v := u.teks(t, `SELECT REINSTYPEID FROM {s}.M_PRODUCTNAME_LIFE_OUTWARD WHERE PRODUCTID = '100044' AND URUT = 1`); v != "10200" {
		t.Errorf("baris OUTWARD pertama REINSTYPEID: %q", v)
	}
	if v := u.teks(t, `SELECT TO_CHAR(TRANSACTIONYEAR) FROM {s}.M_PRODUCTNAME_LIFE_OUTWARD WHERE PRODUCTID = '100044' AND URUT = 1`); v != "2025" {
		t.Errorf("TRANSACTIONYEAR ← TREATYYEAR: %q", v)
	}
}
