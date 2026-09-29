package handlers

// Uji pintu HTTP keputusan penawaran - GILIRAN-14 butir bq.

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"nusantarare/modul/premiumlist/models"
)

// TestRutePenggolongManualTidakAdaLagi - Decision3 tidak ditanyakan.
//
// ⛔ `Offer`/`Premium` hasil decision table `IsFlagOnGoingPolicy` atas bendera
// kasus; rute yang menerimanya dari pengguna membuka jalur yang Pega tidak
// punya.
func TestRutePenggolongManualTidakAdaLagi(t *testing.T) {
	// Refactor bentuk B (30-09-2026): tabel rute modul ini kini
	// `rute_premiumlist.go` (DaftarkanRute), bukan Router bersama.
	isi, err := os.ReadFile("rute_premiumlist.go")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(isi), "/penggolong") {
		t.Error("rute penggolong manual masih terdaftar di Router")
	}
	if !strings.Contains(string(isi), `"POST /api/polis-life/{id}/keputusan"`) {
		t.Error("rute keputusan hilang")
	}
}

// Bendera di luar tabel: keputusannya sah, jalurnya yang tidak ada - 409.
func TestBenderaTanpaKonektorDijawab409(t *testing.T) {
	w := httptest.NewRecorder()
	if !jawabGalatPolis(w, fmt.Errorf("polis %q: %w", "UJI-PL-A",
		fmt.Errorf("%w: bendera %q", models.ErrBenderaTanpaKonektor, ""))) {
		t.Fatal("galat tidak dijawab")
	}
	if w.Code != http.StatusConflict {
		t.Errorf("kode = %d, mau 409", w.Code)
	}
}

// Jawaban akibat tidak lagi membawa `menungguPenggolong`: sesudah bq tidak ada
// keadaan "menunggu" yang dapat dilihat layar.
func TestJawabanAkibatTanpaPenandaMenunggu(t *testing.T) {
	w := httptest.NewRecorder()
	tulisAkibat(w, models.AkibatKeputusan{TahapTujuan: models.TahapPolisDetail})
	var isi map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &isi); err != nil {
		t.Fatal(err)
	}
	if _, ada := isi["menungguPenggolong"]; ada {
		t.Errorf("jawaban masih membawa menungguPenggolong: %s", w.Body.String())
	}
	if isi["tahapTujuan"] != models.TahapPolisDetail {
		t.Errorf("tahapTujuan = %v", isi["tahapTujuan"])
	}
}
