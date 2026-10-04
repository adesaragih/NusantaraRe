// Package models memuat bentuk data kasus endorsement Fac In di memori.
//
// Untuk apa berkas ini: agregat penawaran (`OfferFacIn`) beserta daftar
// barisnya, sebatas properti yang DIBACA atau DITULIS rule before-image
// (spec `docs/04-spec/05-spec-edm-before-image.md`, tiket E06-E11).
//
// ⛔ Ini BUKAN skema lengkap halaman Pega. Struktur JSON `JSON_POLIS.DATA_JSON`
// dan isi penuh tiap kelas belum terverifikasi (CLAUDE.md §4.5) - pembacanya
// milik tiket E17 yang masih ter-block. Properti yang tidak disentuh rule
// before-image TIDAK dimodelkan, dan karena itu tidak ikut tersalin.
//
// ⛔ Properti yang hanya DISALIN (tidak dihitung) bertipe `string`: nilainya
// dibawa apa adanya, dan tipenya di Pega tidak ditebak. Hanya yang masuk
// aritmetika diberi tipe: uang = `uang.Money`, rasio = `uang.Ratio`
// (K-010/K-012), tanggal yang dikurangkan = `time.Time`.
//
// Nama medan = nama properti Pega apa adanya, termasuk salah ejanya
// (`RICommision`, `EndorsmentReason`), supaya ketertelusuran ke korpus tidak
// butuh kamus.
package models

import (
	"encoding/json"
	"time"

	"nusantarare/inti/backend/uang"
)

// NilaiOld - nilai literal penanda baris warisan lapis C.
//
// [terverifikasi] Tersimpan BERKUTIP di korpus:
// `<PropertiesValue>"old"</PropertiesValue>` (tujuh varian
// `Activity/SetOLDValueToEDMWork_<LOB>.xml`). Isi stringnya `old`.
const NilaiOld = "old"

// NilaiProrate - nilai `OfferFacIn.IsProRate` yang disetel varian FIRE.
//
// [terverifikasi] `Activity/SetOLDValueToEDMWork_FIRE.xml` langkah 1.1:
// `newWorkPage.OfferFacIn.IsProRate = "Prorate"`.
const NilaiProrate = "Prorate"

// NilaiShortPeriod - `OfferFacIn.IsProRate` yang membuka 1.3
// `CountPaymentEdmTSIObj_Act` ([terverifikasi] `.IsProRate=="ShortPeriod"`).
const NilaiShortPeriod = "ShortPeriod"

// NilaiFixRate - `ProtectSpreading.CARI2` yang menandai tarif tetap
// ([terverifikasi] `CountPaymentEdm_Act` langkah 9 menyetel `"FIX RATE"`).
const NilaiFixRate = "FIX RATE"

// FlagDihapus - `.FlagDelete` baris yang dihapus. [terverifikasi]
// `CalculatePremiFire` membandingkannya sebagai teks (`=="1"`, 7.x/10.x) DAN
// sebagai angka (`!=1`, 6.2.x).
const FlagDihapus = "1"

// FlagPolisBerjalan - `pyWorkPage.FlagOnGoingPolicy` polis yang periodenya
// sudah berjalan ([terverifikasi] deskripsi langkah 10 `CalculatePremiFire`:
// "loop saat periode polis sudah berjalan", prakondisi `==1`).
const FlagPolisBerjalan = "1"

// NilaiBenar - literal boolean berkutip `"true"` (`.IsAdjustableFlag`).
const NilaiBenar = "true"

// CalculateMethodProRata - `.CalculateMethod` yang disetel
// `SetLocalNonMbuProrate` langkah 1. [terverifikasi] deskripsi langkah 9.1:
// "Perhitungan dengan menggunakan PRO RATA (1)".
const CalculateMethodProRata = "1"

// StatusBusinessEDM - `QuotationData.StatusBusiness` kasus endorsement:
// `DataToEDM` menyetel `"3"`, dan `When/IsEDM` menguji `= 3`. Keluarga
// diskriminator siklus K-029: 1 = NB, 2 = RNW, 3 = EDM.
const StatusBusinessEDM = "3"

