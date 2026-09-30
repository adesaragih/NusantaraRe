package services

import (
	"fmt"
	"strings"
	"testing"
)

// Keputusan work owner 30-09-2026: kalimat galat Treaty berbahasa Inggris,
// termasuk sentinel `inti/` yang ikut di rantai galatnya.
func TestTeksInggrisTCO(t *testing.T) {
	for _, g := range pesanIntiInggrisTCO {
		p := TeksInggrisTCO(fmt.Sprintf("%s: UJI detail", g.asal))
		if strings.Contains(p, g.asal) || !strings.HasPrefix(p, g.inggris) || !strings.HasSuffix(p, "UJI detail") {
			t.Errorf("%q -> %q, mau diawali %q", g.asal, p, g.inggris)
		}
		// Tidak ada kalimat asal yang memuat kalimat asal lain (urutan ganti tidak berpengaruh).
		for _, lain := range pesanIntiInggrisTCO {
			if lain.asal != g.asal && strings.Contains(g.asal, lain.asal) {
				t.Errorf("%q memuat %q", g.asal, lain.asal)
			}
		}
	}
	if s := TeksInggrisTCO("sudah Inggris"); s != "sudah Inggris" {
		t.Errorf("teks lain berubah: %q", s)
	}
}
