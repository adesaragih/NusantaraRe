package models

// Untuk apa berkas ini: SELISIH PROPORSIONAL - port `Activity/EDMTCalculateTreatyDifference` (kelas
// Data-PolicyTreatyIn, ruleset 01-01-74; 10 langkah aktif, nol `//`). Dipicu tombol "Calculate Value Difference"
// (`Section/DetailPolicyTreatyInPropNewData2` S24, refresh thisSection) dan defer-load
// `Section/DetailPolicyTreatyInPropValueDifference` (tab Value Difference disegarkan setiap sel uang tab New Data
// berubah - `refresh otherSection`).
//
//	langkah 1     (PRE `.OldData.EDMNo==""`, salah -> lompat HasEDMNo) uang: Diff.X = New.X - Old.X;
//	              Installment, RiCommOgp, OveriddingCommOgp, RiCommOnp, OveriddingCommOnp,
//	              TotalSharePercentagePremium/Claim DISALIN
//	langkah 2.1   per baris spreading baru (.pxListSubscript): ClaimSpreaded / PremiumSpreaded = baru - lama(i);
//	              ClaimPercentage, SharePercentage, TreatyType, TreatyName disalin
//	langkah 3.1   per baris angsuran baru: PaymentTotal / Premium = baru - lama(i); DueDate, InstallmentNo,
//	              InstallmentPercentage disalin
//	langkah 4-6   (label HasEDMNo, `.OldData.EDMNo != ""`) rumus SAMA tetapi pengurangnya `OldData.TreatyDifference.X`
//
// ⛔ PENYIMPANGAN SADAR `[keputusan work owner 23-09-2026]` spec-penyimpanan ID-28 / ID-30 (AC 16): SATU RUMUS
// UNTUK SEMUA generasi - `selisih.X = baris_ini.X - baris(OLD_POLIS_ID).X`. Varian langkah 4-6 (mengurangi
// terhadap SELISIH generasi lampau; 180 - 50 = 130, padahal 180 - 150 = 30) TIDAK dibangun. Terbukti dari data
// produksi 23-09: `OldData.TreatyDifference` dokumen nyata kosong.
//
// ⚠️ Pasangan baris menurut POSISI (`.pxListSubscript`), bukan kunci dagang (ID-34): baris baru ke-i lawan baris
// lama ke-i; baris lama yang tidak ada (baris tambahan, Add) bernilai 0 dalam ekspresi Pega.
// ⚠️ Daftar selisih disusun ULANG sepanjang daftar baru - Pega menulis per indeks sehingga baris selisih lama di
// luar panjang daftar baru tertinggal; di EDM spreading tidak dapat dihapus (ID-16) sehingga hanya angsuran yang
// dapat memendek (FillPaymentInstallment), dan sisa baris itu nol pembaca produksi.

import "github.com/cockroachdb/apd/v3"

const od = HalamanPolis + ".OldData."

// medanUangSelisih - medan yang DIKURANGI `EDMTCalculateTreatyDifference` langkah 1 (urut XML), tanpa
// TotalPremium / TotalClaim (turunan baris, `HitungTotalSelisih`).
var medanUangSelisih = []string{"GrossPremium", "PremiOgp", "ResultOgp1", "ResultOgp2", "Claim", "SalvageValue",
	"ExcessLoss", "NetPremium", "BalanceDueTo", "PremiOnp", "ResultOnp1", "ResultOnp2", "Deduction1", "Deduction2",
	"PPNValue", "PPHValue", "BalanceBeforeTax", "BalanceBeforePPH"}

// medanSalinSelisih - medan yang DISALIN langkah 1 (persen dan jumlah termin).
var medanSalinSelisih = []string{"Installment", "RiCommOgp", "OveriddingCommOgp", "RiCommOnp", "OveriddingCommOnp"}

