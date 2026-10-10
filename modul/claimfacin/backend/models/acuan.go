package models

// Untuk apa berkas ini: ACUAN - bacaan basis data yang dibutuhkan port activity Claim Fac In (RDB-List / Report
// Definition baca-saja). Port activity di paket ini tetap murni: setiap bacaan lewat antarmuka `Acuan`, diisi services
// dari repository (Oracle) atau dari gudang tiruan (uji). Pola disalin dari `modul/claimnonprop/backend/models/acuan.go`.

import "context"

// BarisPolisCari - satu baris grid pop-up Choose Polis (`GetPolisForClaim_SQL`, FACINPRODUCTION; alias Pega dipakai
// sebagai nama properti baris).
type BarisPolisCari struct {
	PolicyNo             string `json:"policyNo"`
	CustomerName         string `json:"customerName"`
	SourceOfBusinessName string `json:"sourceOfBusinessName"`
	CedingCoName         string `json:"cedingCoName"`
	QQ                   string `json:"qq"`
	StartDateTime        string `json:"startDateTime"`
	EndDateTime          string `json:"endDateTime"`
	Prodke               string `json:"prodke"`
	BusinessName         string `json:"businessName"`
}

// Jenis pencarian polis (`SearchPolis_act` langkah 1-4: Primary.SearchType).
const (
	CariNoPolis = "1" // A.NOPOLIS = :teks
	CariCeding  = "2" // A.CEDINGCO LIKE :teks || '%'
	CariInsured = "3" // A.INSUREDNAME LIKE :teks || '%'
	CariQQ      = "4" // A.QQNAME LIKE :teks || '%'
)

// BatasCariPolis - batas baris grid Choose Polis (FACINPRODUCTION ±1,88 juta baris; prompt §4: kueri wajib berkunci).
const BatasCariPolis = 500

// RiwayatKlaimPolis - satu klaim pada polis yang sama (BrowseHistoryClaim / CariHistoryClaim_SQL atas JSON_KLAIM,
// ditambah klaim FACIN sistem baru - JSON_KLAIM sistem baru tidak berisi DATA_JSON).
type RiwayatKlaimPolis struct {
	// KunciKasus - IDPEGA JSON_KLAIM (`CariHistoryClaim_SQL` BRANCH_CODE) / ID kasus sistem baru.
	KunciKasus    string
	ClaimNo       string
	DateOfLoss    string // tanggal halaman
	CauseOfLoss   string
	ClaimEstimate string
}

// BarisWilayah - satu baris pencarian kode pos (BrowseRW_SQL; alias diluruskan, pola Claim Prop / Non Prop).
type BarisWilayah struct {
	RWID, RW, DistrictID, District, CityID, City, ProvinceID, Province string
}

// Pilihan - satu opsi daftar pilihan.
type Pilihan struct {
	Nilai string `json:"nilai"`
	Label string `json:"label"`
	// Tambahan - kolom tambahan autocomplete.
	Tambahan map[string]string `json:"tambahan,omitempty"`
}

// AnggotaKomite - satu baris roster `FilterEmailKomiteWithLimit` (EMAILKOMITE STS_KLAIM FACIN).
type AnggotaKomite struct {
	ID          string
	OperatorID  string
	Email       string
	Jabatan     string
	Degree      string
	LimitBottom string
	LimitTop    string
	// Approval, Comment, TanggalSetuju - keputusan anggota tangga kasus komite (tampilan); kosong pada roster calon.
	Approval, Comment, TanggalSetuju string
}

// RekeningBank - satu baris BANKACCOUNT.
type RekeningBank struct {
	ClientName, NameOfBank, BranchOfBank, AccountNo, SwiftCode, IDOfBank, CurrencyID, Currency string
}

