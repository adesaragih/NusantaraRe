package handlers

// Pintu HTTP modul Komite Claim Life - tiket 01.
//
//	GET /api/komite         Inbox Komite milik pelaku (worklist `KomiteRouter`)
//	GET /api/komite/{id}    satu kasus beserta tangganya
//	POST /api/komite/{id}/keputusan  {"keputusan":"1|2","komentar":"…"} (tiket 02)
//
// Nol aturan dagang di sini; siapa melihat apa diputuskan
// `services/komite_inbox.go`.
//
// ⚠️ Berkas ini MILIK modul Komite. Penyerahan dari Claim Life (`komite.go`,
// A2) tidak disentuh.

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"nusantarare/internal/services"
)

// inboxKomite melayani GET /api/komite.
func inboxKomite(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat(w, http.StatusServiceUnavailable, "database belum dikonfigurasi")
			return
		}
		halaman, _ := strconv.Atoi(r.URL.Query().Get("halaman"))
		ukuran, _ := strconv.Atoi(r.URL.Query().Get("ukuran"))
		hal, err := svc.InboxKomite().Ambil(r.Context(), pelakuDari(r, stubPelaku), halaman, ukuran)
		if jawabGalatKomite(w, err) {
			return
		}
		tulisJSONPolis(w, hal)
	}
}

// kasusKomite melayani GET /api/komite/{id}.
func kasusKomite(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat(w, http.StatusServiceUnavailable, "database belum dikonfigurasi")
			return
		}
		k, err := svc.InboxKomite().Kasus(r.Context(), pelakuDari(r, stubPelaku), r.PathValue("id"))
		if jawabGalatKomite(w, err) {
			return
		}
		tulisJSONPolis(w, k)
	}
}

// isiKeputusanKomite adalah badan `POST .../keputusan`.
type isiKeputusanKomite struct {
	Keputusan string `json:"keputusan"`
	Komentar  string `json:"komentar"`
}

// putuskanKomite melayani POST /api/komite/{id}/keputusan - tombol `Submit`
// `ShowTransfer` b34722.
func putuskanKomite(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat(w, http.StatusServiceUnavailable, "database belum dikonfigurasi")
			return
		}
		var isi isiKeputusanKomite
		if err := json.NewDecoder(r.Body).Decode(&isi); err != nil {
			galat(w, http.StatusBadRequest, "badan permintaan bukan JSON yang sah")
			return
		}
		hasil, err := svc.KeputusanKomite().
			DenganJejak(services.PerekamJejakOracle(svc)).
			Putuskan(r.Context(), pelakuDari(r, stubPelaku), r.PathValue("id"),
				isi.Keputusan, isi.Komentar, time.Now())
		if jawabGalatKomite(w, err) {
			return
		}
		tulisJSONPolis(w, hasil)
	}
}

// jawabGalatKomite menerjemahkan galat services menjadi kode HTTP.
//
// Mengembalikan true bila permintaan SUDAH dijawab.
func jawabGalatKomite(w http.ResponseWriter, err error) bool {
	switch {
	case err == nil:
		return false
	case errors.Is(err, services.ErrTanpaIdentitas):
		galat(w, http.StatusUnauthorized, "permintaan tanpa identitas pelaku ditolak")
	case errors.Is(err, services.ErrTanpaWewenang):
		// 403, bukan 404: kasusnya ada, pelakunya bukan anggota tangganya.
		galat(w, http.StatusForbidden, "Anda bukan anggota tangga komite kasus ini")
	case errors.Is(err, services.ErrKasusKomiteTakDitemukan):
		galat(w, http.StatusNotFound, "kasus komite tidak ditemukan")
	case errors.Is(err, services.ErrKeputusanKomiteTidakDikenal),
		errors.Is(err, services.ErrPermintaanTidakSah):
		galat(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, services.ErrTanggaKomiteBerhenti),
		errors.Is(err, services.ErrKeputusanKomiteBersamaan),
		errors.Is(err, services.ErrKasusSudahTertutup):
		galat(w, http.StatusConflict, err.Error())
	case errors.Is(err, services.ErrPenyelesaianAkhirBelumAda):
		// 501: permintaannya sah, bagian sistemnya yang belum dibangun.
		galat(w, http.StatusNotImplemented, err.Error())
	default:
		galat(w, http.StatusInternalServerError, "gagal memproses permintaan komite")
	}
	return true
}
