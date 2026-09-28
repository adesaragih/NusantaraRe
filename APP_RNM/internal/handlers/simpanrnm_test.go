package handlers

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

func TestRuteSimpanRNMTerdaftar(t *testing.T) {
	isi, err := os.ReadFile("handlers.go")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(isi), `"POST /api/klaim-life/{id}/outstanding"`) {
		t.Error("rute Save to RNM tidak terdaftar di Router")
	}
}

// TestPelanggaranRNMDijawab422DenganKalimatXML - kalimat apa adanya.
func TestPelanggaranRNMDijawab422DenganKalimatXML(t *testing.T) {
	w := httptest.NewRecorder()
	jawabGalatSimpanRNM(w, &services.PelanggaranRNM{Langkah: "11.12", Pesan: "DOB cannnot be blank No 1"})
	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("kode = %d, mau 422", w.Code)
	}
	var isi struct{ Galat, Langkah string }
	if err := json.Unmarshal(w.Body.Bytes(), &isi); err != nil {
		t.Fatal(err)
	}
	if isi.Galat != "DOB cannnot be blank No 1" || isi.Langkah != "11.12" {
		t.Errorf("badan = %+v", isi)
	}
}

func TestGalatSimpanRNMDipetakanKeKodeYangBenar(t *testing.T) {
	for _, u := range []struct {
		nama string
		err  error
		mau  int
	}{
		{"identitas", services.ErrTanpaIdentitas, http.StatusUnauthorized},
		{"wewenang", services.ErrTanpaWewenang, http.StatusForbidden},
		{"tertutup", services.ErrKasusSudahTertutup, http.StatusConflict},
		{"tahap", services.ErrSimpanRNMBukanOutstanding, http.StatusConflict},
		{"type", services.ErrTypeTidakDikenal, http.StatusUnprocessableEntity},
		{"ambang", services.ErrAmbangSTNCBelumDiketahui, http.StatusUnprocessableEntity},
		{"kategori", services.ErrKategoriWajibBelumDiketahui, http.StatusUnprocessableEntity},
		{"produk", services.ErrAmbangProdukTakDitemukan, http.StatusUnprocessableEntity},
		{"polis", services.ErrPolisNomorTakDitemukan, http.StatusUnprocessableEntity},
		{"permintaan", services.ErrPermintaanTidakSah, http.StatusBadRequest},
		{"lain", errors.New("x"), http.StatusInternalServerError},
	} {
		w := httptest.NewRecorder()
		if !jawabGalatSimpanRNM(w, u.err) {
			t.Errorf("%s: tidak dijawab", u.nama)
			continue
		}
		if w.Code != u.mau {
			t.Errorf("%s: kode = %d, mau %d", u.nama, w.Code, u.mau)
		}
	}
	if jawabGalatSimpanRNM(httptest.NewRecorder(), nil) {
		t.Error("nil dijawab")
	}
}

func TestSimpanRNMTanpaDatabaseMenjawab503(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/api/klaim-life/K-1/outstanding", nil)
	r.SetPathValue("id", "K-1")
	simpanRNM(services.New(nil), true)(w, r)
	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("kode = %d, mau 503", w.Code)
	}
}
