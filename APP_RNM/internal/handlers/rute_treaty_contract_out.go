package handlers

// Pintu HTTP modul Treaty Contract Out.
//
//	GET  /api/treaty-contract-out/jenis-reasuransi  jenis reasuransi non-life tersaring (tiket 02)
//	GET  /api/treaty-contract-out/grup-treaty       master grup treaty, dibaca saja (tiket 03)
//	GET  /api/treaty-contract-out/tahun             daftar tahun treaty, ID DESC (tiket 03)
//	POST /api/treaty-contract-out/tahun             tahun treaty BARU - tombol `Save` b10332 (tiket 03)
//	GET  /api/treaty-contract-out/tahun/{id}        satu tahun treaty (tiket 03)
//	PUT  /api/treaty-contract-out/tahun/{id}        perbarui SELURUH medan - `Edit` b19939 → `Save` (tiket 03)
//
// Nol aturan dagang di sini; saringan master di repository (SQL), gerbang
// tahun di models/services. Implementasi Oracle DISUNTIK di sini
// (penyuntikan_test.go).
//
// ⛔ POST dan PUT terpisah walau Pega punya satu `Save` ber-upsert: identitas
// baris baru TIDAK PERNAH datang dari klien (AC 5), dan pembaruan menimpa
// baris yang ID-nya disebut di jalur (AC 8). Badan PUT yang membawa `id`
// berbeda dari jalurnya DITOLAK, bukan salah satu dipilih diam-diam.
//
// ⚠️ Berkas ini MILIK sesi Treaty Contract Out. `handlers.go` disentuh hanya
// dengan SATU baris pemanggilan `daftarkanRuteTreatyContractOut`.

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"nusantarare/internal/models"
	"nusantarare/internal/services"
)

// daftarkanRuteTreatyContractOut mendaftarkan seluruh rute modul ini.
func daftarkanRuteTreatyContractOut(mux *http.ServeMux, svc *services.Service, stubPelaku bool) {
	mux.HandleFunc("GET /api/treaty-contract-out/jenis-reasuransi",
		jenisReasuransiTreaty(svc, stubPelaku))
	mux.HandleFunc("GET /api/treaty-contract-out/grup-treaty",
		grupTreaty(svc, stubPelaku))
	mux.HandleFunc("GET /api/treaty-contract-out/tahun",
		daftarTahunTreaty(svc, stubPelaku))
	mux.HandleFunc("POST /api/treaty-contract-out/tahun",
		simpanTahunTreaty(svc, stubPelaku, false))
	mux.HandleFunc("GET /api/treaty-contract-out/tahun/{id}",
		satuTahunTreaty(svc, stubPelaku))
	mux.HandleFunc("PUT /api/treaty-contract-out/tahun/{id}",
		simpanTahunTreaty(svc, stubPelaku, true))
	// Tiket 12: lampiran tahun treaty (tco_lampiran.go).
	daftarkanRuteLampiranTCO(mux, svc, stubPelaku)
	// Tiket 04: kontrak treaty di dalam tahun (tco_kontrak.go).
	daftarkanRuteKontrakTCO(mux, svc, stubPelaku)
}

// jawabanDaftarJenisReasuransi adalah badan jawaban daftar jenis reasuransi.
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

// jawabanDaftarGrupTreaty adalah badan jawaban daftar grup treaty.
type jawabanDaftarGrupTreaty struct {
	Daftar []services.GrupTreaty `json:"daftar"`
	Total  int                   `json:"total"`
}

// grupTreaty melayani GET /api/treaty-contract-out/grup-treaty.
func grupTreaty(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat(w, http.StatusServiceUnavailable, "database belum dikonfigurasi")
			return
		}
		daftar, err := svc.GrupTreaty().
			DenganPembaca(services.PembacaGrupTreatyOracle(svc)).
			Daftar(r.Context(), pelakuDari(r, stubPelaku))
		if jawabGalatTreatyContractOut(w, err) {
			return
		}
		tulisJSONPolis(w, jawabanDaftarGrupTreaty{Daftar: daftar, Total: len(daftar)})
	}
}

// daftarTahunTreaty melayani GET /api/treaty-contract-out/tahun.
func daftarTahunTreaty(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat(w, http.StatusServiceUnavailable, "database belum dikonfigurasi")
			return
		}
		halaman, _ := strconv.Atoi(r.URL.Query().Get("halaman"))
		ukuran, _ := strconv.Atoi(r.URL.Query().Get("ukuran"))
		hal, err := svc.TahunTreatyTCO().
			DenganGudang(services.GudangTahunTreatyOracle(svc)).
			Daftar(r.Context(), pelakuDari(r, stubPelaku), halaman, ukuran)
		if jawabGalatTreatyContractOut(w, err) {
			return
		}
		tulisJSONPolis(w, hal)
	}
}

// satuTahunTreaty melayani GET /api/treaty-contract-out/tahun/{id}.
func satuTahunTreaty(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat(w, http.StatusServiceUnavailable, "database belum dikonfigurasi")
			return
		}
		t, err := svc.TahunTreatyTCO().
			DenganGudang(services.GudangTahunTreatyOracle(svc)).
			Ambil(r.Context(), pelakuDari(r, stubPelaku), r.PathValue("id"))
		if jawabGalatTreatyContractOut(w, err) {
			return
		}
		tulisJSONPolis(w, t)
	}
}

