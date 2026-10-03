package handlers

// Pintu HTTP tombol portal `Input Offer` / `Input Premium` - GILIRAN-13 paket 1.

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/modul/premiumlistlife/backend/services"
)

func TestRuteBuatKasusPolisTerdaftar(t *testing.T) {
	// Refactor bentuk B (30-09-2026): tabel rute modul ini kini
	// `rute_premiumlist.go` (DaftarkanRute), bukan Router bersama.
	isi, err := os.ReadFile("rute_premiumlist.go")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(isi), `"POST /api/polis-life",`) {
		t.Error("rute POST /api/polis-life tidak terdaftar di Router")
	}
}

func TestBuatKasusPolisTanpaDatabaseMenjawab503(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/api/polis-life", strings.NewReader(`{"flag":"0"}`))
	buatKasusPolis(services.New(nil), true)(w, r)
	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("kode = %d, mau 503", w.Code)
	}
}

func TestGalatBuatKasusPolisDipetakan(t *testing.T) {
	for _, u := range []struct {
		nama string
		err  error
		mau  int
	}{
		{"identitas", inti.ErrTanpaIdentitas, http.StatusUnauthorized},
		{"bendera", galat.ErrPermintaanTidakSah, http.StatusBadRequest},
	} {
		w := httptest.NewRecorder()
		if !jawabGalatPolis(w, u.err) || w.Code != u.mau {
			t.Errorf("%s: kode = %d, mau %d", u.nama, w.Code, u.mau)
		}
	}
}
