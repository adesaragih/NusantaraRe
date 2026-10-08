package services

// Untuk apa berkas ini: PERAKIT ISI EFEK - email `SendEmailKlaim_KMT` dan dokumen `PrintFileAcceptance_TKMT` dirakit
// SAAT DIKIRIM dari pengenal di MUATAN outbox (MUATAN T_LOG_SERVICE_RNM hanya pengenal dan angka, tanpa nama atau alamat
// - claimlife/015). Dipanggil pelaksana outbox (`efek.go`) di produksi sebelum panggilan nyata yang masih menunggu
// persetujuan manusia.

import (
	"context"
	"fmt"
	"time"

	"nusantarare/modul/komiteclaimprop/backend/models"
)

// Kunci isi MUATAN efek email-komite / dokumen-akseptasi.
const (
	IsiJenis    = "jenis"    // EmailPenyetujuBerikut / EmailPembuatSetuju / EmailPembuatTolak
	IsiPenerima = "penerima" // ID akun penerima (S12 anggota berikut; S13-S15 pembuat kasus)
	IsiAnggota  = "anggota"  // ID baris tangga yang diputuskan Submit ini (S6 / S7)
	IsiPelaku   = "pelaku"   // ID akun penyetuju (`OperatorID`)
	IsiSaat     = "saat"     // PrintFileAcceptance_TKMT S4 `NameInput.CARI24` (RFC 3339)
)

// DokumenAkseptasi - berkas PrintFileAcceptance_TKMT S9-S11 yang S13 `InsertDocument_Act` simpan (KATEGORI_1,
// NAMAFILE, MIME "pdf", isi; Folder "Claim" InsertDocument_Act S4).
type DokumenAkseptasi struct {
	NamaBerkas, Kategori, Folder, MIME string
	Isi                                []byte
}

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

// SusunDokumenAkseptasi = PrintFileAcceptance_TKMT S5-S11 (PDF) untuk satu baris efek dokumen-akseptasi.
func (l *Layanan) SusunDokumenAkseptasi(ctx context.Context, komiteID string, isi map[string]string) (DokumenAkseptasi,
	error) {
	_, kl, err := l.muat(ctx, nil, komiteID, false)
	if err != nil {
		return DokumenAkseptasi{}, err
	}
	saat, err := time.Parse(time.RFC3339, isi[IsiSaat])
	if err != nil {
		return DokumenAkseptasi{}, fmt.Errorf("%w: saat cetak %q", ErrPermintaanTidakSah, isi[IsiSaat])
	}
	pembuat, err := l.a.NamaPelaku(ctx, isi[IsiPelaku]) // stream: Create by OperatorID.pyUserName
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
	return DokumenAkseptasi{NamaBerkas: models.NamaBerkasAkseptasi(adj["AcceptedNo"]),
		Kategori: models.KategoriDokumenAkseptasi, Folder: models.FolderDokumenKlaim, MIME: MIMEDokumenAkseptasi,
		Isi: isiPDF}, nil
}
