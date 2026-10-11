package handlers

// Tab Object keluarga Aneka (tiket 73 P1): GET/PUT /api/nbfacin/kasus/{caseId}/objek-aneka - kontrak `LokasiAneka` frontend
// (modul/nbfacin/frontend/api.ts). Uang / angka teks desimal (ADR-0003), tanggal DD-MM-YYYY.

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/cockroachdb/apd/v3"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/services"
)

type daftarAneka struct {
	Lokasi []lokasiAnekaKabel `json:"lokasi"`
	// ProRatePercent - % Pro Rate tingkat CASE (`pyWorkPage.OfferFacIn.ProRatePercent`, tiket 73 P2, K73-12) =
	// T_GENERAL_POLIS.PRO_RATE_PERCENT (migrasi 255). GET: tersimpan, NULL -> dari tanggal polis (CPCA langkah 4, Day 365);
	// PUT: ditulis (kosong = NULL).
	ProRatePercent string `json:"proRatePercent"`
}

type lokasiAnekaKabel struct {
	ObjectNo      string              `json:"objectNo"`
	IsTopRisk     bool                `json:"isTopRisk"`
	RoadType      string              `json:"roadType"`
	RoadName      string              `json:"roadName"`
	BuildingNo    string              `json:"buildingNo"`
	ZipCode       string              `json:"zipCode"`
	Country       string              `json:"country"`
	RiskLocation  string              `json:"riskLocation"`
	Territory     string              `json:"territory"`
	City          string              `json:"city"`
	District      string              `json:"district"`
	Province      string              `json:"province"`
	RiskAddressID string              `json:"riskAddressId"`
	Okupasi       []okupasiAnekaKabel `json:"okupasi"`
	ObjekGolf     []objekAnekaKabel   `json:"objekGolf"`
	// TotalPerCurrency - tiket 73 P2, BACA-SAJA (diabaikan PUT).
	TotalPerCurrency []totalCoverageKabel `json:"totalPerCurrency"`
}

type okupasiAnekaKabel struct {
	OccupationID   string            `json:"occupationId"`
	OccupationName string            `json:"occupationName"`
	Objek          []objekAnekaKabel `json:"objek"`
}

type objekAnekaKabel struct {
	ObjectName         string            `json:"objectName"`
	Section            string            `json:"section"`
	NoKontrak          string            `json:"noKontrak"`
	StartDate          string            `json:"startDate"`
	EndDate            string            `json:"endDate"`
	NoKtp              string            `json:"noKtp"`
	NoNpwp             string            `json:"noNpwp"`
	Currency           string            `json:"currency"`
	TSI                string            `json:"tsi"`
	TSIOld             string            `json:"tsiOld"`
	DraftWordingID     string            `json:"draftWordingId"`
	MaintenancePeriod  string            `json:"maintenancePeriod"`
	Others             string            `json:"others"`
	Year               string            `json:"year"`
	Quantity           string            `json:"quantity"`
	UnitQuantity       string            `json:"unitQuantity"`
	SerialNo           string            `json:"serialNo"`
	BuiltIn            string            `json:"builtIn"`
	Flag               string            `json:"flag"`
	GeographicalLimits string            `json:"geographicalLimits"`
	InvoiceNo          string            `json:"invoiceNo"`
	NoOfTree           string            `json:"noOfTree"`
	AreaHectar         string            `json:"areaHectar"`
	LimitIdemnityDay   string            `json:"limitIdemnityDay"`
	LimitIdemnityMonth string            `json:"limitIdemnityMonth"`
	LimitOfLiability   string            `json:"limitOfLiability"`
	AlmRiskID          string            `json:"almRiskId"`
	Briguna            brigunaKabel      `json:"briguna"`
	MarineHull         marineHullKabel   `json:"marineHull"`
	AviationHull       aviationHullKabel `json:"aviationHull"`
	VehicleHE          vehicleHEKabel    `json:"vehicleHe"`
	Memorandum         memorandumKabel   `json:"memorandum"`
	// Coverages - tiket 73 P2. PUT: kunci absen / null = coverage tersimpan dipertahankan; larik = ganti utuh.
	Coverages []coverageAnekaKabel `json:"coverages"`
}

