package models

// Untuk apa berkas ini: PORT ACTIVITY AKSEPTASI Non Prop (tab Acceptation, AdjustmentDetailNP): AddAkseptasiCNP_Act,
// DeleteAkseptasi_Act, AdjClaimCNP_Act (+ CountSpreadingXOL, ProtectNilaiClaim), SetInterimXOL_Act,
// SetPayableTreatyNP_Act, SetAccoutNo_Act, ProteksiSendKomiteCNP_Act, HitServiceToKasir_Act (muatan Kasir).

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/cockroachdb/apd/v3"
)

// Pesan VERBATIM.
const (
	PesanAdjusterKosong   = "Adjuster / Professional tidak boleh kosong"                // AddAkseptasiCNP_Act 2
	PesanConsultantKosong = "Consultant tidak boleh kosong"                             // AddAkseptasiCNP_Act 2
	PesanPilihPayable     = "Please choose payable first"                               // ProteksiSendKomiteCNP_Act 2
	PesanPilihPayment     = "Please choose Payment Type first"                          // ProteksiSendKomiteCNP_Act 2
	PesanRekeningKosong   = "Data Bank Account Can't NULL"                              // ProteksiSendKomiteCNP_Act 2
	PesanOkupasiKosong    = "Occupation cannot be empty"                                // ProteksiSendKomiteCNP_Act 2
	PesanPolisKomite      = "Please Input Policy No"                                    // ProteksiSendKomiteCNP_Act 2
	PesanEmailKosong      = "Email Profile Can't NULL"                                  // ProteksiSendKomiteCNP_Act 2
	PesanHarusCancel      = "Payment type must be Cancellation"                         // ProteksiSendKomiteCNP_Act 2
	PesanMataUangKosong   = "Currency cant be NULL"                                     // ProteksiSendKomiteCNP_Act 2
	PesanIDBankKosong     = "ID of Bank Can't be NULL"                                  // ProteksiSendKomiteCNP_Act 2
	PesanRekeningSalah    = "Error No Account "                                         // ProteksiSendKomiteCNP_Act 7.3
	PesanMelebihiEstimasi = "Adjustment value should not be more than estimation value" // 23.1
	PesanNetNegatif       = "Spreading In RNM Net Claim Can not be Negative"            // ProteksiNilaiClaim 2
	TeksTambahAkseptasi   = "Add Adjustment"                                            // AddAkseptasiCNP_Act 12
	TeksKirimKomite       = "Send Adjustment to Committe (Acceptation)"                 // ProteksiSendKomiteCNP_Act 2
	TeksKasirMasuk        = "Akseptasi Sudah Masuk ke Kasir"                            // HitServiceToKasir_Act 9.8
)

// PaymentTypeCancellation - Payment Type 7 (ProteksiSendKomiteCNP_Act langkah 4.1; label prompt values tidak diekspor,
// OQ-CNP-14).
const PaymentTypeCancellation = "7"

// adj - baris AdjustmentList ke-n.
func adj(h *Halaman, n int) (Baris, error) { return barisDaftar(h, DaftarAdjustment, n) }

// ---------------------------------------------------------------- tambah / hapus

