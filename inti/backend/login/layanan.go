package login

// Aturan login - masuk, sesi, keluar, ganti sandi, buat pengguna.

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	// ErrKredensial - SATU jawaban untuk akun tidak ada, nonaktif, dan sandi
	// salah: jawaban tidak boleh membocorkan akun mana yang ada.
	ErrKredensial = errors.New("login: akun atau sandi salah")
	// ErrTerkunci - akun terkunci sementara sesudah BatasGagal sandi salah.
	ErrTerkunci = errors.New("login: akun terkunci sementara")
	// ErrTanpaRahasia - SESI_RAHASIA kosong; sesi tidak dapat diterbitkan.
	ErrTanpaRahasia = errors.New("login: SESI_RAHASIA belum disetel")
	// ErrSandiSama - sandi baru sama dengan sandi lama.
	ErrSandiSama = errors.New("login: sandi baru sama dengan sandi lama")
	// ErrAkunTidakSah - bentuk LOGIN_ID atau NAME tidak dapat dipakai.
	ErrAkunTidakSah = errors.New("login: akun atau nama tidak sah")
	// ErrAkunSudahAda - LOGIN_ID sudah dipakai.
	ErrAkunSudahAda = errors.New("login: akun sudah ada")
	// ErrJenjangTidakCocok - unit bukan milik divisinya, divisi bukan milik
	// organisasinya, atau tingkat bawah diisi tanpa tingkat atasnya.
	ErrJenjangTidakCocok = errors.New("login: jenjang organisasi tidak cocok")
)

// Layanan menjalankan aturan login di atas sebuah Gudang.
type Layanan struct {
	gudang   Gudang
	rahasia  []byte
	sekarang func() time.Time
}

// NewLayanan menyusunnya. `rahasia` kosong = sesi tidak dapat diterbitkan
// (cukup untuk `-buat-pengguna`).
func NewLayanan(g Gudang, rahasia []byte) *Layanan {
	return &Layanan{gudang: g, rahasia: rahasia, sekarang: time.Now}
}

// DenganJam mengganti jam - untuk uji.
func (l *Layanan) DenganJam(f func() time.Time) *Layanan {
	salin := *l
	salin.sekarang = f
	return &salin
}

func (l *Layanan) profil(ctx context.Context, a Akun) (Profil, error) {
	peran, err := l.gudang.Workbasket(ctx, a.ID)
	if err != nil {
		return Profil{}, err
	}
	if peran == nil {
		peran = []string{}
	}
	menu, err := l.gudang.Menu(ctx, a.ID)
	if err != nil {
		return Profil{}, err
	}
	if menu == nil {
		menu = []string{}
	}
	lihat, err := l.gudang.MenuLihat(ctx, a.ID)
	if err != nil {
		return Profil{}, err
	}
	return Profil{AkunID: a.ID, Nama: a.Nama, Peran: peran, Organisasi: a.Organisasi,
		Divisi: a.Divisi, Unit: a.Unit, WajibGantiSandi: a.WajibGantiSandi, Menu: menu, MenuLihat: irisan(lihat, menu)}, nil
}

// Masuk memeriksa akun dan sandi lalu menerbitkan token sesi.
func (l *Layanan) Masuk(ctx context.Context, akun, sandi string) (Profil, Token, error) {
	if len(l.rahasia) == 0 {
		return Profil{}, Token{}, ErrTanpaRahasia
	}
	akun = strings.TrimSpace(akun)
	if !AkunSah(akun) {
		CocokSandi(hashTiruan(), sandi)
		return Profil{}, Token{}, ErrKredensial
	}
	a, err := l.gudang.AmbilAkun(ctx, akun)
	if errors.Is(err, ErrAkunTidakAda) || (err == nil && !a.Aktif) {
		CocokSandi(hashTiruan(), sandi) // waktu jawab sama dengan akun yang ada
		return Profil{}, Token{}, ErrKredensial
	}
	if err != nil {
		return Profil{}, Token{}, err
	}
	// ⛔ Terkunci diperiksa SEBELUM sandi: sandi benar pun ditolak, dan
	// percobaan selama terkunci tidak memperpanjang kuncinya.
	if a.Terkunci {
		return Profil{}, Token{}, ErrTerkunci
	}
	if !CocokSandi(a.HashSandi, sandi) {
		if err := l.gudang.CatatGagal(ctx, a.ID); err != nil {
			return Profil{}, Token{}, err
		}
		return Profil{}, Token{}, ErrKredensial
	}
	if err := l.gudang.CatatBerhasil(ctx, a.ID); err != nil {
		return Profil{}, Token{}, err
	}
	p, err := l.profil(ctx, a)
	if err != nil {
		return Profil{}, Token{}, err
	}
	return p, TokenBaru(a.ID, a.VersiSesi, l.sekarang()), nil
}

