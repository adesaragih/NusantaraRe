// Package dokumenpolis adalah lampiran kasus berkategori "Reas" - padanan kelas Pega `ASM-FW-GISFW-Int-DOCUMENT_POLIS`
// dan grid `AttachmentGridReas` (kelas `ASM-FW-GISFW-Data-OfferFacIn`), yang di Pega dipakai NB FacIn dan dipinjam NB
// Treaty In / EDM Treaty In (`InputPolicyTreatyInPre_Act` -> `SetCategoryAttach`). Keputusan work owner 08-10-2026:
// "bisa buatkan untuk NB dan EDM treaty?" - panel di bawah layar kasus; Upload / Delete "semua bisa asal belum resolve".
//
// Sumber XML (folder korpus NB FacIn - grid dan aksinya TIDAK ada di ekspor NB / EDM Treaty In):
//
//	grid kategori   `Section/AttachmentGridReas` (deferred `SetCategoryAttach`): Category b1340, Count b1493,
//	                Upload File b1643, View File b1791; kategori = `CategoryAttach_SQL` (`CATEGORY_ATTACH_REAS`
//	                urut NOTE), Count = `InputParamUploadReas_act` (Obj-Browse IDPEGA = pzInsKey, KATEGORI_2 = NOTE)
//	popup berkas    `Harness/ViewPictureList` -> `Section/ReasViewAttachment`: File b1363 (`DownloadDocumentPolis`),
//	                View Office Online b2670 (MIME xls..pptx dan T_STORAGE_ID, b2951), Note b1618 (KATEGORI_2),
//	                Upload Date b1762 (TANGGAL), Delete b3434 (`DeleteDocumentPolis_Act`)
//	unggah          `InsertDocument_Act`: ID `yyyyMMddhhmmssSSS` Asia/Jakarta, MIME = ekstensi, KATEGORI_1 "Reas",
//	                `InsertGoogleStorage_Act` (Folder "Policy", Durasi 1800), Obj-Save bila T_STORAGE_ID terisi
//	unduh / office  `DownloadDocumentPolis`: `GetUrlGoogleStorage_Act` Durasi 1800 (+ penampil Office)
//	hapus           `DeleteDocumentPolis_Act`: `DeleteGoogleStorage_Act` bila T_STORAGE_ID terisi, lalu Obj-Delete
//
// ⛔ Kunci kasus: lampiran BARU ber-IDPEGA = ID T_WORK_POLIS apa adanya (`KunciInstans`, keputusan work owner
// 06-10-2026 "key dari T_WORK_POLIS jangan diubah"); lampiran Pega lama berkunci pzInsKey `ASM-FW-GISFW-WORK <pyID>`
// (6 dokumen NB Treaty In di DEV 08-10-2026) - keduanya DIBACA (`KasusDari`).
//
// ⚠️ PENYIMPANGAN SADAR: Delete Pega bersyarat `FlagOnGoingPolicy = '0' || IT` (b3654) - di sini keputusan work owner
// (belum Resolve, ditentukan modul lewat `Kasus.BolehUbah`); berkas yang `InsertGoogleStorage_Act` lewati DITOLAK
// (Pega tidak menyimpan barisnya, tanpa pesan); hapus jarak jauh yang gagal menahan barisnya; unduh dan View dibaca
// backend (penjaga `unduhdokumen.test.ts`).
package dokumenpolis

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/penyimpanan"
	"nusantarare/inti/backend/unggah"
)

const (
	// KelasKasusPega - kelas dalam pzInsKey kasus Pega (`<kelas> <pyID>`), IDPEGA lampiran Pega lama.
	KelasKasusPega = "ASM-FW-GISFW-WORK"
	// KategoriReas - `KATEGORI_1` (`SetCategoryAttach` 4.1 `Param.category = "Reas"`).
	KategoriReas = "Reas"
	// FolderBerkas / DurasiBerkas - `InsertDocument_Act` 4 dan `DownloadDocumentPolis` 2.
	FolderBerkas = "Policy"
	DurasiBerkas = 1800
	// batasNama - `DOCUMENT_POLIS.NAMAFILE VARCHAR2(1000)` (katalog DEV 08-10-2026).
	batasNama = 1000
	// cobaIDMaksimum - ID bentrok (PK `ID`) dicoba milidetik berikutnya.
	cobaIDMaksimum = 1000
)

