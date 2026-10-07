// Package handlers adalah pintu HTTP modul EDM Treaty In. Tugasnya hanya menerima permintaan, membaca pelaku,
// memanggil services, dan menerjemahkan galat ke kode HTTP. Nol aturan dagang, nol SQL. Pola: salinan
// `modul/nbtreatyin/backend/handlers/rute.go` (06-10-2026).
//
// Rute - prefix `/api/edm-treaty-in` (MODUL.md):
//
//	GET  /api/edm-treaty-in/kasus                     daftar portal (SFAPortal_Endorsement_Treaty, RD InboxEDM_RD2; ?cari=)
//	GET  /api/edm-treaty-in/kotak-masuk               cacah berkas yang menunggu akun per workbasket (Beranda)
//	GET  /api/edm-treaty-in/kotak-masuk/kasus         daftar berkas yang menunggu akun (?workbasket=; Beranda)
//	GET  /api/edm-treaty-in/periksa-polis             sel "No Polis Treaty" TreatyCreateEdm (?nopolis=; TrtEdmCheckPolicyError +
//	                                                  CheckNopolisAvailability)
//	POST /api/edm-treaty-in/kasus                     tombol Create TreatyCreateEdm (CreateEDMT)
//	GET  /api/edm-treaty-in/kasus/{id}                buka assignment (pra-proses)
//	PUT  /api/edm-treaty-in/kasus/{id}                Save layar admin
//	POST /api/edm-treaty-in/kasus/{id}/hitung         refresh berhitung (Count*, EDMTCalculateTreatyDifference, ...)
//	POST /api/edm-treaty-in/kasus/{id}/bisnis         grid popup BusinessAndSOBListEDM (tanpa simpan)
//	POST /api/edm-treaty-in/kasus/{id}/pilih-bisnis   Choose popup BusinessAndSOBListEDM (EDMChooseBusiness_Act)
//	POST /api/edm-treaty-in/kasus/{id}/kirim          finishAssignment
//	GET  /api/edm-treaty-in/acuan                     daftar pilihan layar
//	GET  /api/edm-treaty-in/hak                       hak layar portal akun ({"copyOld": bool})
//	GET  /api/edm-treaty-in/lama                      popup Copy Old - superadmin (WO 07-10-2026)
//	POST /api/edm-treaty-in/lama/salin                Process Copy {"ids": [...]} - superadmin
package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/inti/backend/menu"
	"nusantarare/modul/edmtreatyin/backend/models"
	"nusantarare/modul/edmtreatyin/backend/services"
)

// Prefix adalah awalan rute modul ini.
const Prefix = "/api/edm-treaty-in"

// KodeMenu - KODE menu modul ini (`M_LOGIN_GO_MENU.MENU_KODE`, sama dengan nama modul `backend.Nama`).
const KodeMenu = "edmtreatyin"

// batasBadan - badan permintaan terbesar (halaman kerja beserta daftarnya).
const batasBadan = 4 << 20

// DaftarkanRute memasang seluruh rute modul ke mux bersama.
func DaftarkanRute(mux *http.ServeMux, l *services.Layanan, stubPelaku bool) {
	h := &rute{l: l, stub: stubPelaku}
	mux.HandleFunc("GET "+Prefix+"/kasus", h.daftar)
	mux.HandleFunc("GET "+Prefix+"/kotak-masuk", h.kotakMasuk)
	mux.HandleFunc("GET "+Prefix+"/kotak-masuk/kasus", h.daftarMenunggu)
	mux.HandleFunc("GET "+Prefix+"/periksa-polis", h.periksaPolis)
	mux.HandleFunc("POST "+Prefix+"/kasus", h.buat)
	mux.HandleFunc("GET "+Prefix+"/kasus/{id}", h.buka)
	mux.HandleFunc("PUT "+Prefix+"/kasus/{id}", h.simpan)
	mux.HandleFunc("POST "+Prefix+"/kasus/{id}/hitung", h.hitung)
	mux.HandleFunc("POST "+Prefix+"/kasus/{id}/bisnis", h.bisnis)
	mux.HandleFunc("POST "+Prefix+"/kasus/{id}/pilih-bisnis", h.pilihBisnis)
	mux.HandleFunc("POST "+Prefix+"/kasus/{id}/kirim", h.kirim)
	mux.HandleFunc("GET "+Prefix+"/acuan", h.acuan)
	mux.HandleFunc("GET "+Prefix+"/hak", h.hak)
	mux.HandleFunc("GET "+Prefix+"/lama", h.daftarLama)
	mux.HandleFunc("POST "+Prefix+"/lama/salin", h.salinLama)
}

