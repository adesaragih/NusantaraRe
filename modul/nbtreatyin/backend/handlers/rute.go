// Package handlers adalah pintu HTTP modul NB Treaty In - seam 1 spec §6.2.
// Tugasnya hanya menerima permintaan, membaca pelaku, memanggil services, dan
// menerjemahkan galat ke kode HTTP. Nol aturan dagang, nol SQL (AC 60).
//
// Rute - prefix `/api/nb-treaty-in` (MODUL.md):
//
//	GET  /api/nb-treaty-in/kasus                      daftar portal (SFAPortal_OpportunitiesList; gerbang antrean, 403)
//	GET  /api/nb-treaty-in/kotak-masuk                cacah berkas yang menunggu akun per workbasket (Beranda)
//	GET  /api/nb-treaty-in/kotak-masuk/kasus          daftar berkas yang menunggu akun (?workbasket=; Beranda)
//	POST /api/nb-treaty-in/kasus                      Create (createWork)
//	GET  /api/nb-treaty-in/kasus/{id}                 buka assignment (pra-proses)
//	PUT  /api/nb-treaty-in/kasus/{id}                 Save layar admin
//	POST /api/nb-treaty-in/kasus/{id}/hitung          refresh berhitung (Count*, ...)
//	POST /api/nb-treaty-in/kasus/{id}/bisnis          grid popup BusinessAndSOBList (RD BrowseTreatyJoinEDM; tanpa simpan)
//	POST /api/nb-treaty-in/kasus/{id}/pilih-bisnis    Choose popup BusinessAndSOBList
//	POST /api/nb-treaty-in/kasus/{id}/pilih-sumber-bisnis  klik baris popup SOB (XOL Retro) - PostDT tanpa simpan (F4)
//	POST /api/nb-treaty-in/kasus/{id}/nomor-polis     GeneratePolicyNoTreaty_Act
//	POST /api/nb-treaty-in/kasus/{id}/kirim           finishAssignment
//	GET  /api/nb-treaty-in/kasus/{id}/riwayat         HISTORYAKSEPTASIPEGA
//	GET  /api/nb-treaty-in/sumber-bisnis              grid popup SOB (BrowseAgentHierarkiList_RD)
//	GET  /api/nb-treaty-in/acuan                      daftar pilihan layar
//	GET  /api/nb-treaty-in/hak                        hak layar portal akun ({"copyOld": bool})
//	GET  /api/nb-treaty-in/lama                       popup Copy Old - superadmin (WO 07-10-2026)
//	POST /api/nb-treaty-in/lama/salin                 Process Copy {"ids": [...]} - superadmin
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
	"nusantarare/modul/nbtreatyin/backend/models"
	"nusantarare/modul/nbtreatyin/backend/services"
)

// Prefix adalah awalan rute modul ini.
const Prefix = "/api/nb-treaty-in"

// KodeMenu - KODE menu modul ini (`M_LOGIN_GO_MENU.MENU_KODE`, sama dengan nama modul `backend.Nama`).
const KodeMenu = "nbtreatyin"

// batasBadan - badan permintaan terbesar (halaman kerja beserta daftarnya).
const batasBadan = 4 << 20

// DaftarkanRute memasang seluruh rute modul ke mux bersama.
func DaftarkanRute(mux *http.ServeMux, l *services.Layanan, stubPelaku bool) {
	h := &rute{l: l, stub: stubPelaku}
	mux.HandleFunc("GET "+Prefix+"/kasus", h.daftar)
	mux.HandleFunc("GET "+Prefix+"/kotak-masuk", h.kotakMasuk)
	mux.HandleFunc("GET "+Prefix+"/kotak-masuk/kasus", h.daftarMenunggu)
	mux.HandleFunc("POST "+Prefix+"/kasus", h.buat)
	mux.HandleFunc("GET "+Prefix+"/kasus/{id}", h.buka)
	mux.HandleFunc("PUT "+Prefix+"/kasus/{id}", h.simpan)
	mux.HandleFunc("POST "+Prefix+"/kasus/{id}/hitung", h.hitung)
	mux.HandleFunc("POST "+Prefix+"/kasus/{id}/bisnis", h.bisnis)
	mux.HandleFunc("POST "+Prefix+"/kasus/{id}/pilih-bisnis", h.pilihBisnis)
	mux.HandleFunc("POST "+Prefix+"/kasus/{id}/pilih-sumber-bisnis", h.pilihSumberBisnis)
	mux.HandleFunc("POST "+Prefix+"/kasus/{id}/nomor-polis", h.nomorPolis)
	mux.HandleFunc("POST "+Prefix+"/kasus/{id}/kirim", h.kirim)
	mux.HandleFunc("GET "+Prefix+"/kasus/{id}/riwayat", h.riwayat)
	mux.HandleFunc("GET "+Prefix+"/sumber-bisnis", h.sumberBisnis)
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

// superadmin - Copy Old (perintah work owner 07-10-2026, khusus superuser): pemegang menu Kelola User dengan menu NB
// Treaty In ber-hak PENUH - pola Copy Old Data Bordereaux (View only berlaku juga bagi superadmin, 05-10-2026).
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
	case galatData(err) != nil:
		// AC 37: kegagalan membaca data kontrak DITAMPILKAN kepada pengguna -
		// teks galat dasarnya saja. Rincian bungkusan (nama objek berskema,
		// kolom katalog, nomor baris dokumen) hanya ke log, sama dengan 500.
		log.Printf("nbtreatyin: %v", err)
		galat.Tulis(w, http.StatusUnprocessableEntity, galatData(err).Error())
	case errors.Is(err, services.ErrKasusTertutup),
		errors.Is(err, services.ErrTahapBerubah),
		errors.Is(err, services.ErrGenerasiTertutup),
		errors.Is(err, services.ErrNomorPolisSudahAda),
		errors.Is(err, services.ErrTindakanTakAdaDiPosisi):
		galat.Tulis(w, http.StatusConflict, err.Error())
	default:
		// Teks galat basis data (ORA-, nama skema) tidak dikirim ke klien.
		log.Printf("nbtreatyin: %v", err)
		galat.Tulis(w, http.StatusInternalServerError, "gagal memproses permintaan NB Treaty In")
	}
}

