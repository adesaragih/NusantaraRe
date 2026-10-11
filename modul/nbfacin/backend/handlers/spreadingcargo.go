package handlers

// Tab Spreading MARINE CARGO (tiket 74 P3): GET / PUT /api/nbfacin/kasus/{caseId}/spreading-cargo, POST …/spreading-cargo/
// hitung-share, …/salin (Copy To All), …/hitung-baris - pola spreading-aneka. Badan PUT / hitung-baris = TampilanSpreadingCargo
// dari layar: yang dipakai percentShare, `coverage` (kode, pemeriksa bentuk 409) dan per baris treatyType / treatyName /
// sharePercentage + tsiGrossSpreaded / claimEstimation (sel isian "TSI 100% Spreaded" / "Limit Of Liability"
// InputDtlSpreadingCoverage_FacIn L3216 / L3912 - tidak dihitung rantai Marine, jadi dibawa apa adanya); kunci lain diabaikan.

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

type coverageSpreadingCargoKabel struct {
	Coverage         string                `json:"coverage"`
	CoverageNote     string                `json:"coverageNote"`
	Currency         string                `json:"currency"`
	Rate             string                `json:"rate"`
	Premium          string                `json:"premium"`
	Discount         string                `json:"discount"`
	TSI              string                `json:"tsi"`
	TSILiability     string                `json:"tsiLiability"`
	LimitOfLiability string                `json:"limitOfLiability"`
	TSINusantaraRe   string                `json:"tsiNusantaraRe"`
	PremiNusantaraRe string                `json:"premiNusantaraRe"`
	Spreading        []barisSpreadingKabel `json:"spreading"`
}

type cargoSpreadingKabel struct {
	No             int                           `json:"no"`
	TradingNote    string                        `json:"tradingNote"`
	GoodNote       string                        `json:"goodNote"`
	PackingNote    string                        `json:"packingNote"`
	ConveyanceNote string                        `json:"conveyanceNote"`
	Coverages      []coverageSpreadingCargoKabel `json:"coverages"`
}

type tampilanSpreadingCargoKabel struct {
	PercentShare      string                   `json:"percentShare"`
	Treaty            []treatySpreadingKabel   `json:"treaty"`
	Template          []templateSpreadingKabel `json:"template"`
	Cargo             []cargoSpreadingKabel    `json:"cargo"`
	RingkasanTreaty   []totalSpreadingKabel    `json:"ringkasanTreaty"`
	RingkasanMataUang []ringkasanMataUangKabel `json:"ringkasanMataUang"`
	Pesan             []string                 `json:"pesan"`
	// LabelLimitOfLiability - ProtectSpreading.CARI15 "Total Limit of Liability(…%)" (SumTSIPremiSpreadedRNM_MARINECARGO_Act L4318);
	// kosong bila tidak ada LoL.
	LabelLimitOfLiability string `json:"labelLimitOfLiability"`
}

