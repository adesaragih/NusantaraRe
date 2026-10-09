package handlers_test

// Uji rute rumus tab Installment — `POST /hitung/angsuran`.

import (
	"net/http"
	"strings"
	"testing"

	"nusantarare/modul/treatyin/backend/handlers"
	"nusantarare/modul/treatyin/backend/services"
)

func TestRuteHitungAngsuran(t *testing.T) {
	h := handlers.RouterDengan(services.LayananDengan(gudangTiruan{}), true, true)
	const jalur = "/api/treaty-in/hitung/angsuran"
	isi := `{"aksi":"nilai","installmentNo":"4","netPremium":[{"Currency":"IDR","Value":"644674819.59"}]}`
	if w := kirim(t, h, http.MethodPost, jalur, "", badan(isi)); w.Code != http.StatusUnauthorized {
		t.Errorf("tanpa identitas: %d", w.Code)
	}
	if w := kirim(t, h, http.MethodPost, jalur, "AKUN-UJI", badan(`{`)); w.Code != http.StatusBadRequest {
		t.Errorf("badan rusak: %d", w.Code)
	}
	if w := kirim(t, h, http.MethodPost, jalur, "AKUN-UJI", badan(`{"aksi":"hapus"}`)); w.Code != http.StatusUnprocessableEntity {
		// `ErrMasukanTidakSah` = 422 di modul ini (JSON sah, isinya ditolak).
		t.Errorf("aksi tak dikenal: %d", w.Code)
	}
	w := kirim(t, h, http.MethodPost, jalur, "AKUN-UJI", badan(isi))
	for _, mau := range []string{`"Amount":"161168704.8975"`, `"TotalInstallmentNP":[{"Currency":"IDR","CurrencyID":"","Value":"644674819.59"}]`} {
		if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), mau) {
			t.Errorf("%d %s — mau %s", w.Code, w.Body.String(), mau)
		}
	}
}
