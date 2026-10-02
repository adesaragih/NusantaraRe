package premium

import (
	"errors"
	"testing"

	"nusantarare/inti/backend/uang"
	"nusantarare/inti/backend/utils"
)

// Angka `mau` di berkas ini adalah CONTOH HITUNG manual dari rumus literal
// CalculatePremiPA_FacIn, dipilih supaya tiap kasus membedakan satu kesalahan
// port. Rekonsiliasi terhadap sistem lama ada di `rekonsiliasi_test.go`.

// kasusDasar - satu coverage PA metode 3 yang sah; tiap uji mengubah medannya.
func kasusDasar() Input {
	return Input{LiniBisnis: LiniPA, CalculateMethod: "3", MataUang: "IDR", TSI: "1000000", Rate: "2.5", ProRatePercent: "100", Discount: "0"}
}

// TestPremiPAContohHitung - langkah 4 (metode '1') dan 6 (metode '3').
func TestPremiPAContohHitung(t *testing.T) {
	for _, u := range []struct {
		nama string
		ubah func(*Input)
		mau  string
	}{
		{
			// L1003: round4(1.000.000 × 2,5 / 1.000) = 2.500,0000; − 10.
			nama: "fixrate dengan diskon",
			ubah: func(in *Input) { in.Discount = "10" },
			mau:  "2490.0000",
		},
		{
			// Diskon kosong dikurangi nol (empat kasus PA nyata tanpa tag Discount).
			nama: "fixrate diskon kosong",
			ubah: func(in *Input) { in.Discount = "" },
			mau:  "2500.0000",
		},
		{
			// L1003 tidak memakai ProRatePercent: kosong pun sah.
			nama: "fixrate tanpa pro rata",
			ubah: func(in *Input) { in.ProRatePercent = "" },
			mau:  "2500.0000",
		},
		{
			// L713: round4(1.000.000 × 2,5 / 100.000) = 25,0000; × 50 = 1.250,0000.
			// Pembagi 1.000 saja akan memberi 125.000.
			nama: "prorata pembagi komposit",
			ubah: func(in *Input) { in.CalculateMethod, in.ProRatePercent = "1", "50" },
			mau:  "1250.0000",
		},
		{
			// L713 urutan round(a,4) × b: round4(4 × 1 / 100.000 = 0,00004) = 0,0000;
			// × 3 = 0,0000. Urutan terbalik round4(0,00012) memberi 0,0001.
			nama: "prorata bulatkan dulu baru kali",
			ubah: func(in *Input) { in.CalculateMethod, in.TSI, in.Rate, in.ProRatePercent = "1", "4", "1", "3" },
			mau:  "0.0000",
		},
	} {
		in := kasusDasar()
		u.ubah(&in)
		got, err := Calculate(in)
		if err != nil {
			t.Fatalf("%s: %v", u.nama, err)
		}
		if s := utils.FormatDecimal(got.Amount); s != u.mau {
			t.Errorf("%s: premi %s, mau %s", u.nama, s, u.mau)
		}
		if got.Currency != in.MataUang {
			t.Errorf("%s: mata uang %q, mau %q", u.nama, got.Currency, in.MataUang)
		}
	}
}

// TestPremiPAMataUangKosongTidakDiberiBawaan - mata uang kosong ikut ke hasil
// apa adanya: sistem lama tidak pernah menetapkan mata uang bawaan, dan
// penanganannya masih pertanyaan terbuka.
func TestPremiPAMataUangKosongTidakDiberiBawaan(t *testing.T) {
	in := kasusDasar()
	in.MataUang = ""
	got, err := Calculate(in)
	if err != nil {
		t.Fatal(err)
	}
	if got.Currency != "" {
		t.Errorf("mata uang %q, mau tetap kosong", got.Currency)
	}
}

// TestPremiPAPembulatanSeriKeAtas - 1 × 0,05 / 1.000 = 0,00005: tepat di tengah.
// `@Math.divide` = setengah-ke-atas (A37, mengganti butir 14): 0,0001, bukan 0,0000.
// Bukti modenya: dua coverage kasus nyata EDM (TestBerkasKasusP5, rekonsiliasi).
func TestPremiPAPembulatanSeriKeAtas(t *testing.T) {
	in := kasusDasar()
	in.TSI, in.Rate = "1", "0.05"
	got, err := Calculate(in)
	if err != nil || utils.FormatDecimal(got.Amount) != "0.0001" {
		t.Fatalf("dapat %v (%v), mau 0.0001", got, err)
	}
}

// TestPremiPAAngkaDiLuarKonvensiDitolak - semua medan bertitik desimal
// (`utils.ParseDecimal`). Koma K-027 adalah format kolom Oracle, urusan pemuat.
func TestPremiPAAngkaDiLuarKonvensiDitolak(t *testing.T) {
	for _, u := range []struct {
		nama string
		ubah func(*Input)
	}{
		{"rate berkoma", func(in *Input) { in.Rate = "2,5" }},
		{"TSI berkoma", func(in *Input) { in.TSI = "1,5" }},
		{"diskon ribuan", func(in *Input) { in.Discount = "1.000,5" }},
		{"pro rata huruf", func(in *Input) { in.CalculateMethod, in.ProRatePercent = "1", "seratus" }},
	} {
		in := kasusDasar()
		u.ubah(&in)
		if _, err := Calculate(in); !errors.Is(err, ErrAngkaTakTerbaca) {
			t.Errorf("%s: galat %v, mau ErrAngkaTakTerbaca", u.nama, err)
		}
	}
}

