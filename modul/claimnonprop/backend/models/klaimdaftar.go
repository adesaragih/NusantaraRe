package models

// Untuk apa berkas ini: PORT ACTIVITY DAFTAR KLAIM Non Prop - Insured Interests (AddInterestListCNP_Act,
// CountTotalInterest_Act, SetTPLNote_Act), Claim Amount per mata uang (AddListClaimNP_Act, SetCurrency_Act),
// deductible (SetFormat_Act, SetDataDeductible_act), Loss Allocation (AddLossAlocation_Act), dan hapus baris grid
// (`deleteRow` bawaan Pega).

import (
	"github.com/cockroachdb/apd/v3"
)

// Teks riwayat VERBATIM.
const (
	TeksTambahKlaim     = "Add Insured Insterest" // AddListClaimNP_Act 2.2 (salah ketik ikut)
	TeksTambahAlokasi   = "Add Claim Amount"      // AddListClaimNP_Act 3.2
	TeksGantiMataUang   = "CHANGE CURRENCY"       // SetCurrency_Act 3
	TeksTambahLossAlloc = "ADD LOSS ALLOCATION"   // AddLossAlocation_Act 3
	TeksGantiDeductible = "CHANGE DEDUCTIBLE"     // SetFormat_Act 2
)

// Tabel AddListClaimNP_Act (`Param.Table`).
const (
	TabelTambahKlaim     = "Claim"
	TabelTambahAlokasi   = "Alokasi"
	TabelTambahSpreading = "Spreading"
)

// barisDaftar - baris ke-n (berbasis satu) daftar `jalur`.
func barisDaftar(h *Halaman, jalur string, n int) (Baris, error) {
	d := h.AmbilDaftar(jalur)
	if n < 1 || n > len(d) {
		return nil, ErrBarisTidakAda
	}
	return d[n-1], nil
}

// AddInterest = `AddInterestListCNP_Act`: baris InterestList baru.
func AddInterest(h *Halaman) {
	h.TambahBaris(DaftarInterest, Baris{})
}

// AddListClaim = `AddListClaimNP_Act` (Param.Table). IsSaveToOs = 0 (1). Tabel tanpa nilai yang dikenal (tombol Add grid
// Claim Acceptation akseptasi memanggilnya tanpa parameter) hanya menurunkan IsSaveToOs.
// ⚠️ Langkah 7 InsertChronology_DT jalan pula tanpa Table dengan `DataChronology.CARI1` sisa aksi sebelumnya (halaman
// requestor); di sini riwayat hanya ditulis bila teksnya ditetapkan activity ini (PARITAS).
func AddListClaim(k *Konteks, h *Halaman, tabel string) {
	h.Setel("IsSaveToOs", "0") // 1
	jejak := func(b Baris) Baris {
		b["pxCreateDateTime"], b["pxCreateOperator"], b["pxCreateOpName"] = k.Waktu(), k.Pelaku, k.Pelaku
		return b
	}
	switch tabel {
	case TabelTambahKlaim: // 2
		h.TambahBaris(DaftarClaimAmount, jejak(Baris{}))
		k.Riwayat(h, TeksTambahKlaim)
	case TabelTambahAlokasi: // 3
		h.TambahBaris(DaftarXOL, jejak(Baris{}))
		k.Riwayat(h, TeksTambahAlokasi)
	case TabelTambahSpreading: // 5 - satu baris per mata uang Summary XOL
		for _, s := range h.AmbilDaftar(DaftarSummaryXOL) {
			h.TambahBaris(DaftarSpreading, jejak(Baris{"Currency": s["Currency"]}))
		}
	}
}

