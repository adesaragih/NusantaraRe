package handlers_test

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/jejak"
	"nusantarare/modul/endorsementlife/backend/handlers"
	"nusantarare/modul/endorsementlife/backend/services"
	"nusantarare/modul/endorsementlife/backend/tiruan"
)

type jejakDiam struct{}

func (jejakDiam) Rekam(context.Context, *db.Tx, jejak.CatatanJejak) error { return nil }

func routerBuat(g *tiruan.Gudang) http.Handler {
	l := services.BaruLayanan(g, tiruan.Transaksi, nil).DenganJejak(jejakDiam{})
	return handlers.RouterDengan(l, true, true)
}

func gudangNB() *tiruan.Gudang {
	g := tiruan.Baru()
	g.Polis["NBLF-1"] = &tiruan.Polis{ID: "NBLF-1", ProdKe: 1, Kepala: map[string]string{"TYPE": "QP"}}
	g.Peserta = append(g.Peserta, &tiruan.Peserta{ID: "UJI-D1", PolisID: "NBLF-1", PLNumber: "UJI-PL-1",
		Nilai: map[string]string{"CERTIFICATE_NO": "UJI-C1"}})
	return g
}

func TestRuteKelayakanDanBuatKasus(t *testing.T) {
	h := routerBuat(gudangNB())
	w := minta(t, h, "POST", handlers.Prefix+"/kelayakan", `{"policyNo":"","edmType":"1"}`, true)
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), `"pesan":["Invalid Policy No !","Policy no Not Found !"]`) ||
		!strings.Contains(w.Body.String(), `"boleh":false`) {
		t.Fatalf("kelayakan kosong %d %s", w.Code, w.Body)
	}
	w = minta(t, h, "POST", handlers.Prefix+"/kelayakan", `{"policyNo":"UJI-PL-1","edmType":"1"}`, true)
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), `"pesan":[],"boleh":true`) {
		t.Fatalf("kelayakan sah %d %s", w.Code, w.Body)
	}
	w = minta(t, h, "POST", handlers.Prefix+"/kasus", `{"policyNo":"UJI-PL-1","edmType":"1","edmDate":"2026-10-01","description":"UJI"}`, true)
	if w.Code != http.StatusCreated || w.Header().Get("Location") != handlers.Prefix+"/kasus/EDMLF-1" || !strings.Contains(w.Body.String(), `"peserta":1`) {
		t.Fatalf("buat %d %s %v", w.Code, w.Body, w.Header())
	}
	// Kasus kedua atas polis yang sama: 422 berpesan gerbang 3 VERBATIM.
	w = minta(t, h, "POST", handlers.Prefix+"/kasus", `{"policyNo":"UJI-PL-1","edmType":"1"}`, true)
	if w.Code != http.StatusUnprocessableEntity || !strings.Contains(w.Body.String(), "There's EDM with this policy no that haven't finish yet!") {
		t.Fatalf("buat kedua %d %s", w.Code, w.Body)
	}
	w = minta(t, h, "GET", handlers.Prefix+"/kasus/EDMLF-1/polis-lama", "", true)
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), `"sumber":"aplikasi"`) || !strings.Contains(w.Body.String(), `"type":"QP"`) {
		t.Fatalf("polis lama %d %s", w.Code, w.Body)
	}
	w = minta(t, h, "GET", handlers.Prefix+"/kasus/EDMLF-1/peserta/UJI-TIDAK-ADA", "", true)
	if w.Code != http.StatusNotFound {
		t.Fatalf("rincian tak ada %d %s", w.Code, w.Body)
	}
}

func TestRuteTulisMenolakBadanRusakDanTanpaIdentitas(t *testing.T) {
	h := routerBuat(gudangNB())
	for _, k := range []struct {
		jalur, badan string
		ident        bool
		kode         int
	}{
		{"/kasus", `{"policyNo":"UJI-PL-1","edmType":"1"}`, false, http.StatusUnauthorized},
		{"/kasus", `{bukan json`, true, http.StatusBadRequest},
		{"/kasus", `{"policyNo":"UJI-PL-1","edmType":"2"}`, true, http.StatusUnprocessableEntity},
		{"/kelayakan", `{"policyNo":"` + strings.Repeat("X", 70<<10) + `"}`, true, http.StatusRequestEntityTooLarge},
	} {
		if w := minta(t, h, "POST", handlers.Prefix+k.jalur, k.badan, k.ident); w.Code != k.kode || !strings.Contains(w.Body.String(), `"galat"`) {
			t.Errorf("%s %.30s: %d %s, mau %d", k.jalur, k.badan, w.Code, w.Body, k.kode)
		}
	}
}
