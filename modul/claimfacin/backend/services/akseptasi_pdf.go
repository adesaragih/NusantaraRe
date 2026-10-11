package services

// Untuk apa berkas ini: DOKUMEN AKSEPTASI tombol "Acceptation" - SaveAcceptation 12 `PrintPDFAccep_MultiAksep`: salinan
// halaman klaim diambil di langkah 12 (SEBELUM langkah 13 menyetel IsPrintAccept - saringan S9.3.11 VERBATIM); SESUDAH
// aksi tersimpan, PDF stream `AcceptanceNotePDF` disusun, digambar, diunggah lewat `inti/backend/penyimpanan` (S26
// InsertDocument_Act: InsertGoogleStorage_Act Folder "Claim", Durasi 1800) dan dicatat T_STORAGE_IMAGE + DOCUMENT_CLAIM
// di satu transaksi. `[penyimpangan sadar]` pola Komite Claim Prop: Pega menyimpan dokumen di tengah activity; di sini
// sesudah commit (penggambaran tidak menahan kunci kasus), dan kegagalan PDF tidak membatalkan akseptasi (pelaku diberi
// `models.PesanDokumenGagal`). Di setiap lingkungan (Pega tanpa gerbang IsPEGAPROD).

import (
	"context"
	"fmt"
	"time"

	"nusantarare/inti/backend/penyimpanan"
	"nusantarare/modul/claimfacin/backend/models"
)

// dokumenTunda - bahan PDF akseptasi yang diambil di langkah 12, diproses sesudah aksi tersimpan.
type dokumenTunda struct {
	h                    *models.Halaman // salinan halaman sebelum langkah 13
	o                    int             // Param.idxObj
	no, claimNo, kasusID string          // Param.NoAkseptasi, ClaimData.ClaimNo, ID kasus klaim
	pelaku               string          // OperatorID (Technical PIC)
	saat                 time.Time       // NameInput.CARI24
}

// tundaDokumenAkseptasi = PrintPDFAccep_MultiAksep atas adjustment (o, i, a) yang diklik: lini tanpa stream terekspor
// (S14-S20) -> info OQ-CFI-20; selainnya bahan PDF disimpan untuk sesudah commit.
func (j *jalanAksi) tundaDokumenAkseptasi(o, i, a int) error {
	if s := models.StreamAkseptasiLini(j.h); s != models.StreamAkseptasi {
		j.info = models.InfoTanpaStream(s)
		return nil
	}
	b, err := models.Adj(j.h, o, i, a)
	if err != nil {
		return err
	}
	claimNo := j.h.Ambil(models.CD + "ClaimNo") // CallActivityInputRegister 9 `ClaimData.ClaimNo := pyID`
	if claimNo == "" {
		claimNo = j.kasus.ID
	}
	j.dokumen = &dokumenTunda{h: j.h.Salin(), o: o, no: b["AcceptedNo"], claimNo: claimNo, kasusID: j.kasus.ID,
		pelaku: j.k.Pelaku, saat: j.k.Sekarang}
	return nil
}

// simpanDokumenAkseptasi = PrintPDFAccep_MultiAksep S1-S26 sesudah aksi tersimpan: susun (Technical PIC = penekan
// tombol `OperatorID.pyUserName`, TreatyName = REINSURANCETYPE S11.2), gambar, unggah, catat. Panik penggambar
// dijadikan galat - akseptasi sudah tersimpan.
func (l *Layanan) simpanDokumenAkseptasi(ctx context.Context, d dokumenTunda) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("services: dokumen akseptasi panik: %v", r)
		}
	}()
	if l.berkas == nil {
		return ErrPenyimpananBelumDipasang
	}
	teknik, err := l.a.NamaPelaku(ctx, d.pelaku)
	if err != nil {
		return err
	}
	reas := map[string]string{}
	for i := range d.h.AmbilDaftar(models.DaftarItem(d.o)) {
		for _, s := range d.h.AmbilDaftar(models.DaftarDiItem(d.o, i+1, models.AnakSpreadKlaim)) {
			if _, ada := reas[s["TreatyType"]]; ada {
				continue
			}
			if reas[s["TreatyType"]], err = l.a.NamaJenisReas(ctx, s["TreatyType"]); err != nil {
				return err
			}
		}
	}
	doc, err := models.SusunAcceptanceNote(d.h, d.o, d.no, d.claimNo, d.saat, teknik, reas)
	if err != nil {
		return err
	}
	isi, err := models.PDFDokumen(doc, d.saat)
	if err != nil {
		return err
	}
	nama := models.NamaBerkasAkseptasi(d.claimNo, d.no)
	o, err := l.berkas.Unggah(ctx, penyimpanan.MasukUnggah{Folder: models.FolderLampiranKlaim, NamaFile: nama, Isi: isi,
		Durasi: models.DurasiLampiranKlaim, Ext: models.MIMEDokumenAkseptasi, Pengguna: d.pelaku})
	if err != nil {
		return err
	}
	_, err = l.catatDokumenKlaim(ctx, o, models.BarisDokumenKlaim{Tanggal: l.jam(), IDPega: models.KunciInstans(d.kasusID),
		NamaFile: nama, MIME: models.MIMEDokumenAkseptasi, Kategori1: models.KategoriDokumenAkseptasi,
		StorageID: o.ImageID, Operator: d.pelaku})
	return err
}
