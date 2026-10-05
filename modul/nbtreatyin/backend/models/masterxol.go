package models

// Untuk apa berkas ini: DATA MASTER XOL - bentuk halaman `pyWorkPage.TreatyIn`
// yang dibaca jalur NB NonProporsional (XOL), dan SATU daftar medan yang boleh
// dibaca dari dokumen master.
//
// ⭐ `[keputusan work owner]` K8 (PROMPT-NB-TREATY-IN-PUTARAN-2 bab 2, PESAN-KOREKSI
// bagian D, 03-10-2026): data master jalur XOL TIDAK ada di view
// `TREATYINDETAILJOINEDM`, dan tabel master relasional modul `treatyin`
// (`KONTRAK`, `LAYER`, `BAGIAN`, `PEMULIHAN_LIMIT`, `TERMIN`, `POTONGAN`) masih
// nol baris. Maka - pengecualian SEMPIT atas P29 - medan master di bawah ini
// dibaca BACA-SAJA dari `JSONDATA` `M_TREATY_IN` / `M_TREATY_IN_EDM` oleh SATU
// fungsi repository (`repository.MasterXOLDariJSON`), lewat antarmuka
// `services.PembacaMasterTreaty` yang kelak diganti kontrak modul `treatyin`.
// `[penyimpangan sadar]` atas P29 - dicatat di tiket 01 dan PERMINTAAN-TIM-INTI E.
//
// ⛔ HANYA medan yang dibaca rule terjangkau jalur NonProp (nomor langkah di
// komentar tiap medan) - SATU ukuran K8 di seluruh modul, `[keputusan work
// owner]` F1 (PROMPT putaran 3 bab 2, disetujui 04-10-2026): lebih luas dari
// daftar harfiah K8, yang memuat keluaran (`TreatyXOLList`) dan medan polis
// (`FlagPPH`, `TypeTax`). Daftarnya TERTUTUP (`MedanMasterXOL`, dikunci uji
// `TestDaftarMedanMasterXOLTertutup`); medan dokumen lain tidak pernah sampai
// ke halaman (`repository.TestUraiMasterXOLHanyaMedanDaftarTertutup`).
// Nol penulisan: halaman `TreatyIn` tidak disimpan (katalog hanya memuat
// `TreatyIn.ID`), dan tidak satu pun berkas menulis JSON.
//
// Bentuknya sama dengan halaman kerja: kunci RELATIF `TreatyIn` - skalar di
// `Nilai` ("RNMShare"), PageList di `Daftar` ("Share"), daftar bersarang lewat
// `JalurAnak` ("Share(1).GrossPremiumList").

import (
	"sort"
	"strings"
)

// MasterXOL adalah isi master satu kontrak treaty untuk jalur XOL.
type MasterXOL struct {
	Nilai  map[string]string  `json:"nilai"`
	Daftar map[string][]Baris `json:"daftar"`
}

// SkemaDaftarMaster - satu PageList master: anggota skalar dan daftar bersarang
// (nama -> anggota skalarnya).
type SkemaDaftarMaster struct {
	Anggota []string
	Anak    map[string][]string
}

// Singkatan kutipan: "NonProp" = Activity `InputPolicyTreatyInDetail_NonProp`,
// "preACT" = `InputPolicyTreatyInDetail_preACT`, "Section" = Section
// `DetailPolicyTreatyInNonProportional` (grid `pyPageListProperty =
// pyWorkPage.TreatyIn.<daftar>`). Diperiksa ulang ke korpus XML 04-10-2026 (F1).

// mataNilai - anggota `(Currency, Value)` daftar total per mata uang: kolom
// grid total Section (`.Currency`, `.Value`).
var mataNilai = []string{"Currency", "Value"}

