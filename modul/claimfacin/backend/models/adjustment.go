package models

// Untuk apa berkas ini: PORT PERHITUNGAN ADJUSTMENT (layar Choose Surveyor / "Input Adjustment": Adjusment_SC grid
// Adjustment, panel InputAdjustment) - tambah baris (`SetIndexObject_ACT` + `CountTotalEstimasi_Act`), hapus
// (`DisableSendComite_Act`), Payment Type (`SetAdjTypePayment_act` + DT `SetTempAdjustID`), Gross Adjustment
// (`SetGrossAdjustment_act`), deductible (`SetNilaiResikoSendiri`), adjuster fee (`SetValueAdjusterFee`), salvage
// (`SetSalvageValue`), spreading adjustment (`SetSpreadingAjsutement_Act`, `CountSpreadingAdjustment_ACT`), mata uang
// (`CheckCurrency_ACT` + DT `SetCurrencyAdjustment_DT`, `ProtectCurrAdjustment`), Ex Gratia (`CekExGratia`), payable
// (`SetPayable_Act`, `SetPayableTo_act`), DLA ceding / SOB (`SetDLACedingSOB`), daftar komite tampilan
// (`SetKomiteList_ACT` / `SetListKomite_act`).
//
// Semua rumus uang lewat `Kalkulator` (apd, NUMBER(38,10)); `@divide(x,100,20)` = Persen.

import (
	"strconv"

	"github.com/cockroachdb/apd/v3"
)

// Pesan VERBATIM.
const (
	PesanIsiAdjuster        = "Silahkan isi/pilih ulang Adjuster dan Consultan "                     // CountTotalEstimasi_Act 20
	PesanAdjFinalKomite     = "Adjustment sudah final dan diakseptasi oleh komite"                   // 21.2
	PesanAdjSesudahFinal    = "Tidak bisa menambah adjustment setelah ada payment type: final"       // 21.3
	PesanSpreadGanda        = "Your claim spreading double, please check and contact your IT!"       // ProtetDuplicateSpread 1
	PesanSalvagePertama     = "Can not start adjustment with salvage value"                          // SetAdjTypePayment_act 4.1
	PesanSudahFinal         = "There is already final adjustment"                                    // SetAdjTypePayment_act 8
	PesanGrossNol           = "Gross adjustment value can not be filled by Zero"                     // SetGrossAdjustment_act 4
	PesanMataUangAdj        = "Error Currency Adjustment, Please Check Again !"                      // ProtectCurrAdjustment 1
	PesanAdjLebihEstimasi   = "Total Adjustment RNM should not be greater than total estimation RNM" // SetNilaiResikoSendiri 5
	PesanSalvageMinus       = "Adjustment RNM for Salvage should be minus"                           // SetNilaiResikoSendiri 5
	PesanAdjNol             = "Adjustment RNM can not be filled by Zero"                             // SetNilaiResikoSendiri 5
	PesanFeeLebihEstimasi   = "Adjuster fee value should not be greater than value estimation"       // SetValueAdjusterFee 2
	PesanFeeNol             = "Profesional Fee can not be filled by Zero"                            // SetValueAdjusterFee 2
	PesanSalvageNol         = "Salvage value can not be filled by Zero"                              // SetSalvageValue 2
	PesanSalvageLebih       = "Salvage value  should not be greater than value estimation"           // SetSalvageValue 2
	PesanSalvageNegatif     = "salvage must be negative"                                             // SetSalvageValue 2
	PesanTotalLebihEstimasi = "Total adjustment is greater than the value estimation"                // SetSalvageValue 2
	PesanSharePersenLebih   = "Total Percentage Not More Than 100"                                   // CountSpreadingAdjustment_ACT 1
	PesanSharePersenKurang  = "Total Percentage Not Less Than 100"                                   // CountSpreadingAdjustment_ACT 1
	TeksBatalAdjustment     = "Cancel Input Adjustment - "                                           // DisableSendComite_Act 2
)

// Payment Type adjustment (kode `.PaymentType`; label prompt values tidak diekspor).
const (
	BayarFinal   = "1"
	BayarInterim = "2"
	BayarSalvage = "3"
	BayarFee     = "4"
	BayarAdjust  = "5"
	BayarExpense = "6"
)

// DaftarMataUangAdj - `.CurencyAdjustment` baris adjustment (pilihan "Choose Currency").
const DaftarMataUangAdj = "CurencyAdjustment"

// Adj - baris adjustment ke-a item (o, i).
func Adj(h *Halaman, o, i, a int) (Baris, error) { return barisDi(h, DaftarAdj(o, i), a) }

func bayarDasar(b Baris) bool { // PaymentType 1, 2, 5 (bukan 3, 4, 6)
	p := b["PaymentType"]
	return p != BayarSalvage && p != BayarFee && p != BayarExpense
}

// ---------------------------------------------------------------- tambah / hapus

