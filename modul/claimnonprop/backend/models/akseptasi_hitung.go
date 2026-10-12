package models

// Untuk apa berkas ini: HITUNG DAN TAMBAH BARIS grid Claim Acceptation / Loss Allocation panel akseptasi
// (AdjustmentDetailNP).
//
// ⭐ PENYIMPANGAN DARI EKSPOR, keputusan pemakai 12-10-2026 ("lanjut semua"; kasus CLMNP-000022): di ekspor, sel dan
// tombol Add kedua grid memanggil `AddListClaimNP_Act`, `SetCurrency_Act`, `CountClaimTNP_Act`, `AddLossAlocation_Act`,
// `CountLossAllocation_act` dengan kelas Data-Adjustment yang tidak punya activity itu (OQ-CNP-33) - Claim Amount in
// IDR / Cedant tidak pernah berubah dan Add tidak menambah baris. Di sini rumus versi kelas Work dipakai atas baris
// AKSEPTASI ke-n (tingkat klaim tidak disentuh):
//
//	Claim Acceptation  CountClaimTNP_Act 13.2.3 (`hitungBarisKlaim`): nilai = Value - CNPDeductible + TPL + Fee +
//	                   Adjuster Fee - Salvage; USD (Claim Amount in IDR) = nilai × Rate of Exchange; Claim Amount
//	                   Cedant = nilai × Share Ceding %; ClaimAmountAdjust = (Value - CNPDeductible [+ TPL bila IsTPL])
//	                   × Share Ceding %. TPL baris sendiri (versi klaim: `max(TPL, TPL baris terakhir)`).
//	Loss Allocation    CountClaimTNP_Act 13.2.4 / CountLossAllocation_act 4 (`HitungPersen`): Claim Amount = ClaimAmountAdjust
//	                   klaim bermata uang sama × Share %, Adjuster Fee / Salvage / Fee = nilai klaim × Share Ceding % ×
//	                   Share %. Alokasi XOL akseptasi tetap lewat Claim Amount grid XOL Allocation (`AdjClaimCNP`).
//	Add                AddListClaimNP_Act (baris kosong) / AddLossAlocation_Act (satu baris per mata uang Claim
//	                   Acceptation), ke daftar akseptasi ke-n.
//	XOL Allocation     ⭐ 12-10-2026 (laporan pemakai: Claim Amount Loss Allocation berubah, XOL Allocation tetap):
//	                   CountLossAllocation_act 11-17 (`alokasiXOL`, rumus tingkat klaim) atas Loss Allocation akseptasi
//	                   ber-To XOL - UR = deductible, layer berikut sampai limitnya (Claim Amount layer NETO) - lalu
//	                   nilai layer dengan rumus AdjClaimCNP_Act 2-3 (`nilaiLayerAkseptasi`: Claim Amount RNM dari
//	                   nilai BRUTO) dan langkah 11, 15-17 (reinstatement, Spreading In / Out, ProtectNilaiClaim).
//	                   `ClaimEstimation` lama tiap layer dipertahankan, selisihnya ke `AdjClaimValue` (Gross Value OS
//	                   = keduanya).
//
// Baris terkunci (`CNPFlagOuts = 1`) tidak dihitung ulang.

import "github.com/cockroachdb/apd/v3"

// TambahKlaimAkseptasi - tombol Add grid Claim Acceptation akseptasi ke-n: satu baris kosong.
func TambahKlaimAkseptasi(k *Konteks, h *Halaman, n int) error {
	if _, err := barisDaftar(h, DaftarAdjustment, n); err != nil {
		return err
	}
	h.TambahBaris(JalurAdj(n, AnakClaimAccept), Baris{"pxCreateDateTime": k.Waktu(), "pxCreateOperator": k.Pelaku,
		"pxCreateOpName": k.Pelaku})
	return nil
}