// coverageAnekaKabel - kontrak `CoverageAneka` (tiket 73 P2); angka teks desimal. oldId, minPremium, tsiLiability, tsiNusantaraRe,
// premiNusantaraRe dihitung / diisi server.
type coverageAnekaKabel struct {
	Coverage             string            `json:"coverage"`
	OldID                string            `json:"oldId"`
	CoverageNote         string            `json:"coverageNote"`
	CalculateMethodFacIn string            `json:"calculateMethodFacIn"`
	CalculateMethod      string            `json:"calculateMethod"`
	ProRatePercent       string            `json:"proRatePercent"`
	PctShortPeriod       string            `json:"pctShortPeriod"`
	Day                  string            `json:"day"`
	Currency             string            `json:"currency"`
	TSI                  string            `json:"tsi"`
	Loading              string            `json:"loading"`
	LostLimit            string            `json:"lostLimit"`
	Rate                 string            `json:"rate"`
	Indemnity            string            `json:"indemnity"`
	Unit                 string            `json:"unit"`
	IndemnityPercentage  string            `json:"indemnityPercentage"`
	Premium              string            `json:"premium"`
	DiscountPercentage   string            `json:"discountPercentage"`
	Discount             string            `json:"discount"`
	LimitOfLiability     string            `json:"limitOfLiability"`
	PctLoL               string            `json:"pctLol"`
	SubLimitNote         string            `json:"subLimitNote"`
	MinPremium           string            `json:"minPremium"`
	TSILiability         string            `json:"tsiLiability"`
	TSINusantaraRe       string            `json:"tsiNusantaraRe"`
	PremiNusantaraRe     string            `json:"premiNusantaraRe"`
	Deductibles          []deductibleKabel `json:"deductibles"`
}

// desimalCoverageAneka - medan desimal coverage (kabel <-> model), satu daftar untuk kedua arah; `server` = diabaikan saat PUT.
var desimalCoverageAneka = []struct {
	jalur  string
	server bool
	kabel  func(*coverageAnekaKabel) *string
	model  func(*models.CoverageAneka) **apd.Decimal
}{
	{"proRatePercent", false, func(k *coverageAnekaKabel) *string { return &k.ProRatePercent },
		func(c *models.CoverageAneka) **apd.Decimal { return &c.ProRatePercent }},
	{"pctShortPeriod", false, func(k *coverageAnekaKabel) *string { return &k.PctShortPeriod },
		func(c *models.CoverageAneka) **apd.Decimal { return &c.PctShortPeriod }},
	{"tsi", false, func(k *coverageAnekaKabel) *string { return &k.TSI }, func(c *models.CoverageAneka) **apd.Decimal { return &c.TSI }},
	{"loading", false, func(k *coverageAnekaKabel) *string { return &k.Loading },
		func(c *models.CoverageAneka) **apd.Decimal { return &c.Loading }},
	{"lostLimit", false, func(k *coverageAnekaKabel) *string { return &k.LostLimit },
		func(c *models.CoverageAneka) **apd.Decimal { return &c.LostLimit }},
	{"rate", false, func(k *coverageAnekaKabel) *string { return &k.Rate }, func(c *models.CoverageAneka) **apd.Decimal { return &c.Rate }},
	{"indemnityPercentage", false, func(k *coverageAnekaKabel) *string { return &k.IndemnityPercentage },
		func(c *models.CoverageAneka) **apd.Decimal { return &c.IndemnityPercentage }},
	{"premium", false, func(k *coverageAnekaKabel) *string { return &k.Premium },
		func(c *models.CoverageAneka) **apd.Decimal { return &c.Premium }},
	{"discountPercentage", false, func(k *coverageAnekaKabel) *string { return &k.DiscountPercentage },
		func(c *models.CoverageAneka) **apd.Decimal { return &c.DiscountPercentage }},
	{"discount", false, func(k *coverageAnekaKabel) *string { return &k.Discount },
		func(c *models.CoverageAneka) **apd.Decimal { return &c.Discount }},
	{"limitOfLiability", false, func(k *coverageAnekaKabel) *string { return &k.LimitOfLiability },
		func(c *models.CoverageAneka) **apd.Decimal { return &c.LimitOfLiability }},
	{"pctLol", false, func(k *coverageAnekaKabel) *string { return &k.PctLoL },
		func(c *models.CoverageAneka) **apd.Decimal { return &c.PctLoL }},
	{"minPremium", true, func(k *coverageAnekaKabel) *string { return &k.MinPremium },
		func(c *models.CoverageAneka) **apd.Decimal { return &c.MinPremium }},
	{"tsiLiability", true, func(k *coverageAnekaKabel) *string { return &k.TSILiability },
		func(c *models.CoverageAneka) **apd.Decimal { return &c.TSILiability }},
	{"tsiNusantaraRe", true, func(k *coverageAnekaKabel) *string { return &k.TSINusantaraRe },
		func(c *models.CoverageAneka) **apd.Decimal { return &c.TSINusantaraRe }},
	{"premiNusantaraRe", true, func(k *coverageAnekaKabel) *string { return &k.PremiNusantaraRe },
		func(c *models.CoverageAneka) **apd.Decimal { return &c.PremiNusantaraRe }},
}

