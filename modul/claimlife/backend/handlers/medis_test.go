package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/inti/backend/kontrak"
	"nusantarare/modul/claimlife/backend/services"
)

// 11-10-2026 - tombol Save layar Medical Check.

func TestRuteMedisTerdaftar(t *testing.T) {
	isi, err := os.ReadFile("handlers.go")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(isi), `"PUT /api/klaim-life/{id}/medis", simpanMedis(`) {
		t.Error("rute Save layar Medical Check tidak terdaftar")
	}
}

func TestGalatMedisDipetakan(t *testing.T) {
	for _, u := range []struct {
		nama string
		err  error
		mau  int
	}{
		{"identitas", inti.ErrTanpaIdentitas, http.StatusUnauthorized},
		{"bukan Medical Advisor", inti.ErrTanpaWewenang, http.StatusForbidden},
		{"tertutup", kontrak.ErrKasusSudahTertutup, http.StatusConflict},
		{"bukan Medical Check", services.ErrMedisBukanMedicalCheck, http.StatusConflict},
		{"peserta diputus", services.ErrRekomendasiTerkunci, http.StatusConflict},
		{"kepanjangan", services.ErrRekomendasiKepanjangan, http.StatusUnprocessableEntity},
		{"klaim tidak ada", services.ErrKlaimTidakAda, http.StatusNotFound},
		{"permintaan", galat.ErrPermintaanTidakSah, http.StatusBadRequest},
		{"lain", errors.New("x"), http.StatusInternalServerError},
	} {
		w := httptest.NewRecorder()
		jawabGalatMedis(w, u.err)
		if w.Code != u.mau {
			t.Errorf("%s: kode %d, mau %d", u.nama, w.Code, u.mau)
		}
	}
}
