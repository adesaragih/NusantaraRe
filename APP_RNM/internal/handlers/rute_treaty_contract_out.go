package handlers

// Pintu HTTP modul Treaty Contract Out.
//
//	GET /api/treaty-contract-out/jenis-reasuransi   daftar jenis reasuransi
//	                                                 non-life tersaring (tiket 02)
//
// Nol aturan dagang di sini; saringannya di repository (SQL), kegagalan
// master kosong di services.
//
// ⚠️ Berkas ini MILIK sesi Treaty Contract Out. `handlers.go` disentuh hanya
// dengan SATU baris pemanggilan `daftarkanRuteTreatyContractOut` - perubahan
// bersama yang aditif, dan dilaporkan.

import (
	"errors"
	"net/http"

	"nusantarare/internal/services"
)

// daftarkanRuteTreatyContractOut mendaftarkan seluruh rute modul ini.
func daftarkanRuteTreatyContractOut(mux *http.ServeMux, svc *services.Service, stubPelaku bool) {
	// Tiket 02 - master dibaca, disaring persis RD dominan (12 awalan
	// blacklist, Flag active, Type 1/2/3). GET: ia MEMBACA.
	mux.HandleFunc("GET /api/treaty-contract-out/jenis-reasuransi",
		jenisReasuransiTreaty(svc, stubPelaku))
}

// jawabanDaftarJenisReasuransi adalah badan jawaban daftar.
type jawabanDaftarJenisReasuransi struct {
	Daftar []services.JenisReasuransi `json:"daftar"`
	Total  int                        `json:"total"`
}

// jenisReasuransiTreaty melayani GET /api/treaty-contract-out/jenis-reasuransi.
func jenisReasuransiTreaty(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat(w, http.StatusServiceUnavailable, "database belum dikonfigurasi")
			return
		}
		daftar, err := svc.JenisReasuransiTreaty().
			DenganPembaca(services.PembacaJenisReasuransiOracle(svc)).
			Daftar(r.Context(), pelakuDari(r, stubPelaku))
		if jawabGalatTreatyContractOut(w, err) {
			return
		}
		tulisJSONPolis(w, jawabanDaftarJenisReasuransi{Daftar: daftar, Total: len(daftar)})
	}
}

// jawabGalatTreatyContractOut menerjemahkan galat services menjadi kode HTTP.
//
// Mengembalikan true bila permintaan SUDAH dijawab.
func jawabGalatTreatyContractOut(w http.ResponseWriter, err error) bool {
	switch {
	case err == nil:
		return false
	case errors.Is(err, services.ErrTanpaIdentitas):
		galat(w, http.StatusUnauthorized, "permintaan tanpa identitas pelaku ditolak")
	case errors.Is(err, services.ErrTanpaWewenang):
		galat(w, http.StatusForbidden, "wewenang tidak mencukupi")
	case errors.Is(err, services.ErrMasterJenisReasuransiKosong):
		// ⛔ 503, dan pesannya MENYEBUT MASTERNYA: keadaan server yang belum
		// siap - master rujukan kosong atau tersaring habis - bukan
		// permintaan yang salah, dan bukan daftar kosong yang diam (ADR-0015).
		galat(w, http.StatusServiceUnavailable, err.Error())
	case errors.Is(err, services.ErrPermintaanTidakSah):
		galat(w, http.StatusBadRequest, err.Error())
	default:
		galat(w, http.StatusInternalServerError, "gagal memproses permintaan treaty contract out")
	}
	return true
}
