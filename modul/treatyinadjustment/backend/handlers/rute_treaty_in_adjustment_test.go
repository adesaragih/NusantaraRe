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
