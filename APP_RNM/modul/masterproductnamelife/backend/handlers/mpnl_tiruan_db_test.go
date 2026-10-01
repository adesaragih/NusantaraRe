//go:build db

package handlers_test

// Tiruan tabel WARISAN Master Product Name Life di skema uji Oracle (uji
// bertag `db`). Tanpa ORACLE_DSN seluruhnya MELEWATI. ⛔ POOLDATA bukan
// sasaran uji `db` (brief bab 1) - skema uji dipagari `uji/skemauji`.
//
// ⛔ Kolom dan tipe PERSIS katalog DEV `ALL_TAB_COLUMNS` 01-10-2026
// (`repository/testdata/katalog-dev.json`): `ID VARCHAR2(6)`, `JSONDATA CLOB` +
// constraint `IS JSON`, kolom datar `RIRISKID VARCHAR2(10)`, `RIRISK VARCHAR2(100)` -
// hanya dua (lanjutan 1 L1: `PRODUCTNAME`/`BEGIN_DATE` tidak ada di DEV); NOL PK.
//
// Jalankan per paket: `go test -tags db -p 1 ./modul/masterproductnamelife/...`.

import (
	"context"
	"database/sql"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/masterproductnamelife/backend/handlers"
	"nusantarare/modul/masterproductnamelife/backend/services"
	"nusantarare/uji/skemauji"
)

// ddlTiruan - kolom dan tipe per tabel (urutan pembuatan).
var ddlTiruan = []struct{ nama, kolom string }{
	{"M_PRODUCT_LIFE", `ID VARCHAR2(6), JSONDATA CLOB CONSTRAINT UJI_MPL_JSON CHECK (JSONDATA IS JSON),
		RIRISKID VARCHAR2(10), RIRISK VARCHAR2(100)`},
	{"M_PRODUCTINWARD_LIFE", `ID VARCHAR2(6), JSONDATA CLOB CONSTRAINT UJI_MPIL_JSON CHECK (JSONDATA IS JSON)`},
	// Master pendukung - kolom yang dibaca modul ini saja, VARCHAR2 generik (paket 2).
	{"AGENT", `ID VARCHAR2(100), CLIENTNAME VARCHAR2(1000), STATUSACTIVE VARCHAR2(10)`},
	{"CLIENT", `ID VARCHAR2(100), NAME VARCHAR2(1000), BU_NOTE VARCHAR2(100)`},
	{"CURRENCY", `ID VARCHAR2(100), CURRENCY VARCHAR2(100)`},
	{"RIRISK_LIFE_SUMMARY", `ID VARCHAR2(100), USEDBY VARCHAR2(1000)`},
	{"CAUSEOFLOSS_LIFE", `ID VARCHAR2(100), CAUSEOFLOSS VARCHAR2(1000)`},
	{"PRODUCT_TYPE_LIFE", `ID VARCHAR2(100), COVERNAME VARCHAR2(1000), BUSINESS VARCHAR2(1000), BENEFIT VARCHAR2(4000)`},
}

// sequenceTiruan - sequence yang dipakai penulis modul (P6).
var sequenceTiruan = []string{"M_PRODUCT_LIFE_SEQ"}

type ujiDB struct {
	srv    *httptest.Server
	mentah *db.DB
	skema  string
	ctx    context.Context
}

// pasangDB membangun tiruan dari nol dan server HTTP modul di atas Oracle.
//
// ⚠️ Satu koneksi saja - `skemauji.BukaRepositori()` - untuk tiruan dan
// server; `BolehDilewati` tetap diperiksa.
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
	for _, d := range ddlTiruan {
		u.exec(t, fmt.Sprintf(`CREATE TABLE %s.%s (%s)`, u.skema, d.nama, d.kolom))
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
	for _, d := range ddlTiruan {
		if _, err := u.mentah.ExecContext(u.ctx, fmt.Sprintf(`DROP TABLE %s.%s PURGE`, u.skema, d.nama)); err != nil &&
			!strings.Contains(err.Error(), "ORA-00942") {
			t.Fatalf("membongkar %s: %v", d.nama, err)
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
	if _, err := u.mentah.ExecContext(u.ctx, strings.ReplaceAll(q, "{s}", u.skema), args...); err != nil {
		t.Fatalf("%s: %v", strings.Fields(q)[0], err)
	}
}

// teks membaca satu nilai sebagai teks.
func (u *ujiDB) teks(t *testing.T, q string, args ...any) string {
	t.Helper()
	var v sql.NullString
	if err := u.mentah.QueryRowContext(u.ctx, strings.ReplaceAll(q, "{s}", u.skema), args...).Scan(&v); err != nil {
		t.Fatalf("membaca %q: %v", q, err)
	}
	return v.String
}