// Sesi membaca cookie, mencocokkannya dengan M_LOGIN_GO, dan
// memperpanjangnya. ErrSesiTidakSah bila rusak, kedaluwarsa, dicabut, atau
// akunnya nonaktif.
func (l *Layanan) Sesi(ctx context.Context, nilai string) (Profil, Token, error) {
	tok, err := BacaToken(l.rahasia, nilai, l.sekarang())
	if err != nil {
		return Profil{}, Token{}, err
	}
	a, err := l.gudang.AmbilAkun(ctx, tok.Akun)
	if errors.Is(err, ErrAkunTidakAda) {
		return Profil{}, Token{}, ErrSesiTidakSah
	}
	if err != nil {
		return Profil{}, Token{}, err
	}
	if !a.Aktif || a.VersiSesi != tok.Versi {
		return Profil{}, Token{}, ErrSesiTidakSah
	}
	p, err := l.profil(ctx, a)
	if err != nil {
		return Profil{}, Token{}, err
	}
	return p, tok.Perpanjang(l.sekarang()), nil
}

// terbitkanUlang memperbarui token sesi YANG SEDANG DIPAKAI ke versi sesi
// akun kini - sesudah admin mengganti password-nya sendiri di Kelola User
// (versinya naik), supaya sesi ini tidak ikut tercabut. Hanya dipanggil rute
// yang sesinya sudah diperiksa di permintaan yang sama.
func (l *Layanan) terbitkanUlang(ctx context.Context, tok Token) (Token, error) {
	a, err := l.gudang.AmbilAkun(ctx, tok.Akun)
	if err != nil {
		return Token{}, err
	}
	if !a.Aktif {
		return Token{}, ErrSesiTidakSah
	}
	tok.Versi = a.VersiSesi
	return tok.Perpanjang(l.sekarang()), nil
}

// Keluar mencabut SELURUH sesi akun itu.
func (l *Layanan) Keluar(ctx context.Context, akun string) error {
	return l.gudang.NaikkanVersi(ctx, akun)
}

// GantiSandi menuntut sandi lama, menegakkan aturan sandi baru, dan
// menerbitkan token berversi baru - sesi lain akun itu tercabut, sesi ini lanjut.
func (l *Layanan) GantiSandi(ctx context.Context, tok Token, lama, baru string) (Profil, Token, error) {
	a, err := l.gudang.AmbilAkun(ctx, tok.Akun)
	if err != nil {
		return Profil{}, Token{}, err
	}
	if !CocokSandi(a.HashSandi, lama) {
		return Profil{}, Token{}, ErrKredensial
	}
	if err := PeriksaSandiBaru(baru); err != nil {
		return Profil{}, Token{}, err
	}
	if baru == lama {
		return Profil{}, Token{}, ErrSandiSama
	}
	hash, err := HashSandi(baru)
	if err != nil {
		return Profil{}, Token{}, err
	}
	if err := l.gudang.GantiSandi(ctx, a.ID, hash, tok.Versi); err != nil {
		return Profil{}, Token{}, err
	}
	a.WajibGantiSandi, a.VersiSesi = false, tok.Versi+1
	p, err := l.profil(ctx, a)
	if err != nil {
		return Profil{}, Token{}, err
	}
	tok.Versi = a.VersiSesi
	return p, tok.Perpanjang(l.sekarang()), nil
}

