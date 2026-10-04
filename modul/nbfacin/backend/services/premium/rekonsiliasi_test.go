package premium

import (
	"encoding/json"
	"os"
	"testing"

	"nusantarare/inti/backend/utils"
)

// TestRekonsiliasiDaftarIzin - REKONSILIASI tahap 1 (ADR-F-0001): premi sistem
// lama dari kasus nyata PA dan MBU (`testdata/daftarizin/premi.json`, asal dan
// daftar-izinnya di README di sana), dibandingkan sampai digit terakhir tanpa
// toleransi. Kerangka lintas-kasus memakai berkas yang sama:
// `services/rekonsiliasi` (tiket 16).
func TestRekonsiliasiDaftarIzin(t *testing.T) {
	isi, err := os.ReadFile("testdata/daftarizin/premi.json")
	if err != nil {
		t.Fatal(err)
	}
	var kasus []struct {
		Kasus string
		Input Input
		Premi string
	}
	if err := json.Unmarshal(isi, &kasus); err != nil {
		t.Fatal(err)
	}
	if len(kasus) != 8 {
		t.Fatalf("%d kasus, mau 8 (4 PA + 4 MBU)", len(kasus))
	}
	for _, u := range kasus {
		got, err := Calculate(u.Input)
		if err != nil {
			t.Errorf("%s: %v", u.Kasus, err)
			continue
		}
		if s := utils.FormatDecimal(got.Amount); s != u.Premi {
			t.Errorf("%s: premi %s, sistem lama %s", u.Kasus, s, u.Premi)
		}
	}
}
