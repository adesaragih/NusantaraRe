package services

// Lampiran tahun treaty - tiket 12 Treaty Contract Out (FITUR BARU).
//
// Untuk apa berkas ini: unggah / unduh / hapus / ulangi / periksa-selaras
// lampiran sebuah tahun treaty, lewat outbox bersama `T_LOG_SERVICE_RNM`.
//
// ⚠️ Penyimpangan sadar 9 `[keputusan work owner]`: jalur lampiran Pega
// (`TreatyOutSaveAttachment` -> `InsertAtatchment_Sql` -> `PEGA_M_ATTACHMENT`,
// dibaca dengan `GetAllAttachment2_Sql`) belum rampung: prosedurnya hanya
// meng-upsert simpanan JSON `M_ATTACHMENTTREATY`, daftarnya membaca
// `M_ATTACHMENTTREATY_2` (OQ-TCO-24 ditutup; simpanan JSON tidak diteruskan).
// Korpus adalah sumber RANTAI TEKNIS dan LABEL; perilakunya dari AC.
//
// ⛔ Taruhannya berbeda dari Komite Claim Life: lampiran OPSIONAL. Yang wajib
// adalah kegagalannya TERLIHAT dan DAPAT DIULANG (ADR-0015), bukan berhasil.
// Unggahan yang gagal tidak membatalkan apa pun - hanya statusnya.
//
// ⛔ PEKERJA MILIK MODUL INI SENDIRI, bukan `PekerjaEfek`. Jalur menyerah
// `PekerjaEfek` menulis jejak Claim Life (`PerekamJejakOracle` ->
// `KlaimLife.SisipJejak`), dan itu penulisan lintas modul. Yang dipakai ulang:
// fungsi outbox repository (`AntreEfek`, `PungutEfek`, `TuntaskanEfek`),
// `Backoff`, `LayakDicobaUlang`, `percobaanMaksimum`, gerbang batas ukuran
// `Unggahan.tulisBerkas`, `models.MimeDokumen`, `models.ImageIDBaru`.
//
// ⛔ Endpoint penyimpanan nyata TIDAK dipanggil. Di DEV penyimpanannya stub
// lokal di bawah `UNGGAHAN_DIR` (tco_penyimpanan.go).
//
// Dibaca sesudah: unggahan.go, antrean.go, tco_tahun.go.

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"nusantarare/inti"
	"nusantarare/inti/db"
	"nusantarare/inti/galat"
	"nusantarare/inti/layanan"
	"nusantarare/inti/outbox"
	"nusantarare/inti/unggah"
	"nusantarare/inti/utils"
	"nusantarare/modul/treaty/models"
	"nusantarare/modul/treaty/repository"
)

const (
	// ModulTreatyContractOut mengisi kolom `MODUL` outbox; pekerja modul ini
	// hanya memungut baris bermodul ini.
	ModulTreatyContractOut = "TREATYCONTRACTOUT"
	// LiniNonLife mengisi kolom `LINI` outbox. `[keputusan kami]` - kolom
	// milik sistem baru, Treaty Contract Out adalah treaty non-life
	// (saringan jenis reasuransi tiket 02).
	LiniNonLife = "NONLIFE"
	// PesanTanpaBerkasTCO - VERBATIM `TreatyOutSaveAttachment.xml` b376.
	PesanTanpaBerkasTCO = "Tidak ada file yg diattach"

	folderModulTCO       = "treaty-contract-out"
	putaranSeusaiAksiTCO = 10
	batasGalatTCO        = 1000
)

var (
	// ErrGudangLampiranBelumDisuntik - gudang rekam lampiran belum dipasang.
	ErrGudangLampiranBelumDisuntik = errors.New("services: gudang lampiran tahun treaty belum disuntik")
	// ErrAntreanLampiranBelumDisuntik - outbox lampiran belum dipasang.
	ErrAntreanLampiranBelumDisuntik = errors.New("services: antrean lampiran tahun treaty belum disuntik")
	// ErrPenyimpananLampiranBelumDisuntik - klien penyimpanan belum dipasang.
	ErrPenyimpananLampiranBelumDisuntik = errors.New("services: penyimpanan lampiran belum disuntik")
	// ErrPembacaKategoriLampiranBelumDisuntik - master kategori belum dipasang.
	ErrPembacaKategoriLampiranBelumDisuntik = errors.New("services: pembaca kategori lampiran belum disuntik")
	// ErrKategoriLampiranKosong - master kategori kosong: keadaan server (503).
	ErrKategoriLampiranKosong = errors.New(
		"services: master kategori lampiran CATEGORY_ATTACH_REAS kosong; lampiran belum dapat diberi kategori")
	// ErrKategoriLampiranTidakDikenal - kategori di luar master (422).
	ErrKategoriLampiranTidakDikenal = errors.New("services: kategori lampiran tidak ada di master CATEGORY_ATTACH_REAS")
	// ErrLampiranTidakAda - bukan lampiran tahun treaty itu (404).
	ErrLampiranTidakAda = repository.ErrLampiranTidakAda
	// ErrLampiranBelumTerkirim - berkasnya belum dipastikan ada di penyimpanan (409).
	ErrLampiranBelumTerkirim = errors.New("services: lampiran belum terkirim ke penyimpanan")
	// ErrLampiranSudahTerkirim - "ulangi" atas lampiran yang berkasnya ada (409).
	ErrLampiranSudahTerkirim = errors.New("services: lampiran sudah terkirim dan berkasnya ada di penyimpanan")
	// ErrLampiranTanpaBerkas - rekam terkirim tetapi berkasnya tidak ada (AC 61).
	ErrLampiranTanpaBerkas = errors.New(
		"services: rekam lampiran ada tetapi berkasnya tidak ada di penyimpanan; periksa keselarasan")
	// ErrBerkasSumberLampiranHilang - berkas antrean tidak ada lagi dan
	// penyimpanan tidak memilikinya: tidak dapat diunggah ulang (permanen).
	ErrBerkasSumberLampiranHilang = errors.New(
		"services: berkas sumber lampiran tidak ada lagi; hapus rekamnya lalu unggah ulang")
	// ErrBerkasTidakAdaDiPenyimpanan - dijawab klien penyimpanan untuk kunci
	// yang tidak ada. Penghapusan menganggapnya selesai.
	ErrBerkasTidakAdaDiPenyimpanan = errors.New("services: berkas tidak ada di penyimpanan")

	errMuatanLampiranRusak = errors.New("services: muatan outbox lampiran tidak terbaca")
	errJenisEfekAsingTCO   = errors.New("services: jenis efek bukan milik pelaksana lampiran")
)

