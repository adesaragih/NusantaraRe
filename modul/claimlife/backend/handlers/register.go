package handlers

// Pintu HTTP pendaftaran klaim - tiket 02.
//
// Nol aturan dagang di sini: berkas ini menerjemahkan HTTP ke panggilan
// services dan sebaliknya. Yang memutuskan boleh atau tidak adalah services.
//
// Dibaca sesudah: handlers.go dan services/pendaftaran.go.

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/modul/claimlife/backend/models"
	"nusantarare/modul/claimlife/backend/services"
)

// permintaanDaftarJSON adalah bentuk badan permintaan pendaftaran.
//
// ⛔ Uang tidak muncul di sini sama sekali: pendaftaran hanya memilih peserta,
// dan angka klaim lahir di tiket 03. Bila kelak ada, ia TEKS (ADR-U-0003).
//
// Dua bentuk (keputusan work owner 10-10-2026, satu case per PL + mata uang): `kelompok` (layar Register baru) atau
// medan tunggal nomorPremiList/nomorPolis/mataUang/sertifikat (bentuk lama = satu kelompok). Pengelompokan dan mata
// uang klien TIDAK dipercaya - server membentuk ulang case dari baris sumber (services.KelompokkanPeserta).
// `type` dan `kodeBisnis` tidak lagi dibaca (10-10-2026): bila klien lama masih mengirimnya, decoder JSON
// mengabaikannya - Type dan kode bisnis per case dari polis PL-nya.
type permintaanDaftarJSON struct {
	NomorPremiList string               `json:"nomorPremiList"`
	NomorPolis     string               `json:"nomorPolis"`
	MataUang       string               `json:"mataUang"`
	Sertifikat     []string             `json:"sertifikat"`
	Kelompok       []kelompokDaftarJSON `json:"kelompok"`
}

// kelompokDaftarJSON - pilihan peserta satu premium list.
type kelompokDaftarJSON struct {
	NomorPremiList string   `json:"nomorPremiList"`
	NomorPolis     string   `json:"nomorPolis"`
	MataUang       string   `json:"mataUang"`
	Sertifikat     []string `json:"sertifikat"`
}

// kasusDaftarJSON - satu case yang terbentuk.
type kasusDaftarJSON struct {
	ID             string `json:"id"`
	NomorKlaim     string `json:"nomorKlaim"`
	NomorPremiList string `json:"nomorPremiList"`
	MataUang       string `json:"mataUang"`
	JumlahPeserta  int    `json:"jumlahPeserta"`
}

// keLayanan menerjemahkan badan JSON ke permintaan layanan.
func (m permintaanDaftarJSON) keLayanan() services.PermintaanDaftar {
	p := services.PermintaanDaftar{
		NomorPremiList: m.NomorPremiList,
		NomorPolis:     m.NomorPolis,
		MataUang:       m.MataUang,
		Sertifikat:     m.Sertifikat,
	}
	for _, k := range m.Kelompok {
		p.Kelompok = append(p.Kelompok, services.KelompokDaftar{NomorPremiList: k.NomorPremiList,
			NomorPolis: k.NomorPolis, MataUang: k.MataUang, Sertifikat: k.Sertifikat})
	}
	return p
}

// jawabanGalatDaftar - status dan kalimat untuk galat pendaftaran; ada=false = bukan galat.
//
// 400 = permintaan cacat bentuk (medan kosong, sertifikat ganda di dalam permintaan, batas terlampaui); 422 = pilihan
// tidak dapat menjadi case (peserta tidak ada / batal, nomor polis tidak cocok dengan PL, satu PL banyak polis);
// 409 = peserta masih aktif di klaim lain (aturan in-app 10-10-2026) - kalimat menyebut PL dan sertifikat, tanpa
// awalan teknis.
func jawabanGalatDaftar(err error) (status int, pesan string, ada bool) {
	switch {
	case err == nil:
		return 0, "", false
	case errors.Is(err, galat.ErrPermintaanTidakSah):
		return http.StatusBadRequest, err.Error(), true
	case errors.Is(err, services.ErrPesertaSudahDiklaim):
		// 409: bentrok dengan data yang SUDAH ADA (peserta aktif di klaim lain), bukan permintaan cacat bentuk.
		return http.StatusConflict, err.Error(), true
	case errors.Is(err, services.ErrKelompokTidakSah):
		return http.StatusUnprocessableEntity,
			strings.NewReplacer("services: ", "", "repository: ", "").Replace(err.Error()), true
	case errors.Is(err, inti.ErrTanpaIdentitas):
		return http.StatusUnauthorized, "permintaan tanpa identitas pelaku ditolak", true
	case errors.Is(err, inti.ErrTanpaWewenang):
		return http.StatusForbidden, "peran tidak mencukupi", true
	}
	return http.StatusInternalServerError, "gagal mendaftarkan klaim", true
}