func keCoverageAnekaKabel(c models.CoverageAneka) coverageAnekaKabel {
	k := coverageAnekaKabel{Coverage: c.Coverage, OldID: c.OldID, CoverageNote: c.CoverageNote,
		CalculateMethodFacIn: c.CalculateMethodFacIn, CalculateMethod: c.CalculateMethod, Day: c.Day, Currency: c.Currency,
		Indemnity: c.Indemnity, Unit: c.Unit, SubLimitNote: c.SubLimitNote, Deductibles: keDeductibleKabel(c.Deductibles)}
	for _, d := range desimalCoverageAneka {
		*d.kabel(&k) = teks(*d.model(&c))
	}
	return k
}

// keCoverageAnekaModel - medan server tidak diurai (diabaikan); `jalur` = "...objek[n].coverages[k]".
func keCoverageAnekaModel(jalur string, k coverageAnekaKabel, masalah *[]string) models.CoverageAneka {
	c := models.CoverageAneka{Coverage: k.Coverage, OldID: k.OldID, CoverageNote: k.CoverageNote,
		CalculateMethodFacIn: k.CalculateMethodFacIn, CalculateMethod: k.CalculateMethod, Day: k.Day, Currency: k.Currency,
		Indemnity: k.Indemnity, Unit: k.Unit, SubLimitNote: k.SubLimitNote,
		Deductibles: keDeductibleModel(jalur+".", k.Deductibles, masalah)}
	for _, d := range desimalCoverageAneka {
		if !d.server {
			*d.model(&c) = services.UraiDesimalIsian(jalur+"."+d.jalur, *d.kabel(&k), masalah)
		}
	}
	return c
}

type brigunaKabel struct {
	LoanAccountNumber string `json:"loanAccountNumber"`
	LoanName          string `json:"loanName"`
	Job               string `json:"job"`
	DepositAccountNo  string `json:"depositAccountNo"`
	IDCardNo          string `json:"idCardNo"`
	Institution       string `json:"institution"`
	LoanType          string `json:"loanType"`
	DateOfBirth       string `json:"dateOfBirth"`
	Salary            string `json:"salary"`
}

type marineHullKabel struct {
	TypeOfVessel   string `json:"typeOfVessel"`
	NRTRT          string `json:"nrtrt"`
	LengthM        string `json:"lengthm"`
	Classification string `json:"classification"`
	GRTRT          string `json:"grtrt"`
	BreadthM       string `json:"breadthm"`
	Construction   string `json:"construction"`
	DWTRT          string `json:"dwtrt"`
	DepthM         string `json:"depthm"`
}

type aviationHullKabel struct {
	RegistrationMarks string `json:"registrationMarks"`
	MaxPassengers     string `json:"maxPassengers"`
	Pilots            string `json:"pilots"`
	StandardUses      string `json:"standardUses"`
	SpecialUses       string `json:"specialUses"`
	SpecialRentalUses string `json:"specialRentalUses"`
}

type vehicleHEKabel struct {
	ObjectNameHE  string `json:"objectNameHe"`
	ObjectIDHE    string `json:"objectIdHe"`
	BrandName     string `json:"brandName"`
	Brand         string `json:"brand"`
	TypeName      string `json:"typeName"`
	Type          string `json:"type"`
	ChassisNumber string `json:"chassisNumber"`
	EngineNumber  string `json:"engineNumber"`
}

