package models

// Untuk apa berkas ini: PORT ACTIVITY BARIS ESTIMASI (grid `.EstimationList` section Estimasi / EstimasiPA /
// EstimasiMarine, kelas Data-ObjectItem / Data-Estimasi): tambah (`ValidateInputEstimate_act`), tanggal
// (`ProtectionDate_Act`), mata uang (`SetConvertValueKurs_Estimation`, `GetNameCurrency_Act`), nilai
// (`CheckEstimateValue`), hapus (`DeleteValueEstimation` + `CountSpreadingClaim_ACT`), deductible (`CountTSI_Act`).
//
// ⚠️ `Set7Hours` (change Estimation Date, `Param.Out := addCalendar(In, 0,0,0,0,7,0,0)`) menggeser stempel GMT tengah
// malam ke WIB. Halaman ini menyimpan tanggal sudah dalam zona Jakarta - pergeseran itu tanpa padanan (PARITAS).

import (
	"strconv"

	"github.com/cockroachdb/apd/v3"
)

// Estimasi - baris estimasi ke-e item (o, i).
func Estimasi(h *Halaman, o, i, e int) (Baris, error) {
	return barisDi(h, DaftarDiItem(o, i, AnakEstimasi), e)
}

// bandingTanggal meniru `@defaultCompareDates(a, b)`: "1" bila a sesudah b, "-1" bila sebelum, "0" sama (tingkat
// hari); tanggal kosong = "".
func bandingTanggal(a, b string) string {
	ta, oka := UraiTanggal(a)
	tb, okb := UraiTanggal(b)
	if !oka || !okb {
		return ""
	}
	switch c := hariSaja(ta).Compare(hariSaja(tb)); {
	case c > 0:
		return "1"
	case c < 0:
		return "-1"
	}
	return "0"
}

// periksaTanggalEstimasi = ProtectionDate_Act 2-8 / ValidateInputEstimate_act 7.1-7.7 / CheckEstimateValue 12.1-12.7
// atas satu baris: estimasi sebelum DOL / sesudah hari ini -> pesan pada EstimationDate baris, IsError 2; tanpa galat ->
// IsError 0. Mengembalikan (hasilDOL, tidakLewatHariIni).
func periksaTanggalEstimasi(k *Konteks, h *Halaman, jalurBaris string, b Baris) (string, bool) {
	dol := bandingTanggal(h.Ambil(CD+"DateOfLoss"), b["EstimationDate"])
	if dol == "1" {
		h.TambahPesan(jalurBaris+".EstimationDate", PesanEstimasiKurangDOL)
	}
	hariOK := !SesudahTanggal(b["EstimationDate"], k.Hari())
	if !hariOK {
		h.TambahPesan(jalurBaris+".EstimationDate", PesanEstimasiLebihHari)
	}
	switch {
	case dol == "1" || !hariOK:
		h.Setel(JalurIsError, "2")
	case (dol == "-1" || dol == "0") && hariOK:
		h.Setel(JalurIsError, "0")
	}
	return dol, hariOK
}

// ProtectionDate = `ProtectionDate_Act` (kelas Data-Estimasi; change Estimation Date, DeleteTest 2.4.1).
func ProtectionDate(k *Konteks, h *Halaman, o, i, e int) {
	b, err := Estimasi(h, o, i, e)
	if err != nil {
		return
	}
	periksaTanggalEstimasi(k, h, JalurBaris(DaftarDiItem(o, i, AnakEstimasi), e), b)
}

