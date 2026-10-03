package models

import "github.com/cockroachdb/apd/v3"

// ObjekFire - satu baris tab Object kasus FIRE (tiket 35): .LocationList(n) beserta
// .Property, .Property.RiskLocation, .Property.BuildingConstruction, .Property.SurroundingRisk
// (tiket 38). Teks apa adanya;
// boolean Pega disimpan teks "true"/"false" (kolom IS_* VARCHAR2(10), bentuk data
// `[terverifikasi]` fixture).
type ObjekFire struct {
	ObjectNo         string // T_PROPERTY.OBJECT_NO          (.Property.ObjectNo, sel 5)
	ObjectType       string // T_PROPERTY.OBJECT_TYPE        (sel 6)
	ObjectName       string // T_PROPERTY.OBJECT_NAME        (sel 13)
	IsMaterialDamage bool   // T_PROPERTY.IS_MATERIAL_DAMAGE (sel 11)
	IsTopRisk        bool   // T_PROPERTY.IS_TOP_RISK        (sel 12)
	RoadType         string // T_PROPERTY.ROAD_TYPE          (sel 28 Type)
	RoadName         string // T_PROPERTY.ROAD_NAME          (sel 29 Address)
	BuildingNo       string // T_PROPERTY.BUILDING_NO        (sel 30)
	ZipCode          string // T_RISKLOCATION.ASM_ZIP_CODE   (sel 31)
	Country          string // T_PROPERTY.COUNTRY            (sel 32)
	RiskLocation     string // T_RISKLOCATION.ASM_ADDRESS    (sel 40 Risk Location)
	Territory        string // T_RISKLOCATION.ASMRW          (sel 35 Territory)
	City             string // T_RISKLOCATION.ASM_CITY       (sel 36)
	District         string // T_RISKLOCATION.ASM_DISTRICT   (sel 37)
	Province         string // T_PROPERTY.PROVINCE           (sel 38)
	RiskAddressID    string // T_PROPERTY.ALM_RISK_ID        (sel 39)
	NumberOfFloor    string // T_BUILDINGCONSTRUCTION.NUMBER_OF_FLOOR (sel 49; rancangan VARCHAR2(50))
	RoofType         string // T_BUILDINGCONSTRUCTION.ROOF_TYPE  (sel 50)
	WallType         string // T_BUILDINGCONSTRUCTION.WALL_TYPE  (sel 51)
	FloorType        string // T_BUILDINGCONSTRUCTION.FLOOR_TYPE (sel 52)
	PartitionType    string // T_BUILDINGCONSTRUCTION.PARTITION_TYPE    (sel 56, kolom baru 186)
	SupportWallType  string // T_BUILDINGCONSTRUCTION.SUPPORT_WALL_TYPE (sel 57, kolom baru 186)
	OtherType        string // T_BUILDINGCONSTRUCTION.OTHERS_TYPE       (sel 58, kolom baru 186)
	// Tiket 38 - `Section\RiskAround.xml` (migrasi 187).
	Ownership           string          // T_PROPERTY.OWNERSHIP                  (sel 65)
	IsProductionProcess bool            // T_PROPERTY.IS_PRODUCTION_PROCESS_FLAG (sel 77)
	IsHotWorkProcess    bool            // T_PROPERTY.IS_HOT_WORK_PROCESS_FLAG   (sel 78)
	IsFlammableItem     bool            // T_PROPERTY.IS_FLAMMABLE_ITEM_FLAG     (sel 79)
	SurroundingRisk     SurroundingRisk // T_SURROUNDINGRISK, satu baris per Property
	// Tiket 39 - .Property.PropertyItemList, T_PROPERTYITEMLIST urut SEQ_NO (migrasi 188).
	Items []ItemObjek
	// Tiket 40 - .Property.OccupationList (BUKAN .Property.RiskLocation.OccupationList), T_OCCUPATIONLIST urut
	// SEQ_NO + T_TABLEOFLIMIT (migrasi 189).
	Occupations []OkupasiObjek
	// Tiket 41 - .LocationList(n).FEAList (milik baris lokasi, bukan Property), T_FEALIST urut SEQ_NO (migrasi 190).
	FEA []BarisFEA
	// Tiket 42 - .Property.ListCauseOfLoss -> T_LISTCAUSEOFLOSS + T_COINSDATA urut SEQ_NO (migrasi 191); loss ratio
	// baris lokasi T_LOCATIONLIST.LOSS_RATIO* - BACA-SAJA bagi layar, DIHITUNG services saat simpan (W-4,
	// `DDL\SetLossRatio_Act.xml`; N-1 digantikan).
	LossRecords []CatatanKerugian
	LossRatio   LossRatio
	// Tiket 43 - .Property.TotalTSIPremiGrossList (SumTotalTSIPremiGross_Act cabang FIRE): DIHITUNG saat baca, tidak
	// disimpan.
	TotalPerCurrency []TotalCoverage
}

