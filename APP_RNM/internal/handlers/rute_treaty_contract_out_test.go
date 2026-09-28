package handlers

// Uji pintu HTTP Treaty Contract Out - TANPA Oracle (tiket 02).

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"nusantarare/internal/services"
)

func TestRuteTreatyContractOutTerdaftarSatuBaris(t *testing.T) {
	isi, err := os.ReadFile("handlers.go")
	if err != nil {
		t.Fatal(err)
	}
	// ⛔ SATU sentuhan ke handlers.go: pemanggilan pendaftar modul.
	if strings.Count(string(isi), "daftarkanRuteTreatyContractOut(mux, svc, stubPelaku)") != 1 {
		t.Error("Router harus memanggil daftarkanRuteTreatyContractOut tepat sekali")
	}
	if strings.Contains(string(isi), "/api/treaty-contract-out") {
		t.Error("rute modul ditulis di handlers.go; ia milik rute_treaty_contract_out.go")
	}
	rute, err := os.ReadFile("rute_treaty_contract_out.go")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(rute), `"GET /api/treaty-contract-out/jenis-reasuransi"`) {
		t.Error("rute jenis reasuransi tidak terdaftar sebagai GET")
	}
	if strings.Contains(string(rute), `"POST /api/treaty-contract-out/jenis-reasuransi"`) {
		t.Error("daftar master terdaftar sebagai POST; ia membaca")
	}
}

func TestJenisReasuransiTanpaDatabaseMenjawab503(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/api/treaty-contract-out/jenis-reasuransi", nil)
	jenisReasuransiTreaty(services.New(nil), true)(w, r)
	if w.Code != http.StatusServiceUnavailable {
		t.Fatalf("kode = %d, mau 503", w.Code)
	}
	var isi map[string]any
	if err := json.NewDecoder(w.Body).Decode(&isi); err != nil {
		t.Fatal(err)
	}
	if _, ada := isi["galat"]; !ada {
		t.Error("envelope tanpa kunci galat")
	}
}

func TestJawabGalatTreatyContractOut(t *testing.T) {
	kasus := []struct {
		err  error
		kode int
	}{
		{services.ErrTanpaIdentitas, http.StatusUnauthorized},
		{services.ErrTanpaWewenang, http.StatusForbidden},
		{services.ErrMasterJenisReasuransiKosong, http.StatusServiceUnavailable},
		{services.ErrPermintaanTidakSah, http.StatusBadRequest},
		{errors.New("UJI: galat lain"), http.StatusInternalServerError},
	}
	for _, k := range kasus {
		w := httptest.NewRecorder()
		if !jawabGalatTreatyContractOut(w, k.err) {
			t.Errorf("%v: tidak dijawab", k.err)
		}
		if w.Code != k.kode {
			t.Errorf("%v: kode %d, mau %d", k.err, w.Code, k.kode)
		}
	}
	// Master kosong: pesannya menyebut masternya (ADR-0015).
	w := httptest.NewRecorder()
	jawabGalatTreatyContractOut(w, services.ErrMasterJenisReasuransiKosong)
	if !strings.Contains(w.Body.String(), "REINSURANCETYPE") {
		t.Errorf("pesan 503 tidak menyebut masternya: %s", w.Body.String())
	}
	if jawabGalatTreatyContractOut(httptest.NewRecorder(), nil) {
		t.Error("nil dijawab sebagai galat")
	}
}
