package services

// Untuk apa berkas ini: LAMPIRAN KLAIM - `GCNMSaveAttachments` (Work-, diekspor work owner 08-10-2026) memanggil
// `InsertDocument_Act` untuk setiap berkas: InsertGoogleStorage_Act (Folder "Claim", Durasi 1800) lalu baris dokumen
// klaim (S5 Obj-Save hanya bila T_STORAGE_ID terisi). Pola lampiran Bordereaux: unggah ke `inti/backend/penyimpanan`,
// lalu T_STORAGE_IMAGE + baris dokumen klaim di SATU transaksi. Kategori = master `T_KATEGORI_DOC_KLAIM` TYPE_KLAIM
// FAC (18 baris di DEV, prompt §2 butir 9; pola Claim Prop). Cacah per kategori = `AttachCategory.pxResults` gerbang
// "Send Claim to Committee" (AttachmentProtect_ACT).

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
	"nusantarare/modul/claimfacin/backend/models"
	"nusantarare/modul/claimfacin/backend/repository"
)

// PenyimpananBerkas - bagian `inti/backend/penyimpanan` yang dipakai modul ini (InsertGoogleStorage_Act +
// Insert_T_Storage_SQL; tiruan di uji).
type PenyimpananBerkas interface {
	Unggah(ctx context.Context, m penyimpanan.MasukUnggah) (penyimpanan.Objek, error)
	Catat(ctx context.Context, tx *db.Tx, o penyimpanan.Objek) error
	Buka(ctx context.Context, imageID string, durasi int, pengguna string) (io.ReadCloser, error)
	Tautan(ctx context.Context, imageID string, durasi int, pengguna string) (string, error)
	HapusObjek(ctx context.Context, imageID, pengguna string) error
	HapusCatatan(ctx context.Context, tx *db.Tx, imageID string) error
}

// ErrPenyimpananBelumDipasang - layanan dirakit tanpa penyimpanan berkas (tanpa Oracle / salah rakit).
var ErrPenyimpananBelumDipasang = fmt.Errorf("%w: penyimpanan lampiran Claim Fac In belum dipasang",
	penyimpanan.ErrStorageBelumSiap)

// cobaIDDokumen - batas percobaan ID dokumen klaim bentrok (+1 milidetik setiap kali).
const cobaIDDokumen = 50

// DenganPenyimpanan memasang penyimpanan berkas lampiran (Oracle: `penyimpanan.Oracle`; uji: tiruan).
func (l *Layanan) DenganPenyimpanan(p PenyimpananBerkas) *Layanan {
	salin := *l
	salin.berkas = p
	return &salin
}

// LampiranKasus - kategori (AttachCategory) dan dokumen klaim satu kasus.
type LampiranKasus struct {
	Kategori []models.KategoriLampiran `json:"kategori"`
	Lampiran []models.Lampiran         `json:"lampiran"`
	// BolehUnggah - pelaku pemegang assignment kasus yang belum ditutup.
	BolehUnggah bool `json:"bolehUnggah"`
}

// DaftarLampiran - kategori + dokumen klaim kasus `id` (siapa pun yang boleh membuka kasus).
func (l *Layanan) DaftarLampiran(ctx context.Context, p inti.Pelaku, id string) (LampiranKasus, error) {
	if err := l.periksaPelaku(p); err != nil {
		return LampiranKasus{}, err
	}
	k, err := l.g.Keadaan(ctx, nil, id)
	if err != nil {
		return LampiranKasus{}, err
	}
	kat, err := l.g.KategoriLampiran(ctx, k.ID)
	if err != nil {
		return LampiranKasus{}, err
	}
	dok, err := l.g.DaftarLampiran(ctx, k.ID)
	if err != nil {
		return LampiranKasus{}, err
	}
	return LampiranKasus{Kategori: kat, Lampiran: dok, BolehUnggah: Pemegang(p, k)}, nil
}

// BerkasUnggahan - satu `dragDropFileUpload.pxResults` (.pyFileName, .pyFileSource, .pyCategory).
type BerkasUnggahan struct {
	Nama string
	Isi  []byte
	// Kategori - `.pyCategory` berkas ini (Add attachment: dipilih per berkas); dipakai bila kategori bersama kosong.
	Kategori string
}