// GudangLampiranTCO membaca dan menulis rekam lampiran + catatan objeknya (tco4).
type GudangLampiranTCO interface {
	Daftar(ctx context.Context, tahunID string) ([]repository.BarisLampiranTCO, error)
	Ambil(ctx context.Context, tahunID, id string) (repository.BarisLampiranTCO, error)
	AmbilUntukKirim(ctx context.Context, tx *db.Tx, id string) (models.LampiranTCO, error)
	Sisip(ctx context.Context, tx *db.Tx, l models.LampiranTCO) (string, error)
	Hapus(ctx context.Context, tx *db.Tx, tahunID, id string) error
	// SimpanObjek - `Insert_T_Storage_SQL`: objek berkas tercatat = terkirim (tco4).
	SimpanObjek(ctx context.Context, tx *db.Tx, o models.ObjekPenyimpananTCO) error
	// HapusObjek - `DeleteStorage_SQL`: catatan objek dibuang.
	HapusObjek(ctx context.Context, tx *db.Tx, imageID string) error
}

// AntreanLampiranTCO adalah outbox yang dipakai lampiran, sudah bermodul.
type AntreanLampiranTCO interface {
	Antre(ctx context.Context, tx *db.Tx, jenis, rujukan, muatan string, saat time.Time) error
	Pungut(ctx context.Context, tx *db.Tx, saat time.Time) (outbox.BarisEfekKeluar, error)
	// PungutRujukan - satu efek jatuh tempo milik SATU rujukan (lampiran).
	PungutRujukan(ctx context.Context, tx *db.Tx, rujukan string, saat time.Time) (outbox.BarisEfekKeluar, error)
	Tuntaskan(ctx context.Context, tx *db.Tx, id, status string, jadwal time.Time,
		galat string, saat time.Time) error
}

// PembacaKategoriLampiranTCO membaca master kategori.
type PembacaKategoriLampiranTCO interface {
	Daftar(ctx context.Context) ([]string, error)
}

// PemeriksaTahunTCO memastikan tahun treaty induknya ada.
type PemeriksaTahunTCO interface {
	Ambil(ctx context.Context, id string) (models.TahunTreaty, error)
}

// KlienPenyimpananTCO adalah penyimpanan berkas di balik antarmuka (AC 62).
//
// ⛔ `Hapus` dan `Buka` atas kunci yang tidak ada menjawab
// `ErrBerkasTidakAdaDiPenyimpanan`; toleransinya tinggal di berkas ini, bukan
// diserahkan ke tiap klien.
type KlienPenyimpananTCO interface {
	// ekstensi - dari nama berkas asli, huruf kecil tanpa titik (`ext` UploadDoc).
	// Hasilnya yang dicatat ke `T_STORAGE_IMAGE` (tco4).
	Simpan(ctx context.Context, kunci string, isi io.Reader, mime, ekstensi string) (models.ObjekPenyimpananTCO, error)
	Buka(ctx context.Context, kunci string) (io.ReadCloser, error)
	Hapus(ctx context.Context, kunci string) error
	Ada(ctx context.Context, kunci string) (bool, error)
}

// --- bawaan yang gagal terang ----------------------------------------------

type gudangLampiranBelumDisuntik struct{}

func (gudangLampiranBelumDisuntik) Daftar(context.Context, string) ([]repository.BarisLampiranTCO, error) {
	return nil, ErrGudangLampiranBelumDisuntik
}
func (gudangLampiranBelumDisuntik) Ambil(context.Context, string, string) (repository.BarisLampiranTCO, error) {
	return repository.BarisLampiranTCO{}, ErrGudangLampiranBelumDisuntik
}
func (gudangLampiranBelumDisuntik) AmbilUntukKirim(context.Context, *db.Tx, string) (models.LampiranTCO, error) {
	return models.LampiranTCO{}, ErrGudangLampiranBelumDisuntik
}
func (gudangLampiranBelumDisuntik) Sisip(context.Context, *db.Tx, models.LampiranTCO) (string, error) {
	return "", ErrGudangLampiranBelumDisuntik
}
func (gudangLampiranBelumDisuntik) Hapus(context.Context, *db.Tx, string, string) error {
	return ErrGudangLampiranBelumDisuntik
}
func (gudangLampiranBelumDisuntik) SimpanObjek(context.Context, *db.Tx, models.ObjekPenyimpananTCO) error {
	return ErrGudangLampiranBelumDisuntik
}
func (gudangLampiranBelumDisuntik) HapusObjek(context.Context, *db.Tx, string) error {
	return ErrGudangLampiranBelumDisuntik
}

type antreanLampiranBelumDisuntik struct{}

func (antreanLampiranBelumDisuntik) Antre(context.Context, *db.Tx, string, string, string, time.Time) error {
	return ErrAntreanLampiranBelumDisuntik
}
func (antreanLampiranBelumDisuntik) Pungut(context.Context, *db.Tx, time.Time) (outbox.BarisEfekKeluar, error) {
	return outbox.BarisEfekKeluar{}, ErrAntreanLampiranBelumDisuntik
}
func (antreanLampiranBelumDisuntik) PungutRujukan(context.Context, *db.Tx, string, time.Time) (outbox.BarisEfekKeluar, error) {
	return outbox.BarisEfekKeluar{}, ErrAntreanLampiranBelumDisuntik
}
func (antreanLampiranBelumDisuntik) Tuntaskan(context.Context, *db.Tx, string, string, time.Time, string, time.Time) error {
	return ErrAntreanLampiranBelumDisuntik
}

type penyimpananBelumDisuntik struct{}

func (penyimpananBelumDisuntik) Simpan(context.Context, string, io.Reader, string, string) (models.ObjekPenyimpananTCO, error) {
	return models.ObjekPenyimpananTCO{}, ErrPenyimpananLampiranBelumDisuntik
}
func (penyimpananBelumDisuntik) Buka(context.Context, string) (io.ReadCloser, error) {
	return nil, ErrPenyimpananLampiranBelumDisuntik
}
func (penyimpananBelumDisuntik) Hapus(context.Context, string) error {
	return ErrPenyimpananLampiranBelumDisuntik
}
func (penyimpananBelumDisuntik) Ada(context.Context, string) (bool, error) {
	return false, ErrPenyimpananLampiranBelumDisuntik
}

