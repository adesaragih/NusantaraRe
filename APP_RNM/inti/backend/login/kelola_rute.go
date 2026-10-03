package login

// Rute Kelola User - `/api/admin/*` (keputusan work owner 01-10-2026).
//
// ⛔ Setiap rute menuntut SESI LOGIN yang memegang menu Kelola User - juga
// saat AUTH_STUB=true: header stub dapat ditulis siapa saja, dan rute ini
// mengubah akun sungguhan. Tanpa sesi 401; sesi tanpa menu itu 403.
//
// ⛔ Password hanya di BADAN permintaan buat dan atur sandi - tidak pernah di
// jalur, log, atau jawaban.

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/inti/backend/menu"
)

// BatasBadanKelola - badan permintaan Kelola User: puluhan workbasket dan
// menu muat jauh di bawahnya.
const BatasBadanKelola = 16 << 10

// DenganKelola memasang layanan Kelola User. nil = rutenya menjawab 503.
func (r *Rute) DenganKelola(k *Kelola) *Rute {
	salin := *r
	salin.kelola = k
	return &salin
}

func (r *Rute) pasangKelola(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/admin/pengguna", r.daftarPengguna)
	mux.HandleFunc("POST /api/admin/pengguna", r.buatPengguna)
	mux.HandleFunc("GET /api/admin/pengguna/{id}", r.rinciPengguna)
	mux.HandleFunc("PUT /api/admin/pengguna/{id}", r.ubahPengguna)
	mux.HandleFunc("DELETE /api/admin/pengguna/{id}", r.hapusPengguna)
	mux.HandleFunc("POST /api/admin/pengguna/{id}/aktif", r.aktifPengguna)
	mux.HandleFunc("POST /api/admin/pengguna/{id}/buka-kunci", r.bukaKunciPengguna)
	mux.HandleFunc("POST /api/admin/pengguna/{id}/sandi", r.sandiPengguna)
	// Di luar `/pengguna/` supaya tidak menutupi akun bernama "pilihan".
	mux.HandleFunc("GET /api/admin/pilihan-pengguna", r.pilihanPengguna)
}

// admin menjawab pelaku rute Kelola User, atau menulis penolakannya.
func (r *Rute) admin(w http.ResponseWriter, req *http.Request) (string, bool) {
	if r.kelola == nil {
		galat.Tulis(w, http.StatusServiceUnavailable,
			"Kelola User belum dapat dipakai: ORACLE_DSN dan SESI_RAHASIA wajib terisi")
		return "", false
	}
	s, ok := sesiDari(req.Context())
	if !ok {
		galat.Tulis(w, http.StatusUnauthorized, "belum login atau sesi sudah berakhir")
		return "", false
	}
	if s.profil.WajibGantiSandi {
		galat.Tulis(w, http.StatusForbidden, "ganti sandi lebih dulu")
		return "", false
	}
	if !inti.PunyaMenu(s.profil.Menu, menu.KodeKelolaUser) {
		galat.Tulis(w, http.StatusForbidden, "akun ini tidak memegang menu Kelola User")
		return "", false
	}
	return s.profil.AkunID, true
}

// tulisGalatKelola menerjemahkan galat aturan ke kode HTTP. Galat lain ke log
// server - pesannya dapat memuat nilai kolom.
func tulisGalatKelola(w http.ResponseWriter, err error, apa string) {
	pesan := strings.TrimPrefix(err.Error(), "login: ")
	switch {
	case errors.Is(err, ErrAkunTidakAda):
		galat.Tulis(w, http.StatusNotFound, "akun tidak ada")
	case errors.Is(err, ErrAkunSudahAda):
		galat.Tulis(w, http.StatusConflict, "username sudah dipakai")
	case errors.Is(err, ErrDiriSendiri), errors.Is(err, ErrAdminTerakhir):
		galat.Tulis(w, http.StatusConflict, pesan)
	case errors.Is(err, ErrAkunTidakSah):
		galat.Tulis(w, http.StatusBadRequest,
			"username hanya huruf, angka, titik, garis bawah, @, atau tanda hubung (maks. 64); nama wajib (maks. 150)")
	case errors.Is(err, ErrSandiTerlaluPendek), errors.Is(err, ErrSandiTerlaluPanjang),
		errors.Is(err, ErrMenuTidakDikenal), errors.Is(err, ErrJenjangTidakCocok), errors.Is(err, ErrMasterTidakAda),
		errors.Is(err, ErrEmailTidakSah), errors.Is(err, ErrTeleponTidakSah), errors.Is(err, ErrNIKTidakSah),
		errors.Is(err, ErrJabatanTidakSah):
		galat.Tulis(w, http.StatusBadRequest, pesan)
	default:
		log.Printf("login: kelola user - %s: %v", apa, err)
		galat.Tulis(w, http.StatusInternalServerError, "Kelola User: "+apa+" gagal; rinciannya di log server")
	}
}

