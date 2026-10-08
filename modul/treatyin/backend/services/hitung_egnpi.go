package services

import (
	"github.com/cockroachdb/apd/v3"

	inti "nusantarare/inti/backend"
)

// RUMUS TAB EGNPI (cabang NON-PROPORSIONAL) - disalin dari ekspor Pega.
//
// ---------------------------------------------------------------------
// ⭐ EMPAT ATURAN, DARI TIGA ACTIVITY
// ---------------------------------------------------------------------
//
//	Activity/SetAmountConversion.xml     .AmountIDR = .Amount x kurs
//	Activity/TreatyInEGNPIListValue.xml  jalankan di atas untuk SETIAP baris
//	                                     (tombol `Update EGNPI Value`)
//	Activity/TreatyInNPSetTotal.xml      param.type == "egnpi"
//	                                     (tombol `Update Total`)
//	Activity/TreatyInNonAddItem.xml      param.Type == "egnpi" (tombol `Add`)
//
// ---------------------------------------------------------------------
// ⛔ YANG TIDAK DAPAT DIBANGUN DI SINI, DAN SEBABNYA
// ---------------------------------------------------------------------
// `SetAmountConversion` punya DUA sumber kurs:
//
//	[1] `TreatyIn.CurrencyList` - grid `Rate of Exchange` kontrak ini
//	[2] bila [1] nol: laporan `BrowseTreatyExchangeYearly_RD` atas
//	    `ASM-FW-GISFW-Int-TREATYEXCHANGEYEARLY`, disaring
//	    `TreatyYear = @substring(@CurrentDateTime(),0,4)`
//
// ⚠️ Perhatikan tahunnya: cadangan [2] memakai **tahun kalender BERJALAN**,
// bukan tahun treaty kontraknya. Jadi dua kontrak bertahun treaty berbeda
// yang mata uangnya tidak ada di grid akan memakai kurs yang SAMA. Itu bunyi
// ekspornya; dinyatakan di sini supaya tidak terbaca sebagai kekeliruan port.
//
// ⛔ Cadangan [2] TIDAK dibangun: berkas ini murni menghitung - nol baca
// basis data - dan rute `/hitung/` menyatakan dirinya begitu. Bila mata
// uangnya nol di grid, `.AmountIDR` DIBIARKAN apa adanya dan satu pesan
// menyebut mata uangnya. Mengisinya nol akan terbaca sebagai "kursnya nol",
// padahal yang benar "kursnya tidak diketahui" - dua hal yang berbeda.

// Aksi tab EGNPI.
const (
	// AksiEgnpiKonversi - `SetAmountConversion` SATU baris; ekspor
	// memanggilnya saat `.Currency` atau `.Amount` berubah.
	AksiEgnpiKonversi = "konversi"
	// AksiEgnpiNilai - tombol `Update EGNPI Value` (`TreatyInEGNPIListValue`).
	AksiEgnpiNilai = "nilai"
	// AksiEgnpiTotal - tombol `Update Total` (`TreatyInNPSetTotal` type=egnpi).
	AksiEgnpiTotal = "total"
	// AksiEgnpiTambah - tombol `Add` (`TreatyInNonAddItem` Type=egnpi).
	AksiEgnpiTambah = "tambah"
	// AksiEgnpiHapus - tombol `Delete` pada baris grid.
	AksiEgnpiHapus = "hapus"
)

// PesanEgnpiTotalNol - total nol, pembagian dibatalkan.
//
// ⚠️ TEKSNYA dari `Param.errmsg` langkah 1 `TreatyInNPSetTotal`
// (`"Error Divide by Zero"`), BUKAN dari langkah yang memicunya: langkah 6
// ber-`Property-Set-Messages` dengan parameter KOSONG, jadi Pega memunculkan
// pesan hampa. Teks yang dideklarasikan dipakai di sini dengan sengaja -
// pesan hampa tidak memberi tahu pemakai apa pun.
const PesanEgnpiTotalNol = "Error Divide by Zero"

// MasukanEgnpi - satu aksi tab EGNPI beserta seluruh isian layarnya.
type MasukanEgnpi struct {
	Aksi  string    `json:"aksi"`
	EGNPI []EgnpiNP `json:"egnpi"`
	Kurs  []KursNP  `json:"kurs"`
	// Retensi - `TreatyIn.Retention`. ⭐ HANYA baris PERTAMA yang dipakai, dan
	// hanya oleh `Add`: ekspor menulis `TreatyIn.Retention(1).Currency`.
	Retensi []NilaiMataUang `json:"retensi"`
	// Indeks baris (mulai 0) untuk `konversi` dan `hapus`.
	Indeks int `json:"indeks"`
}

