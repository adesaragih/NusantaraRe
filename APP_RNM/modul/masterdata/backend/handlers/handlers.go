// Package handlers memasang rute HTTP modul Master Data (`/api/masterdata`).
package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/modul/masterdata/backend/models"
	"nusantarare/modul/masterdata/backend/services"
)

// batasBadan - badan permintaan terbesar yang dibaca.
const batasBadan = 1 << 16

// DaftarkanRute memasang seluruh rute modul ini.
func DaftarkanRute(mux *http.ServeMux, svc *services.Service, stubPelaku bool) {
	mux.HandleFunc("GET /api/masterdata", daftarTabel())
	mux.HandleFunc("GET /api/masterdata/{tabel}", daftarBaris(svc))
	mux.HandleFunc("POST /api/masterdata/{tabel}", tambah(svc, stubPelaku))
	mux.HandleFunc("PUT /api/masterdata/{tabel}/{id}", ubah(svc, stubPelaku))
	mux.HandleFunc("PUT /api/masterdata/{tabel}/{id}/status", ubahStatus(svc, stubPelaku))
}

func tulisGalat(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, services.ErrMasukanMaster):
		galat.Tulis(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, inti.ErrTanpaIdentitas):
		galat.Tulis(w, http.StatusUnauthorized, err.Error())
	case errors.Is(err, services.ErrMasterTidakDikenal), errors.Is(err, services.ErrBarisTidakAda):
		galat.Tulis(w, http.StatusNotFound, err.Error())
	case errors.Is(err, services.ErrSudahAda):
		galat.Tulis(w, http.StatusConflict, err.Error())
	case errors.Is(err, services.ErrMasterTanpaDatabase):
		galat.Tulis(w, http.StatusServiceUnavailable, err.Error())
	default:
		log.Printf("masterdata: galat server: %v", err)
		galat.Tulis(w, http.StatusInternalServerError, "galat server")
	}
}

// kolomKabel / tabelKabel - GET /api/masterdata: metadata setiap master untuk layar generik.
type kolomKabel struct {
	Kunci   string `json:"kunci"`
	Kolom   string `json:"kolom"`
	Lebar   int    `json:"lebar"`
	Wajib   bool   `json:"wajib"`
	Turunan bool   `json:"turunan"`
}

type tabelKabel struct {
	Kunci      string       `json:"kunci"`
	Judul      string       `json:"judul"`
	IDOtomatis bool         `json:"idOtomatis"`
	Kolom      []kolomKabel `json:"kolom"`
}

func daftarTabel() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		hasil := make([]tabelKabel, 0, len(models.DaftarMaster))
		for _, t := range models.DaftarMaster {
			tk := tabelKabel{Kunci: t.Kunci, Judul: t.Judul, IDOtomatis: t.IDOtomatis, Kolom: []kolomKabel{}}
			for _, k := range t.Kolom {
				tk.Kolom = append(tk.Kolom, kolomKabel{Kunci: k.JSON, Kolom: k.Nama, Lebar: k.Lebar, Wajib: k.Wajib, Turunan: k.Turunan})
			}
			hasil = append(hasil, tk)
		}
		galat.TulisJSON(w, struct {
			Tabel []tabelKabel `json:"tabel"`
		}{hasil})
	}
}

// daftarBaris - GET /api/masterdata/{tabel}?q=&status=&halaman= -> {baris, total, halaman, ukuran}; setiap baris
// memuat kunci kolom + "aktif" (boolean).
func daftarBaris(svc *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		nomor := 1
		if v := q.Get("halaman"); v != "" {
			n, err := strconv.Atoi(v)
			if err != nil {
				galat.Tulis(w, http.StatusBadRequest, "halaman harus bilangan bulat")
				return
			}
			nomor = n
		}
		h, err := svc.Daftar(r.Context(), r.PathValue("tabel"), q.Get("q"), q.Get("status"), nomor)
		if err != nil {
			tulisGalat(w, err)
			return
		}
		baris := make([]map[string]any, 0, len(h.Baris))
		for i, b := range h.Baris {
			m := map[string]any{"aktif": h.Aktif[i]}
			for k, v := range b {
				m[k] = v
			}
			baris = append(baris, m)
		}
		galat.TulisJSON(w, struct {
			Baris   []map[string]any `json:"baris"`
			Total   int              `json:"total"`
			Halaman int              `json:"halaman"`
			Ukuran  int              `json:"ukuran"`
		}{baris, h.Total, h.Nomor, h.Ukuran})
	}
}

// bacaBadan - objek JSON bernilai teks saja (nilai bukan teks -> 400).
func bacaBadan(w http.ResponseWriter, r *http.Request) (map[string]string, bool) {
	var m map[string]string
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, batasBadan)).Decode(&m); err != nil || m == nil {
		galat.Tulis(w, http.StatusBadRequest, "badan permintaan harus objek JSON bernilai teks")
		return nil, false
	}
	return m, true
}

// tambah - POST /api/masterdata/{tabel} -> 201 {id}.
func tambah(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m, ok := bacaBadan(w, r)
		if !ok {
			return
		}
		id, err := svc.Tambah(r.Context(), inti.PelakuDari(r, stubPelaku), r.PathValue("tabel"), m)
		if err != nil {
			tulisGalat(w, err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(struct {
			ID string `json:"id"`
		}{id})
	}
}

// ubah - PUT /api/masterdata/{tabel}/{id} -> {id}.
func ubah(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m, ok := bacaBadan(w, r)
		if !ok {
			return
		}
		id := r.PathValue("id")
		if err := svc.Ubah(r.Context(), inti.PelakuDari(r, stubPelaku), r.PathValue("tabel"), id, m); err != nil {
			tulisGalat(w, err)
			return
		}
		galat.TulisJSON(w, struct {
			ID string `json:"id"`
		}{id})
	}
}

// ubahStatus - PUT /api/masterdata/{tabel}/{id}/status badan {aktif} -> {id, aktif}.
func ubahStatus(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var b struct {
			Aktif *bool `json:"aktif"`
		}
		dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, batasBadan))
		dec.DisallowUnknownFields()
		if err := dec.Decode(&b); err != nil || b.Aktif == nil {
			galat.Tulis(w, http.StatusBadRequest, `badan permintaan harus {"aktif": true|false}`)
			return
		}
		id := r.PathValue("id")
		if err := svc.UbahStatus(r.Context(), inti.PelakuDari(r, stubPelaku), r.PathValue("tabel"), id, *b.Aktif); err != nil {
			tulisGalat(w, err)
			return
		}
		galat.TulisJSON(w, struct {
			ID    string `json:"id"`
			Aktif bool   `json:"aktif"`
		}{id, *b.Aktif})
	}
}
