package models

// Untuk apa berkas ini: ACUAN - bacaan basis data yang dibutuhkan port activity Claim Non Prop (RDB-List / Report
// Definition baca-saja). Port activity di paket ini tetap murni: setiap bacaan lewat antarmuka `Acuan`, diisi services
// dari repository (Oracle) atau dari gudang tiruan (uji).

import "context"

// MasterTreaty - halaman `pyWorkPage.TreatyInMaster` Non Prop: JSON master treaty (`GetLimitsTreatyIn_SQL`: M_TREATY_IN
// UNION ALL M_TREATY_IN_EDM `where id = IDMaster`, lalu Java `adoptJSONObject` - SetValueClaimTNP_Act langkah 3-4).
// Hanya medan yang DIBACA rule Claim Non Prop (sensus `TreatyInMaster.*` 09-10-2026; bentuk JSON dibaca dari DEV
// M_TREATY_IN NonProportional).
type MasterTreaty struct {
	ID                    string
	ProportionType        string
	Ceding                string
	CedingID              string
	LeadingReinsSource    string
	LeadingReinsSourceID  string
	Bordeaux              string
	BordereauxNote        string
	AccountingMode        string
	AccountingModeNonProp string
	TeritorialScope       string
	Commencement          string // tanggal halaman ("2006-01-02")
	Termination           string
	TreatyYear            string
	RNMShare              string
	EDMState              string
	StatusAkseptasi       string
	Limits                []LimitXOL
	LimitSummaryList      []RingkasanLimit
	CurrencyList          []KursMaster
	Share                 []ShareXOL
}

// LimitXOL - `TreatyInMaster.Limits(n)` Non Prop (layer XoL).
type LimitXOL struct {
	Layer, LayerType, LayerPart, LayerPartType string
	Currency, Currency2                        string
	Limit, Limit2                              string
	Deductible, Deductible2                    string
	ReinstatementPct                           string
	IsCombineMDP, NoRIPCalculation             string
	// TreatyGroups - `.TreatyGroupList(k).TreatyGroup` (CountLossAllocation_act langkah 11.1).
	TreatyGroups []string
	MDPList      []NilaiMataUang
}

// NilaiMataUang - baris `{Currency, Value}` (MDPList).
type NilaiMataUang struct {
	Currency, Value string
}

// RingkasanLimit - `TreatyInMaster.LimitSummaryList(n)` (MDP gabungan, CountLossAllocation_act langkah 17.2.10).
type RingkasanLimit struct {
	Layer, LayerType    string
	MDP, MDP2           string
	Currency, Currency2 string
}

// KursMaster - `TreatyInMaster.CurrencyList(n)` (`.Currency`, `.Conversion`).
type KursMaster struct {
	Currency, Conversion string
}

// ShareXOL - `TreatyInMaster.Share(n)` (spreading XoL, CountLossAllocation_act langkah 21).
type ShareXOL struct {
	SpreadingTypeIDXOL, SpreadingTypeXOL, SpreadingTotalPctXOL string
	SpreadingListXOL                                           []SpreadingMaster
}

// SpreadingMaster - `Share(1).SpreadingListXOL(k)`.
type SpreadingMaster struct {
	ReinsTypeID, ReinsTypeName, Pct string
}

// BarisMaster - satu baris grid popup ChooseMasterTNonProp (RDB `BrowseDtlTreatyNP` / `BrowseDtlTreatyNPEDM` /
// `BrowseTreatyNP(_EDM)`): TREATYINDETAIL / TREATYINDETAILEDM `PROPORTIONTYPE = 'NonProportional'`.
type BarisMaster struct {
	TreatyID           string `json:"treatyId"`
	TreatyContractName string `json:"treatyContractName"`
	ProportionType     string `json:"proportionType"`
	ClassOfBusinessID  string `json:"classOfBusinessId"`
	ClassOfBusiness    string `json:"classOfBusiness"`
	SOB                string `json:"sob"`
	Ceding             string `json:"ceding"`
	TreatyYear         string `json:"treatyYear"`
	TreatyGroup        string `json:"treatyGroup"`
}

// Popup ChooseMasterTNonProp: `PROPORTIONTYPE = 'NonProportional'` (RDB BrowseDtlTreatyNP). Saringan per kolom dicari di
// server lalu batas 500 - pola popup Claim Prop (keputusan work owner Claim Prop 08-10-2026, prompt §8 P2).
const (
	ProporsiMaster = "NonProportional"
	BatasMaster    = 500
)

