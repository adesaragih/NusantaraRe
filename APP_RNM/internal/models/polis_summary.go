package models

// Rekap uang per mata uang - tiket 05a PremiumList Life.
//
// Untuk apa berkas ini: rumus `DataTransform/AppendCurrencySummary_DT.xml`,
// ditulis ulang MURNI. Nol Oracle, nol float.
//
// `[terverifikasi]` DT dibaca utuh 28-09-2026 (5.992 baris, 121 penetapan)
// lewat URUTAN AKSI, bukan pasangan nama/nilai:
//
//	grep -o "<pyActionName>[^<]*\|<pyPropertiesName>[^<]*\|<pyPropertiesValue>[^<]*" \
//	  DataTransform/AppendCurrencySummary_DT.xml
//
// ⛔ KENAPA URUTAN AKSI DAN BUKAN PASANGAN. Di ekspor Pega, kondisi sebuah
// `WHEN` menempati medan `pyPropertiesName` - medan yang sama dengan nama
// properti yang di-SET. Memasangkan nama dengan nilai karena itu MERATAKAN
// struktur percabangannya, dan empat cabang yang saling meniadakan terbaca
// seperti empat penjumlahan berurutan. Bacaan itu sudah pernah terjadi, dan
// akibatnya melipatgandakan premi.
//
// Bentuk sesungguhnya:
//
//	WHEN Param.Currency == .CURRENCY            <- pengelompokan mata uang
//	  WHEN pyWorkPage.Type == "QR"  -> += .GROSS_PREMIUM
//	  WHEN pyWorkPage.Type == "QP"  -> += .GROSS_PREMIUM_REFUND
//	  WHEN pyWorkPage.Type == "TP"  -> += .GROSS_PREMIUM_RETRO
//	  WHEN pyWorkPage.Type == "TR"  -> += .GROSS_PREMIUM_REFUND_RETRO
//
// Satu polis punya SATU `Type`, jadi TEPAT SATU cabang menyala.
//
// Dibaca sesudah: polis_unggah.go (sumber barisnya), polis_nomor.go.

import (
	"fmt"
	"strings"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/pkg/utils"
)

// AngkaDesimalBalance adalah angka desimal pembulatan akhir `BALANCE`.
//
// `[terverifikasi]` `@divide(Param.Balance,1,4)` - pembagian dengan SATU,
// yaitu pembulatan ke empat angka yang ditulis sebagai pembagian.
//
// ⛔ KEANEHAN WARISAN, disalin apa adanya. Menuliskannya sebagai pembulatan
// yang jujur akan menghilangkan jejak bahwa rule aslinya menulis begini - dan
// jejak itu yang menjelaskan kenapa angkanya berhenti di empat.
const AngkaDesimalBalance = 4

// premiumMenurutTipe memetakan `Type` polis ke kolom yang menyusun `PREMIUM`.
//
// ⛔ SATU KOLOM, BUKAN EMPAT. Lihat kepala berkas.
var premiumMenurutTipe = map[string]string{
	TipePLQuotationRealisasi: "GROSS_PREMIUM",
	TipePLQuotationProposal:  "GROSS_PREMIUM_REFUND",
	TipePLTreatyProposal:     "GROSS_PREMIUM_RETRO",
	TipePLTreatyRealisasi:    "GROSS_PREMIUM_REFUND_RETRO",
}

// KolomPremiumUntukTipe menjawab kolom mana yang menyusun `PREMIUM`.
func KolomPremiumUntukTipe(tipe string) (string, error) {
	k, ada := premiumMenurutTipe[strings.TrimSpace(tipe)]
	if !ada {
		return "", fmt.Errorf("%w: %q", ErrTipePLTanpaCabang, tipe)
	}
	return k, nil
}

