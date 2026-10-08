package handlers

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/inti/backend/unggah"
	"nusantarare/modul/treatyin/backend/services"
)

// daftarkanLampiran - panel Attachment (`Section/WorkAttachments.xml`).
//
//	GET  /kontrak/{id}/lampiran  Refresh (`GetMasterTreatyCategory_Act`) — panel
//	                             saja, tanpa membaca ulang form
//	POST /kontrak/{id}/lampiran  Attach modal `ASM Attach Content`
//	                             (`TreatySaveAttachment`): multipart, medan
//	                             `kategori` (CATEGORY_ID) dan `berkas` (satu atau
//	                             lebih)
func daftarkanLampiran(pasang func(string, rute)) {
	dasar := Prefix + "/kontrak/{id}/lampiran"
	pasang("GET "+dasar, func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		panel, err := l.BacaPanelLampiran(r.Context(), p, r.PathValue("id"))
		tulis(w, panel, err)
	})
	pasang("POST "+dasar, func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		r.Body = http.MaxBytesReader(w, r.Body, unggah.BatasUkuranUnggahan+galat.BatasFormulir)
		var besar *http.MaxBytesError
		switch err := r.ParseMultipartForm(galat.BatasFormulir); {
		case errors.As(err, &besar):
			galat.Tulis(w, http.StatusRequestEntityTooLarge, fmt.Sprintf("files exceed %d MB", unggah.BatasUkuranUnggahan>>20))
			return
		case err != nil:
			galat.Tulis(w, http.StatusBadRequest, "request body is not a valid multipart form")
			return
		}
		m := services.MasukanUnggahLampiran{IDKontrak: r.PathValue("id"), KodeKategori: r.FormValue("kategori")}
		for _, hdr := range r.MultipartForm.File["berkas"] {
			f, err := hdr.Open()
			if err != nil {
				galat.Tulis(w, http.StatusBadRequest, "a file part cannot be read")
				return
			}
			isi, err := io.ReadAll(io.LimitReader(f, unggah.BatasUkuranUnggahan+1))
			_ = f.Close()
			if err != nil {
				galat.Tulis(w, http.StatusBadRequest, "a file part cannot be read")
				return
			}
			m.Berkas = append(m.Berkas, services.BerkasUnggah{Nama: hdr.Filename, Isi: isi})
		}
		hasil, err := l.UnggahLampiran(r.Context(), p, m)
		tulis(w, hasil, err)
	})
}
