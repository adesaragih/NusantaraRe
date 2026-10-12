package handlers

// Pintu HTTP grid diagnosa - butir bd, §2 A bagian 2.
//
//	POST   /api/klaim-life/{id}/peserta/{pesertaId}/diagnosa
//	PUT    /api/klaim-life/{id}/peserta/{pesertaId}/diagnosa/{diagId}
//	DELETE /api/klaim-life/{id}/peserta/{pesertaId}/diagnosa/{diagId}
//
// Ketiganya padanan tombol grid `.DiagnoseList`: `Add` b4690, `Choose`
// b2509 (lewat `SetDisease`), dan `Delete` b6160.
//
// Nol aturan dagang di sini - ketiga gerbangnya (tahap, pemegang,
// `STS_REJECT` peserta) hidup di `services/diagnosa.go`. Berkas ini hanya
// menerjemahkan HTTP ke panggilan services dan sebaliknya.
//
// ⛔ Jalurnya BERSARANG di bawah pesertanya, dan itu bukan hiasan REST:
// diagnosa tidak punya hidup di luar peserta yang memuatnya
// (`SetDisease.xml` b389 `Obj-Save pyWorkPage`). Jalur `/api/diagnosa/{id}`
// akan membuat pengenal peserta menjadi opsional, dan gerbang yang
// memeriksa peserta menjadi gerbang yang dapat dilewati.
//
// Dibaca sesudah: penyakit.go (pencarian katalognya) dan tutup.go.

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/inti/backend/kontrak"
	"nusantarare/modul/claimlife/backend/models"
	"nusantarare/modul/claimlife/backend/services"
)

// isiDiagnosa adalah badan permintaan `PUT`.
//
// ⚠️ Nama medan sama dengan yang dipakai `Penyakit` pada pencarian
// (`kodeIcd`, `nama`), dan itu disengaja: layar menyalin satu baris hasil
// pencarian ke satu baris diagnosa, dan dua kosakata untuk satu nilai adalah
// satu tempat lagi untuk tertukar. Cacat lintas-lapis yang sama sudah lima
// kali terjadi di modul ini.
type isiDiagnosa struct {
	KodeIcd       string `json:"kodeIcd"`
	Nama          string `json:"nama"`
	GroupDiagnose string `json:"groupDiagnose"`
}

// tambahDiagnosa melayani POST - `Add` b4690.
func tambahDiagnosa(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat.Tulis(w, http.StatusServiceUnavailable, "The database is not configured.")
			return
		}
		lahir, err := svc.Diagnosa().Tambah(r.Context(), inti.PelakuDari(r, stubPelaku),
			r.PathValue("id"), r.PathValue("pesertaId"))
		if jawabGalatDiagnosa(w, err) {
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(lahir)
	}
}

// ubahDiagnosa melayani PUT - `Choose` b2509 -> `SetDisease`.
// kelompokDiagnosa melayani GET /api/klaim-life/kelompok-diagnosa - daftar pilihan GROUP DIAGNOSE (keputusan work owner
// 12-10-2026), SATU sumber: models.KelompokDiagnosa. Baca saja, tanpa database.
func kelompokDiagnosa() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(struct {
			Daftar []string `json:"daftar"`
		}{Daftar: models.KelompokDiagnosa()})
	}
}

