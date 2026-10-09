package services

// Baris DETAIL kontrak tuntas — `SaveTreatyInDetail_Act` (Treaty In →
// `TREATYINDETAIL`) dan `SaveTreatyInDetailEdm_Act` (Adjustment →
// `TREATYINDETAILEDM`). Perintah WO 9 Oktober 2026; lihat
// `models/detail_treaty.go`.
//
// ---------------------------------------------------------------------
// ⭐ TERJEMAHAN LANGKAH DEMI LANGKAH, BUKAN RINGKASAN
// ---------------------------------------------------------------------
// Halaman `SaveData` Pega TIDAK dikosongkan antarbaris: ia dibuang sekali di
// langkah [2] lalu ditimpa sepanjang perulangan, dan setiap `RDB-List
// SaveTreatyInDetail` menyisipkan ISINYA SAAT ITU. Nilai dari baris
// sebelumnya ikut terbawa ke baris berikutnya bila langkah berikutnya tidak
// menimpanya — perilaku yang ditiru di sini dengan SATU peta `saveData`
// yang hidup sepanjang susunan, dan salinannya per sisipan.
//
// Kode aksi precondition ekspor (disapu dari seluruh korpus):
//
//	2  Continue Whens  — jalankan langkah
//	3  Skip Step       — lewati langkah (beserta anaknya)
//	4  Exit Iteration  — keluar dari iterasi perulangan terdalam
//
// Langkah ber-`pyStepsPreCondition = false` mengabaikan WHEN-nya. Langkah
// tanpa metode (kolom metode kosong) TIDAK menjalankan parameter
// `Property-Set` sisa yang masih tersimpan di XML-nya — `[5.4.2.1.1.2.2.1]`
// dan `[5.4.3.3.2.2.2.1]` (`DEDUCTION1/2 := 0`) adalah dua contohnya.
//
// ---------------------------------------------------------------------
// ⚠️ TIGA HAL YANG DIIKUTI APA ADANYA, DAN DILAPORKAN
// ---------------------------------------------------------------------
//  1. Treaty In NonProp "spreading type 1" (`[5.4.2]`): tidak satu langkah
//     pun mengisi `SaveData.LIMITCURRENCY`, sehingga syarat Gross Premium
//     `.Currency == SaveData.LIMITCURRENCY` (gagal → Exit Iteration) tidak
//     pernah benar untuk mata uang terisi — cabang ini menyisipkan NOL baris.
//     Cabang EDM yang setara (`[7.2.2.1]`) dan cabang "type > 1" Treaty In
//     (`[5.4.3.2.1]`) mempunyai langkah `RNM Limit` yang mengisinya.
//  2. EDM Proportional (`[6]`): empat `RDB-List` sisipan ber-EPI memanggil
//     `SaveTreatyInDetail` (tabel Treaty In), bukan `…Edm`. Perintah WO
//     menyebut `TREATYINDETAILEDM` untuk Adjustment; seluruh baris EDM
//     ditulis ke sana (data DEV: 363 baris EDM Prop ber-EPI memang di sana).
//  3. `CLASSOFBUSINESS`/`CLASSOFBUSINESSID` tidak diisi: cabang ber-COB
//     (`[4]`/`[5]` EDM) ber-`//`, dan cabang Treaty In yang hidup bernama
//     "No COB".

import (
	"fmt"
	"strings"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/modul/treatyin/backend/models"
)

// saveData - halaman `SaveData`: properti → teks Pega. Kuncinya NAMA KOLOM
// (`ADJ_RATE` = `SaveData.ADJUSMENT_RATE`).
type saveData map[string]string

// penyusunDetail - satu jalannya activity detail.
type penyusunDetail struct {
	doc   map[string]any
	sd    saveData
	baris []saveData
	galat error
}

// sisip - `RDB-List SaveTreatyInDetail(Edm)`: isi `SaveData` SAAT INI.
func (p *penyusunDetail) sisip() {
	salin := make(saveData, len(p.sd))
	for k, v := range p.sd {
		salin[k] = v
	}
	p.baris = append(p.baris, salin)
}

