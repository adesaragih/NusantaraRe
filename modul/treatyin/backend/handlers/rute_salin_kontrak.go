package handlers

import (
	"encoding/json"
	"net/http"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/modul/treatyin/backend/services"
)

// daftarkanSalin - tombol `Copy` daftar kontrak (`Section/InputTreatyInOffer.xml`
// cell 993 → `SetTreatyIn_Act(ID=.ID)` → `Activity/TreatyInCopy.xml`).
//
//	GET  /kontrak-warisan/{id}/salin  draf salinan untuk form — ⛔ NOL tulisan:
//	                          `TreatyInCopy` hanya `Property-Set` dan
//	                          `RDB-List GetCurrentDate`, tanpa `SaveTreatyIn`
//	POST /kontrak/salin       Save (atau Submit/Decline) draf itu — ID baru
//	                          lahir di sini, persis `UnknownId` Pega
//
// Lihat `services/salin_kontrak.go` untuk rantai ekspornya.
func daftarkanSalin(pasang func(string, rute)) {
	pasang("GET "+Prefix+"/kontrak-warisan/{id}/salin", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		k, err := l.BacaDraftSalinan(r.Context(), p, r.PathValue("id"))
		tulis(w, k, err)
	})
	pasang("POST "+Prefix+"/kontrak/salin", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		r.Body = http.MaxBytesReader(w, r.Body, batasBadanJSON)
		var m services.MasukanSalin
		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			galat.Tulis(w, http.StatusBadRequest, "request body must be valid JSON")
			return
		}
		hasil, err := l.SimpanSalinan(r.Context(), p, m)
		tulis(w, hasil, err)
	})
}
