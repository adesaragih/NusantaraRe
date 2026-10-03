// Package handlers adalah pintu HTTP modul NB Treaty In - seam 1 spec §6.2.
// Tugasnya hanya menerima permintaan, membaca pelaku, memanggil services, dan
// menerjemahkan galat ke kode HTTP. Nol aturan dagang, nol SQL (AC 60).
//
// Rute - prefix `/api/nb-treaty-in` (MODUL.md):
//
//	GET  /api/nb-treaty-in/kasus                      daftar portal (SFAPortal_OpportunitiesList)
//	POST /api/nb-treaty-in/kasus                      Create (createWork)
//	GET  /api/nb-treaty-in/kasus/{id}                 buka assignment (pra-proses)
//	PUT  /api/nb-treaty-in/kasus/{id}                 Save layar admin
//	POST /api/nb-treaty-in/kasus/{id}/hitung          refresh berhitung (Count*, ...)
//	POST /api/nb-treaty-in/kasus/{id}/pilih-bisnis    Choose popup BusinessAndSOBList
//	POST /api/nb-treaty-in/kasus/{id}/pilih-sumber-bisnis  klik baris popup SOB (XOL Retro)
//	POST /api/nb-treaty-in/kasus/{id}/nomor-polis     GeneratePolicyNoTreaty_Act
//	POST /api/nb-treaty-in/kasus/{id}/kirim           finishAssignment
//	GET  /api/nb-treaty-in/kasus/{id}/riwayat         HISTORYAKSEPTASIPEGA
//	GET  /api/nb-treaty-in/bisnis                     grid popup BusinessAndSOBList
//	GET  /api/nb-treaty-in/sumber-bisnis              grid popup SOB (BrowseAgentHierarkiList_RD)
//	GET  /api/nb-treaty-in/acuan                      daftar pilihan layar
package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/modul/nbtreatyin/backend/models"
	"nusantarare/modul/nbtreatyin/backend/services"
)

// Prefix adalah awalan rute modul ini.
const Prefix = "/api/nb-treaty-in"

// batasBadan - badan permintaan terbesar (halaman kerja beserta daftarnya).
const batasBadan = 4 << 20

// DaftarkanRute memasang seluruh rute modul ke mux bersama.
func DaftarkanRute(mux *http.ServeMux, l *services.Layanan, stubPelaku bool) {
	h := &rute{l: l, stub: stubPelaku}
	mux.HandleFunc("GET "+Prefix+"/kasus", h.daftar)
	mux.HandleFunc("POST "+Prefix+"/kasus", h.buat)
	mux.HandleFunc("GET "+Prefix+"/kasus/{id}", h.buka)
	mux.HandleFunc("PUT "+Prefix+"/kasus/{id}", h.simpan)
	mux.HandleFunc("POST "+Prefix+"/kasus/{id}/hitung", h.hitung)
	mux.HandleFunc("POST "+Prefix+"/kasus/{id}/pilih-bisnis", h.pilihBisnis)
	mux.HandleFunc("POST "+Prefix+"/kasus/{id}/pilih-sumber-bisnis", h.pilihSumberBisnis)
	mux.HandleFunc("POST "+Prefix+"/kasus/{id}/nomor-polis", h.nomorPolis)
	mux.HandleFunc("POST "+Prefix+"/kasus/{id}/kirim", h.kirim)
	mux.HandleFunc("GET "+Prefix+"/kasus/{id}/riwayat", h.riwayat)
	mux.HandleFunc("GET "+Prefix+"/bisnis", h.bisnis)
	mux.HandleFunc("GET "+Prefix+"/sumber-bisnis", h.sumberBisnis)
	mux.HandleFunc("GET "+Prefix+"/acuan", h.acuan)
}

// Router menyusun mux tersendiri - untuk uji.
func Router(l *services.Layanan, stubPelaku bool) http.Handler {
	mux := http.NewServeMux()
	DaftarkanRute(mux, l, stubPelaku)
	return mux
}

type rute struct {
	l    *services.Layanan
	stub bool
}

func (h *rute) pelaku(r *http.Request) inti.Pelaku { return inti.PelakuDari(r, h.stub) }

// tulisGalat menerjemahkan galat services ke kode HTTP.
func tulisGalat(w http.ResponseWriter, err error) {
	var v *services.GalatValidasi
	switch {
	case errors.As(err, &v):
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusUnprocessableEntity)
		_ = json.NewEncoder(w).Encode(struct {
			Galat string   `json:"galat"`
			Pesan []string `json:"pesan"`
		}{"validasi layar gagal: " + strings.Join(v.Pesan, "; "), v.Pesan})
	case errors.Is(err, services.ErrTanpaOracle):
		galat.Tulis(w, http.StatusServiceUnavailable, "database belum dikonfigurasi")
	case errors.Is(err, inti.ErrTanpaIdentitas):
		galat.Tulis(w, http.StatusUnauthorized, err.Error())
	case errors.Is(err, inti.ErrTanpaWewenang):
		galat.Tulis(w, http.StatusForbidden, err.Error())
	case errors.Is(err, services.ErrKasusTidakAda):
		galat.Tulis(w, http.StatusNotFound, err.Error())
	case errors.Is(err, services.ErrPermintaanTidakSah):
		galat.Tulis(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, services.ErrDataKontrakTidakAda),
		errors.Is(err, services.ErrTipeNomorKosong),
		errors.Is(err, services.ErrOJKKosong),
		errors.Is(err, services.ErrMasterXOLTidakAda),
		errors.Is(err, services.ErrMasterXOLRusak):
		// AC 37: kegagalan membaca data kontrak DITAMPILKAN kepada pengguna.
		galat.Tulis(w, http.StatusUnprocessableEntity, err.Error())
	case errors.Is(err, services.ErrKasusTertutup),
		errors.Is(err, services.ErrTahapBerubah),
		errors.Is(err, services.ErrGenerasiTertutup),
		errors.Is(err, services.ErrNomorPolisSudahAda),
		errors.Is(err, services.ErrTindakanTakAdaDiPosisi):
		galat.Tulis(w, http.StatusConflict, err.Error())
	default:
		// Teks galat basis data (ORA-, nama skema) tidak dikirim ke klien.
		log.Printf("nbtreatyin: %v", err)
		galat.Tulis(w, http.StatusInternalServerError, "gagal memproses permintaan NB Treaty In")
	}
}

