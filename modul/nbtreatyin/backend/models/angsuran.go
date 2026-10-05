package models

// Untuk apa berkas ini: JADWAL ANGSURAN dan PENYEBARAN RISIKO - port
// `FillPaymentInstallment`, `SetValidateInstallment_Act`,
// `CountPctInstallment_Act`, `CountSpreading_Act` (tiket 07, 13; AC 79).
//
// Daftar yang disentuh: `PolicyTreatyIn.ListInstallment` (anggota
// InstallmentNo, DueDate, InstallmentPercentage, Premium, PaymentTotal) dan
// `PolicyTreatyIn.SpreadingRiskList` (TreatyType, TreatyName, SharePercentage,
// ClaimPercentage, PremiumSpreaded, ClaimSpreaded).

import (
	"strconv"
	"strings"
	"time"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/utils"
)

// Jalur kedua daftar di halaman.
const (
	DaftarAngsuran  = HalamanPolis + ".ListInstallment"
	DaftarSpreading = HalamanPolis + ".SpreadingRiskList"
	DaftarUsulan    = HalamanPolis + ".SuggestList"
)

// presisiBagiRata - pembagian rata `100 / jumlah baris` spreading NB dibulatkan
// pada 10 desimal (spec-penyimpanan ID-28, AC 36 `[terverifikasi]`).
const presisiBagiRata = 10

// PesanAngsuranLebih100 - VERBATIM `SetValidateInstallment_Act` langkah 2.
const PesanAngsuranLebih100 = "Installment percentage is more than 100"

func angkaBaris(k *kalkulator, b Baris, m string) *apd.Decimal {
	d, err := AngkaTeks(m, b[m])
	if err != nil {
		k.catat(err)
		return apd.New(0, 0)
	}
	return d
}

// FillPaymentInstallment = `Activity/FillPaymentInstallment` (Param.Installment
// = `.Installment`). `sekarang` menggantikan `@CurrentDateTime()` supaya
// fungsi ini murni.
//
// Langkah 4 dan 5 rule berlabel `//` - DINONAKTIFKAN, tidak diport.
func FillPaymentInstallment(h *Halaman, sekarang time.Time) error {
	k := &kalkulator{}
	p := polis{h, k}
	// langkah 1: Property-Remove .ListInstallment
	h.SetelDaftar(DaftarAngsuran, nil)
	// langkah 2
	no := 0
	pct := apd.New(0, 0)
	premi := p.n("BalanceDueTo")
	n, err := strconv.Atoi(strings.TrimSpace(p.teks("Installment")))
	if err != nil || n <= 0 {
		// Batas ulang `Param.Installment-1` < 0: nol putaran.
		return k.err
	}
	jumlah := apd.New(int64(n), 0)
	bagian := k.BagiBulat(seratus, jumlah, 4) // @Math.divide(100, Param.Installment, 4)
	var daftar []Baris
	for i := 0; i < n; i++ { // langkah 3: ulang 0..Installment-1
		no++                                                    // 3.1
		pct = k.Tambah(pct, bagian)                             // 3.2
		b := Baris{}                                            // .ListInstallment(<APPEND>)
		if no == n && k.Tambah(pct, bagian).Cmp(seratus) != 0 { // 3.3
			b["InstallmentPercentage"] = formatAngka(k.Tambah(bagian, k.Kurang(seratus, pct)))
		}
		if no < n { // 3.4
			b["InstallmentPercentage"] = formatAngka(bagian)
		}
		// 3.5
		b["InstallmentNo"] = strconv.Itoa(no)
		persen := angkaBaris(k, b, "InstallmentPercentage")
		nilai := k.BagiBulat(k.Kali(premi, persen), seratus, 4)
		b["Premium"] = formatAngka(nilai)
		b["DueDate"] = utils.FormatTanggal(sekarang)
		b["PaymentTotal"] = formatAngka(nilai)
		daftar = append(daftar, b)
	}
	if k.err != nil {
		return k.err
	}
	h.SetelDaftar(DaftarAngsuran, daftar)
	return nil
}

