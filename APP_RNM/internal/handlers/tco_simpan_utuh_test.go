package handlers

// Uji pintu HTTP simpan utuh - TANPA Oracle (tiket 09).

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"nusantarare/internal/services"
)

func TestRuteSimpanUtuhTerdaftarDanTanpaDatabase503(t *testing.T) {
	router := Router(services.New(nil), true)
	for _, r := range []struct{ metode, jalur string }{
		{http.MethodPost, "/api/treaty-contract-out/tahun/1000001/kontrak-utuh"},
		{http.MethodPut, "/api/treaty-contract-out/tahun/1000001/kontrak/1000003/utuh"},
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
	if strings.Count(string(isi), "daftarkanRuteSimpanUtuhTCO(mux, svc, stubPelaku)") != 1 {
		t.Error("rute simpan utuh harus didaftarkan tepat sekali")
	}
}

// AC 38: kode dari galat baris; pesan menyebut baris itu; galat server tidak bocor.
func TestJawabGalatUtuhTCO(t *testing.T) {
	for _, k := range []struct {
		err        error
		kode       int
		ada, tiada string
	}{
		{services.GalatSimpanUtuhTCO{Bagian: "klausul ke-3", Galat: services.GalatKlausulDobel{IDLain: "1000009", Jenis: "EPI"}},
			http.StatusConflict, "klausul ke-3", ""},
		{services.GalatSimpanUtuhTCO{Bagian: "reinsurer ke-2", Galat: services.ErrReinsurerTidakAda},
			http.StatusNotFound, "reinsurer ke-2", ""},
		{services.GalatSimpanUtuhTCO{Bagian: "jejak", Galat: errors.New("ORA-00001 rincian internal")},
			http.StatusInternalServerError, "gagal pada jejak (galat server)", "ORA-00001"},
		{services.ErrSimpanUtuhTidakSah, http.StatusBadRequest, "tidak sah", ""},
		{services.GalatSimpanUtuhTCO{Bagian: "klausul ke-2", Galat: services.ErrMasterKursRusak},
			http.StatusServiceUnavailable, "klausul ke-2", "(galat server)"},
	} {
		w := httptest.NewRecorder()
		if !jawabGalatUtuhTCO(w, k.err) {
			t.Fatal("tidak dijawab")
		}
		b := w.Body.String()
		if w.Code != k.kode || !strings.Contains(b, k.ada) || (k.tiada != "" && strings.Contains(b, k.tiada)) {
			t.Errorf("%v: %d %s", k.err, w.Code, b)
		}
	}
}