// HasilEgnpi - baris sesudah aksi, berikut ketiga total layar.
type HasilEgnpi struct {
	EGNPI []EgnpiNP `json:"egnpi"`
	// TotalEgnpiAmount - `TreatyIn.TotalEgnpiAmount`, panel `Total Amount in IDR`.
	TotalEgnpiAmount string `json:"TotalEgnpiAmount"`
	// TotalEgnpiProportion - `TreatyIn.TotalEgnpiProportion`, panel `Total Proportion %`.
	TotalEgnpiProportion string `json:"TotalEgnpiProportion"`
	// TotalEgnpiAmountNP - grid `Total EGNPI Amount` | `Value`, per mata uang.
	TotalEgnpiAmountNP []NilaiMataUang `json:"TotalEgnpiAmountNP"`
	Pesan              []string        `json:"pesan"`
}

// HitungEgnpi - bentuk ber-pelaku untuk handler.
func (l *Layanan) HitungEgnpi(p inti.Pelaku, m MasukanEgnpi) (HasilEgnpi, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return HasilEgnpi{}, err
	}
	return HitungEgnpi(m), nil
}

// HitungEgnpi menjalankan satu aksi tab EGNPI. Murni; masukan tidak diubah.
func HitungEgnpi(m MasukanEgnpi) HasilEgnpi {
	h := HasilEgnpi{EGNPI: append([]EgnpiNP{}, m.EGNPI...), Pesan: []string{}}
	ada := m.Indeks >= 0 && m.Indeks < len(h.EGNPI)

	switch m.Aksi {
	case AksiEgnpiKonversi:
		if ada {
			h.Pesan = append(h.Pesan, SetAmountConversion(&h.EGNPI[m.Indeks], m.Kurs)...)
		}
	case AksiEgnpiNilai:
		h.Pesan = append(h.Pesan, TreatyInEGNPIListValue(h.EGNPI, m.Kurs)...)
	case AksiEgnpiTambah:
		h.EGNPI = TambahBarisEgnpi(h.EGNPI, m.Retensi)
	case AksiEgnpiHapus:
		if ada {
			h.EGNPI = append(h.EGNPI[:m.Indeks], h.EGNPI[m.Indeks+1:]...)
		}
	}

	// ⭐ Ketiga total SELALU dihitung ulang, apa pun aksinya. Di Pega tombol
	// `Update Total` yang melakukannya, dan layar yang lupa menekannya
	// memperlihatkan total basi di sebelah baris yang baru berubah.
	tot, prop, perMataUang, pesan := NPSetTotalEgnpi(h.EGNPI)
	h.TotalEgnpiAmount, h.TotalEgnpiProportion, h.TotalEgnpiAmountNP = tot, prop, perMataUang

	// ⛔ CACAT YANG DIPERBAIKI 7 Oktober 2026 — pesan bocor ke aksi lain.
	//
	// Menekan `Add` melahirkan baris KOSONG, total menjadi nol, dan pagar
	// bagi-nol langsung meneriakkan `Error Divide by Zero` — di layar yang
	// pemakainya belum sempat mengetik apa pun. Pemilik proses
	// melaporkannya dengan satu tangkapan layar dan satu pertanyaan: "ini
	// kenapa?".
	//
	// ⭐ Di Pega pagar itu MILIK SATU TOMBOL. `TreatyInNPSetTotal` langkah 6
	// hanya berjalan lewat `Update Total` (`param.type == "egnpi"`);
	// `TreatyInNonAddItem` dan `TreatyInEGNPIListValue` nol memanggilnya.
	//
	// ⚠️ Yang BUKAN jalan keluar: berhenti menghitung ulang total pada aksi
	// lain. Totalnya memang harus ikut berubah — yang tidak boleh ikut
	// hanyalah TERIAKANNYA. Jadi hitungannya tetap, pesannya disaring.
	if m.Aksi == AksiEgnpiTotal {
		h.Pesan = append(h.Pesan, pesan...)
	}
	return h
}