// AddAkseptasi = `AddAkseptasiCNP_Act` (tombol Add grid Acceptation List). Mengembalikan false bila keluar karena
// proteksi (pesan terpasang).
func AddAkseptasi(k *Konteks, h *Halaman, m MasterTreaty) (bool, error) {
	no := []rune(h.Ambil(CD + "NoClaim"))
	titikKe6 := len(no) > 5 && no[5] == '.'
	if h.Ambil(OQ+"BusinessName") == "" || h.Ambil(CD+"QuotationData.BusinessOldId") == "" || titikKe6 { // 3 (trans T:6)
		h.TambahPesan("", PesanBisnisKosong)
		return false, nil
	}
	keluar := false
	if h.Ambil(CD+"AppointedADJID") == "" { // 4, 6
		h.TambahPesan("", PesanAdjusterKosong)
		keluar = true
	}
	if h.Ambil(CD+"ConsultantID") == "" { // 5, 7
		h.TambahPesan("", PesanConsultantKosong)
		keluar = true
	}
	if keluar { // 8
		return false, nil
	}
	b := Baris{"Type": "", "PersenRNM": h.Ambil(TM + "RNMShare"), "pxCreateOpName": k.Pelaku,
		"pxCreateOperator": k.Pelaku, "pxCreateDateTime": k.Waktu(), "Payable": h.Ambil(CD + "Payable"),
		"PayableTo": h.Ambil(CD + "PayableTo"), "NoAccount": h.Ambil(RC + "NoAccount"),
		"NameOfBank": h.Ambil(RC + "NameOfBank"), "BranchOfBank": h.Ambil(RC + "BranchOfBank"), "DirectToKasir": "true"}
	if t := h.Ambil(CD + "TypeDeductible"); t == "" || t == "0" { // 13
		b["IndividualRiskType"], b["IndividualRiskPercentage"], b["IndividualRiskValue"] = "3", "0", ""
	}
	n := h.TambahBaris(DaftarAdjustment, b) // 12
	h.Setel("AktifButton", "1")
	h.SetelDaftar(JalurAdj(n, AnakSpreadIn), SalinDaftar(h.AmbilDaftar(DaftarSpreading)))
	h.SetelDaftar(JalurAdj(n, AnakXOLLama), SalinDaftar(h.AmbilDaftar(DaftarXOL)))
	h.SetelDaftar(JalurAdj(n, AnakSpreadOut), SalinDaftar(h.AmbilDaftar(DaftarBreakQS)))
	k.Riwayat(h, TeksTambahAkseptasi)                                          // 14
	h.SetelDaftar(JalurAdj(n, AnakXOL), SalinDaftar(h.AmbilDaftar(DaftarXOL))) // 15
	var kal Kalkulator
	var acc []Baris
	for _, r := range h.AmbilDaftar(DaftarClaimAmount) { // 16
		var t Baris
		for _, a := range acc {
			if a["Currency"] == r["Currency"] {
				t = a
			}
		}
		if t == nil {
			t = r.Salin()
			t["CNPFlagOuts"] = ""
			acc = append(acc, t)
			continue
		}
		for _, p := range []string{"AdjusterFee", "ClaimAmountCedant", "CNPOthersFee", "Salvage", "USD", "Value"} {
			t[p] = Teks(kal.Tambah(kal.B(t, p), kal.B(r, p)))
		}
	}
	h.SetelDaftar(JalurAdj(n, AnakClaimAccept), acc)
	var loss []Baris
	for _, l := range h.AmbilDaftar(DaftarLossAlloc) { // 17
		var t Baris
		for _, a := range loss {
			if a["Currency"] == l["Currency"] && a["TreatyName"] == l["TreatyName"] {
				t = a
			}
		}
		if t == nil {
			t = l.Salin()
			t["CNPFlagOuts"] = ""
			loss = append(loss, t)
			continue
		}
		for _, p := range []string{"AdjusterFee", "CNPOthersFee", "Salvage", "ClaimAmountAdjust"} {
			t[p] = Teks(kal.Tambah(kal.B(t, p), kal.B(l, p)))
		}
	}
	h.SetelDaftar(JalurAdj(n, AnakLossAlloc), loss)
	if err := kal.Galat(); err != nil {
		return false, err
	}
	// 18 CountLossAllocation_act dengan halaman akseptasi sebagai primary - rumusnya menulis jalur absolut
	// `pyWorkPage.ClaimData`, sehingga yang dihitung ulang alokasi XoL tingkat KLAIM (ditiru apa adanya).
	return true, HitungXOL(k, h, m, HitungXOLMod, 0, "")
}

// DeleteAkseptasi = `DeleteAkseptasi_Act` (tombol Delete grid Acceptation List, nonaktif bila AcceptanceStatus terisi):
// baris dibuang (1), AktifButton = 0 (2); Spreading Claim Out tingkat klaim disusun ulang dari seluruh akseptasi (3) -
// medan turunan `SusunSpreadingAkseptasiKlaim`.
func DeleteAkseptasi(h *Halaman, idx int) error {
	if _, err := adj(h, idx); err != nil {
		return err
	}
	h.HapusBaris(DaftarAdjustment, idx)
	h.Setel("AktifButton", "0")
	return nil
}

// SusunSpreadingAkseptasiKlaim - grid "Spreading Claim Out" tab Acceptation (`ClaimData.SpreadingAdjustment` /
// `SpreadingAdjustmentQS`): Σ Spreading In / Out seluruh akseptasi per mata uang x treaty (DeleteAkseptasi_Act langkah 3).
// ⚠️ CountSpreadingXOL (AdjClaimCNP_Act langkah 16) menyusunnya dari akseptasi yang sedang disunting saja; di sini satu
// rumus (Σ seluruh akseptasi) sebagai medan turunan (PARITAS).
func SusunSpreadingAkseptasiKlaim(h *Halaman) error {
	var kal Kalkulator
	gabung := func(anak, tujuan string) {
		var out []Baris
		for i := range h.AmbilDaftar(DaftarAdjustment) {
			for _, s := range h.AmbilDaftar(JalurAdj(i+1, anak)) {
				var t Baris
				for _, o := range out {
					if o["Currency"] == s["Currency"] && o["TreatyType"] == s["TreatyType"] {
						t = o
					}
				}
				if t == nil {
					t = s.Salin()
					out = append(out, t)
					continue
				}
				for _, p := range []string{"ClaimSpreaded", "AdjusterFee", "PremiumSpreaded", "CNPOthersFee", "Salvage"} {
					t[p] = Teks(kal.Tambah(kal.B(t, p), kal.B(s, p)))
				}
				t["TotalClaim"] = Teks(kal.Kurang(kal.Tambah(kal.Tambah(kal.B(t, "ClaimSpreaded"), kal.B(t, "AdjusterFee")),
					kal.B(t, "CNPOthersFee")), kal.B(t, "Salvage")))
			}
		}
		h.SetelDaftar(tujuan, out)
	}
	gabung(AnakSpreadIn, DaftarSpreadAdj)
	gabung(AnakSpreadOut, DaftarSpreadAdjQS)
	return kal.Galat()
}

// ---------------------------------------------------------------- penyesuaian XoL akseptasi