// TambahEstimasi = tombol Add grid estimasi (`ValidateInputEstimate_act`): baris baru bertanggal hari aksi + kronologi
// "Input Estimation - n"; tanpa spreading DI SELURUH KLAIM -> pesan, baris dibuang, KELUAR (6); tanggal semua baris
// diperiksa; Deductible baris 1 = NetDeductibleValue item; `CopyCurrency` (PersenRNM bagi baris ber-EstimationValue 0);
// `GetNameCurrency_Act`; `SetConvertCurrencyValue_Act`.
func TambahEstimasi(k *Konteks, h *Halaman, o, i int) error {
	it, err := Item(h, o, i)
	if err != nil {
		return err
	}
	daftar := DaftarDiItem(o, i, AnakEstimasi)
	n := h.TambahBaris(daftar, Baris{"pxCreateOperator": k.Pelaku, "pxCreateOpName": k.Pelaku,
		"EstimationDate": k.Waktu()}) // 1
	k.Kronologi(h, "Input Estimation - "+strconv.Itoa(n)) // 1-2
	spread := 0
	for oo := range h.AmbilDaftar(DaftarObjek) { // 3
		for ii := range h.AmbilDaftar(DaftarItem(oo + 1)) {
			spread += len(h.AmbilDaftar(DaftarDiItem(oo+1, ii+1, AnakSpreadPolis)))
		}
	}
	if spread == 0 { // 4, 6
		h.TambahPesan(JalurAnak(DaftarItem(o), i, "CoverageID"), PesanEstimasiTanpaSpread)
		h.HapusBaris(daftar, n)
		return nil // 6 transisi pasca-langkah `true` -> keluar (RALAT 10-10-2026: langkah 7+ tidak berjalan)
	}
	for e, b := range h.AmbilDaftar(daftar) { // 7
		periksaTanggalEstimasi(k, h, JalurBaris(daftar, e+1), b)
		if e == 0 { // 7.8
			b["Deductible"] = it["NetDeductibleValue"]
		} else {
			b["Deductible"] = "0"
		}
	}
	for _, b := range h.AmbilDaftar(daftar) { // 8 CopyCurrency
		if AngkaNol(b["EstimationValue"]) {
			b["PersenRNM"] = h.Ambil(AwalanPolis + ".PercentShare")
		}
	}
	if err := NamaMataUangEstimasi(k, h, o, i); err != nil { // 9
		return err
	}
	return KursEstimasi(k, h, o, i) // 10
}

// AngkaNol - teks bernilai nol (`.EstimationValue == 0`; kosong = 0).
func AngkaNol(s string) bool {
	d, err := AngkaTeks("", s)
	return err == nil && d.IsZero()
}

// KursEstimasi = `SetConvertCurrencyValue_Act` / `SetConvertValueKurs_Estimation` 1: KursValue setiap baris yang belum
// tercetak CFS (`.PrintFaceClaim != 1`) = kurs standar CurrencyID-nya.
func KursEstimasi(k *Konteks, h *Halaman, o, i int) error {
	for _, b := range h.AmbilDaftar(DaftarDiItem(o, i, AnakEstimasi)) {
		if b["PrintFaceClaim"] == "1" || b["CurrencyID"] == "" {
			continue
		}
		v, err := k.Acuan.KursStandar(k.Ctxt(), b["CurrencyID"])
		if err != nil {
			return err
		}
		b["KursValue"] = v
	}
	return nil
}

// NamaMataUangEstimasi = `GetNameCurrency_Act`: Currency setiap baris = nama mata uang CurrencyID (RD BrowseCurrency_RD).
func NamaMataUangEstimasi(k *Konteks, h *Halaman, o, i int) error {
	for _, b := range h.AmbilDaftar(DaftarDiItem(o, i, AnakEstimasi)) {
		if b["CurrencyID"] == "" {
			b["Currency"] = ""
			continue
		}
		v, err := k.Acuan.NamaMataUang(k.Ctxt(), b["CurrencyID"])
		if err != nil {
			return err
		}
		b["Currency"] = v
	}
	return nil
}

// UbahMataUangEstimasi = change Currency baris estimasi: `SetConvertValueKurs_Estimation` lalu `GetNameCurrency_Act`.
func UbahMataUangEstimasi(k *Konteks, h *Halaman, o, i int) error {
	if err := KursEstimasi(k, h, o, i); err != nil {
		return err
	}
	return NamaMataUangEstimasi(k, h, o, i)
}

