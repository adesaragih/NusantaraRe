// Package handlers memasang API modul R/I Comm Life - `/api/ri-comm-life` (keputusan work owner 06-10-2026 butir 9).
// Nol aturan dagang di sini; nol impor repository.
//
//	GET    /api/ri-comm-life?id=&usedby=&urut=&arah=&halaman=   grid `BrowseRICommSummary` (50 per halaman)
//	GET    /api/ri-comm-life/{id}                               satu ringkasan + jumlah rincian (dialog Delete)
//	POST   /api/ri-comm-life                                    Save - Add (`AddToListSummary_Act`)
//	PUT    /api/ri-comm-life/{id}                               Save - Edit (`EditListSummary_DT`)
//	DELETE /api/ri-comm-life/{id}                               Delete (`DeleteSummaryDetail`) beserta rinciannya
//	GET    /api/ri-comm-life/{id}/detail?halaman=               R/I COMM DETAIL (`InboxRIComm`, 50 per halaman)
//	POST   /api/ri-comm-life/{id}/detail                        R/I COMM DETAIL Save - tambah baris (`AddToList_Act`)
//	PUT    /api/ri-comm-life/{id}/detail/{detailId}             R/I COMM DETAIL Save - ubah baris (`EditList_DT`)
//	POST   /api/ri-comm-life/unggah/pratinjau                   View Upload (`ViewCSVResult_RIComm`) - tanpa menulis
//	POST   /api/ri-comm-life/unggah                             Simpan Upload (`SubmitRIComm_Act`)
//
// Gerbang menu `ricommlife` dan gerbang tulis View only dipasang `cmd/api` (HakLihat tanpa pola bebas: View Upload pun
// tertutup bagi View only); layanan menolak lagi lewat Aktor.Penuh.
package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/inti/backend/menu"
	"nusantarare/modul/ricommlife/backend/models"
	"nusantarare/modul/ricommlife/backend/services"
)

// Prefix rute API modul ini (`MODUL.md`).
const Prefix = "/api/ri-comm-life"

// KodeMenu - KODE menu modul ini (`M_LOGIN_GO_MENU.MENU_KODE`, sama dengan nama modul).
const KodeMenu = "ricommlife"

// Batas badan: form kecil; unggahan = teks CSV (models.MaksBytesCSV) dibungkus JSON (pelolosan bisa menggandakan).
const (
	batasBadan       = 16 << 10
	batasBadanUnggah = 2*models.MaksBytesCSV + 4<<10
)

// Router - mux berisi rute modul (Oracle).
func Router(svc *services.Service, stubPelaku bool) http.Handler {
	mux := http.NewServeMux()
	DaftarkanRute(mux, svc, stubPelaku)
	return mux
}

// RouterDengan - mux di atas layanan tertentu (uji dengan tiruan).
func RouterDengan(l *services.Layanan, adaDB, stubPelaku bool) http.Handler {
	mux := http.NewServeMux()
	daftarkan(mux, func() *services.Layanan { return l }, func() bool { return adaDB }, stubPelaku)
	return mux
}

// DaftarkanRute memasang rute ke mux bersama.
func DaftarkanRute(mux *http.ServeMux, svc *services.Service, stubPelaku bool) {
	daftarkan(mux, func() *services.Layanan { return services.LayananOracle(svc) }, svc.PunyaDatabase, stubPelaku)
}

type rute func(w http.ResponseWriter, r *http.Request, l *services.Layanan, a services.Aktor)

// aktorDari - pelaku sesi (atau stub) dan hak menunya (Full = bukan View only).
func aktorDari(r *http.Request, stub bool) services.Aktor {
	p := inti.PelakuDari(r, stub)
	return services.Aktor{AkunID: strings.TrimSpace(p.AkunID), Penuh: menu.BolehUbah(r.Context(), KodeMenu)}
}

func halaman(r *http.Request) int {
	n, err := strconv.Atoi(strings.TrimSpace(r.URL.Query().Get("halaman")))
	if err != nil || n < 1 {
		return 1
	}
	return n
}