// TambahAdjustment = tombol "+" grid Adjustment (`SetIndexObject_ACT` lalu `CountTotalEstimasi_Act`).
//
// `[inferensi]` (PARITAS, OQ-CFI-23): tombol itu tidak ber-addRow; CountTotalEstimasi_Act menulis `.Adjustment(<LAST>)`
// - di sini baris baru ditambahkan lalu diisi. Pemeriksaan 20-21 (adjuster / consultant kosong; adjustment terakhir final)
// dijalankan atas baris terakhir YANG SUDAH ADA sebelum penambahan; ada pesan = baris tidak ditambahkan.
func TambahAdjustment(k *Konteks, h *Halaman, o, i int) error {
	it, err := Item(h, o, i)
	if err != nil {
		return err
	}
	if (h.Ambil(CD+"ConsultantID") == "" || h.Ambil(CD+"AppointedADJID") == "") && h.Ambil(JalurIsAnyAccept) != "1" { // 20
		h.TambahPesan(DaftarAdj(o, i), PesanIsiAdjuster)
	}
	if adj := h.AmbilDaftar(DaftarAdj(o, i)); len(adj) > 0 { // 21
		akhir := adj[len(adj)-1]
		switch {
		case akhir["IsApproved"] == "1" && akhir["PaymentType"] == BayarFinal: // 21.2
			h.TambahPesan(DaftarAdj(o, i), PesanAdjFinalKomite)
		case akhir["PaymentType"] == BayarFinal && akhir["IsKomite"] == "": // 21.3
			h.TambahPesan(DaftarAdj(o, i), PesanAdjSesudahFinal)
		}
	}
	if h.AdaPesan() { // 22
		return nil
	}
	TurunkanIndeksItem(h) // SetIndexObject_ACT
	it["IsKomite"] = "1"  // 2
	est := h.AmbilDaftar(DaftarDiItem(o, i, AnakEstimasi))
	for _, mu := range mataUangEstimasi(est) { // 3-5: SpreadingClaim tanpa mata uang
		for _, s := range h.AmbilDaftar(DaftarDiItem(o, i, AnakSpreadKlaim)) {
			if s["CurrencyID"] == "" {
				s["Currency"], s["CurrencyID"] = mu["Currency"], mu["ID"]
			}
		}
	}
	if n := len(h.AmbilDaftar(DaftarDiItem(o, i, AnakBreakQS))); n%2 != 0 { // 6 "JIKA SPREADING BERMASALAH, TARIK ULANG"
		if err := CheckLimit(k, h, o, i); err != nil {
			return err
		}
	}
	var spread []Baris
	for _, s := range h.AmbilDaftar(DaftarDiItem(o, i, AnakSpreadKlaim)) { // 7
		if !AngkaNol(s["ClaimSpreaded"]) {
			spread = append(spread, Baris{"TreatyType": s["TreatyType"], "TreatyName": s["TreatyName"],
				"SharePercentage": s["SharePercentage"], "CurrencyID": s["CurrencyID"], "Currency": s["Currency"],
				"ClaimSpreaded": s["ClaimSpreaded"]})
		}
	}
	if err := protekSpreadGanda(h, o, i, spread); err != nil { // 8
		return err
	}
	if h.AdaPesan() {
		return nil
	}
	muSpread := map[string]bool{} // 9-10
	for _, s := range spread {
		muSpread[s["CurrencyID"]] = true
	}
	var qs []Baris // 11
	for _, b := range h.AmbilDaftar(DaftarDiItem(o, i, AnakBreakQS)) {
		if muSpread[b["CurrencyID"]] {
			qs = append(qs, Baris{"TreatyType": b["TreatyType"], "TreatyName": b["TreatyName"],
				"SharePercentage": b["SharePercentage"], "Currency": b["Currency"], "CurrencyID": b["CurrencyID"]})
		}
	}
	var muAdj []Baris // 12-14
	lihat := map[string]bool{}
	for _, e := range est {
		if e["PrintFaceClaim"] == "1" && !lihat[e["CurrencyID"]] {
			lihat[e["CurrencyID"]] = true
			muAdj = append(muAdj, Baris{"CurrencyID": e["CurrencyID"], "Currency": e["Currency"], "KursValue": e["KursValue"]})
		}
	}
	h.Hapus(CD + "Payable") // 15
	h.Hapus(CD + "PayableTo")
	nama, err := k.Acuan.NamaPelaku(k.Ctxt(), k.Pelaku)
	if err != nil {
		return err
	}
	baru := Baris{"pxCreateOperator": k.Pelaku, "pxCreateOpName": nama, "pxCreateDateTime": k.Waktu()}
	a := h.TambahBaris(DaftarAdj(o, i), baru)
	mataUangAkhir := ""
	if len(est) > 0 {
		mataUangAkhir = est[len(est)-1]["Currency"]
	}
	// 17 - PERBAIKAN (PARITAS `[penyimpangan sadar]`): TotalEstimasiValue = `.TotalEstimasiConvertIDR` item, yang ditulis
	// CLaimFaceSheet_Act 23.4 dari akumulator lokal yang tidak pernah dinolkan antar item (item ke-2 dst. memuat estimasi
	// item sebelumnya); di sini Σ ConvertValue item itu sendiri (TotalEstimationValueinIDR).
	baru["EstimationValue"], baru["TotalEstimasiValue"] = it["TotalEstimasi"], it["TotalEstimationValueinIDR"]
	baru["CurrencyEstimasi"] = mataUangAkhir
	h.SetelDaftar(JalurAnak(DaftarAdj(o, i), a, DaftarMataUangAdj), muAdj) // 17.1
	baru["Payable"] = "2"
	h.Setel(CD+"Payable", "2")
	h.Setel(CD+"PayableTo", h.Ambil(OQ+"SobName"))
	baru["PayableTo"] = h.Ambil(OQ + "SobName")
	baru["DirectToKasir"] = "true"
	h.SetelDaftar(DaftarDiAdj(o, i, a, AnakAdjSpread), spread) // 17.2
	h.SetelDaftar(DaftarDiAdj(o, i, a, AnakAdjQS), qs)
	for _, s := range h.AmbilDaftar(DaftarDiAdj(o, i, a, AnakAdjSpread)) { // 17.3.2
		s["ClaimSpreaded"] = ""
		if s["TreatyName"] == "" {
			nama, err := k.Acuan.NamaJenisReas(k.Ctxt(), s["TreatyType"])
			if err != nil {
				return err
			}
			s["TreatyName"] = nama
		}
	}
	for _, s := range h.AmbilDaftar(DaftarDiItem(o, i, AnakSpreadPolis)) { // 18-19
		if s["TreatyType"] == TreatyFacRetro {
			baru["IsFacRetro"], baru["DirectToKasir"] = "1", "false"
		}
	}
	return nil
}

