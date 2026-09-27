package handlers

// Pintu HTTP gerbang `Close Claim` - kelompok Detail & Tutup.
//
//	`GET /api/klaim-life/{id}/boleh-tutup`
//
// Meniru `Section/CloseClaim_Section.xml` b1081 `Close Claim` ->
// `pyActivity` b1101 `ProtectCloseClaim_act`.
//
// ⛔ GET, bukan POST, dan itu keputusan yang disengaja: activity Pega
// memeriksa LALU menyelesaikan penugasan, sedangkan yang dibangun di sini
// baru pemeriksaannya. Memberinya POST akan menjanjikan penutupan yang belum
// terjadi - dan tombol yang menjanjikan lebih daripada yang ia lakukan adalah
// cacat yang paling mahal ditemukan belakangan.
//
// ⚠️ Sisi "menyelesaikan penugasan" (`Call FinishAssignment` b836) menunggu
// pembacaan alurnya: tahap mana yang menjadi tujuan sesudah tutup belum
// dibaca dari `Flow/`, dan menebaknya akan memindahkan kasus ke tempat yang
// salah tanpa satu pun galat.
//
// Dibaca sesudah: tahap.go.

import (
	"encoding/json"
	"errors"
	"net/http"

	"nusantarare/internal/services"
)

// bolehTutup melayani GET /api/klaim-life/{id}/boleh-tutup.
func bolehTutup(svc *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat(w, http.StatusServiceUnavailable, "database belum dikonfigurasi")
			return
		}
		hasil, err := svc.Tutup().Periksa(r.Context(), r.PathValue("id"))
		switch {
		case err == nil:
		case errors.Is(err, services.ErrKlaimTidakAda):
			galat(w, http.StatusNotFound, "klaim tidak ditemukan")
			return
		default:
			galat(w, http.StatusInternalServerError, "gagal memeriksa kesiapan tutup klaim")
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(hasil)
	}
}