// SetCurrencyAkseptasi - sel Currency baris Claim Acceptation ke-b akseptasi ke-n: nama mata uang dan Rate of Exchange
// master (`TreatyInMaster.CurrencyList.Conversion`; IDR = 1), lalu hitung ulang.
func SetCurrencyAkseptasi(k *Konteks, h *Halaman, m MasterTreaty, n, b int) error {
	r, err := barisDaftar(h, JalurAdj(n, AnakClaimAccept), b)
	if err != nil {
		return err
	}
	nama, err := k.Acuan.NamaMataUang(k.Ctxt(), r["CurrencyID"])
	if err != nil {
		return err
	}
	kurs := ""
	for _, c := range m.CurrencyList {
		if c.Currency == nama {
			kurs = c.Conversion
		}
	}
	if r["CurrencyID"] == IDMataUangIDR {
		kurs = "1"
	}
	r["Currency"], r["AltValue"] = nama, kurs
	// Prorata baris baru = baris akseptasi lain bermata uang sama, selain itu 100 (alokasi XOL membaca
	// `PctProrateClaim` baris TERAKHIR mata uang itu - kosong membuat limit layer 0).
	if r["PctProrateClaim"] == "" {
		r["PctProrateClaim"] = "100"
		for _, x := range h.AmbilDaftar(JalurAdj(n, AnakClaimAccept)) {
			if x["Currency"] == nama && x["PctProrateClaim"] != "" {
				r["PctProrateClaim"] = x["PctProrateClaim"]
				break
			}
		}
	}
	return HitungKlaimAkseptasi(h, n)
}

// HitungKlaimAkseptasi - sel angka grid Claim Acceptation akseptasi ke-n: Claim Amount in IDR, Claim Amount Cedant,
// ClaimAmountAdjust setiap baris, lalu Loss Allocation akseptasi bermata uang sama.
func HitungKlaimAkseptasi(h *Halaman, n int) error {
	if _, err := barisDaftar(h, DaftarAdjustment, n); err != nil {
		return err
	}
	var kal Kalkulator
	share := kal.BagiPega(kal.H(h, CD+"ShareCeding"), k100())
	isTPL := h.Ambil(CD+"IsTPL") == "true"
	klaim := h.AmbilDaftar(JalurAdj(n, AnakClaimAccept))
	for _, r := range klaim {
		if r["CNPFlagOuts"] == "1" {
			continue
		}
		bersih := kal.Kurang(kal.B(r, "Value"), kal.B(r, "CNPDeductible"))
		dasar := bersih
		if isTPL {
			dasar = kal.Tambah(dasar, kal.B(r, "TPL"))
		}
		fee := kal.Kurang(kal.Tambah(kal.B(r, "CNPOthersFee"), kal.B(r, "AdjusterFee")), kal.B(r, "Salvage"))
		nilai := kal.Tambah(kal.Tambah(bersih, kal.B(r, "TPL")), fee)
		r["USD"] = Teks(kal.Kali(nilai, kal.B(r, "AltValue")))
		r["ClaimAmountCedant"] = Teks(kal.Kali(nilai, share))
		r["ClaimAmountAdjust"] = Teks(kal.Kali(dasar, share))
	}
	for _, l := range h.AmbilDaftar(JalurAdj(n, AnakLossAlloc)) {
		if l["CNPFlagOuts"] != "1" {
			sebarLossAkseptasi(&kal, l, klaimMataUang(klaim, l["Currency"]), share)
		}
	}
	return kal.Galat()
}

// TambahLossAkseptasi - tombol Add grid Loss Allocation akseptasi ke-n: satu baris per mata uang Claim Acceptation
// (AddLossAlocation_Act, sama dengan tab Estimation).
func TambahLossAkseptasi(h *Halaman, n int) error {
	if _, err := barisDaftar(h, DaftarAdjustment, n); err != nil {
		return err
	}
	sudah := map[string]bool{}
	for _, r := range h.AmbilDaftar(JalurAdj(n, AnakClaimAccept)) {
		if r["Currency"] == "" || sudah[r["Currency"]] {
			continue
		}
		sudah[r["Currency"]] = true
		h.TambahBaris(JalurAdj(n, AnakLossAlloc), Baris{"Currency": r["Currency"], "CurrencyID": r["CurrencyID"]})
	}
	return nil
}

// HitungLossAkseptasi - sel Share (%) baris Loss Allocation ke-b akseptasi ke-n (`HitungPersen`).
func HitungLossAkseptasi(h *Halaman, n, b int) error {
	l, err := barisDaftar(h, JalurAdj(n, AnakLossAlloc), b)
	if err != nil {
		return err
	}
	if l["CNPFlagOuts"] == "1" {
		return nil
	}
	var kal Kalkulator
	share := kal.BagiPega(kal.H(h, CD+"ShareCeding"), k100())
	sebarLossAkseptasi(&kal, l, klaimMataUang(h.AmbilDaftar(JalurAdj(n, AnakClaimAccept)), l["Currency"]), share)
	return kal.Galat()
}

