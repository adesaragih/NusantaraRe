package layanan_test

// Token penyimpanan - TANPA Oracle dan TANPA garam sungguhan.
//
// Pemilik: A2 (butir an).

import (
	"errors"
	"strings"
	"testing"
	"time"

	"nusantarare/inti/layanan"
)

// TestRakitTokenGagalTerangTanpaGaram - garam kosong bukan token kosong.
func TestRakitTokenGagalTerangTanpaGaram(t *testing.T) {
	saat := time.Date(2026, 9, 27, 10, 30, 0, 0, time.UTC)
	for _, garam := range []string{"", "   "} {
		if _, err := layanan.RakitToken(garam, saat); !errors.Is(
			err, layanan.ErrGaramTokenKosong) {
			t.Errorf("garam %q: galat = %v, mau ErrGaramTokenKosong", garam, err)
		}
	}
}

// TestTokenBerubahTiapMilidetik - stempel waktunya sampai milidetik.
//
// ⛔ `[data DBA]` procedure memakai `TO_CHAR(SYSTIMESTAMP,'YYYYMMDDHH24MISSFF3')`.
// Memangkasnya ke detik membuat dua token yang terbit dalam detik yang sama
// menjadi identik - dan token yang dapat ditebak bukan token.
func TestTokenBerubahTiapMilidetik(t *testing.T) {
	const garamPalsu = "UJI-GARAM-BUKAN-YANG-SEBENARNYA"
	dasar := time.Date(2026, 9, 27, 10, 30, 0, 0, time.UTC)

	a, err := layanan.RakitToken(garamPalsu, dasar)
	if err != nil {
		t.Fatal(err)
	}
	b, err := layanan.RakitToken(garamPalsu, dasar.Add(time.Millisecond))
	if err != nil {
		t.Fatal(err)
	}
	if a == b {
		t.Error("dua token berselisih satu milidetik identik; stempelnya dipangkas")
	}
	// Bentuknya MD5-hex: 32 karakter, huruf besar.
	if len(a) != 32 {
		t.Errorf("panjang token = %d, mau 32 (MD5 hex)", len(a))
	}
	if a != strings.ToUpper(a) {
		t.Errorf("token bukan huruf besar: %q", a)
	}
	// ⛔ Garamnya TIDAK muncul di tokennya - itu inti hash.
	if strings.Contains(a, garamPalsu) {
		t.Error("garam muncul di token")
	}
}

// TestTokenSamaUntukSaatYangSama - fungsi murni, dapat diuji.
func TestTokenSamaUntukSaatYangSama(t *testing.T) {
	const garamPalsu = "UJI-GARAM"
	saat := time.Date(2026, 9, 27, 10, 30, 0, 123000000, time.UTC)
	a, _ := layanan.RakitToken(garamPalsu, saat)
	b, _ := layanan.RakitToken(garamPalsu, saat)
	if a != b {
		t.Error("token tidak dapat diulang untuk saat yang sama")
	}
	// Garam berbeda -> token berbeda.
	c, _ := layanan.RakitToken("UJI-GARAM-LAIN", saat)
	if a == c {
		t.Error("dua garam berbeda menghasilkan token yang sama")
	}
}
