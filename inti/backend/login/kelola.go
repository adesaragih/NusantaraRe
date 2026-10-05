package login

// Kelola User - daftar, buat, ubah, aktif/nonaktif, buka kunci, atur
// password, dan hapus permanen akun M_LOGIN_GO beserta workbasket dan
// menunya (keputusan work owner 01-10-2026).
//
// Password DIKETIK admin - saat membuat akun dan di tab Security (permintaan
// work owner 01-10-2026: "tambahkan tab security untuk mengelola password,
// jadi bisa ganti password sendiri dan password akun lain", menggantikan
// perintah sebelumnya "reset sandi jangan dulu"). Centang "Change Password
// Next Login" mengisi MUST_CHANGE_PASSWORD: dicentang = wajib ganti saat login
// berikutnya, tidak dicentang = tidak perlu.
//
// Penjaga (keputusan work owner):
//   - admin tidak dapat mencabut Kelola User dari dirinya, menonaktifkan, atau
//     menghapus dirinya sendiri (`ErrDiriSendiri`, diperiksa di sini)
//   - harus tersisa minimal satu akun AKTIF yang memegang Kelola User
//     (`ErrAdminTerakhir`, ditegakkan gudang DI DALAM transaksi perubahannya:
//     dua admin yang saling mencabut tidak dapat sama-sama lolos)

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/menu"
)

var (
	// ErrDiriSendiri - admin menonaktifkan, menghapus, atau mencabut Kelola
	// User dari dirinya sendiri.
	ErrDiriSendiri = errors.New("login: admin tidak dapat menonaktifkan, menghapus, atau mencabut Kelola User dari dirinya sendiri")
	// ErrAdminTerakhir - perubahan menyisakan nol akun aktif ber-Kelola User.
	ErrAdminTerakhir = errors.New("login: harus tersisa minimal satu akun aktif yang memegang Kelola User")
	// ErrMenuTidakDikenal - KODE menu bukan baris modul tabel menu maupun menu aplikasi.
	ErrMenuTidakDikenal = errors.New("login: menu tidak dikenal")
	// ErrHakLihatTidakSah - hak LIHAT untuk menu yang tidak dipegang akun, atau milik modul yang belum mendaftar
	// akses Lihat (`inti.Pendaftaran.HakLihat`).
	ErrHakLihatTidakSah = errors.New("login: akses View only hanya untuk menu yang dipilih dan mendukungnya")
)

// RingkasAkun adalah satu baris daftar Kelola User.
//
// ⛔ Nol hash, nol versi sesi.
type RingkasAkun struct {
	AkunID string `json:"akunId"`
	// IDKontak - `CONTACT_ID` `CON-n` (migrasi 905): diberi saat akun dibuat, tidak pernah berubah.
	IDKontak        string `json:"contactId"`
	Nama            string `json:"nama"`
	Organisasi      string `json:"organisasi"`
	Divisi          string `json:"divisi"`
	Unit            string `json:"unit"`
	Aktif           bool   `json:"aktif"`
	Terkunci        bool   `json:"terkunci"`
	WajibGantiSandi bool   `json:"wajibGantiSandi"`
	// LoginTerakhir - `LAST_LOGIN`, `YYYY-MM-DD HH:MI` jam Oracle; kosong = belum pernah.
	LoginTerakhir string `json:"loginTerakhir"`
	// Kontak - `email`, `telepon`, `nik`, `jabatan` (migrasi 904, Kelola User 03-10-2026).
	Kontak
}

// RinciAkun adalah satu akun beserta workbasket dan menunya - isi form ubah.
type RinciAkun struct {
	RingkasAkun
	Workbasket []string `json:"workbasket"`
	Menu       []string `json:"menu"`
	// MenuLihat - bagian dari Menu yang ber-hak LIHAT (migrasi 914); selalu terisi.
	MenuLihat []string `json:"menuLihat"`
}

