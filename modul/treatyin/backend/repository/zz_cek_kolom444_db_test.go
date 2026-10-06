//go:build db

package repository_test

// SEMENTARA — apakah keenam belas kolom migrasi 444 sungguh terisi.

import (
	"testing"

	"nusantarare/inti/backend/config"
	"nusantarare/inti/backend/db"
)

func TestCekKolom444Terisi(t *testing.T) {
	cfg, _ := config.Load()
	if !cfg.PunyaOracle() {
		t.Skip("ORACLE_DSN kosong")
	}
	d, err := db.Open(cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = d.Close() }()
	nama, _ := d.Qualify("T_TREATY_REVISION")
	for _, k := range []string{
		"REPORTINGSTART", "REPORTINGEND", "REPORTINGPERIOD", "REPORTINGINTERVAL",
		"REPORTINGSUBMISSION", "REPORTINGCONFIRMATION", "REPORTINGSETTLEMENT",
		"CEDINGID", "LEADINGREINSSOURCEID", "LEADINGREINSID", "INFORMATION",
		"POSITION", "POSITIONUSERNAME", "CHOOSESTATUSAKSEPTASI",
		"ACCUMULATIONPERIOD", "COMMENTTEKS",
	} {
		var n int
		if err := d.QueryRowContext(t.Context(),
			"SELECT COUNT(*) FROM "+nama+" WHERE "+k+" IS NOT NULL AND "+k+" <> ' '").Scan(&n); err != nil {
			t.Fatalf("%s: %v", k, err)
		}
		t.Logf("  %-24s terisi di %d baris", k, n)
	}
	var total int
	_ = d.QueryRowContext(t.Context(), "SELECT COUNT(*) FROM "+nama).Scan(&total)
	t.Logf("baris T_TREATY_REVISION: %d", total)
}
