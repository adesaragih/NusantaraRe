package handlers

// Tab Object FIRE (tiket 35, Surrounding Risk tiket 38): GET/PUT /api/nbfacin/kasus/{caseId}/objek - kontrak `ObjekFire`
// frontend (modul/nbfacin/frontend/api.ts).

import (
	"encoding/json"
	"net/http"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/services"
)

type objekKabel struct {
	ObjectNo         string `json:"objectNo"`
	ObjectType       string `json:"objectType"`
	ObjectName       string `json:"objectName"`
	IsMaterialDamage bool   `json:"isMaterialDamage"`
	IsTopRisk        bool   `json:"isTopRisk"`
	RoadType         string `json:"roadType"`
	RoadName         string `json:"roadName"`
	BuildingNo       string `json:"buildingNo"`
	ZipCode          string `json:"zipCode"`
	Country          string `json:"country"`
	RiskLocation     string `json:"riskLocation"`
	Territory        string `json:"territory"`
	City             string `json:"city"`
	District         string `json:"district"`
	Province         string `json:"province"`
	RiskAddressID    string `json:"riskAddressId"`
	NumberOfFloor    string `json:"numberOfFloor"`
	RoofType         string `json:"roofType"`
	WallType         string `json:"wallType"`
	FloorType        string `json:"floorType"`
	PartitionType    string `json:"partitionType"`
	SupportWallType  string `json:"supportWallType"`
	OtherType        string `json:"otherType"`
	// tiket 38
	Ownership           string       `json:"ownership"`
	IsProductionProcess bool         `json:"isProductionProcess"`
	IsHotWorkProcess    bool         `json:"isHotWorkProcess"`
	IsFlammableItem     bool         `json:"isFlammableItem"`
	SurroundingRisk     sekitarKabel `json:"surroundingRisk"`
	// tiket 39
	Items []itemKabel `json:"items"`
}

// itemKabel - kontrak `ItemObjek` frontend (tiket 39); uang/persen = teks desimal.
type itemKabel struct {
	ItemTypeID     string `json:"itemTypeId"`
	ItemType       string `json:"itemType"`
	Note           string `json:"note"`
	PropertyYear   string `json:"propertyYear"`
	Unit           string `json:"unit"`
	Condition      string `json:"condition"`
	Currency       string `json:"currency"`
	TSI            string `json:"tsi"`
	YearOfPlanting string `json:"yearOfPlanting"`
	NoOfTree       string `json:"noOfTree"`
	AreaHectar     string `json:"areaHectar"`
	Remark         string `json:"remark"`
	IsAdjustable   bool   `json:"isAdjustable"`
	PctAdjust2     string `json:"pctAdjust2"`
	PctAdjustOther string `json:"pctAdjustOther"`
}

// sekitarKabel - kontrak `SurroundingRisk` frontend (tiket 38).
type sekitarKabel struct {
	Front              sisiKabel `json:"front"`
	Left               sisiKabel `json:"left"`
	Back               sisiKabel `json:"back"`
	Right              sisiKabel `json:"right"`
	HousekeepingStatus string    `json:"housekeepingStatus"`
	FloodAreaStatus    string    `json:"floodAreaStatus"`
	FloodArea          string    `json:"floodArea"`
	HousekeepingRemark string    `json:"housekeepingRemark"`
}

// sisiKabel - kontrak `SisiRisiko` frontend.
type sisiKabel struct {
	Occupation   string `json:"occupation"`
	Construction string `json:"construction"`
	Distance     string `json:"distance"`
	Note         string `json:"note"`
}

type daftarObjek struct {
	Baris []objekKabel `json:"baris"`
}

// batasBadanObjek - badan PUT terbesar yang diterima (A112).
const batasBadanObjek = 1 << 20

