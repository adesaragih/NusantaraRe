package models

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

// BarisOccupation - satu saran Occupation Surrounding Risk (tiket 38, OCCUPATION TYPE FIRE).
type BarisOccupation struct {
	OldID string // OCCUPATION.OLDID - nilai yang disimpan
	Name  string // OCCUPATION.NAME
}
