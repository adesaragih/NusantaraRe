package handlers_test

// Uji jalur HTTP tiket 15: 503 tanpa basis data, 401 tanpa identitas, 404
// himpunan tak dikenal, 200 berbadan JSON.

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"nusantarare/modul/treatyin/backend/handlers"
	"nusantarare/modul/treatyin/backend/models"
	"nusantarare/modul/treatyin/backend/services"
)

type gudangTiruan struct {
	isi     []models.Acuan
	kontrak models.KontrakDenganVersi
	warisan []models.Kontrak
	adaQS   bool
	galat   error
	// Tiket 40 - nil berarti versinya yang pertama.
	versiDasar *models.VersiKontrak
	// Tiket 42.
	bukti models.BuktiArsip
	// Layar daftar kontrak.
	daftar []models.BarisDaftarKontrak
	// Layar daftar WARISAN (`TREATY_IN`).
	barisWarisan   []models.BarisDaftarWarisan
	cacahWarisan   int
	kontrakWarisan models.KontrakWarisan
}

// Tiket 32.
func (g gudangTiruan) CatatPemulihanLimit(context.Context, int64, []models.PemulihanLimit) error {
	return g.galat
}

// Tiket 42.
func (g gudangTiruan) SimpanArsipMuatanKeluar(context.Context, models.ArsipMuatanKeluar) error {
	return g.galat
}

func (g gudangTiruan) BuktiArsipKontrak(context.Context, int64) (models.BuktiArsip, error) {
	if g.galat != nil {
		return models.BuktiArsip{}, g.galat
	}
	return g.bukti, nil
}

// Tiket 40.
func (g gudangTiruan) BacaVersiDasar(context.Context, int64) (*models.VersiKontrak, error) {
	if g.galat != nil {
		return nil, g.galat
	}
	return g.versiDasar, nil
}

func (g gudangTiruan) DaftarAcuan(context.Context, models.Himpunan) ([]models.Acuan, error) {
	return g.isi, nil
}

func (g gudangTiruan) BuatKontrakDenganVersiPertama(context.Context, models.Kontrak, models.VersiKontrak) (int64, int64, error) {
	if g.galat != nil {
		return 0, 0, g.galat
	}
	return 777, 888, nil
}

func (g gudangTiruan) BacaKontrak(context.Context, int64) (models.KontrakDenganVersi, error) {
	if g.galat != nil {
		return models.KontrakDenganVersi{}, g.galat
	}
	return g.kontrak, nil
}

func minta(t *testing.T, h http.Handler, jalur, pelaku string) *httptest.ResponseRecorder {
	t.Helper()
	r := httptest.NewRequest(http.MethodGet, jalur, nil)
	if pelaku != "" {
		r.Header.Set("X-Pelaku", pelaku)
	}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	return w
}

// Tanpa basis data setiap rute menjawab 503 - dan ia dijawab SEBELUM identitas
// diperiksa, sebab proses tanpa Oracle tidak dapat melayani siapa pun.
func TestTanpaBasisDataMenjawab503(t *testing.T) {
	h := handlers.RouterDengan(services.LayananDengan(gudangTiruan{}), false, true)

	w := minta(t, h, handlers.Prefix+"/acuan/jenis-potongan", "AKUN-UJI")

	if w.Code != http.StatusServiceUnavailable {
		t.Fatalf("mau 503, dapat %d", w.Code)
	}
}

func TestTanpaIdentitasMenjawab401(t *testing.T) {
	h := handlers.RouterDengan(services.LayananDengan(gudangTiruan{}), true, true)

	w := minta(t, h, handlers.Prefix+"/acuan/jenis-potongan", "")

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("mau 401, dapat %d (badan %s)", w.Code, w.Body.String())
	}
}