// UnggahLampiran = GCNMSaveAttachments: kategori setiap berkas = S1.2 `@if(TempInputParam.pyCategory = "", .pyCategory,
// TempInputParam.pyCategory)` - `kategori` bersama (Upload File baris) atau `.Kategori` berkas (Add attachment), lalu
// InsertDocument_Act (S1.7). Seluruh kategori diperiksa sebelum berkas pertama diunggah; berkas diproses berurutan dan
// yang sudah tersimpan tetap tersimpan bila berkas sesudahnya gagal (Pega: satu Obj-Save per putaran).
func (l *Layanan) UnggahLampiran(ctx context.Context, p inti.Pelaku, id, kategori string, berkas []BerkasUnggahan) (
	[]models.Lampiran, error) {
	if err := l.periksaPelaku(p); err != nil {
		return nil, err
	}
	k, err := l.g.Keadaan(ctx, nil, id)
	if err != nil {
		return nil, err
	}
	if k.Tertutup() {
		return nil, ErrKasusTertutup
	}
	if !Pemegang(p, k) {
		return nil, ErrBukanPemegang
	}
	kat, err := l.g.KategoriLampiran(ctx, k.ID)
	if err != nil {
		return nil, err
	}
	if len(berkas) == 0 {
		return nil, &GalatValidasi{Pesan: []string{"No file attached"}}
	}
	master := models.CacahLampiran(kat)
	kategori = strings.TrimSpace(kategori)
	for i := range berkas {
		if kategori != "" { // S1.2: TempInputParam.pyCategory terisi -> berlaku untuk semua berkas
			berkas[i].Kategori = kategori
		}
		berkas[i].Kategori = strings.TrimSpace(berkas[i].Kategori)
		if _, ada := master[berkas[i].Kategori]; !ada {
			return nil, &GalatValidasi{Pesan: []string{fmt.Sprintf(
				"Category %q of file %q is not in the claim document category master", berkas[i].Kategori, berkas[i].Nama)}}
		}
	}
	if l.berkas == nil {
		return nil, ErrPenyimpananBelumDipasang
	}
	var out []models.Lampiran
	for _, b := range berkas {
		a, err := l.simpanLampiran(ctx, p.AkunID, k.ID, b.Kategori, b)
		if err != nil {
			return out, err
		}
		out = append(out, a)
	}
	return out, nil
}

