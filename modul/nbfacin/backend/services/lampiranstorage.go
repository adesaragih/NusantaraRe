package services

// Penyimpanan berkas lampiran NB FacIn (tiket 67) - SALINAN mesin Treaty Contract Out (keputusan work owner 08-10-2026 K-B
// "Salin ke NB FacIn"; modul/X tidak boleh mengimpor modul/Y), diberi bentuk Pega NB `[terverifikasi]` `NB FacIn\Activity\`:
//
//	InsertGoogleStorage_Act  POST ("Google","upload") {App, Kodestring, Durasi 1800, Folder "Policy/Doc/YYYY/MM/",
//	                         Namafile "yyyyMMdd-hhmmss-S - <nama>", Image base64, ext, MimeType}; URLImage kosong = gagal
//	GetUrlGoogleStorage_Act  URLPUBLIC tersimpan dipakai selama EXPDATE belum lewat; bila lewat POST ("Google","geturl")
//	                         {App, Kodestring, Folder, Namafile, Durasi} lalu Update_T_Storage_SQL; EXPDATE kosong = tanpa URL
//	DeleteGoogleStorage_Act  POST ("Google","delete") {App, Kodestring, Namafile = APPFOLDER tanpa "gs://<App>/"}
//	ServiceGoogle            POST JSON, batas waktu 300000 ms; alamat dari M_LINK_SERVICE saat jalan (ADR-0013)
//
// Dua pelaksana menurut PELAKSANA_STORAGE (config inti yang sudah ada): `stub` (bawaan) = folder lokal di bawah
// UNGGAHAN_DIR; `nyata` = layanan penyimpanan. URL bertanda tangan DIPROKSI backend (Pega membukanya di peramban).
//
// Transport HTTP-nya di lampiranpengirim.go (satu berkas klien keluar). ⛔ Nol URL literal, nol env var di berkas ini.

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/layanan"
	"nusantarare/inti/backend/unggah"
	"nusantarare/modul/nbfacin/backend/models"
)

const (
	// BatasWaktuStorageLampiran - ServiceGoogle `pyResponseTimeout 300000`.
	BatasWaktuStorageLampiran = 300 * time.Second
	// DurasiURLStorageLampiran - InsertDocument_Act / DownloadDocumentPolis `Durasi = 1800`.
	DurasiURLStorageLampiran = 1800
	// MarginTokenStorage - token diperbarui 15 detik sebelum habis (pola Treaty Contract Out, `[keputusan kami]`).
	MarginTokenStorage = 15 * time.Second
	// sisaURLTersimpan - URL tersimpan dipakai hanya bila masih berlaku semenit lagi `[keputusan agent]`.
	sisaURLTersimpan     = time.Minute
	batasJawabanStorage  = 1 << 20
	folderStubLampiran   = "nbfacin"
	subfolderPenyimpanan = "penyimpanan"
)

var (
	// ErrBerkasTidakAdaDiPenyimpanan - objek tidak ada (URL bertanda tangan 404 / berkas stub tidak ada).
	ErrBerkasTidakAdaDiPenyimpanan = errors.New("services: berkas tidak ada di penyimpanan")
	// ErrObjekNyataDiStub - pelaksana stub diminta menghapus objek yang tercatat di layanan nyata (APPFOLDER terisi):
	// tidak dihapus, catatannya tidak dibuang (PERMANEN).
	ErrObjekNyataDiStub = errors.New("services: pelaksana penyimpanan stub tidak dapat menghapus objek di layanan nyata")
	// ErrTokenPenyimpananGagal - token penyimpanan gagal diambil.
	ErrTokenPenyimpananGagal = errors.New("services: token penyimpanan gagal diambil")
	// ErrStorageTakTerjangkau - jaringan / batas waktu; layak dicoba ulang.
	ErrStorageTakTerjangkau = errors.New("services: layanan penyimpanan tidak terjangkau")
	// ErrStorageGagal - layanan menjawab galat (5xx, 404 titik layanan, 3xx, 401, 403, 429); layak dicoba ulang.
	ErrStorageGagal = errors.New("services: layanan penyimpanan menjawab galat")
	// ErrStorageTokenDitolak - 401/403: cache token dikosongkan. Membungkus ErrStorageGagal.
	ErrStorageTokenDitolak = fmt.Errorf("%w: token ditolak", ErrStorageGagal)
	// ErrStorageMenolakPermintaan - 400 / 422 / alamat rusak: PERMANEN.
	ErrStorageMenolakPermintaan = errors.New("services: layanan penyimpanan menolak permintaan")
	// ErrStorageJawabanRusak - URLImage kosong / bukan https / JSON rusak.
	ErrStorageJawabanRusak = errors.New("services: jawaban layanan penyimpanan rusak")
)

