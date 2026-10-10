package handlers

import (
	"net/http"
	"strconv"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/treatyin/backend/services"
)

// daftarkanUnduhSemuaLampiran - tombol `Download All` panel Attachment
// (`DownloadAll_Act`): seluruh lampiran kontrak dalam `AllDocuments.zip`.
// Layar memakai fetch beridentitas, seperti tautan nama berkas.
//
//	GET /kontrak/{id}/lampiran/unduh-semua
func daftarkanUnduhSemuaLampiran(pasang func(string, rute)) {
	pasang("GET "+Prefix+"/kontrak/{id}/lampiran/unduh-semua", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		isi, err := l.UnduhSemuaLampiran(r.Context(), p, r.PathValue("id"))
		if jawabGalat(w, err) {
			return
		}
		w.Header().Set("Content-Type", "application/zip")
		w.Header().Set("Content-Disposition", "attachment; filename="+services.NamaUnduhSemuaLampiran)
		w.Header().Set("Content-Length", strconv.Itoa(len(isi)))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(isi)
	})
}