// protekSpreadGanda = `ProtetDuplicateSpread`: share setiap TreatyType dijumlah; jenis bershare 0 -> pesan (baris
// adjustment baru tidak ditambahkan - langkah 3.3 membuang `.Adjustment(<LAST>)`).
func protekSpreadGanda(h *Halaman, o, i int, spread []Baris) error {
	var kk Kalkulator
	jumlah := map[string]*apd.Decimal{}
	for _, s := range spread {
		if jumlah[s["TreatyType"]] == nil {
			jumlah[s["TreatyType"]] = apd.New(0, 0)
		}
		jumlah[s["TreatyType"]] = kk.Tambah(jumlah[s["TreatyType"]], kk.B(s, "SharePercentage"))
	}
	if err := kk.Galat(); err != nil {
		return err
	}
	for _, v := range jumlah {
		if Nol(v) {
			h.TambahPesan(DaftarAdj(o, i), PesanSpreadGanda)
			return nil
		}
	}
	return nil
}

// HapusAdjustment = tombol hapus baris grid Adjustment (`DisableSendComite_Act`, Param.IdxAdjustment).
func HapusAdjustment(k *Konteks, h *Halaman, o, i, a int) error {
	it, err := Item(h, o, i)
	if err != nil {
		return err
	}
	if _, err := Adj(h, o, i, a); err != nil {
		return err
	}
	h.HapusBaris(DaftarAdj(o, i), a)                    // 1
	it["AdjustmentVal"] = "0"                           // 2
	k.Kronologi(h, TeksBatalAdjustment+strconv.Itoa(a)) // 2-3
	h.Setel(JalurAktifButton, "0")                      // 4
	it["IsKomite"] = ""
	if n := len(h.AmbilDaftar(DaftarAdj(o, i))); n > 0 { // 5 IndexAjustment = subscript baris terakhir
		it[PropIndeksAdj] = strconv.Itoa(n)
	}
	return nil
}

// ---------------------------------------------------------------- Payment Type

// SetAdjTypePayment = change Payment Type (`SetAdjTypePayment_act`).
func SetAdjTypePayment(k *Konteks, h *Halaman, o, i, a int) error {
	it, err := Item(h, o, i)
	if err != nil {
		return err
	}
	b, err := Adj(h, o, i, a)
	if err != nil {
		return err
	}
	b["Currency"] = "" // 3
	pertama := func(item int) string {
		if adj := h.AmbilDaftar(DaftarAdj(o, item)); len(adj) > 0 {
			return adj[0]["PaymentType"]
		}
		return ""
	}
	salvagePertama := false
	switch { // 4-7: lini Marine / MBU membaca item 1, Fire / Aneka item ini
	case IsMarineCargo(h) || IsMBU(h):
		salvagePertama = pertama(1) == BayarSalvage
	case IsFire(h) || IsAneka(h):
		salvagePertama = pertama(i) == BayarSalvage
	}
	if salvagePertama {
		h.Setel(JalurIsError, "2")
		h.TambahPesan(JalurAnak(DaftarAdj(o, i), a, "PaymentType"), PesanSalvagePertama) // 11
	}
	final := 0
	for _, x := range h.AmbilDaftar(DaftarAdj(o, i)) { // 9
		if x["PaymentType"] == BayarFinal {
			final++
		}
	}
	if final > 1 { // 10
		h.TambahPesan(JalurAnak(DaftarAdj(o, i), a, "PaymentType"), PesanSudahFinal)
	}
	kosongkan := func(ps ...string) {
		for _, p := range ps {
			b[p] = ""
		}
	}
	switch b["PaymentType"] {
	case BayarFee: // 13
		kosongkan("AdjustmentValue", "ConvertAdjustmentValue", "IndividualRiskPercentage", "IndividualRiskType",
			"IndividualRiskValue", "ProposeAdjustmentValue", "SalvageValue")
		it["PaymentType"] = BayarFee
	case BayarSalvage: // 15
		kosongkan("AdjusterFeeValue", "ProfessionalFee", "SurveyExpenses", "VAT", "VATValue", "AdjustmentValue",
			"ConvertAdjustmentValue", "IndividualRiskPercentage", "IndividualRiskType", "IndividualRiskValue",
			"ProposeAdjustmentValue")
		it["PaymentType"] = BayarSalvage
	default: // 14
		kosongkan("SalvageValue", "AdjusterFeeValue", "ProfessionalFee", "SurveyExpenses", "VAT", "VATValue")
	}
	// 16-20 (PUCLStatus.IDAdjustment, ShareNusare dari SpreadingList 10007) tidak ditiru: nol pembaca di korpus.
	if b["PaymentType"] != "" && b["PaymentType"] != BayarSalvage { // 21
		h.Setel(JalurIsError, "")
	}
	switch b["PaymentType"] { // DT SetTempAdjustID 5-9
	case BayarFinal, BayarInterim, BayarSalvage, BayarFee, BayarAdjust:
		it["PaymentType"] = b["PaymentType"]
		h.Setel(CD+"PaymentType", b["PaymentType"])
	}
	if b["GrossAdjustment"] != "" { // 23
		return SetGrossAdjustment(k, h, o, i, a)
	}
	return nil
}

// ---------------------------------------------------------------- Gross Adjustment

