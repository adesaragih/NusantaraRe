package handlers

// Tab utama Loss Record kasus FIRE (tiket 55): GET / PUT /api/nbfacin/kasus/{caseId}/loss-record, POST …/loss-record/hitung,
// GET /api/nbfacin/class-of-business/saran. Kontrak frontend modul/nbfacin/frontend/api.ts `TampilanLossRecord`.

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/services"
)

// batasBadanLossRecord - 500 baris × ±4,6 KB.
const batasBadanLossRecord = 4 << 20

type barisLossRecordKabel struct {
	UwYear          string `json:"uwYear"`
	ClassOfBusiness string `json:"classOfBusiness"`
	DateOfLoss      string `json:"dateOfLoss"`
	CauseOfLoss     string `json:"causeOfLoss"`
	Currency        string `json:"currency"`
	LossAmount      string `json:"lossAmount"`
	LossAmountIdr   string `json:"lossAmountIdr"`
}

type totalLossRecordKabel struct {
	ClassOfBusiness string `json:"classOfBusiness"`
	Currency        string `json:"currency"`
	LossAmount      string `json:"lossAmount"`
	LossAmountIdr   string `json:"lossAmountIdr"`
}

type rasioLossRecordKabel struct {
	Gross1 string `json:"gross1"`
	Gross3 string `json:"gross3"`
	Gross5 string `json:"gross5"`
	Net1   string `json:"net1"`
	Net3   string `json:"net3"`
	Net5   string `json:"net5"`
}

type tampilanLossRecordKabel struct {
	PremiPar   string                 `json:"premiPar"`
	PremiMb    string                 `json:"premiMb"`
	PremiPl    string                 `json:"premiPl"`
	PremiOther string                 `json:"premiOther"`
	Baris      []barisLossRecordKabel `json:"baris"`
	TotalCob   []totalLossRecordKabel `json:"totalCob"`
	Rasio      rasioLossRecordKabel   `json:"rasio"`
	Pesan      []string               `json:"pesan"`
}

func keLossRecordKabel(t services.TampilanLossRecord) tampilanLossRecordKabel {
	l, r := t.LossRecord, t.LossRecord.Rasio
	out := tampilanLossRecordKabel{PremiPar: teks(l.PremiPar), PremiMb: teks(l.PremiMb), PremiPl: teks(l.PremiPl),
		PremiOther: teks(l.PremiOther), Baris: []barisLossRecordKabel{}, TotalCob: []totalLossRecordKabel{},
		Rasio: rasioLossRecordKabel{teks(r.Gross1), teks(r.Gross3), teks(r.Gross5), teks(r.Net1), teks(r.Net3), teks(r.Net5)},
		Pesan: t.Pesan}
	for _, b := range l.Baris {
		out.Baris = append(out.Baris, barisLossRecordKabel{b.UwYear, b.ClassOfBusiness, b.DateOfLoss, b.CauseOfLoss, b.Currency,
			teks(b.LossAmount), teks(b.LossAmountIDR)})
	}
	for _, c := range t.TotalCob {
		out.TotalCob = append(out.TotalCob, totalLossRecordKabel{c.ClassOfBusiness, c.Currency, teks(c.LossAmount), teks(c.LossAmountIDR)})
	}
	if out.Pesan == nil {
		out.Pesan = []string{}
	}
	return out
}