// Himpunan di luar KELIMA menjawab 404, dan pesannya MENYEBUT yang diminta -
// penolakan yang tidak menyebut apa yang ditolak membuat pemanggilnya menebak.
func TestHimpunanTakDikenalMenjawab404(t *testing.T) {
	h := handlers.RouterDengan(services.LayananDengan(gudangTiruan{}), true, true)

	w := minta(t, h, handlers.Prefix+"/acuan/VERSI_KONTRAK", "AKUN-UJI")

	if w.Code != http.StatusNotFound {
		t.Fatalf("mau 404, dapat %d", w.Code)
	}
	var badan struct {
		Galat string `json:"galat"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &badan); err != nil {
		t.Fatalf("badan bukan JSON: %v", err)
	}
	if badan.Galat == "" {
		t.Fatal("badan 404 tanpa pesan")
	}
}

// Uji POSITIF: KELIMA himpunan menjawab 200 berbadan daftar.
// `mata-uang` dicabut 4 Oktober 2026 — migrasi 434.
func TestKelimaHimpunanMenjawab200(t *testing.T) {
	h := handlers.RouterDengan(
		services.LayananDengan(gudangTiruan{isi: []models.Acuan{{ID: 7, Kode: "IDR", Nama: "Rupiah", Aktif: "1"}}}),
		true, true)

	for _, himpunan := range []string{
		"jenis-potongan", "kelas-bisnis", "kelompok-treaty", "bahaya", "jenis-reasuransi",
	} {
		w := minta(t, h, handlers.Prefix+"/acuan/"+himpunan, "AKUN-UJI")
		if w.Code != http.StatusOK {
			t.Errorf("%s: mau 200, dapat %d (badan %s)", himpunan, w.Code, w.Body.String())
			continue
		}
		var baris []models.Acuan
		if err := json.Unmarshal(w.Body.Bytes(), &baris); err != nil {
			t.Errorf("%s: badan bukan daftar JSON: %v", himpunan, err)
			continue
		}
		if len(baris) != 1 || baris[0].Kode != "IDR" {
			t.Errorf("%s: badan tidak sesuai: %+v", himpunan, baris)
		}
	}
}

// Daftar kosong terkirim sebagai `[]`, bukan `null`: pembaca JavaScript yang
// menerima null akan jatuh pada `.map`, dan jatuhnya di layar - bukan di sini.
func TestDaftarKosongTerkirimSebagaiLarikKosong(t *testing.T) {
	h := handlers.RouterDengan(services.LayananDengan(gudangTiruan{isi: []models.Acuan{}}), true, true)

	w := minta(t, h, handlers.Prefix+"/acuan/bahaya", "AKUN-UJI")

	if got := w.Body.String(); got != "[]\n" && got != "[]" {
		t.Fatalf("mau badan `[]`, dapat %q", got)
	}
}

func (g gudangTiruan) CariKontrakSerupa(context.Context, models.Kontrak) ([]int64, error) {
	return nil, g.galat
}

func (g gudangTiruan) CariKontrakLewatNomorWarisan(context.Context, string) ([]models.Kontrak, error) {
	if g.galat != nil {
		return nil, g.galat
	}
	return g.warisan, nil
}

func (g gudangTiruan) PerbaruiKontrak(context.Context, models.Kontrak) error { return g.galat }

func (g gudangTiruan) TambahVersi(context.Context, int64, models.VersiKontrak) (int64, error) {
	if g.galat != nil {
		return 0, g.galat
	}
	return 999, nil
}

// Tiket 35 dan 36.
func (g gudangTiruan) AdaQuotaSharePadaVersi(context.Context, int64) (bool, error) {
	if g.galat != nil {
		return false, g.galat
	}
	return g.adaQS, nil
}

func (g gudangTiruan) CatatKetentuanProporsional(_ context.Context, _, _, _ int64, jenis, _ string, _ *int64) error {
	if g.galat != nil {
		return g.galat
	}
	_ = jenis
	return nil
}

// Layar daftar kontrak — ronde layar 1.
func (g gudangTiruan) DaftarKontrak(context.Context) ([]models.BarisDaftarKontrak, error) {
	if g.galat != nil {
		return nil, g.galat
	}
	return g.daftar, nil
}

// Layar daftar WARISAN — jalur baca `TREATY_IN`.
func (g gudangTiruan) CacahKontrakWarisan(context.Context) (int, error) {
	if g.galat != nil {
		return 0, g.galat
	}
	return g.cacahWarisan, nil
}

func (g gudangTiruan) DaftarKontrakWarisan(_ context.Context, offset, batas int) ([]models.BarisDaftarWarisan, error) {
	if g.galat != nil {
		return nil, g.galat
	}
	// Memotong seperti basis data memotong, supaya uji penomoran halaman
	// menguji penomorannya - bukan tiruan yang selalu mengembalikan semua.
	if offset >= len(g.barisWarisan) {
		return []models.BarisDaftarWarisan{}, nil
	}
	akhir := offset + batas
	if akhir > len(g.barisWarisan) {
		akhir = len(g.barisWarisan)
	}
	return append([]models.BarisDaftarWarisan{}, g.barisWarisan[offset:akhir]...), nil
}

// Satu kontrak WARISAN.
func (g gudangTiruan) BacaKontrakWarisan(_ context.Context, id string) (models.KontrakWarisan, error) {
	if g.galat != nil {
		return models.KontrakWarisan{}, g.galat
	}
	return g.kontrakWarisan, nil
}

// Tab yang PINDAH ke tabel pendaratan migrasi 430.
func (g gudangTiruan) BacaPeriodePelaporan(_ context.Context, _ string) ([]models.BarisPeriodeWarisan, error) {
	return nil, nil
}

func (g gudangTiruan) BacaPortofolio(_ context.Context, _ string) ([]models.BarisPortofolioWarisan, error) {
	return nil, nil
}

func (g gudangTiruan) BacaAkumulasi(_ context.Context, _ string) ([]models.BarisAkumulasiWarisan, error) {
	return nil, nil
}

// Empat tab pendaratan berikutnya — tiruan datar; perilakunya diuji di
// `services`, bukan di sini.
func (g gudangTiruan) BacaEgnpi(_ context.Context, _ string) ([]models.BarisEgnpiWarisan, error) {
	return nil, nil
}
func (g gudangTiruan) BacaRetensi(_ context.Context, _ string) ([]models.BarisRetensiWarisan, error) {
	return nil, nil
}
func (g gudangTiruan) BacaAngsuran(_ context.Context, _ string) ([]models.BarisAngsuranWarisan, error) {
	return nil, nil
}
func (g gudangTiruan) BacaCatatan(_ context.Context, _ string) ([]models.BarisCatatanWarisan, error) {
	return nil, nil
}

// Empat tab dari `M_TREATY_IN2` - satu seam untuk keempatnya.
func (g gudangTiruan) BacaLayerWarisan(_ context.Context, _ string) ([]models.BarisLayerWarisan, error) {
	return nil, nil
}

func (g gudangTiruan) BacaSkalaKoasuransi(_ context.Context, _ string) ([]models.BarisSkalaKoasuransiWarisan, error) {
	return nil, nil
}
