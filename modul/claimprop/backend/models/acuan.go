package models

// Untuk apa berkas ini: ACUAN - bacaan basis data yang dibutuhkan port activity (RDB-List / Report Definition
// baca-saja). Port activity di paket ini tetap murni: setiap bacaan lewat antarmuka `Acuan`, diisi services dari
// repository (Oracle) atau dari gudang tiruan (uji).

import "context"

// MasterTreaty - halaman `pyWorkPage.TreatyInMaster`: JSON master treaty (`GetLimitsTreatyIn_SQL`: M_TREATY_IN /
// M_TREATY_IN_EDM `where id = IDMaster`, lalu Java `adoptJSONObject` - SetValueToClaim_Act langkah 3-4). Hanya medan
// yang DIBACA rule Claim Prop terjangkau (sensus `TreatyInMaster.*` 07-10-2026).
type MasterTreaty struct {
	ID                   string
	TreatyContractName   string
	ProportionType       string
	Ceding               string
	CedingID             string
	LeadingReinsSource   string
	LeadingReinsSourceID string
	Bordeaux             string
	BordereauxNote       string
	AccountingMode       string
	TeritorialScope      string
	Commencement         string // tanggal halaman ("2006-01-02")
	Termination          string
	TreatyYear           string
	RNMShareP            string
	StatusAkseptasi      string
	Limits               []LimitMaster
}

// LimitMaster - `TreatyInMaster.Limits(n)`.
type LimitMaster struct {
	TreatyType string
	Detail     []DetailLimit
}

// DetailLimit - `TreatyInMaster.Limits(n).Detail(m)`.
type DetailLimit struct {
	TreatyGroupID string
	RNMShare      string
	CashLossList  []CashLoss
	SpreadingList []SpreadingMaster
}

// CashLoss - `Detail(m).CashLossList(k)`: plafon cash call (DeleteEstimation_Act langkah 6-7).
type CashLoss struct {
	Currency string
	Value    string
}

// SpreadingMaster - `Limits(1).Detail(1).SpreadingList(k)`: anak spreading (SetTreatyNameSpreading_Act langkah 11 ->
// SpreadingBreakQS).
type SpreadingMaster struct {
	ReinsTypeID   string
	ReinsTypeName string
	Pct           string
}

// BarisMaster - satu baris grid popup "Data Master TreatyIn" (RD `BrowseCLAIM_MASTER_TREATY`, view CLAIM_MASTER_TREATY).
type BarisMaster struct {
	TreatyID           string `json:"treatyId"`
	ClassOfBusiness    string `json:"classOfBusiness"`
	ClassOfBusinessID  string `json:"classOfBusinessId"`
	TreatyContractName string `json:"treatyContractName"`
	SOB                string `json:"sob"`
	Ceding             string `json:"ceding"`
	TreatyType         string `json:"treatyType"`
	ProportionType     string `json:"proportionType"`
	TreatyGroup        string `json:"treatyGroup"`
	TreatyGroupID      string `json:"treatyGroupId"`
	TreatyYear         string `json:"treatyYear"`
}

// Popup "Data Master TreatyIn" (Section MasterTreatyInList, RD BrowseCLAIM_MASTER_TREATY): parameter section
// TREATYTYPE = "Proportional" menyaring `.PROPORTIONTYPE`. Filter per kolom, batas, dan urutan terbaru dulu = keputusan
// work owner 08-10-2026 (batas 20 di RD membuat 20 baris pertama satu treaty saja - pola popup NB Treaty In: saring di
// server lalu batas 500, 50 per halaman di layar).
const (
	ProporsiMaster = "Proportional"
	BatasMaster    = 500
)

// SaringanMaster - filter per kolom popup master (kosong = tanpa saringan); dicocokkan tanpa beda huruf, mengandung.
type SaringanMaster struct {
	TreatyID        string
	ClassOfBusiness string
	ContractName    string
	SOB             string
	InsuredName     string
	TreatyType      string
	TreatyGroup     string
	TreatyYear      string
}

// BerkasPolis - berkas NB / EDM Treaty In yang memuat satu nomor polis RNM. Tombol View (Pega: GetDetailPolis_act
// mengurai dokumen polis ke halaman kelas ASM-FW-GISFW-Work = berkas polis Treaty In; harness DetailPolisRealization
// tidak diekspor) membuka layar modul NB / EDM Treaty In di tab baru (keputusan work owner 08-10-2026).
type BerkasPolis struct {
	Modul string `json:"modul"`
	Kasus string `json:"kasus"`
}

// Nama modul frontend tujuan tombol View (`menu.ts` modul NB / EDM Treaty In).
const (
	ModulNBTreatyIn  = "nbtreatyin"
	ModulEDMTreatyIn = "edmtreatyin"
)

