// Package models memuat bentuk data modul Master Product Name Life.
//
// Satu produk = halaman Pega `ProductName` (sisi umum, `M_PRODUCT_LIFE`) +
// `ProductNameInward` (sisi inward, `M_PRODUCTINWARD_LIFE`) + tujuh daftar
// bersarang di halaman umum. Bentuk ini struct BERNAMA; kunci JSON Pega
// (`POLICYHODER`, `UnderwritingLimitList`, `Non_Employee`, ...) hanya dikenal
// repository (P1, `docs/RALAT-DEV-30-09-2026.md`).
//
// ⛔ Angka (uang, persen, usia, hari) TEKS sepanjang jalan - tidak pernah
// float (ADR-0003); services memeriksanya sebagai desimal dan menyimpannya
// apa adanya. Tanggal `YYYY-MM-DD` di API; repository menulisnya `dd/MM/yyyy`
// seperti Pega (`BrowseReinstypeOR_SQL` b84 `TO_DATE(…, 'DD/MM/YYYY')`).
//
// Bukti tiap medan: `docs/PARITAS-LAYAR-DAN-AKSI.md` §3–§5.
package models

// Produk adalah satu produk life utuh.
type Produk struct {
	// ID - `M_PRODUCT_LIFE.ID` (`'1' ‖ LPAD(M_PRODUCT_LIFE_SEQ, 5, '0')`),
	// tampil sebagai `Product Code` b3894. Tidak pernah dari klien untuk produk baru.
	ID     string       `json:"id"`
	Umum   ProdukUmum   `json:"umum"`
	Inward ProdukInward `json:"inward"`

	// Daftar bersarang halaman `ProductName` (PARITAS §5).
	LienClause            []BarisLien     `json:"lienClause"`
	DocumentClaim         []BarisDokumen  `json:"documentClaim"`
	PlanList              []BarisPlan     `json:"planList"`
	FinancialUnderwriting []BarisFinUW    `json:"financialUnderwriting"`
	UnderwritingLimit     []BarisUWLimit  `json:"underwritingLimit"`
	OutwardList           []BarisOutward  `json:"outwardList"`
	CommentList           []BarisKomentar `json:"commentList"`

	// Penanda PERMINTAAN simpan (paket 9) - tidak disimpan, tidak pernah dijawab.
	//
	// SalinanDari - `Copy` b59854 (`CopyProduct` b144/b173 hanya mengosongkan kedua
	// ID): produk baru mewarisi medan milik server produk asal ini.
	SalinanDari string `json:"salinanDari,omitempty"`
	// HitungOutward - checkbox `On Retention` b47312 diubah dalam sesi sunting
	// ini: `GetReinsTypeOR_Life` b47488 mengisi ulang `OutwardList` (OQ-MPNL-09).
	HitungOutward bool `json:"hitungOutward,omitempty"`
}

// ProdukUmum - halaman `ProductName` (`M_PRODUCT_LIFE.JSONDATA`).
type ProdukUmum struct {
	ProductName      string `json:"productName"`      // `Product Name` b3620
	Ceding           string `json:"ceding"`           // `Ceding` b4075
	CedingID         string `json:"cedingId"`         // `setCeding_DT` b191
	SOBName          string `json:"sobName"`          // `SOB` b4463
	SOBID            string `json:"sobId"`            // `setSOB_DT` b191
	RIComm           string `json:"riComm"`           // `Deduction (%)` b7104
	RIRisk           string `json:"riRisk"`           // `R/I Risk Name` b7398
	RIRiskID         string `json:"riRiskId"`         // `setRIRISK_DT` b191
	InwardName       string `json:"inwardName"`       // `Treaty Name` b8719
	TreatyNumber     string `json:"treatyNumber"`     // `Treaty Number` b8900
	Cause            string `json:"cause"`            // `Cause Of Loss` b10661
	CauseID          string `json:"causeId"`          // `setCauseOfLoss_DT` b172
	PolicyHolder     string `json:"policyHolder"`     // JSON `POLICYHODER` - salinan inward (`SaveProductName_Act` 1 b361)
	PolicyHolderName string `json:"policyHolderName"` // JSON `POLICYHODERNAME` - salinan inward
	CreateOp         string `json:"createOp"`         // `SaveProductName_Act` 6 b1370
	UpdateOp         string `json:"updateOp"`         // `SaveProductName_Act` 1 b361
	// IsORS - checkbox `On Retention` b47312 (RALAT R9).
	IsORS bool `json:"isOrs"`
	// Comment - area teks `Comment` popup `SaveProductName_Confirm` b1025.
	Comment string `json:"comment"`

	// Medan layar MATI (`OTHER 1=2`, PARITAS §3.1) - kuncinya dibaca view
	// `PRODUCT_LIFE`. Diisi dari JSON lama saja, tidak pernah dari klien.
	// TypeBasicRider - kunci JSON `TYPE` (`GenerateUpload_Act` b332: `1` = Basic, lainnya Rider).
	// Bukan `Type` klaim Claim Life - nama sengaja berbeda (penjaga `satutype_test.go`).
	TypeBasicRider string `json:"type"`
	TypeCeding     string `json:"typeCeding"`
	Grup           string `json:"grup"`
	ProductCode    string `json:"productCode"`
	ProductType    string `json:"productType"`
	ProductTypeID  string `json:"productTypeId"`
	RIRate         string `json:"riRate"`
	RIRateID       string `json:"riRateId"`
	RICommID       string `json:"riCommId"`
	OutwardName    string `json:"outwardName"`
	OutwardNameID  string `json:"outwardNameId"`
	OutwardRate    string `json:"outwardRate"`
	OutwardRateID  string `json:"outwardRateId"`
	OutwardComm    string `json:"outwardComm"`
	OutwardCommID  string `json:"outwardCommId"`
	Benefit        string `json:"benefit"`
	BenefitID      string `json:"benefitId"`
}

