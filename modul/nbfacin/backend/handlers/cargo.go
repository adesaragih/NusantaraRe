package handlers

// Tab Object + Coverage MARINE CARGO (tiket 74 P1 + P2): GET/PUT /api/nbfacin/kasus/{caseId}/objek-cargo dan POST
// /api/nbfacin/kasus/{caseId}/hitung-coverage-cargo - kontrak `ObjekCargo` / `CoverageCargo` sesi 0f (10-10-2026). Angka teks
// desimal bertitik (ADR-0003), tanggal DD-MM-YYYY (ralat 0f), bool boolean. Medan "server" (catatan master, maxAgeDetail,
// ship.age, flagDelete, currencyId / OldId, discount, tsiLiability, tsiNusantaraRe, premiNusantaraRe, proRatePercent, totalGross)
// dikirim GET dan DIABAIKAN PUT. `coverages` absen / null = coverage lama cargo ber-`id` sama dipertahankan.

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/cockroachdb/apd/v3"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/services"
)

type daftarCargo struct {
	Cargo []objekCargoKabel `json:"cargo"`
}

type objekCargoKabel struct {
	ID             string               `json:"id"`
	IsTopRisk      bool                 `json:"isTopRisk"`
	TradingID      string               `json:"tradingId"`
	TradingNote    string               `json:"tradingNote"`
	GoodID         string               `json:"goodId"`
	GoodNote       string               `json:"goodNote"`
	ConveyanceID   string               `json:"conveyanceId"`
	ConveyanceName string               `json:"conveyanceName"`
	ConveyanceNote string               `json:"conveyanceNote"`
	PackingID      string               `json:"packingId"`
	PackingNote    string               `json:"packingNote"`
	FromRute       string               `json:"fromRute"`
	ToRute         string               `json:"toRute"`
	SailDate       string               `json:"sailDate"`
	BLNumber       string               `json:"blNumber"`
	InvoiceNumber  string               `json:"invoiceNumber"`
	InvoiceDate    string               `json:"invoiceDate"`
	MaxAge         string               `json:"maxAge"`
	MaxAgeDetail   string               `json:"maxAgeDetail"`
	SurveyAgent    agenSurveiKabel      `json:"surveyAgent"`
	LC             lcCargoKabel         `json:"lc"`
	Ship           kapalCargoKabel      `json:"ship"`
	AdditionalShip []kapalTambahanKabel `json:"additionalShip"`
	FlagDelete     string               `json:"flagDelete"`
	Coverages      []coverageCargoKabel `json:"coverages"`
	TotalGross     []totalCoverageKabel `json:"totalGross"`
}

type agenSurveiKabel struct {
	AgentID  string `json:"agentId"`
	FullName string `json:"fullName"`
	Address  string `json:"address"`
	City     string `json:"city"`
	Country  string `json:"country"`
}

type lcCargoKabel struct {
	LC        bool   `json:"lc"`
	Bank      string `json:"bank"`
	Number    string `json:"number"`
	Condition string `json:"condition"`
	Remark    string `json:"remark"`
	EndDate   string `json:"endDate"`
}

type kapalCargoKabel struct {
	ID             string `json:"id"`
	NmShip         string `json:"nmShip"`
	ShipTypeID     string `json:"shipTypeId"`
	ShipTypeName   string `json:"shipTypeName"`
	GRT            string `json:"grt"`
	NRT            string `json:"nrt"`
	DWT            string `json:"dwt"`
	Age            string `json:"age"`
	FlagName       string `json:"flagName"`
	Construction   string `json:"construction"`
	YMake1         string `json:"yMake1"`
	YMake2         string `json:"yMake2"`
	Remark         string `json:"remark"`
	Register       string `json:"register"`
	Imo            string `json:"imo"`
	FormerShipName string `json:"formerShipName"`
	Classification string `json:"classification"`
	Rebuilt        string `json:"rebuilt"`
}

