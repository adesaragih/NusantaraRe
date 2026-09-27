package handlers

// Pintu HTTP dropdown Register - A3 kelompok Register.
//
// `GET /api/rujukan/{jenis}?cari=`
//
// Dibaca sesudah: inbox.go.

import (
	"encoding/json"
	"errors"
	"net/http"

	"nusantarare/internal/services"
)

// cariRujukan melayani GET /api/rujukan/{jenis}.
func cariRujukan(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat(w, http.StatusServiceUnavailable, "database belum dikonfigurasi")
			return
		}
		hasil, err := svc.SumberRujukan().Cari(r.Context(),
			pelakuDari(r, stubPelaku), r.PathValue("jenis"), r.URL.Query().Get("cari"))
		switch {
		case errors.Is(err, services.ErrTanpaIdentitas):
			galat(w, http.StatusUnauthorized, "permintaan tanpa identitas pelaku ditolak")
			return
		case errors.Is(err, services.ErrJenisRujukanTidakDikenal):
			galat(w, http.StatusNotFound, "jenis rujukan tidak dikenal")
			return
		case err != nil:
			galat(w, http.StatusInternalServerError, "gagal membaca daftar rujukan")
			return
		}
		// Go menulis slice kosong sebagai null; layar menginginkan daftar.
		if hasil == nil {
			hasil = []services.BarisRujukan{}
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(hasil)
	}
}