// HitungLossJumlahAkseptasi - sel Claim Amount baris Loss Allocation ke-b akseptasi ke-n (`HitungJumlah`): Share (%) =
// Claim Amount ÷ ClaimAmountAdjust klaim bermata uang sama × 100.
func HitungLossJumlahAkseptasi(h *Halaman, n, b int) error {
	l, err := barisDaftar(h, JalurAdj(n, AnakLossAlloc), b)
	if err != nil {
		return err
	}
	r := klaimMataUang(h.AmbilDaftar(JalurAdj(n, AnakClaimAccept)), l["Currency"])
	if l["CNPFlagOuts"] == "1" || r == nil {
		return nil
	}
	var kal Kalkulator
	l["ClaimPercentage"] = Teks(kal.Kali(kal.BagiPega(kal.B(l, "ClaimAmountAdjust"), kal.B(r, "ClaimAmountAdjust")), k100()))
	return kal.Galat()
}

// HitungXOLAkseptasi - XOL Allocation akseptasi ke-n disusun ulang dari Loss Allocation akseptasi ber-To XOL
// (`alokasiXOL`), lalu reinstatement, Spreading In / Out, ProtectNilaiClaim, Spreading klaim (AdjClaimCNP_Act 11,
// 15-17). Layer yang sudah ada: `ClaimEstimation` lama dipertahankan, `AdjClaimValue` = alokasi baru - estimasi; layer
// baru: estimasi 0; layer lama yang tak teralokasi lagi: alokasi 0 (tidak dibuang).
func HitungXOLAkseptasi(k *Konteks, h *Halaman, m MasterTreaty, n int) error {
	if _, err := barisDaftar(h, DaftarAdjustment, n); err != nil {
		return err
	}
	baru, err := alokasiXOL(k, h, m, h.AmbilDaftar(JalurAdj(n, AnakLossAlloc)), h.AmbilDaftar(JalurAdj(n, AnakClaimAccept)))
	if err != nil {
		return err
	}
	kunci := func(b Baris) string { return b["Currency"] + "|" + b["TreatyName"] }
	estimasi := map[string]string{}
	lama := h.AmbilDaftar(JalurAdj(n, AnakXOL))
	for _, b := range lama {
		if _, ada := estimasi[kunci(b)]; !ada {
			estimasi[kunci(b)] = b["ClaimEstimation"]
		}
	}
	var kal Kalkulator
	sudah := map[string]bool{}
	for _, b := range baru {
		sudah[kunci(b)] = true
		est, ada := estimasi[kunci(b)]
		if !ada || est == "" {
			est = "0"
		}
		b["ClaimEstimation"] = est
		nilaiLayerAkseptasi(&kal, b)
	}
	for _, b := range lama {
		if sudah[kunci(b)] {
			continue
		}
		sudah[kunci(b)] = true
		z := b.Salin()
		z["AdjClaimValue"] = Teks(kal.Kurang(apd.New(0, 0), kal.B(b, "ClaimEstimation")))
		for _, p := range []string{"ClaimAmountAdjust", "TotalClaim", "ClaimSpreaded", "AdjusterFee", "Salvage",
			"CNPOthersFee", "TotalClaimRNM", "TotalSpread", "ClaimAmountIDR"} {
			z[p] = "0"
		}
		baru = append(baru, z)
	}
	reinstatementXOL(&kal, baru)
	if err := kal.Galat(); err != nil {
		return err
	}
	h.SetelDaftar(JalurAdj(n, AnakXOL), baru)
	if err := SebarAkseptasi(h, n); err != nil {
		return err
	}
	ProtectNilaiClaim(h, n)
	return SusunSpreadingAkseptasiKlaim(h)
}