type kategoriLampiranBelumDisuntik struct{}

func (kategoriLampiranBelumDisuntik) Daftar(context.Context) ([]string, error) {
	return nil, ErrPembacaKategoriLampiranBelumDisuntik
}

// --- implementasi Oracle ---------------------------------------------------

type gudangLampiranOracle struct {
	m  *repository.MasterLampiranTCO
	db *db.DB
}

func (g gudangLampiranOracle) Daftar(ctx context.Context, tahunID string) ([]repository.BarisLampiranTCO, error) {
	return g.m.Daftar(ctx, ModulTreatyContractOut, unggah.JenisEfekStorageUnggah, tahunID)
}
func (g gudangLampiranOracle) Ambil(ctx context.Context, tahunID, id string) (repository.BarisLampiranTCO, error) {
	return g.m.Ambil(ctx, ModulTreatyContractOut, unggah.JenisEfekStorageUnggah, tahunID, id)
}
func (g gudangLampiranOracle) AmbilUntukKirim(ctx context.Context, tx *db.Tx, id string) (models.LampiranTCO, error) {
	return g.m.AmbilUntukKirim(ctx, tx, id)
}
func (g gudangLampiranOracle) Sisip(ctx context.Context, tx *db.Tx, l models.LampiranTCO) (string, error) {
	return g.m.Sisip(ctx, tx, l)
}
func (g gudangLampiranOracle) Hapus(ctx context.Context, tx *db.Tx, tahunID, id string) error {
	return g.m.Hapus(ctx, tx, tahunID, id)
}
func (g gudangLampiranOracle) SimpanObjek(ctx context.Context, tx *db.Tx, o models.ObjekPenyimpananTCO) error {
	return g.m.SimpanObjek(ctx, tx, o)
}
func (g gudangLampiranOracle) HapusObjek(ctx context.Context, tx *db.Tx, imageID string) error {
	return g.m.HapusObjek(ctx, tx, imageID)
}

// GudangLampiranOracle menyusun gudang lampiran di atas Oracle.
func GudangLampiranOracle(svc *Service) GudangLampiranTCO {
	return gudangLampiranOracle{m: repository.NewMasterLampiranTCO(svc.DB()), db: svc.DB()}
}

type antreanLampiranOracle struct {
	outbox *outbox.Penyimpan
	db     *db.DB
}

func (a antreanLampiranOracle) Antre(ctx context.Context, tx *db.Tx, jenis, rujukan, muatan string,
	saat time.Time) error {
	_, err := a.outbox.AntreEfek(ctx, tx, LiniNonLife, ModulTreatyContractOut, jenis, rujukan, muatan, saat)
	return err
}
func (a antreanLampiranOracle) Pungut(ctx context.Context, tx *db.Tx, saat time.Time) (
	outbox.BarisEfekKeluar, error) {
	return a.outbox.PungutEfek(ctx, tx, ModulTreatyContractOut, saat)
}
func (a antreanLampiranOracle) PungutRujukan(ctx context.Context, tx *db.Tx, rujukan string, saat time.Time) (
	outbox.BarisEfekKeluar, error) {
	return repository.PungutEfekRujukanTCO(ctx, a.db, tx, ModulTreatyContractOut, rujukan, saat)
}
func (a antreanLampiranOracle) Tuntaskan(ctx context.Context, tx *db.Tx, id, status string,
	jadwal time.Time, galat string, saat time.Time) error {
	return a.outbox.TuntaskanEfek(ctx, tx, id, status, jadwal, galat, saat)
}

// AntreanLampiranOracle menyusun outbox lampiran di atas `T_LOG_SERVICE_RNM`.
//
// ⚠️ Memakai fungsi outbox yang sudah ada di `PohonKlaim` - tabelnya bersama
// lintas modul, dipisah kolom `MODUL`. Tidak ada SQL outbox baru.
func AntreanLampiranOracle(svc *Service) AntreanLampiranTCO {
	return antreanLampiranOracle{outbox: outbox.NewPenyimpan(svc.DB()), db: svc.DB()}
}

type kategoriLampiranOracle struct{ k *repository.KategoriLampiran }

func (k kategoriLampiranOracle) Daftar(ctx context.Context) ([]string, error) { return k.k.Daftar(ctx) }

// KategoriLampiranOracle menyusun pembaca master `CATEGORY_ATTACH_REAS`.
func KategoriLampiranOracle(svc *Service) PembacaKategoriLampiranTCO {
	return kategoriLampiranOracle{k: repository.NewKategoriLampiran(svc.DB())}
}

// --- bentuk jawaban --------------------------------------------------------

// LampiranTampil adalah satu lampiran seperti dikirim ke layar.
//
// ⛔ `IMAGEID` TIDAK dikirim: ia kunci berkas di penyimpanan, dan kunci yang
// beredar di layar bukan lagi kunci.
type LampiranTampil struct {
	ID           string `json:"id"`
	IDTreatyYear string `json:"idTreatyYear"`
	FileName     string `json:"fileName"`
	FileMimeType string `json:"fileMimeType"`
	Category     string `json:"category"`
	UserID       string `json:"userId"`
	TglUpload    string `json:"tglUpload"`
	// Status - `terkirim` / `tertunda` / `gagal` (models.StatusLampiranTCO).
	Status string `json:"status"`
	// Percobaan dan Galat - nasib efek unggah terakhir; galat kosong bila terkirim.
	Percobaan int    `json:"percobaan"`
	Galat     string `json:"galat"`
}

// TampilLampiran menerjemahkan satu baris gudang.
func TampilLampiran(b repository.BarisLampiranTCO) LampiranTampil {
	status := models.StatusLampiranTCO(b.TStorageID, b.StatusEfek == outbox.StatusEfekGagalPermanen)
	galat := ""
	if status != models.StatusLampiranTerkirim {
		galat = b.GalatEfek
	}
	return LampiranTampil{ID: b.ID, IDTreatyYear: b.IDTreatyYear, FileName: b.FileName,
		FileMimeType: b.FileMimeType, Category: b.Category, UserID: b.UserID,
		TglUpload: utils.FormatTanggalWaktu(b.TglUpload), Status: status,
		Percobaan: b.PercobaanEfek, Galat: galat}
}

// HasilLampiranTCO adalah jawaban unggah / ulangi.
//
// ⚠️ `Peringatan` terisi bila antrean pengiriman sendiri tidak dapat
// dijalankan (bukan bila efeknya gagal - itu tampil di `Status`/`Galat`).
// Lampirannya tetap tersimpan; peringatan itu supaya kegagalannya tidak diam.
type HasilLampiranTCO struct {
	Lampiran   LampiranTampil `json:"lampiran"`
	Peringatan string         `json:"peringatan"`
}