// daftarKlaim melayani POST /api/klaim-life.
func daftarKlaim(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat.Tulis(w, http.StatusServiceUnavailable, "database belum dikonfigurasi")
			return
		}
		var masuk permintaanDaftarJSON
		if err := json.NewDecoder(r.Body).Decode(&masuk); err != nil {
			galat.Tulis(w, http.StatusBadRequest, "badan permintaan bukan JSON yang sah")
			return
		}

		// ⛔ Hanya NOMOR sertifikat yang diteruskan. Nilai polis dibaca server
		// sendiri dari sumbernya - lihat repository.AmbilUntukKlaim.
		// RALAT 10-10-2026: tanpa penomor - CLAIM_NO kosong sampai Save to RNM.
		pohon, err := svc.Pendaftaran().Daftar(r.Context(), inti.PelakuDari(r, stubPelaku), masuk.keLayanan())
		if status, pesan, ada := jawabanGalatDaftar(err); ada {
			galat.Tulis(w, status, pesan)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(jawabanDaftar(pohon))
	}
}

// jawabanDaftar - badan 201: daftar case yang terbentuk {id, nomorKlaim, nomorPremiList, mataUang, jumlahPeserta}.
func jawabanDaftar(pohon []models.PohonKlaim) any {
	keluar := struct {
		Kasus []kasusDaftarJSON `json:"kasus"`
	}{Kasus: []kasusDaftarJSON{}}
	for _, p := range pohon {
		pl := ""
		if len(p.Klaim.Peserta) > 0 {
			pl = p.Klaim.Peserta[0].NomorPremiList // seluruh peserta satu case se-PL
		}
		keluar.Kasus = append(keluar.Kasus, kasusDaftarJSON{ID: p.Work.ID, NomorKlaim: p.Klaim.NomorKlaim,
			NomorPremiList: pl, MataUang: p.Klaim.ClaimRetro.Currency, JumlahPeserta: len(p.Klaim.Peserta)})
	}
	return keluar
}

// cariPeserta melayani
// GET /api/peserta-life?pl=<PL_NUMBER>&sertifikat=<teks>&nama=<teks>&n=<batas>.
//
// `sertifikat` dan `nama` adalah kotak pencarian `Find Insured` layar Register
// (`Section/InputRegisterClaimLife.xml:16576` -> `LoadDataPesertaSpesifik_Act`).
// Keduanya boleh kosong; kosong berarti tidak menyaring.
func cariPeserta(svc *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat.Tulis(w, http.StatusServiceUnavailable, "database belum dikonfigurasi")
			return
		}
		pl := r.URL.Query().Get("pl")
		if pl == "" {
			// ⛔ 400, bukan hasil kosong: tabel peserta berisi 66,8 juta baris
			// dan hanya ber-index pada PL_NUMBER. Pencarian tanpa itu bukan
			// "pencarian luas", melainkan pemindaian penuh.
			galat.Tulis(w, http.StatusBadRequest, "parameter pl (nomor premium list) wajib diisi")
			return
		}
		batas, _ := strconv.Atoi(r.URL.Query().Get("n"))

		hasil, err := svc.Peserta().Cari(r.Context(), pl,
			r.URL.Query().Get("sertifikat"), r.URL.Query().Get("nama"), batas)
		if err != nil {
			galat.Tulis(w, http.StatusInternalServerError, "gagal mencari peserta")
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(hasil)
	}
}