// simpanLampiran = InsertDocument_Act untuk satu berkas.
func (l *Layanan) simpanLampiran(ctx context.Context, akun, klaimID, kategori string, b BerkasUnggahan) (models.Lampiran,
	error) {
	nama := strings.TrimSpace(b.Nama)
	if i := strings.LastIndexAny(nama, `/\`); i >= 0 {
		nama = strings.TrimSpace(nama[i+1:])
	}
	switch {
	case nama == "":
		return models.Lampiran{}, &GalatValidasi{Pesan: []string{"No file attached"}}
	case len(b.Isi) > unggah.BatasUkuranUnggahan:
		return models.Lampiran{}, &GalatValidasi{Pesan: []string{fmt.Sprintf("The file is larger than %d MB",
			unggah.BatasUkuranUnggahan>>20)}}
	}
	mime := penyimpanan.Ekstensi(nama) // S2-S3: .pyFileType / ekstensi NAMAFILE, huruf kecil
	o, err := l.berkas.Unggah(ctx, penyimpanan.MasukUnggah{Folder: models.FolderLampiranKlaim, NamaFile: nama,
		Isi: b.Isi, Durasi: models.DurasiLampiranKlaim, Ext: mime, Pengguna: akun})
	if err != nil {
		return models.Lampiran{}, err
	}
	kini := l.jam()
	baru, err := l.catatDokumenKlaim(ctx, o, models.BarisDokumenKlaim{Tanggal: kini, IDPega: models.KunciInstans(klaimID),
		NamaFile: nama, MIME: mime, Kategori1: kategori, StorageID: o.ImageID, Operator: akun})
	if err != nil {
		return models.Lampiran{}, err
	}
	return models.Lampiran{ID: baru.ID, NamaFile: nama, Kategori: kategori, MIME: mime,
		Tanggal: kini.In(models.Jakarta).Format("2006-01-02 15:04:05"), Operator: akun}, nil
}

// catatDokumenKlaim = InsertDocument_Act S4-S5 sesudah objek `o` terunggah: Insert_T_Storage_SQL + baris dokumen klaim
// `baru` di SATU transaksi; ID (`yyyyMMddhhmmssSSS` dari `baru.Tanggal`) yang bentrok dicoba +1 milidetik. Dipakai
// lampiran dan dokumen akseptasi.
func (l *Layanan) catatDokumenKlaim(ctx context.Context, o penyimpanan.Objek, baru models.BarisDokumenKlaim) (
	models.BarisDokumenKlaim, error) {
	err := l.g.Transaksi(ctx, func(tx *db.Tx) error {
		if err := l.berkas.Catat(ctx, tx, o); err != nil {
			return err
		}
		for i := 0; i < cobaIDDokumen; i++ {
			baru.ID = models.IDDokumenKlaim(baru.Tanggal.Add(time.Duration(i) * time.Millisecond))
			err := l.g.SisipDokumenKlaim(ctx, tx, baru)
			if errors.Is(err, repository.ErrIDDokumenTerpakai) {
				continue
			}
			return err
		}
		return fmt.Errorf("services: tidak menemukan ID dokumen klaim kosong (%d percobaan)", cobaIDDokumen)
	})
	if err != nil {
		// Objeknya sudah di penyimpanan tanpa catatan - dicatat di log untuk dibersihkan (pola Bordereaux).
		log.Printf("claimfacin: dokumen %s klaim %s tidak tercatat, objek penyimpanan tertinggal", o.ImageID, baru.IDPega)
	}
	return baru, err
}

// IsiLampiran - isi satu dokumen klaim untuk diunduh.
type IsiLampiran struct {
	Isi  io.ReadCloser
	Nama string
	Mime string
}

// ErrLampiranTidakAda - dokumen tidak ada di klaim ini.
var ErrLampiranTidakAda = fmt.Errorf("%w: lampiran tidak ada di klaim ini", ErrKasusTidakAda)

// ambilLampiran - satu dokumen klaim kasus `id` (pelaku boleh membuka kasus).
func (l *Layanan) ambilLampiran(ctx context.Context, p inti.Pelaku, id, lid string) (models.Kasus, models.Lampiran, error) {
	if err := l.periksaPelaku(p); err != nil {
		return models.Kasus{}, models.Lampiran{}, err
	}
	k, err := l.g.Keadaan(ctx, nil, id)
	if err != nil {
		return models.Kasus{}, models.Lampiran{}, err
	}
	dok, err := l.g.DaftarLampiran(ctx, k.ID)
	if err != nil {
		return models.Kasus{}, models.Lampiran{}, err
	}
	for _, a := range dok {
		if a.ID == strings.TrimSpace(lid) {
			return k, a, nil
		}
	}
	return models.Kasus{}, models.Lampiran{}, ErrLampiranTidakAda
}

// TautanOfficeLampiran - View Office Online (pola NB `DownloadDocumentPolis` 3): URL bertanda tangan untuk penampil
// kantor, hanya xls / xlsx / doc / docx / ppt / pptx yang objeknya tercatat.
func (l *Layanan) TautanOfficeLampiran(ctx context.Context, p inti.Pelaku, id, lid string) (string, error) {
	_, a, err := l.ambilLampiran(ctx, p, id, lid)
	if err != nil {
		return "", err
	}
	if !models.EkstensiOffice[strings.ToLower(a.MIME)] {
		return "", &GalatValidasi{Pesan: []string{"View Office Online is only for Excel, Word, and PowerPoint files"}}
	}
	if !a.AdaObjek {
		return "", penyimpanan.ErrObjekTidakAda
	}
	if l.berkas == nil {
		return "", ErrPenyimpananBelumDipasang
	}
	return l.berkas.Tautan(ctx, a.StorageID, models.DurasiLihatLampiran, p.AkunID)
}

// HapusLampiran - Delete (pola NB `DeleteDocumentPolis_Act`): objek di penyimpanan (bila tercatat), lalu catatan
// T_STORAGE_IMAGE dan baris dokumen klaim di SATU transaksi. Hanya pemegang assignment kasus yang belum ditutup
// (aturan unggah Claim Fac In); hapus di penyimpanan gagal = baris tetap (dapat diulang).
func (l *Layanan) HapusLampiran(ctx context.Context, p inti.Pelaku, id, lid string) error {
	k, a, err := l.ambilLampiran(ctx, p, id, lid)
	if err != nil {
		return err
	}
	if k.Tertutup() {
		return ErrKasusTertutup
	}
	if !Pemegang(p, k) {
		return ErrBukanPemegang
	}
	if a.AdaObjek && l.berkas == nil {
		return ErrPenyimpananBelumDipasang
	}
	// Satu transaksi: baris DOCUMENT_CLAIM dihapus, lalu objek storage (HapusObjek membaca catatan T_STORAGE_IMAGE yang
	// masih ada), lalu catatannya. Objek storage gagal dihapus = transaksi batal, baris tetap utuh (temuan review
	// 10-10-2026: dahulu objek dihapus sebelum transaksi, sehingga transaksi yang gagal meninggalkan baris tanpa objek).
	return l.g.Transaksi(ctx, func(tx *db.Tx) error {
		if err := l.g.HapusDokumenKlaim(ctx, tx, k.ID, a.ID); err != nil {
			return err
		}
		if !a.AdaObjek {
			return nil
		}
		if err := l.berkas.HapusObjek(ctx, a.StorageID, p.AkunID); err != nil {
			return err
		}
		return l.berkas.HapusCatatan(ctx, tx, a.StorageID)
	})
}

// UnduhLampiran - tombol View File: `GetBase64Attachment` S5-S5.1 (`GetUrlGoogleStorage_Act` ImageID = T_STORAGE_ID,
// Durasi 1800); isi dibaca backend (pola Bordereaux - frontend tidak membuka URL bertanda tangan).
func (l *Layanan) UnduhLampiran(ctx context.Context, p inti.Pelaku, id, lid string) (IsiLampiran, error) {
	_, a, err := l.ambilLampiran(ctx, p, id, lid)
	if err != nil {
		return IsiLampiran{}, err
	}
	if !a.AdaObjek { // S5.1 berprasyarat ImageID kosong -> dilewati
		return IsiLampiran{}, penyimpanan.ErrObjekTidakAda
	}
	if l.berkas == nil {
		return IsiLampiran{}, ErrPenyimpananBelumDipasang
	}
	isi, err := l.berkas.Buka(ctx, a.StorageID, models.DurasiLihatLampiran, p.AkunID)
	if err != nil {
		return IsiLampiran{}, err
	}
	return IsiLampiran{Isi: isi, Nama: a.NamaFile, Mime: unggah.MimeDariNamaFile(a.NamaFile)}, nil
}

// PindahKategoriLampiran - Change Category (layar Pega View File, screenshot work owner 09-10-2026): dokumen terpilih
// dipindah ke `kategori` master FAC di SATU transaksi (gagal satu = tidak ada yang pindah). Hanya pemegang assignment
// kasus yang belum ditutup (aturan unggah Claim Fac In).
func (l *Layanan) PindahKategoriLampiran(ctx context.Context, p inti.Pelaku, id, kategori string, lids []string) error {
	if err := l.periksaPelaku(p); err != nil {
		return err
	}
	k, err := l.g.Keadaan(ctx, nil, id)
	if err != nil {
		return err
	}
	if k.Tertutup() {
		return ErrKasusTertutup
	}
	if !Pemegang(p, k) {
		return ErrBukanPemegang
	}
	kat, err := l.g.KategoriLampiran(ctx, k.ID)
	if err != nil {
		return err
	}
	kategori = strings.TrimSpace(kategori)
	if _, ada := models.CacahLampiran(kat)[kategori]; !ada {
		return &GalatValidasi{Pesan: []string{fmt.Sprintf("Category %q is not in the claim document category master",
			kategori)}}
	}
	if len(lids) == 0 {
		return &GalatValidasi{Pesan: []string{"No document selected"}}
	}
	dok, err := l.g.DaftarLampiran(ctx, k.ID)
	if err != nil {
		return err
	}
	ada := make(map[string]bool, len(dok))
	for _, a := range dok {
		ada[a.ID] = true
	}
	for _, lid := range lids {
		if !ada[strings.TrimSpace(lid)] {
			return ErrLampiranTidakAda
		}
	}
	return l.g.Transaksi(ctx, func(tx *db.Tx) error {
		for _, lid := range lids {
			if err := l.g.PindahKategoriDokumen(ctx, tx, k.ID, strings.TrimSpace(lid), kategori); err != nil {
				return err
			}
		}
		return nil
	})
}
