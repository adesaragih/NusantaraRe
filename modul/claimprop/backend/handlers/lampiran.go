package handlers

// Rute lampiran klaim (GCNMSaveAttachments -> InsertDocument_Act):
//
//	GET  /api/claim-prop/kasus/{id}/lampiran   kategori (AttachCategory: master PROP + cacah) dan dokumen klaim
//	POST /api/claim-prop/kasus/{id}/lampiran   unggah, multipart satu atau lebih `berkas` + `kategori` bersama (Upload
//	                                           File baris) ATAU `kategoriBerkas` sebanyak berkas, urutan sama (Add attachment)
//	GET  /api/claim-prop/kasus/{id}/lampiran/{lid}/isi      File / View (GetBase64Attachment -> GetUrlGoogleStorage_Act)
//	GET  /api/claim-prop/kasus/{id}/lampiran/{lid}/office   View Office Online (pola NB Treaty In) - {url}
//	POST /api/claim-prop/kasus/{id}/lampiran/{lid}/hapus    Delete (pola NB Treaty In)
//	POST /api/claim-prop/kasus/{id}/lampiran/kategori       Change Category - JSON {kategori, ids}

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"mime"
	"net/http"

	"nusantarare/inti/backend/galat"
	"nusantarare/inti/backend/unggah"
	"nusantarare/modul/claimprop/backend/services"
)

// batasBadanLampiran - beberapa berkas sekali unggah, masing-masing paling besar `unggah.BatasUkuranUnggahan`.
const batasBadanLampiran = 4*unggah.BatasUkuranUnggahan + 1<<20

func (h *rute) lampiran(w http.ResponseWriter, r *http.Request) {
	out, err := h.l.DaftarLampiran(r.Context(), h.pelaku(r), r.PathValue("id"))
	if err != nil {
		tulisGalat(w, err)
		return
	}
	tulisJSON(w, http.StatusOK, out)
}

func (h *rute) unggahLampiran(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, batasBadanLampiran)
	if err := r.ParseMultipartForm(8 << 20); err != nil {
		var besar *http.MaxBytesError
		if errors.As(err, &besar) {
			galat.Tulis(w, http.StatusRequestEntityTooLarge, "The files are too large")
			return
		}
		galat.Tulis(w, http.StatusBadRequest, "No file attached")
		return
	}
	var berkas []services.BerkasUnggahan
	perBerkas := r.MultipartForm.Value["kategoriBerkas"]
	if len(perBerkas) > 0 && len(perBerkas) != len(r.MultipartForm.File["berkas"]) {
		galat.Tulis(w, http.StatusBadRequest, "kategoriBerkas must have one category per file")
		return
	}
	for i, kepala := range r.MultipartForm.File["berkas"] {
		f, err := kepala.Open()
		if err != nil {
			galat.Tulis(w, http.StatusBadRequest, "The file could not be read")
			return
		}
		isi, err := io.ReadAll(io.LimitReader(f, unggah.BatasUkuranUnggahan+1))
		_ = f.Close()
		if err != nil {
			galat.Tulis(w, http.StatusBadRequest, "The file could not be read")
			return
		}
		b := services.BerkasUnggahan{Nama: kepala.Filename, Isi: isi}
		if len(perBerkas) > 0 {
			b.Kategori = perBerkas[i]
		}
		berkas = append(berkas, b)
	}
	out, err := h.l.UnggahLampiran(r.Context(), h.pelaku(r), r.PathValue("id"), r.FormValue("kategori"), berkas)
	if err != nil {
		tulisGalat(w, err)
		return
	}
	tulisJSON(w, http.StatusOK, struct {
		Lampiran any `json:"lampiran"`
	}{out})
}

func (h *rute) unduhLampiran(w http.ResponseWriter, r *http.Request) {
	f, err := h.l.UnduhLampiran(r.Context(), h.pelaku(r), r.PathValue("id"), r.PathValue("lid"))
	if err != nil {
		tulisGalat(w, err)
		return
	}
	defer func() { _ = f.Isi.Close() }()
	w.Header().Set("Content-Type", f.Mime)
	w.Header().Set("Content-Disposition", mime.FormatMediaType("attachment", map[string]string{"filename": f.Nama}))
	w.Header().Set("X-Content-Type-Options", "nosniff")
	if _, err := io.Copy(w, f.Isi); err != nil {
		log.Printf("claimprop: mengirim lampiran %s: %v", r.PathValue("lid"), err)
	}
}

func (h *rute) officeLampiran(w http.ResponseWriter, r *http.Request) {
	url, err := h.l.TautanOfficeLampiran(r.Context(), h.pelaku(r), r.PathValue("id"), r.PathValue("lid"))
	if err != nil {
		tulisGalat(w, err)
		return
	}
	tulisJSON(w, http.StatusOK, struct {
		URL string `json:"url"`
	}{url})
}

func (h *rute) hapusLampiran(w http.ResponseWriter, r *http.Request) {
	if err := h.l.HapusLampiran(r.Context(), h.pelaku(r), r.PathValue("id"), r.PathValue("lid")); err != nil {
		tulisGalat(w, err)
		return
	}
	tulisJSON(w, http.StatusOK, struct {
		OK bool `json:"ok"`
	}{true})
}

func (h *rute) pindahKategoriLampiran(w http.ResponseWriter, r *http.Request) {
	var m struct {
		Kategori string   `json:"kategori"`
		IDs      []string `json:"ids"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, batasBadan)).Decode(&m); err != nil {
		galat.Tulis(w, http.StatusBadRequest, "badan permintaan tidak sah")
		return
	}
	if err := h.l.PindahKategoriLampiran(r.Context(), h.pelaku(r), r.PathValue("id"), m.Kategori, m.IDs); err != nil {
		tulisGalat(w, err)
		return
	}
	tulisJSON(w, http.StatusOK, struct {
		OK bool `json:"ok"`
	}{true})
}