// AdjClaimCNP = `AdjClaimCNP_Act` (param idx; baris akseptasi `n`): Claim Amount satu layer XOL Allocation akseptasi
// diubah, layer terakhir mata uang itu menyerap selisih terhadap Loss Allocation akseptasi, reinstatement premium per
// layer (11), lalu Spreading In / Out = (alokasi - Previously Calculated) x share (15), CountSpreadingXOL (16),
// ProtectNilaiClaim (17). Hardcode CLMNP-232 / CLMNP-861 / CLMNP-975 dibuang (OQ-CNP-03); langkah 12-14 (Summary XOL
// tingkat klaim) = medan turunan.
func AdjClaimCNP(h *Halaman, n, idx int) error {
	b, err := adj(h, n)
	if err != nil {
		return err
	}
	xol := h.AmbilDaftar(JalurAdj(n, AnakXOL))
	if idx < 1 || idx > len(xol) {
		return ErrBarisTidakAda
	}
	var kal Kalkulator
	r := xol[idx-1]
	pct := kal.B(r, "ClaimPercentage")
	var adjV *apd.Decimal // 2
	if r["TreatyType"] == TreatyUR {
		adjV = kal.Kurang(kal.B(r, "TotalClaim"), kal.B(r, "ClaimEstimation"))
	} else {
		adjV = kal.Kurang(kal.B(r, "TotalClaim"), kal.Kurang(kal.Tambah(kal.Tambah(kal.B(r, "ClaimEstimation"),
			bagiPersen(&kal, kal.B(r, "AdjusterFee"), pct)), bagiPersen(&kal, kal.B(r, "CNPOthersFee"), pct)),
			bagiPersen(&kal, kal.B(r, "Salvage"), pct)))
	}
	r["AdjClaimValue"] = Teks(adjV)
	r["ClaimAmountAdjust"] = Teks(kal.Tambah(adjV, kal.B(r, "ClaimEstimation")))
	r["ClaimSpreaded"] = Teks(kal.Kali(kal.Tambah(kal.B(r, "ClaimEstimation"), adjV), kal.BagiPega(pct, k100()))) // 3
	r["TotalClaimRNM"] = r["ClaimSpreaded"]
	cur, nilaiLama := r["Currency"], r["ClaimEstimation"]
	totalClaim, totalValue := apd.New(0, 0), apd.New(0, 0)
	totAdj, totOth, totSalv := apd.New(0, 0), apd.New(0, 0), apd.New(0, 0)
	for _, l := range h.AmbilDaftar(JalurAdj(n, AnakLossAlloc)) { // 6
		if l["Currency"] != cur {
			continue
		}
		totalClaim = kal.Tambah(totalClaim, kal.B(l, "ClaimAmountAdjust"))
		totalValue = kal.Kurang(kal.Tambah(kal.Tambah(kal.B(l, "ClaimAmountAdjust"), kal.B(l, "AdjusterFee")),
			kal.B(l, "CNPOthersFee")), kal.B(l, "Salvage")) // baris terakhir, bukan jumlah (ditiru)
		totAdj, totOth, totSalv = kal.B(l, "AdjusterFee"), kal.B(l, "CNPOthersFee"), kal.B(l, "Salvage")
	}
	idxAkhir, totalAl := 0, apd.New(0, 0)
	for i, s := range xol { // 7
		if s["Currency"] == cur && SamaAngka(s["ClaimEstimation"], nilaiLama) {
			jumlah := kal.Tambah(kal.B(s, "ClaimEstimation"), kal.B(s, "AdjClaimValue"))
			s["ClaimSpreaded"] = Teks(kal.BagiPega(kal.Kali(jumlah, kal.B(s, "ClaimPercentage")), k100()))
			s["ClaimAmountAdjust"] = Teks(jumlah)
			s["TotalClaimRNM"] = Teks(kal.Kurang(kal.Tambah(kal.Tambah(kal.B(s, "ClaimSpreaded"), kal.B(s, "AdjusterFee")),
				kal.B(s, "CNPOthersFee")), kal.B(s, "Salvage")))
		}
		if s["Currency"] == cur {
			idxAkhir = i + 1
			totalAl = kal.Tambah(totalAl, kal.Tambah(kal.B(s, "ClaimEstimation"), kal.B(s, "AdjClaimValue")))
		}
	}
	if idxAkhir > 0 { // 8
		L := xol[idxAkhir-1]
		pctL := kal.B(L, "ClaimPercentage")
		p := kal.BagiPega(pctL, k100())
		L["AdjClaimValue"] = Teks(kal.Tambah(kal.B(L, "AdjClaimValue"), kal.Kurang(totalClaim, totalAl)))
		jumlah := kal.Tambah(kal.B(L, "ClaimEstimation"), kal.B(L, "AdjClaimValue"))
		if totalClaim.IsZero() {
			L["ClaimSpreaded"] = "0"
		} else {
			L["ClaimSpreaded"] = Teks(kal.Kali(jumlah, p))
		}
		L["ClaimAmountAdjust"] = Teks(jumlah)
		if L["TreatyType"] == TreatyUR {
			L["TotalClaim"] = L["ClaimAmountAdjust"]
		} else {
			L["TotalClaim"] = Teks(kal.Kurang(kal.Tambah(kal.Tambah(jumlah, bagiPersen(&kal, kal.B(L, "AdjusterFee"), pctL)),
				bagiPersen(&kal, kal.B(L, "CNPOthersFee"), pctL)), bagiPersen(&kal, kal.B(L, "Salvage"), pctL)))
		}
		if totalClaim.IsZero() {
			L["TotalClaim"] = Teks(kal.Kurang(totalValue, kal.B(r, "TotalClaim")))
		}
		spread := kal.B(L, "ClaimSpreaded")
		L["AdjusterFee"] = Teks(kal.Kali(totAdj, p))
		L["CNPOthersFee"] = Teks(kal.Kali(totOth, p))
		L["Salvage"] = Teks(kal.Kali(totSalv, p))
		L["TotalClaimRNM"] = Teks(kal.Kurang(kal.Tambah(kal.Tambah(spread, kal.Kali(totAdj, p)), kal.Kali(totOth, p)),
			kal.Kali(totSalv, p)))
	}
	if kal.B(r, "TotalClaim").IsZero() { // 9
		r["ClaimSpreaded"], r["AdjusterFee"], r["Salvage"], r["CNPOthersFee"] = "0", "0", "0", "0"
	}
	for _, s := range xol { // 11
		if s["TreatyType"] == TreatyUR {
			s["CNPReinstatement"] = "0"
		} else {
			s["CNPReinstatement"] = Teks(kal.Kali(kal.Kali(kal.BagiPega(kal.B(s, "TotalClaim"), kal.B(s, "CNPLimit")),
				kal.B(s, "CNPMDP")), kal.BagiPega(kal.B(s, "CNPPctReinstate"), k100())))
		}
		s["CNPReinstatementRNM"] = Teks(kal.Kali(kal.B(s, "CNPReinstatement"), kal.BagiPega(kal.B(s, "ClaimPercentage"), k100())))
	}
	if err := kal.Galat(); err != nil {
		return err
	}
	if err := SebarAkseptasi(h, n); err != nil { // 15
		return err
	}
	ProtectNilaiClaim(h, n) // 17
	_ = b
	return SusunSpreadingAkseptasiKlaim(h) // 16
}

