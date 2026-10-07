package repository

import (
	"os"
	"regexp"
	"strings"
	"testing"

	"nusantarare/modul/nbfacin/backend/models"
)

// TestSQLFEA - tiket 41: FEA dibaca lewat case (JOIN lokasi) urut lokasi lalu SEQ_NO; dihapus SEBELUM lokasi.
func TestSQLFEA(t *testing.T) {
	baca := sqlBacaFEA(tabelObjekUji())
	for _, harus := range []string{"FROM UJI.F f", "JOIN UJI.L l ON l.ID = f.PARENT_ID", "WHERE l.PARENT_ID = :1",
		"ORDER BY f.PARENT_ID, f.SEQ_NO"} {
		if !strings.Contains(baca, harus) {
			t.Errorf("baca FEA tanpa %q", harus)
		}
	}
	hapus := sqlHapusObjek(tabelObjekUji())
	if n := len(hapus); hapus[n-2] != "DELETE FROM UJI.F WHERE PARENT_ID IN (SELECT ID FROM UJI.L WHERE PARENT_ID = :1)" ||
		!strings.HasPrefix(hapus[n-1], "DELETE FROM UJI.L ") {
		t.Errorf("FEA harus dihapus tepat sebelum lokasi: %q", hapus[n-2:])
	}
}

// TestArgFEAMenurutKolom - bind :5..:13 dipasangkan lewat KUNCI (nilai uji = nama kolomnya sendiri); kosong -> NULL.
func TestArgFEAMenurutKolom(t *testing.T) {
	f := models.BarisFEA{APAR: "APAR", Sprinkler: "SPRINKLER", SmokeDetector: "SMOKE_DETECTOR", Hydrant: "HYDRANT",
		PrivateTruckBrigade: "PRIVATE_TRUCK_BRIGADE", PrivateFireBrigade: "PRIVATE_FIRE_BRIGADE", TeamSOPSafety: "TEAM_SOP_SAFETY",
		TeamSOPRiskManagement: "TEAM_SOP_RISK_MANAGEMENT", Info: "INFO_FEA"}
	m := regexp.MustCompile(`\(([^)]*)\) VALUES`).FindStringSubmatch(sqlSisipFEA("UJI.F"))
	kolom := strings.Split(m[1], ", ")[4:]
	arg := argFEA(f)
	if len(arg) != len(kolom) || len(kolom) != 9 {
		t.Fatalf("%d bind, %d kolom, mau 9", len(arg), len(kolom))
	}
	for i, k := range kolom {
		if arg[i] != k {
			t.Errorf("bind kolom %s = %v", k, arg[i])
		}
	}
	for i, a := range argFEA(models.BarisFEA{}) {
		if a != nil {
			t.Errorf("kosong[%d] = %v, mau NULL", i, a)
		}
	}
}

// TestBacaFEAMemakaiKolomBernama - setiap kunci v("...") di fea.go ada di kolomBacaFEA dan sebaliknya.
func TestBacaFEAMemakaiKolomBernama(t *testing.T) {
	b, err := os.ReadFile("fea.go")
	if err != nil {
		t.Fatal(err)
	}
	ada, dibaca := map[string]bool{}, map[string]bool{}
	for _, k := range kolomBacaFEA {
		ada[k] = true
	}
	for _, m := range regexp.MustCompile(`v\("([^"]+)"\)`).FindAllStringSubmatch(string(b), -1) {
		dibaca[m[1]] = true
		if !ada[m[1]] {
			t.Errorf("kunci %q dibaca tetapi tidak dipilih", m[1])
		}
	}
	for _, k := range kolomBacaFEA {
		if !dibaca[k] {
			t.Errorf("kolom %s dipilih tetapi tidak dibaca", k)
		}
	}
}
