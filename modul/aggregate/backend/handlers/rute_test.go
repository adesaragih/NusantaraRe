package handlers_test

// Rute HTTP Aggregate di atas gudang tiruan: bentuk jawaban dan kode status.

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"nusantarare/modul/aggregate/backend/handlers"
	"nusantarare/modul/aggregate/backend/services"
	"nusantarare/modul/aggregate/backend/tiruan"
)

func router(g *tiruan.Gudang, adaDB bool) http.Handler {
	return handlers.RouterDengan(services.BaruLayanan(g, tiruan.Transaksi), adaDB, true)
}

func kirim(t *testing.T, h http.Handler, metode, jalur, badan string) *httptest.ResponseRecorder {
	t.Helper()
	r := httptest.NewRequest(metode, jalur, strings.NewReader(badan))
	r.Header.Set("X-Pelaku", "UJI-ADMIN")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	return w
}

func TestRuteDaftarRingkasanRincianHapus(t *testing.T) {
	h := router(tiruan.Contoh(), true)
	if w := kirim(t, h, "GET", handlers.Prefix+"?q=lama", ""); w.Code != 200 || !strings.Contains(w.Body.String(), `"tanggalInput":"01-03-2025"`) ||
		!strings.Contains(w.Body.String(), `"jumlahBaris":1`) {
		t.Errorf("daftar %d %s", w.Code, w.Body.String())
	}
	if w := kirim(t, h, "GET", handlers.Prefix+"/ringkasan", ""); w.Code != 200 || !strings.Contains(w.Body.String(),
		`{"irisan":[{"cedingCode":"UJI-C9","cedingName":"UJI CEDING LAMA","treatyType":"OR","coverage":"EQVET","rnmValueInUsd":"7"}]}`) {
		t.Errorf("ringkasan %d %s", w.Code, w.Body.String())
	}
	kunci := "tanggalInput=01-03-2025&cedingCode=UJI-C9&cedingName=" + url.QueryEscape("UJI CEDING LAMA") +
		"&treatyType=OR&asAt=30-09-2024&uwYear=2023"
	if w := kirim(t, h, "GET", handlers.Prefix+"/rincian?"+kunci, ""); w.Code != 200 || !strings.Contains(w.Body.String(), `"ID":"AGG-500"`) {
		t.Errorf("rincian %d %s", w.Code, w.Body.String())
	}
	badan := `{"tanggalInput":"01-03-2025","cedingCode":"UJI-C9","cedingName":"UJI CEDING LAMA","treatyType":"OR","asAt":"30-09-2024","uwYear":"2023"}`
	if w := kirim(t, h, "POST", handlers.Prefix+"/hapus", badan); w.Code != 200 || !strings.Contains(w.Body.String(), `"dihapus":1`) {
		t.Errorf("hapus %d %s", w.Code, w.Body.String())
	}
	if w := kirim(t, h, "POST", handlers.Prefix+"/hapus", badan); w.Code != 404 {
		t.Errorf("hapus kedua %d", w.Code)
	}
}

func TestRuteMasterTreatyPratinjauSimpan(t *testing.T) {
	g := tiruan.Contoh()
	h := router(g, true)
	if w := kirim(t, h, "GET", handlers.Prefix+"/master-treaty?q=satu", ""); w.Code != 200 || !strings.Contains(w.Body.String(), `"id":"UJI-V3"`) {
		t.Errorf("master treaty %d %s", w.Code, w.Body.String())
	}
	csv := "kepala\n" + strings.Repeat(";", 0) + "1.1 UJI;OR;EQVET;30/09/2024;2023;;;USD;;;;;;;;;;;;;;;;;;;;;;;;1;1.000;;;;;\n"
	w := kirim(t, h, "POST", handlers.Prefix+"/pratinjau", `{"csv":`+jsonTeks(csv)+`,"masterTreaty":["UJI-V1"]}`)
	if w.Code != 200 || !strings.Contains(w.Body.String(), `"TO_USD":"1"`) || !strings.Contains(w.Body.String(), `"ASSESMENT_ZONE":"Total :"`) {
		t.Fatalf("pratinjau %d %s", w.Code, w.Body.String())
	}
	// Save tanpa Master ID: pesan Pega, 422, dipisah baris baru.
	w = kirim(t, h, "POST", handlers.Prefix+"/simpan", `{"baris":[{"ASSESMENT_ZONE":"","TREATY_TYPE":"OR"}]}`)
	if w.Code != 422 || !strings.Contains(w.Body.String(), "Assesment Zone in list 1 Not Found") ||
		!strings.Contains(w.Body.String(), `\nTo USD in list 1 cannot be empty`) {
		t.Errorf("simpan ditolak %d %s", w.Code, w.Body.String())
	}
	if w := kirim(t, h, "POST", handlers.Prefix+"/simpan", `{"baris":[],"asing":1}`); w.Code != 400 {
		t.Errorf("badan asing %d", w.Code)
	}
}

func TestRuteTanpaDatabase503(t *testing.T) {
	if w := kirim(t, router(tiruan.Contoh(), false), "GET", handlers.Prefix, ""); w.Code != 503 {
		t.Errorf("tanpa database %d", w.Code)
	}
}

// jsonTeks - teks sebagai literal JSON.
func jsonTeks(s string) string {
	r := strings.NewReplacer(`\`, `\\`, `"`, `\"`, "\n", `\n`)
	return `"` + r.Replace(s) + `"`
}
