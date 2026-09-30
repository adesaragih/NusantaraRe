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

	"nusantarare/inti/backend/penomor"
	"nusantarare/inti/backend/utils"
)

// AngkaDesimalRekap adalah angka desimal pembulatan akhir SELURUH rekap.
//
// `[terverifikasi]` `@divide(Param.X,1,4)` - pembagian dengan SATU, yaitu
// pembulatan ke empat angka yang ditulis sebagai pembagian.
//
// ⚠️ RALAT 28-09-2026 (giliran 10): bagian 1 menamainya `AngkaDesimalBalance`
// karena membacanya hanya pada `.BALANCE`. DT langkah 2.39-2.74 menulis
// `@divide(…,1,4)` pada KETIGA PULUH ENAM keluaran - `.PREMIUM` (2.39),
// `.COMMISSION` (2.40), tiga puluh tiga kolom `KolomJumlahSummary`, lalu
// `.BALANCE` (2.74). Namanya diganti supaya tidak berbohong.
//
// ⛔ KEANEHAN WARISAN, disalin apa adanya. Menuliskannya sebagai pembulatan
// yang jujur akan menghilangkan jejak bahwa rule aslinya menulis begini - dan
// jejak itu yang menjelaskan kenapa angkanya berhenti di empat.
const AngkaDesimalRekap = 4

// KolomJumlahSummary adalah kolom yang dijumlah APA ADANYA per mata uang.
//
// `[terverifikasi]` DT langkah 2.38.1.5-2.38.1.38 (penjumlahan `param.X +
// .X`) dan 2.41-2.73 (keluaran `@divide`), urut dokumen. Tiga puluh tiga
// kolom; `COMM` tidak di sini karena ia menyusun `COMMISSION`.
//
// ⛔ BUKAN `KolomUangUnggah`. Bagian 1 menjumlah ke-32 kolom unggahan - daftar
// yang memuat `NET_PREMIUM`, `GROSS_PREMIUM`, `SUM_INSURED`, `EM_PERCENT`,
// `FLEET_DISCOUNT`, `COMM`, `SHARE_NUSANTARA_RE` yang TIDAK pernah dijumlah
// DT, dan MELEWATKAN `DEDUCTION`, `RI_ADMIN_FEE*`, `DEDUCTION_REFUND`,
// `SUM_AT_RISK_GROSS`, `SHARE_NUSANTARA_RE_GROSS` yang dijumlah DT. Rekap
// dengan daftar itu menyimpan nol di kolom yang seharusnya berisi.
var KolomJumlahSummary = []string{
	"BROKERAGE_FEE", "OVR_COMM", "TAX", "PROF_COMM", "CLAIM",
	"NET_PREMIUM_REFUND", "GROSS_PREMIUM_REFUND", "COMM_REFUND",
	"BROKERAGE_FEE_REFUND", "OVR_COMM_REFUND", "TAX_REFUND",
	"SHARE_RETRO", "GROSS_PREMIUM_RETRO", "DISCOUNT_PREMIUM_RETRO",
	"OVR_COMM_RETRO", "BROKERAGE_FEE_RETRO", "NET_PREMIUM_RETRO",
	"GROSS_PREMIUM_REFUND_RETRO", "DISCOUNT_PREMIUM_REFUND_RETRO",
	"OVR_COMM_REFUND_RETRO", "BROKERAGE_FEE_REFUND_RETRO",
	"NET_PREMIUM_REFUND_RETRO", "CLAIM_AMOUNT", "RI_ADMIN_FEE_RETRO",
	"RI_ADMIN_FEE_REFUND_RETRO", "RI_ADMIN_FEE_REFUND", "DEDUCTION_REFUND",
	"RI_ADMIN_FEE", "DEDUCTION", "SUM_REASURED", "SHARE_NUSANTARA_RE_GROSS",
	"SUM_AT_RISK_GROSS", "CEDING_RETENTION",
}

// KolomBacaSummary adalah kolom baris peserta yang dibutuhkan rekap.
//
// Gabungan `KolomJumlahSummary`, `COMM`, kolom `PREMIUM` keempat cabang, dan
// seluruh suku `BALANCE` - tanpa kembar, urut kemunculan. Pembaca repository
// memakainya supaya tidak ada kolom yang dirumuskan di sini tetapi tidak
// pernah dibaca dari basis data (dan diam-diam bernilai nol).
func KolomBacaSummary() []string {
	lihat := map[string]bool{}
	var keluar []string
	tambah := func(k ...string) {
		for _, x := range k {
			if !lihat[x] {
				lihat[x] = true
				keluar = append(keluar, x)
			}
		}
	}
	tambah(KolomJumlahSummary...)
	tambah("COMM")
	for _, tipe := range []string{penomor.TipePLQuotationRealisasi, penomor.TipePLQuotationProposal,
		penomor.TipePLTreatyProposal, penomor.TipePLTreatyRealisasi} {
		tambah(premiumMenurutTipe[tipe])
		tambah(balanceMenurutTipe[tipe].Tambah...)
		tambah(balanceMenurutTipe[tipe].Kurang...)
	}
	return keluar
}

