package backend

import (
	"net/http/httptest"
	"testing"
)

// Pelaku dari sesi login MENDAHULUI stub header - dan tanpa sesi maupun stub,
// pelakunya kosong (gagal tertutup).
func TestPelakuSesiMendahuluiStub(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("X-Pelaku", "UJI-HEADER")
	r.Header.Set("X-Peran", "ReasLifeSPV")
	if p := PelakuDari(r, false); p.AkunID != "" {
		t.Errorf("tanpa sesi dan tanpa stub: %+v", p)
	}
	if p := PelakuDari(r, true); p.AkunID != "UJI-HEADER" {
		t.Errorf("stub: %+v", p)
	}
	r = r.WithContext(DenganPelakuSesi(r.Context(), Pelaku{AkunID: "UJI-SESI", Peran: []string{"ReasLifeAdmin"}}))
	for _, stub := range []bool{false, true} {
		if p := PelakuDari(r, stub); p.AkunID != "UJI-SESI" || len(p.Peran) != 1 || p.Peran[0] != "ReasLifeAdmin" {
			t.Errorf("stub=%v: %+v, mau pelaku sesi", stub, p)
		}
	}
}

// Menu akun hasil login dibawa context TERPISAH dari Pelaku: ia menjawab
// "layar mana yang boleh dibuka", bukan "peran apa yang dipegang". Tanpa sesi
// tidak ada daftar sama sekali - bukan daftar kosong yang sah.
func TestAksesMenuSesi(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	if kode, ada := AksesMenuDari(r.Context()); ada || kode != nil {
		t.Errorf("tanpa sesi: %v %v", kode, ada)
	}
	ctx := DenganAksesMenu(r.Context(), []string{"claimlife", "kelolauser"})
	kode, ada := AksesMenuDari(ctx)
	if !ada || len(kode) != 2 || kode[0] != "claimlife" || kode[1] != "kelolauser" {
		t.Errorf("dengan sesi: %v %v", kode, ada)
	}
	if !PunyaMenu(kode, "kelolauser") || PunyaMenu(kode, "premiumlistlife") {
		t.Errorf("PunyaMenu %v", kode)
	}
	// Sesi tanpa satu pun menu tetap SESI: semua ditolak, bukan semua boleh.
	if kode, ada := AksesMenuDari(DenganAksesMenu(r.Context(), nil)); !ada || len(kode) != 0 {
		t.Errorf("sesi tanpa menu: %v %v", kode, ada)
	}
}
