package models

// Tab **Share** cabang NON-PROPORSIONAL — bentuknya dari
// `Section/TreatyInTabsNonProportional.xml` (@1695720–@3291797) dan
// `Section/Share.xml` (panel rincian baris, `pyEditAction = Share`).
//
// ⭐ Ejaan JSON = ejaan dokumen Pega (`RNMShare`, `SpreadingListXOL`, …),
// seperti `LayerNP`: layar dan rumus memegang nama properti yang sama
// dengan ekspor, jadi satu medan tidak punya dua nama.
//
// ⛔ Sumbernya tabel pendaratan SAJA (nol JSON, nol `M_TREATY_IN`):
//
//	Share                      T_TREATY_SHARE (+ _AMOUNT Gross/Net, _SPREADING, _DEDUCTION)
//	FacultativeShareList       T_TREATY_FAC_SHARE (+ _AMOUNT, _DEDUCTION)
//	ShareReins                 T_TREATY_RETRO_SHARE
//	ShareFacultativeReinsurers T_TREATY_FAC_REINSURER
//	FacultativeShare(+Brokerage) T_TREATY_REVISION
//
// ⚠️ Skalar akar `RNMShare`, `BrokeragePercent`, `RNMShareAcrossTheBoard`,
// `RnmShareDeducted` TIDAK punya kolom pendaratan. Dua yang pertama dibaca
// dari SALINAN yang Activity Pega tulis sendiri (lihat
// `services/share_np_muat.go`); dua sisanya bawaan / turunan.

// NilaiMataUang - satu baris larik bernilai per mata uang (kelas
// `ASM-FW-GISFW-Data-TreatyInTotal`): `RnmLimitList`, `GrossPremiumList`,
// `TotalShareRnmNP`, …
type NilaiMataUang struct {
	Currency   string `json:"Currency"`
	CurrencyID string `json:"CurrencyID"`
	Value      string `json:"Value"`
}

// GrupShareNP - satu baris `TreatyGroupList` baris Share (salinan
// `Limits(n).TreatyGroupList`, `TreatyInNonAddItem` langkah 13.1).
type GrupShareNP struct {
	TreatyGroup   string `json:"TreatyGroup"`
	TreatyGroupID string `json:"TreatyGroupID"`
}

// BarisDeduksiShare - satu baris grid `Deduction Details` panel Share.
type BarisDeduksiShare struct {
	Comment    string `json:"Comment"`
	Currency   string `json:"Currency"`
	CurrencyID string `json:"CurrencyID"`
	Deduction  string `json:"Deduction"`
	// `DeductionPct` dan bendera `Auto Calculate %` — `CalculateDeduction`
	// langkah 4 hanya berlanjut ke mata uang berikutnya bila bendera `true`.
	DeductionPct          string `json:"DeductionPct"`
	DeductionPctCalculate string `json:"DeductionPctCalculate,omitempty"`
}

// BarisSpreadingNP - satu baris `SpreadingListXOL`.
//
// Dua asal: `FetchQSfromMasterXOL` menyalin anak susunan treaty master
// (`QS (OR)` 40 / `QS (R/I)` 60, …); jalur manual (`AddSpreadingXOL` +
// `SetSpreadingXOL`) membiarkan pemakai mengisi `ReinsTypeID` dan `Pct`, lalu
// mengisi kelima larik di bawah sebagai bagian `Pct ÷ RNMShare`.
type BarisSpreadingNP struct {
	ReinsTypeName     string `json:"ReinsTypeName"`
	ReinsTypeID       string `json:"ReinsTypeID"`
	ParentReinsTypeID string `json:"ParentReinsTypeID"`
	Pct               string `json:"Pct"`
	Rp                string `json:"Rp"`
	Usd               string `json:"Usd"`

	RnmLimitList        []NilaiMataUang `json:"RnmLimitList"`
	GrossPremiumList    []NilaiMataUang `json:"GrossPremiumList"`
	GrossPremiumMinList []NilaiMataUang `json:"GrossPremiumMinList"`
	DeductionTotalList  []NilaiMataUang `json:"DeductionTotalList"`
	NetPremiumList      []NilaiMataUang `json:"NetPremiumList"`
}

// BarisShareNP - satu baris `TreatyIn.Share` (dan `FacultativeShareList`,
// kelas yang sama) — SATU per layer Limits, urutan layer.
type BarisShareNP struct {
	LayerType       string        `json:"LayerType"`
	Layer           string        `json:"Layer"`
	LayerPartType   string        `json:"LayerPartType"`
	LayerPart       string        `json:"LayerPart"`
	Cover           string        `json:"Cover"`
	RNMShare        string        `json:"RNMShare"`
	TreatyGroupList []GrupShareNP `json:"TreatyGroupList"`

	// Hanya `FacultativeShareList` — limit layer apa adanya (langkah 11.1).
	Limit  string `json:"Limit"`
	Limit2 string `json:"Limit2"`

	SpreadingTypeXOL     string             `json:"SpreadingTypeXOL"`
	SpreadingTypeIDXOL   string             `json:"SpreadingTypeIDXOL"`
	SpreadingTotalPctXOL string             `json:"SpreadingTotalPctXOL"`
	SpreadingListXOL     []BarisSpreadingNP `json:"SpreadingListXOL"`

	DeductionList []BarisDeduksiShare `json:"DeductionList"`

	RnmLimitList        []NilaiMataUang `json:"RnmLimitList"`
	GrossPremiumList    []NilaiMataUang `json:"GrossPremiumList"`
	GrossPremiumMinList []NilaiMataUang `json:"GrossPremiumMinList"`
	DeductionTotalList  []NilaiMataUang `json:"DeductionTotalList"`
	NetPremiumList      []NilaiMataUang `json:"NetPremiumList"`

	// Bagian OR / R/I — `FetchQSfromMasterXOL`, `TreatyInSetBrokerage`.
	RNMSpreadedListXOL           []NilaiMataUang `json:"RNMSpreadedListXOL"`
	RNMSpreadedListRIXOL         []NilaiMataUang `json:"RNMSpreadedListRIXOL"`
	RNMSpreadedListGrossXOL      []NilaiMataUang `json:"RNMSpreadedListGrossXOL"`
	RNMSpreadedListGrossRIXOL    []NilaiMataUang `json:"RNMSpreadedListGrossRIXOL"`
	RNMSpreadedListGrossMinXOL   []NilaiMataUang `json:"RNMSpreadedListGrossMinXOL"`
	RNMSpreadedListGrossRIMinXOL []NilaiMataUang `json:"RNMSpreadedListGrossRIMinXOL"`
	RNMSpreadedListDeductXOL     []NilaiMataUang `json:"RNMSpreadedListDeductXOL"`
	RNMSpreadedListDeductRIXOL   []NilaiMataUang `json:"RNMSpreadedListDeductRIXOL"`
	RNMSpreadedListNetXOL        []NilaiMataUang `json:"RNMSpreadedListNetXOL"`
	RNMSpreadedListNetRIXOL      []NilaiMataUang `json:"RNMSpreadedListNetRIXOL"`
}

