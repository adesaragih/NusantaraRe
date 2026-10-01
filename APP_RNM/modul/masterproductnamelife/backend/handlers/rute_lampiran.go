package handlers

// Rute lampiran (paket 8, PARITAS §6).
//
//	GET    /api/master-product-name-life/produk/{id}/lampiran                 `Refresh` b65223 / `View`
//	POST   /api/master-product-name-life/produk/{id}/lampiran                 `Add attachment` b64698 → `Attach` (multipart `berkas`)
//	POST   /api/master-product-name-life/produk/{id}/lampiran/{lid}/ulangi    kirim ulang (tiket 09)
//	GET    /api/master-product-name-life/produk/{id}/lampiran/{lid}/unduh     tautan nama berkas b68857
//	GET    /api/master-product-name-life/produk/{id}/lampiran/unduh-semua     `Download All` b67619 (zip)
//	GET    /api/master-product-name-life/produk/{id}/lampiran/{lid}/office    `View Office Online` b69247 - 503 stub
//	DELETE /api/master-product-name-life/produk/{id}/lampiran/{lid}           `Delete` b69663

import (
	"bytes"
	"io"
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
		if err := r.ParseMultipartForm(galat.BatasFormulir); err == nil {
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
		_, _ = io.Copy(w, b.Isi)
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
		jawabGalat(w, l.LihatOffice(r.Context(), p, r.PathValue("id"), r.PathValue("lid")))
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
