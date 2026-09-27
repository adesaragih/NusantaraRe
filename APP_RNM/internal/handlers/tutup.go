package handlers

// Pintu HTTP gerbang `Close Claim` - kelompok Detail & Tutup.
//
// Meniru `Section/CloseClaim_Section.xml` b1081 `Close Claim` ->
// `pyActivity` b1101 `ProtectCloseClaim_act`.
//
// DUA rute, dan pemisahannya disengaja:
//
//	GET  /api/klaim-life/{id}/boleh-tutup  - memeriksa saja
//	POST /api/klaim-life/{id}/tutup        - memeriksa LALU menutup
//
// ⛔ Yang pertama tidak pernah mengubah apa pun: layar memanggilnya untuk
// memutuskan apakah tombol `Close Claim` pantas ditawarkan, dan permintaan
// yang menjawab pertanyaan tidak boleh menjawabnya dengan perbuatan.
//
// ⛔ Yang kedua lahir dari keputusan bb (27-09-2026). Ronde sebelumnya
// menahan POST karena tujuan sesudah tutup belum terbaca dari `Flow/`.
// Alurnya kini sudah dibaca utuh: `Register_Flow.xml` punya 12 konektor dan
// TIDAK SATU PUN bernama `CloseClaim` - ia local action, dan hanya shape
// `End1` (b883) yang menetapkan status kerja `Resolved-Completed` (b899).
// Perilaku mesin Pega untuk `FinishAssignment` dari local action tanpa
// konektor senama tetap TIDAK dapat diturunkan dari ekspor - itu OQ-I. Yang
// ditiru adalah niat nyatanya, dan niatnya tidak ambigu: label `Close
// Claim`, konfirmasi "Are you sure want to Close Claim?", dan gerbang yang
// menahan setiap peserta yang belum diaksep.
//
// Dibaca sesudah: tahap.go.

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

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

// jawabanPenghalang adalah badan 409: pesan DAN daftar penghalangnya.
//
// ⛔ Keduanya, bukan salah satu. Layar memerlukan daftarnya untuk menunjuk
// peserta mana yang menahan; pemakai memerlukan kalimatnya. Mengirim pesan
// saja memaksa layar memanggil `boleh-tutup` lagi hanya untuk menampilkan
// apa yang barusan diputuskan server.
//
// ⚠️ Kunci `galat` SAMA dengan amplop galat biasa (`galat()`), sehingga
// pembaca galat di klien tidak perlu tahu rute mana yang menjawab. Cacat
// `galat` vs `error` yang pernah ada lahir persis dari dua bentuk amplop.
type jawabanPenghalang struct {
	Galat      string                      `json:"galat"`
	Penghalang []services.PenghalangTampil `json:"penghalang"`
}

// tutupKlaim melayani POST /api/klaim-life/{id}/tutup.
func tutupKlaim(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat(w, http.StatusServiceUnavailable, "database belum dikonfigurasi")
			return
		}
		err := svc.Tutup().DenganJejak(services.PerekamJejakOracle(svc)).
			Tutup(r.Context(), pelakuDari(r, stubPelaku), r.PathValue("id"), time.Now())

		var halangan *services.GalatPenghalang
		switch {
		case err == nil:
			// 204: tidak ada badan. Layar kembali ke Inbox sesudahnya -
			// padanan `closeContainer` b1129.
			w.WriteHeader(http.StatusNoContent)
		case errors.As(err, &halangan):
			// ⛔ 409 berisi SELURUH penghalang. Pega memasang pesannya di
			// dalam loop, sekali per peserta yang tertandai; mengirim satu
			// saja memaksa pemakai menutup berulang kali dan menemukan satu
			// penghalang baru setiap kali.
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusConflict)
			_ = json.NewEncoder(w).Encode(jawabanPenghalang{
				Galat:      "klaim belum boleh ditutup: ada peserta yang belum diaksep",
				Penghalang: halangan.Penghalang,
			})
		case errors.Is(err, services.ErrTanpaIdentitas):
			galat(w, http.StatusUnauthorized, "permintaan tanpa identitas pelaku ditolak")
		case errors.Is(err, services.ErrTanpaWewenang):
			galat(w, http.StatusForbidden,
				"hanya pemegang tahap kasus ini yang dapat menutupnya")
		case errors.Is(err, services.ErrTahapTidakMenutup):
			// ⛔ 409, bukan 403: perannya mungkin benar, TAHAPNYA yang tidak
			// menawarkan tombol itu. `pyLocalAction>CloseClaim` hanya ada di
			// InputOSClaimLife dan InputAkseptasiClaimLife.
			galat(w, http.StatusConflict,
				"Close Claim hanya ada pada tahap Outstanding Claim dan Claim Analis")
		case errors.Is(err, services.ErrKasusSudahTertutup):
			galat(w, http.StatusConflict, "kasus ini sudah ditutup")
		case errors.Is(err, services.ErrKlaimTidakAda):
			galat(w, http.StatusNotFound, "klaim tidak ditemukan")
		case errors.Is(err, services.ErrTahapTidakDikenal):
			galat(w, http.StatusConflict, "tahap kasus ini tidak dikenal")
		default:
			galat(w, http.StatusInternalServerError, "gagal menutup klaim")
		}
	}
}