// uraiMasukanLossRecord - badan ketat {premiPar, premiMb, premiPl, premiOther, baris: [{classOfBusiness, dateOfLoss,
// causeOfLoss, currency, lossAmount}]}; angka = teks desimal (≤ 8 desimal), 400 ber-jalur.
func uraiMasukanLossRecord(w http.ResponseWriter, r *http.Request) (models.LossRecord, bool) {
	var b struct {
		PremiPar   string `json:"premiPar"`
		PremiMb    string `json:"premiMb"`
		PremiPl    string `json:"premiPl"`
		PremiOther string `json:"premiOther"`
		Baris      []struct {
			ClassOfBusiness string `json:"classOfBusiness"`
			DateOfLoss      string `json:"dateOfLoss"`
			CauseOfLoss     string `json:"causeOfLoss"`
			Currency        string `json:"currency"`
			LossAmount      string `json:"lossAmount"`
		} `json:"baris"`
	}
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, batasBadanLossRecord))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&b); err != nil {
		galat.Tulis(w, http.StatusBadRequest, `badan permintaan harus {"premiPar", "premiMb", "premiPl", "premiOther", "baris": [{"classOfBusiness", "dateOfLoss", "causeOfLoss", "currency", "lossAmount"}]}`)
		return models.LossRecord{}, false
	}
	var masalah []string
	u := services.UraiDesimalIsian
	m := models.LossRecord{PremiPar: u("premiPar", b.PremiPar, &masalah), PremiMb: u("premiMb", b.PremiMb, &masalah),
		PremiPl: u("premiPl", b.PremiPl, &masalah), PremiOther: u("premiOther", b.PremiOther, &masalah), Baris: []models.BarisLossRecord{}}
	for i, x := range b.Baris {
		m.Baris = append(m.Baris, models.BarisLossRecord{ClassOfBusiness: x.ClassOfBusiness, DateOfLoss: x.DateOfLoss,
			CauseOfLoss: x.CauseOfLoss, Currency: x.Currency, LossAmount: u(fmt.Sprintf("baris[%d].lossAmount", i), x.LossAmount, &masalah)})
	}
	if len(masalah) > 0 {
		tulisGalat(w, fmt.Errorf("%w: %s", services.ErrMasukanLossRecord, strings.Join(masalah, "; ")))
		return m, false
	}
	return m, true
}

// bacaLossRecord - GET /api/nbfacin/kasus/{caseId}/loss-record.
func bacaLossRecord(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t, err := svc.TampilanLossRecordKasus(r.Context(), inti.PelakuDari(r, stubPelaku), r.PathValue("caseId"))
		if err != nil {
			tulisGalat(w, err)
			return
		}
		galat.TulisJSON(w, keLossRecordKabel(t))
	}
}

// hitungLossRecord - POST /api/nbfacin/kasus/{caseId}/loss-record/hitung (tanpa simpan).
func hitungLossRecord(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pelaku := inti.PelakuDari(r, stubPelaku)
		if err := svc.WajibAksesLossRecord(r.Context(), pelaku); err != nil {
			tulisGalat(w, err) // 401 / 403 sebelum badan diurai
			return
		}
		m, ok := uraiMasukanLossRecord(w, r)
		if !ok {
			return
		}
		t, err := svc.HitungLossRecord(r.Context(), pelaku, r.PathValue("caseId"), m)
		if err != nil {
			tulisGalat(w, err)
			return
		}
		galat.TulisJSON(w, keLossRecordKabel(t))
	}
}

// simpanLossRecord - PUT /api/nbfacin/kasus/{caseId}/loss-record.
func simpanLossRecord(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pelaku := inti.PelakuDari(r, stubPelaku)
		if err := svc.WajibAksesLossRecord(r.Context(), pelaku); err != nil {
			tulisGalat(w, err) // 401 / 403 sebelum badan diurai
			return
		}
		m, ok := uraiMasukanLossRecord(w, r)
		if !ok {
			return
		}
		t, err := svc.SimpanLossRecord(r.Context(), pelaku, r.PathValue("caseId"), m)
		if err != nil {
			tulisGalat(w, err)
			return
		}
		galat.TulisJSON(w, keLossRecordKabel(t))
	}
}

// saranKelasBisnis - GET /api/nbfacin/class-of-business/saran?q= -> {baris: [{id, label}]}.
func saranKelasBisnis(svc *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		d, err := svc.SaranKelasBisnis(r.Context(), r.URL.Query().Get("q"))
		tulisSaran(w, d, err)
	}
}
