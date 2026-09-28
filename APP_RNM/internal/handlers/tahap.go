package handlers

// Pintu HTTP perpindahan tahap - A3 kelompok Outstanding, butir aw.
//
// `POST /api/klaim-life/{id}/tahap/{tujuan}`
//
// ⛔ BUTIR aw `[DIPUTUSKAN 27-09-2026, dari maksud XML yang terang; cacat
// rule warisan dilaporkan; veto work owner]`.
//
// Dua tombol layar Outstanding memanggil rute ini:
//
//	`Send Back to Register`  `Section/InputOSClaimLife.xml:21404`
//	  -> `<pyLocalAction>SendtoAdmin` 21433 -> `SendtoAdmin_Act`
//	`Send to Medical Check`  21349 / 21839
//	  -> `SendtoAdmin_Act1` 21863
//
// ⚠️ CACAT RULE WARISAN, ditiru MAKSUDnya bukan hurufnya. Kedua activity itu
// berprasyarat `pyWorkPage.pyPosition=="ReasLifeMedicalAdvisor"` dengan
// `WhenTrue=2` / `WhenFalse=3` (lewati), padahal tombolnya berdiri di layar
// Outstanding yang `pyPosition`-nya `ReasLifeAdmin`. Menurut XML apa adanya,
// kedua tombol itu TIDAK MENULIS APA PUN pada posisi Admin - dan kasus tidak
// pernah dapat kembali ke Input Register.
//
// Itu cacat, bukan maksud: label tombolnya, penyambung `Decision3` menuju
// `Assignment2`, dan ADR-U-0002 ketiganya menyebut jalur balik ini sebagai
// fitur. Cacatnya dilaporkan di `OQ-untuk-tim.md` (OQ-C) dengan barisnya.
//
// Dibaca sesudah: inbox.go.

import (
	"errors"
	"net/http"
	"time"

	"nusantarare/internal/models"
	"nusantarare/internal/services"
)

// tahapTujuan menerjemahkan potongan jalur menjadi tahap.
//
// ⛔ Himpunan TERTUTUP dan memakai KATA, bukan angka: jalur yang berisi
// `/tahap/2` tidak dapat dibaca siapa pun, dan angka yang bergeser bila
// urutan `models.Tahap` berubah akan memindahkan kasus ke tempat yang salah
// tanpa satu pun galat.
func tahapTujuan(potongan string) (models.Tahap, bool) {
	switch potongan {
	case "input-register":
		return models.TahapInputRegister, true
	case "outstanding":
		return models.TahapOutstanding, true
	case "medical-check":
		return models.TahapMedicalCheck, true
	case "claim-analis":
		return models.TahapClaimAnalis, true
	default:
		return models.TahapTidakDikenal, false
	}
}

// pindahTahap melayani POST /api/klaim-life/{id}/tahap/{tujuan}.
func pindahTahap(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat(w, http.StatusServiceUnavailable, "database belum dikonfigurasi")
			return
		}
		tujuan, sah := tahapTujuan(r.PathValue("tujuan"))
		if !sah {
			galat(w, http.StatusNotFound, "tahap tujuan tidak dikenal")
			return
		}
		err := svc.Tahap().DenganJejak(services.PerekamJejakOracle(svc)).
			Pindah(r.Context(), pelakuDari(r, stubPelaku), r.PathValue("id"),
				tujuan, time.Now())

		switch {
		case err == nil:
			w.WriteHeader(http.StatusNoContent)
		case errors.Is(err, services.ErrTanpaIdentitas):
			galat(w, http.StatusUnauthorized, "permintaan tanpa identitas pelaku ditolak")
		case errors.Is(err, services.ErrTanpaWewenang):
			// ⛔ Orang memindahkan pekerjaan yang SEDANG IA PEGANG.
			galat(w, http.StatusForbidden,
				"hanya pemegang tahap asal yang dapat memindahkan kasus ini")
		case errors.Is(err, services.ErrPerpindahanTidakSah):
			// ⛔ 409, bukan 400: permintaannya berbentuk benar, keadaan
			// kasusnya yang tidak mengizinkan. Tangga yang setiap anaknya
			// dapat dilompati bukan tangga.
			galat(w, http.StatusConflict,
				"perpindahan itu tidak ada di tangga kerja klaim")
		case errors.Is(err, services.ErrKasusSudahTertutup):
			galat(w, http.StatusConflict, "kasus sudah ditutup dan tidak dapat diubah")
		case errors.Is(err, services.ErrTahapTidakDikenal):
			galat(w, http.StatusConflict, "tahap kasus ini tidak dikenal")
		default:
			galat(w, http.StatusInternalServerError, "gagal memindahkan tahap klaim")
		}
	}
}
