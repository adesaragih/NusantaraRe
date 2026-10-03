package handlers

// Rute lampiran (paket 8, PARITAS §6).
//
//	GET    /api/master-product-name-life/produk/{id}/lampiran                 `Refresh` b65270 / `View`
//	POST   /api/master-product-name-life/produk/{id}/lampiran                 `Add attachment` b64747 → `Attach` (multipart `berkas`)
//	POST   /api/master-product-name-life/produk/{id}/lampiran/{lid}/ulangi    kirim ulang (tiket 09)
//	GET    /api/master-product-name-life/produk/{id}/lampiran/{lid}/unduh     tautan nama berkas b68903
//	GET    /api/master-product-name-life/produk/{id}/lampiran/unduh-semua     `Download All` b67657 (zip)
//	GET    /api/master-product-name-life/produk/{id}/lampiran/{lid}/office    `View Office Online` b69291 - {url} bertanda tangan
//	DELETE /api/master-product-name-life/produk/{id}/lampiran/{lid}           `Delete` b69714

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/inti/backend/unggah"
	"nusantarare/modul/masterproductnamelife/backend/services"
)

func daftarkanLampiran(pasang func(string, rute)) {
	dasar := Prefix + "/produk/{id}/lampiran"
	pasang("GET "+dasar, func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		d, err := l.DaftarLampiran(r.Context(), p, r.PathValue("id"))
		tulisDaftar(w, d, err)
	})
	pasang("POST "+dasar, func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		r.Body = http.MaxBytesReader(w, r.Body, unggah.BatasUkuranUnggahan+galat.BatasFormulir)
		nama, isi := "", io.Reader(nil)
		var besar *http.MaxBytesError
		switch err := r.ParseMultipartForm(galat.BatasFormulir); {
		case errors.As(err, &besar):
			galat.Tulis(w, http.StatusRequestEntityTooLarge, fmt.Sprintf("file exceeds %d MB", unggah.BatasUkuranUnggahan>>20))
			return
		case err != nil && !errors.Is(err, http.ErrNotMultipart):
			galat.Tulis(w, http.StatusBadRequest, "request body is not a valid multipart form")
			return
		case err == nil:
			// Tanpa bagian `berkas` = "Tidak ada file yg diattach" (`ProductNameSaveAttachment` 1 b292).
			if f, hdr, err := r.FormFile("berkas"); err == nil {
				defer func() { _ = f.Close() }()
				nama, isi = hdr.Filename, f
			}
		}
		a, err := l.UnggahLampiran(r.Context(), p, r.PathValue("id"), nama, isi)
		tulis(w, a, err)
	})
	pasang("POST "+dasar+"/{lid}/ulangi", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		a, err := l.UlangiLampiran(r.Context(), p, r.PathValue("id"), r.PathValue("lid"))
		tulis(w, a, err)
	})
	pasang("GET "+dasar+"/{lid}/unduh", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		b, err := l.UnduhLampiran(r.Context(), p, r.PathValue("id"), r.PathValue("lid"))
		if jawabGalat(w, err) {
			return
		}
		defer func() { _ = b.Isi.Close() }()
		kirimBerkas(w, b.Nama, b.Mime)
		if _, err := io.Copy(w, b.Isi); err != nil {
			// Isi putus di tengah (penyimpanan): sambungan diputus - peramban melihat unduhan gagal, bukan berkas
			// terpotong berstatus 200. Rincian galat tidak dicetak (bisa memuat alamat jaringan).
			log.Printf("master product name life: attachment %s download was interrupted", r.PathValue("lid"))
			panic(http.ErrAbortHandler)
		}
	})
	pasang("GET "+dasar+"/unduh-semua", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		var buf bytes.Buffer
		if _, err := l.UnduhSemuaLampiran(r.Context(), p, r.PathValue("id"), &buf); jawabGalat(w, err) {
			return
		}
		kirimBerkas(w, "lampiran-"+r.PathValue("id")+".zip", "application/zip")
		_, _ = w.Write(buf.Bytes())
	})
	pasang("GET "+dasar+"/{lid}/office", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		u, err := l.LihatOffice(r.Context(), p, r.PathValue("id"), r.PathValue("lid"))
		tulis(w, map[string]string{"url": u}, err)
	})
	pasang("DELETE "+dasar+"/{lid}", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		err := l.HapusLampiran(r.Context(), p, r.PathValue("id"), r.PathValue("lid"))
		tulis(w, map[string]string{"id": r.PathValue("lid")}, err)
	})
}

// kirimBerkas menulis kepala jawaban unduhan (nama berkas UTF-8, RFC 6266).
func kirimBerkas(w http.ResponseWriter, nama, mime string) {
	w.Header().Set("Content-Type", mime)
	w.Header().Set("Content-Disposition", "attachment; filename*=UTF-8"+`''`+url.PathEscape(nama))
	w.WriteHeader(http.StatusOK)
}
