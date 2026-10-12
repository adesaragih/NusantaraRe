package handlers

// Pintu HTTP dokumen pendukung - butir be.
//
//	POST   /api/klaim-life/{id}/peserta/{pesertaId}/dokumen   (multipart)
//	GET    /api/dokumen/{dokId}/isi
//	DELETE /api/klaim-life/{id}/dokumen/{dokId}
//
// Padanan tiga tombol `Section/DocumentLife.xml`: `Add attachment` b1245,
// `View Office Online` b3502 (dan tautan baris), `Delete` b4288.
//
// ⛔ Jalur unduh TIDAK menyebut klaim, dan itu bukan kelalaian: `URLPUBLIC`
// di sistem lama dicari dengan `imageid` saja (`GetLinkStorage_SQL.xml` b91),
// dan nilai kolom itulah yang tersimpan. Batas klaimnya tetap ditegakkan -
// services mencari klaim pemilik dokumennya lebih dulu.
//
// Dibaca sesudah: diagnosa.go dan services/unggahan.go.

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/inti/backend/kontrak"
	"nusantarare/inti/backend/unggah"
	"nusantarare/modul/claimlife/backend/models"
	"nusantarare/modul/claimlife/backend/services"
)

// unggahDokumen melayani POST multipart - `Add attachment` b1245.
func unggahDokumen(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat.Tulis(w, http.StatusServiceUnavailable, "The database is not configured.")
			return
		}
		// ⛔ `ReadForm` DIHINDARI: ia menulis seluruh berkas ke berkas
		// sementara lebih dulu. `MultipartReader` mengalirkannya, sehingga
		// batas ukuran berlaku SEBELUM 25 MiB mendarat dua kali.
		if err := r.ParseMultipartForm(galat.BatasFormulir); err != nil {
			galat.Tulis(w, http.StatusBadRequest, "The request is not a valid multipart upload.")
			return
		}
		berkas, kepala, err := r.FormFile("berkas")
		if err != nil {
			galat.Tulis(w, http.StatusBadRequest, "The request has no file part.")
			return
		}
		defer func() { _ = berkas.Close() }()

		dok, err := svc.Dokumen().
			DenganKategori(services.KategoriWajibOracle(svc)).
			Unggah(r.Context(), inti.PelakuDari(r, stubPelaku),
				r.PathValue("id"), r.PathValue("pesertaId"),
				unggah.BerkasMasuk{
					NamaFile: kepala.Filename,
					Mime:     kepala.Header.Get("Content-Type"),
					Kategori: r.FormValue("kategori"),
					Isi:      berkas,
				}, time.Now())
		if jawabGalatDokumen(w, err) {
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(dok)
	}
}

// isiDokumen melayani GET /api/dokumen/{dokId}/isi.
//
// ⛔ Identitas diperiksa SEBELUM basis data (GILIRAN-12, temuan /code-review):
// pranala tanpa header adalah cacat yang pernah membuat setiap unduhan 401,
// dan urutan ini membuat penolakannya dapat diuji pada handler yang
// sebenarnya, tanpa Oracle (`dokumen_identitas_test.go`).
func isiDokumen(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pelaku := inti.PelakuDari(r, stubPelaku)
		if jawabGalatDokumen(w, inti.WajibIdentitas(pelaku)) {
			return
		}
		if !svc.PunyaDatabase() {
			galat.Tulis(w, http.StatusServiceUnavailable, "The database is not configured.")
			return
		}
		id, ok := pengenalDokumen(w, r)
		if !ok {
			return
		}
		dok, isi, err := svc.Dokumen().DenganKategori(services.KategoriWajibOracle(svc)).BukaDokumenLewatPengenal(r.Context(), pelaku, id)
		if jawabGalatDokumen(w, err) {
			return
		}
		kirimIsiDokumen(w, dok, isi)
	}
}

// isiDokumenKlaim melayani GET /api/klaim-life/{id}/dokumen/{dokId}/isi - "View" berbatas klaim (11-10-2026).
func isiDokumenKlaim(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pelaku := inti.PelakuDari(r, stubPelaku)
		if jawabGalatDokumen(w, inti.WajibIdentitas(pelaku)) {
			return
		}
		if !svc.PunyaDatabase() {
			galat.Tulis(w, http.StatusServiceUnavailable, "The database is not configured.")
			return
		}
		id, ok := pengenalDokumen(w, r)
		if !ok {
			return
		}
		dok, isi, err := svc.Dokumen().BukaDokumen(r.Context(), pelaku, r.PathValue("id"), id)
		if jawabGalatDokumen(w, err) {
			return
		}
		kirimIsiDokumen(w, dok, isi)
	}
}

// kirimIsiDokumen mengalirkan isi dokumen (Google Storage atau salinan lokal) sebagai lampiran unduhan.
func kirimIsiDokumen(w http.ResponseWriter, dok models.Dokumen, isi io.ReadCloser) {
	defer func() { _ = isi.Close() }()
	w.Header().Set("Content-Type", dok.Mime)
	// ⛔ `attachment`, bukan `inline`: isi berkas datang dari pengunggah.
	w.Header().Set("Content-Disposition", "attachment; filename="+strconv.Quote(dok.NamaFile))
	w.Header().Set("X-Content-Type-Options", "nosniff")
	_, _ = io.Copy(w, isi)
}

