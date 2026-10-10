package services

// Lampiran case NB FacIn (tiket 67) - grid "Attachment File" Pega (`Section\AttachmentGridReas.xml`: Category · Count · Upload
// File · View File) dan popup View File (`Harness\ViewPictureList.xml`: File · Note · Upload Date · Delete).
//
// Keputusan 08-10-2026 (issues/67): tabel Pega DOCUMENT_POLIS + T_STORAGE_IMAGE (K67-1); unggah oleh siapa pun yang membuka
// case selama case BELUM jadi polis (K67-6, 09-10-2026, menggantikan "pemegang workbasket POSITION" K67-2), hapus bila
// FlagOnGoingPolicy "0" ATAU akun IT, siapa pun yang membuka case (K67-2 / K-F "Persis Pega"); hapus FISIK seperti Pega (K-A, pengecualian sadar ADR-0031); mesin penyimpan salinan Treaty Contract Out (K-B);
// unggah ASINKRON lewat outbox bersama (K-C, ADR-0008): DOCUMENT_POLIS + efek `storage-unggah` satu transaksi, T_STORAGE_IMAGE
// oleh pelaksana sesudah unggah sukses; ID cap waktu Pega (K-D); jenis berkas di luar GetMimeType 48 baris = 400 (K-E).
//
// Penyimpangan sadar dari Pega (dicatat issues/67): Pega menyimpan baris DOCUMENT_POLIS hanya sesudah unggah sukses
// (InsertDocument_Act langkah 5) - di sini baris lahir dulu dengan status "menunggu"; Pega membatalkan diam-diam berkas tak
// dikenal / kecil (InsertGoogleStorage_Act langkah 5) - di sini 400 berpesan; URL bertanda tangan diproksi backend.

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/layanan"
	"nusantarare/inti/backend/outbox"
	"nusantarare/inti/backend/unggah"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
)

const (
	// ModulOutboxNBFacIn / LiniOutboxNBFacIn - kolom MODUL / LINI outbox bersama T_LOG_SERVICE_RNM `[keputusan agent]`
	// (LINI = lini case NB FacIn di T_WORK_POLIS).
	ModulOutboxNBFacIn = "NBFACIN"
	LiniOutboxNBFacIn  = repository.LiniFacIn
	// Kategori1Lampiran - AttachFileGIS langkah 8 `KATEGORI_1="Reas"`.
	Kategori1Lampiran = "Reas"
	// IntervalPekerjaLampiran - jeda pekerja latar antrean lampiran `[keputusan agent]`: tanpa env baru (inti tidak
	// disentuh, K-B); nol pekerja tanpa Oracle.
	IntervalPekerjaLampiran = time.Minute
	// AkunPekerjaLampiran - pelaku jejak pekerja latar.
	AkunPekerjaLampiran = "SISTEM-PEKERJA-NBFACIN"
	// ukuranMinimumLampiran - InsertGoogleStorage_Act langkah 5 `@length(Param.Image) < 10` (base64 < 10 karakter = berkas
	// < 7 byte) -> keluar.
	ukuranMinimumLampiran = 7
	putaranSeusaiAksi     = 10
	batasPutaranPekerja   = 50
	batasGalatTersimpan   = 1000
	folderAntreLampiran   = "antre"
)

var (
	// ErrLampiranTanpaDatabase - layanan tanpa basis data (503).
	ErrLampiranTanpaDatabase = errors.New("services: basis data tidak dikonfigurasi, lampiran tidak dapat diproses")
	// ErrMasukanLampiran - masukan unggah / daftar tidak sah (400).
	ErrMasukanLampiran = errors.New("services: masukan lampiran tidak sah")
	// ErrKategoriLampiranTidakDikenal - kategori di luar master CATEGORY_ATTACH_REAS (422).
	ErrKategoriLampiranTidakDikenal = errors.New("services: kategori lampiran tidak ada di master CATEGORY_ATTACH_REAS")
	// ErrLampiranTidakAda - dokumen tidak ada / bukan milik case itu (404).
	ErrLampiranTidakAda = errors.New("services: dokumen lampiran tidak ada di case ini")
	// ErrLampiranBelumTersimpan - berkas belum dipastikan ada di penyimpanan (409).
	ErrLampiranBelumTersimpan = errors.New("services: berkas lampiran belum tersimpan di penyimpanan")
	// ErrLampiranSudahTersimpan - "ulangi" atas berkas yang sudah tersimpan (409).
	ErrLampiranSudahTersimpan = errors.New("services: berkas lampiran sudah tersimpan; tidak ada yang perlu diulang")
	// ErrLampiranTanpaBerkas - tercatat tetapi berkasnya tidak ada di penyimpanan (409).
	ErrLampiranTanpaBerkas = errors.New("services: berkas lampiran tidak ditemukan di penyimpanan")
	// ErrBerkasSumberLampiranHilang - berkas antrean hilang dan penyimpanan tidak memilikinya (409; PERMANEN di pelaksana).
	ErrBerkasSumberLampiranHilang = errors.New("services: berkas sumber lampiran sudah tidak ada; hapus lalu unggah ulang")

	errMuatanLampiranRusak = errors.New("services: muatan outbox lampiran tak terbaca")
	errJenisEfekAsing      = errors.New("services: jenis efek bukan milik pelaksana lampiran")
)

