package handlers

// Pintu HTTP dialog Edit Date - DOL (tiket 06) dan tiga tanggal klaim lain
// (sensus §3.1, 28-09-2026).
//
// Nol aturan dagang di sini: berkas ini menerjemahkan HTTP ke panggilan
// services dan sebaliknya. Yang memutuskan sah atau tidak adalah
// services.ValidasiDOL (DOL) dan TanggalKejadian.SetTanggalKlaim (gerbang
// tahap dan peran tiga tanggal lainnya).
//
// Dibaca sesudah: handlers.go dan services/dol.go.

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

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
		case err != nil:
			jawabGalatTanggal(w, err, "gagal menyimpan tanggal kejadian")
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// jawabGalatTanggal menerjemahkan galat bersama KEDUA rute dialog Edit Date.
//
// ⛔ Satu terjemahan, sebab satu dialog. Sebelum 28-09-2026 rute DOL
// menjawab 500 atas kasus yang sudah ditutup - `ErrKasusSudahTertutup` tidak
// ada di daftarnya, dan 500 membuat orang menelepon, bukan membaca.
func jawabGalatTanggal(w http.ResponseWriter, err error, gagal string) {
	switch {
	case errors.Is(err, services.ErrPermintaanTidakSah):
		galat(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, services.ErrTanpaIdentitas):
		galat(w, http.StatusUnauthorized, "permintaan tanpa identitas pelaku ditolak")
	case errors.Is(err, services.ErrTanpaWewenang):
		galat(w, http.StatusForbidden, "peran tidak mencukupi")
	case errors.Is(err, services.ErrKasusSudahTertutup):
		galat(w, http.StatusConflict, "kasus sudah ditutup dan tidak dapat diubah")
	case errors.Is(err, services.ErrTahapTidakBolehUbahTanggal):
		galat(w, http.StatusConflict, "tanggal klaim hanya dapat diubah di tahap Outstanding Claim")
	case errors.Is(err, services.ErrTahapTidakDikenal):
		// 422: datanya yang tidak lengkap - tahap kasus tidak dapat dibaca.
		galat(w, http.StatusUnprocessableEntity, "tahap kasus tidak dikenal")
	default:
		galat(w, http.StatusInternalServerError, gagal)
	}
}

// permintaanTanggalKlaimJSON - ketiga isian dialog selain DOL. Kosong SAH:
// isian yang dikosongkan menjadi NULL, sama seperti di Pega.
type permintaanTanggalKlaimJSON struct {
	TanggalTerimaKlaim    string `json:"tanggalTerimaKlaim"`
	TanggalDokumenLengkap string `json:"tanggalDokumenLengkap"`
	TanggalKonfirmasi     string `json:"tanggalKonfirmasi"`
}

// uraiTanggalOpsional mengurai satu isian; kosong menjadi nil.
func uraiTanggalOpsional(teks string) (*time.Time, error) {
	if strings.TrimSpace(teks) == "" {
		return nil, nil
	}
	t, err := utils.ParseTanggal(strings.TrimSpace(teks))
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// setTanggalKlaim melayani
// PUT /api/klaim-life/{id}/peserta/{pesertaId}/tanggal-klaim.
//
// Tombol `Save` `EditDateClaimLife_Section.xml` b1910 -> `UpdateDateClaimLife_Act`.
func setTanggalKlaim(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat(w, http.StatusServiceUnavailable, "database belum dikonfigurasi")
			return
		}
		var masuk permintaanTanggalKlaimJSON
		if err := json.NewDecoder(r.Body).Decode(&masuk); err != nil {
			galat(w, http.StatusBadRequest, "badan permintaan bukan JSON yang sah")
			return
		}
		var tgl models.TanggalKlaim
		for _, isian := range []struct {
			nama string
			teks string
			ke   **time.Time
		}{
			{"tanggalTerimaKlaim", masuk.TanggalTerimaKlaim, &tgl.TerimaKlaim},
			{"tanggalDokumenLengkap", masuk.TanggalDokumenLengkap, &tgl.DokumenLengkap},
			{"tanggalKonfirmasi", masuk.TanggalKonfirmasi, &tgl.Konfirmasi},
		} {
			t, err := uraiTanggalOpsional(isian.teks)
			if err != nil {
				galat(w, http.StatusBadRequest, isian.nama+" bukan tanggal yang dikenal")
				return
			}
			*isian.ke = t
		}
		if err := svc.TanggalKejadian().SetTanggalKlaim(r.Context(), pelakuDari(r, stubPelaku),
			r.PathValue("id"), r.PathValue("pesertaId"), tgl); err != nil {
			jawabGalatTanggal(w, err, "gagal menyimpan tanggal klaim")
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
