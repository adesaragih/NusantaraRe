package models

// Untuk apa berkas ini: PORT ACTIVITY RINCIAN KERUGIAN - Insured Interest, deductible, claim amount, loss allocation,
// estimasi, spreading klaim (tiket 05-07, 09; Section OutstandingClaim_Intrs / _Est / _Sprd, InputAcceptation_Est).
//
// Setiap port menulis BASIS (baris baru, mata uang, kurs, nilai masukan) persis seperti activity aslinya; nilai
// TURUNAN dihitung `HitungTurunan` sekali jalan sesudahnya (spec §4). Riwayat (`InsertChronology_DT`) memakai teks
// VERBATIM korpus.

import "github.com/cockroachdb/apd/v3"

// ---------------------------------------------------------------- interest

// AddInterest meniru `Activity/AddInterest_act.xml` (tombol Add grid Insured Interests).
func AddInterest(k *Konteks, h *Halaman) {
	h.TambahBaris(DaftarInterest, Baris{"ObjectName": ""})
	k.Riwayat(h, "Add Insured Insterest")
}

// baris mengambil baris ke-n (berbasis satu) satu daftar.
func baris(h *Halaman, jalur string, n int) (Baris, error) {
	d := h.AmbilDaftar(jalur)
	if n < 1 || n > len(d) {
		return nil, ErrBarisTidakAda
	}
	return d[n-1], nil
}

// setMataUang menulis nama mata uang (RD BrowseCurrency_RD) dan kurs standar (RDB CurrencyStandard) satu baris.
func setMataUang(k *Konteks, b Baris, currID, propKurs string) error {
	nama, err := k.Acuan.NamaMataUang(k.Ctxt(), currID)
	if err != nil {
		return err
	}
	kurs, err := k.Acuan.KursStandar(k.Ctxt(), currID)
	if err != nil {
		return err
	}
	b["CurrencyID"] = currID
	b["Currency"] = nama
	b[propKurs] = kurs
	return nil
}

// SetCurencyInterest meniru `Activity/SetCurencyInterest_act.xml` (change Currency interest): mata uang dan kurs baris,
// lalu CountTotalInsterest_Act bila TSI terisi. `lama` = halaman tersimpan sebelum perubahan (padanan snapshot
// `TWorkPage` CopyOldataCurr_act).
func SetCurencyInterest(k *Konteks, lama, h *Halaman, idx int, currID string) error {
	b, err := baris(h, DaftarInterest, idx)
	if err != nil {
		return err
	}
	if err := setMataUang(k, b, currID, "KursObjectItem"); err != nil {
		return err
	}
	if b["TSIPerObject"] != "" {
		return CountTotalInsterest(lama, h, idx)
	}
	return nil
}

// CountTotalInsterest meniru `Activity/CountTotalInsterest_Act.xml` (change TSI / mata uang interest).
//
// Langkah 3-10 = turunan interest (`HitungTurunan`). Langkah 11-17 berjalan hanya bila `Local.PropCount==1`, yaitu bila
// baris interest TERAKHIR (perulangan 7 menimpa) ber-IsAdjVal kosong: rantai claim amount -> loss allocation ->
// estimasi -> spreading. Di sini rantai itu = (a) CopyOldataCurr_act langkah 11.4 untuk perubahan mata uang interest,
// (b) gross estimasi yang belum terkirim mengikuti hasil loss allocation (langkah 13.2.1), (c) seluruh turunan.
// ⚠️ Langkah 11.2.1 (`.Value = TSIPerObject` mentah, mengabaikan share ceding dan deductible) tidak ditiru: `.Value`
// claim amount adalah turunan rumus 2024 (AC 52) dari `.ClaimAmount`.
func CountTotalInsterest(lama, h *Halaman, idx int) error {
	d := h.AmbilDaftar(DaftarInterest)
	if err := HitungTurunan(h); err != nil {
		return err
	}
	if len(d) == 0 || d[len(d)-1]["IsAdjVal"] != "" {
		return nil
	}
	if len(h.AmbilDaftar(DaftarClaimAmount)) > 0 && lama != nil && idx > 0 {
		PropagasiMataUangInterest(lama, h, idx)
	}
	if err := HitungTurunan(h); err != nil {
		return err
	}
	if len(h.AmbilDaftar(DaftarLossAlloc)) > 0 {
		SelaraskanGrossEstimasi(h)
	}
	return HitungTurunan(h)
}

