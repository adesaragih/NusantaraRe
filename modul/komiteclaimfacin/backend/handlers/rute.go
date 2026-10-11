// Package handlers adalah pintu HTTP modul Komite Claim Fac In: menerima permintaan, membaca pelaku, memanggil services,
// dan menerjemahkan galat ke kode HTTP. Nol aturan dagang, nol SQL.
//
// Rute - prefix `/api/komite-claim-fac-in` (`/api/komite` milik Komite Claim Life, `/api/komite-claim-prop` /
// `/api/komite-claim-non-prop` milik Komite Claim Prop / Non Prop). Ketiganya juga rute pinjaman pemegang menu
// `claimfacin` (`cmd/api/rakit.go` `ruteDipinjam`): modul ini TANPA menu.
//
//	GET  /api/komite-claim-fac-in/kasus                  daftar kerja penyetuju (worklist KomiteRouter)
//	GET  /api/komite-claim-fac-in/kasus/{id}             buka kasus (flow action ViewTransferDtl)
//	POST /api/komite-claim-fac-in/kasus/{id}/putuskan    Submit (finishAssignment -> KomitePostAct)
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
	"nusantarare/modul/komiteclaimfacin/backend/models"
	"nusantarare/modul/komiteclaimfacin/backend/services"
)

// Prefix adalah awalan rute modul ini.
const Prefix = "/api/komite-claim-fac-in"

// batasBadan - badan permintaan terbesar (isian Submit).
const batasBadan = 1 << 20

// DaftarkanRute memasang seluruh rute modul ke mux bersama.
func DaftarkanRute(mux *http.ServeMux, l *services.Layanan, stubPelaku bool) {
	h := &rute{l: l, stub: stubPelaku}
	mux.HandleFunc("GET "+Prefix+"/kasus", h.daftar)
	mux.HandleFunc("GET "+Prefix+"/kasus/{id}", h.buka)
	mux.HandleFunc("POST "+Prefix+"/kasus/{id}/putuskan", h.putuskan)
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

func tulisJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("komiteclaimfacin: menulis jawaban: %v", err)
	}
}

// tulisGalat menerjemahkan galat services ke kode HTTP.
func tulisGalat(w http.ResponseWriter, err error) {
	var v *services.GalatValidasi
	switch {
	case errors.As(err, &v):
		tulisJSON(w, http.StatusUnprocessableEntity, struct {
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
	case errors.Is(err, services.ErrKasusTertutup),
		errors.Is(err, services.ErrKeputusanBersamaan),
		errors.Is(err, services.ErrKlaimIndukTidakAda):
		galat.Tulis(w, http.StatusConflict, err.Error())
	default:
		// Teks galat basis data (ORA-, nama skema) tidak dikirim ke klien.
		log.Printf("komiteclaimfacin: %v", err)
		galat.Tulis(w, http.StatusInternalServerError, "gagal memproses permintaan Komite Claim Fac In")
	}
}

func (h *rute) daftar(w http.ResponseWriter, r *http.Request) {
	out, err := h.l.DaftarKerja(r.Context(), h.pelaku(r))
	if err != nil {
		tulisGalat(w, err)
		return
	}
	tulisJSON(w, http.StatusOK, out)
}

func (h *rute) buka(w http.ResponseWriter, r *http.Request) {
	ly, err := h.l.BukaKasus(r.Context(), h.pelaku(r), r.PathValue("id"))
	if err != nil {
		tulisGalat(w, err)
		return
	}
	tulisJSON(w, http.StatusOK, ly)
}

func (h *rute) putuskan(w http.ResponseWriter, r *http.Request) {
	badan, err := io.ReadAll(io.LimitReader(r.Body, batasBadan+1))
	if err != nil || len(badan) > batasBadan {
		galat.Tulis(w, http.StatusBadRequest, "badan permintaan tidak terbaca atau terlalu besar")
		return
	}
	var kep models.Keputusan
	if err := json.Unmarshal(badan, &kep); err != nil {
		galat.Tulis(w, http.StatusBadRequest, "isian keputusan bukan JSON yang sah")
		return
	}
	out, err := h.l.Putuskan(r.Context(), h.pelaku(r), r.PathValue("id"), kep)
	if err != nil {
		tulisGalat(w, err)
		return
	}
	tulisJSON(w, http.StatusOK, out)
}