// Router menyusun mux tersendiri - untuk uji.
func Router(l *services.Layanan, stubPelaku bool) http.Handler {
	mux := http.NewServeMux()
	DaftarkanRute(mux, l, stubPelaku)
	return mux
}

type rute struct {
	l    *services.Layanan
	stub bool
}

func (h *rute) pelaku(r *http.Request) inti.Pelaku { return inti.PelakuDari(r, h.stub) }

// superadmin - Copy Old (perintah work owner 07-10-2026 "KHUSUS BUAT SUPERUSER"): pemegang menu Kelola User dengan menu
// EDM Treaty In ber-hak PENUH - pola Copy Old Data Bordereaux (View only berlaku juga bagi superadmin, 05-10-2026).
func superadmin(r *http.Request) bool {
	kode, _ := inti.AksesMenuDari(r.Context())
	return inti.PunyaMenu(kode, menu.KodeKelolaUser) && menu.BolehUbah(r.Context(), KodeMenu)
}

// tulisGalat menerjemahkan galat services ke kode HTTP.
func tulisGalat(w http.ResponseWriter, err error) {
	var v *services.GalatValidasi
	switch {
	case errors.As(err, &v):
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusUnprocessableEntity)
		_ = json.NewEncoder(w).Encode(struct {
			Galat string   `json:"galat"`
			Pesan []string `json:"pesan"`
		}{"validasi layar gagal: " + strings.Join(v.Pesan, "; "), v.Pesan})
	case errors.Is(err, services.ErrTanpaOracle):
		galat.Tulis(w, http.StatusServiceUnavailable, "database belum dikonfigurasi")
	case errors.Is(err, inti.ErrTanpaIdentitas):
		galat.Tulis(w, http.StatusUnauthorized, err.Error())
	case errors.Is(err, inti.ErrTanpaWewenang):
		galat.Tulis(w, http.StatusForbidden, err.Error())
	case errors.Is(err, services.ErrKasusTidakAda):
		galat.Tulis(w, http.StatusNotFound, err.Error())
	case errors.Is(err, services.ErrPermintaanTidakSah):
		galat.Tulis(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, services.ErrMasterXOLTidakAda), errors.Is(err, services.ErrMasterXOLRusak):
		// kegagalan membaca master kontrak ditampilkan kepada pengguna - teks galat dasarnya saja
		log.Printf("edmtreatyin: %v", err)
		if errors.Is(err, services.ErrMasterXOLTidakAda) {
			galat.Tulis(w, http.StatusUnprocessableEntity, services.ErrMasterXOLTidakAda.Error())
		} else {
			galat.Tulis(w, http.StatusUnprocessableEntity, services.ErrMasterXOLRusak.Error())
		}
	case errors.Is(err, services.ErrKasusTertutup),
		errors.Is(err, services.ErrTahapBerubah),
		errors.Is(err, services.ErrGenerasiTertutup),
		errors.Is(err, services.ErrGenerasiSudahDiendorse),
		errors.Is(err, services.ErrNomorPolisSudahAda),
		errors.Is(err, services.ErrTindakanTakAdaDiPosisi):
		galat.Tulis(w, http.StatusConflict, err.Error())
	default:
		// Teks galat basis data (ORA-, nama skema) tidak dikirim ke klien.
		log.Printf("edmtreatyin: %v", err)
		galat.Tulis(w, http.StatusInternalServerError, "gagal memproses permintaan EDM Treaty In")
	}
}

// bacaJSON membaca badan permintaan JSON (kosong = nilai nol).
func bacaJSON(w http.ResponseWriter, r *http.Request, tujuan any) bool {
	r.Body = http.MaxBytesReader(w, r.Body, batasBadan)
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(tujuan); err != nil && !errors.Is(err, io.EOF) {
		galat.Tulis(w, http.StatusBadRequest, "badan permintaan bukan JSON yang sah: "+err.Error())
		return false
	}
	return true
}

// badanHalaman - badan permintaan yang membawa halaman kerja.
type badanHalaman struct {
	Halaman *models.Halaman `json:"halaman"`
}

func (h *rute) daftar(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	// status=selesai - switch portal Resolved (aturan portal NB, keputusan WO 07-10-2026); bawaan In Progress
	out, err := h.l.DaftarKasus(r.Context(), h.pelaku(r), strings.TrimSpace(q.Get("cari")), q.Get("status") == "selesai")
	if err != nil {
		tulisGalat(w, err)
		return
	}
	if out == nil {
		out = []models.RingkasanKasus{}
	}
	galat.TulisJSON(w, out)
}

func (h *rute) kotakMasuk(w http.ResponseWriter, r *http.Request) {
	out, err := h.l.KotakMasuk(r.Context(), h.pelaku(r))
	if err != nil {
		tulisGalat(w, err)
		return
	}
	galat.TulisJSON(w, out)
}

