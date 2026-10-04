package repository

// Perekam jejak audit - TANPA Oracle.
//
// Pemilik: A2 (butir am).

import (
	"os"
	"strings"
	"testing"
)

// OQ-M5 DITUTUP (GILIRAN-17): alasan penolakan disimpan di kolom komentar
// baru jejak klaim - migrasi 021, lebar dari preseden KOMITE_COMMENT (013),
// sebab korpus tidak mengekspor rule Property `Remarks`/`KomiteComment`.
func TestMigrasi021KolomKomentarJejak(t *testing.T) {
	naik, err := os.ReadFile("../migrations/021_kolom_komentar_jejak.sql")
	if err != nil {
		t.Fatal(err)
	}
	turun, err := os.ReadFile("../migrations/021_kolom_komentar_jejak_down.sql")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(naik), "ALTER TABLE {skema}.T_CLAIMLF_JEJAK ADD (") ||
		!strings.Contains(string(naik), "KOMENTAR VARCHAR2(4000)") {
		t.Errorf("021 tidak menambah KOMENTAR VARCHAR2(4000):\n%s", naik)
	}
	if !strings.Contains(string(turun), "ALTER TABLE {skema}.T_CLAIMLF_JEJAK DROP (") ||
		!strings.Contains(string(turun), "KOMENTAR") {
		t.Errorf("021_down tidak membuang KOMENTAR:\n%s", turun)
	}
}