// keKabel / keModel - pemetaan medan EKSPLISIT (struct bersarang tidak dapat dikonversi langsung).
func keKabel(o models.ObjekFire) objekKabel {
	s := o.SurroundingRisk
	return objekKabel{ObjectNo: o.ObjectNo, ObjectType: o.ObjectType, ObjectName: o.ObjectName,
		IsMaterialDamage: o.IsMaterialDamage, IsTopRisk: o.IsTopRisk, RoadType: o.RoadType, RoadName: o.RoadName,
		BuildingNo: o.BuildingNo, ZipCode: o.ZipCode, Country: o.Country, RiskLocation: o.RiskLocation,
		Territory: o.Territory, City: o.City, District: o.District, Province: o.Province, RiskAddressID: o.RiskAddressID,
		NumberOfFloor: o.NumberOfFloor, RoofType: o.RoofType, WallType: o.WallType, FloorType: o.FloorType,
		PartitionType: o.PartitionType, SupportWallType: o.SupportWallType, OtherType: o.OtherType,
		Ownership: o.Ownership, IsProductionProcess: o.IsProductionProcess, IsHotWorkProcess: o.IsHotWorkProcess,
		IsFlammableItem: o.IsFlammableItem,
		SurroundingRisk: sekitarKabel{Front: sisiKabel(s.Front), Left: sisiKabel(s.Left), Back: sisiKabel(s.Back),
			Right: sisiKabel(s.Right), HousekeepingStatus: s.HousekeepingStatus, FloodAreaStatus: s.FloodAreaStatus,
			FloodArea: s.FloodArea, HousekeepingRemark: s.HousekeepingRemark},
		Items: keItemKabel(o.Items)}
}

// keItemKabel / keItemModel - ItemObjek <-> itemKabel; selalu larik (bukan null) ke luar.
func keItemKabel(d []models.ItemObjek) []itemKabel {
	hasil := make([]itemKabel, 0, len(d))
	for _, i := range d {
		hasil = append(hasil, itemKabel(i))
	}
	return hasil
}

func keItemModel(d []itemKabel) []models.ItemObjek {
	hasil := make([]models.ItemObjek, 0, len(d))
	for _, i := range d {
		hasil = append(hasil, models.ItemObjek(i))
	}
	return hasil
}

func keModel(o objekKabel) models.ObjekFire {
	s := o.SurroundingRisk
	return models.ObjekFire{ObjectNo: o.ObjectNo, ObjectType: o.ObjectType, ObjectName: o.ObjectName,
		IsMaterialDamage: o.IsMaterialDamage, IsTopRisk: o.IsTopRisk, RoadType: o.RoadType, RoadName: o.RoadName,
		BuildingNo: o.BuildingNo, ZipCode: o.ZipCode, Country: o.Country, RiskLocation: o.RiskLocation,
		Territory: o.Territory, City: o.City, District: o.District, Province: o.Province, RiskAddressID: o.RiskAddressID,
		NumberOfFloor: o.NumberOfFloor, RoofType: o.RoofType, WallType: o.WallType, FloorType: o.FloorType,
		PartitionType: o.PartitionType, SupportWallType: o.SupportWallType, OtherType: o.OtherType,
		Ownership: o.Ownership, IsProductionProcess: o.IsProductionProcess, IsHotWorkProcess: o.IsHotWorkProcess,
		IsFlammableItem: o.IsFlammableItem,
		SurroundingRisk: models.SurroundingRisk{Front: models.SisiRisiko(s.Front), Left: models.SisiRisiko(s.Left),
			Back: models.SisiRisiko(s.Back), Right: models.SisiRisiko(s.Right), HousekeepingStatus: s.HousekeepingStatus,
			FloodAreaStatus: s.FloodAreaStatus, FloodArea: s.FloodArea, HousekeepingRemark: s.HousekeepingRemark},
		Items: keItemModel(o.Items)}
}

func keObjekKabel(d []models.ObjekFire) daftarObjek {
	baris := make([]objekKabel, 0, len(d))
	for _, o := range d {
		baris = append(baris, keKabel(o))
	}
	return daftarObjek{Baris: baris}
}

// bacaObjek - GET /api/nbfacin/kasus/{caseId}/objek.
func bacaObjek(svc *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		d, err := svc.BacaObjek(r.Context(), r.PathValue("caseId"))
		if err != nil {
			tulisGalat(w, err)
			return
		}
		galat.TulisJSON(w, keObjekKabel(d))
	}
}

// simpanObjek - PUT /api/nbfacin/kasus/{caseId}/objek (tombol Save): ganti utuh, jawab baca ulang.
func simpanObjek(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var b daftarObjek
		if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, batasBadanObjek)).Decode(&b); err != nil {
			galat.Tulis(w, http.StatusBadRequest, "badan permintaan bukan JSON daftar objek yang sah")
			return
		}
		baris := make([]models.ObjekFire, len(b.Baris))
		for i, o := range b.Baris {
			baris[i] = keModel(o)
		}
		d, err := svc.GantiObjek(r.Context(), inti.PelakuDari(r, stubPelaku), r.PathValue("caseId"), baris)
		if err != nil {
			tulisGalat(w, err)
			return
		}
		galat.TulisJSON(w, keObjekKabel(d))
	}
}
