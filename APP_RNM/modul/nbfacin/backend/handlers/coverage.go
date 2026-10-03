package handlers

// Tab Coverage FIRE tahap C1 (tiket 43): POST hitung-coverage (CountPremi_ACT tanpa menyimpan), GET coverage (popup
// Choose Coverage), GET coverage-otomatis (AddCoverageAutoFire). Kontrak `CoverageObjek` / `BarisCoverage` /
// `TotalCoverage` frontend (modul/nbfacin/frontend/api.ts).

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/galat"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/services"
)

// coverageKabel - kontrak `CoverageObjek`; uang / rate / persen = teks desimal (ADR-0034). tsi, tsiLiability,
// proRatePercent dihitung server - diabaikan saat masuk.
type coverageKabel struct {
	Coverage            string `json:"coverage"`
	OldID               string `json:"oldId"`
	CoverageNote        string `json:"coverageNote"`
	CoverageBasis       string `json:"coverageBasis"`
	Day                 string `json:"day"`
	TSI                 string `json:"tsi"`
	Indemnity           string `json:"indemnity"`
	Rate                string `json:"rate"`
	RateOJK             string `json:"rateOjk"`
	FirstLoss           string `json:"firstLoss"`
	DiscountPercentage  string `json:"discountPercentage"`
	TSILiability        string `json:"tsiLiability"`
	NetRate             string `json:"netRate"`
	LimitOfLiability    string `json:"limitOfLiability"`
	PctLoL              string `json:"pctLol"`
	ProRatePercent      string `json:"proRatePercent"`
	IndemnityPercentage string `json:"indemnityPercentage"`
	FirstScale          string `json:"firstScale"`
	Sublimit            string `json:"sublimit"`
	LostLimit           string `json:"lostLimit"`
	EmlPml              string `json:"emlPml"`
	Discount            string `json:"discount"`
	Premium             string `json:"premium"`
	Conditions          string `json:"conditions"`
	// tiket 45 - selalu larik ke luar; boleh absen saat masuk.
	Deductibles []deductibleKabel `json:"deductibles"`
}

// deductibleKabel - kontrak `Deductible` (tiket 45); kode = teks, persen / uang / hari = teks desimal (ADR-0034).
type deductibleKabel struct {
	TypeDeductible  string `json:"typeDeductible"`
	PctDeductible   string `json:"pctDeductible"`
	MinMax          string `json:"minMax"`
	Currency        string `json:"currency"`
	TypeDeductible2 string `json:"typeDeductible2"`
	PctDeductible2  string `json:"pctDeductible2"`
	Condition       string `json:"condition"`
	Amount          string `json:"amount"`
	InputCondition  string `json:"inputCondition"`
	TimeExcess      string `json:"timeExcess"`
}

func keDeductibleKabel(d []models.Deductible) []deductibleKabel {
	hasil := make([]deductibleKabel, 0, len(d))
	for _, x := range d {
		hasil = append(hasil, deductibleKabel{TypeDeductible: x.TypeDeductible, PctDeductible: teks(x.PctDeductible), MinMax: x.MinMax,
			Currency: x.Currency, TypeDeductible2: x.TypeDeductible2, PctDeductible2: teks(x.PctDeductible2), Condition: x.Condition,
			Amount: teks(x.Amount), InputCondition: x.InputCondition, TimeExcess: teks(x.TimeExcess)})
	}
	return hasil
}

// keDeductibleModel - deductible coverage; jalur pesan "<awalCoverage>deductibles[d].<medan>".
func keDeductibleModel(awalCoverage string, d []deductibleKabel, masalah *[]string) []models.Deductible {
	if len(d) == 0 {
		return nil
	}
	hasil := make([]models.Deductible, 0, len(d))
	for i, x := range d {
		awal := fmt.Sprintf("%sdeductibles[%d].", awalCoverage, i)
		u := func(medan, nilai string) *apd.Decimal { return services.UraiDesimalIsian(awal+medan, nilai, masalah) }
		hasil = append(hasil, models.Deductible{TypeDeductible: x.TypeDeductible, MinMax: x.MinMax, Currency: x.Currency,
			TypeDeductible2: x.TypeDeductible2, Condition: x.Condition, InputCondition: x.InputCondition,
			PctDeductible: u("pctDeductible", x.PctDeductible), PctDeductible2: u("pctDeductible2", x.PctDeductible2),
			Amount: u("amount", x.Amount), TimeExcess: u("timeExcess", x.TimeExcess)})
	}
	return hasil
}