// memorandumKabel - suretyType BACA-SAJA: diisi server dari BusinessCode, isian PUT diabaikan.
type memorandumKabel struct {
	SuretyType           string `json:"suretyType"`
	ApplicationNumber    string `json:"applicationNumber"`
	Accumulation         string `json:"accumulation"`
	Collateral           string `json:"collateral"`
	InsuredSiup          string `json:"insuredSiup"`
	KeyPerson            string `json:"keyPerson"`
	OriginalRate         string `json:"originalRate"`
	PrincipalsOutGo      string `json:"principalsOutGo"`
	TypeOfWording        string `json:"typeOfWording"`
	BuildingStatus       string `json:"buildingStatus"`
	SupportingDocument   string `json:"supportingDocument"`
	PrincipalsData       string `json:"principalsData"`
	ConclusionsProposals string `json:"conclusionsProposals"`
	PrincipalsBackground string `json:"principalsBackground"`
	SummaryFinancial     string `json:"summaryFinancial"`
	PrincipalsGuarantee  string `json:"principalsGuarantee"`
}

// desimalAneka - medan desimal objek (kabel teks <-> model *apd.Decimal), satu daftar untuk kedua arah.
var desimalAneka = []struct {
	jalur string
	kabel func(*objekAnekaKabel) *string
	model func(*models.ObjekAneka) **apd.Decimal
}{
	{"tsi", func(k *objekAnekaKabel) *string { return &k.TSI }, func(o *models.ObjekAneka) **apd.Decimal { return &o.TSI }},
	{"tsiOld", func(k *objekAnekaKabel) *string { return &k.TSIOld }, func(o *models.ObjekAneka) **apd.Decimal { return &o.TSIOld }},
	{"maintenancePeriod", func(k *objekAnekaKabel) *string { return &k.MaintenancePeriod },
		func(o *models.ObjekAneka) **apd.Decimal { return &o.MaintenancePeriod }},
	{"noOfTree", func(k *objekAnekaKabel) *string { return &k.NoOfTree }, func(o *models.ObjekAneka) **apd.Decimal { return &o.NoOfTree }},
	{"areaHectar", func(k *objekAnekaKabel) *string { return &k.AreaHectar },
		func(o *models.ObjekAneka) **apd.Decimal { return &o.AreaHectar }},
	{"limitIdemnityDay", func(k *objekAnekaKabel) *string { return &k.LimitIdemnityDay },
		func(o *models.ObjekAneka) **apd.Decimal { return &o.LimitIdemnityDay }},
	{"limitIdemnityMonth", func(k *objekAnekaKabel) *string { return &k.LimitIdemnityMonth },
		func(o *models.ObjekAneka) **apd.Decimal { return &o.LimitIdemnityMonth }},
	{"limitOfLiability", func(k *objekAnekaKabel) *string { return &k.LimitOfLiability },
		func(o *models.ObjekAneka) **apd.Decimal { return &o.LimitOfLiability }},
	{"briguna.salary", func(k *objekAnekaKabel) *string { return &k.Briguna.Salary },
		func(o *models.ObjekAneka) **apd.Decimal { return &o.Briguna.Salary }},
}

