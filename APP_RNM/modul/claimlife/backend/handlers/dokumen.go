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
	"os"
	"strconv"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/inti/backend/kontrak"
	"nusantarare/inti/backend/unggah"
	"nusantarare/modul/claimlife/backend/services"
)

// unggahDokumen melayani POST multipart - `Add attachment` b1245.
func unggahDokumen(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat.Tulis(w, http.StatusServiceUnavailable, "database belum dikonfigurasi")
			return
		}
		// ⛔ `ReadForm` DIHINDARI: ia menulis seluruh berkas ke berkas
		// sementara lebih dulu. `MultipartReader` mengalirkannya, sehingga
		// batas ukuran berlaku SEBELUM 25 MiB mendarat dua kali.
		if err := r.ParseMultipartForm(galat.BatasFormulir); err != nil {
			galat.Tulis(w, http.StatusBadRequest, "permintaan bukan multipart yang sah")
			return
		}
		berkas, kepala, err := r.FormFile("berkas")
		if err != nil {
			galat.Tulis(w, http.StatusBadRequest, "bagian `berkas` tidak ada di permintaan")
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
			galat.Tulis(w, http.StatusServiceUnavailable, "database belum dikonfigurasi")
			return
		}
		id, ok := pengenalDokumen(w, r)
		if !ok {
			return
		}
		dok, jalur, err := svc.Dokumen().UnduhLewatPengenal(r.Context(), pelaku, id)
		if jawabGalatDokumen(w, err) {
			return
		}
		f, err := os.Open(jalur)
		if err != nil {
			// ⛔ 404, bukan 500. Baris ada tetapi berkasnya tidak: itu
			// keadaan yang pemakai dapat laporkan dengan jelas, dan 500
			// membuatnya terbaca sebagai kerusakan server.
			galat.Tulis(w, http.StatusNotFound, "berkas dokumen tidak ditemukan di penyimpanan")
			return
		}
		defer func() { _ = f.Close() }()

		w.Header().Set("Content-Type", dok.Mime)
		// ⛔ `attachment`, bukan `inline`. Berkas yang dirender di dalam
		// halaman kita berjalan di asal kita - dan isi berkas itu datang
		// dari pengunggah.
		w.Header().Set("Content-Disposition",
			"attachment; filename="+strconv.Quote(dok.NamaFile))
		w.Header().Set("X-Content-Type-Options", "nosniff")
		_, _ = io.Copy(w, f)
	}
}

// hapusDokumen melayani DELETE - `Delete` b4288.
func hapusDokumen(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat.Tulis(w, http.StatusServiceUnavailable, "database belum dikonfigurasi")
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
		galat.Tulis(w, http.StatusBadRequest, "pengenal dokumen bukan angka yang sah")
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
		galat.Tulis(w, http.StatusUnauthorized, "permintaan tanpa identitas pelaku ditolak")
	case errors.Is(err, inti.ErrTanpaWewenang):
		galat.Tulis(w, http.StatusForbidden, "wewenang tidak mencukupi")
	case errors.Is(err, kontrak.ErrKasusSudahTertutup):
		galat.Tulis(w, http.StatusConflict, "kasus sudah ditutup dan tidak dapat diubah")
	case errors.Is(err, unggah.ErrUnggahanDirBelumDisetel):
		// 503: bukan salah pemanggil. Ia keadaan server yang belum siap, dan
		// pesannya menyebut apa yang kurang - bukan "gagal mengunggah".
		galat.Tulis(w, http.StatusServiceUnavailable,
			"UNGGAHAN_DIR belum disetel; unggahan dokumen belum dapat dilayani")
	case errors.Is(err, services.ErrKategoriWajibBelumDiketahui):
		galat.Tulis(w, http.StatusServiceUnavailable,
			"daftar kategori dokumen belum tersedia")
	case errors.Is(err, unggah.ErrBerkasTerlaluBesar):
		galat.Tulis(w, http.StatusRequestEntityTooLarge, err.Error())
	case errors.Is(err, unggah.ErrBerkasKosong):
		galat.Tulis(w, http.StatusBadRequest, "berkas kosong")
	case errors.Is(err, services.ErrKategoriDokumenTidakDikenal):
		galat.Tulis(w, http.StatusUnprocessableEntity, err.Error())
	case errors.Is(err, services.ErrDokumenBelumTerunggah):
		// 409: barisnya ada, berkasnya belum tertaut. Layar dapat berkata
		// "sedang diproses" alih-alih "tidak ditemukan".
		galat.Tulis(w, http.StatusConflict, "berkas belum selesai diunggah")
	case errors.Is(err, services.ErrDokumenTidakAda):
		// 404 - uji asap baca-saja DEV (GILIRAN-12) menjumpai 500 di sini.
		galat.Tulis(w, http.StatusNotFound, "dokumen tidak ada")
	case errors.Is(err, galat.ErrPermintaanTidakSah):
		galat.Tulis(w, http.StatusBadRequest, err.Error())
	default:
		galat.Tulis(w, http.StatusInternalServerError, "gagal memproses dokumen")
	}
	return true
}
