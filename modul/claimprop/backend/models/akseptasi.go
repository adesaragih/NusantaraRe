package models

// Untuk apa berkas ini: PORT ACTIVITY BARIS ADJUSTMENT (tiket 08, 10, 11; Section InputAcceptation_Adjs dan
// AdjustmentDetail): tambah / hapus baris, mata uang, gross dan nilai adjustment, spreading adjustment, roster komite,
// payable dan rekening bank, DLA ceding / SOB.
//
// ⚠️ Beberapa activity menghapus baris PageList sambil mengulanginya maju (`Property-Remove` di dalam REPEAT EMBEDDED:
// SetNameCurrency_Act 12-13, CountGrossAdjTreaty_Act 5.1, FilterLossAllocation_Act 2.1, AddAdjustment_Act 12). Urutan
// iterasi Pega atas daftar yang menyusut tidak terbaca dari XML; yang dibangun adalah MAKSUD langkahnya (saring baris
// yang tidak cocok / hapus baris baru) - `[dugaan]`, dicatat di PARITAS.

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/cockroachdb/apd/v3"
)

// Pesan VERBATIM.
const (
	PesanAdjusterLagi       = "Please select the Adjuster and Consultant again."                // AddAdjustment_Act 9
	PesanPeriodeTBA         = "Policy period is TBA. Please verify the dates."                  // AddAdjustment_Act 10
	PesanSpreadingAdjKosong = "Spreading AdjustmentList Can't Null"                             // SetNameCurrency_Act 6
	PesanAdjLewatEstimasi   = "Adjustment RNM should not be greater than estimation RNM"        // CountValueADJTreaty_Act 2
	PesanSalvageMinus       = "Adjustment RNM for Salvage should be minus"                      // ErrMsg1
	PesanAdjNol             = "Adjustment RNM can not be filled by Zero"                        // ErrMsg3
	PesanTotalAdjLewat      = "Total Adjustment  RNM should not be greater than estimation RNM" // ErrMsg4
	PesanPilihPayable       = "Please choose payable first"                                     // ProteksiInitialandDate_Act 2
	PesanPilihTipeBayar     = "Please choose Payment Type first"                                // ErrMsg
	PesanBankNull           = "Data Bank Account Can't NULL"                                    // ErrMsg1
	PesanKomiteMasihJalan   = "Can not close claim, there is adjustment in comitee!"            // CloseClaimProp 2.1
	PesanKasirBelumSukses   = "Cannot close claim, there is a direct to cashier that has not been successful."
	PesanPremiBelumLunas    = "Akseptasi tidak dapat dilanjutkan dikarenakan Premi belum Lunas" // CekPremiLunas_Act 6
	PesanCetakDLA           = "Please Print DLA"
)

// Teks riwayat VERBATIM.
const (
	TeksTambahAdjustment = "Add Adjustment"
	TeksHitungAdjustment = "Count Value Adjustment"
	TeksNilaiAdjustment  = "Input Data Adjustment and Value"
	TeksKomiteAkseptasi  = "Send Adjustment to Committe (Acceptation)"
	TeksTutupKlaim       = "Finish Adjustment (Close Claim)"
	TeksCetakDLA         = "Print DLA"
)

// StatusKasirSukses - `GetStatusKasir_Act` langkah 5 / `HitServiceToKasir_Act` 13.5.
const StatusKasirSukses = "Akseptasi Sudah Masuk ke Kasir"

// STSKlaimProp - `Param.STS_KLAIM = "PROP"` roster komite (SetKomiteTreaty_ACT 3, AddKomiteTreatyChild_ACT 20).
const STSKlaimProp = "PROP"

// JalurAdj - jalur baris adjustment ke-n.
func JalurAdj(n int, anak string) string { return JalurAnak(DaftarAdjustment, n, anak) }

// adj mengambil baris adjustment ke-n.
func adj(h *Halaman, n int) (Baris, error) { return baris(h, DaftarAdjustment, n) }

// mataUangEstimasi - Kurs.pxResults AddAdjustment_Act langkah 2-3: mata uang estimasi unik (CurrencyID#Currency,
// kemunculan pertama) beserta kurs baris pertamanya. Ini `AdjustmentList(n).CurencyAdjustment` - sumber dropdown
// Currency AdjustmentDetail; turunan, tidak disimpan.
func mataUangEstimasi(h *Halaman) []Baris {
	var out []Baris
	sudah := map[string]bool{}
	for _, e := range h.AmbilDaftar(DaftarEstimasi) {
		k := e["CurrencyID"] + "#" + e["Currency"]
		if sudah[k] {
			continue
		}
		sudah[k] = true
		out = append(out, Baris{"CurrencyID": e["CurrencyID"], "Currency": e["Currency"], "KursValue": e["KursValue"]})
	}
	return out
}

// SusunMataUangAdjustment mengisi `AdjustmentList(n).CurencyAdjustment` setiap baris (turunan).
func SusunMataUangAdjustment(h *Halaman) {
	mu := mataUangEstimasi(h)
	for i := range h.AmbilDaftar(DaftarAdjustment) {
		h.SetelDaftar(JalurAdj(i+1, "CurencyAdjustment"), SalinDaftar(mu))
	}
}