// ringkasLapisan - anggota baris ringkasan layer, kolom grid Section "Share" /
// "Share Facultative": Note, MDP, MDP2 dibaca Section saja; Limit, Limit2,
// Deductible, Deductible2, NetPremi, NetPremi2 juga preACT 18.1
// (`@if(.Limit>0,"IDR","")`, `@divide(.Deductible,...)`, `.NetPremi + .PPNValue`).
var ringkasLapisan = []string{"Note", "Limit", "Limit2", "MDP", "MDP2", "Deductible", "Deductible2", "NetPremi", "NetPremi2"}

// SkalarMasterXOL - medan skalar master yang dibaca.
var SkalarMasterXOL = []string{
	// TreatyNonPropSetSpreading 4.1 (`@divide(.Pct,pyWorkPage.TreatyIn.RNMShare,20)`),
	// CountNetPremi_act 3, CountResult1_Act 4 (bersyarat NonProportional);
	// sel Section "% RNM Share" (tampil bila FacultativeShare = 0).
	"RNMShare",
	// sel Section "% RNM Share" (tampil bila FacultativeShare != 0).
	"RnmShareDeducted",
	// NonProp 14, 21, 22; TreatyNonPropSetSpreading 7; Section SpreadingRiskList
	// (pyCondition / pyReadOnlyCondition `pyWorkPage.TreatyIn.FacultativeShare >0`);
	// wadah Section "Share Facultative" (`FacultativeShare>0`).
	"FacultativeShare",
	// When `TreatyMasterInEDM` (`EDMState` = "1"|"2"|"3"; NonProp 10-11; wadah
	// Section DetailPoliciesNonProportional).
	"EDMState",
	// SetTreatyIn_Act 13 (`TreatyIn.ProportionType=="NonProportional"`).
	"ProportionType",
}

