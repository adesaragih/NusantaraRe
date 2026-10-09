package services

// Lampiran berkas - keputusan work owner 08-10-2026: dibangun "kaya XML nya" di atas penyimpanan bersama
// `inti/backend/penyimpanan` (kelas Pega `T_STORAGE_IMAGE`, dipakai semua modul).
//
//	grid kategori   `AttachmentsBdx` (deferred `GetKategotyDocBdx` b880)          KategoriLampiran
//	popup berkas    `AttachmentDetailBdx` (deferred `getAttcachmentList` b255)    DaftarLampiran
//	Upload File     flow action `BordereauxAttach` -> `AttachDocBdx_Post`       UnggahLampiran
//	nama berkas     `DownloadAttachmentBdx` (ViewOffice=false, b2152-b2175)       UnduhLampiran
//	View Office     `DownloadAttachmentBdx` (ViewOffice=true, b2606)              TautanOffice
//	Delete          `DeleteAttachmentBdx` (b3208-b3231)                           HapusLampiran
//
// ⚠️ PENYIMPANGAN SADAR (bug Pega diperbaiki di Go, keputusan work owner 04-10-2026):
//   - Lampiran hanya untuk berkas TERSIMPAN: BDX_ID di sini lahir saat Save (Pega: `InputNew`).
//   - Berkas yang `InsertGoogleStorage_Act` lewati (jenis tak dikenal, isi kosong) DITOLAK; Pega tetap menyisipkan
//     baris lampiran tanpa `T_STORAGE_ID`.
//   - Delete: hapus di penyimpanan gagal = baris tetap (dapat diulang); Pega tetap menghapus barisnya.
//   - Unduh: isi dibaca backend dari URL bertanda tangan - frontend dilarang membuka jendela (penjaga lintas modul
//     `unduhdokumen.test.ts`); Pega membuka URL itu di jendela popup (`openUrlInWindow` b2217).

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"nusantarare/inti/backend/penyimpanan"
	"nusantarare/inti/backend/unggah"
	"nusantarare/modul/bordereaux/backend/models"
	"nusantarare/modul/bordereaux/backend/repository"
)

const (
	// FolderLampiran / DurasiLampiran - `AttachDocBdx_Post` 1.3 (`Folder = "Contract"`, `Durasi = 3600`);
	// `DownloadAttachmentBdx` 1 memakai Durasi 3600 yang sama.
	FolderLampiran = "Contract"
	DurasiLampiran = 3600
	// batasNamaLampiran - `M_ATTACHMENTBORDEREAUX.FILENAME VARCHAR2(1000)` (katalog DEV 08-10-2026).
	batasNamaLampiran = 1000
)

// ekstensiOffice - syarat tampil `View Office Online` (`AttachmentDetailBdx` b2876): `.HASIL5` (FILEMIMETYPE, berisi
// ekstensi) xls, xlsx, doc, docx, ppt, pptx.
var ekstensiOffice = map[string]bool{"xls": true, "xlsx": true, "doc": true, "docx": true, "ppt": true, "pptx": true}

// ErrLampiranTidakAda - lampiran tidak ada di berkas dan kategori itu.
var ErrLampiranTidakAda = errors.New("services: lampiran tidak ada")

// PenyimpananBerkas - bagian `inti/backend/penyimpanan` yang dipakai modul ini (tiruan di uji).
type PenyimpananBerkas interface {
	Unggah(ctx context.Context, m penyimpanan.MasukUnggah) (penyimpanan.Objek, error)
	Catat(ctx context.Context, tx *dbTx, o penyimpanan.Objek) error
	Buka(ctx context.Context, imageID string, durasi int, pengguna string) (io.ReadCloser, error)
	Tautan(ctx context.Context, imageID string, durasi int, pengguna string) (string, error)
	HapusObjek(ctx context.Context, imageID, pengguna string) error
	HapusCatatan(ctx context.Context, tx *dbTx, imageID string) error
}

