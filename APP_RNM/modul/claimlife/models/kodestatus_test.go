package models_test

// Penjaga statik tiket 04: kode status mentah hanya lahir di `models`.
//
// Pemilik: tiket 04. Dibaca sesudah: klaimlife.go.
//
// Nilai "0", "1", dan "2" tidak berarti apa-apa bila dibaca sendirian - "1"
// justru berarti DIAKSEP, bukan ditolak. Literalnya yang tersebar di banyak
// berkas tidak dapat ditelusuri, dan yang salah menulis "2" ketika maksudnya
// "1" tidak akan pernah tertangkap pembaca.

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// polaKodeLiteral mencocokkan pemberian literal angka ke apa pun yang
// bernama kode atau status - termasuk literal ber-kutip-tunggal di dalam SQL.
//
// ⚠️ Dua percobaan sebelumnya lebih sempit daripada kalimatnya sendiri:
// yang pertama hanya `KodeStatus[:=]"n"` sehingga mutasi lewat variabel lokal
// lolos; yang kedua melewatkan `STS_REJECT = '2'` di teks SQL, yaitu lapisan
// yang justru benar-benar menulisnya.
var polaKodeLiteral = regexp.MustCompile(
	`(?i)(kode|status|STS_REJECT)[A-Za-z_]*\s*[:=]+\s*['"][0-9]['"]`)

// polaKomentar membuang komentar sebelum pencocokan. Tanpa ini, prosa yang
// MENGUTIP kodenya - hal biasa di repositori ini - dituduh sebagai kodenya.
var polaKomentar = regexp.MustCompile(`(?m)^\s*//.*$`)

func TestKodeStatusLiteralHanyaDiModels(t *testing.T) {
	var temuan []string
	diperiksa := 0
	err := filepath.Walk(akarAplikasiPindai, func(jalur string, info os.FileInfo, err error) error {
		if err == nil && info.IsDir() && lewatiFolderPindai(info.Name()) {
			return filepath.SkipDir
		}
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(jalur, ".go") ||
			strings.HasSuffix(jalur, "_test.go") {
			return nil
		}
		isi, err := os.ReadFile(jalur)
		if err != nil {
			return err
		}
		rel := relLapisan(jalur)
		// Refactor bentuk B (30-09-2026): kosakata status baris dibagi Claim
		// Life dan Komite, kini tinggal di kontrak lintas modul - SATU tempat.
		if strings.HasPrefix(rel, "models/") || rel == "inti/kontrak/klaim.go" {
			return nil
		}
		bersih := polaKomentar.ReplaceAllString(string(isi), "")
		if polaKodeLiteral.MatchString(bersih) {
			temuan = append(temuan, rel)
		}
		diperiksa++
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	// ⛔ Pemeriksaan diri: penjaga yang berkas sumbernya nol akan selalu hijau
	// dan tidak ada yang menyadarinya. `satutype_test.go` punya pagar yang
	// sama; yang ini semula tidak.
	if diperiksa < 10 {
		t.Fatalf("hanya %d berkas terbaca; pembacanya yang rusak, bukan kodenya",
			diperiksa)
	}
	if len(temuan) != 0 {
		t.Errorf("literal kode status di luar models: %v. Pakai kontrak.KodeOutstanding, "+
			"KodeAksep, atau KodeDitolak - nilainya menipu bila dibaca sendirian, "+
			"dan literal yang tersebar tidak dapat ditelusuri", temuan)
	}
}