// PesanJenisBerkasTidakDidukung - K-E `[keputusan agent]` "Jenis berkas tidak didukung: .<ext>".
func PesanJenisBerkasTidakDidukung(ekstensi string) string {
	if ekstensi == "" {
		return "Jenis berkas tidak didukung: (tanpa ekstensi)"
	}
	return "Jenis berkas tidak didukung: ." + ekstensi
}

const (
	pesanTanpaBerkas     = "Pilih berkas yang akan diunggah."
	pesanBerkasKecil     = "Berkas kosong atau terlalu kecil."
	pesanKategoriKosong  = "Kategori wajib dipilih."
	pesanAntreanTertunda = "antrean pengiriman lampiran tidak dapat dijalankan; berkas tetap menunggu dan dapat diulang"
)

// AntreanLampiran - outbox bersama yang dipakai lampiran, sudah bermodul NBFACIN.
type AntreanLampiran interface {
	Antre(ctx context.Context, tx *db.Tx, jenis, rujukan, muatan string, saat time.Time) error
	Pungut(ctx context.Context, tx *db.Tx, saat time.Time) (outbox.BarisEfekKeluar, error)
	PungutRujukan(ctx context.Context, tx *db.Tx, rujukan string, saat time.Time) (outbox.BarisEfekKeluar, error)
	Tuntaskan(ctx context.Context, tx *db.Tx, id, status string, jadwal time.Time, galat string, saat time.Time) error
}

type antreanLampiranOracle struct {
	outbox *outbox.Penyimpan
	repo   repository.PenyimpanLampiran
}

func (a antreanLampiranOracle) Antre(ctx context.Context, tx *db.Tx, jenis, rujukan, muatan string, saat time.Time) error {
	_, err := a.outbox.AntreEfek(ctx, tx, LiniOutboxNBFacIn, ModulOutboxNBFacIn, jenis, rujukan, muatan, saat)
	return err
}
func (a antreanLampiranOracle) Pungut(ctx context.Context, tx *db.Tx, saat time.Time) (outbox.BarisEfekKeluar, error) {
	return a.outbox.PungutEfek(ctx, tx, ModulOutboxNBFacIn, saat)
}
func (a antreanLampiranOracle) PungutRujukan(ctx context.Context, tx *db.Tx, rujukan string, saat time.Time) (outbox.BarisEfekKeluar, error) {
	return a.repo.PungutEfekRujukan(ctx, tx, ModulOutboxNBFacIn, rujukan, saat)
}
func (a antreanLampiranOracle) Tuntaskan(ctx context.Context, tx *db.Tx, id, status string, jadwal time.Time, galat string,
	saat time.Time) error {
	return a.outbox.TuntaskanEfek(ctx, tx, id, status, jadwal, galat, saat)
}

// pencatatObjekOracle - Update_T_Storage_SQL dalam transaksi pendeknya sendiri.
type pencatatObjekOracle struct {
	repo      repository.PenyimpanLampiran
	transaksi Transaksi
}

func (p pencatatObjekOracle) PerbaruiObjek(ctx context.Context, o models.ObjekStorage) error {
	return p.transaksi(ctx, func(tx *db.Tx) error { return p.repo.PerbaruiObjek(ctx, tx, o) })
}

// penyimpanTokenOracle - PenyimpanTokenStorage atas repository + SimpanToken bersama inti.
type penyimpanTokenOracle struct {
	repository.PenyimpanLampiran
	simpan *layanan.PenyimpanToken
}

func (p penyimpanTokenOracle) SimpanToken(ctx context.Context, tx *db.Tx, app, token, pengguna string, sampai time.Time) error {
	return p.simpan.SimpanToken(ctx, tx, app, token, pengguna, sampai)
}