// DenganPenyimpanan memasang penyimpanan berkas (Oracle: `LayananOracle`; uji: tiruan).
func (l *Layanan) DenganPenyimpanan(p PenyimpananBerkas) *Layanan {
	salin := *l
	salin.berkas = p
	return &salin
}

// IDLampiran - `@FormatDateTime(@CurrentDateTime(), "yyyyMMddHHmmssSSS", "Asia/Jakarta")` (`AttachDocBdx_Post` b392);
// bentrok ditangani pemanggil.
func IDLampiran(t time.Time) string {
	t = t.In(lokasiJakarta)
	return t.Format("20060102150405") + milidetik(t)
}

// BolehLampiran - `Upload File` (`AttachmentsBdx` b4459) dan `Delete` (`AttachmentDetailBdx` b3376): hanya berkas yang
// boleh di-Edit (`HakAtas(...).Ubah` - di tangan pembuatnya, belum Resolve-Complete, menu PENUH, pembuat atau
// superadmin). Pega: `BORDEREAUX.ViewStage != 1 || OperatorID.pyPosition = 'IT Developer'`; pengecualian IT Developer
// DIBUANG - keputusan work owner 08-10-2026: "kalo cuman view jangan bisa upload dan delete". Mode View
// (`ViewData=1`) ditegakkan layar; menu View only ditolak di sini juga.
func BolehLampiran(a Aktor, h models.Header) bool {
	return HakAtas(a, h).Ubah
}

// berkasLampiran - header berkas yang lampirannya disentuh.
func (l *Layanan) berkasLampiran(ctx context.Context, bdxID string) (models.Header, error) {
	h, err := l.gudang.AmbilHeader(ctx, nil, strings.TrimSpace(bdxID))
	if errors.Is(err, repository.ErrTidakAda) {
		return models.Header{}, ErrTidakAda
	}
	return h, err
}

// KategoriLampiran - grid `AttachmentsBdx`: setiap kategori master dan jumlah lampiran berkas ini.
func (l *Layanan) KategoriLampiran(ctx context.Context, _ Aktor, bdxID string) ([]models.KategoriLampiran, error) {
	h, err := l.berkasLampiran(ctx, bdxID)
	if err != nil {
		return nil, err
	}
	return l.gudang.KategoriLampiran(ctx, h.BdxID)
}

// DaftarLampiran - popup `AttachmentDetailBdx`: lampiran berkas ini di satu kategori.
func (l *Layanan) DaftarLampiran(ctx context.Context, _ Aktor, bdxID, kategoriID string) ([]models.Lampiran, error) {
	h, err := l.berkasLampiran(ctx, bdxID)
	if err != nil {
		return nil, err
	}
	return l.gudang.DaftarLampiran(ctx, h.BdxID, strings.TrimSpace(kategoriID))
}

