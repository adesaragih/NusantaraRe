package handlers

// Pintu HTTP modul PremiumList Life - tiket 01.
//
//	POST /api/polis-life/{id}/keputusan   body {"keputusan":"Confirm|Reject|Decline"}
//	POST /api/polis-life/{id}/penggolong  body {"hasil":"Offer|Premium"}
//	GET  /api/polis-life/{id}/summary     rekap per mata uang (tiket 05a)
//	POST /api/polis-life/{id}/summary     nomor + rekap + warisan (tiket 05a)
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

// ringkasPolis melayani GET /api/polis-life/ringkas?nomorPolis=...
//
// ⛔ BERKUNCI NOMOR POLIS, bukan id kasus. Claim Life mengenal polis lewat
// NOMORNYA - itu yang diketik pemakai di `Choose Policy No` - dan ia tidak
// pernah tahu id kasus PremiumList.
//
// ⚠️ Rute ini BERDIRI SEBELUM `GET /api/polis-life/{id}` di tabel rute, dan
// urutan itu tidak menentukan apa pun di `ServeMux` Go 1.22: pola yang lebih
// SPESIFIK menang, dan `ringkas` harfiah lebih spesifik daripada `{id}`.
// Dicatat supaya tidak ada yang "memperbaikinya" dengan menambah pemeriksaan
// `id == "ringkas"`.
func ringkasPolis(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat(w, http.StatusServiceUnavailable, "database belum dikonfigurasi")
			return
		}
		hasil, err := svc.RingkasPolis().Ambil(r.Context(),
			pelakuDari(r, stubPelaku), r.URL.Query().Get("nomorPolis"))
		if jawabGalatPolis(w, err) {
			return
		}
		tulisJSONPolis(w, hasil)
	}
}

// kepalaPolis melayani GET /api/polis-life/{id}.
//
// ⛔ GET: ia MEMBACA. Nomor PL TIDAK terbit di sini - membuka layar detail
// tidak boleh menggerakkan penghitung, dan rute yang menerbitkan nomor saat
// seseorang menyegarkan halaman akan membakar nomor tiap kali.
func kepalaPolis(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat(w, http.StatusServiceUnavailable, "database belum dikonfigurasi")
			return
		}
		kepala, err := svc.DetailPolis().Kepala(r.Context(),
			pelakuDari(r, stubPelaku), r.PathValue("id"))
		if jawabGalatPolis(w, err) {
			return
		}
		tulisJSONPolis(w, kepala)
	}
}

// pesertaPolis melayani GET /api/polis-life/{id}/peserta.
func pesertaPolis(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat(w, http.StatusServiceUnavailable, "database belum dikonfigurasi")
			return
		}
		halaman, _ := strconv.Atoi(r.URL.Query().Get("halaman"))
		ukuran, _ := strconv.Atoi(r.URL.Query().Get("ukuran"))
		hal, err := svc.DetailPolis().Peserta(r.Context(),
			pelakuDari(r, stubPelaku), r.PathValue("id"), halaman, ukuran)
		if jawabGalatPolis(w, err) {
			return
		}
		tulisJSONPolis(w, hal)
	}
}

// terbitkanNomorPolis melayani POST /api/polis-life/{id}/nomor.
//
// ⛔ POST, dan TANPA badan permintaan. Tidak ada satu pun bahan nomor yang
// boleh datang dari klien: awalan, tipe, kode bisnis, periode, dan urut
// seluruhnya dibaca server dari sumbernya masing-masing. Nomor yang bahannya
// dapat disebut pemanggil adalah nomor yang dapat dipilih pemanggil.
func terbitkanNomorPolis(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat(w, http.StatusServiceUnavailable, "database belum dikonfigurasi")
			return
		}
		hasil, err := svc.NomorPremiumList().Terbitkan(r.Context(),
			pelakuDari(r, stubPelaku), r.PathValue("id"))
		if jawabGalatPolis(w, err) {
			return
		}
		// ⚠️ 200, bukan 201, bahkan saat nomornya baru lahir. Yang dibuat
		// bukan sumber daya baru di alamat baru - ia medan pada polis yang
		// sudah ada, dan alamatnya tetap sama sesudahnya.
		tulisJSONPolis(w, hasil)
	}
}

// rekapPolis melayani GET /api/polis-life/{id}/summary - tiket 05a.
//
// Layar `ShowLifePremiumSummary`: rekap per mata uang dihitung dari peserta,
// TANPA menyimpan apa pun.
func rekapPolis(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat(w, http.StatusServiceUnavailable, "database belum dikonfigurasi")
			return
		}
		hasil, err := svc.SummaryPremiumList().Lihat(r.Context(),
			pelakuDari(r, stubPelaku), r.PathValue("id"))
		if jawabGalatPolis(w, err) {
			return
		}
		tulisJSONPolis(w, hasil)
	}
}

