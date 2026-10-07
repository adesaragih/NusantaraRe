package handlers

// Tab Spreading kasus FIRE (tiket 48): GET / PUT /api/nbfacin/kasus/{caseId}/spreading, POST …/spreading/hitung-share,
// POST …/spreading/salin. Kontrak frontend modul/nbfacin/frontend/api.ts `TampilanSpreading`.
//
// Badan PUT = TampilanSpreading utuh dari layar: kunci tambahan (oldId, rate, tsi, total[], treaty[], ringkasan…) DITERIMA
// dan diabaikan - yang dipakai hanya percentShare, oldId coverage (pemeriksa bentuk, 409) dan treatyType / treatyName /
// sharePercentage baris spreading. Nilai hitungan dihitung ulang server.

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

// batasBadanSpreading - badan PUT berisi seluruh lokasi beserta angka tampilnya.
const batasBadanSpreading = 8 << 20

type treatySpreadingKabel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type templateSpreadingKabel struct {
	TreatyType      string `json:"treatyType"`
	TreatyName      string `json:"treatyName"`
	SharePercentage string `json:"sharePercentage"`
}

type barisSpreadingKabel struct {
	TreatyType       string `json:"treatyType"`
	TreatyName       string `json:"treatyName"`
	SharePercentage  string `json:"sharePercentage"`
	TSIGrossSpreaded string `json:"tsiGrossSpreaded"`
	TSISpreaded      string `json:"tsiSpreaded"`
	ClaimEstimation  string `json:"claimEstimation"`
	PremiumSpreaded  string `json:"premiumSpreaded"`
}

type coverageSpreadingKabel struct {
	OldID            string                `json:"oldId"`
	CoverageBasis    string                `json:"coverageBasis"`
	Rate             string                `json:"rate"`
	Premium          string                `json:"premium"`
	Discount         string                `json:"discount"`
	TSI              string                `json:"tsi"`
	TSILiability     string                `json:"tsiLiability"`
	TSINusantaraRe   string                `json:"tsiNusantaraRe"`
	PremiNusantaraRe string                `json:"premiNusantaraRe"`
	Spreading        []barisSpreadingKabel `json:"spreading"`
}

type itemSpreadingKabel struct {
	ItemType        string                   `json:"itemType"`
	Currency        string                   `json:"currency"`
	TSI             string                   `json:"tsi"`
	TotalGrossPremi string                   `json:"totalGrossPremi"`
	TotalPremiumRnm string                   `json:"totalPremiumRnm"`
	Coverages       []coverageSpreadingKabel `json:"coverages"`
}

type totalSpreadingKabel struct {
	Currency        string  `json:"currency"`
	TreatyType      string  `json:"treatyType"`
	TreatyName      string  `json:"treatyName"`
	SharePercentage string  `json:"sharePercentage"`
	ClaimSpreaded   string  `json:"claimSpreaded"`
	TSISpreaded     string  `json:"tsiSpreaded"`
	ClaimEstimation string  `json:"claimEstimation"`
	ClaimAmountIDR  *string `json:"claimAmountIdr,omitempty"`
	PremiumSpreaded string  `json:"premiumSpreaded"`
}

type lokasiSpreadingKabel struct {
	ObjectNo   string                `json:"objectNo"`
	ObjectName string                `json:"objectName"`
	Location   string                `json:"location"`
	Items      []itemSpreadingKabel  `json:"items"`
	Total      []totalSpreadingKabel `json:"total"`
}

type ringkasanMataUangKabel struct {
	Currency   string `json:"currency"`
	TSITopRisk string `json:"tsiTopRisk"`
	TSI        string `json:"tsi"`
	LoL        string `json:"lol"`
	Premium    string `json:"premium"`
}

type tampilanSpreadingKabel struct {
	PercentShare      string                   `json:"percentShare"`
	Treaty            []treatySpreadingKabel   `json:"treaty"`
	Template          []templateSpreadingKabel `json:"template"`
	Lokasi            []lokasiSpreadingKabel   `json:"lokasi"`
	RingkasanTreaty   []totalSpreadingKabel    `json:"ringkasanTreaty"`
	RingkasanMataUang []ringkasanMataUangKabel `json:"ringkasanMataUang"`
	Pesan             []string                 `json:"pesan"`
}

