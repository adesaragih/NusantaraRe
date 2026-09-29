//go:build db

// Seam HTTP terhadap skema uji Oracle (brief bab 5).
//
// Sistem digerakkan lewat HTTP dan hasilnya diperiksa lewat HTTP - bukan lewat
// mock repository, dan bukan lewat SQL sampingan di dalam test ini.
//
// Jalankan: make test-db   (perlu ORACLE_DSN dan ORACLE_SCHEMA)
package handlers_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"nusantarare/inti/config"
	"nusantarare/inti/db"
	"nusantarare/modul/claimlife/handlers"
	"nusantarare/modul/claimlife/services"
	"nusantarare/uji/skemauji"
)

func server(t *testing.T) (*httptest.Server, func()) {
	t.Helper()
	sqlDB, skema, err := skemauji.Buka()
	if err != nil {
		// Hanya "Oracle belum dikonfigurasi" yang dilewati. Salah
		// konfigurasi, menunjuk produksi, atau menunjuk skema yang bukan
		// skema uji harus MENGGAGALKAN - test ini menghapus tabel.
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
	if err := skemauji.IsiContoh(ctx, sqlDB, skema); err != nil {
		t.Fatalf("mengisi fixture: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		t.Fatal(err)
	}
	db, err := db.Open(cfg)
	if err != nil {
		t.Fatal(err)
	}
	// Refactor bentuk B: pembaca polis PremiumList disambung seperti di cmd/api.
	srv := httptest.NewServer(handlers.Router(services.New(db).DenganPembacaPolis(skemauji.PembacaPolis(db)), false))
	return srv, func() {
		srv.Close()
		_ = skemauji.Bongkar(ctx, sqlDB, skema)
		_ = db.Close()
		_ = sqlDB.Close()
	}
}

func ambil(t *testing.T, srv *httptest.Server, jalur string) (int, string) {
	t.Helper()
	res, err := http.Get(srv.URL + jalur)
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

// Tiket 01 AC-1: GET satu klaim mengembalikan klaim beserta SELURUH baris
// AdjustmentList-nya, bukan hanya baris terakhir.
func TestGetKlaimMengembalikanSeluruhBaris(t *testing.T) {
	srv, bersihkan := server(t)
	defer bersihkan()

	kode, badan := ambil(t, srv, "/api/klaim-life/UJI-KLAIM-1")
	if kode != http.StatusOK {
		t.Fatalf("status = %d, badan = %s", kode, badan)
	}

	var jawab struct {
		ID         string `json:"id"`
		NomorKlaim string `json:"nomorKlaim"`
		CacahBaris int    `json:"cacahBaris"`
		Peserta    []struct {
			ID              string `json:"id"`
			NomorSertifikat string `json:"nomorSertifikat"`
			Baris           []struct {
				ID              string `json:"id"`
				Status          string `json:"status"`
				KodeStatus      string `json:"kodeStatus"`
				StatusDiketahui bool   `json:"statusDiketahui"`
				JumlahKlaim     struct {
					Amount   string `json:"amount"`
					Currency string `json:"currency"`
				} `json:"jumlahKlaim"`
			} `json:"baris"`
		} `json:"peserta"`
	}
	if err := json.Unmarshal([]byte(badan), &jawab); err != nil {
		t.Fatalf("badan bukan JSON yang diharapkan: %v\n%s", err, badan)
	}

	if jawab.CacahBaris != 5 {
		t.Errorf("cacahBaris = %d, mau 5", jawab.CacahBaris)
	}
	if len(jawab.Peserta) != 2 {
		t.Fatalf("cacah peserta = %d, mau 2", len(jawab.Peserta))
	}
	if len(jawab.Peserta[0].Baris) != 3 {
		t.Errorf("peserta 1 punya %d baris, mau 3 - hanya baris terakhir yang lolos?",
			len(jawab.Peserta[0].Baris))
	}
	if jawab.Peserta[0].NomorSertifikat != "006" {
		t.Errorf("nomor sertifikat = %q lewat HTTP, awalan nol hilang",
			jawab.Peserta[0].NomorSertifikat)
	}
}

// Tiket 01 AC-3: status sebagai kata, bukan nama field dan bukan angka.
func TestStatusKeluarSebagaiKata(t *testing.T) {
	srv, bersihkan := server(t)
	defer bersihkan()

	_, badan := ambil(t, srv, "/api/klaim-life/UJI-KLAIM-1")
	for _, kata := range []string{`"status":"Outstanding"`, `"status":"Aksep"`, `"status":"Ditolak"`} {
		if !strings.Contains(badan, kata) {
			t.Errorf("tidak ada %s di badan jawaban", kata)
		}
	}
	if strings.Contains(badan, "STS_REJECT") {
		t.Error("nama field warisan STS_REJECT bocor ke kontrak API")
	}
	// Kode di luar spec tidak ditebak.
	if !strings.Contains(badan, `"status":"Tidak diketahui"`) {
		t.Error("kode 9 seharusnya dilaporkan sebagai tidak diketahui")
	}
}

// Tiket 01 AC-4: nol uang sebagai binary floating point di kontrak API.
func TestUangDiKontrakAPIAdalahTeks(t *testing.T) {
	srv, bersihkan := server(t)
	defer bersihkan()

	_, badan := ambil(t, srv, "/api/klaim-life/UJI-KLAIM-1")
	if !strings.Contains(badan, `"amount":"1234567890.12345678"`) {
		t.Errorf("jumlah besar berubah atau bukan teks:\n%s", badan)
	}
	if !strings.Contains(badan, `"amount":"0.00000001"`) {
		t.Errorf("jumlah kecil berubah atau bukan teks:\n%s", badan)
	}
	if strings.Contains(badan, `"amount":1`) || strings.Contains(badan, `"amount":0.`) {
		t.Errorf("ada jumlah uang yang keluar sebagai angka JSON:\n%s", badan)
	}
}

func TestKlaimTidakAdaMenjawab404(t *testing.T) {
	srv, bersihkan := server(t)
	defer bersihkan()

	kode, _ := ambil(t, srv, "/api/klaim-life/TIDAK-ADA")
	if kode != http.StatusNotFound {
		t.Errorf("status = %d, mau 404", kode)
	}
}