// SebarAkseptasi = AdjClaimCNP_Act langkah 15: per mata uang Claim Acceptation, Σ layer non-UR XOL Allocation dikurangi
// Previously Calculated (`AlokasiXOLPaid`), lalu Spreading In / Out = nilai x share; Total Claim = spread + fee; RNM Net
// Claim = Total Claim - Premium Spreaded (reinstatement RNM x share).
func SebarAkseptasi(h *Halaman, n int) error {
	var kal Kalkulator
	for _, c := range h.AmbilDaftar(JalurAdj(n, AnakClaimAccept)) {
		cur := c["Currency"]
		adjF, oth, salv, value, reins := apd.New(0, 0), apd.New(0, 0), apd.New(0, 0), apd.New(0, 0), apd.New(0, 0)
		for _, s := range h.AmbilDaftar(JalurAdj(n, AnakXOL)) {
			if s["Currency"] == cur && s["TreatyType"] != TreatyUR {
				adjF, oth, salv = kal.Tambah(adjF, kal.B(s, "AdjusterFee")), kal.Tambah(oth, kal.B(s, "CNPOthersFee")),
					kal.Tambah(salv, kal.B(s, "Salvage"))
				value, reins = kal.Tambah(value, kal.B(s, "ClaimSpreaded")), kal.Tambah(reins, kal.B(s, "CNPReinstatementRNM"))
			}
		}
		for _, s := range h.AmbilDaftar(JalurAdj(n, AnakXOLDibayar)) {
			if s["Currency"] == cur && s["TreatyType"] != TreatyUR {
				adjF, oth, salv = kal.Kurang(adjF, kal.B(s, "AdjusterFee")), kal.Kurang(oth, kal.B(s, "CNPOthersFee")),
					kal.Kurang(salv, kal.B(s, "Salvage"))
				value, reins = kal.Kurang(value, kal.B(s, "ClaimSpreaded")), kal.Kurang(reins, kal.B(s, "CNPReinstatementRNM"))
			}
		}
		for _, anak := range []string{AnakSpreadIn, AnakSpreadOut} {
			for _, s := range h.AmbilDaftar(JalurAdj(n, anak)) {
				if s["Currency"] != cur {
					continue
				}
				p := kal.BagiPega(kal.B(s, "SharePercentage"), k100())
				s["ClaimSpreaded"] = Teks(kal.Kali(value, p))
				s["PremiumSpreaded"] = Teks(kal.Kali(reins, p))
				s["AdjusterFee"] = Teks(kal.Kali(adjF, p))
				s["CNPOthersFee"] = Teks(kal.Kali(oth, p))
				s["Salvage"] = Teks(kal.Kali(salv, p))
				total := kal.Kurang(kal.Tambah(kal.Tambah(kal.B(s, "ClaimSpreaded"), kal.B(s, "AdjusterFee")),
					kal.B(s, "CNPOthersFee")), kal.B(s, "Salvage"))
				s["TotalClaim"] = Teks(total)
				s["NetClaim"] = Teks(kal.Kurang(total, kal.B(s, "PremiumSpreaded")))
			}
		}
	}
	return kal.Galat()
}

// ProtectNilaiClaim = `ProtectNilaiClaim` (akseptasi `n`): FlagErrorKasir 0 / DirectToKasir true, lalu Spreading In
// ber-RNM Net Claim negatif -> FlagErrorKasir 1, DirectToKasir false.
func ProtectNilaiClaim(h *Halaman, n int) {
	b, err := adj(h, n)
	if err != nil {
		return
	}
	b["FlagErrorKasir"], b["DirectToKasir"] = "0", "true"
	var kal Kalkulator
	for _, s := range h.AmbilDaftar(JalurAdj(n, AnakSpreadIn)) {
		if Lebih(apd.New(0, 0), kal.B(s, "NetClaim")) {
			b["FlagErrorKasir"], b["DirectToKasir"] = "1", "false"
			return
		}
	}
}