type kapalTambahanKabel struct {
	ID           string `json:"id"`
	Register     string `json:"register"`
	NmShip       string `json:"nmShip"`
	VoyageNo     string `json:"voyageNo"`
	ShipTypeID   string `json:"shipTypeId"`
	FlagName     string `json:"flagName"`
	NRT          string `json:"nrt"`
	DWT          string `json:"dwt"`
	GRT          string `json:"grt"`
	YMake1       string `json:"yMake1"`
	Construction string `json:"construction"`
}

type coverageCargoKabel struct {
	Coverage           string                   `json:"coverage"`
	CoverageNote       string                   `json:"coverageNote"`
	CoverageInitial    string                   `json:"coverageInitial"`
	Currency           string                   `json:"currency"`
	CurrencyID         string                   `json:"currencyId"`
	CurrencyOldID      string                   `json:"currencyOldId"`
	Rate               string                   `json:"rate"`
	LimitOfLiability   string                   `json:"limitOfLiability"`
	DiscountPercentage string                   `json:"discountPercentage"`
	TSI                string                   `json:"tsi"`
	Premium            string                   `json:"premium"`
	MinPremium         string                   `json:"minPremium"`
	Discount           string                   `json:"discount"`
	TSILiability       string                   `json:"tsiLiability"`
	TSINusantaraRe     string                   `json:"tsiNusantaraRe"`
	PremiNusantaraRe   string                   `json:"premiNusantaraRe"`
	ProRatePercent     string                   `json:"proRatePercent"`
	WarrantyExtension  []warrantyExtensionKabel `json:"warrantyExtension"`
	Deductibles        []deductibleCargoKabel   `json:"deductibles"`
}

type warrantyExtensionKabel struct {
	Coverage     string `json:"coverage"`
	CoverageNote string `json:"coverageNote"`
}

// deductibleCargoKabel - `Deductible` (tiket 45) + descriptions (`.Descriptions` per baris deductible, L17652).
type deductibleCargoKabel struct {
	deductibleKabel
	Descriptions string `json:"descriptions"`
}

// ---- model -> kabel -------------------------------------------------------------------------------------------------------

func keDaftarCargo(d []services.CargoTampil) daftarCargo {
	out := daftarCargo{Cargo: make([]objekCargoKabel, 0, len(d))}
	for _, x := range d {
		c := x.Objek
		k := objekCargoKabel{ID: c.ID, IsTopRisk: c.IsTopRisk, TradingID: c.TradingID, TradingNote: c.TradingNote, GoodID: c.GoodID,
			GoodNote: c.GoodNote, ConveyanceID: c.ConveyanceID, ConveyanceName: c.ConveyanceName, ConveyanceNote: c.ConveyanceNote,
			PackingID: c.PackingID, PackingNote: c.PackingNote, FromRute: c.FromRute, ToRute: c.ToRute, SailDate: c.SailDate,
			BLNumber: c.BLNumber, InvoiceNumber: c.InvoiceNumber, InvoiceDate: c.InvoiceDate, MaxAge: c.MaxAge,
			MaxAgeDetail: c.MaxAgeDetail, SurveyAgent: agenSurveiKabel(c.SurveyAgent), LC: lcCargoKabel(c.LC),
			Ship: kapalCargoKabel{ID: c.Ship.ID, NmShip: c.Ship.NmShip, ShipTypeID: c.Ship.ShipTypeID, ShipTypeName: c.Ship.ShipTypeName,
				GRT: c.Ship.GRT, NRT: c.Ship.NRT, DWT: c.Ship.DWT, Age: c.Ship.Age, FlagName: c.Ship.FlagName,
				Construction: c.Ship.Construction, YMake1: c.Ship.YMake1, YMake2: c.Ship.YMake2, Remark: c.Ship.Remark,
				Register: c.Ship.Register, Imo: c.Ship.Imo, FormerShipName: c.Ship.FormerShipName, Classification: c.Ship.Classification,
				Rebuilt: c.Ship.Rebuilt},
			AdditionalShip: make([]kapalTambahanKabel, 0, len(c.Ship.AdditionalShip)), FlagDelete: c.FlagDelete,
			Coverages: make([]coverageCargoKabel, 0, len(c.Coverages)), TotalGross: make([]totalCoverageKabel, 0, len(x.TotalGross))}
		for _, s := range c.Ship.AdditionalShip {
			k.AdditionalShip = append(k.AdditionalShip, kapalTambahanKabel(s))
		}
		for _, v := range c.Coverages {
			k.Coverages = append(k.Coverages, keCoverageCargoKabel(v))
		}
		for _, g := range x.TotalGross {
			k.TotalGross = append(k.TotalGross, totalCoverageKabel{Currency: g.Currency, TSI: teks(g.TSI), Premium: teks(g.Premium),
				Rate: teks(g.Rate)})
		}
		out.Cargo = append(out.Cargo, k)
	}
	return out
}