func keTotalSpreadingKabel(d []models.TotalSpreading, denganFL bool) []totalSpreadingKabel {
	out := make([]totalSpreadingKabel, 0, len(d))
	for _, t := range d {
		k := totalSpreadingKabel{Currency: t.Currency, TreatyType: t.TreatyType, TreatyName: t.TreatyName,
			SharePercentage: teks(t.SharePercentage), ClaimSpreaded: teks(t.ClaimSpreaded), TSISpreaded: teks(t.TSISpreaded),
			ClaimEstimation: teks(t.ClaimEstimation), PremiumSpreaded: teks(t.PremiumSpreaded)}
		if denganFL {
			fl := teks(t.ClaimAmountIDR)
			k.ClaimAmountIDR = &fl
		}
		out = append(out, k)
	}
	return out
}

func keSpreadingKabel(d services.TampilanSpreading) tampilanSpreadingKabel {
	out := tampilanSpreadingKabel{PercentShare: teks(d.PercentShare), Treaty: []treatySpreadingKabel{},
		Template: []templateSpreadingKabel{}, Lokasi: []lokasiSpreadingKabel{},
		RingkasanTreaty: keTotalSpreadingKabel(d.RingkasanTreaty, true), RingkasanMataUang: []ringkasanMataUangKabel{}, Pesan: d.Pesan}
	if out.Pesan == nil {
		out.Pesan = []string{}
	}
	for _, t := range d.Treaty {
		out.Treaty = append(out.Treaty, treatySpreadingKabel{t.ID, t.Name})
	}
	for _, t := range d.Template {
		out.Template = append(out.Template, templateSpreadingKabel{t.TreatyType, t.TreatyName, teks(t.SharePercentage)})
	}
	for _, l := range d.Lokasi {
		lk := lokasiSpreadingKabel{ObjectNo: l.ObjectNo, ObjectName: l.ObjectName, Location: l.Location,
			Items: []itemSpreadingKabel{}, Total: keTotalSpreadingKabel(l.Total, false)}
		for _, it := range l.Items {
			ik := itemSpreadingKabel{ItemType: it.ItemType, Currency: it.Currency, TSI: teks(it.TSI),
				TotalGrossPremi: teks(it.TotalGrossPremi), TotalPremiumRnm: teks(it.TotalPremiumNusantaraRe),
				Coverages: []coverageSpreadingKabel{}}
			for _, c := range it.Coverages {
				ck := coverageSpreadingKabel{OldID: c.OldID, CoverageBasis: c.CoverageBasis, Rate: teks(c.Rate), Premium: teks(c.Premium),
					Discount: teks(c.Discount), TSI: teks(c.TSI), TSILiability: teks(c.TSILiability),
					TSINusantaraRe: teks(c.TSINusantaraRe), PremiNusantaraRe: teks(c.PremiNusantaraRe), Spreading: []barisSpreadingKabel{}}
				for _, b := range c.Spreading {
					ck.Spreading = append(ck.Spreading, barisSpreadingKabel{b.TreatyType, b.TreatyName, teks(b.SharePercentage),
						teks(b.TSIGrossSpreaded), teks(b.TSISpreaded), teks(b.ClaimEstimation), teks(b.PremiumSpreaded)})
				}
				ik.Coverages = append(ik.Coverages, ck)
			}
			lk.Items = append(lk.Items, ik)
		}
		out.Lokasi = append(out.Lokasi, lk)
	}
	for _, r := range d.RingkasanMataUang {
		out.RingkasanMataUang = append(out.RingkasanMataUang, ringkasanMataUangKabel{r.Currency, teks(r.TSITopRisk), teks(r.TSI),
			teks(r.LoL), teks(r.Premium)})
	}
	return out
}