// NetNegatif - ada Spreading In akseptasi `n` ber-RNM Net Claim negatif (`CekError.CARI5`, ProteksiNilaiClaim).
func NetNegatif(h *Halaman, n int) bool {
	var kal Kalkulator
	for _, s := range h.AmbilDaftar(JalurAdj(n, AnakSpreadIn)) {
		if Lebih(apd.New(0, 0), kal.B(s, "NetClaim")) {
			return true
		}
	}
	return false
}

// SetInterimXOL = `SetInterimXOL_Act` (change Payment Type): CNPIndexInterim = indeks interim akseptasi interim
// (Payment Type 2) terakhir yang sudah disetujui + 1.
func SetInterimXOL(h *Halaman, n int) error {
	b, err := adj(h, n)
	if err != nil {
		return err
	}
	terakhir := 0
	for _, a := range h.AmbilDaftar(DaftarAdjustment) {
		if a["PaymentType"] == "2" && a["AcceptanceStatus"] == "1" {
			terakhir, _ = strconv.Atoi(strings.TrimSpace(a["CNPIndexInterim"]))
		}
	}
	b["CNPIndexInterim"] = strconv.Itoa(terakhir + 1)
	return nil
}

// ---------------------------------------------------------------- payable dan rekening

// MultiMataUang = `FlagCurrency.CARI15` SetPayableTreatyNP_Act langkah 1: ada akseptasi bermata uang lebih dari satu.
func MultiMataUang(h *Halaman) bool {
	for i := range h.AmbilDaftar(DaftarAdjustment) {
		if len(MataUangAkseptasi(h, i+1)) > 1 {
			return true
		}
	}
	return false
}

// MataUangAkseptasi - `.CurencyAdjustment` akseptasi `n` (AddAkseptasiCNP_Act 9-10: mata uang unik Claim Amount;
// medan turunan dari Claim Acceptation tersimpan).
func MataUangAkseptasi(h *Halaman, n int) []Baris {
	var out []Baris
	sudah := map[string]bool{}
	for _, c := range h.AmbilDaftar(JalurAdj(n, AnakClaimAccept)) {
		k := c["CurrencyID"] + "#" + c["Currency"]
		if sudah[k] {
			continue
		}
		sudah[k] = true
		out = append(out, Baris{"CurrencyID": c["CurrencyID"], "Currency": c["Currency"]})
	}
	return out
}

// SalinRekening menulis rekening `r` ke medan akseptasi berakhiran `ke` ("" atau "2") dan ke ReceiverClaim.
func SalinRekening(h *Halaman, b Baris, r RekeningBank, ke string) {
	b["NameOfBank"+ke], b["NoAccount"+ke], b["BranchOfBank"+ke] = r.NameOfBank, r.AccountNo, r.BranchOfBank
	b["SwiftCode"+ke], b["IDOfBank"+ke], b["Currency"+ke] = r.SwiftCode, r.IDOfBank, r.Currency
	h.Setel(RC+"NameOfBank"+ke, r.NameOfBank)
	h.Setel(RC+"NoAccount"+ke, r.AccountNo)
	h.Setel(RC+"BranchOfBank"+ke, r.BranchOfBank)
	h.Setel(RC+"SwiftCode"+ke, r.SwiftCode)
	h.Setel(RC+"IDOfBank"+ke, r.IDOfBank)
	h.Setel(RC+"Currency"+ke, r.Currency)
}

