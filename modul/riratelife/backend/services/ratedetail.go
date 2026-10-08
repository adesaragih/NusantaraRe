package services

import (
	"context"
	"errors"
	"strings"

	"nusantarare/modul/riratelife/backend/models"
	"nusantarare/modul/riratelife/backend/repository"
)

// SimpanRate - Save form Rate Detail (`InboxRIRate`, View Detail.xml): idRate kosong = tambah baris (`AddToList_Act`
// b3196), selain itu ubah baris (`EditList_DT` b9853). R/I RATE NAME (`TempIDUsedBy.USEDBY` b1747, disabled) dan
// IDUSEDBY (`TempIDUsedBy.ID` b1526) SELALU dari ringkasan idRingkasan - bukan isian. Rule `AddToList_Act` tidak
// tersedia; aturan di bawahnya ASUMSI (MODUL.md A13-A15):
//   - kombinasi (GENDER, AGE, CONTRACT) tidak boleh kembar di ringkasan yang sama (seperti Upload, A9);
//   - ID rate baru dari SEQ_M_RATE_LIFE (nomor terpakai dilewati);
//   - MODIFY OPERATOR / MODIFY DATE ringkasan ikut diperbarui (seperti Upload, A10).
func (l *Layanan) SimpanRate(ctx context.Context, a Aktor, idRingkasan, idRate string, isi models.IsianRate) (models.Rate, error) {
	if err := wajibPenuh(a); err != nil {
		return models.Rate{}, err
	}
	idRingkasan, idRate = strings.TrimSpace(idRingkasan), strings.TrimSpace(idRate)
	rapi, err := models.PeriksaIsianRate(isi)
	if err != nil {
		return models.Rate{}, tolak("%s", err.Error())
	}
	var hasil models.Rate
	err = l.tx(ctx, func(tx *dbTx) error {
		ring, err := l.gudang.Ambil(ctx, tx, idRingkasan)
		if err != nil {
			return err
		}
		lama, err := l.gudang.RateDari(ctx, tx, []string{idRingkasan})
		if err != nil {
			return err
		}
		kunci := models.KunciRate("", rapi.Gender, rapi.Age, rapi.Contract)
		milik := false
		for _, r := range lama {
			if r.ID == idRate {
				milik = true
				continue
			}
			if models.KunciRate("", r.Gender, r.Age, r.Contract) == kunci {
				return tolak("GENDER %s, AGE %s, CONTRACT %s already exists in R/I RATE NAME %s (rate ID %s)",
					kosongTampil(rapi.Gender), kosongTampil(rapi.Age), rapi.Contract, ring.UsedBy, r.ID)
			}
		}
		if idRate != "" && !milik {
			return ErrRateTidakAda
		}
		hasil = models.Rate{ID: idRate, IDUsedBy: ring.ID, UsedBy: ring.UsedBy, Gender: rapi.Gender,
			Contract: rapi.Contract, Age: rapi.Age, Rate: rapi.Rate}
		if idRate == "" {
			if hasil.ID, err = (&pemberiID{l: l}).baru(ctx, tx); err != nil {
				return err
			}
			if err := l.gudang.SisipRate(ctx, tx, hasil); err != nil {
				return err
			}
		} else if err := l.gudang.UbahRate(ctx, tx, hasil); err != nil {
			if errors.Is(err, repository.ErrTidakAda) {
				return ErrRateTidakAda
			}
			return err
		}
		ring.OperatorID, ring.ModifiedDate = potongByte(a.AkunID, models.BatasNama), FormatWaktuPega(l.jam())
		return l.gudang.UbahRingkasan(ctx, tx, ring)
	})
	if err != nil {
		return models.Rate{}, petaTidakAda(err)
	}
	return hasil, nil
}