// AddAdjustment meniru `Activity/AddAdjustment_Act.xml` (tombol Add grid "Acceptation List").
func AddAdjustment(k *Konteks, h *Halaman) {
	h.BersihkanPesan() // 1
	b := Baris{        // 4-5
		"pxCreateDateTime": k.Waktu(), "pxCreateOperator": k.Pelaku, "pxCreateOpName": k.Pelaku,
		"PersenRNM": h.Ambil(TM + "RNMShareP"), "TotalEstimasiValue": h.Ambil(CD + "TotalEstimasiIDR"),
		"DirectToKasir": "true",
	}
	switch tipe := h.Ambil(CD + "TypeDeductible"); { // 6-8
	case tipe == "1" || tipe == "2":
		b["IndividualRiskType"] = tipe
		b["IndividualRiskPercentage"] = h.Ambil(CD + "Amount")
		b["IndividualRiskValue"] = h.Ambil(CD + "DeductibleValue")
	case tipe == "" || tipe == "0":
		b["IndividualRiskType"] = "3"
		b["IndividualRiskPercentage"] = "0"
		b["IndividualRiskValue"] = ""
	}
	n := h.TambahBaris(DaftarAdjustment, b)
	h.SetelDaftar(JalurAdj(n, AnakSpreadAdj), SalinDaftar(h.AmbilDaftar(DaftarSpreading)))
	h.SetelDaftar(JalurAdj(n, AnakLossAllocation), SalinDaftar(h.AmbilDaftar(DaftarLossAlloc)))
	h.SetelDaftar(JalurAdj(n, AnakQuotaShare), SalinDaftar(h.AmbilDaftar(DaftarBreakQS)))
	h.Setel("AktifButton", "1")
	err := ""
	if (h.Ambil(CD+"ConsultantID") == "" || h.Ambil(CD+"AppointedADJID") == "") && n == 1 { // 9
		err = PesanAdjusterLagi
	}
	if h.Ambil(CD+"PeriodPolicyTBA") == "true" { // 10
		err = PesanPeriodeTBA
	}
	if err != "" { // 11-12: pesan, baris baru dibuang, tombol Add aktif lagi
		h.TambahPesan(CD+"AkseptasiList", err)
		h.Setel("AktifButton", "0")
		h.HapusBaris(DaftarAdjustment, n)
		return
	}
	k.Riwayat(h, TeksTambahAdjustment) // 13
	SusunMataUangAdjustment(h)
}

// DeleteAjsutment meniru `Activity/DeleteAjsutment_Act.xml` (tombol Delete baris adjustment; disabled bila
// `.IsKomite=1 || .IsSubjectivity = true`).
func DeleteAjsutment(h *Halaman, idx int) error {
	b, err := adj(h, idx)
	if err != nil {
		return err
	}
	if b["IsKomite"] == "1" || b["IsSubjectivity"] == "true" {
		return ErrBarisBeku
	}
	h.BersihkanPesan()
	h.HapusBaris(DaftarAdjustment, idx)
	h.Setel("AktifButton", "0")
	return nil
}

// SetNameCurrency meniru `Activity/SetNameCurrency_Act.xml` (change Currency AdjustmentDetail): nama, kurs (bila
// kosong), gross / nilai dari estimasi bermata uang itu, salinan spreading klaim yang disaring mata uangnya, loss
// allocation unik, lalu CountGrossAdjTreaty_Act dan SetPayableTreaty_Act.
func SetNameCurrency(k *Konteks, h *Halaman, idx int, currID string) error {
	b, err := adj(h, idx)
	if err != nil {
		return err
	}
	b["CurrencyID"] = currID // 1
	nama, err := k.Acuan.NamaMataUang(k.Ctxt(), currID)
	if err != nil {
		return err
	}
	b["Currency"] = nama    // 3-4
	if b["KursIDR"] == "" { // 5-6
		kurs, err := k.Acuan.KursStandar(k.Ctxt(), currID)
		if err != nil {
			return err
		}
		b["KursIDR"] = kurs
	}
	var kal Kalkulator // 7-8
	g, gi, ev, evi := apd.New(0, 0), apd.New(0, 0), apd.New(0, 0), apd.New(0, 0)
	for _, e := range h.AmbilDaftar(DaftarEstimasi) {
		if e["CurrencyID"] != currID {
			continue
		}
		g = kal.Tambah(g, kal.B(e, "GrossEstimationPct"))
		gi = kal.Tambah(gi, kal.B(e, "ConvertGrossEstimasi"))
		ev = kal.Tambah(ev, kal.B(e, "EstimationValue"))
		evi = kal.Tambah(evi, kal.B(e, "ConvertValue"))
	}
	if err := kal.Galat(); err != nil {
		return err
	}
	b["GrossAdjustment"], b["GrossAdjustmentIDR"] = Teks(g), Teks(gi)
	b["AdjustmentValue"], b["ValueAdjustment"] = Teks(ev), Teks(evi)
	k.Riwayat(h, TeksHitungAdjustment) // 8 + 18
	saring := func(d []Baris) []Baris {
		var out []Baris
		for _, x := range d {
			if x["CurrencyID"] == currID {
				out = append(out, x.Salin())
			}
		}
		return out
	}
	h.SetelDaftar(JalurAdj(idx, AnakSpreadAdj), saring(h.AmbilDaftar(DaftarSpreading))) // 9, 12
	h.SetelDaftar(JalurAdj(idx, AnakQuotaShare), saring(h.AmbilDaftar(DaftarBreakQS)))  // 9, 13
	var la []Baris                                                                      // 10-11
	sudah := map[string]bool{}
	for _, x := range h.AmbilDaftar(DaftarLossAlloc) {
		kunci := x["TreatyType"] + "#" + x["CurrencyID"] + "#" + x["SharePercentage"]
		if !sudah[kunci] {
			sudah[kunci] = true
			la = append(la, x.Salin())
		}
	}
	h.SetelDaftar(JalurAdj(idx, AnakLossAllocation), la)
	if len(h.AmbilDaftar(JalurAdj(idx, AnakSpreadAdj))) == 0 { // 14
		h.TambahPesan("", PesanSpreadingAdjKosong)
	}
	// 15: ClaimData.SpreadingAdjustment(QS) = turunan (HitungSpreadingAdjustmentTotal).
	if err := CountGrossAdjTreaty(k, h, idx); err != nil { // 16
		return err
	}
	return SetPayableTreaty(k, h, idx) // 17
}

// persenLA - SharePercentage baris loss allocation adjustment bermata uang `currID` (baris terakhir menang); kosong =
// semua baris (CountValueADJTreaty_Act 11.1 berprakondisi nonaktif).
func persenLA(h *Halaman, idx int, currID string, semua bool) string {
	p := ""
	for _, x := range h.AmbilDaftar(JalurAdj(idx, AnakLossAllocation)) {
		if semua || x["CurrencyID"] == currID {
			p = x["SharePercentage"]
		}
	}
	return p
}