// KlienStorageLampiran - penyimpanan berkas di balik antarmuka.
type KlienStorageLampiran interface {
	// Simpan - unggah satu berkas ke lokasinya; hasilnya dicatat ke T_STORAGE_IMAGE.
	Simpan(ctx context.Context, lok models.LokasiObjek, isi io.Reader, ekstensi, mime string) (models.ObjekStorage, error)
	// Buka - isi objek tercatat (link = baris T_STORAGE_IMAGE).
	Buka(ctx context.Context, link models.LinkStorage) (io.ReadCloser, error)
	// Hapus - objek tercatat; ErrBerkasTidakAdaDiPenyimpanan bila sudah tidak ada.
	Hapus(ctx context.Context, link models.LinkStorage) error
	// Ada - objek di lokasi itu ada? Bila ada, objek geturl ikut dijawab (dicatat bila catatannya hilang).
	Ada(ctx context.Context, lok models.LokasiObjek) (bool, models.ObjekStorage, error)
}

// ---- stub lokal --------------------------------------------------------------------------------------------------

// kunciBerkasSah - kunci stub = IMAGEID heksa; apa pun selain huruf dan angka ditolak sebelum menyentuh disk.
var kunciBerkasSah = regexp.MustCompile(`^[A-Za-z0-9]{8,64}$`)

type storageLokal struct{ folder string }

// StorageLampiranLokal - STUB DEV di `<UNGGAHAN_DIR>/nbfacin/penyimpanan`; folder kosong = setiap operasi gagal terang.
func StorageLampiranLokal(unggahanDir string) KlienStorageLampiran {
	if strings.TrimSpace(unggahanDir) == "" {
		return storageLokal{}
	}
	return storageLokal{folder: filepath.Join(unggahanDir, folderStubLampiran, subfolderPenyimpanan)}
}

func (p storageLokal) jalur(kunci string) (string, error) {
	if p.folder == "" {
		return "", unggah.ErrUnggahanDirBelumDisetel
	}
	if !kunciBerkasSah.MatchString(kunci) {
		return "", fmt.Errorf("%w: kunci berkas %q tidak sah", ErrStorageMenolakPermintaan, kunci)
	}
	return filepath.Join(p.folder, kunci), nil
}

// Simpan - lewat berkas sementara lalu ganti nama: kunci yang sama ditimpa utuh. Stub tidak punya URL / folder aplikasi /
// kedaluwarsa: yang dicatat hanya IMAGEID dan Namafile.
func (p storageLokal) Simpan(_ context.Context, lok models.LokasiObjek, isi io.Reader, _, _ string) (models.ObjekStorage, error) {
	tujuan, err := p.jalur(lok.ImageID)
	if err != nil {
		return models.ObjekStorage{}, err
	}
	if err := os.MkdirAll(p.folder, 0o750); err != nil {
		return models.ObjekStorage{}, fmt.Errorf("services: menyiapkan folder stub penyimpanan: %w", err)
	}
	f, err := os.CreateTemp(p.folder, ".unggah-*")
	if err != nil {
		return models.ObjekStorage{}, fmt.Errorf("services: membuat berkas sementara stub: %w", err)
	}
	_, salinErr := io.Copy(f, isi)
	tutupErr := f.Close()
	if salinErr != nil || tutupErr != nil {
		_ = os.Remove(f.Name())
		return models.ObjekStorage{}, fmt.Errorf("services: menulis stub penyimpanan: %w", errors.Join(salinErr, tutupErr))
	}
	if err := os.Rename(f.Name(), tujuan); err != nil {
		_ = os.Remove(f.Name())
		return models.ObjekStorage{}, fmt.Errorf("services: menaruh berkas di stub penyimpanan: %w", err)
	}
	return models.ObjekStorage{ImageID: lok.ImageID, Namafile: lok.Namafile}, nil
}

func (p storageLokal) Buka(_ context.Context, link models.LinkStorage) (io.ReadCloser, error) {
	j, err := p.jalur(link.ImageID)
	if err != nil {
		return nil, err
	}
	f, err := os.Open(j)
	if os.IsNotExist(err) {
		return nil, ErrBerkasTidakAdaDiPenyimpanan
	}
	return f, err
}

// Hapus - objek yang tercatat di layanan NYATA (APPFOLDER terisi: dokumen Pega / generator PDF) tidak disentuh stub.
func (p storageLokal) Hapus(_ context.Context, link models.LinkStorage) error {
	if strings.TrimSpace(link.AppFolder) != "" {
		return ErrObjekNyataDiStub
	}
	j, err := p.jalur(link.ImageID)
	if err != nil {
		return err
	}
	if err := os.Remove(j); os.IsNotExist(err) {
		return ErrBerkasTidakAdaDiPenyimpanan
	} else if err != nil {
		return fmt.Errorf("services: menghapus dari stub penyimpanan: %w", err)
	}
	return nil
}