// DaftarMasterXOL - PageList master yang dibaca, beserta anggotanya.
var DaftarMasterXOL = map[string]SkemaDaftarMaster{
	"Share": {
		// LayerType, Layer, LayerPartType, LayerPart: InsertToTreatyXOLList 3.7.2,
		// RetroShare 2.3.7.2. SpreadingTypeXOL, SpreadingTypeIDXOL:
		// TreatyNonPropSetSpreading 3, 4, 7; RetroShare 2.2.
		Anggota: []string{"LayerType", "Layer", "LayerPartType", "LayerPart", "SpreadingTypeXOL", "SpreadingTypeIDXOL"},
		Anak: map[string][]string{
			// GrossPremiumList, NetPremiumList, DeductionTotalList, RnmLimitList,
			// DeductionList: NonProp 16.1-16.4; InsertToTreatyXOLList 3.2.3.x.
			"GrossPremiumList":   mataNilai,
			"NetPremiumList":     mataNilai,
			"DeductionTotalList": mataNilai,
			"RnmLimitList":       mataNilai,
			"DeductionList":      {"Currency", "Deduction"},
			// TreatyNonPropSetSpreading 2 (ReinsTypeID, ReinsTypeName), 4.1 (.Pct);
			// NonProp 9 (`Share(1).SpreadingListXOL(1).Pct`, `(2).Pct`).
			"SpreadingListXOL": {"ReinsTypeID", "ReinsTypeName", "Pct"},
		},
	},
	"Installment": {
		// NonProp 20.1, 20.3; InsertToTreatyXOLList 3.1, 3.3; RetroShare 2.3.1.
		Anggota: []string{"Currency", "AmountTotal", "PctTotal"},
		Anak: map[string][]string{
			// NonProp 20.4.1 (.Installment, .PaymentDate, .InstallmentPct,
			// .Amount, .Currency).
			"InstallmentList": {"Installment", "PaymentDate", "InstallmentPct", "Amount", "Currency"},
		},
	},
	"FacultativeShareList": {
		Anak: map[string][]string{
			// RetroShare 2.3.7.4.
			"GrossPremiumList": mataNilai,
			// NonProp 17.1.
			"DeductionTotalList": mataNilai,
			// RetroShare 2.3.3.3, 2.3.7.4.4.
			"DeductionList": {"Currency", "Deduction"},
		},
	},
	"Limits": {
		// TreatySetReinstatement 1 (kalang `TreatyIn.Limits`), dipanggil
		// SetTreatyIn_Act 13; SetReinstatementPct 2 (.ReinstatementValue),
		// 3.1 (.ReinstatementNote, .Limit, .Limit2).
		Anggota: []string{"ReinstatementValue", "ReinstatementNote", "Limit", "Limit2"},
		Anak: map[string][]string{
			// SetReinstatementPct 3.2-3.3 (`.MDPList(n).Currency/Value`).
			"MDPList": mataNilai,
			// TreatySetReinstatement 1.1 (`.Reinstatement_List(1).ReinstatementValue == ""`).
			"Reinstatement_List": {"ReinstatementValue"},
		},
	},
	// ringkasLapisan; NonProp 8.1 (FlagRetroTreaty); preACT 18.1; grid Section
	// "Share". NetPremiAfterPPN/PPH(2): sel grid "Share" - ditimpa preACT 18.1
	// HANYA bila FlagPPH == "true"; selain itu Section menampilkan nilai dokumen
	// master (`adoptJSONObject` NonProp 6).
	"LimitShareSummaryList": {Anggota: append(append([]string{}, ringkasLapisan...),
		"NetPremiAfterPPN", "NetPremiAfterPPH", "NetPremiAfterPPN2", "NetPremiAfterPPH2")},
	// ringkasLapisan; NonProp 7 (disalin ke LimitShareSummaryList bila
	// FlagRetroTreaty); grid Section "Share Facultative".
	"LimitFacShareSummaryList": {Anggota: ringkasLapisan},
	// grid Section "Limits".
	"LimitSummaryList": {Anggota: []string{"Note", "Limit", "Limit2", "Deductible", "Deductible2", "MDP", "MDP2"}},
	// Total per mata uang - grid total Section (mataNilai).
	"TotalLimitIOONP":      {Anggota: mataNilai},
	"TotalLimitDeductblNP": {Anggota: mataNilai},
	"TotalLimitMDPNP":      {Anggota: mataNilai},
	"TotalShareRnmNP":      {Anggota: mataNilai},
	"TotalShareGrossNP":    {Anggota: mataNilai},
	// + preACT 18.3 (`.Currency`); TotalPPNValue/TotalPPHValue: sel grid
	// "Total Brokerage", ditimpa preACT 18.3.1-18.3.2 hanya bila FlagPPH == "true".
	"TotalShareDeductionNP": {Anggota: []string{"Currency", "Value", "TotalPPNValue", "TotalPPHValue"}},
	// + preACT 18.2 (`.Currency`); TotalNetPremiAfterPPN/Tax: sel grid "Total Net
	// Premi", ditimpa preACT 18.2 hanya bila FlagPPH == "true". (NonProp 9 membaca
	// SALINAN dari TotalFacShareDeductionNP yang ditimpa NonProp 7.)
	"TotalShareNetNP": {Anggota: []string{"Currency", "Value", "TotalNetPremiAfterPPN", "TotalNetPremiAfterTax"}},
	// grid Section "Total Spreaded" (NonProp 9 MENULIS `(1).Value`, bukan membaca).
	"TotalSpreadedNetPremi":   {Anggota: mataNilai},
	"TotalSpreadedNetPremiRI": {Anggota: mataNilai},
	"TotalFacShareRnmNP":      {Anggota: mataNilai},
	"TotalFacShareGrossNP":    {Anggota: mataNilai},
	// + NonProp 7 (disalin ke TotalShareNetNP bila FlagRetroTreaty).
	"TotalFacShareDeductionNP": {Anggota: mataNilai},
	"TotalFacShareNetNP":       {Anggota: mataNilai},
}

