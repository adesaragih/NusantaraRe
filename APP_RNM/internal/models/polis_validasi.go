package models

// Gerbang simpan penawaran polis - `ProtectAccept`, tiket 01 bagian 2.
//
// Untuk apa berkas ini: tiga belas pemeriksaan HIDUP yang
// `Activity/ProtectAccept.xml` jalankan sebelum sebuah penawaran boleh maju.
// Seluruhnya MURNI.
//
// `[terverifikasi]` 13, dihitung dua cara atas berkas pecahan:
//
//	B  grep -c "<PropertiesValue>Local.Err + "   -> 19, dikurangi 6 milik
//	   langkah ter-remark (4.1, 6.1-6.5)          -> 13
//	A  pengurai XML: langkah daun ber-`Local.Err +` yang dirinya maupun
//	   induknya tidak ber-`//`                    -> 13 hidup, 6 mati
//
// ⛔ PESANNYA VERBATIM, termasuk yang janggal. `System Reinsurance can't null`
// adalah kalimat yang pemakai sistem lama hafal. Memperbaiki ejaannya berarti
// layar baru berbicara dengan kalimat yang tidak pernah ada, dan orang yang
// mencari kalimat lamanya tidak akan menemukannya.
//
// Cara rule itu bekerja, dibaca sebagai pohon 28-09-2026 (beserta
// `pyStepsBlockName`, sensus remark GILIRAN-12):
//
//	b333  `Page-Clear-Messages`
//	b430  `ProtectLife.CARI1 = 0`          <- penanda "belum ada galat"
//	3     b624 `EmailTypePL==1` WhenFalse=6 -> KELUAR (lihat celah di bawah)
//	4     b1207 `Position=="Offer"`   - 4.2 TypeCeding, 4.3 COB
//	      4.1 NoOffer - ⛔ TER-REMARK (`//` b720)
//	5     b2288 `Position=="Premium"` - 5.1 Type ... 5.6 SOB
//	6     per baris detail (sum insured, sum reasured, rate, gross, net) -
//	      ⛔ TER-REMARK (`//` b2340; 6.3 pun b2756)
//	7     per mata uang: `PREMIUM<=0`, `BALANCE<=0`
//	8     b3792 `Premium && CurrencyList kosong` - lihat celah di bawah
//	9-11  `CurrencyList(1)` saat Premium (umum, TP, TR)
//	b4462 `Page-Set-Messages` berprasyarat b4563 `ProtectLife.CARI1==1`
//
// ⛔ RALAT 28-09-2026 (sensus remark GILIRAN-12): ronde pertama meniru 4.1
// (`Please choose no offer !`) dan kelima pemeriksaan per baris langkah 6 -
// keduanya ter-remark, jadi sistem lama tidak pernah menolak karenanya.
// Keduanya DIBUANG. Ronde itu juga memasang gerbang posisi hanya pada COB dan
// SOB; di rule, b1207 dan b2288 menggerbangi SELURUH langkah 4 dan 5 - kini
// ditiru begitu.
//
// ⚠️ CELAH TERCATAT (bukan ditiru): langkah 3 keluar tanpa satu pemeriksaan
// pun bila `pyWorkPage.EmailTypePL != 1`, dan langkah 8 menyalakan penanda
// galat tanpa kalimat bila daftar mata uang kosong di posisi Premium.
// `EmailTypePL` tidak punya sumber di aplikasi ini. Dan gerbang ini BELUM
// TERSAMBUNG ke rute mana pun - pemanggilnya hanya uji.
//
// ⚠️ Jadi rule itu MENGUMPULKAN seluruh galat lalu menampilkannya sekali -
// bukan berhenti pada yang pertama. Ditiru apa adanya: pemakai yang harus
// menekan Simpan tujuh kali untuk menemukan tujuh isian kosong akan berhenti
// memakai layarnya.
//
// Dibaca sesudah: polis_penawaran.go (tahap dan posisinya).

import (
	"strings"

	"github.com/cockroachdb/apd/v3"
)

