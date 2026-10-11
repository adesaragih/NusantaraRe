package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"nusantarare/modul/endorsmentfacin/backend/services"
	"nusantarare/modul/endorsmentfacin/backend/tiruan"
)

func kirim(t *testing.T, h http.Handler, badan string) *httptest.ResponseRecorder {
	t.Helper()
	r := httptest.NewRequest(http.MethodPost, Prefix+"/periksa", strings.NewReader(badan))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	return w
}

func routerUji(g *tiruan.Gudang, adaDB bool) http.Handler {
	return RouterDengan(services.BaruLayanan(g).DenganTransaksi(g.Transaksi), adaDB, true)
}

func TestPeriksaTanpaDatabase503(t *testing.T) {
	if w := kirim(t, routerUji(&tiruan.Gudang{}, false), `{"policyNo":"UJI.01.0001"}`); w.Code != http.StatusServiceUnavailable {
		t.Fatalf("kode %d", w.Code)
	}
}

func TestPeriksaKodeGalat(t *testing.T) {
	h := routerUji(&tiruan.Gudang{JenisBisnisNilai: "Fire"}, true)
	for _, k := range []struct {
		badan string
		kode  int
	}{
		{`{`, http.StatusBadRequest},
		{`{"policyNo":"UJI.01.0001","lain":1}`, http.StatusBadRequest},
		{`{"policyNo":""}`, http.StatusUnprocessableEntity},
		{`{"policyNo":"UJI.01.0001","edmType":"3"}`, http.StatusUnprocessableEntity},
		{`{"policyNo":"UJI.01.0001","endorsementDate":"2025-01-01"}`, http.StatusUnprocessableEntity},
		{`{"policyNo":"UJI.01.0001","edmType":"","endorsementDate":""}`, http.StatusOK},
		{`{"policyNo":"` + strings.Repeat("x", 20<<10) + `"}`, http.StatusRequestEntityTooLarge},
	} {
		if w := kirim(t, h, k.badan); w.Code != k.kode {
			t.Errorf("%.60s: kode %d, ingin %d (%s)", k.badan, w.Code, k.kode, w.Body.String())
		}
	}
	// Pesan 422 tanpa awalan penanda galat.
	w := kirim(t, h, `{"policyNo":" "}`)
	if !strings.Contains(w.Body.String(), "Policy No is required") || strings.Contains(w.Body.String(), "services:") {
		t.Fatalf("badan 422 = %s", w.Body.String())
	}
}

func TestPeriksaGalatGudang500(t *testing.T) {
	w := kirim(t, routerUji(&tiruan.Gudang{GalatStatusEDM: errors.New("ORA-03113 rincian")}, true), `{"policyNo":"UJI.01.0001"}`)
	if w.Code != http.StatusInternalServerError || strings.Contains(w.Body.String(), "ORA-03113") {
		t.Fatalf("kode %d badan %s", w.Code, w.Body.String())
	}
}

// Bentuk kabel final (kontrak sesi utama): kunci camelCase, `gerbang` 6 entri, `tanggal` null tanpa tanggal.
func TestPeriksaBentukJawaban(t *testing.T) {
	w := kirim(t, routerUji(&tiruan.Gudang{JenisBisnisNilai: "Fire"}, true), `{"policyNo":"UJI.01.0001","edmType":"4","edmTypeNew":"1"}`)
	if w.Code != http.StatusOK {
		t.Fatalf("kode %d: %s", w.Code, w.Body.String())
	}
	var j map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &j); err != nil {
		t.Fatal(err)
	}
	for _, kunci := range []string{"ditolak", "lengkap", "pesan", "gerbang", "bisnis", "tanggal"} {
		if _, ada := j[kunci]; !ada {
			t.Errorf("kunci %s tidak ada: %s", kunci, w.Body.String())
		}
	}
	if j["tanggal"] != nil {
		t.Errorf("tanggal = %v, ingin null", j["tanggal"])
	}
	if p, ok := j["pesan"].([]any); !ok || len(p) != 0 {
		t.Errorf("pesan = %v, ingin []", j["pesan"])
	}
	g, _ := j["gerbang"].([]any)
	if len(g) != 6 {
		t.Fatalf("gerbang = %v", j["gerbang"])
	}
	satu, _ := g[0].(map[string]any)
	for _, kunci := range []string{"no", "hasil", "catatan"} {
		if _, ada := satu[kunci]; !ada {
			t.Errorf("gerbang tanpa kunci %s", kunci)
		}
	}
	b, _ := j["bisnis"].(map[string]any)
	if b["hasil"] != "lolos" || b["jenisBisnis"] != "Fire" {
		t.Errorf("bisnis = %v", b)
	}
}

