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
