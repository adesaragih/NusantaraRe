package handlers

// Uji pintu HTTP kontrak - TANPA Oracle (tiket 04).

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"nusantarare/modul/treatycontractout/models"
	"nusantarare/modul/treatycontractout/services"
)

func TestRuteKontrakTerdaftarDanTanpaDatabase503(t *testing.T) {
	router := Router(services.New(nil), true)
	for _, r := range []struct{ metode, jalur string }{
		{http.MethodGet, "/api/treaty-contract-out/tahun/1000001/kontrak"},
		{http.MethodPost, "/api/treaty-contract-out/tahun/1000001/kontrak"},
		{http.MethodPut, "/api/treaty-contract-out/tahun/1000001/kontrak/1000003"},
		{http.MethodGet, "/api/treaty-contract-out/tahun/1000001/kontrak/akhir-bawaan?mulai=2026-01-01"},
	} {
		w := httptest.NewRecorder()
		q := httptest.NewRequest(r.metode, r.jalur, strings.NewReader("{}"))
		q.Header.Set("X-Pelaku", "UJI-ADMIN")
		router.ServeHTTP(w, q)
		if w.Code != http.StatusServiceUnavailable {
			t.Errorf("%s %s: kode %d, mau 503", r.metode, r.jalur, w.Code)
		}
	}
	isi, err := os.ReadFile("rute_treaty_contract_out.go")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Count(string(isi), "daftarkanRuteKontrakTCO(mux, svc, stubPelaku)") != 1 {
		t.Error("rute kontrak harus didaftarkan tepat sekali dari pendaftar modul")
	}
	// Tiket 10 yang membawa jalur hapus kontrak - bukan tiket ini.
	sumber, _ := os.ReadFile("tco_kontrak.go")
	if strings.Contains(string(sumber), `"DELETE `) {
		t.Error("jalur hapus kontrak milik tiket 10 (kaskade), bukan tiket 04")
	}
}

func TestJawabGalatKontrakTCO(t *testing.T) {
	for _, k := range []struct {
		err  error
		kode int
	}{
		{services.ErrKontrakTidakAda, http.StatusNotFound},
		{services.GalatKontrakDobel{IDLain: "1000004", ReinsTypeID: "10003"}, http.StatusConflict},
		{models.ErrKontrakJenisReasuransiKosong, http.StatusUnprocessableEntity},
		{models.ErrKontrakMulaiKosong, http.StatusUnprocessableEntity},
		{models.ErrKontrakAkhirKosong, http.StatusUnprocessableEntity},
		{models.ErrKontrakTahunMulaiBeda, http.StatusUnprocessableEntity},
		{services.ErrJenisReasuransiDiLuarDaftar, http.StatusUnprocessableEntity},
	} {
		w := httptest.NewRecorder()
		jawabGalatTreatyContractOut(w, k.err)
		if w.Code != k.kode {
			t.Errorf("%v: kode %d, mau %d", k.err, w.Code, k.kode)
		}
	}
	w := httptest.NewRecorder()
	jawabGalatTreatyContractOut(w, services.GalatKontrakDobel{IDLain: "1000004", ReinsTypeID: "10003"})
	if !strings.Contains(w.Body.String(), "Data sudah pernah di Input") || !strings.Contains(w.Body.String(), "1000004") {
		t.Errorf("pesan 409: %s", w.Body.String())
	}
}

// AC 5: POST yang membawa id dan PUT yang id-nya berbeda ditolak di pintu.
func TestSimpanKontrakMenolakIdentitasDariKlien(t *testing.T) {
	sumber, err := os.ReadFile("tco_kontrak.go")
	if err != nil {
		t.Fatal(err)
	}
	for _, mau := range []string{`!perbarui && id != ""`, `perbarui && id != "" && id != r.PathValue("kid")`,
		`masuk.ID = r.PathValue("kid")`} {
		if !strings.Contains(string(sumber), mau) {
			t.Errorf("gerbang identitas %q hilang", mau)
		}
	}
}
