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
	// Tiket 58 - tahap Submit dan isian Deductible Acceptance Team Leader.
	Tahap                  string `json:"tahap"`
	DeductibleBedaMataUang bool   `json:"deductibleBedaMataUang"`
	DeductibleAcceptance   bool   `json:"deductibleAcceptance"`
	CommentDeductible      string `json:"commentDeductible"`
	// Tiket 61 - tahap Underwriting: Validity Of Acceptance (hari) + tanggal kabel DD-MM-YYYY; "" tahap lain.
	DaysValidity    string `json:"daysValidity"`
	DateValidity    string `json:"dateValidity"`
	WaitingBindDate string `json:"waitingBindDate"`
	// Tiket 64 - tahap Binding (false tahap lain).
	ConfirmBinding bool `json:"confirmBinding"`
	ReceivedRiSlip bool `json:"receivedRiSlip"`
}

// validitasSubmitKabel - api.ts `ValiditasSubmit` (tiket 61).
type validitasSubmitKabel struct {
	DaysValidity    string `json:"daysValidity"`
	DateValidity    string `json:"dateValidity"`
	WaitingBindDate string `json:"waitingBindDate"`
}

type hasilSubmitKabel struct {
	Galat    string   `json:"galat,omitempty"`
	Pesan    []string `json:"pesan"`
	Position string   `json:"position"`
	NBStatus string   `json:"nbStatus"`
	// PesanBinding - tiket 64b: penahan tombol "Binding" (periksa Accept tahap Binding); [] selain itu (selalu ada).
	PesanBinding []string `json:"pesanBinding"`
	// NomorPolis - tiket 69a-2: nomor polis bila Submit menyelesaikan Policy; "" selain itu (selalu ada).
	NomorPolis string `json:"nomorPolis"`
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
	out.Tahap, out.DeductibleBedaMataUang, out.DeductibleAcceptance, out.CommentDeductible = t.Tahap, t.DeductibleBedaMataUang,
		t.DeductibleAcceptance, t.CommentDeductible
	out.DaysValidity, out.DateValidity, out.WaitingBindDate = t.DaysValidity, t.DateValidity, t.WaitingBindDate
	out.ConfirmBinding, out.ReceivedRiSlip = t.ConfirmBinding, t.ReceivedRiSlip
	return out
}

func keHasilSubmitKabel(h services.HasilSubmit) hasilSubmitKabel {
	return hasilSubmitKabel{Pesan: append([]string{}, h.Pesan...), Position: h.Position, NBStatus: h.NBStatus,
		PesanBinding: append([]string{}, h.PesanBinding...), NomorPolis: h.NomorPolis}
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
	_ = json.NewEncoder(w).Encode(hasilSubmitKabel{Galat: strings.Join(g.Pesan, "; "), Pesan: g.Pesan, PesanBinding: []string{}})
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
			EmailType            string `json:"emailType"`
			DeductibleAcceptance bool   `json:"deductibleAcceptance"`
			CommentDeductible    string `json:"commentDeductible"`
			DaysValidity         string `json:"daysValidity"`
		}
		dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, batasBadanSubmit))
		dec.DisallowUnknownFields()
		if err := dec.Decode(&b); err != nil {
			galat.Tulis(w, http.StatusBadRequest, `badan permintaan harus {"emailType", "deductibleAcceptance", "commentDeductible", "daysValidity"}`)
			return
		}
		h, err := svc.PeriksaSubmit(r.Context(), pelaku, r.PathValue("caseId"), services.IsiPeriksa{EmailType: b.EmailType,
			DeductibleAcceptance: b.DeductibleAcceptance, CommentDeductible: b.CommentDeductible, DaysValidity: b.DaysValidity})
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
			// Tiket 58 - isian Team Leader (marketing: false / "").
			DeductibleAcceptance bool   `json:"deductibleAcceptance"`
			CommentDeductible    string `json:"commentDeductible"`
			// Tiket 61 - tahap Underwriting ("" tahap lain).
			DaysValidity string `json:"daysValidity"`
			// Tiket 64 - tahap Binding (false tahap lain).
			ReceivedRiSlip bool `json:"receivedRiSlip"`
		}
		dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, batasBadanSubmit))
		dec.DisallowUnknownFields()
		if err := dec.Decode(&b); err != nil {
			galat.Tulis(w, http.StatusBadRequest, `badan permintaan harus {"emailType", "comment", "bandingTo", "sebabTc", "sebabNonTc", "deductibleAcceptance", "commentDeductible", "daysValidity", "receivedRiSlip"}`)
			return
		}
		isi := models.IsiSubmit{EmailType: b.EmailType, Comment: b.Comment, BandingTo: b.BandingTo, SebabNonTc: b.SebabNonTc,
			DeductibleAcceptance: b.DeductibleAcceptance, CommentDeductible: b.CommentDeductible, DaysValidity: b.DaysValidity,
			ReceivedRiSlip: b.ReceivedRiSlip}
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

// tekanBinding - POST /api/nbfacin/kasus/{caseId}/submit/binding (tiket 64, badan {}): tombol Binding; jawab TampilanSubmit terbaru.
func tekanBinding(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pelaku := inti.PelakuDari(r, stubPelaku)
		if err := inti.WajibIdentitas(pelaku); err != nil {
			tulisGalat(w, err)
			return
		}
		var b struct{}
		dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, batasBadanSubmit))
		dec.DisallowUnknownFields()
		if err := dec.Decode(&b); err != nil {
			galat.Tulis(w, http.StatusBadRequest, "badan permintaan harus {}")
			return
		}
		t, err := svc.TekanBinding(r.Context(), pelaku, r.PathValue("caseId"))
		if err != nil {
			tulisGalatSubmit(w, err)
			return
		}
		galat.TulisJSON(w, keTampilanSubmitKabel(t))
	}
}

// hitungValiditasSubmit - POST /api/nbfacin/kasus/{caseId}/submit/validitas {daysValidity} (tiket 61, mode edit, tidak
// menyimpan).
func hitungValiditasSubmit(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pelaku := inti.PelakuDari(r, stubPelaku)
		if err := inti.WajibIdentitas(pelaku); err != nil {
			tulisGalat(w, err)
			return
		}
		var b struct {
			DaysValidity string `json:"daysValidity"`
		}
		dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, batasBadanSubmit))
		dec.DisallowUnknownFields()
		if err := dec.Decode(&b); err != nil {
			galat.Tulis(w, http.StatusBadRequest, `badan permintaan harus {"daysValidity"}`)
			return
		}
		v, err := svc.HitungValiditasSubmit(r.Context(), pelaku, r.PathValue("caseId"), b.DaysValidity)
		if err != nil {
			tulisGalat(w, err)
			return
		}
		galat.TulisJSON(w, validitasSubmitKabel(v))
	}
}
