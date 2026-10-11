package handlers

// Tab Scoring Risk kasus FIRE (tiket 54): GET / PUT /api/nbfacin/kasus/{caseId}/scoring dan POST …/scoring/hitung.
// Kontrak usulan sesi 9d (frontend sesi 0f): TampilanScoring; badan {kolom: [{indexLocation, pilihan, persenLossRatio}, …]}.

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

// batasBadanScoring - dua kolom × 24 pilihan.
const batasBadanScoring = 1 << 16

type kolomScoringKabel struct {
	OccupationNote  string            `json:"occupationNote"`
	TopRiskLocation string            `json:"topRiskLocation"`
	IndexLocation   string            `json:"indexLocation"`
	Pilihan         map[string]int    `json:"pilihan"`
	PersenLossRatio string            `json:"persenLossRatio"`
	Skor            map[string]string `json:"skor"`
	Total           string            `json:"total"`
}

type hasilScoringKabel struct {
	AcceptableNotApproval  bool `json:"acceptableNotApproval"`
	AcceptableWithApproval bool `json:"acceptableWithApproval"`
	NotAcceptable          bool `json:"notAcceptable"`
	OpDivHead              bool `json:"opDivHead"`
	OpDirector             bool `json:"opDirector"`
	TechDirector           bool `json:"techDirector"`
	PresDirector           bool `json:"presDirector"`
	DueOccupation          bool `json:"dueOccupation"`
	DueObjectConditions    bool `json:"dueObjectConditions"`
	DueLossRatio           bool `json:"dueLossRatio"`
	DueOthers              bool `json:"dueOthers"`
}

type lokasiDominanKabel struct {
	IndexLocation string `json:"indexLocation"`
	Label         string `json:"label"`
}

type tampilanScoringKabel struct {
	BolehUbah      bool                 `json:"bolehUbah"`
	InsuredName    string               `json:"insuredName"`
	Status         string               `json:"status"`
	Mulai          string               `json:"mulai"`
	Selesai        string               `json:"selesai"`
	Kolom          []kolomScoringKabel  `json:"kolom"`
	FinalScore     string               `json:"finalScore"`
	NoteFinalScore string               `json:"noteFinalScore"`
	Hasil          hasilScoringKabel    `json:"hasil"`
	LokasiDominan  []lokasiDominanKabel `json:"lokasiDominan"`
	Pesan          []string             `json:"pesan"`
}

// keKolomScoringKabel - 24 kunci selalu ada di pilihan / skor; kolom yang tidak dihitung: skor dan total kosong.
func keKolomScoringKabel(k models.KolomScoring) kolomScoringKabel {
	out := kolomScoringKabel{OccupationNote: k.OccupationNote, TopRiskLocation: k.TopRiskLocation, IndexLocation: k.IndexLocation,
		Pilihan: map[string]int{}, PersenLossRatio: teks(k.PersenLossRatio), Skor: map[string]string{}, Total: teks(k.Total)}
	for _, f := range models.DaftarFaktorScoring {
		out.Pilihan[f.Kunci] = k.Pilihan[f.Kunci]
		out.Skor[f.Kunci] = teks(k.Skor[f.Kunci])
	}
	return out
}

func keScoringKabel(t services.TampilanScoring) tampilanScoringKabel {
	h := t.Scoring.Hasil
	out := tampilanScoringKabel{BolehUbah: t.BolehUbah, InsuredName: t.InsuredName, Status: t.Status, Mulai: t.Mulai,
		Selesai: t.Selesai, FinalScore: teks(t.Scoring.FinalScore), NoteFinalScore: t.Scoring.NoteFinalScore,
		Hasil: hasilScoringKabel{AcceptableNotApproval: h.AcceptableNotApproval, AcceptableWithApproval: h.AcceptableWithApproval,
			NotAcceptable: h.NotAcceptable, OpDivHead: h.OperationalDivHead, OpDirector: h.OperationalDirector,
			TechDirector: h.TechnicalDirector, PresDirector: h.PresidentDirector, DueOccupation: h.Occupation,
			DueObjectConditions: h.ObjectConditions, DueLossRatio: h.LossRatio, DueOthers: h.Others},
		LokasiDominan: []lokasiDominanKabel{}, Pesan: t.Pesan}
	for _, k := range t.Scoring.Kolom {
		out.Kolom = append(out.Kolom, keKolomScoringKabel(k))
	}
	for _, l := range t.LokasiDominan {
		out.LokasiDominan = append(out.LokasiDominan, lokasiDominanKabel{l.IndexLocation, l.Label})
	}
	if out.Pesan == nil {
		out.Pesan = []string{}
	}
	return out
}

// uraiMasukanScoring - badan ketat {kolom: [1..2 × {indexLocation, pilihan, persenLossRatio}]}; persen = teks desimal.
func uraiMasukanScoring(w http.ResponseWriter, r *http.Request) ([2]models.KolomScoring, bool) {
	var masukan [2]models.KolomScoring
	var b struct {
		Kolom []struct {
			IndexLocation   string         `json:"indexLocation"`
			Pilihan         map[string]int `json:"pilihan"`
			PersenLossRatio string         `json:"persenLossRatio"`
		} `json:"kolom"`
	}
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, batasBadanScoring))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&b); err != nil || len(b.Kolom) < 1 || len(b.Kolom) > 2 {
		galat.Tulis(w, http.StatusBadRequest, `badan permintaan harus {"kolom": [1..2 × {"indexLocation", "pilihan": {kunci: nomor}, "persenLossRatio"}]}`)
		return masukan, false
	}
	var masalah []string
	for i, k := range b.Kolom {
		masukan[i] = models.KolomScoring{IndexLocation: k.IndexLocation, Pilihan: k.Pilihan,
			PersenLossRatio: services.UraiDesimalIsian(fmt.Sprintf("kolom[%d].persenLossRatio", i), k.PersenLossRatio, &masalah)}
	}
	if len(masalah) > 0 {
		tulisGalat(w, fmt.Errorf("%w: %s", services.ErrMasukanScoring, strings.Join(masalah, "; ")))
		return masukan, false
	}
	return masukan, true
}

// bacaScoring - GET /api/nbfacin/kasus/{caseId}/scoring.
func bacaScoring(svc *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t, err := svc.TampilanScoringKasus(r.Context(), r.PathValue("caseId"))
		if err != nil {
			tulisGalat(w, err)
			return
		}
		galat.TulisJSON(w, keScoringKabel(t))
	}
}

// hitungScoring - POST /api/nbfacin/kasus/{caseId}/scoring/hitung (tanpa simpan).
func hitungScoring(svc *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		masukan, ok := uraiMasukanScoring(w, r)
		if !ok {
			return
		}
		t, err := svc.HitungScoring(r.Context(), r.PathValue("caseId"), masukan)
		if err != nil {
			tulisGalat(w, err)
			return
		}
		galat.TulisJSON(w, keScoringKabel(t))
	}
}

// simpanScoring - PUT /api/nbfacin/kasus/{caseId}/scoring.
func simpanScoring(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		masukan, ok := uraiMasukanScoring(w, r)
		if !ok {
			return
		}
		t, err := svc.SimpanScoring(r.Context(), inti.PelakuDari(r, stubPelaku), r.PathValue("caseId"), masukan)
		if err != nil {
			tulisGalat(w, err)
			return
		}
		galat.TulisJSON(w, keScoringKabel(t))
	}
}