var (
	// ErrTidakAda - dokumen tidak ada di kasus itu (404).
	ErrTidakAda = errors.New("dokumenpolis: attachment not found in this case")
	// ErrDilarang - kasus sudah Resolve: Upload / Delete ditolak (403).
	ErrDilarang = errors.New("dokumenpolis: attachments of a resolved case cannot be changed")
	// ErrMasukanTidakSah - permintaan ditolak; kalimatnya untuk pengguna (422).
	ErrMasukanTidakSah = errors.New("dokumenpolis: invalid input")
	// ErrIDTerpakai - ORA-00001 saat menyisipkan (PK `ID`).
	ErrIDTerpakai = errors.New("dokumenpolis: ID already used")
)

func tolak(format string, a ...any) error {
	return fmt.Errorf("%w: %s", ErrMasukanTidakSah, fmt.Sprintf(format, a...))
}

// Kasus - kasus pemilik lampiran, disusun modul pemakai.
type Kasus struct {
	// Tulis - IDPEGA lampiran baru (`KunciInstans` modul).
	Tulis string
	// Baca - IDPEGA yang dibaca: Tulis dan pzInsKey Pega lama.
	Baca []string
	// BolehUbah - Upload File / Delete (keputusan work owner 08-10-2026: belum Resolve).
	BolehUbah bool
}

// KasusDari menyusun Kasus dari kunci kasus (`KunciInstans`) dan pyID-nya: lampiran Pega lama berkunci
// `ASM-FW-GISFW-WORK <pyID>`. Kunci yang sudah berbentuk pzInsKey (salinan Copy Old) tidak digandakan.
func KasusDari(kunci, pyID string, bolehUbah bool) Kasus {
	kunci, pyID = strings.TrimSpace(kunci), strings.TrimSpace(pyID)
	baca := []string{kunci}
	if lama := KelasKasusPega + " " + pyID; pyID != "" && lama != kunci {
		baca = append(baca, lama)
	}
	return Kasus{Tulis: kunci, Baca: baca, BolehUbah: bolehUbah}
}

// Kategori - satu baris `AttachmentGridReas`: NOTE dan jumlah dokumen kasus di kategori itu (`.CountAttach`).
type Kategori struct {
	Nama  string `json:"nama"`
	Cacah int    `json:"cacah"`
}

// Dokumen - satu baris `DOCUMENT_POLIS` di popup `ReasViewAttachment`.
type Dokumen struct {
	ID       string `json:"id"`
	NamaFile string `json:"namaFile"`
	// Ekstensi - `MIME` (`@toLowerCase(Param.MIME)`, berisi ekstensi; DEV: `pdf`).
	Ekstensi string `json:"ekstensi"`
	// Kategori - `KATEGORI_2` (kolom Note).
	Kategori string `json:"kategori"`
	// Tanggal - `TANGGAL` (kolom Upload Date), `DD-MM-YYYY HH24:MI`.
	Tanggal    string `json:"tanggal"`
	Pengunggah string `json:"pengunggah"`
	// StorageID - `T_STORAGE_ID`; tidak dikirim ke layar.
	StorageID string `json:"-"`
	// AdaObjek - `T_STORAGE_ID` terisi (layar: View Office Online hanya bila ada objek, b2951).
	AdaObjek bool `json:"adaObjek"`
}

// Catatan - `DOCUMENT_POLIS` dan `CATEGORY_ATTACH_REAS` (Oracle di `catatan.go`, tiruan di uji).
type Catatan interface {
	Kategori(ctx context.Context, baca []string) ([]Kategori, error)
	AdaKategori(ctx context.Context, nama string) (bool, error)
	Daftar(ctx context.Context, baca []string, kategori string) ([]Dokumen, error)
	Ambil(ctx context.Context, baca []string, id string) (Dokumen, bool, error)
	// Sisip - `InsertDocument_Act` Obj-Save; ID bentrok = ErrIDTerpakai.
	Sisip(ctx context.Context, tx *db.Tx, d Dokumen, idPega string) error
	Hapus(ctx context.Context, tx *db.Tx, baca []string, id string) error
}

