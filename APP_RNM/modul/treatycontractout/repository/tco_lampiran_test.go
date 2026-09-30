package repository

// Uji SQL lampiran di tabel WARISAN - TANPA Oracle (tiket 12, tco4).

import (
	"context"
	"database/sql"
	"strings"
	"testing"
	"time"

	"nusantarare/inti/db"
	"nusantarare/modul/treatycontractout/models"
)

// Kolom dan saringan VERBATIM RDB warisan.
func TestSQLLampiranWarisanTCO(t *testing.T) {
	d := sqlDaftarLampiranTCO("S.A", "S.I", "S.O", true)
	for _, mau := range []string{"a.ID, a.FILENAME, a.FILEMIMETYPE, a.CATEGORY_ID, a.T_STORAGE_ID",
		"FROM S.A a", "WHERE a.TREATYID = :3 AND a.ID = :4", "ORDER BY a.ID DESC",
		"(SELECT MAX(s.IMAGEID) FROM S.I s WHERE s.IMAGEID = a.T_STORAGE_ID)"} {
		if !strings.Contains(d, mau) {
			t.Errorf("daftar tanpa %q:\n%s", mau, d)
		}
	}
	// InsertAttachment2_Sql b84 (Treaty In) - sembilan kolom, DATA_JSON NULL.
	if s := sqlSisipLampiranTCO("S.A"); !strings.Contains(s,
		"(ID, TREATYID, CATEGORY, FILENAME, FILEMIMETYPE, DATA_JSON, USERNAME, CATEGORY_ID, T_STORAGE_ID)") ||
		!strings.Contains(s, "VALUES (:1, :2, :3, :4, :5, NULL, :6, :7, :8)") {
		t.Errorf("sisip: %s", s)
	}
	// DeleteAttachment2_Sql b84.
	if h := sqlHapusLampiranTCO("S.A"); !strings.Contains(h, "WHERE TREATYID = :1 AND ID = :2") {
		t.Errorf("hapus: %s", h)
	}
	// Insert_T_Storage_SQL b85 / DeleteStorage_SQL b85.
	if o := sqlSimpanObjekTCO("S.I"); !strings.Contains(o, "(IMAGEID, URLPUBLIC, APPFOLDER, EXPDATE, FILENAME, APPNAME, STORAGE)") ||
		!strings.Contains(o, "TO_DATE(:4, 'DD/MM/YYYY HH24:MI:SS')") || !strings.Contains(o, "'standard'") {
		t.Errorf("simpan objek: %s", o)
	}
	if h := sqlHapusObjekTCO("S.I"); !strings.Contains(h, "WHERE IMAGEID = :1") {
		t.Errorf("hapus objek: %s", h)
	}
	for _, q := range []string{d, sqlSisipLampiranTCO("S.A"), sqlHapusLampiranTCO("S.A"), sqlSimpanObjekTCO("S.I"),
		sqlHapusObjekTCO("S.I"), sqlKunciLampiranTCO("S.A"), sqlAmbilUntukKirimLampiranTCO("S.A", "S.I"),
		sqlAdaIDLampiranTCO("S.A"), sqlTreatyYearLampiranTCO("S.Y")} {
		if err := db.PeriksaSQL(q); err != nil {
			t.Errorf("PeriksaSQL: %v", err)
		}
	}
}

// ID = TO_CHAR(SYSTIMESTAMP, 'YYYYMMDDHH24MISSFF3') WIB - dan ia cap waktu unggahnya.
func TestIDLampiranTCODuaArah(t *testing.T) {
	w := time.Date(2026, 9, 29, 2, 3, 4, 56_000_000, time.UTC)
	id := IDLampiranTCO(w)
	if id != "20260929090304056" {
		t.Errorf("ID %q, mau 20260929090304056 (09:03 WIB)", id)
	}
	if b := waktuDariIDLampiranTCO(id); !b.Equal(w) {
		t.Errorf("waktu pulang %v, mau %v", b, w)
	}
	if !waktuDariIDLampiranTCO("1000000001").IsZero() {
		t.Error("ID berbentuk lain harus memberi waktu nol")
	}
}

// TREATYID = TreatyYear + TreatyYearID (TreatyOutSaveAttachment b1402).
func TestKunciTreatyLampiranTCO(t *testing.T) {
	if k := models.KunciTreatyLampiranTCO(" 2026 ", "1000001"); k != "20261000001" {
		t.Errorf("kunci %q", k)
	}
	var n [7]sql.NullString
	for i, v := range []string{"20260929090304056", "a.pdf", "application/pdf", "CLAUSES", "ABCDEF", "ABCDEF", "UJI-ADMIN"} {
		n[i] = sql.NullString{String: v, Valid: true}
	}
	l := pindaiLampiranTCO(n, "1000001")
	if l.Category != "CLAUSES" || l.ImageID != "ABCDEF" || l.TStorageID != "ABCDEF" || l.IDTreatyYear != "1000001" ||
		l.TglUpload.IsZero() {
		t.Errorf("pindai: %+v", l)
	}
}

// OQ-TCO-26 (lanjutan 4): `Update_T_Storage_SQL` b85 ditiru apa adanya - empat
// kolom, kunci IMAGEID, tanpa COMMIT; menuntut transaksi seperti penulis lain.
func TestSQLPerbaruiObjekTCOSepertiUpdateTStorage(t *testing.T) {
	q := sqlPerbaruiObjekTCO("S.I")
	for _, mau := range []string{"UPDATE S.I", "URLPUBLIC = :1", "APPFOLDER = :2",
		"EXPDATE = TO_DATE(:3, 'DD/MM/YYYY HH24:MI:SS')", "TANGGAL_UPLOAD = TO_DATE(:4, 'MM/DD/YYYY HH24:MI:SS')",
		"WHERE IMAGEID = :5"} {
		if !strings.Contains(q, mau) {
			t.Errorf("tanpa %q:\n%s", mau, q)
		}
	}
	if err := db.PeriksaSQL(q); err != nil {
		t.Error(err)
	}
	err := NewMasterLampiranTCO(nil).PerbaruiObjek(context.Background(), nil, models.ObjekPenyimpananTCO{ImageID: "UJI"})
	if err == nil {
		t.Error("tanpa transaksi harus galat")
	}
}
