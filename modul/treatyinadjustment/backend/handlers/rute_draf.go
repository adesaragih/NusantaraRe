package handlers

import (
	"encoding/json"
	"net/http"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/modul/treatyinadjustment/backend/models"
	"nusantarare/modul/treatyinadjustment/backend/services"
)

// daftarkanDraf - tombol `Add Revision` / `Add Adjustment Premium` dan
// `Choose` di picker-nya (`Section/InputTreatyInAdjustment.xml` @500554,
// @515429).
//
// ⛔ KEDUANYA BACA SAJA. `POST /draf` menyusun penyesuaian baru dari dokumen
// master dan MENGEMBALIKANNYA — nol tulisan. Di Pega `Choose` juga menyimpan;
// di sini simpanannya menunggu jalur Save (pemilik proses, 7 Oktober 2026).
// `POST` karena masukannya parameter Activity, bukan pengenal sumber daya.
func daftarkanDraf(pasang func(string, rute)) {
	pasang("GET "+Prefix+"/penyesuaian-warisan/master", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		d, err := l.DaftarMasterPilihan(r.Context(), p, r.URL.Query().Get("jenis"))
		tulis(w, d, err)
	})
	pasang("POST "+Prefix+"/penyesuaian-warisan/draf", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		var m models.MasukanDraf
		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			galat.Tulis(w, http.StatusBadRequest, "request body must be valid JSON")
			return
		}
		d, err := l.DrafPenyesuaian(r.Context(), p, m)
		tulis(w, d, err)
	})
}
