package handlers

// Tab Object FIRE (tiket 35, Surrounding Risk tiket 38): GET/PUT /api/nbfacin/kasus/{caseId}/objek - kontrak `ObjekFire`
// frontend (modul/nbfacin/frontend/api.ts).

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cockroachdb/apd/v3"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/inti/backend/utils"
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
	// tiket 40
	Occupations []okupasiKabel `json:"occupations"`
	// tiket 41
	FEA []feaKabel `json:"fea"`
	// tiket 42 - lossRatio dan internalLossRecords BACA-SAJA: dikirim GET, diabaikan PUT.
	LossRecords         []kerugianKabel `json:"lossRecords"`
	LossRatio           lossRatioKabel  `json:"lossRatio"`
	InternalLossRecords []klaimKabel    `json:"internalLossRecords"`
	// tiket 43 - BACA-SAJA: dihitung saat GET, diabaikan PUT.
	TotalPerCurrency []totalCoverageKabel `json:"totalPerCurrency"`
}

// kerugianKabel - kontrak `CatatanKerugian` frontend (tiket 42); uang teks desimal, dateOfLoss DD-MM-YYYY.
type kerugianKabel struct {
	DateOfLoss       string `json:"dateOfLoss"`
	CoinsName        string `json:"coinsName"`
	LossObject       string `json:"lossObject"`
	Currency         string `json:"currency"`
	Amount           string `json:"amount"`
	Claim            string `json:"claim"`
	PreventionOfLoss string `json:"preventionOfLoss"`
	CauseOfLoss      string `json:"causeOfLoss"`
	Remarks          string `json:"remarks"`
	Detail           string `json:"detail"`
}

// lossRatioKabel - kontrak `LossRatio` frontend.
type lossRatioKabel struct {
	OneYearAmount        string `json:"oneYearAmount"`
	OneYearPercent       string `json:"oneYearPercent"`
	ThreeFiveYearAmount  string `json:"threeFiveYearAmount"`
	ThreeFiveYearPercent string `json:"threeFiveYearPercent"`
}

// klaimKabel - kontrak `KlaimInternal` frontend. Selalu larik kosong (butir 83 / N-5): grid Pega tidak pernah terisi dan
// rancangan tidak punya jalur ListCauseOfLossClaim.
type klaimKabel struct {
	DateOfLoss    string `json:"dateOfLoss"`
	LocationNo    string `json:"locationNo"`
	Location      string `json:"location"`
	Currency      string `json:"currency"`
	Premium       string `json:"premium"`
	OSClaim       string `json:"osClaim"`
	AcceptedClaim string `json:"acceptedClaim"`
	IncurredClaim string `json:"incurredClaim"`
	LossRatio     string `json:"lossRatio"`
	Remark        string `json:"remark"`
}

// feaKabel - kontrak `BarisFEA` frontend (tiket 41).
type feaKabel struct {
	APAR                  string `json:"apar"`
	Sprinkler             string `json:"sprinkler"`
	SmokeDetector         string `json:"smokeDetector"`
	Hydrant               string `json:"hydrant"`
	PrivateTruckBrigade   string `json:"privateTruckBrigade"`
	PrivateFireBrigade    string `json:"privateFireBrigade"`
	TeamSOPSafety         string `json:"teamSopSafety"`
	TeamSOPRiskManagement string `json:"teamSopRiskManagement"`
	Info                  string `json:"info"`
}

// okupasiKabel - kontrak `OkupasiObjek` frontend (tiket 40).
type okupasiKabel struct {
	OccupationID      string `json:"occupationId"`
	OccupationName    string `json:"occupationName"`
	Category          string `json:"category"`
	ConstructionClass string `json:"constructionClass"`
	PctLimit          string `json:"pctLimit"`
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
	// tiket 43 - coverages selalu larik ke luar; totalGrossPremi BACA-SAJA (dihitung ulang saat PUT); totalNetRate
	// diterima dan disimpan apa adanya (A159).
	Coverages       []coverageKabel `json:"coverages"`
	TotalGrossPremi string          `json:"totalGrossPremi"`
	TotalNetRate    string          `json:"totalNetRate"`
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
		Items: keItemKabel(o.Items), Occupations: keOkupasiKabel(o.Occupations), FEA: keFEAKabel(o.FEA),
		LossRecords: keKerugianKabel(o.LossRecords), InternalLossRecords: []klaimKabel{},
		TotalPerCurrency: keTotalKabel(o.TotalPerCurrency),
		LossRatio: lossRatioKabel{OneYearAmount: teks(o.LossRatio.OneYearAmount), OneYearPercent: teks(o.LossRatio.OneYearPercent),
			ThreeFiveYearAmount: teks(o.LossRatio.ThreeFiveYearAmount), ThreeFiveYearPercent: teks(o.LossRatio.ThreeFiveYearPercent)}}
}

// teks - desimal -> teks kabel (ADR-0034: teks hanya di batas JSON); nil = "".
func teks(d *apd.Decimal) string { return utils.FormatDecimal(d) }

// keKerugianKabel / keKerugianModel - CatatanKerugian <-> kerugianKabel; selalu larik ke luar. Uang diurai SEKALI di
// sini lewat services.UraiDesimalIsian (pesan 400 ber-indeks ke masalah).
func keKerugianKabel(d []models.CatatanKerugian) []kerugianKabel {
	hasil := make([]kerugianKabel, 0, len(d))
	for _, c := range d {
		hasil = append(hasil, kerugianKabel{DateOfLoss: c.DateOfLoss, CoinsName: c.CoinsName, LossObject: c.LossObject,
			Currency: c.Currency, Amount: teks(c.Amount), Claim: teks(c.Claim), PreventionOfLoss: teks(c.PreventionOfLoss),
			CauseOfLoss: c.CauseOfLoss, Remarks: c.Remarks, Detail: c.Detail})
	}
	return hasil
}

