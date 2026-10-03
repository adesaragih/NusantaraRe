package handlers

// Copy Old - permintaan work owner 03-10-2026 (tombol di samping `Add`):
//
//	GET  /api/master-product-name-life/produk-lama         isi popup: produk tabel JSON warisan yang belum ada di tabel flat
//	POST /api/master-product-name-life/produk-lama/salin   `Process Copy` - badan {"ids": [...]}, jawaban hasil per ID

import (
	"net/http"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/masterproductnamelife/backend/services"
)

func daftarkanLama(pasang func(string, rute)) {
	pasang("GET "+Prefix+"/produk-lama", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		d, err := l.DaftarProdukLama(r.Context(), p)
		tulisDaftar(w, d, err)
	})
	pasang("POST "+Prefix+"/produk-lama/salin", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		var b struct {
			IDs []string `json:"ids"`
		}
		if bacaBadan(w, r, &b) {
			return
		}
		j, err := l.SalinProdukLama(r.Context(), p, b.IDs)
		tulis(w, j, err)
	})
}