// CountGrossAdjTreaty meniru `Activity/CountGrossAdjTreaty_Act.xml` (change Gross / Individual Risk): GrossValue,
// ProposeAdjustmentValue, loss allocation adjustment bermata uang adjustment beserta spread-nya, lalu
// CountValueADJTreaty_Act. `@divide(x,100,20)` = persen tanpa pembulatan (spec §4).
func CountGrossAdjTreaty(k *Konteks, h *Halaman, idx int) error {
	b, err := adj(h, idx)
	if err != nil {
		return err
	}
	var kal Kalkulator
	cur := b["CurrencyID"]
	sc := kal.H(h, CD+"ShareCeding")
	pla := kal.Teks("SharePercentage", persenLA(h, idx, cur, false)) // 2
	gross := kal.B(b, "GrossAdjustment")
	persenRNM := kal.B(b, "PersenRNM")
	b["GrossValue"] = Teks(kal.Persen(gross, persenRNM, pla, sc)) // 3-4
	propose := gross
	if b["Type"] == "1" {
		propose = kal.Kurang(gross, kal.B(b, "IndividualRiskValue"))
	}
	b["ProposeAdjustmentValue"] = Teks(propose)
	var la []Baris // 5
	for _, x := range h.AmbilDaftar(JalurAdj(idx, AnakLossAllocation)) {
		if x["CurrencyID"] != cur {
			continue
		}
		spread := kal.Persen(propose, sc, kal.B(x, "SharePercentage"))
		x["ClaimSpreaded"] = Teks(spread)
		x["ClaimEstimation"] = Teks(kal.Persen(spread, persenRNM))
		la = append(la, x)
	}
	if err := kal.Galat(); err != nil {
		return err
	}
	h.SetelDaftar(JalurAdj(idx, AnakLossAllocation), la)
	return CountValueADJTreaty(k, h, idx) // 6 (7 ber-remark)
}

// CountValueADJTreaty meniru `Activity/CountValueADJTreaty_Act.xml`: individual risk, nilai adjustment RNM / IDR,
// penanda `IsError`, pesan validasi, lalu CountSpreadingADJ_Act dan SetKomiteTreaty_ACT.
func CountValueADJTreaty(k *Konteks, h *Halaman, idx int) error {
	b, err := adj(h, idx)
	if err != nil {
		return err
	}
	h.BersihkanPesan() // 1
	h.Setel("IsError", "")
	k.Riwayat(h, TeksNilaiAdjustment) // 2-3
	var kal Kalkulator
	gross := kal.B(b, "GrossAdjustment")
	irPct := kal.B(b, "IndividualRiskPercentage")
	tipe := b["Type"]
	switch b["IndividualRiskType"] {
	case "1": // 4-5
		if tipe == "1" {
			b["IndividualRiskValue"] = Teks(kal.Persen(gross, irPct))
		} else {
			b["IndividualRiskValue"] = "0"
		}
	case "2": // 6-8: TSI = .Value claim amount bermata uang adjustment (baris terakhir)
		if tipe == "1" {
			tsi := "0"
			for _, c := range h.AmbilDaftar(DaftarClaimAmount) {
				if c["CurrencyID"] == b["CurrencyID"] {
					tsi = c["Value"]
				}
			}
			b["IndividualRiskValue"] = Teks(kal.Persen(kal.Teks("TSIValue", tsi), irPct))
		} else {
			b["IndividualRiskValue"] = "0"
		}
	case "3": // 9
		b["IndividualRiskPercentage"] = "0"
	}
	propose := kal.Kurang(gross, kal.B(b, "IndividualRiskValue"))
	if b["IndividualRiskType"] == "1" || b["IndividualRiskType"] == "2" || b["IndividualRiskType"] == "3" {
		b["ProposeAdjustmentValue"] = Teks(propose)
	} else {
		propose = kal.B(b, "ProposeAdjustmentValue")
	}
	sc := kal.H(h, CD+"ShareCeding") // 10-12
	pla := kal.Teks("SharePercentage", persenLA(h, idx, "", true))
	persenRNM := kal.B(b, "PersenRNM")
	nilai := kal.Persen(propose, persenRNM, pla, sc)
	b["AdjustmentValue"] = Teks(nilai)
	b["IndividualRiskRNM"] = Teks(kal.Persen(kal.B(b, "IndividualRiskValue"), persenRNM, pla, sc))
	b["ValueAdjustment"] = Teks(kal.Kali(kal.B(b, "KursIDR"), nilai))
	sum := apd.New(0, 0) // 13: adjustment ber-status selain 2 (termasuk baris ini)
	for _, x := range h.AmbilDaftar(DaftarAdjustment) {
		if x["AcceptanceStatus"] != "2" {
			sum = kal.Tambah(sum, kal.B(x, "ValueAdjustment"))
		}
	}
	totEst := kal.B(b, "TotalEstimasiValue")
	nilaiIDR := kal.B(b, "ValueAdjustment")
	if err := kal.Galat(); err != nil {
		return err
	}
	jalur := func(p string) string { return JalurAdj(idx, p) }
	if Nol(propose) { // 14-15: transisi aktif - keluar
		h.Setel("IsError", "2")
		h.TambahPesan(jalur("ProposeAdjustmentValue"), PesanAdjNol)
		return nil
	}
	if Lebih(sum, totEst) { // 16
		h.Setel("IsError", "2")
	}
	if Lebih(nilaiIDR, totEst) && Lebih(sum, totEst) && tipe == "1" { // 17
		h.TambahPesan(jalur("AdjustmentValue"), PesanAdjLewatEstimasi)
	}
	if propose.Sign() > 0 && tipe == "3" { // 18
		h.TambahPesan(jalur("ProposeAdjustmentValue"), PesanSalvageMinus)
	}
	if !Lebih(sum, totEst) { // 19
		h.Setel("IsError", "")
	} else { // 20-21
		h.Setel("IsError", "3")
		h.TambahPesan(jalur("GrossAdjustment"), PesanTotalAdjLewat)
	}
	CountSpreadingADJ(h, idx) // 22
	return SetKomiteTreaty(k, h, idx)
}

