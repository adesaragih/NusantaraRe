package login

import (
	"errors"
	"strings"
	"testing"
)

func TestAturanSandiBaru(t *testing.T) {
	for sandi, mau := range map[string]error{
		"":                        ErrSandiTerlaluPendek,
		"123456789":               ErrSandiTerlaluPendek, // 9
		"1234567890":              nil,                   // 10 - batas yang disetujui work owner
		"ééééééééé":               ErrSandiTerlaluPendek, // 9 karakter, 18 byte: dihitung KARAKTER
		strings.Repeat("a", 72):   nil,
		strings.Repeat("a", 73):   ErrSandiTerlaluPanjang, // bcrypt hanya membaca 72 byte
		strings.Repeat("é", 37):   ErrSandiTerlaluPanjang, // 74 byte
		"sandi dengan spasi 1234": nil,
	} {
		if got := PeriksaSandiBaru(sandi); !errors.Is(got, mau) {
			t.Errorf("PeriksaSandiBaru(%q) = %v, mau %v", sandi, got, mau)
		}
	}
}

func TestHashSandiBcryptDanCocok(t *testing.T) {
	h, err := HashSandi("Sandi-Uji-0001")
	if err != nil {
		t.Fatal(err)
	}
	// bcrypt cost 12: `$2a$12$` - sandi asli tidak pernah ikut tersimpan.
	if !strings.HasPrefix(h, "$2a$12$") || strings.Contains(h, "Sandi-Uji-0001") {
		t.Errorf("hash %q", h)
	}
	if !CocokSandi(h, "Sandi-Uji-0001") {
		t.Error("sandi benar ditolak")
	}
	for _, salah := range []string{"Sandi-Uji-0002", "", "sandi-uji-0001"} {
		if CocokSandi(h, salah) {
			t.Errorf("sandi salah %q diterima", salah)
		}
	}
	if CocokSandi("bukan-hash", "Sandi-Uji-0001") {
		t.Error("hash rusak diterima")
	}
}

func TestSandiSementaraMemenuhiAturanDanAcak(t *testing.T) {
	lihat := map[string]bool{}
	for i := 0; i < 50; i++ {
		s, err := SandiSementara()
		if err != nil {
			t.Fatal(err)
		}
		if len(s) != 16 || PeriksaSandiBaru(s) != nil {
			t.Errorf("sandi sementara %q", s)
		}
		// Tanpa huruf yang mudah tertukar saat dibacakan: 0/O, 1/l/I.
		if strings.ContainsAny(s, "0O1lI") {
			t.Errorf("sandi sementara memuat huruf yang mudah tertukar: %q", s)
		}
		lihat[s] = true
	}
	if len(lihat) != 50 {
		t.Errorf("sandi sementara berulang: %d unik dari 50", len(lihat))
	}
}

func TestPolaAkun(t *testing.T) {
	for akun, sah := range map[string]bool{
		"ade.saragih": true, "UJI-ADMIN": true, "u_1@rnm": true,
		"": false, "ada spasi": false, "a|b": false, strings.Repeat("a", 65): false, strings.Repeat("a", 64): true,
	} {
		if got := AkunSah(akun); got != sah {
			t.Errorf("AkunSah(%q) = %v, mau %v", akun, got, sah)
		}
	}
}