// SusunDetailTreatyIn - `SaveTreatyInDetail_Act` atas dokumen `TreatyIn`
// kontrak `id`; hasilnya baris `TREATYINDETAIL` urut sisipan.
func SusunDetailTreatyIn(doc map[string]any, id string) ([]models.BarisDetailTreaty, error) {
	p := mulaiDetail(doc, id)
	switch teksDok(doc, "ProportionType") {
	case "Proportional": // [4]
		p.cabangProporsional()
	case "NonProportional": // [5]
		p.cabangNonProporsionalTreatyIn()
	}
	return p.selesai(models.KolomDetailTreatyIn)
}

// SusunDetailTreatyInEDM - `SaveTreatyInDetailEdm_Act` atas dokumen
// penyesuaian `id`; hasilnya baris `TREATYINDETAILEDM` urut sisipan.
//
// `[4]` dan `[5]` (cabang ber-COB) ber-`//` — yang hidup `[6]` dan `[7]`.
func SusunDetailTreatyInEDM(doc map[string]any, id string) ([]models.BarisDetailTreaty, error) {
	p := mulaiDetail(doc, id)
	switch teksDok(doc, "ProportionType") {
	case "Proportional": // [6] — identik dengan Treaty In [4]
		p.cabangProporsional()
	case "NonProportional": // [7]
		p.cabangNonProporsionalEDM()
	}
	return p.selesai(models.KolomDetailTreatyInEDM)
}

// mulaiDetail - [1] RemoveTreatyInDetail (di repository), [2] Page-Remove
// SaveData, [3] medan kepala.
func mulaiDetail(doc map[string]any, id string) *penyusunDetail {
	p := &penyusunDetail{doc: doc, sd: saveData{}}
	p.sd["TREATYID"] = id
	p.sd["PROPORTIONTYPE"] = teksDok(doc, "ProportionType")
	p.sd["SOB"] = teksDok(doc, "LeadingReinsSource")
	p.sd["SOBID"] = teksDok(doc, "LeadingReinsSourceID")
	p.sd["CEDING"] = teksDok(doc, "Ceding")
	p.sd["CEDINGID"] = teksDok(doc, "CedingID")
	p.sd["TREATYYEAR"] = teksDok(doc, "TreatyYear")
	p.sd["TREATYCONTRACTNAME"] = teksDok(doc, "TreatyContractName")
	return p
}

// cabangProporsional - Treaty In [4] = EDM [6] ("When Proportional (No COB)").
func (p *penyusunDetail) cabangProporsional() {
	for _, lim := range barisDok(p.doc["Limits"]) {
		// [4.1]
		p.sd["TREATYTYPE"] = teksDok(lim, "TreatyType")
		p.sd["BROKERAGE"] = teksDok(p.doc, "BrokeragePercentP")
		for _, det := range barisDok(lim["Detail"]) {
			// [4.2.1] Treaty group & OGR ONR
			p.sd["TREATYGROUP"] = teksDok(det, "TreatyGroup")
			p.sd["TREATYGROUPID"] = teksDok(det, "TreatyGroupID")
			p.sd["RIOGR"] = p.persenRI(teksDok(det, "RIOGR"), "RIOGR")
			p.sd["RIONR"] = p.persenRI(teksDok(det, "RIONR"), "RIONR")
			p.sd["RNM_SHARE"] = teksDok(det, "RNMShare")
			if p.sd["RNM_SHARE"] == "" {
				p.sd["RNM_SHARE"] = teksDok(p.doc, "RNMShareP")
			}
			p.sd["SPREADINGTYPEID"] = teksDok(det, "SpreadingTypeID")
			p.sd["SPREADINGTYPE"] = teksDok(det, "SpreadingType")
			p.sd["CURRENCY_RSMD"] = teksDok(det, "CurrencyRSMD")
			p.sd["CURRENCY_EQ"] = teksDok(det, "CurrencyEarthquake")
			p.sd["CURRENCY_FLOOD_JABO"] = teksDok(det, "CurrencyFloodJab")
			p.sd["CURRENCY_FLOOD_NATION"] = teksDok(det, "CurrencyFloodNat")
			p.sd["LIMIT_RSMD"] = teksDok(det, "RSMDLimit")
			p.sd["LIMIT_EQ"] = teksDok(det, "Earthquake")
			p.sd["LIMIT_FLOOD_JABO"] = teksDok(det, "FloodJab")
			p.sd["LIMIT_FLOOD_NATION"] = teksDok(det, "FloodNation")
			p.sd["QS_PCT"] = teksDok(det, "QSPct")
			p.sd["SPL_LINE"] = teksDok(det, "Surplus")
			p.detailProporsional(det)
		}
	}
}

