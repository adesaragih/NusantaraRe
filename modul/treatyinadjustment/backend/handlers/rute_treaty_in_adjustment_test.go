package handlers_test

// Uji jalur HTTP tiket 01 dan 05.

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"nusantarare/modul/treatyinadjustment/backend/handlers"
	"nusantarare/modul/treatyinadjustment/backend/models"
	"nusantarare/modul/treatyinadjustment/backend/services"
)

type gudangTiruan struct {
	kontrak []models.Kontrak
	versi   map[int64][]models.Versi

	lampiran        []models.BarisLampiranWarisan
	katalogKategori map[string]string
	riwayat         []models.BarisRiwayatWarisan
	galatLampiran   error

	daftarPenyesuaian []models.BarisPenyesuaian
	penyesuaian       map[string]models.Penyesuaian
	galatPenyesuaian  error
}

func (g gudangTiruan) DaftarKontrak(context.Context) ([]models.Kontrak, error) {
	return g.kontrak, nil
}

func (g gudangTiruan) RantaiVersi(_ context.Context, id int64) ([]models.Versi, error) {
	return g.versi[id], nil
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

func nomor(n int64) *int64 { return &n }

func routerIsi() http.Handler {
	return handlers.RouterDengan(services.LayananDengan(gudangTiruan{
		kontrak: []models.Kontrak{{ID: 7, NomorKontrakWarisan: "TRI-7", SifatProporsi: "PROPORSIONAL",
			TanggalMulai: "2026-01-01", TanggalBerakhir: "2026-12-31"}},
		versi: map[int64][]models.Versi{7: {
			{ID: 1, IDKontrak: 7, NomorUrutVersi: nomor(1), KeadaanSiklusHidup: "DRAFT"},
			{ID: 2, IDKontrak: 7, NomorUrutVersi: nil, KeadaanSiklusHidup: "DRAFT", IDVersiDasar: nomor(1)},
		}},
	}), true, true)
}

func TestTanpaBasisDataMenjawab503(t *testing.T) {
	h := handlers.RouterDengan(services.LayananDengan(gudangTiruan{}), false, true)

	for _, jalur := range []string{handlers.Prefix + "/kontrak", handlers.Prefix + "/kontrak/7/versi"} {
		if w := minta(t, h, jalur, "AKUN-UJI"); w.Code != http.StatusServiceUnavailable {
			t.Errorf("%s: mau 503, dapat %d", jalur, w.Code)
		}
	}
}

func TestTanpaIdentitasMenjawab401(t *testing.T) {
	h := routerIsi()

	for _, jalur := range []string{handlers.Prefix + "/kontrak", handlers.Prefix + "/kontrak/7/versi"} {
		if w := minta(t, h, jalur, ""); w.Code != http.StatusUnauthorized {
			t.Errorf("%s: mau 401, dapat %d", jalur, w.Code)
		}
	}
}

// Pengenal yang bukan angka adalah BENTUK permintaan yang salah (400), bukan
// kontrak yang tidak ada (404). Dua pertanyaan berbeda, dua jawaban berbeda.
func TestPengenalBukanAngkaMenjawab400(t *testing.T) {
	h := routerIsi()

	w := minta(t, h, handlers.Prefix+"/kontrak/tujuh/versi", "AKUN-UJI")

	if w.Code != http.StatusBadRequest {
		t.Fatalf("mau 400, dapat %d (badan %s)", w.Code, w.Body.String())
	}
}

func TestKontrakTidakAdaMenjawab404(t *testing.T) {
	h := routerIsi()

	w := minta(t, h, handlers.Prefix+"/kontrak/999/versi", "AKUN-UJI")

	if w.Code != http.StatusNotFound {
		t.Fatalf("mau 404, dapat %d (badan %s)", w.Code, w.Body.String())
	}
}

// Uji POSITIF: rantai versi terkirim, dan nomor urut maupun dasar yang KOSONG
// terkirim sebagai `null` - bukan 0, yang akan terbaca sebagai versi nomor nol.
func TestRantaiVersiMenjawab200DenganNullYangBenar(t *testing.T) {
	h := routerIsi()

	w := minta(t, h, handlers.Prefix+"/kontrak/7/versi", "AKUN-UJI")

	if w.Code != http.StatusOK {
		t.Fatalf("mau 200, dapat %d (badan %s)", w.Code, w.Body.String())
	}
	var versi []models.Versi
	if err := json.Unmarshal(w.Body.Bytes(), &versi); err != nil {
		t.Fatalf("badan bukan daftar JSON: %v", err)
	}
	if len(versi) != 2 {
		t.Fatalf("mau 2 versi, dapat %d", len(versi))
	}
	if versi[0].IDVersiDasar != nil {
		t.Error("versi pertama: dasar harus null")
	}
	if versi[1].NomorUrutVersi != nil {
		t.Error("baris warisan: nomor urut harus null, bukan 0")
	}
	var mentah []map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &mentah); err != nil {
		t.Fatal(err)
	}
	if n, ada := mentah[1]["nomorUrutVersi"]; !ada || n != nil {
		t.Errorf("nomorUrutVersi kosong harus terkirim sebagai null, dapat %v", n)
	}
}

