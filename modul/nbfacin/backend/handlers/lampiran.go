package handlers

// Lampiran case NB FacIn (tiket 67) - grid "Attachment File" dan popup View File:
//
//	GET    /api/nbfacin/kasus/{caseId}/lampiran                     kategori + jumlah + bolehUnggah / bolehHapus
//	GET    /api/nbfacin/kasus/{caseId}/lampiran/berkas?kategori=    berkas satu kategori
//	POST   /api/nbfacin/kasus/{caseId}/lampiran/berkas              multipart (berkas, kategori) -> 201 {id}
//	GET    /api/nbfacin/kasus/{caseId}/lampiran/berkas/{id}/isi     isi berkas (proksi), attachment + nosniff
//	GET    /api/nbfacin/kasus/{caseId}/lampiran/berkas/{id}/office  View Office Online -> {url} (URL bertanda tangan https)
//	POST   /api/nbfacin/kasus/{caseId}/lampiran/berkas/{id}/ulangi  antre ulang unggahan -> 204
//	DELETE /api/nbfacin/kasus/{caseId}/lampiran/berkas/{id}         hapus FISIK (K-A) -> 204

import (
	"encoding/json"
	"errors"
	"io"
	"mime"
	"net/http"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/inti/backend/unggah"
	"nusantarare/modul/nbfacin/backend/services"
)

// batasPermintaanLampiran - SELURUH badan multipart (berkas 25 MiB inti + formulir); batas berkasnya sendiri ditegakkan
// saat menyalin.
const batasPermintaanLampiran = unggah.BatasUkuranUnggahan + galat.BatasFormulir

func daftarkanRuteLampiran(mux *http.ServeMux, svc *services.Service, stubPelaku bool) {
	mux.HandleFunc("GET /api/nbfacin/kasus/{caseId}/lampiran", tampilanLampiran(svc, stubPelaku))
	mux.HandleFunc("GET /api/nbfacin/kasus/{caseId}/lampiran/berkas", daftarBerkasLampiran(svc, stubPelaku))
	mux.HandleFunc("POST /api/nbfacin/kasus/{caseId}/lampiran/berkas", unggahLampiran(svc, stubPelaku))
	mux.HandleFunc("GET /api/nbfacin/kasus/{caseId}/lampiran/berkas/{id}/isi", unduhLampiran(svc, stubPelaku))
	mux.HandleFunc("GET /api/nbfacin/kasus/{caseId}/lampiran/berkas/{id}/office", officeLampiran(svc, stubPelaku))
	mux.HandleFunc("POST /api/nbfacin/kasus/{caseId}/lampiran/berkas/{id}/ulangi", ulangiLampiran(svc, stubPelaku))
	mux.HandleFunc("DELETE /api/nbfacin/kasus/{caseId}/lampiran/berkas/{id}", hapusLampiran(svc, stubPelaku))
}

type kategoriLampiranKabel struct {
	Kategori string `json:"kategori"`
	Jumlah   int    `json:"jumlah"`
}

type tampilanLampiranKabel struct {
	Kategori    []kategoriLampiranKabel `json:"kategori"`
	BolehUnggah bool                    `json:"bolehUnggah"`
	BolehHapus  bool                    `json:"bolehHapus"`
}

type barisLampiranKabel struct {
	ID       string `json:"id"`
	NamaFile string `json:"namaFile"`
	Catatan  string `json:"catatan"`
	Tanggal  string `json:"tanggal"`
	// Status - "" tersimpan, "menunggu", "gagal" (K-C).
	Status string `json:"status"`
}

// tulisGalatLampiran - galat khusus lampiran, sisanya tulisGalat bersama.
func tulisGalatLampiran(w http.ResponseWriter, err error) {
	var besar *http.MaxBytesError
	switch {
	case errors.Is(err, services.ErrMasukanLampiran), errors.Is(err, unggah.ErrBerkasKosong):
		galat.Tulis(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, unggah.ErrBerkasTerlaluBesar), errors.As(err, &besar):
		galat.Tulis(w, http.StatusRequestEntityTooLarge, "berkas melebihi batas ukuran 25 MiB")
	case errors.Is(err, services.ErrLampiranTidakAda):
		galat.Tulis(w, http.StatusNotFound, err.Error())
	case errors.Is(err, services.ErrLampiranBelumTersimpan), errors.Is(err, services.ErrLampiranSudahTersimpan),
		errors.Is(err, services.ErrLampiranTanpaBerkas), errors.Is(err, services.ErrBerkasSumberLampiranHilang),
		errors.Is(err, services.ErrLampiranOfficeTanpaStorage):
		galat.Tulis(w, http.StatusConflict, err.Error())
	case errors.Is(err, services.ErrKategoriLampiranTidakDikenal):
		galat.Tulis(w, http.StatusUnprocessableEntity, err.Error())
	case errors.Is(err, services.ErrLampiranTanpaDatabase), errors.Is(err, unggah.ErrUnggahanDirBelumDisetel):
		galat.Tulis(w, http.StatusServiceUnavailable, err.Error())
	default:
		tulisGalat(w, err)
	}
}

