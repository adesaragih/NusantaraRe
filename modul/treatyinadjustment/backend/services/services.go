// Package services memegang aturan modul Treaty In Adjustment
// (`treatyinadjustment`).
package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/treatyinadjustment/backend/models"
	"nusantarare/modul/treatyinadjustment/backend/repository"
)

// ErrKontrakTidakAda - kontrak yang diminta tidak ada.
var ErrKontrakTidakAda = errors.New("kontrak tidak ada")

// ErrIDTidakSah - pengenal kontrak bukan bilangan positif.
var ErrIDTidakSah = errors.New("pengenal kontrak tidak sah")

// Service membawa akar layanan modul ini.
type Service struct{ *inti.Dasar }

// DariDasar merakit Service di atas akar yang perakit sediakan.
func DariDasar(d *inti.Dasar) *Service { return &Service{Dasar: d} }

// New membuat Service; db boleh nil (proses tanpa Oracle - rute menjawab 503).
func New(d *db.DB) *Service { return &Service{Dasar: inti.NewDasar(d)} }

// PunyaDatabase aman dipanggil pada nil.
func (s *Service) PunyaDatabase() bool { return s != nil && s.Dasar.PunyaDatabase() }

// Gudang adalah seluruh sentuhan basis data yang Layanan butuhkan.
type Gudang interface {
	DaftarKontrak(ctx context.Context) ([]models.Kontrak, error)
	RantaiVersi(ctx context.Context, idKontrak int64) ([]models.Versi, error)

	// Panel Attachment - `M_ATTACHMENTTREATY_2`, tabel WARISAN, BACA SAJA.
	BacaLampiranKontrak(ctx context.Context, masterID string) ([]models.BarisLampiranWarisan, error)
	BacaKatalogKategoriLampiran(ctx context.Context) (map[string]string, error)

	// Panel History - `T_VIEW_COMMENT`, BACA SAJA. Seam TERPISAH dari
	// lampiran: kontrak yang lampirannya gagal dibaca tetap harus
	// memperlihatkan riwayatnya, dan sebaliknya.
	BacaRiwayatKontrak(ctx context.Context, masterID string) ([]models.BarisRiwayatWarisan, error)

	// Layar Adjustment - `TREATY_IN_EDM` + `M_TREATY_IN_EDM`, BACA SAJA.
	DaftarPenyesuaianWarisan(ctx context.Context) ([]models.BarisPenyesuaian, error)
	// ⛔⛔ `BacaPenyesuaianWarisan` DICABUT dari seam ini 6 Oktober 2026 —
	// ia mengurai `M_TREATY_IN_EDM.JSONDATA`, dan pemilik proses melarang
	// keras menarik nilai dari JSONDATA. Penggantinya membaca TABEL
	// PENDARATAN yang sama dengan modul Treaty In.
	BacaPenyesuaianPendaratan(ctx context.Context, id string) (models.Penyesuaian, error)
}

// Layanan memegang aturan modul ini di atas satu Gudang.
type Layanan struct{ gudang Gudang }

// LayananOracle merakit Layanan di atas Oracle.
func LayananOracle(s *Service) *Layanan { return &Layanan{gudang: repository.Baru(s.DB())} }

// LayananDengan merakit Layanan di atas Gudang mana pun - dipakai uji.
func LayananDengan(g Gudang) *Layanan { return &Layanan{gudang: g} }

// DaftarKontrak membaca kepala seluruh kontrak.
func (l *Layanan) DaftarKontrak(ctx context.Context, p inti.Pelaku) ([]models.Kontrak, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return nil, err
	}
	return l.gudang.DaftarKontrak(ctx)
}

