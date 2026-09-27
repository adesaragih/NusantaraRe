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

	"nusantarare/internal/services"
)

// permintaanDaftarJSON adalah bentuk badan permintaan pendaftaran.
//
// ⛔ Uang tidak muncul di sini sama sekali: pendaftaran hanya memilih peserta,
// dan angka klaim lahir di tiket 03. Bila kelak ada, ia TEKS (ADR-U-0003).
type permintaanDaftarJSON struct {
	NomorPremiList string   `json:"nomorPremiList"`
	NomorPolis     string   `json:"nomorPolis"`
	Type           string   `json:"type"`
	KodeBisnis     string   `json:"kodeBisnis"`
	MataUang       string   `json:"mataUang"`
	Sertifikat     []string `json:"sertifikat"`
}

// daftarKlaim melayani POST /api/klaim-life.
func daftarKlaim(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat(w, http.StatusServiceUnavailable, "database belum dikonfigurasi")
			return
		}
		var masuk permintaanDaftarJSON
		if err := json.NewDecoder(r.Body).Decode(&masuk); err != nil {
			galat(w, http.StatusBadRequest, "badan permintaan bukan JSON yang sah")
			return
		}

		// ⛔ Hanya NOMOR sertifikat yang diteruskan. Nilai polis dibaca server
		// sendiri dari sumbernya - lihat repository.AmbilUntukKlaim.
		pohon, err := svc.Pendaftaran().DenganPenomor(services.PenomorCounterOracle(svc)).Daftar(r.Context(), pelakuDari(r, stubPelaku), services.PermintaanDaftar{
			NomorPremiList: masuk.NomorPremiList,
			NomorPolis:     masuk.NomorPolis,
			Type:           masuk.Type,
			KodeBisnis:     masuk.KodeBisnis,
			MataUang:       masuk.MataUang,
			Sertifikat:     masuk.Sertifikat,
		})
		switch {
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
			galat(w, http.StatusInternalServerError, "gagal mendaftarkan klaim")
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(struct {
			ID         string `json:"id"`
			NomorKlaim string `json:"nomorKlaim"`
		}{pohon.Work.ID, pohon.Klaim.NomorKlaim})
	}
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
			galat(w, http.StatusServiceUnavailable, "database belum dikonfigurasi")
			return
		}
		pl := r.URL.Query().Get("pl")
		if pl == "" {
			// ⛔ 400, bukan hasil kosong: tabel peserta berisi 66,8 juta baris
			// dan hanya ber-index pada PL_NUMBER. Pencarian tanpa itu bukan
			// "pencarian luas", melainkan pemindaian penuh.
			galat(w, http.StatusBadRequest, "parameter pl (nomor premium list) wajib diisi")
			return
		}
		batas, _ := strconv.Atoi(r.URL.Query().Get("n"))

		hasil, err := svc.Peserta().Cari(r.Context(), pl,
			r.URL.Query().Get("sertifikat"), r.URL.Query().Get("nama"), batas)
		if err != nil {
			galat(w, http.StatusInternalServerError, "gagal mencari peserta")
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(hasil)
	}
}
