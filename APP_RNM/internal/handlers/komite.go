package handlers

// Pintu HTTP penyerahan ke Komite - tiket 10.
//
// Nol aturan dagang di sini: berkas ini menerjemahkan HTTP ke panggilan
// services dan sebaliknya. Yang memutuskan boleh atau tidak adalah services.
//
// Dibaca sesudah: handlers.go dan services/komite.go.

import (
	"errors"
	"net/http"
	"time"

	"nusantarare/internal/services"
)

// serahkanKomite melayani
// POST /api/klaim-life/{id}/peserta/{pesertaId}/adjustment/{adjId}/komite.
func serahkanKomite(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat(w, http.StatusServiceUnavailable, "database belum dikonfigurasi")
			return
		}
		err := svc.Komite().Serahkan(r.Context(), pelakuDari(r, stubPelaku),
			r.PathValue("id"), r.PathValue("pesertaId"), r.PathValue("adjId"),
			time.Now())

		switch {
		case errors.Is(err, services.ErrTanpaIdentitas):
			galat(w, http.StatusUnauthorized, "permintaan tanpa identitas pelaku ditolak")
			return
		case errors.Is(err, services.ErrTanpaWewenang):
			// 403: identitasnya ada, perannya yang kurang untuk Type ini.
			galat(w, http.StatusForbidden,
				"peran Anda tidak berwenang menyerahkan baris bertipe ini ke Komite")
			return
		case errors.Is(err, services.ErrTypeTidakDikenal):
			galat(w, http.StatusUnprocessableEntity,
				"Type klaim tidak dikenal; penyerahan ke Komite menuntut Type yang sah")
			return
		case errors.Is(err, services.ErrRekeningBelumLengkap):
			// 422: ⛔ Pesannya diteruskan APA ADANYA. Ia satu-satunya pesan di
			// pintu ini yang BUKAN karangan kita: teksnya persis seperti sistem
			// lama, dan pengguna lama mengenalinya.
			galat(w, http.StatusUnprocessableEntity, err.Error())
			return
		case errors.Is(err, services.ErrBarisSudahDiserahkan):
			galat(w, http.StatusConflict,
				"baris sudah pernah diserahkan ke Komite")
			return
		case errors.Is(err, services.ErrBarisBukanOutstanding):
			galat(w, http.StatusConflict,
				"hanya baris Outstanding yang dapat diserahkan ke Komite")
			return
		case errors.Is(err, services.ErrMataUangKlaimCampur):
			galat(w, http.StatusUnprocessableEntity,
				"baris pada klaim ini bermata uang campur; penyerahan menuntut mata uang tunggal")
			return
		case errors.Is(err, services.ErrRosterKomiteKosong):
			galat(w, http.StatusUnprocessableEntity,
				"tidak ada tingkat komite yang menutup nilai klaim ini")
			return
		case errors.Is(err, services.ErrRosterBelumDiputuskan),
			errors.Is(err, services.ErrKasusKomiteBelumDiputuskan),
			errors.Is(err, services.ErrJejakBelumDiputuskan):
			// 501: bukan kerusakan melainkan keputusan yang belum diambil.
			// Pesan dalamnya tidak diteruskan - ia menyebut nama objek basis
			// data (CLAUDE.md bab 4 butir 10).
			galat(w, http.StatusNotImplemented,
				"penyerahan ke Komite belum dapat disimpan: tempatnya belum diputuskan work owner")
			return
		case errors.Is(err, services.ErrPermintaanTidakSah):
			galat(w, http.StatusBadRequest, err.Error())
			return
		case err != nil:
			galat(w, http.StatusInternalServerError, "gagal menyerahkan baris ke Komite")
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
