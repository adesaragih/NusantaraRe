package handlers

// Tab Spreading keluarga Aneka (tiket 73 P3): GET / PUT /api/nbfacin/kasus/{caseId}/spreading-aneka, POST …/spreading-aneka/
// hitung-share, …/salin (Copy To All), …/hitung-baris. Pola endpoint tab Spreading FIRE (spreading.go); periksa-template
// FIRE dipakai ulang apa adanya.
//
// Badan PUT / hitung-baris = TampilanSpreadingAneka utuh dari layar: yang dipakai hanya percentShare, oldId coverage (pemeriksa
// bentuk lokasi / okupasi / objek / objek Golf / coverage, 409) dan treatyType / treatyName / sharePercentage baris spreading;
// kunci lain diterima dan diabaikan.

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

type coverageSpreadingAnekaKabel struct {
	OldID            string                `json:"oldId"`
	CoverageNote     string                `json:"coverageNote"`
	Loading          string                `json:"loading"`
	Rate             string                `json:"rate"`
	TSI              string                `json:"tsi"`
	Premium          string                `json:"premium"`
	Discount         string                `json:"discount"`
	LimitOfLiability string                `json:"limitOfLiability"`
	TSINusantaraRe   string                `json:"tsiNusantaraRe"`
	PremiNusantaraRe string                `json:"premiNusantaraRe"`
	Spreading        []barisSpreadingKabel `json:"spreading"`
}

type objekSpreadingAnekaKabel struct {
	ObjectName string                        `json:"objectName"`
	Currency   string                        `json:"currency"`
	TSIOld     string                        `json:"tsiOld"`
	TSI        string                        `json:"tsi"`
	Coverages  []coverageSpreadingAnekaKabel `json:"coverages"`
}

type okupasiSpreadingAnekaKabel struct {
	OccupationName string                     `json:"occupationName"`
	Objek          []objekSpreadingAnekaKabel `json:"objek"`
}

type lokasiSpreadingAnekaKabel struct {
	ObjectNo  string                       `json:"objectNo"`
	Location  string                       `json:"location"`
	Okupasi   []okupasiSpreadingAnekaKabel `json:"okupasi"`
	ObjekGolf []objekSpreadingAnekaKabel   `json:"objekGolf"`
	Total     []totalSpreadingKabel        `json:"total"`
}

type tampilanSpreadingAnekaKabel struct {
	PercentShare      string                      `json:"percentShare"`
	Treaty            []treatySpreadingKabel      `json:"treaty"`
	Template          []templateSpreadingKabel    `json:"template"`
	Lokasi            []lokasiSpreadingAnekaKabel `json:"lokasi"`
	RingkasanTreaty   []totalSpreadingKabel       `json:"ringkasanTreaty"`
	RingkasanMataUang []ringkasanMataUangKabel    `json:"ringkasanMataUang"`
	Pesan             []string                    `json:"pesan"`
	TampilTotalLokasi bool                        `json:"tampilTotalLokasi"`
}

func keObjekSpreadingAnekaKabel(objek []models.ObjekAneka) []objekSpreadingAnekaKabel {
	out := make([]objekSpreadingAnekaKabel, 0, len(objek))
	for _, o := range objek {
		ok := objekSpreadingAnekaKabel{ObjectName: o.ObjectName, Currency: o.Currency, TSIOld: teks(o.TSIOld), TSI: teks(o.TSI),
			Coverages: []coverageSpreadingAnekaKabel{}}
		for _, c := range o.Coverages {
			ck := coverageSpreadingAnekaKabel{OldID: c.OldID, CoverageNote: c.CoverageNote, Loading: teks(c.Loading), Rate: teks(c.Rate),
				TSI: teks(c.TSI), Premium: teks(c.Premium), Discount: teks(c.Discount), LimitOfLiability: teks(c.LimitOfLiability),
				TSINusantaraRe: teks(c.TSINusantaraRe), PremiNusantaraRe: teks(c.PremiNusantaraRe), Spreading: []barisSpreadingKabel{}}
			for _, b := range c.Spreading {
				ck.Spreading = append(ck.Spreading, barisSpreadingKabel{b.TreatyType, b.TreatyName, teks(b.SharePercentage),
					teks(b.TSIGrossSpreaded), teks(b.TSISpreaded), teks(b.ClaimEstimation), teks(b.PremiumSpreaded)})
			}
			ok.Coverages = append(ok.Coverages, ck)
		}
		out = append(out, ok)
	}
	return out
}

