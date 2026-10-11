package handlers_test

// Uji jalur HTTP tiket 01 dan 05.

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
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

	// Picker Add dan tombol `Choose`.
	master        []models.BarisMasterPilihan
	dokumenMaster map[string]models.SisiPenyesuaian
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
		// Kesebelas pasangan `M_KATEGORIMASTERTREATY` (8 Oktober 2026).
		katalogKategori: map[string]string{
			"00000": "Others", "00001": "Analysed Email", "00002": "Approval Email",
			"00003": "Binding, signed share Email", "00004": "Info Pack",
			"00005": "Summary Treaty Leader", "00006": "Assessment Inward Treaty Form / Format Analisa Treaty",
			"00007": "Pega Proportional Calculation /Perhitungan Pega Proportional",
			"00008": "Letter of Acknowledgment / LOA", "00009": "Claim Data", "00010": "Offer Email",
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
		t.Errorf("%d kategori, mau 11", len(isi.Kategori))
	}
	if len(isi.Berkas) != 1 {
		t.Errorf("%d berkas, mau 1", len(isi.Berkas))
	}
	// ⭐ Kesebelasnya BERNAMA, urut nama seperti Pega (`order by note`).
	if k := isi.Kategori[0]; k.Nama != "Analysed Email" || k.Cacah != 1 {
		t.Errorf("baris pertama %+v, mau Analysed Email (1)", k)
	}
	for _, k := range isi.Kategori {
		if !k.Dipastikan || k.Nama == "" {
			t.Errorf("kode %s sampai ke layar tanpa nama", k.Kode)
		}
	}
}

// `?jenis=NonProportional` — kode 00007 bernama Non-Prop
// (`GetMasterTreatyCategory_Act` [2.2]).
func TestLampiranNonPropMemakaiNamaNonProp(t *testing.T) {
	w := minta(t, routerLampiran(), "/api/treaty-in-adjustment/kontrak-warisan/1001851/lampiran?jenis=NonProportional", "UJI")
	if w.Code != http.StatusOK {
		t.Fatalf("kode %d (badan %s)", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "Pega Non Proportional Calculation /Perhitungan Pega Non Proportional") {
		t.Errorf("nama Non-Prop tidak sampai: %s", w.Body.String())
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
// Panel `Existing Policy for Master ID` - `TREATYINPRODUCTION`, tabel warisan.
func (g gudangTiruan) BacaPolisMaster(_ context.Context, _ string) ([]models.BarisPolisMaster, error) {
	return []models.BarisPolisMaster{}, g.galatLampiran
}

func (g gudangTiruan) BacaRiwayatKontrak(_ context.Context, _ string) ([]models.BarisRiwayatWarisan, error) {
	return g.riwayat, g.galatLampiran
}