func tampilanLampiran(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t, err := svc.TampilanLampiranKasus(r.Context(), inti.PelakuDari(r, stubPelaku), r.PathValue("caseId"))
		if err != nil {
			tulisGalatLampiran(w, err)
			return
		}
		out := tampilanLampiranKabel{Kategori: make([]kategoriLampiranKabel, 0, len(t.Kategori)),
			BolehUnggah: t.BolehUnggah, BolehHapus: t.BolehHapus}
		for _, k := range t.Kategori {
			out.Kategori = append(out.Kategori, kategoriLampiranKabel{Kategori: k.Kategori, Jumlah: k.Jumlah})
		}
		galat.TulisJSON(w, out)
	}
}

func daftarBerkasLampiran(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		d, err := svc.DaftarBerkasLampiran(r.Context(), inti.PelakuDari(r, stubPelaku), r.PathValue("caseId"),
			r.URL.Query().Get("kategori"))
		if err != nil {
			tulisGalatLampiran(w, err)
			return
		}
		out := struct {
			Baris []barisLampiranKabel `json:"baris"`
		}{Baris: make([]barisLampiranKabel, 0, len(d))}
		for _, b := range d {
			out.Baris = append(out.Baris, barisLampiranKabel{ID: b.ID, NamaFile: b.NamaFile, Catatan: b.Catatan,
				Tanggal: b.Tanggal, Status: b.Status})
		}
		galat.TulisJSON(w, out)
	}
}

func unggahLampiran(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, batasPermintaanLampiran)
		if err := r.ParseMultipartForm(galat.BatasFormulir); err != nil {
			var besar *http.MaxBytesError
			if errors.As(err, &besar) {
				tulisGalatLampiran(w, err)
				return
			}
			galat.Tulis(w, http.StatusBadRequest, "permintaan harus multipart (berkas, kategori)")
			return
		}
		defer func() { _ = r.MultipartForm.RemoveAll() }()
		berkas, kepala, err := r.FormFile("berkas")
		if err != nil {
			galat.Tulis(w, http.StatusBadRequest, "Pilih berkas yang akan diunggah.")
			return
		}
		defer func() { _ = berkas.Close() }()
		id, err := svc.UnggahLampiran(r.Context(), inti.PelakuDari(r, stubPelaku), r.PathValue("caseId"),
			unggah.BerkasMasuk{NamaFile: kepala.Filename, Kategori: r.FormValue("kategori"), Isi: berkas})
		if err != nil {
			tulisGalatLampiran(w, err)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(struct {
			ID string `json:"id"`
		}{ID: id})
	}
}

// officeLampiran - GET …/office: 200 {url}. ⛔ URL tidak ditulis ke log; Cache-Control no-store (URL bertanda tangan).
func officeLampiran(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, err := svc.TautanOfficeLampiran(r.Context(), inti.PelakuDari(r, stubPelaku), r.PathValue("caseId"), r.PathValue("id"))
		if err != nil {
			tulisGalatLampiran(w, err)
			return
		}
		w.Header().Set("Cache-Control", "no-store")
		galat.TulisJSON(w, struct {
			URL string `json:"url"`
		}{u})
	}
}

// kepalaUnduhan - ⛔ `attachment` (isi datang dari pengunggah) + nosniff; nama dikodekan RFC 2231.
func kepalaUnduhan(w http.ResponseWriter, jenis, nama string) {
	if jenis == "" {
		jenis = unggah.MimeBawaan
	}
	w.Header().Set("Content-Type", jenis)
	disposisi := mime.FormatMediaType("attachment", map[string]string{"filename": nama})
	if disposisi == "" {
		disposisi = "attachment"
	}
	w.Header().Set("Content-Disposition", disposisi)
	w.Header().Set("X-Content-Type-Options", "nosniff")
}

func unduhLampiran(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		d, isi, err := svc.UnduhLampiran(r.Context(), inti.PelakuDari(r, stubPelaku), r.PathValue("caseId"), r.PathValue("id"))
		if err != nil {
			tulisGalatLampiran(w, err)
			return
		}
		defer func() { _ = isi.Close() }()
		kepalaUnduhan(w, unggah.MimeDariNamaFile(d.NamaFile), d.NamaFile)
		_, _ = io.Copy(w, isi)
	}
}

func ulangiLampiran(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := svc.UlangiLampiran(r.Context(), inti.PelakuDari(r, stubPelaku), r.PathValue("caseId"), r.PathValue("id")); err != nil {
			tulisGalatLampiran(w, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func hapusLampiran(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := svc.HapusLampiran(r.Context(), inti.PelakuDari(r, stubPelaku), r.PathValue("caseId"), r.PathValue("id")); err != nil {
			tulisGalatLampiran(w, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
