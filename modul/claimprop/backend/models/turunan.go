package models

// Untuk apa berkas ini: SATU FUNGSI PERHITUNGAN TURUNAN klaim (tiket 00/05-09; spec §4 butir 1-4, AC 20-24, 26, 43).
//
// ⚠️ PENYIMPANGAN SADAR spec §4 (keputusan work owner): di Pega nilai turunan ditulis berantai oleh banyak activity
// (CountTotalInsterest_Act langkah 11-17, CopyOldataCurr_act, CountPersen_act langkah 6, CountEstimation_Act
// langkah 3, CountSpreading_act langkah 8-11, CountListClaimAmountIDR, AddListClaimAmount langkah 7) - nilai akhir
// bergantung urutan pemanggilan dan membaca balik hasil antara. Di sini SATU fungsi, dijalankan sesudah setiap aksi,
// menghitung SEMUA turunan sekali jalan dari basis asli. Rumus tiap turunan tetap rumus Pega; yang berubah hanya
// ketelitiannya (presisi penuh) dan urutannya (kali dulu, bagi terakhir).
//
// Basis (diisi layar atau activity pembuat baris): interest (CurrencyID, KursObjectItem, TSIPerObject), claim amount
// (CurrencyID, IDR = kurs, ClaimAmount, Note), header (ShareCeding, deductible, RNMShareP), loss allocation
// (CurrencyID, TreatyType, SharePercentage, PremiumSpreaded = kurs, IsOldData), estimasi (CurrencyID, KursValue,
// GrossEstimationPct, TypeLossID, PrintFaceClaim), spreading klaim dan break QS (TreatyType, SharePercentage,
// CurrencyID).
//
// ⚠️ Baris BEKU tidak dihitung ulang (AC 23): loss allocation ber-IsOldData "Yes" (CountPersen_act 6.2.1,
// CountTotalInsterest_Act 12.2.1), claim amount ber-Note "Yes" (CountTotalInsterest_Act 11.2.1), estimasi
// ber-PrintFaceClaim 1 (sudah terkirim ke OS_AKSEPTASI_KLAIM - CountTotalInsterest_Act 13.2.1, CopyOldataCurr_act 8.1).

import (
	"github.com/cockroachdb/apd/v3"
)

// HitungTurunan menghitung seluruh nilai turunan halaman klaim.
func HitungTurunan(h *Halaman) error {
	var k Kalkulator
	hitungInterest(&k, h)
	hitungDeductible(&k, h)
	hitungClaimAmount(&k, h)
	hitungLossAllocation(&k, h)
	hitungEstimasi(&k, h)
	hitungSpreadingKlaim(&k, h)
	return k.Galat()
}

// urutanMataUang - mata uang berbeda menurut urutan kemunculan pertama (pola Java "remove kurs yg sama" Pega
// menghapus dari BELAKANG, sehingga kemunculan PERTAMA yang tinggal).
func urutanMataUang(d []Baris) []Baris {
	var out []Baris
	sudah := map[string]bool{}
	for _, b := range d {
		id := b["CurrencyID"]
		if id == "" || sudah[id] {
			continue
		}
		sudah[id] = true
		out = append(out, Baris{"CurrencyID": id, "Currency": b["Currency"]})
	}
	return out
}

// hitungInterest - CountTotalInsterest_Act langkah 3-10: TSIPerObjectIDR = kurs x TSI, TotalInterestInsured per mata
// uang, TotalSumInsuredIDR.
func hitungInterest(k *Kalkulator, h *Halaman) {
	d := h.AmbilDaftar(DaftarInterest)
	total := apd.New(0, 0)
	per := map[string]*apd.Decimal{}
	for _, b := range d {
		idr := k.Kali(k.B(b, "KursObjectItem"), k.B(b, "TSIPerObject"))
		b["TSIPerObjectIDR"] = Teks(idr)
		total = k.Tambah(total, idr)
		if id := b["CurrencyID"]; id != "" {
			if per[id] == nil {
				per[id] = apd.New(0, 0)
			}
			per[id] = k.Tambah(per[id], k.B(b, "TSIPerObject"))
		}
	}
	var tot []Baris
	for _, m := range urutanMataUang(d) {
		m["Value"] = Teks(per[m["CurrencyID"]])
		tot = append(tot, m)
	}
	h.SetelDaftar(DaftarTotalTSI, tot)
	if len(d) == 0 {
		h.Hapus(CD + "TotalSumInsuredIDR")
		return
	}
	h.Setel(CD+"TotalSumInsuredIDR", Teks(total))
}

// tsiMataUang - TotalInterestInsured.Value satu mata uang (CountDeductible_Act langkah 2.1.1.1).
func tsiMataUang(h *Halaman, id string) string {
	for _, b := range h.AmbilDaftar(DaftarTotalTSI) {
		if b["CurrencyID"] == id {
			return b["Value"]
		}
	}
	return ""
}