// PropagasiMataUangInterest meniru `Activity/CopyOldataCurr_act.xml` jalur "Interest" (FlagGantiCurrVal.CARI1):
// mata uang lama baris interest `idx` (snapshot `TWorkPage`) diganti mata uang baru di ListClaimAmount (langkah 6),
// SpreadingRisk (7), EstimationList belum terkirim (8), SpreadingClaim (9), SpreadingBreakQS (10) - beserta kursnya.
// Rantai syarat langkah 7-10 (kode 5 = OR, spec §4 tambahan ronde 6) terpenuhi untuk jalur Interest.
func PropagasiMataUangInterest(lama, h *Halaman, idx int) {
	bl, err := baris(lama, DaftarInterest, idx)
	if err != nil || len(lama.AmbilDaftar(DaftarInterest)) != len(h.AmbilDaftar(DaftarInterest)) {
		return // langkah 3 baris 2: panjang daftar lama = baru
	}
	bb, _ := baris(h, DaftarInterest, idx)
	lamaID := bl["CurrencyID"]
	if lamaID == "" || lamaID == bb["CurrencyID"] {
		return
	}
	ganti := func(jalur, propKurs string, saring func(Baris) bool) {
		for _, b := range h.AmbilDaftar(jalur) {
			if b["CurrencyID"] != lamaID || (saring != nil && !saring(b)) {
				continue
			}
			b["CurrencyID"] = bb["CurrencyID"]
			b["Currency"] = bb["Currency"]
			if propKurs != "" {
				b[propKurs] = bb["KursObjectItem"]
			}
		}
	}
	ganti(DaftarClaimAmount, "IDR", nil)
	ganti(DaftarLossAlloc, "PremiumSpreaded", nil)
	ganti(DaftarEstimasi, "KursValue", func(b Baris) bool { return b["PrintFaceClaim"] != "1" })
	ganti(DaftarSpreading, "", nil)
	ganti(DaftarBreakQS, "", nil)
}

// SelaraskanGrossEstimasi meniru `CountTotalInsterest_Act` langkah 13: estimasi belum terkirim (PrintFaceClaim != 1)
// yang mata uangnya dan TypeLossID-nya sama dengan baris loss allocation mengambil `.GrossEstimationPct =
// .ClaimSpreaded` baris itu (perulangan loss allocation - baris terakhir yang cocok menang).
func SelaraskanGrossEstimasi(h *Halaman) {
	for _, la := range h.AmbilDaftar(DaftarLossAlloc) {
		for _, e := range h.AmbilDaftar(DaftarEstimasi) {
			if e["CurrencyID"] == la["CurrencyID"] && e["TypeLossID"] == la["TreatyType"] && e["PrintFaceClaim"] != "1" {
				e["GrossEstimationPct"] = la["ClaimSpreaded"]
			}
		}
	}
}

// DeleteInterest meniru `Activity/DeleteInterest_act.xml` (tombol Delete baris interest). Pengurangan total
// (langkah 2-3, 5) = turunan.
func DeleteInterest(k *Konteks, h *Halaman, idx int) error {
	if _, err := baris(h, DaftarInterest, idx); err != nil {
		return err
	}
	h.HapusBaris(DaftarInterest, idx)
	k.Riwayat(h, "Delete  Insured Interest")
	return nil
}

// ---------------------------------------------------------------- deductible dan claim amount

// SetFormat meniru `Activity/SetFormat_Act.xml` (change Deductible): DeductibleType "false" mengosongkan isian
// deductible.
func SetFormat(k *Konteks, h *Halaman) {
	if h.Ambil(CD+"DeductibleType") == "false" {
		for _, p := range []string{"FormType", "CurrencyDeductible", "DeductibleValue", "Amount", "TypeDeductible", "NetDeductibleValue"} {
			h.Setel(CD+p, "")
		}
	}
	k.Riwayat(h, "CHANGE DEDUCTIBLE")
}