// SetGrossAdjustment = change / enter Gross Adjustment (`SetGrossAdjustment_act`).
func SetGrossAdjustment(k *Konteks, h *Halaman, o, i, a int) error {
	b, err := Adj(h, o, i, a)
	if err != nil {
		return err
	}
	if b["CurrencyID"] != "" && b["Currency"] == "" { // 1
		nama, err := k.Acuan.NamaMataUang(k.Ctxt(), b["CurrencyID"])
		if err != nil {
			return err
		}
		b["Currency"] = nama
	}
	h.Setel(JalurIsError, "") // 4
	if b["PersenRNM"] == "" { // 5
		b["PersenRNM"] = h.Ambil(AwalanPolis + ".PercentShare")
	}
	var kk Kalkulator
	b["IndividualRiskPercentage"] = "0" // 6 (pre=false: berlaku setiap kali)
	b["ProposeAdjustmentValue"] = Teks(kk.Kurang(kk.B(b, "GrossAdjustment"), kk.B(b, "IndividualRiskValue")))
	kurs := kk.B(b, "CurrencyDol")
	switch p := b["PaymentType"]; {
	case bayarDasar(b): // 7-8
		v := kk.Persen(kk.B(b, "ProposeAdjustmentValue"), kk.B(b, "PersenRNM"))
		b["AdjustmentValue"], b["ValueAdjustment"] = Teks(v), Teks(kk.Kali(kurs, v))
		if err := kk.Galat(); err != nil {
			return err
		}
		if err := SetNilaiResikoSendiri(k, h, o, i, a); err != nil {
			return err
		}
	case p == BayarFee || p == BayarExpense: // 9-10
		if p == BayarFee {
			b["ProfessionalFee"] = b["GrossAdjustment"]
			v := kk.Persen(kk.B(b, "ProfessionalFee"), kk.B(b, "PersenRNM"))
			b["AdjusterFeeValue"], b["ValueAdjustment"] = Teks(v), Teks(kk.Kali(kurs, v))
		}
		if err := kk.Galat(); err != nil {
			return err
		}
		if err := SetValueAdjusterFee(k, h, o, i, a); err != nil {
			return err
		}
	case p == BayarSalvage: // 11-12
		v := kk.Persen(kk.B(b, "GrossAdjustment"), kk.B(b, "PersenRNM"))
		b["SalvageValue"], b["ValueAdjustment"] = Teks(v), Teks(kk.Kali(kurs, v))
		if err := kk.Galat(); err != nil {
			return err
		}
		if err := SetSalvageValue(k, h, o, i, a); err != nil {
			return err
		}
	}
	if AngkaNol(b["GrossAdjustment"]) { // 13
		h.TambahPesan(JalurAnak(DaftarAdj(o, i), a, "GrossAdjustment"), PesanGrossNol)
	}
	if err := SetPayable(k, h, o, i, a); err != nil { // 14
		return err
	}
	b["GrossValue"] = Teks(kk.Persen(kk.B(b, "GrossAdjustment"), kk.B(b, "PersenRNM"))) // 15
	if ob, err := Objek(h, o); err == nil {
		ob["IsKomite"] = "0"
	}
	ProtectCurrAdjustment(h, o, i, a) // 17
	return kk.Galat()
}

// ProtectCurrAdjustment = `ProtectCurrAdjustment`: mata uang adjustment wajib ada di estimasi item yang sudah tercetak
// CFS.
func ProtectCurrAdjustment(h *Halaman, o, i, a int) {
	b, err := Adj(h, o, i, a)
	if err != nil {
		return
	}
	ada := false
	for _, e := range h.AmbilDaftar(DaftarDiItem(o, i, AnakEstimasi)) {
		if e["PrintFaceClaim"] == "1" && e["CurrencyID"] == b["CurrencyID"] {
			ada = true
		}
	}
	if !ada || b["CurrencyID"] == "" {
		h.TambahPesan(JalurAnak(DaftarAdj(o, i), a, "UploadLOD"), PesanMataUangAdj)
	}
}

// ---------------------------------------------------------------- deductible, fee, salvage

// jumlahAdjMataUang - Σ ValueAdjustment adjustment item yang bukan AcceptanceStatus 2 dan ber-mata uang `cur`
// (SetNilaiResikoSendiri 13, SetValueAdjusterFee 6.2, SetSalvageValue 7.2).
func jumlahAdjMataUang(h *Halaman, o, i int, cur string) (*apd.Decimal, error) {
	var kk Kalkulator
	total := apd.New(0, 0)
	for _, x := range h.AmbilDaftar(DaftarAdj(o, i)) {
		if x["AcceptanceStatus"] != "2" && x["CurrencyID"] == cur {
			total = kk.Tambah(total, kk.B(x, "ValueAdjustment"))
		}
	}
	return total, kk.Galat()
}

// bulat4 = `@divide(x, 1, 4)`.
func bulat4(x *apd.Decimal) *apd.Decimal {
	c := apd.BaseContext.WithPrecision(60)
	c.Rounding = apd.RoundHalfUp
	r := new(apd.Decimal)
	_, _ = c.Quantize(r, x, -4)
	return r
}

// PropIndeksAdj - `ObjectItemList(i).IndexAjustment`: adjustment terakhir yang dihitung (dibaca CreateKMTNo_Act dan
// ViewKomite_act); disimpan INDEX_ADJUSTMENT.
const PropIndeksAdj = "IndexAjustment"

// setIndeksAdj = `IndexAjustment := .pxListSubscript` (SetNilaiResikoSendiri 26, SetValueAdjusterFee 17,
// SetSalvageValue 16).
func setIndeksAdj(h *Halaman, o, i, a int) {
	if it, err := Item(h, o, i); err == nil {
		it[PropIndeksAdj] = strconv.Itoa(a)
	}
}

// setAdjVal - penanda tombol kirim komite item (AdjustmentVal / IsAdjVal).
func setAdjVal(h *Halaman, o, i int, v string) {
	if it, err := Item(h, o, i); err == nil {
		it["AdjustmentVal"], it["IsAdjVal"] = v, v
	}
}