// Sumber master popup (BrowseDtlMasterTNP_Act Param.TreatyIN): IN = TREATYINDETAIL ∪ TREATYINDETAILEDM, INEDM =
// TREATYINDETAILEDM saja (tombol Choose Master Input Acceptation setelah ada akseptasi).
const (
	SumberMasterIN    = "IN"
	SumberMasterINEDM = "INEDM"
)

// SaringanMaster - filter per kolom popup master (kosong = tanpa saringan); dicocokkan tanpa beda huruf, mengandung.
type SaringanMaster struct {
	Sumber          string
	TreatyID        string
	ContractName    string
	SOB             string
	Ceding          string
	TreatyGroup     string
	ClassOfBusiness string
	TreatyYear      string
}

// BarisPolis - satu baris `GetDataPolisNonProp_SQL` (TREATYINPRODUCTION `NOOFFER = IDMaster OR NOOFFER = 7 karakter
// pertama IDMaster`). Alias CARI diluruskan (OQ-CNP-24 bawaan "ikut SQL").
type BarisPolis struct {
	IDPega       string `json:"idPega"`
	PolicyNo     string `json:"policyNo"`
	BusinessCode string `json:"businessCode"`
	InsuredName  string `json:"insuredName"`
	BeginDate    string `json:"beginDate"`
	EndDate      string `json:"endDate"`
	CedingCo     string `json:"cedingCo"`
	SOB          string `json:"sob"`
	TreatyGroup  string `json:"treatyGroup"`
}

// BerkasPolis - berkas NB / EDM Treaty In yang memuat satu nomor polis (View polis, bawaan OQ-CNP-13 = jendela Modal
// penuh NB/EDM Treaty In, pola Claim Prop).
type BerkasPolis struct {
	Modul string `json:"modul"`
	Kasus string `json:"kasus"`
}

// Nama modul frontend tujuan tombol View.
const (
	ModulNBTreatyIn  = "nbtreatyin"
	ModulEDMTreatyIn = "edmtreatyin"
)

// ModulBerkasPolis - T_GENERAL_POLIS_TREATY.PRODKE 0 = NB Treaty In, >= 1 = EDM Treaty In.
func ModulBerkasPolis(prodke int) string {
	if prodke >= 1 {
		return ModulEDMTreatyIn
	}
	return ModulNBTreatyIn
}

// Pilihan - satu opsi daftar pilihan (RD mata uang, adjuster/consultant, provinsi, Loss Allocation, rekening).
type Pilihan struct {
	Nilai string `json:"nilai"`
	Label string `json:"label"`
	// Tambahan - kolom tambahan autocomplete (mis. nama adjuster, rincian rekening).
	Tambahan map[string]string `json:"tambahan,omitempty"`
}

// BarisWilayah - satu baris `BrowseRW_SQL` (RW + CITY menurut kode pos), alias diluruskan.
type BarisWilayah struct {
	RWID, RW, DistrictID, District, CityID, City, ProvinceID, Province string
}

// RiwayatKlaimPolis - satu klaim lain pada polis yang sama (CekHistoryClaimNonProp_SQL / CariHistoryClaim_SQL).
type RiwayatKlaimPolis struct {
	KunciKasus string
	IDKasus    string
	IDMaster   string
	Bisnis     string
	DateOfLoss string // tanggal halaman
	Ditolak    bool   // ada di CLAIMREJECTED
}

// AnggotaKomite - satu baris roster `FilterEmailKomiteWithLimit` (EMAILKOMITE STS_KLAIM NONPROP).
type AnggotaKomite struct {
	ID         string
	OperatorID string
	Email      string
	Jabatan    string
	Degree     string
	// Approval, Comment, TanggalSetuju - keputusan anggota tangga kasus komite (tampilan); kosong pada roster calon.
	Approval, Comment, TanggalSetuju string
}

// RekeningBank - satu baris BANKACCOUNT.
type RekeningBank struct {
	ClientName, NameOfBank, BranchOfBank, AccountNo, SwiftCode, IDOfBank, CurrencyID, Currency string
}

