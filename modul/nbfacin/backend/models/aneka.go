package models

import "github.com/cockroachdb/apd/v3"

// Tab Object keluarga Aneka (tiket 73 P1, migrasi 214). Jalur Pega `[terverifikasi]`: non-Golf
// `.OfferFacIn.LocationList(n).Property.RiskLocation.OccupationList(m).AnekaList(k)`, Golf
// `.OfferFacIn.LocationList(n).Property.RiskLocation.AnekaList(k)` (riset form tiket 73 bagian D.1). Medan form
// `NB FacIn\Section\InputDtlObjectAneka_FacIn.xml`; uang / angka *apd.Decimal (ADR-0034), tanggal teks Pega di sini.

// LokasiAneka - satu lokasi (T_LOCATIONLIST + T_PROPERTY + T_RISKLOCATION, baris yang sama dengan tab Object FIRE) dengan
// okupasi Aneka dan objek Golf.
type LokasiAneka struct {
	ObjectNo      string // T_PROPERTY.OBJECT_NO
	IsTopRisk     bool   // T_PROPERTY.IS_TOP_RISK (teks "true" / "false")
	RoadType      string // T_PROPERTY.ROAD_TYPE
	RoadName      string // T_PROPERTY.ROAD_NAME
	BuildingNo    string // T_PROPERTY.BUILDING_NO
	ZipCode       string // T_RISKLOCATION.ASM_ZIP_CODE
	Country       string // T_PROPERTY.COUNTRY
	RiskLocation  string // T_RISKLOCATION.ASM_ADDRESS
	Territory     string // T_RISKLOCATION.ASMRW
	City          string // T_RISKLOCATION.ASM_CITY
	District      string // T_RISKLOCATION.ASM_DISTRICT
	Province      string // T_PROPERTY.PROVINCE
	RiskAddressID string // T_PROPERTY.ALM_RISK_ID
	Okupasi       []OkupasiAneka
	ObjekGolf     []ObjekAneka // `.Property.RiskLocation.AnekaList` (T_ANEKALIST ber-PARENT_TABLE T_RISKLOCATION)
	// TotalPerCurrency - `.Property.TotalTSIPremiGrossList` (tiket 73 P2), dihitung saat baca, tidak disimpan.
	TotalPerCurrency []TotalCoverage
}

// OkupasiAneka - `.Property.RiskLocation.OccupationList(m)` (T_OCCUPATIONLIST ber-PARENT_TABLE T_RISKLOCATION).
type OkupasiAneka struct {
	OccupationID   string
	OccupationName string
	Objek          []ObjekAneka // `.AnekaList` (T_ANEKALIST ber-PARENT_TABLE T_OCCUPATIONLIST)
}

// ObjekAneka - satu baris `.AnekaList` (T_ANEKALIST) dan empat halaman anaknya.
type ObjekAneka struct {
	ObjectName         string
	Section            string
	NoKontrak          string
	StartDate          string // teks Pega (START_DATE)
	EndDate            string // teks Pega (END_DATE)
	NoKtp              string
	NoNpwp             string
	Currency           string // `.Currency.Name` -> CURRENCY_CODE
	TSI                *apd.Decimal
	TSIOld             *apd.Decimal
	DraftWordingID     string
	MaintenancePeriod  *apd.Decimal
	Others             string
	Year               string
	Quantity           string
	UnitQuantity       string
	SerialNo           string
	BuiltIn            string
	Flag               string
	GeographicalLimits string
	InvoiceNo          string
	NoOfTree           *apd.Decimal
	AreaHectar         *apd.Decimal
	LimitIdemnityDay   *apd.Decimal
	LimitIdemnityMonth *apd.Decimal
	LimitOfLiability   *apd.Decimal
	AlmRiskID          string
	Briguna            BrigunaAneka
	MarineHull         MarineHull
	AviationHull       AviationHull
	VehicleHE          VehicleHE
	Memorandum         MemorandumAccept
	// Coverages - `.CoverageList` objek (tiket 73 P2, T_COVERAGELIST ber-PARENT_TABLE T_ANEKALIST), urut SEQ_NO. PUT: nil = kunci
	// `coverages` ABSEN (coverage tersimpan dipertahankan per posisi); larik (juga kosong) = ganti utuh. GET: selalu larik.
	Coverages []CoverageAneka
}