// AddListClaimAmount meniru `Activity/AddListClaimAmount.xml` (tombol Add grid claim amount): ShareCeding 100 bila 0
// atau kosong; satu baris per mata uang interest (kemunculan pertama, kurs baris interest pertama mata uang itu) dengan
// ClaimAmount = TotalInterestInsured mata uang itu. ⚠️ Daftar TIDAK dikosongkan lebih dulu - menekan Add dua kali
// menggandakan baris (ditiru apa adanya).
func AddListClaimAmount(k *Konteks, h *Halaman) error {
	if s := h.Ambil(CD + "ShareCeding"); s == "0" || s == "" { // 1
		h.Setel(CD+"ShareCeding", "100")
	}
	if err := HitungTurunan(h); err != nil { // TotalInterestInsured terkini
		return err
	}
	kurs := map[string]string{}
	for _, b := range h.AmbilDaftar(DaftarInterest) {
		if _, ada := kurs[b["CurrencyID"]]; !ada {
			kurs[b["CurrencyID"]] = b["KursObjectItem"]
		}
	}
	for _, m := range urutanMataUang(h.AmbilDaftar(DaftarInterest)) { // 3-6
		h.TambahBaris(DaftarClaimAmount, Baris{
			"CurrencyID":  m["CurrencyID"],
			"Currency":    m["Currency"],
			"IDR":         kurs[m["CurrencyID"]],
			"ClaimAmount": tsiMataUang(h, m["CurrencyID"]),
		})
	}
	k.Riwayat(h, "Add Claim Amount") // 8-9
	return nil
}

// SetCurencyList meniru `Activity/SetCurencyList_act.xml` (change Currency claim amount): nama, kurs standar; baris ke-2
// dan seterusnya mengambil kurs baris lain bermata uang sama (langkah 6.1.1, baris terakhir yang cocok menang).
func SetCurencyList(k *Konteks, h *Halaman, idx int, currID string) error {
	b, err := baris(h, DaftarClaimAmount, idx)
	if err != nil {
		return err
	}
	if err := setMataUang(k, b, currID, "IDR"); err != nil {
		return err
	}
	if idx > 1 {
		for _, x := range h.AmbilDaftar(DaftarClaimAmount) {
			if x["CurrencyID"] == b["CurrencyID"] {
				b["IDR"] = x["IDR"]
			}
		}
	}
	return nil
}

// CountListClaimAmountIDR meniru `Activity/CountListClaimAmountIDR.xml` (change Claim Amount). Langkah 4 berprakondisi
// nonaktif, sehingga teks riwayat SELALU "Add Value Claim Amount". Lalu AddLossAllocation_act bila IsOutstanding 1.
func CountListClaimAmountIDR(k *Konteks, h *Halaman) error {
	if err := HitungTurunan(h); err != nil {
		return err
	}
	k.Riwayat(h, "Add Value Claim Amount")
	if h.Ambil("IsOutstanding") == "1" {
		return AddLossAllocation(k, h)
	}
	return nil
}

// DeleteListClaim meniru `Activity/DeleteListClaim_Act.xml`.
func DeleteListClaim(k *Konteks, h *Halaman, idx int) error {
	if _, err := baris(h, DaftarClaimAmount, idx); err != nil {
		return err
	}
	h.HapusBaris(DaftarClaimAmount, idx)
	h.SetelDaftar(DaftarInterestDtl, nil)
	k.Riwayat(h, "Delete List Claim Amount")
	return nil
}

// ---------------------------------------------------------------- loss allocation

// AddLossAllocation meniru `Activity/AddLossAllocation_act.xml` (tombol Add grid Loss Allocation): treaty dari baris
// loss allocation TERAKHIR (langkah 5 menimpa); satu mata uang claim amount -> satu baris (mata uang dan kurs claim
// amount terakhir), lebih -> satu baris per mata uang (unik CurrencyID#Currency#IDR).
func AddLossAllocation(k *Konteks, h *Halaman) error {
	var uang []Baris
	sudah := map[string]bool{}
	for _, b := range h.AmbilDaftar(DaftarClaimAmount) {
		kunci := b["CurrencyID"] + "#" + b["Currency"] + "#" + b["IDR"]
		if !sudah[kunci] {
			sudah[kunci] = true
			uang = append(uang, b)
		}
	}
	var tn, tt, share string
	for _, b := range h.AmbilDaftar(DaftarLossAlloc) {
		tn, tt, share = b["TreatyName"], b["TreatyType"], b["SharePercentage"]
	}
	ca := h.AmbilDaftar(DaftarClaimAmount)
	switch {
	case len(uang) == 1: // 6
		last := ca[len(ca)-1]
		h.TambahBaris(DaftarLossAlloc, Baris{"CurrencyID": last["CurrencyID"], "Currency": last["Currency"],
			"PremiumSpreaded": last["IDR"], "TreatyType": tt, "TreatyName": tn, "SharePercentage": share})
	case len(uang) > 1: // 7
		for _, u := range uang {
			h.TambahBaris(DaftarLossAlloc, Baris{"CurrencyID": u["CurrencyID"], "Currency": u["Currency"],
				"PremiumSpreaded": u["IDR"], "TreatyType": tt, "TreatyName": tn, "SharePercentage": share})
		}
	}
	teks := "Add Loss Allocation"
	if h.Ambil("IsOutstanding") == "1" { // 9
		teks = "Edit % dan Value Loss Allocation"
	}
	k.Riwayat(h, teks)
	if h.Ambil("IsOutstanding") == "1" { // 11
		return CountPersen(k, h)
	}
	return HitungTurunan(h)
}

