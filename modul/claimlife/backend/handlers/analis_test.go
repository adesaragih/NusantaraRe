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

// 12-10-2026 - tombol Save layar Claim Analis.

func TestRuteAnalisTerdaftar(t *testing.T) {
	isi, err := os.ReadFile("handlers.go")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(isi), `"PUT /api/klaim-life/{id}/analis", simpanAnalis(`) {
		t.Error("rute Save layar Claim Analis tidak terdaftar")
	}
}

func TestGalatAnalisDipetakan(t *testing.T) {
	for _, u := range []struct {
		nama string
		err  error
		mau  int
	}{
		{"identitas", inti.ErrTanpaIdentitas, http.StatusUnauthorized},
		{"bukan Claim Analyst", inti.ErrTanpaWewenang, http.StatusForbidden},
		{"tertutup", kontrak.ErrKasusSudahTertutup, http.StatusConflict},
		{"bukan Claim Analis", services.ErrAnalisBukanClaimAnalis, http.StatusConflict},
		{"peserta tidak aktif", services.ErrCatatanAnalisGagalTulis, http.StatusConflict},
		{"kepanjangan", services.ErrCatatanAnalisKepanjangan, http.StatusUnprocessableEntity},
		{"klaim tidak ada", services.ErrKlaimTidakAda, http.StatusNotFound},
		{"permintaan", galat.ErrPermintaanTidakSah, http.StatusBadRequest},
		{"lain", errors.New("x"), http.StatusInternalServerError},
	} {
		w := httptest.NewRecorder()
		jawabGalatAnalis(w, u.err)
		if w.Code != u.mau {
			t.Errorf("%s: kode %d, mau %d", u.nama, w.Code, u.mau)
		}
	}
}