// IsianAkun adalah isian ubah akun. `LOGIN_ID` tidak dapat diubah.
type IsianAkun struct {
	Nama                     string
	Organisasi, Divisi, Unit string
	Workbasket               []string
	Menu                     []string
	// MenuLihat - bagian dari Menu yang ber-hak LIHAT (migrasi 914).
	MenuLihat []string
	Kontak
}

// OpsiMaster adalah satu pilihan dropdown master. `Induk` = CODE organisasi
// pemilik divisi, atau CODE divisi pemilik unit.
type OpsiMaster struct {
	Kode  string `json:"kode"`
	Nama  string `json:"nama"`
	Induk string `json:"induk,omitempty"`
}

// OpsiMenu adalah satu kotak centang menu.
type OpsiMenu struct {
	Kode      string `json:"kode"`
	Label     string `json:"label"`
	Golongan  string `json:"golongan"`
	Dimigrasi bool   `json:"dimigrasi"`
	// BisaLihat - menu modul yang mendaftar akses Lihat (`inti.Pendaftaran.HakLihat`): Kelola User menawarkan
	// pilihan View only / Full di sampingnya.
	BisaLihat bool `json:"bisaLihat"`
}

// PilihanKelola adalah isi pilihan form Kelola User. Master hanya yang aktif.
type PilihanKelola struct {
	Organisasi []OpsiMaster `json:"organisasi"`
	Divisi     []OpsiMaster `json:"divisi"`
	Unit       []OpsiMaster `json:"unit"`
	Workbasket []OpsiMaster `json:"workbasket"`
	Menu       []OpsiMenu   `json:"menu"`
}

// GudangKelola menambah Gudang dengan perubahan milik admin.
//
// ⛔ UbahAkun, SetelAktif, dan HapusAkun menegakkan `ErrAdminTerakhir` di
// DALAM transaksinya - sesudah perubahan, sebelum COMMIT.
type GudangKelola interface {
	Gudang
	DaftarAkun(ctx context.Context) ([]RingkasAkun, error)
	RingkasAkun(ctx context.Context, id string) (RingkasAkun, error)
	// WorkbasketSemua - SELURUH WORKBASKET_ID akun itu, juga yang nonaktif di master.
	WorkbasketSemua(ctx context.Context, id string) ([]string, error)
	// UbahAkun menulis profil dan MENGGANTI seluruh workbasket dan menunya.
	UbahAkun(ctx context.Context, id string, a IsianAkun) error
	// SetelAktif mengubah IS_ACTIVE dan menaikkan SESSION_VERSION - cookie
	// lama tidak hidup lagi walau akunnya diaktifkan kembali.
	SetelAktif(ctx context.Context, id string, aktif bool) error
	// BukaKunci menolkan FAILED_COUNT dan mencabut LOCKED_UNTIL.
	BukaKunci(ctx context.Context, id string) error
	// AturSandi menulis MUST_CHANGE_PASSWORD; `hash` terisi = juga hash baru,
	// SESSION_VERSION naik (seluruh sesinya dicabut), dan kuncinya dibuka.
	AturSandi(ctx context.Context, id, hash string, wajibGanti bool) error
	// HapusAkun membuang menu, workbasket, lalu akunnya - satu transaksi.
	HapusAkun(ctx context.Context, id string) error
	// Master - organisasi, divisi, unit, dan workbasket AKTIF (tanpa Menu).
	Master(ctx context.Context) (PilihanKelola, error)
}

// Kelola menjalankan aturan Kelola User.
type Kelola struct {
	layanan *Layanan
	gudang  GudangKelola
	menu    menu.PembacaMenu
	// bisaLihat - KODE menu modul yang mendaftar akses Lihat (`DenganHakLihat`).
	bisaLihat map[string]bool
}

// NewKelola menyusunnya. `m` membaca baris modul tabel menu - sumber KODE
// menu yang sah, bersama `menu.MenuAplikasi`.
func NewKelola(g GudangKelola, m menu.PembacaMenu) *Kelola {
	return &Kelola{layanan: NewLayanan(g, nil), gudang: g, menu: m}
}

