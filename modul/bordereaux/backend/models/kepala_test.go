package models

import (
	"strings"
	"testing"
)

// Premium Fire mengikuti kepala sheet `Fire Premium` (format bordereaux 2025): setiap kolom dipetakan, grupnya
// bersebelahan, dan urutan judul daun sama dengan baris kepala Excel (B10:AZ11, tanpa kolom No).
func TestKepalaExcelPremiumFire(t *testing.T) {
	k, err := CariKombinasi(TypePremium, "FIRE")
	if err != nil {
		t.Fatal(err)
	}
	if !AdaKepalaExcel(k) {
		t.Fatal("premi.fire belum memakai kepala Excel")
	}
	var daun, grup []string
	selesai := map[string]bool{}
	sebelum := ""
	for _, c := range k.Kolom {
		if _, ada := kepalaExcel[k.Kode][c.Kolom]; !ada {
			t.Errorf("kolom %s tanpa kepala Excel", c.Kolom)
		}
		kg := KepalaKolom(k, c)
		if kg.Sembunyi {
			continue
		}
		if kg.Grup != sebelum && selesai[kg.Grup] {
			t.Errorf("grup %q terpecah di kolom %s", kg.Grup, c.Kolom)
		}
		if kg.Grup != "" && kg.Grup != sebelum {
			grup = append(grup, kg.Grup)
		}
		selesai[sebelum] = sebelum != ""
		sebelum = kg.Grup
		daun = append(daun, strings.ReplaceAll(kg.Judul, "\n", " "))
	}
	if len(daun) != 50 {
		t.Errorf("%d kolom tampil, sheet Fire Premium 50 (C10:AZ10 tanpa No)", len(daun))
	}
	mauGrup := "*PERIOD OF INSURANCE|*BREAKDOWN OF SI|COORDINATES|*CEDANT'S SHARE|*SPREADING OF RISK"
	if strings.Join(grup, "|") != mauGrup {
		t.Errorf("grup %v", grup)
	}
	// "*Premium (100% Ceded Premium)" diganti "Premium (Reinsurer Share)" (perintah work owner 05-10-2026).
	if daun[0] != "*COB" || daun[6] != "END" || daun[len(daun)-3] != "*SPECIAL ACCEPTANCE (YES / NO)" ||
		daun[len(daun)-4] != "*Premium (Nusantara Re Share)" || daun[len(daun)-5] != "Premium (Reinsurer Share)" {
		t.Errorf("urutan judul %v", daun)
	}
}

// Ke-29 kombinasi mengikuti sheet-nya: setiap kolom berkepala, judul tidak kosong, dan tidak ada lagi kata
// "Indonesia Re" (diganti "Nusantara Re", perintah work owner 05-10-2026). Grup BOLEH muncul dua kali bila sheet-nya
// begitu (MV Premi: MOTOR VEHICLE dan HEAVY EQUIPMENT untuk data kendaraan, lalu lagi untuk sum insured).
func TestKepalaExcelSemuaKombinasi(t *testing.T) {
	if len(kepalaExcel) != len(Kombinasi) {
		t.Errorf("%d kombinasi berkepala Excel, mau %d", len(kepalaExcel), len(Kombinasi))
	}
	for _, k := range Kombinasi {
		if !AdaKepalaExcel(k) {
			t.Errorf("%s tanpa kepala Excel", k.Kode)
			continue
		}
		for _, c := range k.Kolom {
			kg, ada := kepalaExcel[k.Kode][c.Kolom]
			switch {
			case !ada:
				t.Errorf("%s.%s tanpa kepala", k.Kode, c.Kolom)
			case kg.Sembunyi:
			case strings.TrimSpace(kg.Judul) == "":
				t.Errorf("%s.%s judul kosong", k.Kode, c.Kolom)
			case strings.Contains(kg.Grup+kg.Judul, "Indonesia Re"):
				t.Errorf("%s.%s masih memuat Indonesia Re", k.Kode, c.Kolom)
			}
		}
	}
}
