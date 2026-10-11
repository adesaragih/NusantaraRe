package handlers

// Rute lampiran klaim induk di layar komite (perintah work owner 11-10-2026), hanya-baca:
//
//	GET /api/komite-claim-fac-in/kasus/{id}/lampiran               grid kategori berisi + dokumen klaim induk
//	GET /api/komite-claim-fac-in/kasus/{id}/lampiran/{lid}/isi     File / View (GetBase64Attachment)
//	GET /api/komite-claim-fac-in/kasus/{id}/lampiran/{lid}/office  View Office Online - {url}

import (
	"io"
	"log"
	"mime"
	"net/http"
)

func (h *rute) lampiran(w http.ResponseWriter, r *http.Request) {
	out, err := h.l.DaftarLampiran(r.Context(), h.pelaku(r), r.PathValue("id"))
	if err != nil {
		tulisGalat(w, err)
		return
	}
	tulisJSON(w, http.StatusOK, out)
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
		log.Printf("komiteclaimfacin: mengirim lampiran %s: %v", r.PathValue("lid"), err)
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