// totalCoverageKabel - kontrak `TotalCoverage` (BACA-SAJA).
type totalCoverageKabel struct {
	Currency string `json:"currency"`
	TSI      string `json:"tsi"`
	Premium  string `json:"premium"`
	Rate     string `json:"rate"`
}

// barisCoverageKabel - kontrak `BarisCoverage`.
type barisCoverageKabel struct {
	ID    string `json:"id"`
	OldID string `json:"oldId"`
	Nama  string `json:"nama"`
}

// batasBadanCoverage - badan POST hitung-coverage terbesar yang diterima.
const batasBadanCoverage = 64 << 10

func keCoverageKabel(d []models.CoverageObjek) []coverageKabel {
	hasil := make([]coverageKabel, 0, len(d))
	for _, c := range d {
		hasil = append(hasil, satuCoverageKabel(c))
	}
	return hasil
}

func satuCoverageKabel(c models.CoverageObjek) coverageKabel {
	return coverageKabel{Coverage: c.Coverage, OldID: c.OldID, CoverageNote: c.CoverageNote,
		CoverageBasis: c.CoverageBasis, Day: c.Day, TSI: teks(c.TSI), Indemnity: c.Indemnity, Rate: teks(c.Rate),
		RateOJK: teks(c.RateOJK), FirstLoss: teks(c.FirstLoss), DiscountPercentage: teks(c.DiscountPercentage),
		TSILiability: teks(c.TSILiability), NetRate: teks(c.NetRate), LimitOfLiability: teks(c.LimitOfLiability),
		PctLoL: teks(c.PctLoL), ProRatePercent: teks(c.ProRatePercent), IndemnityPercentage: teks(c.IndemnityPercentage),
		FirstScale: teks(c.FirstScale), Sublimit: teks(c.Sublimit), LostLimit: teks(c.LostLimit), EmlPml: teks(c.EmlPml),
		Discount: teks(c.Discount), Premium: teks(c.Premium), Conditions: c.Conditions, Deductibles: keDeductibleKabel(c.Deductibles)}
}

// keCoverageModel - satu coverage kabel -> model; `awal` = awalan jalur pesan 400 (mis. "baris[0].items[1].coverages[2].").
// Medan server (tsi, tsiLiability, proRatePercent) tidak diurai - dihitung ulang services.
func keCoverageModel(awal string, c coverageKabel, masalah *[]string) models.CoverageObjek {
	u := func(medan, nilai string) *apd.Decimal { return services.UraiDesimalIsian(awal+medan, nilai, masalah) }
	return models.CoverageObjek{Coverage: c.Coverage, OldID: c.OldID, CoverageNote: c.CoverageNote,
		CoverageBasis: c.CoverageBasis, Day: c.Day, Indemnity: c.Indemnity, Conditions: c.Conditions,
		Rate: u("rate", c.Rate), RateOJK: u("rateOjk", c.RateOJK), FirstLoss: u("firstLoss", c.FirstLoss),
		DiscountPercentage: u("discountPercentage", c.DiscountPercentage), NetRate: u("netRate", c.NetRate),
		LimitOfLiability: u("limitOfLiability", c.LimitOfLiability), PctLoL: u("pctLol", c.PctLoL),
		IndemnityPercentage: u("indemnityPercentage", c.IndemnityPercentage), FirstScale: u("firstScale", c.FirstScale),
		Sublimit: u("sublimit", c.Sublimit), LostLimit: u("lostLimit", c.LostLimit), EmlPml: u("emlPml", c.EmlPml),
		Discount: u("discount", c.Discount), Premium: u("premium", c.Premium),
		Deductibles: keDeductibleModel(awal, c.Deductibles, masalah)}
}

// keCoverageModelDaftar - coverages item; jalur pesan "<awalItem>coverages[k].".
func keCoverageModelDaftar(awalItem string, d []coverageKabel, masalah *[]string) []models.CoverageObjek {
	if len(d) == 0 {
		return nil
	}
	hasil := make([]models.CoverageObjek, 0, len(d))
	for k, c := range d {
		hasil = append(hasil, keCoverageModel(fmt.Sprintf("%scoverages[%d].", awalItem, k), c, masalah))
	}
	return hasil
}

func keTotalKabel(d []models.TotalCoverage) []totalCoverageKabel {
	hasil := make([]totalCoverageKabel, 0, len(d))
	for _, t := range d {
		hasil = append(hasil, totalCoverageKabel{Currency: t.Currency, TSI: teks(t.TSI), Premium: teks(t.Premium), Rate: teks(t.Rate)})
	}
	return hasil
}