// SetCurrency meniru `Activity/SetCurrency_Act.xml` (change Curr loss allocation, InputAcceptation_Est).
func SetCurrency(k *Konteks, h *Halaman, idx int, currID string) error {
	b, err := baris(h, DaftarLossAlloc, idx)
	if err != nil {
		return err
	}
	if err := setMataUang(k, b, currID, "PremiumSpreaded"); err != nil {
		return err
	}
	if idx > 1 {
		for _, x := range h.AmbilDaftar(DaftarLossAlloc) {
			if x["CurrencyID"] == b["CurrencyID"] {
				b["PremiumSpreaded"] = x["PremiumSpreaded"]
			}
		}
	}
	return nil
}

// TipeJenisLossAllocation - `Type` REINSURANCETYPE grid Loss Allocation (parameter RD `Type = 4`, Section
// InputAcceptation_Adjs / OutstandingClaim_Est).
const TipeJenisLossAllocation = "4"

// SetNameTreaty meniru `Activity/SetNameTreaty_Act.xml` (change Treaty Type loss allocation). Langkah 7 (pesan duplikat
// treaty+mata uang) ber-remark.
//
// ⚠️ Di OutstandingClaim_Est / InputAcceptation_Est dropdown `.TreatyName` bernilai `TreatyInMaster.Limits.TreatyType`
// (NAMA treaty) dan activity dipanggil TANPA Param.Index. `[data DEV 07-10-2026]` baris loss allocation dokumen CLMP
// membawa `TreatyType` = ID REINSURANCETYPE ber-TYPE 4 yang NOTE-nya sama dengan `TreatyName` (10035 "QUOTA SHARE",
// 10042 "SURPLUS"). Maka: baris bernama tanpa TreatyType mendapat ID lewat NOTE (saringan B RD BrowseReinsuranceType_RD
// `.Note = Param.Note`); baris ber-TreatyType mendapat nama lewat ID (langkah 4.1).
func SetNameTreaty(k *Konteks, h *Halaman, idx int) error {
	d := h.AmbilDaftar(DaftarLossAlloc)
	for i, b := range d {
		if idx > 0 && i != idx-1 {
			continue
		}
		switch {
		case b["TreatyType"] != "" && idx > 0:
			nama, err := k.Acuan.NamaJenisReasuransi(k.Ctxt(), b["TreatyType"])
			if err != nil {
				return err
			}
			if nama != "" {
				b["TreatyName"] = nama
			}
		case b["TreatyName"] != "":
			id, err := k.Acuan.IDJenisReasuransi(k.Ctxt(), b["TreatyName"], TipeJenisLossAllocation)
			if err != nil {
				return err
			}
			b["TreatyType"] = id
		}
	}
	k.Riwayat(h, "Add Name Loss Allocation")
	return nil
}

// CountPersen meniru `Activity/CountPersen_act.xml` (change Share(%) loss allocation). Pesan "Total more than 100%"
// (langkah 4.2.4) ber-remark - hanya penanda. Blok 5 ber-remark (perkalian ganda share ceding tidak pernah berjalan,
// spec §4 tambahan 19-09). Riwayat langkah 8 berprakondisi nonaktif ("Add type ..."), langkah 9 menimpanya bila
// IsOutstanding 1. Lalu AddEstimation_Act bila IsOutstanding 1.
func CountPersen(k *Konteks, h *Halaman) error {
	if err := HitungTurunan(h); err != nil {
		return err
	}
	teks := "Add type and Change % Loss Allocation"
	if h.Ambil("IsOutstanding") == "1" {
		teks = "Edit type and Change % Loss Allocation"
	}
	k.Riwayat(h, teks)
	if h.Ambil("IsOutstanding") == "1" {
		return AddEstimation(k, h)
	}
	return nil
}

