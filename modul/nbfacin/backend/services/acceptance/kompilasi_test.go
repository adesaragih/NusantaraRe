package acceptance

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestJabatanAntreanGagalKompilasi - menukar antrean dan kode jabatan gagal SAAT
// KOMPILASI (tiket NB-11). Kontrolnya wajib terkompilasi.
func TestJabatanAntreanGagalKompilasi(t *testing.T) {
	goBin, err := exec.LookPath("go")
	if err != nil {
		t.Skipf("perintah go tidak ditemukan, uji kompilasi dilewati: %v", err)
	}
	bangun := func(paket string) (string, error) {
		k, err := exec.Command(goBin, "build", "-o", filepath.Join(t.TempDir(), "keluaran"), paket).CombinedOutput()
		return string(k), err
	}
	if k, err := bangun("./testdata/kontroltipe"); err != nil {
		t.Fatalf("kontrol wajib terkompilasi: %v\n%s", err, k)
	}
	k, err := bangun("./testdata/tukarjabatanantrean")
	if err == nil {
		t.Fatal("terkompilasi - Antrean dan Jabatan dapat ditukar")
	}
	// Rumusan lengkap pesan kompilator berbeda antar versi Go; yang dicocokkan
	// hanya inti galatnya per arah tukar.
	for _, m := range []string{"acceptance.Antrean) as acceptance.Jabatan", "acceptance.Jabatan) as acceptance.Antrean"} {
		if !strings.Contains(k, m) {
			t.Errorf("galat kompilasi tidak memuat %q:\n%s", m, k)
		}
	}
}
