package handlers

// Pintu HTTP modul Komite Claim Life - tiket 01.
//
//	GET /api/komite         Inbox Komite milik pelaku (worklist `KomiteRouter`)
//	GET /api/komite/{id}    satu kasus beserta tangganya
//	POST /api/komite/{id}/keputusan  {"keputusan":"1|2","komentar":"…"} (tiket 02)
//	POST /api/komite/{id}/eskalasi   naik SATU tingkat, admin saja (tiket 03)
//	GET  /api/komite/laporan-harian  efek "perlu intervensi" hari ini (tiket 08)
//	GET  /api/komite/{id}/riwayat    riwayat tangga + eskalasi, siapa pun (tiket 09)
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

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/inti/backend/jejak"
	"nusantarare/inti/backend/kontrak"
	"nusantarare/modul/komiteclaimlife/backend/services"
)

// inboxKomite melayani GET /api/komite.
func inboxKomite(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat.Tulis(w, http.StatusServiceUnavailable, "database belum dikonfigurasi")
			return
		}
		halaman, _ := strconv.Atoi(r.URL.Query().Get("halaman"))
		ukuran, _ := strconv.Atoi(r.URL.Query().Get("ukuran"))
		hal, err := svc.InboxKomite().Ambil(r.Context(), inti.PelakuDari(r, stubPelaku), halaman, ukuran)
		if jawabGalatKomite(w, err) {
			return
		}
		galat.TulisJSON(w, hal)
	}
}

// kasusKomite melayani GET /api/komite/{id}.
func kasusKomite(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat.Tulis(w, http.StatusServiceUnavailable, "database belum dikonfigurasi")
			return
		}
		k, err := svc.InboxKomite().Kasus(r.Context(), inti.PelakuDari(r, stubPelaku), r.PathValue("id"))
		if jawabGalatKomite(w, err) {
			return
		}
		galat.TulisJSON(w, k)
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
			galat.Tulis(w, http.StatusServiceUnavailable, "database belum dikonfigurasi")
			return
		}
		var isi isiKeputusanKomite
		if err := json.NewDecoder(r.Body).Decode(&isi); err != nil {
			galat.Tulis(w, http.StatusBadRequest, "badan permintaan bukan JSON yang sah")
			return
		}
		hasil, err := svc.KeputusanKomite().
			DenganJejak(jejak.PerekamJejakOracle(svc)).
			DenganPenyelesaiAkhir(services.PenyelesaiAkhirKomiteOracle(svc)).
			Putuskan(r.Context(), inti.PelakuDari(r, stubPelaku), r.PathValue("id"),
				isi.Keputusan, isi.Komentar, time.Now())
		if jawabGalatKomite(w, err) {
			return
		}
		galat.TulisJSON(w, hasil)
	}
}

// eskalasiKomite melayani POST /api/komite/{id}/eskalasi.
//
// ⛔ TANPA badan: tingkat tujuan tidak dapat disebut pemanggil - selalu
// tingkat berjalan + 1. "Turun" atau "lompat" karena itu tidak dapat diminta.
func eskalasiKomite(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat.Tulis(w, http.StatusServiceUnavailable, "database belum dikonfigurasi")
			return
		}
		hasil, err := svc.KeputusanKomite().
			DenganJejak(jejak.PerekamJejakOracle(svc)).
			Eskalasi(r.Context(), inti.PelakuDari(r, stubPelaku), r.PathValue("id"), time.Now())
		if jawabGalatKomite(w, err) {
			return
		}
		galat.TulisJSON(w, hasil)
	}
}

// laporanHarianKomite melayani GET /api/komite/laporan-harian.
func laporanHarianKomite(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat.Tulis(w, http.StatusServiceUnavailable, "database belum dikonfigurasi")
			return
		}
		l, err := svc.InboxKomite().LaporanHarian(r.Context(), inti.PelakuDari(r, stubPelaku), time.Now())
		if jawabGalatKomite(w, err) {
			return
		}
		galat.TulisJSON(w, l)
	}
}

// riwayatKomite melayani GET /api/komite/{id}/riwayat.
func riwayatKomite(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat.Tulis(w, http.StatusServiceUnavailable, "database belum dikonfigurasi")
			return
		}
		rw, err := svc.InboxKomite().Riwayat(r.Context(), inti.PelakuDari(r, stubPelaku), r.PathValue("id"))
		if jawabGalatKomite(w, err) {
			return
		}
		galat.TulisJSON(w, rw)
	}
}

