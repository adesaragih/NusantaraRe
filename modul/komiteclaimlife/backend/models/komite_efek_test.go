package models

import "testing"

// TestKeadaanEfekKasus - tiket 08, AC 23/25 spec.
func TestKeadaanEfekKasus(t *testing.T) {
	for _, u := range []struct {
		kode []string
		mau  string
	}{
		{nil, ""},
		{[]string{"selesai", "selesai", "selesai"}, KataEfekTuntas},
		{[]string{"selesai", "antre"}, KataKeputusanTersimpan},
		{[]string{"selesai", "jalan", "gagal-permanen"}, KataEfekPerluIntervensi},
		{[]string{"selesai", "kode-asing"}, KataEfekPerluIntervensi},
	} {
		if got := KeadaanEfekKasus(u.kode); got != u.mau {
			t.Errorf("%v: %q, mau %q", u.kode, got, u.mau)
		}
	}
}