// RemoveLossAlloction meniru `Activity/RemoveLossAlloction_act.xml`.
func RemoveLossAlloction(k *Konteks, h *Halaman, idx int) error {
	if _, err := baris(h, DaftarLossAlloc, idx); err != nil {
		return err
	}
	h.HapusBaris(DaftarLossAlloc, idx)
	k.Riwayat(h, "Delete Loss Allocation")
	return nil
}

// ---------------------------------------------------------------- estimasi

// AddEstimation meniru `Activity/AddEstimation_Act.xml` (tombol Add grid Estimation List, CountPersen_act langkah 11).
//
// LossAlloc = baris loss allocation unik (CurrencyID#Currency#TreatyType#...): Java langkah 3 menghapus duplikat
// dengan kunci isi baris; di sini kunci CurrencyID#TreatyType#ClaimSpreaded#PremiumSpreaded (kemunculan pertama).
// Langkah 4: estimasi belum terkirim yang TypeLossID-nya = TreatyType LossAlloc dihapus, dan gross-nya DITAMBAHKAN ke
// SETIAP LossAlloc bermata uang sama (langkah 4.2.2.1.2 - tanpa memeriksa jenis treaty; temuan logika butir 8,
// ditiru). Langkah 5: satu estimasi per LossAlloc ber-ClaimSpreaded bukan nol (`.ClaimSpreaded==0` T=3).
func AddEstimation(k *Konteks, h *Halaman) error {
	if err := HitungTurunan(h); err != nil {
		return err
	}
	var la []Baris
	sudah := map[string]bool{}
	for _, b := range h.AmbilDaftar(DaftarLossAlloc) { // 2-3
		kunci := b["CurrencyID"] + "#" + b["Currency"] + "#" + b["TreatyType"] + "#" + b["TreatyName"] + "#" + b["ClaimSpreaded"] + "#" + b["PremiumSpreaded"]
		if sudah[kunci] {
			continue
		}
		sudah[kunci] = true
		la = append(la, Baris{"CurrencyID": b["CurrencyID"], "Currency": b["Currency"], "TreatyType": b["TreatyType"],
			"TreatyName": b["TreatyName"], "ClaimSpreaded": b["ClaimSpreaded"], "PremiumSpreaded": b["PremiumSpreaded"]})
	}
	var kal Kalkulator
	for _, l := range la { // 4
		est := h.AmbilDaftar(DaftarEstimasi)
		var sisa []Baris
		for _, e := range est {
			if e["PrintFaceClaim"] != "1" && e["TypeLossID"] == l["TreatyType"] {
				for _, x := range la {
					if x["CurrencyID"] == e["CurrencyID"] {
						x["ClaimSpreaded"] = Teks(kal.Tambah(kal.B(x, "ClaimSpreaded"), kal.B(e, "GrossEstimationPct")))
					}
				}
				continue // 4.2.3 dihapus
			}
			sisa = append(sisa, e)
		}
		h.SetelDaftar(DaftarEstimasi, sisa)
	}
	if err := kal.Galat(); err != nil {
		return err
	}
	for _, l := range la { // 5
		if l["ClaimSpreaded"] == "" || Nol(kal.B(l, "ClaimSpreaded")) {
			continue
		}
		kurs, err := k.Acuan.KursStandar(k.Ctxt(), l["CurrencyID"])
		if err != nil {
			return err
		}
		h.TambahBaris(DaftarEstimasi, Baris{"CurrencyID": l["CurrencyID"], "GrossEstimationPct": l["ClaimSpreaded"],
			"Currency": l["Currency"], "KursValue": kurs, "EstimationDate": k.Hari(),
			"TypeLoss": l["TreatyName"], "TypeLossID": l["TreatyType"]})
	}
	return CountEstimation(k, h) // 7
}