func keSpreadingCargoKabel(d services.TampilanSpreadingCargo) tampilanSpreadingCargoKabel {
	out := tampilanSpreadingCargoKabel{PercentShare: teks(d.PercentShare), Treaty: []treatySpreadingKabel{},
		Template: []templateSpreadingKabel{}, Cargo: []cargoSpreadingKabel{}, RingkasanTreaty: keTotalSpreadingKabel(d.RingkasanTreaty, true),
		RingkasanMataUang: []ringkasanMataUangKabel{}, Pesan: d.Pesan, LabelLimitOfLiability: d.LabelLoL}
	if out.Pesan == nil {
		out.Pesan = []string{}
	}
	for _, t := range d.Treaty {
		out.Treaty = append(out.Treaty, treatySpreadingKabel{t.ID, t.Name})
	}
	for _, t := range d.Template {
		out.Template = append(out.Template, templateSpreadingKabel{t.TreatyType, t.TreatyName, teks(t.SharePercentage)})
	}
	for _, c := range d.Cargo {
		ck := cargoSpreadingKabel{No: c.No, TradingNote: c.TradingNote, GoodNote: c.GoodNote, PackingNote: c.PackingNote,
			ConveyanceNote: c.ConveyanceNote, Coverages: []coverageSpreadingCargoKabel{}}
		for _, v := range c.Coverages {
			vk := coverageSpreadingCargoKabel{Coverage: v.Coverage, CoverageNote: v.CoverageNote, Currency: v.Currency, Rate: teks(v.Rate),
				Premium: teks(v.Premium), Discount: teks(v.Discount), TSI: teks(v.TSI), TSILiability: teks(v.TSILiability),
				LimitOfLiability: teks(v.LimitOfLiability), TSINusantaraRe: teks(v.TSINusantaraRe), PremiNusantaraRe: teks(v.PremiNusantaraRe),
				Spreading: []barisSpreadingKabel{}}
			for _, b := range v.Spreading {
				vk.Spreading = append(vk.Spreading, barisSpreadingKabel{b.TreatyType, b.TreatyName, teks(b.SharePercentage),
					teks(b.TSIGrossSpreaded), teks(b.TSISpreaded), teks(b.ClaimEstimation), teks(b.PremiumSpreaded)})
			}
			ck.Coverages = append(ck.Coverages, vk)
		}
		out.Cargo = append(out.Cargo, ck)
	}
	for _, r := range d.RingkasanMataUang {
		out.RingkasanMataUang = append(out.RingkasanMataUang, ringkasanMataUangKabel{r.Currency, teks(r.TSITopRisk), teks(r.TSI),
			teks(r.LoL), teks(r.Premium)})
	}
	return out
}

// badanSpreadingCargo - badan POST / PUT; kunci lain diabaikan.
type badanSpreadingCargo struct {
	PercentShare *string                  `json:"percentShare"`
	Template     []templateSpreadingKabel `json:"template"`
	Cargo        []cargoSpreadingKabel    `json:"cargo"`
}

// uraiBadanSpreadingCargo - JSON + percentShare + template; galat 400 sudah ditulis bila ok=false.
func uraiBadanSpreadingCargo(w http.ResponseWriter, r *http.Request) (b badanSpreadingCargo, ps *apd.Decimal, templat []models.TemplateSpreading, ok bool) {
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, batasBadanSpreading)).Decode(&b); err != nil || b.PercentShare == nil {
		galat.Tulis(w, http.StatusBadRequest, "badan permintaan harus memuat percentShare")
		return b, nil, nil, false
	}
	var masalah []string
	ps = services.UraiDesimalIsian("percentShare", *b.PercentShare, &masalah)
	templat = make([]models.TemplateSpreading, 0, len(b.Template))
	for i, t := range b.Template {
		templat = append(templat, models.TemplateSpreading{TreatyType: t.TreatyType, TreatyName: t.TreatyName,
			SharePercentage: services.UraiDesimalIsian(fmt.Sprintf("template[%d].sharePercentage", i), t.SharePercentage, &masalah)})
	}
	if len(masalah) > 0 {
		tulisGalat(w, fmt.Errorf("%w: %s", services.ErrMasukanSpreading, strings.Join(masalah, "; ")))
		return b, nil, nil, false
	}
	return b, ps, templat, true
}

// uraiCargoBadan - cargo badan (wajib) -> model: kode coverage + baris; galat 400 sudah ditulis bila ok=false.
func uraiCargoBadan(w http.ResponseWriter, b badanSpreadingCargo) ([]models.ObjekCargo, bool) {
	if b.Cargo == nil {
		galat.Tulis(w, http.StatusBadRequest, "badan permintaan harus memuat cargo")
		return nil, false
	}
	var masalah []string
	cargo := make([]models.ObjekCargo, 0, len(b.Cargo))
	for i, c := range b.Cargo {
		cm := models.ObjekCargo{Coverages: make([]models.CoverageCargo, 0, len(c.Coverages))}
		for m, v := range c.Coverages {
			vm := models.CoverageCargo{Coverage: v.Coverage, Spreading: make([]models.BarisSpreading, 0, len(v.Spreading))}
			for x, s := range v.Spreading {
				awal := fmt.Sprintf("cargo[%d].coverages[%d].spreading[%d].", i, m, x)
				vm.Spreading = append(vm.Spreading, models.BarisSpreading{TreatyType: s.TreatyType, TreatyName: s.TreatyName,
					SharePercentage:  services.UraiDesimalIsian(awal+"sharePercentage", s.SharePercentage, &masalah),
					TSIGrossSpreaded: services.UraiDesimalIsian(awal+"tsiGrossSpreaded", s.TSIGrossSpreaded, &masalah),
					ClaimEstimation:  services.UraiDesimalIsian(awal+"claimEstimation", s.ClaimEstimation, &masalah)})
			}
			cm.Coverages = append(cm.Coverages, vm)
		}
		cargo = append(cargo, cm)
	}
	if len(masalah) > 0 {
		tulisGalat(w, fmt.Errorf("%w: %s", services.ErrMasukanSpreading, strings.Join(masalah, "; ")))
		return nil, false
	}
	return cargo, true
}

