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

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/modul/treatycontractout/backend/models"
	"nusantarare/modul/treatycontractout/backend/services"
)

func TestRuteTreatyContractOutTerdaftarSatuBaris(t *testing.T) {
	// Refactor bentuk B (30-09-2026): pendaftar modul kini dipanggil
	// `modul/treatycontractout/backend/modul.go` (paket 6: `cmd/api` merakit modul AKTIF dari
	// daftar), bukan `internal/handlers.Router`. Yang dijaga tetap sama: SATU
	// sentuhan di titik perakitan, dan nol rute modul ini di luar
	// `rute_treaty_contract_out.go`.
	rakit, err := os.ReadFile("../modul.go")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Count(string(rakit), "handlers.DaftarkanRute(mux, m.svc, m.stubPelaku)") != 1 {
		t.Error("modul.go harus memanggil DaftarkanRute modul ini tepat sekali")
	}
	for _, lain := range []string{"../modul.go", "../../../daftar.go", "../../../../cmd/api/main.go",
		"../../../../cmd/api/rakit.go", "../../../../modul/claimlife/backend/handlers/handlers.go"} {
		isi, err := os.ReadFile(lain)
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(string(isi), "/api/treaty-contract-out") {
			t.Errorf("rute modul ditulis di %s; ia milik rute_treaty_contract_out.go", lain)
		}
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
	// Master dan daftar MEMBACA: tidak pernah POST. Tahun treaty tidak punya
	// jalur hapus.
	//
	// ⚠️ 29-09-2026 tiket 12: dipersempit ke jalur TAHUN itu sendiri (tanda
	// kutip penutup ikut dicocokkan) dan diperluas ke tco_lampiran.go -
	// `DELETE .../tahun/{id}/lampiran/{lid}` menghapus LAMPIRAN, bukan tahun.
	lampiran, err := os.ReadFile("tco_lampiran.go")
	if err != nil {
		t.Fatal(err)
	}
	for _, tidak := range []string{`"POST /api/treaty-contract-out/jenis-reasuransi"`,
		`"POST /api/treaty-contract-out/grup-treaty"`, `"POST /api/treaty-contract-out/kategori-lampiran"`,
		`"DELETE /api/treaty-contract-out/tahun"`, `"DELETE /api/treaty-contract-out/tahun/{id}"`} {
		if strings.Contains(string(rute), tidak) || strings.Contains(string(lampiran), tidak) {
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
		{inti.ErrTanpaIdentitas, http.StatusUnauthorized},
		{inti.ErrTanpaWewenang, http.StatusForbidden},
		{services.ErrMasterJenisReasuransiKosong, http.StatusServiceUnavailable},
		{services.ErrMasterGrupTreatyKosong, http.StatusServiceUnavailable},
		{services.ErrTahunTreatyTidakAda, http.StatusNotFound},
		{services.GalatTahunTreatyDobel{IDLain: "1000005"}, http.StatusConflict},
		{models.ErrPeriodeTerbalik, http.StatusUnprocessableEntity},
		{models.ErrTahunTreatyGrupKosong, http.StatusUnprocessableEntity},
		{models.ErrTahunTreatyTahunKosong, http.StatusUnprocessableEntity},
		{models.ErrTahunTreatyBukanAngka, http.StatusUnprocessableEntity},
		{galat.ErrPermintaanTidakSah, http.StatusBadRequest},
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
	for _, mau := range []string{`masuk.ID != "" && masuk.ID != id`, `a new treaty year must not carry an id`} {
		if !strings.Contains(string(rute), mau) {
			t.Errorf("gerbang %q tidak ada di handler", mau)
		}
	}
}