const badanKasus = `{"policyNo":"UJI.01.0001","edmType":"4","edmTypeNew":"1","detailType":"2",
 "endorsementDate":"01-06-2025","endorsementInternalRetro":"0","note":"Catatan uji","endorsementSource":"Ceding",
 "endorsementSourceNote":"Surat uji","endorsementStatus":"0","adminCharge":true,"survey":false}`

func kirimKasus(t *testing.T, h http.Handler, badan, akun string) *httptest.ResponseRecorder {
	t.Helper()
	r := httptest.NewRequest(http.MethodPost, Prefix+"/kasus", strings.NewReader(badan))
	if akun != "" {
		r.Header.Set("X-Pelaku", akun)
	}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	return w
}

func gudangKasus() *tiruan.Gudang {
	return tiruan.Bersih()
}

func TestKasusLahir201(t *testing.T) {
	g := gudangKasus()
	w := kirimKasus(t, routerUji(g, true), badanKasus, "uji.akun")
	if w.Code != http.StatusCreated {
		t.Fatalf("kode %d: %s", w.Code, w.Body.String())
	}
	var j map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &j); err != nil {
		t.Fatal(err)
	}
	if j["caseId"] != "EDM-00001" || j["periksa"] != nil {
		t.Fatalf("badan = %s", w.Body.String())
	}
	if p, ok := j["pesan"].([]any); !ok || len(p) != 0 {
		t.Fatalf("pesan = %v, ingin []", j["pesan"])
	}
	if g.WorkPolis["EDM-00001"] != "uji.akun" {
		t.Fatalf("pembuat = %v", g.WorkPolis)
	}
}

func TestKasusDitolak200(t *testing.T) {
	g := gudangKasus()
	g.Renewal = true
	w := kirimKasus(t, routerUji(g, true), badanKasus, "uji.akun")
	if w.Code != http.StatusOK {
		t.Fatalf("kode %d: %s", w.Code, w.Body.String())
	}
	var j struct {
		CaseID  string `json:"caseId"`
		Periksa *struct {
			Ditolak bool `json:"ditolak"`
			Pesan   []struct {
				Medan, Teks string
			} `json:"pesan"`
		} `json:"periksa"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &j); err != nil {
		t.Fatal(err)
	}
	if j.CaseID != "" || j.Periksa == nil || !j.Periksa.Ditolak || len(j.Periksa.Pesan) != 1 {
		t.Fatalf("badan = %s", w.Body.String())
	}
	if len(g.WorkPolis) != 0 {
		t.Fatalf("ditolak tetapi menulis: %v", g.WorkPolis)
	}
}

func dengan(ubah func(*tiruan.Gudang)) *tiruan.Gudang {
	g := tiruan.Bersih()
	ubah(g)
	return g
}

func TestKasusKodeGalat(t *testing.T) {
	for _, k := range []struct {
		nama  string
		g     *tiruan.Gudang
		badan string
		akun  string
		kode  int
		teks  string
	}{
		{"tanpa akun", gudangKasus(), badanKasus, "", http.StatusForbidden, "Log in to create an endorsement"},
		{"note kosong", gudangKasus(), strings.Replace(badanKasus, `"Catatan uji"`, `""`, 1), "uji.akun",
			http.StatusUnprocessableEntity, "Endorsement Note is required"},
		{"medan asing", gudangKasus(), strings.Replace(badanKasus, `"survey"`, `"lain"`, 1), "uji.akun",
			http.StatusBadRequest, ""},
		{"sequence belum ada", dengan(func(g *tiruan.Gudang) {
			g.GalatNomor = fmt.Errorf("%w: ORA-02289", services.ErrBelumAda)
		}), badanKasus, "uji.akun", http.StatusServiceUnavailable, "a required table or sequence is not in this schema"},
		{"galat lain", dengan(func(g *tiruan.Gudang) { g.GalatSisipKasus = errors.New("ORA-12899 rincian") }),
			badanKasus, "uji.akun", http.StatusInternalServerError, "rinciannya di log server"},
	} {
		w := kirimKasus(t, routerUji(k.g, true), k.badan, k.akun)
		if w.Code != k.kode || !strings.Contains(w.Body.String(), k.teks) || strings.Contains(w.Body.String(), "ORA-") {
			t.Errorf("%s: kode %d badan %s", k.nama, w.Code, w.Body.String())
		}
	}
	if w := kirimKasus(t, routerUji(gudangKasus(), false), badanKasus, "uji.akun"); w.Code != http.StatusServiceUnavailable {
		t.Fatalf("tanpa DB: %d", w.Code)
	}
}

func TestPeriksaHanyaPOST(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, Prefix+"/periksa", nil)
	w := httptest.NewRecorder()
	routerUji(&tiruan.Gudang{}, true).ServeHTTP(w, r)
	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("GET = %d", w.Code)
	}
}
