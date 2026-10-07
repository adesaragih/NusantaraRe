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
	"slices"
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
	// DaftarSurvei - QuotationData.SurveyReportList (salinan polis lama, `Protection_Act` langkah 4).
	DaftarSurvei = HalamanPolis + ".QuotationData.SurveyReportList"
	DaftarUsulan = HalamanPolis + ".SuggestList"
)

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

// presisiSpreadingEDM - presisi `@divide` CountSpreading_Act EDM langkah 5.1 (`@divide(..,20)`); NB 10
// (`presisiBagiRata`). spec-penyimpanan ID-39 / AC 35: perbedaan sadar, ditiru apa adanya.
const presisiSpreadingEDM = 20

// CountSpreading = `Activity/CountSpreading_Act` VERSI EDM (korpus `EDM Treaty In`, kelas Data-PolicyTreatyIn,
// ruleset 01-01-95; 5 langkah aktif, 4 `//`). Param.Index hanya dibaca langkah `//` - di sini tidak dipakai.
//
//	1-3 `//`  SharePercentage / ClaimPercentage / Premium-ClaimSpreaded baris Param.Index - tidak diport
//	4         PRE @LengthOfPageList(.SpreadingRiskList)=1 && (.SpreadingRiskList(1).SplitRNMSharePct = "" ||
//	          .SpreadingRiskList(1).SplitRNMSharePct = 0) -> .SpreadingRiskList(1).SplitRNMSharePct = 100
//	5         setiap baris (5.1 PRE `.SharePercentage==""` TIDAK dicentang - selalu jalan):
//	          SharePercentage = @if(.SharePercentage=="",@divide(.SplitRNMSharePct,Primary.RNMShare,20),.SharePercentage)
//	          ClaimPercentage = @if(.ClaimPercentage == "",.SharePercentage,.ClaimPercentage)
//	          PremiumSpreaded = Primary.NetPremium*@divide(.SharePercentage,100,20)
//	          ClaimSpreaded = (Primary.ExcessLoss + Primary.Claim - Primary.SalvageValue)* @divide(.ClaimPercentage ,100,20)
//	          5.2 jumlah empat lokal
//	6         TotalSharePercentagePremium / TotalPremium / TotalSharePercentageClaim / TotalClaim = lokal (juga
//	          tanpa baris: 0)
//	7 `//`    Call BreakDownSpreading_Act - tidak diport
//
// ⛔ Cabang 5.1 `.SharePercentage == ""` TIDAK diport: pembaginya `Primary.RNMShare` (`PolicyTreatyIn.RNMShare`,
// bukan master `TreatyIn.RNMShare`) tidak diisi rule mana pun (korpus EDM dan NB), dan `.SplitRNMSharePct` baris yang
// ditambah juga kosong - hasilnya tidak pernah bermakna. Keputusan work owner 07-10-2026 (rekomendasi b): baris
// ber-%Share kosong - baris yang DITAMBAH lewat Add; baris generasi lama selalu terisi `FillSpreading` - DITOLAK
// dengan pesan di baris itu (`PesanShareSpreadingKosong`), baris lain tetap dihitung, total hanya dari baris
// ber-%Share. Submit Admin tertahan lewat `CountOGPONP` langkah 9 (validasi kirim menjalankannya atas salinan halaman,
// satu panggilan per baris - pesannya karena itu dipasang SEKALI per baris). ⚠️ Berbeda dari NB (bawaan `100/jumlah baris`,
// presisi 10) - ID-39, AC 35-36.
func CountSpreading(h *Halaman, _ int) error {
	k := &kalkulator{}
	p := polis{h, k}
	daftar := h.AmbilDaftar(DaftarSpreading)
	// langkah 4
	if len(daftar) == 1 && pbKosongAtauNol(daftar[0]["SplitRNMSharePct"]) {
		daftar[0]["SplitRNMSharePct"] = "100"
	}
	net := p.n("NetPremium")
	// (Primary.ExcessLoss + Primary.Claim - Primary.SalvageValue)
	klaim := k.Kurang(k.Tambah(p.n("ExcessLoss"), p.n("Claim")), p.n("SalvageValue"))
	shareKlaim, nilaiKlaim := apd.New(0, 0), apd.New(0, 0)
	sharePremi, nilaiPremi := apd.New(0, 0), apd.New(0, 0)
	for i, b := range daftar { // langkah 5
		if b["SharePercentage"] == "" { // 5.1 - keputusan WO 07-10-2026: pesan, bukan pembagian
			if m := PesanShareSpreadingKosong(i + 1); !slices.Contains(h.Pesan[""], m) {
				h.TambahPesan("", m)
			}
			continue
		}
		if b["ClaimPercentage"] == "" {
			b["ClaimPercentage"] = b["SharePercentage"]
		}
		b["PremiumSpreaded"] = formatAngka(k.Kali(net, k.BagiBulat(angkaBaris(k, b, "SharePercentage"), seratus, presisiSpreadingEDM)))
		b["ClaimSpreaded"] = formatAngka(k.Kali(klaim, k.BagiBulat(angkaBaris(k, b, "ClaimPercentage"), seratus, presisiSpreadingEDM)))
		// 5.2
		shareKlaim = k.Tambah(angkaBaris(k, b, "ClaimPercentage"), shareKlaim)
		nilaiKlaim = k.Tambah(angkaBaris(k, b, "ClaimSpreaded"), nilaiKlaim)
		sharePremi = k.Tambah(angkaBaris(k, b, "SharePercentage"), sharePremi)
		nilaiPremi = k.Tambah(angkaBaris(k, b, "PremiumSpreaded"), nilaiPremi)
	}
	if k.err != nil {
		return k.err
	}
	// langkah 6
	p.setel("TotalSharePercentagePremium", sharePremi)
	p.setel("TotalPremium", nilaiPremi)
	p.setel("TotalSharePercentageClaim", shareKlaim)
	p.setel("TotalClaim", nilaiKlaim)
	return nil
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
