package handlers

// View Policy tab Offer Facultative / Treaty / Life - tiket 71 perluasan (permintaan work owner 11-10-2026; kontrak sesi 9d). Semua
// {baris, total, halaman, ukuran}, 20 baris per halaman, filter kosong -> baris [] total 0:
//
//	GET /api/nbfacin/polis/offer?opsi=&filter=&halaman=       baris {caseId, noOffer, insuredName, sourceOfBusiness, qqName "", startPeriod, endPeriod}
//	GET /api/nbfacin/polis/treaty?opsi=&filter=&halaman=      baris {caseId, policyNo, insuredName, sourceOfBusiness, startPeriod, endPeriod}
//	GET /api/nbfacin/polis/life?filter=&halaman=              baris {caseId, policyNo, endorsNo, typePolicy, sourceOfBusiness, cedingCoName, qq, riSlipCeding}
//	GET /api/nbfacin/polis/life-offer?filter=&halaman=        baris {caseId, noOffer, typePolicy, businessName, cedingCoName, qq, lastStatus}
//	GET /api/nbfacin/polis/treaty/berkas?nopolis=            {modul: "nbtreatyin" | "edmtreatyin", kasus} (tombol View tab Treaty)

import (
	"net/http"
	"strconv"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/services"
)

func daftarkanRuteViewPolis(mux *http.ServeMux, svc *services.Service, stubPelaku bool) {
	mux.HandleFunc("GET /api/nbfacin/polis/offer", cariView(stubPelaku, func(r *http.Request, p inti.Pelaku, h int) (services.HalamanView[models.BarisOfferFac], error) {
		return svc.CariOfferFac(r.Context(), p, r.URL.Query().Get("opsi"), r.URL.Query().Get("filter"), h)
	}, kabelOfferFac))
	mux.HandleFunc("GET /api/nbfacin/polis/treaty", cariView(stubPelaku, func(r *http.Request, p inti.Pelaku, h int) (services.HalamanView[models.BarisPolisTreaty], error) {
		return svc.CariPolisTreaty(r.Context(), p, r.URL.Query().Get("opsi"), r.URL.Query().Get("filter"), h)
	}, kabelPolisTreaty))
	mux.HandleFunc("GET /api/nbfacin/polis/life", cariView(stubPelaku, func(r *http.Request, p inti.Pelaku, h int) (services.HalamanView[models.BarisPolisLife], error) {
		return svc.CariPolisLife(r.Context(), p, r.URL.Query().Get("filter"), h)
	}, kabelPolisLife))
	mux.HandleFunc("GET /api/nbfacin/polis/life-offer", cariView(stubPelaku, func(r *http.Request, p inti.Pelaku, h int) (services.HalamanView[models.BarisOfferLife], error) {
		return svc.CariOfferLife(r.Context(), p, r.URL.Query().Get("filter"), h)
	}, kabelOfferLife))
	mux.HandleFunc("GET /api/nbfacin/polis/treaty/berkas", berkasPolisTreaty(svc, stubPelaku))
}

type barisOfferFacKabel struct {
	CaseID           string `json:"caseId"`
	NoOffer          string `json:"noOffer"`
	InsuredName      string `json:"insuredName"`
	SourceOfBusiness string `json:"sourceOfBusiness"`
	QQName           string `json:"qqName"`
	StartPeriod      string `json:"startPeriod"`
	EndPeriod        string `json:"endPeriod"`
}

type barisPolisTreatyKabel struct {
	CaseID           string `json:"caseId"`
	PolicyNo         string `json:"policyNo"`
	InsuredName      string `json:"insuredName"`
	SourceOfBusiness string `json:"sourceOfBusiness"`
	StartPeriod      string `json:"startPeriod"`
	EndPeriod        string `json:"endPeriod"`
}

type barisPolisLifeKabel struct {
	CaseID           string `json:"caseId"`
	PolicyNo         string `json:"policyNo"`
	EndorsNo         string `json:"endorsNo"`
	TypePolicy       string `json:"typePolicy"`
	SourceOfBusiness string `json:"sourceOfBusiness"`
	CedingCoName     string `json:"cedingCoName"`
	QQ               string `json:"qq"`
	RISlipCeding     string `json:"riSlipCeding"`
}

