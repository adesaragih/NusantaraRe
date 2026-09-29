package services

// Penyimpanan berkas lampiran - tiket 12 Treaty Contract Out.
//
// Untuk apa berkas ini: DUA implementasi `KlienPenyimpananTCO`.
//
//  1. `PenyimpananLokalTCO` - STUB DEV. Berkas ditaruh di folder kita sendiri
//     di bawah `UNGGAHAN_DIR`. Bawaan (`PELAKSANA_STORAGE=stub`).
//  2. `PenyimpananJarakJauhTCO` - bentuk penyambungan nyata: alamat
//     di-resolve SAAT JALAN dari `M_LINK_SERVICE` (ADR-0013, `ServiceGoogle`
//     tanpa URL literal), token dari cache yang memperbarui SEBELUM kedaluwarsa
//     (AC 60), transport di balik `PengirimBerkasTCO`.
//
// Transport nyata: `tco_pengirim_storage.go` [keputusan work owner 29-09-2026,
// OQ-TCO-08]; dipilih `PenyimpananLampiranTCO` bila `PELAKSANA_STORAGE=nyata`.
// Tanpa pengirim, rangkaian jarak jauh gagal dengan
// `ErrPenyimpananBelumDisetujui` - permanen, tampil di outbox.
//
// ⛔ Nol URL, nol env var di berkas ini (ADR-0013); penjaga
// `TestNolAlamatLayananDiKode` dan `TestTCOLampiranTanpaAlamatLiteral`.
//
// Dibaca sesudah: tco_lampiran.go, efekkeluar.go.

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"time"
)

// kunciBerkasSah - kunci penyimpanan adalah heksa `IMAGEID`; apa pun selain
// huruf dan angka ditolak sebelum menyentuh jalur di disk.
var kunciBerkasSah = regexp.MustCompile(`^[A-Za-z0-9]{8,64}$`)

func periksaKunciBerkas(kunci string) error {
	if !kunciBerkasSah.MatchString(kunci) {
		return fmt.Errorf("%w: kunci berkas %q tidak berbentuk", ErrPermintaanTidakSah, kunci)
	}
	return nil
}

// --- stub lokal ------------------------------------------------------------

type penyimpananLokalTCO struct{ folder string }

// PenyimpananLokalTCO menyusun stub lokal di bawah `UNGGAHAN_DIR`.
//
// ⛔ INI STUB, dan namanya menyebutnya. `UNGGAHAN_DIR` kosong berarti setiap
// operasi gagal terang (`ErrUnggahanDirBelumDisetel`), bukan jatuh ke folder
// kerja siapa pun.
func PenyimpananLokalTCO(svc *Service) KlienPenyimpananTCO {
	if svc.unggahanDir == "" {
		return penyimpananLokalTCO{}
	}
	return PenyimpananLokalDi(filepath.Join(svc.unggahanDir, folderModulTCO, "penyimpanan"))
}

// PenyimpananLokalDi menyusun stub lokal di folder tertentu - dipakai uji.
func PenyimpananLokalDi(folder string) KlienPenyimpananTCO {
	return penyimpananLokalTCO{folder: folder}
}

func (p penyimpananLokalTCO) jalur(kunci string) (string, error) {
	if p.folder == "" {
		return "", ErrUnggahanDirBelumDisetel
	}
	if err := periksaKunciBerkas(kunci); err != nil {
		return "", err
	}
	return filepath.Join(p.folder, kunci), nil
}

// Simpan menulis lewat berkas sementara lalu mengganti namanya: kunci yang
// SAMA ditimpa utuh, tidak pernah menjadi dua berkas atau berkas setengah.
func (p penyimpananLokalTCO) Simpan(_ context.Context, kunci string, isi io.Reader, _ string) error {
	tujuan, err := p.jalur(kunci)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(p.folder, 0o750); err != nil {
		return fmt.Errorf("services: menyiapkan folder penyimpanan lokal: %w", err)
	}
	f, err := os.CreateTemp(p.folder, ".unggah-*")
	if err != nil {
		return fmt.Errorf("services: membuat berkas sementara penyimpanan: %w", err)
	}
	_, salinErr := io.Copy(f, isi)
	tutupErr := f.Close()
	if salinErr != nil || tutupErr != nil {
		_ = os.Remove(f.Name())
		return fmt.Errorf("services: menulis penyimpanan lokal: %w", errors.Join(salinErr, tutupErr))
	}
	if err := os.Rename(f.Name(), tujuan); err != nil {
		_ = os.Remove(f.Name())
		return fmt.Errorf("services: menaruh berkas di penyimpanan lokal: %w", err)
	}
	return nil
}