// CountSpreadingADJ meniru `Activity/CountSpreadingADJ_Act.xml` atas baris adjustment `idx`: spread-in baris = nilai
// adjustment x share; spread-out (quota share) = spread-in baris TERAKHIR x share (langkah 6.4 berprakondisi
// nonaktif). Daftar total `ClaimData.SpreadingAdjustment(QS)` turunan - `HitungSpreadingAdjustmentTotal`.
func CountSpreadingADJ(h *Halaman, idx int) {
	b, err := adj(h, idx)
	if err != nil {
		return
	}
	var kal Kalkulator
	nilai := kal.B(b, "AdjustmentValue")
	var terakhir *apd.Decimal
	for _, s := range h.AmbilDaftar(JalurAdj(idx, AnakSpreadAdj)) { // 6
		sp := kal.Persen(nilai, kal.B(s, "SharePercentage"))
		s["ClaimSpreaded"] = Teks(sp)
		terakhir = sp
	}
	if terakhir == nil {
		terakhir = apd.New(0, 0)
	}
	for _, q := range h.AmbilDaftar(JalurAdj(idx, AnakQuotaShare)) { // 7
		q["ClaimSpreaded"] = Teks(kal.Persen(terakhir, kal.B(q, "SharePercentage")))
	}
}

// HitungSpreadingAdjustmentTotal - grid "Spreading Adjustment Total" (Section InputAcceptation_Adjs) =
// CountSpreadingADJ_Act langkah 2-9: salinan SpreadingClaim / SpreadingBreakQS, lalu per mata uang adjustment
// spread = share x total AdjustmentValue adjustment bermata uang itu ber-status selain 2. Turunan, dihitung setiap
// halaman dimuat; hanya bila ada baris adjustment (Pega mengisinya pertama kali saat aksi adjustment).
func HitungSpreadingAdjustmentTotal(h *Halaman) error {
	adjs := h.AmbilDaftar(DaftarAdjustment)
	if len(adjs) == 0 {
		h.SetelDaftar(DaftarSpreadAdj, nil)
		h.SetelDaftar(DaftarSpreadAdjQS, nil)
		return nil
	}
	in := SalinDaftar(h.AmbilDaftar(DaftarSpreading))
	out := SalinDaftar(h.AmbilDaftar(DaftarBreakQS))
	var kal Kalkulator
	for _, mu := range urutanMataUang(adjs) {
		tot := apd.New(0, 0)
		for _, a := range adjs {
			if a["CurrencyID"] == mu["CurrencyID"] && a["AcceptanceStatus"] != "2" {
				tot = kal.Tambah(tot, kal.B(a, "AdjustmentValue"))
			}
		}
		for _, d := range [][]Baris{in, out} {
			for _, s := range d {
				if s["CurrencyID"] == mu["CurrencyID"] {
					s["ClaimSpreaded"] = Teks(kal.Persen(tot, kal.B(s, "SharePercentage")))
				}
			}
		}
	}
	if err := kal.Galat(); err != nil {
		return err
	}
	h.SetelDaftar(DaftarSpreadAdj, in)
	h.SetelDaftar(DaftarSpreadAdjQS, out)
	return nil
}

// SetKomiteTreaty meniru `Activity/SetKomiteTreaty_ACT.xml`: ValueAdjustment = nilai x kurs, roster komite
// EMAILKOMITE (`LIMIT_BOTTOM <= nilai`, STS_KLAIM PROP) ke `.ComiteeClaim`, `.TotalKomite` = cacahnya.
// ComiteeClaim dan TotalKomite turunan (katalog) - `SusunKomiteAdjustment` menyusunnya ulang saat dimuat.
func SetKomiteTreaty(k *Konteks, h *Halaman, idx int) error {
	b, err := adj(h, idx)
	if err != nil {
		return err
	}
	var kal Kalkulator
	b["ValueAdjustment"] = Teks(kal.Kali(kal.B(b, "AdjustmentValue"), kal.B(b, "KursIDR"))) // 2
	if err := kal.Galat(); err != nil {
		return err
	}
	return SusunKomiteAdjustment(k, h, idx, nil)
}

// ---------------------------------------------------------------- payable dan rekening

// SetPayableTreaty meniru `Activity/SetPayableTreaty_Act.xml` langkah 1 (dipanggil tanpa `Posisi` dari Section
// AdjustmentDetail dan SetNameCurrency_Act; langkah 3 hanya untuk Posisi Acc/AccEdit yang tidak pernah dikirim).
// Payable 1 = ceding, 2 = SOB (bawaan bila kosong), 3 = penerima klaim. Rekening pertama klien + mata uang
// (GetDataBankAccount_sql) menimpa rekening penerima dan baris (langkah 1.5-1.6).
func SetPayableTreaty(k *Konteks, h *Halaman, idx int) error {
	b, err := adj(h, idx)
	if err != nil {
		return err
	}
	if h.Ambil(CD+"Payable") == "" { // 1.1
		h.Setel(CD+"Payable", "2")
		b["Payable"] = "2"
	}
	p := h.Ambil(CD + "Payable")
	b["Payable"] = p // 1.2
	klien := ""
	switch p {
	case "2": // 1.3
		h.Setel(CD+"PayableTo", h.Ambil(TM+"LeadingReinsSource"))
		b["PayableTo"] = h.Ambil(TM + "LeadingReinsSource")
		klien = h.Ambil(TM + "LeadingReinsSourceID")
	case "1": // 1.4
		h.Setel(CD+"PayableTo", h.Ambil(TM+"Ceding"))
		b["PayableTo"] = h.Ambil(TM + "Ceding")
		klien = h.Ambil(TM + "CedingID")
	case "3": // 1.7.1
		b["PayableTo"] = h.Ambil(CD + "PayableTo")
	}
	var r RekeningBank // 1.5-1.6
	if klien != "" {
		rek, err := k.Acuan.RekeningBank(k.Ctxt(), klien, b["CurrencyID"])
		if err != nil {
			return err
		}
		if len(rek) > 0 {
			r = rek[0]
		}
	}
	h.Setel(CD+"ReceiverClaim.NoAccount", r.AccountNo)
	h.Setel(CD+"ReceiverClaim.BranchOfBank", r.BranchOfBank)
	h.Setel(CD+"ReceiverClaim.NameOfBank", r.NameOfBank)
	h.Setel(CD+"ReceiverClaim.SwiftCode", r.SwiftCode)
	h.Setel(CD+"ReceiverClaim.IDOfBank", r.IDOfBank)
	b["NameOfBank"], b["NoAccount"], b["BranchOfBank"] = r.NameOfBank, r.AccountNo, r.BranchOfBank
	b["SwiftCode"], b["IDOfBank"] = r.SwiftCode, r.IDOfBank
	return nil
}

