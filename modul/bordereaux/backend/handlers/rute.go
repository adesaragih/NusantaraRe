// Package handlers memasang API modul Bordereaux - `/api/bordereaux` (keputusan work owner 04-10-2026).
//
//	GET  /api/bordereaux                       daftar (filter `InboxBordereaux_RD`), berhalaman
//	GET  /api/bordereaux/pilihan               Type, Business per Type, status, kolom grid per kombinasi, hak Input Data
//	GET  /api/bordereaux/berkas/{id}           buka satu berkas: header, detail, Summary, riwayat, hak
//	POST /api/bordereaux/unggah-csv            Upload CSV -> baris detail + Summary (tanpa menyimpan)
//	POST /api/bordereaux/simpan                Save (berkas baru atau ubah)
//	POST /api/bordereaux/berkas/{id}/submit    Submit / Approve / Reject
//	POST /api/bordereaux/berkas/{id}/hapus     Delete
//	GET  /api/bordereaux/chart                 chart daftar: jumlah berkas per Business x Type x Ceding
//	GET  /api/bordereaux/cedant?q=             autocomplete Cedant (popup Choose Master Treaty)
//	GET  /api/bordereaux/master-treaty?ceding= grid Master Treaty satu cedant
//	GET  /api/bordereaux/lama                  popup Copy Old Data - superadmin
//	POST /api/bordereaux/lama/salin            Process Copy {"ids": [...]} - superadmin
//
// Lampiran (`AttachmentsBdx`, `AttachmentDetailBdx`; penyimpanan bersama `inti/backend/penyimpanan`):
//
//	GET  /api/bordereaux/berkas/{id}/lampiran                               kategori + jumlah (`GetKategotyDocBdx`)
//	GET  /api/bordereaux/berkas/{id}/lampiran/{kategori}                    lampiran satu kategori (`getAttcachmentList`)
//	POST /api/bordereaux/berkas/{id}/lampiran/{kategori}                    Upload File, multipart `berkas` (`AttachDocBdx_Post`)
//	GET  /api/bordereaux/berkas/{id}/lampiran/{kategori}/{lid}/isi          unduh (`DownloadAttachmentBdx`)
//	GET  /api/bordereaux/berkas/{id}/lampiran/{kategori}/{lid}/office       URL View Office Online (`ViewOffice=true`)
//	POST /api/bordereaux/berkas/{id}/lampiran/{kategori}/{lid}/hapus        Delete (`DeleteAttachmentBdx`)
//
// Gerbang menu `bordereaux` dipasang `cmd/api` (403). Superadmin = pemegang menu Kelola User.
package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"mime"
	"net/http"
	"strconv"
	"strings"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/inti/backend/menu"
	"nusantarare/inti/backend/penyimpanan"
	"nusantarare/inti/backend/unggah"
	"nusantarare/modul/bordereaux/backend/models"
	"nusantarare/modul/bordereaux/backend/services"
)

// Prefix rute API modul ini (`MODUL.md`).
const Prefix = "/api/bordereaux"

const (
	batasBadanKecil = 1 << 16
	batasUnggah     = 16 << 20
	batasSimpan     = 32 << 20
)

// Router - mux berisi rute modul (uji).
func Router(svc *services.Service, stubPelaku bool) http.Handler {
	mux := http.NewServeMux()
	DaftarkanRute(mux, svc, stubPelaku)
	return mux
}

// RouterDengan - mux di atas layanan tertentu (uji dengan tiruan).
func RouterDengan(l *services.Layanan, adaDB, stubPelaku bool) http.Handler {
	mux := http.NewServeMux()
	daftarkan(mux, func() *services.Layanan { return l }, func() bool { return adaDB }, stubPelaku)
	return mux
}

// DaftarkanRute memasang rute ke mux bersama.
func DaftarkanRute(mux *http.ServeMux, svc *services.Service, stubPelaku bool) {
	daftarkan(mux, func() *services.Layanan { return services.LayananOracle(svc) }, svc.PunyaDatabase, stubPelaku)
}

type rute func(w http.ResponseWriter, r *http.Request, l *services.Layanan, a services.Aktor)

// KodeMenu - KODE menu modul ini (`M_LOGIN_GO_MENU.MENU_KODE`, sama dengan nama modul).
const KodeMenu = "bordereaux"

