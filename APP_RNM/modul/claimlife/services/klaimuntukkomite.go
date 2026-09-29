package services

// Penyedia kontrak Claim Life untuk Komite - refactor bentuk B.
//
// Untuk apa berkas ini: Komite Claim Life dulu memanggil repository Claim Life
// (`NewKlaimLife`, `NewPohonKlaim`) dan `Service.PastikanKasusTerbuka`
// langsung - keduanya satu paket. Kini Komite modul sendiri dan hanya mengenal
// `inti/kontrak.KlaimKomite`; berkas ini implementasinya. Setiap metode
// meneruskan apa adanya ke pemanggilan yang dulu Komite lakukan sendiri.

import (
	"context"
	"time"

	"nusantarare/inti/db"
	"nusantarare/inti/kontrak"
	"nusantarare/modul/claimlife/repository"
)

// KlaimUntukKomite menyusun kontrak Claim Life yang dipakai Komite.
func KlaimUntukKomite(svc *Service) kontrak.KlaimKomite { return klaimKomite{svc: svc} }

type klaimKomite struct{ svc *Service }

func (k klaimKomite) baca() *repository.KlaimLife { return repository.NewKlaimLife(k.svc.DB()) }

func (k klaimKomite) TypeKlaim(ctx context.Context, klaimID string) (string, error) {
	return k.baca().TypeKlaim(ctx, klaimID)
}

// KodeBisnisKlaim - header yang tidak ada dijawab teks kosong, persis cara
// Komite dulu memeriksanya (`hdr == nil || KodeBisnis == ""`).
func (k klaimKomite) KodeBisnisKlaim(ctx context.Context, klaimID string) (string, error) {
	hdr, err := k.baca().AmbilHeader(ctx, klaimID)
	if err != nil || hdr == nil {
		return "", err
	}
	return hdr.KodeBisnis, nil
}

func (k klaimKomite) PerbaruiStatusBaris(ctx context.Context, tx *db.Tx, pesertaID, adjID,
	kodeLama, kode, nomorAksep string, tglAksep time.Time) error {
	return k.baca().PerbaruiStatusBaris(ctx, tx, pesertaID, adjID, kodeLama, kode, nomorAksep, tglAksep)
}

func (k klaimKomite) CerminkanHeader(ctx context.Context, tx *db.Tx, klaimID, kode, nomorAksep string) error {
	return k.baca().CerminkanHeader(ctx, tx, klaimID, kode, nomorAksep)
}

func (k klaimKomite) CabutPenandaDipilih(ctx context.Context, tx *db.Tx, pesertaID string) error {
	return k.baca().CabutPenandaDipilih(ctx, tx, pesertaID)
}

func (k klaimKomite) NomorAkseptasiDipakai(ctx context.Context, tx *db.Tx, nomor string) (bool, error) {
	return repository.NewPohonKlaim(k.svc.DB()).NomorAkseptasiDipakai(ctx, tx, nomor)
}

func (k klaimKomite) PastikanKasusTerbuka(ctx context.Context, klaimID string) error {
	return k.svc.PastikanKasusTerbuka(ctx, klaimID)
}
