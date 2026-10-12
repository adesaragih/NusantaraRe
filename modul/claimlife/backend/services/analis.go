package services

// Tombol `Save` layar Claim Analis - 12-10-2026 (work owner: "CLAIM ANALYST dibuat text area dan hanya bisa diedit oleh
// claim analyst, pada tahap lain hanya read only").
//
// Kolom `T_CLAIMLF_PREMIUMLIST_DETAIL.NOTES VARCHAR2(1000)` (migrasi 003) - label CLAIM ANALYST b15194, sudah dibaca
// sebagai models.Peserta.CatatanAnalis. Tanpa migrasi baru. Pola sama dengan medis.go (Recommendation).
//
// Dibaca sesudah: medis.go.

import (
	"context"
	"errors"
	"fmt"
	"strings"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/galat"
	"nusantarare/modul/claimlife/backend/models"
	"nusantarare/modul/claimlife/backend/repository"
)

// BatasCatatanAnalis - T_CLAIMLF_PREMIUMLIST_DETAIL.NOTES VARCHAR2(1000), migrasi 003. Byte, diperiksa sebelum bind.
const BatasCatatanAnalis = 1000

var (
	// ErrAnalisBukanClaimAnalis - tombol Save catatan Claim Analyst hanya ada di tahap Claim Analis.
	ErrAnalisBukanClaimAnalis = errors.New("services: catatan Claim Analyst hanya dapat disimpan pada tahap Claim Analis")
	// ErrCatatanAnalisKepanjangan - melebihi lebar kolom.
	ErrCatatanAnalisKepanjangan = errors.New("services: catatan Claim Analyst melebihi 1000 byte")
	// ErrCatatanAnalisGagalTulis - peserta tidak lagi aktif saat ditulis.
	ErrCatatanAnalisGagalTulis = errors.New("services: catatan Claim Analyst tidak dapat ditulis untuk peserta itu")
)

// IsianAnalis - isian tombol Save Claim Analis: catatan per pengenal peserta.
type IsianAnalis struct {
	Catatan map[string]string `json:"catatan"`
}

// PeriksaIsianAnalis - bentuk isian, MURNI. `aktif` = pengenal peserta aktif klaim.
func PeriksaIsianAnalis(is IsianAnalis, aktif map[string]bool) error {
	if len(is.Catatan) == 0 {
		return fmt.Errorf("%w: tidak ada isian untuk disimpan", galat.ErrPermintaanTidakSah)
	}
	for id, v := range is.Catatan {
		if !aktif[id] {
			return fmt.Errorf("%w: peserta %q bukan peserta aktif klaim", galat.ErrPermintaanTidakSah, id)
		}
		if len(v) > BatasCatatanAnalis {
			return fmt.Errorf("%w: peserta %q (%d byte)", ErrCatatanAnalisKepanjangan, id, len(v))
		}
	}
	return nil
}

// SimpanAnalis adalah layanan tombol Save catatan Claim Analyst.
type SimpanAnalis struct{ svc *Service }

// SimpanAnalis menyusun layanannya.
func (s *Service) SimpanAnalis() *SimpanAnalis { return &SimpanAnalis{svc: s} }

// Simpan menulis catatan Claim Analyst per peserta dalam SATU transaksi.
//
// Gerbang: identitas, pengenal, kasus terbuka, tahap Claim Analis + pemegangnya (ReasLifeSPV), bentuk isian, lalu tulis.
// ⛔ Status baris, nomor, dan tahap TIDAK disentuh.
func (x *SimpanAnalis) Simpan(ctx context.Context, pelaku inti.Pelaku, klaimID string, is IsianAnalis) error {
	if err := inti.WajibIdentitas(pelaku); err != nil {
		return err
	}
	if strings.TrimSpace(klaimID) == "" {
		return fmt.Errorf("%w: pengenal klaim wajib diisi", galat.ErrPermintaanTidakSah)
	}
	if x == nil || x.svc == nil || !x.svc.PunyaDatabase() {
		return db.ErrTanpaOracle
	}
	if err := x.svc.PastikanKasusTerbuka(ctx, klaimID); err != nil {
		if errors.Is(err, repository.ErrWorkTidakAda) {
			return fmt.Errorf("%w: %w", ErrKlaimTidakAda, err)
		}
		return err
	}
	baca := repository.NewKlaimLife(x.svc.DB())
	tahap, err := tahapKasus(ctx, baca, klaimID)
	if err != nil {
		return err
	}
	if tahap != models.TahapClaimAnalis {
		return fmt.Errorf("%w: tahap %s", ErrAnalisBukanClaimAnalis, tahap)
	}
	if err := WajibPemegangTahap(pelaku, tahap); err != nil {
		return err
	}
	peserta, err := baca.AmbilPeserta(ctx, klaimID)
	if err != nil {
		return err
	}
	aktif := map[string]bool{}
	for _, p := range peserta {
		aktif[p.ID] = true
	}
	if err := PeriksaIsianAnalis(is, aktif); err != nil {
		return err
	}
	return x.svc.DalamTransaksi(ctx, func(tx *db.Tx) error {
		for id, v := range is.Catatan {
			ok, err := baca.SetelCatatanAnalis(ctx, tx, klaimID, id, v)
			if err != nil {
				return err
			}
			if !ok {
				return fmt.Errorf("%w: peserta %q", ErrCatatanAnalisGagalTulis, id)
			}
		}
		return nil
	})
}