// aktorDari - pelaku sesi (atau stub) + superadmin = memegang menu Kelola User + hak menu PENUH (bukan View only).
func aktorDari(r *http.Request, stub bool) services.Aktor {
	p := inti.PelakuDari(r, stub)
	kode, _ := inti.AksesMenuDari(r.Context())
	return services.Aktor{AkunID: strings.TrimSpace(p.AkunID), Peran: p.Peran, Superadmin: inti.PunyaMenu(kode, menu.KodeKelolaUser),
		Penuh: menu.BolehUbah(r.Context(), KodeMenu)}
}

func daftarkan(mux *http.ServeMux, layanan func() *services.Layanan, adaDB func() bool, stub bool) {
	pasang := func(pola string, f rute) {
		mux.HandleFunc(pola, func(w http.ResponseWriter, r *http.Request) {
			if !adaDB() {
				galat.Tulis(w, http.StatusServiceUnavailable, "database is not configured")
				return
			}
			f(w, r, layanan(), aktorDari(r, stub))
		})
	}
	pasang("GET "+Prefix, func(w http.ResponseWriter, r *http.Request, l *services.Layanan, a services.Aktor) {
		q := r.URL.Query()
		f := models.Filter{BdxID: q.Get("id"), Type: q.Get("type"), Business: q.Get("business"), ReffNoSOA: q.Get("reffSoa"),
			ReffNoBDX: q.Get("reffBdx"), Ceding: q.Get("ceding"), Treaty: q.Get("treaty"), ReportStart: q.Get("start"),
			ReportEnd: q.Get("end"), Position: q.Get("position"), Status: q.Get("status")}
		h, err := l.Daftar(r.Context(), a, f, angka(r, "halaman"), angka(r, "ukuran"))
		if jawabGalat(w, err, "membaca daftar") {
			return
		}
		galat.TulisJSON(w, h)
	})
	pasang("GET "+Prefix+"/pilihan", func(w http.ResponseWriter, _ *http.Request, _ *services.Layanan, a services.Aktor) {
		galat.TulisJSON(w, services.DaftarPilihan(a))
	})
	pasang("GET "+Prefix+"/berkas/{id}", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, a services.Aktor) {
		rinci, err := l.Buka(r.Context(), a, r.PathValue("id"))
		if jawabGalat(w, err, "membuka berkas") {
			return
		}
		galat.TulisJSON(w, rinci)
	})
	pasang("POST "+Prefix+"/unggah-csv", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, _ services.Aktor) {
		var badan struct {
			Type     string `json:"type"`
			Business string `json:"business"`
			CSV      string `json:"csv"`
		}
		if !bacaBadan(w, r, &badan, batasUnggah) {
			return
		}
		p, err := l.UnggahCSV(r.Context(), badan.Type, badan.Business, badan.CSV)
		if jawabGalat(w, err, "membaca CSV") {
			return
		}
		galat.TulisJSON(w, p)
	})
	pasang("POST "+Prefix+"/simpan", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, a services.Aktor) {
		var badan services.PermintaanSimpan
		if !bacaBadan(w, r, &badan, batasSimpan) {
			return
		}
		id, err := l.Simpan(r.Context(), a, badan)
		if jawabGalat(w, err, "menyimpan") {
			return
		}
		galat.TulisJSON(w, struct {
			BdxID string `json:"bdxId"`
		}{id})
	})
	pasang("POST "+Prefix+"/berkas/{id}/submit", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, a services.Aktor) {
		var badan struct {
			Setuju   bool   `json:"setuju"`
			Komentar string `json:"komentar"`
		}
		if !bacaBadan(w, r, &badan, batasBadanKecil) {
			return
		}
		if jawabGalat(w, l.Submit(r.Context(), a, r.PathValue("id"), badan.Setuju, badan.Komentar), "submit") {
			return
		}
		galat.TulisJSON(w, struct {
			OK bool `json:"ok"`
		}{true})
	})
	pasang("POST "+Prefix+"/berkas/{id}/hapus", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, a services.Aktor) {
		if jawabGalat(w, l.Hapus(r.Context(), a, r.PathValue("id")), "menghapus") {
			return
		}
		galat.TulisJSON(w, struct {
			OK bool `json:"ok"`
		}{true})
	})
	pasang("GET "+Prefix+"/chart", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, _ services.Aktor) {
		d, err := l.Chart(r.Context())
		if jawabGalat(w, err, "membaca chart") {
			return
		}
		galat.TulisJSON(w, struct {
			Irisan []models.IrisanChart `json:"irisan"`
		}{d})
	})
	pasang("GET "+Prefix+"/cedant", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, _ services.Aktor) {
		d, err := l.CariCedant(r.Context(), r.URL.Query().Get("q"))
		if jawabGalat(w, err, "mencari cedant") {
			return
		}
		galat.TulisJSON(w, struct {
			Daftar []models.Cedant `json:"daftar"`
		}{d})
	})
	pasang("GET "+Prefix+"/master-treaty", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, _ services.Aktor) {
		d, err := l.CariTreaty(r.Context(), r.URL.Query().Get("ceding"))
		if jawabGalat(w, err, "membaca master treaty") {
			return
		}
		galat.TulisJSON(w, struct {
			Daftar []models.MasterTreaty `json:"daftar"`
		}{d})
	})
	pasang("GET "+Prefix+"/lama", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, a services.Aktor) {
		d, err := l.DaftarLama(r.Context(), a)
		if jawabGalat(w, err, "membaca data lama") {
			return
		}
		galat.TulisJSON(w, struct {
			Daftar []models.BerkasLama `json:"daftar"`
		}{d})
	})
	pasang("POST "+Prefix+"/lama/salin", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, a services.Aktor) {
		var badan struct {
			IDs []string `json:"ids"`
		}
		if !bacaBadan(w, r, &badan, batasBadanKecil) {
			return
		}
		j, err := l.SalinLama(r.Context(), a, badan.IDs)
		if jawabGalat(w, err, "menyalin data lama") {
			return
		}
		galat.TulisJSON(w, j)
	})
	daftarkanLampiran(pasang)
}

