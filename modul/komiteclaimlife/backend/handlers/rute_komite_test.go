package handlers

import (
	"os"
	"regexp"
	"strings"
	"testing"
)

// TestPesanGalatHandlerBerbahasaInggris - work owner 12-10-2026 ("perbaiki tampilannya"): setiap pesan galat yang
// sampai ke layar kasus komite berbahasa Inggris, kalimat tetap (bukan teks galat services/inti yang diteruskan).
func TestPesanGalatHandlerBerbahasaInggris(t *testing.T) {
	src, err := os.ReadFile("rute_komite.go")
	if err != nil {
		t.Fatal(err)
	}
	s := string(src)
	if strings.Contains(s, "err.Error()") {
		t.Error("teks galat services/inti masih diteruskan ke layar (err.Error())")
	}
	pesan := regexp.MustCompile(`galat\.Tulis\(w, [^,]+, "([^"]*)"\)`).FindAllStringSubmatch(s, -1)
	if len(pesan) < 10 {
		t.Fatalf("pesan galat handler terlalu sedikit terbaca: %d", len(pesan))
	}
	indonesia := regexp.MustCompile(`(?i)\b(belum|tidak|bukan|gagal|kasus|permintaan|ditolak|tangga|dikonfigurasi|anda)\b`)
	for _, m := range pesan {
		if indonesia.MatchString(m[1]) {
			t.Errorf("pesan galat berbahasa Indonesia: %q", m[1])
		}
		if !strings.HasSuffix(m[1], ".") || m[1][0] < 'A' || m[1][0] > 'Z' {
			t.Errorf("pesan galat bukan kalimat (huruf besar di awal, titik di akhir): %q", m[1])
		}
	}
}
