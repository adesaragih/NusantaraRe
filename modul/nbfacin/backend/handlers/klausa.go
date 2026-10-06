package handlers

// Tab Clauses kasus FIRE (tiket 47): GET / PUT /api/nbfacin/kasus/{caseId}/klausa, GET /api/nbfacin/klausa/{id}/argumen,
// dan GET /api/nbfacin/klausa (pencarian). Kontrak frontend modul/nbfacin/frontend/api.ts.

import (
	"encoding/json"
	"net/http"
	"strconv"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/services"
)

// batasBadanKlausa - badan PUT ClauseList (isi klausa sampai 4000 bita per baris, 500 baris).
const batasBadanKlausa = 8 << 20

// argumenKlausaKabel - kontrak `{ argumentNumber, argumentDescription, argumentValue }`.
type argumenKlausaKabel struct {
	ArgumentNumber      string `json:"argumentNumber"`
	ArgumentDescription string `json:"argumentDescription"`
	ArgumentValue       string `json:"argumentValue"`
}

// klausaKasusKabel - kontrak `KlausaKasus`.
type klausaKasusKabel struct {
	ClauseCode        string               `json:"clauseCode"`
	ClauseTitle       string               `json:"clauseTitle"`
	ClauseDescription string               `json:"clauseDescription"`
	ClauseLanguage    string               `json:"clauseLanguage"`
	ClauseLanguageID  string               `json:"clauseLanguageId"`
	ClauseContent     string               `json:"clauseContent"`
	ClauseContentTemp string               `json:"clauseContentTemp"`
	ArgumentCount     string               `json:"argumentCount"`
	ArgumentList      []argumenKlausaKabel `json:"argumentList"`
}

func keArgumenKabel(d []models.ArgumenKlausa) []argumenKlausaKabel {
	out := make([]argumenKlausaKabel, 0, len(d))
	for _, a := range d {
		out = append(out, argumenKlausaKabel{a.Number, a.Description, a.Value})
	}
	return out
}

func keKlausaKabel(d []models.KlausaKasus) any {
	baris := make([]klausaKasusKabel, 0, len(d))
	for _, k := range d {
		baris = append(baris, klausaKasusKabel{k.Code, k.Title, k.Description, k.Language, k.LanguageID, k.Content,
			k.ContentTemp, k.ArgumentCount, keArgumenKabel(k.Arguments)})
	}
	return struct {
		Baris []klausaKasusKabel `json:"baris"`
	}{baris}
}

// bacaKlausa - GET /api/nbfacin/kasus/{caseId}/klausa.
func bacaKlausa(svc *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		d, err := svc.BacaKlausa(r.Context(), r.PathValue("caseId"))
		if err != nil {
			tulisGalat(w, err)
			return
		}
		galat.TulisJSON(w, keKlausaKabel(d))
	}
}

// simpanKlausa - PUT /api/nbfacin/kasus/{caseId}/klausa badan {baris: KlausaKasus[]} (kunci asing -> 400): ganti utuh,
// jawab baca ulang.
func simpanKlausa(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var b struct {
			Baris []klausaKasusKabel `json:"baris"`
		}
		dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, batasBadanKlausa))
		dec.DisallowUnknownFields()
		if err := dec.Decode(&b); err != nil || b.Baris == nil {
			galat.Tulis(w, http.StatusBadRequest, "badan permintaan harus {\"baris\": [klausa]}")
			return
		}
		baris := make([]models.KlausaKasus, 0, len(b.Baris))
		for _, k := range b.Baris {
			arg := make([]models.ArgumenKlausa, 0, len(k.ArgumentList))
			for _, a := range k.ArgumentList {
				arg = append(arg, models.ArgumenKlausa{Number: a.ArgumentNumber, Description: a.ArgumentDescription, Value: a.ArgumentValue})
			}
			baris = append(baris, models.KlausaKasus{Code: k.ClauseCode, Title: k.ClauseTitle, Description: k.ClauseDescription,
				Language: k.ClauseLanguage, LanguageID: k.ClauseLanguageID, Content: k.ClauseContent, ContentTemp: k.ClauseContentTemp,
				ArgumentCount: k.ArgumentCount, Arguments: arg})
		}
		d, err := svc.GantiKlausa(r.Context(), inti.PelakuDari(r, stubPelaku), r.PathValue("caseId"), baris)
		if err != nil {
			tulisGalat(w, err)
			return
		}
		galat.TulisJSON(w, keKlausaKabel(d))
	}
}

// argumenKlausa - GET /api/nbfacin/klausa/{id}/argumen -> {baris}.
func argumenKlausa(svc *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		d, err := svc.ArgumenKlausa(r.Context(), r.PathValue("id"))
		if err != nil {
			tulisGalat(w, err)
			return
		}
		galat.TulisJSON(w, struct {
			Baris []argumenKlausaKabel `json:"baris"`
		}{keArgumenKabel(d)})
	}
}

// hasilKlausaKabel - kontrak `HasilKlausa` `{ id, title, info, text, language, argumentCount }`.
type hasilKlausaKabel struct {
	ID            string `json:"id"`
	Title         string `json:"title"`
	Info          string `json:"info"`
	Text          string `json:"text"`
	Language      string `json:"language"`
	ArgumentCount string `json:"argumentCount"`
}

// cariKlausa - GET /api/nbfacin/klausa?bahasa=&q=&halaman= -> {baris, total, halaman, ukuran: 10}.
func cariKlausa(svc *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		nomor := 1
		if v := q.Get("halaman"); v != "" {
			n, err := strconv.Atoi(v)
			if err != nil {
				galat.Tulis(w, http.StatusBadRequest, "halaman harus bilangan bulat")
				return
			}
			nomor = n
		}
		h, err := svc.CariKlausa(r.Context(), q.Get("bahasa"), q.Get("q"), nomor)
		if err != nil {
			tulisGalat(w, err)
			return
		}
		baris := make([]hasilKlausaKabel, 0, len(h.Baris))
		for _, k := range h.Baris {
			baris = append(baris, hasilKlausaKabel(k))
		}
		galat.TulisJSON(w, struct {
			Baris   []hasilKlausaKabel `json:"baris"`
			Total   int                `json:"total"`
			Halaman int                `json:"halaman"`
			Ukuran  int                `json:"ukuran"`
		}{baris, h.Total, h.Nomor, h.Ukuran})
	}
}
