//go:build db

package handlers_test

// Seam HTTP business dan `Copy to all Reinstype` terhadap Oracle NYATA (paket 6).

import (
	"net/http"
	"strings"
	"testing"
)

func TestDBSalinSemuaKeJenisBerbedaSatuTransaksi(t *testing.T) {
	u := pasangDB(t)
	s := u.skema
	u.exec(t, `INSERT INTO `+s+`.TREATYYEAR_LIFE (ID, TREATYYEAR, STARTDATE, ENDDATE)
		VALUES ('1000001', '2026', DATE '2026-01-01', DATE '2026-12-31')`)
	for _, k := range [][2]string{{"1000002", "10196"}, {"1000003", "10197"}, {"1000004", "10196"}} {
		u.exec(t, `INSERT INTO `+s+`.TREATYCONTRACT_LIFE (ID, IDTREATYYEAR, REINSTYPEID, REINSTYPENAME)
			VALUES (:1, '1000001', :2, 'UJI')`, k[0], k[1])
	}
	u.exec(t, `INSERT INTO `+s+`.BUSINESS (ID, NOTE, OLDID) VALUES ('UJI-B01', 'UJI BUSINESS', 'L01')`)
	u.exec(t, `INSERT INTO `+s+`.M_RATE_LIFE_SUMMARY (ID, USEDBY) VALUES ('UJI-RATE', 'UJI R')`)
	kode, badan := u.kirim(t, "POST", "/api/master-contract-retro-life/kontrak/1000002/business",
		`{"bizCode":"UJI-B01","bizName":"x","riRateId":"UJI-RATE","riRate":" UJI R, 0,5% "}`)
	if kode != http.StatusOK || !strings.Contains(badan, `"id":"1000044"`) {
		t.Fatalf("POST business: %d %s", kode, badan)
	}
	if got := u.teks(t, `SELECT RIRATE FROM {s}.TREATYBUSINESS_LIFE WHERE ID = '1000044'`); got != " UJI R, 0,5% " {
		t.Errorf("RIRATE tidak apa adanya: %q", got)
	}
	kode, badan = u.get(t, "/api/master-contract-retro-life/business/1000044/salin-semua")
	if kode != http.StatusOK || !strings.Contains(badan, `"id":"1000003"`) || strings.Contains(badan, `"id":"1000004"`) {
		t.Fatalf("pratinjau (hanya jenis berbeda): %d %s", kode, badan)
	}
	kode, badan = u.kirim(t, "POST", "/api/master-contract-retro-life/business/1000044/salin-semua", `{"sasaran":["1000003"]}`)
	if kode != http.StatusOK {
		t.Fatalf("salin-semua: %d %s", kode, badan)
	}
	if n := u.cacah(t, "TREATYBUSINESS_LIFE", "TREATYCONTRACTID = '1000003' AND REINSTYPEID = '10197' "+
		"AND TREATYYEAR = '2026' AND BIZCODE = 'UJI-B01' AND TGLUPDATE IS NOT NULL"); n != 1 {
		t.Errorf("baris salinan: %d", n)
	}
}