// SetNilaiResikoSendiri = change Deductible Type / % / Value (`SetNilaiResikoSendiri`): kurs standar mata uang
// adjustment, deductible menurut tipe (1 % gross, 2 % TSI item, 3 nilai tetap), Adjustment RNM, Σ adjustment item
// dibanding estimasi, penanda kirim komite, spreading adjustment, daftar komite.
func SetNilaiResikoSendiri(k *Konteks, h *Halaman, o, i, a int) error {
	it, err := Item(h, o, i)
	if err != nil {
		return err
	}
	b, err := Adj(h, o, i, a)
	if err != nil {
		return err
	}
	kurs, err := k.Acuan.KursStandar(k.Ctxt(), b["CurrencyID"]) // 2-4
	if err != nil {
		return err
	}
	b["KursIDR"], b["CurrencyDol"] = kurs, kurs
	h.Setel(JalurIsError, "") // 5
	var kk Kalkulator
	if AngkaNol(b["ProposeAdjustmentValue"]) { // 6
		it["IsAdjVal"] = ""
	}
	gross := kk.B(b, "GrossAdjustment")
	switch b["IndividualRiskType"] {
	case "1": // 7
		v := kk.Persen(gross, kk.B(b, "IndividualRiskPercentage"))
		b["IndividualRiskValue"] = Teks(v)
		b["ProposeAdjustmentValue"] = Teks(kk.Kurang(gross, v))
	case "2": // 8-9
		tsi := kk.Kali(kk.B(it, "TSIPerObject"), kk.B(it, "KursObjectItem"))
		v := kk.Persen(kk.BagiPega(tsi, kk.B(b, "CurrencyDol")), kk.B(b, "IndividualRiskPercentage"))
		b["IndividualRiskValue"] = Teks(v)
		b["ProposeAdjustmentValue"] = Teks(kk.Kurang(gross, v))
	case "3": // 11
		b["IndividualRiskPercentage"] = "0"
		b["ProposeAdjustmentValue"] = Teks(kk.Kurang(gross, kk.B(b, "IndividualRiskValue")))
	}
	adjRNM := kk.Persen(kk.B(b, "ProposeAdjustmentValue"), kk.B(b, "PersenRNM")) // 12
	b["AdjustmentValue"] = Teks(adjRNM)
	b["IndividualRiskRNM"] = Teks(kk.Persen(kk.B(b, "IndividualRiskValue"), kk.B(b, "PersenRNM")))
	b["ValueAdjustment"] = Teks(kk.Kali(kk.B(b, "CurrencyDol"), adjRNM))
	if err := kk.Galat(); err != nil {
		return err
	}
	jumlah, err := jumlahAdjMataUang(h, o, i, b["UploadLOD"]) // 13
	if err != nil {
		return err
	}
	tipe := b["PaymentType"]
	bukanSF := tipe != BayarSalvage && tipe != BayarFee
	propNol := AngkaNol(b["ProposeAdjustmentValue"])
	if bukanSF && propNol { // 14-15
		h.Setel(JalurIsError, "2")
		h.TambahPesan(JalurAnak(DaftarAdj(o, i), a, "AdjustmentValue"), PesanAdjNol)
	}
	if bukanSF && Lebih(jumlah, kk.B(b, "TotalEstimasiValue")) { // 16
		h.Setel(JalurIsError, "2")
	}
	if bukanSF && Lebih(bulat4(jumlah), bulat4(kk.B(b, "EstimationValue"))) { // 17
		h.TambahPesan(JalurAnak(DaftarAdj(o, i), a, "AdjustmentValue"), PesanAdjLebihEstimasi)
	}
	if tipe == BayarSalvage && !Lebih(apd.New(0, 0), kk.B(b, "ProposeAdjustmentValue")) { // 23
		h.TambahPesan(JalurAnak(DaftarAdj(o, i), a, "ProposeAdjustmentValue"), PesanSalvageMinus)
	}
	melebihi := Lebih(jumlah, kk.B(b, "EstimationValue"))
	if !propNol && !melebihi { // 24
		setAdjVal(h, o, i, "1")
	}
	if propNol || melebihi { // 25
		setAdjVal(h, o, i, "0")
	}
	if err := kk.Galat(); err != nil {
		return err
	}
	setIndeksAdj(h, o, i, a) // 26
	// 29 SetKomiteList_ACT: daftar komite = turunan tampilan (SusunKomiteAdjustment). Langkah 2-nya menghitung ulang
	// ValueAdjustment SEMUA baris adjustment item dari AdjustmentValue x CurrencyDol (kondisi `PaymentType!="4" ||
	// PaymentType!="3"` selalu benar) - baris fee / salvage (AdjustmentValue kosong) menjadi 0 dan jumlah per mata uang
	// berikutnya salah. `[penyimpangan sadar]` (PARITAS, OQ-CFI-03): tidak ditiru; baris 1 / 2 / 5 sudah memenuhi rumus
	// itu dari langkah 12.
	return SetSpreadingAdjustment(h, o, i, a) // 28
}

// SetValueAdjusterFee = change Survey Expenses / VAT (`SetValueAdjusterFee`, Payment Type 4 / 6).
func SetValueAdjusterFee(_ *Konteks, h *Halaman, o, i, a int) error {
	b, err := Adj(h, o, i, a)
	if err != nil {
		return err
	}
	var kk Kalkulator
	h.Setel(JalurIsError, "")             // 2
	sub := kk.B(b, "GrossAdjustment")     // 2 Local.SubTotal
	vat := kk.Persen(sub, kk.B(b, "VAT")) // 3
	b["VATValue"], b["ProfessionalFee"] = Teks(vat), Teks(kk.Tambah(vat, sub))
	fee := kk.Persen(kk.B(b, "ProfessionalFee"), kk.B(b, "PersenRNM")) // 4
	b["AdjusterFeeValue"] = Teks(fee)
	b["ValueAdjustment"] = Teks(kk.Kali(kk.B(b, "CurrencyDol"), fee))
	b["IndividualRiskRNM"] = Teks(kk.Persen(vat, kk.B(b, "PersenRNM")))
	if err := kk.Galat(); err != nil {
		return err
	}
	feeNol := AngkaNol(b["ProfessionalFee"])
	if feeNol { // 5 / 9
		h.TambahPesan(JalurAnak(DaftarAdj(o, i), a, "ProfessionalFee"), PesanFeeNol)
	}
	jumlah, err := jumlahAdjMataUang(h, o, i, b["UploadLOD"]) // 6
	if err != nil {
		return err
	}
	if feeNol { // 10
		h.Setel(JalurIsError, "2")
		setAdjVal(h, o, i, "0")
	}
	if Lebih(jumlah, kk.B(b, "EstimationValue")) { // 11
		h.Setel(JalurIsError, "2")
		setAdjVal(h, o, i, "0")
	}
	if Lebih(bulat4(jumlah), bulat4(kk.B(b, "EstimationValue"))) { // 12
		h.TambahPesan(JalurAnak(DaftarAdj(o, i), a, "AdjusterFeeValue"), PesanFeeLebihEstimasi)
	}
	// 13 `pyWorkPage.IsError<2 || CountIDR<EstimationValue`: IsError hanya "" (= 0) atau 2.
	if h.Ambil(JalurIsError) != "2" || Lebih(kk.B(b, "EstimationValue"), jumlah) {
		if it, err := Item(h, o, i); err == nil {
			it["AdjustmentVal"] = "1"
		}
	}
	setIndeksAdj(h, o, i, a) // 17
	if err := kk.Galat(); err != nil {
		return err
	}
	if err := SetSpreadingAdjustment(h, o, i, a); err != nil { // 15
		return err
	}
	if it, err := Item(h, o, i); err == nil { // 16
		it["IsAdjVal"] = "1"
	}
	return nil
}