// mataUangEstimasi - daftar mata uang baris estimasi (urut kemunculan, tanpa kembar - langkah Java "remove value yg
// sama"): ID -> {Currency, KursValue}.
func mataUangEstimasi(rows []Baris) []Baris {
	var out []Baris
	lihat := map[string]bool{}
	for _, b := range rows {
		if lihat[b["CurrencyID"]] {
			continue
		}
		lihat[b["CurrencyID"]] = true
		out = append(out, Baris{"ID": b["CurrencyID"], "Currency": b["Currency"], "KURS": b["KursValue"]})
	}
	return out
}

// susunSpreadKlaim = CheckEstimateValue 13-19 / CountSpreadingClaim_ACT 9-14: SpreadingClaim = SpreadingList x mata
// uang estimasi; ClaimSpreaded = total EstimationValue mata uang itu x Share / 100; TotalClaimSpreaded item.
func susunSpreadKlaim(h *Halaman, o, i int, it Baris) error {
	var k Kalkulator
	est := h.AmbilDaftar(DaftarDiItem(o, i, AnakEstimasi))
	var rows []Baris
	for _, mu := range mataUangEstimasi(est) { // 14
		total := apd.New(0, 0) // 16
		for _, b := range est {
			if b["CurrencyID"] == mu["ID"] {
				total = k.Tambah(total, k.B(b, "EstimationValue"))
			}
		}
		for _, s := range h.AmbilDaftar(DaftarDiItem(o, i, AnakSpreadPolis)) { // 14.2, 17.2
			rows = append(rows, Baris{"TreatyType": s["TreatyType"], "TreatyName": s["TreatyName"],
				"SharePercentage": s["SharePercentage"], "Currency": mu["Currency"], "CurrencyID": mu["ID"],
				"ClaimSpreaded": Teks(k.Persen(total, k.B(s, "SharePercentage")))})
		}
	}
	sum := apd.New(0, 0) // 18
	for _, b := range rows {
		sum = k.Tambah(sum, k.B(b, "ClaimSpreaded"))
	}
	if err := k.Galat(); err != nil {
		return err
	}
	h.SetelDaftar(DaftarDiItem(o, i, AnakSpreadKlaim), rows) // 15
	it["TotalClaimSpreaded"] = Teks(sum)                     // 19
	return nil
}