// detailProporsional - [4.2.2]–[4.2.4] satu baris `Detail`.
func (p *penyusunDetail) detailProporsional(det map[string]any) {
	ioo := barisDok(det["IOOLimitList"])
	if len(ioo) < 1 { // [4.2.2]–[4.2.3]
		p.sisip()
	}
	for _, l := range ioo { // [4.2.4]
		p.sd["LIMITCURRENCY"] = teksDok(l, "Currency")
		p.sd["LIMITVALUE"] = teksDok(l, "Value")
		retensi := barisDok(det["RetentionList"])
		if len(retensi) < 1 {
			p.sisip()
		}
		for _, r := range retensi { // [4.2.4.4]
			if p.sd["LIMITCURRENCY"] != teksDok(r, "Currency") { // CEK CURR, F=3
				continue
			}
			p.sd["RETENTIONCURRENCY"] = teksDok(r, "Currency")
			p.sd["RETENTIONVALUE"] = teksDok(r, "Value")
			cession := barisDok(det["CessionList"])
			if len(cession) < 1 {
				p.sisip()
			}
			for _, c := range cession {
				if p.sd["LIMITCURRENCY"] != teksDok(c, "Currency") {
					continue
				}
				p.sd["CESSIONCURRENCY"] = teksDok(c, "Currency")
				p.sd["CESSIONVALUE"] = teksDok(c, "Value")
				bagian := barisDok(det["RNMShareList"])
				if len(bagian) < 1 {
					p.sisip()
				}
				for _, s := range bagian {
					if p.sd["LIMITCURRENCY"] != teksDok(s, "Currency") {
						continue
					}
					p.sd["SHARECURRENCY"] = teksDok(s, "Currency")
					p.sd["SHAREVALUE"] = teksDok(s, "Value")
					p.spreadingProporsional(det)
				}
			}
		}
	}
}

// spreadingProporsional - [4.2.4.4.1.4.1.4.1.2] (type 1) dan [….3]
// (type > 1). Syarat keduanya membaca `Detail.SpreadingTypeID`, bukan
// `SaveData` — jadi tepat satu yang jalan.
func (p *penyusunDetail) spreadingProporsional(det map[string]any) {
	sebar := barisDok(det["SpreadingList"])
	epi := barisDok(det["EPIList"])
	if teksDok(det, "SpreadingTypeID") != "" { // UNTUK SPREADING TYPE 1
		for _, s := range sebar {
			switch teksDok(s, "ReinsTypeID") {
			case "10004": // QS (R/I)
				p.sd["QS_RI"] = teksDok(s, "Pct")
			case "10028": // QS (OR)
				p.sd["QS_OR"] = teksDok(s, "Pct")
			}
		}
		p.sisipEPI(epi)
		return
	}
	// UNTUK SPREADING TYPE > 1 — metode langkah perulangan jalan di tiap
	// iterasi SEBELUM anak-anaknya.
	for _, s := range sebar {
		p.sd["SPREADINGTYPEID"] = teksDok(s, "ReinsTypeID")
		p.sd["SPREADINGTYPE"] = teksDok(s, "ReinsTypeName")
		p.sd["SPREAD_RNM_SHARE_VALUE"] = "0"
		p.sd["SPREAD_RNM_SHARE_PCT"] = teksDok(s, "Pct")
		for _, b := range barisDok(s["BreakDownSprdList"]) {
			p.sd["SPREAD_RNM_SHARE_VALUE"] = p.jumlah(p.sd["SPREAD_RNM_SHARE_VALUE"], teksDok(b, "Amount"), "BreakDownSprdList.Amount")
			switch teksDok(b, "ReinsID") {
			case "10004":
				p.sd["QS_RI"] = teksDok(b, "SharePct")
			case "10028":
				p.sd["QS_OR"] = teksDok(b, "SharePct")
			}
		}
		p.sisipEPI(epi)
	}
}

