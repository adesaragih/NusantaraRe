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

type gudangTiruan struct{ isi []models.Acuan }

func (g gudangTiruan) DaftarAcuan(context.Context, models.Himpunan) ([]models.Acuan, error) {
	return g.isi, nil
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

	w := minta(t, h, handlers.Prefix+"/acuan/mata-uang", "AKUN-UJI")

	if w.Code != http.StatusServiceUnavailable {
		t.Fatalf("mau 503, dapat %d", w.Code)
	}
}

func TestTanpaIdentitasMenjawab401(t *testing.T) {
	h := handlers.RouterDengan(services.LayananDengan(gudangTiruan{}), true, true)

	w := minta(t, h, handlers.Prefix+"/acuan/mata-uang", "")

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("mau 401, dapat %d (badan %s)", w.Code, w.Body.String())
	}
}

// Himpunan di luar keenam menjawab 404, dan pesannya MENYEBUT yang diminta -
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

// Uji POSITIF: keenam himpunan menjawab 200 berbadan daftar.
func TestKeenamHimpunanMenjawab200(t *testing.T) {
	h := handlers.RouterDengan(
		services.LayananDengan(gudangTiruan{isi: []models.Acuan{{ID: 7, Kode: "IDR", Nama: "Rupiah", Aktif: "1"}}}),
		true, true)

	for _, himpunan := range []string{
		"mata-uang", "jenis-potongan", "kelas-bisnis", "kelompok-treaty", "bahaya", "jenis-reasuransi",
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