func keSpreadingAnekaKabel(d services.TampilanSpreadingAneka) tampilanSpreadingAnekaKabel {
	out := tampilanSpreadingAnekaKabel{PercentShare: teks(d.PercentShare), Treaty: []treatySpreadingKabel{},
		Template: []templateSpreadingKabel{}, Lokasi: []lokasiSpreadingAnekaKabel{},
		RingkasanTreaty: keTotalSpreadingKabel(d.RingkasanTreaty, true), RingkasanMataUang: []ringkasanMataUangKabel{}, Pesan: d.Pesan,
		TampilTotalLokasi: d.TampilTotalLokasi}
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
		lk := lokasiSpreadingAnekaKabel{ObjectNo: l.ObjectNo, Location: l.Location, Okupasi: []okupasiSpreadingAnekaKabel{},
			ObjekGolf: keObjekSpreadingAnekaKabel(l.ObjekGolf), Total: keTotalSpreadingKabel(l.Total, false)}
		for _, o := range l.Okupasi {
			lk.Okupasi = append(lk.Okupasi, okupasiSpreadingAnekaKabel{OccupationName: o.OccupationName,
				Objek: keObjekSpreadingAnekaKabel(o.Objek)})
		}
		out.Lokasi = append(out.Lokasi, lk)
	}
	for _, r := range d.RingkasanMataUang {
		out.RingkasanMataUang = append(out.RingkasanMataUang, ringkasanMataUangKabel{r.Currency, teks(r.TSITopRisk), teks(r.TSI),
			teks(r.LoL), teks(r.Premium)})
	}
	return out
}

// badanSpreadingAneka - badan POST / PUT; kunci lain diabaikan (lihat kepala berkas).
type badanSpreadingAneka struct {
	PercentShare *string                     `json:"percentShare"`
	Template     []templateSpreadingKabel    `json:"template"`
	Lokasi       []lokasiSpreadingAnekaKabel `json:"lokasi"`
}

// uraiBadanSpreadingAneka - JSON + percentShare + template; galat 400 sudah ditulis bila ok=false.
func uraiBadanSpreadingAneka(w http.ResponseWriter, r *http.Request) (b badanSpreadingAneka, ps *apd.Decimal, templat []models.TemplateSpreading, ok bool) {
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

// uraiLokasiBadanAneka - lokasi badan (wajib) -> model: oldId coverage + treatyType / treatyName / sharePercentage baris; galat
// 400 sudah ditulis bila ok=false.
func uraiLokasiBadanAneka(w http.ResponseWriter, b badanSpreadingAneka) ([]models.LokasiAneka, bool) {
	if b.Lokasi == nil {
		galat.Tulis(w, http.StatusBadRequest, "badan permintaan harus memuat lokasi")
		return nil, false
	}
	var masalah []string
	objek := func(jalur string, ok []objekSpreadingAnekaKabel) []models.ObjekAneka {
		hasil := make([]models.ObjekAneka, 0, len(ok))
		for n, o := range ok {
			om := models.ObjekAneka{Coverages: make([]models.CoverageAneka, 0, len(o.Coverages))}
			for m, c := range o.Coverages {
				cm := models.CoverageAneka{OldID: c.OldID, Spreading: make([]models.BarisSpreading, 0, len(c.Spreading))}
				for x, s := range c.Spreading {
					cm.Spreading = append(cm.Spreading, models.BarisSpreading{TreatyType: s.TreatyType, TreatyName: s.TreatyName,
						SharePercentage: services.UraiDesimalIsian(fmt.Sprintf("%s[%d].coverages[%d].spreading[%d].sharePercentage", jalur, n, m, x),
							s.SharePercentage, &masalah)})
				}
				om.Coverages = append(om.Coverages, cm)
			}
			hasil = append(hasil, om)
		}
		return hasil
	}
	lokasi := make([]models.LokasiAneka, 0, len(b.Lokasi))
	for i, l := range b.Lokasi {
		lm := models.LokasiAneka{Okupasi: make([]models.OkupasiAneka, 0, len(l.Okupasi)),
			ObjekGolf: objek(fmt.Sprintf("lokasi[%d].objekGolf", i), l.ObjekGolf)}
		for j, o := range l.Okupasi {
			lm.Okupasi = append(lm.Okupasi, models.OkupasiAneka{Objek: objek(fmt.Sprintf("lokasi[%d].okupasi[%d].objek", i, j), o.Objek)})
		}
		lokasi = append(lokasi, lm)
	}
	if len(masalah) > 0 {
		tulisGalat(w, fmt.Errorf("%w: %s", services.ErrMasukanSpreading, strings.Join(masalah, "; ")))
		return nil, false
	}
	return lokasi, true
}

// bacaSpreadingAneka - GET /api/nbfacin/kasus/{caseId}/spreading-aneka.
func bacaSpreadingAneka(svc *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		d, err := svc.TampilanSpreadingAneka(r.Context(), r.PathValue("caseId"))
		if err != nil {
			tulisGalat(w, err)
			return
		}
		galat.TulisJSON(w, keSpreadingAnekaKabel(d))
	}
}