// Berkas - bagian `inti/backend/penyimpanan` yang dipakai paket ini.
type Berkas interface {
	Unggah(ctx context.Context, m penyimpanan.MasukUnggah) (penyimpanan.Objek, error)
	Catat(ctx context.Context, tx *db.Tx, o penyimpanan.Objek) error
	Buka(ctx context.Context, imageID string, durasi int, pengguna string) (io.ReadCloser, error)
	Tautan(ctx context.Context, imageID string, durasi int, pengguna string) (string, error)
	HapusObjek(ctx context.Context, imageID, pengguna string) error
	HapusCatatan(ctx context.Context, tx *db.Tx, imageID string) error
}

// Transaksi menjalankan fn di dalam satu transaksi.
type Transaksi func(ctx context.Context, fn func(tx *db.Tx) error) error

// Layanan - aturan lampiran "Reas".
type Layanan struct {
	catatan Catatan
	berkas  Berkas
	tx      Transaksi
	jam     func() time.Time
}

// Baru menyusun layanan; jam nil = `time.Now`.
func Baru(c Catatan, b Berkas, tx Transaksi, jam func() time.Time) *Layanan {
	if jam == nil {
		jam = time.Now
	}
	return &Layanan{catatan: c, berkas: b, tx: tx, jam: jam}
}

// Oracle menyusun layanan di atas Oracle akar `akar`; garam token penyimpanan dari `config.StorageTokenSalt`.
func Oracle(akar inti.Akar, garam string) *Layanan {
	return Baru(catatanOracle{db: akar.DB()}, penyimpanan.Oracle(akar, garam), akar.DalamTransaksi, nil)
}

// zonaJakarta - `@CurrentDate(…, "Asia/Jakarta")`.
var zonaJakarta = time.FixedZone("WIB", 7*3600)

// IDDokumen - `@CurrentDate("yyyyMMddhhmmssSSS","Asia/Jakarta")` (`InsertDocument_Act` 3) - format XML apa adanya,
// jam 12-an (`hh`); bentrok ditangani pemanggil.
func IDDokumen(t time.Time) string {
	t = t.In(zonaJakarta)
	return t.Format("20060102030405") + fmt.Sprintf("%03d", t.Nanosecond()/int(time.Millisecond))
}

// Kategori - grid `AttachmentGridReas`.
func (l *Layanan) Kategori(ctx context.Context, k Kasus) ([]Kategori, error) {
	return l.catatan.Kategori(ctx, k.Baca)
}

// Daftar - popup `ReasViewAttachment`: dokumen kasus di satu kategori.
func (l *Layanan) Daftar(ctx context.Context, k Kasus, kategori string) ([]Dokumen, error) {
	return l.catatan.Daftar(ctx, k.Baca, strings.TrimSpace(kategori))
}