// hapusDokumen melayani DELETE - `Delete` b4288.
func hapusDokumen(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat.Tulis(w, http.StatusServiceUnavailable, "The database is not configured.")
			return
		}
		id, ok := pengenalDokumen(w, r)
		if !ok {
			return
		}
		err := svc.Dokumen().Hapus(r.Context(), inti.PelakuDari(r, stubPelaku),
			r.PathValue("id"), id, time.Now())
		if jawabGalatDokumen(w, err) {
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

// pengenalDokumen membaca {dokId} sebagai angka.
//
// ⛔ Pengenal dokumen ANGKA, dan nilainya CAP WAKTU
// (`@CurrentDate("yyyyMMddhhmmssSSS")` b648) - bukan nomor urut. 17 angka
// muat di int64; yang lebih panjang ditolak di sini, bukan dibiarkan menjadi
// pencarian yang pasti gagal.
func pengenalDokumen(w http.ResponseWriter, r *http.Request) (int64, bool) {
	n, err := strconv.ParseInt(r.PathValue("dokId"), 10, 64)
	if err != nil || n <= 0 {
		galat.Tulis(w, http.StatusBadRequest, "The document ID is not a valid number.")
		return 0, false
	}
	return n, true
}

// jawabGalatDokumen menerjemahkan galat services menjadi kode HTTP.
func jawabGalatDokumen(w http.ResponseWriter, err error) bool {
	switch {
	case err == nil:
		return false
	case errors.Is(err, inti.ErrTanpaIdentitas):
		galat.Tulis(w, http.StatusUnauthorized, "The request has no user identity.")
	case errors.Is(err, inti.ErrTanpaWewenang):
		galat.Tulis(w, http.StatusForbidden, "You are not authorized to do this.")
	case errors.Is(err, kontrak.ErrKasusSudahTertutup):
		galat.Tulis(w, http.StatusConflict, "The case is already closed and cannot be changed.")
	case errors.Is(err, unggah.ErrUnggahanDirBelumDisetel):
		// 503: bukan salah pemanggil. Ia keadaan server yang belum siap, dan
		// pesannya menyebut apa yang kurang - bukan "gagal mengunggah".
		galat.Tulis(w, http.StatusServiceUnavailable,
			"Document upload is not available: UNGGAHAN_DIR is not set.")
	case errors.Is(err, services.ErrKategoriWajibBelumDiketahui):
		galat.Tulis(w, http.StatusServiceUnavailable,
			"The document category list is not available yet.")
	case errors.Is(err, unggah.ErrBerkasTerlaluBesar):
		galat.Tulis(w, http.StatusRequestEntityTooLarge, err.Error())
	case errors.Is(err, unggah.ErrBerkasKosong):
		galat.Tulis(w, http.StatusBadRequest, "The file is empty.")
	case errors.Is(err, services.ErrStorageBelumSiap):
		// 503: "Storage is not configured" + apa yang kurang (tanpa alamat/token/garam).
		galat.Tulis(w, http.StatusServiceUnavailable, services.PesanStorage(err))
	case errors.Is(err, services.ErrBerkasTidakDiStorage):
		galat.Tulis(w, http.StatusConflict, services.PesanStorage(err))
	case errors.Is(err, services.ErrStorageGagal):
		galat.Tulis(w, http.StatusBadGateway, services.PesanStorage(err))
	case errors.Is(err, services.ErrBerkasDokumenTidakAda):
		galat.Tulis(w, http.StatusNotFound, "The document file was not found in storage.")
	case errors.Is(err, services.ErrKategoriDokumenTerlaluPanjang):
		galat.Tulis(w, http.StatusUnprocessableEntity, err.Error())
	case errors.Is(err, services.ErrKlaimTidakAda):
		galat.Tulis(w, http.StatusNotFound, "The claim does not exist.")
	case errors.Is(err, services.ErrKategoriDokumenTidakDikenal):
		galat.Tulis(w, http.StatusUnprocessableEntity, err.Error())
	case errors.Is(err, services.ErrDokumenBelumTerunggah):
		// 409: barisnya ada, berkasnya belum tertaut. Layar dapat berkata
		// "sedang diproses" alih-alih "tidak ditemukan".
		galat.Tulis(w, http.StatusConflict, "The file has not finished uploading.")
	case errors.Is(err, services.ErrDokumenTidakAda):
		// 404 - uji asap baca-saja DEV (GILIRAN-12) menjumpai 500 di sini.
		galat.Tulis(w, http.StatusNotFound, "The document does not exist.")
	case errors.Is(err, galat.ErrPermintaanTidakSah):
		galat.Tulis(w, http.StatusBadRequest, err.Error())
	default:
		galat.Tulis(w, http.StatusInternalServerError, "Failed to process the document.")
	}
	return true
}

// kategoriDokumen melayani GET /api/klaim-life/{id}/dokumen/kategori - daftar kategori lampiran (CATEGORY_ATTACH_CLAIMLIFE,
// butir ar1) untuk dropdown "Category" (11-10-2026: dulu isian teks bebas, sehingga hampir setiap unggahan ditolak 422
// "kategori tidak dikenal"). Sumbernya SAMA dengan pemeriksa Unggah, jadi pilihan dropdown pasti diterima.
func kategoriDokumen(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if jawabGalatDokumen(w, inti.WajibIdentitas(inti.PelakuDari(r, stubPelaku))) {
			return
		}
		if !svc.PunyaDatabase() {
			galat.Tulis(w, http.StatusServiceUnavailable, "The database is not configured.")
			return
		}
		daftar, err := services.KategoriWajibOracle(svc).KategoriWajib(r.Context(), r.PathValue("id"))
		if jawabGalatDokumen(w, err) {
			return
		}
		if daftar == nil {
			daftar = []string{}
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(struct {
			Kategori []string `json:"kategori"`
		}{daftar})
	}
}