// PilihanRekening - isi autocomplete "Name of Bank" (pageList `Result.pxResults`): Payable 1/2 = rekening klien +
// mata uang (GetDataBankAccount_sql, langkah 1.5), Payable 3 = rekening mata uang (GetDataBankAccount2_sql, 1.7.2).
func PilihanRekening(k *Konteks, h *Halaman, idx int) ([]RekeningBank, error) {
	b, err := adj(h, idx)
	if err != nil {
		return nil, err
	}
	switch h.Ambil(CD + "Payable") {
	case "3":
		return k.Acuan.RekeningBankMataUang(k.Ctxt(), b["CurrencyID"])
	case "1":
		return k.Acuan.RekeningBank(k.Ctxt(), h.Ambil(TM+"CedingID"), b["CurrencyID"])
	default:
		return k.Acuan.RekeningBank(k.Ctxt(), h.Ambil(TM+"LeadingReinsSourceID"), b["CurrencyID"])
	}
}

// PilihRekening - "set saat pilih" autocomplete Name of Bank (CLIENTNAME->.PayableTo, BRANCHOFBANK, ACCOUNTNO,
// SWIFTCODE, IDOFBANK) lalu `Activity/SetPayableTo_act.xml` (langkah 5-6: Payable 3 menyalin ke penerima klaim).
func PilihRekening(h *Halaman, idx int, r RekeningBank) error {
	b, err := adj(h, idx)
	if err != nil {
		return err
	}
	b["NameOfBank"], b["PayableTo"], b["BranchOfBank"] = r.NameOfBank, r.ClientName, r.BranchOfBank
	b["NoAccount"], b["SwiftCode"], b["IDOfBank"] = r.AccountNo, r.SwiftCode, r.IDOfBank
	h.BersihkanPesan()
	if b["Payable"] == "3" {
		h.Setel(CD+"ReceiverClaim.NameOfBank", b["NameOfBank"])
		h.Setel(CD+"ReceiverClaim.BranchOfBank", b["BranchOfBank"])
		h.Setel(CD+"ReceiverClaim.NoAccount", b["NoAccount"])
		h.Setel(CD+"ReceiverClaim.SwiftCode", b["SwiftCode"])
	}
	return nil
}

// SetDLACedingSOB meniru `Activity/SetDLACedingSOB.xml`: nomor DLA ceding / SOB baris disalin ke klaim bila terisi.
func SetDLACedingSOB(h *Halaman, idx int) error {
	b, err := adj(h, idx)
	if err != nil {
		return err
	}
	if b["DLANoCeding"] != "" {
		h.Setel(CD+"DLANoCeding", b["DLANoCeding"])
	}
	if b["DLANoSOB"] != "" {
		h.Setel(CD+"DLANoSOB", b["DLANoSOB"])
	}
	return nil
}

// FilterLossAllocation meniru `Activity/FilterLossAllocation_Act.xml` (change Allocation): loss allocation adjustment
// tinggal baris ber-TreatyName dan Share sama dengan pilihan.
func FilterLossAllocation(h *Halaman, idx int) error {
	b, err := adj(h, idx)
	if err != nil {
		return err
	}
	var la []Baris
	for _, x := range h.AmbilDaftar(JalurAdj(idx, AnakLossAllocation)) {
		if x["TreatyName"] == b["TreatyName"] && x["SharePercentage"] == b["ShareLossAllocation"] {
			la = append(la, x)
		}
	}
	h.SetelDaftar(JalurAdj(idx, AnakLossAllocation), la)
	return nil
}

// PilihAllocation - "set saat pilih" autocomplete Allocation: `.SharePercentage -> .ShareLossAllocation`.
func PilihAllocation(h *Halaman, idx int, treatyName string) error {
	b, err := adj(h, idx)
	if err != nil {
		return err
	}
	b["TreatyName"] = treatyName
	for _, x := range h.AmbilDaftar(JalurAdj(idx, AnakLossAllocation)) {
		if x["TreatyName"] == treatyName {
			b["ShareLossAllocation"] = x["SharePercentage"]
		}
	}
	return FilterLossAllocation(h, idx)
}

// ---------------------------------------------------------------- komite

// Kepala pop-up "Komite klaim Treaty" (ProteksiInitialandDate_Act 2): halaman sementara, tanpa kolom - ikut `ModeLayar`.
const (
	JalurTanggalKomite = "TempCommiteClaim.DateOfComitee"
	JalurPICKomite     = "TreatyExchangeYearly.UserName"
)