func daftarkan(mux *http.ServeMux, layanan func() *services.Layanan, adaDB func() bool, stub bool) {
	pasang := func(pola string, f rute) {
		mux.HandleFunc(pola, func(w http.ResponseWriter, r *http.Request) {
			if !adaDB() {
				galat.Tulis(w, http.StatusServiceUnavailable, "database is not configured")
				return
			}
			f(w, r, layanan(), aktorDari(r, stub))
		})
	}
	pasang("GET "+Prefix, func(w http.ResponseWriter, r *http.Request, l *services.Layanan, _ services.Aktor) {
		q := r.URL.Query()
		d, err := l.Daftar(r.Context(), models.Saringan{ID: q.Get("id"), UsedBy: q.Get("usedby"), Urut: q.Get("urut"),
			Turun: strings.EqualFold(q.Get("arah"), "desc"), Halaman: halaman(r)})
		if jawabGalat(w, err, "membaca daftar") {
			return
		}
		galat.TulisJSON(w, d)
	})
	pasang("GET "+Prefix+"/{id}", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, _ services.Aktor) {
		j, err := l.Buka(r.Context(), r.PathValue("id"))
		if jawabGalat(w, err, "membaca") {
			return
		}
		galat.TulisJSON(w, j)
	})
	simpan := func(w http.ResponseWriter, r *http.Request, l *services.Layanan, a services.Aktor, id string) {
		var isi models.Isian
		if !bacaBadan(w, r, &isi, batasBadan) {
			return
		}
		hasil, err := l.Simpan(r.Context(), a, id, isi)
		if jawabGalat(w, err, "menyimpan") {
			return
		}
		galat.TulisJSON(w, hasil)
	}
	pasang("POST "+Prefix, func(w http.ResponseWriter, r *http.Request, l *services.Layanan, a services.Aktor) {
		simpan(w, r, l, a, "")
	})
	pasang("PUT "+Prefix+"/{id}", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, a services.Aktor) {
		simpan(w, r, l, a, strings.TrimSpace(r.PathValue("id")))
	})
	pasang("DELETE "+Prefix+"/{id}", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, a services.Aktor) {
		h, err := l.Hapus(r.Context(), a, r.PathValue("id"))
		if jawabGalat(w, err, "menghapus") {
			return
		}
		galat.TulisJSON(w, h)
	})
	pasang("GET "+Prefix+"/{id}/detail", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, _ services.Aktor) {
		d, err := l.DaftarKomisi(r.Context(), r.PathValue("id"), halaman(r))
		if jawabGalat(w, err, "membaca detail") {
			return
		}
		galat.TulisJSON(w, d)
	})
	simpanDetail := func(w http.ResponseWriter, r *http.Request, l *services.Layanan, a services.Aktor, idDetail string) {
		var isi models.IsianKomisi
		if !bacaBadan(w, r, &isi, batasBadan) {
			return
		}
		hasil, err := l.SimpanKomisi(r.Context(), a, r.PathValue("id"), idDetail, isi)
		if jawabGalat(w, err, "menyimpan detail") {
			return
		}
		galat.TulisJSON(w, hasil)
	}
	pasang("POST "+Prefix+"/{id}/detail", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, a services.Aktor) {
		simpanDetail(w, r, l, a, "")
	})
	pasang("PUT "+Prefix+"/{id}/detail/{detailId}", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, a services.Aktor) {
		if strings.TrimSpace(r.PathValue("detailId")) == "" {
			galat.Tulis(w, http.StatusNotFound, "R/I comm detail row not found in this R/I comm")
			return
		}
		simpanDetail(w, r, l, a, r.PathValue("detailId"))
	})
	pasang("POST "+Prefix+"/unggah/pratinjau", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, a services.Aktor) {
		var p services.PermintaanUnggah
		if !bacaBadan(w, r, &p, batasBadanUnggah) {
			return
		}
		h, err := l.Pratinjau(r.Context(), a, p)
		if jawabGalat(w, err, "membaca unggahan") {
			return
		}
		galat.TulisJSON(w, h)
	})
	pasang("POST "+Prefix+"/unggah", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, a services.Aktor) {
		var p services.PermintaanUnggah
		if !bacaBadan(w, r, &p, batasBadanUnggah) {
			return
		}
		h, err := l.SimpanUnggah(r.Context(), a, p)
		if jawabGalat(w, err, "menyimpan unggahan") {
			return
		}
		galat.TulisJSON(w, h)
	})
}

// bacaBadan membaca badan JSON; galat = 400 / 413.
func bacaBadan(w http.ResponseWriter, r *http.Request, ke any, batas int64) bool {
	d := json.NewDecoder(http.MaxBytesReader(w, r.Body, batas))
	d.DisallowUnknownFields()
	if err := d.Decode(ke); err != nil {
		var besar *http.MaxBytesError
		if errors.As(err, &besar) {
			galat.Tulis(w, http.StatusRequestEntityTooLarge, "request body is too large")
			return false
		}
		galat.Tulis(w, http.StatusBadRequest, "request body is not valid JSON for this form")
		return false
	}
	return true
}

func pesan(err, penanda error) string { return strings.TrimPrefix(err.Error(), penanda.Error()+": ") }

// jawabGalat menulis galat; true = sudah dijawab.
func jawabGalat(w http.ResponseWriter, err error, apa string) bool {
	switch {
	case err == nil:
		return false
	case errors.Is(err, services.ErrTidakAda):
		galat.Tulis(w, http.StatusNotFound, "R/I comm summary not found")
	case errors.Is(err, services.ErrKomisiTidakAda):
		galat.Tulis(w, http.StatusNotFound, "R/I comm detail row not found in this R/I comm")
	case errors.Is(err, services.ErrDilarang):
		galat.Tulis(w, http.StatusForbidden, pesan(err, services.ErrDilarang))
	case errors.Is(err, services.ErrMasukanTidakSah):
		galat.Tulis(w, http.StatusUnprocessableEntity, pesan(err, services.ErrMasukanTidakSah))
	case errors.Is(err, services.ErrBelumAda):
		log.Printf("ricommlife: %s: %v", apa, err)
		galat.Tulis(w, http.StatusServiceUnavailable,
			"R/I Comm Life: table or sequence (M_RICOMM_LIFE_SUMMARY, M_RICOMM_LIFE, M_SITE_DATABASE, M_RICOMM_LIFE_SUMMARY_SEQ, M_RICOMM_LIFE_SEQ) is not in this schema")
	default:
		log.Printf("ricommlife: %s: %v", apa, err)
		galat.Tulis(w, http.StatusInternalServerError, "R/I Comm Life: "+apa+" gagal; rinciannya di log server")
	}
	return true
}