func daftarkanLampiran(pasang func(pola string, f rute)) {
	const dasar = Prefix + "/berkas/{id}/lampiran"
	pasang("GET "+dasar, func(w http.ResponseWriter, r *http.Request, l *services.Layanan, a services.Aktor) {
		d, err := l.KategoriLampiran(r.Context(), a, r.PathValue("id"))
		if jawabGalat(w, err, "membaca kategori lampiran") {
			return
		}
		galat.TulisJSON(w, struct {
			Daftar []models.KategoriLampiran `json:"daftar"`
		}{d})
	})
	pasang("GET "+dasar+"/{kategori}", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, a services.Aktor) {
		d, err := l.DaftarLampiran(r.Context(), a, r.PathValue("id"), r.PathValue("kategori"))
		if jawabGalat(w, err, "membaca lampiran") {
			return
		}
		galat.TulisJSON(w, struct {
			Daftar []models.Lampiran `json:"daftar"`
		}{d})
	})
	pasang("POST "+dasar+"/{kategori}", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, a services.Aktor) {
		nama, isi, ok := bacaBerkas(w, r)
		if !ok {
			return
		}
		b, err := l.UnggahLampiran(r.Context(), a, r.PathValue("id"), r.PathValue("kategori"), nama, isi)
		if jawabGalat(w, err, "mengunggah lampiran") {
			return
		}
		galat.TulisJSON(w, b)
	})
	pasang("GET "+dasar+"/{kategori}/{lid}/isi", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, a services.Aktor) {
		f, err := l.UnduhLampiran(r.Context(), a, r.PathValue("id"), r.PathValue("kategori"), r.PathValue("lid"))
		if jawabGalat(w, err, "mengunduh lampiran") {
			return
		}
		defer func() { _ = f.Isi.Close() }()
		w.Header().Set("Content-Type", f.Mime)
		w.Header().Set("Content-Disposition", mime.FormatMediaType("attachment", map[string]string{"filename": f.Nama}))
		w.Header().Set("X-Content-Type-Options", "nosniff")
		if _, err := io.Copy(w, f.Isi); err != nil {
			log.Printf("bordereaux: mengirim lampiran %s: %v", r.PathValue("lid"), err)
		}
	})
	pasang("GET "+dasar+"/{kategori}/{lid}/office", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, a services.Aktor) {
		u, err := l.TautanOffice(r.Context(), a, r.PathValue("id"), r.PathValue("kategori"), r.PathValue("lid"))
		if jawabGalat(w, err, "membuka View Office Online") {
			return
		}
		galat.TulisJSON(w, struct {
			URL string `json:"url"`
		}{u})
	})
	pasang("POST "+dasar+"/{kategori}/{lid}/hapus", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, a services.Aktor) {
		err := l.HapusLampiran(r.Context(), a, r.PathValue("id"), r.PathValue("kategori"), r.PathValue("lid"))
		if jawabGalat(w, err, "menghapus lampiran") {
			return
		}
		galat.TulisJSON(w, struct {
			OK bool `json:"ok"`
		}{true})
	})
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

