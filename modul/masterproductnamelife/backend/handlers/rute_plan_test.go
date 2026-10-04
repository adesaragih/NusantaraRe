package handlers_test

import (
	"net/http"
	"strings"
	"testing"

	"nusantarare/modul/masterproductnamelife/backend/models"
)

func TestHTTPMasterPlanDanViewRate(t *testing.T) {
	u := server(t, true)
	u.g.Plan = []models.JenisPlan{{ID: "P1", CoverName: "UJI COVER", Business: "UJI BIZ", Benefit: "UJI MANFAAT"}}
	kode, badan := u.minta(t, "GET", pre+"/master-plan?cari=biz", "", true)
	if kode != http.StatusOK || !strings.Contains(badan, `"coverName":"UJI COVER"`) {
		t.Errorf("master plan: %d %s", kode, badan)
	}
	// K1 01-10-2026 (OQ-MPNL-03): `View Rate` yang dulu 503 kini 200 berisi view `RATE_LIFE`.
	u.g.Rate["R1"] = []models.BarisRate{{ID: "UJI-1", UsedBy: "UJI RATE", Gender: "U", Contract: "10", Age: "30", Rate: "0,5"}}
	kode, badan = u.minta(t, "GET", pre+"/rate?riRateId=R1", "", true)
	if kode != http.StatusOK || !strings.Contains(badan, `"rate":"0,5"`) || !strings.Contains(badan, `"terpotong":false`) ||
		!strings.Contains(badan, `"total":1`) {
		t.Errorf("View Rate: %d %s", kode, badan)
	}
	if kode, badan := u.minta(t, "GET", pre+"/rate?riRateId=", "", true); kode != http.StatusUnprocessableEntity {
		t.Errorf("View Rate tanpa RIRATEID: %d %s", kode, badan)
	}
}