// SetPayableNP = `SetPayableTreatyNP_Act` dengan Posisi "Acc" (change Payable To, langkah 7 + 9) atau tanpa Posisi
// (change Specify, langkah 3). Payable 1 = ceding, 2 = SOB, 3 = lainnya (penerima `ReceiverClaim.Name`). Rekening dari
// BANKACCOUNT klien (satu mata uang: GetDataBankAccount_sql CLIENTID + mata uang Claim Acceptation pertama; multi mata
// uang: GetDatabyClientName), rekening tanpa nomor / mata uang dibuang; rekening 1 dan 2 dari baris hasil 1 dan 2.
func SetPayableNP(k *Konteks, h *Halaman, n int, posisi string) error {
	b, err := adj(h, n)
	if err != nil {
		return err
	}
	multi := MultiMataUang(h)
	curAcc := ""
	if acc := h.AmbilDaftar(JalurAdj(n, AnakClaimAccept)); len(acc) > 0 {
		curAcc = acc[0]["CurrencyID"]
	}
	payable := h.Ambil(CD + "Payable")
	if posisi != "Acc" && posisi != "AccEdit" { // 3
		b["Payable"], b["PayableTo"] = payable, h.Ambil(CD+"PayableTo")
		var rek []RekeningBank
		if multi {
			rek, err = k.Acuan.RekeningBankNama(k.Ctxt(), h.Ambil(CD+"PayableTo"), "")
		} else {
			rek, err = k.Acuan.RekeningBankNama(k.Ctxt(), h.Ambil(CD+"PayableTo"), curAcc)
		}
		if err != nil {
			return err
		}
		rek = RekeningSah(rek)
		if len(rek) > 0 {
			SalinRekening(h, b, rek[0], "")
		}
		if multi && len(rek) > 1 {
			SalinRekening(h, b, rek[1], "2")
		}
		if payable == "3" { // 6
			h.Setel(RC+"Name", h.Ambil(CD+"PayableTo"))
		}
		return nil
	}
	// 7 (Acc)
	klien := h.Ambil(TM + "LeadingReinsSourceID")
	h.Setel(CD+"PayableTo", h.Ambil(TM+"LeadingReinsSource"))
	b["PayableTo"], b["Payable"] = h.Ambil(TM+"LeadingReinsSource"), payable
	if payable == "1" {
		klien = h.Ambil(TM + "CedingID")
		h.Setel(CD+"PayableTo", h.Ambil(TM+"Ceding"))
		b["PayableTo"] = h.Ambil(TM + "Ceding")
	}
	var rek []RekeningBank
	if multi {
		rek, err = k.Acuan.RekeningBankKlien(k.Ctxt(), klien)
	} else {
		rek, err = k.Acuan.RekeningBank(k.Ctxt(), klien, curAcc)
	}
	if err != nil {
		return err
	}
	rek = RekeningSah(rek)
	if len(rek) > 0 { // 7.6
		SalinRekening(h, b, rek[0], "")
	}
	if multi && len(rek) > 1 { // 9.6
		SalinRekening(h, b, rek[1], "2")
	}
	if payable == "3" { // 5-6, 7.8
		h.Setel(CD+"PayableTo", h.Ambil(RC+"Name"))
		b["PayableTo"] = h.Ambil(RC + "Name")
		for _, p := range []string{"NameOfBank", "NoAccount", "BranchOfBank", "IDOfBank", "SwiftCode", "Currency"} {
			b[p] = h.Ambil(RC + p)
		}
	}
	return nil
}

// RekeningSah = SetPayableTreatyNP_Act 3.4 / 7.5: rekening tanpa ACCOUNTNO atau CURRENCY dibuang.
func RekeningSah(rek []RekeningBank) []RekeningBank {
	var out []RekeningBank
	for _, r := range rek {
		if r.AccountNo != "" && r.Currency != "" {
			out = append(out, r)
		}
	}
	return out
}

// SetAccountNo = `SetAccoutNo_Act` (change Name of Bank / Payable To): rekening Spreading In akseptasi `n` per mata
// uang - baris bermata uang rekening 2 memakai rekening 2, selainnya rekening 1.
// ⚠️ DISEDERHANAKAN `[inferensi]` (PARITAS, OQ): activity mencocokkan baris `Result` (halaman RDB SetPayable yang hidup
// di requestor) lewat indeks posisi dan cabang (3.4-3.11); maksudnya - rekening per mata uang Spreading In - yang dibangun.
func SetAccountNo(h *Halaman, n int) error {
	b, err := adj(h, n)
	if err != nil {
		return err
	}
	for _, s := range h.AmbilDaftar(JalurAdj(n, AnakSpreadIn)) {
		if b["Currency2"] != "" && s["Currency"] == b["Currency2"] && b["NoAccount2"] != "" {
			s["NoAccount"], s["IDOfBank"] = b["NoAccount2"], b["IDOfBank2"]
		} else {
			s["NoAccount"], s["IDOfBank"] = b["NoAccount"], b["IDOfBank"]
		}
	}
	return nil
}

// ---------------------------------------------------------------- proteksi kirim komite

// CekKomite - halaman `CekError` ProteksiSendKomiteCNP_Act (gerbang tombol "Send Claim to Committee" popup KomiteCLMNP).
type CekKomite struct {
	MelebihiEstimasi bool   // CARI1
	PolicyNo         string // CARI2
	EmailPelaku      string // CARI3
	RekeningKosong   bool   // CARI4
	NetNegatif       bool   // CARI5
	MataUangKosong   bool   // CARI6
	RekeningSalah    bool   // CARI8
	IDBankKosong     bool   // CARI9
	HarusCancel      bool   // CARI10
}

// Lolos = gerbang pyDisabledWhen tombol Send Claim to Committee (Section KomiteCLMNP).
func (c CekKomite) Lolos(h *Halaman) bool {
	return h.Ambil(CD+"Payable") != "" && !c.MelebihiEstimasi && c.PolicyNo != "" && c.EmailPelaku != "" &&
		!c.RekeningKosong && !c.NetNegatif && !c.MataUangKosong && !c.RekeningSalah && !c.IDBankKosong &&
		!c.HarusCancel && h.Ambil(CD+"Occupation") != ""
}

