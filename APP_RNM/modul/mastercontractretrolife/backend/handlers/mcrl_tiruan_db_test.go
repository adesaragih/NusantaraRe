//go:build db

package handlers_test

// Tiruan tabel WARISAN dan MASTER Master Contract Retro Life di skema uji
// Oracle (uji bertag `db`). Tanpa ORACLE_DSN seluruhnya MELEWATI.
//
// ⛔ Tipe kolom PERSIS DDL `[data DBA]` (`docs/ddl-tables-from-dba.md`): uang
// dan persen NUMBER tanpa presisi, tanggal DATE, `RIRATE` VARCHAR2(1000), dan
// NOL constraint (K1 - DEV nol P/R/U). Master ditiru dengan kolom yang dibaca
// modul ini saja, VARCHAR2 generik.
//
// ⚠️ `uji/skemauji` juga meniru `TREATYYEAR_LIFE` untuk Claim Life (bentuk
// sama). Jalankan uji `db` per paket: `go test -tags db -p 1 ./modul/mastercontractretrolife/...`.

import (
	"context"
	"database/sql"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/mastercontractretrolife/backend/handlers"
	"nusantarare/modul/mastercontractretrolife/backend/repository"
	"nusantarare/modul/mastercontractretrolife/backend/services"
	"nusantarare/uji/skemauji"
)

// ddlTiruan - kolom dan tipe per tabel.
var ddlTiruan = map[string]string{
	repository.TabelTahun: `ID VARCHAR2(100), TREATYYEAR VARCHAR2(100), UNDERWRITINGYEAR VARCHAR2(100),
		USERID VARCHAR2(100), TGLUPDATE DATE, STARTDATE DATE, ENDDATE DATE`,
	repository.TabelKontrak: `ID VARCHAR2(100), IDTREATYYEAR VARCHAR2(100), REINSTYPEID VARCHAR2(100),
		REINSTYPENAME VARCHAR2(1000), USERID VARCHAR2(100), TGLUPDATE DATE, IDR NUMBER, USD NUMBER, B_IDR NUMBER,
		B_USD NUMBER, IDR_SELISIH NUMBER, USD_SELISIH NUMBER, TREATYENDDATE DATE, TREATYSTARTDATE DATE`,
	repository.TabelReinsurer: `ID VARCHAR2(100), TREATYYEARID VARCHAR2(100), TREATYCONTRACTID VARCHAR2(100),
		REINSTYPEID VARCHAR2(100), REINSTYPENAME VARCHAR2(1000), REINSURERNAME VARCHAR2(1000), PCTSHARE NUMBER,
		COMMISION NUMBER, OVR_COMM NUMBER, USERID VARCHAR2(1000), TGLUPDATE DATE, REINSURERID VARCHAR2(100)`,
	repository.TabelSecurity: `ID VARCHAR2(100), TREATYYEARID VARCHAR2(100), TREATYCONTRACTID VARCHAR2(100),
		TREATYREINSURERID VARCHAR2(100), REINSURERID VARCHAR2(100), REINSURERNAME VARCHAR2(1000), PCTSHARE NUMBER,
		USERID VARCHAR2(1000), TGLUPDATE DATE`,
	repository.TabelBusiness: `ID VARCHAR2(100), TREATYYEARID VARCHAR2(100), TREATYYEAR VARCHAR2(100),
		REINSTYPEID VARCHAR2(100), REINSTYPENAME VARCHAR2(1000), BIZCODE VARCHAR2(100), BIZNAME VARCHAR2(1000),
		USERID VARCHAR2(100), TGLUPDATE DATE, TREATYCONTRACTID VARCHAR2(100), RIRATEID VARCHAR2(100),
		RIRATE VARCHAR2(1000)`,
	repository.MasterJenisReasuransi: `ID VARCHAR2(100), NOTE VARCHAR2(1000), FLAG VARCHAR2(100)`,
	repository.MasterReinsurer:       `ID VARCHAR2(100), CLIENTNAME VARCHAR2(1000), STATUSACTIVE VARCHAR2(10)`,
	repository.MasterBusiness:        `ID VARCHAR2(100), NOTE VARCHAR2(1000), OLDID VARCHAR2(100)`,
}

