package handlers

// Following Policy Number (tiket 70) - popup Search tab General:
//
//	GET  /api/nbfacin/kasus/{caseId}/following?opsi=&filter=&halaman=   {baris, total, halaman, ukuran}
//	POST /api/nbfacin/kasus/{caseId}/following/nopolis                   {caseIdSumber, policyNo} -> KasusNB
//	POST /api/nbfacin/kasus/{caseId}/following/salin-semua               {caseIdSumber, policyNo} -> KasusNB

import (
	"encoding/json"
	"net/http"
	"strconv"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/modul/nbfacin/backend/services"
)

// batasBadanFollowing - dua medan pendek.
const batasBadanFollowing = 4 << 10

func daftarkanRuteFollowing(mux *http.ServeMux, svc *services.Service, stubPelaku bool) {
	mux.HandleFunc("GET /api/nbfacin/kasus/{caseId}/following", cariFollowing(svc, stubPelaku))
	mux.HandleFunc("POST /api/nbfacin/kasus/{caseId}/following/nopolis", salinFollowing(svc, stubPelaku, false))
	mux.HandleFunc("POST /api/nbfacin/kasus/{caseId}/following/salin-semua", salinFollowing(svc, stubPelaku, true))
}

// barisFollowingKabel - satu baris grid (api.ts BarisFollowing).
type barisFollowingKabel struct {
	CaseID           string `json:"caseId"`
	PolicyNo         string `json:"policyNo"`
	ReffNumber       string `json:"reffNumber"`
	InsuredName      string `json:"insuredName"`
	SourceOfBusiness string `json:"sourceOfBusiness"`
	QQName           string `json:"qqName"`
	StartPeriod      string `json:"startPeriod"`
	EndPeriod        string `json:"endPeriod"`
}

type halamanFollowingKabel struct {
	Baris   []barisFollowingKabel `json:"baris"`
	Total   int                   `json:"total"`
	Halaman int                   `json:"halaman"`
	Ukuran  int                   `json:"ukuran"`
}

type isianFollowing struct {
	CaseIDSumber string `json:"caseIdSumber"`
	PolicyNo     string `json:"policyNo"`
}

// cariFollowing - GET …/following; halaman kosong = 1.
func cariFollowing(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		halaman := 1
		if v := q.Get("halaman"); v != "" {
			n, err := strconv.Atoi(v)
			if err != nil {
				galat.Tulis(w, http.StatusBadRequest, "halaman wajib bilangan bulat")
				return
			}
			halaman = n
		}
		h, err := svc.CariFollowing(r.Context(), inti.PelakuDari(r, stubPelaku), r.PathValue("caseId"), q.Get("opsi"), q.Get("filter"), halaman)
		if err != nil {
			tulisGalat(w, err)
			return
		}
		out := halamanFollowingKabel{Baris: []barisFollowingKabel{}, Total: h.Total, Halaman: h.Halaman, Ukuran: h.Ukuran}
		for _, b := range h.Baris {
			out.Baris = append(out.Baris, barisFollowingKabel{CaseID: b.CaseID, PolicyNo: b.PolicyNo, ReffNumber: b.ReffNumber,
				InsuredName: b.InsuredName, SourceOfBusiness: b.SourceOfBusiness, QQName: b.QQName, StartPeriod: b.StartPeriod,
				EndPeriod: b.EndPeriod})
		}
		galat.TulisJSON(w, out)
	}
}

// salinFollowing - POST …/following/nopolis (semua=false) dan …/salin-semua (semua=true); jawaban = GET kasus.
func salinFollowing(svc *services.Service, stubPelaku, semua bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var b isianFollowing
		if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, batasBadanFollowing)).Decode(&b); err != nil {
			galat.Tulis(w, http.StatusBadRequest, "badan permintaan bukan JSON {caseIdSumber, policyNo} yang sah")
			return
		}
		pelaku, id := inti.PelakuDari(r, stubPelaku), r.PathValue("caseId")
		salin := svc.SalinNopolisFollowing
		if semua {
			salin = svc.SalinSemuaFollowing
		}
		if err := salin(r.Context(), pelaku, id, b.CaseIDSumber, b.PolicyNo); err != nil {
			tulisGalat(w, err)
			return
		}
		k, err := svc.BacaKasus(r.Context(), id)
		if err != nil {
			tulisGalat(w, err)
			return
		}
		galat.TulisJSON(w, denganAksesTab(keKasusKabel(k), svc, r, stubPelaku))
	}
}