// StatusEDMPolisBatal - `OutputData1.pxResults(1).CARI20` (`GetEdmStatus`) yang
// membuat `SetValueToEDMWork` langkah 4-6 menolak: polis sudah dibatalkan
// lewat endorsement. Arti nilai lain kolom itu: OQ-020.
const StatusEDMPolisBatal = "1"

// JenisEndorsemen - `EdmType`, enumerasi TERKUNCI K-029. Nilai teks apa
// adanya dari Pega.
type JenisEndorsemen string

const (
	EdmBatalSejakSemula JenisEndorsemen = "1" // K-029
	EdmBatalProrata     JenisEndorsemen = "2" // K-029
	// EdmUsang3 - kode usang; cabangnya tetap diport apa adanya (K-029, K-046).
	EdmUsang3 JenisEndorsemen = "3"
	// EdmPerubahan - penambahan / pengurangan / perubahan (K-029).
	EdmPerubahan JenisEndorsemen = "4"
)

// JenisPenyesuaian - nilai `QuotationData.Type` sebuah kasus endorsement,
// apa adanya dari Pega (teks; arti tiap kode: OQ-020).
//
// ⛔ Tipe bernama, bukan `string`: ini BUKAN `Type` klaim maupun `Type` polis
// Life yang dijaga `claimlife` (`models/satutype_test.go`, satu rumah
// tersimpan untuk `Type`). Tipe tersendiri membuat ketiganya tidak dapat
// tertukar di kompilator.
type JenisPenyesuaian string

// KasusEndorsement - `newWorkPage`, objek kerja endorsement
// (kelas `ASM-FW-GISFW-Work-Endorsement`).
type KasusEndorsement struct {
	// Kelas, PrefiksID - kelas kerja dan awalan ID kasus (`SetValueToEDMWork`
	// langkah 7: `param.classname`, `param.IDPrefix`).
	Kelas, PrefiksID string
	// PxInsName, PzInsKey - identitas kasus (dari pembuat kasus, `svcAddWorkObject`).
	PxInsName, PzInsKey string
	// PyLabel - `= pxInsName` (`DataToEDM`).
	PyLabel string
	// EndorsementID - handle kasus PORTAL (`DataToEDM`: `= .pzInsKey` portal).
	EndorsementID string
	// PolicyNumber - `curWorkPage.PolicyNumber = .PolicyNo` (langkah 7).
	PolicyNumber string
	// Lima properti konektor `Start2 → Assignment7` flow
	// `InputAddendumFacultativeIn` (+ `NBStatusNew`).
	FlagOnGoingPolicy, IsCedingConfirm, Position, PositionNote, NBStatus, NBStatusNew string

	OfferFacIn OfferFacIn
	// Quotation - `newWorkPage.Quotation`, salinan halaman
	// `OfferFacIn.QuotationData` (SetValueToEDMWork 14.3).
	Quotation Quotation
	// Policy - `newWorkPage.Policy`; hanya `Payment` yang disentuh.
	Policy   Polis
	IsB2B    string
	PPnCheck string
	// TradingList, GoodsList, ConveyanceList - disalin varian Marine Cargo
	// dari `OldData.CargoList(1)`. Bentuk barisnya belum terverifikasi; dibawa
	// sebagai JSON mentah.
	TradingList    json.RawMessage
	GoodsList      json.RawMessage
	ConveyanceList json.RawMessage
}

// Polis - `newWorkPage.Policy`.
type Polis struct {
	// PolicyNo - `.Policy.PolicyNo = .PolicyNo` (`DataToEDM`).
	PolicyNo string
	Payment  Pembayaran
}