// CheckEstimateValue = `CheckEstimateValue` (change Estimation Gross (100%) baris estimasi): langkah 2 nama treaty
// spreading, 6-9 rumus baris (Net = Gross - Deductible; RNM = Net x PercentShare / 100; IDR = kurs x RNM), 12 pemeriksaan
// baris (tanggal; Gross dan RNM keduanya 0 -> pesan + baris dibuang), 13-19 SpreadingClaim per mata uang, 21-22 batas
// LoL (Fire) / TSI, 24-25 penanda CFS objek, 26 total item.
//
// Langkah 25 bergerbang satu pzInsKey kasus tertulis mati - DIBUANG (bawaan c prompt §3).
func CheckEstimateValue(k *Konteks, h *Halaman, o, i int) error {
	it, err := Item(h, o, i)
	if err != nil {
		return err
	}
	h.BersihkanPesan()                                                     // 1
	for _, s := range h.AmbilDaftar(DaftarDiItem(o, i, AnakSpreadPolis)) { // 2
		if s["TreatyName"] == "" {
			nama, err := k.Acuan.NamaJenisReas(k.Ctxt(), s["TreatyType"])
			if err != nil {
				return err
			}
			s["TreatyName"] = nama
		}
	}
	if IsMarineCargo(h) && it["ValueTSINusareIDR"] == "" { // 4
		if c := h.AmbilDaftar(DaftarDiItem(o, i, "CoverageList")); len(c) > 0 {
			it["ValueTSINusareIDR"] = c[0]["TSIinIDR"]
		}
	}
	TurunkanIndeksItem(h) // 5
	var kk Kalkulator
	share := kk.H(h, AwalanPolis+".PercentShare") // 6
	daftar := DaftarDiItem(o, i, AnakEstimasi)
	est := h.AmbilDaftar(daftar)
	lebihTSI := len(est) > 0 && est[len(est)-1]["EstimasiMoreThanTSI"] == "true"
	if len(est) > 0 { // 7-8
		if Lebih(apd.New(0, 0), kk.B(est[0], "GrossEstimationPct")) {
			h.TambahPesan(JalurAnak(DaftarItem(o), i, "Test"), PesanEstimasiNegatif)
		}
		if Lebih(apd.New(0, 0), kk.B(est[0], "EstimationValue")) {
			h.TambahPesan(JalurAnak(DaftarItem(o), i, "Test"), PesanEstimasiNegatif)
		}
	}
	rnm, gross, grossIDR, idr := apd.New(0, 0), apd.New(0, 0), apd.New(0, 0), apd.New(0, 0)
	for _, b := range est { // 9
		net := kk.Kurang(kk.B(b, "GrossEstimationPct"), kk.B(b, "Deductible")) // 9.1
		b["NetEstimationValue"] = Teks(net)
		v := kk.Persen(net, share) // 9.2
		b["EstimationValue"] = Teks(v)
		rnm = kk.Tambah(rnm, v)
		it["TotalEstimasi"] = Teks(rnm)
		gross = kk.Tambah(gross, kk.B(b, "GrossEstimationPct"))
		it["TotalGrossEstimasi"] = Teks(gross)
		cg := kk.Kali(kk.B(b, "KursValue"), kk.B(b, "GrossEstimationPct"))
		b["ConvertGrossEstimasi"] = Teks(cg)
		grossIDR = kk.Tambah(grossIDR, cg)
		it["TotalGrossEstimasiIDR"] = Teks(grossIDR)
		cv := kk.Kali(kk.B(b, "KursValue"), v)
		b["ConvertValue"] = Teks(cv)
		idr = kk.Tambah(idr, cv)
		it["EstimastionReserve"] = Teks(idr)
	}
	if err := kk.Galat(); err != nil {
		return err
	}
	errSts, dolTerakhir, hariTerakhir := false, "", true
	var buang []int
	for e, b := range est { // 12
		dolTerakhir, hariTerakhir = periksaTanggalEstimasi(k, h, JalurBaris(daftar, e+1), b)
		if AngkaNol(b["GrossEstimationPct"]) && AngkaNol(b["EstimationValue"]) { // 12.6, 12.8
			h.TambahPesan(JalurAnak(DaftarItem(o), i, "Test"), PesanEstimasiNol)
			errSts = true
			buang = append(buang, e+1)
		}
	}
	for n := len(buang) - 1; n >= 0; n-- { // 12.9 (dibuang sesudah loop, dari belakang)
		h.HapusBaris(daftar, buang[n])
	}
	if err := susunSpreadKlaim(h, o, i, it); err != nil { // 13-19
		return err
	}
	hasil := idr // 20
	lol := kk.B(it, "LimitofLiability")
	if IsFire(h) && !Nol(lol) && Lebih(hasil, kk.Kali(lol, kk.B(it, "KursObjectItem"))) { // 21
		h.TambahPesan(JalurAnak(DaftarItem(o), i, "Test"), PesanEstimasiLebihLoL)
	}
	tsi := kk.B(it, "ValueTSINusareIDR")
	if !lebihTSI && Lebih(hasil, tsi) { // 22-23
		h.TambahPesan(JalurAnak(DaftarItem(o), i, "Test"), PesanEstimasiLebihTSI)
		errSts = true
	}
	if len(est) > 0 && Lebih(apd.New(0, 0), kk.B(est[0], "EstimationValue")) { // 23
		errSts = true
	}
	if err := kk.Galat(); err != nil {
		return err
	}
	ob, err := Objek(h, o)
	if err != nil {
		return err
	}
	switch {
	case h.AdaPesan() || errSts || dolTerakhir == "1" || !hariTerakhir: // 24
		ob["CFS"], ob["PrintFaceClaim"] = "0", "1"
	case lebihTSI || !Lebih(hasil, tsi): // 25
		ob["CFS"], ob["PrintFaceClaim"] = "1", "0"
	}
	it["TotalEstimationValueinIDR"] = Teks(idr) // 26
	it["TotalGrossEstimasiIDR"] = Teks(grossIDR)
	it["TotalGrossEstimasi"] = Teks(gross)
	return nil
}