func (h *rute) daftarMenunggu(w http.ResponseWriter, r *http.Request) {
	out, err := h.l.DaftarMenunggu(r.Context(), h.pelaku(r), strings.TrimSpace(r.URL.Query().Get("workbasket")))
	if err != nil {
		tulisGalat(w, err)
		return
	}
	if out == nil {
		out = []models.RingkasanKasus{}
	}
	galat.TulisJSON(w, out)
}

func (h *rute) periksaPolis(w http.ResponseWriter, r *http.Request) {
	out, err := h.l.PeriksaPolis(r.Context(), h.pelaku(r), r.URL.Query().Get("nopolis"))
	if err != nil {
		tulisGalat(w, err)
		return
	}
	galat.TulisJSON(w, out)
}

func (h *rute) buat(w http.ResponseWriter, r *http.Request) {
	var b services.PermintaanBuat
	if !bacaJSON(w, r, &b) {
		return
	}
	k, err := h.l.BuatKasus(r.Context(), h.pelaku(r), b)
	if err != nil {
		tulisGalat(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(k)
}

func (h *rute) buka(w http.ResponseWriter, r *http.Request) {
	ly, err := h.l.BukaKasus(r.Context(), h.pelaku(r), r.PathValue("id"))
	if err != nil {
		tulisGalat(w, err)
		return
	}
	galat.TulisJSON(w, ly)
}

func (h *rute) simpan(w http.ResponseWriter, r *http.Request) {
	var b badanHalaman
	if !bacaJSON(w, r, &b) {
		return
	}
	ly, err := h.l.SimpanDraf(r.Context(), h.pelaku(r), r.PathValue("id"), b.Halaman)
	if err != nil {
		tulisGalat(w, err)
		return
	}
	galat.TulisJSON(w, ly)
}

func (h *rute) hitung(w http.ResponseWriter, r *http.Request) {
	var b services.PermintaanHitung
	if !bacaJSON(w, r, &b) {
		return
	}
	ly, err := h.l.Hitung(r.Context(), h.pelaku(r), r.PathValue("id"), b)
	if err != nil {
		tulisGalat(w, err)
		return
	}
	galat.TulisJSON(w, ly)
}

func (h *rute) bisnis(w http.ResponseWriter, r *http.Request) {
	var b badanHalaman
	if !bacaJSON(w, r, &b) {
		return
	}
	out, err := h.l.DaftarBisnis(r.Context(), h.pelaku(r), r.PathValue("id"), b.Halaman)
	if err != nil {
		tulisGalat(w, err)
		return
	}
	if out == nil {
		out = []models.Baris{}
	}
	galat.TulisJSON(w, out)
}

func (h *rute) pilihBisnis(w http.ResponseWriter, r *http.Request) {
	var b badanHalaman
	if !bacaJSON(w, r, &b) {
		return
	}
	ly, err := h.l.PilihBisnis(r.Context(), h.pelaku(r), r.PathValue("id"), b.Halaman)
	if err != nil {
		tulisGalat(w, err)
		return
	}
	galat.TulisJSON(w, ly)
}

func (h *rute) kirim(w http.ResponseWriter, r *http.Request) {
	var b badanHalaman
	if !bacaJSON(w, r, &b) {
		return
	}
	k, err := h.l.Kirim(r.Context(), h.pelaku(r), r.PathValue("id"), b.Halaman)
	if err != nil {
		tulisGalat(w, err)
		return
	}
	galat.TulisJSON(w, k)
}

func (h *rute) acuan(w http.ResponseWriter, r *http.Request) {
	a, err := h.l.DaftarAcuan(r.Context(), h.pelaku(r))
	if err != nil {
		tulisGalat(w, err)
		return
	}
	galat.TulisJSON(w, a)
}

func (h *rute) hak(w http.ResponseWriter, r *http.Request) {
	galat.TulisJSON(w, h.l.HakPortal(h.pelaku(r), superadmin(r)))
}

func (h *rute) daftarLama(w http.ResponseWriter, r *http.Request) {
	out, err := h.l.DaftarDokumenLama(r.Context(), h.pelaku(r), superadmin(r))
	if err != nil {
		tulisGalat(w, err)
		return
	}
	galat.TulisJSON(w, out)
}

func (h *rute) salinLama(w http.ResponseWriter, r *http.Request) {
	var b struct {
		IDs []string `json:"ids"`
	}
	if !bacaJSON(w, r, &b) {
		return
	}
	j, err := h.l.SalinDokumenLama(r.Context(), h.pelaku(r), superadmin(r), b.IDs)
	if err != nil {
		tulisGalat(w, err)
		return
	}
	galat.TulisJSON(w, j)
}