// SetCurrency = `SetCurrency_Act` (params IDCurr, idx, Note): nama mata uang (GetCurrency), kurs master
// (`TreatyInMaster.CurrencyList.Conversion`).
//   - Note "Interest" (7): KursObjectItem = 1 untuk IDR, selain itu kurs master.
//   - Note "Claim" (8): AltValue = kurs master; Value = Σ TSI Insured Interests mata uang itu, TPL = Σ TPLAmount, kurs =
//     KursObjectItem interest (8.4-8.5); deductible baris (8.5); lalu (9) mata uang SETIAP baris Loss Allocation =
//     mata uang baris ListClaimAmount BERINDEKS SAMA (ditiru apa adanya).
func SetCurrency(k *Konteks, h *Halaman, m MasterTreaty, idx int, catatan string) error {
	jalur := DaftarClaimAmount
	if catatan == "Interest" {
		jalur = DaftarInterest
	}
	r, err := barisDaftar(h, jalur, idx)
	if err != nil {
		return err
	}
	currID := r["CurrencyID"] // 1
	nama, err := k.Acuan.NamaMataUang(k.Ctxt(), currID)
	if err != nil {
		return err
	}
	k.Riwayat(h, TeksGantiMataUang) // 3-4
	var kal Kalkulator
	kursROE := ""
	for _, c := range m.CurrencyList { // 6
		if c.Currency == nama {
			kursROE = c.Conversion
		}
	}
	switch catatan {
	case "Interest": // 7
		r["KursObjectItem"] = kursROE
		if nama == "IDR" || currID == IDMataUangIDR {
			r["KursObjectItem"] = "1"
		}
		r["Currency"] = nama
	case "Claim": // 8
		r["AltValue"] = kursROE
		r["Currency"] = nama
		if currID == IDMataUangIDR { // 8.3
			r["AltValue"] = "1"
		}
		dedValue := kal.H(h, CD+"DeductibleValue")
		kursD := kal.Teks(".KURSROE", kursROE)
		switch {
		case h.Ambil(CD+"CurrencyDeductible") == nama:
		case h.Ambil(CD+"CurrencyDeductible") == IDMataUangIDR:
			if Lebih(kursD, apd.New(0, 0)) {
				dedValue = kal.BagiPega(dedValue, kursD)
			} else {
				dedValue = apd.New(0, 0)
			}
		default:
			dedValue = kal.Kali(dedValue, kursD)
		}
		totalTSI, totalTPL := apd.New(0, 0), apd.New(0, 0)
		isTPL := false
		kursAkhir := kursROE                               // 8.5 `AltValue = Local.KURSROE` (menimpa "1" IDR langkah 8.3 bila tak ada interest)
		for _, it := range h.AmbilDaftar(DaftarInterest) { // 8.4
			if it["Currency"] != nama {
				continue
			}
			totalTSI = kal.Tambah(totalTSI, kal.B(it, "TSIPerObject"))
			totalTPL = kal.Tambah(totalTPL, kal.B(it, "TPLAmount"))
			isTPL = isTPL || it["IsTPL"] == "true"
			kursAkhir = it["KursObjectItem"]
			if it["IsTPL"] == "true" { // 8.4.2
				h.Setel(CD+"TPLFormat", it["TPLFormat"])
				h.Setel(CD+"TPLType", it["TPLType"])
				h.Setel(CD+"PctTPL", it["TPLPct"])
			}
		}
		pctValue := kal.BagiPega(kal.Kali(totalTSI, kal.H(h, CD+"Amount")), k100()) // 8.5
		r["Value"] = Teks(totalTSI)
		r["CNPTSI"] = Teks(totalTSI)
		r["TPL"] = Teks(totalTPL)
		r["ClaimAmountCedant"] = Teks(kal.Tambah(totalTSI, totalTPL))
		r["AltValue"] = kursAkhir
		if isTPL {
			h.Setel(CD+"IsTPL", "true")
		} else {
			h.Setel(CD+"IsTPL", "false")
		}
		var ded *apd.Decimal
		minMax := h.Ambil(CD + "CNPDeducMinMax")
		if h.Ambil(CD+"FormType") == "1" {
			if minMax == "1" {
				ded = maks(dedValue, pctValue)
			} else {
				ded = pctValue
			}
		} else {
			if minMax != "1" {
				ded = maks(dedValue, pctValue)
			} else {
				ded = dedValue
			}
		}
		r["CNPDeductible"] = Teks(ded)
		klaim := h.AmbilDaftar(DaftarClaimAmount) // 9
		for i, l := range h.AmbilDaftar(DaftarLossAlloc) {
			if i < len(klaim) {
				l["Currency"], l["CurrencyID"] = klaim[i]["Currency"], klaim[i]["CurrencyID"]
			} else {
				l["Currency"], l["CurrencyID"] = "", ""
			}
		}
	}
	return kal.Galat()
}