// Acuan - seluruh bacaan basis data port activity Claim Fac In.
type Acuan interface {
	// RekeningBank = GetDataBankAccount_sql (`CLIENTID AND CURRENCYID AND STS_CODE = '1'`, SetPayable_Act 6).
	RekeningBank(ctx context.Context, klien, cur string) ([]RekeningBank, error)
	// RekeningBankMataUang = GetDataBankAccount2_sql (`CURRENCYID AND STS_CODE = '1'`, SetPayable_Act 8.2).
	RekeningBankMataUang(ctx context.Context, cur string) ([]RekeningBank, error)
	// CariPolis = GetPolisForClaim_SQL (FACINPRODUCTION) - teks cari DIIKAT, tidak disisipkan (perbaikan prompt §5 butir 1).
	CariPolis(ctx context.Context, jenis, teks string) ([]BarisPolisCari, error)
	// DokumenPolis = GetCopyNBForClaim (`DATA_JSON FROM JSON_POLIS WHERE nopolis AND prodke`).
	DokumenPolis(ctx context.Context, nopolis, prodke string) ([]byte, bool, error)
	// RiwayatKlaimPolis = BrowseHistoryClaim / CariHistoryClaim_SQL (+ klaim FACIN sistem baru).
	RiwayatKlaimPolis(ctx context.Context, nopolis string) ([]RiwayatKlaimPolis, error)
	// KlaimDitolak = RD RejectedClaim_RD (CLAIMREJECTED `.INSKEY`) - kunci kasus ada di daftar klaim ditolak.
	KlaimDitolak(ctx context.Context, kunci string) (bool, error)
	// KursStandar = CurrencyStandard (`POOLDATA.GETCURRENCYSTANDARD(:cur, SYSDATE)`).
	KursStandar(ctx context.Context, currencyID string) (string, error)
	// NamaJenisReas = RD BrowseReinsuranceType_RD (`REINSURANCETYPE.NOTE` menurut ID) - CheckEstimateValue 2.
	NamaJenisReas(ctx context.Context, id string) (string, error)
	// NamaMataUang = GetCurrency (`CURRENCY FROM CURRENCY WHERE ID`).
	NamaMataUang(ctx context.Context, currencyID string) (string, error)
	// Wilayah = BrowseRW_SQL (RW + CITY menurut ZIPCODE).
	Wilayah(ctx context.Context, kodePos string) ([]BarisWilayah, error)
	// NamaAdjuster = RD BrowseAdjusterConsultant `.NAME` where `.ID`, baris pertama.
	NamaAdjuster(ctx context.Context, id string) (string, error)
	// AgenKlien = GetLeaderReport (`id FROM agent WHERE clientname = :nama`), baris terakhir (loop menimpa).
	AgenKlien(ctx context.Context, namaKlien string) (string, error)
	// AlamatAgen = GetAddressCeding (empat bagian alamat klien milik agen).
	AlamatAgen(ctx context.Context, agentID string) ([4]string, error)
	// RosterKomite - EMAILKOMITE STS_KLAIM FACIN aktif urut DEGREE (FilterEmailKomiteWithLimit).
	RosterKomite(ctx context.Context) ([]AnggotaKomite, error)
	// TingkatPelaku - JABATAN baris roster FACIN aktif pelaku ("" bila tidak ada).
	TingkatPelaku(ctx context.Context, operatorID string) (string, error)
	// NamaPelaku = `OperatorID.pyLabel` (M_LOGIN_GO.NAME menurut LOGIN_ID); akun tanpa nama = akun itu sendiri.
	NamaPelaku(ctx context.Context, akun string) (string, error)
	// EmailPelaku = `OperatorID.pyAddress` (M_LOGIN_GO.EMAIL).
	EmailPelaku(ctx context.Context, akun string) (string, error)
	// EstimasiKasusTerbuka = GetDataEstimation (`SUM(CLAIM_VALUE) REINSURANCE.TRLOSS_DETAIL_T WHERE NO_SPK = pyID AND
	// NO_AKSEP IS NULL`) - HANYA dipanggil di produksi (skema luar tidak terbaca dari akun DEV, pola OQ-CP-11).
	EstimasiKasusTerbuka(ctx context.Context, kasusID string) (string, error)
	// EstimasiPolisTerbuka = GetAllClaimWithSameNopolis_Sql (Σ CLAIM_VALUE TRLOSS_DETAIL_T STS_AKSEP kosong klaim
	// berpolis sama) - HANYA dipanggil di produksi.
	EstimasiPolisTerbuka(ctx context.Context, nopolis string) (string, error)
	// BatasTreaty = GetDataTreatyLimit_Sql (baris pertama; ada = false bila tidak ada baris). mulaiPolis "2006-01-02".
	BatasTreaty(ctx context.Context, kodeBisnis, jenisTreaty, mulaiPolis string) (BatasTreaty, bool, error)
	// QuotaShare = GetQuotaShare (PROPORTIONALARRG TREATYDESCID 10001 menurut tahun, grup, induk).
	QuotaShare(ctx context.Context, tahun, grup, induk string) ([]BarisQuotaShare, error)
	// ReasuradurTreaty = GetListRetro_Sql (TREATYREINSURER menurut jenis, tahun, grup).
	ReasuradurTreaty(ctx context.Context, jenis, tahun, grup string) ([]ReasuradurTreaty, error)
	// MarketingOfficer = SetMOClaim_Act langkah 4 (Obj-Browse MARKETINGOFFICER menurut ID): ID, CLIENTID, CLIENTNAME,
	// TEAMGROUP, BRANCHDETAILID, BRANCHDETAILNAME.
	MarketingOfficer(ctx context.Context, id string) ([6]string, bool, error)
}
