//go:build db

package handlers_test

// Seam HTTP reinsurer dan laporan total share terhadap Oracle NYATA (paket 4).

import (
	"net/http"
	"strings"
	"testing"
)

func TestDBReinsurerSalinanDanLaporan(t *testing.T) {
	u := pasangDB(t)
	s := u.skema
	u.exec(t, `INSERT INTO `+s+`.TREATYYEAR_LIFE (ID, TREATYYEAR, STARTDATE, ENDDATE)
		VALUES ('1000001', '2026', DATE '2026-01-01', DATE '2026-12-31')`)
	u.exec(t, `INSERT INTO `+s+`.TREATYCONTRACT_LIFE (ID, IDTREATYYEAR, REINSTYPEID, REINSTYPENAME)
		VALUES ('1000002', '1000001', '10196', 'QS')`)
	u.exec(t, `INSERT INTO `+s+`.TREATYCONTRACT_LIFE (ID, IDTREATYYEAR, REINSTYPEID, REINSTYPENAME)
		VALUES ('1000003', '1000001', '10197', '2ND QS')`)
	u.exec(t, `INSERT INTO `+s+`.AGENT (ID, CLIENTNAME, STATUSACTIVE) VALUES ('UJI-L01', 'UJI REASURANSI', '1')`)
	for _, share := range []string{"33.3333333333333333333", "66.6666666666666666667"} {
		kode, badan := u.kirim(t, "POST", "/api/master-contract-retro-life/kontrak/1000002/reinsurer",
			`{"reinsurerName":"x","reinsurerId":"UJI-L01","pctShare":"`+share+`","komisi":"10","ovrComm":"1"}`)
		if kode != http.StatusOK {
			t.Fatalf("POST reinsurer: %d %s", kode, badan)
		}
	}
	if n := u.cacah(t, "TREATYREINSURER_LIFE", "TREATYYEARID = '1000001' AND REINSTYPEID = '10196' "+
		"AND REINSTYPENAME = 'QS' AND REINSURERNAME = 'UJI REASURANSI' AND USERID = 'UJI-PELAKU'"); n != 2 {
		t.Errorf("salinan kontrak / nama master: %d baris", n)
	}
	kode, badan := u.get(t, "/api/master-contract-retro-life/laporan/total-share-bukan-100?tahun=1000001")
	// Tepat 100 (desimal penuh Oracle) tidak muncul; kontrak tanpa reinsurer muncul bertotal 0.
	if kode != http.StatusOK || strings.Contains(badan, `"kontrakId":"1000002"`) ||
		!strings.Contains(badan, `"kontrakId":"1000003"`) || !strings.Contains(badan, `"selisih":"100"`) {
		t.Errorf("laporan: %d %s", kode, badan)
	}
}
