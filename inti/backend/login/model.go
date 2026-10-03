package login

import (
	"context"
	"errors"
)

// Akun adalah satu baris M_LOGIN_GO.
type Akun struct {
	ID        string
	Nama      string
	HashSandi string
	// CODE master organisasi - bukan ID (keputusan work owner 01-10-2026).
	Organisasi, Divisi, Unit string
	Aktif                    bool
	// Terkunci dihitung OLEH ORACLE (`LOCKED_UNTIL > SYSDATE`): satu jam,
	// jam basis data, untuk mengunci dan untuk memeriksa kunci.
	Terkunci        bool
	WajibGantiSandi bool
	VersiSesi       int64
}

// AkunBaru adalah isian `-buat-pengguna` dan Kelola User.
type AkunBaru struct {
	ID, Nama                 string
	Organisasi, Divisi, Unit string
	Workbasket               []string
	// Menu - KODE menu (`M_LOGIN_GO_MENU`); kosong = akun tanpa satu layar pun.
	Menu []string
	// Kontak - email, nomor HP, NIK, jabatan (migrasi 904); opsional.
	Kontak
}

// Profil adalah identitas yang dikirim ke layar (`GET /api/auth/saya`).
//
// ⛔ Nol sandi, nol hash, nol versi sesi di sini.
type Profil struct {
	AkunID          string   `json:"akunId"`
	Nama            string   `json:"nama"`
	Peran           []string `json:"peran"`
	Organisasi      string   `json:"organisasi"`
	Divisi          string   `json:"divisi"`
	Unit            string   `json:"unit"`
	WajibGantiSandi bool     `json:"wajibGantiSandi"`
	// Menu - KODE menu yang boleh dibuka (`M_LOGIN_GO_MENU`, Kelola User
	// 01-10-2026): modul, dan `kelolauser` bagi admin.
	Menu []string `json:"menu"`
}

var (
	// ErrAkunTidakAda - LOGIN_ID tidak ada di M_LOGIN_GO.
	ErrAkunTidakAda = errors.New("login: akun tidak ada")
	// ErrMasterTidakAda - CODE organisasi atau WORKBASKET_ID tidak ada, atau nonaktif.
	ErrMasterTidakAda = errors.New("login: data master tidak ada atau nonaktif")
	// ErrMenuBelumDimigrasi - tabel M_LOGIN_GO_MENU belum ada: migrasi 903
	// belum dijalankan. Login menjawab 503 yang menyebutnya.
	ErrMenuBelumDimigrasi = errors.New("login: tabel M_LOGIN_GO_MENU belum ada - migrasi 903 belum dijalankan (-migrate, oleh work owner)")
)

// Gudang membaca dan menulis M_LOGIN_GO, M_LOGIN_GO_WORKBASKET, dan membaca
// master organisasi serta workbasket. Antarmuka supaya layanan diuji tanpa Oracle.
type Gudang interface {
	AmbilAkun(ctx context.Context, id string) (Akun, error)
	// Workbasket - WORKBASKET_ID akun itu yang masih aktif di M_WORKBASKET.
	Workbasket(ctx context.Context, id string) ([]string, error)
	// Menu - KODE menu akun itu (M_LOGIN_GO_MENU), berurutan.
	Menu(ctx context.Context, id string) ([]string, error)
	// CatatGagal menaikkan FAILED_COUNT dan mengunci sesudah BatasGagal.
	CatatGagal(ctx context.Context, id string) error
	// CatatBerhasil menolkan FAILED_COUNT, mencabut kunci, mengisi LAST_LOGIN.
	CatatBerhasil(ctx context.Context, id string) error
	// GantiSandi menulis hash baru dan menaikkan SESSION_VERSION - hanya bila
	// versinya masih versiLama (ErrSesiTidakSah bila tidak).
	GantiSandi(ctx context.Context, id, hash string, versiLama int64) error
	// NaikkanVersi mencabut seluruh sesi akun itu (logout).
	NaikkanVersi(ctx context.Context, id string) error
	// InfoUnit - CODE divisi pemilik unit itu, dan apakah unitnya aktif.
	InfoUnit(ctx context.Context, code string) (divisi string, aktif bool, err error)
	// InfoDivisi - CODE organisasi pemilik divisi itu, dan apakah aktif.
	InfoDivisi(ctx context.Context, code string) (organisasi string, aktif bool, err error)
	InfoOrganisasi(ctx context.Context, code string) (aktif bool, err error)
	WorkbasketAktif(ctx context.Context, id string) (bool, error)
	// BuatAkun menulis akun, workbasket, dan menunya dalam SATU transaksi;
	// `wajibGanti` mengisi MUST_CHANGE_PASSWORD.
	BuatAkun(ctx context.Context, a AkunBaru, hash string, wajibGanti bool) error
	// PemakaiUsername - LOGIN_ID akun yang username-nya sama dengan `id`, tanpa beda huruf (migrasi 905); urut,
	// kosong = belum dipakai.
	PemakaiUsername(ctx context.Context, id string) ([]string, error)
	// PemakaiEmail - LOGIN_ID akun yang EMAIL-nya sama dengan `email`, tanpa beda huruf (migrasi 905); urut,
	// kosong = belum dipakai.
	PemakaiEmail(ctx context.Context, email string) ([]string, error)
}
