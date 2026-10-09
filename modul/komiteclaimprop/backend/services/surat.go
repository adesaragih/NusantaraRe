package services

// Untuk apa berkas ini: EMAIL DAN DOKUMEN AKSEPTASI.
//
//   - Email `SendEmailKlaim_KMT` dirakit SAAT DIKIRIM dari pengenal di MUATAN outbox (MUATAN T_LOG_SERVICE_RNM hanya
//     pengenal dan angka, tanpa nama atau alamat - claimlife/015); dipanggil pelaksana outbox (`efek.go`) di produksi.
//   - Dokumen `PrintFileAcceptance_TKMT` S9-S13 (keputusan work owner 08-10-2026, pola lampiran Bordereaux): sesudah
//     Submit tersimpan, PDF diunggah lewat `inti/backend/penyimpanan` (InsertGoogleStorage_Act, folder "Claim") lalu
//     T_STORAGE_IMAGE + DOCUMENT_CLAIM dicatat di SATU transaksi. Seperti Pega (tanpa gerbang IsPEGAPROD), di setiap
//     lingkungan. `[penyimpangan sadar]` urutan: Pega membuat PDF SEBELUM menyimpan kasus (S21 < S41) - tidak ditiru
//     (keputusan work owner 2026-09-18); kegagalan PDF tidak membatalkan keputusan, penyetuju diberi `GalatDokumen`.

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/penyimpanan"
	"nusantarare/modul/komiteclaimprop/backend/models"
	"nusantarare/modul/komiteclaimprop/backend/repository"
)

// Kunci isi MUATAN efek email-komite.
const (
	IsiJenis    = "jenis"    // EmailPenyetujuBerikut / EmailPembuatSetuju / EmailPembuatTolak
	IsiPenerima = "penerima" // ID akun penerima (S12 anggota berikut; S13-S15 pembuat kasus)
	IsiAnggota  = "anggota"  // ID baris tangga yang diputuskan Submit ini (S6 / S7)
)

// DokumenAkseptasi - berkas PrintFileAcceptance_TKMT S9-S11 yang S13 `InsertDocument_Act` simpan (IDPEGA, KATEGORI_1,
// NAMAFILE, MIME "pdf", isi; Folder "Claim" InsertDocument_Act S4).
type DokumenAkseptasi struct {
	IDPega, NamaBerkas, Kategori, Folder, MIME string
	Isi                                        []byte
}

// ErrPenyimpananBelumDipasang - layanan dirakit tanpa penyimpanan berkas (tanpa Oracle / salah rakit).
var ErrPenyimpananBelumDipasang = fmt.Errorf("%w: penyimpanan dokumen akseptasi belum dipasang",
	penyimpanan.ErrStorageBelumSiap)

// cobaIDDokumen - batas percobaan ID DOCUMENT_CLAIM bentrok (+1 milidetik setiap kali).
const cobaIDDokumen = 50

// MIMEDokumenAkseptasi - S13 `MIME="pdf"`.
const MIMEDokumenAkseptasi = "pdf"

