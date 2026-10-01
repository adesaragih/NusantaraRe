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
	kode, badan = u.minta(t, "GET", pre+"/rate?riRateId=R1", "", true)
	if kode != http.StatusServiceUnavailable || !strings.Contains(badan, "OQ-MPNL-03") {
		t.Errorf("View Rate menunggu OQ-MPNL-03: %d %s", kode, badan)
	}
}