// namaBerkasUnggahan - nama berkas tanpa jalur folder peramban.
func namaBerkasUnggahan(nama string) string {
	nama = strings.TrimSpace(nama)
	if i := strings.LastIndexAny(nama, `/\`); i >= 0 {
		nama = nama[i+1:]
	}
	return strings.TrimSpace(nama)
}

// Unggah - `InsertDocument_Act` untuk satu berkas: `InsertGoogleStorage_Act` (Folder "Policy", Durasi 1800), lalu objek
// (`Insert_T_Storage_SQL`) dan baris `DOCUMENT_POLIS` di SATU transaksi.
func (l *Layanan) Unggah(ctx context.Context, k Kasus, akun, kategori, nama string, isi []byte) (Dokumen, error) {
	if !k.BolehUbah {
		return Dokumen{}, ErrDilarang
	}
	kategori = strings.TrimSpace(kategori)
	ada, err := l.catatan.AdaKategori(ctx, kategori)
	if err != nil {
		return Dokumen{}, err
	}
	if !ada {
		return Dokumen{}, tolak("attachment category %q is not in CATEGORY_ATTACH_REAS", kategori)
	}
	nama = namaBerkasUnggahan(nama)
	switch {
	case nama == "":
		return Dokumen{}, tolak("No file attached")
	case len(nama) > batasNama:
		return Dokumen{}, tolak("the file name is longer than %d characters", batasNama)
	case len(isi) > unggah.BatasUkuranUnggahan:
		return Dokumen{}, tolak("the file is larger than %d MB", unggah.BatasUkuranUnggahan>>20)
	}
	ext := penyimpanan.Ekstensi(nama) // 2: MIME dari nama berkas
	o, err := l.berkas.Unggah(ctx, penyimpanan.MasukUnggah{Folder: FolderBerkas, NamaFile: nama, Isi: isi,
		Durasi: DurasiBerkas, Ext: ext, Pengguna: akun})
	if err != nil {
		return Dokumen{}, err
	}
	d := Dokumen{NamaFile: nama, Ekstensi: ext, Kategori: kategori, Pengunggah: akun, StorageID: o.ImageID, AdaObjek: true}
	jam := l.jam()
	err = l.tx(ctx, func(tx *db.Tx) error {
		if err := l.berkas.Catat(ctx, tx, o); err != nil {
			return err
		}
		for i := 0; ; i++ {
			if i >= cobaIDMaksimum {
				return errors.New("dokumenpolis: no free document ID")
			}
			d.ID = IDDokumen(jam.Add(time.Duration(i) * time.Millisecond))
			err := l.catatan.Sisip(ctx, tx, d, k.Tulis)
			if errors.Is(err, ErrIDTerpakai) {
				continue
			}
			return err
		}
	})
	if err != nil {
		log.Printf("dokumenpolis: dokumen %s kasus %s tidak tercatat, objek penyimpanan tertinggal: %v", o.ImageID, k.Tulis, err)
		return Dokumen{}, err
	}
	return d, nil
}

func (l *Layanan) ambil(ctx context.Context, k Kasus, id string) (Dokumen, error) {
	d, ada, err := l.catatan.Ambil(ctx, k.Baca, strings.TrimSpace(id))
	if err != nil {
		return Dokumen{}, err
	}
	if !ada {
		return Dokumen{}, ErrTidakAda
	}
	return d, nil
}

// Isi - isi satu dokumen untuk diunduh / View.
type Isi struct {
	Isi  io.ReadCloser
	Nama string
	Mime string
}

// Unduh - tautan File (`DownloadDocumentPolis`, `GetUrlGoogleStorage_Act` Durasi 1800); isinya dibaca backend.
func (l *Layanan) Unduh(ctx context.Context, k Kasus, akun, id string) (Isi, error) {
	d, err := l.ambil(ctx, k, id)
	if err != nil {
		return Isi{}, err
	}
	if d.StorageID == "" {
		return Isi{}, penyimpanan.ErrObjekTidakAda
	}
	r, err := l.berkas.Buka(ctx, d.StorageID, DurasiBerkas, akun)
	if err != nil {
		return Isi{}, err
	}
	return Isi{Isi: r, Nama: d.NamaFile, Mime: unggah.MimeDariNamaFile(d.NamaFile)}, nil
}

// ekstensiOffice - syarat View Office Online (`ReasViewAttachment` b2951).
var ekstensiOffice = map[string]bool{"xls": true, "xlsx": true, "doc": true, "docx": true, "ppt": true, "pptx": true}

// TautanOffice - View Office Online: URL bertanda tangan untuk penampil kantor (`DownloadDocumentPolis` 3).
func (l *Layanan) TautanOffice(ctx context.Context, k Kasus, akun, id string) (string, error) {
	d, err := l.ambil(ctx, k, id)
	if err != nil {
		return "", err
	}
	if !ekstensiOffice[strings.ToLower(d.Ekstensi)] {
		return "", tolak("View Office Online is only for Excel, Word, and PowerPoint files")
	}
	if d.StorageID == "" {
		return "", penyimpanan.ErrObjekTidakAda
	}
	return l.berkas.Tautan(ctx, d.StorageID, DurasiBerkas, akun)
}

// Hapus - `DeleteDocumentPolis_Act`: objek di penyimpanan (bila tercatat), lalu `DeleteStorage_SQL` dan baris
// `DOCUMENT_POLIS` di SATU transaksi.
func (l *Layanan) Hapus(ctx context.Context, k Kasus, akun, id string) error {
	if !k.BolehUbah {
		return ErrDilarang
	}
	d, err := l.ambil(ctx, k, id)
	if err != nil {
		return err
	}
	if d.StorageID != "" {
		if err := l.berkas.HapusObjek(ctx, d.StorageID, akun); err != nil {
			return err
		}
	}
	return l.tx(ctx, func(tx *db.Tx) error {
		if d.StorageID != "" {
			if err := l.berkas.HapusCatatan(ctx, tx, d.StorageID); err != nil {
				return err
			}
		}
		return l.catatan.Hapus(ctx, tx, k.Baca, d.ID)
	})
}
