package premium

import (
	"errors"
	"testing"

	"nusantarare/inti/backend/utils"
)

// kasusMBU - satu coverage MBU.
func kasusMBU(tsi, rate, loading, proRata string) Input {
	return Input{LiniBisnis: LiniMBU, MataUang: "IDR", TSI: tsi, Rate: rate, Loading: loading, ProRatePercentCoverage: proRata}
}

// TestPremiMBUBersarang - contoh hitung yang MEMBEDAKAN bentuk bersarang L1144
// dari bentuk yang diratakan menjadi satu pembagian ÷10.000, dan yang
// membedakan urutan .Loading. Dicari sistematis (bukan dikarang), tanpa kasus seri.
func TestPremiMBUBersarang(t *testing.T) {
	for _, u := range []struct {
		nama string
		in   Input
		mau  string
	}{
		// round4(12345 × 0.0137 / 100 = 1.691265) = 1.6913; × 33.333 / 100 =
		// 0.56376… → 0.5638. Diratakan: 12345 × 0.0137 × 33.333 / 10000 =
		// 0.56374… → 0.5637.
		{"bersarang, bukan ÷10000", kasusMBU("12345", "0.0137", "", "33.333"), "0.5638"},
		{"bersarang, pro-rata 150", kasusMBU("12345", "0.0137", "", "150"), "2.5370"},
		// .Loading DIJUMLAHKAN ke rate sebelum dikali: 1.000.000 × (2,5 + 0,5) / 100
		// = 30.000. Ditambahkan sesudah akan memberi 25.000,5.
		{"loading sebelum perkalian", kasusMBU("1000000", "2.5", "0.5", "100"), "30000.0000"},
		// `@if(.ProRatePercent=="",100,.ProRatePercent)`.
		{"pro-rata kosong = 100", kasusMBU("1000000", "2.5", "", ""), "25000.0000"},
	} {
		got, err := Calculate(u.in)
		if err != nil {
			t.Errorf("%s: %v", u.nama, err)
			continue
		}
		if s := utils.FormatDecimal(got.Amount); s != u.mau {
			t.Errorf("%s: premi %s, mau %s", u.nama, s, u.mau)
		}
	}
}

// TestPremiMBUMedanKosong - Rate kosong ditolak (Loading kosong = 0 dibuktikan
// kasus MBU nyata di TestRekonsiliasiDaftarIzin).
func TestPremiMBUMedanKosong(t *testing.T) {
	if _, err := Calculate(kasusMBU("1000000", "", "", "100")); !errors.Is(err, ErrRasioKosong) {
		t.Errorf("rate kosong: galat %v, mau ErrRasioKosong", err)
	}
}