func keKerugianModel(n int, d []kerugianKabel, masalah *[]string) []models.CatatanKerugian {
	hasil := make([]models.CatatanKerugian, 0, len(d))
	for m, c := range d {
		awal := fmt.Sprintf("baris[%d].lossRecords[%d].", n, m)
		hasil = append(hasil, models.CatatanKerugian{DateOfLoss: c.DateOfLoss, CoinsName: c.CoinsName, LossObject: c.LossObject,
			Currency: c.Currency, Amount: services.UraiDesimalIsian(awal+"amount", c.Amount, masalah),
			Claim:            services.UraiDesimalIsian(awal+"claim", c.Claim, masalah),
			PreventionOfLoss: services.UraiDesimalIsian(awal+"preventionOfLoss", c.PreventionOfLoss, masalah),
			CauseOfLoss:      c.CauseOfLoss, Remarks: c.Remarks, Detail: c.Detail})
	}
	return hasil
}

// keFEAKabel / keFEAModel - BarisFEA <-> feaKabel; selalu larik ke luar.
func keFEAKabel(d []models.BarisFEA) []feaKabel {
	hasil := make([]feaKabel, 0, len(d))
	for _, f := range d {
		hasil = append(hasil, feaKabel(f))
	}
	return hasil
}

func keFEAModel(d []feaKabel) []models.BarisFEA {
	hasil := make([]models.BarisFEA, 0, len(d))
	for _, f := range d {
		hasil = append(hasil, models.BarisFEA(f))
	}
	return hasil
}

// keOkupasiKabel / keOkupasiModel - OkupasiObjek <-> okupasiKabel; selalu larik ke luar.
func keOkupasiKabel(d []models.OkupasiObjek) []okupasiKabel {
	hasil := make([]okupasiKabel, 0, len(d))
	for _, o := range d {
		hasil = append(hasil, okupasiKabel(o))
	}
	return hasil
}

func keOkupasiModel(d []okupasiKabel) []models.OkupasiObjek {
	hasil := make([]models.OkupasiObjek, 0, len(d))
	for _, o := range d {
		hasil = append(hasil, models.OkupasiObjek(o))
	}
	return hasil
}

// keItemKabel / keItemModel - ItemObjek <-> itemKabel; selalu larik (bukan null) ke luar. Uang/persen diurai SEKALI di
// sini lewat services.UraiDesimalIsian.
func keItemKabel(d []models.ItemObjek) []itemKabel {
	hasil := make([]itemKabel, 0, len(d))
	for _, i := range d {
		hasil = append(hasil, itemKabel{ItemTypeID: i.ItemTypeID, ItemType: i.ItemType, Note: i.Note, PropertyYear: i.PropertyYear,
			Unit: i.Unit, Condition: i.Condition, Currency: i.Currency, TSI: teks(i.TSI), YearOfPlanting: i.YearOfPlanting,
			NoOfTree: i.NoOfTree, AreaHectar: i.AreaHectar, Remark: i.Remark, IsAdjustable: i.IsAdjustable,
			PctAdjust2: teks(i.PctAdjust2), PctAdjustOther: teks(i.PctAdjustOther), Coverages: keCoverageKabel(i.Coverages),
			TotalGrossPremi: teks(i.TotalGrossPremi), TotalNetRate: teks(i.TotalNetRate)})
	}
	return hasil
}

func keItemModel(n int, d []itemKabel, masalah *[]string) []models.ItemObjek {
	hasil := make([]models.ItemObjek, 0, len(d))
	for m, i := range d {
		awal := fmt.Sprintf("baris[%d].items[%d].", n, m)
		hasil = append(hasil, models.ItemObjek{ItemTypeID: i.ItemTypeID, ItemType: i.ItemType, Note: i.Note,
			PropertyYear: i.PropertyYear, Unit: i.Unit, Condition: i.Condition, Currency: i.Currency,
			TSI:            services.UraiDesimalIsian(awal+"tsi", i.TSI, masalah),
			YearOfPlanting: i.YearOfPlanting, NoOfTree: i.NoOfTree, AreaHectar: i.AreaHectar, Remark: i.Remark,
			IsAdjustable:   i.IsAdjustable,
			PctAdjust2:     services.UraiDesimalIsian(awal+"pctAdjust2", i.PctAdjust2, masalah),
			PctAdjustOther: services.UraiDesimalIsian(awal+"pctAdjustOther", i.PctAdjustOther, masalah),
			Coverages:      keCoverageModelDaftar(awal, i.Coverages, masalah),
			TotalNetRate:   services.UraiDesimalIsian(awal+"totalNetRate", i.TotalNetRate, masalah)})
	}
	return hasil
}

// keModel - objek ke-`n` badan PUT -> model; masalah urai desimal (ADR-0034) dikumpulkan ke `masalah`.
func keModel(n int, o objekKabel, masalah *[]string) models.ObjekFire {
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
		Items: keItemModel(n, o.Items, masalah), Occupations: keOkupasiModel(o.Occupations), FEA: keFEAModel(o.FEA),
		LossRecords: keKerugianModel(n, o.LossRecords, masalah)}
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
		var masalah []string
		for i, o := range b.Baris {
			baris[i] = keModel(i, o, &masalah)
		}
		if err := services.GalatIsianObjek(masalah); err != nil {
			tulisGalat(w, err)
			return
		}
		d, err := svc.GantiObjek(r.Context(), inti.PelakuDari(r, stubPelaku), r.PathValue("caseId"), baris)
		if err != nil {
			tulisGalat(w, err)
			return
		}
		galat.TulisJSON(w, keObjekKabel(d))
	}
}
