package handlers

// Rute kontrak - tiket 14.
//
// ⛔ Dua rute, dan keduanya API. TIDAK ada layar baru: `L-4` menyatakan tidak
// ada spesifikasi layar di mana pun, dan papan melarang mengarangnya. Yang
// mendarat di sini lapisan APLIKASI tiket 14 - jalur simpan tempat INV-53 dan
// INV-29 akhirnya ditegakkan - bukan lapisan layarnya.

import (
	"encoding/json"
	"net/http"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/modul/treatyin/backend/services"
)

// batasBadanJSON - badan permintaan terbesar (1 MiB). Kepala kontrak tidak
// pernah sebesar itu; yang dibatasi badan yang tak berujung.
const batasBadanJSON = 1 << 20

func daftarkanKontrak(pasang func(string, rute)) {
	pasang("POST "+Prefix+"/kontrak", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		r.Body = http.MaxBytesReader(w, r.Body, batasBadanJSON)
		var m services.MasukanKontrak
		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			galat.Tulis(w, http.StatusBadRequest, "request body must be valid JSON")
			return
		}
		h, err := l.BuatKontrak(r.Context(), p, m)
		if jawabGalat(w, err) {
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(h)
	})
	pasang("GET "+Prefix+"/kontrak/{id}", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		id, ok := pengenal(w, r)
		if !ok {
			return
		}
		k, err := l.BacaKontrak(r.Context(), p, id)
		tulis(w, k, err)
	})
}