// AddLossAllocation = `AddLossAlocation_Act`: satu baris Loss Allocation per mata uang (unik) ListClaimAmount.
// `TempListCurr` di Pega hidup di requestor dan tidak pernah dibuang; di sini dihitung dari ListClaimAmount saat ini.
func AddLossAllocation(k *Konteks, h *Halaman) {
	sudah := map[string]bool{}
	for _, r := range h.AmbilDaftar(DaftarClaimAmount) {
		if sudah[r["Currency"]] {
			continue
		}
		sudah[r["Currency"]] = true
		h.TambahBaris(DaftarLossAlloc, Baris{"Currency": r["Currency"], "CurrencyID": r["CurrencyID"]})
	}
	k.Riwayat(h, TeksTambahLossAlloc)
}

// HapusBarisDaftar - tombol Delete grid (`deleteRow` bawaan Pega).
func HapusBarisDaftar(h *Halaman, jalur string, idx int) error {
	if _, err := barisDaftar(h, jalur, idx); err != nil {
		return err
	}
	h.HapusBaris(jalur, idx)
	return nil
}

// SetFormat = `SetFormat_Act` (change centang Deductible): Deductible tidak dicentang -> Format, mata uang, nilai,
// persen, dan dasar deductible dikosongkan (1); riwayat (2-3).
func SetFormat(k *Konteks, h *Halaman) {
	if h.Ambil(CD+"DeductibleType") != "true" {
		for _, p := range []string{"FormType", "CurrencyDeductible", "DeductibleValue", "Amount", "TypeDeductible"} {
			h.Setel(CD+p, "")
		}
	}
	k.Riwayat(h, TeksGantiDeductible)
}

// SetDataDeductible = `SetDataDeductible_act` (params Posisi, TypeDeduct). Dari dropdown Format (tanpa parameter)
// TypeDeductible ikut dikosongkan (langkah 1 `= Param.TypeDeduct`, ditiru apa adanya).
func SetDataDeductible(h *Halaman, posisi, typeDeduct string) error {
	var kal Kalkulator
	h.Setel(CD+"TypeDeductible", typeDeduct) // 1
	interest := h.AmbilDaftar(DaftarInterest)
	if h.Ambil(CD+"CurrencyDeductible") == "" && len(interest) > 0 { // 2
		h.Setel(CD+"CurrencyDeductible", interest[0]["CurrencyID"])
		h.Setel(CD+"TSIDeductible", Teks(kal.BagiPega(kal.H(h, CD+"TotalSumInsuredIDR"), kal.B(interest[0], "KursObjectItem"))))
	}
	kurs := apd.New(1, 0)
	for _, it := range interest { // 3
		if it["CurrencyID"] == h.Ambil(CD+"CurrencyDeductible") {
			kurs = kal.B(it, "KursObjectItem")
		}
	}
	if posisi == "TypeDeductible" && h.Ambil(CD+"FormType") == "1" { // 4
		if h.Ambil(CD+"CurrencyDeductible") == IDMataUangIDR {
			h.Setel(CD+"TSIDeductible", h.Ambil(CD+"TotalSumInsuredIDR"))
		} else {
			h.Setel(CD+"TSIDeductible", Teks(kal.BagiPega(kal.H(h, CD+"TotalSumInsuredIDR"), kurs)))
		}
	}
	return kal.Galat()
}

