package dokumenpolis

import (
	"context"
	"errors"
	"io"
	"log"
	"mime"
	"net/http"
	"strings"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/galat"
	"nusantarare/inti/backend/penyimpanan"
	"nusantarare/inti/backend/unggah"
)

// SumberKasus - modul menerjemahkan `{id}` jalur dan pelaku menjadi Kasus (ada, boleh diubah); galatnya dijawab
// `galatModul` modul itu (identitas, kasus tidak ada, tanpa Oracle).
type SumberKasus func(ctx context.Context, p inti.Pelaku, id string) (Kasus, error)

// Pasang memasang rute lampiran "Reas" satu modul di bawah `dasar` (memuat `{id}` kasus, mis.
// `/api/nb-treaty-in/kasus/{id}/lampiran`). Kategori lewat kueri - NOTE memuat garis miring (`R/I SLIP`).
//
//	GET  <dasar>                           kategori + jumlah (`AttachmentGridReas`)
//	GET  <dasar>/dokumen?kategori=         dokumen satu kategori (`ReasViewAttachment`)
//	POST <dasar>/dokumen?kategori=         Upload File, multipart `berkas` (`InsertDocument_Act`)
//	GET  <dasar>/dokumen/{did}/isi         unduh / View (`DownloadDocumentPolis`)
//	GET  <dasar>/dokumen/{did}/office      URL View Office Online (`DownloadDocumentPolis` ViewOffice)
//	POST <dasar>/dokumen/{did}/hapus       Delete (`DeleteDocumentPolis_Act`)
func Pasang(mux *http.ServeMux, dasar string, layanan func() *Layanan, sumber SumberKasus, stub bool,
	galatModul func(http.ResponseWriter, error)) {
	pasang := func(pola string, f func(w http.ResponseWriter, r *http.Request, l *Layanan, k Kasus, akun string)) {
		mux.HandleFunc(pola, func(w http.ResponseWriter, r *http.Request) {
			p := inti.PelakuDari(r, stub)
			k, err := sumber(r.Context(), p, r.PathValue("id"))
			if err != nil {
				galatModul(w, err)
				return
			}
			f(w, r, layanan(), k, strings.TrimSpace(p.AkunID))
		})
	}
	jawab := func(w http.ResponseWriter, err error, apa string) bool {
		if err == nil {
			return false
		}
		if !jawabGalat(w, err, apa) {
			galatModul(w, err)
		}
		return true
	}
	pasang("GET "+dasar, func(w http.ResponseWriter, r *http.Request, l *Layanan, k Kasus, _ string) {
		d, err := l.Kategori(r.Context(), k)
		if jawab(w, err, "membaca kategori lampiran") {
			return
		}
		galat.TulisJSON(w, struct {
			Daftar    []Kategori `json:"daftar"`
			BolehUbah bool       `json:"bolehUbah"`
		}{d, k.BolehUbah})
	})
	pasang("GET "+dasar+"/dokumen", func(w http.ResponseWriter, r *http.Request, l *Layanan, k Kasus, _ string) {
		d, err := l.Daftar(r.Context(), k, r.URL.Query().Get("kategori"))
		if jawab(w, err, "membaca lampiran") {
			return
		}
		galat.TulisJSON(w, struct {
			Daftar []Dokumen `json:"daftar"`
		}{d})
	})
	pasang("POST "+dasar+"/dokumen", func(w http.ResponseWriter, r *http.Request, l *Layanan, k Kasus, akun string) {
		nama, isi, ok := bacaBerkas(w, r)
		if !ok {
			return
		}
		d, err := l.Unggah(r.Context(), k, akun, r.URL.Query().Get("kategori"), nama, isi)
		if jawab(w, err, "mengunggah lampiran") {
			return
		}
		galat.TulisJSON(w, d)
	})
	pasang("GET "+dasar+"/dokumen/{did}/isi", func(w http.ResponseWriter, r *http.Request, l *Layanan, k Kasus, akun string) {
		f, err := l.Unduh(r.Context(), k, akun, r.PathValue("did"))
		if jawab(w, err, "mengunduh lampiran") {
			return
		}
		defer func() { _ = f.Isi.Close() }()
		w.Header().Set("Content-Type", f.Mime)
		w.Header().Set("Content-Disposition", mime.FormatMediaType("attachment", map[string]string{"filename": f.Nama}))
		w.Header().Set("X-Content-Type-Options", "nosniff")
		if _, err := io.Copy(w, f.Isi); err != nil {
			log.Printf("dokumenpolis: mengirim lampiran %s: %v", r.PathValue("did"), err)
		}
	})
	pasang("GET "+dasar+"/dokumen/{did}/office", func(w http.ResponseWriter, r *http.Request, l *Layanan, k Kasus, akun string) {
		u, err := l.TautanOffice(r.Context(), k, akun, r.PathValue("did"))
		if jawab(w, err, "membuka View Office Online") {
			return
		}
		galat.TulisJSON(w, struct {
			URL string `json:"url"`
		}{u})
	})
	pasang("POST "+dasar+"/dokumen/{did}/hapus", func(w http.ResponseWriter, r *http.Request, l *Layanan, k Kasus, akun string) {
		if jawab(w, l.Hapus(r.Context(), k, akun, r.PathValue("did")), "menghapus lampiran") {
			return
		}
		galat.TulisJSON(w, struct {
			OK bool `json:"ok"`
		}{true})
	})
}

