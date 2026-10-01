package handlers_test

// Seam HTTP rute baca (paket 1) di atas gudang tiruan - tanpa Oracle. Tiruan
// menyimpan JSONDATA mentah dan memakai kodek repository sungguhan.

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"nusantarare/modul/masterproductnamelife/backend/handlers"
	"nusantarare/modul/masterproductnamelife/backend/services"
	"nusantarare/modul/masterproductnamelife/backend/tiruan"
)

func gudangHTTP() *tiruan.Gudang {
	g := tiruan.Baru()
	g.Umum["UJI-002"] = `{"ID":"UJI-002","CEDING":"UJI CEDING B","TREATYNUMBER":"UJI/2","INWARDNAME":"UJI B",
		"CREATEOP":"UJI-OP","UPDATEOP":"UJI-OP2","POLICYHODERNAME":"UJI PEMEGANG"}`
	g.Umum["UJI-001"] = `{"ID":"UJI-001","CEDING":"UJI CEDING A","PlanList":[{"Plan":"UJI PLAN","RIRATE":"UJI RATE"}]}`
	g.Inward["UJI-001"] = `{"ID":"UJI-001","PRODUCTID":"UJI-001","BEGIN":"01/03/2026","CEDINGLIMIT":"150000000.5"}`
	return g
}

type uji struct {
	srv *httptest.Server
	g   *tiruan.Gudang
}

func server(t *testing.T, adaDB bool) *uji {
	t.Helper()
	g := gudangHTTP()
	srv := httptest.NewServer(handlers.RouterDengan(services.BaruLayanan(g, g.Transaksi, nil), adaDB, true))
	t.Cleanup(srv.Close)
	return &uji{srv: srv, g: g}
}

func (u *uji) minta(t *testing.T, metode, jalur, badan string, pelaku bool) (int, string) {
	t.Helper()
	var isi io.Reader
	if badan != "" {
		isi = bytes.NewBufferString(badan)
	}
	req, err := http.NewRequest(metode, u.srv.URL+jalur, isi)
	if err != nil {
		t.Fatal(err)
	}
	if pelaku {
		req.Header.Set("X-Pelaku", "UJI-PELAKU")
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = res.Body.Close() }()
	b, _ := io.ReadAll(res.Body)
	return res.StatusCode, string(b)
}

const pre = "/api/master-product-name-life"

func TestHTTPGridProdukUrutIDDenganKolomGrid(t *testing.T) {
	u := server(t, true)
	kode, badan := u.minta(t, "GET", pre+"/produk", "", true)
	if kode != http.StatusOK {
		t.Fatalf("%d %s", kode, badan)
	}
	var j struct {
		Daftar []map[string]string `json:"daftar"`
		Total  int                 `json:"total"`
	}
	if err := json.Unmarshal([]byte(badan), &j); err != nil {
		t.Fatal(err)
	}
	if j.Total != 2 || j.Daftar[0]["id"] != "UJI-001" || j.Daftar[1]["id"] != "UJI-002" {
		t.Errorf("urut ID ASC (BrowseProduct_Life b1094): %s", badan)
	}
	if j.Daftar[1]["ceding"] != "UJI CEDING B" || j.Daftar[1]["treatyNumber"] != "UJI/2" ||
		j.Daftar[1]["inwardName"] != "UJI B" || j.Daftar[1]["updateOp"] != "UJI-OP2" {
		t.Errorf("kolom grid: %s", badan)
	}
}

func TestHTTPProdukUtuhKeduaSisi(t *testing.T) {
	u := server(t, true)
	kode, badan := u.minta(t, "GET", pre+"/produk/UJI-001", "", true)
	if kode != http.StatusOK {
		t.Fatalf("%d %s", kode, badan)
	}
	for _, w := range []string{`"id":"UJI-001"`, `"ceding":"UJI CEDING A"`, `"begin":"2026-03-01"`,
		`"cedingLimit":"150000000.5"`, `"plan":"UJI PLAN"`, `"productId":"UJI-001"`} {
		if !strings.Contains(badan, w) {
			t.Errorf("jawaban tanpa %s: %s", w, badan)
		}
	}
}

func TestHTTPGalatBerkalimat(t *testing.T) {
	u := server(t, true)
	if kode, badan := u.minta(t, "GET", pre+"/produk/UJI-TIDAK-ADA", "", true); kode != http.StatusNotFound ||
		!strings.Contains(badan, "product not found") {
		t.Errorf("404: %d %s", kode, badan)
	}
	if kode, _ := u.minta(t, "GET", pre+"/produk", "", false); kode != http.StatusUnauthorized {
		t.Errorf("tanpa identitas harus 401, dapat %d", kode)
	}
	u.g.Umum["UJI-RUSAK"] = `{"PRODUCTNAME":`
	if kode, badan := u.minta(t, "GET", pre+"/produk/UJI-RUSAK", "", true); kode != http.StatusInternalServerError ||
		!strings.Contains(badan, "JSONDATA cannot be read") {
		t.Errorf("JSON rusak harus 500 berkalimat: %d %s", kode, badan)
	}
	mati := server(t, false)
	if kode, badan := mati.minta(t, "GET", pre+"/produk", "", true); kode != http.StatusServiceUnavailable ||
		!strings.Contains(badan, "database is not configured") {
		t.Errorf("tanpa Oracle harus 503: %d %s", kode, badan)
	}
}