// sisipEPI - satu sisipan bila `EPIList` kosong, selainnya satu per baris EPI.
func (p *penyusunDetail) sisipEPI(epi []map[string]any) {
	if len(epi) < 1 {
		p.sisip()
	}
	for _, e := range epi {
		p.sd["EPICURRENCY"] = teksDok(e, "Currency")
		p.sd["EPIVALUE"] = teksDok(e, "Value")
		p.sisip()
	}
}

// cabangNonProporsionalTreatyIn - Treaty In [5] ("When Non Proportional
// (No COB)").
func (p *penyusunDetail) cabangNonProporsionalTreatyIn() {
	limits := barisDok(p.doc["Limits"])
	for i, sh := range barisDok(p.doc["Share"]) {
		// [5.1] — `TreatyIn.Limits(Local.idx_limits)`: layer BERINDEKS SAMA.
		lim := elemenKe(limits, i)
		p.kepalaShare(sh)
		p.sd["SPREADINGTYPE"] = teksDok(sh, "SpreadingTypeXOL")
		p.sd["SPREADINGTYPEID"] = teksDok(sh, "SpreadingTypeIDXOL")
		p.sd["BROKERAGE"] = teksDok(p.doc, "BrokeragePercent")
		p.sd["MDP_PCT"] = teksDok(lim, "MDPPct")
		p.sd["ROL_PCT"] = teksDok(lim, "ROLPct")
		p.sd["DEDUCTIBLE"] = teksDok(lim, "Deductible")
		p.sd["DEDUCTIBLE2"] = teksDok(lim, "Deductible2")
		p.sd["ADJ_RATE"] = teksDok(lim, "AdjRate")
		for _, x := range barisDok(lim["PremiumEarnedList"]) { // [5.2]
			p.sd["PREMIUM_EARNED"] = teksDok(x, "Value")
		}
		for _, x := range barisDok(lim["MDPList"]) { // [5.3]
			p.sd["MDP"] = teksDok(x, "Value")
		}
		for _, g := range barisDok(sh["TreatyGroupList"]) { // [5.4]
			p.sd["TREATYGROUP"] = teksDok(g, "TreatyGroup")
			p.sd["TREATYGROUPID"] = teksDok(g, "TreatyGroupID")
			if teksDok(sh, "SpreadingTypeXOL") != "" { // [5.4.2] UNTUK SPREADING TYPE 1
				// ⚠️ Nol langkah `RNM Limit` di sini (lihat kepala berkas, 1).
				for range barisDok(sh["RnmLimitList"]) {
					p.premiShare(sh)
				}
			} else { // [5.4.3] UNTUK SPREADING TYPE > 1
				p.spreadingShare(sh)
			}
		}
	}
}