// badanSpreading - badan POST / PUT; kunci lain diabaikan (lihat kepala berkas).
type badanSpreading struct {
	PercentShare *string                  `json:"percentShare"`
	Template     []templateSpreadingKabel `json:"template"`
	Lokasi       []lokasiSpreadingKabel   `json:"lokasi"`
}

// uraiBadanSpreading - JSON + percentShare + template; galat 400 sudah ditulis bila ok=false.
func uraiBadanSpreading(w http.ResponseWriter, r *http.Request) (b badanSpreading, ps *apd.Decimal, templat []models.TemplateSpreading, ok bool) {
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

// bacaSpreading - GET /api/nbfacin/kasus/{caseId}/spreading.
func bacaSpreading(svc *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		d, err := svc.TampilanSpreading(r.Context(), r.PathValue("caseId"))
		if err != nil {
			tulisGalat(w, err)
			return
		}
		galat.TulisJSON(w, keSpreadingKabel(d))
	}
}

// hitungShareSpreading - POST …/spreading/hitung-share badan {percentShare}: tidak menyimpan.
func hitungShareSpreading(svc *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, ps, _, ok := uraiBadanSpreading(w, r)
		if !ok {
			return
		}
		d, err := svc.HitungShareSpreading(r.Context(), r.PathValue("caseId"), ps)
		if err != nil {
			tulisGalat(w, err)
			return
		}
		galat.TulisJSON(w, keSpreadingKabel(d))
	}
}

// salinSpreading - POST …/spreading/salin badan {percentShare, template}: tidak menyimpan.
func salinSpreading(svc *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, ps, templat, ok := uraiBadanSpreading(w, r)
		if !ok {
			return
		}
		d, err := svc.SalinSpreading(r.Context(), r.PathValue("caseId"), ps, templat)
		if err != nil {
			tulisGalat(w, err)
			return
		}
		galat.TulisJSON(w, keSpreadingKabel(d))
	}
}

// simpanSpreading - PUT …/spreading badan TampilanSpreading (lokasi wajib): hitung ulang, simpan, jawab baca ulang.
func simpanSpreading(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b, ps, _, ok := uraiBadanSpreading(w, r)
		if !ok {
			return
		}
		if b.Lokasi == nil {
			galat.Tulis(w, http.StatusBadRequest, "badan permintaan harus memuat lokasi")
			return
		}
		var masalah []string
		lokasi := make([]models.LokasiSpreading, 0, len(b.Lokasi))
		for i, l := range b.Lokasi {
			lm := models.LokasiSpreading{Items: make([]models.ItemSpreading, 0, len(l.Items))}
			for j, it := range l.Items {
				im := models.ItemSpreading{Coverages: make([]models.CoverageSpreading, 0, len(it.Coverages))}
				for m, c := range it.Coverages {
					cm := models.CoverageSpreading{OldID: c.OldID, Spreading: make([]models.BarisSpreading, 0, len(c.Spreading))}
					for n, s := range c.Spreading {
						jalur := fmt.Sprintf("lokasi[%d].items[%d].coverages[%d].spreading[%d].sharePercentage", i, j, m, n)
						cm.Spreading = append(cm.Spreading, models.BarisSpreading{TreatyType: s.TreatyType, TreatyName: s.TreatyName,
							SharePercentage: services.UraiDesimalIsian(jalur, s.SharePercentage, &masalah)})
					}
					im.Coverages = append(im.Coverages, cm)
				}
				lm.Items = append(lm.Items, im)
			}
			lokasi = append(lokasi, lm)
		}
		if len(masalah) > 0 {
			tulisGalat(w, fmt.Errorf("%w: %s", services.ErrMasukanSpreading, strings.Join(masalah, "; ")))
			return
		}
		d, err := svc.SimpanSpreading(r.Context(), inti.PelakuDari(r, stubPelaku), r.PathValue("caseId"), ps, lokasi)
		if err != nil {
			tulisGalat(w, err)
			return
		}
		galat.TulisJSON(w, keSpreadingKabel(d))
	}
}
