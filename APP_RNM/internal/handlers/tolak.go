package handlers

// Pintu HTTP Reject Outstanding - tiket 05.
//
// Nol aturan dagang di sini: berkas ini menerjemahkan HTTP ke panggilan
// services dan sebaliknya. Yang memutuskan boleh atau tidak adalah services.
//
// ⭐ Pintu ini sekaligus menutup dua AC sisa tiket 04: pencerminan status ke
// peserta dan ke header baru benar-benar terjadi ketika ada yang memanggilnya.
//
// Dibaca sesudah: handlers.go dan services/tolak.go.

import (
	"errors"
	"net/http"
	"time"

	"nusantarare/internal/services"
)

// tolakBaris melayani
// POST /api/klaim-life/{id}/adjustment/{adjId}/tolak.
func tolakBaris(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat(w, http.StatusServiceUnavailable, "database belum dikonfigurasi")
			return
		}
		err := svc.Status().DenganJejak(services.PerekamJejakOracle(svc)).Tolak(r.Context(), pelakuDari(r, stubPelaku),
			r.PathValue("id"), r.PathValue("adjId"), time.Now())

		switch {
		case errors.Is(err, services.ErrTanpaIdentitas):
			// 401: yang kurang identitasnya. Dipisah dari 403 karena keduanya
			// menjawab pertanyaan yang berbeda - dan ronde pertama berkas ini
			// menjawab "bukan ReasLifeAdmin" kepada pemanggil yang sebenarnya
			// hanya belum menyebut dirinya.
			galat(w, http.StatusUnauthorized, "permintaan tanpa identitas pelaku ditolak")
			return
		case errors.Is(err, services.ErrTanpaWewenang):
			// 403: identitasnya ada, perannya yang kurang.
			galat(w, http.StatusForbidden,
				"hanya ReasLifeAdmin yang dapat menolak baris Outstanding")
			return
		case errors.Is(err, services.ErrBarisSudahFinal):
			// 409: baris itu sudah punya keputusan. Bukan galat permintaan,
			// melainkan bentrokan dengan keadaan yang sudah ada.
			galat(w, http.StatusConflict,
				"baris sudah diputus dan tidak dapat ditolak lagi")
			return
		case errors.Is(err, services.ErrKlaimBelumBernomor):
			// 422: datanya yang belum lengkap, bukan permintaannya yang salah.
			//
			// ⚠️ Sampai butir o diputuskan, TIDAK ADA klaim yang bernomor,
			// sehingga jalur nyata selalu berhenti di sini. Itu keadaan yang
			// benar - bukan kerusakan - dan dinyatakan di tiket.
			galat(w, http.StatusUnprocessableEntity,
				"klaim belum bernomor; penolakan baris menunggu nomor klaim")
			return
		case errors.Is(err, services.ErrPermintaanTidakSah):
			galat(w, http.StatusBadRequest, err.Error())
			return
		case err != nil:
			galat(w, http.StatusInternalServerError, "gagal menolak baris")
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