// DenganLampiran memasang lampiran atas Oracle: repository, outbox bersama, dan pelaksana penyimpanan menurut
// PELAKSANA_STORAGE (`nyata` = layanan; selain itu stub lokal di bawah UNGGAHAN_DIR). ⛔ Garam tidak pernah dicetak.
func (s *Service) DenganLampiran(d *inti.Dasar, nyata bool, garam string) *Service {
	if d == nil || !d.PunyaDatabase() {
		return s
	}
	repo := repository.NewLampiranOracle(d.DB())
	s.lampiran, s.unggahanDir = repo, d.UnggahanDir()
	s.antreanLampiran = antreanLampiranOracle{outbox: outbox.NewPenyimpan(d.DB()), repo: repo}
	s.storageNyata = nyata
	if !nyata {
		s.storageLampiran = StorageLampiranLokal(s.unggahanDir)
		return s
	}
	cache := NewCacheTokenStorage(sumberTokenStorage{transaksi: d.DalamTransaksi,
		penyimp: penyimpanTokenOracle{PenyimpanLampiran: repo, simpan: layanan.NewPenyimpanToken(d.DB())}, garam: garam},
		time.Now, MarginTokenStorage)
	s.storageLampiran = NewStorageLampiranJarakJauh(layanan.ResolverLinkServiceOracle(d), cache, repo.AppStorage, nil,
		pencatatObjekOracle{repo: repo, transaksi: d.DalamTransaksi}, nil)
	return s
}

// DenganLampiranUji memasang bagian lampiran langsung - dipakai uji.
func (s *Service) DenganLampiranUji(repo repository.PenyimpanLampiran, antrean AntreanLampiran, storage KlienStorageLampiran,
	unggahanDir string) *Service {
	s.lampiran, s.antreanLampiran, s.storageLampiran, s.unggahanDir = repo, antrean, storage, unggahanDir
	return s
}

// PenyimpananLampiranNyata - pelaksana nyata terpasang (dicatat modul.go).
func (s *Service) PenyimpananLampiranNyata() bool { return s.storageNyata }

// PunyaLampiran - lampiran terpasang (pekerja latar hanya menyala bila ya).
func (s *Service) PunyaLampiran() bool { return s.siapLampiran() == nil }

// muatanLampiran - muatan JSON outbox: PENGENAL SAJA (nol nama orang, nol alamat).
type muatanLampiran struct {
	DokumenID   string `json:"dokumenId"`
	ImageID     string `json:"imageId"`
	BerkasAntre string `json:"berkasAntre"`
}

func (m muatanLampiran) teks() (string, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return "", fmt.Errorf("services: menyusun muatan lampiran: %w", err)
	}
	return string(b), nil
}

// TampilanLampiran - GET …/lampiran.
type TampilanLampiran struct {
	Kategori    []models.KategoriLampiran
	BolehUnggah bool
	BolehHapus  bool
}

// BarisLampiran - satu berkas popup View File.
type BarisLampiran struct {
	ID       string
	NamaFile string
	// Catatan - kolom "Note" ViewPictureList = `.KATEGORI_2` (b3993) `[terverifikasi]`.
	Catatan string
	Tanggal string
	Status  string
}

func (s *Service) siapLampiran() error {
	if s.lampiran == nil || s.antreanLampiran == nil || s.storageLampiran == nil || s.transaksi == nil || s.tahap == nil {
		return ErrLampiranTanpaDatabase
	}
	return nil
}

// tahapLampiran - identitas (401), layanan siap (503), case ada (404); jawab POSITION / FLAG_ONGOING_POLICY.
func (s *Service) tahapLampiran(ctx context.Context, pelaku inti.Pelaku, id string) (models.TahapKasus, error) {
	if err := inti.WajibIdentitas(pelaku); err != nil {
		return models.TahapKasus{}, err
	}
	if err := s.siapLampiran(); err != nil {
		return models.TahapKasus{}, err
	}
	if !idKasusSah(id) {
		return models.TahapKasus{}, ErrKasusTidakAda
	}
	t, err := s.tahap.TahapKasus(ctx, id)
	if err != nil {
		return models.TahapKasus{}, galatKasus(err)
	}
	return t, nil
}

// statusKunciUnggah - K67-6 + ralat (keputusan work owner 09-10-2026 "Resolved-Decline, Resolved-Reject JUGA DIKUNCI JUGA";
// ejaan 10-10-2026: "Resolved-Rejected" mengikuti Pega, `NB FacIn\ReportDefinition\GetListOpportunity*.xml`): daftar
// EKSPLISIT, bukan awalan "Resolved-". Sistem baru belum menulis status reject.
var statusKunciUnggah = map[string]bool{statusSelesaiPolicy: true, statusDecline: true, statusRejected: true}

// statusRejected - STATUS_WORK reject ejaan Pega (K67-6, ralat ejaan 10-10-2026).
const statusRejected = "Resolved-Rejected"

