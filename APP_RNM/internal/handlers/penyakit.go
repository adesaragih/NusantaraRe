package handlers

// Pintu HTTP pencarian diagnosa - kelompok Medis.
//
//	`GET /api/penyakit-life?icd=&nama=&batas=`
//
// Meniru `ReportDefinition/BrowseDiseaseLife_RD.xml` lewat
// `services.Diagnosa`. GET: ia membaca, tidak mengubah apa pun.
//
// ⛔ `batas` yang diminta klien TIDAK dipercaya apa adanya - ia dijepit
// `models.BatasPenyakit` ke `pyMaxRecords` 500. Tabelnya 97.586 baris, dan
// `?batas=1000000` adalah permintaan yang memuat seluruh tabel ke memori.
//
// Dibaca sesudah: tutup.go.

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"nusantarare/internal/services"
)

// cariPenyakit melayani GET /api/penyakit-life.
func cariPenyakit(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat(w, http.StatusServiceUnavailable, "database belum dikonfigurasi")
			return
		}
		// Batas yang tidak terbaca menjadi 0, dan 0 dijepit ke ukuran
		// halaman - bukan menjadi galat. Pencarian yang ditolak karena satu
		// parameter salah ketik lebih menjengkelkan daripada berguna.
		batas, _ := strconv.Atoi(r.URL.Query().Get("batas"))

		hasil, err := svc.Penyakit().Cari(r.Context(), pelakuDari(r, stubPelaku),
			r.URL.Query().Get("icd"), r.URL.Query().Get("nama"), batas)
		switch {
		case err == nil:
		case errors.Is(err, services.ErrTanpaIdentitas):
			galat(w, http.StatusUnauthorized, "permintaan tanpa identitas pelaku ditolak")
			return
		case errors.Is(err, services.ErrTanpaWewenang):
			galat(w, http.StatusForbidden, "wewenang tidak mencukupi")
			return
		default:
			galat(w, http.StatusInternalServerError, "gagal mencari diagnosa")
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(hasil)
	}
}