// Pesan galat penawaran - VERBATIM `Local.Err` di `ProtectAccept.xml`.
const (
	PesanTypeCedingKosong   = "System Reinsurance can't null" // b943
	PesanBusinessCodeKosong = "COB can't null"                // b1105
	PesanTypeKosong         = "Type can't null"               // b1368
	PesanProductNameKosong  = "Product Name can't null"       // b1530
	PesanProRateTypeKosong  = "Prorate Type can't null"       // b1692
	PesanMarketingKosong    = "Marketing can't null"          // b1854
	PesanBelumUnggahCSV     = "Please upload CSV !"           // b2016
	PesanSOBKosong          = "SOB can't null"                // b2186
	PesanPremiumNol         = "Premium is 0"                  // b3394
	PesanBalanceNol         = "Balance is 0"                  // b3556, b3906, b4068, b4300
)

// PenawaranPolis adalah medan header yang `ProtectAccept` periksa.
//
// ⚠️ Seluruhnya TEKS, dan itu disengaja: yang diperiksa rule aslinya adalah
// `== ""`, bukan nilai. Kolom yang kosong dan kolom yang bernilai nol adalah
// dua keadaan berbeda (ADR-U-0027), dan teks dapat membedakannya.
type PenawaranPolis struct {
	// Posisi - `Offer` atau `Premium`; lihat PosisiOffer/PosisiPremium.
	Posisi           string
	TypeCeding       string
	BusinessCode     string
	Type             string
	ProductName      string
	ProRateType      string
	MarketingName    string
	SourceOfBusiness string
	// CacahDetail - `@LengthOfPageList(PremiumListSummary.PremiumListDetail)`
	// (5.5). Isi barisnya tidak diperiksa: langkah 6 ter-remark.
	CacahDetail int
	// MataUang adalah `PremiumListSummary.CurrencyList`.
	MataUang []BarisMataUangPenawaran
}

// BarisMataUangPenawaran adalah satu baris `CurrencyList`.
type BarisMataUangPenawaran struct {
	Premium Money
	Balance Money
	// DiscountPremiumRetro dipakai cabang `Type=="TP"` b4190.
	DiscountPremiumRetro Money
	// DiscountPremiumRefundRetro dipakai cabang `Type=="TR"` b4422.
	DiscountPremiumRefundRetro Money
}

