package services

import (
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"nusantarare/modul/nbfacin/backend/models"
)

// TestPeriksaFEA - tiket 41: lima jumlah unit kosong atau bulat >= 0 (M-2), tiga dropdown tanpa enumerasi, lebar
// kolom; pesan ber-indeks baris dan FEA.
func TestPeriksaFEA(t *testing.T) {
	if m := periksaFEA(0, []models.BarisFEA{{APAR: "2", Sprinkler: "0", PrivateFireBrigade: "UJI", Info: "UJI"}, {}}); len(m) != 0 {
		t.Fatalf("sah ditolak: %v", m)
	}
	for nama, u := range map[string]struct {
		f     models.BarisFEA
		pesan string
	}{
		"apar minus":      {models.BarisFEA{APAR: "-1"}, "baris[4].fea[1].apar harus bilangan bulat >= 0"},
		"hydrant desimal": {models.BarisFEA{Hydrant: "1.5"}, "baris[4].fea[1].hydrant"},
		"truk teks":       {models.BarisFEA{PrivateTruckBrigade: "dua"}, "baris[4].fea[1].privateTruckBrigade"},
		"info 501":        {models.BarisFEA{Info: strings.Repeat("U", 501)}, "baris[4].fea[1].info paling banyak 500"},
		"dropdown 51":     {models.BarisFEA{TeamSOPSafety: strings.Repeat("U", 51)}, "baris[4].fea[1].teamSopSafety paling banyak 50"},
	} {
		if m := strings.Join(periksaFEA(4, []models.BarisFEA{{}, u.f}), "; "); !strings.Contains(m, u.pesan) {
			t.Errorf("%s: %q, mau %q", nama, m, u.pesan)
		}
	}
}

// TestLebarFEASamaDenganMigrasi - lebar validasi = lebar kolom migrasi 190 (dipasangkan lewat NAMA kolom).
func TestLebarFEASamaDenganMigrasi(t *testing.T) {
	b, err := os.ReadFile("../migrations/190_t_fealist.sql")
	if err != nil {
		t.Fatal(err)
	}
	lebar := map[string]int{}
	for _, m := range regexp.MustCompile(`(?m)^\s+([A-Z_]+)\s+VARCHAR2\((\d+)\)`).FindAllStringSubmatch(string(b), -1) {
		lebar[m[1]], _ = strconv.Atoi(m[2])
	}
	kolom := map[string]string{"apar": "APAR", "sprinkler": "SPRINKLER", "smokeDetector": "SMOKE_DETECTOR", "hydrant": "HYDRANT",
		"privateTruckBrigade": "PRIVATE_TRUCK_BRIGADE", "privateFireBrigade": "PRIVATE_FIRE_BRIGADE",
		"teamSopSafety": "TEAM_SOP_SAFETY", "teamSopRiskManagement": "TEAM_SOP_RISK_MANAGEMENT", "info": "INFO_FEA"}
	for _, l := range lebarFEA {
		if n, ada := lebar[kolom[l.nama]]; !ada || n != l.n {
			t.Errorf("%s: lebar %d, kolom %s migrasi 190 = %d (ada=%v)", l.nama, l.n, kolom[l.nama], n, ada)
		}
	}
	if len(lebarFEA) != len(kolom) {
		t.Errorf("%d medan, peta %d", len(lebarFEA), len(kolom))
	}
}
