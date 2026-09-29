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
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"nusantarare/internal/services"
	"nusantarare/inti"
	"nusantarare/inti/galat"
	"nusantarare/inti/jejak"
)

// permintaanTolakJSON - isian dialog Reject Outstanding (OQ-M5, GILIRAN-17).
// Hanya `Remarks` (b1687) yang disimpan; `Date` (b790) tidak pernah ditulis
// `RejectOSClaimLife_Act`, dan `PIC` (b975) adalah pelaku itu sendiri.
type permintaanTolakJSON struct {
	Komentar string `json:"komentar"`
}

// uraiAlasanTolak membaca `Remarks` dari badan JSON; wajib-isinya diperiksa
// services (satu tempat).
func uraiAlasanTolak(r *http.Request) (string, error) {
	var masuk permintaanTolakJSON
	if err := json.NewDecoder(r.Body).Decode(&masuk); err != nil {
		return "", err
	}
	return masuk.Komentar, nil
}

// tolakBaris melayani
// POST /api/klaim-life/{id}/adjustment/{adjId}/tolak.
func tolakBaris(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat.Tulis(w, http.StatusServiceUnavailable, "database belum dikonfigurasi")
			return
		}
		alasan, err := uraiAlasanTolak(r)
		if err != nil {
			galat.Tulis(w, http.StatusBadRequest, "badan permintaan bukan JSON yang sah")
			return
		}
		err = svc.Status().DenganJejak(jejak.PerekamJejakOracle(svc)).Tolak(r.Context(), inti.PelakuDari(r, stubPelaku),
			r.PathValue("id"), r.PathValue("adjId"), alasan, time.Now())

		switch {
		case errors.Is(err, inti.ErrTanpaIdentitas):
			// 401: yang kurang identitasnya. Dipisah dari 403 karena keduanya
			// menjawab pertanyaan yang berbeda - dan ronde pertama berkas ini
			// menjawab "bukan ReasLifeAdmin" kepada pemanggil yang sebenarnya
			// hanya belum menyebut dirinya.
			galat.Tulis(w, http.StatusUnauthorized, "permintaan tanpa identitas pelaku ditolak")
			return
		case errors.Is(err, inti.ErrTanpaWewenang):
			// 403: identitasnya ada, perannya yang kurang.
			galat.Tulis(w, http.StatusForbidden,
				"hanya ReasLifeAdmin yang dapat menolak baris Outstanding")
			return
		case errors.Is(err, services.ErrBarisSudahFinal):
			// 409: baris itu sudah punya keputusan. Bukan galat permintaan,
			// melainkan bentrokan dengan keadaan yang sudah ada.
			galat.Tulis(w, http.StatusConflict,
				"baris sudah diputus dan tidak dapat ditolak lagi")
			return
		case errors.Is(err, services.ErrKlaimBelumBernomor):
			// 422: datanya yang belum lengkap, bukan permintaannya yang salah.
			//
			// ⚠️ Sampai butir o diputuskan, TIDAK ADA klaim yang bernomor,
			// sehingga jalur nyata selalu berhenti di sini. Itu keadaan yang
			// benar - bukan kerusakan - dan dinyatakan di tiket.
			galat.Tulis(w, http.StatusUnprocessableEntity,
				"klaim belum bernomor; penolakan baris menunggu nomor klaim")
			return
		case errors.Is(err, services.ErrKasusSudahTertutup):
			galat.Tulis(w, http.StatusConflict, "kasus sudah ditutup dan tidak dapat diubah")
			return
		case errors.Is(err, galat.ErrPermintaanTidakSah):
			galat.Tulis(w, http.StatusBadRequest, err.Error())
			return
		case err != nil:
			galat.Tulis(w, http.StatusInternalServerError, "gagal menolak baris")
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
