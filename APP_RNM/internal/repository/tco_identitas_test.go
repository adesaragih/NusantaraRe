package repository

import (
	"errors"
	"testing"
)

func TestFormatIdentitasTCO(t *testing.T) {
	kasus := []struct {
		n     int64
		lebar int
		mau   string
	}{
		{1, LebarIdentitasTCO, "1000001"},
		{42, LebarIdentitasTCO, "1000042"},
		{999999, LebarIdentitasTCO, "1999999"},
		{1, LebarIdentitasKlausulTCO, "10000001"},
	}
	for _, k := range kasus {
		dapat, err := FormatIdentitasTCO(k.n, k.lebar)
		if err != nil || dapat != k.mau {
			t.Errorf("(%d,%d) -> %q %v, mau %q", k.n, k.lebar, dapat, err, k.mau)
		}
	}
}

// ⛔ LPAD Oracle memotong diam-diam; di sini gagal terang.
func TestFormatIdentitasTCOMenolakYangTidakMuat(t *testing.T) {
	if _, err := FormatIdentitasTCO(1000000, LebarIdentitasTCO); !errors.Is(err, ErrIdentitasMelampauiLebar) {
		t.Errorf("tujuh digit di lebar enam harus ditolak, dapat %v", err)
	}
	if _, err := FormatIdentitasTCO(-1, LebarIdentitasTCO); !errors.Is(err, ErrIdentitasMelampauiLebar) {
		t.Errorf("negatif harus ditolak, dapat %v", err)
	}
}

// Bentuk yang diterbitkan terbaca kembali oleh pengurai ekor - dua sisi.
func TestIdentitasTCOPulangPergi(t *testing.T) {
	id, err := FormatIdentitasTCO(3188, LebarIdentitasKlausulTCO)
	if err != nil {
		t.Fatal(err)
	}
	n, ok := EkorIdentitasTCO(id, LebarIdentitasKlausulTCO)
	if !ok || n != 3188 {
		t.Errorf("%s -> %d %v", id, n, ok)
	}
	if _, ok := EkorIdentitasTCO(id, LebarIdentitasTCO); ok {
		t.Error("identitas lebar 7 tidak boleh terbaca sebagai lebar 6")
	}
}