func keCoverageCargoKabel(v models.CoverageCargo) coverageCargoKabel {
	k := coverageCargoKabel{Coverage: v.Coverage, CoverageNote: v.CoverageNote, CoverageInitial: v.CoverageInitial,
		Currency: v.Currency, CurrencyID: v.CurrencyID, CurrencyOldID: v.CurrencyOldID, Rate: teks(v.Rate),
		LimitOfLiability: teks(v.LimitOfLiability), DiscountPercentage: teks(v.DiscountPercentage), TSI: teks(v.TSI),
		Premium: teks(v.Premium), MinPremium: teks(v.MinPremium), Discount: teks(v.Discount), TSILiability: teks(v.TSILiability),
		TSINusantaraRe: teks(v.TSINusantaraRe), PremiNusantaraRe: teks(v.PremiNusantaraRe), ProRatePercent: teks(v.ProRatePercent),
		WarrantyExtension: make([]warrantyExtensionKabel, 0, len(v.WarrantyExtension)),
		Deductibles:       make([]deductibleCargoKabel, 0, len(v.Deductibles))}
	for _, w := range v.WarrantyExtension {
		k.WarrantyExtension = append(k.WarrantyExtension, warrantyExtensionKabel(w))
	}
	for _, d := range v.Deductibles {
		k.Deductibles = append(k.Deductibles, deductibleCargoKabel{deductibleKabel: keDeductibleKabel([]models.Deductible{d.Deductible})[0],
			Descriptions: d.Descriptions})
	}
	return k
}

// ---- kabel -> model (medan server diabaikan) ------------------------------------------------------------------------------

func keObjekCargoModel(jalur string, k objekCargoKabel, masalah *[]string) models.ObjekCargo {
	c := models.ObjekCargo{ID: k.ID, IsTopRisk: k.IsTopRisk, TradingID: k.TradingID, GoodID: k.GoodID, ConveyanceID: k.ConveyanceID,
		PackingID: k.PackingID, FromRute: k.FromRute, ToRute: k.ToRute, SailDate: k.SailDate, BLNumber: k.BLNumber,
		InvoiceNumber: k.InvoiceNumber, InvoiceDate: k.InvoiceDate, MaxAge: k.MaxAge, SurveyAgent: models.AgenSurvei(k.SurveyAgent),
		LC: models.LCCargo(k.LC),
		Ship: models.KapalCargo{ID: k.Ship.ID, NmShip: k.Ship.NmShip, ShipTypeID: k.Ship.ShipTypeID, ShipTypeName: k.Ship.ShipTypeName,
			GRT: k.Ship.GRT, NRT: k.Ship.NRT, DWT: k.Ship.DWT, FlagName: k.Ship.FlagName, Construction: k.Ship.Construction,
			YMake1: k.Ship.YMake1, YMake2: k.Ship.YMake2, Remark: k.Ship.Remark, Register: k.Ship.Register, Imo: k.Ship.Imo,
			FormerShipName: k.Ship.FormerShipName, Classification: k.Ship.Classification, Rebuilt: k.Ship.Rebuilt}}
	for _, s := range k.AdditionalShip {
		c.Ship.AdditionalShip = append(c.Ship.AdditionalShip, models.KapalTambahan(s))
	}
	if k.Coverages != nil {
		c.Coverages = make([]models.CoverageCargo, 0, len(k.Coverages))
		for i, v := range k.Coverages {
			c.Coverages = append(c.Coverages, keCoverageCargoModel(fmt.Sprintf("%s.coverages[%d]", jalur, i), v, masalah))
		}
	}
	return c
}