func keObjekAnekaKabel(o models.ObjekAneka) objekAnekaKabel {
	k := objekAnekaKabel{ObjectName: o.ObjectName, Section: o.Section, NoKontrak: o.NoKontrak, StartDate: o.StartDate,
		EndDate: o.EndDate, NoKtp: o.NoKtp, NoNpwp: o.NoNpwp, Currency: o.Currency, DraftWordingID: o.DraftWordingID,
		Others: o.Others, Year: o.Year, Quantity: o.Quantity, UnitQuantity: o.UnitQuantity, SerialNo: o.SerialNo,
		BuiltIn: o.BuiltIn, Flag: o.Flag, GeographicalLimits: o.GeographicalLimits, InvoiceNo: o.InvoiceNo, AlmRiskID: o.AlmRiskID,
		Briguna: brigunaKabel{LoanAccountNumber: o.Briguna.LoanAccountNumber, LoanName: o.Briguna.LoanName, Job: o.Briguna.Job,
			DepositAccountNo: o.Briguna.DepositAccountNo, IDCardNo: o.Briguna.IDCardNo, Institution: o.Briguna.Institution,
			LoanType: o.Briguna.LoanType, DateOfBirth: o.Briguna.DateOfBirth},
		MarineHull: marineHullKabel(o.MarineHull), AviationHull: aviationHullKabel(o.AviationHull),
		VehicleHE: vehicleHEKabel{ObjectNameHE: o.VehicleHE.ObjectNameHE, ObjectIDHE: o.VehicleHE.ObjectIDHE,
			BrandName: o.VehicleHE.BrandName, Brand: o.VehicleHE.Brand, TypeName: o.VehicleHE.TypeName, Type: o.VehicleHE.Type,
			ChassisNumber: o.VehicleHE.ChassisNumber, EngineNumber: o.VehicleHE.EngineNumber},
		Memorandum: memorandumKabel(o.Memorandum)}
	for _, d := range desimalAneka {
		*d.kabel(&k) = teks(*d.model(&o))
	}
	k.Coverages = make([]coverageAnekaKabel, 0, len(o.Coverages))
	for _, c := range o.Coverages {
		k.Coverages = append(k.Coverages, keCoverageAnekaKabel(c))
	}
	return k
}

func keObjekAnekaModel(jalur string, k objekAnekaKabel, masalah *[]string) models.ObjekAneka {
	o := models.ObjekAneka{ObjectName: k.ObjectName, Section: k.Section, NoKontrak: k.NoKontrak, StartDate: k.StartDate,
		EndDate: k.EndDate, NoKtp: k.NoKtp, NoNpwp: k.NoNpwp, Currency: k.Currency, DraftWordingID: k.DraftWordingID,
		Others: k.Others, Year: k.Year, Quantity: k.Quantity, UnitQuantity: k.UnitQuantity, SerialNo: k.SerialNo,
		BuiltIn: k.BuiltIn, Flag: k.Flag, GeographicalLimits: k.GeographicalLimits, InvoiceNo: k.InvoiceNo, AlmRiskID: k.AlmRiskID,
		Briguna: models.BrigunaAneka{LoanAccountNumber: k.Briguna.LoanAccountNumber, LoanName: k.Briguna.LoanName,
			Job: k.Briguna.Job, DepositAccountNo: k.Briguna.DepositAccountNo, IDCardNo: k.Briguna.IDCardNo,
			Institution: k.Briguna.Institution, LoanType: k.Briguna.LoanType, DateOfBirth: k.Briguna.DateOfBirth},
		MarineHull: models.MarineHull(k.MarineHull), AviationHull: models.AviationHull(k.AviationHull),
		VehicleHE: models.VehicleHE{ObjectNameHE: k.VehicleHE.ObjectNameHE, ObjectIDHE: k.VehicleHE.ObjectIDHE,
			BrandName: k.VehicleHE.BrandName, Brand: k.VehicleHE.Brand, TypeName: k.VehicleHE.TypeName, Type: k.VehicleHE.Type,
			ChassisNumber: k.VehicleHE.ChassisNumber, EngineNumber: k.VehicleHE.EngineNumber},
		Memorandum: models.MemorandumAccept(k.Memorandum)}
	for _, d := range desimalAneka {
		*d.model(&o) = services.UraiDesimalIsian(jalur+"."+d.jalur, *d.kabel(&k), masalah)
	}
	if k.Coverages != nil { // kunci absen / null -> nil = coverage tersimpan dipertahankan
		o.Coverages = make([]models.CoverageAneka, 0, len(k.Coverages))
		for i, c := range k.Coverages {
			o.Coverages = append(o.Coverages, keCoverageAnekaModel(fmt.Sprintf("%s.coverages[%d]", jalur, i), c, masalah))
		}
	}
	return o
}

