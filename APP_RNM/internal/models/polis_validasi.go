package models

// Gerbang simpan penawaran polis - `ProtectAccept`, tiket 01 bagian 2.
//
// Untuk apa berkas ini: dua puluh pemeriksaan yang `Activity/ProtectAccept.xml`
// jalankan sebelum sebuah penawaran boleh maju. Seluruhnya MURNI.
//
// ⛔ PESANNYA VERBATIM, termasuk yang janggal. `System Reinsurance can't null`
// dan `Please choose no offer !` - beserta spasi sebelum tanda serunya - adalah
// kalimat yang pemakai sistem lama hafal. Memperbaiki ejaannya berarti layar
// baru berbicara dengan kalimat yang tidak pernah ada, dan orang yang mencari
// kalimat lamanya tidak akan menemukannya.
//
// Cara rule itu bekerja, dibaca sebagai pohon 28-09-2026:
//
//	b333  `Page-Clear-Messages`
//	b430  `ProtectLife.CARI1 = 0`          <- penanda "belum ada galat"
//	...   tiap pemeriksaan: prasyaratnya MUNCUL SESUDAH parameternya di DOM,
//	      `WhenTrue=2` LANJUT (galat dicatat), `WhenFalse=3` LEWATI
//	      `ProtectLife.CARI1 = 1` dan `Local.Err = Local.Err + "\n" + <pesan>`
//	b4462 `Page-Set-Messages` berprasyarat b4563 `ProtectLife.CARI1==1`
//
// ⚠️ Jadi rule itu MENGUMPULKAN seluruh galat lalu menampilkannya sekali -
// bukan berhenti pada yang pertama. Ditiru apa adanya: pemakai yang harus
// menekan Simpan tujuh kali untuk menemukan tujuh isian kosong akan berhenti
// memakai layarnya.
//
// Dibaca sesudah: polis_penawaran.go (tahap dan posisinya).

import (
	"strconv"
	"strings"

	"github.com/cockroachdb/apd/v3"
)