// TotalCoverage - satu mata uang di Total per objek: TSI = Σ TSIObjectItem, Premium = Σ TotalGrossPremi item,
// Rate = Premium / TSI × 1000 (TSI 0 -> 0).
type TotalCoverage struct {
	Currency           string
	TSI, Premium, Rate *apd.Decimal
}

// CoverageObjek - satu .Property.PropertyItemList(m).CoverageList(k) (kelas Data-Coverage, tiket 43) ->
// T_COVERAGELIST. Uang / rate / persen = *apd.Decimal (ADR-0034); nil = kosong. FirstLoss / FirstScale / LostLimit /
// EmlPml dihitung sebagai angka tetapi tersimpan teks di kolom rancangan VARCHAR2(50).
type CoverageObjek struct {
	Coverage, OldID, CoverageNote, CoverageBasis, Day, Indemnity, Conditions                           string
	TSI, Rate, RateOJK, FirstLoss, DiscountPercentage, TSILiability, NetRate, LimitOfLiability, PctLoL *apd.Decimal
	ProRatePercent, IndemnityPercentage, FirstScale, Sublimit, LostLimit, EmlPml, Discount, Premium    *apd.Decimal
	// PctAdjustment - .PctAdjustment hasil CountPremi langkah 11-15 (bukan medan kontrak; disimpan PCT_ADJUSTMENT).
	PctAdjustment *apd.Decimal
}

// CatatanKerugian - satu .Property.ListCauseOfLoss(n) (kelas Data-CauseOfLoss). Uang = *apd.Decimal, nil = kosong
// (ADR-0003/0034).
// DateOfLoss: DD-MM-YYYY di batas services <-> handler; teks Pega di batas services <-> repository.
type CatatanKerugian struct {
	DateOfLoss       string       // DATE_OF_LOSS       (.DateOfLoss, kolom baru 191)
	CoinsName        string       // T_COINSDATA.COINS_NAME (.CoinsData.CoinsName - kolom grid "Insured Name")
	LossObject       string       // LOSS_OBJECT        (kolom baru 191)
	Currency         string       // CURRENCY           (wajib, K-012)
	Amount           *apd.Decimal // AMOUNT             uang "Total of Loss" (kolom baru 191)
	Claim            *apd.Decimal // CLAIM              uang "Total Claim (100%)"
	PreventionOfLoss *apd.Decimal // PREVENTION_OF_LOSS uang (kolom baru 191)
	CauseOfLoss      string       // CAUSE_OF_LOSS      (kolom baru 191)
	Remarks          string       // REMARKS
	Detail           string       // DETAIL             "Loss Detail" (dilebarkan 500, A146)
}

// LossRatio - .LossRatio1Year* / .LossRatio35Year* baris lokasi; *apd.Decimal, nil = kosong (ADR-0034). Amount = uang
// NUMBER(38,8); percent tersimpan teks desimal di kolom rancangan VARCHAR2(50).
type LossRatio struct {
	OneYearAmount, OneYearPercent, ThreeFiveYearAmount, ThreeFiveYearPercent *apd.Decimal
}

// BarisFEA - satu .FEAList(n) (kelas Data-OfferFacIn-OfferFEAList); empat medan dari halaman tertanam .DataFEA.
// Jumlah unit = teks angka bulat >= 0 (M-2).
type BarisFEA struct {
	APAR                  string // APAR                     (.APAR)
	Sprinkler             string // SPRINKLER                (.Sprinkler)
	SmokeDetector         string // SMOKE_DETECTOR           (.SmokeDetector)
	Hydrant               string // HYDRANT                  (.Hydrant)
	PrivateTruckBrigade   string // PRIVATE_TRUCK_BRIGADE    (.DataFEA.PrivateTruckBrigade)
	PrivateFireBrigade    string // PRIVATE_FIRE_BRIGADE     (.DataFEA.PrivateFireBrigade)
	TeamSOPSafety         string // TEAM_SOP_SAFETY          (.DataFEA.TeamSOPSafety)
	TeamSOPRiskManagement string // TEAM_SOP_RISK_MANAGEMENT (.DataFEA.TeamSOPRiskManagement)
	Info                  string // INFO_FEA                 (.InfoFEA)
}

// OkupasiObjek - satu .Property.OccupationList(n) (kelas Data-Occupation) beserta halaman .TableOfLimit-nya.
type OkupasiObjek struct {
	OccupationID      string // T_OCCUPATIONLIST.OCCUPATION_ID  (.OccupationId = OCCUPATION.OLDID)
	OccupationName    string // T_OCCUPATIONLIST.OCCUPATION_NAME (.OccupationName = OCCUPATION.NAME)
	Category          string // T_TABLEOFLIMIT.CATEGORY    (.TableOfLimit.Category - I/II/III dari KDRiskExposure)
	ConstructionClass string // T_TABLEOFLIMIT.DESCRIPTION (.TableOfLimit.Description - Class of Construction)
	PctLimit          string // T_TABLEOFLIMIT.PCT_LIMIT   (.TableOfLimit.PctLimit - TEKS apa adanya, butir 68.1)
}