// hitungDeductible - CountDeductible_Act (hanya FormType 2): Net = MAX(persen x basis, nilai flat) (AC 50, 51).
// TypeDeductible 1 ("Jika Klaim"): basis = TotalInterestInsured mata uang deductible; 2 ("Jika TSI"): TSIDeductible.
// FormType 1: NetDeductibleValue diisi langsung di layar ("Amount").
func hitungDeductible(k *Kalkulator, h *Halaman) {
	if h.Ambil(CD+"FormType") != "2" {
		return
	}
	var basis string
	switch h.Ambil(CD + "TypeDeductible") {
	case "1":
		if h.Ambil(CD+"CurrencyDeductible") == "" || h.Ambil(CD+"Amount") == "" {
			return
		}
		basis = tsiMataUang(h, h.Ambil(CD+"CurrencyDeductible"))
	case "2":
		if h.Ambil(CD+"Amount") == "" {
			return
		}
		basis = h.Ambil(CD + "TSIDeductible")
	default:
		return
	}
	persen := k.Persen(k.Teks("basis", basis), k.H(h, CD+"Amount"))
	flat := k.H(h, CD+"DeductibleValue")
	if Lebih(persen, flat) {
		h.Setel(CD+"NetDeductibleValue", Teks(persen))
	} else {
		h.Setel(CD+"NetDeductibleValue", Teks(flat))
	}
}

// hitungClaimAmount - rumus 2024 (AC 52, keputusan work owner; CountListClaimAmountIDR langkah 1.1):
// .NetDeductibleValue = header; .Value = ShareCeding% x ClaimAmount - NetDeductibleValue; .USD = .Value x .IDR.
// Rumus 2022 AddListClaimAmount langkah 7.1 (deductible SEBELUM share) ditinggalkan.
func hitungClaimAmount(k *Kalkulator, h *Halaman) {
	d := h.AmbilDaftar(DaftarClaimAmount)
	tot, totIDR := apd.New(0, 0), apd.New(0, 0)
	for _, b := range d {
		if b["Note"] != "Yes" {
			b["NetDeductibleValue"] = h.Ambil(CD + "NetDeductibleValue")
			nilai := k.Kurang(k.Persen(k.B(b, "ClaimAmount"), k.H(h, CD+"ShareCeding")), k.H(h, CD+"NetDeductibleValue"))
			b["Value"] = Teks(nilai)
			b["USD"] = Teks(k.Kali(nilai, k.B(b, "IDR")))
		}
		tot = k.Tambah(tot, k.B(b, "Value"))
		totIDR = k.Tambah(totIDR, k.B(b, "USD"))
	}
	if len(d) == 0 {
		h.Hapus(CD + "TotalListClaimAmount")
		h.Hapus(CD + "TotalListClaimAmountIDR")
		return
	}
	h.Setel(CD+"TotalListClaimAmount", Teks(tot))
	h.Setel(CD+"TotalListClaimAmountIDR", Teks(totIDR))
}

// nilaiClaimMataUang - `.Value` baris claim amount TERAKHIR bermata uang `id` (CountPersen_act langkah 6: perulangan
// claim amount menimpa Local.Value, baris terakhir menang).
func nilaiClaimMataUang(h *Halaman, id string) (string, bool) {
	v, ada := "", false
	for _, b := range h.AmbilDaftar(DaftarClaimAmount) {
		if b["CurrencyID"] == id {
			v, ada = b["Value"], true
		}
	}
	return v, ada
}

// hitungLossAllocation - CountPersen_act langkah 6.2.1: .ClaimSpreaded = Share% x Value(claim amount mata uang sama);
// .ClaimEstimation = .ClaimSpreaded x .PremiumSpreaded (kurs). Baris IsOldData "Yes" beku.
func hitungLossAllocation(k *Kalkulator, h *Halaman) {
	for _, b := range h.AmbilDaftar(DaftarLossAlloc) {
		if b["IsOldData"] == "Yes" {
			continue
		}
		v, ada := nilaiClaimMataUang(h, b["CurrencyID"])
		if !ada {
			continue
		}
		spread := k.Persen(k.Teks("Value", v), k.B(b, "SharePercentage"))
		b["ClaimSpreaded"] = Teks(spread)
		b["ClaimEstimation"] = Teks(k.Kali(spread, k.B(b, "PremiumSpreaded")))
	}
}

