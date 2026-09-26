package handlers

// Pintu HTTP penghapusan klaim - tiket 15.
//
// Dua pintu, dan pemisahannya disengaja: satu hanya MEMBACA (pratinjau
// dampak), satu MENGHAPUS. Pengguna yang menekan Batal berhenti di pintu
// pertama, yang tidak punya satu pun tulisan untuk dibatalkan.
//
// Dibaca sesudah: handlers.go dan services/hapus.go.

import (
	"encoding/json"
	"errors"
	"net/http"

	"nusantarare/internal/models"
	"nusantarare/internal/services"
)

// dampakJSON adalah bentuk jawaban cacah dampak.
//
// Tiap jenis punya medannya sendiri - layar menampilkannya baris per baris,
// dan satu angka total tidak dapat dipecah kembali.
type dampakJSON struct {
	Header            int `json:"header"`
	Peserta           int `json:"peserta"`
	Adjustment        int `json:"adjustment"`
	Spreading         int `json:"spreading"`
	SpreadingRetro    int `json:"spreadingRetro"`
	Dokumen           int `json:"dokumen"`
	WorkClaim         int `json:"barisWork"`
	BarisDatarWarisan int `json:"barisDatarWarisan"`
	// Total SENGAJA tidak memuat barisDatarWarisan - nasibnya belum diputuskan.
	Total int `json:"total"`
}

func dariDampak(d models.DampakHapus) dampakJSON {
	return dampakJSON{
		Header:  d.Header,
		Peserta: d.Peserta, Adjustment: d.Adjustment, Spreading: d.Spreading,
		SpreadingRetro: d.SpreadingRetro, Dokumen: d.Dokumen,
		WorkClaim: d.WorkClaim, BarisDatarWarisan: d.BarisDatarWarisan,
		Total: d.Total(),
	}
}

// jawabGalatHapus menerjemahkan galat penghapusan ke kode HTTP.
//
// Satu tempat untuk kedua pintu: dua salinan pemetaan galat berarti dua
// kesempatan untuk menjawab berbeda atas sebab yang sama - persis yang terjadi
// pada ErrTanpaWewenang sebelum ia dipisah.
func jawabGalatHapus(w http.ResponseWriter, err error) bool {
	switch {
	case err == nil:
		return false
	case errors.Is(err, services.ErrTanpaIdentitas):
		galat(w, http.StatusUnauthorized, "permintaan tanpa identitas pelaku ditolak")
	case errors.Is(err, services.ErrTanpaWewenang):
		galat(w, http.StatusForbidden, "hanya ReasLifeAdmin yang dapat menghapus klaim")
	case errors.Is(err, services.ErrKlaimTidakAda):
		galat(w, http.StatusNotFound, "klaim tidak ada")
	case errors.Is(err, services.ErrKlaimSudahDiKomite):
		galat(w, http.StatusConflict,
			"klaim sudah diserahkan ke Komite dan tidak dapat dihapus")
	case errors.Is(err, services.ErrHapusFisikDilarang):
		// 501: bukan kerusakan melainkan keputusan yang belum diambil.
		// ADR-U-0031 menetapkan penghapusan berupa PENANDA; kolomnya belum ada.
		galat(w, http.StatusNotImplemented,
			"penghapusan klaim adalah penanda, bukan hapus fisik (ADR-U-0031); "+
				"kolom penandanya belum diputuskan work owner")
	case errors.Is(err, services.ErrPermintaanTidakSah):
		galat(w, http.StatusBadRequest, err.Error())
	default:
		galat(w, http.StatusInternalServerError, "gagal memproses penghapusan klaim")
	}
	return true
}

// dampakHapus melayani GET /api/klaim-life/{id}/dampak-hapus.
func dampakHapus(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat(w, http.StatusServiceUnavailable, "database belum dikonfigurasi")
			return
		}
		d, err := svc.Penghapusan().Dampak(r.Context(), pelakuDari(r, stubPelaku),
			r.PathValue("id"))
		if jawabGalatHapus(w, err) {
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(dariDampak(d))
	}
}

// hapusKlaim melayani DELETE /api/klaim-life/{id}.
//
// Jawabannya membawa cacah yang BENAR-BENAR terhapus, bukan sekadar 204:
// pemanggil dapat membandingkannya dengan angka yang ia tampilkan di
// peringatan, dan selisih apa pun terlihat.
func hapusKlaim(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat(w, http.StatusServiceUnavailable, "database belum dikonfigurasi")
			return
		}
		d, err := svc.Penghapusan().Hapus(r.Context(), pelakuDari(r, stubPelaku),
			r.PathValue("id"))
		if jawabGalatHapus(w, err) {
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(dariDampak(d))
	}
}