// bolehUnggahLampiran - K67-6 (keputusan work owner 09-10-2026, menggantikan bagian unggah K67-2): siapa pun yang membuka case,
// kecuali STATUS_WORK salah satu statusKunciUnggah (polis jadi "Resolved-Completed" K69-7, "Resolved-Decline", "Resolved-Rejected").
// Pega: tombol Upload `Section\AttachmentGridReas.xml` sel L2330 pyVisible ALWAYS / pyReadOnly false, tanpa syarat
// `[terverifikasi]` sesi 0f; View Policy / View Following ditangani frontend (hanyaLihat).
func bolehUnggahLampiran(t models.TahapKasus) bool {
	return !statusKunciUnggah[strings.TrimSpace(t.StatusWork)]
}

// pesanUnggahPolis - 403 unggah / ulangi pada case terkunci (teks agent).
const pesanUnggahPolis = "unggah lampiran tidak boleh: case sudah selesai (polis / decline / reject)"

// bolehHapusLampiran - AttachmentGridReas / ReasViewAttachment b4605 `FlagOnGoingPolicy = '0' || OperatorID.pyOrgDivision =
// 'IT'` untuk siapa pun yang membuka case (K-F "Persis Pega"); IT = penanda yang sama dengan aksesLossRecord (JOB_POSITION).
func (s *Service) bolehHapusLampiran(ctx context.Context, t models.TahapKasus, pelaku inti.Pelaku) bool {
	return t.FlagOnGoingPolicy == flagOffer || s.AksesTabIT(ctx, pelaku)
}

// TampilanLampiranKasus - grid Attachment File: kategori master (urut NOTE) + Count dokumen case + hak tombol.
func (s *Service) TampilanLampiranKasus(ctx context.Context, pelaku inti.Pelaku, id string) (TampilanLampiran, error) {
	t, err := s.tahapLampiran(ctx, pelaku, id)
	if err != nil {
		return TampilanLampiran{}, err
	}
	master, err := s.lampiran.KategoriLampiran(ctx)
	if err != nil {
		return TampilanLampiran{}, err
	}
	jumlah, err := s.lampiran.JumlahLampiran(ctx, idPega(id))
	if err != nil {
		return TampilanLampiran{}, err
	}
	kat := make([]models.KategoriLampiran, 0, len(master))
	for _, m := range master {
		kat = append(kat, models.KategoriLampiran{Kategori: m, Jumlah: jumlah[m]})
	}
	return TampilanLampiran{Kategori: kat, BolehUnggah: bolehUnggahLampiran(t),
		BolehHapus: s.bolehHapusLampiran(ctx, t, pelaku)}, nil
}

func barisLampiran(d models.DokumenPolis) BarisLampiran {
	return BarisLampiran{ID: d.ID, NamaFile: d.NamaFile, Catatan: d.Kategori2, Tanggal: d.Tanggal,
		Status: models.StatusDokumen(d, outbox.StatusEfekGagalPermanen)}
}

// DaftarBerkasLampiran - popup View File satu kategori (GetAllDocument).
func (s *Service) DaftarBerkasLampiran(ctx context.Context, pelaku inti.Pelaku, id, kategori string) ([]BarisLampiran, error) {
	if _, err := s.tahapLampiran(ctx, pelaku, id); err != nil {
		return nil, err
	}
	if strings.TrimSpace(kategori) == "" {
		return nil, fmt.Errorf("%w: %s", ErrMasukanLampiran, pesanKategoriKosong)
	}
	dok, err := s.lampiran.DaftarLampiran(ctx, idPega(id), kategori, ModulOutboxNBFacIn, unggah.JenisEfekStorageUnggah)
	if err != nil {
		return nil, err
	}
	out := make([]BarisLampiran, 0, len(dok))
	for _, d := range dok {
		out = append(out, barisLampiran(d))
	}
	return out, nil
}

func (s *Service) folderAntre() string {
	return filepath.Join(s.unggahanDir, folderStubLampiran, folderAntreLampiran)
}

func (s *Service) jalurAntre(imageID, nama string) string {
	return filepath.Join(s.folderAntre(), models.NamaBerkasAntreLampiran(imageID, nama))
}

