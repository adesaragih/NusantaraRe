package services

import (
	"context"
	"errors"
	"strings"

	"nusantarare/modul/ricommlife/backend/models"
	"nusantarare/modul/ricommlife/backend/repository"
)

// SimpanKomisi - Save form R/I COMM DETAIL (`InboxRIComm`): idKomisi kosong = tambah baris (`AddToList_Act` b2955),
// selain itu ubah baris (`EditList_DT` b9323). USEDBY (`TempIDUsedBy.USEDBY` b1718, disabled) dan IDUSEDBY
// (`TempIDUsedBy.ID` b1509) SELALU dari ringkasan idRingkasan - bukan isian. Rule `AddToList_Act` tidak tersedia;
// aturan di bawahnya ASUMSI (MODUL.md):
//   - kombinasi (CONTRACT, YEAR) tidak boleh kembar di ringkasan yang sama (seperti Upload);
//   - MODIFY OPERATOR / MODIFY DATE ringkasan ikut diperbarui (seperti Upload).
//
// Satu transaksi: rincian + ringkasan (butir 5).
func (l *Layanan) SimpanKomisi(ctx context.Context, a Aktor, idRingkasan, idKomisi string, isi models.IsianKomisi) (models.Komisi, error) {
	if err := wajibPenuh(a); err != nil {
		return models.Komisi{}, err
	}
	idRingkasan, idKomisi = strings.TrimSpace(idRingkasan), strings.TrimSpace(idKomisi)
	rapi, err := models.PeriksaIsianKomisi(isi)
	if err != nil {
		return models.Komisi{}, tolak("%s", err.Error())
	}
	var hasil models.Komisi
	err = l.tx(ctx, func(tx *dbTx) error {
		ring, err := l.gudang.Ambil(ctx, tx, idRingkasan)
		if err != nil {
			return err
		}
		lama, err := l.gudang.KomisiDari(ctx, tx, []string{idRingkasan})
		if err != nil {
			return err
		}
		kunci := models.KunciKomisi(rapi.Contract, rapi.Year)
		milik := false
		for _, k := range lama {
			if k.ID == idKomisi {
				milik = true
				continue
			}
			if models.KunciKomisi(k.Contract, k.Year) == kunci {
				return tolak("CONTRACT %s, YEAR %s already exists in R/I COMM NAME %s (ID %s)",
					rapi.Contract, rapi.Year, ring.UsedBy, k.ID)
			}
		}
		if idKomisi != "" && !milik {
			return ErrKomisiTidakAda
		}
		hasil = models.Komisi{ID: idKomisi, IDUsedBy: ring.ID, UsedBy: ring.UsedBy, Contract: rapi.Contract,
			Year: rapi.Year, Comm: rapi.Comm}
		if idKomisi == "" {
			if hasil.ID, err = (&pemberiID{l: l}).baru(ctx, tx); err != nil {
				return err
			}
			if err := l.gudang.SisipKomisi(ctx, tx, hasil); err != nil {
				return err
			}
		} else if err := l.gudang.UbahKomisi(ctx, tx, hasil); err != nil {
			if errors.Is(err, repository.ErrTidakAda) {
				return ErrKomisiTidakAda
			}
			return err
		}
		ring.OperatorID, ring.ModifiedDate = potongByte(a.AkunID, models.BatasNama), FormatWaktuPega(l.jam())
		return l.gudang.UbahRingkasan(ctx, tx, ring)
	})
	if err != nil {
		return models.Komisi{}, petaTidakAda(err)
	}
	return hasil, nil
}