func tulisBarisCoverage(w http.ResponseWriter, d []models.BarisCoverage) {
	baris := make([]barisCoverageKabel, 0, len(d))
	for _, b := range d {
		baris = append(baris, barisCoverageKabel(b))
	}
	galat.TulisJSON(w, struct {
		Baris []barisCoverageKabel `json:"baris"`
	}{baris})
}

// cariCoverage - GET /api/nbfacin/coverage?cari= (popup Choose Coverage).
func cariCoverage(svc *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		d, err := svc.CariCoverage(r.Context(), r.URL.Query().Get("cari"))
		if err != nil {
			tulisGalat(w, err)
			return
		}
		tulisBarisCoverage(w, d)
	}
}

// coverageOtomatis - GET /api/nbfacin/coverage-otomatis (lima coverage AddCoverageAutoFire).
func coverageOtomatis(svc *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		d, err := svc.CoverageOtomatis(r.Context())
		if err != nil {
			tulisGalat(w, err)
			return
		}
		tulisBarisCoverage(w, d)
	}
}

// badanHitungNetRate - badan POST hitung-net-rate (tiket 44); totalNetRate kosong = 0 (CalculateNetRate langkah 1).
type badanHitungNetRate struct {
	TSI            string          `json:"tsi"`
	TotalNetRate   string          `json:"totalNetRate"`
	Coverages      []coverageKabel `json:"coverages"`
	IsAdjustable   bool            `json:"isAdjustable"`
	PctAdjustOther string          `json:"pctAdjustOther"`
}

// hitungNetRate - POST /api/nbfacin/kasus/{caseId}/hitung-net-rate: bagi Total Net Rate ke coverage, tanpa menyimpan.
func hitungNetRate(svc *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var b badanHitungNetRate
		if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, batasBadanCoverage)).Decode(&b); err != nil {
			galat.Tulis(w, http.StatusBadRequest, "badan permintaan bukan JSON hitung net rate yang sah")
			return
		}
		var masalah []string
		cov := keCoverageModelDaftar("", b.Coverages, &masalah)
		tsi := services.UraiDesimalIsian("tsi", b.TSI, &masalah)
		tnr := services.UraiDesimalIsian("totalNetRate", b.TotalNetRate, &masalah)
		it := services.PenyesuaianItem{IsAdjustable: b.IsAdjustable,
			PctAdjustOther: services.UraiDesimalIsian("pctAdjustOther", b.PctAdjustOther, &masalah)}
		if err := services.GalatIsianObjek(masalah); err != nil {
			tulisGalat(w, err)
			return
		}
		h, err := svc.HitungNetRateKasus(r.Context(), r.PathValue("caseId"), cov, tsi, tnr, it)
		if err != nil {
			tulisGalat(w, err)
			return
		}
		galat.TulisJSON(w, struct {
			Coverages    []coverageKabel `json:"coverages"`
			TotalNetRate string          `json:"totalNetRate"`
			FlagNetRate  bool            `json:"flagNetRate"`
		}{keCoverageKabel(h.Coverages), teks(h.TotalNetRate), h.Flag})
	}
}

// badanHitungCoverage - badan POST hitung-coverage. modeDiskon kosong = diturunkan (A157).
type badanHitungCoverage struct {
	Coverage       coverageKabel `json:"coverage"`
	TSI            string        `json:"tsi"`
	Mode           string        `json:"mode"`
	ModeDiskon     string        `json:"modeDiskon"`
	IsAdjustable   bool          `json:"isAdjustable"`
	PctAdjustOther string        `json:"pctAdjustOther"`
}

// hitungCoverage - POST /api/nbfacin/kasus/{caseId}/hitung-coverage: hitung satu coverage tanpa menyimpan.
func hitungCoverage(svc *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var b badanHitungCoverage
		if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, batasBadanCoverage)).Decode(&b); err != nil {
			galat.Tulis(w, http.StatusBadRequest, "badan permintaan bukan JSON hitung coverage yang sah")
			return
		}
		var masalah []string
		c := keCoverageModel("coverage.", b.Coverage, &masalah)
		tsi := services.UraiDesimalIsian("tsi", b.TSI, &masalah)
		it := services.PenyesuaianItem{IsAdjustable: b.IsAdjustable,
			PctAdjustOther: services.UraiDesimalIsian("pctAdjustOther", b.PctAdjustOther, &masalah)}
		if err := services.GalatIsianObjek(masalah); err != nil {
			tulisGalat(w, err)
			return
		}
		hasil, err := svc.HitungCoverageKasus(r.Context(), r.PathValue("caseId"), c, tsi, b.Mode, b.ModeDiskon, it)
		if err != nil {
			tulisGalat(w, err)
			return
		}
		galat.TulisJSON(w, satuCoverageKabel(hasil))
	}
}
