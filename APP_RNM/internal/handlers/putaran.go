package handlers

// Pintu HTTP `Add` grid adjustment - putaran berikutnya (tiket 11) dan baris
// PERTAMA (GILIRAN-13 butir bo). Satu rute: yang membedakan keduanya keadaan
// grid, dan itu diputuskan services.
//
// Nol aturan dagang di sini. Yang memutuskan boleh atau tidak adalah services.
//
// Dibaca sesudah: handlers.go dan services/hasilkomite.go.

import (
	"errors"
	"net/http"
	"time"

	"nusantarare/internal/services"
)

// tambahPutaran melayani
// POST /api/klaim-life/{id}/peserta/{pesertaId}/putaran.
func tambahPutaran(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat(w, http.StatusServiceUnavailable, "database belum dikonfigurasi")
			return
		}
		err := svc.Putaran().DenganJejak(services.PerekamJejakOracle(svc)).Tambah(r.Context(), pelakuDari(r, stubPelaku),
			r.PathValue("id"), r.PathValue("pesertaId"), time.Now())
		if jawabGalatPutaran(w, err) {
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

// jawabGalatPutaran menerjemahkan galat `Add` ke kode HTTP; false = tidak
// ada galat, pemanggil meneruskan ke jalur berhasil.
func jawabGalatPutaran(w http.ResponseWriter, err error) bool {
	switch {
	case err == nil:
		return false
	case errors.Is(err, services.ErrTanpaIdentitas):
		galat(w, http.StatusUnauthorized, "permintaan tanpa identitas pelaku ditolak")
	case errors.Is(err, services.ErrTanpaWewenang):
		galat(w, http.StatusForbidden,
			"hanya ReasLifeSPV yang dapat menambah baris adjustment")
	case errors.Is(err, services.ErrBukanPenolakan):
		// 409: bukan permintaan yang salah, melainkan bentrokan dengan
		// keadaan yang ada - baris terakhirnya belum ditolak.
		galat(w, http.StatusConflict,
			"putaran berikutnya hanya lahir sesudah baris terakhir ditolak")
	case errors.Is(err, services.ErrTahapTanpaAddAdjustment):
		// Padanan syarat tampil b18160 - di Pega tombolnya tidak ada, di sini
		// permintaannya ditolak dengan kalimat yang menyebut sebabnya.
		galat(w, http.StatusConflict,
			"baris adjustment pertama hanya dapat ditambahkan di tahap Claim Analis")
	case errors.Is(err, services.ErrPesertaSudahDiputus):
		galat(w, http.StatusConflict,
			"peserta sudah diputus; isian layar Detail-nya tidak dapat diubah lagi")
	case errors.Is(err, services.ErrTahapTidakDikenal):
		galat(w, http.StatusConflict, "tahap kasus ini tidak dikenal")
	case errors.Is(err, services.ErrKasusSudahTertutup):
		galat(w, http.StatusConflict, "kasus sudah ditutup dan tidak dapat diubah")
	case errors.Is(err, services.ErrPermintaanTidakSah):
		galat(w, http.StatusBadRequest, err.Error())
	default:
		galat(w, http.StatusInternalServerError, "gagal menambah baris adjustment")
	}
	return true
}
