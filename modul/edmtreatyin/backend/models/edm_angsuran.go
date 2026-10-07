package models

// Untuk apa berkas ini: ANGSURAN ENDORSEMEN - port `Activity/FillPaymentInstallmentEDMT` (kelas
// Data-PolicyTreatyIn, ruleset 01-01-81, Param.Installment; 16 langkah aktif, nol `//`). Pemanggil:
// `EDMChooseBusiness_Act` langkah 12 (kasus ber-`pyWorkIDPrefix` "EDMT-", Param.Installment = jumlah
// `TreatyIn.Installment(1).InstallmentList` master) dan sel `.Installment` S12 `Section/DetailPolicyTreatyInAddendum`
// (NonProp baru; change -> refresh, Param.Installment = .Installment).
//
//	1        Property-Remove .ListInstallment
//	2.1      per baris TreatyXOLDifferenceList (c): ListInstallment(c).Currency / IDCurrency = .Currency /
//	         .IDCurrency; Premium = PaymentTotal = .NetPremi; PremiumAfterPPN / PPH / Tax = .NetPremiAfter*;
//	2.1.2    TreatyDifference.TotalPremium = .NetPremi (mata uang terakhir menang - nol pembaca, tidak disimpan)
//	2.2      per baris ListInstallment: rincian InstallmentList sebanyak Param.Installment -
//	         persen = @Math.divide(100, N, 4), baris terakhir menampung sisa (2.2.2.3: q + (100 - jumlah));
//	         Premium = PaymentTotal = @Math.divide(premi x persen, 100, 4); After* sama; DueDate = .StatementDate
//	3        pyExpanded (keadaan layar) ; 4 Obj-Save (penyimpanan oleh services)
//
// Algoritma persen sama dengan NB `FillPaymentInstallment` (models/angsuran.go), tetapi dua tingkat per mata uang
// dan DueDate = StatementDate (bukan sekarang).

import (
	"strconv"
	"strings"

	"github.com/cockroachdb/apd/v3"
)

// FillPaymentInstallmentEDMT menyusun ulang `ListInstallment` (dan rinciannya) dari `TreatyXOLDifferenceList`.
// `installment` = Param.Installment (teks; bukan bilangan positif = nol putaran rincian).
func FillPaymentInstallmentEDMT(h *Halaman, installment string) error {
	k := &kalkulator{}
	// langkah 1
	for i := range h.AmbilDaftar(DaftarAngsuran) {
		h.SetelDaftar(JalurAnak(DaftarAngsuran, i+1, "InstallmentList"), nil)
	}
	h.SetelDaftar(DaftarAngsuran, nil)
	// langkah 2.1
	var daftar []Baris
	for _, x := range h.AmbilDaftar(DaftarSelisihXOL) {
		daftar = append(daftar, Baris{
			"Currency":        x["Currency"],
			"IDCurrency":      x["IDCurrency"],
			"Premium":         x["NetPremi"],
			"PaymentTotal":    x["NetPremi"],
			"PremiumAfterPPN": x["NetPremiAfterPPN"],
			"PremiumAfterPPH": x["NetPremiAfterPPH"],
			"PremiumAfterTax": x["NetPremiAfterTax"],
		})
	}
	h.SetelDaftar(DaftarAngsuran, buangKosong(daftar))
	// langkah 2.2
	n, err := strconv.Atoi(strings.TrimSpace(installment))
	if err != nil || n <= 0 {
		return nil // batas ulang Param.Installment-1 < 0: nol putaran
	}
	bagian := k.BagiBulat(seratus, k.d(strconv.Itoa(n)), 4)
	// DueDate kolom DATE (T_POLIS_INSTALMENT_DETAIL.DUE_DATE): bagian tanggal StatementDate.
	statement := tanggalSaja(h.Ambil(pt + "StatementDate"))
	for i, b := range daftar {
		premi := angkaBaris(k, b, "Premium")
		ppn := angkaBaris(k, b, "PremiumAfterPPN")
		pph := angkaBaris(k, b, "PremiumAfterPPH")
		pajak := angkaBaris(k, b, "PremiumAfterTax")
		no := 0
		pct := k.d("0")
		var rinci []Baris
		for j := 0; j < n; j++ {
			no++                        // 2.2.2.1
			pct = k.Tambah(pct, bagian) // 2.2.2.2
			r := Baris{}
			if no == n && k.Tambah(pct, bagian).Cmp(seratus) != 0 { // 2.2.2.3
				r["InstallmentPercentage"] = formatAngka(k.Tambah(bagian, k.Kurang(seratus, pct)))
			}
			if no < n { // 2.2.2.4
				r["InstallmentPercentage"] = formatAngka(bagian)
			}
			persen := angkaBaris(k, r, "InstallmentPercentage")
			bagi := func(x *apd.Decimal) string { return formatAngka(k.BagiBulat(k.Kali(x, persen), seratus, 4)) }
			r["InstallmentNo"] = strconv.Itoa(no) // 2.2.2.5
			r["Premium"] = bagi(premi)
			r["DueDate"] = statement
			r["PaymentTotal"] = bagi(premi)
			r["Currency"] = b["Currency"]
			r["IDCurrency"] = b["IDCurrency"]
			r["PremiumAfterPPN"] = bagi(ppn)
			r["PremiumAfterPPH"] = bagi(pph)
			r["PremiumAfterTax"] = bagi(pajak)
			rinci = append(rinci, r)
		}
		h.SetelDaftar(JalurAnak(DaftarAngsuran, i+1, "InstallmentList"), rinci)
	}
	return k.err
}
