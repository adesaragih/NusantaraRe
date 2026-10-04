//go:build db

package handlers_test

// Seam HTTP lampiran terhadap skema uji Oracle NYATA (paket 8): rekam
// `M_ATTACHMENTPRODUCTNAME`, objek `T_STORAGE_IMAGE`, outbox `T_LOG_SERVICE_RNM`.
//
// ⚠️ `T_LOG_SERVICE_RNM` (migrasi Claim Life 015) dan `T_STORAGE_IMAGE` (tiruan
// Treaty Contract Out) dipakai bersama modul lain di skema uji: dibuat HANYA bila
// belum ada (dan dibongkar hanya bila uji ini yang membuatnya); baris milik uji ini
// dibersihkan menurut penandanya.

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"nusantarare/modul/masterproductnamelife/backend/handlers"
	"nusantarare/modul/masterproductnamelife/backend/services"
)

var ddlBersamaLampiran = []struct{ nama, ddl string }{
	{"M_ATTACHMENTPRODUCTNAME", `ID VARCHAR2(100), TREATYID VARCHAR2(100), CATEGORY VARCHAR2(100), FILENAME VARCHAR2(1000),
		FILEMIMETYPE VARCHAR2(100), DATA_JSON CLOB, USERNAME VARCHAR2(100), T_STORAGE_ID VARCHAR2(100)`},
	{"T_STORAGE_IMAGE", `IMAGEID VARCHAR2(100), URLPUBLIC VARCHAR2(4000), APPFOLDER VARCHAR2(1000), EXPDATE DATE,
		FILENAME VARCHAR2(1000), APPNAME VARCHAR2(100), STORAGE VARCHAR2(100), TANGGAL_UPLOAD DATE`},
	{"T_FOLDER_IMAGE", `APPNAME VARCHAR2(100)`},
	// Bentuk migrasi Claim Life 015 (tanpa index).
	{"T_LOG_SERVICE_RNM", `ID VARCHAR2(40) NOT NULL, LINI VARCHAR2(16) NOT NULL, MODUL VARCHAR2(32) NOT NULL,
		JENIS_EFEK VARCHAR2(32) NOT NULL, RUJUKAN VARCHAR2(40) NOT NULL, MUATAN CLOB, STATUS VARCHAR2(16) NOT NULL,
		PERCOBAAN NUMBER(5) DEFAULT 0 NOT NULL, JADWAL_BERIKUT TIMESTAMP, GALAT_TERAKHIR VARCHAR2(4000),
		DIBUAT TIMESTAMP NOT NULL, DIPERBARUI TIMESTAMP`},
}

func (u *ujiDB) ada(t *testing.T, jenis, nama string) bool {
	t.Helper()
	q := `SELECT COUNT(*) FROM ALL_TABLES WHERE OWNER = :1 AND TABLE_NAME = :2`
	if jenis == "SEQUENCE" {
		q = `SELECT COUNT(*) FROM ALL_SEQUENCES WHERE SEQUENCE_OWNER = :1 AND SEQUENCE_NAME = :2`
	}
	var n int
	if err := u.mentah.QueryRowContext(u.ctx, q, strings.ToUpper(u.skema), nama).Scan(&n); err != nil {
		t.Fatalf("katalog %s: %v", nama, err)
	}
	return n > 0
}