// ProteksiInitialandDate meniru `Activity/ProteksiInitialandDate_Act.xml` (pembuka pop-up "Komite klaim Treaty"):
// salin catatan komite baris sebelumnya, tanggal dan PIC pop-up, riwayat, lalu pesan bila tipe / bank / payable belum
// dipilih.
func ProteksiInitialandDate(k *Konteks, h *Halaman, idx int) error {
	b, err := adj(h, idx)
	if err != nil {
		return err
	}
	if idx > 1 { // 1
		prev, _ := adj(h, idx-1)
		b["DataCommitteeTreaty.Remarks"] = prev["DataCommitteeTreaty.Remarks"]
		b["DataCommitteeTreaty.CircumCauseOfLoss"] = prev["DataCommitteeTreaty.CircumCauseOfLoss"]
	}
	nama, err := k.Acuan.NamaPelaku(k.Ctxt(), k.Pelaku) // 2 UserName = OperatorID.pyLabel
	if err != nil {
		return err
	}
	h.Setel(JalurPICKomite, nama)
	h.Setel(JalurTanggalKomite, k.Hari())               // 2 DateOfComitee = @substring(CurrentDateTime, 0, 8)
	k.Riwayat(h, TeksKomiteAkseptasi)                   // 2-3
	for i, x := range h.AmbilDaftar(DaftarAdjustment) { // 4
		if x["Type"] == "" {
			h.Setel("IsError", "2")
		}
		if x["AcceptanceStatus"] == "" && (x["NoAccount"] == "" || x["NameOfBank"] == "") {
			h.TambahPesan(JalurAdj(i+1, "BranchOfBank"), PesanBankNull)
		}
	}
	if v, _ := strconv.Atoi(h.Ambil("IsError")); v > 1 { // 5
		h.TambahPesan(JalurAdj(idx, "PayableTo"), PesanPilihTipeBayar)
	}
	if h.Ambil(CD+"Payable") == "" { // 7
		h.TambahPesan(JalurAdj(idx, "Type"), PesanPilihPayable)
	}
	return nil
}

// Kategori lampiran wajib (AttachmentProtect_ACT langkah 3).
const (
	LampiranLOD     = "LOD"
	LampiranDLA     = "DLA"
	LampiranSPGR    = "SPGR"
	LampiranADU     = "ADU"
	LampiranInvoice = "Invoice"
	LampiranSalvage = "Salvage"
)

// Pesan lampiran VERBATIM (AttachmentProtect_ACT langkah 2).
var pesanLampiran = map[string]string{
	LampiranLOD: "Please Upload Attachment LOD", LampiranDLA: "Please Upload Attachment DLA",
	LampiranSPGR: "Please Upload Attachment SPGR", LampiranADU: "Please Upload Attachment ADU",
	LampiranInvoice: "Please Upload Attachment Invoice", LampiranSalvage: "Please Upload Attachment Salvage",
}

// Pesan proteksi komite VERBATIM.
const (
	PesanBankKosongUlang = "Any Bank Data is Null, Please Check Again / Re-Select Bank Account"
	PesanCekSpreading    = "Please check your spreading !!!"
	PesanOccupation      = "Occupation cannot empty"
)

// AttachmentProtect meniru `Activity/AttachmentProtect_ACT.xml` (prasyarat tombol "Send to Committe": `Protect.CARI1
// = "1" And Protect.CARI2 = "1"`). `lampiran` = cacah berkas per kategori (`AttachCategory.pxResults.CountAttach` -
// sumbernya tidak terekspor, OQ-CP-12). `limitDirut` = LIMIT_BOTTOM baris Direktur Utama. Mengembalikan
// (lolosLampiran, lolosPremi) sebelum CekPremiLunas; lolosPremi dihitung pemanggil lewat `PremiLunas`.
func AttachmentProtect(h *Halaman, idx int, lampiran map[string]int, limitDirut string) (bool, error) {
	b, err := adj(h, idx)
	if err != nil {
		return false, err
	}
	h.BersihkanPesan() // 1
	ada := func(kat string) bool { return lampiran[kat] > 0 }
	var kal Kalkulator
	total := apd.New(0, 0) // 4.1
	if b["AcceptanceStatus"] != "2" {
		total = kal.B(b, "ValueAdjustment")
	}
	tipe := ""
	if b["AcceptanceStatus"] == "" || b["IsSubjectivity"] == "true" { // 4.2
		tipe = b["Type"]
	}
	bank := !(b["NameOfBank"] == "" && b["NoAccount"] == "" && b["IDOfBank"] == "")                                         // 4.3
	spread := len(h.AmbilDaftar(JalurAdj(idx, AnakSpreadAdj))) > 0 && len(h.AmbilDaftar(JalurAdj(idx, AnakQuotaShare))) > 0 // 4.4
	lewat := false                                                                                                          // 5-6
	if limitDirut != "" {
		lewat = !Lebih(kal.Teks("LIMIT_BOTTOM", limitDirut), total)
	}
	if err := kal.Galat(); err != nil {
		return false, err
	}
	pesan := func(kat string) { h.TambahPesan("IsError", pesanLampiran[kat]) }
	if tipe == "1" && !ada(LampiranLOD) { // 7
		pesan(LampiranLOD)
	}
	if (tipe == "1" || tipe == "2" || tipe == "3" || tipe == "4") && !ada(LampiranDLA) { // 8
		pesan(LampiranDLA)
	}
	if tipe == "1" && !ada(LampiranSPGR) { // 9
		pesan(LampiranSPGR)
	}
	if lewat && tipe == "1" && !ada(LampiranADU) { // 10
		pesan(LampiranADU)
	}
	if (tipe == "2" || tipe == "4") && !ada(LampiranInvoice) { // 11
		pesan(LampiranInvoice)
	}
	if tipe == "3" && !ada(LampiranSalvage) { // 12
		pesan(LampiranSalvage)
	}
	if !bank { // 13
		h.TambahPesan("IsError", PesanBankKosongUlang)
	}
	okupasiKosong := h.Ambil(CD+"Occupation") == ""
	if okupasiKosong { // 14
		h.TambahPesan("IsError", PesanOccupation)
	}
	if !spread { // 15
		h.TambahPesan("IsError", PesanCekSpreading)
	}
	lolos := false
	switch tipe { // 16-18 (syarat ADU langkah 16 selalu benar: `(L&&A) || (!L&&!A) || A`... ditiru apa adanya)
	case "1":
		lolos = ada(LampiranLOD) && ada(LampiranDLA) && ada(LampiranSPGR) &&
			((lewat && ada(LampiranADU)) || (!lewat && !ada(LampiranADU)) || ada(LampiranADU))
	case "2", "4":
		lolos = ada(LampiranInvoice) && ada(LampiranDLA)
	case "3":
		lolos = ada(LampiranSalvage) && ada(LampiranDLA)
	}
	if lolos { // 19-20: CARI1 yang sudah 1 ditimpa penanda bank lalu penanda spreading
		lolos = bank
	}
	if lolos {
		lolos = spread
	}
	// `[keputusan work owner 09-10-2026, penyimpangan sadar]` "Occupation cannot empty itu masih ada protek, tapi popup
	// komite masih muncul!": di XML langkah 14 hanya pesan (CARI1 langkah 16-20 tidak membacanya, dan action set
	// tombol membuka harness CommitteeTreaty tanpa syarat). Di sini Occupation kosong menggagalkan proteksi: popup
	// tidak dibuka dan Send Claim to Committee ditolak (aksiKomite), untuk setiap Type seperti pesannya.
	if okupasiKosong {
		lolos = false
	}
	return lolos, nil
}