// TemuanSelarasTCO adalah satu rekam yang tidak sejalan dengan penyimpanan.
type TemuanSelarasTCO struct {
	LampiranID string `json:"lampiranId"`
	FileName   string `json:"fileName"`
	Masalah    string `json:"masalah"`
	// Perbaikan - `ulangi` (berkas antrean masih ada) atau `hapus`.
	Perbaikan string `json:"perbaikan"`
}

// muatanLampiranTCO adalah muatan JSON outbox lampiran.
//
// ⛔ PENGENAL SAJA. Nol nama orang, nol alamat, nol kredensial; nama berkas
// antrean adalah kunci + akhiran (`models.NamaBerkasAntreLampiranTCO`), bukan
// nama unggahan.
type muatanLampiranTCO struct {
	TahunID     string `json:"tahunId"`
	LampiranID  string `json:"lampiranId"`
	ImageID     string `json:"imageId"`
	BerkasAntre string `json:"berkasAntre"`
}

func (m muatanLampiranTCO) teks() (string, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return "", fmt.Errorf("services: merakit muatan lampiran: %w", err)
	}
	return string(b), nil
}

// --- layanan ---------------------------------------------------------------

// LampiranTahunTCO melayani panel lampiran tahun treaty.
type LampiranTahunTCO struct {
	svc         *Service
	gudang      GudangLampiranTCO
	antrean     AntreanLampiranTCO
	penyimpanan KlienPenyimpananTCO
	kategori    PembacaKategoriLampiranTCO
	tahun       PemeriksaTahunTCO
	folder      string
	jam         func() time.Time
	transaksi   func(ctx context.Context, fn func(tx *db.Tx) error) error
	batas       int64
}

// LampiranTahunTCO menyusun layanannya; seluruh bawaannya gagal terang.
func (s *Service) LampiranTahunTCO() *LampiranTahunTCO {
	return &LampiranTahunTCO{svc: s, gudang: gudangLampiranBelumDisuntik{},
		antrean: antreanLampiranBelumDisuntik{}, penyimpanan: penyimpananBelumDisuntik{},
		kategori: kategoriLampiranBelumDisuntik{}, tahun: gudangTahunTreatyBelumDisuntik{},
		folder: s.UnggahanDir(), jam: time.Now, transaksi: s.DalamTransaksi, batas: unggah.BatasUkuranUnggahan}
}

func (l *LampiranTahunTCO) salin() *LampiranTahunTCO { s := *l; return &s }

// DenganGudang memasang gudang rekam lampiran.
func (l *LampiranTahunTCO) DenganGudang(g GudangLampiranTCO) *LampiranTahunTCO {
	s := l.salin()
	s.gudang = g
	return s
}

// DenganAntrean memasang outbox lampiran.
func (l *LampiranTahunTCO) DenganAntrean(a AntreanLampiranTCO) *LampiranTahunTCO {
	s := l.salin()
	s.antrean = a
	return s
}

// DenganPenyimpanan memasang klien penyimpanan berkas.
func (l *LampiranTahunTCO) DenganPenyimpanan(p KlienPenyimpananTCO) *LampiranTahunTCO {
	s := l.salin()
	s.penyimpanan = p
	return s
}

// DenganKategori memasang pembaca master kategori.
func (l *LampiranTahunTCO) DenganKategori(k PembacaKategoriLampiranTCO) *LampiranTahunTCO {
	s := l.salin()
	s.kategori = k
	return s
}

// DenganTahun memasang pemeriksa tahun treaty induk.
func (l *LampiranTahunTCO) DenganTahun(t PemeriksaTahunTCO) *LampiranTahunTCO {
	s := l.salin()
	s.tahun = t
	return s
}

// DenganFolder mengganti `UNGGAHAN_DIR` - dipakai uji.
func (l *LampiranTahunTCO) DenganFolder(f string) *LampiranTahunTCO {
	s := l.salin()
	s.folder = f
	return s
}

// DenganJam mengganti sumber waktu - dipakai uji.
func (l *LampiranTahunTCO) DenganJam(j func() time.Time) *LampiranTahunTCO {
	s := l.salin()
	s.jam = j
	return s
}

// DenganTransaksi mengganti pelaksana transaksi - dipakai uji.
func (l *LampiranTahunTCO) DenganTransaksi(f func(ctx context.Context, fn func(tx *db.Tx) error) error) *LampiranTahunTCO {
	s := l.salin()
	s.transaksi = f
	return s
}

// DenganBatas mengecilkan batas ukuran - dipakai uji; nol/negatif diabaikan.
func (l *LampiranTahunTCO) DenganBatas(n int64) *LampiranTahunTCO {
	if n <= 0 {
		return l
	}
	s := l.salin()
	s.batas = n
	return s
}

func (l *LampiranTahunTCO) folderAntre() string {
	return filepath.Join(l.folder, folderModulTCO, "antre")
}

func (l *LampiranTahunTCO) jalurAntre(imageID, nama string) string {
	return filepath.Join(l.folderAntre(), models.NamaBerkasAntreLampiranTCO(imageID, nama))
}

func (l *LampiranTahunTCO) periksaTahun(ctx context.Context, tahunID string) error {
	if strings.TrimSpace(tahunID) == "" {
		return fmt.Errorf("%w: pengenal tahun treaty wajib diisi", galat.ErrPermintaanTidakSah)
	}
	_, err := l.tahun.Ambil(ctx, tahunID)
	return err
}

// Kategori membaca master kategori; kosong adalah kegagalan (503).
func (l *LampiranTahunTCO) Kategori(ctx context.Context, pelaku inti.Pelaku) ([]string, error) {
	if err := inti.WajibIdentitas(pelaku); err != nil {
		return nil, err
	}
	daftar, err := l.kategori.Daftar(ctx)
	if err != nil {
		return nil, err
	}
	if len(daftar) == 0 {
		return nil, ErrKategoriLampiranKosong
	}
	return daftar, nil
}