// bacaJSON membaca badan permintaan JSON (kosong = nilai nol).
func bacaJSON(w http.ResponseWriter, r *http.Request, tujuan any) bool {
	r.Body = http.MaxBytesReader(w, r.Body, batasBadan)
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(tujuan); err != nil && !errors.Is(err, io.EOF) {
		galat.Tulis(w, http.StatusBadRequest, "badan permintaan bukan JSON yang sah: "+err.Error())
		return false
	}
	return true
}

// badanHalaman - badan permintaan yang membawa halaman kerja.
type badanHalaman struct {
	Halaman *models.Halaman `json:"halaman"`
	// IDDetail - ID baris view kontrak (pilih bisnis).
	IDDetail string `json:"idDetail"`
}

func (h *rute) daftar(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	out, err := h.l.DaftarKasus(r.Context(), h.pelaku(r), models.SaringanKasus{
		Cari: strings.TrimSpace(q.Get("cari")), Posisi: strings.TrimSpace(q.Get("posisi")),
	})
	if err != nil {
		tulisGalat(w, err)
		return
	}
	if out == nil {
		out = []models.RingkasanKasus{}
	}
	galat.TulisJSON(w, out)
}

func (h *rute) buat(w http.ResponseWriter, r *http.Request) {
	k, err := h.l.BuatKasus(r.Context(), h.pelaku(r))
	if err != nil {
		tulisGalat(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(k)
}

func (h *rute) buka(w http.ResponseWriter, r *http.Request) {
	ly, err := h.l.BukaKasus(r.Context(), h.pelaku(r), r.PathValue("id"))
	if err != nil {
		tulisGalat(w, err)
		return
	}
	galat.TulisJSON(w, ly)
}

func (h *rute) simpan(w http.ResponseWriter, r *http.Request) {
	var b badanHalaman
	if !bacaJSON(w, r, &b) {
		return
	}
	ly, err := h.l.SimpanDraf(r.Context(), h.pelaku(r), r.PathValue("id"), b.Halaman)
	if err != nil {
		tulisGalat(w, err)
		return
	}
	galat.TulisJSON(w, ly)
}

func (h *rute) hitung(w http.ResponseWriter, r *http.Request) {
	var b services.PermintaanHitung
	if !bacaJSON(w, r, &b) {
		return
	}
	ly, err := h.l.Hitung(r.Context(), h.pelaku(r), r.PathValue("id"), b)
	if err != nil {
		tulisGalat(w, err)
		return
	}
	galat.TulisJSON(w, ly)
}

func (h *rute) pilihBisnis(w http.ResponseWriter, r *http.Request) {
	var b badanHalaman
	if !bacaJSON(w, r, &b) {
		return
	}
	ly, err := h.l.PilihBisnis(r.Context(), h.pelaku(r), r.PathValue("id"), b.IDDetail, b.Halaman)
	if err != nil {
		tulisGalat(w, err)
		return
	}
	galat.TulisJSON(w, ly)
}

func (h *rute) nomorPolis(w http.ResponseWriter, r *http.Request) {
	var b badanHalaman
	if !bacaJSON(w, r, &b) {
		return
	}
	n, err := h.l.TerbitkanNomor(r.Context(), h.pelaku(r), r.PathValue("id"), b.Halaman)
	if err != nil {
		tulisGalat(w, err)
		return
	}
	galat.TulisJSON(w, n)
}

func (h *rute) kirim(w http.ResponseWriter, r *http.Request) {
	var b badanHalaman
	if !bacaJSON(w, r, &b) {
		return
	}
	k, err := h.l.Kirim(r.Context(), h.pelaku(r), r.PathValue("id"), b.Halaman)
	if err != nil {
		tulisGalat(w, err)
		return
	}
	galat.TulisJSON(w, k)
}

func (h *rute) riwayat(w http.ResponseWriter, r *http.Request) {
	out, err := h.l.RiwayatKasus(r.Context(), h.pelaku(r), r.PathValue("id"))
	if err != nil {
		tulisGalat(w, err)
		return
	}
	if out == nil {
		galat.TulisJSON(w, []struct{}{})
		return
	}
	galat.TulisJSON(w, out)
}

func (h *rute) bisnis(w http.ResponseWriter, r *http.Request) {
	out, err := h.l.DaftarBisnis(r.Context(), h.pelaku(r), r.URL.Query().Get("cari"))
	if err != nil {
		tulisGalat(w, err)
		return
	}
	if out == nil {
		out = []models.BarisKontrak{}
	}
	galat.TulisJSON(w, out)
}

func (h *rute) acuan(w http.ResponseWriter, r *http.Request) {
	a, err := h.l.DaftarAcuan(r.Context(), h.pelaku(r))
	if err != nil {
		tulisGalat(w, err)
		return
	}
	galat.TulisJSON(w, a)
}
