package services

// Untuk apa berkas ini: EMAIL `SendEmailKlaim_KMT` - dirakit SAAT DIKIRIM dari pengenal di MUATAN outbox (MUATAN
// T_LOG_SERVICE_RNM hanya pengenal dan angka, tanpa nama atau alamat - claimlife/015); dipanggil pelaksana outbox
// (`efek.go`) di produksi. Pola Komite Claim Prop. PDF persetujuan (`GenerateAccCNP_act`, stream `AccClaimKomite_HTML`)
// tidak dibangun: stream-nya tidak diekspor (OQ-CNP-22).

import (
	"context"
	"fmt"
	"strings"

	"nusantarare/modul/komiteclaimnonprop/backend/models"
)

// Kunci isi MUATAN efek email-komite.
const (
	IsiJenis    = "jenis"    // EmailPenyetujuBerikut / EmailPembuatSetuju
	IsiPenerima = "penerima" // ID akun penerima (S12 anggota berikut; S13-S15 pembuat kasus)
	IsiAnggota  = "anggota"  // ID baris tangga yang diputuskan Submit ini (S6 / S7)
)

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
		if kepada == "" { // KomiteID workbasket (migrasi claimnonprop 611): semua anggotanya (keputusan 09-10-2026)
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
