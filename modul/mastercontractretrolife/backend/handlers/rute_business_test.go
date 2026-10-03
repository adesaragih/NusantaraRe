package handlers_test

// Rute business dan `Copy to all Reinstype` (paket 6).

import (
	"net/http"
	"strings"
	"testing"

	"nusantarare/modul/mastercontractretrolife/backend/handlers"
	"nusantarare/modul/mastercontractretrolife/backend/models"
)

func serverBusiness(t *testing.T) *uji {
	u := server(t, true)
	u.g.MasterBiz = []models.MasterBusiness{{ID: "UJI-B01", Note: "UJI BUSINESS", OldID: "L01"}}
	u.g.Kontrak["UJI-K2"] = models.Kontrak{ID: "UJI-K2", IDTreatyYear: "UJI-T1", ReinsTypeID: "10197", ReinsTypeName: "2ND QS"}
	return u
}

func TestBusinessBaruLaluSalinSemuaDenganKonfirmasi(t *testing.T) {
	u := serverBusiness(t)
	kode, badan := u.kirim(t, "POST", handlers.Prefix+"/kontrak/UJI-K1/business",
		`{"bizCode":"UJI-B01","bizName":"x","riRateId":"UJI-RATE","riRate":"UJI R"}`)
	if kode != http.StatusOK || !strings.Contains(badan, `"bizName":"UJI BUSINESS"`) || !strings.Contains(badan, `"treatyYear":"2026"`) {
		t.Fatalf("POST business: %d %s", kode, badan)
	}
	kode, badan = u.minta(t, "GET", handlers.Prefix+"/business/1000044/salin-semua", true)
	if kode != http.StatusOK || !strings.Contains(badan, `"id":"UJI-K2"`) || strings.Contains(badan, `"sasaran":[{"id":"UJI-K1"`) {
		t.Fatalf("pratinjau: %d %s", kode, badan)
	}
	kode, _ = u.kirim(t, "POST", handlers.Prefix+"/business/1000044/salin-semua", `{"sasaran":[]}`)
	if kode != http.StatusConflict {
		t.Errorf("konfirmasi basi: %d, mau 409", kode)
	}
	kode, badan = u.kirim(t, "POST", handlers.Prefix+"/business/1000044/salin-semua", `{"sasaran":["UJI-K2"]}`)
	if kode != http.StatusOK || !strings.Contains(badan, `"pesan":"Copied to all reins types."`) || !strings.Contains(badan, `"jumlah":1`) {
		t.Errorf("salin-semua: %d %s", kode, badan)
	}
}
