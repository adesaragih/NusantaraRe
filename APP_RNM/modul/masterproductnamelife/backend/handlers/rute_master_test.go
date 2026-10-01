package handlers_test

// Seam HTTP ketujuh pemilih master (paket 2, tiket 04) di atas gudang tiruan.

import (
	"net/http"
	"strings"
	"testing"

	"nusantarare/modul/masterproductnamelife/backend/models"
	"nusantarare/modul/masterproductnamelife/backend/tiruan"
)

func isiMaster(g *tiruan.Gudang) {
	g.Master[models.MasterCeding] = []models.NilaiMaster{{ID: "L0UJI1", Nama: "UJI CEDING SATU"}, {ID: "L0UJI2", Nama: "UJI LAIN"}}
	g.Master[models.MasterSOB] = []models.NilaiMaster{{ID: "L0SOB", Nama: "UJI SOB"}}
	g.Master[models.MasterPemegangPolis] = []models.NilaiMaster{{ID: "UJI-ORG-1", Nama: "UJI PEMEGANG"}}
	g.Master[models.MasterMataUang] = []models.NilaiMaster{{ID: "1", Nama: "IDR"}}
	g.Master[models.MasterRIRisk] = []models.NilaiMaster{{ID: "1000117", Nama: "UJI RISK"}}
	g.Master[models.MasterPenyebab] = []models.NilaiMaster{{ID: "100004", Nama: "ANY CAUSE"}}
}

func TestHTTPPemilihMasterMencariHurufBesar(t *testing.T) {
	u := server(t, true)
	isiMaster(u.g)
	kode, badan := u.minta(t, "GET", pre+"/master/ceding?cari=satu", "", true)
	if kode != http.StatusOK || !strings.Contains(badan, `"id":"L0UJI1"`) || strings.Contains(badan, "L0UJI2") ||
		!strings.Contains(badan, `"total":1`) {
		t.Errorf("cari ceding (SearchPolicyHolder_act b236 huruf besar): %d %s", kode, badan)
	}
	if u.g.CariTerakhir != "SATU" {
		t.Errorf("kata cari harus dikirim huruf besar ke gudang: %q", u.g.CariTerakhir)
	}
	for _, j := range []string{"sob", "pemegang-polis", "mata-uang", "ri-risk", "penyebab"} {
		if kode, badan := u.minta(t, "GET", pre+"/master/"+j, "", true); kode != http.StatusOK {
			t.Errorf("%s: %d %s", j, kode, badan)
		}
	}
}

func TestHTTPRIRateMenungguPersetujuan(t *testing.T) {
	u := server(t, true)
	kode, badan := u.minta(t, "GET", pre+"/master/ri-rate?cari=x", "", true)
	if kode != http.StatusServiceUnavailable || !strings.Contains(badan, "OQ-MPNL-03") {
		t.Errorf("R/I Rate harus 503 berkalimat: %d %s", kode, badan)
	}
}

func TestHTTPMasterGalatBerkalimat(t *testing.T) {
	u := server(t, true)
	if kode, _ := u.minta(t, "GET", pre+"/master/tidak-dikenal", "", true); kode != http.StatusNotFound {
		t.Errorf("jenis tak dikenal: %d", kode)
	}
	u.g.GagalMaster = tiruan.GalatMasterUji("AGENT")
	kode, badan := u.minta(t, "GET", pre+"/master/ceding", "", true)
	if kode != http.StatusServiceUnavailable || !strings.Contains(badan, "AGENT") || strings.Contains(badan, "ORA-") {
		t.Errorf("master tak terbaca: 503 yang menyebut objek, tanpa sebab Oracle: %d %s", kode, badan)
	}
	if kode, _ := u.minta(t, "GET", pre+"/master/ceding", "", false); kode != http.StatusUnauthorized {
		t.Errorf("tanpa identitas: %d", kode)
	}
}
