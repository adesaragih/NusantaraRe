package handlers

// Uji pintu HTTP reinsurer - TANPA Oracle (tiket 05).

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"nusantarare/internal/models"
	"nusantarare/internal/services"
)

func TestRuteReinsurerTerdaftarDanTanpaDatabase503(t *testing.T) {
	router := Router(services.New(nil), true)
	for _, r := range []struct{ metode, jalur string }{
		{http.MethodGet, "/api/treaty-contract-out/reinsurer-master?cari=uji"},
		{http.MethodGet, "/api/treaty-contract-out/tahun/1000001/kontrak/1000003/reinsurer"},
		{http.MethodPost, "/api/treaty-contract-out/tahun/1000001/kontrak/1000003/reinsurer"},
		{http.MethodPut, "/api/treaty-contract-out/tahun/1000001/kontrak/1000003/reinsurer/1000005"},
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
	if strings.Count(string(isi), "daftarkanRuteReinsurerTCO(mux, svc, stubPelaku)") != 1 {
		t.Error("rute reinsurer harus didaftarkan tepat sekali")
	}
	sumber, _ := os.ReadFile("tco_reinsurer.go")
	if strings.Contains(string(sumber), `"DELETE `) {
		t.Error("hapus reinsurer milik tiket 10 (DeleteTreatyReins_Act), bukan tiket 05")
	}
}

func TestJawabGalatReinsurerTCO(t *testing.T) {
	for _, k := range []struct {
		err  error
		kode int
	}{
		{services.ErrReinsurerTidakAda, http.StatusNotFound},
		{services.ErrReinsurerDiLuarMaster, http.StatusUnprocessableEntity},
		{models.ErrReinsurerKosong, http.StatusUnprocessableEntity},
		{models.ErrPersenKosong, http.StatusUnprocessableEntity},
		{models.ErrPersenBukanDesimal, http.StatusUnprocessableEntity},
		{models.ErrPersenDiLuarRentang, http.StatusUnprocessableEntity},
		{models.ErrTotalShareMelebihi100, http.StatusUnprocessableEntity},
	} {
		w := httptest.NewRecorder()
		jawabGalatTreatyContractOut(w, k.err)
		if w.Code != k.kode {
			t.Errorf("%v: kode %d, mau %d", k.err, w.Code, k.kode)
		}
	}
}
