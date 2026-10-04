package handlers_test

// Rute security dan eksposur (paket 5).

import (
	"net/http"
	"strings"
	"testing"

	"nusantarare/modul/mastercontractretrolife/backend/handlers"
	"nusantarare/modul/mastercontractretrolife/backend/models"
)

func TestPostSecurityLaluEksposurDikunciID(t *testing.T) {
	u := serverReinsurer(t)
	u.g.Reinsurer["UJI-R1"] = models.Reinsurer{ID: "UJI-R1", TreatyYearID: "UJI-T1", TreatyContractID: "UJI-K1",
		PctShare: desimal(t, "40")}
	kode, badan := u.kirim(t, "POST", handlers.Prefix+"/reinsurer/UJI-R1/security",
		`{"reinsurerName":"x","reinsurerId":"UJI-L01","pctShare":"10"}`)
	if kode != http.StatusOK || !strings.Contains(badan, `"treatyReinsurerId":"UJI-R1"`) {
		t.Fatalf("POST security: %d %s", kode, badan)
	}
	kode, badan = u.minta(t, "GET", handlers.Prefix+"/reinsurer/UJI-R1/security", true)
	if kode != http.StatusOK || !strings.Contains(badan, `"eksposur":{"1000044":"4"}`) ||
		!strings.Contains(badan, `"pctShare":"10"`) {
		t.Errorf("GET security: %d %s", kode, badan)
	}
}
