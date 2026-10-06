package handlers

// GET /api/nbfacin/opportunity?cari=&halaman= - daftar case NB di portal Opportunity (tiket 32).

import (
	"net/http"
	"strconv"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/modul/nbfacin/backend/services"
)

// barisPortal - satu baris daftar; teks apa adanya, NULL = "".
type barisPortal struct {
	CaseID        string `json:"caseId"`
	Name          string `json:"name"`
	GroupBusiness string `json:"groupBusiness"`
	InsuredName   string `json:"insuredName"`
	Marketing     string `json:"marketing"`
	Status        string `json:"status"`
}

func cariPortal(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// halaman bukan angka diteruskan sebagai 0 (tak sah): services memeriksa identitas
		// lebih dulu, jadi permintaan tanpa identitas tetap 401.
		halaman := 1
		if h := r.URL.Query().Get("halaman"); h != "" {
			n, err := strconv.Atoi(h)
			if err != nil {
				n = 0
			}
			halaman = n
		}
		hasil, err := svc.CariPortal(r.Context(), inti.PelakuDari(r, stubPelaku), r.URL.Query().Get("cari"), halaman)
		if err != nil {
			tulisGalat(w, err)
			return
		}
		baris := make([]barisPortal, 0, len(hasil.Baris))
		for _, b := range hasil.Baris {
			baris = append(baris, barisPortal{CaseID: b.CaseID, Name: b.Name, GroupBusiness: b.GroupBusiness,
				InsuredName: b.InsuredName, Marketing: b.Marketing, Status: b.Status})
		}
		galat.TulisJSON(w, struct {
			Baris   []barisPortal `json:"baris"`
			Total   int           `json:"total"`
			Halaman int           `json:"halaman"`
			Ukuran  int           `json:"ukuran"`
		}{baris, hasil.Total, hasil.Halaman, hasil.Ukuran})
	}
}
