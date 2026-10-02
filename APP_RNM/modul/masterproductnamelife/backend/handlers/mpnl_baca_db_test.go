//go:build db

package handlers_test

// Seam HTTP rute baca terhadap skema uji Oracle NYATA (paket 1) - sejak 02-10-2026 atas tabel FLAT.
// Tanpa ORACLE_DSN seluruhnya MELEWATI dengan pesan.

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"testing"
)

func (u *ujiDB) kirim(t *testing.T, metode, jalur, badan string) (int, string) {
	t.Helper()
	var isi io.Reader
	if badan != "" {
		isi = bytes.NewBufferString(badan)
	}
	req, err := http.NewRequest(metode, u.srv.URL+jalur, isi)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("X-Pelaku", "UJI-PELAKU")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = res.Body.Close() }()
	b, _ := io.ReadAll(res.Body)
	return res.StatusCode, string(b)
}

func TestDBBacaProdukDariTabelFlat(t *testing.T) {
	u := pasangDB(t)
	u.exec(t, `INSERT INTO {s}.M_PRODUCTNAME_LIFE (ID, CEDING, TREATYNUMBER) VALUES ('100002', 'UJI CEDING B', 'UJI/2')`)
	u.exec(t, `INSERT INTO {s}.M_PRODUCTNAME_LIFE (ID, CEDING, BEGIN_DATE, MAXEXPIREDCLAIM, RICOMM, IS_ORS)
		VALUES ('100001', 'UJI CEDING A', DATE '2026-03-01', 180, 0.5, 1)`)
	u.exec(t, `INSERT INTO {s}.M_PRODUCTNAME_LIFE_UWLIMIT (PRODUCTID, URUT, MAXINSURED, MINAGE)
		VALUES ('100001', 1, 1175000000.12345678, 18)`)
	u.exec(t, `INSERT INTO {s}.M_PRODUCTNAME_LIFE_COMMENT (PRODUCTID, URUT, TANGGAL, OPERATORNAME)
		VALUES ('100001', 1, TIMESTAMP '2024-12-02 06:54:50.847', 'UJI-A')`)

	kode, badan := u.kirim(t, "GET", pre+"/produk", "")
	if kode != http.StatusOK || strings.Index(badan, "100001") > strings.Index(badan, "100002") ||
		!strings.Contains(badan, "UJI CEDING B") {
		t.Fatalf("grid urut ID dari induk: %d %s", kode, badan)
	}
	kode, badan = u.kirim(t, "GET", pre+"/produk/100001", "")
	if kode != http.StatusOK {
		t.Fatalf("produk: %d %s", kode, badan)
	}
	// Bentuk kanonik dari Oracle: desimal di bawah satu berawalan 0 (TM9 menulis `.5`), tanggal YYYY-MM-DD, ID inward =
	// ID produk, stempel komentar Pega GMT.
	for _, w := range []string{`"maxInsured":"1175000000.12345678"`, `"minAge":"18"`, `"id":"100001"`, `"productId":"100001"`,
		`"begin":"2026-03-01"`, `"maxExpiredClaim":"180"`, `"riComm":"0.5"`, `"isOrs":true`, `"date":"20241202T065450.847 GMT"`} {
		if !strings.Contains(badan, w) {
			t.Errorf("tanpa %s: %s", w, badan)
		}
	}
	if kode, _ := u.kirim(t, "GET", pre+"/produk/100777", ""); kode != http.StatusNotFound {
		t.Errorf("produk tidak ada: %d", kode)
	}
}