// TestPremiPAMedanKosong - kosong = kosong, bukan nol (ADR-U-0022), untuk medan
// yang dipakai rumus. Diskon pengecualiannya (TestPremiPAContohHitung).
func TestPremiPAMedanKosong(t *testing.T) {
	for _, u := range []struct {
		nama string
		ubah func(*Input)
		mau  error
	}{
		{"TSI", func(in *Input) { in.TSI = "" }, uang.ErrUangKosong},
		{"Rate", func(in *Input) { in.Rate = "" }, ErrRasioKosong},
		{"ProRatePercent metode 1", func(in *Input) { in.CalculateMethod, in.ProRatePercent = "1", "" }, ErrRasioKosong},
	} {
		in := kasusDasar()
		u.ubah(&in)
		if _, err := Calculate(in); !errors.Is(err, u.mau) {
			t.Errorf("%s kosong: galat %v, mau %v", u.nama, err, u.mau)
		}
	}
}

// TestPremiPAMetode - nilai di luar '1', '2', '3' tidak menulis premi di sistem
// lama dan ditolak di sini. Perbandingannya TEKS, seperti
// langkah 4-6: "01" bukan '1'.
func TestPremiPAMetode(t *testing.T) {
	for _, u := range []struct {
		metode string
		mau    error
	}{
		{"", ErrMetodeTakDikenal},
		{"4", ErrMetodeTakDikenal},
		{"01", ErrMetodeTakDikenal},
	} {
		in := kasusDasar()
		in.CalculateMethod = u.metode
		if _, err := Calculate(in); !errors.Is(err, u.mau) {
			t.Errorf("metode %q: galat %v, mau %v", u.metode, err, u.mau)
		}
	}
}

// TestLiniDiLuarPetaSkalaPanic - lini bisnis di luar peta K-018 → panic dari
// pintu Calculate juga, sebelum masukan dibaca (resolver diuji sendiri di
// resolver_test.go).
func TestLiniDiLuarPetaSkalaPanic(t *testing.T) {
	for _, lini := range []LiniBisnis{"", "pa", "TRAVEL"} {
		in := kasusDasar()
		in.LiniBisnis = lini
		harusPanic(t, "lini "+string(lini), func() { _, _ = Calculate(in) })
	}
}

// TestPremiPAShortPeriod - tiket NB-04: langkah 5 L858 (metode '2'):
// (round4(TSI × Rate / 1.000) × round4(PctShortPeriod / 100)) − Discount.
// Contoh hitung; kasus nyata metode 2 tidak ada di DDL\CONTOH.
func TestPremiPAShortPeriod(t *testing.T) {
	for _, u := range []struct {
		nama, pct, diskon, mau string
	}{
		// round4(2.500) = 2500.0000; round4(50 / 100) = 0.5000; × → 1250.00000000.
		{"setengah periode", "50", "0", "1250.00000000"},
		// round4(33.33333 / 100 = 0.3333333) = 0.3333 - dibulatkan DULU, baru dikali.
		{"pembulatan faktor periode", "33.33333", "0", "833.25000000"},
		{"dengan diskon", "50", "10", "1240.00000000"},
	} {
		in := kasusDasar()
		in.CalculateMethod, in.PctShortPeriod, in.Discount = "2", u.pct, u.diskon
		got, err := Calculate(in)
		if err != nil {
			t.Errorf("%s: %v", u.nama, err)
			continue
		}
		if s := utils.FormatDecimal(got.Amount); s != u.mau {
			t.Errorf("%s: premi %s, mau %s", u.nama, s, u.mau)
		}
	}
	in := kasusDasar()
	in.CalculateMethod, in.PctShortPeriod = "2", ""
	if _, err := Calculate(in); !errors.Is(err, ErrRasioKosong) {
		t.Errorf("PctShortPeriod kosong: galat %v, mau ErrRasioKosong", err)
	}
}

// TestPremiPADiskonPersen - tiket NB-04: langkah 2 (`param.DiscountType ==
// "Percent"`) menghitung ulang `.Discount = @Math.divide((.DiscountPercentage *
// .Premium),100,20)` dari `.Premium` SEBELUM rumus premi berjalan, lalu rumus
// mengurangkannya. Langkah 3 ("Amount") hanya menulis .DiscountPercentage.
func TestPremiPADiskonPersen(t *testing.T) {
	in := kasusDasar()
	// 10 × 2.500 / 100 = 250 (20 desimal); fixrate 2500.0000 − 250 = 2250.
	in.DiscountType, in.DiscountPercentage, in.PremiSebelumnya, in.Discount = "Percent", "10", "2500", "999"
	got, err := Calculate(in)
	if err != nil {
		t.Fatal(err)
	}
	if s := utils.FormatDecimal(got.Amount); s != "2250.00000000000000000000" {
		t.Errorf("premi %s, mau 2250.00000000000000000000 (diskon dihitung ulang, bukan 999)", s)
	}
	// "Amount": .Discount masukan dipakai apa adanya.
	in.DiscountType = "Amount"
	if got, err = Calculate(in); err != nil || utils.FormatDecimal(got.Amount) != "1501.0000" {
		t.Errorf("Amount: premi %v (%v), mau 1501.0000", got.Amount, err)
	}
	// "Percent" tanpa premi sebelumnya: kosong bukan nol.
	in.DiscountType, in.PremiSebelumnya = "Percent", ""
	if _, err := Calculate(in); !errors.Is(err, uang.ErrUangKosong) {
		t.Errorf("Percent tanpa premi sebelumnya: galat %v, mau ErrUangKosong", err)
	}
}