func (p storageLokal) Ada(_ context.Context, lok models.LokasiObjek) (bool, models.ObjekStorage, error) {
	j, err := p.jalur(lok.ImageID)
	if err != nil {
		return false, models.ObjekStorage{}, err
	}
	if _, err := os.Stat(j); os.IsNotExist(err) {
		return false, models.ObjekStorage{}, nil
	} else if err != nil {
		return false, models.ObjekStorage{}, err
	}
	return true, models.ObjekStorage{ImageID: lok.ImageID, Namafile: lok.Namafile}, nil
}

// ---- token ---------------------------------------------------------------------------------------------------------

// SumberTokenStorage menerbitkan token baru beserta kedaluwarsanya.
type SumberTokenStorage interface {
	TokenBaru(ctx context.Context) (token string, kedaluwarsa time.Time, err error)
}

// CacheTokenStorage memakai ulang token sampai `margin` sebelum kedaluwarsa; gagal = cache dikosongkan.
type CacheTokenStorage struct {
	sumber      SumberTokenStorage
	jam         func() time.Time
	margin      time.Duration
	mu          sync.Mutex
	token       string
	kedaluwarsa time.Time
}

// NewCacheTokenStorage menyusun cache-nya.
func NewCacheTokenStorage(sumber SumberTokenStorage, jam func() time.Time, margin time.Duration) *CacheTokenStorage {
	if jam == nil {
		jam = time.Now
	}
	return &CacheTokenStorage{sumber: sumber, jam: jam, margin: margin}
}

// Lupakan mengosongkan cache - dipanggil saat layanan menolak token.
func (c *CacheTokenStorage) Lupakan() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.token, c.kedaluwarsa = "", time.Time{}
}

// Token - token yang masih cukup umur, atau yang baru.
func (c *CacheTokenStorage) Token(ctx context.Context) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.token != "" && c.jam().Add(c.margin).Before(c.kedaluwarsa) {
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
		return "", fmt.Errorf("%w: sumber memberi token kosong", ErrTokenPenyimpananGagal)
	}
	c.token, c.kedaluwarsa = tok, exp
	return tok, nil
}

// PenyimpanTokenStorage - bacaan dan tulisan token GCP_IMAGE.
type PenyimpanTokenStorage interface {
	AppStorage(ctx context.Context) (string, error)
	TokenStorageBerlaku(ctx context.Context, tx *db.Tx, app string, saat time.Time, sisaMinimum time.Duration) (string, time.Duration, error)
	SimpanToken(ctx context.Context, tx *db.Tx, app, token, pengguna string, sampai time.Time) error
}

// sumberTokenStorage - GET_TOKEN_STORAGE ditiru tanpa CALL (ADR-0043; `inti/backend/layanan`): token berlaku dipakai ulang
// beserta sisa umurnya (hanya yang sisanya > margin - penyimpangan sadar kecil pola Treaty Contract Out); bila tidak ada,
// dirakit dengan garam STORAGE_TOKEN_SALT dan disimpan berumur satu menit.
type sumberTokenStorage struct {
	transaksi Transaksi
	penyimp   PenyimpanTokenStorage
	garam     string
	jam       func() time.Time
}

func (s sumberTokenStorage) TokenBaru(ctx context.Context) (string, time.Time, error) {
	if strings.TrimSpace(s.garam) == "" {
		return "", time.Time{}, layanan.ErrGaramTokenKosong
	}
	if s.penyimp == nil || s.transaksi == nil {
		return "", time.Time{}, db.ErrTanpaOracle
	}
	jam := s.jam
	if jam == nil {
		jam = time.Now
	}
	app, err := s.penyimp.AppStorage(ctx)
	if err != nil {
		return "", time.Time{}, err
	}
	if strings.TrimSpace(app) == "" {
		return "", time.Time{}, layanan.ErrAppNameKosong
	}
	var tok string
	var sampai time.Time
	err = s.transaksi(ctx, func(tx *db.Tx) error {
		saat := jam()
		lama, sisa, err := s.penyimp.TokenStorageBerlaku(ctx, tx, app, saat, MarginTokenStorage)
		if err != nil {
			return err
		}
		if lama != "" {
			tok, sampai = lama, saat.Add(sisa)
			return nil
		}
		baru, err := layanan.RakitToken(s.garam, saat)
		if err != nil {
			return err
		}
		if err := s.penyimp.SimpanToken(ctx, tx, app, baru, layanan.PenggunaTokenBawaan, saat.Add(layanan.UmurToken)); err != nil {
			return err
		}
		tok, sampai = baru, saat.Add(layanan.UmurToken)
		return nil
	})
	if err != nil {
		return "", time.Time{}, err
	}
	return tok, sampai, nil
}