// keCoverageCargoModel - medan isian coverage (server diabaikan, kecuali currencyId di hitung - dihitung ulang juga).
func keCoverageCargoModel(jalur string, k coverageCargoKabel, masalah *[]string) models.CoverageCargo {
	u := func(medan, nilai string) *apd.Decimal {
		return services.UraiDesimalIsian(jalur+"."+medan, nilai, masalah)
	}
	v := models.CoverageCargo{Coverage: k.Coverage, CoverageNote: k.CoverageNote, CoverageInitial: k.CoverageInitial,
		Currency: k.Currency, Rate: u("rate", k.Rate), LimitOfLiability: u("limitOfLiability", k.LimitOfLiability),
		DiscountPercentage: u("discountPercentage", k.DiscountPercentage), TSI: u("tsi", k.TSI), Premium: u("premium", k.Premium),
		MinPremium: u("minPremium", k.MinPremium)}
	for _, w := range k.WarrantyExtension {
		v.WarrantyExtension = append(v.WarrantyExtension, models.WarrantyExtension(w))
	}
	if v.WarrantyExtension == nil {
		v.WarrantyExtension = []models.WarrantyExtension{}
	}
	ded := make([]deductibleKabel, len(k.Deductibles))
	for i, d := range k.Deductibles {
		ded[i] = d.deductibleKabel
	}
	v.Deductibles = []models.DeductibleCargo{}
	for i, d := range keDeductibleModel(jalur+".", ded, masalah) {
		v.Deductibles = append(v.Deductibles, models.DeductibleCargo{Deductible: d, Descriptions: k.Deductibles[i].Descriptions})
	}
	return v
}

// ---- rute -----------------------------------------------------------------------------------------------------------------

// bacaObjekCargo - GET /api/nbfacin/kasus/{caseId}/objek-cargo.
func bacaObjekCargo(svc *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		d, err := svc.BacaObjekCargo(r.Context(), r.PathValue("caseId"))
		if err != nil {
			tulisGalat(w, err)
			return
		}
		galat.TulisJSON(w, keDaftarCargo(d))
	}
}

// simpanObjekCargo - PUT /api/nbfacin/kasus/{caseId}/objek-cargo: daftar diganti utuh, jawab baca ulang.
func simpanObjekCargo(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pelaku := inti.PelakuDari(r, stubPelaku)
		if err := inti.WajibIdentitas(pelaku); err != nil {
			tulisGalat(w, err)
			return
		}
		var b daftarCargo
		if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, batasBadanObjek)).Decode(&b); err != nil {
			galat.Tulis(w, http.StatusBadRequest, "badan permintaan bukan JSON daftar cargo yang sah")
			return
		}
		var masalah []string
		cargo := make([]models.ObjekCargo, 0, len(b.Cargo))
		for i, k := range b.Cargo {
			cargo = append(cargo, keObjekCargoModel(fmt.Sprintf("cargo[%d]", i), k, &masalah))
		}
		if len(masalah) > 0 {
			galat.Tulis(w, http.StatusBadRequest, strings.Join(masalah, "; "))
			return
		}
		d, err := svc.GantiObjekCargo(r.Context(), pelaku, r.PathValue("caseId"), cargo)
		if err != nil {
			tulisGalat(w, err)
			return
		}
		galat.TulisJSON(w, keDaftarCargo(d))
	}
}

