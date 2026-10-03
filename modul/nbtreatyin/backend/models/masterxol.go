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
// komentar tiap medan) - SATU ukuran K8 di seluruh modul, `[menunggu konfirmasi
// WO]` (tiket 01 bab P9, PERMINTAAN-TIM-INTI F1): lebih luas dari daftar harfiah
// K8, yang memuat keluaran (`TreatyXOLList`) dan medan polis (`FlagPPH`,
// `TypeTax`). Medan dokumen lain tidak pernah sampai ke halaman.
// Nol penulisan: halaman `TreatyIn` tidak disimpan (katalog hanya memuat
// `TreatyIn.ID`), dan tidak satu pun berkas menulis JSON.
//
// Bentuknya sama dengan halaman kerja: kunci RELATIF `TreatyIn` - skalar di
// `Nilai` ("RNMShare"), PageList di `Daftar` ("Share"), daftar bersarang lewat
// `JalurAnak` ("Share(1).GrossPremiumList").

import "strings"

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

// mataNilai - anggota `(Currency, Value)` daftar total per mata uang.
var mataNilai = []string{"Currency", "Value"}

// ringkasLapisan - anggota baris ringkasan layer (`LimitShareSummaryList`,
// `LimitFacShareSummaryList`): grid section `DetailPolicyTreatyInNonProportional`
// "Share" / "Share Facultative" dan `InputPolicyTreatyInDetail_preACT` 18.1.
var ringkasLapisan = []string{"Note", "Limit", "Limit2", "MDP", "MDP2", "Deductible", "Deductible2", "NetPremi", "NetPremi2"}

// SkalarMasterXOL - medan skalar master yang dibaca.
var SkalarMasterXOL = []string{
	// TreatyNonPropSetSpreading 4.1 (`@divide(.Pct,TreatyIn.RNMShare,20)`),
	// CountNetPremi_act 3, CountResult1_Act 4; section "% RNM Share".
	"RNMShare",
	// section DetailPolicyTreatyInNonProportional "% RNM Share" (FacultativeShare != 0).
	"RnmShareDeducted",
	// NonProp 14, 21, 22; TreatyNonPropSetSpreading 7; section "Share Facultative".
	"FacultativeShare",
	// When `TreatyMasterInEDM` (NonProp 10-11; DetailPoliciesNonProportional).
	"EDMState",
	// SetTreatyIn_Act 13 (`TreatyIn.ProportionType=="NonProportional"`).
	"ProportionType",
}

// DaftarMasterXOL - PageList master yang dibaca, beserta anggotanya.
var DaftarMasterXOL = map[string]SkemaDaftarMaster{
	// NonProp 16; InsertToTreatyXOLList 3.2, 3.7; RetroShare 2.2, 2.3.3, 2.3.7;
	// TreatyNonPropSetSpreading 2-4, 7.
	"Share": {
		Anggota: []string{"LayerType", "Layer", "LayerPartType", "LayerPart", "SpreadingTypeXOL", "SpreadingTypeIDXOL"},
		Anak: map[string][]string{
			"GrossPremiumList":   mataNilai,
			"NetPremiumList":     mataNilai,
			"DeductionTotalList": mataNilai,
			"RnmLimitList":       mataNilai,
			"DeductionList":      {"Currency", "Deduction"},
			"SpreadingListXOL":   {"ReinsTypeID", "ReinsTypeName", "Pct"},
		},
	},
	// NonProp 20; InsertToTreatyXOLList 3; RetroShare 2.3.
	"Installment": {
		Anggota: []string{"Currency", "AmountTotal", "PctTotal"},
		Anak: map[string][]string{
			"InstallmentList": {"Installment", "PaymentDate", "InstallmentPct", "Amount", "Currency"},
		},
	},
	// NonProp 17; RetroShare 2.3.3.3, 2.3.7.4.
	"FacultativeShareList": {
		Anak: map[string][]string{
			"GrossPremiumList":   mataNilai,
			"DeductionTotalList": mataNilai,
			"DeductionList":      {"Currency", "Deduction"},
		},
	},
	// TreatySetReinstatement 1.1 / SetReinstatementPct (dipanggil SetTreatyIn_Act 13).
	"Limits": {
		Anggota: []string{"ReinstatementValue", "ReinstatementNote", "Limit", "Limit2"},
		Anak: map[string][]string{
			"MDPList":            mataNilai,
			"Reinstatement_List": {"ReinstatementValue"},
		},
	},
	// preACT 18 (FlagPPH); NonProp 7-8 (FlagRetroTreaty); section grid "Share".
	"LimitShareSummaryList":    {Anggota: ringkasLapisan},
	"LimitFacShareSummaryList": {Anggota: ringkasLapisan},
	// section grid "Limits".
	"LimitSummaryList": {Anggota: []string{"Note", "Limit", "Limit2", "Deductible", "Deductible2", "MDP", "MDP2"}},
	// section: total per mata uang. TotalShareNetNP / TotalShareDeductionNP juga
	// dibaca preACT 18.2-18.3; TotalShareNetNP NonProp 9; TotalFacShareDeductionNP NonProp 7.
	"TotalLimitIOONP":          {Anggota: mataNilai},
	"TotalLimitDeductblNP":     {Anggota: mataNilai},
	"TotalLimitMDPNP":          {Anggota: mataNilai},
	"TotalShareRnmNP":          {Anggota: mataNilai},
	"TotalShareGrossNP":        {Anggota: mataNilai},
	"TotalShareDeductionNP":    {Anggota: mataNilai},
	"TotalShareNetNP":          {Anggota: mataNilai},
	"TotalSpreadedNetPremi":    {Anggota: mataNilai},
	"TotalSpreadedNetPremiRI":  {Anggota: mataNilai},
	"TotalFacShareRnmNP":       {Anggota: mataNilai},
	"TotalFacShareGrossNP":     {Anggota: mataNilai},
	"TotalFacShareDeductionNP": {Anggota: mataNilai},
	"TotalFacShareNetNP":       {Anggota: mataNilai},
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
