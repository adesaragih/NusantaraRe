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
	// tampil sebagai `Product Code` b3928. Tidak pernah dari klien untuk produk baru.
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
	// SalinanDari - `Copy` b59812 (`CopyProduct` b138/b161 hanya mengosongkan kedua
	// ID): produk baru mewarisi medan milik server produk asal ini.
	SalinanDari string `json:"salinanDari,omitempty"`
	// HitungOutward - checkbox `On Retention` b47303 diubah dalam sesi sunting
	// ini: `GetReinsTypeOR_Life` b47476 mengisi ulang `OutwardList` (OQ-MPNL-09).
	HitungOutward bool `json:"hitungOutward,omitempty"`
}

// ProdukUmum - halaman `ProductName` (`M_PRODUCT_LIFE.JSONDATA`).
type ProdukUmum struct {
	ProductName      string `json:"productName"`      // `Product Name` b3652
	Ceding           string `json:"ceding"`           // `Ceding` b4107
	CedingID         string `json:"cedingId"`         // `setCeding_DT` b179
	SOBName          string `json:"sobName"`          // `SOB` b4495
	SOBID            string `json:"sobId"`            // `setSOB_DT` b179
	RIComm           string `json:"riComm"`           // `Deduction (%)` b7137
	RIRisk           string `json:"riRisk"`           // `R/I Risk Name` b7430
	RIRiskID         string `json:"riRiskId"`         // `setRIRISK_DT` b179
	InwardName       string `json:"inwardName"`       // `Treaty Name` b8753
	TreatyNumber     string `json:"treatyNumber"`     // `Treaty Number` b8934
	Cause            string `json:"cause"`            // `Cause Of Loss` b10693
	CauseID          string `json:"causeId"`          // `setCauseOfLoss_DT` b168
	PolicyHolder     string `json:"policyHolder"`     // JSON `POLICYHODER` - salinan inward (`SaveProductName_Act` 1 b359)
	PolicyHolderName string `json:"policyHolderName"` // JSON `POLICYHODERNAME` - salinan inward
	CreateOp         string `json:"createOp"`         // `SaveProductName_Act` 6 b1368
	UpdateOp         string `json:"updateOp"`         // `SaveProductName_Act` 1 b359
	// IsORS - checkbox `On Retention` b47303 (RALAT R9).
	IsORS bool `json:"isOrs"`
	// Comment - area teks `Comment` popup `SaveProductName_Confirm` b1058.
	Comment string `json:"comment"`

	// Medan layar MATI (`OTHER 1=2`, PARITAS §3.1) - kuncinya dibaca view
	// `PRODUCT_LIFE`. Diisi dari JSON lama saja, tidak pernah dari klien.
	// TypeBasicRider - kunci JSON `TYPE` (`GenerateUpload_Act` b330: `1` = Basic, lainnya Rider).
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
	// ID baris inward - sama dengan ID produk untuk produk baru (R14, OQ-MPNL-02).
	ID        string `json:"id"`
	ProductID string `json:"productId"` // `SaveProductName_Act` 13 b2528

	PolicyHolder        string `json:"policyHolder"`        // JSON `POLICYHODER` - `setPolicyHolder_DT` b168
	PolicyHolderName    string `json:"policyHolderName"`    // `Policy Holder` b17129
	Insured             string `json:"insured"`             // `Insured` b18376
	AddendumNo          string `json:"addendumNo"`          // `Addendum No.` b18866
	AddendumWord        string `json:"addendumWord"`        // `Addendum` b19466
	AmandementNo        string `json:"amandementNo"`        // `Amandement No.` b20189
	AmandementSchd      string `json:"amandementSchd"`      // `Amandement` b20794
	MaxExpiredClaim     string `json:"maxExpiredClaim"`     // `Max Notification Claim Expired` b21813
	Begin               string `json:"begin"`               // `Begin Date` b22001
	STNC                string `json:"stnc"`                // `STNC` b22336
	CedingRetentionNum  string `json:"cedingRetentionNum"`  // `Ceding Retention (%)` b22526
	CedingLimit         string `json:"cedingLimit"`         // `Ceding's Limit` b22689
	Brokerage           string `json:"brokerage"`           // `Brokerage Fee (%)` b22947
	MinAge              string `json:"minAge"`              // `Minimum Age (Years)` b23110
	MaxAge              string `json:"maxAge"`              // `Maximum Age (Years)` b23297
	ExpiryAge           string `json:"expiryAge"`           // `Expiry Age (Years)` b23484
	ExtraPremi          string `json:"extraPremi"`          // `Extra Premium` b23667
	MinSumInsured       string `json:"minSumInsured"`       // `Min Sum Insured` b23828
	MaxSumInsured       string `json:"maxSumInsured"`       // `Max Sum Insured` b23988
	MaxSumReasured      string `json:"maxSumReasured"`      // `Max Sum Reasured` b24246
	RNMShare            string `json:"rnmShare"`            // `Nusantara Re Share (%)` b24705
	RNMLimitNum         string `json:"rnmLimitNum"`         // `Nusantara Re's Limit` b25268
	PremiumFactor       string `json:"premiumFactor"`       // `Premium Factor (%)` b25430
	Payment             string `json:"payment"`             // `Premium Payment Method` b25642
	SubjectTo           string `json:"subjectTo"`           // `Subject To` b25952
	AnnuityInterest     string `json:"annuityInterest"`     // `Annuity Interest (%)` b26124
	PremiumRefundFactor string `json:"premiumRefundFactor"` // `Premium Refund Factor (%)` b26338
	MaxDataReceive      string `json:"maxDataReceive"`      // `Max Production Data Receive` b27094
	Mature              string `json:"mature"`              // `Expired Date` b27282
	Birthday            string `json:"birthday"`            // `Birthday` b27992
	Currency            string `json:"currency"`            // `Currency` b28173
	CurrencyID          string `json:"currencyId"`          // `setCurrency_DT` b168
	ExtraMortality      string `json:"extraMortality"`      // `Extra Mortality (%)` b29288
	MaxContract         string `json:"maxContract"`         // `Max Contract (year)` b29474
	ProportionalTable   string `json:"proportionalTable"`   // `Proportional Table` b29658

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

// BarisLien - `ProductName.LienClause` (grid b12204, RALAT R10).
type BarisLien struct {
	Usia    string `json:"usia"`    // `Usia saat Klaim` b12745
	Manfaat string `json:"manfaat"` // `% Manfaat yang dibayarkan` b12894
	Asli    Asli   `json:"asli,omitempty"`
}

// BarisDokumen - `ProductName.DocumentClaim` (grid b14604).
type BarisDokumen struct {
	Document string `json:"document"` // `Document List` b15129
	Asli     Asli   `json:"asli,omitempty"`
}

// BarisPlan - `ProductName.PlanList` (grid b31560).
type BarisPlan struct {
	Plan     string `json:"plan"`     // `Plan Name` b31849
	PlanID   string `json:"planId"`   // autocomplete `.ID` → `.PlanID`
	Name     string `json:"name"`     // `Bussines` b31998
	Benefit  string `json:"benefit"`  // `Benefit` b32147
	RIRate   string `json:"riRate"`   // `R/I Rate` b32300
	RIRateID string `json:"riRateId"` // `SetRIRate` 1 b247
	Asli     Asli   `json:"asli,omitempty"`
}

// BarisFinUW - `ProductName.FinancialUnderwritingList` (grid b37151).
type BarisFinUW struct {
	MinInsured  string `json:"minInsured"`  // `Min Insured` b37674
	MaxInsured  string `json:"maxInsured"`  // `Max Insured` b37822
	Employee    string `json:"employee"`    // `Employee` b37970
	NonEmployee string `json:"nonEmployee"` // `Non-Employee` b38119 → JSON `Non_Employee`
	Asli        Asli   `json:"asli,omitempty"`
}

// BarisUWLimit - `ProductName.UnderwritingLimitList` (grid b42078).
type BarisUWLimit struct {
	MinInsured  string `json:"minInsured"`  // `Min Insured` b42598
	MaxInsured  string `json:"maxInsured"`  // `Max Insured` b42746
	MinAge      string `json:"minAge"`      // `Min Age` b42894
	MaxAge      string `json:"maxAge"`      // `Max Age` b43042
	Medical     string `json:"medical"`     // `Medical` b43191 - teks bebas (R16)
	Description string `json:"description"` // `Description` b43340
	Asli        Asli   `json:"asli,omitempty"`
}

// BarisOutward - `ProductName.OutwardList` (diisi `GetReinsTypeOR_Life` 4.1 b768).
type BarisOutward struct {
	ReinsTypeID      string `json:"reinsTypeId"`
	ReinsTypeName    string `json:"reinsTypeName"`
	TransactionYear  string `json:"transactionYear"` // ← `TREATYYEAR`
	TreatyContractID string `json:"treatyContractId"`
	UnderwritingYear string `json:"underwritingYear"`
	// OvrComm - dibaca view `PRODUCT_LIFE` (`OutwardList[0].OVR_COMM`); korpus
	// tidak punya penulisnya (OQ-MPNL-09).
	OvrComm string `json:"ovrComm"`
	Asli    Asli   `json:"asli,omitempty"`
}

// BarisKomentar - `ProductName.CommentList` (`AddCommentList_Act` 1 b233).
type BarisKomentar struct {
	Date         string `json:"date"`         // `Date` b62099 - `@CurrentDateTime()`
	OperatorName string `json:"operatorName"` // `PIC` b62250 - `OperatorID.pxInsName`
	IsApproved   string `json:"isApproved"`   // `param.status` (tidak dikirim)
	Suggest      string `json:"suggest"`      // `Comment` b62403
	Asli         Asli   `json:"asli,omitempty"`
}

// RingkasanProduk - satu baris grid daftar (`BrowseProduct_Life`, PARITAS §2).
type RingkasanProduk struct {
	ID           string `json:"id"`           // `ID` b72406
	Ceding       string `json:"ceding"`       // `Ceding` b72554
	TreatyNumber string `json:"treatyNumber"` // `Treaty Number` b72711
	InwardName   string `json:"inwardName"`   // `Treaty Name` b72870
	CreateOp     string `json:"createOp"`     // `Create Operator` b73029
	UpdateOp     string `json:"updateOp"`     // `Last Updated Operator` b73188
}
