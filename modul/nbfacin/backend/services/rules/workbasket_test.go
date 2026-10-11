package rules

import (
	"os"
	"regexp"
	"strings"
	"testing"
)

// namaJUWLama - workbasket lama ReasFacInJuniorUnderwriting TANPA akhiran A / B (ADR-0032, keputusan work owner 07-10-2026
// dan 08-10-2026: When yang memakainya dibaca ReasFacInJuniorUnderwritingB).
var namaJUWLama = regexp.MustCompile(`ReasFacInJuniorUnderwriting([^AB]|$)`)

// TestRegistryTanpaWorkbasketJUWLama - dua cara: (1) setiap teks registry terurai (logika, kiri, kanan, rujukan, catatan);
// (2) teks berkas registry_gen.go, termasuk komentar. Instrumen diuji atas contoh yang sudah diketahui.
func TestRegistryTanpaWorkbasketJUWLama(t *testing.T) {
	for teks, mau := range map[string]bool{`"ReasFacInJuniorUnderwriting"`: true, "ReasFacInJuniorUnderwriting": true,
		"ReasFacInJuniorUnderwriting ": true, `"ReasFacInJuniorUnderwritingB"`: false, "ReasFacInJuniorUnderwritingA": false} {
		if namaJUWLama.MatchString(teks) != mau {
			t.Fatalf("instrumen %q: mau %v", teks, mau)
		}
	}
	for nama, pr := range registry {
		teks := []string{pr.logika, pr.catatan, pr.panik}
		for _, kd := range pr.kondisi {
			teks = append(teks, kd.kiri, kd.kanan.teks, kd.rujukan)
		}
		for _, s := range teks {
			if namaJUWLama.MatchString(s) {
				t.Errorf("%s memuat workbasket lama: %q", nama, s)
			}
		}
	}
	isi, err := os.ReadFile("registry_gen.go")
	if err != nil {
		t.Fatal(err)
	}
	n := 0
	for i, baris := range strings.Split(strings.ReplaceAll(string(isi), "\r\n", "\n"), "\n") {
		if strings.Contains(baris, "ReasFacInJuniorUnderwritingB") {
			n++
		}
		if namaJUWLama.MatchString(baris) {
			t.Errorf("registry_gen.go:%d memuat workbasket lama: %s", i+1, strings.TrimSpace(baris))
		}
	}
	// Pembanding: pemetaan terpakai (ISUWSPREAD kondisi D; korpus NB FacIn\When: hanya IsUWSpread.xml memuat nama lama).
	if n == 0 || registry["ISUWSPREAD"].kondisi["D"].kanan.teks != "ReasFacInJuniorUnderwritingB" {
		t.Errorf("ReasFacInJuniorUnderwritingB di registry_gen.go %d baris; ISUWSPREAD D %q", n, registry["ISUWSPREAD"].kondisi["D"].kanan.teks)
	}
}
