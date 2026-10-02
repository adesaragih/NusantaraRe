//go:build db

package handlers_test

// Seam HTTP rute baca terhadap skema uji Oracle NYATA (paket 1).
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

func TestDBBacaProdukDariJSONDATAKeduaTabel(t *testing.T) {
	u := pasangDB(t)
	u.exec(t, `INSERT INTO {s}.M_PRODUCT_LIFE (ID, JSONDATA) VALUES ('100002', :1)`,
		`{"ID":"100002","CEDING":"UJI CEDING B","TREATYNUMBER":"UJI/2"}`)
	u.exec(t, `INSERT INTO {s}.M_PRODUCT_LIFE (ID, JSONDATA) VALUES ('100001', :1)`,
		`{"ID":"100001","CEDING":"UJI CEDING A","UnderwritingLimitList":[{"MaxInsured":1175000000.123456789}]}`)
	// Inward lama ber-ID sequence sendiri, bertaut lewat PRODUCTID (Pega).
	u.exec(t, `INSERT INTO {s}.M_PRODUCTINWARD_LIFE (ID, JSONDATA) VALUES ('100009', :1)`,
		`{"PRODUCTID":"100001","BEGIN":"01/03/2026","MAXEXPIREDCLAIM":"180"}`)

	kode, badan := u.kirim(t, "GET", pre+"/produk", "")
	if kode != http.StatusOK || strings.Index(badan, "100001") > strings.Index(badan, "100002") {
		t.Fatalf("grid urut ID: %d %s", kode, badan)
	}
	kode, badan = u.kirim(t, "GET", pre+"/produk/100001", "")
	if kode != http.StatusOK {
		t.Fatalf("produk: %d %s", kode, badan)
	}
	for _, w := range []string{`"maxInsured":"1175000000.123456789"`, `"id":"100009"`, `"begin":"2026-03-01"`,
		`"maxExpiredClaim":"180"`} {
		if !strings.Contains(badan, w) {
			t.Errorf("tanpa %s: %s", w, badan)
		}
	}
	if kode, _ := u.kirim(t, "GET", pre+"/produk/100777", ""); kode != http.StatusNotFound {
		t.Errorf("produk tidak ada: %d", kode)
	}
}