// namaBerkasUnggahan - nama berkas tanpa jalur folder peramban.
func namaBerkasUnggahan(nama string) string {
	nama = strings.TrimSpace(nama)
	if i := strings.LastIndexAny(nama, `/\`); i >= 0 {
		nama = nama[i+1:]
	}
	return strings.TrimSpace(nama)
}

// UnggahLampiran - `AttachDocBdx_Post` untuk satu berkas: `InsertGoogleStorage_Act` (Folder "Contract", Durasi
// 3600), lalu objek (`Insert_T_Storage_SQL`) dan baris lampiran (`AttachDocumentBdx_SQL` save) di SATU transaksi.
// CATEGORY_ID = kategori grid (`Primary.MASTERID` b445), CATEGORY = namanya (`Primary.TYPE` b466), FILEMIMETYPE =
// ekstensi (`.pyFileMimeType` b508), USERNAME = pengunggah.
func (l *Layanan) UnggahLampiran(ctx context.Context, a Aktor, bdxID, kategoriID, nama string, isi []byte) (models.Lampiran, error) {
	h, err := l.berkasLampiran(ctx, bdxID)
	if err != nil {
		return models.Lampiran{}, err
	}
	if !BolehLampiran(a, h) {
		return models.Lampiran{}, dilarang("you may not add attachments to this bordereaux")
	}
	kategoriID = strings.TrimSpace(kategoriID)
	kategori, ada, err := l.gudang.NamaKategoriLampiran(ctx, kategoriID)
	if err != nil {
		return models.Lampiran{}, err
	}
	if !ada {
		return models.Lampiran{}, tolak("attachment category %q is not in the category master", kategoriID)
	}
	nama = namaBerkasUnggahan(nama)
	switch {
	case nama == "":
		return models.Lampiran{}, tolak("No file attached")
	case len(nama) > batasNamaLampiran:
		return models.Lampiran{}, tolak("the file name is longer than %d characters", batasNamaLampiran)
	case len(isi) > unggah.BatasUkuranUnggahan:
		return models.Lampiran{}, tolak("the file is larger than %d MB", unggah.BatasUkuranUnggahan>>20)
	}
	ext := penyimpanan.Ekstensi(nama)
	o, err := l.penyimpanan().Unggah(ctx, penyimpanan.MasukUnggah{Folder: FolderLampiran, NamaFile: nama, Isi: isi,
		Durasi: DurasiLampiran, Ext: ext, Pengguna: a.AkunID})
	if err != nil {
		return models.Lampiran{}, err
	}
	baru := models.Lampiran{BdxID: h.BdxID, KategoriID: kategoriID, Kategori: kategori, FileName: nama, Ekstensi: ext,
		Username: a.AkunID, StorageID: o.ImageID}
	jam := l.sekarang()
	err = l.tx(ctx, func(tx *dbTx) error {
		if err := l.penyimpanan().Catat(ctx, tx, o); err != nil {
			return err
		}
		for i := 0; ; i++ {
			if i >= cobaIDMaksimum {
				return fmt.Errorf("services: tidak menemukan ID lampiran kosong")
			}
			baru.ID = IDLampiran(jam.Add(time.Duration(i) * time.Millisecond))
			err := l.gudang.SisipLampiran(ctx, tx, baru)
			if errors.Is(err, repository.ErrIDTerpakai) {
				continue
			}
			return err
		}
	})
	if err != nil {
		// Objeknya sudah di penyimpanan tanpa catatan - dicatat di log untuk dibersihkan.
		log.Printf("bordereaux: lampiran %s berkas %s tidak tercatat, objek penyimpanan tertinggal: %v", o.ImageID, h.BdxID, err)
		return models.Lampiran{}, err
	}
	return baru, nil
}

// ambilLampiran - lampiran berkas dan kategori itu.
func (l *Layanan) ambilLampiran(ctx context.Context, bdxID, kategoriID, id string) (models.Header, models.Lampiran, error) {
	h, err := l.berkasLampiran(ctx, bdxID)
	if err != nil {
		return models.Header{}, models.Lampiran{}, err
	}
	a, ada, err := l.gudang.AmbilLampiran(ctx, h.BdxID, strings.TrimSpace(kategoriID), strings.TrimSpace(id))
	if err != nil {
		return models.Header{}, models.Lampiran{}, err
	}
	if !ada {
		return models.Header{}, models.Lampiran{}, ErrLampiranTidakAda
	}
	return h, a, nil
}

// IsiLampiran - isi satu lampiran untuk diunduh.
type IsiLampiran struct {
	Isi  io.ReadCloser
	Nama string
	Mime string
}

// UnduhLampiran - tautan nama berkas: `DownloadAttachmentBdx` (`GetUrlGoogleStorage_Act`, ImageID = `.HASIL4`,
// Durasi 3600), isinya dibaca backend.
func (l *Layanan) UnduhLampiran(ctx context.Context, a Aktor, bdxID, kategoriID, id string) (IsiLampiran, error) {
	_, b, err := l.ambilLampiran(ctx, bdxID, kategoriID, id)
	if err != nil {
		return IsiLampiran{}, err
	}
	if strings.TrimSpace(b.StorageID) == "" {
		return IsiLampiran{}, penyimpanan.ErrObjekTidakAda
	}
	isi, err := l.penyimpanan().Buka(ctx, b.StorageID, DurasiLampiran, a.AkunID)
	if err != nil {
		return IsiLampiran{}, err
	}
	return IsiLampiran{Isi: isi, Nama: b.FileName, Mime: unggah.MimeDariNamaFile(b.FileName)}, nil
}

// TautanOffice - `View Office Online`: URL bertanda tangan untuk penampil kantor (`DownloadAttachmentBdx` 2 b388,
// pembungkusan penampil di layar). Hanya xls/xlsx/doc/docx/ppt/pptx (b2876).
func (l *Layanan) TautanOffice(ctx context.Context, a Aktor, bdxID, kategoriID, id string) (string, error) {
	_, b, err := l.ambilLampiran(ctx, bdxID, kategoriID, id)
	if err != nil {
		return "", err
	}
	if !ekstensiOffice[strings.ToLower(b.Ekstensi)] {
		return "", tolak("View Office Online is only for Excel, Word, and PowerPoint files")
	}
	if strings.TrimSpace(b.StorageID) == "" {
		return "", penyimpanan.ErrObjekTidakAda
	}
	return l.penyimpanan().Tautan(ctx, b.StorageID, DurasiLampiran, a.AkunID)
}

// HapusLampiran - `DeleteAttachmentBdx`: `DeleteGoogleStorage_Act` (objek di penyimpanan), lalu `DeleteStorage_SQL`
// dan `AttachDocumentBdx_SQL` delete di SATU transaksi. Baris lama tanpa `T_STORAGE_ID` langsung dihapus.
func (l *Layanan) HapusLampiran(ctx context.Context, a Aktor, bdxID, kategoriID, id string) error {
	h, b, err := l.ambilLampiran(ctx, bdxID, kategoriID, id)
	if err != nil {
		return err
	}
	if !BolehLampiran(a, h) {
		return dilarang("you may not delete attachments of this bordereaux")
	}
	if b.StorageID != "" {
		if err := l.penyimpanan().HapusObjek(ctx, b.StorageID, a.AkunID); err != nil {
			return err
		}
	}
	return l.tx(ctx, func(tx *dbTx) error {
		if b.StorageID != "" {
			if err := l.penyimpanan().HapusCatatan(ctx, tx, b.StorageID); err != nil {
				return err
			}
		}
		return l.gudang.HapusLampiran(ctx, tx, h.BdxID, b.ID)
	})
}

// penyimpanan - penyimpanan terpasang; belum dipasang = gagal terang.
func (l *Layanan) penyimpanan() PenyimpananBerkas {
	if l.berkas == nil {
		return penyimpananBelumDipasang{}
	}
	return l.berkas
}

type penyimpananBelumDipasang struct{}

var errPenyimpananBelumDipasang = fmt.Errorf("%w: bordereaux file storage is not wired", penyimpanan.ErrStorageBelumSiap)

func (penyimpananBelumDipasang) Unggah(context.Context, penyimpanan.MasukUnggah) (penyimpanan.Objek, error) {
	return penyimpanan.Objek{}, errPenyimpananBelumDipasang
}
func (penyimpananBelumDipasang) Catat(context.Context, *dbTx, penyimpanan.Objek) error {
	return errPenyimpananBelumDipasang
}
func (penyimpananBelumDipasang) Buka(context.Context, string, int, string) (io.ReadCloser, error) {
	return nil, errPenyimpananBelumDipasang
}
func (penyimpananBelumDipasang) Tautan(context.Context, string, int, string) (string, error) {
	return "", errPenyimpananBelumDipasang
}
func (penyimpananBelumDipasang) HapusObjek(context.Context, string, string) error {
	return errPenyimpananBelumDipasang
}
func (penyimpananBelumDipasang) HapusCatatan(context.Context, *dbTx, string) error {
	return errPenyimpananBelumDipasang
}
