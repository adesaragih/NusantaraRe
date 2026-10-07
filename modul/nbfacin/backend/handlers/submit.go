package handlers

// Submit tahap marketing (tiket 57): GET /api/nbfacin/kasus/{caseId}/submit, POST …/submit/periksa, POST …/submit.
// Kontrak frontend api.ts `TampilanSubmit`, `IsiSubmit`, `HasilSubmit`. 422 proteksi Accept = `{galat, pesan, position,
// nbStatus}`: `pesan` larik verbatim, `galat` gabungannya (amplop galat bersama dibaca klien untuk setiap jawaban tidak ok).

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/services"
)

// batasBadanSubmit - Comment + Cause of Decline.
const batasBadanSubmit = 1 << 20

type riwayatAkseptasiKabel struct {
	Tanggal  string `json:"tanggal"`
	PIC      string `json:"pic"`
	Approval string `json:"approval"`
	Status   string `json:"status"`
	Comment  string `json:"comment"`
}

type barisSebabTcKabel struct {
	TermCondition string   `json:"termCondition"`
	Opsi          []string `json:"opsi"`
}

type pilihanSebabTcKabel struct {
	TermCondition string `json:"termCondition"`
	Prioritas     string `json:"prioritas"`
	Description1  string `json:"description1"`
	Description2  string `json:"description2"`
	PctRNM        string `json:"pctRnm"`
	PctSOB        string `json:"pctSob"`
}

type tampilanSubmitKabel struct {
	BolehSubmit    bool                    `json:"bolehSubmit"`
	Alasan         string                  `json:"alasan"`
	Riwayat        []riwayatAkseptasiKabel `json:"riwayat"`
	OpsiBanding    []string                `json:"opsiBanding"`
	SebabTc        []barisSebabTcKabel     `json:"sebabTc"`
	SebabNonTc     []string                `json:"sebabNonTc"`
	PctRnmOtomatis struct {
		Rate   string `json:"rate"`
		RIComm string `json:"riComm"`
	} `json:"pctRnmOtomatis"`
}

type hasilSubmitKabel struct {
	Galat    string   `json:"galat,omitempty"`
	Pesan    []string `json:"pesan"`
	Position string   `json:"position"`
	NBStatus string   `json:"nbStatus"`
}

func keTampilanSubmitKabel(t services.TampilanSubmit) tampilanSubmitKabel {
	out := tampilanSubmitKabel{BolehSubmit: t.BolehSubmit, Alasan: t.Alasan, Riwayat: []riwayatAkseptasiKabel{},
		OpsiBanding: append([]string{}, t.OpsiBanding...), SebabTc: []barisSebabTcKabel{}, SebabNonTc: append([]string{}, t.SebabNonTc...)}
	for _, r := range t.Riwayat {
		out.Riwayat = append(out.Riwayat, riwayatAkseptasiKabel(r))
	}
	for _, s := range t.SebabTc {
		out.SebabTc = append(out.SebabTc, barisSebabTcKabel{TermCondition: s.TermCondition, Opsi: append([]string{}, s.Opsi...)})
	}
	out.PctRnmOtomatis.Rate, out.PctRnmOtomatis.RIComm = t.PctRnmRate, t.PctRnmRIComm
	return out
}

func keHasilSubmitKabel(h services.HasilSubmit) hasilSubmitKabel {
	return hasilSubmitKabel{Pesan: append([]string{}, h.Pesan...), Position: h.Position, NBStatus: h.NBStatus}
}

// tulisGalatSubmit - 422 proteksi dengan pesan larik; galat lain lewat tulisGalat.
func tulisGalatSubmit(w http.ResponseWriter, err error) {
	var g services.GalatProteksi
	if !errors.As(err, &g) {
		tulisGalat(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusUnprocessableEntity)
	_ = json.NewEncoder(w).Encode(hasilSubmitKabel{Galat: strings.Join(g.Pesan, "; "), Pesan: g.Pesan})
}

// tampilanSubmit - GET /api/nbfacin/kasus/{caseId}/submit.
func tampilanSubmit(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t, err := svc.TampilanSubmit(r.Context(), inti.PelakuDari(r, stubPelaku), r.PathValue("caseId"))
		if err != nil {
			tulisGalat(w, err)
			return
		}
		galat.TulisJSON(w, keTampilanSubmitKabel(t))
	}
}

// periksaSubmit - POST /api/nbfacin/kasus/{caseId}/submit/periksa {emailType}.
func periksaSubmit(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pelaku := inti.PelakuDari(r, stubPelaku)
		if err := inti.WajibIdentitas(pelaku); err != nil {
			tulisGalat(w, err) // 401 sebelum badan diurai
			return
		}
		var b struct {
			EmailType string `json:"emailType"`
		}
		dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, batasBadanSubmit))
		dec.DisallowUnknownFields()
		if err := dec.Decode(&b); err != nil {
			galat.Tulis(w, http.StatusBadRequest, `badan permintaan harus {"emailType"}`)
			return
		}
		h, err := svc.PeriksaSubmit(r.Context(), pelaku, r.PathValue("caseId"), b.EmailType)
		if err != nil {
			tulisGalat(w, err)
			return
		}
		galat.TulisJSON(w, keHasilSubmitKabel(h))
	}
}

// ajukanSubmit - POST /api/nbfacin/kasus/{caseId}/submit (IsiSubmit ketat).
func ajukanSubmit(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pelaku := inti.PelakuDari(r, stubPelaku)
		if err := inti.WajibIdentitas(pelaku); err != nil {
			tulisGalat(w, err) // 401 sebelum badan diurai
			return
		}
		var b struct {
			EmailType  string                `json:"emailType"`
			Comment    string                `json:"comment"`
			BandingTo  string                `json:"bandingTo"`
			SebabTc    []pilihanSebabTcKabel `json:"sebabTc"`
			SebabNonTc string                `json:"sebabNonTc"`
		}
		dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, batasBadanSubmit))
		dec.DisallowUnknownFields()
		if err := dec.Decode(&b); err != nil {
			galat.Tulis(w, http.StatusBadRequest, `badan permintaan harus {"emailType", "comment", "bandingTo", "sebabTc", "sebabNonTc"}`)
			return
		}
		isi := models.IsiSubmit{EmailType: b.EmailType, Comment: b.Comment, BandingTo: b.BandingTo, SebabNonTc: b.SebabNonTc}
		for _, s := range b.SebabTc {
			isi.SebabTc = append(isi.SebabTc, models.PilihanSebabTc(s))
		}
		h, err := svc.AjukanSubmit(r.Context(), pelaku, r.PathValue("caseId"), isi)
		if err != nil {
			tulisGalatSubmit(w, err)
			return
		}
		galat.TulisJSON(w, keHasilSubmitKabel(h))
	}
}
