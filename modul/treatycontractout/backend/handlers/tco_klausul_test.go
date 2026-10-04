package handlers

// Uji pintu HTTP klausul - TANPA Oracle (tiket 08).

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"nusantarare/modul/treatycontractout/backend/models"
	"nusantarare/modul/treatycontractout/backend/services"
)

func TestRuteKlausulTerdaftarDanTanpaDatabase503(t *testing.T) {
	router := Router(services.New(nil), true)
	for _, r := range []struct{ metode, jalur string }{
		{http.MethodGet, "/api/treaty-contract-out/jenis-klausul?isXol=0"},
		{http.MethodGet, "/api/treaty-contract-out/klausul-pilihan/occupation?cari=uji"},
		{http.MethodGet, "/api/treaty-contract-out/tahun/1000001/klausul?descId=10009"},
		{http.MethodPost, "/api/treaty-contract-out/tahun/1000001/klausul"},
		{http.MethodPut, "/api/treaty-contract-out/tahun/1000001/klausul/10000001"},
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
	if strings.Count(string(isi), "daftarkanRuteKlausulTCO(mux, svc, stubPelaku)") != 1 {
		t.Error("rute klausul harus didaftarkan tepat sekali")
	}
	// Klausul tidak dihapus lewat konteks ini (kaskade tiket 10 pun tidak menyentuhnya).
	sumber, _ := os.ReadFile("tco_klausul.go")
	if strings.Contains(string(sumber), `"DELETE `) {
		t.Error("jalur hapus klausul tidak ada di tiket mana pun")
	}
}

func TestJawabGalatKlausulTCO(t *testing.T) {
	for _, k := range []struct {
		err  error
		kode int
	}{
		{services.ErrKlausulTidakAda, http.StatusNotFound},
		{services.GalatKlausulDobel{IDLain: "10000009", Jenis: "EPI"}, http.StatusConflict},
		{services.ErrKlausulIndukBeranak, http.StatusConflict},
		{services.ErrKlausulJenisBerubah, http.StatusBadRequest},
		{models.ErrKlausulJenisTakDikenal, http.StatusUnprocessableEntity},
		{models.ErrKlausulDitahan, http.StatusUnprocessableEntity},
		{models.ErrKlausulMedanWajib, http.StatusUnprocessableEntity},
		{models.ErrMedanBukanMilikJenis, http.StatusUnprocessableEntity},
		{models.ErrTotalPctAnakMelebihi100, http.StatusUnprocessableEntity},
		{services.ErrJenisKlausulDiLuarMaster, http.StatusUnprocessableEntity},
		{services.ErrPilihanDiLuarMaster, http.StatusUnprocessableEntity},
	} {
		w := httptest.NewRecorder()
		jawabGalatTreatyContractOut(w, k.err)
		if w.Code != k.kode {
			t.Errorf("%v: kode %d, mau %d", k.err, w.Code, k.kode)
		}
	}
}