// sequenceTiruan - lima sequence yang ditiru penulis modul (K3).
var sequenceTiruan = []string{"TREATYYEAR_LIFE_SEQ", "TREATYCONTRACT_LIFE_SEQ", "TREATYREINSURER_LIFE_SEQ",
	"TREATYSECURITYREINSURER_LIFE_SEQ", "TREATYBUSINESS_LIFE_SEQ"}

type ujiDB struct {
	srv    *httptest.Server
	mentah *db.DB
	skema  string
	ctx    context.Context
}

// pasangDB membangun tiruan dari nol dan server HTTP modul di atas Oracle.
//
// ⚠️ Satu koneksi saja - `skemauji.BukaRepositori()` - untuk tiruan dan untuk
// server: pemanggil `skemauji.Buka()` dicacah penjaga Claim Life
// (`TestSetiapPemanggilBukaMemeriksaBolehDilewati`), dan tidak ada yang perlu
// dibuka dua kali. `BolehDilewati` tetap diperiksa.
func pasangDB(t *testing.T) *ujiDB {
	t.Helper()
	repoDB, err := skemauji.BukaRepositori()
	if err != nil {
		if !skemauji.BolehDilewati(err) {
			t.Fatalf("skema uji menolak: %v", err)
		}
		t.Skipf("lewati: %v", err)
	}
	ctx := context.Background()
	if err := repoDB.Ping(ctx); err != nil {
		t.Skipf("lewati: oracle tidak terjangkau: %v", err)
	}
	u := &ujiDB{mentah: repoDB, skema: repoDB.Skema(), ctx: ctx}
	u.bongkar(t)
	for nama, kolom := range ddlTiruan {
		u.exec(t, fmt.Sprintf(`CREATE TABLE %s.%s (%s)`, u.skema, nama, kolom))
	}
	for _, s := range sequenceTiruan {
		u.exec(t, fmt.Sprintf(`CREATE SEQUENCE %s.%s START WITH 44 NOCACHE`, u.skema, s))
	}
	u.srv = httptest.NewServer(handlers.Router(services.New(repoDB), true))
	t.Cleanup(func() {
		u.srv.Close()
		u.bongkar(t)
		_ = repoDB.Close()
	})
	return u
}

func (u *ujiDB) bongkar(t *testing.T) {
	t.Helper()
	for nama := range ddlTiruan {
		if _, err := u.mentah.ExecContext(u.ctx, fmt.Sprintf(`DROP TABLE %s.%s PURGE`, u.skema, nama)); err != nil &&
			!strings.Contains(err.Error(), "ORA-00942") {
			t.Fatalf("membongkar %s: %v", nama, err)
		}
	}
	for _, s := range sequenceTiruan {
		if _, err := u.mentah.ExecContext(u.ctx, fmt.Sprintf(`DROP SEQUENCE %s.%s`, u.skema, s)); err != nil &&
			!strings.Contains(err.Error(), "ORA-02289") {
			t.Fatalf("membongkar %s: %v", s, err)
		}
	}
}

func (u *ujiDB) exec(t *testing.T, q string, args ...any) {
	t.Helper()
	if _, err := u.mentah.ExecContext(u.ctx, q, args...); err != nil {
		t.Fatalf("%s: %v", strings.Fields(q)[0], err)
	}
}

// cacah menghitung baris sebuah tabel tiruan dengan saringan opsional.
func (u *ujiDB) cacah(t *testing.T, tabel, saring string, args ...any) int {
	t.Helper()
	q := fmt.Sprintf(`SELECT COUNT(*) FROM %s.%s`, u.skema, tabel)
	if saring != "" {
		q += " WHERE " + saring
	}
	var n int
	if err := u.mentah.QueryRowContext(u.ctx, q, args...).Scan(&n); err != nil {
		t.Fatalf("mencacah %s: %v", tabel, err)
	}
	return n
}

// teks membaca satu nilai sebagai teks (angka lewat TO_CHAR TM9).
func (u *ujiDB) teks(t *testing.T, q string, args ...any) string {
	t.Helper()
	var v sql.NullString
	if err := u.mentah.QueryRowContext(u.ctx, strings.ReplaceAll(q, "{s}", u.skema), args...).Scan(&v); err != nil {
		t.Fatalf("membaca %q: %v", q, err)
	}
	return v.String
}