func pesan(err, penanda error) string { return strings.TrimPrefix(err.Error(), penanda.Error()+": ") }

// jawabGalat menjawab galat paket ini dan `inti/backend/penyimpanan`; false = bukan miliknya (dijawab modul).
func jawabGalat(w http.ResponseWriter, err error, apa string) bool {
	switch {
	case errors.Is(err, ErrTidakAda):
		galat.Tulis(w, http.StatusNotFound, "Attachment not found")
	case errors.Is(err, ErrDilarang):
		galat.Tulis(w, http.StatusForbidden, "This case is resolved; its attachments can no longer be changed")
	case errors.Is(err, ErrMasukanTidakSah):
		galat.Tulis(w, http.StatusUnprocessableEntity, pesan(err, ErrMasukanTidakSah))
	case errors.Is(err, penyimpanan.ErrBerkasDitolak):
		galat.Tulis(w, http.StatusUnprocessableEntity, "The file cannot be stored: "+pesan(err, penyimpanan.ErrBerkasDitolak))
	case errors.Is(err, penyimpanan.ErrObjekTidakAda), errors.Is(err, penyimpanan.ErrBerkasTidakDiStorage):
		galat.Tulis(w, http.StatusConflict, "The attachment file is not in storage")
	case errors.Is(err, penyimpanan.ErrStorageBelumSiap):
		log.Printf("dokumenpolis: %s: %v", apa, err)
		galat.Tulis(w, http.StatusServiceUnavailable, penyimpanan.PesanBelumSiap(err))
	case errors.Is(err, penyimpanan.ErrStorageGagal):
		log.Printf("dokumenpolis: %s: %v", apa, err)
		galat.Tulis(w, http.StatusBadGateway, "File storage failed: "+pesan(err, penyimpanan.ErrStorageGagal))
	case errors.Is(err, db.ErrTanpaOracle):
		galat.Tulis(w, http.StatusServiceUnavailable, "database is not configured")
	default:
		return false
	}
	return true
}

// bacaBerkas membaca satu berkas multipart `berkas`; galat = 400 / 413.
func bacaBerkas(w http.ResponseWriter, r *http.Request) (string, []byte, bool) {
	r.Body = http.MaxBytesReader(w, r.Body, unggah.BatasUkuranUnggahan+1<<20)
	if err := r.ParseMultipartForm(8 << 20); err != nil {
		var besar *http.MaxBytesError
		if errors.As(err, &besar) {
			galat.Tulis(w, http.StatusRequestEntityTooLarge, "The file is larger than 25 MB")
			return "", nil, false
		}
		galat.Tulis(w, http.StatusBadRequest, "No file attached")
		return "", nil, false
	}
	f, kepala, err := r.FormFile("berkas")
	if err != nil {
		galat.Tulis(w, http.StatusBadRequest, "No file attached")
		return "", nil, false
	}
	defer func() { _ = f.Close() }()
	isi, err := io.ReadAll(io.LimitReader(f, unggah.BatasUkuranUnggahan+1))
	if err != nil {
		galat.Tulis(w, http.StatusBadRequest, "The file could not be read")
		return "", nil, false
	}
	if len(isi) > unggah.BatasUkuranUnggahan {
		galat.Tulis(w, http.StatusRequestEntityTooLarge, "The file is larger than 25 MB")
		return "", nil, false
	}
	return kepala.Filename, isi, true
}