// UnggahLampiran - satu berkas (AttachContentGIS -> AttachFileGIS -> InsertDocument_Act). Urutan: gerbang (identitas, case,
// workbasket, folder, nama, jenis, kategori) -> berkas ke folder antrean -> SATU transaksi DOCUMENT_POLIS + efek outbox ->
// sesudah commit efek ini dijalankan seketika (gagalnya hanya status). Jawab ID dokumen.
func (s *Service) UnggahLampiran(ctx context.Context, pelaku inti.Pelaku, id string, berkas unggah.BerkasMasuk) (string, error) {
	t, err := s.tahapLampiran(ctx, pelaku, id)
	if err != nil {
		return "", err
	}
	if !bolehUnggahLampiran(t) {
		return "", fmt.Errorf("%w: %s", inti.ErrTanpaWewenang, pesanUnggahPolis)
	}
	if strings.TrimSpace(s.unggahanDir) == "" {
		return "", unggah.ErrUnggahanDirBelumDisetel
	}
	nama := models.NamaDasarBerkas(berkas.NamaFile)
	if nama == "" || berkas.Isi == nil {
		return "", fmt.Errorf("%w: %s", ErrMasukanLampiran, pesanTanpaBerkas)
	}
	ext := models.EkstensiBerkas(nama)
	if ext == "" || unggah.MimeDariNamaFile(nama) == unggah.MimeBawaan {
		return "", fmt.Errorf("%w: %s", ErrMasukanLampiran, PesanJenisBerkasTidakDidukung(ext))
	}
	if strings.TrimSpace(berkas.Kategori) == "" {
		return "", fmt.Errorf("%w: %s", ErrMasukanLampiran, pesanKategoriKosong)
	}
	master, err := s.lampiran.KategoriLampiran(ctx)
	if err != nil {
		return "", err
	}
	kategori, sah := models.KategoriLampiranSah(master, berkas.Kategori)
	if !sah {
		return "", fmt.Errorf("%w: %q", ErrKategoriLampiranTidakDikenal, berkas.Kategori)
	}
	saat := s.sekarang()
	imageID, err := unggah.ImageIDBaru(saat)
	if err != nil {
		return "", err
	}
	jalur := s.jalurAntre(imageID, nama)
	if err := unggah.TulisBerkas(jalur, berkas.Isi, 0); err != nil {
		if errors.Is(err, unggah.ErrBerkasKosong) {
			return "", fmt.Errorf("%w: %s", ErrMasukanLampiran, pesanBerkasKecil)
		}
		return "", err
	}
	if info, err := os.Stat(jalur); err != nil || info.Size() < ukuranMinimumLampiran {
		_ = os.Remove(jalur)
		return "", fmt.Errorf("%w: %s", ErrMasukanLampiran, pesanBerkasKecil)
	}
	dok := models.DokumenPolis{IDPega: idPega(id), NamaFile: nama, Mime: ext, Kategori1: Kategori1Lampiran,
		Kategori2: kategori, TStorageID: imageID, Operator: pelaku.AkunID}
	var dokID string
	err = s.transaksi(ctx, func(tx *db.Tx) error {
		i, err := s.lampiran.SisipLampiran(ctx, tx, dok, saat)
		if err != nil {
			return err
		}
		dokID = i
		muatan, err := muatanLampiran{DokumenID: i, ImageID: imageID, BerkasAntre: filepath.Base(jalur)}.teks()
		if err != nil {
			return err
		}
		// ⛔ Outbox DI DALAM transaksi dokumennya: efek yang diantre terpisah dapat hilang sendirian.
		return s.antreanLampiran.Antre(ctx, tx, unggah.JenisEfekStorageUnggah, i, muatan, saat)
	})
	if err != nil {
		if e := os.Remove(jalur); e != nil && !os.IsNotExist(e) {
			return "", fmt.Errorf("%w (dan membuang berkas antrean gagal: %v)", err, e)
		}
		return "", err
	}
	s.kirimSekarang(ctx, dokID)
	return dokID, nil
}

// ambilDokumen - satu dokumen DIBATASI case-nya.
func (s *Service) ambilDokumen(ctx context.Context, id, dokID string) (models.DokumenPolis, error) {
	d, err := s.lampiran.AmbilLampiran(ctx, idPega(id), dokID, ModulOutboxNBFacIn, unggah.JenisEfekStorageUnggah)
	if errors.Is(err, repository.ErrLampiranTidakAda) {
		return models.DokumenPolis{}, ErrLampiranTidakAda
	}
	return d, err
}