// SetValidateInstallment = `Activity/SetValidateInstallment_Act`.
//
// ⛔ LANGKAH 1 (`Page-Clear-Messages` atas `pyWorkPage.PolicyTreatyIn`) TIDAK
// DIBANGUN - AC 78 `[terbuka]`. Aktivitas ini dipanggil di ujung
// `CountOGPONP_Act` (langkah 10), SESUDAH `CountResult*` memasang `ErrorMsg1`;
// langkah 1 itulah yang menghapus pesan validasi yang baru saja terpasang.
// Penggolongannya di tiket 12 (RALAT) dan INVENTARIS bab 6.
func SetValidateInstallment(h *Halaman) error {
	k := &kalkulator{}
	p := polis{h, k}
	balance := p.n("BalanceDueTo")
	total := apd.New(0, 0)
	daftar := h.AmbilDaftar(DaftarAngsuran)
	for _, b := range daftar { // langkah 3
		persen := angkaBaris(k, b, "InstallmentPercentage")
		// 3.1: .Premium = BalanceDueTo* @divide(.InstallmentPercentage,100,4); .PaymentTotal sama
		nilai := k.Kali(balance, k.BagiBulat(persen, seratus, 4))
		b["Premium"] = formatAngka(nilai)
		b["PaymentTotal"] = formatAngka(nilai)
		total = k.Tambah(total, persen) // 3.2
	}
	if k.err != nil {
		return k.err
	}
	if total.Cmp(seratus) > 0 { // langkah 4
		h.TambahPesan(HalamanPolis, PesanAngsuranLebih100)
	}
	return nil
}

// CountPctInstallment = `Activity/CountPctInstallment_Act` (Param.idx, berbasis 1).
func CountPctInstallment(h *Halaman, idx int) error {
	k := &kalkulator{}
	p := polis{h, k}
	daftar := h.AmbilDaftar(DaftarAngsuran)
	if idx < 1 || idx > len(daftar) {
		return nil
	}
	b := daftar[idx-1]
	// .InstallmentPercentage = @Math.divide(@toDecimal(.Premium),@toDecimal(.BalanceDueTo),4)*100
	b["InstallmentPercentage"] = formatAngka(k.Kali(k.BagiBulat(angkaBaris(k, b, "Premium"), p.n("BalanceDueTo"), 4), seratus))
	// .PaymentTotal = .Premium
	b["PaymentTotal"] = b["Premium"]
	return k.err
}