// OfferFacIn - agregat penawaran (kelas `ASM-FW-GISFW-Data-OfferFacIn`).
//
// `OldData` adalah LAPIS A: dokumen polis versi terakhir, kelasnya sama
// dengan `OfferFacIn` sendiri (`GetEDMOldData_SQL`). Ia pointer karena
// rekursif; nil berarti halaman itu belum pernah dimuat.
type OfferFacIn struct {
	QuotationData Quotation
	PolicyData    PolicyData

	Currency     string
	CurrencyList []BarisMataUang
	// Parameters - halaman yang bentuknya belum terverifikasi; JSON mentah.
	Parameters json.RawMessage

	InwardScale          string
	PercentShare         string
	AdditionalCapital    string
	MaxPctTreatyCapacity string
	MaxTreatyCapacity    string
	CurrentYear          string
	ProRatePercent       string
	ProRateType          string
	ShareCedantType      string
	IsSpecialAcceptance  string
	BinderRNM            string
	PPnCheck             string
	PolicyMasterNumber   string
	IsB2B                string

	CedingCedantList []Cedant

	LocationList []Lokasi
	CargoList    []Kargo
	VehicleList  []Kendaraan
	PersonList   []Orang

	// IsProRate - penanda tingkat agregat; HANYA varian FIRE yang
	// menyetelnya, dan ia TIDAK termasuk 54 salinan lapis A.
	IsProRate string

	// ProrateStartEDM, ProrateEDMEnd - porsi periode sebelum dan sesudah
	// tanggal endorsement (SetValueToEDMWork langkah 15).
	//
	// ⚠️ K-048: SEMENTARA. Modul selisih menghitung ulang dan menimpanya
	// sebelum jalur produksi membacanya. Konsumen mana pun dilarang
	// memperlakukannya sebagai final.
	ProrateStartEDM uang.Ratio
	ProrateEDMEnd   uang.Ratio

	EndorsmentReason string
	// IsBanding - `DataToEDM`: `.OfferFacIn.IsBanding = false` (teks apa adanya).
	IsBanding string

	OldData *OfferFacIn
}

// Quotation - kelas `ASM-FW-GISFW-Data-Quotation`, dipakai sebagai
// `OfferFacIn.QuotationData` dan `newWorkPage.Quotation`.
type Quotation struct {
	BusinessOldId    string
	BusinessCode     string
	BusinessFac      string
	BusinessName     string
	BusinessType     string
	CedingCo         string
	CedingCoName     string
	SourceOfBusiness string
	SobName          string
	MarketingCode    string
	MarketingName    string
	MOID             string
	TeamGroup        string
	InsuredName      string
	InsuredID        string
	NoOfferSlip      string
	IsGroup          string
	QQName           string
	PolicyType       string
	EDMDay           string

	// StatusBusiness - dibaca predikat `IsEDM` (`= 3`).
	StatusBusiness string
	// Type - `QuotationData.Type`, jenis penyesuaian endorsement; `"4"`
	// disebut "EDM ADJ SPREADING" di deskripsi langkah lapis C. Arti kode
	// lain: OQ-020.
	Type JenisPenyesuaian
	// EdmDate - tanggal endorsement, dibaca prorata langkah 15. Waktu nol =
	// tidak terisi.
	EdmDate time.Time
	// OldPolicyNo - kunci pemuatan lapis A (`GetEDMOldData_SQL`).
	OldPolicyNo string

	// Diisi `DataToEDM` dari kasus portal (arti kode: OQ-020; `EdmType`
	// terkunci K-029).
	EdmSource, EdmSourceNote, EdmNote string
	EdmType                           JenisEndorsemen
	// EdmTypeNew - pilihan `SetEdmType`: 1 penambahan, 2 pengurangan, 3
	// perubahan (deskripsi daftar pilihannya); BUKAN enumerasi K-029.
	EdmTypeNew, EdmStatus    string
	EdmChargeFee, EdmSurvey  string
	EndorsementInternalRetro string
}

// PolicyData - `OfferFacIn.PolicyData`.
type PolicyData struct {
	// StartDateTime, EndDateTime - periode polis; dikurangkan di prorata.
	// Waktu nol = tidak terisi.
	StartDateTime time.Time
	EndDateTime   time.Time
	OfferingDate  string
	ProdDateTime  string
	// EndorsementNo - dikosongkan varian Marine Cargo sesudah PolicyData
	// disalin utuh dari OldData.
	EndorsementNo string
	Payment       Pembayaran
}

// Pembayaran - `PolicyData.Payment` dan `newWorkPage.Policy.Payment`.
type Pembayaran struct {
	Installment     string
	RICommision     string
	PctBrokerageFee string
}

// BarisMataUang - satu baris kelas `ASM-FW-GISFW-Data-OfferFacIn-Currency`:
// `CurrencyList`, `TotalTSIList`, `TotalTSIPremiGrossList`, dan
// `CedingCedantList(i).CurrencyList`.
type BarisMataUang struct {
	TSI     uang.Money
	Premium uang.Money
	Rate    uang.Ratio

	TSIOld     uang.Money
	PremiumOld uang.Money
	RateOld    uang.Ratio
}