// SetAmountConversion - `Activity/SetAmountConversion.xml`.
//
//	.AmountIDR = .Amount x kurs(.Currency)
//
// ⚠️ Kursnya diambil dari kecocokan TERAKHIR, bukan pertama: ekspor menapaki
// seluruh `CurrencyList` dan setiap kecocokan MENIMPA `local.ConvValue`.
// Nol `break` di sana, jadi nol `break` di sini.
func SetAmountConversion(b *EgnpiNP, kurs []KursNP) []string {
	conv := apd.New(0, 0)
	for _, k := range kurs {
		if k.Currency == b.Currency {
			conv = angka(k.Conversion)
		}
	}
	if conv.Sign() > 0 {
		b.AmountIDR = teks(kali(angka(b.Amount), conv))
		return nil
	}
	// ⛔ Cadangan `BrowseTreatyExchangeYearly_RD` tidak dibangun - lihat
	// kepala berkas. `.AmountIDR` dibiarkan, tidak dinolkan.
	if b.Currency == "" {
		return nil
	}
	return []string{"Rate of Exchange " + b.Currency + " not found"}
}

// TreatyInEGNPIListValue - `Activity/TreatyInEGNPIListValue.xml`: satu
// langkah, `SetAmountConversion` untuk SETIAP baris EGNPI.
func TreatyInEGNPIListValue(rows []EgnpiNP, kurs []KursNP) []string {
	var pesan []string
	for i := range rows {
		pesan = append(pesan, SetAmountConversion(&rows[i], kurs)...)
	}
	return pesan
}

// NPSetTotalEgnpi - `Activity/TreatyInNPSetTotal.xml` cabang
// `param.type == "egnpi"` (langkah 4-7):
//
//	[4] kosongkan TotalEgnpiAmountNP, TotalEgnpiAmount, TotalEgnpiProportion
//	[5] TotalEgnpiAmount = SIGMA .AmountIDR
//	[6] bila TotalEgnpiAmount == 0 -> BERHENTI (pagar bagi-nol)
//	[7] per baris: .Proportion = @divide(.AmountIDR, TotalEgnpiAmount, 20)*100
//	    TotalEgnpiProportion += .Proportion
//	    TotalEgnpiAmountNP: SIGMA .Amount per MATA UANG
//
// ⛔ Perhatikan langkah 7 menjumlah `.Amount` (mata uang asli) sementara
// langkah 5 menjumlah `.AmountIDR`. Dua kolom berbeda di baris yang sama -
// menyeragamkannya akan membuat grid per mata uang memperlihatkan rupiah.
//
// ⚠️ Skala 20 bukan hiasan: `@divide(...,20)` itulah yang membuat jumlah
// seluruh proporsi jatuh TEPAT di 100. Memendekkannya membuat totalnya
// meleset di desimal belakang.
func NPSetTotalEgnpi(rows []EgnpiNP) (total, proporsi string, perMataUang []NilaiMataUang, pesan []string) {
	perMataUang = []NilaiMataUang{}

	jumlahIDR := apd.New(0, 0)
	for _, r := range rows {
		jumlahIDR = tambah(jumlahIDR, angka(r.AmountIDR))
	}
	total = teks(jumlahIDR)

	if jumlahIDR.IsZero() {
		// ⛔ Pagar langkah 6: tanpa ini langkah 7 membagi dengan nol.
		// Proporsi baris DIBIARKAN - menolkannya akan menghapus angka yang
		// masih sah dari hitungan sebelumnya.
		return total, "0", perMataUang, []string{PesanEgnpiTotalNol}
	}

	jumlahProp := apd.New(0, 0)
	for i := range rows {
		p := kali(bagiBulatDes(angka(rows[i].AmountIDR), jumlahIDR, 20), apd.New(100, 0))
		rows[i].Proportion = teks(p)
		jumlahProp = tambah(jumlahProp, p)
		perMataUang = tambahPerMataUang(perMataUang, rows[i].Currency, rows[i].CurrencyID, angka(rows[i].Amount))
	}
	return total, teks(jumlahProp), perMataUang, nil
}

// TambahBarisEgnpi - `Activity/TreatyInNonAddItem.xml` cabang
// `param.Type == "egnpi"`: satu baris kosong yang mata uangnya MEWARISI
// baris PERTAMA tab Maximum Retention (`TreatyIn.Retention(1)`).
//
// ⚠️ Nol baris retensi berarti baris baru tanpa mata uang. Itu bunyi
// ekspornya, dan ia tidak dikarang menjadi "IDR" di sini: mata uang karangan
// akan ikut ke hitungan konversi tanpa seorang pun memilihnya.
func TambahBarisEgnpi(rows []EgnpiNP, retensi []NilaiMataUang) []EgnpiNP {
	baru := EgnpiNP{}
	if len(retensi) > 0 {
		baru.Currency = retensi[0].Currency
		baru.CurrencyID = retensi[0].CurrencyID
	}
	return append(rows, baru)
}
