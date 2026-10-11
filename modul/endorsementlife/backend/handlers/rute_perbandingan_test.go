package handlers_test

import (
	"net/http"
	"strings"
	"testing"

	"nusantarare/modul/endorsementlife/backend/handlers"
)

// TestRutePerbandingan - GET .../perbandingan (08-10-2026): gerbang baca yang sama.
func TestRutePerbandingan(t *testing.T) {
	g := gudangSatuKasus()
	g.Peserta[0].Nilai["CURRENCY"] = "IDR"
	g.Peserta[0].Nilai["GROSS_PREMIUM"] = "10"
	h := routerBuat(g)
	jalur := handlers.Prefix + "/kasus/EDML-1/perbandingan"
	if w := minta(t, h, "GET", jalur, "", false); w.Code != http.StatusUnauthorized {
		t.Fatalf("tanpa identitas %d", w.Code)
	}
	w := minta(t, h, "GET", jalur, "", true)
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), `"currency":"IDR"`) ||
		!strings.Contains(w.Body.String(), `"before":{"cacah":1,`) || !strings.Contains(w.Body.String(), `"PREMIUM":"10"`) {
		t.Fatalf("perbandingan %d %s", w.Code, w.Body)
	}
	if w := minta(t, h, "GET", handlers.Prefix+"/kasus/EDML-9/perbandingan", "", true); w.Code != http.StatusNotFound {
		t.Fatalf("kasus tak ada %d", w.Code)
	}
}
