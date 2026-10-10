package services

// Untuk apa berkas ini: DOKUMEN AKSEPTASI - `PrintPDFAccep_MultiAksep_KMT` (KomitePost_Adjustment S9): sesudah Submit
// tersimpan, PDF stream `AcceptanceNotePDF` diunggah lewat `inti/backend/penyimpanan` (InsertGoogleStorage_Act, folder
// "Claim") lalu T_STORAGE_IMAGE + DOCUMENT_CLAIM dicatat di SATU transaksi. Pola Komite Claim Prop (disalin, bukan
// diimpor): `[penyimpangan sadar]` urutan - Pega membuat PDF SEBELUM menyimpan kasus; di sini sesudah keputusan
// tersimpan, dan kegagalan PDF tidak membatalkan keputusan (penyetuju diberi `PesanDokumenGagal`). Di setiap
// lingkungan (Pega tanpa gerbang IsPEGAPROD).

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/penyimpanan"
	"nusantarare/modul/komiteclaimfacin/backend/models"
	"nusantarare/modul/komiteclaimfacin/backend/repository"
)

// dokumenAkseptasi - berkas S24-S26 yang S27 `InsertDocument_Act` simpan.
type dokumenAkseptasi struct {
	IDPega, NamaBerkas string
	Isi                []byte
}

// ErrPenyimpananBelumDipasang - layanan dirakit tanpa penyimpanan berkas (tanpa Oracle / salah rakit).
var ErrPenyimpananBelumDipasang = fmt.Errorf("%w: penyimpanan dokumen akseptasi belum dipasang",
	penyimpanan.ErrStorageBelumSiap)

// PesanDokumenGagal - keputusan TERSIMPAN, tetapi PDF akseptasi gagal dibuat / diunggah / dicatat (teks sama dengan
// Komite Claim Prop).
const PesanDokumenGagal = "The decision was saved, but the acceptance note PDF could not be stored. Please contact the administrator."

// cobaIDDokumen - batas percobaan ID DOCUMENT_CLAIM bentrok (+1 milidetik setiap kali).
const cobaIDDokumen = 50

// susunDokumenAkseptasi = PrintPDFAccep_MultiAksep_KMT S1-S26 (PDF) kasus komite `komiteID` yang baru disetujui
// `pelaku` (Technical PIC = `OperatorID.pyUserName`); `saat` = S3 / S14 `NameInput.CARI24`. Kasus dan klaim dibaca
// ulang sesudah commit (keadaan tersimpan, pola Komite Claim Prop).
func (l *Layanan) susunDokumenAkseptasi(ctx context.Context, komiteID, pelaku string, saat time.Time) (
	dokumenAkseptasi, error) {
	k, kl, err := l.muat(ctx, nil, komiteID, false)
	if err != nil {
		return dokumenAkseptasi{}, err
	}
	teknik, err := l.a.NamaPelaku(ctx, pelaku)
	if err != nil {
		return dokumenAkseptasi{}, err
	}
	reas, err := l.a.NamaJenisReas(ctx) // S12.2 CallSpreadingView (TreatyName)
	if err != nil {
		return dokumenAkseptasi{}, err
	}
	no := models.Adjustment(kl)["AcceptedNo"]
	d, err := models.SusunAcceptanceNote(kl, no, saat, teknik, reas)
	if err != nil {
		return dokumenAkseptasi{}, err
	}
	isi, err := models.PDFDokumen(d, saat)
	if err != nil {
		return dokumenAkseptasi{}, err
	}
	return dokumenAkseptasi{IDPega: models.KunciInstans(k.KlaimID), NamaBerkas: models.NamaBerkasAkseptasi(k.KlaimID, no),
		Isi: isi}, nil
}

// simpanDokumenAkseptasi = S24-S27 sesudah Submit tersimpan: PDF -> InsertGoogleStorage_Act (Folder "Claim", Durasi
// 1800, Ext = MIME) -> Insert_T_Storage_SQL + DOCUMENT_CLAIM di SATU transaksi. Panik penggambar dijadikan galat -
// keputusan sudah tersimpan.
func (l *Layanan) simpanDokumenAkseptasi(ctx context.Context, komiteID, akun string, saat time.Time) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("services: dokumen akseptasi panik: %v", r)
		}
	}()
	if l.berkas == nil {
		return ErrPenyimpananBelumDipasang
	}
	d, err := l.susunDokumenAkseptasi(ctx, komiteID, akun, saat)
	if err != nil {
		return err
	}
	o, err := l.berkas.Unggah(ctx, penyimpanan.MasukUnggah{Folder: models.FolderDokumenKlaim, NamaFile: d.NamaBerkas,
		Isi: d.Isi, Durasi: models.DurasiDokumenKlaim, Ext: models.MIMEDokumenAkseptasi, Pengguna: akun})
	if err != nil {
		return err
	}
	kini := l.jam()
	err = l.g.Transaksi(ctx, func(tx *db.Tx) error {
		if err := l.berkas.Catat(ctx, tx, o); err != nil {
			return err
		}
		for i := 0; i < cobaIDDokumen; i++ {
			err := l.g.SisipDokumenKlaim(ctx, tx, models.BarisDokumenKlaim{
				ID: models.IDDokumenKlaim(kini.Add(time.Duration(i) * time.Millisecond)), Tanggal: kini,
				IDPega: d.IDPega, NamaFile: d.NamaBerkas, MIME: models.MIMEDokumenAkseptasi,
				Kategori1: models.KategoriDokumenAkseptasi, StorageID: o.ImageID,
				Operator: akun})
			if errors.Is(err, repository.ErrIDDokumenTerpakai) {
				continue
			}
			return err
		}
		return fmt.Errorf("services: tidak menemukan ID DOCUMENT_CLAIM kosong (%d percobaan)", cobaIDDokumen)
	})
	if err != nil {
		// Objeknya sudah di penyimpanan tanpa catatan - dicatat di log untuk dibersihkan (pola Bordereaux).
		log.Printf("komiteclaimfacin: dokumen akseptasi %s kasus %s tidak tercatat, objek penyimpanan tertinggal",
			o.ImageID, komiteID)
	}
	return err
}