// ProteksiKirimKomite = `ProteksiSendKomiteCNP_Act` (akseptasi `n`) + `ProteksiNilaiClaim`.
//
// PERBAIKAN OQ-CNP-05 butir 3 (PARITAS `[penyimpangan sadar]`):
//   - "Error No Account" (7.3 / 7.4): syarat XML memuat `Curr1 == .Currency` DAN `.Currency != Curr1` sehingga tidak
//     pernah aktif. Konjungsi yang saling meniadakan dibuang; sisanya dipertahankan: baris Spreading In bermata uang
//     rekening 1 (2) yang nomor rekeningnya bukan rekening 1 (2) dan bukan rekening lainnya -> error.
//   - `pyWorkPage.IsError` (Payment Type kosong = 2) DIRESET setiap proteksi dijalankan (XML tidak pernah meresetnya).
func ProteksiKirimKomite(k *Konteks, h *Halaman, n int, email string) (CekKomite, error) {
	b, err := adj(h, n)
	if err != nil {
		return CekKomite{}, err
	}
	h.BersihkanPesan() // 1
	h.Setel("IsError", "")
	c := CekKomite{PolicyNo: h.Ambil(CD + "PolicyData.PolicyNo"), EmailPelaku: email}
	h.Setel(JalurAdjProp(n, "DataCommitteeTreaty.CircumCauseOfLoss"), h.Ambil(CD+"CNPCircumtances")) // 2
	k.Riwayat(h, TeksKirimKomite)                                                                    // 3
	semuaUR := len(h.AmbilDaftar(JalurAdj(n, AnakXOL))) > 0                                          // 4
	for _, s := range h.AmbilDaftar(JalurAdj(n, AnakXOL)) {
		if s["TreatyName"] != TreatyUR {
			semuaUR = false
		}
	}
	if semuaUR && b["PaymentType"] != PaymentTypeCancellation { // 4.1, 5
		c.HarusCancel = true
		h.TambahPesan("", PesanHarusCancel)
	}
	for _, s := range h.AmbilDaftar(JalurAdj(n, AnakSpreadIn)) { // 7 (perbaikan butir 3)
		lain := func(no string) bool { return s["NoAccount"] != no }
		if b["Currency"] != "" && s["Currency"] == b["Currency"] && lain(b["NoAccount"]) && lain(b["NoAccount2"]) {
			c.RekeningSalah = true
		}
		if b["Currency2"] != "" && s["Currency"] == b["Currency2"] && lain(b["NoAccount2"]) && lain(b["NoAccount"]) {
			c.RekeningSalah = true
		}
	}
	if b["NoAccount"] == "" || b["NameOfBank"] == "" { // 8
		b["BranchOfBank"], b["NameOfBank"], b["NoAccount"] = h.Ambil(RC+"BranchOfBank"), h.Ambil(RC+"NameOfBank"),
			h.Ambil(RC+"NoAccount")
	}
	if b["PaymentType"] == "" { // 9
		h.Setel("IsError", "2")
	}
	kosongRek := b["NoAccount"] == "" || b["NameOfBank"] == ""
	if b["AcceptanceStatus"] == "" && (kosongRek || b["BranchOfBank"] == "") { // 10
		c.RekeningKosong = true
	}
	c.MataUangKosong = b["Currency"] == ""        // 11
	c.IDBankKosong = b["IDOfBank"] == ""          // 12
	if b["AcceptanceStatus"] == "" && kosongRek { // 13
		h.TambahPesan("", PesanRekeningKosong)
	}
	if h.Ambil(CD+"Occupation") == "" { // 14
		h.TambahPesan("", PesanOkupasiKosong)
	}
	if h.Ambil("IsError") == "2" { // 15
		h.TambahPesan("", PesanPilihPayment)
	}
	if c.PolicyNo == "" { // 16
		h.TambahPesan("", PesanPolisKomite)
	}
	if c.MataUangKosong { // 17
		h.TambahPesan("", PesanMataUangKosong)
	}
	if c.IDBankKosong { // 18
		h.TambahPesan("", PesanIDBankKosong)
	}
	if c.EmailPelaku == "" { // 19
		h.TambahPesan("", PesanEmailKosong)
	}
	if c.RekeningSalah { // 20
		h.TambahPesan("", PesanRekeningSalah)
	}
	if h.Ambil(CD+"Payable") == "" { // 22
		h.TambahPesan("", PesanPilihPayable)
	}
	melebihi, err := MelebihiEstimasi(h) // 23
	if err != nil {
		return c, err
	}
	if melebihi {
		c.MelebihiEstimasi = true
		h.TambahPesan("", PesanMelebihiEstimasi)
	}
	c.NetNegatif = NetNegatif(h, n) // ProteksiNilaiClaim
	if c.NetNegatif {
		h.TambahPesan("", PesanNetNegatif)
	}
	return c, nil
}

// JalurAdjProp - jalur satu properti baris AdjustmentList(n) (`AdjustmentList(n).prop`).
func JalurAdjProp(n int, prop string) string {
	return fmt.Sprintf("%s(%d).%s", DaftarAdjustment, n, prop)
}

// MelebihiEstimasi = ProteksiSendKomiteCNP_Act langkah 23: per mata uang Summary XOL, total klaim Loss Allocation
// ber-To XOL (`TempTotalClaim` CountLossAllocation_act: Value + fee - Salvage; TPL tidak ada di halaman itu = 0) melebihi
// CNPTotalClaim Summary XOL.
func MelebihiEstimasi(h *Halaman) (bool, error) {
	if err := SusunSummaryXOL(h); err != nil {
		return false, err
	}
	var kal Kalkulator
	for _, s := range h.AmbilDaftar(DaftarSummaryXOL) {
		total := apd.New(0, 0)
		ada := false
		for _, l := range h.AmbilDaftar(DaftarLossAlloc) {
			if l["CNPFlagXOL"] != "true" || l["Currency"] != s["Currency"] {
				continue
			}
			ada = true
			total = kal.Tambah(total, kal.Kurang(kal.Tambah(kal.Tambah(kal.B(l, "ClaimAmountAdjust"),
				kal.B(l, "CNPOthersFee")), kal.B(l, "AdjusterFee")), kal.B(l, "Salvage")))
		}
		if ada && Lebih(total, kal.B(s, "CNPTotalClaim")) {
			return true, kal.Galat()
		}
	}
	return false, kal.Galat()
}

