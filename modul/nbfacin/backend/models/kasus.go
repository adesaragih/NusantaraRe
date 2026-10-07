package models

// Kasus - satu case NB untuk layar Inward Facultative (tiket 31): baris T_WORK_POLIS
// (milik premiumlistlife, K-064) + T_NB_OPPORTUNITY + nama tertanggung T_M_ACCOUNT +
// blok General (T_GENERAL_POLIS / T_QUOTATIONDATA). Teks apa adanya; NULL = "".
type Kasus struct {
	CaseID      string
	Position    string
	StatusWork  string
	Opportunity Opportunity
	InsuredName string
	General     General
}

// General - isian blok General (section Periode) sebagaimana TERSIMPAN. Tanggal = teks
// bentuk Pega (butir 78.1): OfferingDate 'YYYYMMDD', StartDateTime/EndDateTime
// 'YYYYMMDDTHHMMSS.mmm GMT'. Tersimpan = baris T_GENERAL_POLIS sudah ada.
type General struct {
	ReffNumber         string // T_QUOTATIONDATA.NO_OFFER_SLIP  (.QuotationData.NoOfferSlip, sel 9)
	QQName             string // T_QUOTATIONDATA.QQ_NAME        (.QuotationData.QQName, sel 20)
	StartDateTime      string // T_GENERAL_POLIS.START_DATE_TIME (.PolicyData.StartDateTime, sel 21)
	OfferingDate       string // T_GENERAL_POLIS.OFFERING_DATE   (.PolicyData.OfferingDate, sel 22)
	EndDateTime        string // T_GENERAL_POLIS.END_DATE_TIME   (.PolicyData.EndDateTime, sel 60)
	PolicyType         string // T_QUOTATIONDATA.POLICY_TYPE    (.QuotationData.PolicyType, sel 26)
	MarketingID        string // T_QUOTATIONDATA.MOID           (.QuotationData.MOID, sel 75)
	Day                string // T_QUOTATIONDATA.EDM_DAY        (.QuotationData.EDMDay, sel 78)
	TypeFacultative    string // T_QUOTATIONDATA.TYPE_FACULTATIVE (.QuotationData.TypeFacultative, sel 43)
	SourceOfBusinessID string // T_QUOTATIONDATA.SOURCE_OF_BUSINESS (.QuotationData.SourceOfBusiness, kode SOB - tiket 33)
	SourceOfBusiness   string // T_QUOTATIONDATA.SOB_NAME       (.QuotationData.SobName, sel 48) - ditulis server dari AGENT (E-4)
	CedingCoName       string // T_QUOTATIONDATA.CEDING_CO_NAME (.QuotationData.CedingCoName, sel 49) - gabungan ";" nama dari AGENT (tiket 34)
	GroupName          string // T_QUOTATIONDATA.GROUP_NAME     (.QuotationData.GroupName, sel 56) - tampil saja
	OldPolicyNumber    string // T_GENERAL_POLIS.FOLLOWING      (.Following, sel 72) - tampil saja
	// CedingList - baca: T_CEDINGCOLIST urut SEQ_NO (tiket 34). CedingIDs - tulis: kode urut
	// pilih; nama diambil server dari AGENT; kosong = daftar dikosongkan.
	CedingList []Ceding
	CedingIDs  []string
	Tersimpan  bool
}

// MarketingOfficer - satu pilihan dropdown Marketing Name (sel 75): ID disimpan ke MOID,
// Nama = CLIENTNAME (pyPrompt `.ClientName`).
type MarketingOfficer struct {
	ID   string
	Nama string
}

// BarisPortal - satu baris daftar case NB di portal Opportunity (tiket 32). NULL = "".
type BarisPortal struct {
	CaseID        string // T_WORK_POLIS.ID
	Name          string // T_NB_OPPORTUNITY.BUSINESS_PROSPECT_NAME
	GroupBusiness string // T_NB_OPPORTUNITY.GROUP_BUSINESS
	InsuredName   string // T_M_ACCOUNT.INSUREDNAME lewat ACCOUNT_ID
	Marketing     string // MARKETINGOFFICER.CLIENTNAME lewat T_QUOTATIONDATA.MOID; kosong bila belum ada
	Status        string // T_WORK_POLIS.STATUS_WORK
}

// SOB - satu pilihan popup Change SOB (tiket 33): baris tabel AGENT.
type SOB struct {
	ID       string
	ClientID string
	Name     string
}

// Ceding - satu baris daftar Ceding Co (tiket 34): kode AGENT.ID dan namanya.
type Ceding struct {
	ID   string
	Name string
}

// RiskAddress - satu baris hasil popup Choose Risk Address (tiket 36): POOLDATA.RISKADDRESS.
type RiskAddress struct {
	ID, Title, Address, NationName, ProvinceName, CityName, DistrictName, TerritoryName, PostalCode string
}

// SaringRisk - saringan popup Choose Risk Address (kosong = dilewati); urutan = filter RD
// BrowseRisksAddress_RD A..G (tiket 36).
type SaringRisk struct {
	Address, ZipCode, Country, Province, City, District, Territory string
}

// BarisRW - satu saran Zip Code (tiket 37): baris POOLDATA.RW (NOTE = territory, NATION = negara).
type BarisRW struct {
	ZipCode, TerritoryName, DistrictName, CityName, ProvinceName, NationName string
}

// AlamatBaru - isian popup Add alamat risiko (tiket 37) = parameter P_* prosedur
// InsertUpdateRISKADDRESS tanpa P_ID.
type AlamatBaru struct {
	NationName, ProvinceName, DistrictName, CityName, TerritoryName, Title, Address, PostalCode string
}
