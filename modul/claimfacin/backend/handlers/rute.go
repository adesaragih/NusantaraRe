// Package handlers adalah pintu HTTP modul Claim Fac In. Tugasnya hanya menerima permintaan, membaca pelaku, memanggil
// services, dan menerjemahkan galat ke kode HTTP. Nol aturan dagang, nol SQL.
//
// Rute - prefix `/api/claim-fac-in` (MODUL.md):
//
//	GET  /api/claim-fac-in/kasus?daftar=saya|workbasket|selesai&cari=       halaman awal (worklist / workbasket / selesai)
//	POST /api/claim-fac-in/kasus                                            Add Claim (Start -> Input Register)
//	GET  /api/claim-fac-in/kasus/{id}?lihat=1                               buka assignment (pra-proses) / lihat saja
//	POST /api/claim-fac-in/kasus/{id}/aksi                                  aksi layar (refresh ber-activity, tombol, Choose)
//	GET  /api/claim-fac-in/kasus/{id}/pilihan/{jenis}?konteks=&indeks=&cari=&jenisCari=   isi pop-up / autocomplete
//	GET  /api/claim-fac-in/acuan                                            daftar pilihan bersama
//	GET  /api/claim-fac-in/hak                                              tab workbasket (anggota ReasKlaimTeknik)
//	GET  /api/claim-fac-in/kasus/{id}/lampiran                              kategori + dokumen klaim (lampiran.go)
//	POST /api/claim-fac-in/kasus/{id}/lampiran                              unggah lampiran (GCNMSaveAttachments)
//	GET  /api/claim-fac-in/kasus/{id}/lampiran/{lid}/isi                    View File
//	GET  /api/claim-fac-in/kasus/{id}/lampiran/{lid}/office                 View Office Online
//	POST /api/claim-fac-in/kasus/{id}/lampiran/{lid}/hapus                  Delete
//	POST /api/claim-fac-in/kasus/{id}/lampiran/kategori                     Change Category
package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/inti/backend/penyimpanan"
	"nusantarare/modul/claimfacin/backend/services"
)

// Prefix adalah awalan rute modul ini.
const Prefix = "/api/claim-fac-in"

// batasBadan - badan permintaan terbesar.
const batasBadan = 4 << 20

// DaftarkanRute memasang seluruh rute modul ke mux bersama.
func DaftarkanRute(mux *http.ServeMux, l *services.Layanan, stubPelaku bool) {
	h := &rute{l: l, stub: stubPelaku}
	mux.HandleFunc("GET "+Prefix+"/kasus", h.daftar)
	mux.HandleFunc("POST "+Prefix+"/kasus", h.buat)
	mux.HandleFunc("GET "+Prefix+"/kasus/{id}", h.buka)
	mux.HandleFunc("POST "+Prefix+"/kasus/{id}/aksi", h.aksi)
	mux.HandleFunc("GET "+Prefix+"/kasus/{id}/pilihan/{jenis}", h.pilihan)
	mux.HandleFunc("GET "+Prefix+"/acuan", h.acuan)
	mux.HandleFunc("GET "+Prefix+"/hak", h.hak)
	mux.HandleFunc("GET "+Prefix+"/kasus/{id}/lampiran", h.lampiran)
	mux.HandleFunc("POST "+Prefix+"/kasus/{id}/lampiran", h.unggahLampiran)
	mux.HandleFunc("GET "+Prefix+"/kasus/{id}/lampiran/{lid}/isi", h.unduhLampiran)
	mux.HandleFunc("GET "+Prefix+"/kasus/{id}/lampiran/{lid}/office", h.officeLampiran)
	mux.HandleFunc("POST "+Prefix+"/kasus/{id}/lampiran/{lid}/hapus", h.hapusLampiran)
	mux.HandleFunc("POST "+Prefix+"/kasus/{id}/lampiran/kategori", h.pindahKategoriLampiran)
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
		log.Printf("claimfacin: menulis jawaban: %v", err)
	}
}

