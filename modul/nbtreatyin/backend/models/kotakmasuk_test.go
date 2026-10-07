package models

// NamaKotakMasuk - nama di NBStatus "NB IS IN <nama>'S INBOX" saat Submit (`[keputusan work owner 06-10-2026]`):
// satu pemegang aktif workbasket tujuan = M_LOGIN_GO.NAME; nol atau lebih dari satu = M_WORKBASKET.NAME.

import "testing"

func TestNamaKotakMasuk(t *testing.T) {
	for _, tt := range []struct {
		nama  string
		p     PemegangKotakMasuk
		harap string
	}{
		{"satu pemegang aktif: nama akun", PemegangKotakMasuk{Jumlah: 1, NamaAkun: "UJI Satu", NamaWorkbasket: "UJI WB"}, "UJI Satu"},
		{"lebih dari satu: nama workbasket", PemegangKotakMasuk{Jumlah: 3, NamaAkun: "UJI Satu", NamaWorkbasket: "UJI WB"}, "UJI WB"},
		{"nol pemegang: nama workbasket", PemegangKotakMasuk{NamaWorkbasket: "UJI WB"}, "UJI WB"},
		{"satu pemegang tanpa nama: nama workbasket", PemegangKotakMasuk{Jumlah: 1, NamaWorkbasket: "UJI WB"}, "UJI WB"},
		{"workbasket tanpa nama: ID workbasket", PemegangKotakMasuk{Jumlah: 2}, PosisiSecHead},
	} {
		if got := NamaKotakMasuk(PosisiSecHead, tt.p); got != tt.harap {
			t.Errorf("%s: %q, harap %q", tt.nama, got, tt.harap)
		}
	}
}
