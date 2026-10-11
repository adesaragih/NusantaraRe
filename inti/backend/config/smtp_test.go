package config

import (
	"errors"
	"strings"
	"testing"
)

// TestBacaSMTP - NBFACIN_SMTP_*: port bukan bilangan / di luar 1..65535 ditolak; kosong = 0; String tidak memuat akun /
// sandi; Dikonfigurasi = host dan alamat pengirim terisi; FROM bukan alamat polos ditolak; TANPA_TLS bawaan false (TLS wajib),
// "true" dibaca, selain boolean ditolak. Nilai uji tiruan (tanpa host / sandi nyata).
func TestBacaSMTP(t *testing.T) {
	t.Setenv("UJISMTP_HOST", "uji.host")
	t.Setenv("UJISMTP_USER", "uji-akun")
	t.Setenv("UJISMTP_PASSWORD", "uji-sandi-rahasia")
	t.Setenv("UJISMTP_FROM", "uji@contoh.test")
	s, err := bacaSMTP("UJISMTP")
	if err != nil || s.Port != 0 || !s.Dikonfigurasi() || s.Dari != "uji@contoh.test" || s.TanpaTLS {
		t.Fatalf("%v %+v", err, s.String())
	}
	if teks := s.String(); strings.Contains(teks, "uji-sandi-rahasia") || strings.Contains(teks, "uji-akun") || strings.Contains(teks, "uji.host") {
		t.Errorf("String membocorkan: %s", teks)
	}
	for _, port := range []string{"abc", "0", "70000"} {
		t.Setenv("UJISMTP_PORT", port)
		if _, err := bacaSMTP("UJISMTP"); !errors.Is(err, ErrKonfigurasi) {
			t.Errorf("port %q: %v", port, err)
		}
	}
	t.Setenv("UJISMTP_PORT", "")
	for _, dari := range []string{"bukan alamat", "UJI <uji@contoh.test>", "a@contoh.test, b@contoh.test"} {
		t.Setenv("UJISMTP_FROM", dari)
		if _, err := bacaSMTP("UJISMTP"); !errors.Is(err, ErrKonfigurasi) {
			t.Errorf("from %q: %v", dari, err)
		}
	}
	t.Setenv("UJISMTP_FROM", "uji@contoh.test")
	t.Setenv("UJISMTP_TANPA_TLS", "true")
	if s, err := bacaSMTP("UJISMTP"); err != nil || !s.TanpaTLS || !strings.Contains(s.String(), "tanpaTLS:true") {
		t.Errorf("tanpa TLS: %v %s", err, s.String())
	}
	t.Setenv("UJISMTP_TANPA_TLS", "ya")
	if _, err := bacaSMTP("UJISMTP"); !errors.Is(err, ErrKonfigurasi) {
		t.Errorf("TANPA_TLS bukan boolean: %v", err)
	}
	t.Setenv("UJISMTP_TANPA_TLS", "")
	t.Setenv("UJISMTP_PORT", "465")
	t.Setenv("UJISMTP_FROM", "")
	if s, err := bacaSMTP("UJISMTP"); err != nil || s.Port != 465 || s.Dikonfigurasi() {
		t.Errorf("465 / tanpa from: %v %s", err, s.String())
	}
}
