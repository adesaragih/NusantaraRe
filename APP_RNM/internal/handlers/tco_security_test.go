package handlers

// Uji pintu HTTP security - TANPA Oracle (tiket 06).

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"nusantarare/internal/models"
	"nusantarare/internal/services"
)

func TestRuteSecurityTerdaftarDanTanpaDatabase503(t *testing.T) {
	router := Router(services.New(nil), true)
	dasar := "/api/treaty-contract-out/tahun/1000001/kontrak/1000003/reinsurer/1000007/security"
	for _, r := range []struct{ metode, jalur string }{
		{http.MethodGet, dasar},
		{http.MethodPost, dasar},
		{http.MethodPut, dasar + "/1000011"},
		{http.MethodDelete, dasar + "/1000011"},
	} {
		w := httptest.NewRecorder()
		q := httptest.NewRequest(r.metode, r.jalur, strings.NewReader("{}"))
		q.Header.Set("X-Pelaku", "UJI-ADMIN")
		router.ServeHTTP(w, q)
		if w.Code != http.StatusServiceUnavailable {
			t.Errorf("%s %s: kode %d, mau 503", r.metode, r.jalur, w.Code)
		}
	}
	isi, _ := os.ReadFile("rute_treaty_contract_out.go")
	if strings.Count(string(isi), "daftarkanRuteSecurityTCO(mux, svc, stubPelaku)") != 1 {
		t.Error("rute security harus didaftarkan tepat sekali")
	}
}

func TestJawabGalatSecurityTCO(t *testing.T) {
	for _, k := range []struct {
		err  error
		kode int
	}{
		{services.ErrSecurityTidakAda, http.StatusNotFound},
		{services.GalatSecurityDobel{IDLain: "1000011", ReasSecurity: "UJI-R1"}, http.StatusConflict},
		{models.ErrSecurityKosong, http.StatusUnprocessableEntity},
		{models.ErrSecurityTanpaReinsurer, http.StatusUnprocessableEntity},
		{models.ErrSecurityMelampauiLebar, http.StatusUnprocessableEntity},
	} {
		w := httptest.NewRecorder()
		jawabGalatTreatyContractOut(w, k.err)
		if w.Code != k.kode {
			t.Errorf("%v: kode %d, mau %d", k.err, w.Code, k.kode)
		}
	}
}