func angka(r *http.Request, nama string) int {
	n, _ := strconv.Atoi(strings.TrimSpace(r.URL.Query().Get(nama)))
	return n
}

// bacaBadan membaca badan JSON; galat = 400 / 413.
func bacaBadan(w http.ResponseWriter, r *http.Request, ke any, batas int64) bool {
	d := json.NewDecoder(http.MaxBytesReader(w, r.Body, batas))
	d.DisallowUnknownFields()
	if err := d.Decode(ke); err != nil {
		var besar *http.MaxBytesError
		if errors.As(err, &besar) {
			galat.Tulis(w, http.StatusRequestEntityTooLarge, "request body is too large")
			return false
		}
		galat.Tulis(w, http.StatusBadRequest, "request body is not valid JSON for this form")
		return false
	}
	return true
}

func pesan(err, penanda error) string { return strings.TrimPrefix(err.Error(), penanda.Error()+": ") }

// jawabGalat menulis galat; true = sudah dijawab.
func jawabGalat(w http.ResponseWriter, err error, apa string) bool {
	var csv services.GalatCSV
	switch {
	case err == nil:
		return false
	case errors.As(err, &csv):
		galat.Tulis(w, http.StatusUnprocessableEntity, csv.Error())
	case errors.Is(err, services.ErrTidakAda):
		galat.Tulis(w, http.StatusNotFound, "Bordereaux not found")
	case errors.Is(err, services.ErrDilarang):
		galat.Tulis(w, http.StatusForbidden, pesan(err, services.ErrDilarang))
	case errors.Is(err, services.ErrSudahDiproses):
		galat.Tulis(w, http.StatusConflict, "This bordereaux has just been processed by someone else; reopen it")
	case errors.Is(err, services.ErrMasukanTidakSah):
		galat.Tulis(w, http.StatusUnprocessableEntity, pesan(err, services.ErrMasukanTidakSah))
	case errors.Is(err, services.ErrLampiranTidakAda):
		galat.Tulis(w, http.StatusNotFound, "Attachment not found")
	case errors.Is(err, penyimpanan.ErrBerkasDitolak):
		galat.Tulis(w, http.StatusUnprocessableEntity, "The file cannot be stored: "+pesan(err, penyimpanan.ErrBerkasDitolak))
	case errors.Is(err, penyimpanan.ErrObjekTidakAda), errors.Is(err, penyimpanan.ErrBerkasTidakDiStorage):
		galat.Tulis(w, http.StatusConflict, "The attachment file is not in storage")
	case errors.Is(err, penyimpanan.ErrStorageBelumSiap):
		log.Printf("bordereaux: %s: %v", apa, err)
		galat.Tulis(w, http.StatusServiceUnavailable, penyimpanan.PesanBelumSiap(err))
	case errors.Is(err, penyimpanan.ErrStorageGagal):
		log.Printf("bordereaux: %s: %v", apa, err)
		galat.Tulis(w, http.StatusBadGateway, "File storage failed: "+pesan(err, penyimpanan.ErrStorageGagal))
	case errors.Is(err, services.ErrBelumDimigrasi):
		log.Printf("bordereaux: %s: %v", apa, err)
		galat.Tulis(w, http.StatusServiceUnavailable, "Bordereaux is not migrated yet (run -migrate: inti 913, 890, 891, 998)")
	default:
		log.Printf("bordereaux: %s: %v", apa, err)
		galat.Tulis(w, http.StatusInternalServerError, "Bordereaux: "+apa+" gagal; rinciannya di log server")
	}
	return true
}