// PremiLunas meniru `Activity/CekPremiLunas_Act.xml` (AttachmentProtect_ACT langkah 21, hanya bila `.IsKomite`
// kosong): saldo premi polis (tanpa titik) bermata uang adjustment > 0 = BELUM LUNAS, kecuali proteksi dibuka
// (OPENPROTEKSI_EDM). Belum lunas -> pesan.
func PremiLunas(k *Konteks, h *Halaman, idx int) (bool, error) {
	b, err := adj(h, idx)
	if err != nil {
		return false, err
	}
	if b["IsKomite"] != "" {
		return true, nil // 21.1 tidak dijalankan: ParamData.HASIL1 kosong -> 21.2 menyetel CARI2 = 1
	}
	nopol := h.Ambil(CD + "PolicyData.PolicyNo")
	saldo, err := k.Acuan.SaldoPremi(k.Ctxt(), strings.ReplaceAll(nopol, ".", ""), b["CurrencyID"])
	if err != nil {
		return false, err
	}
	saldo = strings.ReplaceAll(saldo, ",", ".") // 5
	d, err := AngkaTeks("HASIL1", saldo)
	if err != nil || d.Sign() <= 0 {
		return true, nil
	}
	buka, err := k.Acuan.AdaProteksiPremi(k.Ctxt(), nopol) // 7
	if err != nil {
		return false, err
	}
	if buka {
		return true, nil
	}
	h.TambahPesan("", PesanPremiBelumLunas) // 8
	return false, nil
}

// ---------------------------------------------------------------- kasir

// MuatanKasir - badan `SendAcceptationToKasir` (`TAllPaymentData`, HitServiceToKasir_Act langkah 13.1.3-13.3,
// jalur CLMP). Disimpan di outbox (efek "kasir"), dikirim hanya di produksi.
type MuatanKasir struct {
	NoTrans, NoKlaim, LbuId, NoPolis, AcceptType, Kepada, AccountNo, TglAksep string
	Nett, Deductible, KaliDeduct                                              string
	StsSyariah, CompanyName, LjtdId, LdcId, StsAp, LkuId, LbgID               string
	TglBolehBayar, Email, UserInput                                           string
}

// Kode tetap muatan kasir (langkah 13.1.3; bukan data orang).
const (
	kasirCompany = "NUSARE"
	kasirLjtd    = "D0031"
	kasirLdc     = "100081"
	kasirLdcSyr  = "100115" // 13.1.4 IsPEGASyariah
)

var bukanAngka = regexp.MustCompile(`[^0-9]`)

// PanjangNoAksepCLMP - HitServiceToKasir_Act langkah 13.1: `@length(.AcceptedNo) = 23 || 24`; selainnya keluar.
func PanjangNoAksepCLMP(no string) bool { return len(no) == 23 || len(no) == 24 }

// tglAksepTeks = `@substring(d,6,8)+"-"+@substring(d,4,6)+"-"+@substring(d,0,4)` atas "yyyyMMdd".
func tglAksepTeks(ymd string) string {
	if len(ymd) < 8 {
		return ""
	}
	return ymd[6:8] + "-" + ymd[4:6] + "-" + ymd[0:4]
}

// TglBolehBayar = langkah 13.1.5: hari > 25 -> tanggal 01 dua bulan sesudahnya, selainnya hari yang sama bulan
// berikutnya; bulan 13 -> "1"; tahun naik bila bulan hasil 01 DAN bulan berjalan Desember.
func TglBolehBayar(ymd string, bulanKini int) string {
	if len(ymd) < 8 {
		return ""
	}
	hari, _ := strconv.Atoi(ymd[6:8])
	bulan, _ := strconv.Atoi(ymd[4:6])
	tahun, _ := strconv.Atoi(ymd[0:4])
	bulan++
	if hari > 25 {
		bulan++
	}
	if bulan >= 13 {
		bulan = 1
	}
	hs := strconv.Itoa(hari)
	if hari > 25 {
		hs = "01"
	}
	if len(hs) < 2 {
		hs = "0" + hs
	}
	bs := strconv.Itoa(bulan)
	if len(bs) < 2 {
		bs = "0" + bs
	}
	if bulan == 1 && bulanKini == 12 {
		tahun++
	}
	return hs + "-" + bs + "-" + strconv.Itoa(tahun)
}

// SusunMuatanKasir = HitServiceToKasir_Act langkah 7, 13.1.1-13.1.5 (jalur CLMP). `email` = GL.F_GET_EMAIL ceding.
func SusunMuatanKasir(k *Konteks, h *Halaman, b Baris, email string, syariah bool) MuatanKasir {
	b["NoAccount"] = bukanAngka.ReplaceAllString(b["NoAccount"], "") // 7
	ldc := kasirLdc
	if syariah {
		ldc = kasirLdcSyr
	}
	tgl := YMD(b["AcceptedDate"])
	user := b["pxCreateOperator"]
	if user == "" {
		user = k.Pelaku
	}
	return MuatanKasir{
		NoTrans: b["AcceptedNo"], NoKlaim: h.Ambil(CD + "NoClaim"), LbuId: h.Ambil(OQ + "BusinessOldId"),
		NoPolis: h.Ambil(CD + "PolicyData.PolicyNo"), AcceptType: b["Type"], Kepada: b["PayableTo"],
		AccountNo: b["NoAccount"], TglAksep: tglAksepTeks(tgl), Nett: b["AdjustmentValue"],
		Deductible: b["IndividualRiskRNM"], KaliDeduct: b["IndividualRiskPercentage"], StsSyariah: "0",
		CompanyName: kasirCompany, LjtdId: kasirLjtd, LdcId: ldc, StsAp: "0", LkuId: b["CurrencyID"],
		LbgID: b["IDOfBank"], TglBolehBayar: TglBolehBayar(tgl, int(k.Sekarang.In(Jakarta).Month())),
		Email: email, UserInput: user,
	}
}