// sukuBalance adalah rumus `BALANCE` satu cabang: yang ditambah dan dikurangi.
//
// ⛔ DISALIN APA ADANYA, termasuk tandanya yang janggal. TIGA keanehan
// warisan, dan tidak satu pun "diperbaiki":
//
//  1. `BROKERAGE_FEE_RETRO` DITAMBAH di cabang TP/TR, sedangkan
//     `BROKERAGE_FEE` DIKURANGI di cabang QR. Biaya yang menambah saldo di
//     satu cabang dan mengurangi di cabang lain.
//  2. Cabang QP mengurangi `TAX`, `PROF_COMM`, `CLAIM` - BUKAN padanan
//     `*_REFUND`-nya, padahal ketiga biaya lain di cabang itu memakai
//     `*_REFUND`.
//  3. `CLAIM_AMOUNT` MENAMBAH di cabang QP; di cabang lain ia tidak muncul.
//
// Menormalkan tandanya mengubah angka uang yang sudah beredar.
type sukuBalance struct {
	Tambah []string
	Kurang []string
}

// balanceMenurutTipe adalah keempat cabang, VERBATIM.
var balanceMenurutTipe = map[string]sukuBalance{
	TipePLQuotationRealisasi: {
		Tambah: []string{"GROSS_PREMIUM"},
		Kurang: []string{"DEDUCTION", "RI_ADMIN_FEE", "BROKERAGE_FEE", "TAX",
			"PROF_COMM", "CLAIM"},
	},
	TipePLQuotationProposal: {
		Tambah: []string{"GROSS_PREMIUM_REFUND", "CLAIM_AMOUNT"},
		Kurang: []string{"DEDUCTION_REFUND", "BROKERAGE_FEE_REFUND",
			"RI_ADMIN_FEE_REFUND", "TAX", "PROF_COMM", "CLAIM"},
	},
	TipePLTreatyProposal: {
		// ⚠️ BROKERAGE_FEE_RETRO di sisi TAMBAH - keanehan nomor 1.
		Tambah: []string{"GROSS_PREMIUM_RETRO", "BROKERAGE_FEE_RETRO"},
		Kurang: []string{"DISCOUNT_PREMIUM_RETRO", "RI_ADMIN_FEE_RETRO"},
	},
	TipePLTreatyRealisasi: {
		Tambah: []string{"GROSS_PREMIUM_REFUND_RETRO", "BROKERAGE_FEE_REFUND_RETRO"},
		Kurang: []string{"DISCOUNT_PREMIUM_REFUND_RETRO", "RI_ADMIN_FEE_REFUND_RETRO"},
	},
}

// SukuBalanceUntukTipe mengembalikan kolom penambah dan pengurang `BALANCE`.
//
// ⚠️ Disediakan supaya uji dapat menyebut kolomnya satu per satu - rumus uang
// yang hanya dapat diperiksa lewat hasil akhirnya adalah rumus yang salahnya
// baru terlihat saat angkanya sudah dipakai.
func SukuBalanceUntukTipe(tipe string) (tambah, kurang []string, err error) {
	s, ada := balanceMenurutTipe[strings.TrimSpace(tipe)]
	if !ada {
		return nil, nil, fmt.Errorf("%w: %q", ErrTipePLTanpaCabang, tipe)
	}
	return s.Tambah, s.Kurang, nil
}

// BarisUang adalah nilai uang satu baris peserta, berkunci nama kolom.
//
// ⚠️ Kolom yang ABSEN dibaca sebagai NOL dalam penjumlahan - bukan galat.
// Rule aslinya menjumlahkan properti yang kosong sebagai nol, dan baris yang
// tidak punya kolom retro memang bukan baris yang rusak.
type BarisUang map[string]*apd.Decimal

// ambil mengembalikan nilai sebuah kolom, atau nol bila absen.
func (b BarisUang) ambil(kolom string) *apd.Decimal {
	if d, ada := b[kolom]; ada && d != nil {
		return d
	}
	return apd.New(0, 0)
}

// RekapMataUang adalah rekap satu mata uang.
type RekapMataUang struct {
	Currency string
	// Premium, Commission, Balance - ketiganya turunan.
	Premium    *apd.Decimal
	Commission *apd.Decimal
	Balance    *apd.Decimal
	// Jumlah memuat ke-32 kolom uang yang dijumlah apa adanya.
	Jumlah map[string]*apd.Decimal
	// CacahBaris adalah cacah peserta yang masuk rekap ini.
	CacahBaris int
}