// tulisGalat menerjemahkan galat services ke kode HTTP.
func tulisGalat(w http.ResponseWriter, err error) {
	var v *services.GalatValidasi
	switch {
	case errors.As(err, &v):
		tulisJSON(w, http.StatusUnprocessableEntity, struct {
			Galat string          `json:"galat"`
			Pesan []string        `json:"pesan"`
			Layar *services.Layar `json:"layar,omitempty"`
		}{"validasi layar gagal: " + strings.Join(v.Pesan, "; "), v.Pesan, v.Layar})
	case errors.Is(err, services.ErrTanpaOracle):
		galat.Tulis(w, http.StatusServiceUnavailable, "database belum dikonfigurasi")
	case errors.Is(err, penyimpanan.ErrObjekTidakAda):
		galat.Tulis(w, http.StatusConflict, err.Error())
	case errors.Is(err, penyimpanan.ErrBerkasDitolak):
		galat.Tulis(w, http.StatusUnprocessableEntity, err.Error())
	case errors.Is(err, penyimpanan.ErrStorageBelumSiap):
		galat.Tulis(w, http.StatusServiceUnavailable, err.Error())
	case errors.Is(err, penyimpanan.ErrStorageGagal):
		galat.Tulis(w, http.StatusBadGateway, err.Error())
	case errors.Is(err, inti.ErrTanpaIdentitas):
		galat.Tulis(w, http.StatusUnauthorized, err.Error())
	case errors.Is(err, inti.ErrTanpaWewenang):
		galat.Tulis(w, http.StatusForbidden, err.Error())
	case errors.Is(err, services.ErrKasusTidakAda):
		galat.Tulis(w, http.StatusNotFound, err.Error())
	case errors.Is(err, services.ErrPermintaanTidakSah):
		galat.Tulis(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, services.ErrKasusTertutup),
		errors.Is(err, services.ErrTahapBerubah),
		errors.Is(err, services.ErrAksiTertutup),
		errors.Is(err, services.ErrAdjustmentDiKomite),
		errors.Is(err, services.ErrTertunda):
		galat.Tulis(w, http.StatusConflict, err.Error())
	default:
		// Teks galat basis data (ORA-, nama skema) tidak dikirim ke klien.
		log.Printf("claimfacin: %v", err)
		galat.Tulis(w, http.StatusInternalServerError, "gagal memproses permintaan Claim Fac In")
	}
}

// hak - tab workbasket halaman awal (anggota ReasKlaimTeknik).
func (h *rute) hak(w http.ResponseWriter, r *http.Request) {
	out, err := h.l.Hak(r.Context(), h.pelaku(r))
	if err != nil {
		tulisGalat(w, err)
		return
	}
	tulisJSON(w, http.StatusOK, out)
}

func (h *rute) daftar(w http.ResponseWriter, r *http.Request) {
	out, err := h.l.DaftarKasus(r.Context(), h.pelaku(r), r.URL.Query().Get("daftar"), r.URL.Query().Get("cari"))
	if err != nil {
		tulisGalat(w, err)
		return
	}
	tulisJSON(w, http.StatusOK, out)
}

func (h *rute) buat(w http.ResponseWriter, r *http.Request) {
	ly, err := h.l.BuatKasus(r.Context(), h.pelaku(r))
	if err != nil {
		tulisGalat(w, err)
		return
	}
	tulisJSON(w, http.StatusCreated, ly)
}

func (h *rute) buka(w http.ResponseWriter, r *http.Request) {
	ly, err := h.l.BukaKasus(r.Context(), h.pelaku(r), r.PathValue("id"), r.URL.Query().Get("lihat") == "1")
	if err != nil {
		tulisGalat(w, err)
		return
	}
	tulisJSON(w, http.StatusOK, ly)
}

func (h *rute) aksi(w http.ResponseWriter, r *http.Request) {
	var req services.PermintaanAksi
	isi, err := io.ReadAll(io.LimitReader(r.Body, batasBadan+1))
	if err != nil || len(isi) > batasBadan {
		galat.Tulis(w, http.StatusBadRequest, "badan permintaan tidak terbaca")
		return
	}
	if err := json.Unmarshal(isi, &req); err != nil {
		galat.Tulis(w, http.StatusBadRequest, "badan permintaan bukan JSON aksi")
		return
	}
	ly, err := h.l.Aksi(r.Context(), h.pelaku(r), r.PathValue("id"), req)
	if err != nil {
		tulisGalat(w, err)
		return
	}
	tulisJSON(w, http.StatusOK, ly)
}

func (h *rute) pilihan(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	n, _ := strconv.Atoi(q.Get("indeks"))
	out, err := h.l.Pilihan(r.Context(), h.pelaku(r), r.PathValue("id"), services.PermintaanPilihan{
		Jenis: r.PathValue("jenis"), Konteks: q.Get("konteks"), N: n, Cari: q.Get("cari"), JenisCari: q.Get("jenisCari")})
	if err != nil {
		tulisGalat(w, err)
		return
	}
	tulisJSON(w, http.StatusOK, out)
}

func (h *rute) acuan(w http.ResponseWriter, r *http.Request) {
	out, err := h.l.Acuan(r.Context(), h.pelaku(r))
	if err != nil {
		tulisGalat(w, err)
		return
	}
	tulisJSON(w, http.StatusOK, out)
}