// Pesan galat penawaran - VERBATIM `Local.Err` di `ProtectAccept.xml`.
const (
	PesanNoOfferKosong      = "Please choose no offer !"      // b781
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

// Pesan per baris detail - `" +local.idx+ "` di tengahnya.
//
// ⚠️ `local.idx` NOMOR URUT baris, 1-based, sebagaimana `.pxListSubscript` -
// bukan pengenal barisnya. Nomor urut dipertahankan supaya pesannya dapat
// dibandingkan dengan sistem lama.
const (
	pesanSumInsuredNol   = "Sum insured number %s is 0"   // b2493
	pesanSumReasuredNol  = "Sum reasured number %s is 0"  // b2655
	pesanRateNol         = "Rate number %s is 0"          // b2817
	pesanGrossPremiumNol = "Gross premium number %s is 0" // b2979
	pesanNetPremiumNol   = "Net premium number %s is 0"   // b3141
)

// PenawaranPolis adalah medan header yang `ProtectAccept` periksa.
//
// ⚠️ Seluruhnya TEKS, dan itu disengaja: yang diperiksa rule aslinya adalah
// `== ""`, bukan nilai. Kolom yang kosong dan kolom yang bernilai nol adalah
// dua keadaan berbeda (ADR-U-0027), dan teks dapat membedakannya.
type PenawaranPolis struct {
	// Posisi - `Offer` atau `Premium`; lihat PosisiOffer/PosisiPremium.
	Posisi           string
	NoOffer          string
	TypeCeding       string
	BusinessCode     string
	Type             string
	ProductName      string
	ProRateType      string
	MarketingName    string
	SourceOfBusiness string
	// Detail adalah baris `PremiumListSummary.PremiumListDetail`.
	Detail []BarisDetailPenawaran
	// MataUang adalah `PremiumListSummary.CurrencyList`.
	MataUang []BarisMataUangPenawaran
}

// BarisDetailPenawaran adalah satu baris detail yang diperiksa.
//
// ⛔ Kelima nilainya `Money`/`Ratio`, bukan float. Yang diperiksa rule aslinya
// `== 0` dan `< 0`; membandingkan float dengan nol adalah cara termudah
// meloloskan angka yang sebenarnya 0,000001 (ADR-U-0003).
type BarisDetailPenawaran struct {
	SumInsured   Money
	SumReasured  Money
	Rate         Ratio
	GrossPremium Money
	NetPremium   Money
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

// ValidasiPenawaran menjalankan seluruh pemeriksaan `ProtectAccept`.
//
// Mengembalikan SELURUH pesan, berurutan seperti rule aslinya. Daftar kosong
// berarti penawaran boleh maju.
//
// ⛔ Mengumpulkan, bukan berhenti pada yang pertama - lihat kepala berkas.
func ValidasiPenawaran(p PenawaranPolis) []string {
	var galat []string
	tambah := func(pesan string) { galat = append(galat, pesan) }

	// Pemeriksaan header. Urutannya urutan langkah di rule.
	if kosong(p.NoOffer) {
		tambah(PesanNoOfferKosong) // b831 -> b781
	}
	if kosong(p.TypeCeding) {
		tambah(PesanTypeCedingKosong) // b993 -> b943
	}
	// ⛔ DUA syarat, dan keduanya harus benar: b1155 `BusinessCode==""` DAN
	// b1207 `Position=="Offer"`. COB hanya wajib di layar penawaran.
	if kosong(p.BusinessCode) && p.Posisi == PosisiOffer {
		tambah(PesanBusinessCodeKosong) // b1105
	}
	if kosong(p.Type) {
		tambah(PesanTypeKosong) // b1418 -> b1368
	}
	if kosong(p.ProductName) {
		tambah(PesanProductNameKosong) // b1580 -> b1530
	}
	if kosong(p.ProRateType) {
		tambah(PesanProRateTypeKosong) // b1742 -> b1692
	}
	if kosong(p.MarketingName) {
		tambah(PesanMarketingKosong) // b1904 -> b1854
	}
	// b2066 `@LengthOfPageList(...PremiumListDetail) = 0`.
	if len(p.Detail) == 0 {
		tambah(PesanBelumUnggahCSV) // b2016
	}
	// b2236 `SourceOfBusiness==""` DAN b2288 `Position=="Premium"`.
	if kosong(p.SourceOfBusiness) && p.Posisi == PosisiPremium {
		tambah(PesanSOBKosong) // b2186
	}

	// Pemeriksaan per baris detail. `local.idx` 1-based.
	for i, b := range p.Detail {
		idx := strconv.Itoa(i + 1)
		// b2543 `SUM_INSURED==0 || SUM_INSURED<0`
		if nolAtauKurang(b.SumInsured.Amount) {
			tambah(sisipkanIndeks(pesanSumInsuredNol, idx))
		}
		// b2705 `SUM_REASURED==0` - HANYA nol, tanpa "atau kurang".
		//
		// ⚠️ Dipertahankan apa adanya. Rule aslinya memang tidak menguji
		// negatif di sini sedangkan di empat pemeriksaan lain ia menguji;
		// "memperbaikinya" berarti menolak baris yang sistem lama terima.
		if nolSaja(b.SumReasured.Amount) {
			tambah(sisipkanIndeks(pesanSumReasuredNol, idx))
		}
		// b2867 `RATE==0` - juga hanya nol.
		if nolSaja(b.Rate.Value) {
			tambah(sisipkanIndeks(pesanRateNol, idx))
		}
		// b3029 `GROSS_PREMIUM==0 || GROSS_PREMIUM<0`
		if nolAtauKurang(b.GrossPremium.Amount) {
			tambah(sisipkanIndeks(pesanGrossPremiumNol, idx))
		}
		// b3191 `NET_PREMIUM==0 || NET_PREMIUM<0`
		if nolAtauKurang(b.NetPremium.Amount) {
			tambah(sisipkanIndeks(pesanNetPremiumNol, idx))
		}
	}

	// Pemeriksaan tingkat mata uang - seluruhnya hanya saat `Premium`.
	//
	// ⚠️ b3444 `PREMIUM<=0` dan b3606 `BALANCE<=0` TIDAK bersyarat posisi di
	// rule aslinya; yang bersyarat `Position=="Premium"` adalah pemeriksaan
	// atas `CurrencyList(1)` (b3958, b4144, b4376). Perbedaan itu ditiru apa
	// adanya, dan ia terbaca dari dua kelompok di bawah.
	for _, m := range p.MataUang {
		if nolAtauKurang(m.Premium.Amount) {
			tambah(PesanPremiumNol) // b3444 -> b3394
		}
		if nolAtauKurang(m.Balance.Amount) {
			tambah(PesanBalanceNol) // b3606 -> b3556
		}
	}
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

// sisipkanIndeks menaruh nomor urut ke tengah pesan.
func sisipkanIndeks(pola, idx string) string {
	return strings.Replace(pola, "%s", idx, 1)
}

// Pembanding desimal untuk gerbang `ProtectAccept`.
//
// ⛔ KOSONG DIPERLAKUKAN SEBAGAI NOL DI SINI, dan HANYA di sini.
// ADR-U-0027 menyatakan kosong bukan nol, dan itu tetap berlaku di seluruh
// model ini - `Money.Kosong()` ada justru untuk membedakannya. Tetapi rule
// yang ditiru menulis `.SUM_INSURED==0`, dan di Pega perbandingan itu bernilai
// BENAR untuk properti yang belum diisi.
//
// Menolak menyamakannya di sini berarti isian yang dikosongkan pemakai lolos
// gerbang yang di sistem lama menahannya - selisih yang justru merugikan.
// Penyimpangan karena itu dinyatakan, disempitkan ke berkas ini, dan tidak
// merembet ke `Money` sendiri.

// nolSaja meniru `x == 0` - tanpa "atau kurang".
func nolSaja(d *apd.Decimal) bool {
	return d == nil || d.IsZero()
}

// nolAtauKurang meniru `x == 0 || x < 0`.
func nolAtauKurang(d *apd.Decimal) bool {
	return d == nil || d.IsZero() || d.Negative
}

// kurangDariNol meniru `x < 0` - nol dan kosong LOLOS.
//
// ⚠️ Berbeda dengan kedua di atas: `DISCOUNT_PREMIUM_RETRO<0` b4190
// hanya menolak yang negatif. Diskon nol adalah diskon yang sah.
func kurangDariNol(d *apd.Decimal) bool {
	return d != nil && d.Negative && !d.IsZero()
}