// spreadingShare - Treaty In [5.4.3], satu baris `Share`.
func (p *penyusunDetail) spreadingShare(sh map[string]any) {
	for _, sp := range barisDok(sh["SpreadingListXOL"]) {
		p.sd["SPREAD_RNM_SHARE_PCT"] = teksDok(sp, "Pct")
		p.sd["SPREAD_RNM_SHARE_VALUE"] = teksDok(sp, "Value")
		p.sd["SPREADINGTYPE"] = teksDok(sp, "ReinsTypeName")
		p.sd["SPREADINGTYPEID"] = teksDok(sp, "ReinsTypeID")
		for _, b := range barisDok(sp["BreakDownSprdListXOL"]) { // [5.4.3.1]
			switch teksDok(b, "ReinsTypeID") {
			case "10004":
				p.sd["QS_RI"] = teksDok(b, "SharePct")
			case "10028":
				p.sd["QS_OR"] = teksDok(b, "SharePct")
			}
		}
		for _, r := range barisDok(sh["RnmLimitList"]) { // [5.4.3.2]
			p.batasRNM(r)
			for _, n := range barisDok(sh["NetPremiumList"]) { // [5.4.3.2.2]
				if teksDok(n, "Currency") != p.sd["LIMITCURRENCY"] { // F=4
					continue
				}
				p.sd["NETPREMICURRENCY"] = teksDok(n, "Currency")
				p.sd["NETPREMIVALUE"] = teksDok(n, "Value")
			}
		}
		for _, r := range barisDok(sp["RnmLimitList"]) { // [5.4.3.3]
			p.sd["LIMITCURRENCY"] = teksDok(r, "Currency")
			for _, gp := range barisDok(sp["GrossPremiumList"]) { // [5.4.3.3.2]
				if teksDok(gp, "Currency") != p.sd["LIMITCURRENCY"] { // F=4
					continue
				}
				p.sd["MDPCURRENCY"] = teksDok(gp, "Currency")
				p.sd["MDPVALUE"] = teksDok(gp, "Value")
				for _, n := range barisDok(sp["NetPremiumList"]) { // [5.4.3.3.2.2]
					if teksDok(n, "Currency") != p.sd["LIMITCURRENCY"] { // F=4
						continue
					}
					// ⚠️ Langkah ini hanya me-nol-kan deduksi — Net Premium
					// tetap dari [5.4.3.2.2].
					p.sd["DEDUCTION1"] = "0"
					p.sd["DEDUCTION2"] = "0"
					p.deduksiDanSisip(sh)
				}
			}
		}
	}
}

// cabangNonProporsionalEDM - EDM [7] ("When Non Proportional (No COB)").
func (p *penyusunDetail) cabangNonProporsionalEDM() {
	for _, sh := range barisDok(p.doc["Share"]) {
		p.kepalaShare(sh) // [7.1]
		p.sd["SPREADINGTYPE"] = teksDok(sh, "SpreadingTypeXOL")
		p.sd["SPREADINGTYPEID"] = teksDok(sh, "SpreadingTypeIDXOL")
		p.sd["BROKERAGE"] = teksDok(p.doc, "BrokeragePercent")
		for _, g := range barisDok(sh["TreatyGroupList"]) { // [7.2]
			p.sd["TREATYGROUP"] = teksDok(g, "TreatyGroup")
			p.sd["TREATYGROUPID"] = teksDok(g, "TreatyGroupID")
			for _, r := range barisDok(sh["RnmLimitList"]) { // [7.2.2]
				p.batasRNM(r) // [7.2.2.1] RNM Limit
				p.premiShare(sh)
			}
		}
	}
}

// kepalaShare - [5.1] / [7.1]: medan layer dan `RNM_SHARE` akar.
func (p *penyusunDetail) kepalaShare(sh map[string]any) {
	p.sd["LAYER"] = teksDok(sh, "Layer")
	p.sd["LAYERPART"] = teksDok(sh, "LayerPart")
	p.sd["LAYERPARTTYPE"] = teksDok(sh, "LayerPartType")
	p.sd["LAYERTYPE"] = teksDok(sh, "LayerType")
	p.sd["RNM_SHARE"] = teksDok(p.doc, "RNMShare")
}

