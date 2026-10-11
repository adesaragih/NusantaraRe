// Package handlers adalah pintu HTTP modul Claim Non Prop. Tugasnya hanya menerima permintaan, membaca pelaku, memanggil
// services, dan menerjemahkan galat ke kode HTTP. Nol aturan dagang, nol SQL.
//
// Rute - prefix `/api/claim-non-prop` (MODUL.md):
//
//	GET  /api/claim-non-prop/kasus?daftar=saya|workbasket|selesai&cari=   halaman awal (worklist / workbasket / selesai)
//	POST /api/claim-non-prop/kasus                                        Add Claim (Start -> Assignment2)
//	GET  /api/claim-non-prop/kasus/{id}                                   buka assignment (pra-proses)
//	POST /api/claim-non-prop/kasus/{id}/aksi                              aksi layar (refresh ber-activity, tombol, Choose)
//	GET  /api/claim-non-prop/kasus/{id}/pilihan/{jenis}?indeks=&cari=     isi pop-up / autocomplete
//	GET  /api/claim-non-prop/acuan                                        daftar pilihan bersama
//	GET  /api/claim-non-prop/hak                                          switch Teknik (anggota ReasKlaimTeknik)
//	GET  /api/claim-non-prop/berkas-polis?nopolis=                        berkas NB / EDM Treaty In (View polis)
//	GET|POST /api/claim-non-prop/kasus/{id}/lampiran...                   lampiran klaim pola Claim Prop (lampiran.go)
//
// Lampiran klaim mengikuti Claim Prop (perintah work owner 09-10-2026); korpus Non Prop sendiri tanpa section unggahan.
// ViewAttachmentNP (lampiran invoice pembayaran) dibaca lewat pilihan `lampiranBayar`.
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
	"nusantarare/modul/claimnonprop/backend/models"
	"nusantarare/modul/claimnonprop/backend/services"
)

// Prefix adalah awalan rute modul ini.
const Prefix = "/api/claim-non-prop"

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
	mux.HandleFunc("GET "+Prefix+"/berkas-polis", h.berkasPolis)
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
		log.Printf("claimnonprop: menulis jawaban: %v", err)
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
		errors.Is(err, services.ErrAkseptasiDiKomite),
		errors.Is(err, services.ErrTertunda):
		galat.Tulis(w, http.StatusConflict, err.Error())
	default:
		// Teks galat basis data (ORA-, nama skema) tidak dikirim ke klien.
		log.Printf("claimnonprop: %v", err)
		galat.Tulis(w, http.StatusInternalServerError, "gagal memproses permintaan Claim Non Prop")
	}
}

// berkasPolis - tombol View: modul dan ID berkas NB / EDM Treaty In untuk nomor polis.
func (h *rute) berkasPolis(w http.ResponseWriter, r *http.Request) {
	out, err := h.l.BerkasPolis(r.Context(), h.pelaku(r), r.URL.Query().Get("nopolis"))
	if err != nil {
		tulisGalat(w, err)
		return
	}
	tulisJSON(w, http.StatusOK, out)
}

// hak - switch Teknik halaman awal (anggota ReasKlaimTeknik).
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
	// `?lihat=1` - tampilan saja (pola Claim Prop; untuk tabel komite tahap 2).
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
	n, _ := strconv.Atoi(r.URL.Query().Get("indeks"))
	q := r.URL.Query()
	sm := models.SaringanMaster{Sumber: q.Get("sumber"), TreatyID: q.Get("treatyId"),
		ClassOfBusiness: q.Get("classOfBusiness"), ContractName: q.Get("contractName"), SOB: q.Get("sob"),
		Ceding: q.Get("ceding"), TreatyGroup: q.Get("treatyGroup"), TreatyYear: q.Get("treatyYear")}
	out, err := h.l.Pilihan(r.Context(), h.pelaku(r), r.PathValue("id"), r.PathValue("jenis"), n, q.Get("cari"), sm)
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
