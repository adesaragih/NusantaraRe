package services

import (
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"nusantarare/modul/nbfacin/backend/models"
)

// TestPeriksaOkupasi - tiket 40: tanpa enumerasi / wajib; hanya lebar kolom (pesan ber-indeks) dan jumlah baris.
func TestPeriksaOkupasi(t *testing.T) {
	sah := []models.OkupasiObjek{{OccupationID: "UJI01", OccupationName: "UJI PABRIK", Category: "III", ConstructionClass: "UJI KELAS",
		PctLimit: "70,000 "}, {}}
	if m := periksaOkupasi(0, sah); len(m) != 0 {
		t.Fatalf("sah ditolak: %v", m)
	}
	if m := strings.Join(periksaOkupasi(2, []models.OkupasiObjek{{}, {PctLimit: strings.Repeat("9", 51)}}), "; "); !strings.Contains(m, "baris[2].occupations[1].pctLimit paling banyak 50") {
		t.Errorf("lebar: %q", m)
	}
	if m := periksaOkupasi(0, make([]models.OkupasiObjek, batasOkupasi+1)); len(m) != 1 || !strings.Contains(m[0], "paling banyak 99999") {
		t.Errorf("jumlah: %v", m)
	}
}

// TestLebarOkupasiSamaDenganMigrasi - lebar validasi = lebar kolom migrasi 189 (dipasangkan lewat NAMA kolom).
func TestLebarOkupasiSamaDenganMigrasi(t *testing.T) {
	b, err := os.ReadFile("../migrations/189_t_occupationlist.sql")
	if err != nil {
		t.Fatal(err)
	}
	lebar := map[string]int{}
	for _, m := range regexp.MustCompile(`(?m)^\s+([A-Z_]+)\s+VARCHAR2\((\d+)\)`).FindAllStringSubmatch(string(b), -1) {
		lebar[m[1]], _ = strconv.Atoi(m[2])
	}
	kolom := map[string]string{"occupationId": "OCCUPATION_ID", "occupationName": "OCCUPATION_NAME", "category": "CATEGORY",
		"constructionClass": "DESCRIPTION", "pctLimit": "PCT_LIMIT"}
	for _, l := range lebarOkupasi {
		if n, ada := lebar[kolom[l.nama]]; !ada || n != l.n {
			t.Errorf("%s: lebar %d, kolom %s migrasi 189 = %d (ada=%v)", l.nama, l.n, kolom[l.nama], n, ada)
		}
	}
	if len(lebarOkupasi) != len(kolom) {
		t.Errorf("%d medan, peta %d", len(lebarOkupasi), len(kolom))
	}
}