// UnduhLampiran - DownloadDocumentPolis -> GetUrlGoogleStorage_Act (proksi). ⛔ Mengembalikan PEMBACA; pemanggil menutupnya.
func (s *Service) UnduhLampiran(ctx context.Context, pelaku inti.Pelaku, id, dokID string) (models.DokumenPolis, io.ReadCloser, error) {
	if _, err := s.tahapLampiran(ctx, pelaku, id); err != nil {
		return models.DokumenPolis{}, nil, err
	}
	d, err := s.ambilDokumen(ctx, id, dokID)
	if err != nil {
		return models.DokumenPolis{}, nil, err
	}
	// Tersimpan = catatan T_STORAGE_IMAGE ada (GetLinkStorage_SQL); tanpa catatan = belum tersimpan.
	link, ada, err := s.lampiran.LinkStorage(ctx, nil, d.TStorageID)
	if err != nil {
		return models.DokumenPolis{}, nil, err
	}
	if !ada {
		return models.DokumenPolis{}, nil, fmt.Errorf("%w: dokumen %s", ErrLampiranBelumTersimpan, dokID)
	}
	rc, err := s.storageLampiran.Buka(ctx, link)
	if errors.Is(err, ErrBerkasTidakAdaDiPenyimpanan) {
		return models.DokumenPolis{}, nil, fmt.Errorf("%w: dokumen %s", ErrLampiranTanpaBerkas, dokID)
	}
	if err != nil {
		return models.DokumenPolis{}, nil, err
	}
	return d, rc, nil
}

// HapusLampiran - DeleteDocumentPolis_Act: hapus FISIK (K-A). Satu transaksi: DELETE DOCUMENT_POLIS + efek `storage-hapus`
// (DeleteGoogleStorage_Act + DeleteStorage_SQL lewat outbox) bila T_STORAGE_ID terisi (langkah 2 dilewati bila kosong).
func (s *Service) HapusLampiran(ctx context.Context, pelaku inti.Pelaku, id, dokID string) error {
	t, err := s.tahapLampiran(ctx, pelaku, id)
	if err != nil {
		return err
	}
	if !s.bolehHapusLampiran(ctx, t, pelaku) {
		return fmt.Errorf("%w: hapus lampiran hanya di tahap Offer (FlagOnGoingPolicy 0) atau untuk akun IT", inti.ErrTanpaWewenang)
	}
	d, err := s.ambilDokumen(ctx, id, dokID)
	if err != nil {
		return err
	}
	err = s.transaksi(ctx, func(tx *db.Tx) error {
		if err := s.lampiran.HapusLampiran(ctx, tx, idPega(id), dokID); err != nil {
			if errors.Is(err, repository.ErrLampiranTidakAda) {
				return ErrLampiranTidakAda
			}
			return err
		}
		if strings.TrimSpace(d.TStorageID) == "" {
			return nil
		}
		muatan, err := muatanLampiran{DokumenID: dokID, ImageID: d.TStorageID,
			BerkasAntre: models.NamaBerkasAntreLampiran(d.TStorageID, d.NamaFile)}.teks()
		if err != nil {
			return err
		}
		return s.antreanLampiran.Antre(ctx, tx, unggah.JenisEfekStorageHapus, dokID, muatan, s.sekarang())
	})
	if err != nil {
		return err
	}
	s.kirimSekarang(ctx, dokID)
	return nil
}

// UlangiLampiran - antre ulang unggahan berkas menunggu / gagal (pola Treaty Contract Out). Ditolak bila sudah tersimpan,
// dan bila tidak ada yang dapat diunggah (berkas antrean hilang DAN penyimpanan tidak memilikinya).
func (s *Service) UlangiLampiran(ctx context.Context, pelaku inti.Pelaku, id, dokID string) error {
	t, err := s.tahapLampiran(ctx, pelaku, id)
	if err != nil {
		return err
	}
	if !bolehUnggahLampiran(t) {
		return fmt.Errorf("%w: %s", inti.ErrTanpaWewenang, pesanUnggahPolis)
	}
	d, err := s.ambilDokumen(ctx, id, dokID)
	if err != nil {
		return err
	}
	if d.Tersimpan {
		return fmt.Errorf("%w: dokumen %s", ErrLampiranSudahTersimpan, dokID)
	}
	lok, ok := models.LokasiUnggah(d.ID, d.TStorageID, d.NamaFile)
	if !ok || strings.TrimSpace(d.TStorageID) == "" {
		return fmt.Errorf("%w: dokumen %s bukan unggahan sistem ini", ErrBerkasSumberLampiranHilang, dokID)
	}
	if _, err := os.Stat(s.jalurAntre(d.TStorageID, d.NamaFile)); err != nil {
		ada, _, errAda := s.storageLampiran.Ada(ctx, lok)
		if errAda != nil {
			return errAda
		}
		if !ada {
			return fmt.Errorf("%w: dokumen %s", ErrBerkasSumberLampiranHilang, dokID)
		}
	}
	err = s.transaksi(ctx, func(tx *db.Tx) error {
		muatan, err := muatanLampiran{DokumenID: dokID, ImageID: d.TStorageID,
			BerkasAntre: models.NamaBerkasAntreLampiran(d.TStorageID, d.NamaFile)}.teks()
		if err != nil {
			return err
		}
		return s.antreanLampiran.Antre(ctx, tx, unggah.JenisEfekStorageUnggah, dokID, muatan, s.sekarang())
	})
	if err != nil {
		return err
	}
	s.kirimSekarang(ctx, dokID)
	return nil
}