// DenganHakLihat menyebut KODE menu modul yang mendukung akses LIHAT (`daftar.HakLihat`, keputusan work owner
// 04-10-2026). Menu lain selalu penuh.
func (k *Kelola) DenganHakLihat(kode []string) *Kelola {
	k.bisaLihat = map[string]bool{}
	for _, x := range kode {
		k.bisaLihat[x] = true
	}
	return k
}

// periksaHakLihat - setiap menu LIHAT harus ikut dipilih dan milik modul yang mendaftar.
func (k *Kelola) periksaHakLihat(menuDipilih, lihat []string) error {
	for _, x := range lihat {
		if !inti.PunyaMenu(menuDipilih, x) || !k.bisaLihat[x] {
			return fmt.Errorf("%w: %s", ErrHakLihatTidakSah, x)
		}
	}
	return nil
}

// bersih membuang isian kosong dan ganda, lalu mengurutkan.
func bersih(daftar []string) []string {
	ada := map[string]bool{}
	out := []string{}
	for _, v := range daftar {
		if v = strings.TrimSpace(v); v != "" && !ada[v] {
			ada[v] = true
			out = append(out, v)
		}
	}
	sort.Strings(out)
	return out
}

// selisih - isi `a` yang tidak ada di `b`.
func selisih(a, b []string) []string {
	ada := map[string]bool{}
	for _, x := range b {
		ada[x] = true
	}
	var out []string
	for _, x := range a {
		if !ada[x] {
			out = append(out, x)
		}
	}
	return out
}

// barisMenu - baris modul tabel menu (`KODE = MODUL`).
func (k *Kelola) barisMenu(ctx context.Context) ([]menu.Baris, error) {
	if k.menu == nil {
		return nil, errors.New("login: pembaca menu tidak tersedia")
	}
	baris, err := k.menu.Baca(ctx)
	if err != nil {
		return nil, fmt.Errorf("login: membaca menu: %w", err)
	}
	var out []menu.Baris
	for _, b := range baris {
		if b.Kode == b.Modul {
			out = append(out, b)
		}
	}
	return out, nil
}

// PeriksaMenu menolak KODE menu yang bukan baris modul tabel menu maupun
// menu aplikasi (`ErrMenuTidakDikenal`) - juga dipakai `-buat-pengguna -menu`.
func (k *Kelola) PeriksaMenu(ctx context.Context, kode []string) error {
	baris, err := k.barisMenu(ctx)
	if err != nil {
		return err
	}
	sah := map[string]bool{}
	for _, b := range baris {
		sah[b.Kode] = true
	}
	for _, a := range menu.MenuAplikasi {
		sah[a.Kode] = true
	}
	for _, x := range kode {
		if !sah[x] {
			return fmt.Errorf("%w: %s", ErrMenuTidakDikenal, x)
		}
	}
	return nil
}

// Daftar membaca seluruh akun, urut LOGIN_ID.
func (k *Kelola) Daftar(ctx context.Context) ([]RingkasAkun, error) {
	return k.gudang.DaftarAkun(ctx)
}

// Rinci membaca satu akun beserta workbasket aktif dan menunya.
func (k *Kelola) Rinci(ctx context.Context, id string) (RinciAkun, error) {
	r, err := k.gudang.RingkasAkun(ctx, id)
	if err != nil {
		return RinciAkun{}, err
	}
	wb, err := k.gudang.Workbasket(ctx, id)
	if err != nil {
		return RinciAkun{}, err
	}
	mn, err := k.gudang.Menu(ctx, id)
	if err != nil {
		return RinciAkun{}, err
	}
	ml, err := k.gudang.MenuLihat(ctx, id)
	if err != nil {
		return RinciAkun{}, err
	}
	mn = bersih(mn)
	return RinciAkun{RingkasAkun: r, Workbasket: bersih(wb), Menu: mn, MenuLihat: irisan(bersih(ml), mn)}, nil
}