// periksaJenjang - setiap CODE ada dan aktif, dan jenjangnya cocok.
//
// ⛔ Hanya saat MENYIMPAN akun (keputusan work owner 01-10-2026): login tidak
// memeriksanya - orang tidak boleh terkunci karena master berubah.
func (l *Layanan) periksaJenjang(ctx context.Context, a AkunBaru) error {
	if (a.Unit != "" && a.Divisi == "") || (a.Divisi != "" && a.Organisasi == "") {
		return fmt.Errorf("%w: tingkat bawah diisi tanpa tingkat atasnya", ErrJenjangTidakCocok)
	}
	if a.Organisasi != "" {
		aktif, err := l.gudang.InfoOrganisasi(ctx, a.Organisasi)
		if err != nil {
			return fmt.Errorf("organisasi %s: %w", a.Organisasi, err)
		}
		if !aktif {
			return fmt.Errorf("organisasi %s: %w", a.Organisasi, ErrMasterTidakAda)
		}
	}
	if a.Divisi != "" {
		org, aktif, err := l.gudang.InfoDivisi(ctx, a.Divisi)
		if err != nil {
			return fmt.Errorf("divisi %s: %w", a.Divisi, err)
		}
		if !aktif {
			return fmt.Errorf("divisi %s: %w", a.Divisi, ErrMasterTidakAda)
		}
		if org != a.Organisasi {
			return fmt.Errorf("%w: divisi %s milik organisasi %s, bukan %s", ErrJenjangTidakCocok, a.Divisi, org, a.Organisasi)
		}
	}
	if a.Unit != "" {
		div, aktif, err := l.gudang.InfoUnit(ctx, a.Unit)
		if err != nil {
			return fmt.Errorf("unit %s: %w", a.Unit, err)
		}
		if !aktif {
			return fmt.Errorf("unit %s: %w", a.Unit, ErrMasterTidakAda)
		}
		if div != a.Divisi {
			return fmt.Errorf("%w: unit %s milik divisi %s, bukan %s", ErrJenjangTidakCocok, a.Unit, div, a.Divisi)
		}
	}
	for _, w := range a.Workbasket {
		aktif, err := l.gudang.WorkbasketAktif(ctx, w)
		if err != nil {
			return fmt.Errorf("workbasket %s: %w", w, err)
		}
		if !aktif {
			return fmt.Errorf("workbasket %s: %w", w, ErrMasterTidakAda)
		}
	}
	return nil
}

// periksaEmailBebas - email (sudah dirapikan) tidak dipakai akun LAIN, tanpa beda huruf (permintaan work owner
// 03-10-2026 "email sudah terdaftar"). Email kosong selalu bebas.
func (l *Layanan) periksaEmailBebas(ctx context.Context, pemilik, email string) error {
	if email == "" {
		return nil
	}
	pemakai, err := l.gudang.PemakaiEmail(ctx, email)
	if err != nil {
		return err
	}
	for _, id := range pemakai {
		if id != pemilik {
			return ErrEmailSudahTerdaftar
		}
	}
	return nil
}

// BuatPengguna membuat akun dengan sandi sementara - dikembalikan SEKALI
// untuk dicetak; akunnya wajib ganti sandi saat login pertama.
func (l *Layanan) BuatPengguna(ctx context.Context, a AkunBaru) (string, error) {
	sandi, err := SandiSementara()
	if err != nil {
		return "", err
	}
	if err := l.buat(ctx, a, sandi, true); err != nil {
		return "", err
	}
	return sandi, nil
}

// BuatPenggunaDenganSandi membuat akun dengan sandi yang DITETAPKAN operator
// (`-sandi-dari-stdin`) - aturan sandi tetap berlaku, dan karena sandinya
// bukan sementara, akunnya tidak wajib ganti sandi.
func (l *Layanan) BuatPenggunaDenganSandi(ctx context.Context, a AkunBaru, sandi string) error {
	if err := PeriksaSandiBaru(sandi); err != nil {
		return err
	}
	return l.buat(ctx, a, sandi, false)
}

func (l *Layanan) buat(ctx context.Context, a AkunBaru, sandi string, wajibGanti bool) error {
	a.ID, a.Nama = strings.TrimSpace(a.ID), strings.TrimSpace(a.Nama)
	if !AkunSah(a.ID) || a.Nama == "" || len(a.Nama) > 150 {
		return ErrAkunTidakSah
	}
	// Username sudah ada (permintaan work owner 03-10-2026): sama persis atau beda huruf saja (migrasi 905).
	if pemakai, err := l.gudang.PemakaiUsername(ctx, a.ID); err != nil {
		return err
	} else if len(pemakai) > 0 {
		return fmt.Errorf("%w: %s", ErrAkunSudahAda, strings.Join(pemakai, ", "))
	}
	if err := l.periksaEmailBebas(ctx, a.ID, a.Email); err != nil {
		return err
	}
	if err := l.periksaJenjang(ctx, a); err != nil {
		return err
	}
	hash, err := HashSandi(sandi)
	if err != nil {
		return err
	}
	return l.gudang.BuatAkun(ctx, a, hash, wajibGanti)
}

// irisan - isi `a` yang juga ada di `b`, urutan `a`; selalu terisi (bukan nil).
func irisan(a, b []string) []string {
	ada := map[string]bool{}
	for _, x := range b {
		ada[x] = true
	}
	out := []string{}
	for _, x := range a {
		if ada[x] {
			out = append(out, x)
		}
	}
	return out
}