// Daftar membaca lampiran satu tahun treaty - tombol `Refresh` b1023.
func (l *LampiranTahunTCO) Daftar(ctx context.Context, pelaku inti.Pelaku, tahunID string) ([]LampiranTampil, error) {
	if err := inti.WajibIdentitas(pelaku); err != nil {
		return nil, err
	}
	baris, err := l.gudang.Daftar(ctx, tahunID)
	if err != nil {
		return nil, err
	}
	if err := l.periksaTahun(ctx, tahunID); err != nil {
		return nil, err
	}
	hasil := make([]LampiranTampil, 0, len(baris))
	for _, b := range baris {
		hasil = append(hasil, TampilLampiran(b))
	}
	return hasil, nil
}

// namaDasarBerkas membuang jalur yang dikirim sebagian peramban.
func namaDasarBerkas(nama string) string {
	d := path.Base(strings.ReplaceAll(strings.TrimSpace(nama), `\`, "/"))
	if d == "." || d == "/" || d == ".." {
		return ""
	}
	return d
}

// Unggah menerima satu berkas - `Add attachment` b578.
//
// Urutan dan sebabnya:
//
//  1. gerbang identitas, tahun, folder, nama, kategori - sebelum satu byte pun
//     mendarat;
//  2. berkas ke folder antrean LEBIH DULU, dinamai kunci berkasnya - byte yang
//     mendarat tidak pernah tanpa rekam;
//  3. satu transaksi: rekam `M_ATTACHMENTTREATY_2` + efek outbox (tco4: nol jejak modul);
//  4. sesudah commit, antrean dijalankan. Kegagalannya TIDAK membatalkan
//     langkah 3 (AC 55) - ia tampil sebagai status.
func (l *LampiranTahunTCO) Unggah(ctx context.Context, pelaku inti.Pelaku, tahunID string,
	berkas unggah.BerkasMasuk) (HasilLampiranTCO, error) {

	if err := inti.WajibIdentitas(pelaku); err != nil {
		return HasilLampiranTCO{}, err
	}
	if strings.TrimSpace(tahunID) == "" {
		return HasilLampiranTCO{}, fmt.Errorf("%w: pengenal tahun treaty wajib diisi", galat.ErrPermintaanTidakSah)
	}
	if strings.TrimSpace(l.folder) == "" {
		return HasilLampiranTCO{}, unggah.ErrUnggahanDirBelumDisetel
	}
	nama := namaDasarBerkas(berkas.NamaFile)
	if nama == "" || berkas.Isi == nil {
		return HasilLampiranTCO{}, fmt.Errorf("%w: %s", unggah.ErrBerkasKosong, PesanTanpaBerkasTCO)
	}
	master, err := l.Kategori(ctx, pelaku)
	if err != nil {
		return HasilLampiranTCO{}, err
	}
	kategori, sah := models.KategoriLampiranSah(master, berkas.Kategori)
	if !sah {
		return HasilLampiranTCO{}, fmt.Errorf("%w: %q", ErrKategoriLampiranTidakDikenal, berkas.Kategori)
	}
	if err := l.periksaTahun(ctx, tahunID); err != nil {
		return HasilLampiranTCO{}, err
	}

	saat := l.jam()
	imageID, err := unggah.ImageIDBaru(saat)
	if err != nil {
		return HasilLampiranTCO{}, err
	}
	jalur := l.jalurAntre(imageID, nama)
	// ⚠️ Gerbang ukuran DIPAKAI ULANG dari `Unggahan.tulisBerkas`: batasnya
	// ditegakkan saat menyalin, berkas kosong dan berlebih dibuang lagi.
	if err := unggah.TulisBerkas(jalur, berkas.Isi, l.batas); err != nil {
		if errors.Is(err, unggah.ErrBerkasKosong) {
			return HasilLampiranTCO{}, fmt.Errorf("%w: %s", unggah.ErrBerkasKosong, PesanTanpaBerkasTCO)
		}
		return HasilLampiranTCO{}, err
	}
	// tco4: ukuran tidak disimpan - `M_ATTACHMENTTREATY_2` tidak punya kolomnya.
	rekam := models.LampiranTCO{IDTreatyYear: tahunID, FileName: nama,
		FileMimeType: unggah.MimeDokumen(berkas.Mime, nama), Category: kategori, ImageID: imageID,
		UserID: pelaku.AkunID, TglUpload: saat}

	err = l.transaksi(ctx, func(tx *db.Tx) error {
		id, err := l.gudang.Sisip(ctx, tx, rekam)
		if err != nil {
			return err
		}
		rekam.ID = id
		muatan, err := muatanLampiranTCO{TahunID: tahunID, LampiranID: id, ImageID: imageID,
			BerkasAntre: filepath.Base(jalur)}.teks()
		if err != nil {
			return err
		}
		// ⛔ Outbox DI DALAM transaksi rekamnya: efek yang diantre terpisah
		// dapat hilang sendirian, dan berkas antrean tidak akan pernah terkirim.
		if err := l.antrean.Antre(ctx, tx, unggah.JenisEfekStorageUnggah, id, muatan, saat); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		// Transaksi batal: tidak ada rekam yang merujuk berkas antrean ini.
		if e := os.Remove(jalur); e != nil && !os.IsNotExist(e) {
			return HasilLampiranTCO{}, fmt.Errorf("%w (dan berkas antrean gagal dibuang: %v)", err, e)
		}
		return HasilLampiranTCO{}, err
	}
	return l.sesudahAksi(ctx, pelaku.AkunID, tahunID, rekam), nil
}

// sesudahAksi menjalankan antrean lalu membaca ulang lampirannya.
func (l *LampiranTahunTCO) sesudahAksi(ctx context.Context, akunID, tahunID string,
	rekam models.LampiranTCO) HasilLampiranTCO {

	peringatan := l.kirimSekarang(ctx, akunID, rekam.ID)
	b, err := l.gudang.Ambil(ctx, tahunID, rekam.ID)
	if err != nil {
		if peringatan == "" {
			peringatan = "status lampiran belum dapat dibaca ulang; muat ulang daftar lampiran"
		}
		return HasilLampiranTCO{Lampiran: TampilLampiran(repository.BarisLampiranTCO{LampiranTCO: rekam}),
			Peringatan: peringatan}
	}
	return HasilLampiranTCO{Lampiran: TampilLampiran(b), Peringatan: peringatan}
}

// kirimSekarang menjalankan antrean modul ini sesudah sebuah aksi.
//
// ⚠️ Efek aksi ini dijalankan SEKETIKA supaya status tampil di jawaban yang
// sama; sisanya dipungut pekerja latar `JalankanPekerja` (OQ-TCO-09) bila
// `cmd/api` menyalakannya. Keduanya boleh berjalan bersamaan: `FOR UPDATE SKIP
// LOCKED` pemungutan membuat satu baris outbox dijalankan satu pihak saja, dan
// unggah ulang menulis ke objek yang sama (`IMAGEID`). Kegagalan EFEK tercatat
// di outbox dan tampil sebagai status; yang dilaporkan di sini hanya kegagalan
// menjalankan antreannya sendiri.
//
// ⛔ Temuan /code-review: yang dijalankan HANYA efek milik lampiran aksi ini
// (`rujukan`) - bukan antrean pemakai lain, yang dulu ikut terkirim di dalam
// permintaan ini dan jejak menyerahnya tercatat atas nama pemakai yang salah.
func (l *LampiranTahunTCO) kirimSekarang(ctx context.Context, akunID, rujukan string) string {
	for n := 0; n < putaranSeusaiAksiTCO; n++ {
		err := l.satuPutaran(ctx, akunID, l.jam(), func(tx *db.Tx, saat time.Time) (outbox.BarisEfekKeluar, error) {
			return l.antrean.PungutRujukan(ctx, tx, rujukan, saat)
		})
		if errors.Is(err, outbox.ErrEfekTidakAda) {
			return ""
		}
		if err != nil {
			return "antrean pengiriman lampiran belum dapat dijalankan; lampiran tetap tertunda dan dapat diulang"
		}
	}
	return ""
}

// JalankanAntrean memungut sampai `maks` efek modul ini yang jatuh tempo.
func (l *LampiranTahunTCO) JalankanAntrean(ctx context.Context, akunID string, maks int) (int, error) {
	n := 0
	for n < maks {
		err := l.SatuPutaran(ctx, akunID, l.jam())
		if errors.Is(err, outbox.ErrEfekTidakAda) {
			return n, nil
		}
		if err != nil {
			return n, err
		}
		n++
	}
	return n, nil
}

// SatuPutaran memungut SATU efek, menjalankannya, dan menuntaskannya.
//
// ⛔ Satu transaksi: pungut, jalankan, tuntaskan - pola `PekerjaEfek`. Yang
// berbeda hanya jalur menyerahnya: tco4 tanpa jejak modul - status outbox
// gagal-permanen itulah catatannya.
func (l *LampiranTahunTCO) SatuPutaran(ctx context.Context, akunID string, saat time.Time) error {
	return l.satuPutaran(ctx, akunID, saat, func(tx *db.Tx, saat time.Time) (outbox.BarisEfekKeluar, error) {
		return l.antrean.Pungut(ctx, tx, saat)
	})
}

// satuPutaran - pungut (lewat `pungut`), jalankan, tuntaskan: satu transaksi.
func (l *LampiranTahunTCO) satuPutaran(ctx context.Context, akunID string, saat time.Time,
	pungut func(tx *db.Tx, saat time.Time) (outbox.BarisEfekKeluar, error)) error {
	return l.transaksi(ctx, func(tx *db.Tx) error {
		b, err := pungut(tx, saat)
		if err != nil {
			return err
		}
		jalanErr := l.laksanakan(ctx, tx, b)
		if jalanErr == nil {
			return l.antrean.Tuntaskan(ctx, tx, b.ID, outbox.StatusEfekSelesai, time.Time{}, "", saat)
		}
		menyerah := !layakUlangLampiranTCO(jalanErr) || b.Percobaan >= outbox.PercobaanMaksimum
		status, jadwal := outbox.StatusEfekAntre, saat.Add(outbox.Backoff(b.Percobaan+1))
		if menyerah {
			status, jadwal = outbox.StatusEfekGagalPermanen, time.Time{}
		}
		if err := l.antrean.Tuntaskan(ctx, tx, b.ID, status, jadwal, ringkasGalatTCO(jalanErr), saat); err != nil {
			return err
		}
		// tco4: nol jejak modul (Pega tidak mencatatnya) - menyerah terlihat
		// dari status gagal-permanen dan galatnya di outbox.
		return nil
	})
}

// layakUlangLampiranTCO - `LayakDicobaUlang` + keadaan permanen milik lampiran.
func layakUlangLampiranTCO(err error) bool {
	for _, permanen := range []error{ErrBerkasSumberLampiranHilang, errMuatanLampiranRusak,
		errJenisEfekAsingTCO, unggah.ErrUnggahanDirBelumDisetel, ErrPenyimpananLampiranBelumDisuntik,
		// OQ-TCO-08 pelaksana nyata: keadaan yang tidak berubah karena dicoba lagi.
		ErrStorageMenolakPermintaanTCO, layanan.ErrGaramTokenKosong, layanan.ErrAppNameKosong,
		repository.ErrAppStorageKosongTCO} {
		if errors.Is(err, permanen) {
			return false
		}
	}
	return outbox.LayakDicobaUlang(err)
}

func ringkasGalatTCO(err error) string {
	r := []rune(err.Error())
	if len(r) > batasGalatTCO {
		return string(r[:batasGalatTCO])
	}
	return string(r)
}

// laksanakan menjalankan satu baris outbox lampiran.
func (l *LampiranTahunTCO) laksanakan(ctx context.Context, tx *db.Tx, b outbox.BarisEfekKeluar) error {
	switch b.Jenis {
	case unggah.JenisEfekStorageUnggah:
		return l.kirimUnggah(ctx, tx, b)
	case unggah.JenisEfekStorageHapus:
		return l.kirimHapus(ctx, tx, b)
	default:
		return fmt.Errorf("%w: %q", errJenisEfekAsingTCO, b.Jenis)
	}
}

// kirimUnggah menyalin berkas antrean ke penyimpanan - IDEMPOTEN.
//
// ⛔ Tiga jalan selesai-tanpa-kerja, dan ketiganya sengaja:
//   - rekamnya sudah dihapus: efek hapus yang membersihkan;
//   - objeknya sudah tercatat di `T_STORAGE_IMAGE`: percobaan sebelumnya berhasil;
//   - berkas antrean hilang TETAPI penyimpanan memiliki kuncinya: percobaan
//     sebelumnya menulis berkasnya lalu gagal dicatat.
//
// Kuncinya `IMAGEID` rekam itu sendiri, jadi pengulangan menulis ke tempat
// yang SAMA - tidak menggandakan berkas (AC 58).
func (l *LampiranTahunTCO) kirimUnggah(ctx context.Context, tx *db.Tx, b outbox.BarisEfekKeluar) error {
	rekam, err := l.gudang.AmbilUntukKirim(ctx, tx, b.Rujukan)
	if errors.Is(err, repository.ErrLampiranTidakAda) {
		return nil
	}
	if err != nil {
		return err
	}
	if strings.TrimSpace(rekam.TStorageID) != "" {
		return nil
	}
	jalur := l.jalurAntre(rekam.ImageID, rekam.FileName)
	f, err := os.Open(jalur)
	if os.IsNotExist(err) {
		ada, errAda := l.penyimpanan.Ada(ctx, rekam.ImageID)
		if errAda != nil {
			return errAda
		}
		if !ada {
			return fmt.Errorf("%w: lampiran %s", ErrBerkasSumberLampiranHilang, rekam.ID)
		}
		// Jawaban unggah yang dulu hilang tidak dapat dipulihkan: yang dicatat
		// hanya kunci objeknya.
		return l.gudang.SimpanObjek(ctx, tx, models.ObjekPenyimpananTCO{ImageID: rekam.ImageID, Namafile: rekam.ImageID})
	}
	if err != nil {
		return fmt.Errorf("services: membuka berkas antrean lampiran %s: %w", rekam.ID, err)
	}
	objek, simpanErr := l.penyimpanan.Simpan(ctx, rekam.ImageID, f, rekam.FileMimeType,
		strings.ToLower(strings.TrimPrefix(filepath.Ext(rekam.FileName), ".")))
	_ = f.Close()
	if simpanErr != nil {
		return fmt.Errorf("services: mengunggah lampiran %s: %w", rekam.ID, simpanErr)
	}
	objek.ImageID = rekam.ImageID
	if objek.Namafile == "" {
		objek.Namafile = rekam.ImageID
	}
	if err := l.gudang.SimpanObjek(ctx, tx, objek); err != nil {
		return err
	}
	// ⚠️ Berkas antrean dibuang sesudah penandaan. Bila transaksinya kemudian
	// batal, percobaan berikutnya menempuh jalan ketiga di atas.
	_ = os.Remove(jalur)
	return nil
}

// kirimHapus membuang berkas di penyimpanan, catatan objeknya
// (`DeleteStorage_SQL`, sesudah REST delete seperti `DeleteGoogleStorage_Act`
// b1582-b1635), dan berkas antreannya.
//
// ⛔ Berkas yang SUDAH tidak ada bukan kegagalan: penghapusan rekam tidak
// boleh gagal karena penyimpanan lebih dulu kehilangan berkasnya.
func (l *LampiranTahunTCO) kirimHapus(ctx context.Context, tx *db.Tx, b outbox.BarisEfekKeluar) error {
	var m muatanLampiranTCO
	if err := json.Unmarshal([]byte(b.Muatan), &m); err != nil || strings.TrimSpace(m.ImageID) == "" {
		return fmt.Errorf("%w: efek %s", errMuatanLampiranRusak, b.ID)
	}
	// ⛔ Urutan: `Hapus` penyimpanan (geturl + `Update_T_Storage_SQL` dalam
	// transaksinya SENDIRI, OQ-TCO-26) SEBELUM `HapusObjek` di `tx`. Dibalik,
	// transaksi pendek penyegaran menunggu kunci baris milik `tx` ini sendiri
	// - Oracle tidak mendeteksinya, efek menggantung sampai konteks habis.
	if err := l.penyimpanan.Hapus(ctx, m.ImageID); err != nil && !errors.Is(err, ErrBerkasTidakAdaDiPenyimpanan) {
		return fmt.Errorf("services: menghapus berkas lampiran %s: %w", m.LampiranID, err)
	}
	if err := l.gudang.HapusObjek(ctx, tx, m.ImageID); err != nil {
		return err
	}
	if m.BerkasAntre != "" && m.BerkasAntre == filepath.Base(m.BerkasAntre) {
		if err := os.Remove(filepath.Join(l.folderAntre(), m.BerkasAntre)); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("services: membuang berkas antrean lampiran %s: %w", m.LampiranID, err)
		}
	}
	return nil
}

// Unduh menyerahkan isi satu lampiran - tautan nama berkas b3470
// (`TreatyOutDownloadOne`).
//
// ⛔ Mengembalikan PEMBACA, bukan isi: berkas 25 MiB tidak dibaca ke memori.
// Pemanggil wajib menutupnya.
func (l *LampiranTahunTCO) Unduh(ctx context.Context, pelaku inti.Pelaku, tahunID, id string) (
	LampiranTampil, io.ReadCloser, error) {

	if err := inti.WajibIdentitas(pelaku); err != nil {
		return LampiranTampil{}, nil, err
	}
	b, err := l.gudang.Ambil(ctx, tahunID, id)
	if err != nil {
		return LampiranTampil{}, nil, err
	}
	tampil := TampilLampiran(b)
	if tampil.Status != models.StatusLampiranTerkirim {
		return LampiranTampil{}, nil, fmt.Errorf("%w: lampiran %s berstatus %s", ErrLampiranBelumTerkirim,
			id, tampil.Status)
	}
	rc, err := l.penyimpanan.Buka(ctx, b.ImageID)
	if errors.Is(err, ErrBerkasTidakAdaDiPenyimpanan) {
		return LampiranTampil{}, nil, fmt.Errorf("%w: lampiran %s", ErrLampiranTanpaBerkas, id)
	}
	if err != nil {
		return LampiranTampil{}, nil, err
	}
	return tampil, rc, nil
}

// UnduhSemua menyerahkan seluruh lampiran terkirim - `Download All` b2659.
//
// ⛔ Keselarasan diperiksa LEBIH DULU: arsip yang dialirkan tidak dapat
// ditarik kembali, jadi rekam tanpa berkas menggagalkan permintaannya sebelum
// satu byte pun terkirim, dengan pesan yang menyebut lampiran mana.
func (l *LampiranTahunTCO) UnduhSemua(ctx context.Context, pelaku inti.Pelaku, tahunID string,
	tulis func(namaEntri string, isi io.Reader) error) error {

	daftar, err := l.Daftar(ctx, pelaku, tahunID)
	if err != nil {
		return err
	}
	var terkirim []LampiranTampil
	kunci := map[string]string{}
	for _, t := range daftar {
		if t.Status == models.StatusLampiranTerkirim {
			terkirim = append(terkirim, t)
		}
	}
	if len(terkirim) == 0 {
		return fmt.Errorf("%w: tidak ada lampiran terkirim pada tahun treaty %s", ErrLampiranBelumTerkirim, tahunID)
	}
	var hilang []string
	for _, t := range terkirim {
		b, err := l.gudang.Ambil(ctx, tahunID, t.ID)
		if err != nil {
			return err
		}
		kunci[t.ID] = b.ImageID
		ada, err := l.penyimpanan.Ada(ctx, b.ImageID)
		if err != nil {
			return err
		}
		if !ada {
			hilang = append(hilang, t.ID)
		}
	}
	if len(hilang) > 0 {
		return fmt.Errorf("%w: lampiran %s", ErrLampiranTanpaBerkas, strings.Join(hilang, ", "))
	}
	for _, t := range terkirim {
		rc, err := l.penyimpanan.Buka(ctx, kunci[t.ID])
		if err != nil {
			return err
		}
		err = tulis(models.NamaEntriZipLampiranTCO(t.ID, t.FileName), rc)
		_ = rc.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

// Hapus membuang satu lampiran - `Delete` b3897 (`DeleteAttachmentTreaty`).
//
// ⛔ Rekam dulu, penyimpanan menyusul lewat outbox - SELALU, termasuk bila
// belum terkirim: percobaan unggah yang menulis berkasnya lalu gagal dicatat
// meninggalkan berkas yang rekamnya tidak tahu.
func (l *LampiranTahunTCO) Hapus(ctx context.Context, pelaku inti.Pelaku, tahunID, id string) (string, error) {
	if err := inti.WajibIdentitas(pelaku); err != nil {
		return "", err
	}
	b, err := l.gudang.Ambil(ctx, tahunID, id)
	if err != nil {
		return "", err
	}
	saat := l.jam()
	err = l.transaksi(ctx, func(tx *db.Tx) error {
		if err := l.gudang.Hapus(ctx, tx, tahunID, id); err != nil {
			return err
		}
		muatan, err := muatanLampiranTCO{TahunID: tahunID, LampiranID: id, ImageID: b.ImageID,
			BerkasAntre: models.NamaBerkasAntreLampiranTCO(b.ImageID, b.FileName)}.teks()
		if err != nil {
			return err
		}
		if err := l.antrean.Antre(ctx, tx, unggah.JenisEfekStorageHapus, id, muatan, saat); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	return l.kirimSekarang(ctx, pelaku.AkunID, id), nil
}

// Ulangi mengantre ulang unggahan satu lampiran (AC 58, AC 61).
//
// ⛔ Ditolak bila tidak ada yang perlu diperbaiki (berkas ada di penyimpanan)
// dan bila tidak ada yang dapat diunggah (berkas antrean hilang dan
// penyimpanan tidak memilikinya) - yang kedua diperbaiki dengan menghapus.
func (l *LampiranTahunTCO) Ulangi(ctx context.Context, pelaku inti.Pelaku, tahunID, id string) (HasilLampiranTCO, error) {
	if err := inti.WajibIdentitas(pelaku); err != nil {
		return HasilLampiranTCO{}, err
	}
	b, err := l.gudang.Ambil(ctx, tahunID, id)
	if err != nil {
		return HasilLampiranTCO{}, err
	}
	ada, err := l.penyimpanan.Ada(ctx, b.ImageID)
	if err != nil {
		return HasilLampiranTCO{}, err
	}
	terkirim := strings.TrimSpace(b.TStorageID) != ""
	if terkirim && ada {
		return HasilLampiranTCO{}, fmt.Errorf("%w: lampiran %s", ErrLampiranSudahTerkirim, id)
	}
	if _, err := os.Stat(l.jalurAntre(b.ImageID, b.FileName)); err != nil && !ada {
		return HasilLampiranTCO{}, fmt.Errorf("%w: lampiran %s", ErrBerkasSumberLampiranHilang, id)
	}
	saat := l.jam()
	err = l.transaksi(ctx, func(tx *db.Tx) error {
		if terkirim {
			// Catatan objek tanpa berkas dibuang dulu; unggah ulang mencatatnya lagi.
			if err := l.gudang.HapusObjek(ctx, tx, b.ImageID); err != nil {
				return err
			}
		}
		muatan, err := muatanLampiranTCO{TahunID: tahunID, LampiranID: id, ImageID: b.ImageID,
			BerkasAntre: models.NamaBerkasAntreLampiranTCO(b.ImageID, b.FileName)}.teks()
		if err != nil {
			return err
		}
		if err := l.antrean.Antre(ctx, tx, unggah.JenisEfekStorageUnggah, id, muatan, saat); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return HasilLampiranTCO{}, err
	}
	return l.sesudahAksi(ctx, pelaku.AkunID, tahunID, b.LampiranTCO), nil
}

// PeriksaSelaras membandingkan rekam dengan penyimpanan (AC 61).
//
// Dua keadaan tidak sejalan yang dicari:
//   - terkirim, tetapi penyimpanan tidak memiliki berkasnya;
//   - belum terkirim, berkas antreannya hilang, dan penyimpanan pun tidak
//     memilikinya - unggahannya tidak akan pernah selesai.
func (l *LampiranTahunTCO) PeriksaSelaras(ctx context.Context, pelaku inti.Pelaku, tahunID string) (
	[]TemuanSelarasTCO, error) {

	if err := inti.WajibIdentitas(pelaku); err != nil {
		return nil, err
	}
	baris, err := l.gudang.Daftar(ctx, tahunID)
	if err != nil {
		return nil, err
	}
	if err := l.periksaTahun(ctx, tahunID); err != nil {
		return nil, err
	}
	temuan := []TemuanSelarasTCO{}
	for _, b := range baris {
		ada, err := l.penyimpanan.Ada(ctx, b.ImageID)
		if err != nil {
			return nil, err
		}
		if ada {
			continue
		}
		_, errAntre := os.Stat(l.jalurAntre(b.ImageID, b.FileName))
		adaAntre := errAntre == nil
		perbaikan := "hapus"
		if adaAntre {
			perbaikan = "ulangi"
		}
		switch {
		case strings.TrimSpace(b.TStorageID) != "":
			temuan = append(temuan, TemuanSelarasTCO{LampiranID: b.ID, FileName: b.FileName,
				Masalah: "rekam tanpa berkas di penyimpanan", Perbaikan: perbaikan})
		case !adaAntre:
			temuan = append(temuan, TemuanSelarasTCO{LampiranID: b.ID, FileName: b.FileName,
				Masalah: "berkas sumber hilang; unggahan tidak akan pernah selesai", Perbaikan: perbaikan})
		}
	}
	return temuan, nil
}
