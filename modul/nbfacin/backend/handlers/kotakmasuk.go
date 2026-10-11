package handlers

// Kotak masuk Beranda NB Fac In (permintaan work owner 07-10-2026): GET /api/nbfacin/kotak-masuk -> [{workbasket, nama,
// jumlah}], GET /api/nbfacin/kotak-masuk/kasus?workbasket= -> [baris portal]. Kontrak frontend api.ts
// `AntreanKotakMasukNbFacIn`, `BarisCaseNB`.

import (
	"net/http"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/modul/nbfacin/backend/services"
)

type antreanKotakMasuk struct {
	Workbasket string `json:"workbasket"`
	Nama       string `json:"nama"`
	Jumlah     int    `json:"jumlah"`
}

// kotakMasuk - GET /api/nbfacin/kotak-masuk.
func kotakMasuk(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		d, err := svc.KotakMasuk(r.Context(), inti.PelakuDari(r, stubPelaku))
		if err != nil {
			tulisGalat(w, err)
			return
		}
		out := make([]antreanKotakMasuk, 0, len(d))
		for _, a := range d {
			out = append(out, antreanKotakMasuk(a))
		}
		galat.TulisJSON(w, out)
	}
}

// barisInbox - api.ts `BarisKotakMasukNbFacIn`.
type barisInbox struct {
	CaseID           string `json:"caseId"`
	InsuredName      string `json:"insuredName"`
	BusinessName     string `json:"businessName"`
	SourceOfBusiness string `json:"sourceOfBusiness"`
	CedingCompany    string `json:"cedingCompany"`
	InceptionDate    string `json:"inceptionDate"`
	ShareRNM         string `json:"shareRnm"`
	Marketing        string `json:"marketing"`
	Status           string `json:"status"`
	CreatedDate      string `json:"createdDate"`
	OfferingDate     string `json:"offeringDate"`
	ValidityDate     string `json:"validityDate"`
	StatusInbox      string `json:"statusInbox"`
	AksiTerakhir     string `json:"aksiTerakhir"`
	TglUpdate        string `json:"tglUpdate"`
	SubmitToUW       string `json:"submitToUw"`
}

// daftarMenunggu - GET /api/nbfacin/kotak-masuk/kasus?workbasket=.
func daftarMenunggu(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		d, err := svc.DaftarMenunggu(r.Context(), inti.PelakuDari(r, stubPelaku), r.URL.Query().Get("workbasket"))
		if err != nil {
			tulisGalat(w, err)
			return
		}
		out := make([]barisInbox, 0, len(d))
		for _, b := range d {
			out = append(out, barisInbox(b))
		}
		galat.TulisJSON(w, out)
	}
}