// jawabGalatKomite menerjemahkan galat services menjadi kode HTTP.
//
// Mengembalikan true bila permintaan SUDAH dijawab.
func jawabGalatKomite(w http.ResponseWriter, err error) bool {
	switch {
	case err == nil:
		return false
	case errors.Is(err, inti.ErrTanpaIdentitas):
		galat.Tulis(w, http.StatusUnauthorized, "permintaan tanpa identitas pelaku ditolak")
	case errors.Is(err, inti.ErrTanpaWewenang):
		// 403, bukan 404: kasusnya ada, pelakunya bukan anggota tangganya.
		galat.Tulis(w, http.StatusForbidden, "Anda bukan anggota tangga komite kasus ini")
	case errors.Is(err, services.ErrKasusKomiteTakDitemukan):
		galat.Tulis(w, http.StatusNotFound, "kasus komite tidak ditemukan")
	case errors.Is(err, services.ErrKeputusanKomiteTidakDikenal),
		errors.Is(err, galat.ErrPermintaanTidakSah):
		galat.Tulis(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, services.ErrTanggaKomiteBerhenti),
		errors.Is(err, services.ErrEskalasiTanpaTingkatAtas),
		errors.Is(err, services.ErrKeputusanKomiteBersamaan),
		errors.Is(err, kontrak.ErrKasusSudahTertutup):
		galat.Tulis(w, http.StatusConflict, err.Error())
	case errors.Is(err, kontrak.ErrNomorAkseptasiBerganda),
		errors.Is(err, kontrak.ErrKodeBisnisBelumTersimpan):
		// 409: keadaan DATA (nomor bertabrakan / kode bisnis kosong).
		galat.Tulis(w, http.StatusConflict, err.Error())
	case errors.Is(err, services.ErrPenyelesaianAkhirBelumAda):
		// 501: permintaannya sah, bagian sistemnya yang belum dibangun.
		galat.Tulis(w, http.StatusNotImplemented, err.Error())
	default:
		galat.Tulis(w, http.StatusInternalServerError, "gagal memproses permintaan komite")
	}
	return true
}

// Router menyusun rute modul ini SAJA, di mux sendiri.
//
// Refactor bentuk B (30-09-2026): dipakai uji HTTP modul ini, yang dulu
// memakai `Router` bersama milik seluruh aplikasi. Produksi tidak memakainya:
// `cmd/api` mendaftarkan `DaftarkanRute` ke mux yang sama dengan modul lain.
func Router(svc *services.Service, stubPelaku bool) http.Handler {
	mux := http.NewServeMux()
	DaftarkanRute(mux, svc, stubPelaku)
	return mux
}

// DaftarkanRute mendaftarkan seluruh rute modul Komite Claim Life.
//
// Refactor bentuk B (30-09-2026): rute ini dulu ditulis di
// `internal/handlers.Router`; kini dipanggil `modul/komiteclaimlife/backend/modul.go`, yang
// dipasang `cmd/api` bila modul ini aktif (MODUL_AKTIF).
func DaftarkanRute(mux *http.ServeMux, svc *services.Service, stubPelaku bool) {
	// Komite Claim Life tiket 01 - Inbox Komite dan satu kasus. Keduanya GET:
	// membaca saja; keputusan komite menyusul di tiket 02.
	mux.HandleFunc("GET /api/komite", inboxKomite(svc, stubPelaku))
	// Tiket 08 - literal mendahului `{id}` (pola paling spesifik menang).
	mux.HandleFunc("GET /api/komite/laporan-harian", laporanHarianKomite(svc, stubPelaku))
	mux.HandleFunc("GET /api/komite/{id}", kasusKomite(svc, stubPelaku))
	// Tiket 09 - riwayat tangga; membaca, bukan memutuskan.
	mux.HandleFunc("GET /api/komite/{id}/riwayat", riwayatKomite(svc, stubPelaku))
	// Tiket 02 - keputusan satu tingkat (`ShowTransfer` Submit).
	mux.HandleFunc("POST /api/komite/{id}/keputusan", putuskanKomite(svc, stubPelaku))
	// Tiket 03 - eskalasi naik satu tingkat (admin).
	mux.HandleFunc("POST /api/komite/{id}/eskalasi", eskalasiKomite(svc, stubPelaku))
}