// Pilihan menyusun isi dropdown dan kotak centang form.
func (k *Kelola) Pilihan(ctx context.Context) (PilihanKelola, error) {
	p, err := k.gudang.Master(ctx)
	if err != nil {
		return PilihanKelola{}, err
	}
	baris, err := k.barisMenu(ctx)
	if err != nil {
		return PilihanKelola{}, err
	}
	urutGolongan := map[string]int{}
	for i, g := range menu.Golongan {
		urutGolongan[g] = i
	}
	sort.SliceStable(baris, func(i, j int) bool {
		a, b := baris[i], baris[j]
		if urutGolongan[a.Golongan] != urutGolongan[b.Golongan] {
			return urutGolongan[a.Golongan] < urutGolongan[b.Golongan]
		}
		if a.Urutan != b.Urutan {
			return a.Urutan < b.Urutan
		}
		return a.ID < b.ID
	})
	p.Menu = []OpsiMenu{}
	for _, b := range baris {
		p.Menu = append(p.Menu, OpsiMenu{Kode: b.Kode, Label: b.Label, Golongan: b.Golongan, Dimigrasi: b.Dimigrasi,
			BisaLihat: k.bisaLihat[b.Kode]})
	}
	for _, a := range menu.MenuAplikasi {
		p.Menu = append(p.Menu, OpsiMenu{Kode: a.Kode, Label: a.Label, Golongan: menu.GolonganAdmin, Dimigrasi: true})
	}
	return p, nil
}

// Buat membuat akun dengan sandi yang DIKETIK admin - aturan sandi berlaku.
// `wajibGanti` = centang "Change Password Next Login".
func (k *Kelola) Buat(ctx context.Context, aktor string, a AkunBaru, sandi string, wajibGanti bool) error {
	if err := PeriksaSandiBaru(sandi); err != nil {
		return err
	}
	a.Workbasket, a.Menu, a.MenuLihat = bersih(a.Workbasket), bersih(a.Menu), bersih(a.MenuLihat)
	a.Organisasi, a.Divisi, a.Unit = strings.TrimSpace(a.Organisasi), strings.TrimSpace(a.Divisi), strings.TrimSpace(a.Unit)
	a.Kontak = RapikanKontak(a.Kontak)
	if err := PeriksaKontak(a.Kontak); err != nil {
		return err
	}
	if err := k.PeriksaMenu(ctx, a.Menu); err != nil {
		return err
	}
	if err := k.periksaHakLihat(a.Menu, a.MenuLihat); err != nil {
		return err
	}
	return k.layanan.buat(ctx, a, sandi, wajibGanti)
}

// AturSandi - tab Security. `sandi` terisi = password baru (aturan sandi
// berlaku; seluruh sesi akun itu dicabut dan kuncinya dibuka); kosong = hanya
// centang "Change Password Next Login" yang berubah. Password akun SENDIRI
// boleh: rute Kelola User menerbitkan ulang cookie admin itu.
func (k *Kelola) AturSandi(ctx context.Context, aktor, id, sandi string, wajibGanti bool) (RinciAkun, error) {
	hash := ""
	if sandi != "" {
		if err := PeriksaSandiBaru(sandi); err != nil {
			return RinciAkun{}, err
		}
		h, err := HashSandi(sandi)
		if err != nil {
			return RinciAkun{}, err
		}
		hash = h
	}
	if _, err := k.gudang.AmbilAkun(ctx, id); err != nil {
		return RinciAkun{}, err
	}
	if err := k.gudang.AturSandi(ctx, id, hash, wajibGanti); err != nil {
		return RinciAkun{}, err
	}
	return k.Rinci(ctx, id)
}

