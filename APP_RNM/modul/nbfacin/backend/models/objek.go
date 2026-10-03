package models

// ObjekFire - satu baris tab Object kasus FIRE (tiket 35): .LocationList(n) beserta
// .Property, .Property.RiskLocation, .Property.BuildingConstruction. Teks apa adanya;
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
}