// CoverageAneka - satu `.CoverageList(k)` objek Aneka (kelas Data-Coverage, form `NB FacIn\Section\InputCoverageAneka_FacIn.xml`).
// Uang / persen / hari angka = *apd.Decimal (nil = kosong); kode = teks apa adanya.
type CoverageAneka struct {
	Coverage, OldID, CoverageNote string
	// CalculateMethodFacIn - '1' PRO RATA / '2' SHORT PERIOD / '3' FIX RATE (`DDL\CalculateMethod_FacIn.xml`); CalculateMethod -
	// `.CalculateMethod` (dibaca ChangeCoveragePremium).
	CalculateMethodFacIn, CalculateMethod string
	ProRatePercent, PctShortPeriod        *apd.Decimal
	Day                                   string // `.Day` 365 / 366 / 360 (teks, kolom DAY VARCHAR2 rancangan)
	Currency                              string // CURRENCY_CODE
	TSI, Loading, LostLimit, Rate         *apd.Decimal
	Indemnity, Unit                       string
	IndemnityPercentage                   *apd.Decimal
	Premium, DiscountPercentage, Discount *apd.Decimal
	LimitOfLiability, PctLoL              *apd.Decimal
	SubLimitNote                          string
	// Dihitung server (diabaikan saat PUT).
	MinPremium, TSILiability, TSINusantaraRe, PremiNusantaraRe *apd.Decimal
	Deductibles                                                []Deductible
	// Spreading - `.SpreadingList` coverage (tiket 73 P3, T_SPREADINGLIST berinduk coverage). Internal: tidak dikirim kabel
	// objek-aneka; Save tab Object / Coverage mempertahankannya per posisi (pola K48-7), tab Spreading Aneka menulisnya.
	Spreading []BarisSpreading
}

// BrigunaAneka - blok `IsCreditBriguna` (kolom T_ANEKALIST, K73-7).
type BrigunaAneka struct {
	LoanAccountNumber string
	LoanName          string
	Job               string
	DepositAccountNo  string
	IDCardNo          string
	Institution       string
	LoanType          string
	DateOfBirth       string // teks Pega
	Salary            *apd.Decimal
}

// MarineHull - `.MarineHull` (T_MARINEHULL).
type MarineHull struct {
	TypeOfVessel, NRTRT, LengthM, Classification, GRTRT, BreadthM, Construction, DWTRT, DepthM string
}

// AviationHull - `.AviationHull` (T_AVIATIONHULL).
type AviationHull struct {
	RegistrationMarks, MaxPassengers, Pilots, StandardUses, SpecialUses, SpecialRentalUses string
}

// VehicleHE - `.VehicleHE` (T_VEHICLEHE); ObjectIDHE / Brand / Type = ID terpilih lookup.
type VehicleHE struct {
	ObjectNameHE, ObjectIDHE, BrandName, Brand, TypeName, Type, ChassisNumber, EngineNumber string
}

// MemorandumAccept - `.MemorandumAccept` (T_MEMORANDUMACCEPT); SuretyType diisi server (InputDtlObjectAneka_PreAct).
type MemorandumAccept struct {
	SuretyType, ApplicationNumber, Accumulation, Collateral, InsuredSiup, KeyPerson, OriginalRate, PrincipalsOutGo,
	TypeOfWording, BuildingStatus, SupportingDocument, PrincipalsData, ConclusionsProposals, PrincipalsBackground,
	SummaryFinancial, PrincipalsGuarantee string
}

// Kosong - tidak satu medan pun terisi (halaman anak tidak ditulis bila kosong).
func (m MarineHull) Kosong() bool { return m == MarineHull{} }

// Kosong - lihat MarineHull.Kosong.
func (a AviationHull) Kosong() bool { return a == AviationHull{} }

// Kosong - lihat MarineHull.Kosong.
func (v VehicleHE) Kosong() bool { return v == VehicleHE{} }

// Kosong - lihat MarineHull.Kosong.
func (m MemorandumAccept) Kosong() bool { return m == MemorandumAccept{} }

// BarisKapal - satu pilihan Vessel Name (lookup `BrowseShip_RD`, POOLDATA.SHIP); semua kolom sumber VARCHAR2, teks apa adanya.
// Pemetaan ke objek saat dipilih (`InputDtlObjectAneka_FacIn.xml` L10490-10564): ID -> AlmRiskID, Nama -> ObjectName, Flag -> Flag,
// Tahun (YEARMAKE1) -> Year, NRT / GRT / DWT / Klasifikasi / Konstruksi -> MarineHull.NRTRT / GRTRT / DWTRT / Classification /
// Construction.
type BarisKapal struct {
	ID, Nama, Flag, Tahun, NRT, GRT, DWT, Klasifikasi, Konstruksi string
}

// BarisTipeKendaraan - satu pilihan Type HE (lookup `BrowseType_RD`, view POOLDATA.BRANDDETAIL): ID -> VehicleHE.Type,
// TypeName -> VehicleHE.TypeName.
type BarisTipeKendaraan struct {
	ID, TypeName string
}

// PilihanOkupasiAneka - satu baris popup Choose Occupation Aneka (`BrowseOccupationANEKA_RD`, POOLDATA.OCCUPATION): ID ->
// OccupationId, Name -> OccupationName (`GetOccupation_act` Property-Set `.OccupationId = Param.ID`, `.OccupationName = Param.Name`).
type PilihanOkupasiAneka struct {
	ID, Name string
}