func TestDaftarKontrakMenjawab200(t *testing.T) {
	h := routerIsi()

	w := minta(t, h, handlers.Prefix+"/kontrak", "AKUN-UJI")

	if w.Code != http.StatusOK {
		t.Fatalf("mau 200, dapat %d", w.Code)
	}
	var k []models.Kontrak
	if err := json.Unmarshal(w.Body.Bytes(), &k); err != nil {
		t.Fatalf("badan bukan daftar JSON: %v", err)
	}
	if len(k) != 1 || k[0].NomorKontrakWarisan != "TRI-7" {
		t.Errorf("badan tidak sesuai: %+v", k)
	}
}

// Panel Attachment - `M_ATTACHMENTTREATY_2`, tabel warisan.
func (g gudangTiruan) BacaLampiranKontrak(_ context.Context, _ string) ([]models.BarisLampiranWarisan, error) {
	return g.lampiran, g.galatLampiran
}

func (g gudangTiruan) BacaKatalogKategoriLampiran(_ context.Context) (map[string]string, error) {
	return g.katalogKategori, g.galatLampiran
}

// ---------------------------------------------------------------------
// Rute panel Attachment - `GET /kontrak-warisan/{id}/lampiran`.
// ---------------------------------------------------------------------

func routerLampiran() http.Handler {
	return handlers.RouterDengan(services.LayananDengan(gudangTiruan{
		katalogKategori: map[string]string{
			"00000": "Others", "00001": "Analysed Email", "00002": "Approval Email",
			"00005": "Summary Treaty Leader", "00006": "Assessment Inward Treaty Form",
			"00007": "Pega Proportional Calculation", "00010": "Offer Email",
		},
		lampiran: []models.BarisLampiranWarisan{
			{ID: "1", KodeKategori: "00001", NamaBerkas: "a.pdf", JenisMime: "application/pdf"},
		},
	}), true, true)
}

func TestLampiranMenjawab200DenganKategoriLengkap(t *testing.T) {
	w := minta(t, routerLampiran(), "/api/treaty-in-adjustment/kontrak-warisan/1001851/lampiran", "UJI")
	if w.Code != http.StatusOK {
		t.Fatalf("kode %d, mau 200 (badan %s)", w.Code, w.Body.String())
	}
	var isi struct {
		Kategori []struct {
			Kode       string `json:"kode"`
			Nama       string `json:"nama"`
			Cacah      int    `json:"cacah"`
			Dipastikan bool   `json:"dipastikan"`
		} `json:"kategori"`
		Berkas []struct {
			NamaBerkas string `json:"namaBerkas"`
		} `json:"berkas"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &isi); err != nil {
		t.Fatalf("badan bukan JSON yang diharapkan: %v - %s", err, w.Body.String())
	}
	if len(isi.Kategori) != 11 {
		t.Errorf("%d kategori, mau 11 (7 terbukti + 4 belum dipastikan)", len(isi.Kategori))
	}
	if len(isi.Berkas) != 1 {
		t.Errorf("%d berkas, mau 1", len(isi.Berkas))
	}
	// ⛔ Keempat kode tanpa nama sampai ke layar TANPA nama tebakan.
	belum := 0
	for _, k := range isi.Kategori {
		if k.Dipastikan {
			continue
		}
		belum++
		if k.Nama != "" {
			t.Errorf("kode %s sampai ke layar bernama %q; pasangannya belum diketahui", k.Kode, k.Nama)
		}
	}
	if belum != 4 {
		t.Errorf("%d kategori belum dipastikan sampai ke layar, mau 4", belum)
	}
}

// Tanpa identitas -> 401, dan rute barunya ikut dijaga pagar yang sama.
func TestLampiranTanpaIdentitasMenjawab401(t *testing.T) {
	h := handlers.RouterDengan(services.LayananDengan(gudangTiruan{}), true, false)
	w := minta(t, h, "/api/treaty-in-adjustment/kontrak-warisan/1001851/lampiran", "")
	if w.Code != http.StatusUnauthorized {
		t.Errorf("kode %d, mau 401", w.Code)
	}
}

// Tanpa basis data -> 503, sama seperti rute lain.
func TestLampiranTanpaBasisDataMenjawab503(t *testing.T) {
	h := handlers.RouterDengan(services.LayananDengan(gudangTiruan{}), false, true)
	w := minta(t, h, "/api/treaty-in-adjustment/kontrak-warisan/1001851/lampiran", "UJI")
	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("kode %d, mau 503", w.Code)
	}
}

// Panel History - `T_VIEW_COMMENT`, tabel warisan.
func (g gudangTiruan) BacaRiwayatKontrak(_ context.Context, _ string) ([]models.BarisRiwayatWarisan, error) {
	return g.riwayat, g.galatLampiran
}