// CountSpreading = `Activity/CountSpreading_Act` (Param.Index berbasis 1).
//
// ⛔ RALAT audit silang P3 (7.4): langkah 1-3 rule ini berlabel `//`
// (dinonaktifkan) - port lama menjalankannya untuk baris `Param.Index`. Tidak ada
// beda hasil (langkah 4.1 menulis keempat medan setiap baris dengan rumus yang
// sama), tetapi langkah nonaktif tidak diport: `idx` hanya parameter tanda
// tangan aksi sel (`.pxListSubscript`), tidak dibaca.
//
//	1 `//`  .SpreadingRiskList(Param.Index).SharePercentage == "" -> 100/@LengthOfPageList(..)
//	2 `//`  .ClaimPercentage == "" -> .SharePercentage
//	3 `//`  .PremiumSpreaded / .ClaimSpreaded baris Param.Index
//	4       setiap baris: 4.1 isi %Share kosong, PremiumSpreaded, ClaimSpreaded; 4.2 jumlah
//	5       empat total
//
// ⛔ Langkah 6 (`Call BreakDownSpreading_Act`) TIDAK DIBANGUN -
// `[keputusan work owner]` 23-09-2026 butir 3b: "Itu tidak usah di migrasi
// perhitungannya" (KEPUTUSAN-RONDE-12 butir 3/3b).
func CountSpreading(h *Halaman, _ int) error {
	k := &kalkulator{}
	p := polis{h, k}
	daftar := h.AmbilDaftar(DaftarSpreading)
	if len(daftar) == 0 {
		return nil
	}
	panjang := apd.New(int64(len(daftar)), 0)
	net := p.n("NetPremium")
	// (Primary.ExcessLoss + Primary.Claim - Primary.SalvageValue)
	klaim := k.Kurang(k.Tambah(p.n("ExcessLoss"), p.n("Claim")), p.n("SalvageValue"))
	// langkah 4: setiap baris
	shareKlaim, nilaiKlaim := apd.New(0, 0), apd.New(0, 0)
	sharePremi, nilaiPremi := apd.New(0, 0), apd.New(0, 0)
	for _, b := range daftar {
		if b["SharePercentage"] == "" { // @if(.SharePercentage=="",(100/...),.SharePercentage)
			b["SharePercentage"] = formatAngka(k.BagiBulat(seratus, panjang, presisiBagiRata))
		}
		if b["ClaimPercentage"] == "" { // @if(.ClaimPercentage == "",.SharePercentage,.ClaimPercentage)
			b["ClaimPercentage"] = b["SharePercentage"]
		}
		b["PremiumSpreaded"] = formatAngka(k.Kali(net, k.BagiBulat(angkaBaris(k, b, "SharePercentage"), seratus, 10)))
		b["ClaimSpreaded"] = formatAngka(k.Kali(klaim, k.BagiBulat(angkaBaris(k, b, "ClaimPercentage"), seratus, 10)))
		// 4.2
		shareKlaim = k.Tambah(angkaBaris(k, b, "ClaimPercentage"), shareKlaim)
		nilaiKlaim = k.Tambah(angkaBaris(k, b, "ClaimSpreaded"), nilaiKlaim)
		sharePremi = k.Tambah(angkaBaris(k, b, "SharePercentage"), sharePremi)
		nilaiPremi = k.Tambah(angkaBaris(k, b, "PremiumSpreaded"), nilaiPremi)
	}
	// langkah 5
	p.setel("TotalSharePercentagePremium", sharePremi)
	p.setel("TotalPremium", nilaiPremi)
	p.setel("TotalSharePercentageClaim", shareKlaim)
	p.setel("TotalClaim", nilaiKlaim)
	return k.err
}

// HitungTotalSpreading = penjumlahan `CountSpreading_Act` langkah 4.2 dan 5
// SAJA - tanpa menghitung ulang baris. Dipakai saat membaca polis tersimpan:
// keempat total adalah TURUNAN baris SpreadingRiskList dan tidak disimpan
// (models/katalog.go).
func HitungTotalSpreading(h *Halaman) error {
	k := &kalkulator{}
	p := polis{h, k}
	shareKlaim, nilaiKlaim := apd.New(0, 0), apd.New(0, 0)
	sharePremi, nilaiPremi := apd.New(0, 0), apd.New(0, 0)
	for _, b := range h.AmbilDaftar(DaftarSpreading) {
		shareKlaim = k.Tambah(angkaBaris(k, b, "ClaimPercentage"), shareKlaim)
		nilaiKlaim = k.Tambah(angkaBaris(k, b, "ClaimSpreaded"), nilaiKlaim)
		sharePremi = k.Tambah(angkaBaris(k, b, "SharePercentage"), sharePremi)
		nilaiPremi = k.Tambah(angkaBaris(k, b, "PremiumSpreaded"), nilaiPremi)
	}
	p.setel("TotalSharePercentagePremium", sharePremi)
	p.setel("TotalPremium", nilaiPremi)
	p.setel("TotalSharePercentageClaim", shareKlaim)
	p.setel("TotalClaim", nilaiKlaim)
	return k.err
}
