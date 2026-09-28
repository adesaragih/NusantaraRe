//go:build db

// Seam HTTP Treaty Contract Out terhadap skema uji Oracle NYATA.
//
// Jalankan: make test-db (perlu ORACLE_DSN, ORACLE_SCHEMA, ORACLE_SKEMA_UJI).
// Tanpa Oracle seluruhnya MELEWATI dengan pesan.
//
// Tiket 02: saringan jenis reasuransi hanya terbukti benar terhadap master
// yang benar-benar memuat ID yang dikecualikan.
package handlers_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"nusantarare/internal/config"
	"nusantarare/internal/handlers"
	"nusantarare/internal/repository"
	"nusantarare/internal/repository/skemauji"
	"nusantarare/internal/services"
)

// serverTCO memasang skema uji dan Router BER-STUB identitas: rute modul ini
// bergerbang identitas (401 tanpa X-Pelaku), jadi header harus terbaca.
func serverTCO(t *testing.T) (*httptest.Server, func(), func(string, []skemauji.JenisReasuransiUji)) {
	t.Helper()
	sqlDB, skema, err := skemauji.Buka()
	if err != nil {
		if !skemauji.BolehDilewati(err) {
			t.Fatalf("skema uji menolak: %v", err)
		}
		t.Skipf("lewati: %v", err)
	}
	ctx := context.Background()
	if err := sqlDB.PingContext(ctx); err != nil {
		t.Skipf("lewati: oracle tidak terjangkau: %v", err)
	}
	if err := skemauji.Pasang(ctx, sqlDB, skema); err != nil {
		t.Fatalf("memasang skema uji: %v", err)
	}
	cfg, err := config.Load()
	if err != nil {
		t.Fatal(err)
	}
	db, err := repository.Open(cfg)
	if err != nil {
		t.Fatal(err)
	}
	srv := httptest.NewServer(handlers.Router(services.New(db), true))
	isi := func(_ string, baris []skemauji.JenisReasuransiUji) {
		if err := skemauji.IsiJenisReasuransiTCO(ctx, sqlDB, skema, baris); err != nil {
			t.Fatalf("mengisi master jenis reasuransi: %v", err)
		}
	}
	return srv, func() {
		srv.Close()
		_ = skemauji.Bongkar(ctx, sqlDB, skema)
		_ = db.Close()
		_ = sqlDB.Close()
	}, isi
}

func ambilTCO(t *testing.T, srv *httptest.Server, jalur string, beridentitas bool) (int, string) {
	t.Helper()
	r, err := http.NewRequest(http.MethodGet, srv.URL+jalur, nil)
	if err != nil {
		t.Fatal(err)
	}
	if beridentitas {
		r.Header.Set("X-Pelaku", "UJI-ADMIN")
	}
	res, err := http.DefaultClient.Do(r)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = res.Body.Close() }()
	b, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	return res.StatusCode, string(b)
}

// Tiket 02 AC 11: dua belas awalan blacklist, Flag active, Type 1/2/3 -
// ditegakkan terhadap master yang memuat yang harus tersingkir.
func TestJenisReasuransiDisaringPersisRD(t *testing.T) {
	srv, bersihkan, isi := serverTCO(t)
	defer bersihkan()
	isi("", []skemauji.JenisReasuransiUji{
		{ID: "10003", Note: "UJI QUOTA SHARE", Tipe: "1", Flag: "active"},
		{ID: "10005", Note: "UJI SURPLUS", Tipe: "2", Flag: "active"},
		{ID: "10006", Note: "UJI XOL", Tipe: "3", Flag: "active"},
		{ID: "10004", Note: "UJI BLACKLIST", Tipe: "1", Flag: "active"},    // blacklist
		{ID: "100041", Note: "UJI AWALAN", Tipe: "1", Flag: "active"},      // berawalan blacklist
		{ID: "10217", Note: "UJI BLACKLIST 12", Tipe: "2", Flag: "active"}, // blacklist terakhir
		{ID: "10007", Note: "UJI TIPE 4", Tipe: "4", Flag: "active"},       // tipe di luar 1-3
		{ID: "10008", Note: "UJI NONAKTIF", Tipe: "1", Flag: "inactive"},   // flag
		{ID: "10009", Note: "UJI FLAG LIFE", Tipe: "1", Flag: "1"},         // ejaan Life
	})

	kode, badan := ambilTCO(t, srv, "/api/treaty-contract-out/jenis-reasuransi", true)
	if kode != http.StatusOK {
		t.Fatalf("status = %d, badan = %s", kode, badan)
	}
	var jawab struct {
		Daftar []struct{ ID, Note, Tipe string } `json:"daftar"`
		Total  int                               `json:"total"`
	}
	if err := json.Unmarshal([]byte(badan), &jawab); err != nil {
		t.Fatalf("badan bukan JSON yang diharapkan: %v\n%s", err, badan)
	}
	if jawab.Total != 3 || len(jawab.Daftar) != 3 {
		t.Fatalf("total %d / %d baris, mau 3:\n%s", jawab.Total, len(jawab.Daftar), badan)
	}
	// Urutan .Note ASC (b667).
	mau := []string{"10003", "10005", "10006"}
	for i, id := range mau {
		if jawab.Daftar[i].ID != id {
			t.Errorf("baris %d = %s, mau %s (urutan NOTE ASC)", i, jawab.Daftar[i].ID, id)
		}
	}
	// Dan tabel kebenaran Go sepakat dengan SQL, baris demi baris.
	for _, r := range jawab.Daftar {
		if !repository.LolosSaringanNonLifeTCO(r.ID, "active", r.Tipe) {
			t.Errorf("SQL meloloskan %s tetapi tabel kebenaran Go menolaknya", r.ID)
		}
	}
}

func TestJenisReasuransiMasterKosongMenjawab503(t *testing.T) {
	srv, bersihkan, _ := serverTCO(t)
	defer bersihkan()
	kode, badan := ambilTCO(t, srv, "/api/treaty-contract-out/jenis-reasuransi", true)
	if kode != http.StatusServiceUnavailable {
		t.Fatalf("master kosong: status = %d, mau 503 (ADR-0015); badan = %s", kode, badan)
	}
}

func TestJenisReasuransiTanpaIdentitasMenjawab401(t *testing.T) {
	srv, bersihkan, _ := serverTCO(t)
	defer bersihkan()
	if kode, _ := ambilTCO(t, srv, "/api/treaty-contract-out/jenis-reasuransi", false); kode != http.StatusUnauthorized {
		t.Fatalf("tanpa identitas: status = %d, mau 401", kode)
	}
}