type barisOfferLifeKabel struct {
	CaseID       string `json:"caseId"`
	NoOffer      string `json:"noOffer"`
	TypePolicy   string `json:"typePolicy"`
	BusinessName string `json:"businessName"`
	CedingCoName string `json:"cedingCoName"`
	QQ           string `json:"qq"`
	LastStatus   string `json:"lastStatus"`
}

type halamanViewKabel[K any] struct {
	Baris   []K `json:"baris"`
	Total   int `json:"total"`
	Halaman int `json:"halaman"`
	Ukuran  int `json:"ukuran"`
}

func kabelOfferFac(b models.BarisOfferFac) barisOfferFacKabel {
	return barisOfferFacKabel{CaseID: b.CaseID, NoOffer: b.NoOffer, InsuredName: b.InsuredName, SourceOfBusiness: b.SourceOfBusiness,
		StartPeriod: b.StartPeriod, EndPeriod: b.EndPeriod}
}

func kabelPolisTreaty(b models.BarisPolisTreaty) barisPolisTreatyKabel {
	return barisPolisTreatyKabel{CaseID: b.CaseID, PolicyNo: b.PolicyNo, InsuredName: b.InsuredName, SourceOfBusiness: b.SourceOfBusiness,
		StartPeriod: b.StartPeriod, EndPeriod: b.EndPeriod}
}

func kabelPolisLife(b models.BarisPolisLife) barisPolisLifeKabel {
	return barisPolisLifeKabel{CaseID: b.CaseID, PolicyNo: b.PolicyNo, EndorsNo: b.EndorsNo, TypePolicy: b.TypePolicy,
		SourceOfBusiness: b.SourceOfBusiness, CedingCoName: b.CedingCoName, QQ: b.QQ, RISlipCeding: b.RISlipCeding}
}

func kabelOfferLife(b models.BarisOfferLife) barisOfferLifeKabel {
	return barisOfferLifeKabel{CaseID: b.CaseID, NoOffer: b.NoOffer, TypePolicy: b.TypePolicy, BusinessName: b.BusinessName,
		CedingCoName: b.CedingCoName, QQ: b.QQ, LastStatus: b.LastStatus}
}

// cariView - halaman kosong = 1, bukan bilangan bulat 400 (seperti cariPolis); baris dipetakan ke kabel.
func cariView[T, K any](stubPelaku bool, cari func(*http.Request, inti.Pelaku, int) (services.HalamanView[T], error),
	kabel func(T) K) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		halaman := 1
		if v := r.URL.Query().Get("halaman"); v != "" {
			n, err := strconv.Atoi(v)
			if err != nil {
				galat.Tulis(w, http.StatusBadRequest, "halaman wajib bilangan bulat")
				return
			}
			halaman = n
		}
		h, err := cari(r, inti.PelakuDari(r, stubPelaku), halaman)
		if err != nil {
			tulisGalat(w, err)
			return
		}
		out := halamanViewKabel[K]{Baris: []K{}, Total: h.Total, Halaman: h.Halaman, Ukuran: h.Ukuran}
		for _, b := range h.Baris {
			out.Baris = append(out.Baris, kabel(b))
		}
		galat.TulisJSON(w, out)
	}
}

type berkasPolisTreatyKabel struct {
	Modul string `json:"modul"`
	Kasus string `json:"kasus"`
}

// berkasPolisTreaty - GET /api/nbfacin/polis/treaty/berkas?nopolis=.
func berkasPolisTreaty(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := svc.BerkasPolisTreaty(r.Context(), inti.PelakuDari(r, stubPelaku), r.URL.Query().Get("nopolis"))
		if err != nil {
			tulisGalat(w, err)
			return
		}
		galat.TulisJSON(w, berkasPolisTreatyKabel{Modul: b.Modul, Kasus: b.Kasus})
	}
}
