package models

// Kontrak hilir ke Claim Life (tiket 11, E2): kosakata status peserta yang
// ditulis modul ini dinilai dengan penyaring Claim Life APA ADANYA - teksnya
// dibaca dari sumber Claim Life, bukan disalin (nol impor lintas modul).

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// penyaringClaimLife membaca konstanta `penyaringHidup` Claim Life.
func penyaringClaimLife(t *testing.T) (bolehNull bool, mati map[string]bool) {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("..", "..", "..", "claimlife", "backend", "repository", "pesertapolis.go"))
	if err != nil {
		t.Fatalf("sumber Claim Life: %v", err)
	}
	m := regexp.MustCompile("const penyaringHidup = `([^`]+)`").FindStringSubmatch(string(b))
	if m == nil {
		t.Fatal("konstanta penyaringHidup Claim Life tidak ditemukan - kontrak berubah, perbarui uji ini dengan sadar")
	}
	p := m[1]
	daftar := regexp.MustCompile(`NOT IN \(([^)]*)\)`).FindStringSubmatch(p)
	if daftar == nil || !strings.Contains(p, "TRIM(") {
		t.Fatalf("bentuk penyaring tak dikenal: %s", p)
	}
	mati = map[string]bool{}
	for _, v := range strings.Split(daftar[1], ",") {
		mati[strings.Trim(strings.TrimSpace(v), "'")] = true
	}
	return strings.Contains(p, " IS NULL OR "), mati
}

// TestStatusPesertaSejalanPenyaringClaimLife - AC 48/48a: `Delete`/`Batal`
// tersaring, NULL (new business), `Old`, `New` tetap tampil; dan
// `StatusHidup` modul ini memberi keputusan yang SAMA untuk setiap nilai.
func TestStatusPesertaSejalanPenyaringClaimLife(t *testing.T) {
	bolehNull, mati := penyaringClaimLife(t)
	if !bolehNull {
		t.Fatal("penyaring Claim Life membuang peserta new business (status NULL) - AC 48a")
	}
	claim := func(v string) bool { return v == "" || !mati[strings.TrimSpace(v)] }
	for _, v := range []string{"", StatusOld, StatusNew, StatusDelete, StatusBatal} {
		if claim(v) != StatusHidup(v) {
			t.Errorf("status %q: Claim Life %v, modul ini %v", v, claim(v), StatusHidup(v))
		}
	}
	if claim(StatusDelete) || claim(StatusBatal) || !claim(StatusOld) || !claim(StatusNew) {
		t.Errorf("kontrak dilanggar: %v", mati)
	}
}

// TestStatusJenisBukanPenandaHidup - AC 48b: `STATUS`/`STATUS_OLD` adalah
// jenis transaksi dan penanda versi lama, bukan hidup/mati.
func TestStatusJenisBukanPenandaHidup(t *testing.T) {
	_, mati := penyaringClaimLife(t)
	for _, v := range []string{StatusJenis(TypeQR), StatusJenis(TypeTR), StatusLama(StatusOld), StatusLama(StatusNew)} {
		if mati[v] {
			t.Errorf("nilai %q dipakai penyaring hidup/mati", v)
		}
	}
}
