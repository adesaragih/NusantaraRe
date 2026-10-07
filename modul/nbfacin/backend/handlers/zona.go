package handlers

// Zone dan 4.2 Construction form coverage FIRE (tiket 52): GET /api/nbfacin/zona dan GET /api/nbfacin/konstruksi-gempa.
// Kontrak frontend modul/nbfacin/frontend/api.ts (saranZona, saranKonstruksiGempa) - bentuk {baris:[{id,label}]}.

import (
	"net/http"

	"nusantarare/inti/backend/galat"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/services"
)

// tulisSaran - {baris:[{id,label}]} (ekstra kosong tidak dikirim).
func tulisSaran(w http.ResponseWriter, d []models.SaranAkumulasi, err error) {
	if err != nil {
		tulisGalat(w, err)
		return
	}
	baris := make([]saranAkumulasiKabel, 0, len(d))
	for _, b := range d {
		baris = append(baris, saranAkumulasiKabel(b))
	}
	galat.TulisJSON(w, struct {
		Baris []saranAkumulasiKabel `json:"baris"`
	}{baris})
}

// saranZona - GET /api/nbfacin/zona?coverage=&oldId=&q=.
func saranZona(svc *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		d, err := svc.SaranZona(r.Context(), q.Get("coverage"), q.Get("oldId"), q.Get("q"))
		tulisSaran(w, d, err)
	}
}

// saranKonstruksiGempa - GET /api/nbfacin/konstruksi-gempa?q=.
func saranKonstruksiGempa(svc *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		d, err := svc.SaranKonstruksiGempa(r.Context(), r.URL.Query().Get("q"))
		tulisSaran(w, d, err)
	}
}
