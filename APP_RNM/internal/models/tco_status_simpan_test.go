package models

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// AC 39: nilai selain "1" - termasuk kosong dan NULL - selalu gagal.
func TestStatusSimpanSuksesTCO(t *testing.T) {
	teks := func(s string) *string { return &s }
	if !StatusSimpanSuksesTCO(teks("1")) {
		t.Error(`"1" harus sukses`)
	}
	for _, s := range []*string{nil, teks(""), teks("0"), teks(" 1"), teks("1 "), teks("01"), teks("2"), teks("sukses")} {
		if StatusSimpanSuksesTCO(s) {
			t.Errorf("%v dibaca sukses", s)
		}
	}
}

// AC 40 (penyimpangan sadar 8): teks galat salin-tempel prosedur (JSON + _KLAIM)
// tidak ada di sumber PRODUKSI modul - termasuk isi string pesan, yang tidak
// diperiksa penjaga pengenal `TestTCONolNamaBohongDiPengenal`.
func TestTCOTanpaJSONKLAIM(t *testing.T) {
	larang := "JSON" + "_KLAIM"
	diperiksa := 0
	for _, pola := range []string{"../*/tco_*.go", "../repository/migrations/3*.sql", "../../frontend/src/*/*treaty-contract-out*",
		"../../frontend/src/components/treaty-contract-out/*", "../../frontend/src/pages/treaty-contract-out/*"} {
		berkas, _ := filepath.Glob(pola)
		for _, b := range berkas {
			if info, err := os.Stat(b); err != nil || info.IsDir() || strings.Contains(b, "_test.") || strings.Contains(b, ".test.") {
				continue
			}
			isi, err := os.ReadFile(b)
			if err != nil {
				t.Fatal(err)
			}
			diperiksa++
			if strings.Contains(string(isi), larang) {
				t.Errorf("%s memuat %s", b, larang)
			}
		}
	}
	if diperiksa < 50 {
		t.Fatalf("hanya %d berkas terbaca; pencarinya yang rusak", diperiksa)
	}
}