// ModulBerkasPolis - T_GENERAL_POLIS_TREATY.PRODKE 0 = NB Treaty In, >= 1 = generasi endorsemen EDM Treaty In.
func ModulBerkasPolis(prodke int) string {
	if prodke >= 1 {
		return ModulEDMTreatyIn
	}
	return ModulNBTreatyIn
}

// BarisPolis - satu baris grid popup "Data Polis" (RDB `SetPolicyTreatyProp` atas TREATYINPRODUCTION).
// ⚠️ Alias berbohong `TREATYGROUP AS "BusinessName"` diluruskan (AC 110): medannya bernama TreatyGroup.
type BarisPolis struct {
	PolicyNo             string `json:"policyNo"`
	NoOffer              string `json:"noOffer"`
	TreatyGroup          string `json:"treatyGroup"`
	SourceOfBusinessName string `json:"sourceOfBusinessName"`
	TreatyYear           string `json:"treatyYear"`
	Prodke               string `json:"prodke"`
	Quarter              string `json:"quarter"`
}

// Pilihan - satu opsi daftar pilihan (RD mata uang, jenis reasuransi, adjuster/consultant, provinsi).
type Pilihan struct {
	Nilai string `json:"nilai"`
	Label string `json:"label"`
	// Tambahan - kolom tambahan autocomplete (mis. nama adjuster).
	Tambahan map[string]string `json:"tambahan,omitempty"`
}

// BarisWilayah - satu baris `BrowseRW_SQL` (RW + CITY menurut kode pos). Alias berbohong (BRANCH_CODE, POLICY_NO,
// END_DATE, ...) diluruskan (AC 106).
type BarisWilayah struct {
	RWID, RW, DistrictID, District, CityID, City, ProvinceID, Province string
}

// RiwayatKlaimPolis - satu klaim lain pada polis yang sama (pembanding duplikat Date of Loss, CheckDateDOL_Act
// langkah 19-20). Alias berbohong `CariHistoryClaim_SQL` diluruskan (AC 106, §16).
type RiwayatKlaimPolis struct {
	KunciKasus string // INSKEY / IDPEGA kasus lama, atau ID kasus sistem baru
	IDKasus    string
	DateOfLoss string // tanggal halaman
	Ditolak    bool   // ada di CLAIMREJECTED (RejectedClaim_RD) - pengecualian AC 103
}

// AnggotaKomite - satu baris roster `FilterEmailKomiteWithLimit` (EMAILKOMITE).
type AnggotaKomite struct {
	ID         string
	OperatorID string
	Email      string
	Jabatan    string
	Degree     string
	// Approval, Comment, TanggalSetuju - keputusan anggota tangga kasus komite (dibaca untuk tampilan); kosong pada
	// roster calon.
	Approval, Comment, TanggalSetuju string
}

// RekeningBank - satu baris BANKACCOUNT.
type RekeningBank struct {
	ClientName, NameOfBank, BranchOfBank, AccountNo, SwiftCode, IDOfBank, CurrencyID string
}

// Retro - satu baris TREATYREINSURER (`GetListRetro_Sql`).
type Retro struct {
	ReinsurerID, ReinsurerName, PctShare, RiComm, AdditionalInfo string
}

// SpreadingPolis - satu baris spreading polis di TREATYINPRODUCTION: JN_REAS (TreatyType), PCT_SHARE_PREMI
// (SharePercentage), CURR_ID (Currency) dan CURRENCY.ID-nya (CurrencyID). Nilai DB apa adanya.
type SpreadingPolis struct {
	TreatyType, SharePercentage, CurrencyID, Currency string
}