func (r *Rute) daftarPengguna(w http.ResponseWriter, req *http.Request) {
	if _, ok := r.admin(w, req); !ok {
		return
	}
	daftar, err := r.kelola.Daftar(req.Context())
	if err != nil {
		tulisGalatKelola(w, err, "membaca daftar akun")
		return
	}
	galat.TulisJSON(w, daftar)
}

func (r *Rute) rinciPengguna(w http.ResponseWriter, req *http.Request) {
	if _, ok := r.admin(w, req); !ok {
		return
	}
	rinci, err := r.kelola.Rinci(req.Context(), req.PathValue("id"))
	if err != nil {
		tulisGalatKelola(w, err, "membaca akun")
		return
	}
	galat.TulisJSON(w, rinci)
}

func (r *Rute) pilihanPengguna(w http.ResponseWriter, req *http.Request) {
	if _, ok := r.admin(w, req); !ok {
		return
	}
	p, err := r.kelola.Pilihan(req.Context())
	if err != nil {
		tulisGalatKelola(w, err, "membaca pilihan")
		return
	}
	galat.TulisJSON(w, p)
}

// isianPengguna - badan buat dan ubah. `akunId` dan `sandi` hanya dibaca saat buat.
type isianPengguna struct {
	AkunID     string   `json:"akunId"`
	Nama       string   `json:"nama"`
	Organisasi string   `json:"organisasi"`
	Divisi     string   `json:"divisi"`
	Unit       string   `json:"unit"`
	Workbasket []string `json:"workbasket"`
	Menu       []string `json:"menu"`
	Sandi      string   `json:"sandi"`
	// WajibGanti - centang "Change Password Next Login"; tidak dikirim = wajib.
	WajibGanti *bool `json:"wajibGanti"`
	// Kontak - `email`, `telepon`, `nik`, `jabatan` (Kelola User 03-10-2026); opsional.
	Kontak
}