// CurencyEstimation meniru `Activity/CurencyEstimation_Act.xml` (change Currency estimasi): mata uang, kurs standar,
// lalu CountEstimation_Act. Langkah 6 (salin kurs baris lain) ber-remark.
func CurencyEstimation(k *Konteks, h *Halaman, idx int, currID string) error {
	b, err := baris(h, DaftarEstimasi, idx)
	if err != nil {
		return err
	}
	if err := setMataUang(k, b, currID, "KursValue"); err != nil {
		return err
	}
	return CountEstimation(k, h)
}

// Pesan estimasi.
const (
	PesanEstimasiLewatTSI    = "Value Estimation RNM more than TSI Interest Insured"
	PesanTanggalEstimasiDOL  = "Estimation Date should not be less than Date of Loss"
	PesanTanggalEstimasiKini = "Estimation Date should not be more than todays date"
)

// CountEstimation meniru `Activity/CountEstimation_Act.xml` (change Gross Estimate, AddEstimation_Act langkah 7).
// Nilai = turunan; penanda dan pesan menurut langkah 17-21; riwayat 23-26; CountSpreading_act bila IsOutstanding 1.
// Cash call langkah 13-16 hanya mengisi Local.CashCall yang tidak dipakai activity ini - tidak ditulis.
func CountEstimation(k *Konteks, h *Halaman) error {
	if err := HitungTurunan(h); err != nil {
		return err
	}
	var kal Kalkulator
	totEst := kal.H(h, CD+"TotalEstimasiIDR")
	tsi := kal.H(h, CD+"TotalSumInsuredIDR")
	if err := kal.Galat(); err != nil {
		return err
	}
	if Lebih(totEst, tsi) { // 17-18 (TotalInsured = jumlah TSIPerObjectIDR = TotalSumInsuredIDR)
		h.TambahPesan("", PesanEstimasiLewatTSI)
		h.Setel("IsCFS", "")
	} else { // 19
		h.Setel("IsCFS", "1")
	}
	adaBaru := false // 20
	for _, b := range h.AmbilDaftar(DaftarEstimasi) {
		if b["PrintFaceClaim"] != "1" {
			adaBaru = true
		}
	}
	if adaBaru && h.Ambil("IsOutstanding") == "1" { // 21
		h.Setel("ReCFS", "1")
	}
	teks := "Add Value Estimation" // 23 (prakondisi nonaktif)
	for _, b := range h.AmbilDaftar(DaftarEstimasi) {
		if b["PrintFaceClaim"] == "1" { // 24
			teks = "Edit  Value Estimation"
		}
	}
	if h.Ambil("IsOutstanding") == "1" { // 25
		teks = "Edit  Value Estimation"
	}
	k.Riwayat(h, teks)
	if h.Ambil("IsOutstanding") == "1" { // 27
		return CountSpreading(k, h)
	}
	return nil
}

// CheckEstimateDate meniru `Activity/CheckEstimateDate_Act.xml`: tanggal estimasi di antara Date of Loss dan hari ini
// (AC 35). `@defaultCompareDates(DOL, tgl)=="1"` = DOL lebih kemudian.
func CheckEstimateDate(k *Konteks, h *Halaman) {
	for i, b := range h.AmbilDaftar(DaftarEstimasi) {
		jalur := JalurAnak(DaftarEstimasi, i+1, "EstimationDate")
		if SesudahTanggal(h.Ambil(CD+"DateOfLoss"), b["EstimationDate"]) {
			h.TambahPesan(jalur, PesanTanggalEstimasiDOL)
		}
		if b["EstimationDate"] != "" && SesudahTanggal(b["EstimationDate"], k.Hari()) {
			h.TambahPesan(jalur, PesanTanggalEstimasiKini)
		}
	}
}

// CashCall - DeleteEstimation_Act langkah 6: untuk setiap baris loss allocation, `Limits` bernama sama, `Detail`
// ber-TreatyGroupID klaim, SEMUA `CashLossList` (langkah 6.2.1.1.1 berprakondisi nonaktif - saringan treaty group
// tidak dievaluasi); `.Value` baris terakhir yang dilalui menang. Mata uang cash call tidak dipakai pembanding.
func CashCall(h *Halaman, m MasterTreaty) string {
	nilai := ""
	for _, la := range h.AmbilDaftar(DaftarLossAlloc) {
		for _, l := range m.Limits {
			if l.TreatyType != la["TreatyName"] {
				continue
			}
			for _, d := range l.Detail {
				if d.TreatyGroupID != h.Ambil(CD+"TreatyGroupID") {
					continue
				}
				for _, c := range d.CashLossList {
					nilai = c.Value
				}
			}
		}
	}
	return nilai
}