// simpanTahunTreaty melayani POST (baru) dan PUT /{id} (perbarui).
func simpanTahunTreaty(svc *services.Service, stubPelaku bool, perbarui bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat(w, http.StatusServiceUnavailable, "database belum dikonfigurasi")
			return
		}
		var masuk services.TahunTreatyMasuk
		if err := json.NewDecoder(r.Body).Decode(&masuk); err != nil {
			galat(w, http.StatusBadRequest, "badan permintaan bukan JSON yang sah")
			return
		}
		if perbarui {
			id := r.PathValue("id")
			if masuk.ID != "" && masuk.ID != id {
				galat(w, http.StatusBadRequest, "id di badan berbeda dari id di jalur")
				return
			}
			masuk.ID = id
		} else if masuk.ID != "" {
			// AC 5: identitas baris baru tidak pernah diketik pengguna.
			galat(w, http.StatusBadRequest, "tahun treaty baru tidak membawa id; identitas dibuat server")
			return
		}
		hasil, err := svc.TahunTreatyTCO().
			DenganGudang(services.GudangTahunTreatyOracle(svc)).
			Simpan(r.Context(), pelakuDari(r, stubPelaku), masuk)
		if jawabGalatTreatyContractOut(w, err) {
			return
		}
		tulisJSONPolis(w, hasil)
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
	case errors.Is(err, services.ErrMasterJenisReasuransiKosong),
		errors.Is(err, services.ErrMasterGrupTreatyKosong):
		// ⛔ 503, dan pesannya MENYEBUT MASTERNYA: keadaan server yang belum
		// siap - master rujukan kosong atau tersaring habis - bukan
		// permintaan yang salah, dan bukan daftar kosong yang diam (ADR-0015).
		galat(w, http.StatusServiceUnavailable, err.Error())
	case errors.Is(err, services.ErrKategoriLampiranKosong),
		errors.Is(err, services.ErrUnggahanDirBelumDisetel):
		// Tiket 12: keadaan server - master kategori kosong atau folder
		// unggahan belum disetel - bukan salah pemanggil.
		galat(w, http.StatusServiceUnavailable, err.Error())
	case errors.Is(err, services.ErrTahunTreatyTidakAda):
		galat(w, http.StatusNotFound, "tahun treaty tidak ditemukan")
	case errors.Is(err, services.ErrKontrakTidakAda):
		galat(w, http.StatusNotFound, "kontrak treaty tidak ditemukan pada tahun treaty ini")
	case errors.Is(err, services.ErrKontrakDobel):
		// 409: "Data sudah pernah di Input" + kontrak mana yang memegang jenisnya.
		galat(w, http.StatusConflict, err.Error())
	case errors.Is(err, models.ErrKontrakJenisReasuransiKosong),
		errors.Is(err, models.ErrKontrakMulaiKosong),
		errors.Is(err, models.ErrKontrakAkhirKosong),
		errors.Is(err, models.ErrKontrakTahunMulaiBeda),
		errors.Is(err, services.ErrJenisReasuransiDiLuarDaftar):
		galat(w, http.StatusUnprocessableEntity, err.Error())
	case errors.Is(err, services.ErrLampiranTidakAda):
		galat(w, http.StatusNotFound, "lampiran tidak ditemukan pada tahun treaty ini")
	case errors.Is(err, services.ErrLampiranBelumTerkirim),
		errors.Is(err, services.ErrLampiranSudahTerkirim),
		errors.Is(err, services.ErrLampiranTanpaBerkas),
		errors.Is(err, services.ErrBerkasSumberLampiranHilang):
		// 409: keadaan DATA lampiran; pesannya menyebut lampiran mana dan
		// perbaikannya.
		galat(w, http.StatusConflict, err.Error())
	case errors.Is(err, services.ErrBerkasTerlaluBesar):
		galat(w, http.StatusRequestEntityTooLarge, err.Error())
	case errors.Is(err, services.ErrBerkasKosong):
		galat(w, http.StatusBadRequest, services.PesanTanpaBerkasTCO)
	case errors.Is(err, services.ErrKategoriLampiranTidakDikenal):
		galat(w, http.StatusUnprocessableEntity, err.Error())
	case errors.Is(err, services.ErrTahunTreatyDobel):
		// 409: keadaan DATA - kombinasi periode + grup sudah dipakai baris
		// lain - dan pesannya menyebut baris mana (AC 73).
		galat(w, http.StatusConflict, err.Error())
	case errors.Is(err, models.ErrPeriodeTerbalik),
		errors.Is(err, models.ErrTahunTreatyGrupKosong),
		errors.Is(err, models.ErrTahunTreatyTahunKosong),
		errors.Is(err, models.ErrTahunTreatyBukanAngka):
		// 422: JSON-nya sah, isinya yang ditolak gerbang - pesan menyebut medannya.
		galat(w, http.StatusUnprocessableEntity, err.Error())
	case errors.Is(err, services.ErrPermintaanTidakSah):
		galat(w, http.StatusBadRequest, err.Error())
	default:
		galat(w, http.StatusInternalServerError, "gagal memproses permintaan treaty contract out")
	}
	return true
}