// kirimSekarang - jalankan efek MILIK aksi ini (rujukan) seketika; sisanya pekerja latar. Kegagalan efek tercatat di
// outbox (status); yang dicatat di sini hanya kegagalan menjalankan antreannya.
func (s *Service) kirimSekarang(ctx context.Context, rujukan string) {
	for n := 0; n < putaranSeusaiAksi; n++ {
		err := s.satuPutaranLampiran(ctx, s.sekarang(), func(tx *db.Tx, saat time.Time) (outbox.BarisEfekKeluar, error) {
			return s.antreanLampiran.PungutRujukan(ctx, tx, rujukan, saat)
		})
		if errors.Is(err, outbox.ErrEfekTidakAda) {
			return
		}
		if err != nil {
			logLampiran("%s (dokumen %s): %v", pesanAntreanTertunda, rujukan, err)
			return
		}
	}
}

// JalankanAntreanLampiran memungut sampai `maks` efek NBFACIN yang jatuh tempo.
func (s *Service) JalankanAntreanLampiran(ctx context.Context, maks int) (int, error) {
	if err := s.siapLampiran(); err != nil {
		return 0, err
	}
	n := 0
	for n < maks {
		err := s.satuPutaranLampiran(ctx, s.sekarang(), func(tx *db.Tx, saat time.Time) (outbox.BarisEfekKeluar, error) {
			return s.antreanLampiran.Pungut(ctx, tx, saat)
		})
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

// JalankanPekerjaLampiran - antrean setiap `interval` sampai ctx selesai; interval <= 0 = mati.
func (s *Service) JalankanPekerjaLampiran(ctx context.Context, interval time.Duration, catat func(string)) {
	if interval <= 0 {
		return
	}
	if catat == nil {
		catat = func(string) {}
	}
	t := time.NewTicker(interval)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			n, err := s.JalankanAntreanLampiran(ctx, batasPutaranPekerja)
			if err != nil && !errors.Is(err, context.Canceled) {
				catat(fmt.Sprintf("nbfacin lampiran: %d efek dijalankan, lalu berhenti: %v", n, err))
			} else if n > 0 {
				catat(fmt.Sprintf("nbfacin lampiran: %d efek dijalankan", n))
			}
		}
	}
}

// satuPutaranLampiran - pungut, jalankan, tuntaskan: SATU transaksi (pola PekerjaEfek / Treaty Contract Out).
func (s *Service) satuPutaranLampiran(ctx context.Context, saat time.Time,
	pungut func(tx *db.Tx, saat time.Time) (outbox.BarisEfekKeluar, error)) error {
	return s.transaksi(ctx, func(tx *db.Tx) error {
		b, err := pungut(tx, saat)
		if err != nil {
			return err
		}
		jalanErr := s.laksanakanLampiran(ctx, tx, b)
		if jalanErr == nil {
			return s.antreanLampiran.Tuntaskan(ctx, tx, b.ID, outbox.StatusEfekSelesai, time.Time{}, "", saat)
		}
		status, jadwal := outbox.StatusEfekAntre, saat.Add(outbox.Backoff(b.Percobaan+1))
		if !layakUlangLampiran(jalanErr) || b.Percobaan >= outbox.PercobaanMaksimum {
			status, jadwal = outbox.StatusEfekGagalPermanen, time.Time{}
		}
		return s.antreanLampiran.Tuntaskan(ctx, tx, b.ID, status, jadwal, ringkasGalatLampiran(jalanErr), saat)
	})
}

// layakUlangLampiran - LayakDicobaUlang + keadaan permanen milik lampiran.
func layakUlangLampiran(err error) bool {
	for _, permanen := range []error{ErrBerkasSumberLampiranHilang, errMuatanLampiranRusak, errJenisEfekAsing,
		unggah.ErrUnggahanDirBelumDisetel, ErrStorageMenolakPermintaan, ErrObjekNyataDiStub, layanan.ErrGaramTokenKosong,
		layanan.ErrAppNameKosong, repository.ErrAppStorageKosong} {
		if errors.Is(err, permanen) {
			return false
		}
	}
	return outbox.LayakDicobaUlang(err)
}

func ringkasGalatLampiran(err error) string {
	r := []rune(err.Error())
	if len(r) > batasGalatTersimpan {
		return string(r[:batasGalatTersimpan])
	}
	return string(r)
}

