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
	"log"
	"net/http"
	"strconv"

	"nusantarare/inti"
	"nusantarare/inti/galat"
	"nusantarare/inti/unggah"
	"nusantarare/modul/treatycontractout/models"
	"nusantarare/modul/treatycontractout/services"
)

// Router menyusun rute modul ini SAJA, di mux sendiri.
//
// Refactor bentuk B (30-09-2026): dipakai uji HTTP modul ini, yang dulu
// memakai `Router` bersama milik seluruh aplikasi. Produksi tidak memakainya:
// `cmd/api` mendaftarkan `DaftarkanRute` ke mux yang sama dengan modul lain.
func Router(svc *services.Service, stubPelaku bool) http.Handler {
	mux := http.NewServeMux()
	DaftarkanRute(mux, svc, stubPelaku)
	return mux
}

// DaftarkanRute mendaftarkan seluruh rute modul ini.
//
// Dulu `daftarkanRuteTreatyContractOut`, dipanggil `internal/handlers.Router`;
// kini dipanggil `cmd/api` (refactor bentuk B).
func DaftarkanRute(mux *http.ServeMux, svc *services.Service, stubPelaku bool) {
	mux.HandleFunc("GET /api/treaty-contract-out/jenis-reasuransi",
		jenisReasuransiTreaty(svc, stubPelaku))
	mux.HandleFunc("GET /api/treaty-contract-out/grup-treaty",
		grupTreaty(svc, stubPelaku))
	mux.HandleFunc("GET /api/treaty-contract-out/tahun",
		daftarTahunTreaty(svc, stubPelaku))
	mux.HandleFunc("POST /api/treaty-contract-out/tahun",
		simpanTahunTreaty(svc, stubPelaku, false))
	// End Date bawaan tahun BARU (belum ber-ID) - keputusan work owner 30-09-2026.
	mux.HandleFunc("GET /api/treaty-contract-out/tahun/akhir-bawaan",
		akhirBawaanTahunTCO(stubPelaku))
	mux.HandleFunc("GET /api/treaty-contract-out/tahun/{id}",
		satuTahunTreaty(svc, stubPelaku))
	mux.HandleFunc("PUT /api/treaty-contract-out/tahun/{id}",
		simpanTahunTreaty(svc, stubPelaku, true))
	// Tiket 12: lampiran tahun treaty (tco_lampiran.go).
	daftarkanRuteLampiranTCO(mux, svc, stubPelaku)
	// Tiket 04: kontrak treaty di dalam tahun (tco_kontrak.go).
	daftarkanRuteKontrakTCO(mux, svc, stubPelaku)
	// Tiket 05: reinsurer pada kombinasi kontrak (tco_reinsurer.go).
	daftarkanRuteReinsurerTCO(mux, svc, stubPelaku)
	// Tiket 07: business pada kombinasi kontrak (tco_business.go).
	daftarkanRuteBusinessTCO(mux, svc, stubPelaku)
	// Tiket 08: klausul - satu tabel, 25 jenis (tco_klausul.go).
	daftarkanRuteKlausulTCO(mux, svc, stubPelaku)
	// Tiket 06: security di bawah reinsurer (tco_security.go).
	daftarkanRuteSecurityTCO(mux, svc, stubPelaku)
	// Tiket 11: kurs USD -> IDR (tco_kurs.go).
	daftarkanRuteKursTCO(mux, svc, stubPelaku)
	// Tiket 10: kaskade hapus kontrak/reinsurer + popup (tco_kaskade.go).
	daftarkanRuteKaskadeTCO(mux, svc, stubPelaku)
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
			galat.Tulis(w, http.StatusServiceUnavailable, "database is not configured")
			return
		}
		daftar, err := svc.JenisReasuransiTreaty().
			DenganPembaca(services.PembacaJenisReasuransiOracle(svc)).
			Daftar(r.Context(), inti.PelakuDari(r, stubPelaku))
		if jawabGalatTreatyContractOut(w, err) {
			return
		}
		galat.TulisJSON(w, jawabanDaftarJenisReasuransi{Daftar: daftar, Total: len(daftar)})
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
			galat.Tulis(w, http.StatusServiceUnavailable, "database is not configured")
			return
		}
		daftar, err := svc.GrupTreaty().
			DenganPembaca(services.PembacaGrupTreatyOracle(svc)).
			Daftar(r.Context(), inti.PelakuDari(r, stubPelaku))
		if jawabGalatTreatyContractOut(w, err) {
			return
		}
		galat.TulisJSON(w, jawabanDaftarGrupTreaty{Daftar: daftar, Total: len(daftar)})
	}
}