// batasRNM - langkah `RNM Limit` ([5.4.3.2.1] / [7.2.2.1]).
func (p *penyusunDetail) batasRNM(r map[string]any) {
	p.sd["LIMITCURRENCY"] = teksDok(r, "Currency")
	p.sd["LIMITVALUE"] = teksDok(r, "Value")
	p.sd["SHARECURRENCY"] = teksDok(r, "Currency")
	p.sd["SHAREVALUE"] = teksDok(r, "Value")
}

// premiShare - Gross Premium → Net Premium → Deduction → sisip
// ([5.4.2.1.1] / [7.2.2.2]), semuanya larik baris `Share` yang sama.
func (p *penyusunDetail) premiShare(sh map[string]any) {
	for _, gp := range barisDok(sh["GrossPremiumList"]) {
		if teksDok(gp, "Currency") != p.sd["LIMITCURRENCY"] { // Gross Premium, F=4
			continue
		}
		p.sd["MDPCURRENCY"] = teksDok(gp, "Currency")
		p.sd["MDPVALUE"] = teksDok(gp, "Value")
		// NetPremiumList — `pyStepsPreCondition = false`: WHEN-nya tidak dipakai.
		for _, n := range barisDok(sh["NetPremiumList"]) {
			if teksDok(n, "Currency") != p.sd["LIMITCURRENCY"] { // Net Premium, F=4
				continue
			}
			p.sd["NETPREMICURRENCY"] = teksDok(n, "Currency")
			p.sd["NETPREMIVALUE"] = teksDok(n, "Value")
			p.sd["DEDUCTION1"] = "0"
			p.sd["DEDUCTION2"] = "0"
			p.deduksiDanSisip(sh)
		}
	}
}

// deduksiDanSisip - `DeductionList` (Deduction, Installment, sisip).
//
// `brokerage` (huruf kecil, tanpa trim — `@toLowerCase(.Comment)`) ke
// `DEDUCTION1`; selainnya ke `DEDUCTION2` (T=3: dilewati bila brokerage).
// Hanya baris bermata uang `MDPCURRENCY` (F=3).
func (p *penyusunDetail) deduksiDanSisip(sh map[string]any) {
	for _, d := range barisDok(sh["DeductionList"]) {
		if teksDok(d, "Currency") != p.sd["MDPCURRENCY"] {
			continue
		}
		if strings.ToLower(teksDok(d, "Comment")) == "brokerage" {
			p.sd["DEDUCTION1"] = p.jumlah(p.sd["DEDUCTION1"], teksDok(d, "Deduction"), "DeductionList.Deduction")
		} else {
			p.sd["DEDUCTION2"] = p.jumlah(p.sd["DEDUCTION2"], teksDok(d, "Deduction"), "DeductionList.Deduction")
		}
	}
	// Installment — `SaveData.INSTALLMENT` bukan parameter prosedur.
	p.sd["INSTALLMENTNO"] = teksDok(p.doc, "InstallmentNo")
	p.sisip()
}

// persenRI - `@if(@toDecimal(x)>100, @divide(@toDecimal(x),10,2),
// @toDecimal(x))`. Kosong = 0 (`@toDecimal`).
func (p *penyusunDetail) persenRI(x, medan string) string {
	d, ok := p.desimal(x, medan)
	if !ok {
		return ""
	}
	if d.Cmp(apd.New(100, 0)) <= 0 {
		return d.Text('f')
	}
	var hasil apd.Decimal
	k := apd.BaseContext.WithPrecision(40)
	k.Rounding = apd.RoundHalfUp
	if _, err := k.Quo(&hasil, d, apd.New(10, 0)); err != nil {
		p.catat(fmt.Errorf("%s %q: %w", medan, x, err))
		return ""
	}
	if _, err := k.Quantize(&hasil, &hasil, -2); err != nil {
		p.catat(fmt.Errorf("%s %q: %w", medan, x, err))
		return ""
	}
	return hasil.Text('f')
}

