package handlers

// Popup spreading SATU peserta (keputusan work owner 07-10-2026).

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/premiumlistlife/backend/services"
)

// Rute semua-peserta dihapus; rute per peserta dan rekap tetap.
func TestRuteSpreadingPerPeserta(t *testing.T) {
	isi, err := os.ReadFile("rute_premiumlist.go")
	if err != nil {
		t.Fatal(err)
	}
	s := string(isi)
	if strings.Contains(s, `"GET /api/polis-life/{id}/spreading"`) {
		t.Error("rute spreading semua peserta masih terdaftar")
	}
	for _, k := range []string{`"GET /api/polis-life/{id}/peserta/{detailId}/spreading"`, `"GET /api/polis-life/{id}/spreading/rekap"`} {
		if !strings.Contains(s, k) {
			t.Errorf("rute %s hilang", k)
		}
	}
	w := httptest.NewRecorder()
	Router(services.New(nil), true).ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/polis-life/NBLF-1/spreading", nil))
	if w.Code == http.StatusServiceUnavailable {
		t.Error("GET .../spreading masih dilayani")
	}
}

// Peserta milik polis lain → 404, bukan 500.
func TestPesertaBukanMilikPolisDijawab404(t *testing.T) {
	w := httptest.NewRecorder()
	if !jawabGalatPolis(w, fmt.Errorf("%w: UJI-D9", services.ErrPesertaTakDitemukan)) {
		t.Fatal("galat tidak dijawab")
	}
	if w.Code != http.StatusNotFound {
		t.Errorf("kode = %d, mau 404", w.Code)
	}
}

// Gerbang sama dengan grid peserta: identitas, lalu database.
func TestGerbangSpreadingPeserta(t *testing.T) {
	d := services.New(nil).DetailPolis()
	ctx := context.Background()
	if _, err := d.SpreadingPeserta(ctx, inti.Pelaku{}, "NBLF-1", "UJI-D1"); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Errorf("anonim: %v", err)
	}
	if _, err := d.SpreadingPeserta(ctx, inti.Pelaku{AkunID: "UJI-USER"}, "NBLF-1", "UJI-D1"); !errors.Is(err, db.ErrTanpaOracle) {
		t.Errorf("tanpa Oracle: %v", err)
	}
	if _, err := d.Peserta(ctx, inti.Pelaku{}, "NBLF-1", 1, 20); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Errorf("grid peserta anonim: %v", err)
	}
}