// nilaiLayerAkseptasi - satu baris XOL Allocation akseptasi dari `TotalClaim` (Claim Amount layer, neto) dengan rumus
// AdjClaimCNP_Act 2-3 - bentuk yang Pega tampilkan di akseptasi:
//
//	UR      ClaimAmountAdjust = TotalClaim
//	layer   ClaimAmountAdjust = TotalClaim - (Adjuster Fee + Fee - Salvage) ÷ RNM %   (bruto, fee kembali ke 100 %)
//	        ClaimSpreaded (Claim Amount RNM) = ClaimAmountAdjust × RNM %
//	        TotalClaimRNM = ClaimSpreaded + Adjuster Fee + Fee - Salvage
//	AdjClaimValue = ClaimAmountAdjust - ClaimEstimation
//
// ⭐ 12-10-2026 (laporan pemakai, tangkapan Pega): Claim 636,4 jt, Salvage 256 jt, RNM 40 % -> XL 1ST LAYER Claim Amount
// 55,4 jt, Claim Amount RNM 124,56 jt (= 311,4 jt × 40 %), Salvage 102,4 jt, Total Claim RNM 22,16 jt. Hasil
// `alokasiXOL` (rumus TINGKAT KLAIM: Claim RNM = estimasi neto × RNM) memberi 22,16 jt - beda bentuk dengan akseptasi.
func nilaiLayerAkseptasi(kal *Kalkulator, b Baris) {
	if b["TreatyType"] == TreatyUR {
		b["ClaimAmountAdjust"] = b["TotalClaim"]
	} else {
		pct := kal.B(b, "ClaimPercentage")
		bruto := kal.Kurang(kal.B(b, "TotalClaim"), kal.Kurang(kal.Tambah(bagiPersen(kal, kal.B(b, "AdjusterFee"), pct),
			bagiPersen(kal, kal.B(b, "CNPOthersFee"), pct)), bagiPersen(kal, kal.B(b, "Salvage"), pct)))
		b["ClaimAmountAdjust"] = Teks(bruto)
		b["ClaimSpreaded"] = Teks(kal.Kali(bruto, kal.BagiPega(pct, k100())))
	}
	b["AdjClaimValue"] = Teks(kal.Kurang(kal.B(b, "ClaimAmountAdjust"), kal.B(b, "ClaimEstimation")))
	b["TotalClaimRNM"] = Teks(kal.Kurang(kal.Tambah(kal.Tambah(kal.B(b, "ClaimSpreaded"), kal.B(b, "AdjusterFee")),
		kal.B(b, "CNPOthersFee")), kal.B(b, "Salvage")))
	b["TotalSpread"] = b["TotalClaimRNM"]
}

// klaimMataUang - baris Claim Acceptation TERAKHIR bermata uang `cur` (13.2.4: tiap baris klaim menimpa), nil bila tak
// ada.
func klaimMataUang(klaim []Baris, cur string) Baris {
	var r Baris
	for _, k := range klaim {
		if cur != "" && k["Currency"] == cur {
			r = k
		}
	}
	return r
}

// sebarLossAkseptasi - satu baris Loss Allocation dari baris klaim bermata uang sama (13.2.4 / `HitungPersen`).
func sebarLossAkseptasi(kal *Kalkulator, l, r Baris, share *apd.Decimal) {
	if r == nil {
		return
	}
	p := kal.BagiPega(kal.B(l, "ClaimPercentage"), k100())
	l["ClaimAmountAdjust"] = Teks(kal.Kali(kal.B(r, "ClaimAmountAdjust"), p))
	l["AdjusterFee"] = Teks(kal.Kali(kal.Kali(kal.B(r, "AdjusterFee"), share), p))
	l["CNPOthersFee"] = Teks(kal.Kali(kal.Kali(kal.B(r, "CNPOthersFee"), share), p))
	l["Salvage"] = Teks(kal.Kali(kal.Kali(kal.B(r, "Salvage"), share), p))
}

// Medan turunan baris akseptasi untuk kaki grid Claim Acceptation (`TotalKlaimAkseptasi`, tidak disimpan).
const (
	PropTotalKlaimIDR  = "TotalClaimAmountIDR"
	PropTotalCedantIDR = "TotalClaimAmountCedantIDR"
)

// TotalKlaimAkseptasi - kaki grid Claim Acceptation akseptasi ke-n (turunan, tidak disimpan): Σ Claim Amount in IDR
// dan Σ Claim Amount Cedant dalam IDR (× Rate of Exchange baris). Di ekspor kakinya `TotalListClaimAmount(IDR)` tanpa
// penulis (OQ-CNP-37) - selalu kosong.
func TotalKlaimAkseptasi(h *Halaman, n int) (idr, cedantIDR string, err error) {
	var kal Kalkulator
	ti, tc := apd.New(0, 0), apd.New(0, 0)
	for _, r := range h.AmbilDaftar(JalurAdj(n, AnakClaimAccept)) {
		ti = kal.Tambah(ti, kal.B(r, "USD"))
		tc = kal.Tambah(tc, kal.Kali(kal.B(r, "ClaimAmountCedant"), kal.B(r, "AltValue")))
	}
	return Teks(ti), Teks(tc), kal.Galat()
}