// submitRekapPolis melayani POST /api/polis-life/{id}/summary - tiket 05a.
//
// Padanan tombol `Submit` (`ShowLifePremiumSummary` b27471) sisi penyimpanan:
// nomor + rekap + salinan peserta warisan dalam SATU transaksi.
//
// ⛔ TANPA badan permintaan, dengan alasan yang sama dengan penomoran: tidak
// ada satu pun angka rekap yang boleh datang dari klien. Rekap yang dikirim
// layar adalah rekap yang dapat diketik ulang orang.
//
// ⚠️ Tiket 05b: sesudah tersimpan, kasus DITUTUP Resolved-Completed -
// `finishAssignment` b26442 → `Transition2` → `END52`. Hanya dari tahap
// Input Premium Summary; tahap lain dijawab 409.
func submitRekapPolis(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat(w, http.StatusServiceUnavailable, "database belum dikonfigurasi")
			return
		}
		hasil, err := svc.SummaryPremiumList().
			DenganJejak(services.PerekamJejakOracle(svc)).
			Submit(r.Context(), pelakuDari(r, stubPelaku), r.PathValue("id"), time.Now())
		if jawabGalatPolis(w, err) {
			return
		}
		tulisJSONPolis(w, hasil)
	}
}

// tulisJSONPolis menulis satu jawaban JSON 200.
func tulisJSONPolis(w http.ResponseWriter, isi any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(isi)
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
	case errors.Is(err, services.ErrPolisNomorTakDitemukan):
		// 404: nomor polisnya memang tidak ada di PremiumList Life.
		galat(w, http.StatusNotFound, err.Error())
	case errors.Is(err, services.ErrPolisTakDitemukan):
		// 404: polisnya memang tidak ada. 500 akan membuat orang mencari
		// kerusakan di server padahal id-nya yang salah.
		galat(w, http.StatusNotFound, "polis tidak ditemukan")
	case errors.Is(err, services.ErrPolisTanpaPeserta):
		// 409: permintaannya sah, keadaan polisnya yang belum siap - dan
		// pesannya menyebut apa yang harus dikerjakan lebih dahulu.
		galat(w, http.StatusConflict, err.Error())
	case errors.Is(err, services.ErrNomorPLTerbitBersamaan):
		// 409: permintaan lain mendahului. Penghitungnya TIDAK bergerak -
		// transaksinya batal - jadi yang perlu dikerjakan pemanggil hanya
		// membaca ulang nomornya.
		galat(w, http.StatusConflict, err.Error())
	case errors.Is(err, services.ErrRekapKosong),
		errors.Is(err, services.ErrSubmitBukanTahapSummary):
		// 409: keadaan DATA polis yang belum siap, bukan permintaan yang salah.
		galat(w, http.StatusConflict, err.Error())
	case errors.Is(err, services.ErrNomorPLBerbedaAntarPeserta):
		// 409: data yang tidak sepakat dengan dirinya sendiri. 500 akan
		// menyembunyikan bahwa yang rusak adalah barisnya, bukan kodenya.
		galat(w, http.StatusConflict, err.Error())
	case errors.Is(err, models.ErrTipePLTanpaCabang),
		errors.Is(err, models.ErrKodeBisnisKosong),
		errors.Is(err, models.ErrAwalanProduksiKosong),
		errors.Is(err, models.ErrPeriodeNomorPLTakBerbentuk):
		// 409: seluruhnya bahan nomor yang belum lengkap di DATA, bukan di
		// permintaan. 400 akan menyalahkan pemanggil atas kolom yang kosong.
		galat(w, http.StatusConflict, err.Error())
	case errors.Is(err, models.ErrKeputusanTidakDikenal),
		errors.Is(err, models.ErrTahapPolisTidakDikenal),
		errors.Is(err, services.ErrPermintaanTidakSah):
		galat(w, http.StatusBadRequest, err.Error())
	default:
		// ⚠️ Kalimatnya NETRAL sejak tiket 03. Sebelumnya ia berbunyi "gagal
		// memproses keputusan polis" - benar saat hanya dua rute keputusan
		// memakainya, menyesatkan sejak rute detail dan penomoran ikut.
		galat(w, http.StatusInternalServerError, "gagal memproses permintaan polis")
	}
	return true
}
