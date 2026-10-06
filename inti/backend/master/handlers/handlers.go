// Package handlers memasang rute HTTP SATU master di prefix modul pemiliknya (`/api/master-<nama>`, pemecahan
// Master Data jadi delapan modul 04-10-2026). Kontrak: `modul/masterprovince/docs/issues/03-pemecahan-delapan-modul.md`.
package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/inti/backend/master/models"
	"nusantarare/inti/backend/master/services"
)

// batasBadan - badan permintaan terbesar yang dibaca.
const batasBadan = 1 << 16

// Pasang memasang rute master `kunci` di bawah `prefix`. Pemilik rute (gerbang menu `cmd/api`) = modul yang
// memanggilnya, jadi aksesnya per menu modul itu.
func Pasang(mux *http.ServeMux, prefix, kunci string, svc *services.Service, stubPelaku bool) {
	mux.HandleFunc("GET "+prefix+"/meta", meta(kunci))
	mux.HandleFunc("GET "+prefix, daftarBaris(svc, kunci))
	mux.HandleFunc("POST "+prefix, tambah(svc, kunci, stubPelaku))
	mux.HandleFunc("PUT "+prefix+"/{id}", ubah(svc, kunci, stubPelaku))
	mux.HandleFunc("PUT "+prefix+"/{id}/status", ubahStatus(svc, kunci, stubPelaku))
	mux.HandleFunc("GET "+prefix+"/rujukan/{kolom}", rujukan(svc, kunci))
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
		log.Printf("master: galat server: %v", err)
		galat.Tulis(w, http.StatusInternalServerError, "galat server")
	}
}

// kolomKabel / saranKabel / metaKabel - GET <prefix>/meta: kolom master untuk layar generik.
type kolomKabel struct {
	Kunci   string `json:"kunci"`
	Kolom   string `json:"kolom"`
	Lebar   int    `json:"lebar"`
	Wajib   bool   `json:"wajib"`
	Turunan bool   `json:"turunan"`
}

type saranKabel struct {
	Kunci string `json:"kunci"`
	Judul string `json:"judul"`
	Nilai string `json:"nilai"`
	Nama  string `json:"nama"`
}

type metaKabel struct {
	Kunci      string       `json:"kunci"`
	Judul      string       `json:"judul"`
	IDOtomatis bool         `json:"idOtomatis"`
	Kolom      []kolomKabel `json:"kolom"`
	Rujukan    []saranKabel `json:"rujukan"`
}

// meta - kolom data lalu keempat kolom jejak ubah (turunan, baca-saja; MD-7), dan kolom berujukan master.
func meta(kunci string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		t, ada := models.CariMaster(kunci)
		if !ada {
			tulisGalat(w, services.ErrMasterTidakDikenal)
			return
		}
		m := metaKabel{Kunci: t.Kunci, Judul: t.Judul, IDOtomatis: t.IDOtomatis, Kolom: []kolomKabel{}, Rujukan: []saranKabel{}}
		for _, k := range t.SeluruhKolom() {
			m.Kolom = append(m.Kolom, kolomKabel{Kunci: k.JSON, Kolom: k.Nama, Lebar: k.Lebar, Wajib: k.Wajib, Turunan: k.Turunan})
		}
		for _, s := range t.Saran() {
			m.Rujukan = append(m.Rujukan, saranKabel{Kunci: s.Kunci, Judul: s.Judul, Nilai: s.Nilai, Nama: s.Nama})
		}
		galat.TulisJSON(w, m)
	}
}

// halaman - nomor halaman `?halaman=` (bawaan 1); false = sudah dijawab 400.
func halaman(w http.ResponseWriter, r *http.Request) (int, bool) {
	v := r.URL.Query().Get("halaman")
	if v == "" {
		return 1, true
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		galat.Tulis(w, http.StatusBadRequest, "halaman harus bilangan bulat")
		return 0, false
	}
	return n, true
}

// tulisHalaman - {baris, total, halaman, ukuran}; setiap baris memuat kunci kolom + "aktif" (boolean).
func tulisHalaman(w http.ResponseWriter, h models.Halaman) {
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

// daftarBaris - GET <prefix>?q=&status=&halaman=.
func daftarBaris(svc *services.Service, kunci string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nomor, ok := halaman(w, r)
		if !ok {
			return
		}
		q := r.URL.Query()
		h, err := svc.Daftar(r.Context(), kunci, q.Get("q"), q.Get("status"), nomor)
		if err != nil {
			tulisGalat(w, err)
			return
		}
		tulisHalaman(w, h)
	}
}

// rujukan - GET <prefix>/rujukan/{kolom}?q=&halaman=: baris master yang dirujuk kolom itu, aktif saja.
func rujukan(svc *services.Service, kunci string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nomor, ok := halaman(w, r)
		if !ok {
			return
		}
		h, err := svc.Rujukan(r.Context(), kunci, r.PathValue("kolom"), r.URL.Query().Get("q"), nomor)
		if err != nil {
			tulisGalat(w, err)
			return
		}
		tulisHalaman(w, h)
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

// tambah - POST <prefix> -> 201 {id}.
func tambah(svc *services.Service, kunci string, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m, ok := bacaBadan(w, r)
		if !ok {
			return
		}
		id, err := svc.Tambah(r.Context(), inti.PelakuDari(r, stubPelaku), kunci, m)
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

// ubah - PUT <prefix>/{id} -> {id}.
func ubah(svc *services.Service, kunci string, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m, ok := bacaBadan(w, r)
		if !ok {
			return
		}
		id := r.PathValue("id")
		if err := svc.Ubah(r.Context(), inti.PelakuDari(r, stubPelaku), kunci, id, m); err != nil {
			tulisGalat(w, err)
			return
		}
		galat.TulisJSON(w, struct {
			ID string `json:"id"`
		}{id})
	}
}

// ubahStatus - PUT <prefix>/{id}/status badan {aktif} -> {id, aktif}.
func ubahStatus(svc *services.Service, kunci string, stubPelaku bool) http.HandlerFunc {
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
		if err := svc.UbahStatus(r.Context(), inti.PelakuDari(r, stubPelaku), kunci, id, *b.Aktif); err != nil {
			tulisGalat(w, err)
			return
		}
		galat.TulisJSON(w, struct {
			ID    string `json:"id"`
			Aktif bool   `json:"aktif"`
		}{id, *b.Aktif})
	}
}