// HapusEstimasi = tombol Delete grid estimasi (`DeleteValueEstimation`, Param.IdxEstimasi): baris dibuang, kronologi
// "Cancel Input Estimastion - n" (salah ketik verbatim), pemeriksaan tanggal baris tersisa, SpreadingClaim disaring ke
// mata uang yang masih ada lalu `CountSpreadingClaim_ACT`, total item, penanda CFS objek (17).
func HapusEstimasi(k *Konteks, h *Halaman, o, i, e int) error {
	it, err := Item(h, o, i)
	if err != nil {
		return err
	}
	daftar := DaftarDiItem(o, i, AnakEstimasi)
	if _, err := Estimasi(h, o, i, e); err != nil {
		return err
	}
	h.HapusBaris(daftar, e)                                       // 1
	k.Kronologi(h, "Cancel Input Estimastion - "+strconv.Itoa(e)) // 2-3
	est := h.AmbilDaftar(daftar)
	dol := ""
	for n, b := range est { // 4
		dol, _ = periksaTanggalEstimasi(k, h, JalurBaris(daftar, n+1), b)
	}
	ada := map[string]bool{} // 6-9
	for _, b := range est {
		ada[b["CurrencyID"]] = true
	}
	var tetap []Baris
	for _, s := range h.AmbilDaftar(DaftarDiItem(o, i, AnakSpreadKlaim)) {
		if ada[s["CurrencyID"]] {
			tetap = append(tetap, s)
		}
	}
	h.SetelDaftar(DaftarDiItem(o, i, AnakSpreadKlaim), tetap)
	if err := CountSpreadingClaim(h, o, i); err != nil { // 10
		return err
	}
	var kk Kalkulator
	rnm, gross, idr, grossIDR := apd.New(0, 0), apd.New(0, 0), apd.New(0, 0), apd.New(0, 0)
	for _, b := range est { // 4.1-4.3
		rnm = kk.Tambah(rnm, kk.B(b, "EstimationValue"))
		gross = kk.Tambah(gross, kk.B(b, "GrossEstimationPct"))
		if IsMBU(h) {
			gross = kk.Tambah(gross, kk.B(b, "GrossEstimationPctMBU"))
		}
		idr = kk.Tambah(idr, kk.B(b, "ConvertValue"))
		grossIDR = kk.Tambah(grossIDR, kk.B(b, "ConvertGrossEstimasi"))
	}
	if err := kk.Galat(); err != nil {
		return err
	}
	it["TotalEstimationValueinIDR"], it["TotalGrossEstimasi"] = Teks(idr), Teks(gross) // 11 (pre=false)
	it["TotalGrossEstimasiIDR"], it["TotalClaimSpreaded"] = Teks(grossIDR), Teks(rnm)
	ob, err := Objek(h, o)
	if err != nil {
		return err
	}
	if dol == "1" { // 12
		ob["CFS"], ob["PrintFaceClaim"] = "0", "1"
	}
	if Nol(idr) && Nol(rnm) { // 13
		it["TotalClaimSpreaded"], it["TotalEstimationValueinIDR"] = "", ""
	}
	if len(est) == 0 { // 15-16 (estimasi pertama kosong)
		ob["PrintFaceClaim"] = "1"
	}
	for _, b := range est { // 17
		if b["PrintFaceClaim"] == "1" {
			ob["CFS"], ob["PrintFaceClaim"] = "0", "1"
		} else {
			ob["CFS"], ob["PrintFaceClaim"] = "1", "0"
		}
	}
	return nil
}