// SetSalvageValue = Payment Type 3 (`SetSalvageValue`).
func SetSalvageValue(_ *Konteks, h *Halaman, o, i, a int) error {
	b, err := Adj(h, o, i, a)
	if err != nil {
		return err
	}
	var kk Kalkulator
	h.Setel(JalurIsError, "")                             // 2
	if Lebih(kk.B(b, "GrossAdjustment"), apd.New(0, 0)) { // 3
		h.TambahPesan(JalurAnak(DaftarAdj(o, i), a, "GrossAdjustment"), PesanSalvageNegatif)
	}
	v := kk.Persen(kk.B(b, "GrossAdjustment"), kk.B(b, "PersenRNM")) // 4
	b["SalvageValue"], b["ValueAdjustment"] = Teks(v), Teks(kk.Kali(kk.B(b, "CurrencyDol"), v))
	if err := kk.Galat(); err != nil {
		return err
	}
	if b["PaymentType"] == BayarSalvage && AngkaNol(b["SalvageValue"]) { // 5
		h.TambahPesan(JalurAnak(DaftarAdj(o, i), a, "SalvageValue"), PesanSalvageNol)
	}
	jumlah, err := jumlahAdjMataUang(h, o, i, b["UploadLOD"]) // 7
	if err != nil {
		return err
	}
	est := kk.B(b, "EstimationValue")
	if Lebih(kk.B(b, "SalvageValue"), est) { // 10-11
		h.Setel(JalurIsError, "2")
		h.TambahPesan(JalurAnak(DaftarAdj(o, i), a, "SalvageValue"), PesanSalvageLebih)
	}
	if Lebih(jumlah, est) { // 12
		h.Setel(JalurIsError, "2")
	}
	if Lebih(bulat4(jumlah), bulat4(est)) { // 13
		h.TambahPesan(JalurAnak(DaftarAdj(o, i), a, "SalvageValue"), PesanTotalLebihEstimasi)
	}
	nol := AngkaNol(b["SalvageValue"])
	if !nol && !Lebih(jumlah, est) { // 14
		setAdjVal(h, o, i, "1")
	}
	if nol || Lebih(jumlah, est) { // 15
		setAdjVal(h, o, i, "0")
	}
	if err := kk.Galat(); err != nil {
		return err
	}
	setIndeksAdj(h, o, i, a)                  // 16
	return SetSpreadingAdjustment(h, o, i, a) // 17
}

// ---------------------------------------------------------------- spreading adjustment

// SetSpreadingAdjustment = `SetSpreadingAjsutement_Act`: share SpreadingAdjustment dijumlah per TreatyType (langkah 3-6;
// langkah 4 Java tidak diekspor - `[inferensi]` membuang baris TreatyType kembar), ClaimSpreaded = nilai dasar menurut
// Payment Type x share (7 / 9 / 11), Break QS dari Claim Spreaded baris "QS"+"FAC" (8 / 10 / 12), total.
func SetSpreadingAdjustment(h *Halaman, o, i, a int) error {
	b, err := Adj(h, o, i, a)
	if err != nil {
		return err
	}
	it, err := Item(h, o, i)
	if err != nil {
		return err
	}
	var kk Kalkulator
	daftar := DaftarDiAdj(o, i, a, AnakAdjSpread)
	var gabung []Baris // 2-6
	for _, s := range h.AmbilDaftar(daftar) {
		var sama Baris
		for _, g := range gabung {
			if g["TreatyType"] == s["TreatyType"] {
				sama = g
			}
		}
		if sama == nil {
			sama = s.Salin()
			sama["SharePercentage"] = "0"
			gabung = append(gabung, sama)
		}
		sama["SharePercentage"] = Teks(kk.Tambah(kk.B(sama, "SharePercentage"), kk.B(s, "SharePercentage")))
	}
	h.SetelDaftar(daftar, gabung)
	var dasar *apd.Decimal
	switch p := b["PaymentType"]; p {
	case BayarFinal, BayarInterim, BayarAdjust:
		dasar = kk.B(b, "AdjustmentValue")
	case BayarFee, BayarExpense:
		dasar = kk.B(b, "AdjusterFeeValue")
	case BayarSalvage:
		dasar = kk.B(b, "SalvageValue")
	}
	totAdj, totPersen, totBreak := apd.New(0, 0), apd.New(0, 0), apd.New(0, 0)
	breakDasar := apd.New(0, 0)
	if dasar != nil {
		for _, s := range h.AmbilDaftar(daftar) {
			v := kk.Persen(dasar, kk.B(s, "SharePercentage"))
			s["ClaimSpreaded"] = Teks(v)
			totAdj = kk.Tambah(totAdj, v)
			s["TotalSpread"] = Teks(totAdj)
			totPersen = kk.Tambah(totPersen, kk.B(s, "SharePercentage"))
			if qsFac(s["TreatyName"]) {
				breakDasar = v
			}
		}
		// 8 / 10 / 12 - kondisi `1 || 2 || 5 && CountSpreadQS!=0` (dst.) berulang atas SpreadingQuotaShare: daftar
		// kosong = tidak ada yang dihitung, jadi setiap Payment Type 1-6 menghitung bila barisnya ada.
		for _, q := range h.AmbilDaftar(DaftarDiAdj(o, i, a, AnakAdjQS)) {
			v := kk.Persen(breakDasar, kk.B(q, "SharePercentage"))
			q["ClaimSpreaded"] = Teks(v)
			totBreak = kk.Tambah(totBreak, v)
			q["TotalSpread"] = Teks(totBreak)
		}
	}
	if err := kk.Galat(); err != nil {
		return err
	}
	it["TotalHasilClaim"] = Teks(totAdj) // 13
	b["TotalSharePersen"], b["TotalSpreadBreakQs"], b["TotalSpreadAdjustment"] = Teks(totPersen), Teks(totBreak), Teks(totAdj)
	return nil
}