// ---------------------------------------------------------------- Kasir

// KonfigKasir - nilai tetap muatan Kasir (HitServiceToKasir_Act 9.3: CompanyName, LjtdId, LdcId, StsAp, StsSyariah,
// Deductible, KaliDeduct). Isinya dibaca dari konfigurasi ber-embed `services/konfigurasi/kasir.json` (pola Komite Claim
// Prop), bukan literal di sini.
type KonfigKasir struct {
	CompanyName, LjtdID, LdcID, LdcIDSyariah, StsAp string
}

// MuatanKasir - satu baris `TAllPaymentData` (Java langkah 9.6).
type MuatanKasir struct {
	NoTrans, NoKlaim, LbuID, NoPolis, AcceptType, Kepada, AccountNo, TglAksep, Nett string
	Deductible, KaliDeduct                                                          string
	StsSyariah, CompanyName, LjtdID, LdcID, StsAp, LkuID, LbgID, TglBolehBayar      string
	Email, UserInput                                                                string
}

var hanyaDigit = regexp.MustCompile(`[^0-9]`)

// PanjangNoAksepCNP = HitServiceToKasir_Act langkah 9 baris 2: nomor akseptasi Non Prop sepanjang 23 / 24 karakter.
func PanjangNoAksepCNP(no string) bool { n := len([]rune(no)); return n == 23 || n == 24 }

// TanggalBolehBayar = HitServiceToKasir_Act langkah 9.5: tanggal akseptasi + 1 bulan, +1 bulan lagi bila tanggal > 25,
// tanggal 01 bila > 25; hasil "dd-MM-yyyy".
//
// PERBAIKAN OQ-CNP-05 butir 5 (PARITAS `[penyimpangan sadar]`): XML menaikkan tahun hanya bila bulan hasil = 01 DAN
// bulan SEKARANG = 12 (`@CurrentDate("MM")`), dan bulan 14 (Des tgl > 25) menjadi "1" tanpa menaikkan tahun bila bulan
// sekarang bukan Desember. Di sini aritmetika tanggal sungguhan: tahun ikut bergulir bersama bulan.
func TanggalBolehBayar(akseptasi time.Time) string {
	hari := akseptasi.Day()
	tambahBulan := 1
	if hari > 25 {
		tambahBulan = 2
		hari = 1
	}
	awal := time.Date(akseptasi.Year(), akseptasi.Month(), 1, 0, 0, 0, 0, time.UTC).AddDate(0, tambahBulan, 0)
	return fmt.Sprintf("%02d-%02d-%d", hari, int(awal.Month()), awal.Year())
}

// SusunMuatanKasir = HitServiceToKasir_Act langkah 9.1-9.6 untuk setiap Spreading In akseptasi `n`: Nett = Total Claim -
// Premium Spreaded baris spreading. ⚠️ LbuID = `pyWorkPage.OfferFacIn.QuotationData.BusinessOldId` (langkah 9.3) - properti
// tanpa penulis di korpus Non Prop, jadi kosong (persis XML; PARITAS, OQ).
func SusunMuatanKasir(h *Halaman, n int, email, pelaku string, cfg KonfigKasir) ([]MuatanKasir, error) {
	b, err := adj(h, n)
	if err != nil {
		return nil, err
	}
	var kal Kalkulator
	tgl, _ := UraiTanggal(b["AcceptedDate"])
	tglAksep, boleh := "", ""
	if !tgl.IsZero() {
		tglAksep, boleh = tgl.Format("02-01-2006"), TanggalBolehBayar(tgl)
	}
	user := b["pxCreateOperator"]
	if user == "" {
		user = pelaku
	}
	var out []MuatanKasir
	for _, s := range h.AmbilDaftar(JalurAdj(n, AnakSpreadIn)) {
		akun := hanyaDigit.ReplaceAllString(s["NoAccount"], "")
		if akun == "" {
			akun = strings.ReplaceAll(b["NoAccount"], "-", "")
		}
		out = append(out, MuatanKasir{NoTrans: b["AcceptedNo"], NoKlaim: h.Ambil(CD + "NoClaim"),
			LbuID: h.Ambil(OQ + "BusinessOldId"), NoPolis: h.Ambil(CD + "PolicyData.PolicyNo"),
			AcceptType: b["PaymentType"], Kepada: b["PayableTo"], AccountNo: akun, TglAksep: tglAksep,
			Nett: Teks(kal.Kurang(kal.B(s, "TotalClaim"), kal.B(s, "PremiumSpreaded"))), Deductible: "0", KaliDeduct: "0", StsSyariah: "0",
			CompanyName: cfg.CompanyName, LjtdID: cfg.LjtdID, LdcID: cfg.LdcID, StsAp: cfg.StsAp,
			LkuID: s["CurrencyID"], LbgID: s["IDOfBank"], TglBolehBayar: boleh, Email: email, UserInput: user})
	}
	return out, kal.Galat()
}