// premiumMenurutTipe memetakan `Type` polis ke kolom yang menyusun `PREMIUM`.
//
// ⛔ SATU KOLOM, BUKAN EMPAT. Lihat kepala berkas.
var premiumMenurutTipe = map[string]string{
	penomor.TipePLQuotationRealisasi: "GROSS_PREMIUM",
	penomor.TipePLQuotationProposal:  "GROSS_PREMIUM_REFUND",
	penomor.TipePLTreatyProposal:     "GROSS_PREMIUM_RETRO",
	penomor.TipePLTreatyRealisasi:    "GROSS_PREMIUM_REFUND_RETRO",
}

// KolomPremiumUntukTipe menjawab kolom mana yang menyusun `PREMIUM`.
func KolomPremiumUntukTipe(tipe string) (string, error) {
	k, ada := premiumMenurutTipe[strings.TrimSpace(tipe)]
	if !ada {
		return "", fmt.Errorf("%w: %q", penomor.ErrTipePLTanpaCabang, tipe)
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
	penomor.TipePLQuotationRealisasi: {
		Tambah: []string{"GROSS_PREMIUM"},
		Kurang: []string{"DEDUCTION", "RI_ADMIN_FEE", "BROKERAGE_FEE", "TAX",
			"PROF_COMM", "CLAIM"},
	},
	penomor.TipePLQuotationProposal: {
		Tambah: []string{"GROSS_PREMIUM_REFUND", "CLAIM_AMOUNT"},
		Kurang: []string{"DEDUCTION_REFUND", "BROKERAGE_FEE_REFUND",
			"RI_ADMIN_FEE_REFUND", "TAX", "PROF_COMM", "CLAIM"},
	},
	penomor.TipePLTreatyProposal: {
		// ⚠️ BROKERAGE_FEE_RETRO di sisi TAMBAH - keanehan nomor 1.
		Tambah: []string{"GROSS_PREMIUM_RETRO", "BROKERAGE_FEE_RETRO"},
		Kurang: []string{"DISCOUNT_PREMIUM_RETRO", "RI_ADMIN_FEE_RETRO"},
	},
	penomor.TipePLTreatyRealisasi: {
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
		return nil, nil, fmt.Errorf("%w: %q", penomor.ErrTipePLTanpaCabang, tipe)
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
	// Jumlah memuat ke-33 kolom `KolomJumlahSummary`, dijumlah apa adanya.
	Jumlah map[string]*apd.Decimal
	// CacahBaris adalah cacah peserta yang masuk rekap ini.
	CacahBaris int
}

// RekapPerMataUang menjumlahkan baris peserta menjadi rekap per mata uang.
//
// `tipe` adalah `T_PREMIUM_LIST.TYPE` - ia yang memilih cabang `PREMIUM` dan
// `BALANCE`, dan ia milik POLIS, bukan milik baris.
//
// ⛔ PEMBULATAN HANYA DI AKHIR. SELURUH keluaran dibulatkan ke empat angka
// SESUDAH seluruh baris dijumlah, persis `@divide(Param.X,1,4)` langkah
// 2.39-2.74 yang berdiri di luar perulangan. Membulatkan per baris lalu menjumlah menghasilkan angka
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
			for _, k := range KolomJumlahSummary {
				r.Jumlah[k] = apd.New(0, 0)
			}
			rekap[cur] = r
			urut = append(urut, cur)
		}
		r.CacahBaris++

		for _, k := range KolomJumlahSummary {
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
		// Pembulatan akhir, dan HANYA di sini - untuk SETIAP keluaran.
		bulatkan := func(nama string, v *apd.Decimal) (*apd.Decimal, error) {
			bulat := apd.New(0, 0)
			if _, err := ctx.Quantize(bulat, v, -AngkaDesimalRekap); err != nil {
				return nil, fmt.Errorf("models: membulatkan %s %s: %w", nama, cur, err)
			}
			return bulat, nil
		}
		var err error
		if r.Premium, err = bulatkan("PREMIUM", r.Premium); err != nil {
			return nil, err
		}
		if r.Commission, err = bulatkan("COMMISSION", r.Commission); err != nil {
			return nil, err
		}
		for _, k := range KolomJumlahSummary {
			if r.Jumlah[k], err = bulatkan(k, r.Jumlah[k]); err != nil {
				return nil, err
			}
		}
		if r.Balance, err = bulatkan("BALANCE", r.Balance); err != nil {
			return nil, err
		}
		keluar = append(keluar, *r)
	}
	return keluar, nil
}