// ProdukInward - halaman `ProductNameInward` (`M_PRODUCTINWARD_LIFE.JSONDATA`).
type ProdukInward struct {
	// ID baris inward - sama dengan ID produk untuk produk baru (R14; OQ-MPNL-02 ditutup data DEV).
	ID        string `json:"id"`
	ProductID string `json:"productId"` // `SaveProductName_Act` 13 b2530

	PolicyHolder        string `json:"policyHolder"`        // JSON `POLICYHODER` - `setPolicyHolder_DT` b172
	PolicyHolderName    string `json:"policyHolderName"`    // `Policy Holder` b17097
	Insured             string `json:"insured"`             // `Insured` b18341
	AddendumNo          string `json:"addendumNo"`          // `Addendum No.` b18834
	AddendumWord        string `json:"addendumWord"`        // `Addendum` b19432
	AmandementNo        string `json:"amandementNo"`        // `Amandement No.` b20155
	AmandementSchd      string `json:"amandementSchd"`      // `Amandement` b20760
	MaxExpiredClaim     string `json:"maxExpiredClaim"`     // `Max Notification Claim Expired` b21781
	Begin               string `json:"begin"`               // `Begin Date` b21969
	STNC                string `json:"stnc"`                // `STNC` b22304
	CedingRetentionNum  string `json:"cedingRetentionNum"`  // `Ceding Retention (%)` b22494
	CedingLimit         string `json:"cedingLimit"`         // `Ceding's Limit` b22657
	Brokerage           string `json:"brokerage"`           // `Brokerage Fee (%)` b22915
	MinAge              string `json:"minAge"`              // `Minimum Age (Years)` b23078
	MaxAge              string `json:"maxAge"`              // `Maximum Age (Years)` b23265
	ExpiryAge           string `json:"expiryAge"`           // `Expiry Age (Years)` b23452
	ExtraPremi          string `json:"extraPremi"`          // `Extra Premium` b23635
	MinSumInsured       string `json:"minSumInsured"`       // `Min Sum Insured` b23796
	MaxSumInsured       string `json:"maxSumInsured"`       // `Max Sum Insured` b23956
	MaxSumReasured      string `json:"maxSumReasured"`      // `Max Sum Reasured` b24214
	RNMShare            string `json:"rnmShare"`            // `Nusantara Re Share (%)` b24674
	RNMLimitNum         string `json:"rnmLimitNum"`         // `Nusantara Re's Limit` b25236
	PremiumFactor       string `json:"premiumFactor"`       // `Premium Factor (%)` b25398
	Payment             string `json:"payment"`             // `Premium Payment Method` b25611
	SubjectTo           string `json:"subjectTo"`           // `Subject To` b25924
	AnnuityInterest     string `json:"annuityInterest"`     // `Annuity Interest (%)` b26092
	PremiumRefundFactor string `json:"premiumRefundFactor"` // `Premium Refund Factor (%)` b26306
	MaxDataReceive      string `json:"maxDataReceive"`      // `Max Production Data Receive` b27062
	Mature              string `json:"mature"`              // `Expired Date` b27250
	Birthday            string `json:"birthday"`            // `Birthday` b27960
	Currency            string `json:"currency"`            // `Currency` b28140
	CurrencyID          string `json:"currencyId"`          // `setCurrency_DT` b172
	ExtraMortality      string `json:"extraMortality"`      // `Extra Mortality (%)` b29256
	MaxContract         string `json:"maxContract"`         // `Max Contract (year)` b29442
	ProportionalTable   string `json:"proportionalTable"`   // `Proportional Table` b29626

	// Kunci view `PRODUCTINWARD_LIFE` tanpa medan form (PARITAS §3.2) -
	// dibaca dari JSON lama, tidak pernah dari klien.
	Ceding             string `json:"ceding"`
	TreatyNumber       string `json:"treatyNumber"`
	InwardTreatyNm     string `json:"inwardTreatyNm"`
	CedingRetentionPct string `json:"cedingRetentionPct"`
	CedingLimitXPN     string `json:"cedingLimitXpn"`
	RNMLimitPct        string `json:"rnmLimitPct"`
	LienClause         string `json:"lienClause"`
	Months             string `json:"months"`
}