// pasangLampiran menyiapkan tabel lampiran bersama dan server dengan stub folder sementara.
func (u *ujiDB) pasangLampiran(t *testing.T) *httptest.Server {
	t.Helper()
	for _, d := range ddlBersamaLampiran {
		if u.ada(t, "TABLE", d.nama) {
			continue
		}
		u.exec(t, fmt.Sprintf(`CREATE TABLE %s.%s (%s)`, u.skema, d.nama, d.ddl))
		nama := d.nama
		t.Cleanup(func() { _, _ = u.mentah.ExecContext(u.ctx, fmt.Sprintf(`DROP TABLE %s.%s PURGE`, u.skema, nama)) })
	}
	if !u.ada(t, "SEQUENCE", "SEQ_LOG_SERVICE_RNM") {
		u.exec(t, fmt.Sprintf(`CREATE SEQUENCE %s.SEQ_LOG_SERVICE_RNM START WITH 1 NOCACHE`, u.skema))
		t.Cleanup(func() {
			_, _ = u.mentah.ExecContext(u.ctx, fmt.Sprintf(`DROP SEQUENCE %s.SEQ_LOG_SERVICE_RNM`, u.skema))
		})
	}
	bersih := func() {
		_, _ = u.mentah.ExecContext(u.ctx, fmt.Sprintf(`DELETE FROM %s.T_STORAGE_IMAGE WHERE IMAGEID IN
			(SELECT T_STORAGE_ID FROM %s.M_ATTACHMENTPRODUCTNAME WHERE TREATYID LIKE 'UJI%%' OR TREATYID = '100044')`, u.skema, u.skema))
		_, _ = u.mentah.ExecContext(u.ctx, fmt.Sprintf(`DELETE FROM %s.M_ATTACHMENTPRODUCTNAME WHERE TREATYID = '100044'`, u.skema))
		_, _ = u.mentah.ExecContext(u.ctx, fmt.Sprintf(`DELETE FROM %s.T_LOG_SERVICE_RNM WHERE MODUL = 'MASTERPRODUCTNAMELIFE'`, u.skema))
		_, _ = u.mentah.ExecContext(u.ctx, fmt.Sprintf(`DELETE FROM %s.T_FOLDER_IMAGE WHERE APPNAME = 'UJI-APP'`, u.skema))
	}
	bersih()
	t.Cleanup(bersih)
	u.exec(t, `INSERT INTO {s}.T_FOLDER_IMAGE (APPNAME) VALUES ('UJI-APP')`)
	l := services.LayananOracle(services.New(u.mentah)).DenganPenyimpanan(services.PenyimpananLokal(t.TempDir()))
	srv := httptest.NewServer(handlers.RouterDengan(l, true, true))
	t.Cleanup(srv.Close)
	return srv
}

func TestDBLampiranRekamObjekOutboxLaluHapus(t *testing.T) {
	u := pasangDB(t)
	u.isiMasterUji(t)
	if kode, badan := u.kirim(t, "POST", pre+"/produk", badanLengkap); kode != http.StatusOK {
		t.Fatalf("produk: %d %s", kode, badan)
	}
	srv := u.pasangLampiran(t)
	var buf bytes.Buffer
	m := multipart.NewWriter(&buf)
	w, _ := m.CreateFormFile("berkas", "UJI nota.pdf")
	_, _ = w.Write([]byte("ISI"))
	_ = m.Close()
	req, _ := http.NewRequest("POST", srv.URL+pre+"/produk/100044/lampiran", &buf)
	req.Header.Set("Content-Type", m.FormDataContentType())
	req.Header.Set("X-Pelaku", "UJI-PELAKU")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	b, _ := io.ReadAll(res.Body)
	_ = res.Body.Close()
	if res.StatusCode != http.StatusOK || !strings.Contains(string(b), `"status":"terunggah"`) {
		t.Fatalf("unggah: %d %s", res.StatusCode, b)
	}
	var a struct{ ID, StorageID string }
	_ = json.Unmarshal(b, &a)
	if n := u.cacah(t, "M_ATTACHMENTPRODUCTNAME", "ID = :1 AND TREATYID = '100044' AND CATEGORY = 'File' "+
		"AND FILEMIMETYPE = 'pdf' AND USERNAME = 'UJI-PELAKU' AND T_STORAGE_ID = :2", a.ID, a.StorageID); n != 1 {
		t.Errorf("rekam InsertAttachProdName_Sql: %d", n)
	}
	if n := u.cacah(t, "T_STORAGE_IMAGE", "IMAGEID = :1 AND URLPUBLIC IS NULL AND STORAGE = 'standard' "+
		"AND APPNAME = 'UJI-APP' AND APPFOLDER LIKE 'Contract/Doc/%'", a.StorageID); n != 1 {
		t.Errorf("objek Insert_T_Storage_SQL tanpa alamat: %d", n)
	}
	if n := u.cacah(t, "T_LOG_SERVICE_RNM", "MODUL = 'MASTERPRODUCTNAMELIFE' AND RUJUKAN = :1 AND STATUS = 'selesai'", a.ID); n != 1 {
		t.Errorf("efek keluar tuntas: %d", n)
	}
	req, _ = http.NewRequest("DELETE", srv.URL+pre+"/produk/100044/lampiran/"+a.ID, nil)
	req.Header.Set("X-Pelaku", "UJI-PELAKU")
	if res, err := http.DefaultClient.Do(req); err != nil || res.StatusCode != http.StatusOK {
		t.Fatalf("hapus: %v %v", res, err)
	}
	if u.cacah(t, "M_ATTACHMENTPRODUCTNAME", "ID = :1", a.ID) != 0 || u.cacah(t, "T_STORAGE_IMAGE", "IMAGEID = :1", a.StorageID) != 0 {
		t.Error("rekam dan objek terhapus")
	}
}
