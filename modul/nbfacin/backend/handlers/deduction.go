package handlers

// Tab Premium Deduction / Brokerage Fee kasus FIRE (tiket 50): GET / PUT /api/nbfacin/kasus/{caseId}/deduction dan
// POST …/deduction/hitung. Kontrak frontend modul/nbfacin/frontend/api.ts `TampilanDeduction` / `BarisDeduction` /
// `AksiDeduction`.

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

// batasBadanDeduction - batas bita badan POST / PUT (jumlah baris dibatasi services, banyakMataUang).
const batasBadanDeduction = 1 << 20

type barisDeductionKabel struct {
	Currency     string `json:"currency"`
	Premium      string `json:"premium"`
	Commision    string `json:"commision"`
	Deduction2   string `json:"deduction2"`
	BrokerageFee string `json:"brokerageFee"`
	PPh          string `json:"pph"`
	PPN          string `json:"ppn"`
}

type tampilanDeductionKabel struct {
	PctPremiumDeduction string                `json:"pctPremiumDeduction"`
	PctBrokerageFee     string                `json:"pctBrokerageFee"`
	PPnCheck            bool                  `json:"ppnCheck"`
	TarifPph            string                `json:"tarifPph"`
	TarifPpn            string                `json:"tarifPpn"`
	MataUang            []barisDeductionKabel `json:"mataUang"`
	Pesan               []string              `json:"pesan"`
}

func keDeductionKabel(d services.TampilanDeduction) tampilanDeductionKabel {
	out := tampilanDeductionKabel{PctPremiumDeduction: teks(d.PctPremiumDeduction), PctBrokerageFee: teks(d.PctBrokerageFee),
		PPnCheck: d.PPnCheck, TarifPph: teks(d.TarifPPH), TarifPpn: teks(d.TarifPPN), MataUang: []barisDeductionKabel{},
		Pesan: append([]string{}, d.Pesan...)}
	for _, b := range d.MataUang {
		out.MataUang = append(out.MataUang, barisDeductionKabel{b.Currency, teks(b.Premium), teks(b.Commision), teks(b.Deduction2),
			teks(b.BrokerageFee), teks(b.PPh), teks(b.PPN)})
	}
	return out
}

// isianBarisKabel - baris badan: hanya kunci yang dikirim `badanDeduction` (Premium / PPh / PPN dihitung server).
type isianBarisKabel struct {
	Currency     string `json:"currency"`
	Commision    string `json:"commision"`
	Deduction2   string `json:"deduction2"`
	BrokerageFee string `json:"brokerageFee"`
}

type isianDeductionKabel struct {
	Aksi                *string           `json:"aksi"`
	Idx                 *int              `json:"idx"`
	PctPremiumDeduction string            `json:"pctPremiumDeduction"`
	PctBrokerageFee     string            `json:"pctBrokerageFee"`
	PPnCheck            *bool             `json:"ppnCheck"`
	MataUang            []isianBarisKabel `json:"mataUang"`
}

// uraiIsianDeduction - JSON (kunci asing -> 400) + angka; galat sudah ditulis bila ok=false. ppnCheck diterima tetapi
// diabaikan: server menyetel true (SetRIComm_Act langkah 5).
func uraiIsianDeduction(w http.ResponseWriter, r *http.Request, denganAksi bool) (isianDeductionKabel, services.IsianDeduction, bool) {
	var b isianDeductionKabel
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, batasBadanDeduction))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&b); err != nil || b.MataUang == nil || (denganAksi && (b.Aksi == nil || b.Idx == nil)) ||
		(!denganAksi && (b.Aksi != nil || b.Idx != nil)) {
		mau := `{"pctPremiumDeduction", "pctBrokerageFee", "ppnCheck", "mataUang": [{currency, commision, deduction2, brokerageFee}]}`
		if denganAksi {
			mau = `{"aksi", "idx", ` + mau[1:]
		}
		galat.Tulis(w, http.StatusBadRequest, "badan permintaan harus "+mau)
		return b, services.IsianDeduction{}, false
	}
	var masalah []string
	u := func(medan, nilai string) *apd.Decimal {
		return services.UraiDesimalIsian(medan, strings.TrimSpace(nilai), &masalah)
	}
	// SetRIComm_Act langkah 5 / TotalBrokerageFee_Act langkah 3: `@replaceAll(RICommision, ",", ".")` - koma desimal % Premium
	// Deduction diterima (hanya medan itu, seperti Pega).
	isi := services.IsianDeduction{PctPremiumDeduction: u("pctPremiumDeduction", strings.ReplaceAll(b.PctPremiumDeduction, ",", ".")),
		PctBrokerageFee: u("pctBrokerageFee", b.PctBrokerageFee), MataUang: make([]models.BarisDeduction, 0, len(b.MataUang))}
	for i, m := range b.MataUang {
		p := fmt.Sprintf("mataUang[%d].", i)
		isi.MataUang = append(isi.MataUang, models.BarisDeduction{Currency: strings.TrimSpace(m.Currency),
			Commision: u(p+"commision", m.Commision), Deduction2: u(p+"deduction2", m.Deduction2), BrokerageFee: u(p+"brokerageFee", m.BrokerageFee)})
	}
	if len(masalah) > 0 {
		tulisGalat(w, fmt.Errorf("%w: %s", services.ErrMasukanDeduction, strings.Join(masalah, "; ")))
		return b, services.IsianDeduction{}, false
	}
	return b, isi, true
}

// bacaDeduction - GET /api/nbfacin/kasus/{caseId}/deduction.
func bacaDeduction(svc *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		d, err := svc.TampilanDeduction(r.Context(), r.PathValue("caseId"))
		if err != nil {
			tulisGalat(w, err)
			return
		}
		galat.TulisJSON(w, keDeductionKabel(d))
	}
}

// hitungDeduction - POST …/deduction/hitung {aksi, idx, …}: tanpa simpan.
func hitungDeduction(svc *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b, isi, ok := uraiIsianDeduction(w, r, true)
		if !ok {
			return
		}
		d, err := svc.HitungDeduction(r.Context(), r.PathValue("caseId"), *b.Aksi, *b.Idx, isi)
		if err != nil {
			tulisGalat(w, err)
			return
		}
		galat.TulisJSON(w, keDeductionKabel(d))
	}
}

// simpanDeduction - PUT …/deduction (badan tanpa aksi / idx): hitung ulang, simpan, jawab baca ulang.
func simpanDeduction(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, isi, ok := uraiIsianDeduction(w, r, false)
		if !ok {
			return
		}
		d, err := svc.SimpanDeduction(r.Context(), inti.PelakuDari(r, stubPelaku), r.PathValue("caseId"), isi)
		if err != nil {
			tulisGalat(w, err)
			return
		}
		galat.TulisJSON(w, keDeductionKabel(d))
	}
}