// CountSpreadingClaim = `CountSpreadingClaim_ACT`: satu mata uang -> ClaimSpreaded SpreadingClaim = total RNM x Share /
// 100 (8); multi mata uang / SpreadingClaim kosong -> disusun ulang per mata uang (9-13); total item (14-16).
func CountSpreadingClaim(h *Halaman, o, i int) error {
	it, err := Item(h, o, i)
	if err != nil {
		return err
	}
	TurunkanIndeksItem(h)
	est := h.AmbilDaftar(DaftarDiItem(o, i, AnakEstimasi))
	mu := mataUangEstimasi(est)
	var kk Kalkulator
	rnm, gross, grossIDR, idr := apd.New(0, 0), apd.New(0, 0), apd.New(0, 0), apd.New(0, 0)
	for _, b := range est { // 5
		rnm = kk.Tambah(rnm, kk.B(b, "EstimationValue"))
		gross = kk.Tambah(gross, kk.B(b, "GrossEstimationPct"))
		grossIDR = kk.Tambah(grossIDR, kk.B(b, "ConvertGrossEstimasi"))
		idr = kk.Tambah(idr, kk.B(b, "ConvertValue"))
	}
	spread := h.AmbilDaftar(DaftarDiItem(o, i, AnakSpreadKlaim))
	if len(mu) <= 1 && len(spread) > 0 { // 8
		sum := apd.New(0, 0)
		for _, s := range spread {
			v := kk.Persen(rnm, kk.B(s, "SharePercentage"))
			s["ClaimSpreaded"] = Teks(v)
			sum = kk.Tambah(sum, v)
			if len(mu) == 1 {
				s["Currency"], s["CurrencyID"] = mu[0]["Currency"], mu[0]["ID"]
			}
		}
		it["TotalClaimSpreaded"] = Teks(sum)
	} else if err := susunSpreadKlaim(h, o, i, it); err != nil { // 9-14
		return err
	}
	if err := kk.Galat(); err != nil {
		return err
	}
	it["TotalEstimationValueinIDR"], it["TotalGrossEstimasiIDR"], it["TotalGrossEstimasi"] = Teks(idr), Teks(grossIDR),
		Teks(gross) // 16
	return nil
}

// CountTSI = `CountTSI_Act` (change % / basis / Amount deductible bentuk 2): basis 1 (Claim Amount) -> nilai = Amount%
// x ClaimDeductible; basis 2 (TSI) -> TSIDeductible = TSINusare, nilai = Amount% x TSI; NetDeductibleValue = yang lebih
// besar antara nilai itu dan DeductibleValue; CurrencyDeductible = mata uang item.
func CountTSI(h *Halaman, o, i int) error {
	it, err := Item(h, o, i)
	if err != nil {
		return err
	}
	var kk Kalkulator
	var nilai *apd.Decimal
	switch it["TypeDeductible"] {
	case "1": // 1.1
		nilai = kk.Persen(kk.B(it, "ClaimDeductible"), kk.B(it, "Amount"))
	case "2": // 2.1
		it["TSIDeductible"] = it["TSINusare"]
		nilai = kk.Persen(kk.B(it, "TSIDeductible"), kk.B(it, "Amount"))
	default:
		return nil
	}
	ded := kk.B(it, "DeductibleValue")
	if err := kk.Galat(); err != nil {
		return err
	}
	if Lebih(nilai, ded) {
		it["NetDeductibleValue"] = Teks(nilai)
	} else {
		it["NetDeductibleValue"] = Teks(ded)
	}
	it["CurrencyDeductible"] = it["CurrencyID"]
	return nil
}