// CountSpreadingAdjustment = change Share % Spreading Adjustment Ex Gratia (`CountSpreadingAdjustment_ACT`).
func CountSpreadingAdjustment(h *Halaman, o, i, a int) error {
	b, err := Adj(h, o, i, a)
	if err != nil {
		return err
	}
	var kk Kalkulator
	var dasar *apd.Decimal
	switch b["PaymentType"] { // 2-4
	case BayarFinal, BayarInterim, BayarAdjust:
		dasar = kk.B(b, "AdjustmentValue")
	case BayarSalvage:
		dasar = kk.B(b, "SalvageValue")
	case BayarFee, BayarExpense:
		dasar = kk.B(b, "AdjusterFeeValue")
	default:
		dasar = apd.New(0, 0)
	}
	persen, total := apd.New(0, 0), apd.New(0, 0)
	daftar := DaftarDiAdj(o, i, a, AnakAdjSpread)
	for n, s := range h.AmbilDaftar(daftar) { // 5
		persen = kk.Tambah(persen, kk.B(s, "SharePercentage"))
		if Lebih(persen, apd.New(100, 0)) { // 5.2
			h.TambahPesan(JalurAnak(daftar, n+1, "SharePercentage"), PesanSharePersenLebih)
		}
		v := kk.Bagi(kk.Kali(kk.B(s, "SharePercentage"), dasar), apd.New(100, 0))
		s["ClaimSpreaded"] = Teks(v)
		total = kk.Tambah(total, v)
	}
	if err := kk.Galat(); err != nil {
		return err
	}
	b["TotalSharePersen"], b["TotalSpreadAdjustment"] = Teks(persen), Teks(total) // 6
	if Lebih(apd.New(100, 0), persen) {                                           // 7
		h.TambahPesan(JalurAnak(DaftarAdj(o, i), a, "TotalSharePersen"), PesanSharePersenKurang)
	}
	return nil
}

// spreadDariKlaim - SpreadingAdjustment / SpreadingQuotaShare dari SpreadingClaim dan Break QS item ber-mata uang
// `cur` (CheckCurrency_ACT 4 / 7, CekExGratia 6 / 8).
func spreadDariKlaim(k *Konteks, h *Halaman, o, i, a int, cur string, namaUlang bool) error {
	var spread []Baris
	for _, s := range h.AmbilDaftar(DaftarDiItem(o, i, AnakSpreadKlaim)) {
		if s["CurrencyID"] == cur && !AngkaNol(s["ClaimSpreaded"]) {
			spread = append(spread, Baris{"TreatyType": s["TreatyType"], "TreatyName": s["TreatyName"],
				"SharePercentage": s["SharePercentage"], "Currency": s["Currency"], "CurrencyID": s["CurrencyID"]})
		}
	}
	h.SetelDaftar(DaftarDiAdj(o, i, a, AnakAdjSpread), spread)
	if namaUlang { // CheckCurrency_ACT 5 (pre=false: semua baris)
		for _, s := range spread {
			s["ClaimSpreaded"] = ""
			nama, err := k.Acuan.NamaJenisReas(k.Ctxt(), s["TreatyType"])
			if err != nil {
				return err
			}
			s["TreatyName"] = nama
		}
	}
	var qs []Baris
	if len(spread) != 0 {
		for _, b := range h.AmbilDaftar(DaftarDiItem(o, i, AnakBreakQS)) {
			if b["CurrencyID"] == cur {
				qs = append(qs, Baris{"TreatyType": b["TreatyType"], "TreatyName": b["TreatyName"],
					"SharePercentage": b["SharePercentage"], "Currency": b["Currency"], "CurrencyID": b["CurrencyID"]})
			}
		}
	}
	h.SetelDaftar(DaftarDiAdj(o, i, a, AnakAdjQS), qs)
	return nil
}

// CheckCurrency = change "Choose Currency" (pre-DT `SetCurrencyAdjustment_DT` lalu `CheckCurrency_ACT`): mata uang, kurs
// (CurrencyDol) dari CurencyAdjustment baris terpilih, spreading adjustment sesuai mata uang, EstimationValue = Σ
// estimasi RNM item bermata uang itu x kurs.
func CheckCurrency(k *Konteks, h *Halaman, o, i, a int) error {
	b, err := Adj(h, o, i, a)
	if err != nil {
		return err
	}
	pilih := b["UploadLOD"]
	var kurs, nama, id string // DT 1
	for _, m := range h.AmbilDaftar(JalurAnak(DaftarAdj(o, i), a, DaftarMataUangAdj)) {
		if m["CurrencyID"] == pilih {
			kurs, nama, id = m["KursValue"], m["Currency"], m["CurrencyID"]
		}
	}
	b["CurrencyID"], b["CurrencyDol"], b["Currency"], b["CurrencyName"] = id, kurs, nama, nama // DT 2-5
	if b["CurrencyID"] != "" && b["Currency"] == "" {                                          // 1
		if b["Currency"], err = k.Acuan.NamaMataUang(k.Ctxt(), b["CurrencyID"]); err != nil {
			return err
		}
	}
	if err := spreadDariKlaim(k, h, o, i, a, pilih, true); err != nil { // 2-7
		return err
	}
	var kk Kalkulator
	total := apd.New(0, 0)
	for _, e := range h.AmbilDaftar(DaftarDiItem(o, i, AnakEstimasi)) { // 8
		if e["CurrencyID"] == pilih {
			total = kk.Tambah(total, kk.Kali(kk.B(e, "EstimationValue"), kk.B(e, "KursValue")))
		}
	}
	b["EstimationValue"] = Teks(total) // 9
	return kk.Galat()
}