// buatPengguna - sandi diketik admin; ⛔ tidak pernah ditulis ke log atau jawaban.
func (r *Rute) buatPengguna(w http.ResponseWriter, req *http.Request) {
	aktor, ok := r.admin(w, req)
	if !ok {
		return
	}
	var m isianPengguna
	if !bacaJSONBatas(w, req, &m, BatasBadanKelola) {
		return
	}
	a := AkunBaru{ID: m.AkunID, Nama: m.Nama, Organisasi: m.Organisasi, Divisi: m.Divisi, Unit: m.Unit,
		Workbasket: m.Workbasket, Menu: m.Menu, Kontak: m.Kontak}
	wajib := m.WajibGanti == nil || *m.WajibGanti
	if err := r.kelola.Buat(req.Context(), aktor, a, m.Sandi, wajib); err != nil {
		tulisGalatKelola(w, err, "membuat akun")
		return
	}
	id := strings.TrimSpace(m.AkunID)
	log.Printf("login: kelola user - %s membuat akun %s (wajib ganti password=%v)", aktor, id, wajib)
	rinci, err := r.kelola.Rinci(req.Context(), id)
	if err != nil {
		tulisGalatKelola(w, err, "membaca akun baru")
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(rinci)
}

func (r *Rute) ubahPengguna(w http.ResponseWriter, req *http.Request) {
	aktor, ok := r.admin(w, req)
	if !ok {
		return
	}
	var m isianPengguna
	if !bacaJSONBatas(w, req, &m, BatasBadanKelola) {
		return
	}
	id := req.PathValue("id")
	rinci, err := r.kelola.Ubah(req.Context(), aktor, id, IsianAkun{Nama: m.Nama, Organisasi: m.Organisasi,
		Divisi: m.Divisi, Unit: m.Unit, Workbasket: m.Workbasket, Menu: m.Menu, Kontak: m.Kontak})
	if err != nil {
		tulisGalatKelola(w, err, "mengubah akun")
		return
	}
	log.Printf("login: kelola user - %s mengubah akun %s", aktor, id)
	galat.TulisJSON(w, rinci)
}

func (r *Rute) aktifPengguna(w http.ResponseWriter, req *http.Request) {
	aktor, ok := r.admin(w, req)
	if !ok {
		return
	}
	var m struct {
		Aktif *bool `json:"aktif"`
	}
	if !bacaJSON(w, req, &m) {
		return
	}
	if m.Aktif == nil {
		galat.Tulis(w, http.StatusBadRequest, "badan wajib memuat aktif: true atau false")
		return
	}
	id := req.PathValue("id")
	rinci, err := r.kelola.SetelAktif(req.Context(), aktor, id, *m.Aktif)
	if err != nil {
		tulisGalatKelola(w, err, "mengubah status akun")
		return
	}
	log.Printf("login: kelola user - %s menyetel aktif=%v akun %s", aktor, *m.Aktif, id)
	galat.TulisJSON(w, rinci)
}

func (r *Rute) bukaKunciPengguna(w http.ResponseWriter, req *http.Request) {
	aktor, ok := r.admin(w, req)
	if !ok {
		return
	}
	id := req.PathValue("id")
	rinci, err := r.kelola.BukaKunci(req.Context(), aktor, id)
	if err != nil {
		tulisGalatKelola(w, err, "membuka kunci akun")
		return
	}
	log.Printf("login: kelola user - %s membuka kunci akun %s", aktor, id)
	galat.TulisJSON(w, rinci)
}

// sandiPengguna - tab Security: `sandi` terisi = password baru (sesi akun itu
// dicabut, kuncinya dibuka), kosong = hanya centang "Change Password Next
// Login". Password akun SENDIRI: cookie admin diterbitkan ulang, sesi ini
// tetap hidup.
func (r *Rute) sandiPengguna(w http.ResponseWriter, req *http.Request) {
	aktor, ok := r.admin(w, req)
	if !ok {
		return
	}
	var m struct {
		Sandi      string `json:"sandi"`
		WajibGanti *bool  `json:"wajibGanti"`
	}
	if !bacaJSON(w, req, &m) {
		return
	}
	if m.WajibGanti == nil {
		galat.Tulis(w, http.StatusBadRequest, "badan wajib memuat wajibGanti: true atau false")
		return
	}
	id := req.PathValue("id")
	rinci, err := r.kelola.AturSandi(req.Context(), aktor, id, m.Sandi, *m.WajibGanti)
	if err != nil {
		tulisGalatKelola(w, err, "mengatur password akun")
		return
	}
	if m.Sandi != "" {
		log.Printf("login: kelola user - %s MENGGANTI PASSWORD akun %s (wajib ganti=%v)", aktor, id, *m.WajibGanti)
	} else {
		log.Printf("login: kelola user - %s menyetel wajib ganti password=%v akun %s", aktor, *m.WajibGanti, id)
	}
	if m.Sandi != "" && id == aktor {
		s, _ := sesiDari(req.Context())
		tok, err := r.layanan.terbitkanUlang(req.Context(), s.token)
		if err != nil {
			tulisGalatKelola(w, err, "memperbarui sesi")
			return
		}
		r.pasangCookie(w, tok)
	}
	galat.TulisJSON(w, rinci)
}

// hapusPengguna - HAPUS PERMANEN (keputusan work owner): akun, workbasket,
// dan menunya. Layar meminta konfirmasi lebih dulu.
func (r *Rute) hapusPengguna(w http.ResponseWriter, req *http.Request) {
	aktor, ok := r.admin(w, req)
	if !ok {
		return
	}
	id := req.PathValue("id")
	if err := r.kelola.Hapus(req.Context(), aktor, id); err != nil {
		tulisGalatKelola(w, err, "menghapus akun")
		return
	}
	log.Printf("login: kelola user - %s MENGHAPUS PERMANEN akun %s", aktor, id)
	w.WriteHeader(http.StatusNoContent)
}
