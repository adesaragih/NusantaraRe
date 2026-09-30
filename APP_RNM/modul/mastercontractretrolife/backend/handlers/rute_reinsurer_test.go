package handlers_test

// Rute reinsurer dan laporan total share (paket 4).

import (
	"net/http"
	"strings"
	"testing"

	"nusantarare/modul/mastercontractretrolife/backend/handlers"
	"nusantarare/modul/mastercontractretrolife/backend/models"
)

func serverReinsurer(t *testing.T) *uji {
	u := server(t, true)
	u.g.MasterRe = []models.MasterReinsurer{{ID: "UJI-L01", ClientName: "UJI REASURANSI SATU", StatusActive: models.StatusMasterReinsurerAktif}}
	return u
}

func TestPostReinsurerLaluTotalShareMencolok(t *testing.T) {
	u := serverReinsurer(t)
	kode, badan := u.kirim(t, "POST", handlers.Prefix+"/kontrak/UJI-K1/reinsurer",
		`{"reinsurerName":"x","reinsurerId":"UJI-L01","pctShare":"40","komisi":"25","ovrComm":"0"}`)
	if kode != http.StatusOK || !strings.Contains(badan, `"reinsurerName":"UJI REASURANSI SATU"`) {
		t.Fatalf("POST reinsurer: %d %s", kode, badan)
	}
	kode, badan = u.minta(t, "GET", handlers.Prefix+"/kontrak/UJI-K1/reinsurer", true)
	if kode != http.StatusOK || !strings.Contains(badan, `"totalShare":"40"`) || !strings.Contains(badan, `"totalBukan100":true`) {
		t.Errorf("GET reinsurer: %d %s", kode, badan)
	}
}

func TestReinsurerShareDiLuarRentang422(t *testing.T) {
	kode, badan := serverReinsurer(t).kirim(t, "POST", handlers.Prefix+"/kontrak/UJI-K1/reinsurer",
		`{"reinsurerName":"x","reinsurerId":"UJI-L01","pctShare":"140","komisi":"25","ovrComm":"0"}`)
	if kode != http.StatusUnprocessableEntity || !strings.Contains(badan, "(%) SHARE must be between 0 and 100") {
		t.Errorf("share 140: %d %s", kode, badan)
	}
}

func TestLaporanTotalShareRuteBaca(t *testing.T) {
	kode, badan := serverReinsurer(t).minta(t, "GET", handlers.Prefix+"/laporan/total-share-bukan-100?tahun=UJI-T1", true)
	if kode != http.StatusOK || !strings.Contains(badan, `"kontrakId":"UJI-K1"`) || !strings.Contains(badan, `"totalShare":"0"`) {
		t.Errorf("laporan: %d %s", kode, badan)
	}
}
