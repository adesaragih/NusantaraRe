package handlers

// Menu View Policy (Facultative) - tiket 71:
//
//	GET /api/nbfacin/polis?opsi=&filter=&halaman=        {baris, total, halaman, ukuran} (baris = BarisFollowing tiket 70)
//	GET /api/nbfacin/polis/versi?nopolis=                 {versi: [{prodke, caseId, sistemBaru}]}
//	GET /api/nbfacin/polis/ringkasan?nopolis=&prodke=     {kepala, baris, ceding}

import (
	"net/http"
	"strconv"

	"github.com/cockroachdb/apd/v3"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/inti/backend/utils"
	"nusantarare/modul/nbfacin/backend/services"
)

func daftarkanRutePolis(mux *http.ServeMux, svc *services.Service, stubPelaku bool) {
	mux.HandleFunc("GET /api/nbfacin/polis", cariPolis(svc, stubPelaku))
	mux.HandleFunc("GET /api/nbfacin/polis/versi", versiPolis(svc, stubPelaku))
	mux.HandleFunc("GET /api/nbfacin/polis/ringkasan", ringkasanPolis(svc, stubPelaku))
}

type versiPolisKabel struct {
	ProdKe     string `json:"prodke"`
	CaseID     string `json:"caseId"`
	SistemBaru bool   `json:"sistemBaru"`
}

type kepalaPolisKabel struct {
	CaseID           string `json:"caseId"`
	PolicyNo         string `json:"policyNo"`
	NoEndors         string `json:"noEndors"`
	InsuredName      string `json:"insuredName"`
	QQName           string `json:"qqName"`
	SourceOfBusiness string `json:"sourceOfBusiness"`
	CedingCo         string `json:"cedingCo"`
	StartPeriod      string `json:"startPeriod"`
	EndPeriod        string `json:"endPeriod"`
	PolicyType       string `json:"policyType"`
	ReffNumber       string `json:"reffNumber"`
	GroupName        string `json:"groupName"`
	StatusBusiness   string `json:"statusBusiness"`
	NoMasterPolis    string `json:"noMasterPolis"`
}

type barisPolisKabel struct {
	ObjectNo         string `json:"objectNo"`
	RiskLocation     string `json:"riskLocation"`
	ObjectItem       string `json:"objectItem"`
	Coverage         string `json:"coverage"`
	AccumulationCode string `json:"accumulationCode"`
	Currency         string `json:"currency"`
	TSI100           string `json:"tsi100"`
	TSIRNM           string `json:"tsiRnm"`
	Rate             string `json:"rate"`
	PremiRNM         string `json:"premiRnm"`
	Treaty           string `json:"treaty"`
	PctShare         string `json:"pctShare"`
	RIComm           string `json:"riComm"`
}

type cedingPolisKabel struct {
	CedingName string `json:"cedingName"`
	PctShare   string `json:"pctShare"`
}

type ringkasanPolisKabel struct {
	Kepala kepalaPolisKabel   `json:"kepala"`
	Baris  []barisPolisKabel  `json:"baris"`
	Ceding []cedingPolisKabel `json:"ceding"`
}

// angkaPolis - desimal -> teks bertitik; nil -> "".
func angkaPolis(d *apd.Decimal) string {
	if d == nil {
		return ""
	}
	return utils.FormatDecimal(d)
}

// cariPolis - GET /api/nbfacin/polis; halaman kosong = 1.
func cariPolis(svc *services.Service, stubPelaku bool) http.HandlerFunc {
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
		h, err := svc.CariPolis(r.Context(), inti.PelakuDari(r, stubPelaku), q.Get("opsi"), q.Get("filter"), halaman)
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

// versiPolis - GET /api/nbfacin/polis/versi.
func versiPolis(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		v, err := svc.VersiPolis(r.Context(), inti.PelakuDari(r, stubPelaku), r.URL.Query().Get("nopolis"))
		if err != nil {
			tulisGalat(w, err)
			return
		}
		out := struct {
			Versi []versiPolisKabel `json:"versi"`
		}{Versi: []versiPolisKabel{}}
		for _, x := range v {
			out.Versi = append(out.Versi, versiPolisKabel{ProdKe: x.ProdKe, CaseID: x.CaseID, SistemBaru: x.SistemBaru})
		}
		galat.TulisJSON(w, out)
	}
}

// ringkasanPolis - GET /api/nbfacin/polis/ringkasan.
func ringkasanPolis(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		p, err := svc.RingkasanPolis(r.Context(), inti.PelakuDari(r, stubPelaku), q.Get("nopolis"), q.Get("prodke"))
		if err != nil {
			tulisGalat(w, err)
			return
		}
		k := p.Kepala
		out := ringkasanPolisKabel{Kepala: kepalaPolisKabel{CaseID: k.CaseID, PolicyNo: k.PolicyNo, NoEndors: k.NoEndors,
			InsuredName: k.InsuredName, QQName: k.QQName, SourceOfBusiness: k.SourceOfBusiness, CedingCo: k.CedingCo,
			StartPeriod: k.StartPeriod, EndPeriod: k.EndPeriod, PolicyType: k.PolicyType, ReffNumber: k.ReffNumber,
			GroupName: k.GroupName, StatusBusiness: k.StatusBusiness, NoMasterPolis: k.NoMasterPolis},
			Baris: []barisPolisKabel{}, Ceding: []cedingPolisKabel{}}
		for _, b := range p.Baris {
			out.Baris = append(out.Baris, barisPolisKabel{ObjectNo: b.ObjectNo, RiskLocation: b.RiskLocation, ObjectItem: b.ObjectItem,
				Coverage: b.Coverage, AccumulationCode: b.AccumulationCode, Currency: b.Currency, TSI100: angkaPolis(b.TSI100),
				TSIRNM: angkaPolis(b.TSIRNM), Rate: b.Rate, PremiRNM: angkaPolis(b.PremiRNM), Treaty: b.Treaty,
				PctShare: angkaPolis(b.PctShare), RIComm: angkaPolis(b.RIComm)})
		}
		for _, c := range p.Ceding {
			out.Ceding = append(out.Ceding, cedingPolisKabel{CedingName: c.CedingName, PctShare: angkaPolis(c.PctShare)})
		}
		galat.TulisJSON(w, out)
	}
}
