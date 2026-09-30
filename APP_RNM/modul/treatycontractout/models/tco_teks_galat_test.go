package models

// Penjaga teks galat - tiket 09 Treaty Contract Out. Sejak OQ-TCO-19
// (keputusan work owner 29-09-2026) simpan utuh dibuang; AC 40 tetap berlaku
// untuk setiap pesan galat modul.

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// AC 40 (penyimpangan sadar 8): teks galat salin-tempel prosedur (JSON + _KLAIM)
// tidak ada di sumber PRODUKSI modul - termasuk isi string pesan, yang tidak
// diperiksa penjaga pengenal `TestTCONolNamaBohongDiPengenal`.
func TestTCOTanpaJSONKLAIM(t *testing.T) {
	larang := "JSON" + "_KLAIM"
	diperiksa := 0
	// Refactor bentuk B (30-09-2026): modul ini kini di modul/treatycontractout/; rentang
	// migrasinya (3xx) - bila kelak ada - tinggal di modul/treatycontractout/migrations.
	// Refactor bentuk B paket 7: frontend modul ini kini frontend/src/modul/treaty-contract-out/.
	// Paket 8: seluruh akar folder itu (api.ts, menu.ts, rute.tsx, labels.ts) -
	// dulu klien backendnya di services/api.ts bersama dan tidak terbaca.
	for _, pola := range []string{"../*/tco_*.go", "../migrations/3*.sql", "../../../frontend/src/modul/treaty-contract-out/*",
		"../../../frontend/src/modul/treaty-contract-out/components/*", "../../../frontend/src/modul/treaty-contract-out/pages/*"} {
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