// DeleteEstimation meniru `Activity/DeleteEstimation_Act.xml` (tombol Delete estimasi).
// ⚠️ Langkah 7 MENGOSONGKAN IsPLA saat estimasi melebihi cash call - kebalikan SaveOutstanding_Act langkah 36 (temuan
// logika butir 10; ditiru apa adanya).
func DeleteEstimation(k *Konteks, h *Halaman, idx int, m MasterTreaty) error {
	if _, err := baris(h, DaftarEstimasi, idx); err != nil {
		return err
	}
	h.HapusBaris(DaftarEstimasi, idx) // 1-2
	if err := HitungTurunan(h); err != nil {
		return err // 3-4
	}
	var kal Kalkulator
	tot := kal.H(h, CD+"TotalEstimasiIDR")
	cc := kal.Teks("CashCall", CashCall(h, m))
	tsi := kal.H(h, CD+"TotalSumInsuredIDR")
	if err := kal.Galat(); err != nil {
		return err
	}
	if Lebih(tot, cc) { // 7
		h.Setel(CD+"IsPLA", "")
	}
	// 5 dan 8 (IsCFS) selalu ditimpa langkah 14-15.
	if Lebih(tot, cc) || !Lebih(tot, tsi) { // 9 Page-Clear-Messages
		h.BersihkanPesan()
	}
	k.Riwayat(h, "Delete Estimation")    // 10
	if h.Ambil("IsOutstanding") == "1" { // 11
		if err := CountSpreading(k, h); err != nil {
			return err
		}
	}
	terkirim := false // 13
	for _, b := range h.AmbilDaftar(DaftarEstimasi) {
		if b["PrintFaceClaim"] == "1" {
			terkirim = true
		}
	}
	if !terkirim { // 14
		h.Setel("IsCFS", "1")
		h.Setel("IsOutstanding", "0")
	} else { // 15
		h.Setel("IsCFS", "")
		h.Setel("IsOutstanding", "1")
	}
	return nil
}

// ---------------------------------------------------------------- spreading klaim

// PesanShareLebih100 - CountSpreading_act langkah 6.2.4 (VERBATIM).
const PesanShareLebih100 = "Persen Share tidak boleh lebih besar dari 100"

// PesanShareKurang100 - ⚠️ PENYIMPANGAN SADAR AC 30 / spec §6 (keputusan work owner): total share reinsurer wajib
// TEPAT 100% - Pega hanya menjaga sisi > 100. Teks baru, sejajar teks Pega.
const PesanShareKurang100 = "Persen Share tidak boleh lebih kecil dari 100"

// CountSpreading meniru `Activity/CountSpreading_act` (berkas CountSpreading_Act.xml; change Share(%) spreading klaim,
// CountEstimation_Act langkah 27): pesan bila total share satu mata uang > 100, turunan, riwayat.
func CountSpreading(k *Konteks, h *Halaman) error {
	if err := HitungTurunan(h); err != nil {
		return err
	}
	tot, err := TotalSharePerMataUang(h)
	if err != nil {
		return err
	}
	seratus := k100()
	for i, b := range h.AmbilDaftar(DaftarSpreading) {
		if Lebih(tot[b["CurrencyID"]], seratus) {
			h.TambahPesan(JalurAnak(DaftarSpreading, i+1, "SharePercentage"), PesanShareLebih100)
		}
	}
	teks := "Add % spreading Claim" // 12 (prakondisi nonaktif)
	if h.Ambil("IsOutstanding") == "1" {
		teks = "Edit % spreading Claim"
	}
	k.Riwayat(h, teks)
	return nil
}

func k100() *apd.Decimal { return apd.New(100, 0) }

// PeriksaShareTepat100 - penjaga dua sisi AC 30: total SharePercentage SpreadingClaim setiap mata uang tepat 100.
// Mengembalikan pesan untuk tiap mata uang yang menyimpang (kosong = lolos).
func PeriksaShareTepat100(h *Halaman) ([]string, error) {
	tot, err := TotalSharePerMataUang(h)
	if err != nil {
		return nil, err
	}
	var pesan []string
	for _, m := range urutanMataUang(h.AmbilDaftar(DaftarSpreading)) {
		t := tot[m["CurrencyID"]]
		switch Banding(t, k100()) {
		case 1:
			pesan = append(pesan, PesanShareLebih100+" ("+m["Currency"]+")")
		case -1:
			pesan = append(pesan, PesanShareKurang100+" ("+m["Currency"]+")")
		}
	}
	return pesan, nil
}