// ValidasiPenawaran menjalankan seluruh pemeriksaan `ProtectAccept` yang hidup.
//
// Mengembalikan SELURUH pesan, berurutan seperti rule aslinya. Daftar kosong
// berarti penawaran boleh maju.
//
// ⛔ Mengumpulkan, bukan berhenti pada yang pertama - lihat kepala berkas.
func ValidasiPenawaran(p PenawaranPolis) []string {
	var galat []string
	tambah := func(pesan string) { galat = append(galat, pesan) }

	// Langkah 4 - HANYA di posisi Offer (b1207 menggerbangi seluruh langkah).
	if p.Posisi == PosisiOffer {
		if kosong(p.TypeCeding) {
			tambah(PesanTypeCedingKosong) // 4.2 b993 -> b943
		}
		if kosong(p.BusinessCode) {
			tambah(PesanBusinessCodeKosong) // 4.3 b1155 -> b1105
		}
	}
	// Langkah 5 - HANYA di posisi Premium (b2288 menggerbangi seluruh langkah).
	if p.Posisi == PosisiPremium {
		if kosong(p.Type) {
			tambah(PesanTypeKosong) // 5.1 b1418 -> b1368
		}
		if kosong(p.ProductName) {
			tambah(PesanProductNameKosong) // 5.2 b1580 -> b1530
		}
		if kosong(p.ProRateType) {
			tambah(PesanProRateTypeKosong) // 5.3 b1742 -> b1692
		}
		if kosong(p.MarketingName) {
			tambah(PesanMarketingKosong) // 5.4 b1904 -> b1854
		}
		// 5.5 b2066 `@LengthOfPageList(...PremiumListDetail) = 0`.
		if p.CacahDetail == 0 {
			tambah(PesanBelumUnggahCSV) // b2016
		}
		if kosong(p.SourceOfBusiness) {
			tambah(PesanSOBKosong) // 5.6 b2236 -> b2186
		}
	}

	// Langkah 7 - per mata uang, TANPA gerbang posisi.
	for _, m := range p.MataUang {
		if nolAtauKurang(m.Premium.Amount) {
			tambah(PesanPremiumNol) // 7.1 b3444 -> b3394
		}
		if nolAtauKurang(m.Balance.Amount) {
			tambah(PesanBalanceNol) // 7.2 b3606 -> b3556
		}
	}
	// Langkah 9-11 - `CurrencyList(1)`, hanya saat Premium.
	if p.Posisi == PosisiPremium && len(p.MataUang) > 0 {
		pertama := p.MataUang[0]
		// b3958 `CurrencyList(1).BALANCE<=0`
		if nolAtauKurang(pertama.Balance.Amount) {
			tambah(PesanBalanceNol) // b3906
		}
		// b4121 `Type=="TP"` + b4144/b4167/b4190
		if p.Type == "TP" && (nolAtauKurang(pertama.Balance.Amount) ||
			nolAtauKurang(pertama.Premium.Amount) ||
			kurangDariNol(pertama.DiscountPremiumRetro.Amount)) {
			tambah(PesanBalanceNol) // b4068
		}
		// b4353 `Type=="TR"` + b4376/b4399/b4422
		if p.Type == "TR" && (nolAtauKurang(pertama.Balance.Amount) ||
			nolAtauKurang(pertama.Premium.Amount) ||
			kurangDariNol(pertama.DiscountPremiumRefundRetro.Amount)) {
			tambah(PesanBalanceNol) // b4300
		}
	}
	return galat
}

// PenawaranBolehMaju menjawab gerbang b4563 `ProtectLife.CARI1==1` terbalik.
func PenawaranBolehMaju(p PenawaranPolis) bool {
	return len(ValidasiPenawaran(p)) == 0
}

// GabungPesanPenawaran merangkai pesan seperti `Local.Err` merangkainya.
//
// ⛔ Pemisahnya BARIS BARU, sebagaimana `Local.Err + "\n" + <pesan>`.
func GabungPesanPenawaran(pesan []string) string {
	return strings.Join(pesan, "\n")
}

// kosong membaca medan teks apa adanya - padanan `== ""` di rule.
func kosong(s string) bool { return strings.TrimSpace(s) == "" }

// Pembanding desimal untuk gerbang `ProtectAccept`.
//
// ⛔ KOSONG DIPERLAKUKAN SEBAGAI NOL DI SINI, dan HANYA di sini.
// ADR-U-0027 menyatakan kosong bukan nol, dan itu tetap berlaku di seluruh
// model ini - `Money.Kosong()` ada justru untuk membedakannya. Tetapi rule
// yang ditiru menulis `.PREMIUM<=0`, dan di Pega perbandingan itu bernilai
// BENAR untuk properti yang belum diisi.
//
// Menolak menyamakannya di sini berarti isian yang dikosongkan pemakai lolos
// gerbang yang di sistem lama menahannya - selisih yang justru merugikan.
// Penyimpangan karena itu dinyatakan, disempitkan ke berkas ini, dan tidak
// merembet ke `Money` sendiri.

// nolAtauKurang meniru `x == 0 || x < 0`.
func nolAtauKurang(d *apd.Decimal) bool {
	return d == nil || d.IsZero() || d.Negative
}

// kurangDariNol meniru `x < 0` - nol dan kosong LOLOS.
//
// ⚠️ Berbeda dengan yang di atas: `DISCOUNT_PREMIUM_RETRO<0` b4190
// hanya menolak yang negatif. Diskon nol adalah diskon yang sah.
func kurangDariNol(d *apd.Decimal) bool {
	return d != nil && d.Negative && !d.IsZero()
}