func keDaftarAneka(lokasi []models.LokasiAneka) daftarAneka {
	hasil := daftarAneka{Lokasi: make([]lokasiAnekaKabel, 0, len(lokasi))}
	for _, l := range lokasi {
		k := lokasiAnekaKabel{ObjectNo: l.ObjectNo, IsTopRisk: l.IsTopRisk, RoadType: l.RoadType, RoadName: l.RoadName,
			BuildingNo: l.BuildingNo, ZipCode: l.ZipCode, Country: l.Country, RiskLocation: l.RiskLocation, Territory: l.Territory,
			City: l.City, District: l.District, Province: l.Province, RiskAddressID: l.RiskAddressID,
			Okupasi: make([]okupasiAnekaKabel, 0, len(l.Okupasi)), ObjekGolf: make([]objekAnekaKabel, 0, len(l.ObjekGolf)),
			TotalPerCurrency: keTotalKabel(l.TotalPerCurrency)}
		for _, o := range l.Okupasi {
			ok := okupasiAnekaKabel{OccupationID: o.OccupationID, OccupationName: o.OccupationName,
				Objek: make([]objekAnekaKabel, 0, len(o.Objek))}
			for _, x := range o.Objek {
				ok.Objek = append(ok.Objek, keObjekAnekaKabel(x))
			}
			k.Okupasi = append(k.Okupasi, ok)
		}
		for _, x := range l.ObjekGolf {
			k.ObjekGolf = append(k.ObjekGolf, keObjekAnekaKabel(x))
		}
		hasil.Lokasi = append(hasil.Lokasi, k)
	}
	return hasil
}

// uraiBadanAneka - badan `{lokasi: LokasiAneka[], proRatePercent}`; galat 400 sudah ditulis bila ok=false.
func uraiBadanAneka(w http.ResponseWriter, r *http.Request) ([]models.LokasiAneka, *apd.Decimal, bool) {
	var b daftarAneka
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, batasBadanObjek)).Decode(&b); err != nil {
		galat.Tulis(w, http.StatusBadRequest, "badan permintaan bukan JSON daftar lokasi Aneka yang sah")
		return nil, nil, false
	}
	var masalah []string
	proRate := services.UraiDesimalIsian("proRatePercent", b.ProRatePercent, &masalah)
	lokasi := make([]models.LokasiAneka, len(b.Lokasi))
	for i, l := range b.Lokasi {
		m := models.LokasiAneka{ObjectNo: l.ObjectNo, IsTopRisk: l.IsTopRisk, RoadType: l.RoadType, RoadName: l.RoadName,
			BuildingNo: l.BuildingNo, ZipCode: l.ZipCode, Country: l.Country, RiskLocation: l.RiskLocation, Territory: l.Territory,
			City: l.City, District: l.District, Province: l.Province, RiskAddressID: l.RiskAddressID}
		for j, o := range l.Okupasi {
			ok := models.OkupasiAneka{OccupationID: o.OccupationID, OccupationName: o.OccupationName}
			for n, x := range o.Objek {
				ok.Objek = append(ok.Objek, keObjekAnekaModel(fmt.Sprintf("lokasi[%d].okupasi[%d].objek[%d]", i, j, n), x, &masalah))
			}
			m.Okupasi = append(m.Okupasi, ok)
		}
		for n, x := range l.ObjekGolf {
			m.ObjekGolf = append(m.ObjekGolf, keObjekAnekaModel(fmt.Sprintf("lokasi[%d].objekGolf[%d]", i, n), x, &masalah))
		}
		lokasi[i] = m
	}
	if err := services.GalatIsianObjek(masalah); err != nil {
		tulisGalat(w, err)
		return nil, nil, false
	}
	return lokasi, proRate, true
}

// bacaObjekAneka - GET /api/nbfacin/kasus/{caseId}/objek-aneka.
func bacaObjekAneka(svc *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		d, err := svc.BacaObjekAneka(r.Context(), r.PathValue("caseId"))
		if err != nil {
			tulisGalat(w, err)
			return
		}
		out := keDaftarAneka(d)
		out.ProRatePercent = teks(svc.ProRateKasusAneka(r.Context(), r.PathValue("caseId")))
		galat.TulisJSON(w, out)
	}
}

// pelanggaranKabel - satu pesan proteksi objek (422); okupasi null = objek Golf; indeks 0-based.
type pelanggaranKabel struct {
	Lokasi  int    `json:"lokasi"`
	Okupasi *int   `json:"okupasi"`
	Objek   int    `json:"objek"`
	Medan   string `json:"medan"`
	Pesan   string `json:"pesan"`
}