// SusunEmailKomite = SendEmailKlaim_KMT S1-S18 untuk satu baris efek email-komite.
func (l *Layanan) SusunEmailKomite(ctx context.Context, komiteID string, isi map[string]string) (models.SurelKomite,
	error) {
	k, kl, err := l.muat(ctx, nil, komiteID, false)
	if err != nil {
		return models.SurelKomite{}, err
	}
	var putus, berikut *models.Anggota
	for i := range k.Tangga {
		if k.Tangga[i].ID == isi[IsiAnggota] {
			putus = &k.Tangga[i]
		}
	}
	if putus == nil {
		return models.SurelKomite{}, fmt.Errorf("%w: baris tangga %q kasus %s", ErrKasusTidakAda, isi[IsiAnggota], komiteID)
	}
	for i := range k.Tangga {
		if k.Tangga[i].Urut == putus.Urut+1 {
			berikut = &k.Tangga[i]
		}
	}
	jenis := isi[IsiJenis]
	b := models.BahanEmail{Jenis: jenis, KomiteID: k.ID, Keputusan: putus.Keputusan, Komentar: putus.Komentar,
		AcceptedNo: models.AdjustmentKlaim(kl)["AcceptedNo"]}
	if b.Pengirim, err = l.a.NamaPelaku(ctx, putus.OperatorID); err != nil { // S10 CARI13 OperatorID.pyUserName
		return models.SurelKomite{}, err
	}
	kepada := ""
	if jenis == models.EmailPenyetujuBerikut { // S12 ComiteeClaim(KomiteCount+1).IDKomite / .KomiteEmail
		if berikut == nil {
			return models.SurelKomite{}, fmt.Errorf("%w: penyetuju berikut kasus %s", ErrKasusTidakAda, komiteID)
		}
		b.Sapaan, kepada = berikut.Jabatan, berikut.Email
		if kepada == "" { // KomiteID workbasket (migrasi claimprop 537): semua anggotanya (keputusan 09-10-2026)
			anggota, err := l.a.EmailAnggotaWorkbasket(ctx, berikut.OperatorID)
			if err != nil {
				return models.SurelKomite{}, err
			}
			kepada = strings.Join(anggota, ", ")
		}
	} else { // S13-S15 Obj-Browse Data-Admin-Operator-ID pembuat: pyUserName / pyEmailAddress
		if b.Sapaan, err = l.a.NamaPelaku(ctx, isi[IsiPenerima]); err != nil {
			return models.SurelKomite{}, err
		}
		if kepada, err = l.a.EmailPelaku(ctx, isi[IsiPenerima]); err != nil {
			return models.SurelKomite{}, err
		}
	}
	d, err := models.SusunDataEmail(kl, b)
	if err != nil {
		return models.SurelKomite{}, err
	}
	h, err := models.RenderEmailKomite(d) // S16
	if err != nil {
		return models.SurelKomite{}, err
	}
	return models.SurelKomite{Akun: l.surel.AkunUntuk(kepada), Kepada: kepada, CC: l.surel.CC,
		Subjek: models.SubjekEmail(jenis, kl, d.KomiteNo), HTML: h}, nil
}

// SusunDokumenAkseptasi = PrintFileAcceptance_TKMT S5-S11 (PDF) kasus komite `komiteID` yang baru disetujui `pelaku`;
// `saat` = S4 `NameInput.CARI24`.
func (l *Layanan) SusunDokumenAkseptasi(ctx context.Context, komiteID, pelaku string, saat time.Time) (DokumenAkseptasi,
	error) {
	k, kl, err := l.muat(ctx, nil, komiteID, false)
	if err != nil {
		return DokumenAkseptasi{}, err
	}
	pembuat, err := l.a.NamaPelaku(ctx, pelaku) // stream: Create by OperatorID.pyUserName
	if err != nil {
		return DokumenAkseptasi{}, err
	}
	adj := models.AdjustmentKlaim(kl)
	d, err := models.SusunAcceptanceNote(kl, adj, saat, pembuat)
	if err != nil {
		return DokumenAkseptasi{}, err
	}
	isiPDF, err := models.PDFAcceptanceNote(d, saat)
	if err != nil {
		return DokumenAkseptasi{}, err
	}
	return DokumenAkseptasi{IDPega: models.KunciInstans(k.KlaimID), NamaBerkas: models.NamaBerkasAkseptasi(adj["AcceptedNo"]),
		Kategori: models.KategoriDokumenAkseptasi, Folder: models.FolderDokumenKlaim, MIME: MIMEDokumenAkseptasi,
		Isi: isiPDF}, nil
}

// simpanDokumenAkseptasi = PrintFileAcceptance_TKMT S9-S13 sesudah Submit tersimpan: PDF -> InsertGoogleStorage_Act
// (Folder "Claim", Durasi 1800, Ext = MIME) -> Insert_T_Storage_SQL + DOCUMENT_CLAIM di SATU transaksi.
func (l *Layanan) simpanDokumenAkseptasi(ctx context.Context, komiteID, akun string, saat time.Time) error {
	if l.berkas == nil {
		return ErrPenyimpananBelumDipasang
	}
	d, err := l.SusunDokumenAkseptasi(ctx, komiteID, akun, saat)
	if err != nil {
		return err
	}
	o, err := l.berkas.Unggah(ctx, penyimpanan.MasukUnggah{Folder: d.Folder, NamaFile: d.NamaBerkas, Isi: d.Isi,
		Durasi: models.DurasiDokumenKlaim, Ext: d.MIME, Pengguna: akun})
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
				IDPega: d.IDPega, NamaFile: d.NamaBerkas, MIME: d.MIME, Kategori1: d.Kategori, StorageID: o.ImageID,
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
		log.Printf("komiteclaimprop: dokumen akseptasi %s kasus %s tidak tercatat, objek penyimpanan tertinggal",
			o.ImageID, komiteID)
	}
	return err
}
