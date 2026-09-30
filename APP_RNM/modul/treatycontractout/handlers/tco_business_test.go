package handlers

// Uji pintu HTTP business - TANPA Oracle (tiket 07).

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"nusantarare/modul/treatycontractout/models"
	"nusantarare/modul/treatycontractout/services"
)

func TestRuteBusinessTerdaftarDanTanpaDatabase503(t *testing.T) {
	router := Router(services.New(nil), true)
	for _, r := range []struct{ metode, jalur string }{
		{http.MethodGet, "/api/treaty-contract-out/business-master"},
		{http.MethodGet, "/api/treaty-contract-out/tahun/1000001/kontrak/1000003/business"},
		{http.MethodPost, "/api/treaty-contract-out/tahun/1000001/kontrak/1000003/business"},
		{http.MethodPut, "/api/treaty-contract-out/tahun/1000001/kontrak/1000003/business/1000002"},
		{http.MethodDelete, "/api/treaty-contract-out/tahun/1000001/kontrak/1000003/business/1000002"},
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
	if strings.Count(string(isi), "daftarkanRuteBusinessTCO(mux, svc, stubPelaku)") != 1 {
		t.Error("rute bisnis harus didaftarkan tepat sekali")
	}
}

func TestJawabGalatBusinessTCO(t *testing.T) {
	for _, k := range []struct {
		err  error
		kode int
	}{
		{services.ErrBusinessTidakAda, http.StatusNotFound},
		{services.GalatBusinessDobel{IDLain: "1000004", BizCode: "UJI-B1"}, http.StatusConflict},
		{services.ErrBusinessDiLuarMaster, http.StatusUnprocessableEntity},
		{models.ErrBusinessKodeKosong, http.StatusUnprocessableEntity},
		{models.ErrBusinessAktifTakSah, http.StatusUnprocessableEntity},
	} {
		w := httptest.NewRecorder()
		jawabGalatTreatyContractOut(w, k.err)
		if w.Code != k.kode {
			t.Errorf("%v: kode %d, mau %d", k.err, w.Code, k.kode)
		}
	}
}
