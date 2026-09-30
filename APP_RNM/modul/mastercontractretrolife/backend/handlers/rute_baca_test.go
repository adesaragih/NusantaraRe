package handlers_test

// Seam HTTP rute baca (paket 1) di atas gudang tiruan - tanpa Oracle.

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/modul/mastercontractretrolife/backend/handlers"
	"nusantarare/modul/mastercontractretrolife/backend/models"
	"nusantarare/modul/mastercontractretrolife/backend/services"
	"nusantarare/modul/mastercontractretrolife/backend/tiruan"
)

func desimal(t *testing.T, s string) *apd.Decimal {
	t.Helper()
	d, _, err := apd.NewFromString(s)
	if err != nil {
		t.Fatal(err)
	}
	return d
}

func gudangHTTP(t *testing.T) *tiruan.Gudang {
	g := tiruan.Baru()
	g.Tahun["UJI-T1"] = models.TahunTreaty{ID: "UJI-T1", TreatyYear: "2026", UnderwritingYear: "2026",
		StartDate: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC), EndDate: time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC)}
	g.Kontrak["UJI-K1"] = models.Kontrak{ID: "UJI-K1", IDTreatyYear: "UJI-T1", ReinsTypeID: "10196", ReinsTypeName: "QS",
		BIDR: desimal(t, "0"), IDR: desimal(t, "1500000000.123456789"), BUSD: desimal(t, "0")}
	g.Jenis = []models.JenisReasuransi{{ID: "10200", Note: "OR"}, {ID: "10196", Note: "QS"}}
	return g
}

type uji struct {
	srv *httptest.Server
	g   *tiruan.Gudang
}

func server(t *testing.T, adaDB bool) *uji {
	t.Helper()
	g := gudangHTTP(t)
	srv := httptest.NewServer(handlers.RouterDengan(services.BaruLayanan(g, g.Transaksi, nil), adaDB, true))
	t.Cleanup(srv.Close)
	return &uji{srv: srv, g: g}
}

func (u *uji) minta(t *testing.T, metode, jalur string, pelaku bool) (int, string) {
	t.Helper()
	req, err := http.NewRequest(metode, u.srv.URL+jalur, nil)
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

func TestTanpaOracleSetiapRute503BerGalat(t *testing.T) {
	u := server(t, false)
	for _, j := range []string{"/tahun", "/tahun/UJI-T1/kontrak", "/kontrak/UJI-K1/reinsurer", "/reinsurer/X/security",
		"/kontrak/UJI-K1/business", "/jenis-reasuransi", "/master-reinsurer", "/master-business", "/ringkasan-rate",
		"/rate?idusedby=1"} {
		kode, badan := u.minta(t, "GET", handlers.Prefix+j, true)
		if kode != http.StatusServiceUnavailable || !strings.Contains(badan, `"galat"`) {
			t.Errorf("%s tanpa Oracle: %d %s", j, kode, badan)
		}
	}
}

func TestTanpaIdentitas401(t *testing.T) {
	kode, _ := server(t, true).minta(t, "GET", handlers.Prefix+"/tahun", false)
	if kode != http.StatusUnauthorized {
		t.Errorf("tanpa X-Pelaku: %d, mau 401", kode)
	}
}

func TestGridTahunDanKontrak(t *testing.T) {
	u := server(t, true)
	kode, badan := u.minta(t, "GET", handlers.Prefix+"/tahun", true)
	if kode != http.StatusOK {
		t.Fatalf("GET tahun: %d %s", kode, badan)
	}
	var d struct {
		Daftar []map[string]string `json:"daftar"`
		Total  int                 `json:"total"`
	}
	if err := json.Unmarshal([]byte(badan), &d); err != nil {
		t.Fatal(err)
	}
	if d.Total != 1 || d.Daftar[0]["startDate"] != "2026-01-01" || d.Daftar[0]["treatyYear"] != "2026" {
		t.Errorf("tahun: %s", badan)
	}
	kode, badan = u.minta(t, "GET", handlers.Prefix+"/tahun/UJI-T1/kontrak", true)
	if kode != http.StatusOK {
		t.Fatalf("GET kontrak: %d %s", kode, badan)
	}
	// ⛔ ADR-0003: uang sebagai TEKS, persis - tanpa pembulatan.
	if !strings.Contains(badan, `"idr":"1500000000.123456789"`) || !strings.Contains(badan, `"usd":""`) {
		t.Errorf("uang kontrak bukan teks persis: %s", badan)
	}
}

func TestIndukTidakAda404Berkalimat(t *testing.T) {
	kode, badan := server(t, true).minta(t, "GET", handlers.Prefix+"/tahun/UJI-TAK-ADA/kontrak", true)
	if kode != http.StatusNotFound || !strings.Contains(badan, "treaty year not found") {
		t.Errorf("tahun tak ada: %d %s", kode, badan)
	}
	if strings.Contains(badan, "services:") {
		t.Errorf("awalan lapisan bocor ke layar: %s", badan)
	}
}

func TestDaftarKosongAdalahLarikBukanNull(t *testing.T) {
	kode, badan := server(t, true).minta(t, "GET", handlers.Prefix+"/kontrak/UJI-K1/business", true)
	if kode != http.StatusOK || !strings.Contains(badan, `"daftar":[]`) {
		t.Errorf("business kosong: %d %s", kode, badan)
	}
}

func TestMasterTidakTerbaca503MenyebutObjek(t *testing.T) {
	u := server(t, true)
	u.g.Jenis = nil
	kode, badan := u.minta(t, "GET", handlers.Prefix+"/jenis-reasuransi", true)
	if kode != http.StatusServiceUnavailable || !strings.Contains(badan, "REINSURANCETYPE") {
		t.Errorf("master jenis kosong: %d %s", kode, badan)
	}
}

func TestJenisReasuransiUrutRD(t *testing.T) {
	kode, badan := server(t, true).minta(t, "GET", handlers.Prefix+"/jenis-reasuransi", true)
	if kode != http.StatusOK || strings.Index(badan, "10200") > strings.Index(badan, "10196") {
		t.Errorf("jenis (urut ID DESC dari gudang): %d %s", kode, badan)
	}
}