// hitungCoverageCargo - POST /api/nbfacin/kasus/{caseId}/hitung-coverage-cargo: `{baris: [{coverage, pemicu}]}` -> `{baris:
// CoverageCargo[]}` urutan sama, tanpa menyimpan.
func hitungCoverageCargo(svc *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var b struct {
			Baris []struct {
				Coverage coverageCargoKabel `json:"coverage"`
				Pemicu   string             `json:"pemicu"`
			} `json:"baris"`
		}
		if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, batasBadanObjek)).Decode(&b); err != nil {
			galat.Tulis(w, http.StatusBadRequest, "badan permintaan bukan JSON hitung coverage cargo yang sah")
			return
		}
		var masalah []string
		baris := make([]services.BarisHitungCargo, 0, len(b.Baris))
		for i, x := range b.Baris {
			baris = append(baris, services.BarisHitungCargo{Coverage: keCoverageCargoModel(fmt.Sprintf("baris[%d].coverage", i),
				x.Coverage, &masalah), Pemicu: x.Pemicu})
		}
		if len(masalah) > 0 {
			galat.Tulis(w, http.StatusBadRequest, strings.Join(masalah, "; "))
			return
		}
		hasil, err := svc.HitungCoverageCargo(r.Context(), r.PathValue("caseId"), baris)
		if err != nil {
			tulisGalat(w, err)
			return
		}
		out := struct {
			Baris []coverageCargoKabel `json:"baris"`
		}{Baris: make([]coverageCargoKabel, 0, len(hasil))}
		for _, c := range hasil {
			out.Baris = append(out.Baris, keCoverageCargoKabel(c))
		}
		galat.TulisJSON(w, out)
	}
}

// barisMasterCargoKabel - satu pilihan master Marine (api.ts BarisMasterCargo).
type barisMasterCargoKabel struct {
	ID          string `json:"id"`
	Description string `json:"description"`
}

// barisAgenKabel / barisMaxAgeKabel / barisKondisiKabel - api.ts BarisAgentSurvei / cariMaxAge / BarisKondisiMarine.
type barisAgenKabel struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	City    string `json:"city"`
	Country string `json:"country"`
	Contact string `json:"contact"`
}

type barisMaxAgeKabel struct {
	ID       string `json:"id"`
	Warranty string `json:"warranty"`
}

type barisKondisiKabel struct {
	ID            string `json:"id"`
	ConditionName string `json:"conditionName"`
}

// tulisBaris - {baris: [...]} atau galat.
func tulisBaris[T any](w http.ResponseWriter, baris []T, err error) {
	if err != nil {
		tulisGalat(w, err)
		return
	}
	if baris == nil {
		baris = []T{}
	}
	galat.TulisJSON(w, struct {
		Baris []T `json:"baris"`
	}{baris})
}

// petaBaris - larik model -> larik kabel.
func petaBaris[M, K any](d []M, f func(M) K) []K {
	k := make([]K, 0, len(d))
	for _, x := range d {
		k = append(k, f(x))
	}
	return k
}

// daftarMasterCargo - GET /api/nbfacin/marine/{jenis}: trading / goods / conveyance / packing -> {baris: [{id, description}]};
// agent?q= -> {baris: [{id, name, address, city, country, contact}]}; max-age?q= -> {baris: [{id, warranty}]}; coverage?q= ->
// {baris: [{id, conditionName}]} (kondisi marine).
func daftarMasterCargo(svc *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")
		switch r.PathValue("jenis") {
		case "agent":
			d, err := svc.CariAgenSurvei(r.Context(), q)
			tulisBaris(w, petaBaris(d, func(x models.BarisAgenSurvei) barisAgenKabel { return barisAgenKabel(x) }), err)
			return
		case "max-age":
			d, err := svc.CariMaxAge(r.Context(), q)
			tulisBaris(w, petaBaris(d, func(x models.BarisMaxAge) barisMaxAgeKabel { return barisMaxAgeKabel(x) }), err)
			return
		case "coverage":
			d, err := svc.CariKondisiMarine(r.Context(), q)
			tulisBaris(w, petaBaris(d, func(x models.BarisKondisiMarine) barisKondisiKabel { return barisKondisiKabel(x) }), err)
			return
		}
		d, err := svc.DaftarMasterCargo(r.Context(), r.PathValue("jenis"))
		if err != nil {
			tulisGalat(w, err)
			return
		}
		out := struct {
			Baris []barisMasterCargoKabel `json:"baris"`
		}{Baris: make([]barisMasterCargoKabel, 0, len(d))}
		for _, x := range d {
			out.Baris = append(out.Baris, barisMasterCargoKabel(x))
		}
		galat.TulisJSON(w, out)
	}
}
