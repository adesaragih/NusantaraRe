package handlers_test

import (
	"net/http"
	"strings"
	"testing"

	"nusantarare/modul/endorsementlife/backend/handlers"
	"nusantarare/modul/endorsementlife/backend/models"
)

func TestRuteSimpan(t *testing.T) {
	g := gudangSatuKasus()
	g.Peserta[0].Nilai["CURRENCY"] = "IDR"
	g.Peserta[0].Nilai["GROSS_PREMIUM"] = "10"
	h := routerBuat(g)
	jalur := handlers.Prefix + "/kasus/EDMLF-1/simpan"
	if w := minta(t, h, "POST", jalur, `{"pilih":["UJI-D1"]}`, false); w.Code != http.StatusUnauthorized {
		t.Fatalf("tanpa identitas %d", w.Code)
	}
	if w := minta(t, h, "POST", jalur, `{"pilih":`, true); w.Code != http.StatusBadRequest {
		t.Fatalf("badan rusak %d", w.Code)
	}
	w := minta(t, h, "POST", jalur, `{"pilih":["UJI-D1"]}`, true)
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), `"ditandai":1,"status":"Delete"`) ||
		!strings.Contains(w.Body.String(), `"PREMIUM":"-10"`) {
		t.Fatalf("simpan %d %s", w.Code, w.Body)
	}
	if g.Peserta[0].EdmStatus != models.StatusDelete {
		t.Error("peserta tidak ditandai")
	}
	if w := minta(t, h, "POST", jalur, `{}`, true); w.Code != http.StatusConflict || !strings.Contains(w.Body.String(), `"galat"`) {
		t.Fatalf("simpan kedua %d %s", w.Code, w.Body)
	}
	g.Polis["EDMLF-1"].Status = models.StatusKasusSelesai
	if w := minta(t, h, "POST", jalur, `{}`, true); w.Code != http.StatusConflict {
		t.Fatalf("kasus tertutup %d %s", w.Code, w.Body)
	}
	if w := minta(t, h, "POST", handlers.Prefix+"/kasus/EDMLF-9/simpan", `{}`, true); w.Code != http.StatusNotFound {
		t.Fatalf("kasus tak ada %d", w.Code)
	}
	k := gudangSatuKasus()
	k.Peserta = nil
	if w := minta(t, routerBuat(k), "POST", jalur, `{}`, true); w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("tanpa peserta %d %s", w.Code, w.Body)
	}
	if w := minta(t, routerUji(gudangSatuKasus(), false), "POST", jalur, `{}`, true); w.Code != http.StatusServiceUnavailable {
		t.Fatalf("tanpa Oracle %d", w.Code)
	}
}
