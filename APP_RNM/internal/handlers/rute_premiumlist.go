package handlers

// Pintu HTTP modul PremiumList Life - tiket 01.
//
//	POST /api/polis-life/{id}/keputusan   body {"keputusan":"Confirm|Reject|Decline"}
//	POST /api/polis-life/{id}/penggolong  body {"hasil":"Offer|Premium"}
//
// Padanan ketiga konektor keputusan `InputPolicyHolder.xml` dan penggolong
// `Decision3`-nya.
//
// ⛔ DUA rute, bukan satu ber-`keputusan` lima nilai. `Offer`/`Premium`
// BUKAN keputusan pengguna atas penawaran - ia hasil penggolong yang hanya
// sah SESUDAH `Confirm` di tahap penawaran. Satu rute untuk keduanya membuat
// gerbang "hanya sesudah Confirm" menjadi gerbang yang harus diingat, bukan
// gerbang yang ada.
//
// Nol aturan dagang di sini; seluruh gerbangnya di `services/polis_penawaran.go`.
//
// ⚠️ Berkas ini MILIK sesi PremiumList Life. `handlers.go` disentuh hanya
// dengan penambahan rute - perubahan bersama yang aditif, dan dilaporkan.

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"nusantarare/internal/models"
	"nusantarare/internal/services"
)

// isiKeputusanPolis adalah badan permintaan keputusan.
type isiKeputusanPolis struct {
	Keputusan string `json:"keputusan"`
}

// isiPenggolongPolis adalah badan permintaan penggolong.
type isiPenggolongPolis struct {
	Hasil string `json:"hasil"`
}

// jawabanAkibat adalah apa yang layar perlu tahu sesudah sebuah keputusan.
//
// ⚠️ Ketiga medannya menyeberang BERNAMA dan lengkap. Layar harus dapat
// membedakan tiga hasil yang berbeda - berpindah, tertutup, menunggu
// penggolong - dan jawaban yang hanya berkata "berhasil" memaksa layar
// membaca ulang seluruh polis untuk menebak yang mana.
type jawabanAkibat struct {
	TahapTujuan        string `json:"tahapTujuan"`
	StatusWork         string `json:"statusWork"`
	MenungguPenggolong bool   `json:"menungguPenggolong"`
}

// putuskanPenawaran melayani POST .../keputusan.
func putuskanPenawaran(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat(w, http.StatusServiceUnavailable, "database belum dikonfigurasi")
			return
		}
		var isi isiKeputusanPolis
		if err := json.NewDecoder(r.Body).Decode(&isi); err != nil {
			galat(w, http.StatusBadRequest, "badan permintaan bukan JSON yang sah")
			return
		}
		akibat, err := svc.Penawaran().
			DenganJejak(services.PerekamJejakOracle(svc)).
			Putuskan(r.Context(), pelakuDari(r, stubPelaku),
				r.PathValue("id"), isi.Keputusan, time.Now())
		if jawabGalatPolis(w, err) {
			return
		}
		tulisAkibat(w, akibat)
	}
}

// golongkanPenawaran melayani POST .../penggolong.
func golongkanPenawaran(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat(w, http.StatusServiceUnavailable, "database belum dikonfigurasi")
			return
		}
		var isi isiPenggolongPolis
		if err := json.NewDecoder(r.Body).Decode(&isi); err != nil {
			galat(w, http.StatusBadRequest, "badan permintaan bukan JSON yang sah")
			return
		}
		akibat, err := svc.Penawaran().
			DenganJejak(services.PerekamJejakOracle(svc)).
			Golongkan(r.Context(), pelakuDari(r, stubPelaku),
				r.PathValue("id"), isi.Hasil, time.Now())
		if jawabGalatPolis(w, err) {
			return
		}
		tulisAkibat(w, akibat)
	}
}

func tulisAkibat(w http.ResponseWriter, a models.AkibatKeputusan) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(jawabanAkibat{
		TahapTujuan:        a.TahapTujuan,
		StatusWork:         a.StatusWork,
		MenungguPenggolong: a.MenungguPenggolong,
	})
}

// jawabGalatPolis menerjemahkan galat services menjadi kode HTTP.
//
// Mengembalikan true bila permintaan SUDAH dijawab.
func jawabGalatPolis(w http.ResponseWriter, err error) bool {
	switch {
	case err == nil:
		return false
	case errors.Is(err, services.ErrTanpaIdentitas):
		galat(w, http.StatusUnauthorized, "permintaan tanpa identitas pelaku ditolak")
	case errors.Is(err, services.ErrTanpaWewenang):
		galat(w, http.StatusForbidden, "wewenang tidak mencukupi")
	case errors.Is(err, services.ErrKasusPolisTertutup):
		// 409: kasus tertutup tidak dapat diputus ulang (AC 3).
		galat(w, http.StatusConflict, "kasus polis sudah ditutup")
	case errors.Is(err, models.ErrKeputusanTidakAdaDiTahapIni):
		// 409, bukan 400: keputusannya SAH, tahapnya yang tidak punya
		// jalurnya. 400 akan membuat orang mengira ia salah ketik.
		galat(w, http.StatusConflict, err.Error())
	case errors.Is(err, services.ErrPenggolongBelumSaatnya):
		galat(w, http.StatusConflict, err.Error())
	case errors.Is(err, models.ErrKeputusanTidakDikenal),
		errors.Is(err, models.ErrTahapPolisTidakDikenal),
		errors.Is(err, services.ErrPermintaanTidakSah):
		galat(w, http.StatusBadRequest, err.Error())
	default:
		galat(w, http.StatusInternalServerError, "gagal memproses keputusan polis")
	}
	return true
}