func (p penyimpananLokalTCO) Buka(_ context.Context, kunci string) (io.ReadCloser, error) {
	j, err := p.jalur(kunci)
	if err != nil {
		return nil, err
	}
	f, err := os.Open(j)
	if os.IsNotExist(err) {
		return nil, ErrBerkasTidakAdaDiPenyimpanan
	}
	return f, err
}

func (p penyimpananLokalTCO) Hapus(_ context.Context, kunci string) error {
	j, err := p.jalur(kunci)
	if err != nil {
		return err
	}
	if err := os.Remove(j); os.IsNotExist(err) {
		return ErrBerkasTidakAdaDiPenyimpanan
	} else if err != nil {
		return fmt.Errorf("services: menghapus dari penyimpanan lokal: %w", err)
	}
	return nil
}

func (p penyimpananLokalTCO) Ada(_ context.Context, kunci string) (bool, error) {
	j, err := p.jalur(kunci)
	if err != nil {
		return false, err
	}
	_, err = os.Stat(j)
	if os.IsNotExist(err) {
		return false, nil
	}
	return err == nil, err
}

// --- rangkaian jarak jauh --------------------------------------------------

// ErrTokenPenyimpananGagal - token penyimpanan gagal diambil; efeknya gagal
// TERLIHAT di outbox, bukan unggahan yang diam saja tidak terjadi (AC 60).
var ErrTokenPenyimpananGagal = errors.New("services: token penyimpanan gagal diambil")

// MarginTokenTCO - token diperbarui selagi masih sisa sekian.
//
// `[data DBA]` token hidup satu menit (`umurToken`); diperbarui 15 detik
// sebelum habis supaya unggahan yang panjang tidak berangkat dengan token
// yang mati di tengah jalan. `[keputusan kami]` untuk angkanya.
const MarginTokenTCO = 15 * time.Second

// SumberTokenTCO menerbitkan token baru beserta kedaluwarsanya.
type SumberTokenTCO interface {
	TokenBaru(ctx context.Context) (token string, kedaluwarsa time.Time, err error)
}

// CacheTokenTCO memakai ulang token sampai `margin` sebelum kedaluwarsa.
type CacheTokenTCO struct {
	sumber      SumberTokenTCO
	jam         func() time.Time
	margin      time.Duration
	mu          sync.Mutex
	token       string
	kedaluwarsa time.Time
}

// NewCacheTokenTCO menyusun cache-nya.
func NewCacheTokenTCO(sumber SumberTokenTCO, jam func() time.Time, margin time.Duration) *CacheTokenTCO {
	if jam == nil {
		jam = time.Now
	}
	return &CacheTokenTCO{sumber: sumber, jam: jam, margin: margin}
}

// Token mengembalikan token yang masih cukup umur, atau menerbitkan yang baru.
//
// ⛔ Kegagalan MENGOSONGKAN cache: token lama yang tersisa tidak boleh dipakai
// sesudah sumbernya berkata ada yang salah.
func (c *CacheTokenTCO) Token(ctx context.Context) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	kini := c.jam()
	if c.token != "" && kini.Add(c.margin).Before(c.kedaluwarsa) {
		return c.token, nil
	}
	c.token, c.kedaluwarsa = "", time.Time{}
	if c.sumber == nil {
		return "", fmt.Errorf("%w: sumber token belum dipasang", ErrTokenPenyimpananGagal)
	}
	tok, exp, err := c.sumber.TokenBaru(ctx)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrTokenPenyimpananGagal, err)
	}
	if tok == "" {
		return "", fmt.Errorf("%w: sumber mengembalikan token kosong", ErrTokenPenyimpananGagal)
	}
	c.token, c.kedaluwarsa = tok, exp
	return tok, nil
}

// PengirimBerkasTCO adalah transport ke penyimpanan jarak jauh.
//
// Implementasinya `NewPengirimBerkasHTTPTCO` (OQ-TCO-08). ⛔ Galatnya TIDAK
// BOLEH memuat alamat, token, atau garam.
type PengirimBerkasTCO interface {
	Kirim(ctx context.Context, alamat, token, kunci string, isi io.Reader, mime string) error
	Ambil(ctx context.Context, alamat, token, kunci string) (io.ReadCloser, error)
	Buang(ctx context.Context, alamat, token, kunci string) error
	Periksa(ctx context.Context, alamat, token, kunci string) (bool, error)
}