// CekExGratia = click Ex Gratia (`CekExGratia`): 1 -> QS breakdown dibuang (spreading adjustment dapat disunting); 0 ->
// spreading adjustment disusun ulang dari SpreadingClaim mata uang terpilih.
func CekExGratia(k *Konteks, h *Halaman, o, i, a int) error {
	b, err := Adj(h, o, i, a)
	if err != nil {
		return err
	}
	if b["ExGratia"] == "1" { // 2-3
		h.SetelDaftar(DaftarDiAdj(o, i, a, AnakAdjQS), nil)
		return nil
	}
	return spreadDariKlaim(k, h, o, i, a, b["UploadLOD"], false) // 4-8
}

// ---------------------------------------------------------------- payable

// JumlahCeding - `@SizeOfPropertyList(OfferFacIn.QuotationData.CedingCoList)`.
func JumlahCeding(h *Halaman) int { return len(h.AmbilDaftar(DaftarCedingCo)) }

// RencanaRekening - rekening yang perlu dibaca SetPayable_Act: klien (ParamID.CARI1) dan mata uang (CARI2).
type RencanaRekening struct {
	Klien, MataUang string
	// SemuaKlien - Payable 3: GetDataBankAccount2_sql (menurut mata uang saja).
	SemuaKlien bool
}

// SetPayable = `SetPayable_Act` bagian murni (langkah 1-5, 8.1): payable adjustment dan klaim menurut
// `ClaimData.Payable` (1 ceding / cedant panel, 2 SOB, 3 lain), mengembalikan rekening yang harus dibaca.
func SetPayable(_ *Konteks, h *Halaman, o, i, a int) error {
	_, err := RencanaPayable(h, o, i, a)
	return err
}

// RencanaPayable - langkah 1-5 / 8.1 SetPayable_Act; rekening dibaca services (`TerapkanRekening`).
func RencanaPayable(h *Halaman, o, i, a int) (RencanaRekening, error) {
	b, err := Adj(h, o, i, a)
	if err != nil {
		return RencanaRekening{}, err
	}
	b["NameOfBank"], b["NoAccount"], b["BranchOfBank"] = "", "", "" // 1 (pre=false: selalu)
	pay := h.Ambil(CD + "Payable")
	r := RencanaRekening{MataUang: b["CurrencyID"]} // 2
	if JumlahCeding(h) > 1 && pay != "3" {          // 3
		b["PayableTo"] = AmbilJalur(h, JalurBaris(DaftarCedingCo, 1)+".CedingCoName")
		b["Payable"] = pay
		h.Setel(CD+"PayableTo", b["PayableTo"])
		r.Klien = h.Ambil(OQ + "CedingCo")
	}
	if pay == "1" { // 4
		cedant := h.AmbilDaftar(DaftarDiAdj(o, i, a, AnakCedingCedant))
		switch len(cedant) {
		case 0: // 4.1
			h.Setel(CD+"PayableTo", h.Ambil(OQ+"CedingCoName"))
			b["Payable"], b["PayableTo"] = pay, h.Ambil(OQ+"CedingCoName")
			r.Klien = h.Ambil(OQ + "CedingCo")
		case 1: // 4.2
			h.Setel(CD+"PayableTo", cedant[0]["CedingCoName"])
			b["Payable"], b["PayableTo"] = pay, cedant[0]["CedingCoName"]
			r.Klien = cedant[0]["CedingCo"]
		}
	}
	if pay == "2" { // 5
		h.Setel(CD+"PayableTo", h.Ambil(OQ+"SobName"))
		b["Payable"], b["PayableTo"] = pay, h.Ambil(OQ+"SobName")
		r.Klien = h.Ambil(OQ + "SourceOfBusiness")
	}
	if pay == "3" { // 8.1
		b["Payable"], b["PayableTo"] = pay, h.Ambil(CD+"PayableTo")
		r.SemuaKlien = true
	}
	return r, nil
}

// TerapkanRekening = SetPayable_Act langkah 6-7 / 8.2: rekening pertama hasil bacaan ke adjustment - Payable 3 hanya
// membaca daftar (pilihan Name of Bank). `ClaimData.ReceiverClaim(1)` (langkah 7) tidak ditiru: nol pembaca di korpus
// (hanya ditulis SetPayable_Act / SetPayableTo_act dan dibuang BackToRegister_act).
func TerapkanRekening(h *Halaman, o, i, a int, rek []RekeningBank) error {
	b, err := Adj(h, o, i, a)
	if err != nil {
		return err
	}
	if h.Ambil(CD+"Payable") == "3" || len(rek) == 0 {
		return nil
	}
	r := rek[0]
	b["NameOfBank"], b["NoAccount"], b["BranchOfBank"], b["SwiftCode"], b["IDOfBank"] = r.NameOfBank, r.AccountNo,
		r.BranchOfBank, r.SwiftCode, r.IDOfBank
	return nil
}

// KunciRekening - nilai pilihan Name of Bank (dibaca ulang server).
func KunciRekening(r RekeningBank) string { return r.AccountNo + "|" + r.NameOfBank }

// PilihRekening = autocomplete Name of Bank (pageList Result.pxResults; isi BRANCHOFBANK, ACCOUNTNO, SWIFTCODE,
// IDOFBANK) lalu `SetPayableTo_act` (langkah 1-4 ber-remark; 6 Payable 3 -> ReceiverClaim(1)).
func PilihRekening(h *Halaman, o, i, a int, r RekeningBank) error {
	b, err := Adj(h, o, i, a)
	if err != nil {
		return err
	}
	b["NameOfBank"], b["BranchOfBank"], b["NoAccount"], b["SwiftCode"], b["IDOfBank"] = r.NameOfBank, r.BranchOfBank,
		r.AccountNo, r.SwiftCode, r.IDOfBank
	return nil
}

// SetDLACedingSOB = change DLA No Ceding / SOB (`SetDLACedingSOB`): nilai terisi disalin ke klaim.
func SetDLACedingSOB(h *Halaman, o, i, a int) error {
	b, err := Adj(h, o, i, a)
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