// jumlah - `a + b` Property-Set; kosong = 0.
func (p *penyusunDetail) jumlah(a, b, medan string) string {
	x, ok := p.desimal(a, medan)
	if !ok {
		return a
	}
	y, ok := p.desimal(b, medan)
	if !ok {
		return a
	}
	var hasil apd.Decimal
	if _, err := apd.BaseContext.WithPrecision(40).Add(&hasil, x, y); err != nil {
		p.catat(fmt.Errorf("%s: %w", medan, err))
		return a
	}
	return hasil.Text('f')
}

// desimal - teks → desimal untuk aritmetika; kosong = 0. Yang tidak terbaca
// DICATAT.
func (p *penyusunDetail) desimal(s, medan string) (*apd.Decimal, bool) {
	if strings.TrimSpace(s) == "" {
		return apd.New(0, 0), true
	}
	d, err := angkaDetail(s)
	if err != nil {
		p.catat(fmt.Errorf("%s %q bukan angka", medan, s))
		return nil, false
	}
	return d, true
}

// angkaDetail - teks → desimal: bentuk kabel lebih dulu, lalu bentuk
// ketikan Indonesia (`kabelDariKetikan`) — aturan yang SAMA dengan `angka`
// rumus Limits. Tabel pendaratan memuat keduanya (kontrak uji DEV
// menyimpan `RIOGR` "22,5"); keduanya diubah menjadi NUMBER yang sama.
func angkaDetail(s string) (*apd.Decimal, error) {
	s = strings.TrimSpace(s)
	d, _, err := apd.NewFromString(s)
	if err != nil {
		d, _, err = apd.NewFromString(kabelDariKetikan(s))
	}
	return d, err
}

func (p *penyusunDetail) catat(err error) {
	if p.galat == nil {
		p.galat = err
	}
}

// selesai - `SaveData` tersisip → baris siap tulis menurut kolom tabelnya.
//
// ⛔ Angka yang tidak terbaca DITOLAK dengan nama kolom dan nomor barisnya:
// prosedur Pega menelan galat konversi (`StsSave := 0`, baris hilang tanpa
// jejak); kontrak tuntas yang detailnya diam-diam bolong lebih buruk daripada
// tombol yang menolak.
func (p *penyusunDetail) selesai(kolom []models.KolomDetail) ([]models.BarisDetailTreaty, error) {
	if p.galat != nil {
		return nil, ditolak("Detail treaty tidak dapat disusun: " + p.galat.Error() + ".")
	}
	out := make([]models.BarisDetailTreaty, 0, len(p.baris))
	for i, sd := range p.baris {
		b := models.BarisDetailTreaty{Teks: map[string]string{}, Angka: map[string]*apd.Decimal{}}
		for _, k := range kolom {
			v := strings.TrimSpace(sd[k.Nama])
			if !k.Angka {
				b.Teks[k.Nama] = sd[k.Nama]
				continue
			}
			if v == "" {
				b.Angka[k.Nama] = nil
				continue
			}
			d, err := angkaDetail(v)
			if err != nil {
				return nil, ditolak(fmt.Sprintf("Detail treaty baris ke-%d: kolom %s %q bukan angka.", i+1, k.Nama, sd[k.Nama]))
			}
			b.Angka[k.Nama] = d
		}
		out = append(out, b)
	}
	return out, nil
}

// barisDok - larik dokumen (`[]any` dari kabel/pendaratan, atau
// `[]map[string]any` dari rumus) sebagai daftar baris.
func barisDok(v any) []map[string]any {
	switch x := v.(type) {
	case []map[string]any:
		return x
	case []any:
		out := make([]map[string]any, 0, len(x))
		for _, e := range x {
			if m, ok := e.(map[string]any); ok {
				out = append(out, m)
			}
		}
		return out
	}
	return nil
}

// elemenKe - `Larik(n)` Pega (indeks dari 0 di sini); di luar jangkauan =
// halaman kosong, semua properti terbaca kosong.
func elemenKe(xs []map[string]any, i int) map[string]any {
	if i >= 0 && i < len(xs) {
		return xs[i]
	}
	return map[string]any{}
}
