// Package pembayaran memuat perhitungan halaman `pyWorkPage.Policy.Payment` (tiket
// 19, disusun agent dari XML atas perintah work owner). Sumbernya jalur unggah CSV
// anggota polis: `UploadCSVPolicyMember_PostAct` langkah 22 dan
// `UploadCSVPolicyMemberEDM_PostAct` langkah 26 → `InputDtlPayment_PreAct` langkah
// 45 → `PremiPaymentMarine` (semua `D:\migrasi\RNM\NB FacIn\Activity\`).
package pembayaran

import (
	"errors"
	"fmt"

	"nusantarare/inti/backend/uang"
	"nusantarare/inti/backend/utils"
	"nusantarare/modul/nbfacin/backend/services/rules"
)

var (
	// ErrCabangEDM - PremiPaymentMarine langkah 2 (`IsEDM`: selisih terhadap
	// `.OldCoverage(1)`, komisi, total cicilan) milik siklus endorsement; tidak
	// diport di sini.
	ErrCabangEDM = errors.New("pembayaran: cabang EDM PremiPaymentMarine belum diport")
	// ErrTanpaCabang - kasus MARINE CARGO yang bukan `IsNB` dan bukan `IsEDM`:
	// [terverifikasi] tidak satu langkah pun menulis `.Policy.Payment`, sehingga
	// nilai lama dibiarkan - tidak ditiru diam-diam.
	ErrTanpaCabang = errors.New("pembayaran: PremiPaymentMarine tanpa cabang yang terbuka")
	// ErrMasukan - premi/diskon kosong atau tak terbaca, atau tanpa coverage.
	// [pertanyaan terbuka] perilaku `@sum` Pega atas nilai kosong dan daftar kosong.
	ErrMasukan = errors.New("pembayaran: masukan coverage kargo tidak sah")
)

// Coverage - satu `.CargoList(n).CoverageList(m)`, angka berupa teks desimal
// bertitik seperti data kerja Pega.
type Coverage struct {
	MataUang string // `.Currency.Name`
	Premium  string // `.Premium`
	Discount string // `.Discount`
}

// Hasil - `.Policy.Payment.Premium` dan `.Diskon`. Berlaku = false bila langkah 45
// tidak terbuka (bukan MARINE CARGO): halaman Payment tidak disentuh.
type Hasil struct {
	Berlaku bool
	Premium uang.Money
	Diskon  uang.Money
}

// Marine - InputDtlPayment_PreAct langkah 45 (`IsMarineCargo`) lalu
// PremiPaymentMarine (ASM-FW-GISFW-WORK!PREMIPAYMENTMARINE). `kargo` = CargoList,
// tiap unsur CoverageList-nya.
//
// [terverifikasi] Langkah 1 (`IsNB`) 1.1: `@sum(.CargoList().CoverageList().Premium)`
// dan `@sum(...Discount)`; 1.2-1.4 berlabel `//` (remark, butir 43). Langkah 2
// (`IsEDM`) berjalan SESUDAH 1 bila terbuka - `IsNB` membaca
// `pyWorkPage.Quotation.StatusBusiness`, `IsEDM` membaca
// `pyWorkPage.OfferFacIn.QuotationData.StatusBusiness`, jadi keduanya bisa benar.
//
// Keputusan agent A45 (menunggu konfirmasi work owner): `@sum` lintas mata uang ditolak
// (`uang.ErrMataUangBerbeda`); Pega menjumlahkan angka tanpa melihat mata uang.
func Marine(k rules.Kasus, kargo [][]Coverage) (Hasil, error) {
	marine, err := rules.Eval("IsMarineCargo", k)
	if err != nil || !marine {
		return Hasil{}, err
	}
	nb, err := rules.Eval("IsNB", k)
	if err != nil {
		return Hasil{}, err
	}
	edm, err := rules.Eval("IsEDM", k)
	if err != nil {
		return Hasil{}, err
	}
	if edm {
		return Hasil{}, ErrCabangEDM
	}
	if !nb {
		return Hasil{}, ErrTanpaCabang
	}
	var hasil Hasil
	for i, cov := range kargo {
		for j, c := range cov {
			premi, err := uangDari(c.MataUang, c.Premium)
			if err != nil {
				return Hasil{}, fmt.Errorf("CargoList(%d).CoverageList(%d).Premium: %w", i+1, j+1, err)
			}
			diskon, err := uangDari(c.MataUang, c.Discount)
			if err != nil {
				return Hasil{}, fmt.Errorf("CargoList(%d).CoverageList(%d).Discount: %w", i+1, j+1, err)
			}
			if !hasil.Berlaku {
				hasil = Hasil{Berlaku: true, Premium: premi, Diskon: diskon}
				continue
			}
			if hasil.Premium, err = hasil.Premium.Add(premi); err != nil {
				return Hasil{}, err
			}
			if hasil.Diskon, err = hasil.Diskon.Add(diskon); err != nil {
				return Hasil{}, err
			}
		}
	}
	if !hasil.Berlaku {
		return Hasil{}, fmt.Errorf("%w: tanpa coverage kargo", ErrMasukan)
	}
	return hasil, nil
}

// uangDari - satu nilai coverage. Keputusan agent A46 (menunggu konfirmasi): mata
// uang kosong ditolak - dua coverage tanpa mata uang lolos Money.Add (mata uangnya
// "sama"), sehingga ketiadaan mata uang tidak boleh lewat diam-diam.
// [terverifikasi] fixture MARINE nyata 104/104 coverage ber-`Currency.Name` "IDR".
func uangDari(mataUang, teks string) (uang.Money, error) {
	if mataUang == "" {
		return uang.Money{}, fmt.Errorf("%w: mata uang kosong", ErrMasukan)
	}
	if teks == "" {
		return uang.Money{}, fmt.Errorf("%w: kosong", ErrMasukan)
	}
	d, err := utils.ParseDecimal(teks)
	if err != nil {
		return uang.Money{}, fmt.Errorf("%w: %v", ErrMasukan, err)
	}
	return uang.Money{Amount: d, Currency: mataUang}, nil
}