// daftarTahunTreaty melayani GET /api/treaty-contract-out/tahun.
func daftarTahunTreaty(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat.Tulis(w, http.StatusServiceUnavailable, "database is not configured")
			return
		}
		halaman, _ := strconv.Atoi(r.URL.Query().Get("halaman"))
		ukuran, _ := strconv.Atoi(r.URL.Query().Get("ukuran"))
		hal, err := svc.TahunTreatyTCO().
			DenganGudang(services.GudangTahunTreatyOracle(svc)).
			DenganGrup(services.PembacaGrupTreatyOracle(svc)).
			Daftar(r.Context(), inti.PelakuDari(r, stubPelaku), halaman, ukuran)
		if jawabGalatTreatyContractOut(w, err) {
			return
		}
		galat.TulisJSON(w, hal)
	}
}

// satuTahunTreaty melayani GET /api/treaty-contract-out/tahun/{id}.
func satuTahunTreaty(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat.Tulis(w, http.StatusServiceUnavailable, "database is not configured")
			return
		}
		t, err := svc.TahunTreatyTCO().
			DenganGudang(services.GudangTahunTreatyOracle(svc)).
			DenganGrup(services.PembacaGrupTreatyOracle(svc)).
			Ambil(r.Context(), inti.PelakuDari(r, stubPelaku), r.PathValue("id"))
		if jawabGalatTreatyContractOut(w, err) {
			return
		}
		galat.TulisJSON(w, t)
	}
}

// simpanTahunTreaty melayani POST (baru) dan PUT /{id} (perbarui).
func simpanTahunTreaty(svc *services.Service, stubPelaku bool, perbarui bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat.Tulis(w, http.StatusServiceUnavailable, "database is not configured")
			return
		}
		var masuk services.TahunTreatyMasuk
		if err := json.NewDecoder(r.Body).Decode(&masuk); err != nil {
			galat.Tulis(w, http.StatusBadRequest, "request body is not valid JSON")
			return
		}
		if perbarui {
			id := r.PathValue("id")
			if masuk.ID != "" && masuk.ID != id {
				galat.Tulis(w, http.StatusBadRequest, "id in the body differs from id in the path")
				return
			}
			masuk.ID = id
		} else if masuk.ID != "" {
			// AC 5: identitas baris baru tidak pernah diketik pengguna.
			galat.Tulis(w, http.StatusBadRequest, "a new treaty year must not carry an id; the server assigns it")
			return
		}
		hasil, err := svc.TahunTreatyTCO().
			DenganGudang(services.GudangTahunTreatyOracle(svc)).
			DenganGrup(services.PembacaGrupTreatyOracle(svc)).
			Simpan(r.Context(), inti.PelakuDari(r, stubPelaku), masuk)
		if jawabGalatTreatyContractOut(w, err) {
			return
		}
		galat.TulisJSON(w, hasil)
	}
}