// EDMTCalculateTreatyDifference menghitung `PolicyTreatyIn.TreatyDifference` dari data baru dan `OldData`.
func EDMTCalculateTreatyDifference(h *Halaman) error {
	k := &kalkulator{}
	for _, m := range medanUangSelisih {
		h.SetelAngka(sd+m, k.Kurang(angkaJalur(k, h, pt+m), angkaJalur(k, h, od+m)))
	}
	for _, m := range medanSalinSelisih {
		h.Setel(sd+m, h.Ambil(pt+m))
	}
	lamaS := h.AmbilDaftar(od + "SpreadingRiskList")
	var spreading []Baris
	for i, b := range h.AmbilDaftar(DaftarSpreading) {
		l := barisKe(lamaS, i)
		spreading = append(spreading, Baris{
			"ClaimPercentage": b["ClaimPercentage"],
			"SharePercentage": b["SharePercentage"],
			"ClaimSpreaded":   formatAngka(k.Kurang(angkaBaris(k, b, "ClaimSpreaded"), angkaBaris(k, l, "ClaimSpreaded"))),
			"PremiumSpreaded": formatAngka(k.Kurang(angkaBaris(k, b, "PremiumSpreaded"), angkaBaris(k, l, "PremiumSpreaded"))),
			"TreatyType":      b["TreatyType"],
			"TreatyName":      b["TreatyName"],
		})
	}
	lamaA := h.AmbilDaftar(od + "ListInstallment")
	var angsuran []Baris
	for i, b := range h.AmbilDaftar(DaftarAngsuran) {
		l := barisKe(lamaA, i)
		angsuran = append(angsuran, Baris{
			"DueDate":               b["DueDate"],
			"InstallmentNo":         b["InstallmentNo"],
			"InstallmentPercentage": b["InstallmentPercentage"],
			"PaymentTotal":          formatAngka(k.Kurang(angkaBaris(k, b, "PaymentTotal"), angkaBaris(k, l, "PaymentTotal"))),
			"Premium":               formatAngka(k.Kurang(angkaBaris(k, b, "Premium"), angkaBaris(k, l, "Premium"))),
		})
	}
	if k.err != nil {
		return k.err
	}
	h.SetelDaftar(DaftarSelisihSpreading, buangKosong(spreading))
	h.SetelDaftar(DaftarSelisihAngsuran, buangKosong(angsuran))
	return HitungTotalSelisih(h)
}

// HitungTotalSelisih - empat total `TreatyDifference` yang TIDAK disimpan (katalog_selisih.go), dihitung saat
// dibaca: TotalPremium / TotalClaim = langkah 1 `New.Total - Old.Total`, sama dengan jumlah baris spreading
// selisih bila tidak ada baris lama yang hilang (ID-15, ID-16); TotalSharePercentagePremium / Claim disalin dari
// data baru.
func HitungTotalSelisih(h *Halaman) error {
	k := &kalkulator{}
	premi, klaim := k.d("0"), k.d("0")
	for _, b := range h.AmbilDaftar(DaftarSelisihSpreading) {
		premi = k.Tambah(premi, angkaBaris(k, b, "PremiumSpreaded"))
		klaim = k.Tambah(klaim, angkaBaris(k, b, "ClaimSpreaded"))
	}
	if k.err != nil {
		return k.err
	}
	if len(h.AmbilDaftar(DaftarSelisihSpreading)) == 0 && h.Ambil(sd+"NetPremium") == "" {
		return nil // belum pernah dihitung
	}
	h.SetelAngka(sd+"TotalPremium", premi)
	h.SetelAngka(sd+"TotalClaim", klaim)
	h.Setel(sd+"TotalSharePercentagePremium", h.Ambil(pt+"TotalSharePercentagePremium"))
	h.Setel(sd+"TotalSharePercentageClaim", h.Ambil(pt+"TotalSharePercentageClaim"))
	return nil
}

// angkaJalur membaca satu jalur halaman sebagai desimal (kosong = 0), galat dicatat kalkulator.
func angkaJalur(k *kalkulator, h *Halaman, jalur string) *apd.Decimal {
	d, err := h.Angka(jalur)
	if err != nil {
		k.catat(err)
		return k.d("0")
	}
	return d
}

// barisKe - anggota ke-i (berbasis 0) daftar, atau baris kosong bila tidak ada (properti tak bernilai = 0).
func barisKe(b []Baris, i int) Baris {
	if i < len(b) && b[i] != nil {
		return b[i]
	}
	return Baris{}
}

// buangKosong - nil untuk daftar tanpa anggota (daftar tak pernah ditulis).
func buangKosong(b []Baris) []Baris {
	if len(b) == 0 {
		return nil
	}
	return b
}
