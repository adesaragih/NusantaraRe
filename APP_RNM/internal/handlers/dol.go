package handlers

// Pintu HTTP tanggal kejadian - tiket 06.
//
// Nol aturan dagang di sini: berkas ini menerjemahkan HTTP ke panggilan
// services dan sebaliknya. Yang memutuskan sah atau tidak adalah
// services.ValidasiDOL.
//
// Dibaca sesudah: handlers.go dan services/dol.go.

import (
	"encoding/json"
	"errors"
	"net/http"

	"nusantarare/internal/models"
	"nusantarare/internal/services"
	"nusantarare/pkg/utils"
)

// permintaanDOLJSON adalah bentuk badan permintaan pengisian tanggal kejadian.
//
// Tanggal dikirim sebagai TEKS berformat tetap, tidak pernah sebagai angka
// epoch: bentuk yang dapat dibaca manusia itulah yang dipulangkan pembaca dan
// yang tertulis di basis data (ADR-U-0022).
type permintaanDOLJSON struct {
	TanggalKejadian string `json:"tanggalKejadian"`
}

// setTanggalKejadian melayani
// PUT /api/klaim-life/{id}/peserta/{pesertaId}/tanggal-kejadian.
func setTanggalKejadian(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat(w, http.StatusServiceUnavailable, "database belum dikonfigurasi")
			return
		}
		var masuk permintaanDOLJSON
		if err := json.NewDecoder(r.Body).Decode(&masuk); err != nil {
			galat(w, http.StatusBadRequest, "badan permintaan bukan JSON yang sah")
			return
		}
		dol, err := utils.ParseTanggal(masuk.TanggalKejadian)
		if err != nil {
			galat(w, http.StatusBadRequest, "tanggalKejadian bukan tanggal yang dikenal")
			return
		}

		err = svc.TanggalKejadian().Set(r.Context(), pelakuDari(r, stubPelaku),
			r.PathValue("id"), r.PathValue("pesertaId"), dol)

		// ⭐ GalatDOL ditangani LEBIH DULU daripada kerabatnya yang umum:
		// kalimatnya adalah kalimat XML apa adanya, dan peserta yang
		// bermasalah dipulangkan di medannya sendiri - bukan disisipkan ke
		// dalam kalimat itu.
		var dolSalah services.GalatDOL
		switch {
		case errors.As(err, &dolSalah):
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusUnprocessableEntity)
			_ = json.NewEncoder(w).Encode(struct {
				Galat     string `json:"galat"`
				PesertaID string `json:"pesertaId"`
			}{dolSalah.Error(), dolSalah.PesertaID})
			return
		case errors.Is(err, services.ErrDOLKosong):
			galat(w, http.StatusBadRequest, "tanggal kejadian belum diisi")
			return
		case errors.Is(err, services.ErrTypeTidakDikenal):
			// ⛔ 422, bukan 500: datanya yang belum lengkap, bukan aplikasinya
			// yang rusak. Type klaim di luar keempat nilai yang dikenal berarti
			// jendela valuasi mana yang berlaku belum dapat ditentukan.
			galat(w, http.StatusUnprocessableEntity,
				"Type klaim tidak dikenal; jendela valuasi tidak dapat ditentukan")
			return
		case errors.Is(err, models.ErrValuasiKosong):
			galat(w, http.StatusUnprocessableEntity,
				"tanggal valuasi peserta kosong; DOL tidak dapat divalidasi")
			return
		case errors.Is(err, services.ErrPermintaanTidakSah):
			galat(w, http.StatusBadRequest, err.Error())
			return
		case errors.Is(err, services.ErrTanpaIdentitas):
			galat(w, http.StatusUnauthorized, "permintaan tanpa identitas pelaku ditolak")
			return
		case errors.Is(err, services.ErrTanpaWewenang):
			galat(w, http.StatusForbidden, "peran tidak mencukupi")
			return
		case err != nil:
			galat(w, http.StatusInternalServerError, "gagal menyimpan tanggal kejadian")
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
