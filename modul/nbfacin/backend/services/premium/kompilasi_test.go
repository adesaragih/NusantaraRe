package premium

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// bangun menjalankan `go build` dengan argumen tambahan `arg`.
func bangun(t *testing.T, arg ...string) (string, error) {
	t.Helper()
	goBin, err := exec.LookPath("go")
	if err != nil {
		t.Skipf("perintah go tidak ditemukan, uji kompilasi dilewati: %v", err)
	}
	semua := append([]string{"build", "-o", filepath.Join(t.TempDir(), "keluaran")}, arg...)
	keluaran, err := exec.Command(goBin, semua...).CombinedOutput()
	return string(keluaran), err
}

// TestUangTambahRasioGagalKompilasi - menjumlahkan uang dengan rasio gagal
// SAAT KOMPILASI, bukan saat berjalan (tiket NB-01, ADR-F-0004). Setiap
// pasangan punya kontrol yang wajib terkompilasi: pembangun yang rusak tidak
// boleh terbaca sebagai bukti.
func TestUangTambahRasioGagalKompilasi(t *testing.T) {
	for _, u := range []struct {
		nama           string
		kontrol, gagal []string
		mau            []string
	}{
		{
			// Tipe `inti`. ⚠️ Tempatnya kelak di `inti/backend/uang`
			// (`docs/USULAN-PR-TIM-INTI-NB01.md`); dihapus dari sini setelah
			// PR itu diterima.
			nama:    "uang.Money dan uang.Ratio",
			kontrol: []string{"./testdata/kontrol"},
			gagal:   []string{"./testdata/uangtambahrasio"},
			mau: []string{
				"mismatched types uang.Money and uang.Ratio",
				// Rumusan "variable of [struct ]type" berbeda antar versi Go.
				"uang.Ratio) as uang.Money value",
			},
		},
		{
			// Tipe lokal `rasio`: berkas `kompilasi_gagal_rasio.go`.
			nama:    "uang.Money dan rasio lokal",
			kontrol: []string{"."},
			gagal:   []string{"-tags", "gagalkompilasi", "."},
			mau: []string{
				"mismatched types uang.Money and rasio",
				"rasio) as uang.Money value",
			},
		},
	} {
		if keluaran, err := bangun(t, u.kontrol...); err != nil {
			t.Fatalf("%s: kontrol wajib terkompilasi: %v\n%s", u.nama, err, keluaran)
		}
		keluaran, err := bangun(t, u.gagal...)
		if err == nil {
			t.Errorf("%s: terkompilasi - jaminan tipe ADR-F-0004 hilang", u.nama)
			continue
		}
		for _, m := range u.mau {
			if !strings.Contains(keluaran, m) {
				t.Errorf("%s: galat kompilasi tidak memuat %q:\n%s", u.nama, m, keluaran)
			}
		}
	}
}