func (s *Service) laksanakanLampiran(ctx context.Context, tx *db.Tx, b outbox.BarisEfekKeluar) error {
	var m muatanLampiran
	if err := json.Unmarshal([]byte(b.Muatan), &m); err != nil || strings.TrimSpace(m.ImageID) == "" {
		return fmt.Errorf("%w: efek %s", errMuatanLampiranRusak, b.ID)
	}
	switch b.Jenis {
	case unggah.JenisEfekStorageUnggah:
		return s.kirimUnggahLampiran(ctx, tx, b.Rujukan, m)
	case unggah.JenisEfekStorageHapus:
		return s.kirimHapusLampiran(ctx, tx, m)
	}
	return fmt.Errorf("%w: %q", errJenisEfekAsing, b.Jenis)
}

// kirimUnggahLampiran - InsertGoogleStorage_Act + Insert_T_Storage_SQL, IDEMPOTEN: dokumen sudah dihapus / objek sudah
// tercatat = selesai; berkas antrean hilang tetapi objeknya ada di penyimpanan = dicatat dari jawaban geturl.
func (s *Service) kirimUnggahLampiran(ctx context.Context, tx *db.Tx, rujukan string, m muatanLampiran) error {
	d, err := s.lampiran.KunciLampiran(ctx, tx, rujukan)
	if errors.Is(err, repository.ErrLampiranTidakAda) {
		return nil
	}
	if err != nil {
		return err
	}
	if d.Tersimpan {
		return nil
	}
	lok, ok := models.LokasiUnggah(d.ID, d.TStorageID, d.NamaFile)
	if !ok || d.TStorageID != m.ImageID {
		return fmt.Errorf("%w: dokumen %s", errMuatanLampiranRusak, d.ID)
	}
	jalur := s.jalurAntre(d.TStorageID, d.NamaFile)
	f, err := os.Open(jalur)
	if os.IsNotExist(err) {
		ada, objek, errAda := s.storageLampiran.Ada(ctx, lok)
		if errAda != nil {
			return errAda
		}
		if !ada {
			return fmt.Errorf("%w: dokumen %s", ErrBerkasSumberLampiranHilang, d.ID)
		}
		objek.ImageID, objek.Namafile = d.TStorageID, lok.Namafile
		return s.lampiran.SimpanObjek(ctx, tx, objek)
	}
	if err != nil {
		return fmt.Errorf("services: membuka berkas antrean dokumen %s: %w", d.ID, err)
	}
	objek, simpanErr := s.storageLampiran.Simpan(ctx, lok, f, d.Mime, unggah.MimeDariNamaFile(d.NamaFile))
	_ = f.Close()
	if simpanErr != nil {
		return fmt.Errorf("services: mengunggah dokumen %s: %w", d.ID, simpanErr)
	}
	objek.ImageID = d.TStorageID
	if objek.Namafile == "" {
		objek.Namafile = lok.Namafile
	}
	if err := s.lampiran.SimpanObjek(ctx, tx, objek); err != nil {
		return err
	}
	// Berkas antrean dibuang sesudah dicatat; bila transaksinya batal, percobaan berikutnya menempuh jalan Ada di atas.
	_ = os.Remove(jalur)
	return nil
}

// kirimHapusLampiran - DeleteGoogleStorage_Act: GetLinkStorage_SQL -> delete -> DeleteStorage_SQL; berkas antrean dibuang.
// Objek yang sudah tidak ada bukan kegagalan. ⛔ Hapus penyimpanan (yang dapat menyegarkan catatan di transaksinya sendiri)
// SEBELUM HapusObjek di tx ini - dibalik, transaksi pendek itu menunggu kunci baris milik tx ini.
func (s *Service) kirimHapusLampiran(ctx context.Context, tx *db.Tx, m muatanLampiran) error {
	link, ada, err := s.lampiran.LinkStorage(ctx, tx, m.ImageID)
	if err != nil {
		return err
	}
	if ada {
		if err := s.storageLampiran.Hapus(ctx, link); err != nil && !errors.Is(err, ErrBerkasTidakAdaDiPenyimpanan) {
			return fmt.Errorf("services: menghapus berkas dokumen %s: %w", m.DokumenID, err)
		}
		if err := s.lampiran.HapusObjek(ctx, tx, m.ImageID); err != nil {
			return err
		}
	}
	if m.BerkasAntre != "" && m.BerkasAntre == filepath.Base(m.BerkasAntre) {
		if err := os.Remove(filepath.Join(s.folderAntre(), m.BerkasAntre)); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("services: membuang berkas antrean dokumen %s: %w", m.DokumenID, err)
		}
	}
	return nil
}

func logLampiran(format string, a ...any) { log.Printf("nbfacin lampiran: "+format, a...) }