// galatDataDasar - galat data kontrak/master/nomor yang dijawab 422 (AC 37).
var galatDataDasar = []error{
	services.ErrDataKontrakTidakAda,
	services.ErrTipeNomorKosong,
	services.ErrOJKKosong,
	services.ErrMasterXOLTidakAda,
	services.ErrMasterXOLRusak,
}

// galatData - galat dasar `galatDataDasar` yang dibungkus `err`, atau nil.
func galatData(err error) error {
	for _, d := range galatDataDasar {
		if errors.Is(err, d) {
			return d
		}
	}
	return nil
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
	// IDDetail - ID baris view kontrak (pilih bisnis).
	IDDetail string `json:"idDetail"`
	// Saringan - saringan kolom popup Choose Business (nama kolom view -> teks).
	Saringan map[string]string `json:"saringan"`
}

func (h *rute) daftar(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	out, err := h.l.DaftarKasus(r.Context(), h.pelaku(r), models.SaringanKasus{
		Cari: strings.TrimSpace(q.Get("cari")), Posisi: strings.TrimSpace(q.Get("posisi")),
		Selesai: q.Get("status") == "selesai", // switch Proses / Resolved; bawaan Proses
	})
	if err != nil {
		tulisGalat(w, err)
		return
	}
	if out == nil {
		out = []models.RingkasanKasus{}
	}
	galat.TulisJSON(w, out)
}

// kotakMasuk - kotak masuk Beranda (keputusan work owner 06-10-2026).
func (h *rute) kotakMasuk(w http.ResponseWriter, r *http.Request) {
	out, err := h.l.KotakMasuk(r.Context(), h.pelaku(r))
	if err != nil {
		tulisGalat(w, err)
		return
	}
	galat.TulisJSON(w, out)
}

// daftarMenunggu - daftar berkas kotak masuk Beranda (keputusan work owner 06-10-2026).
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

func (h *rute) buat(w http.ResponseWriter, r *http.Request) {
	k, err := h.l.BuatKasus(r.Context(), h.pelaku(r))
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

func (h *rute) pilihBisnis(w http.ResponseWriter, r *http.Request) {
	var b badanHalaman
	if !bacaJSON(w, r, &b) {
		return
	}
	ly, err := h.l.PilihBisnis(r.Context(), h.pelaku(r), r.PathValue("id"), b.IDDetail, b.Halaman)
	if err != nil {
		tulisGalat(w, err)
		return
	}
	galat.TulisJSON(w, ly)
}

func (h *rute) nomorPolis(w http.ResponseWriter, r *http.Request) {
	var b badanHalaman
	if !bacaJSON(w, r, &b) {
		return
	}
	n, err := h.l.TerbitkanNomor(r.Context(), h.pelaku(r), r.PathValue("id"), b.Halaman)
	if err != nil {
		tulisGalat(w, err)
		return
	}
	galat.TulisJSON(w, n)
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

func (h *rute) riwayat(w http.ResponseWriter, r *http.Request) {
	out, err := h.l.RiwayatKasus(r.Context(), h.pelaku(r), r.PathValue("id"))
	if err != nil {
		tulisGalat(w, err)
		return
	}
	if out == nil {
		galat.TulisJSON(w, []struct{}{})
		return
	}
	galat.TulisJSON(w, out)
}

func (h *rute) bisnis(w http.ResponseWriter, r *http.Request) {
	var b badanHalaman
	if !bacaJSON(w, r, &b) {
		return
	}
	out, err := h.l.DaftarBisnis(r.Context(), h.pelaku(r), r.PathValue("id"), b.Halaman, b.Saringan)
	if err != nil {
		tulisGalat(w, err)
		return
	}
	if out == nil {
		out = []models.BarisKontrak{}
	}
	galat.TulisJSON(w, out)
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
