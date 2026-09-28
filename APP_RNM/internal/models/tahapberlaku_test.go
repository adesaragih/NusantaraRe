package models

import "testing"

// TestTahapBerlakuKolomMenangPeranCadangan - satu sumber aturan cadangan tahap.
func TestTahapBerlakuKolomMenangPeranCadangan(t *testing.T) {
	for _, u := range []struct {
		kolom, peran string
		mau          Tahap
	}{
		{"Input Register", PeranAdminLife, TahapInputRegister}, // kolom menang
		{"", PeranAdminLife, TahapOutstanding},                 // cadangan baris lama
		{"", PeranMedicalLife, TahapMedicalCheck},
		{"Claim Analis", "", TahapClaimAnalis},
		{"asing", "asing", TahapTidakDikenal},
	} {
		if got := TahapBerlaku(u.kolom, u.peran); got != u.mau {
			t.Errorf("(%q, %q) = %v, mau %v", u.kolom, u.peran, got, u.mau)
		}
	}
}
