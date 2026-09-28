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
	"strconv"
	"time"

	"nusantarare/internal/models"
	"nusantarare/internal/services"
)

// kotakMasukPolis melayani GET /api/polis-life.
//
// ⛔ GET: ia MEMBACA. Kotak masuk yang mengubah sesuatu adalah kotak masuk
// yang berubah karena seseorang menyegarkan halamannya.
func kotakMasukPolis(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat(w, http.StatusServiceUnavailable, "database belum dikonfigurasi")
			return
		}
		ukuran, _ := strconv.Atoi(r.URL.Query().Get("ukuran"))
		// ⚠️ Halaman yang tidak terbaca menjadi 0, dan services menjepitnya
		// ke 1 - bukan menjadi galat. Daftar yang ditolak karena satu
		// parameter salah ketik lebih menjengkelkan daripada berguna.
		halaman, _ := strconv.Atoi(r.URL.Query().Get("halaman"))
		hal, err := svc.InboxPolis().Ambil(r.Context(), pelakuDari(r, stubPelaku),
			r.URL.Query().Get("posisi"), halaman, ukuran)
		if jawabGalatPolis(w, err) {
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(hal)
	}
}

// jawabanPeriode adalah periode produksi yang berlaku saat ini.
//
// ⛔ Ia DITAMPILKAN sebelum pemakai menyimpan (AC tiket 02), bukan
// tersimpan diam-diam. Transaksi yang mendarat di bulan yang salah karena
// seseorang menyimpannya lewat tengah malam adalah kekeliruan yang hanya
// dapat dicegah dengan menunjukkannya lebih dulu.
type jawabanPeriode struct {
	Periode string `json:"periode"`
}

// periodeProduksi melayani GET /api/polis-life/periode.
func periodeProduksi(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat(w, http.StatusServiceUnavailable, "database belum dikonfigurasi")
			return
		}
		p, err := svc.Periode().Sekarang(r.Context(), pelakuDari(r, stubPelaku))
		switch {
		case err == nil:
		case errors.Is(err, models.ErrTanggalTutupBukuKosong),
			errors.Is(err, models.ErrTanggalTutupBukuTidakMasukAkal):
			// ⛔ 503, dan pesannya MENYEBUT TABEL SUMBERNYA. Ia keadaan
			// server yang belum siap - tabel rujukan yang kosong - bukan
			// permintaan yang salah.
			galat(w, http.StatusServiceUnavailable, err.Error())
			return
		default:
			if jawabGalatPolis(w, err) {
				return
			}
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(jawabanPeriode{Periode: models.PeriodeTeks(p)})
	}
}

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