func ubahDiagnosa(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat.Tulis(w, http.StatusServiceUnavailable, "The database is not configured.")
			return
		}
		diagID, ok := pengenalDiagnosa(w, r)
		if !ok {
			return
		}
		var isi isiDiagnosa
		if err := json.NewDecoder(r.Body).Decode(&isi); err != nil {
			galat.Tulis(w, http.StatusBadRequest, "The request body is not valid JSON.")
			return
		}
		err := svc.Diagnosa().Ubah(r.Context(), inti.PelakuDari(r, stubPelaku),
			r.PathValue("id"), r.PathValue("pesertaId"), diagID,
			isi.KodeIcd, isi.Nama, isi.GroupDiagnose)
		if jawabGalatDiagnosa(w, err) {
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

// hapusDiagnosa melayani DELETE - `Delete` b6160.
func hapusDiagnosa(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat.Tulis(w, http.StatusServiceUnavailable, "The database is not configured.")
			return
		}
		diagID, ok := pengenalDiagnosa(w, r)
		if !ok {
			return
		}
		err := svc.Diagnosa().Hapus(r.Context(), inti.PelakuDari(r, stubPelaku),
			r.PathValue("id"), r.PathValue("pesertaId"), diagID)
		if jawabGalatDiagnosa(w, err) {
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

// pengenalDiagnosa membaca {diagId} sebagai angka, dan menjawab 400 bila
// bukan.
//
// ⛔ Pengenal diagnosa ANGKA - berbeda dengan pengenal klaim dan peserta,
// yang teks (ADR-U-0022). Ia milik kita, dari `SEQ_CLAIMLF_DIAGNOSE`, dan
// tidak pernah datang dari sistem lama. Angka yang tidak terbaca dijawab 400
// di sini alih-alih menjadi 0 yang diteruskan ke services: pengenal 0 akan
// dicari, tidak ketemu, dan dijawab "bukan milik peserta" - kalimat yang
// benar tentang hal yang salah.
func pengenalDiagnosa(w http.ResponseWriter, r *http.Request) (int64, bool) {
	n, err := strconv.ParseInt(r.PathValue("diagId"), 10, 64)
	if err != nil || n <= 0 {
		galat.Tulis(w, http.StatusBadRequest, "The diagnosis ID is not a valid number.")
		return 0, false
	}
	return n, true
}

// jawabGalatDiagnosa menerjemahkan galat services menjadi kode HTTP.
//
// Mengembalikan true bila permintaan SUDAH dijawab.
//
// ⛔ Satu terjemahan untuk ketiga rute. Tiga salinan berarti tiga kesempatan
// untuk menjawab 500 atas hal yang sebenarnya 409 - dan 500 adalah jawaban
// yang membuat orang menelepon, bukan memperbaiki isiannya.
func jawabGalatDiagnosa(w http.ResponseWriter, err error) bool {
	switch {
	case err == nil:
		return false
	case errors.Is(err, inti.ErrTanpaIdentitas):
		galat.Tulis(w, http.StatusUnauthorized, "The request has no user identity.")
	case errors.Is(err, inti.ErrTanpaWewenang):
		galat.Tulis(w, http.StatusForbidden,
			"Only the holder of this stage can change the diagnosis.")
	case errors.Is(err, services.ErrDiagnosaTerkunci):
		// 409: bukan permintaan yang salah, melainkan bentrokan dengan
		// keadaan yang sudah ada. Padanan `pyDisabledWhen` b4682/b5059/
		// b5870/b6152 - di Pega tombolnya mati, di sini permintaannya
		// ditolak dengan kalimat yang menyebut sebabnya.
		galat.Tulis(w, http.StatusConflict,
			"This insured has been decided; the diagnosis cannot be changed.")
	case errors.Is(err, kontrak.ErrKasusSudahTertutup):
		galat.Tulis(w, http.StatusConflict, "The case is already closed and cannot be changed.")
	case errors.Is(err, services.ErrTahapTidakBergridPeserta):
		galat.Tulis(w, http.StatusConflict,
			"This stage does not open the insured detail.")
	case errors.Is(err, services.ErrDiagnosaBukanMedicalCheck):
		galat.Tulis(w, http.StatusConflict,
			"the diagnosis can only be changed at the Medical Check stage")
	case errors.Is(err, services.ErrKelompokDiagnosaTidakDikenal):
		galat.Tulis(w, http.StatusUnprocessableEntity, "Group Diagnose must be one of the listed groups.")
	case errors.Is(err, services.ErrNilaiDiagnosaKepanjangan):
		galat.Tulis(w, http.StatusUnprocessableEntity, err.Error())
	case errors.Is(err, galat.ErrPermintaanTidakSah):
		galat.Tulis(w, http.StatusBadRequest, err.Error())
	default:
		galat.Tulis(w, http.StatusInternalServerError, "Failed to change the diagnosis.")
	}
	return true
}
