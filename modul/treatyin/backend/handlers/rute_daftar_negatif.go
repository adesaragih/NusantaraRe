package handlers

import (
	"net/http"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/treatyin/backend/services"
)

// daftarkanDaftarNegatifAgen - `TreatyInCheckCedingBlacklist` ("This name is
// on Agent Negative List"), dijalankan layar saat kontrak dibuka lewat
// tombol `Edit`. BACA SAJA atas `AGENT`.
//
//	GET /agen/daftar-negatif?cedant=<CedingID>&asalBisnis=<LeadingReinsSourceID>
func daftarkanDaftarNegatifAgen(pasang func(string, rute)) {
	pasang("GET "+Prefix+"/agen/daftar-negatif", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		q := r.URL.Query()
		hasil, err := l.PeriksaDaftarNegatifAgen(r.Context(), p, q.Get("cedant"), q.Get("asalBisnis"))
		tulis(w, hasil, err)
	})
}
