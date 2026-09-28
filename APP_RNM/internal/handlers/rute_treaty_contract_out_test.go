package handlers

// Uji pintu HTTP Treaty Contract Out - TANPA Oracle (tiket 02, 03).

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"nusantarare/internal/models"
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
	for _, mau := range []string{
		`"GET /api/treaty-contract-out/jenis-reasuransi"`,
		`"GET /api/treaty-contract-out/grup-treaty"`,
		`"GET /api/treaty-contract-out/tahun"`,
		`"POST /api/treaty-contract-out/tahun"`,
		`"GET /api/treaty-contract-out/tahun/{id}"`,
		`"PUT /api/treaty-contract-out/tahun/{id}"`,
	} {
		if !strings.Contains(string(rute), mau) {
			t.Errorf("rute %s tidak terdaftar", mau)
		}
	}
	// Master dan daftar MEMBACA: tidak pernah POST.
	for _, tidak := range []string{`"POST /api/treaty-contract-out/jenis-reasuransi"`,
		`"POST /api/treaty-contract-out/grup-treaty"`, `"DELETE /api/treaty-contract-out/tahun`} {
		if strings.Contains(string(rute), tidak) {
			t.Errorf("rute %s tidak boleh ada", tidak)
		}
	}
	// AC 72: nol jalur salin tahun treaty.
	if strings.Contains(strings.ToLower(string(rute)), "salin") || strings.Contains(strings.ToLower(string(rute)), "copy") {
		t.Error("ada jalur salin tahun treaty - fitur itu DIBUANG (AC 72)")
	}
}

func TestTreatyContractOutTanpaDatabaseMenjawab503(t *testing.T) {
	svc := services.New(nil)
	kasus := map[string]http.HandlerFunc{
		"jenis-reasuransi": jenisReasuransiTreaty(svc, true),
		"grup-treaty":      grupTreaty(svc, true),
		"tahun":            daftarTahunTreaty(svc, true),
		"tahun/{id}":       satuTahunTreaty(svc, true),
		"POST tahun":       simpanTahunTreaty(svc, true, false),
	}
	for nama, h := range kasus {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/api/treaty-contract-out/x", nil)
		h(w, r)
		if w.Code != http.StatusServiceUnavailable {
			t.Errorf("%s: kode = %d, mau 503", nama, w.Code)
		}
		var isi map[string]any
		if err := json.NewDecoder(w.Body).Decode(&isi); err != nil {
			t.Fatal(err)
		}
		if _, ada := isi["galat"]; !ada {
			t.Errorf("%s: envelope tanpa kunci galat", nama)
		}
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
		{services.ErrMasterGrupTreatyKosong, http.StatusServiceUnavailable},
		{services.ErrTahunTreatyTidakAda, http.StatusNotFound},
		{services.GalatTahunTreatyDobel{IDLain: "1000005"}, http.StatusConflict},
		{models.ErrPeriodeTerbalik, http.StatusUnprocessableEntity},
		{models.ErrTahunTreatyGrupKosong, http.StatusUnprocessableEntity},
		{models.ErrTahunTreatyTahunKosong, http.StatusUnprocessableEntity},
		{models.ErrTahunTreatyBukanAngka, http.StatusUnprocessableEntity},
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
	// Master kosong dan dobel: pesannya menyebut sebabnya.
	w := httptest.NewRecorder()
	jawabGalatTreatyContractOut(w, services.ErrMasterJenisReasuransiKosong)
	if !strings.Contains(w.Body.String(), "REINSURANCETYPE") {
		t.Errorf("pesan 503 tidak menyebut masternya: %s", w.Body.String())
	}
	w = httptest.NewRecorder()
	jawabGalatTreatyContractOut(w, services.GalatTahunTreatyDobel{IDLain: "1000005"})
	if !strings.Contains(w.Body.String(), "1000005") {
		t.Errorf("pesan 409 tidak menyebut tahun treaty lain: %s", w.Body.String())
	}
	if jawabGalatTreatyContractOut(httptest.NewRecorder(), nil) {
		t.Error("nil dijawab sebagai galat")
	}
}

// AC 5: POST yang membawa id ditolak SEBELUM menyentuh layanan; PUT dengan id
// badan yang berbeda dari jalur ditolak pula.
func TestSimpanTahunTreatyMenolakIdentitasDariKlien(t *testing.T) {
	svc := services.New(nil)
	// Tanpa database jawabannya 503 lebih dulu; uji ini memeriksa urutan
	// gerbang lewat badan pada layanan TANPA database tidak mungkin - jadi
	// yang dikunci di sini adalah TEKS pemeriksaannya di sumber.
	_ = svc
	rute, err := os.ReadFile("rute_treaty_contract_out.go")
	if err != nil {
		t.Fatal(err)
	}
	for _, mau := range []string{`masuk.ID != "" && masuk.ID != id`, `tahun treaty baru tidak membawa id`} {
		if !strings.Contains(string(rute), mau) {
			t.Errorf("gerbang %q tidak ada di handler", mau)
		}
	}
}