// ItemObjek - satu .Property.PropertyItemList(n) (tiket 39, `Section\PropertyItemFacIn_Section.xml`).
// Uang dan persen = desimal berskala tetap *apd.Decimal, nil = kosong (ADR-0003/0016/0034, butir 94): teks hanya di
// batas JSON (handler) dan bind/hasil Oracle (repository).
type ItemObjek struct {
	ItemTypeID     string       // ITEM_TYPE_ID       (.ItemTypeID = V_JN_OBJ_ITEM.MJOI_KODE)
	ItemType       string       // ITEM_TYPE          (.ItemType = V_JN_OBJ_ITEM.JN_OBJ_ITEM)
	Note           string       // PROPERTI_ITEM_NOTE (.PropertiItemNote = V_JN_OBJ_ITEM.KETERANGAN)
	PropertyYear   string       // PROPERTY_YEAR      (kolom baru 188)
	Unit           string       // UNIT               (kolom baru 188)
	Condition      string       // CONDITION          (kolom baru 188)
	Currency       string       // CURRENCY           (CURRENCY.CURRENCY; wajib, A133)
	TSI            *apd.Decimal // TSI_OBJECT_ITEM uang NUMBER(38,8)
	YearOfPlanting string       // YEAR               (.Year, kolom baru 188)
	NoOfTree       string       // NO_OF_TREE         (kolom baru 188)
	AreaHectar     string       // AREA_HECTAR        (kolom baru 188)
	Remark         string       // REMARK
	IsAdjustable   bool         // IS_ADJUSTABLE_FLAG teks "true"/"false"
	PctAdjust2     *apd.Decimal // PCT_ADJUST2 persen NUMBER(38,8)
	PctAdjustOther *apd.Decimal // PCT_ADJUST_OTHER persen NUMBER(38,8)
	// Tiket 43 - .CoverageList item (T_COVERAGELIST, induk T_PROPERTYITEMLIST lewat PARENT_TABLE) dan total item.
	Coverages       []CoverageObjek
	TotalGrossPremi *apd.Decimal // TOTAL_GROSS_PREMI = Σ Premium coverage (CountPremi langkah 55-56)
	TotalNetRate    *apd.Decimal // TOTAL_NET_RATE - dihitung tahap C2 (CalculateNetRate); tahap C1 diteruskan apa adanya
}

// BarisTableOfLimit - satu pilihan Class of Construction (tiket 40, TABLEOFLIMIT); PctLimit teks apa adanya.
type BarisTableOfLimit struct {
	Description string // TABLEOFLIMIT.DESCRIPTION
	PctLimit    string // TABLEOFLIMIT.PCTLIMIT
}

// BarisCoverage - satu pilihan coverage (tiket 43) dari tabel COVERAGE (popup dan otomatis).
type BarisCoverage struct {
	ID    string // ID -> .Coverage
	OldID string // OLDID - kode tampil
	Nama  string // COVERAGE.NAME -> .CoverageNote
}

// JenisItem - satu pilihan Object Item Type (tiket 39, V_JN_OBJ_ITEM).
type JenisItem struct {
	Kode       string // MJOI_KODE
	Nama       string // JN_OBJ_ITEM
	Keterangan string // KETERANGAN -> Object Item Note
}

// SurroundingRisk - .Property.SurroundingRisk (tiket 38). Teks apa adanya; tidak dicocokkan ke
// daftar pilihan.
type SurroundingRisk struct {
	Front, Left, Back, Right SisiRisiko // {FRONT|LEFT|BACK|RIGHT}_* (sel 9-14 / 23-28 / 37-42 / 51-56)
	HousekeepingStatus       string     // HOUSEKEEPING_STATUS (sel 66, kolom rancangan)
	FloodAreaStatus          string     // FLOOD_AREA_STATUS   (sel 67, kolom rancangan)
	FloodArea                string     // FLOOD_AREA          (sel 68, kolom baru 187)
	HousekeepingRemark       string     // HOUSEKEEPING_REMARK (sel 69, kolom baru 187)
}

// SisiRisiko - satu sisi .{Front|Left|Back|Right}{Occupation|Construction|Distance|Note}
// (kolom baru 187, A130).
type SisiRisiko struct {
	Occupation   string // OCCUPATION.OLDID terpilih (autocomplete menyimpan .OldID)
	Construction string
	Distance     string // teks angka 0..100000, <= 2 desimal (A129)
	Note         string // OCCUPATION.NAME saat Occupation dipilih
}

// BarisOccupation - satu okupasi FIRE: saran Surrounding Risk (tiket 38) dan popup Choose Occupation (tiket 40).
type BarisOccupation struct {
	OldID          string // OCCUPATION.OLDID - nilai yang disimpan
	Name           string // OCCUPATION.NAME
	KdRiskExposure string // OCCUPATION.KDRISKEXPOSURE - SetDataOccupation: 03 -> III, 02 -> II, 01 -> I (tiket 40)
}