// RantaiVersi membaca rantai versi satu kontrak.
//
// Kontrak yang tidak ada menghasilkan ErrKontrakTidakAda, BUKAN daftar kosong.
// Keduanya berbeda arti: rantai kosong tidak mungkin - sebuah kontrak selalu
// punya sekurangnya versi pertamanya (tiket 14) - sehingga daftar kosong yang
// dikembalikan apa adanya akan terbaca di layar sebagai "kontrak ini tidak
// punya versi", padahal kontraknya sendiri yang tidak ada.
func (l *Layanan) RantaiVersi(ctx context.Context, p inti.Pelaku, idKontrak int64) ([]models.Versi, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return nil, err
	}
	if idKontrak <= 0 {
		return nil, fmt.Errorf("%w: %d", ErrIDTidakSah, idKontrak)
	}
	versi, err := l.gudang.RantaiVersi(ctx, idKontrak)
	if err != nil {
		return nil, err
	}
	if len(versi) == 0 {
		return nil, fmt.Errorf("%w: %d", ErrKontrakTidakAda, idKontrak)
	}
	return versi, nil
}

// Pesan mengembalikan kalimat yang layak sampai ke layar.
func Pesan(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

// LampiranKontrakWarisan menyusun panel Attachment satu kontrak.
//
// ⛔ Pengenalnya TEKS. `M_ATTACHMENTTREATY_2.TREATYID` adalah `VARCHAR2(100)`
// - berbeda dari `idKontrak` model baru yang `int64`. Mengubahnya menjadi
// angka di sini akan menolak pengenal warisan yang sah.
func (l *Layanan) LampiranKontrakWarisan(ctx context.Context, p inti.Pelaku, masterID string) (LampiranKontrak, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return LampiranKontrak{}, err
	}
	// ⛔ Pengenal KOSONG ditolak DI SINI, bukan diteruskan: kueri dengan
	// pengenal kosong mengembalikan nol baris, dan "tidak ada" adalah
	// jawaban yang berbeda dari "tidak ditanyakan".
	if strings.TrimSpace(masterID) == "" {
		return LampiranKontrak{}, fmt.Errorf("%w: pengenal kontrak kosong", ErrIDTidakSah)
	}

	berkas, err := l.gudang.BacaLampiranKontrak(ctx, masterID)
	if err != nil {
		return LampiranKontrak{}, err
	}
	katalog, err := l.gudang.BacaKatalogKategoriLampiran(ctx)
	if err != nil {
		return LampiranKontrak{}, err
	}
	return LampiranKontrak{
		Kategori: SusunKategoriLampiran(katalog, berkas),
		Berkas:   berkas,
	}, nil
}

// RiwayatKontrakWarisan membaca panel History satu kontrak.
//
// ⛔ Seam TERPISAH dari lampiran, dan itu disengaja: keduanya dibaca dari
// tabel yang berbeda, dan kegagalan salah satunya tidak boleh mengosongkan
// yang lain.
func (l *Layanan) RiwayatKontrakWarisan(ctx context.Context, p inti.Pelaku, masterID string) ([]models.BarisRiwayatWarisan, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return nil, err
	}
	if strings.TrimSpace(masterID) == "" {
		return nil, fmt.Errorf("%w: pengenal kontrak kosong", ErrIDTidakSah)
	}
	// ⛔ NOL terjemahan di sini, dan itu DIUKUR — bukan kelalaian.
	//
	// `T_VIEW_COMMENT.TANGGAL` memuat cap waktu Pega utuh, bentuk
	// `20190422T101800.000 GMT` — bukan `YYYYMMDD`. Penerjemah tanggal
	// modul Treaty In (`TanggalTampil`) hanya mengubah yang DELAPAN ANGKA
	// dan mengembalikan sisanya apa adanya, jadi pada kolom ini ia tidak
	// mengubah apa pun. Panel Information & Submit modul sebelah
	// menampilkan kolom yang sama apa adanya pula.
	//
	// Menulis penerjemah di sini karena itu akan menambah penerjemah
	// KELIMA yang nol bedanya — dan penerjemah yang tidak mengubah apa pun
	// adalah tempat yang menunggu seseorang membuatnya berbeda diam-diam.
	return l.gudang.BacaRiwayatKontrak(ctx, masterID)
}