// hitungShareSpreadingAneka - POST …/spreading-aneka/hitung-share badan {percentShare}: tidak menyimpan.
func hitungShareSpreadingAneka(svc *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, ps, _, ok := uraiBadanSpreadingAneka(w, r)
		if !ok {
			return
		}
		d, err := svc.HitungShareSpreadingAneka(r.Context(), r.PathValue("caseId"), ps)
		if err != nil {
			tulisGalat(w, err)
			return
		}
		galat.TulisJSON(w, keSpreadingAnekaKabel(d))
	}
}

// salinSpreadingAneka - POST …/spreading-aneka/salin badan {percentShare, template}: tidak menyimpan; tahap uw saja.
func salinSpreadingAneka(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pelaku := inti.PelakuDari(r, stubPelaku)
		if err := inti.WajibIdentitas(pelaku); err != nil {
			tulisGalat(w, err)
			return
		}
		_, ps, templat, ok := uraiBadanSpreadingAneka(w, r)
		if !ok {
			return
		}
		d, err := svc.SalinSpreadingAneka(r.Context(), pelaku, r.PathValue("caseId"), ps, templat)
		if err != nil {
			tulisGalat(w, err)
			return
		}
		galat.TulisJSON(w, keSpreadingAnekaKabel(d))
	}
}

// simpanSpreadingAneka - PUT …/spreading-aneka badan TampilanSpreadingAneka (lokasi wajib): hitung ulang, simpan, baca ulang.
func simpanSpreadingAneka(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b, ps, _, ok := uraiBadanSpreadingAneka(w, r)
		if !ok {
			return
		}
		lokasi, ok := uraiLokasiBadanAneka(w, b)
		if !ok {
			return
		}
		d, err := svc.SimpanSpreadingAneka(r.Context(), inti.PelakuDari(r, stubPelaku), r.PathValue("caseId"), ps, lokasi)
		if err != nil {
			tulisGalat(w, err)
			return
		}
		galat.TulisJSON(w, keSpreadingAnekaKabel(d))
	}
}

// hitungBarisSpreadingAneka - POST …/spreading-aneka/hitung-baris badan {percentShare, lokasi}: hitungan Save tahap uw tanpa
// menulis.
func hitungBarisSpreadingAneka(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pelaku := inti.PelakuDari(r, stubPelaku)
		if err := inti.WajibIdentitas(pelaku); err != nil {
			tulisGalat(w, err)
			return
		}
		b, ps, _, ok := uraiBadanSpreadingAneka(w, r)
		if !ok {
			return
		}
		lokasi, ok := uraiLokasiBadanAneka(w, b)
		if !ok {
			return
		}
		d, err := svc.HitungBarisSpreadingAneka(r.Context(), pelaku, r.PathValue("caseId"), ps, lokasi)
		if err != nil {
			tulisGalat(w, err)
			return
		}
		galat.TulisJSON(w, keSpreadingAnekaKabel(d))
	}
}
