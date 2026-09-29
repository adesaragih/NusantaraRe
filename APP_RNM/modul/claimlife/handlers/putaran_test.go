package handlers

// Uji pintu HTTP `Add` grid adjustment - putaran berikutnya.

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"nusantarare/inti"
	"nusantarare/inti/galat"
	"nusantarare/inti/kontrak"
	"nusantarare/modul/claimlife/services"
)

func TestGalatPutaranDipetakanKeKodeYangBenar(t *testing.T) {
	for _, u := range []struct {
		nama string
		err  error
		mau  int
	}{
		{"tanpa identitas", inti.ErrTanpaIdentitas, http.StatusUnauthorized},
		{"tanpa wewenang", inti.ErrTanpaWewenang, http.StatusForbidden},
		{"baris terakhir belum ditolak", services.ErrBukanPenolakan, http.StatusConflict},
		{"kasus tertutup", kontrak.ErrKasusSudahTertutup, http.StatusConflict},
		{"bukan Claim Analis", services.ErrTahapTanpaAddAdjustment, http.StatusConflict},
		{"tahap tidak dikenal", services.ErrTahapTidakDikenal, http.StatusConflict},
		{"permintaan tidak sah", galat.ErrPermintaanTidakSah, http.StatusBadRequest},
		// Dibungkus sebab-sebabnya - seperti layanan membungkusnya.
		{"dibungkus", fmt.Errorf("%w: peserta \"P-1\" tidak ada baris",
			services.ErrBukanPenolakan), http.StatusConflict},
	} {
		w := httptest.NewRecorder()
		if !jawabGalatPutaran(w, u.err) {
			t.Errorf("%s: galat tidak dijawab", u.nama)
			continue
		}
		if w.Code != u.mau {
			t.Errorf("%s: kode = %d, mau %d", u.nama, w.Code, u.mau)
		}
	}
	if jawabGalatPutaran(httptest.NewRecorder(), nil) {
		t.Error("galat nil dijawab; jalur berhasil tidak akan tercapai")
	}
}

// Tanpa Oracle: 503, bukan 500 - dan tidak ada yang ditulis.
func TestPutaranTanpaOracle503(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/api/klaim-life/CLM-1/peserta/P-1/putaran", nil)
	r.SetPathValue("id", "CLM-1")
	r.SetPathValue("pesertaId", "P-1")
	tambahPutaran(services.New(nil), true)(w, r)
	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("kode = %d, mau 503", w.Code)
	}
}