// BarisReinsShare - satu baris grid `Reinsurer Name` (`ShareReins`) atau
// `Facultative Reinsurers` (`ShareFacultativeReinsurers`).
type BarisReinsShare struct {
	ID        string `json:"ID"`
	ReinsID   string `json:"ReinsID"`
	ReinsName string `json:"ReinsName"`
	Layer     string `json:"Layer"`
	SharePct  string `json:"SharePct"`
}

// RingkasanShareNP - satu baris `LimitShareSummaryList` (grid "Summarry of
// RNM Share"), `TreatyInSummaryLimitShare`.
type RingkasanShareNP struct {
	LayerType     string `json:"LayerType"`
	Layer         string `json:"Layer"`
	LayerPartType string `json:"LayerPartType"`
	LayerPart     string `json:"LayerPart"`
	Note          string `json:"Note"`
	Limit         string `json:"Limit"`
	Limit2        string `json:"Limit2"`
	MDP           string `json:"MDP"`
	MDP2          string `json:"MDP2"`
	Deductible    string `json:"Deductible"`
	Deductible2   string `json:"Deductible2"`
	NetPremi      string `json:"NetPremi"`
	NetPremi2     string `json:"NetPremi2"`
}

// ShareNP - seluruh isi tab Share Non-Prop.
type ShareNP struct {
	RNMShare                  string `json:"RNMShare"`
	BrokeragePercent          string `json:"BrokeragePercent"`
	RNMShareAcrossTheBoard    string `json:"RNMShareAcrossTheBoard"`
	FacultativeShare          string `json:"FacultativeShare"`
	FacultativeShareBrokerage string `json:"FacultativeShareBrokerage"`
	RnmShareDeducted          string `json:"RnmShareDeducted"`
	IsProRate                 string `json:"IsProRate"`

	ShareReins                 []BarisReinsShare `json:"ShareReins"`
	ShareFacultativeReinsurers []BarisReinsShare `json:"ShareFacultativeReinsurers"`

	Share                []BarisShareNP `json:"Share"`
	FacultativeShareList []BarisShareNP `json:"FacultativeShareList"`

	LimitShareSummaryList    []RingkasanShareNP `json:"LimitShareSummaryList"`
	LimitFacShareSummaryList []RingkasanShareNP `json:"LimitFacShareSummaryList"`

	// Kesembilan grid "Total All Layers RNM Share" — kuncinya nama properti
	// (`TotalShareRnmNP`, `TotalSpreadedRnmProp`, …); selalu kesembilannya.
	Total map[string][]NilaiMataUang `json:"Total"`
}

// SharePendaratan - isi tabel pendaratan tab Share, APA ADANYA (simpul
// berkunci ejaan dokumen); services yang menafsirkannya.
type SharePendaratan struct {
	Share                      []map[string]any
	FacultativeShareList       []map[string]any
	ShareReins                 []map[string]any
	ShareFacultativeReinsurers []map[string]any
}

// SusunanSpreading - satu baris `PROPORTIONALARRG` yang RD spreading
// kembalikan: induk (`BrowseTreatyArrangement_ParentReinsMasterTrt`) atau
// anak (`BrowseTreatyArrangement_Limit_RD`).
type SusunanSpreading struct {
	ReinsTypeID       string `json:"reinsTypeId"`
	ReinsTypeName     string `json:"reinsTypeName"`
	ParentReinsTypeID string `json:"parentReinsTypeId"`
	TreatyYearID      string `json:"treatyYearId"`
	// TreatyYear - `.TreatyYear` baris susunan; `FetchQSfromMaster` (Prop)
	// langkah 6.1 menyalinnya ke `Local.TreatyYear` untuk RD anak.
	TreatyYear string `json:"treatyYear"`
	Pct        string `json:"pct"`
	Rp         string `json:"rp"`
	Usd        string `json:"usd"`
}

// ShareDetailWarisan - `RNM_SHARE` / `BROKERAGE` kontrak dari `TREATYINDETAIL`
// (tabel datar yang `SaveTreatyInDetail_Act` Pega tulis) — cadangan terakhir
// skalar akar Share.
type ShareDetailWarisan struct {
	RNMShare         string
	BrokeragePercent string
}