// Asli - kunci baris JSON lama yang tidak dikelola layar (mis. `pxObjClass`,
// jejak `pxCreate*`), dibawa bolak-balik tanpa dibaca siapa pun selain
// repository. Kosong untuk baris baru.
type Asli = string

// BarisLien - `ProductName.LienClause` (grid b12201, RALAT R10).
type BarisLien struct {
	Usia    string `json:"usia"`    // `Usia saat Klaim` b12741
	Manfaat string `json:"manfaat"` // `% Manfaat yang dibayarkan` b12890
	Asli    Asli   `json:"asli,omitempty"`
}

// BarisDokumen - `ProductName.DocumentClaim` (grid b14601).
type BarisDokumen struct {
	Document string `json:"document"` // `Document List` b15125
	Asli     Asli   `json:"asli,omitempty"`
}

// BarisPlan - `ProductName.PlanList` (grid b31557).
type BarisPlan struct {
	Plan     string `json:"plan"`     // `Plan Name` b31845
	PlanID   string `json:"planId"`   // autocomplete `.ID` → `.PlanID`
	Name     string `json:"name"`     // `Bussines` b31994
	Benefit  string `json:"benefit"`  // `Benefit` b32143
	RIRate   string `json:"riRate"`   // `R/I Rate` b32296
	RIRateID string `json:"riRateId"` // `SetRIRate` 1 b249
	Asli     Asli   `json:"asli,omitempty"`
}

// BarisFinUW - `ProductName.FinancialUnderwritingList` (grid b37148).
type BarisFinUW struct {
	MinInsured  string `json:"minInsured"`  // `Min Insured` b37670
	MaxInsured  string `json:"maxInsured"`  // `Max Insured` b37818
	Employee    string `json:"employee"`    // `Employee` b37966
	NonEmployee string `json:"nonEmployee"` // `Non-Employee` b38115 → JSON `Non_Employee`
	Asli        Asli   `json:"asli,omitempty"`
}

// BarisUWLimit - `ProductName.UnderwritingLimitList` (grid b42075).
type BarisUWLimit struct {
	MinInsured  string `json:"minInsured"`  // `Min Insured` b42594
	MaxInsured  string `json:"maxInsured"`  // `Max Insured` b42742
	MinAge      string `json:"minAge"`      // `Min Age` b42890
	MaxAge      string `json:"maxAge"`      // `Max Age` b43038
	Medical     string `json:"medical"`     // `Medical` b43187 - teks bebas (R16)
	Description string `json:"description"` // `Description` b43336
	Asli        Asli   `json:"asli,omitempty"`
}

// BarisOutward - `ProductName.OutwardList` (diisi `GetReinsTypeOR_Life` 4.1 b770).
type BarisOutward struct {
	ReinsTypeID      string `json:"reinsTypeId"`
	ReinsTypeName    string `json:"reinsTypeName"`
	TransactionYear  string `json:"transactionYear"` // ← `TREATYYEAR`
	TreatyContractID string `json:"treatyContractId"`
	UnderwritingYear string `json:"underwritingYear"`
	// OvrComm - dibaca view `PRODUCT_LIFE` (`OutwardList[0].OVR_COMM`); korpus
	// tidak punya penulisnya; DEV: 0 dari 153 `OutwardList[0]` terisi - ditulis kosong (OQ-MPNL-09 ditutup).
	OvrComm string `json:"ovrComm"`
	Asli    Asli   `json:"asli,omitempty"`
}

// BarisKomentar - `ProductName.CommentList` (`AddCommentList_Act` 1 b235).
type BarisKomentar struct {
	Date         string `json:"date"`         // `Date` b62095 - `@CurrentDateTime()`
	OperatorName string `json:"operatorName"` // `PIC` b62246 - `OperatorID.pxInsName`
	IsApproved   string `json:"isApproved"`   // `param.status` (tidak dikirim)
	Suggest      string `json:"suggest"`      // `Comment` b62399
	Asli         Asli   `json:"asli,omitempty"`
}

// RingkasanProduk - satu baris grid daftar (`BrowseProduct_Life`, PARITAS §2).
type RingkasanProduk struct {
	ID           string `json:"id"`           // `ID` b72403
	Ceding       string `json:"ceding"`       // `Ceding` b72551
	TreatyNumber string `json:"treatyNumber"` // `Treaty Number` b72707
	InwardName   string `json:"inwardName"`   // `Treaty Name` b72866
	CreateOp     string `json:"createOp"`     // `Create Operator` b73025
	UpdateOp     string `json:"updateOp"`     // `Last Updated Operator` b73184
}