// BolehAkseptasi - tombol "Acceptation" AdjustmentDetail: tampil `.AcceptanceStatus = 1`, nonaktif
// `.DirectToKasir = 'false' || .StatusKasir != ”`; HitServiceToKasir_Act langkah 2 menuntut DirectToKasir "true".
func BolehAkseptasi(b Baris) bool {
	return b["AcceptanceStatus"] == "1" && b["DirectToKasir"] == "true" && b["StatusKasir"] == ""
}

// TerapkanStatusKasir = GetStatusKasir_Act (defer load blok kasir): hanya baris diterima ber-DirectToKasir yang belum
// sukses; KET log "Success" -> teks sukses, selainnya KET apa adanya.
func TerapkanStatusKasir(b Baris, ket string, ada bool) {
	if !(b["AcceptanceStatus"] == "1" && b["DirectToKasir"] == "true") || b["StatusKasir"] == StatusKasirSukses || !ada {
		return
	}
	if ket == "Success" {
		b["StatusKasir"] = StatusKasirSukses
		return
	}
	b["StatusKasir"] = ket
}

// ---------------------------------------------------------------- penutupan

// PeriksaTutupKlaim = CloseClaimProp langkah 1-4: ditolak bila ada adjustment berstatus kosong / 0, atau diterima
// ber-DirectToKasir yang belum masuk kasir (pesan terakhir yang menang).
func PeriksaTutupKlaim(h *Halaman) bool {
	pesan := ""
	for _, b := range h.AmbilDaftar(DaftarAdjustment) {
		if b["AcceptanceStatus"] == "" || b["AcceptanceStatus"] == "0" {
			pesan = PesanKomiteMasihJalan
		}
		if b["AcceptanceStatus"] == "1" && b["DirectToKasir"] == "true" && b["StatusKasir"] != StatusKasirSukses {
			pesan = PesanKasirBelumSukses
		}
	}
	if pesan != "" {
		h.TambahPesan("", pesan)
		return false
	}
	return true
}

// SelesaiTutupKlaim = CloseClaimProp langkah 5 dan 8 (`DataChronology.CARI12` - lihat PARITAS: InsertChronology_DT
// membaca CARI1, teks riwayat Pega adalah teks aksi sebelumnya; di sini teks langkah 5 `[dugaan]`).
func SelesaiTutupKlaim(k *Konteks, h *Halaman, remarks string) {
	h.Setel(CD+"Remark", remarks)
	h.Setel(CD+"Remark_Close", remarks)
	k.Riwayat(h, TeksTutupKlaim)
}

// BarisOSTutup = CloseClaimProp langkah 6 (PEGA_JSON_OS_AKSEP_KLAIMTNP, STATUS_RJT "4"): kolom datar, DATA_JSON
// kosong (keputusan work owner).
func BarisOSTutup(k *Konteks, h *Halaman, kunciKasus string) BarisOS {
	return BarisOS{CaseID: kunciKasus, NoClaim: h.Ambil(CD + "NoClaim"), NoPolis: h.Ambil(CD + "PolicyData.PolicyNo"),
		MasterID: h.Ambil(TM + "ID"), StsReject: StsOSTutupBerkas, CauseOfLossID: h.Ambil(CD + "CauseOfLossID"),
		CauseOfLoss: h.Ambil(CD + "CauseOfLoss"), EstimationDate: k.Sekarang, InsertOp: k.Pelaku,
		// 6.1-6.2: DATA_JSON = halaman InputParamOs (GetPageJSONString).
		DataJSON: JSONHalamanPega(map[string]string{
			"CauseOfLoss":    h.Ambil(CD + "CauseOfLoss"),
			"CauseOfLossID":  h.Ambil(CD + "CauseOfLossID"),
			"NoClaim":        h.Ambil(CD + "NoClaim"),
			"IDMasterTreaty": h.Ambil(TM + "ID"),
			"pxObjClass":     KelasOSAkseptasi,
		})}
}

// ---------------------------------------------------------------- PLA dan DLA

// RakitNomorPLA = TryMakePLA_Act langkah 12 / PrintDLATreatyIn langkah 11: `jenis + BusinessOldId + "." + MM.YYYY +
// "." + urut` (tanpa "T").
func RakitNomorPLA(jenis, oldID, mmYYYY string, urut int) string {
	return jenis + oldID + "." + mmYYYY + "." + lpad5(urut)
}

func lpad5(n int) string {
	s := strconv.Itoa(n)
	for len(s) < 5 {
		s = "0" + s
	}
	return s
}

// TerapkanNomorPLA = TryMakePLA_Act langkah 14 dan 18.1: NoPla kosong diisi; estimasi terkirim tanpa NoPLA mendapat
// nomornya. Dokumen PDF (langkah 23-30) = OQ-CP-05.
func TerapkanNomorPLA(h *Halaman, nomor string) {
	if h.Ambil(CD+"NoPla") == "" && nomor != "" {
		h.Setel(CD+"NoPla", nomor)
	}
	for _, e := range h.AmbilDaftar(DaftarEstimasi) {
		if e["PrintFaceClaim"] == "1" && e["NoPLA"] == "" {
			e["NoPLA"] = h.Ambil(CD + "NoPla")
		}
	}
}

// BolehDLA - tombol "Generate DLA": tampil `.AcceptanceStatus = 1`, nonaktif `.IsFacRetro != 1 || .DLA_No != ”`.
func BolehDLA(b Baris) bool {
	return b["AcceptanceStatus"] == "1" && b["IsFacRetro"] == "1" && b["DLA_No"] == ""
}