// jawabGalatTreatyContractOut menerjemahkan galat services menjadi kode HTTP.
//
// Mengembalikan true bila permintaan SUDAH dijawab.
func jawabGalatTreatyContractOut(w http.ResponseWriter, err error) bool {
	switch {
	case err == nil:
		return false
	case errors.Is(err, inti.ErrTanpaIdentitas):
		galat.Tulis(w, http.StatusUnauthorized, "request without user identity is rejected")
	case errors.Is(err, inti.ErrTanpaWewenang):
		galat.Tulis(w, http.StatusForbidden, "insufficient permission")
	case errors.Is(err, services.ErrMasterJenisReasuransiKosong),
		errors.Is(err, services.ErrMasterGrupTreatyKosong):
		// ⛔ 503, dan pesannya MENYEBUT MASTERNYA: keadaan server yang belum
		// siap - master rujukan kosong atau tersaring habis - bukan
		// permintaan yang salah, dan bukan daftar kosong yang diam (ADR-0015).
		galat.Tulis(w, http.StatusServiceUnavailable, pesanTCO(err))
	case errors.Is(err, services.ErrKategoriLampiranKosong),
		errors.Is(err, unggah.ErrUnggahanDirBelumDisetel):
		// Tiket 12: keadaan server - master kategori kosong atau folder
		// unggahan belum disetel - bukan salah pemanggil.
		galat.Tulis(w, http.StatusServiceUnavailable, pesanTCO(err))
	case errors.Is(err, services.ErrTahunTreatyTidakAda):
		galat.Tulis(w, http.StatusNotFound, "treaty year not found")
	case errors.Is(err, services.ErrKontrakTidakAda):
		galat.Tulis(w, http.StatusNotFound, "treaty contract not found in this treaty year")
	case errors.Is(err, services.ErrKontrakDobel):
		// 409: "Data sudah pernah di Input" + kontrak mana yang memegang jenisnya.
		galat.Tulis(w, http.StatusConflict, pesanTCO(err))
	case errors.Is(err, models.ErrKontrakJenisReasuransiKosong),
		errors.Is(err, models.ErrKontrakMulaiKosong),
		errors.Is(err, models.ErrKontrakAkhirKosong),
		errors.Is(err, models.ErrKontrakTahunMulaiBeda),
		errors.Is(err, services.ErrJenisReasuransiDiLuarDaftar):
		galat.Tulis(w, http.StatusUnprocessableEntity, pesanTCO(err))
	case errors.Is(err, services.ErrReinsurerTidakAda):
		galat.Tulis(w, http.StatusNotFound, "reinsurer not found in this contract combination")
	case errors.Is(err, models.ErrReinsurerKosong),
		errors.Is(err, models.ErrPersenKosong),
		errors.Is(err, models.ErrPersenBukanDesimal),
		errors.Is(err, models.ErrPersenDiLuarRentang),
		errors.Is(err, models.ErrTotalShareMelebihi100),
		errors.Is(err, services.ErrReinsurerDiLuarMaster):
		galat.Tulis(w, http.StatusUnprocessableEntity, pesanTCO(err))
	case errors.Is(err, services.ErrGrupTreatyDiLuarMaster):
		galat.Tulis(w, http.StatusUnprocessableEntity, pesanTCO(err))
	case errors.Is(err, services.ErrKontrakBeranak),
		errors.Is(err, services.ErrTahunBeranak):
		// 409: kombinasi beranak tidak boleh diganti kuncinya (temuan /code-review).
		galat.Tulis(w, http.StatusConflict, pesanTCO(err))
	case errors.Is(err, services.ErrDampakBerubah),
		errors.Is(err, services.ErrKaskadeTidakUtuh):
		// 409: keadaan DATA berubah sejak popup - tinjau ulang, tidak ada yang terhapus.
		galat.Tulis(w, http.StatusConflict, pesanTCO(err))
	case errors.Is(err, services.ErrKursTidakAda):
		// 422 + pesan `NewTreatyArrEpi.xml` b870, diterjemahkan [keputusan work owner 30-09-2026: bahasa Inggris] (ADR-0015).
		galat.Tulis(w, http.StatusUnprocessableEntity, pesanTCO(err))
	case errors.Is(err, services.ErrMasterKursRusak):
		// 503: master kurs / mata uang tidak dapat dipakai - keadaan server.
		galat.Tulis(w, http.StatusServiceUnavailable, pesanTCO(err))
	case errors.Is(err, services.ErrKonversiKursTidakSah):
		galat.Tulis(w, http.StatusBadRequest, pesanTCO(err))
	case errors.Is(err, services.ErrSecurityTidakAda):
		galat.Tulis(w, http.StatusNotFound, "security not found for this reinsurer")
	case errors.Is(err, services.ErrSecurityDobel):
		galat.Tulis(w, http.StatusConflict, pesanTCO(err))
	case errors.Is(err, models.ErrSecurityKosong),
		errors.Is(err, models.ErrSecurityTanpaReinsurer),
		errors.Is(err, models.ErrSecurityMelampauiLebar):
		galat.Tulis(w, http.StatusUnprocessableEntity, pesanTCO(err))
	case errors.Is(err, services.ErrBusinessTidakAda):
		galat.Tulis(w, http.StatusNotFound, "business row not found in this contract combination")
	case errors.Is(err, services.ErrBusinessDobel):
		galat.Tulis(w, http.StatusConflict, pesanTCO(err))
	case errors.Is(err, models.ErrBusinessKodeKosong),
		errors.Is(err, models.ErrBusinessAktifTakSah),
		errors.Is(err, services.ErrBusinessDiLuarMaster):
		galat.Tulis(w, http.StatusUnprocessableEntity, pesanTCO(err))
	case errors.Is(err, services.ErrKlausulTidakAda):
		galat.Tulis(w, http.StatusNotFound, pesanTCO(err))
	case errors.Is(err, services.ErrKlausulDobel),
		errors.Is(err, services.ErrKlausulIndukBeranak):
		galat.Tulis(w, http.StatusConflict, pesanTCO(err))
	case errors.Is(err, services.ErrKlausulJenisBerubah):
		galat.Tulis(w, http.StatusBadRequest, pesanTCO(err))
	case errors.Is(err, models.ErrKlausulJenisTakDikenal),
		errors.Is(err, models.ErrKlausulDitahan),
		errors.Is(err, models.ErrKlausulMedanWajib),
		errors.Is(err, models.ErrMedanBukanMilikJenis),
		errors.Is(err, models.ErrTotalPctAnakMelebihi100),
		errors.Is(err, services.ErrJenisKlausulDiLuarMaster),
		errors.Is(err, services.ErrPilihanDiLuarMaster):
		galat.Tulis(w, http.StatusUnprocessableEntity, pesanTCO(err))
	case errors.Is(err, services.ErrLampiranTidakAda):
		galat.Tulis(w, http.StatusNotFound, "attachment not found in this treaty year")
	case errors.Is(err, services.ErrLampiranBelumTerkirim),
		errors.Is(err, services.ErrLampiranSudahTerkirim),
		errors.Is(err, services.ErrLampiranTanpaBerkas),
		errors.Is(err, services.ErrBerkasSumberLampiranHilang):
		// 409: keadaan DATA lampiran; pesannya menyebut lampiran mana dan
		// perbaikannya.
		galat.Tulis(w, http.StatusConflict, pesanTCO(err))
	case errors.Is(err, unggah.ErrBerkasTerlaluBesar):
		galat.Tulis(w, http.StatusRequestEntityTooLarge, pesanTCO(err))
	case errors.Is(err, unggah.ErrBerkasKosong):
		galat.Tulis(w, http.StatusBadRequest, services.PesanTanpaBerkasTCO)
	case errors.Is(err, services.ErrKategoriLampiranTidakDikenal):
		galat.Tulis(w, http.StatusUnprocessableEntity, pesanTCO(err))
	case errors.Is(err, services.ErrTahunTreatyDobel):
		// 409: keadaan DATA - kombinasi periode + grup sudah dipakai baris
		// lain - dan pesannya menyebut baris mana (AC 73).
		galat.Tulis(w, http.StatusConflict, pesanTCO(err))
	case errors.Is(err, models.ErrPeriodeTerbalik),
		errors.Is(err, models.ErrTahunTreatyGrupKosong),
		errors.Is(err, models.ErrTahunTreatyTahunKosong),
		errors.Is(err, models.ErrTahunTreatyBukanAngka):
		// 422: JSON-nya sah, isinya yang ditolak gerbang - pesan menyebut medannya.
		galat.Tulis(w, http.StatusUnprocessableEntity, pesanTCO(err))
	case errors.Is(err, galat.ErrPermintaanTidakSah):
		galat.Tulis(w, http.StatusBadRequest, pesanTCO(err))
	default:
		// ⛔ Sebab aslinya DICATAT: tanpa baris ini galat tak terduga hilang
		// di layar DAN di konsol backend (tco_galat500_test.go).
		log.Printf("treaty contract out: %v", err)
		galat.Tulis(w, http.StatusInternalServerError, "failed to process the treaty contract out request")
	}
	return true
}

// pesanTCO - kalimat galat untuk layar Treaty, BERBAHASA INGGRIS [keputusan
// work owner 30-09-2026]; sentinel `inti/` diganti `services.TeksInggrisTCO`.
func pesanTCO(err error) string {
	return services.TeksInggrisTCO(err.Error())
}

// akhirBawaanTahunTCO - End Date bawaan tahun treaty BARU (belum ber-ID):
// aturan yang SAMA dengan kontrak (`models.AkhirKontrakBawaanTCO`) [keputusan
// work owner 30-09-2026]. Murni hitungan - tanpa basis data.
func akhirBawaanTahunTCO(stub bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		akhir, err := services.AkhirTahunBawaan(inti.PelakuDari(r, stub), r.URL.Query().Get("mulai"))
		if jawabGalatTreatyContractOut(w, err) {
			return
		}
		galat.TulisJSON(w, struct {
			EndDate string `json:"endDate"`
		}{akhir})
	}
}