// BarisReinsType - satu baris `GetReinsuranceTypeBYName_SQL` (REINSURANCETYPE type '4' FLAG 'active' NOTE = nama).
type BarisReinsType struct {
	ID, Note string
}

// Acuan - seluruh bacaan basis data port activity Claim Non Prop.
type Acuan interface {
	// MasterTreaty = GetLimitsTreatyIn_SQL + adoptJSONObject.
	MasterTreaty(ctx context.Context, id string) (MasterTreaty, bool, error)
	// JenisReasXOL = GetReinsuranceTypeBYName_SQL (REINSURANCETYPE `type='4' and FLAG='active' and note = nama`).
	JenisReasXOL(ctx context.Context, nama string) ([]BarisReinsType, error)
	// NamaMataUang = GetCurrency (`CURRENCY where id = CurrID`).
	NamaMataUang(ctx context.Context, currencyID string) (string, error)
	// DaftarPolis = GetDataPolisNonProp_SQL.
	DaftarPolis(ctx context.Context, idMaster string) ([]BarisPolis, error)
	// AdaPolisMaster = GetNopolis_SQL (`TREATYINPRODUCTION where NOOFFER = IDMaster`) - ada setidaknya satu baris.
	AdaPolisMaster(ctx context.Context, idMaster string) (bool, error)
	// NamaTreatySpreading = GetTreatyName_SQL (PROPORTIONALARRG 'TREATY LIMIT' menurut BusinessCode dan tanggal mulai
	// polis `Track.CARI33`) - REINSTYPENAME (CARI2), sumber `ListTreaty` (InputOutStandingClmTNP_PreAct langkah 3).
	NamaTreatySpreading(ctx context.Context, bizCode, tanggalMulai string) ([]string, error)
	// Wilayah = BrowseRW_SQL (RW + CITY where ZIPCODE).
	Wilayah(ctx context.Context, kodePos string) ([]BarisWilayah, error)
	// NamaAdjuster = RD BrowseAdjusterConsultant `.NAME` where `.ID`, baris pertama.
	NamaAdjuster(ctx context.Context, id string) (string, error)
	// MarketingPolis = GetMObyNopol_SQL, baris pertama: ID, ClientID, ClientName, TeamGroup, BranchDetailID, BranchDetailName.
	MarketingPolis(ctx context.Context, nopolis string) ([6]string, bool, error)
	// RiwayatKlaimPolis = CekHistoryClaimNonProp_SQL / CariHistoryClaim_SQL - klaim lain pada polis yang sama.
	RiwayatKlaimPolis(ctx context.Context, nopolis string) ([]RiwayatKlaimPolis, error)
	// KodeLamaBisnis = GetDataBusiness_SQL (`BUSINESS.OLDID where NOTE = nama bisnis`).
	KodeLamaBisnis(ctx context.Context, namaBisnis string) (string, error)
	// RekeningBank = GetDataBankAccount_sql (CLIENTID + CURRENCYID + STS_CODE='1').
	RekeningBank(ctx context.Context, clientID, currencyID string) ([]RekeningBank, error)
	// RekeningBankKlien = GetDatabyClientName (CLIENTID + STS_CODE='1').
	RekeningBankKlien(ctx context.Context, clientID string) ([]RekeningBank, error)
	// RekeningBankNama = GetDatabyClientName2 / 3 (CLIENTNAME [+ CURRENCYID]).
	RekeningBankNama(ctx context.Context, clientName, currencyID string) ([]RekeningBank, error)
	// RosterKomite - EMAILKOMITE STS_KLAIM NONPROP aktif urut DEGREE (FilterEmailKomiteWithLimit; batas dipakai hanya
	// bila `hanyaTingkat1`, OQ-CNP-01).
	RosterKomite(ctx context.Context, hanyaTingkat1 bool) ([]AnggotaKomite, error)
	// AlamatAgen = GetAddressCeding (GetReportStatus_Act langkah 8): alamat klien agen `InputData.CARI1` - empat bagian
	// alamat (kosong bila tidak ada baris).
	AlamatAgen(ctx context.Context, agentID string) ([4]string, error)
	// NamaPelaku = `OperatorID.pyLabel` (M_LOGIN_GO.NAME menurut LOGIN_ID); akun tanpa nama = akun itu sendiri.
	NamaPelaku(ctx context.Context, akun string) (string, error)
}
