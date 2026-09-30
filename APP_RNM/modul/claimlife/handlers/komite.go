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

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/inti/backend/jejak"
	"nusantarare/inti/backend/kontrak"
	"nusantarare/modul/claimlife/services"
)

// serahkanKomite melayani
// POST /api/klaim-life/{id}/peserta/{pesertaId}/adjustment/{adjId}/komite.
func serahkanKomite(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat.Tulis(w, http.StatusServiceUnavailable, "database belum dikonfigurasi")
			return
		}
		err := svc.Komite().DenganJejak(jejak.PerekamJejakOracle(svc)).
			DenganRoster(services.RosterKomiteOracle(svc)).
			DenganKasus(services.KasusKomiteOracle(svc)).
			DenganPenyalur(services.PenyalurClaimLifeOracle(svc)).
			Serahkan(r.Context(), inti.PelakuDari(r, stubPelaku),
				r.PathValue("id"), r.PathValue("pesertaId"), r.PathValue("adjId"),
				time.Now())

		switch {
		case errors.Is(err, inti.ErrTanpaIdentitas):
			galat.Tulis(w, http.StatusUnauthorized, "permintaan tanpa identitas pelaku ditolak")
			return
		case errors.Is(err, inti.ErrTanpaWewenang):
			// 403: identitasnya ada, perannya yang kurang untuk Type ini.
			galat.Tulis(w, http.StatusForbidden,
				"peran Anda tidak berwenang menyerahkan baris bertipe ini ke Komite")
			return
		case errors.Is(err, services.ErrTypeTidakDikenal):
			galat.Tulis(w, http.StatusUnprocessableEntity,
				"Type klaim tidak dikenal; penyerahan ke Komite menuntut Type yang sah")
			return
		case errors.Is(err, services.ErrRekeningBelumLengkap):
			// 422: ⛔ Pesannya diteruskan APA ADANYA. Ia satu-satunya pesan di
			// pintu ini yang BUKAN karangan kita: teksnya persis seperti sistem
			// lama, dan pengguna lama mengenalinya.
			galat.Tulis(w, http.StatusUnprocessableEntity, err.Error())
			return
		case errors.Is(err, services.ErrBarisSudahDiserahkan):
			galat.Tulis(w, http.StatusConflict,
				"baris sudah pernah diserahkan ke Komite")
			return
		case errors.Is(err, services.ErrBarisBukanOutstanding):
			galat.Tulis(w, http.StatusConflict,
				"hanya baris Outstanding yang dapat diserahkan ke Komite")
			return
		case errors.Is(err, services.ErrMataUangKlaimCampur):
			galat.Tulis(w, http.StatusUnprocessableEntity,
				"baris pada klaim ini bermata uang campur; penyerahan menuntut mata uang tunggal")
			return
		case errors.Is(err, services.ErrRosterKomiteKosong):
			galat.Tulis(w, http.StatusUnprocessableEntity,
				"tidak ada tingkat komite yang menutup nilai klaim ini")
			return
		case errors.Is(err, kontrak.ErrKasusSudahTertutup):
			galat.Tulis(w, http.StatusConflict, "kasus sudah ditutup dan tidak dapat diubah")
			return
		case errors.Is(err, galat.ErrPermintaanTidakSah):
			galat.Tulis(w, http.StatusBadRequest, err.Error())
			return
		case err != nil:
			galat.Tulis(w, http.StatusInternalServerError, "gagal menyerahkan baris ke Komite")
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