// PenyimpananJarakJauhTCO merangkai resolver, cache token, dan transport.
type PenyimpananJarakJauhTCO struct {
	resolver ResolverEndpoint
	token    *CacheTokenTCO
	pengirim PengirimBerkasTCO
}

// NewPenyimpananJarakJauhTCO menyusun rangkaiannya; `pengirim` nil berarti
// penyambungan belum disetujui.
func NewPenyimpananJarakJauhTCO(resolver ResolverEndpoint, token *CacheTokenTCO,
	pengirim PengirimBerkasTCO) *PenyimpananJarakJauhTCO {
	return &PenyimpananJarakJauhTCO{resolver: resolver, token: token, pengirim: pengirim}
}

// siapkan meresolve alamat SAAT JALAN lalu mengambil token.
//
// ⛔ Setiap panggilan meresolve ulang: isi `M_LINK_SERVICE` yang diganti DBA
// berlaku tanpa proses dijalankan ulang, dan alamat tidak pernah tersimpan di
// medan struct ini.
func (p *PenyimpananJarakJauhTCO) siapkan(ctx context.Context, kunci KunciLayanan) (string, string, error) {
	if p.pengirim == nil {
		return "", "", ErrPenyimpananBelumDisetujui
	}
	if p.resolver == nil {
		return "", "", ErrResolverBelumDiputuskan
	}
	alamat, err := p.resolver.Resolve(ctx, kunci)
	if err != nil {
		return "", "", err
	}
	if p.token == nil {
		return "", "", fmt.Errorf("%w: cache token belum dipasang", ErrTokenPenyimpananGagal)
	}
	tok, err := p.token.Token(ctx)
	if err != nil {
		return "", "", err
	}
	return alamat, tok, nil
}

func galatJarakJauh(kunci KunciLayanan, err error) error {
	// ⛔ Menyebut KUNCI layanan, tidak pernah alamatnya.
	return fmt.Errorf("services: penyimpanan jarak jauh (%s/%s): %w", kunci.Kategori1, kunci.Kategori2, err)
}

// Simpan - kunci `("Google", "upload")`.
func (p *PenyimpananJarakJauhTCO) Simpan(ctx context.Context, kunci string, isi io.Reader, mime string) error {
	alamat, tok, err := p.siapkan(ctx, KunciUnggahBerkas)
	if err != nil {
		return err
	}
	if err := p.pengirim.Kirim(ctx, alamat, tok, kunci, isi, mime); err != nil {
		return galatJarakJauh(KunciUnggahBerkas, err)
	}
	return nil
}

// Buka - kunci `("Google", "geturl")`.
func (p *PenyimpananJarakJauhTCO) Buka(ctx context.Context, kunci string) (io.ReadCloser, error) {
	alamat, tok, err := p.siapkan(ctx, KunciURLBerkas)
	if err != nil {
		return nil, err
	}
	rc, err := p.pengirim.Ambil(ctx, alamat, tok, kunci)
	if err != nil {
		if errors.Is(err, ErrBerkasTidakAdaDiPenyimpanan) {
			return nil, err
		}
		return nil, galatJarakJauh(KunciURLBerkas, err)
	}
	return rc, nil
}

// Hapus - kunci `("Google", "delete")`.
func (p *PenyimpananJarakJauhTCO) Hapus(ctx context.Context, kunci string) error {
	alamat, tok, err := p.siapkan(ctx, KunciHapusBerkas)
	if err != nil {
		return err
	}
	if err := p.pengirim.Buang(ctx, alamat, tok, kunci); err != nil {
		if errors.Is(err, ErrBerkasTidakAdaDiPenyimpanan) {
			return err
		}
		return galatJarakJauh(KunciHapusBerkas, err)
	}
	return nil
}

// Ada - kunci `("Google", "geturl")`.
func (p *PenyimpananJarakJauhTCO) Ada(ctx context.Context, kunci string) (bool, error) {
	alamat, tok, err := p.siapkan(ctx, KunciURLBerkas)
	if err != nil {
		return false, err
	}
	ada, err := p.pengirim.Periksa(ctx, alamat, tok, kunci)
	if err != nil {
		return false, galatJarakJauh(KunciURLBerkas, err)
	}
	return ada, nil
}