// MedanMasterXOL - DAFTAR TERTUTUP medan master yang boleh dibaca dari dokumen
// JSON (F1): setiap skalar `SkalarMasterXOL`, setiap anggota PageList
// `DaftarMasterXOL` ("Share.LayerType"), dan setiap anggota daftar bersarang
// ("Share.GrossPremiumList.Currency"), berurut. Uji
// `TestDaftarMedanMasterXOLTertutup` mengunci isinya, dan
// `repository.TestUraiMasterXOLHanyaMedanDaftarTertutup` gagal bila pembaca
// JSON mengembalikan medan di luar daftar ini.
func MedanMasterXOL() []string {
	out := append([]string{}, SkalarMasterXOL...)
	for nama, sk := range DaftarMasterXOL {
		for _, a := range sk.Anggota {
			out = append(out, nama+"."+a)
		}
		for anak, anggota := range sk.Anak {
			for _, a := range anggota {
				out = append(out, nama+"."+anak+"."+a)
			}
		}
	}
	sort.Strings(out)
	return out
}

// TanggalMasterXOL - anggota master bertipe tanggal (masuk halaman sebagai
// "2006-01-02"): `Installment().InstallmentList().PaymentDate` -> NonProp 20.4.1
// `.DueDate` -> kolom DATE `T_POLIS_INSTALMENT_DETAIL.DUE_DATE`.
var TanggalMasterXOL = map[string]bool{"Installment.InstallmentList.PaymentDate": true}

const jMaster = HalamanMaster + "."

// TerapkanMasterXOL menaruh master di halaman `TreatyIn` = NonProp langkah 6 /
// TreatyRealizationCheckXOLList langkah 5 (`adoptJSONObject` / `Page-Copy` ke
// `pyWorkPage.TreatyIn`).
//
// ⚠️ `[tafsiran]` Medan XOL master sebelumnya DIHAPUS lebih dulu (daftar
// diganti, bukan digabung) supaya pilih bisnis kedua tidak menyisakan baris
// kontrak lama. `TreatyIn.ID`, `Commencement`, `Termination` tetap - ketiganya
// dari view (preACT 6, `TerapkanMasterKontrak`).
func TerapkanMasterXOL(h *Halaman, m MasterXOL) {
	h.pastikan()
	for _, s := range SkalarMasterXOL {
		h.Hapus(jMaster + s)
	}
	for k := range h.Daftar {
		if strings.HasPrefix(k, jMaster) {
			delete(h.Daftar, k)
		}
	}
	for k, v := range m.Nilai {
		h.Setel(jMaster+k, v)
	}
	for k, v := range m.Daftar {
		h.SetelDaftar(jMaster+k, salinBaris(v))
	}
}

// TreatyMasterInEDM = `When/TreatyMasterInEDM`: `pyWorkPage.TreatyIn.EDMState`
// = "1" ATAU "2" ATAU "3".
func TreatyMasterInEDM(h *Halaman) bool {
	switch h.Ambil(jMaster + "EDMState") {
	case "1", "2", "3":
		return true
	}
	return false
}

// mDaftar membaca PageList master `TreatyIn.<nama>`.
func mDaftar(h *Halaman, nama string) []Baris { return h.AmbilDaftar(jMaster + nama) }

// mAnak membaca daftar bersarang `TreatyIn.<daftar>(<i>).<anak>` (i berbasis 1);
// baris induk yang tidak ada = daftar kosong.
func mAnak(h *Halaman, daftar string, i int, anak string) []Baris {
	return h.AmbilDaftar(JalurAnak(jMaster+daftar, i, anak))
}

// MataUangAngsuranMaster - nama mata uang yang ID-nya dicari RDB
// `GetDataCurrencyByName_SQL` (NonProp 20.1-20.2, InsertToTreatyXOLList 3.3-3.4,
// RetroShare 2.3.4-2.3.5): `TreatyIn.Installment().Currency`, urut tanpa kembar.
func MataUangAngsuranMaster(h *Halaman) []string {
	var out []string
	ada := map[string]bool{}
	for _, b := range mDaftar(h, "Installment") {
		if c := b["Currency"]; !ada[c] {
			ada[c] = true
			out = append(out, c)
		}
	}
	return out
}

// IDMataUang - hasil `GetDataCurrencyByName_SQL` per nama mata uang
// (`CurrencySearch.pxResults(1).ID`; nama tak dikenal = "").
type IDMataUang map[string]string
