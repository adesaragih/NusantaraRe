package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

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
//
// Modal `View File` (`ShowAttachmentTreaty`):
//
//	GET  …/lampiran/{lid}/isi                tautan nama berkas: isi berkas
//	                                         dialirkan (DownloadAttachmentTreaty)
//	GET  …/lampiran/{lid}/tautan[?office=1]  DownloadAttachmentTreaty → {"url"}
//	POST …/lampiran/{lid}/hapus              Delete_act
//	POST …/lampiran/kategori                 ChangeDokument_Act("Save"):
//	                                         {"perubahan":[{"id","kategori"}]}
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
	// Tautan nama berkas — isi berkas DIALIRKAN (layar memakai fetch
	// beridentitas, bukan membuka URL; `unduhdokumen.test.ts`).
	pasang("GET "+dasar+"/{lid}/isi", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		b, err := l.IsiLampiran(r.Context(), p, r.PathValue("id"), r.PathValue("lid"))
		if jawabGalat(w, err) {
			return
		}
		defer func() { _ = b.Isi.Close() }()
		w.Header().Set("Content-Type", b.Mime)
		w.Header().Set("Content-Disposition", "attachment; filename*=UTF-8''"+url.PathEscape(b.Nama))
		w.WriteHeader(http.StatusOK)
		if _, err := io.Copy(w, b.Isi); err != nil {
			// Isi putus di tengah: sambungan diputus — peramban melihat
			// unduhan gagal, bukan berkas terpotong berstatus 200.
			log.Printf("treaty in: attachment %s download was interrupted", r.PathValue("lid"))
			panic(http.ErrAbortHandler)
		}
	})
	pasang("GET "+dasar+"/{lid}/tautan", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		u, err := l.TautanLampiran(r.Context(), p, r.PathValue("id"), r.PathValue("lid"), r.URL.Query().Get("office") == "1")
		tulis(w, map[string]string{"url": u}, err)
	})
	pasang("POST "+dasar+"/{lid}/hapus", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		panel, err := l.HapusLampiran(r.Context(), p, r.PathValue("id"), r.PathValue("lid"))
		tulis(w, panel, err)
	})
	pasang("POST "+dasar+"/kategori", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		var m struct {
			Perubahan []services.MasukanUbahKategori `json:"perubahan"`
		}
		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			galat.Tulis(w, http.StatusBadRequest, "request body must be valid JSON")
			return
		}
		panel, err := l.UbahKategoriLampiran(r.Context(), p, r.PathValue("id"), m.Perubahan)
		tulis(w, panel, err)
	})
}