// hitungEstimasi - CountEstimation_Act langkah 3-10, 28: .PersenRNM = RNMShareP; .EstimationValue = RNMShareP% x
// .GrossEstimationPct; .ConvertValue = .EstimationValue x .KursValue; .ConvertGrossEstimasi = .GrossEstimationPct x
// .KursValue; total IDR; total mata uang asal bila hanya satu mata uang; ListTotalEstimation per mata uang.
func hitungEstimasi(k *Kalkulator, h *Halaman) {
	d := h.AmbilDaftar(DaftarEstimasi)
	totIDR, totGrossIDR := apd.New(0, 0), apd.New(0, 0)
	totNilai, totGross := apd.New(0, 0), apd.New(0, 0)
	perNilai, perGross := map[string]*apd.Decimal{}, map[string]*apd.Decimal{}
	for _, b := range d {
		if b["PrintFaceClaim"] != "1" {
			b["PersenRNM"] = h.Ambil(TM + "RNMShareP")
			nilai := k.Persen(k.B(b, "GrossEstimationPct"), k.H(h, TM+"RNMShareP"))
			b["EstimationValue"] = Teks(nilai)
			b["ConvertValue"] = Teks(k.Kali(nilai, k.B(b, "KursValue")))
			b["ConvertGrossEstimasi"] = Teks(k.Kali(k.B(b, "GrossEstimationPct"), k.B(b, "KursValue")))
		}
		totIDR = k.Tambah(totIDR, k.B(b, "ConvertValue"))
		totGrossIDR = k.Tambah(totGrossIDR, k.B(b, "ConvertGrossEstimasi"))
		totNilai = k.Tambah(totNilai, k.B(b, "EstimationValue"))
		totGross = k.Tambah(totGross, k.B(b, "GrossEstimationPct"))
		if id := b["CurrencyID"]; id != "" {
			if perNilai[id] == nil {
				perNilai[id], perGross[id] = apd.New(0, 0), apd.New(0, 0)
			}
			perNilai[id] = k.Tambah(perNilai[id], k.B(b, "EstimationValue"))
			perGross[id] = k.Tambah(perGross[id], k.B(b, "GrossEstimationPct"))
		}
	}
	mu := urutanMataUang(d)
	var lte []Baris
	for _, m := range mu {
		m["Value"] = Teks(perNilai[m["CurrencyID"]])
		m["IDR"] = Teks(perGross[m["CurrencyID"]])
		lte = append(lte, m)
	}
	h.SetelDaftar(DaftarTotalEst, lte)
	if len(d) == 0 {
		for _, p := range []string{"TotalEstimasiIDR", "TotalGrossEstimateIDR", "TotalEstimasi", "TotalGrossEstimateTreaty"} {
			h.Hapus(CD + p)
		}
		return
	}
	h.Setel(CD+"TotalEstimasiIDR", Teks(totIDR))
	h.Setel(CD+"TotalGrossEstimateIDR", Teks(totGrossIDR))
	if len(mu) < 2 { // langkah 8: hanya bila mata uangnya satu
		h.Setel(CD+"TotalEstimasi", Teks(totNilai))
		h.Setel(CD+"TotalGrossEstimateTreaty", Teks(totGross))
	} else {
		h.Hapus(CD + "TotalEstimasi")
		h.Hapus(CD + "TotalGrossEstimateTreaty")
	}
}

// estimasiMataUang - jumlah EstimationValue satu mata uang (CountSpreading_act langkah 8: `Estimate.EstimationValue`).
func estimasiMataUang(k *Kalkulator, h *Halaman, id string) *apd.Decimal {
	s := apd.New(0, 0)
	for _, b := range h.AmbilDaftar(DaftarEstimasi) {
		if b["CurrencyID"] == id {
			s = k.Tambah(s, k.B(b, "EstimationValue"))
		}
	}
	return s
}

// hitungSpreadingKlaim - CountSpreading_act langkah 8-11: SpreadingClaim.ClaimSpreaded = Share% x estimasi mata uang
// sama; SpreadingBreakQS.ClaimSpreaded = Share% x ClaimSpreaded baris SpreadingClaim TERAKHIR bermata uang sama
// (langkah 11 perulangan menimpa). Langkah 8 dan 11 berprakondisi nonaktif - berjalan untuk seluruh baris.
func hitungSpreadingKlaim(k *Kalkulator, h *Halaman) {
	terakhir := map[string]string{}
	for _, b := range h.AmbilDaftar(DaftarSpreading) {
		if b["CurrencyID"] == "" {
			continue
		}
		b["ClaimSpreaded"] = Teks(k.Persen(estimasiMataUang(k, h, b["CurrencyID"]), k.B(b, "SharePercentage")))
		terakhir[b["CurrencyID"]] = b["ClaimSpreaded"]
	}
	for _, b := range h.AmbilDaftar(DaftarBreakQS) {
		v, ada := terakhir[b["CurrencyID"]]
		if !ada {
			continue
		}
		b["ClaimSpreaded"] = Teks(k.Persen(k.Teks("ClaimSpreaded", v), k.B(b, "SharePercentage")))
	}
}

// TotalSharePerMataUang - jumlah SharePercentage SpreadingClaim per mata uang (CountSpreading_act langkah 6.2.2).
func TotalSharePerMataUang(h *Halaman) (map[string]*apd.Decimal, error) {
	var k Kalkulator
	out := map[string]*apd.Decimal{}
	for _, b := range h.AmbilDaftar(DaftarSpreading) {
		id := b["CurrencyID"]
		if out[id] == nil {
			out[id] = apd.New(0, 0)
		}
		out[id] = k.Tambah(out[id], k.B(b, "SharePercentage"))
	}
	return out, k.Galat()
}
