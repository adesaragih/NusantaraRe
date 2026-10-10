package handlers_test

// Uji rute rumus tab Value Difference — `POST /hitung/selisih`.

import (
	"net/http"
	"strings"
	"testing"

	"nusantarare/modul/treatyin/backend/handlers"
	"nusantarare/modul/treatyin/backend/services"
)

func TestRuteHitungSelisih(t *testing.T) {
	h := handlers.RouterDengan(services.LayananDengan(gudangTiruan{}), true, true)
	const jalur = "/api/treaty-in/hitung/selisih"
	isi := `{"aksi":"edm","akar":{"medan":{"RNMShare":"50"},"larik":{}},"lama":{"medan":{"RNMShare":"40"},"larik":{}},"actual":{"medan":{},"larik":{}}}`
	if w := kirim(t, h, http.MethodPost, jalur, "", badan(isi)); w.Code != http.StatusUnauthorized {
		t.Errorf("tanpa identitas: %d", w.Code)
	}
	if w := kirim(t, h, http.MethodPost, jalur, "AKUN-UJI", badan(`{`)); w.Code != http.StatusBadRequest {
		t.Errorf("badan rusak: %d", w.Code)
	}
	if w := kirim(t, h, http.MethodPost, jalur, "AKUN-UJI", badan(`{"aksi":"hapus"}`)); w.Code != http.StatusUnprocessableEntity {
		t.Errorf("aksi tak dikenal: %d", w.Code)
	}
	w := kirim(t, h, http.MethodPost, jalur, "AKUN-UJI", badan(isi))
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), `"RNMShare":"10"`) || !strings.Contains(w.Body.String(), `"sebelumProrata":null`) {
		t.Errorf("%d %s", w.Code, w.Body.String())
	}
}