// SetTreatyNameSpreading meniru `Activity/SetTreatyNameSpreading_Act.xml` (change Treaty Type spreading klaim):
// nama treaty baris `idx` dari REINSURANCETYPE, lalu SpreadingBreakQS DIBANGUN ULANG dari
// `TreatyInMaster.Limits(1).Detail(1).SpreadingList` untuk setiap mata uang estimasi (satu mata uang: mata uang
// estimasi baris 1). Langkah 10 (GetBreakDownTreaty) ber-remark.
// Langkah 16-17 (FacRetroList dari TREATYREINSURER) memakai FindData.CARI3/CARI6 dari DataF.CARI11/CARI13 yang TIDAK
// disetel activity ini - bacaan dengan tahun dan treaty group kosong tidak pernah menemukan baris; tidak ditulis.
func SetTreatyNameSpreading(k *Konteks, h *Halaman, idx int, m MasterTreaty) error {
	b, err := baris(h, DaftarSpreading, idx)
	if err != nil {
		return err
	}
	if b["TreatyType"] != "" {
		nama, err := k.Acuan.NamaJenisReasuransi(k.Ctxt(), b["TreatyType"])
		if err != nil {
			return err
		}
		if nama != "" {
			b["TreatyName"] = nama
		}
	}
	var daftar []SpreadingMaster
	if len(m.Limits) > 0 && len(m.Limits[0].Detail) > 0 {
		daftar = m.Limits[0].Detail[0].SpreadingList
	}
	mu := urutanMataUang(h.AmbilDaftar(DaftarEstimasi))
	var qs []Baris
	sudah := map[string]bool{}
	tambah := func(s SpreadingMaster, id, nama string) {
		kunci := s.ReinsTypeID + "#" + id
		if sudah[kunci] {
			return // langkah 14 Java: unik TreatyType#CurrencyID
		}
		sudah[kunci] = true
		pct := s.Pct
		if d, err := AngkaTeks("Pct", s.Pct); err == nil {
			pct = Teks(d)
		}
		qs = append(qs, Baris{"TreatyType": s.ReinsTypeID, "TreatyName": s.ReinsTypeName, "SharePercentage": pct,
			"CurrencyID": id, "Currency": nama})
	}
	switch {
	case len(mu) == 1: // 12
		e := h.AmbilDaftar(DaftarEstimasi)[0]
		for _, s := range daftar {
			tambah(s, e["CurrencyID"], e["Currency"])
		}
	case len(mu) > 1: // 13
		for _, u := range mu {
			for _, s := range daftar {
				tambah(s, u["CurrencyID"], u["Currency"])
			}
		}
	}
	h.SetelDaftar(DaftarBreakQS, qs) // 15
	return HitungTurunan(h)
}

// TambahBarisSpreading - ikon grid standar Pega (`pzPegaDefaultGridIcons`) grid "Spreading List" pertama: baris kosong
// baru, mata uang mengikuti estimasi baris 1. Tombol Add kustom grid ini memanggil `AddSpreading_Act` yang TIDAK ada di
// ekspor (nonaktif, OQ); satu-satunya jalan baris spreading lahir adalah penambah baris bawaan grid.
func TambahBarisSpreading(h *Halaman) {
	b := Baris{"TreatyType": "", "SharePercentage": ""}
	if est := h.AmbilDaftar(DaftarEstimasi); len(est) > 0 {
		b["CurrencyID"], b["Currency"] = est[0]["CurrencyID"], est[0]["Currency"]
	}
	h.TambahBaris(DaftarSpreading, b)
}

// HapusBarisSpreading - ikon hapus baris bawaan grid; baris IsOldData "Yes" ditolak (tombol Delete disabled jika
// `.IsOldData='Yes'`).
func HapusBarisSpreading(h *Halaman, idx int) error {
	b, err := baris(h, DaftarSpreading, idx)
	if err != nil {
		return err
	}
	if b["IsOldData"] == "Yes" {
		return ErrBarisBeku
	}
	h.HapusBaris(DaftarSpreading, idx)
	return HitungTurunan(h)
}