// bacaSpreadingCargo - GET /api/nbfacin/kasus/{caseId}/spreading-cargo.
func bacaSpreadingCargo(svc *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		d, err := svc.TampilanSpreadingCargo(r.Context(), r.PathValue("caseId"))
		if err != nil {
			tulisGalat(w, err)
			return
		}
		galat.TulisJSON(w, keSpreadingCargoKabel(d))
	}
}

// hitungShareSpreadingCargo - POST …/spreading-cargo/hitung-share badan {percentShare}: tidak menyimpan.
func hitungShareSpreadingCargo(svc *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, ps, _, ok := uraiBadanSpreadingCargo(w, r)
		if !ok {
			return
		}
		d, err := svc.HitungShareSpreadingCargo(r.Context(), r.PathValue("caseId"), ps)
		if err != nil {
			tulisGalat(w, err)
			return
		}
		galat.TulisJSON(w, keSpreadingCargoKabel(d))
	}
}

// salinSpreadingCargo - POST …/spreading-cargo/salin badan {percentShare, template}: tidak menyimpan; tahap uw saja.
func salinSpreadingCargo(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pelaku := inti.PelakuDari(r, stubPelaku)
		if err := inti.WajibIdentitas(pelaku); err != nil {
			tulisGalat(w, err)
			return
		}
		_, ps, templat, ok := uraiBadanSpreadingCargo(w, r)
		if !ok {
			return
		}
		d, err := svc.SalinSpreadingCargo(r.Context(), pelaku, r.PathValue("caseId"), ps, templat)
		if err != nil {
			tulisGalat(w, err)
			return
		}
		galat.TulisJSON(w, keSpreadingCargoKabel(d))
	}
}

// simpanSpreadingCargo - PUT …/spreading-cargo badan TampilanSpreadingCargo (cargo wajib): hitung ulang, simpan, baca ulang.
func simpanSpreadingCargo(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b, ps, _, ok := uraiBadanSpreadingCargo(w, r)
		if !ok {
			return
		}
		cargo, ok := uraiCargoBadan(w, b)
		if !ok {
			return
		}
		d, err := svc.SimpanSpreadingCargo(r.Context(), inti.PelakuDari(r, stubPelaku), r.PathValue("caseId"), ps, cargo)
		if err != nil {
			tulisGalat(w, err)
			return
		}
		galat.TulisJSON(w, keSpreadingCargoKabel(d))
	}
}

// hitungBarisSpreadingCargo - POST …/spreading-cargo/hitung-baris badan {percentShare, cargo}: hitungan Save tahap uw tanpa menulis.
func hitungBarisSpreadingCargo(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pelaku := inti.PelakuDari(r, stubPelaku)
		if err := inti.WajibIdentitas(pelaku); err != nil {
			tulisGalat(w, err)
			return
		}
		b, ps, _, ok := uraiBadanSpreadingCargo(w, r)
		if !ok {
			return
		}
		cargo, ok := uraiCargoBadan(w, b)
		if !ok {
			return
		}
		d, err := svc.HitungBarisSpreadingCargo(r.Context(), pelaku, r.PathValue("caseId"), ps, cargo)
		if err != nil {
			tulisGalat(w, err)
			return
		}
		galat.TulisJSON(w, keSpreadingCargoKabel(d))
	}
}