// tulisGalatAneka - 422 proteksi objek `{galat, pelanggaran}`; galat lain lewat tulisGalat.
func tulisGalatAneka(w http.ResponseWriter, err error) {
	var g services.GalatProteksiAneka
	if !errors.As(err, &g) {
		tulisGalat(w, err)
		return
	}
	p := make([]pelanggaranKabel, 0, len(g.Pelanggaran))
	pesan := make([]string, 0, len(g.Pelanggaran))
	for _, x := range g.Pelanggaran {
		p = append(p, pelanggaranKabel{Lokasi: x.Lokasi, Okupasi: x.Okupasi, Objek: x.Objek, Medan: x.Medan, Pesan: x.Pesan})
		pesan = append(pesan, x.Pesan)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusUnprocessableEntity)
	_ = json.NewEncoder(w).Encode(struct {
		Galat       string             `json:"galat"`
		Pelanggaran []pelanggaranKabel `json:"pelanggaran"`
	}{strings.Join(pesan, "; "), p})
}

// simpanObjekAneka - PUT /api/nbfacin/kasus/{caseId}/objek-aneka (tombol Save): ganti utuh, jawab baca ulang.
func simpanObjekAneka(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pelaku := inti.PelakuDari(r, stubPelaku)
		if err := inti.WajibIdentitas(pelaku); err != nil {
			tulisGalat(w, err)
			return
		}
		lokasi, proRate, ok := uraiBadanAneka(w, r)
		if !ok {
			return
		}
		d, err := svc.GantiObjekAneka(r.Context(), pelaku, r.PathValue("caseId"), lokasi, proRate)
		if err != nil {
			tulisGalatAneka(w, err)
			return
		}
		out := keDaftarAneka(d)
		out.ProRatePercent = teks(svc.ProRateKasusAneka(r.Context(), r.PathValue("caseId")))
		galat.TulisJSON(w, out)
	}
}

// barisHitungAnekaKabel - satu baris badan POST …/hitung-coverage-aneka.
type barisHitungAnekaKabel struct {
	Coverage coverageAnekaKabel `json:"coverage"`
	Pemicu   string             `json:"pemicu"`
}

// hitungCoverageAneka - POST /api/nbfacin/kasus/{caseId}/hitung-coverage-aneka (tiket 73 P2): badan `{proRatePercent, baris:
// [{coverage, pemicu}]}` -> `{proRatePercent, baris: CoverageAneka[]}` urutan sama, tanpa menyimpan. proRatePercent = nilai CASE
// yang dipegang layar (masuk) dan sesudah rantai (keluar).
func hitungCoverageAneka(svc *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var b struct {
			ProRatePercent string                  `json:"proRatePercent"`
			Baris          []barisHitungAnekaKabel `json:"baris"`
		}
		if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, batasBadanObjek)).Decode(&b); err != nil {
			galat.Tulis(w, http.StatusBadRequest, "badan permintaan bukan JSON hitung coverage Aneka yang sah")
			return
		}
		var masalah []string
		pr := services.UraiDesimalIsian("proRatePercent", b.ProRatePercent, &masalah)
		baris := make([]services.BarisHitungAneka, 0, len(b.Baris))
		for i, x := range b.Baris {
			baris = append(baris, services.BarisHitungAneka{
				Coverage: keCoverageAnekaModel(fmt.Sprintf("baris[%d].coverage", i), x.Coverage, &masalah), Pemicu: x.Pemicu})
		}
		if len(masalah) > 0 {
			galat.Tulis(w, http.StatusBadRequest, strings.Join(masalah, "; "))
			return
		}
		hasil, prBaru, err := svc.HitungCoverageAnekaKasus(r.Context(), r.PathValue("caseId"), pr, baris)
		if err != nil {
			tulisGalat(w, err)
			return
		}
		out := struct {
			ProRatePercent string               `json:"proRatePercent"`
			Baris          []coverageAnekaKabel `json:"baris"`
		}{ProRatePercent: teks(prBaru), Baris: make([]coverageAnekaKabel, 0, len(hasil))}
		for _, c := range hasil {
			out.Baris = append(out.Baris, keCoverageAnekaKabel(c))
		}
		galat.TulisJSON(w, out)
	}
}