// SetTPLNote = `SetTPLNote_Act` (baris InterestList `idx`, detail InputDtlInterest): TPL ditulis ke SETIAP baris
// ListClaimAmount (4.2-4.3, tanpa saringan mata uang - ditiru), Loss Allocation bermata uang sama dihitung ulang (5),
// lalu mesin XoL (8). Obj-Save langkah 4.5 / 7 = penyimpanan aksi.
func SetTPLNote(k *Konteks, h *Halaman, m MasterTreaty, idx int) error {
	it, err := barisDaftar(h, DaftarInterest, idx)
	if err != nil {
		return err
	}
	if it["TPLFormat"] == "1" { // 1
		it["TPLAmount"], it["TPLPct"], it["TPLCurrency"] = it["TPLAmount2"], it["TPLPct2"], it["CurrencyTPL2"]
		it["TPLType"], it["CNPMinMax"] = it["TPLType2"], it["CNPMinMax2"]
	}
	var kal Kalkulator
	pct, amount := kal.B(it, "TPLPct"), kal.B(it, "TPLAmount") // 2
	tipe := "TSI"
	if it["TPLType"] == "1" {
		tipe = "Claim"
	}
	it["TPLNote"] = Teks(pct) + "% of " + tipe + " min." + Teks(amount)
	cur := it["Currency"]
	pctAmount := apd.New(0, 0)
	if it["TPLType"] != "1" {
		pctAmount = kal.BagiPega(kal.Kali(pct, kal.B(it, "TSIPerObject")), k100())
	}
	klaim := apd.New(0, 0)
	for _, r := range h.AmbilDaftar(DaftarClaimAmount) { // 3
		if r["Currency"] == cur {
			klaim = kal.Tambah(klaim, kal.B(r, "Value"))
		}
	}
	lossClaim := apd.New(0, 0)
	for _, r := range h.AmbilDaftar(DaftarClaimAmount) { // 4
		if it["TPLType"] == "1" { // 4.1
			pctAmount = kal.BagiPega(kal.Kali(pct, klaim), k100())
		}
		var fin *apd.Decimal // 4.2
		if it["TPLFormat"] == "1" {
			if it["CNPMinMax"] == "1" {
				fin = maks(amount, pctAmount)
			} else {
				fin = pctAmount
			}
		} else {
			if it["CNPMinMax"] != "1" {
				fin = maks(amount, pctAmount)
			} else {
				fin = amount
			}
		}
		r["TPL"] = Teks(fin)
		bersih := kal.Kurang(kal.Tambah(kal.B(r, "Value"), fin), kal.B(r, "CNPDeductible"))
		r["ClaimAmountCedant"] = Teks(bersih)
		r["USD"] = Teks(kal.Kali(bersih, kal.B(r, "AltValue")))
		if r["Currency"] == cur { // 4.4
			lossClaim = kal.B(r, "ClaimAmountCedant")
		}
	}
	analis := h.Ambil(CD+"IsAnalisTransfer") == "1"
	for _, l := range h.AmbilDaftar(DaftarLossAlloc) { // 5
		if l["Currency"] != cur {
			continue
		}
		if l["CNPFlagOuts"] == "1" && analis {
			lossClaim = kal.Kurang(lossClaim, kal.B(l, "ClaimAmountAdjust"))
		}
		if l["CNPFlagOuts"] != "1" && !analis {
			l["ClaimAmountAdjust"] = Teks(kal.BagiPega(kal.Kali(lossClaim, kal.B(l, "ClaimPercentage")), k100()))
		}
	}
	if err := kal.Galat(); err != nil {
		return err
	}
	return HitungXOL(k, h, m, HitungXOLMod, 0, "") // 8
}