// Acuan - seluruh bacaan basis data port activity Claim Prop.
type Acuan interface {
	// SpreadingPolis - spreading polis (TREATYINPRODUCTION per NOPOLIS, unik, urut JN_REAS): pengisi SpreadingClaim
	// saat polis dipilih dan sumber dropdown Treaty Type (keputusan work owner 08-10-2026; AddSpreading_Act tidak
	// diekspor).
	SpreadingPolis(ctx context.Context, nopolis string) ([]SpreadingPolis, error)
	// KursStandar = `POOLDATA.GETCURRENCYSTANDARD(CurrID, SYSDATE)` (RDB CurrencyStandard).
	KursStandar(ctx context.Context, currencyID string) (string, error)
	// NamaMataUang = RD BrowseCurrency_RD `.Currency` where `.ID = CurrID`.
	NamaMataUang(ctx context.Context, currencyID string) (string, error)
	// NamaJenisReasuransi = RD BrowseReinsuranceType_RD `.Note` where `.ID`.
	NamaJenisReasuransi(ctx context.Context, id string) (string, error)
	// IDJenisReasuransi = RD BrowseReinsuranceType_RD `.ID` where `.Note = Param.Note` and `.Type`, baris pertama
	// urut ID (SetNameTreaty_Act, nama treaty dari dropdown master).
	IDJenisReasuransi(ctx context.Context, nama, tipe string) (string, error)
	// MasterTreaty = GetLimitsTreatyIn_SQL + adoptJSONObject.
	MasterTreaty(ctx context.Context, id string) (MasterTreaty, bool, error)
	// AdaPolisMaster = GetNopolis_SQL (`TREATYINPRODUCTION where NOOFFER = IDMaster`) - ada setidaknya satu baris.
	AdaPolisMaster(ctx context.Context, idMaster string) (bool, error)
	// AdaPolisRealisasi = CekPolicyNumber_SQL (`TREATYINPRODUCTION where NOPOLIS = PolicyNo`).
	AdaPolisRealisasi(ctx context.Context, nopolis string) (bool, error)
	// TreatyGroupBisnis = GetTreatyGroupID (`TREATYBUSINESS` BIZCODE + TREATYYEAR + ISACTIVE='1'), baris pertama.
	TreatyGroupBisnis(ctx context.Context, bizCode, treatyYear string) (string, error)
	// YearOfQuartal = GetYearofQuartal (JSON_POLIS `DATA_JSON.YearOfQuartal` per NOPOLIS + PRODKE), baris pertama.
	YearOfQuartal(ctx context.Context, nopolis, prodke string) (string, error)
	// Wilayah = BrowseRW_SQL (RW + CITY where ZIPCODE).
	Wilayah(ctx context.Context, kodePos string) ([]BarisWilayah, error)
	// AlamatKlien = GetAddressCeding (M_CLIENT alamat klien milik agen `agentID`), empat bagian alamat.
	AlamatKlien(ctx context.Context, agentID string) ([4]string, bool, error)
	// NamaAdjuster = RD BrowseAdjusterConsultant `.NAME` where `.ID`, baris pertama.
	NamaAdjuster(ctx context.Context, id string) (string, error)
	// MarketingPolis = GetMObyNopol_SQL, baris pertama: ID, ClientID, ClientName, TeamGroup, BranchDetailID, BranchDetailName.
	MarketingPolis(ctx context.Context, nopolis string) ([6]string, bool, error)
	// RiwayatKlaimPolis = CariHistoryClaim_SQL + RejectedClaim_RD - klaim lain pada polis yang sama.
	RiwayatKlaimPolis(ctx context.Context, nopolis string) ([]RiwayatKlaimPolis, error)
	// KodeLamaBisnis = GetOldIDBusiness (`BUSINESS.OLDID where ID = BusinessCode`).
	KodeLamaBisnis(ctx context.Context, bizCode string) (string, error)
	// TahunTreaty = TreatyYearTreatyin_SQL (TREATYYEAR menurut treaty group dan tanggal `yyyyMMdd`).
	TahunTreaty(ctx context.Context, treatyGroupID, tanggalYMD string) (string, error)
	// LimitPLA = GetLimitPLATreatyin (PROPORTIONALARRG `TREATYDESCID='10001'`): kolom RP baris pertama.
	LimitPLA(ctx context.Context, treatyYear, treatyGroupID, reinsTypeID string) (string, bool, error)
	// DaftarRetro = GetListRetro_Sql (TREATYREINSURER).
	DaftarRetro(ctx context.Context, reinsTypeID, treatyYear, treatyGroupID string) ([]Retro, error)
	// RosterKomite = RD FilterEmailKomiteWithLimit (`LIMIT_BOTTOM <= nilai AND STS_KLAIM AND STS_AKTIF='1'`, urut DEGREE).
	RosterKomite(ctx context.Context, nilaiIDR, stsKlaim string) ([]AnggotaKomite, error)
	// LimitDirekturUtama = GetLimitDirekturUtama_SQL - batas Direktur Utama (AC 56; nama literal dibuang, AC 55).
	LimitDirekturUtama(ctx context.Context) (string, bool, error)
	// RekeningBank = GetDataBankAccount_sql (CLIENTID + CURRENCYID + STS_CODE='1').
	RekeningBank(ctx context.Context, clientID, currencyID string) ([]RekeningBank, error)
	// RekeningBankMataUang = GetDataBankAccount2_sql (CURRENCYID + STS_CODE='1').
	RekeningBankMataUang(ctx context.Context, currencyID string) ([]RekeningBank, error)
	// RekeningBankKlien = GetDatabyClientName (CLIENTID + STS_CODE='1').
	RekeningBankKlien(ctx context.Context, clientID string) ([]RekeningBank, error)
	// SaldoPremi = CekLunasPremi_Sql (`sum(ivd_trans_sign * ivd_total)` ARASAPAS.INVOICE x DETAIL_INVOICE).
	SaldoPremi(ctx context.Context, invoiceNo, currencyID string) (string, error)
	// AdaProteksiPremi = CekProteksiKlaim (OPENPROTEKSI_EDM TYPE='5' STS_AKSEP='1' POLICY_NO).
	AdaProteksiPremi(ctx context.Context, nopolis string) (bool, error)
}
