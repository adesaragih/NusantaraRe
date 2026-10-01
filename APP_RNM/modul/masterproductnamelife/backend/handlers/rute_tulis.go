package handlers

// Rute tulis modul Master Product Name Life.
//
//	POST /api/master-product-name-life/produk        produk BARU - `Add` b71816 → … → `Save` (paket 3)
//	PUT  /api/master-product-name-life/produk/{id}   ubah - `View` → `Edit` b59443 → `Save` (paket 3)
//
// ⛔ POST dan PUT terpisah walau Pega punya satu `Save` ber-upsert: identitas
// baru tidak pernah datang dari klien (ADR-0006), dan badan PUT yang membawa
// `id` berbeda dari jalurnya DITOLAK, bukan salah satu dipilih diam-diam.

import (
	"encoding/json"
	"net/http"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/modul/masterproductnamelife/backend/models"
	"nusantarare/modul/masterproductnamelife/backend/services"
)

// bacaBadan mengurai badan JSON; gagal = 400 dan true.
func bacaBadan(w http.ResponseWriter, r *http.Request, ke any) bool {
	if err := json.NewDecoder(r.Body).Decode(ke); err != nil {
		galat.Tulis(w, http.StatusBadRequest, "request body is not valid JSON")
		return true
	}
	return false
}

// idJalur menyamakan id badan dengan id jalur (PUT); beda = 400 dan true.
func idJalur(w http.ResponseWriter, r *http.Request, id *string) bool {
	dariJalur := r.PathValue("id")
	if *id != "" && *id != dariJalur {
		galat.Tulis(w, http.StatusBadRequest, "id in the body differs from id in the path")
		return true
	}
	*id = dariJalur
	return false
}

func daftarkanTulis(pasang func(string, rute)) {
	pasang("POST "+Prefix+"/produk", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		var m models.Produk
		if bacaBadan(w, r, &m) {
			return
		}
		hasil, err := l.SimpanProduk(r.Context(), p, m, true)
		tulis(w, hasil, err)
	})
	pasang("PUT "+Prefix+"/produk/{id}", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		var m models.Produk
		if bacaBadan(w, r, &m) || idJalur(w, r, &m.ID) {
			return
		}
		hasil, err := l.SimpanProduk(r.Context(), p, m, false)
		tulis(w, hasil, err)
	})
}
