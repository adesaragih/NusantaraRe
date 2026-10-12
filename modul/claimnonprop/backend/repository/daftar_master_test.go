package repository

import (
	"strings"
	"testing"

	"nusantarare/modul/claimnonprop/backend/models"
)

var tabelPopupUji = tabelPopupMaster{detail: "S.TREATYINDETAIL", edm: "S.TREATYINDETAILEDM", master: "S.TREATY_IN",
	masterEDM: "S.TREATY_IN_EDM", tahun: "S.TREATYYEAR", bisnis: "S.TREATYBUSINESS"}

// TestDaftarMasterCoBDariTreatyBusiness - Class of Business popup dari TREATYBUSINESS (CLASSOFBUSINESS TREATYINDETAIL
// kosong untuk semua master Non Prop, laporan pemakai 11-10-2026); sumber IN = detail + EDM, INEDM = EDM saja.
func TestDaftarMasterCoBDariTreatyBusiness(t *testing.T) {
	in := sqlDaftarMaster(tabelPopupUji, models.SaringanMaster{ClassOfBusiness: "motor"})
	if len(in) != 2 || !strings.Contains(in[0].q, "FROM S.TREATYINDETAIL WHERE") ||
		!strings.Contains(in[0].q, "JOIN S.TREATY_IN d") || !strings.Contains(in[1].q, "JOIN S.TREATY_IN_EDM d") {
		t.Fatalf("sumber IN harus detail+TREATY_IN lalu EDM+TREATY_IN_EDM: %d kueri", len(in))
	}
	for _, k := range in {
		for _, harus := range []string{"JOIN S.TREATYBUSINESS b", "JOIN S.TREATYYEAR c", "NVL(k.BIZNAME, a.CLASSOFBUSINESS) CLASSOFBUSINESS",
			"UPPER(CLASSOFBUSINESS) LIKE :3"} {
			if !strings.Contains(k.q, harus) {
				t.Errorf("kueri tanpa %q", harus)
			}
		}
		if len(k.args) != 3 || k.args[0] != models.ProporsiMaster || k.args[2] != "%MOTOR%" {
			t.Errorf("argumen = %v", k.args)
		}
	}
	edm := sqlDaftarMaster(tabelPopupUji, models.SaringanMaster{Sumber: models.SumberMasterINEDM})
	if len(edm) != 1 || strings.Contains(edm[0].q, "S.TREATYINDETAIL ") || !strings.Contains(edm[0].q, "S.TREATYINDETAILEDM") {
		t.Fatalf("sumber INEDM harus EDM saja")
	}
}

// TestUrutkanBarisMasterSepertiOracle - urutan gabungan = `TREATYYEAR DESC, TREATYID DESC, CLASSOFBUSINESS` Oracle.
func TestUrutkanBarisMasterSepertiOracle(t *testing.T) {
	b := []models.BarisMaster{
		{TreatyYear: "2024", TreatyID: "1002320", ClassOfBusiness: "MOTORCYCLE"},
		{TreatyYear: "2025", TreatyID: "1002318", ClassOfBusiness: "MOTOR VEHICLE"},
		{TreatyYear: "2024", TreatyID: "1002320", ClassOfBusiness: ""},
		{TreatyYear: "2024", TreatyID: "1002320", ClassOfBusiness: "MOTOR VEHICLE"},
		{TreatyYear: "2024", TreatyID: "1001789", ClassOfBusiness: "FIRE"},
	}
	urutkanBarisMaster(b)
	var got []string
	for _, r := range b {
		got = append(got, r.TreatyYear+"/"+r.TreatyID+"/"+r.ClassOfBusiness)
	}
	harap := "2025/1002318/MOTOR VEHICLE|2024/1002320/MOTOR VEHICLE|2024/1002320/MOTORCYCLE|2024/1002320/|2024/1001789/FIRE"
	if strings.Join(got, "|") != harap {
		t.Errorf("urutan = %s\nharap   %s", strings.Join(got, "|"), harap)
	}
}