// RekapPerMataUang menjumlahkan baris peserta menjadi rekap per mata uang.
//
// `tipe` adalah `T_PREMIUM_LIST.TYPE` - ia yang memilih cabang `PREMIUM` dan
// `BALANCE`, dan ia milik POLIS, bukan milik baris.
//
// ⛔ PEMBULATAN HANYA DI AKHIR. `BALANCE` dibulatkan ke empat angka SESUDAH
// seluruh baris dijumlah, persis `@divide(Param.Balance,1,4)` yang berdiri di
// luar perulangan. Membulatkan per baris lalu menjumlah menghasilkan angka
// yang BERBEDA, dan selisihnya tumbuh bersama cacah peserta.
//
// ⛔ NOL FLOAT (ADR-U-0003, ADR-U-0016).
func RekapPerMataUang(tipe string, baris []BarisUang, mataUang []string) (
	[]RekapMataUang, error) {

	kolomPremium, err := KolomPremiumUntukTipe(tipe)
	if err != nil {
		return nil, err
	}
	tambah, kurang, err := SukuBalanceUntukTipe(tipe)
	if err != nil {
		return nil, err
	}
	if len(baris) != len(mataUang) {
		return nil, fmt.Errorf(
			"models: %d baris uang tetapi %d mata uang - keduanya harus sejajar",
			len(baris), len(mataUang))
	}

	ctx := utils.DecimalContext()
	urut := []string{}
	rekap := map[string]*RekapMataUang{}

	for i, b := range baris {
		// ⛔ Pengelompokannya VERBATIM `WHEN Param.Currency == .CURRENCY`.
		cur := strings.TrimSpace(mataUang[i])
		r, ada := rekap[cur]
		if !ada {
			r = &RekapMataUang{
				Currency:   cur,
				Premium:    apd.New(0, 0),
				Commission: apd.New(0, 0),
				Balance:    apd.New(0, 0),
				Jumlah:     map[string]*apd.Decimal{},
			}
			for _, k := range KolomUangUnggah {
				r.Jumlah[k] = apd.New(0, 0)
			}
			rekap[cur] = r
			urut = append(urut, cur)
		}
		r.CacahBaris++

		for _, k := range KolomUangUnggah {
			if _, err := ctx.Add(r.Jumlah[k], r.Jumlah[k], b.ambil(k)); err != nil {
				return nil, fmt.Errorf("models: menjumlah %s: %w", k, err)
			}
		}
		if _, err := ctx.Add(r.Premium, r.Premium, b.ambil(kolomPremium)); err != nil {
			return nil, fmt.Errorf("models: menjumlah PREMIUM: %w", err)
		}
		// ⛔ `COMMISSION` = Σ `.COMM` SAJA, tanpa cabang `Type`. `PROF_COMM`
		// dan `OVR_COMM` parameter tersendiri - menggabungkannya di sini
		// menjumlahkan komisi dua kali.
		if _, err := ctx.Add(r.Commission, r.Commission, b.ambil("COMM")); err != nil {
			return nil, fmt.Errorf("models: menjumlah COMMISSION: %w", err)
		}
		for _, k := range tambah {
			if _, err := ctx.Add(r.Balance, r.Balance, b.ambil(k)); err != nil {
				return nil, fmt.Errorf("models: menjumlah BALANCE (+%s): %w", k, err)
			}
		}
		for _, k := range kurang {
			if _, err := ctx.Sub(r.Balance, r.Balance, b.ambil(k)); err != nil {
				return nil, fmt.Errorf("models: menjumlah BALANCE (-%s): %w", k, err)
			}
		}
	}

	keluar := make([]RekapMataUang, 0, len(urut))
	for _, cur := range urut {
		r := rekap[cur]
		// Pembulatan akhir, dan HANYA di sini.
		bulat := apd.New(0, 0)
		if _, err := ctx.Quantize(bulat, r.Balance, -AngkaDesimalBalance); err != nil {
			return nil, fmt.Errorf("models: membulatkan BALANCE %s: %w", cur, err)
		}
		r.Balance = bulat
		keluar = append(keluar, *r)
	}
	return keluar, nil
}
