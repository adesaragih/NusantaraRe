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

	"nusantarare/internal/services"
)

// batasFormulir membatasi bagian NON-berkas sebuah multipart.
//
// ⚠️ Ini BUKAN batas ukuran berkas - itu `services.BatasUkuranUnggahan`, dan
// ia ditegakkan saat menyalin. Yang di sini membatasi berapa banyak formulir
// yang ditahan di MEMORI sebelum bagian berkasnya dialirkan ke disk.
const batasFormulir = 1 << 20

// unggahDokumen melayani POST multipart - `Add attachment` b1245.
func unggahDokumen(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat(w, http.StatusServiceUnavailable, "database belum dikonfigurasi")
			return
		}
		// ⛔ `ReadForm` DIHINDARI: ia menulis seluruh berkas ke berkas
		// sementara lebih dulu. `MultipartReader` mengalirkannya, sehingga
		// batas ukuran berlaku SEBELUM 25 MiB mendarat dua kali.
		if err := r.ParseMultipartForm(batasFormulir); err != nil {
			galat(w, http.StatusBadRequest, "permintaan bukan multipart yang sah")
			return
		}
		berkas, kepala, err := r.FormFile("berkas")
		if err != nil {
			galat(w, http.StatusBadRequest, "bagian `berkas` tidak ada di permintaan")
			return
		}
		defer func() { _ = berkas.Close() }()

		dok, err := svc.Dokumen().
			DenganKategori(services.KategoriWajibOracle(svc)).
			Unggah(r.Context(), pelakuDari(r, stubPelaku),
				r.PathValue("id"), r.PathValue("pesertaId"),
				services.BerkasMasuk{
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
func isiDokumen(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat(w, http.StatusServiceUnavailable, "database belum dikonfigurasi")
			return
		}
		id, ok := pengenalDokumen(w, r)
		if !ok {
			return
		}
		dok, jalur, err := svc.Dokumen().UnduhLewatPengenal(r.Context(),
			pelakuDari(r, stubPelaku), id)
		if jawabGalatDokumen(w, err) {
			return
		}
		f, err := os.Open(jalur)
		if err != nil {
			// ⛔ 404, bukan 500. Baris ada tetapi berkasnya tidak: itu
			// keadaan yang pemakai dapat laporkan dengan jelas, dan 500
			// membuatnya terbaca sebagai kerusakan server.
			galat(w, http.StatusNotFound, "berkas dokumen tidak ditemukan di penyimpanan")
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
			galat(w, http.StatusServiceUnavailable, "database belum dikonfigurasi")
			return
		}
		id, ok := pengenalDokumen(w, r)
		if !ok {
			return
		}
		err := svc.Dokumen().Hapus(r.Context(), pelakuDari(r, stubPelaku),
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
		galat(w, http.StatusBadRequest, "pengenal dokumen bukan angka yang sah")
		return 0, false
	}
	return n, true
}

// jawabGalatDokumen menerjemahkan galat services menjadi kode HTTP.
func jawabGalatDokumen(w http.ResponseWriter, err error) bool {
	switch {
	case err == nil:
		return false
	case errors.Is(err, services.ErrTanpaIdentitas):
		galat(w, http.StatusUnauthorized, "permintaan tanpa identitas pelaku ditolak")
	case errors.Is(err, services.ErrTanpaWewenang):
		galat(w, http.StatusForbidden, "wewenang tidak mencukupi")
	case errors.Is(err, services.ErrKasusSudahTertutup):
		galat(w, http.StatusConflict, "kasus sudah ditutup dan tidak dapat diubah")
	case errors.Is(err, services.ErrUnggahanDirBelumDisetel):
		// 503: bukan salah pemanggil. Ia keadaan server yang belum siap, dan
		// pesannya menyebut apa yang kurang - bukan "gagal mengunggah".
		galat(w, http.StatusServiceUnavailable,
			"UNGGAHAN_DIR belum disetel; unggahan dokumen belum dapat dilayani")
	case errors.Is(err, services.ErrKategoriWajibBelumDiketahui):
		galat(w, http.StatusServiceUnavailable,
			"daftar kategori dokumen belum tersedia")
	case errors.Is(err, services.ErrBerkasTerlaluBesar):
		galat(w, http.StatusRequestEntityTooLarge, err.Error())
	case errors.Is(err, services.ErrBerkasKosong):
		galat(w, http.StatusBadRequest, "berkas kosong")
	case errors.Is(err, services.ErrKategoriDokumenTidakDikenal):
		galat(w, http.StatusUnprocessableEntity, err.Error())
	case errors.Is(err, services.ErrDokumenBelumTerunggah):
		// 409: barisnya ada, berkasnya belum tertaut. Layar dapat berkata
		// "sedang diproses" alih-alih "tidak ditemukan".
		galat(w, http.StatusConflict, "berkas belum selesai diunggah")
	case errors.Is(err, services.ErrPermintaanTidakSah):
		galat(w, http.StatusBadRequest, err.Error())
	default:
		galat(w, http.StatusInternalServerError, "gagal memproses dokumen")
	}
	return true
}