// Ubah menulis profil dan mengganti seluruh workbasket dan menu akun itu.
//
// ⛔ Yang TIDAK terlihat di form tidak dibuang dan tidak menghalangi
// (temuan /code-review): workbasket yang kini nonaktif di master tidak
// ditampilkan form, jadi dipertahankan apa adanya; dan jenjang organisasi
// yang TIDAK diubah tidak diperiksa ulang - master yang belakangan
// dinonaktifkan tidak boleh mengunci akun dari perubahan lain. Yang BARU
// diisi (jenjang yang berubah, workbasket yang ditambah) tetap harus ada
// dan aktif.
func (k *Kelola) Ubah(ctx context.Context, aktor, id string, isi IsianAkun) (RinciAkun, error) {
	lama, err := k.gudang.AmbilAkun(ctx, id)
	if err != nil {
		return RinciAkun{}, err
	}
	isi.Nama = strings.TrimSpace(isi.Nama)
	if isi.Nama == "" || len(isi.Nama) > 150 {
		return RinciAkun{}, ErrAkunTidakSah
	}
	isi.Organisasi, isi.Divisi, isi.Unit = strings.TrimSpace(isi.Organisasi), strings.TrimSpace(isi.Divisi), strings.TrimSpace(isi.Unit)
	isi.Workbasket, isi.Menu, isi.MenuLihat = bersih(isi.Workbasket), bersih(isi.Menu), bersih(isi.MenuLihat)
	isi.Kontak = RapikanKontak(isi.Kontak)
	if err := PeriksaKontak(isi.Kontak); err != nil {
		return RinciAkun{}, err
	}
	if err := k.layanan.periksaEmailBebas(ctx, id, isi.Email); err != nil {
		return RinciAkun{}, err
	}
	if id == aktor && !inti.PunyaMenu(isi.Menu, menu.KodeKelolaUser) {
		return RinciAkun{}, ErrDiriSendiri
	}
	wbSemua, err := k.gudang.WorkbasketSemua(ctx, id)
	if err != nil {
		return RinciAkun{}, err
	}
	wbAktif, err := k.gudang.Workbasket(ctx, id)
	if err != nil {
		return RinciAkun{}, err
	}
	periksa := AkunBaru{Workbasket: selisih(isi.Workbasket, wbSemua)}
	if isi.Organisasi != lama.Organisasi || isi.Divisi != lama.Divisi || isi.Unit != lama.Unit {
		periksa.Organisasi, periksa.Divisi, periksa.Unit = isi.Organisasi, isi.Divisi, isi.Unit
	}
	if err := k.layanan.periksaJenjang(ctx, periksa); err != nil {
		return RinciAkun{}, err
	}
	isi.Workbasket = bersih(append(isi.Workbasket, selisih(wbSemua, wbAktif)...))
	if err := k.PeriksaMenu(ctx, isi.Menu); err != nil {
		return RinciAkun{}, err
	}
	if err := k.periksaHakLihat(isi.Menu, isi.MenuLihat); err != nil {
		return RinciAkun{}, err
	}
	if err := k.gudang.UbahAkun(ctx, id, isi); err != nil {
		return RinciAkun{}, err
	}
	return k.Rinci(ctx, id)
}

// SetelAktif mengaktifkan atau menonaktifkan akun - nonaktif mencabut seluruh
// sesinya pada permintaan berikutnya.
func (k *Kelola) SetelAktif(ctx context.Context, aktor, id string, aktif bool) (RinciAkun, error) {
	if id == aktor && !aktif {
		return RinciAkun{}, ErrDiriSendiri
	}
	if _, err := k.gudang.AmbilAkun(ctx, id); err != nil {
		return RinciAkun{}, err
	}
	if err := k.gudang.SetelAktif(ctx, id, aktif); err != nil {
		return RinciAkun{}, err
	}
	return k.Rinci(ctx, id)
}

// BukaKunci mencabut kunci sesudah sandi salah beruntun.
func (k *Kelola) BukaKunci(ctx context.Context, aktor, id string) (RinciAkun, error) {
	if _, err := k.gudang.AmbilAkun(ctx, id); err != nil {
		return RinciAkun{}, err
	}
	if err := k.gudang.BukaKunci(ctx, id); err != nil {
		return RinciAkun{}, err
	}
	return k.Rinci(ctx, id)
}

// Hapus membuang akun beserta workbasket dan menunya, PERMANEN.
func (k *Kelola) Hapus(ctx context.Context, aktor, id string) error {
	if id == aktor {
		return ErrDiriSendiri
	}
	if _, err := k.gudang.AmbilAkun(ctx, id); err != nil {
		return err
	}
	return k.gudang.HapusAkun(ctx, id)
}