// Cedant - satu baris `CedingCedantList` (kelas `ASM-FW-GISFW-Data-Quotation`
// menurut `pyStepsClassName` SetOldData 4.2).
type Cedant struct {
	CurrencyList []BarisMataUang
}

// Lokasi - satu baris `LocationList`
// (kelas `ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance`).
type Lokasi struct {
	IsOldData string
	Property  PropertiLokasi
}

// PropertiLokasi - `LocationList(i).Property`.
type PropertiLokasi struct {
	PropertyItemList       []ItemProperti
	TotalTSIList           []BarisMataUang
	TotalTSIPremiGrossList []BarisMataUang
	RiskLocation           RisikoLokasi
}

// RisikoLokasi - `Property.RiskLocation`.
type RisikoLokasi struct {
	// AnekaList - langsung di bawah RiskLocation untuk lini Golf.
	AnekaList []Aneka
	// OccupationList - untuk lini Aneka.
	OccupationList []Okupasi
}

// ItemProperti - kelas `ASM-FW-GISFW-Data-PropertyItem` (lini kebakaran).
type ItemProperti struct {
	IsOldData string

	TSIObjectItem           uang.Money
	TotalGrossPremi         uang.Money
	TotalPremiumNusantaraRe uang.Money

	TSIObjectItemOld           uang.Money
	TotalGrossPremiOld         uang.Money
	TotalPremiumNusantaraReOld uang.Money

	CoverageList []Coverage
}

// Okupasi - kelas `ASM-FW-GISFW-Data-Occupation`.
type Okupasi struct {
	IsOldData string
	AnekaList []Aneka
}

// Aneka - kelas `ASM-FW-GISFW-Data-Aneka`.
type Aneka struct {
	IsOldData    string
	TSI          uang.Money
	TSIOld       uang.Money
	CoverageList []Coverage
}

// Kargo - kelas `ASM-FW-GISFW-Data-Cargo`.
type Kargo struct {
	FlagOldData            string
	CoverageList           []Coverage
	TotalTSIPremiGrossList []BarisMataUang
	TradingList            json.RawMessage
	GoodList               json.RawMessage
	ConveyanceList         json.RawMessage
}

// Kendaraan - kelas `ASM-FW-GISFW-Data-Vehicle`.
type Kendaraan struct {
	FlagOldData            string
	CoverageList           []Coverage
	TotalTSIPremiGrossList []BarisMataUang
}

// Orang - kelas `Data-Party-Person`: tertanggung PA, peserta Travel, dan
// tertanggung Life.
type Orang struct {
	FlagOldData string
	// ASMCoverage - coverage PA dan Travel.
	ASMCoverage []Coverage
	// CoverageList - coverage Life (varian `_LIFE`).
	CoverageList           []Coverage
	TotalTSIPremiGrossList []BarisMataUang
}

// Coverage - kelas `ASM-FW-GISFW-Data-Coverage`.
type Coverage struct {
	IsOldData string
	// EDM - dikosongkan varian FIRE (`.EDM = ""`); arti isinya belum
	// terverifikasi.
	EDM string

	TSI                       uang.Money
	Premium                   uang.Money
	PremiNusantaraRe          uang.Money
	PremiRp                   uang.Money
	PremiumGrossDiscountFleet uang.Money

	TSIOld                       uang.Money
	PremiumOld                   uang.Money
	PremiNusantaraReOld          uang.Money
	PremiRpOld                   uang.Money
	PremiumGrossDiscountFleetOld uang.Money

	SpreadingList []Spreading
	// AdditionalCoverage - coverage tambahan lini kendaraan bermotor.
	AdditionalCoverage []Coverage
}

// Spreading - kelas `ASM-FW-GISFW-Data-SpreadingRisk`.
//
// Hanya lima properti yang dimodelkan: kelimanya yang dibawa ke baris hasil
// penggabungan lapis C (`TempSpread.pxResults(<APPEND>)...`). Baris hasil
// penggabungan di Pega pun HANYA memuat kelima properti ini.
type Spreading struct {
	IsOldData       string
	TreatyType      string
	TSISpreaded     uang.Money
	PremiumSpreaded uang.Money
	SharePercentage uang.Ratio
}
