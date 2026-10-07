package handlers

// GET /api/nbfacin/sob?cari=&halaman= - pilihan SOB popup Change SOB (tiket 33).

import (
	"net/http"
	"strconv"

	"nusantarare/inti/backend/galat"
	"nusantarare/modul/nbfacin/backend/services"
)

// barisSOB - satu pilihan: id = kode yang dikirim balik lewat sourceOfBusinessId.
type barisSOB struct {
	ID       string `json:"id"`
	ClientID string `json:"clientId"`
	Name     string `json:"name"`
}

func cariSOB(svc *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		halaman := 1
		if h := r.URL.Query().Get("halaman"); h != "" {
			n, err := strconv.Atoi(h)
			if err != nil {
				n = 0 // tak sah -> 400 dari services
			}
			halaman = n
		}
		hasil, err := svc.CariSOB(r.Context(), r.URL.Query().Get("cari"), halaman)
		if err != nil {
			tulisGalat(w, err)
			return
		}
		baris := make([]barisSOB, 0, len(hasil.Baris))
		for _, b := range hasil.Baris {
			baris = append(baris, barisSOB{ID: b.ID, ClientID: b.ClientID, Name: b.Name})
		}
		galat.TulisJSON(w, struct {
			Baris   []barisSOB `json:"baris"`
			Total   int        `json:"total"`
			Halaman int        `json:"halaman"`
			Ukuran  int        `json:"ukuran"`
		}{baris, hasil.Total, hasil.Halaman, hasil.Ukuran})
	}
}
