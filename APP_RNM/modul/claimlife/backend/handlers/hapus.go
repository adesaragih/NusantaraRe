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

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/inti/backend/kontrak"
	"nusantarare/modul/claimlife/backend/models"
	"nusantarare/modul/claimlife/backend/services"
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

// metodeKlaimDiizinkan mengisi header `Allow` pada jawaban 405.
//
// ⚠️ Daftarnya harus cocok dengan rute yang benar-benar terdaftar untuk
// `/api/klaim-life/{id}` di handlers.go. Penjaga statik menegakkannya:
// header `Allow` yang berbohong lebih buruk daripada tidak ada.
const metodeKlaimDiizinkan = "GET"

// jawabGalatHapus menerjemahkan galat penghapusan ke kode HTTP.
//
// Satu tempat untuk kedua pintu: dua salinan pemetaan galat berarti dua
// kesempatan untuk menjawab berbeda atas sebab yang sama - persis yang terjadi
// pada ErrTanpaWewenang sebelum ia dipisah.
func jawabGalatHapus(w http.ResponseWriter, err error) bool {
	switch {
	case err == nil:
		return false
	case errors.Is(err, inti.ErrTanpaIdentitas):
		galat.Tulis(w, http.StatusUnauthorized, "permintaan tanpa identitas pelaku ditolak")
	case errors.Is(err, inti.ErrTanpaWewenang):
		galat.Tulis(w, http.StatusForbidden, "hanya ReasLifeAdmin yang dapat menghapus klaim")
	case errors.Is(err, services.ErrKlaimTidakAda):
		galat.Tulis(w, http.StatusNotFound, "klaim tidak ada")
	case errors.Is(err, services.ErrKlaimSudahDiKomite):
		galat.Tulis(w, http.StatusConflict,
			"klaim sudah diserahkan ke Komite dan tidak dapat dihapus")
	case errors.Is(err, services.ErrHapusFisikDilarang):
		// ⛔ 405, BUKAN 501. Ronde sebelumnya menjawab 501 - dan 501
		// berarti "server ini belum bisa", yaitu janji bahwa suatu hari ia
		// bisa. ADR-U-0031 justru menetapkan sebaliknya: penghapusan klaim
		// SELALU berupa penanda + nilai pembalik, TIDAK PERNAH fisik. Jadi
		// DELETE atas sumber daya ini bukan metode yang belum siap,
		// melainkan metode yang memang tidak berlaku - dan itu 405.
		//
		// ⚠️ Header `Allow` WAJIB menyertai 405 (RFC 9110 §15.5.6).
		// Tanpanya, klien tidak diberi tahu apa yang boleh ia lakukan
		// sebagai gantinya, dan 405 menjadi penolakan buta.
		w.Header().Set("Allow", metodeKlaimDiizinkan)
		galat.Tulis(w, http.StatusMethodNotAllowed,
			"penghapusan klaim selalu berupa penanda dan nilai pembalik, "+
				"tidak pernah hapus fisik (ADR-U-0031); DELETE tidak berlaku "+
				"atas sumber daya ini")
	case errors.Is(err, kontrak.ErrKasusSudahTertutup):
		galat.Tulis(w, http.StatusConflict, "kasus sudah ditutup dan tidak dapat diubah")
	case errors.Is(err, galat.ErrPermintaanTidakSah):
		galat.Tulis(w, http.StatusBadRequest, err.Error())
	default:
		galat.Tulis(w, http.StatusInternalServerError, "gagal memproses penghapusan klaim")
	}
	return true
}

// dampakHapus melayani GET /api/klaim-life/{id}/dampak-hapus.
func dampakHapus(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat.Tulis(w, http.StatusServiceUnavailable, "database belum dikonfigurasi")
			return
		}
		d, err := svc.Penghapusan().Dampak(r.Context(), inti.PelakuDari(r, stubPelaku),
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
			galat.Tulis(w, http.StatusServiceUnavailable, "database belum dikonfigurasi")
			return
		}
		d, err := svc.Penghapusan().Hapus(r.Context(), inti.PelakuDari(r, stubPelaku),
			r.PathValue("id"))
		if jawabGalatHapus(w, err) {
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(dariDampak(d))
	}
}
