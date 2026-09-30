package handlers

// Uji pintu HTTP kurs - TANPA Oracle (tiket 11).

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"nusantarare/modul/treatycontractout/backend/models"
	"nusantarare/modul/treatycontractout/backend/services"
)

func TestRuteKursTerdaftarDanTanpaDatabase503(t *testing.T) {
	router := Router(services.New(nil), true)
	for _, jalur := range []string{
		"/api/treaty-contract-out/tahun/1000001/kurs",
		"/api/treaty-contract-out/tahun/1000001/kurs/konversi?dari=Rp&nilai=1",
	} {
		w := httptest.NewRecorder()
		q := httptest.NewRequest(http.MethodGet, jalur, nil)
		q.Header.Set("X-Pelaku", "UJI-ADMIN")
		router.ServeHTTP(w, q)
		if w.Code != http.StatusServiceUnavailable {
			t.Errorf("%s: kode %d, mau 503", jalur, w.Code)
		}
	}
	isi, _ := os.ReadFile("rute_treaty_contract_out.go")
	if strings.Count(string(isi), "daftarkanRuteKursTCO(mux, svc, stubPelaku)") != 1 {
		t.Error("rute kurs harus didaftarkan tepat sekali")
	}
	// Master kurs tidak ditulis: nol rute penulis kurs.
	sumber, _ := os.ReadFile("tco_kurs.go")
	for _, m := range []string{`"POST `, `"PUT `, `"DELETE `} {
		if strings.Contains(string(sumber), m) {
			t.Errorf("rute penulis kurs %s", m)
		}
	}
}

func TestJawabGalatKursTCO(t *testing.T) {
	for _, k := range []struct {
		err   error
		kode  int
		pesan string
	}{
		{models.GalatKursTidakAda{TreatyYear: "2026"}, http.StatusUnprocessableEntity, "No exchange rate for Treaty Year : 2026"},
		{errors.Join(services.ErrMasterKursRusak), http.StatusServiceUnavailable, ""},
		{services.ErrKonversiKursTidakSah, http.StatusBadRequest, ""},
	} {
		w := httptest.NewRecorder()
		jawabGalatTreatyContractOut(w, k.err)
		if w.Code != k.kode || !strings.Contains(w.Body.String(), k.pesan) {
			t.Errorf("%v: kode %d %s, mau %d", k.err, w.Code, w.Body.String(), k.kode)
		}
	}
}
